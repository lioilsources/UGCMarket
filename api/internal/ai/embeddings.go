// Package ai: embeddings and chat via an OpenAI-compatible API (AiStack LiteLLM
// gateway on Spark: models "shop" and "shop-embed" = bge-m3, 1024 dim).
// Without a configured URL a mock runs — ported from the shop (services/brain/brain/embeddings.py).
package ai

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const Dim = 1024

type Embedder interface {
	Embed(ctx context.Context, text string) ([]float32, error)
}

func NewEmbedder(baseURL, model string) Embedder {
	if baseURL == "" {
		return MockEmbedder{}
	}
	return &HTTPEmbedder{base: baseURL, model: model, http: &http.Client{Timeout: 60 * time.Second}}
}

type HTTPEmbedder struct {
	base, model string
	http        *http.Client
}

func (e *HTTPEmbedder) Embed(ctx context.Context, text string) ([]float32, error) {
	// set encoding_format explicitly: LiteLLM otherwise sends None and vLLM validation rejects it
	body, _ := json.Marshal(map[string]any{"model": e.model, "input": text, "encoding_format": "float"})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, e.base+"/embeddings", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := e.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("embeddings: %s", resp.Status)
	}
	var out struct {
		Data []struct {
			Embedding []float32 `json:"embedding"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	if len(out.Data) == 0 || len(out.Data[0].Embedding) != Dim {
		return nil, fmt.Errorf("embeddings: čekám 1 vektor o %d dimenzích", Dim)
	}
	return out.Data[0].Embedding, nil
}

// MockEmbedder builds a vector from hash-seeded word vectors — overlapping
// words yield similar vectors, so similarity search works even without a GPU.
type MockEmbedder struct{}

var wordRe = regexp.MustCompile(`[\p{L}\p{N}_]+`)

func (MockEmbedder) Embed(_ context.Context, text string) ([]float32, error) {
	words := wordRe.FindAllString(strings.ToLower(text), -1)
	if len(words) == 0 {
		words = []string{"empty"}
	}
	acc := make([]float64, Dim)
	for _, w := range words {
		for i, v := range wordVector(w) {
			acc[i] += v
		}
	}
	var norm float64
	for _, v := range acc {
		norm += v * v
	}
	norm = math.Sqrt(norm)
	if norm == 0 {
		norm = 1
	}
	out := make([]float32, Dim)
	for i, v := range acc {
		out[i] = float32(v / norm)
	}
	return out, nil
}

func wordVector(word string) []float64 {
	out := make([]float64, 0, Dim)
	var ctr [4]byte
	for counter := uint32(0); len(out) < Dim; counter++ {
		binary.BigEndian.PutUint32(ctr[:], counter)
		digest := sha256.Sum256(append([]byte(word), ctr[:]...))
		for _, b := range digest {
			out = append(out, float64(b)/255.0-0.5)
		}
	}
	return out[:Dim]
}
