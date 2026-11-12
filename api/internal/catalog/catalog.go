// Package catalog: public catalog — filters over JSONB attributes, facets and
// vector index for chat. The filter logic is a port of the shop
// (apps/web/src/routes/api/products/+server.ts), but facets are not computed
// at ingest into filter_schemas, but live from the listed products.
package catalog

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgvector/pgvector-go"

	"github.com/lioilsources/RobloxMarket/api/internal/ai"
	"github.com/lioilsources/RobloxMarket/api/internal/db"
)

// Facet = a filter over a single attributes key. The list is a “vertical config”
// for 3D assets; to add a key, just add it here (and populate it in ingest).
type Facet struct {
	Key    string         `json:"attr_key"`
	Label  string         `json:"label"`
	Kind   string         `json:"kind"` // set | range | bool
	Config map[string]any `json:"config"`
}

// "style" is deliberately absent: in ugc it is a prompt style (“catalog photo, glossy”), not a catalog category.
var Facets = []Facet{
	{Key: "category", Label: "Category", Kind: "set"},
	{Key: "collection", Label: "Collection", Kind: "set"},
	{Key: "formats", Label: "Formats", Kind: "set"},
	{Key: "tags", Label: "Tags", Kind: "set"},
	{Key: "polycount", Label: "Triangles", Kind: "range"},
	{Key: "watertight", Label: "Watertight mesh", Kind: "bool"},
}

func facetKind(key string) string {
	for _, f := range Facets {
		if f.Key == key {
			return f.Kind
		}
	}
	return ""
}

// Card is the public shape of a product for the gallery, chat and detail.
type Card struct {
	ID          int64           `json:"id"`
	Slug        string          `json:"slug"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	PriceCents  int32           `json:"price_cents"`
	Currency    string          `json:"currency"`
	Tier        string          `json:"tier"`
	Status      string          `json:"status"`
	Attributes  json.RawMessage `json:"attributes"`
	ThumbURL    string          `json:"thumb_url,omitempty"`
	ViewerURL   string          `json:"viewer_url,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
}

func MediaURLs(id int64, hasThumb, hasViewer bool) (thumb, viewer string) {
	if hasThumb {
		thumb = fmt.Sprintf("/api/media/%d/thumb", id)
	}
	if hasViewer {
		viewer = fmt.Sprintf("/api/media/%d/viewer.glb", id)
	}
	return
}

func NewCard(p db.Product, media []db.ProductMedium) Card {
	var hasThumb, hasViewer bool
	for _, m := range media {
		hasThumb = hasThumb || m.Kind == "thumb"
		hasViewer = hasViewer || m.Kind == "viewer_glb"
	}
	c := Card{ID: p.ID, Slug: p.Slug, Title: p.Title, Description: p.Description, PriceCents: p.PriceCents,
		Currency: p.Currency, Tier: p.Tier, Status: p.Status, Attributes: p.Attributes, CreatedAt: p.CreatedAt}
	c.ThumbURL, c.ViewerURL = MediaURLs(p.ID, hasThumb, hasViewer)
	return c
}

type Query struct {
	Text    string
	Tier    string // "" | free | paid
	Filters map[string]json.RawMessage
	Limit   int
}

const cardColumns = `p.id, p.slug, p.title, p.description, p.price_cents, p.currency, p.tier, p.status,
	p.attributes, p.created_at,
	EXISTS (SELECT 1 FROM product_media m WHERE m.product_id = p.id AND m.kind = 'thumb'),
	EXISTS (SELECT 1 FROM product_media m WHERE m.product_id = p.id AND m.kind = 'viewer_glb')`

func scanCards(rows pgx.Rows) ([]Card, error) {
	defer rows.Close()
	cards := []Card{}
	for rows.Next() {
		var c Card
		var hasThumb, hasViewer bool
		if err := rows.Scan(&c.ID, &c.Slug, &c.Title, &c.Description, &c.PriceCents, &c.Currency, &c.Tier,
			&c.Status, &c.Attributes, &c.CreatedAt, &hasThumb, &hasViewer); err != nil {
			return nil, err
		}
		c.ThumbURL, c.ViewerURL = MediaURLs(c.ID, hasThumb, hasViewer)
		cards = append(cards, c)
	}
	return cards, rows.Err()
}

// List returns listed products matching the filters. Filter keys are
// restricted to Facets and values are always passed as parameters.
func List(ctx context.Context, pool *pgxpool.Pool, q Query) ([]Card, error) {
	conds := []string{"p.status = 'listed'"}
	var args []any
	arg := func(v any) string {
		args = append(args, v)
		return fmt.Sprintf("$%d", len(args))
	}
	if q.Tier == "free" || q.Tier == "paid" {
		conds = append(conds, "p.tier = "+arg(q.Tier))
	}
	if t := strings.TrimSpace(q.Text); t != "" {
		like := arg("%" + t + "%")
		conds = append(conds, fmt.Sprintf("(p.title ILIKE %s OR p.description ILIKE %s OR p.attributes::text ILIKE %s)", like, like, like))
	}
	for key, raw := range q.Filters {
		switch facetKind(key) {
		case "set":
			var values []string
			if json.Unmarshal(raw, &values) != nil || len(values) == 0 {
				continue
			}
			k, v := arg(key), arg(values)
			conds = append(conds, fmt.Sprintf("(p.attributes->>%s = ANY(%s) OR p.attributes->%s ?| %s)", k, v, k, v))
		case "range":
			var r struct{ Min, Max *float64 }
			if json.Unmarshal(raw, &r) != nil {
				continue
			}
			k := arg(key)
			expr := fmt.Sprintf("(p.attributes->>%s)::numeric", k)
			conds = append(conds, fmt.Sprintf("jsonb_typeof(p.attributes->%s) = 'number'", k))
			if r.Min != nil {
				conds = append(conds, expr+" >= "+arg(*r.Min))
			}
			if r.Max != nil {
				conds = append(conds, expr+" <= "+arg(*r.Max))
			}
		case "bool":
			var b bool
			if json.Unmarshal(raw, &b) != nil || !b {
				continue
			}
			conds = append(conds, fmt.Sprintf("p.attributes->%s = 'true'::jsonb", arg(key)))
		}
	}
	limit := q.Limit
	if limit <= 0 || limit > 200 {
		limit = 120
	}
	sql := fmt.Sprintf("SELECT %s FROM products p WHERE %s ORDER BY p.created_at DESC, p.id DESC LIMIT %d",
		cardColumns, strings.Join(conds, " AND "), limit)
	rows, err := pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return scanCards(rows)
}

// AdminList = all products including draft/archived.
func AdminList(ctx context.Context, pool *pgxpool.Pool) ([]Card, error) {
	rows, err := pool.Query(ctx, "SELECT "+cardColumns+" FROM products p ORDER BY p.created_at DESC, p.id DESC")
	if err != nil {
		return nil, err
	}
	return scanCards(rows)
}

// CardsByIDs preserves the order of ids (relevance order from vector search).
func CardsByIDs(ctx context.Context, pool *pgxpool.Pool, ids []int64) ([]Card, error) {
	rows, err := pool.Query(ctx, "SELECT "+cardColumns+" FROM products p WHERE p.id = ANY($1)", ids)
	if err != nil {
		return nil, err
	}
	cards, err := scanCards(rows)
	if err != nil {
		return nil, err
	}
	byID := map[int64]Card{}
	for _, c := range cards {
		byID[c.ID] = c
	}
	out := make([]Card, 0, len(ids))
	for _, id := range ids {
		if c, ok := byID[id]; ok {
			out = append(out, c)
		}
	}
	return out, nil
}

// LiveFacets computes facet values/ranges from the currently listed products.
// Facets without values are omitted.
func LiveFacets(ctx context.Context, pool *pgxpool.Pool) ([]Facet, error) {
	out := []Facet{}
	for _, f := range Facets {
		f.Config = map[string]any{}
		switch f.Kind {
		case "set":
			rows, err := pool.Query(ctx, `
				SELECT DISTINCT v FROM products p,
				LATERAL jsonb_array_elements_text(
					CASE WHEN jsonb_typeof(p.attributes->$1) = 'array' THEN p.attributes->$1
					     ELSE jsonb_build_array(p.attributes->>$1) END) AS v
				WHERE p.status = 'listed' AND v IS NOT NULL AND v <> ''
				ORDER BY v`, f.Key)
			if err != nil {
				return nil, err
			}
			values, err := pgx.CollectRows(rows, pgx.RowTo[string])
			if err != nil {
				return nil, err
			}
			if len(values) == 0 {
				continue
			}
			f.Config["values"] = values
		case "range":
			var lo, hi *float64
			err := pool.QueryRow(ctx, `
				SELECT min((attributes->>$1)::numeric)::float8, max((attributes->>$1)::numeric)::float8
				FROM products WHERE status = 'listed' AND jsonb_typeof(attributes->$1) = 'number'`, f.Key).Scan(&lo, &hi)
			if err != nil {
				return nil, err
			}
			if lo == nil {
				continue
			}
			f.Config["min"], f.Config["max"] = *lo, *hi
		case "bool":
			var exists bool
			if err := pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM products
				WHERE status = 'listed' AND jsonb_typeof(attributes->$1) = 'boolean')`, f.Key).Scan(&exists); err != nil {
				return nil, err
			}
			if !exists {
				continue
			}
		}
		out = append(out, f)
	}
	return out, nil
}

// EmbeddingText = what gets indexed about a product for chat.
func EmbeddingText(p db.Product) string {
	var attrs map[string]any
	_ = json.Unmarshal(p.Attributes, &attrs)
	parts := []string{p.Title, p.Description}
	for _, k := range []string{"category", "style", "collection", "prompt"} {
		if s, ok := attrs[k].(string); ok && s != "" {
			parts = append(parts, s)
		}
	}
	if tags, ok := attrs["tags"].([]any); ok {
		for _, t := range tags {
			parts = append(parts, fmt.Sprint(t))
		}
	}
	return strings.Join(parts, " | ")
}

// Index recomputes a product's embedding only when the indexed text has changed.
func Index(ctx context.Context, q *db.Queries, emb ai.Embedder, p db.Product) error {
	text := EmbeddingText(p)
	sum := sha256.Sum256([]byte(text))
	hash := hex.EncodeToString(sum[:])
	prev, err := q.GetEmbeddingHash(ctx, p.ID)
	if err == nil && prev == hash {
		return nil
	}
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}
	vec, err := emb.Embed(ctx, text)
	if err != nil {
		return err
	}
	return q.UpsertEmbedding(ctx, db.UpsertEmbeddingParams{ProductID: p.ID, Embedding: pgvector.NewVector(vec), TextHash: hash})
}
