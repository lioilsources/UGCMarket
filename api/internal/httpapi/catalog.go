package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"

	"github.com/lioilsources/RobloxMarket/api/internal/catalog"
	"github.com/lioilsources/RobloxMarket/api/internal/db"
)

func (s *Server) handleFacets(w http.ResponseWriter, r *http.Request) {
	facets, err := catalog.LiveFacets(r.Context(), s.pool)
	if err != nil {
		serverErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, facets)
}

// GET /api/products?q=&tier=free|paid&filters=<JSON {attr_key: value}>
func (s *Server) handleListProducts(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()
	query := catalog.Query{Text: qs.Get("q"), Tier: qs.Get("tier")}
	if raw := qs.Get("filters"); raw != "" {
		if err := json.Unmarshal([]byte(raw), &query.Filters); err != nil {
			httpErr(w, http.StatusBadRequest, "invalid filters")
			return
		}
	}
	cards, err := catalog.List(r.Context(), s.pool, query)
	if err != nil {
		serverErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, cards)
}

type productFile struct {
	Format string `json:"format"`
	Bytes  int64  `json:"bytes"`
}

// Detail: listed product (admin also sees drafts), formats and whether the signed-in user already owns it.
func (s *Server) handleProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	p, err := s.q.GetProductBySlug(ctx, r.PathValue("slug"))
	if errors.Is(err, pgx.ErrNoRows) {
		httpErr(w, http.StatusNotFound, "asset not found")
		return
	}
	if err != nil {
		serverErr(w, err)
		return
	}
	cust := s.auth.Customer(r)
	if p.Status != "listed" && (cust == nil || !s.cfg.IsAdmin(cust.Email)) {
		httpErr(w, http.StatusNotFound, "asset not found")
		return
	}
	media, err := s.q.ListProductMedia(ctx, p.ID)
	if err != nil {
		serverErr(w, err)
		return
	}
	files, err := s.q.ListProductFiles(ctx, p.ID)
	if err != nil {
		serverErr(w, err)
		return
	}
	out := struct {
		catalog.Card
		Files []productFile `json:"files"`
		Owned bool          `json:"owned"`
	}{Card: catalog.NewCard(p, media), Files: []productFile{}}
	for _, f := range files {
		out.Files = append(out.Files, productFile{Format: f.Format, Bytes: f.Bytes})
	}
	if cust != nil {
		owned, err := s.q.CustomerOwnedProductIDs(ctx, db.CustomerOwnedProductIDsParams{CustomerID: cust.ID, Ids: []int64{p.ID}})
		if err != nil {
			serverErr(w, err)
			return
		}
		out.Owned = len(owned) > 0
	}
	writeJSON(w, http.StatusOK, out)
}

// GET /api/media/{id}/thumb | viewer.glb — previews only, never files for sale.
func (s *Server) handleMedia(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, ok := pathID(r, "id")
	kind := map[string]string{"thumb": "thumb", "viewer.glb": "viewer_glb"}[r.PathValue("name")]
	if !ok || kind == "" {
		httpErr(w, http.StatusNotFound, "not found")
		return
	}
	p, err := s.q.GetProductByID(ctx, id)
	if err == nil && p.Status != "listed" {
		if cust := s.auth.Customer(r); cust == nil || !s.cfg.IsAdmin(cust.Email) {
			err = pgx.ErrNoRows
		}
	}
	var m db.ProductMedium
	if err == nil {
		m, err = s.q.GetProductMedia(ctx, db.GetProductMediaParams{ProductID: id, Kind: kind, Position: 0})
	}
	if errors.Is(err, pgx.ErrNoRows) {
		httpErr(w, http.StatusNotFound, "not found")
		return
	}
	if err != nil {
		serverErr(w, err)
		return
	}
	f, err := s.fs.Open(m.StorageKey)
	if err != nil {
		serverErr(w, err)
		return
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil {
		serverErr(w, err)
		return
	}
	if kind == "viewer_glb" {
		w.Header().Set("Content-Type", "model/gltf-binary")
	}
	w.Header().Set("Cache-Control", "public, max-age=3600")
	http.ServeContent(w, r, "", st.ModTime(), f)
}
