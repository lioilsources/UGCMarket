// Package shop: orders, payments, licenses and downloads.
package shop

import (
	"bytes"
	"context"
	"crypto/rand"
	"embed"
	"encoding/base32"
	"errors"
	"fmt"
	"log"
	"strings"
	"text/template"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/lioilsources/RobloxMarket/api/internal/db"
	"github.com/lioilsources/RobloxMarket/api/internal/mail"
	"github.com/lioilsources/RobloxMarket/api/internal/storage"
)

// UserError is an error that can be shown to the buyer (HTTP 400).
type UserError struct{ msg string }

func (e UserError) Error() string { return e.msg }

func errBadRequest(msg string) error { return UserError{msg} }

//go:embed templates/*.txt
var templateFS embed.FS

var templates = template.Must(template.ParseFS(templateFS, "templates/*.txt"))

const (
	downloadTTL     = 24 * time.Hour
	downloadMaxUses = 5
)

type Shop struct {
	pool      *pgxpool.Pool
	q         *db.Queries
	fs        *storage.FS
	mail      mail.Mailer
	publicURL string
	providers map[string]PaymentProvider
}

func New(pool *pgxpool.Pool, fs *storage.FS, m mail.Mailer, publicURL string, providers ...PaymentProvider) *Shop {
	s := &Shop{pool: pool, q: db.New(pool), fs: fs, mail: m, publicURL: publicURL, providers: map[string]PaymentProvider{}}
	for _, p := range providers {
		s.providers[p.Name()] = p
	}
	return s
}

// PaidProviders = payment paths available for paid orders (for the UI).
func (s *Shop) PaidProviders() []string {
	out := []string{}
	for name := range s.providers {
		if name != "free" {
			out = append(out, name)
		}
	}
	return out
}

func (s *Shop) Provider(name string) (PaymentProvider, bool) {
	p, ok := s.providers[name]
	return p, ok
}

type CheckoutRequest struct {
	ProductIDs []int64 `json:"product_ids"`
	Provider   string  `json:"provider"`
	TaxCountry string  `json:"tax_country"`
}

type CheckoutResult struct {
	OrderID     int64  `json:"order_id"`
	Status      string `json:"status"`
	RedirectURL string `json:"redirect_url,omitempty"`
}

// Checkout creates an order and starts the payment. Already owned products are
// skipped; a cart with only free products goes through the Free provider.
func (s *Shop) Checkout(ctx context.Context, cust db.Customer, req CheckoutRequest) (CheckoutResult, error) {
	if len(req.ProductIDs) == 0 || len(req.ProductIDs) > 50 {
		return CheckoutResult{}, errBadRequest("your cart is empty")
	}
	products, err := s.q.ListListedProductsByIDs(ctx, req.ProductIDs)
	if err != nil {
		return CheckoutResult{}, err
	}
	owned, err := s.q.CustomerOwnedProductIDs(ctx, db.CustomerOwnedProductIDsParams{CustomerID: cust.ID, Ids: req.ProductIDs})
	if err != nil {
		return CheckoutResult{}, err
	}
	isOwned := map[int64]bool{}
	for _, id := range owned {
		isOwned[id] = true
	}
	var items []LineItem
	var total int32
	currency := ""
	for _, p := range products {
		if isOwned[p.ID] {
			continue
		}
		price := p.PriceCents
		if p.Tier == "free" {
			price = 0
		}
		if currency != "" && currency != p.Currency {
			return CheckoutResult{}, errBadRequest("the cart mixes currencies")
		}
		currency = p.Currency
		total += price
		items = append(items, LineItem{Product: p, Price: price, License: "personal"})
	}
	if len(items) == 0 {
		return CheckoutResult{}, errBadRequest("you already own these assets or they are no longer available")
	}

	providerName := req.Provider
	if total == 0 {
		providerName = "free"
	}
	provider, ok := s.providers[providerName]
	if !ok {
		return CheckoutResult{}, errBadRequest("payments are not live yet — only free assets can be claimed for now")
	}

	var taxCountry *string
	if c := strings.ToUpper(strings.TrimSpace(req.TaxCountry)); len(c) == 2 {
		taxCountry = &c
	}
	order, err := s.createOrder(ctx, cust.ID, total, currency, provider.Name(), taxCountry, items)
	if err != nil {
		return CheckoutResult{}, err
	}
	start, err := provider.Start(ctx, order, items)
	if err != nil {
		_ = s.q.MarkOrderFailed(ctx, order.ID)
		return CheckoutResult{}, err
	}
	if start.ProviderRef != "" {
		if err := s.q.SetOrderProviderRef(ctx, db.SetOrderProviderRefParams{ID: order.ID, ProviderRef: &start.ProviderRef}); err != nil {
			return CheckoutResult{}, err
		}
	}
	if start.Paid {
		if err := s.Fulfill(ctx, order.ID); err != nil {
			return CheckoutResult{}, err
		}
		return CheckoutResult{OrderID: order.ID, Status: "paid"}, nil
	}
	return CheckoutResult{OrderID: order.ID, Status: "pending", RedirectURL: start.RedirectURL}, nil
}

func (s *Shop) createOrder(ctx context.Context, customerID int64, total int32, currency, provider string, taxCountry *string, items []LineItem) (db.Order, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return db.Order{}, err
	}
	defer tx.Rollback(ctx)
	q := s.q.WithTx(tx)
	order, err := q.CreateOrder(ctx, db.CreateOrderParams{CustomerID: customerID, TotalCents: total, Currency: currency, Provider: provider, TaxCountry: taxCountry})
	if err != nil {
		return db.Order{}, err
	}
	for _, it := range items {
		if err := q.AddOrderItem(ctx, db.AddOrderItemParams{OrderID: order.ID, ProductID: it.Product.ID, PriceCents: it.Price, LicenseType: it.License}); err != nil {
			return db.Order{}, err
		}
	}
	return order, tx.Commit(ctx)
}

// Fulfill marks the order as paid and issues licenses. Idempotent:
// a repeated call (duplicate webhook) changes nothing.
func (s *Shop) Fulfill(ctx context.Context, orderID int64) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	q := s.q.WithTx(tx)
	order, err := q.MarkOrderPaid(ctx, orderID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil // already paid (or failed/refunded) — do nothing
	}
	if err != nil {
		return err
	}
	items, err := q.ListOrderItems(ctx, orderID)
	if err != nil {
		return err
	}
	for _, it := range items {
		if _, err := q.CreateLicense(ctx, db.CreateLicenseParams{OrderID: orderID, ProductID: it.ProductID,
			CustomerID: order.CustomerID, LicenseType: it.LicenseType, LicenseKey: newLicenseKey()}); err != nil {
			return err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	go s.sendReceipt(order, len(items))
	return nil
}

// HandlePaymentEvent: common idempotent handling of webhooks from all providers.
func (s *Shop) HandlePaymentEvent(ctx context.Context, provider string, ev PaymentEvent) error {
	_, err := s.q.RecordPaymentEvent(ctx, db.RecordPaymentEventParams{Provider: provider, EventID: ev.ID, Payload: ev.Payload})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil // duplicate delivery
	}
	if err != nil {
		return err
	}
	switch ev.Status {
	case "paid":
		return s.Fulfill(ctx, ev.OrderID)
	case "failed":
		return s.q.MarkOrderFailed(ctx, ev.OrderID)
	}
	return nil
}

func (s *Shop) sendReceipt(order db.Order, n int) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cust, err := s.customerEmail(ctx, order.CustomerID)
	if err != nil {
		log.Printf("receipt order %d: %v", order.ID, err)
		return
	}
	text := fmt.Sprintf("Thanks! Order #%d is complete (%d asset(s)).\n\nDownload your files and licenses any time at:\n%s/me\n", order.ID, n, s.publicURL)
	if err := s.mail.Send(ctx, cust, fmt.Sprintf("Ol1n Market — order #%d", order.ID), text); err != nil {
		log.Printf("receipt order %d: %v", order.ID, err)
	}
}

func (s *Shop) customerEmail(ctx context.Context, id int64) (string, error) {
	var email string
	err := s.pool.QueryRow(ctx, "SELECT email FROM customers WHERE id = $1", id).Scan(&email)
	return email, err
}

func newLicenseKey() string {
	b := make([]byte, 10)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	k := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b) // 16 chars
	return fmt.Sprintf("OM-%s-%s-%s-%s", k[0:4], k[4:8], k[8:12], k[12:16])
}

func render(name string, data any) []byte {
	var buf bytes.Buffer
	if err := templates.ExecuteTemplate(&buf, name, data); err != nil {
		log.Printf("template %s: %v", name, err)
	}
	return buf.Bytes()
}
