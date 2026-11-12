-- name: CreateOrder :one
INSERT INTO orders (customer_id, total_cents, currency, provider, tax_country)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: AddOrderItem :exec
INSERT INTO order_items (order_id, product_id, price_cents, license_type) VALUES ($1, $2, $3, $4);

-- name: SetOrderProviderRef :exec
UPDATE orders SET provider_ref = $2 WHERE id = $1;

-- name: MarkOrderPaid :one
UPDATE orders SET status = 'paid', paid_at = now()
WHERE id = $1 AND status = 'pending'
RETURNING *;

-- name: MarkOrderFailed :exec
UPDATE orders SET status = 'failed' WHERE id = $1 AND status = 'pending';

-- name: GetOrder :one
SELECT * FROM orders WHERE id = $1;

-- name: ListOrderItems :many
SELECT * FROM order_items WHERE order_id = $1;

-- name: AdminListOrders :many
SELECT o.*, c.email,
       (SELECT count(*) FROM order_items i WHERE i.order_id = o.id)::int AS item_count
FROM orders o
JOIN customers c ON c.id = o.customer_id
ORDER BY o.created_at DESC
LIMIT $1;

-- name: RecordPaymentEvent :one
-- Returns a row only on the first write — a repeated webhook gets pgx.ErrNoRows.
INSERT INTO payment_events (provider, event_id, payload) VALUES ($1, $2, $3)
ON CONFLICT DO NOTHING
RETURNING provider;

-- name: CustomerOwnedProductIDs :many
SELECT product_id FROM licenses WHERE customer_id = $1 AND product_id = ANY(@ids::bigint[]);

-- name: CreateLicense :one
INSERT INTO licenses (order_id, product_id, customer_id, license_type, license_key)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (order_id, product_id) DO UPDATE SET license_key = licenses.license_key
RETURNING *;

-- name: GetLicense :one
SELECT * FROM licenses WHERE id = $1;

-- name: ListCustomerLicenses :many
SELECT l.*, p.title, p.slug, p.tier
FROM licenses l
JOIN products p ON p.id = l.product_id
WHERE l.customer_id = $1
ORDER BY l.issued_at DESC;

-- name: CreateDownloadToken :exec
INSERT INTO download_tokens (token_hash, license_id, expires_at, max_uses) VALUES ($1, $2, $3, $4);

-- name: UseDownloadToken :one
UPDATE download_tokens SET uses = uses + 1
WHERE token_hash = $1 AND uses < max_uses AND expires_at > now()
RETURNING license_id;

-- name: CountDownloads :one
SELECT coalesce(sum(uses), 0)::bigint FROM download_tokens;
