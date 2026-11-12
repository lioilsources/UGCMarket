package httpapi

// End-to-end test of the whole Phase 1 flow against a real Postgres (pgvector)
// and a fake ugc-backend:
//
//	MARKET_TEST_DATABASE_URL=postgres://market:PASSWORD@NAS_HOST:5438/market_test?sslmode=disable go test ./...
//
// The test wipes the public schema, so it only runs against a database whose name ends in _test.

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	pgxvec "github.com/pgvector/pgvector-go/pgx"

	"github.com/lioilsources/RobloxMarket/api/internal/ai"
	"github.com/lioilsources/RobloxMarket/api/internal/auth"
	"github.com/lioilsources/RobloxMarket/api/internal/config"
	"github.com/lioilsources/RobloxMarket/api/internal/db"
	"github.com/lioilsources/RobloxMarket/api/internal/ingest"
	"github.com/lioilsources/RobloxMarket/api/internal/shop"
	"github.com/lioilsources/RobloxMarket/api/internal/storage"
	"github.com/lioilsources/RobloxMarket/api/internal/ugc"
)

type captureMailer struct {
	mu   sync.Mutex
	last map[string]string
}

func (m *captureMailer) Send(_ context.Context, to, _, text string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.last[strings.ToLower(to)] = text
	return nil
}

func (m *captureMailer) token(t *testing.T, to string) string {
	m.mu.Lock()
	defer m.mu.Unlock()
	match := regexp.MustCompile(`token=([A-Za-z0-9_-]+)`).FindStringSubmatch(m.last[strings.ToLower(to)])
	if match == nil {
		t.Fatalf("no magic link for %s in %q", to, m.last[to])
	}
	return match[1]
}

// fakeUGC mimics the ugc-backend endpoints that market uses.
type fakeUGC struct {
	mu   sync.Mutex
	jobs []ugc.Job
}

func (f *fakeUGC) handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /jobs", func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		defer f.mu.Unlock()
		_ = json.NewEncoder(w).Encode(f.jobs)
	})
	mux.HandleFunc("GET /items", func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		defer f.mu.Unlock()
		items := []map[string]any{}
		for _, j := range f.jobs {
			items = append(items, map[string]any{"id": j.ID, "name": "golden-crown-with-gems", "category": j.Category,
				"tags": []string{"crown", "gold"}, "state": "packed", "updated_at": j.UpdatedAt})
		}
		_ = json.NewEncoder(w).Encode(items)
	})
	mux.HandleFunc("GET /jobs/{id}/glb", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("full") == "1" {
			_, _ = io.WriteString(w, "glTF-full")
			return
		}
		_, _ = io.WriteString(w, "glTF-preview")
	})
	mux.HandleFunc("GET /jobs/{id}/preview", func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "\x89PNG-concept")
	})
	mux.HandleFunc("GET /packed/{id}/download", func(w http.ResponseWriter, r *http.Request) {
		zw := zip.NewWriter(w)
		for name, body := range map[string]string{"golden-crown.fbx": "FBX-data", "golden-crown_tex.png": "PNG-tex", "item.json": "{}"} {
			fw, _ := zw.Create(name)
			_, _ = io.WriteString(fw, body)
		}
		_ = zw.Close()
	})
	return mux
}

type client struct {
	t    *testing.T
	base string
	http *http.Client
}

func newClient(t *testing.T, base string) *client {
	jar, _ := cookiejar.New(nil)
	return &client{t: t, base: base, http: &http.Client{Jar: jar}}
}

func (c *client) do(method, path string, body any, out any) int {
	c.t.Helper()
	var rd io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		rd = bytes.NewReader(b)
	}
	req, _ := http.NewRequest(method, c.base+path, rd)
	resp, err := c.http.Do(req)
	if err != nil {
		c.t.Fatal(err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if out != nil {
		if err := json.Unmarshal(raw, out); err != nil {
			c.t.Fatalf("%s %s: %v (%s)", method, path, err, raw)
		}
	}
	return resp.StatusCode
}

func (c *client) login(m *captureMailer, email string) {
	c.t.Helper()
	if code := c.do("POST", "/api/auth/login", map[string]string{"email": email}, nil); code != 200 {
		c.t.Fatalf("login: %d", code)
	}
	if code := c.do("GET", "/api/auth/verify?token="+m.token(c.t, email), nil, nil); code != 200 {
		c.t.Fatalf("verify: %d", code)
	}
}

func TestEndToEnd(t *testing.T) {
	dbURL := os.Getenv("MARKET_TEST_DATABASE_URL")
	if dbURL == "" {
		t.Skip("MARKET_TEST_DATABASE_URL not set")
	}
	ctx := context.Background()
	reset, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		t.Fatal(err)
	}
	// The test wipes the schema, so it only ever runs against a *_test database.
	var dbName string
	if err := reset.QueryRow(ctx, "SELECT current_database()").Scan(&dbName); err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(dbName, "_test") {
		reset.Close(ctx)
		t.Fatalf("refusing to wipe database %q: MARKET_TEST_DATABASE_URL must point to a *_test database", dbName)
	}
	_, err = reset.Exec(ctx, "DROP SCHEMA public CASCADE; CREATE SCHEMA public;")
	reset.Close(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Migrate(dbURL); err != nil {
		t.Fatal(err)
	}
	pcfg, _ := pgxpool.ParseConfig(dbURL)
	pcfg.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error { return pgxvec.RegisterTypes(ctx, conn) }
	pool, err := pgxpool.NewWithConfig(ctx, pcfg)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	report, _ := json.Marshal(map[string]any{"tri_count": 3870, "texture_size": 1024, "watertight": true, "attachment": "HatAttachment"})
	fake := &fakeUGC{jobs: []ugc.Job{
		{ID: "crown01", Status: "packed", Prompt: "golden crown with gems", Category: "hat", Style: "cartoon",
			Verdict: "WARN", Report: report, UpdatedAt: time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)},
		{ID: "broken01", Status: "packed", Prompt: "broken", Category: "hat", Verdict: "FAIL", UpdatedAt: time.Now()},
	}}
	ugcSrv := httptest.NewServer(fake.handler())
	defer ugcSrv.Close()

	cfg := config.Config{Currency: "USD", AdminEmails: []string{"admin@example.com"}, PublicURL: "http://market.test"}
	fs, _ := storage.NewFS(t.TempDir())
	mailer := &captureMailer{last: map[string]string{}}
	emb := ai.MockEmbedder{}
	q := db.New(pool)
	in := ingest.New(pool, fs, ugc.New(ugcSrv.URL), emb, 300, "USD")
	api := New(cfg, pool, fs, auth.New(q, mailer, cfg.PublicURL, false), shop.New(pool, fs, mailer, cfg.PublicURL, shop.Free{}), in, emb, nil)
	srv := httptest.NewServer(api.Routes())
	defer srv.Close()

	// 1) ingest: FAIL is skipped, a new piece is a draft and not publicly visible
	admin := newClient(t, srv.URL)
	admin.login(mailer, "Admin@Example.com")
	var res ingest.Result
	if code := admin.do("POST", "/api/admin/ingest", nil, &res); code != 200 || res.Created != 1 || res.SkippedFail != 1 {
		t.Fatalf("ingest: %d %+v", code, res)
	}
	var public []map[string]any
	anon := newClient(t, srv.URL)
	anon.do("GET", "/api/products", nil, &public)
	if len(public) != 0 {
		t.Fatalf("draft must not be public: %v", public)
	}

	// 2) admin: list as free
	var adminList []map[string]any
	admin.do("GET", "/api/admin/products", nil, &adminList)
	if len(adminList) != 1 || adminList[0]["title"] != "Golden Crown with Gems" {
		t.Fatalf("admin list: %v", adminList)
	}
	id := int64(adminList[0]["id"].(float64))
	slug := adminList[0]["slug"].(string)
	if code := admin.do("PATCH", fmt.Sprintf("/api/admin/products/%d", id), map[string]any{"tier": "free", "status": "listed"}, nil); code != 200 {
		t.Fatalf("patch: %d", code)
	}
	if code := anon.do("PATCH", fmt.Sprintf("/api/admin/products/%d", id), map[string]any{"status": "draft"}, nil); code != 401 {
		t.Fatalf("anon admin: %d", code)
	}

	// 3) catalog, filters, facets, detail, media
	anon.do("GET", `/api/products?filters={"category":["hat"],"polycount":{"min":1000}}`, nil, &public)
	if len(public) != 1 || public[0]["viewer_url"] == "" {
		t.Fatalf("filtered list: %v", public)
	}
	anon.do("GET", `/api/products?filters={"category":["back"]}`, nil, &public)
	if len(public) != 0 {
		t.Fatalf("category filter: %v", public)
	}
	var facets []catalogFacet
	anon.do("GET", "/api/catalog/facets", nil, &facets)
	if !hasFacet(facets, "category") || !hasFacet(facets, "polycount") {
		t.Fatalf("facets: %+v", facets)
	}
	var detail map[string]any
	if code := anon.do("GET", "/api/products/"+slug, nil, &detail); code != 200 || len(detail["files"].([]any)) != 3 {
		t.Fatalf("detail: %d %v", code, detail)
	}
	if resp, _ := http.Get(srv.URL + fmt.Sprintf("/api/media/%d/viewer.glb", id)); resp.StatusCode != 200 {
		t.Fatalf("media: %d", resp.StatusCode)
	}

	// 4) chat (mock) returns products + text
	resp, err := http.Post(srv.URL+"/api/chat", "application/json", strings.NewReader(`{"messages":[{"role":"user","content":"golden crown"}]}`))
	if err != nil {
		t.Fatal(err)
	}
	stream, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if !strings.Contains(string(stream), "event: products") || !strings.Contains(string(stream), slug) || !strings.Contains(string(stream), "event: done") {
		t.Fatalf("chat stream: %s", stream)
	}

	// 5) customer: free checkout → license → download
	buyer := newClient(t, srv.URL)
	if code := buyer.do("POST", "/api/checkout", map[string]any{"product_ids": []int64{id}}, nil); code != 401 {
		t.Fatalf("checkout without login: %d", code)
	}
	buyer.login(mailer, "buyer@example.com")
	var co shop.CheckoutResult
	if code := buyer.do("POST", "/api/checkout", map[string]any{"product_ids": []int64{id}}, &co); code != 200 || co.Status != "paid" {
		t.Fatalf("checkout: %d %+v", code, co)
	}
	if code := buyer.do("POST", "/api/checkout", map[string]any{"product_ids": []int64{id}}, nil); code != 400 {
		t.Fatalf("second checkout of owned product: %d", code)
	}
	var lics []map[string]any
	buyer.do("GET", "/api/me/licenses", nil, &lics)
	if len(lics) != 1 || lics[0]["license_type"] != "personal" {
		t.Fatalf("licenses: %v", lics)
	}
	licID := int64(lics[0]["id"].(float64))
	if code := admin.do("POST", fmt.Sprintf("/api/me/licenses/%d/download", licID), nil, nil); code != 404 {
		t.Fatalf("foreign license download: %d", code)
	}
	var dl map[string]string
	buyer.do("POST", fmt.Sprintf("/api/me/licenses/%d/download", licID), nil, &dl)
	resp, err = http.Get(srv.URL + dl["url"])
	if err != nil {
		t.Fatal(err)
	}
	zipBody, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	zr, err := zip.NewReader(bytes.NewReader(zipBody), int64(len(zipBody)))
	if err != nil {
		t.Fatalf("download zip: %v (%s)", err, zipBody)
	}
	got := map[string]string{}
	for _, f := range zr.File {
		rc, _ := f.Open()
		b, _ := io.ReadAll(rc)
		rc.Close()
		got[f.Name] = string(b)
	}
	for _, name := range []string{slug + ".fbx", slug + "_tex.png", slug + ".glb", "LICENSE.txt", "README.txt"} {
		if _, ok := got[slug+"/"+name]; !ok {
			t.Fatalf("zip missing %s: %v", name, keys(got))
		}
	}
	if got[slug+"/"+slug+".fbx"] != "FBX-data" || !strings.Contains(got[slug+"/LICENSE.txt"], lics[0]["license_key"].(string)) {
		t.Fatalf("zip content wrong")
	}

	// 6) paid product: checkout is rejected without a provider
	admin.do("PATCH", fmt.Sprintf("/api/admin/products/%d", id), map[string]any{"tier": "paid", "price_cents": 300}, nil)
	other := newClient(t, srv.URL)
	other.login(mailer, "other@example.com")
	if code := other.do("POST", "/api/checkout", map[string]any{"product_ids": []int64{id}, "provider": "stripe"}, nil); code != 400 {
		t.Fatalf("paid checkout without provider: %d", code)
	}

	// 7) re-ingest unchanged; disappearing from packed → archived; return → draft
	admin.do("POST", "/api/admin/ingest", nil, &res)
	if res.Unchanged != 1 || res.Created+res.Updated != 0 {
		t.Fatalf("re-ingest: %+v", res)
	}
	fake.mu.Lock()
	saved := fake.jobs
	fake.jobs = nil
	fake.mu.Unlock()
	admin.do("POST", "/api/admin/ingest", nil, &res)
	if res.Archived != 1 {
		t.Fatalf("archive: %+v", res)
	}
	fake.mu.Lock()
	fake.jobs = saved
	fake.mu.Unlock()
	admin.do("POST", "/api/admin/ingest", nil, &res)
	admin.do("GET", "/api/admin/products", nil, &adminList)
	if res.Updated != 1 || adminList[0]["status"] != "draft" {
		t.Fatalf("restore: %+v %v", res, adminList[0]["status"])
	}

	// 8) admin order overview
	var orders map[string]any
	admin.do("GET", "/api/admin/orders", nil, &orders)
	if len(orders["orders"].([]any)) != 1 || orders["downloads"].(float64) != 1 {
		t.Fatalf("orders: %v", orders)
	}
}

type catalogFacet struct {
	Key string `json:"attr_key"`
}

func hasFacet(fs []catalogFacet, key string) bool {
	for _, f := range fs {
		if f.Key == key {
			return true
		}
	}
	return false
}

func keys(m map[string]string) []string {
	out := []string{}
	for k := range m {
		out = append(out, k)
	}
	return out
}
