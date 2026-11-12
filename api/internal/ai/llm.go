package ai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type LLM struct {
	base, model, key string
	http             *http.Client
}

// NewLLM returns nil when no URL is set — the caller then runs in mock mode.
func NewLLM(baseURL, model, key string) *LLM {
	if baseURL == "" {
		return nil
	}
	return &LLM{base: baseURL, model: model, key: key, http: &http.Client{}}
}

// Stream calls /chat/completions with stream=true and invokes onDelta for each chunk of text.
func (l *LLM) Stream(ctx context.Context, messages []Message, onDelta func(string) error) error {
	body, _ := json.Marshal(map[string]any{"model": l.model, "messages": messages, "stream": true})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, l.base+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+l.key)
	resp, err := l.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("llm: %s", resp.Status)
	}
	sc := bufio.NewScanner(resp.Body)
	sc.Buffer(make([]byte, 64*1024), 1024*1024)
	for sc.Scan() {
		data, ok := strings.CutPrefix(sc.Text(), "data: ")
		if !ok {
			continue
		}
		if data == "[DONE]" {
			break
		}
		var chunk struct {
			Choices []struct {
				Delta struct {
					Content string `json:"content"`
				} `json:"delta"`
			} `json:"choices"`
		}
		if json.Unmarshal([]byte(data), &chunk) != nil || len(chunk.Choices) == 0 {
			continue
		}
		if d := chunk.Choices[0].Delta.Content; d != "" {
			if err := onDelta(d); err != nil {
				return err
			}
		}
	}
	return sc.Err()
}
