// Package mail sends transactional e-mails (magic link, order confirmation).
// The shop had no mailer; Resend = a single HTTP POST, no SMTP server.
// Without an API key the e-mail is just logged (dev / before the domain is set up).
package mail

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Mailer interface {
	Send(ctx context.Context, to, subject, text string) error
}

func New(resendKey, from string) Mailer {
	if resendKey == "" {
		return LogMailer{}
	}
	return &Resend{key: resendKey, from: from, http: &http.Client{Timeout: 20 * time.Second}}
}

type LogMailer struct{}

func (LogMailer) Send(_ context.Context, to, subject, text string) error {
	log.Printf("mail (log only) to=%s subject=%q\n%s", to, subject, text)
	return nil
}

type Resend struct {
	key, from string
	http      *http.Client
}

func (r *Resend) Send(ctx context.Context, to, subject, text string) error {
	body, _ := json.Marshal(map[string]any{"from": r.from, "to": []string{to}, "subject": subject, "text": text})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.resend.com/emails", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+r.key)
	req.Header.Set("Content-Type", "application/json")
	resp, err := r.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("resend: %s", resp.Status)
	}
	return nil
}
