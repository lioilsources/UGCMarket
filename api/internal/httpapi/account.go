package httpapi

import (
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/lioilsources/RobloxMarket/api/internal/auth"
	"github.com/lioilsources/RobloxMarket/api/internal/db"
	"github.com/lioilsources/RobloxMarket/api/internal/shop"
)

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
		Next  string `json:"next"`
	}
	if err := decode(r, &req); err != nil {
		httpErr(w, http.StatusBadRequest, "invalid request")
		return
	}
	email, err := auth.NormalizeEmail(req.Email)
	if err != nil {
		httpErr(w, http.StatusBadRequest, "%v", err)
		return
	}
	if !s.loginLimit.Allow("email:"+email) || !s.loginLimit.Allow("ip:"+clientIP(r)) {
		httpErr(w, http.StatusTooManyRequests, "too many attempts, try again in a few minutes")
		return
	}
	if err := s.auth.SendLink(r.Context(), email, req.Next); err != nil {
		serverErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "sent"})
}

// GET /api/auth/verify?token=… — the web route /auth/verify calls it and forwards Set-Cookie.
func (s *Server) handleVerify(w http.ResponseWriter, r *http.Request) {
	cust, session, err := s.auth.Verify(r.Context(), r.URL.Query().Get("token"))
	if errors.Is(err, auth.ErrInvalidLink) {
		httpErr(w, http.StatusBadRequest, "%v", err)
		return
	}
	if err != nil {
		serverErr(w, err)
		return
	}
	s.auth.SetCookie(w, session)
	writeJSON(w, http.StatusOK, s.meBody(cust))
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	s.auth.Logout(w, r)
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) meBody(c db.Customer) map[string]any {
	return map[string]any{"id": c.ID, "email": c.Email, "is_admin": s.cfg.IsAdmin(c.Email)}
}

func (s *Server) handleMe(w http.ResponseWriter, r *http.Request, cust db.Customer) {
	writeJSON(w, http.StatusOK, s.meBody(cust))
}

func (s *Server) handleLicenses(w http.ResponseWriter, r *http.Request, cust db.Customer) {
	rows, err := s.q.ListCustomerLicenses(r.Context(), cust.ID)
	if err != nil {
		serverErr(w, err)
		return
	}
	if rows == nil {
		rows = []db.ListCustomerLicensesRow{}
	}
	writeJSON(w, http.StatusOK, rows)
}

func (s *Server) handleIssueDownload(w http.ResponseWriter, r *http.Request, cust db.Customer) {
	id, ok := pathID(r, "id")
	if !ok {
		httpErr(w, http.StatusNotFound, "not found")
		return
	}
	token, err := s.shop.IssueDownload(r.Context(), cust.ID, id)
	if err != nil {
		serverErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"url": "/api/download/" + token})
}

func (s *Server) handleDownload(w http.ResponseWriter, r *http.Request) {
	if err := s.shop.ServeDownload(w, r, r.PathValue("token")); err != nil {
		if errors.Is(err, shop.ErrNotFound) {
			httpErr(w, http.StatusNotFound, "this download link is invalid or expired — get a fresh one under My assets")
			return
		}
		log.Printf("download: %v", err)
	}
}

func (s *Server) handleCheckout(w http.ResponseWriter, r *http.Request, cust db.Customer) {
	var req shop.CheckoutRequest
	if err := decode(r, &req); err != nil {
		httpErr(w, http.StatusBadRequest, "invalid request")
		return
	}
	if req.TaxCountry == "" {
		req.TaxCountry = r.Header.Get("CF-IPCountry") // Cloudflare; record-keeping only, not a tax calculation
	}
	res, err := s.shop.Checkout(r.Context(), cust, req)
	if err != nil {
		serverErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

// POST /api/webhooks/{provider} — in Phase 1 no WebhookProvider is registered.
func (s *Server) handleWebhook(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("provider")
	p, ok := s.shop.Provider(name)
	wp, isWebhook := p.(shop.WebhookProvider)
	if !ok || !isWebhook {
		httpErr(w, http.StatusNotFound, "unknown provider")
		return
	}
	ev, err := wp.ParseWebhook(r)
	if err != nil {
		httpErr(w, http.StatusBadRequest, "%v", err)
		return
	}
	if err := s.shop.HandlePaymentEvent(r.Context(), name, ev); err != nil {
		serverErr(w, err) // 5xx → the provider retries the webhook
		return
	}
	_, _ = io.WriteString(w, "ok")
}

func clientIP(r *http.Request) string {
	if ip := r.Header.Get("CF-Connecting-IP"); ip != "" {
		return ip
	}
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return strings.TrimSpace(strings.Split(xff, ",")[0])
	}
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	return host
}
