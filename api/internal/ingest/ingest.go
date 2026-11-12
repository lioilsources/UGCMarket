// Package ingest pulls finished pieces from ugc-backend into the market catalog.
//
// “Finished” = a job in the packed state (the owner approved it a second time in triage) and
// the conversion verdict is not FAIL. A new piece is created as a draft — only an
// admin lists it. Files are copied into our own storage because reject/reconvert
// moves and overwrites them in ugc. A piece that disappears from packed is archived
// (attributes.ugc_missing) and goes back to draft when it returns.
package ingest

import (
	"archive/zip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/lioilsources/RobloxMarket/api/internal/ai"
	"github.com/lioilsources/RobloxMarket/api/internal/catalog"
	"github.com/lioilsources/RobloxMarket/api/internal/db"
	"github.com/lioilsources/RobloxMarket/api/internal/storage"
	"github.com/lioilsources/RobloxMarket/api/internal/ugc"
)

type Result struct {
	StartedAt   time.Time `json:"started_at"`
	FinishedAt  time.Time `json:"finished_at"`
	Seen        int       `json:"seen"`
	Created     int       `json:"created"`
	Updated     int       `json:"updated"`
	Unchanged   int       `json:"unchanged"`
	SkippedFail int       `json:"skipped_fail"`
	Archived    int       `json:"archived"`
	Errors      []string  `json:"errors"`
}

type Ingester struct {
	pool       *pgxpool.Pool
	q          *db.Queries
	fs         *storage.FS
	ugc        *ugc.Client
	emb        ai.Embedder
	priceCents int32
	currency   string

	mu   sync.Mutex // one run at a time (poller vs. manual trigger from admin)
	last *Result
}

func New(pool *pgxpool.Pool, fs *storage.FS, client *ugc.Client, emb ai.Embedder, priceCents int32, currency string) *Ingester {
	return &Ingester{pool: pool, q: db.New(pool), fs: fs, ugc: client, emb: emb, priceCents: priceCents, currency: currency}
}

func (in *Ingester) Last() *Result {
	in.mu.Lock()
	defer in.mu.Unlock()
	return in.last
}

// Poll runs Run at an interval until ctx is done.
func (in *Ingester) Poll(ctx context.Context, every time.Duration) {
	for {
		res, err := in.Run(ctx)
		if err != nil {
			log.Printf("ingest: %v", err)
		} else if res.Created+res.Updated+res.Archived > 0 || len(res.Errors) > 0 {
			log.Printf("ingest: new=%d updated=%d archived=%d errors=%d", res.Created, res.Updated, res.Archived, len(res.Errors))
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(every):
		}
	}
}

func (in *Ingester) Run(ctx context.Context) (*Result, error) {
	in.mu.Lock()
	defer in.mu.Unlock()
	res := &Result{StartedAt: time.Now(), Errors: []string{}}

	jobs, err := in.ugc.PackedJobs(ctx)
	if err != nil {
		return nil, err
	}
	items, err := in.ugc.Items(ctx)
	if err != nil {
		return nil, err
	}
	itemByID := map[string]ugc.Item{}
	for _, it := range items {
		itemByID[it.ID] = it
	}
	existing, err := in.q.ListIngestedProducts(ctx)
	if err != nil {
		return nil, err
	}
	byRef := map[string]db.ListIngestedProductsRow{}
	for _, p := range existing {
		byRef[p.AssetRef] = p
	}

	present := map[string]bool{}
	for _, job := range jobs {
		res.Seen++
		if job.Verdict == "FAIL" {
			res.SkippedFail++
			continue
		}
		present[job.ID] = true
		item := itemByID[job.ID]
		if err := in.upsert(ctx, job, item, res); err != nil {
			res.Errors = append(res.Errors, fmt.Sprintf("%s: %v", job.ID, err))
		}
	}

	for _, p := range existing {
		if present[p.AssetRef] || p.Status == "archived" {
			continue
		}
		attrs := withAttr(p.Attributes, "ugc_missing", true)
		full, err := in.q.GetProductByID(ctx, p.ID)
		if err == nil {
			err = in.q.UpdateProductSnapshot(ctx, db.UpdateProductSnapshotParams{
				ID: p.ID, Attributes: attrs, SourceUpdatedAt: full.SourceUpdatedAt, Status: "archived"})
		}
		if err != nil {
			res.Errors = append(res.Errors, fmt.Sprintf("archive %s: %v", p.AssetRef, err))
			continue
		}
		res.Archived++
	}

	res.FinishedAt = time.Now()
	in.last = res
	return res, nil
}

func (in *Ingester) upsert(ctx context.Context, job ugc.Job, item ugc.Item, res *Result) error {
	srcUpdated := job.UpdatedAt
	if item.UpdatedAt.After(srcUpdated) {
		srcUpdated = item.UpdatedAt
	}
	p, err := in.q.GetProductByAssetRef(ctx, job.ID)
	isNew := errors.Is(err, pgx.ErrNoRows)
	if err != nil && !isNew {
		return err
	}
	if !isNew {
		missing := attrBool(p.Attributes, "ugc_missing")
		if !missing && p.SourceUpdatedAt != nil && p.SourceUpdatedAt.Equal(srcUpdated) {
			res.Unchanged++
			return nil
		}
	}

	if isNew {
		title := titleFrom(item.Name, job.Prompt, job.ID)
		slug, err := in.uniqueSlug(ctx, slugify(title), job.ID)
		if err != nil {
			return err
		}
		p, err = in.q.InsertProduct(ctx, db.InsertProductParams{
			Slug: slug, Title: title, Description: item.DescriptionEN,
			PriceCents: in.priceCents, Currency: in.currency, Tier: "paid",
			AssetRef: job.ID, Attributes: json.RawMessage(`{}`), SourceUpdatedAt: nil,
		})
		if err != nil {
			return err
		}
	}

	formats, err := in.copyFiles(ctx, p.ID, job.ID)
	if err != nil {
		return err
	}
	attrs, err := json.Marshal(attributes(job, item, formats))
	if err != nil {
		return err
	}
	status := p.Status
	if status == "archived" && attrBool(p.Attributes, "ugc_missing") {
		status = "draft" // the piece returned to packed (typically after a reconvert) — admin will check it
	}
	if err := in.q.UpdateProductSnapshot(ctx, db.UpdateProductSnapshotParams{
		ID: p.ID, Attributes: attrs, SourceUpdatedAt: &srcUpdated, Status: status}); err != nil {
		return err
	}
	if p, err = in.q.GetProductByID(ctx, p.ID); err != nil {
		return err
	}
	if err := catalog.Index(ctx, in.q, in.emb, p); err != nil {
		// without an embedding the product just doesn't show up in chat — ingest doesn't fail because of it
		res.Errors = append(res.Errors, fmt.Sprintf("%s: embedding: %v", job.ID, err))
	}
	if isNew {
		res.Created++
	} else {
		res.Updated++
	}
	return nil
}

// copyFiles downloads previews and files for sale. Returns the list of sale formats.
func (in *Ingester) copyFiles(ctx context.Context, productID int64, jobID string) ([]string, error) {
	dir := fmt.Sprintf("products/%d", productID)
	put := func(src, key string) (int64, string, error) {
		body, err := in.ugc.Open(ctx, src)
		if err != nil {
			return 0, "", err
		}
		defer body.Close()
		return in.fs.Put(key, body)
	}

	// Previews: preview.glb (512px textures) for the 3D viewer, 2D concept as thumbnail.
	// The concept is optional (missing after reject/without preview) — its absence is not an error.
	if _, _, err := put("/jobs/"+jobID+"/glb", dir+"/viewer.glb"); err != nil {
		return nil, fmt.Errorf("viewer glb: %w", err)
	}
	if err := in.q.UpsertProductMedia(ctx, db.UpsertProductMediaParams{ProductID: productID, Kind: "viewer_glb", StorageKey: dir + "/viewer.glb"}); err != nil {
		return nil, err
	}
	if _, _, err := put("/jobs/"+jobID+"/preview", dir+"/thumb.png"); err == nil {
		if err := in.q.UpsertProductMedia(ctx, db.UpsertProductMediaParams{ProductID: productID, Kind: "thumb", StorageKey: dir + "/thumb.png"}); err != nil {
			return nil, err
		}
	}

	var formats []string
	addFile := func(format, key string, n int64, sum string) error {
		formats = append(formats, format)
		return in.q.UpsertProductFile(ctx, db.UpsertProductFileParams{ProductID: productID, Format: format, StorageKey: key, Bytes: n, Sha256: sum})
	}

	// Roblox-ready FBX + baked texture from the packed bundle.
	tmp, err := os.CreateTemp("", "ugc-pack-*.zip")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmp.Name())
	defer tmp.Close()
	body, err := in.ugc.Open(ctx, "/packed/"+jobID+"/download")
	if err != nil {
		return nil, fmt.Errorf("packed zip: %w", err)
	}
	size, err := io.Copy(tmp, body)
	body.Close()
	if err != nil {
		return nil, err
	}
	zr, err := zip.NewReader(tmp, size)
	if err != nil {
		return nil, fmt.Errorf("packed zip: %w", err)
	}
	for _, f := range zr.File {
		name := strings.ToLower(path.Base(f.Name))
		var format, key string
		switch {
		case strings.HasSuffix(name, ".fbx"):
			format, key = "fbx", dir+"/model.fbx"
		case strings.HasSuffix(name, ".png"):
			format, key = "png", dir+"/model_tex.png"
		default:
			continue // item.json is internal (prices in Robux etc.)
		}
		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		n, sum, err := in.fs.Put(key, rc)
		rc.Close()
		if err != nil {
			return nil, err
		}
		if err := addFile(format, key, n, sum); err != nil {
			return nil, err
		}
	}
	if len(formats) == 0 {
		return nil, errors.New("packed zip neobsahuje FBX")
	}

	// Raw GLB (higher detail, original scale) as a universal format.
	n, sum, err := put("/jobs/"+jobID+"/glb?full=1", dir+"/model.glb")
	if err != nil {
		return nil, fmt.Errorf("glb: %w", err)
	}
	if err := addFile("glb", dir+"/model.glb", n, sum); err != nil {
		return nil, err
	}
	return formats, nil
}

func attributes(job ugc.Job, item ugc.Item, formats []string) map[string]any {
	var rep ugc.Report
	_ = json.Unmarshal(job.Report, &rep)
	tags := []string{}
	_ = json.Unmarshal(item.Tags, &tags)
	// Sale formats = what the buyer actually finds in the bundle (png is the texture for the FBX).
	sale := []string{}
	for _, f := range formats {
		if f != "png" {
			sale = append(sale, f)
		}
	}
	a := map[string]any{
		"ugc_id":            job.ID,
		"prompt":            job.Prompt,
		"category":          job.Category,
		"style":             job.Style,
		"collection":        job.Collection,
		"backend":           job.Backend,
		"source_pipeline":   job.Backend,
		"verdict":           job.Verdict,
		"tags":              tags,
		"formats":           sale,
		"polycount":         rep.TriCount,
		"texture_px":        rep.TextureSize,
		"uv_ok":             rep.UVOK,
		"watertight":        rep.Watertight,
		"materials":         rep.MaterialCount,
		"roblox_attachment": rep.Attachment,
		"bbox_studs":        rep.BBoxTargetStuds,
		"verdict_reasons":   rep.VerdictReasons,
	}
	for k, v := range a {
		if s, ok := v.(string); ok && s == "" {
			delete(a, k)
		}
	}
	return a
}

func (in *Ingester) uniqueSlug(ctx context.Context, base, id string) (string, error) {
	if base == "" {
		base = "asset"
	}
	slug := base
	for i := 0; ; i++ {
		taken, err := in.q.SlugTaken(ctx, slug)
		if err != nil || !taken {
			return slug, err
		}
		suffix := strings.ToLower(id)
		if len(suffix) > 6 {
			suffix = suffix[len(suffix)-6:]
		}
		slug = fmt.Sprintf("%s-%s", base, suffix)
		if i > 0 {
			slug = fmt.Sprintf("%s-%s-%d", base, suffix, i)
		}
	}
}

// titleFrom: a manually renamed item in ugc takes precedence; otherwise the title
// is derived from the prompt (item.name in ugc is only the first 4 words of the prompt, typically
// truncated like “a-cute-classic-school”). Up to the first comma, without an article, max 7 words.
func titleFrom(name, prompt, id string) string {
	src := name
	if prompt != "" && (name == "" || name == autoName(prompt)) {
		src, _, _ = strings.Cut(prompt, ",")
	}
	words := strings.FieldsFunc(src, func(r rune) bool { return r == '-' || r == '_' || r == ' ' })
	if len(words) > 1 {
		switch strings.ToLower(words[0]) {
		case "a", "an", "the":
			words = words[1:]
		}
	}
	if len(words) > 7 {
		words = words[:7]
	}
	small := map[string]bool{"with": true, "in": true, "of": true, "and": true, "on": true, "for": true, "a": true, "an": true, "the": true}
	for i, w := range words {
		if i > 0 && small[strings.ToLower(w)] {
			words[i] = strings.ToLower(w)
			continue
		}
		r := []rune(w)
		words[i] = strings.ToUpper(string(r[:1])) + string(r[1:])
	}
	// don't end on a conjunction (“… Sword with”)
	for len(words) > 1 && small[words[len(words)-1]] {
		words = words[:len(words)-1]
	}
	if len(words) == 0 {
		return id
	}
	return strings.Join(words, " ")
}

// autoName = ugc safeName (nas/main.go): first 4 words of the prompt, a-z0-9 and hyphens.
func autoName(prompt string) string {
	words := strings.Fields(strings.ToLower(prompt))
	if len(words) > 4 {
		words = words[:4]
	}
	var b strings.Builder
	for _, r := range strings.Join(words, "-") {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == '-' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func slugify(s string) string {
	var b strings.Builder
	dash := false
	for _, r := range strings.ToLower(s) {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			b.WriteRune(r)
			dash = false
		} else if !dash && b.Len() > 0 {
			b.WriteByte('-')
			dash = true
		}
	}
	return strings.Trim(b.String(), "-")
}

func attrBool(raw json.RawMessage, key string) bool {
	var m map[string]any
	_ = json.Unmarshal(raw, &m)
	b, _ := m[key].(bool)
	return b
}

func withAttr(raw json.RawMessage, key string, v any) json.RawMessage {
	m := map[string]any{}
	_ = json.Unmarshal(raw, &m)
	m[key] = v
	out, _ := json.Marshal(m)
	return out
}
