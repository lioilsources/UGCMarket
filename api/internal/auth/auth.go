// Package auth: sign-in via e-mail magic link and session cookie.
// Account = e-mail; the customer is created on first link verification.
package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"net/url"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/lioilsources/RobloxMarket/api/internal/db"
	mailer "github.com/lioilsources/RobloxMarket/api/internal/mail"
)

const (
	CookieName = "om_session"
	linkTTL    = 30 * time.Minute
	sessionTTL = 60 * 24 * time.Hour
)

var ErrInvalidLink = errors.New("the sign-in link is invalid or expired")

// NewToken returns a random token (for URL/cookie) and its hash (for the DB).
func NewToken() (token, hash string) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	token = base64.RawURLEncoding.EncodeToString(b)
	return token, Hash(token)
}

func Hash(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

type Service struct {
	q         *db.Queries
	mail      mailer.Mailer
	publicURL string
	secure    bool
}

func New(q *db.Queries, m mailer.Mailer, publicURL string, secureCookie bool) *Service {
	return &Service{q: q, mail: m, publicURL: publicURL, secure: secureCookie}
}

func NormalizeEmail(raw string) (string, error) {
	addr, err := mail.ParseAddress(strings.TrimSpace(raw))
	if err != nil || addr.Name != "" || len(addr.Address) > 254 {
		return "", errors.New("invalid e-mail address")
	}
	return strings.ToLower(addr.Address), nil
}

// SendLink sends a magic link. `next` is the relative path to return to after sign-in.
func (s *Service) SendLink(ctx context.Context, email, next string) error {
	token, hash := NewToken()
	if err := s.q.CreateLoginToken(ctx, db.CreateLoginTokenParams{TokenHash: hash, Email: email, ExpiresAt: time.Now().Add(linkTTL)}); err != nil {
		return err
	}
	link := fmt.Sprintf("%s/auth/verify?token=%s", s.publicURL, token)
	if next = SafeNext(next); next != "/" {
		link += "&next=" + url.QueryEscape(next)
	}
	text := fmt.Sprintf("Sign in to Ol1n Market:\n\n%s\n\nThe link is valid for 30 minutes and works once.\nIf you didn't request it, ignore this e-mail.\n", link)
	return s.mail.Send(ctx, email, "Your Ol1n Market sign-in link", text)
}

// Verify consumes a magic link, creates the customer and returns a session token.
func (s *Service) Verify(ctx context.Context, token string) (db.Customer, string, error) {
	email, err := s.q.ConsumeLoginToken(ctx, Hash(token))
	if errors.Is(err, pgx.ErrNoRows) {
		return db.Customer{}, "", ErrInvalidLink
	}
	if err != nil {
		return db.Customer{}, "", err
	}
	cust, err := s.q.UpsertCustomer(ctx, email)
	if err != nil {
		return db.Customer{}, "", err
	}
	session, hash := NewToken()
	if err := s.q.CreateSession(ctx, db.CreateSessionParams{TokenHash: hash, CustomerID: cust.ID, ExpiresAt: time.Now().Add(sessionTTL)}); err != nil {
		return db.Customer{}, "", err
	}
	return cust, session, nil
}

func (s *Service) SetCookie(w http.ResponseWriter, session string) {
	http.SetCookie(w, &http.Cookie{Name: CookieName, Value: session, Path: "/", HttpOnly: true,
		Secure: s.secure, SameSite: http.SameSiteLaxMode, MaxAge: int(sessionTTL.Seconds())})
}

func (s *Service) ClearCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{Name: CookieName, Value: "", Path: "/", HttpOnly: true,
		Secure: s.secure, SameSite: http.SameSiteLaxMode, MaxAge: -1})
}

// Customer returns the signed-in customer, or nil.
func (s *Service) Customer(r *http.Request) *db.Customer {
	c, err := r.Cookie(CookieName)
	if err != nil || c.Value == "" {
		return nil
	}
	cust, err := s.q.GetSessionCustomer(r.Context(), Hash(c.Value))
	if err != nil {
		return nil
	}
	return &cust
}

func (s *Service) Logout(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie(CookieName); err == nil {
		_ = s.q.DeleteSession(r.Context(), Hash(c.Value))
	}
	s.ClearCookie(w)
}

// SafeNext allows only a relative path on our own site (no open redirect).
func SafeNext(next string) string {
	if !strings.HasPrefix(next, "/") || strings.HasPrefix(next, "//") || strings.ContainsAny(next, "\\\r\n") {
		return "/"
	}
	return next
}
