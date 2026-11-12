// Package config reads market-api settings from environment variables.
// Empty values for optional integrations = mock/disabled (as in the shop):
// without LLM_BASE_URL chat answers from a template, without RESEND_API_KEY e-mails are only logged.
package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Addr        string
	DatabaseURL string
	DataDir     string // asset file storage (copy from ugc-backend)
	PublicURL   string // https://market.ol1n.com — for links in e-mails

	UGCURL          string        // http://ugc-api:8095; empty = ingest disabled
	UGCPollInterval time.Duration // 0 = manual ingest via admin only

	DefaultPriceCents int32 // uniform price of newly ingested assets
	Currency          string

	AdminEmails []string // customers with these e-mails can see /admin

	LLMBaseURL        string
	LLMModel          string
	LLMAPIKey         string
	EmbeddingsBaseURL string
	EmbeddingsModel   string

	ResendAPIKey string
	MailFrom     string

	CookieSecure bool
}

func Load() Config {
	return Config{
		Addr:              env("MARKET_ADDR", ":8096"),
		DatabaseURL:       env("DATABASE_URL", "postgres://market:market@localhost:5435/market?sslmode=disable"),
		DataDir:           env("MARKET_DATA", "./data"),
		PublicURL:         strings.TrimRight(env("PUBLIC_URL", "http://localhost:5173"), "/"),
		UGCURL:            strings.TrimRight(env("UGC_URL", ""), "/"),
		UGCPollInterval:   duration("UGC_POLL_INTERVAL", 5*time.Minute),
		DefaultPriceCents: int32(intEnv("DEFAULT_PRICE_CENTS", 300)),
		Currency:          env("CURRENCY", "USD"),
		AdminEmails:       list("ADMIN_EMAILS"),
		LLMBaseURL:        strings.TrimRight(env("LLM_BASE_URL", ""), "/"),
		LLMModel:          env("LLM_MODEL", "shop"),
		LLMAPIKey:         env("LLM_API_KEY", "not-needed"),
		EmbeddingsBaseURL: strings.TrimRight(env("EMBEDDINGS_BASE_URL", ""), "/"),
		EmbeddingsModel:   env("EMBEDDINGS_MODEL", "shop-embed"),
		ResendAPIKey:      env("RESEND_API_KEY", ""),
		MailFrom:          env("MAIL_FROM", "Ol1n Market <market@ol1n.com>"),
		CookieSecure:      env("COOKIE_SECURE", "false") == "true",
	}
}

func (c Config) IsAdmin(email string) bool {
	for _, a := range c.AdminEmails {
		if strings.EqualFold(a, email) {
			return true
		}
	}
	return false
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func intEnv(key string, def int) int {
	if v, err := strconv.Atoi(os.Getenv(key)); err == nil {
		return v
	}
	return def
}

func duration(key string, def time.Duration) time.Duration {
	if v, err := time.ParseDuration(os.Getenv(key)); err == nil {
		return v
	}
	return def
}

func list(key string) []string {
	var out []string
	for _, s := range strings.Split(os.Getenv(key), ",") {
		if s = strings.TrimSpace(s); s != "" {
			out = append(out, s)
		}
	}
	return out
}
