package shop

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/lioilsources/RobloxMarket/api/internal/db"
)

// PaymentProvider is one payment path. The checkout flow is the same for all:
// order.create → provider.Start → (webhook) paid → fulfill.
// Free is a provider that returns Paid immediately; Stripe and BTCPay will be added in Phase 2
// as a WebhookProvider.
type PaymentProvider interface {
	Name() string // value of orders.provider
	Start(ctx context.Context, order db.Order, items []LineItem) (StartResult, error)
}

// WebhookProvider can verify and parse a webhook (signature, payload).
// Idempotency (payment_events) is handled by the common handler, not the provider.
type WebhookProvider interface {
	PaymentProvider
	ParseWebhook(r *http.Request) (PaymentEvent, error)
}

type LineItem struct {
	Product db.Product
	Price   int32
	License string
}

type StartResult struct {
	Paid        bool   // payment completed immediately (free)
	RedirectURL string // hosted checkout (Stripe Checkout, BTCPay invoice)
	ProviderRef string // session id / invoice id
}

type PaymentEvent struct {
	ID         string // Stripe event.id, BTCPay invoiceId+type
	OrderID    int64
	Status     string // paid | failed | ignore
	TaxCountry string
	Payload    json.RawMessage
}

type Free struct{}

func (Free) Name() string { return "free" }

func (Free) Start(_ context.Context, order db.Order, _ []LineItem) (StartResult, error) {
	if order.TotalCents != 0 {
		return StartResult{}, errBadRequest("the free provider only handles zero-total orders")
	}
	return StartResult{Paid: true}, nil
}
