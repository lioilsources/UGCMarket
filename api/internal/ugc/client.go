// Package ugc is a client for ugc-backend (UGC Factory, JODA :8095).
// Market only reads from it — see docs/UGC_API.md. On the LAN/compose network the API is
// unauthenticated; the optional CF Access headers are for running outside JODA.
package ugc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Job struct {
	ID         string          `json:"id"`
	Status     string          `json:"status"`
	Prompt     string          `json:"prompt"`
	Category   string          `json:"category"`
	Style      string          `json:"style"`
	Backend    string          `json:"backend"`
	Collection string          `json:"collection"`
	Symmetry   string          `json:"symmetry"`
	Verdict    string          `json:"verdict"`
	Report     json.RawMessage `json:"report"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

type Item struct {
	ID            string          `json:"id"`
	Name          string          `json:"name"`
	DescriptionCS string          `json:"description_cs"`
	DescriptionEN string          `json:"description_en"`
	Category      string          `json:"category"`
	Tags          json.RawMessage `json:"tags"`
	Collection    string          `json:"collection"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

// Report is the output of convert.py (blender_scripts/convert.py) — only the fields market shows.
type Report struct {
	TriCount        int       `json:"tri_count"`
	TextureSize     int       `json:"texture_size"`
	UVOK            bool      `json:"uv_ok"`
	Watertight      bool      `json:"watertight"`
	MaterialCount   int       `json:"material_count"`
	Attachment      string    `json:"attachment"`
	BBoxTargetStuds []float64 `json:"bbox_target_studs"`
	VerdictReasons  []string  `json:"verdict_reasons"`
}

type Client struct {
	base string
	http *http.Client
}

func New(base string) *Client {
	return &Client{base: base, http: &http.Client{Timeout: 5 * time.Minute}}
}

func (c *Client) getJSON(ctx context.Context, path string, out any) error {
	body, err := c.Open(ctx, path)
	if err != nil {
		return err
	}
	defer body.Close()
	return json.NewDecoder(body).Decode(out)
}

// Open returns the response body; the caller must close it.
func (c *Client) Open(ctx context.Context, path string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.base+path, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("ugc GET %s: %s", path, resp.Status)
	}
	return resp.Body, nil
}

// PackedJobs = pieces the owner approved a second time in triage (packed).
func (c *Client) PackedJobs(ctx context.Context) ([]Job, error) {
	var jobs []Job
	return jobs, c.getJSON(ctx, "/jobs?status=packed&limit=1000", &jobs)
}

func (c *Client) Items(ctx context.Context) ([]Item, error) {
	var items []Item
	return items, c.getJSON(ctx, "/items", &items)
}
