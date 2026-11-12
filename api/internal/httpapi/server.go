// Package httpapi: REST + SSE API of market-api. The web (SvelteKit) calls it through the
// /api/* proxy, so cookies and CORS are same-origin.
package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/lioilsources/RobloxMarket/api/internal/ai"
	"github.com/lioilsources/RobloxMarket/api/internal/auth"
	"github.com/lioilsources/RobloxMarket/api/internal/config"
	"github.com/lioilsources/RobloxMarket/api/internal/db"
	"github.com/lioilsources/RobloxMarket/api/internal/ingest"
	"github.com/lioilsources/RobloxMarket/api/internal/shop"
	"github.com/lioilsources/RobloxMarket/api/internal/storage"
)

type Server struct {
	cfg    config.Config
	pool   *pgxpool.Pool
	q      *db.Queries
	fs     *storage.FS
	auth   *auth.Service
	shop   *shop.Shop
	ingest *ingest.Ingester // nil = UGC_URL not set
	emb    ai.Embedder
	llm    *ai.LLM // nil = mock mode

	loginLimit *limiter
}

func New(cfg config.Config, pool *pgxpool.Pool, fs *storage.FS, a *auth.Service, s *shop.Shop,
	in *ingest.Ingester, emb ai.Embedder, llm *ai.LLM) *Server {
	return &Server{cfg: cfg, pool: pool, q: db.New(pool), fs: fs, auth: a, shop: s, ingest: in, emb: emb, llm: llm,
		loginLimit: newLimiter(5, 10*time.Minute)}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", s.handleHealth)

	mux.HandleFunc("GET /api/config", s.handleConfig)
	mux.HandleFunc("GET /api/catalog/facets", s.handleFacets)
	mux.HandleFunc("GET /api/products", s.handleListProducts)
	mux.HandleFunc("GET /api/products/{slug}", s.handleProduct)
	mux.HandleFunc("GET /api/media/{id}/{name}", s.handleMedia)
	mux.HandleFunc("POST /api/chat", s.handleChat)

	mux.HandleFunc("POST /api/auth/login", s.handleLogin)
	mux.HandleFunc("GET /api/auth/verify", s.handleVerify)
	mux.HandleFunc("POST /api/auth/logout", s.handleLogout)
	mux.HandleFunc("GET /api/me", s.customerOnly(s.handleMe))
	mux.HandleFunc("GET /api/me/licenses", s.customerOnly(s.handleLicenses))
	mux.HandleFunc("POST /api/me/licenses/{id}/download", s.customerOnly(s.handleIssueDownload))
	mux.HandleFunc("POST /api/checkout", s.customerOnly(s.handleCheckout))
	mux.HandleFunc("GET /api/download/{token}", s.handleDownload)
	mux.HandleFunc("POST /api/webhooks/{provider}", s.handleWebhook)

	mux.HandleFunc("GET /api/admin/products", s.adminOnly(s.handleAdminProducts))
	mux.HandleFunc("PATCH /api/admin/products/{id}", s.adminOnly(s.handleAdminUpdateProduct))
	mux.HandleFunc("GET /api/admin/ingest", s.adminOnly(s.handleAdminIngestStatus))
	mux.HandleFunc("POST /api/admin/ingest", s.adminOnly(s.handleAdminIngestRun))
	mux.HandleFunc("GET /api/admin/orders", s.adminOnly(s.handleAdminOrders))
	return logRequests(mux)
}

type ctxHandler func(w http.ResponseWriter, r *http.Request, cust db.Customer)

func (s *Server) customerOnly(h ctxHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cust := s.auth.Customer(r)
		if cust == nil {
			httpErr(w, http.StatusUnauthorized, "please sign in")
			return
		}
		h(w, r, *cust)
	}
}

func (s *Server) adminOnly(h ctxHandler) http.HandlerFunc {
	return s.customerOnly(func(w http.ResponseWriter, r *http.Request, cust db.Customer) {
		if !s.cfg.IsAdmin(cust.Email) {
			httpErr(w, http.StatusForbidden, "admins only")
			return
		}
		h(w, r, cust)
	})
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if err := s.pool.Ping(r.Context()); err != nil {
		httpErr(w, http.StatusServiceUnavailable, "db: %v", err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleConfig(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"currency":       s.cfg.Currency,
		"paid_providers": s.shop.PaidProviders(),
		"chat_mode":      map[bool]string{true: "llm", false: "mock"}[s.llm != nil],
	})
}

// --- helpers --------------------------------------------------------------

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func httpErr(w http.ResponseWriter, status int, format string, args ...any) {
	writeJSON(w, status, map[string]string{"error": fmt.Sprintf(format, args...)})
}

// serverErr logs the detail and returns a generic message to the client; UserError passes through.
func serverErr(w http.ResponseWriter, err error) {
	var ue shop.UserError
	if errors.As(err, &ue) {
		httpErr(w, http.StatusBadRequest, "%s", ue.Error())
		return
	}
	if errors.Is(err, shop.ErrNotFound) {
		httpErr(w, http.StatusNotFound, "not found")
		return
	}
	log.Printf("error: %v", err)
	httpErr(w, http.StatusInternalServerError, "internal error")
}

func decode(r *http.Request, v any) error {
	r.Body = http.MaxBytesReader(nil, r.Body, 1<<20)
	return json.NewDecoder(r.Body).Decode(v)
}

func pathID(r *http.Request, name string) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue(name), 10, 64)
	return id, err == nil && id > 0
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

// Flush passes SSE through the wrapper.
func (w *statusWriter) Flush() {
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(sw, r)
		path := r.URL.Path
		if strings.HasPrefix(path, "/api/download/") {
			path = "/api/download/…" // the token doesn't belong in the log
		}
		if path != "/healthz" {
			log.Printf("%s %s %d %s", r.Method, path, sw.status, time.Since(start).Round(time.Millisecond))
		}
	})
}

// limiter: at most n events per key per window (against magic-link spamming).
type limiter struct {
	mu     sync.Mutex
	n      int
	window time.Duration
	hits   map[string][]time.Time
}

func newLimiter(n int, window time.Duration) *limiter {
	return &limiter{n: n, window: window, hits: map[string][]time.Time{}}
}

func (l *limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	recent := l.hits[key][:0]
	for _, t := range l.hits[key] {
		if now.Sub(t) < l.window {
			recent = append(recent, t)
		}
	}
	if len(recent) >= l.n {
		l.hits[key] = recent
		return false
	}
	l.hits[key] = append(recent, now)
	return true
}
