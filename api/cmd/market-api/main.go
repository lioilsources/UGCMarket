// market-api: backend of Ol1n Market (market.ol1n.com) — catalog of 3D assets
// from UGC Factory, orders, licenses, downloads and the chat assistant.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	pgxvec "github.com/pgvector/pgvector-go/pgx"

	"github.com/lioilsources/RobloxMarket/api/internal/ai"
	"github.com/lioilsources/RobloxMarket/api/internal/auth"
	"github.com/lioilsources/RobloxMarket/api/internal/config"
	"github.com/lioilsources/RobloxMarket/api/internal/db"
	"github.com/lioilsources/RobloxMarket/api/internal/httpapi"
	"github.com/lioilsources/RobloxMarket/api/internal/ingest"
	"github.com/lioilsources/RobloxMarket/api/internal/mail"
	"github.com/lioilsources/RobloxMarket/api/internal/shop"
	"github.com/lioilsources/RobloxMarket/api/internal/storage"
	"github.com/lioilsources/RobloxMarket/api/internal/ugc"
)

func main() {
	cfg := config.Load()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := db.Migrate(cfg.DatabaseURL); err != nil {
		log.Fatal(err)
	}
	pcfg, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	pcfg.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error { return pgxvec.RegisterTypes(ctx, conn) }
	pool, err := pgxpool.NewWithConfig(ctx, pcfg)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	fs, err := storage.NewFS(cfg.DataDir)
	if err != nil {
		log.Fatal(err)
	}
	q := db.New(pool)
	mailer := mail.New(cfg.ResendAPIKey, cfg.MailFrom)
	emb := ai.NewEmbedder(cfg.EmbeddingsBaseURL, cfg.EmbeddingsModel)
	llm := ai.NewLLM(cfg.LLMBaseURL, cfg.LLMModel, cfg.LLMAPIKey)
	authSvc := auth.New(q, mailer, cfg.PublicURL, cfg.CookieSecure)
	// Phase 1: Free only. Stripe / BTCPay will be added here as further PaymentProviders.
	shopSvc := shop.New(pool, fs, mailer, cfg.PublicURL, shop.Free{})

	var ingester *ingest.Ingester
	if cfg.UGCURL != "" {
		ingester = ingest.New(pool, fs, ugc.New(cfg.UGCURL), emb, cfg.DefaultPriceCents, cfg.Currency)
		if cfg.UGCPollInterval > 0 {
			go ingester.Poll(ctx, cfg.UGCPollInterval)
		}
	}
	go purgeLoop(ctx, q)

	srv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           httpapi.New(cfg, pool, fs, authSvc, shopSvc, ingester, emb, llm).Routes(),
		ReadHeaderTimeout: 10 * time.Second,
	}
	go func() {
		<-ctx.Done()
		shutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdown)
	}()
	log.Printf("market-api on %s (ugc=%q llm=%v embeddings=%v mail=%v)", cfg.Addr, cfg.UGCURL,
		cfg.LLMBaseURL != "", cfg.EmbeddingsBaseURL != "", cfg.ResendAPIKey != "")
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

func purgeLoop(ctx context.Context, q *db.Queries) {
	for {
		if err := q.PurgeExpiredTokens(ctx); err != nil && ctx.Err() == nil {
			log.Printf("purge: %v", err)
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Hour):
		}
	}
}
