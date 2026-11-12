package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pgvector/pgvector-go"

	"github.com/lioilsources/RobloxMarket/api/internal/ai"
	"github.com/lioilsources/RobloxMarket/api/internal/catalog"
	"github.com/lioilsources/RobloxMarket/api/internal/db"
)

// Port of the shop (services/brain/brain/chat.py). SSE events:
//
//	products [Card, ...]   — relevant assets (vector search)
//	token    {"text": ...} — chunk of the answer
//	done     {}
const systemPrompt = `You are the shopping assistant of Ol1n Market, a store selling 3D models
(FBX/GLB) made for Roblox Studio and other engines. Answer briefly and in the user's language.
Below are assets from the store relevant to the question. Recommend the ones that fit and
mention their name and price. If nothing fits, say so honestly — never invent assets.
Ol1n Market is not affiliated with Roblox Corporation; never claim otherwise.`

const (
	maxChatMessages = 20
	maxChatChars    = 2000
)

func (s *Server) handleChat(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Messages []ai.Message `json:"messages"`
	}
	if err := decode(r, &req); err != nil || len(req.Messages) == 0 {
		httpErr(w, http.StatusBadRequest, "expected {messages: [...]}")
		return
	}
	if len(req.Messages) > maxChatMessages {
		req.Messages = req.Messages[len(req.Messages)-maxChatMessages:]
	}
	var query string
	msgs := make([]ai.Message, 0, len(req.Messages))
	for _, m := range req.Messages {
		if m.Role != "user" && m.Role != "assistant" {
			continue
		}
		if len(m.Content) > maxChatChars {
			m.Content = m.Content[:maxChatChars]
		}
		if m.Role == "user" {
			query = m.Content
		}
		msgs = append(msgs, m)
	}

	ctx := r.Context()
	cards, err := s.searchCards(r, query, 5)
	if err != nil {
		serverErr(w, err)
		return
	}

	flusher, _ := w.(http.Flusher)
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("X-Accel-Buffering", "no")
	send := func(event string, data any) error {
		b, _ := json.Marshal(data)
		if _, err := fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event, b); err != nil {
			return err
		}
		if flusher != nil {
			flusher.Flush()
		}
		return nil
	}

	_ = send("products", cards)
	token := func(t string) error { return send("token", map[string]string{"text": t}) }
	if s.llm != nil {
		full := append([]ai.Message{{Role: "system", Content: systemPrompt + "\n\nAssets:\n" + productContext(cards)}}, msgs...)
		if err := s.llm.Stream(ctx, full, token); err != nil && ctx.Err() == nil {
			_ = token("\n\n(The assistant is unavailable right now — the matching assets are shown above.)")
		}
	} else {
		answer := mockAnswer(query, cards)
		for i := 0; i < len(answer); i += 12 {
			if token(answer[i:min(i+12, len(answer))]) != nil {
				return
			}
			time.Sleep(20 * time.Millisecond)
		}
	}
	_ = send("done", struct{}{})
}

func (s *Server) searchCards(r *http.Request, query string, limit int32) ([]catalog.Card, error) {
	ctx := r.Context()
	vec, err := s.emb.Embed(ctx, query)
	if err != nil {
		return nil, err
	}
	products, err := s.q.SearchListedProducts(ctx, db.SearchListedProductsParams{Query: pgvector.NewVector(vec), Lim: limit})
	if err != nil {
		return nil, err
	}
	ids := make([]int64, len(products))
	for i, p := range products {
		ids[i] = p.ID
	}
	return catalog.CardsByIDs(ctx, s.pool, ids)
}

func productContext(cards []catalog.Card) string {
	var b strings.Builder
	for _, c := range cards {
		var attrs map[string]any
		_ = json.Unmarshal(c.Attributes, &attrs)
		fmt.Fprintf(&b, "- %s | %s | category: %v | style: %v | triangles: %v | formats: %v",
			c.Title, priceLabel(c), attrs["category"], attrs["style"], attrs["polycount"], attrs["formats"])
		if c.Description != "" {
			fmt.Fprintf(&b, " | %.200s", c.Description)
		}
		b.WriteByte('\n')
	}
	return b.String()
}

func priceLabel(c catalog.Card) string {
	if c.Tier == "free" {
		return "free"
	}
	return fmt.Sprintf("%d.%02d %s", c.PriceCents/100, c.PriceCents%100, c.Currency)
}

func mockAnswer(query string, cards []catalog.Card) string {
	if len(cards) == 0 {
		return "I couldn't find anything matching that in the store. Try describing it differently."
	}
	names := make([]string, 0, 3)
	for _, c := range cards[:min(3, len(cards))] {
		names = append(names, fmt.Sprintf("%s (%s)", c.Title, priceLabel(c)))
	}
	return fmt.Sprintf("(mock mode — no LLM connected) For “%s” I found: %s. The asset cards are below.",
		query, strings.Join(names, ", "))
}
