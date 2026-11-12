package httpapi

import (
	"errors"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/lioilsources/RobloxMarket/api/internal/catalog"
	"github.com/lioilsources/RobloxMarket/api/internal/db"
)

func (s *Server) handleAdminProducts(w http.ResponseWriter, r *http.Request, _ db.Customer) {
	cards, err := catalog.AdminList(r.Context(), s.pool)
	if err != nil {
		serverErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, cards)
}

// PATCH /api/admin/products/{id} — partial update: price, tier, texts, status.
func (s *Server) handleAdminUpdateProduct(w http.ResponseWriter, r *http.Request, _ db.Customer) {
	ctx := r.Context()
	id, ok := pathID(r, "id")
	if !ok {
		httpErr(w, http.StatusNotFound, "nenalezeno")
		return
	}
	var req struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
		PriceCents  *int32  `json:"price_cents"`
		Tier        *string `json:"tier"`
		Status      *string `json:"status"`
	}
	if err := decode(r, &req); err != nil {
		httpErr(w, http.StatusBadRequest, "neplatný požadavek")
		return
	}
	p, err := s.q.GetProductByID(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		httpErr(w, http.StatusNotFound, "nenalezeno")
		return
	}
	if err != nil {
		serverErr(w, err)
		return
	}
	upd := db.AdminUpdateProductParams{ID: p.ID, Title: p.Title, Description: p.Description,
		PriceCents: p.PriceCents, Tier: p.Tier, Status: p.Status}
	if req.Title != nil {
		upd.Title = strings.TrimSpace(*req.Title)
	}
	if req.Description != nil {
		upd.Description = strings.TrimSpace(*req.Description)
	}
	if req.PriceCents != nil {
		upd.PriceCents = *req.PriceCents
	}
	if req.Tier != nil {
		upd.Tier = *req.Tier
	}
	if req.Status != nil {
		upd.Status = *req.Status
	}
	switch {
	case upd.Title == "":
		httpErr(w, http.StatusBadRequest, "název nesmí být prázdný")
		return
	case upd.Tier != "free" && upd.Tier != "paid":
		httpErr(w, http.StatusBadRequest, "tier musí být free nebo paid")
		return
	case upd.Status != "draft" && upd.Status != "listed" && upd.Status != "archived":
		httpErr(w, http.StatusBadRequest, "neplatný status")
		return
	case upd.PriceCents < 0:
		httpErr(w, http.StatusBadRequest, "cena nesmí být záporná")
		return
	case upd.Tier == "paid" && upd.PriceCents == 0:
		httpErr(w, http.StatusBadRequest, "placený asset potřebuje cenu > 0")
		return
	}
	p, err = s.q.AdminUpdateProduct(ctx, upd)
	if err != nil {
		serverErr(w, err)
		return
	}
	if err := catalog.Index(ctx, s.q, s.emb, p); err != nil {
		serverErr(w, err)
		return
	}
	media, err := s.q.ListProductMedia(ctx, p.ID)
	if err != nil {
		serverErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, catalog.NewCard(p, media))
}

func (s *Server) handleAdminIngestStatus(w http.ResponseWriter, r *http.Request, _ db.Customer) {
	if s.ingest == nil {
		writeJSON(w, http.StatusOK, map[string]any{"enabled": false})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"enabled": true, "last": s.ingest.Last()})
}

func (s *Server) handleAdminIngestRun(w http.ResponseWriter, r *http.Request, _ db.Customer) {
	if s.ingest == nil {
		httpErr(w, http.StatusBadRequest, "ingest je vypnutý (UGC_URL není nastavené)")
		return
	}
	res, err := s.ingest.Run(r.Context())
	if err != nil {
		httpErr(w, http.StatusBadGateway, "ugc-backend: %v", err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (s *Server) handleAdminOrders(w http.ResponseWriter, r *http.Request, _ db.Customer) {
	rows, err := s.q.AdminListOrders(r.Context(), 200)
	if err != nil {
		serverErr(w, err)
		return
	}
	downloads, err := s.q.CountDownloads(r.Context())
	if err != nil {
		serverErr(w, err)
		return
	}
	if rows == nil {
		rows = []db.AdminListOrdersRow{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"orders": rows, "downloads": downloads})
}
