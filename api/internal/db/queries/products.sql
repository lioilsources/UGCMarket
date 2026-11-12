-- name: GetProductByID :one
SELECT * FROM products WHERE id = $1;

-- name: GetProductBySlug :one
SELECT * FROM products WHERE slug = $1;

-- name: GetProductByAssetRef :one
SELECT * FROM products WHERE asset_ref = $1;

-- name: SlugTaken :one
SELECT EXISTS (SELECT 1 FROM products WHERE slug = $1);

-- name: InsertProduct :one
INSERT INTO products (slug, title, description, price_cents, currency, tier, asset_ref, attributes, source_updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: UpdateProductSnapshot :exec
UPDATE products
SET attributes = $2, source_updated_at = $3, status = $4, updated_at = now()
WHERE id = $1;

-- name: ListIngestedProducts :many
SELECT id, asset_ref, status, attributes FROM products;

-- name: AdminListProducts :many
SELECT * FROM products ORDER BY created_at DESC, id DESC;

-- name: AdminUpdateProduct :one
UPDATE products
SET title = $2, description = $3, price_cents = $4, tier = $5, status = $6, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: ListListedProductsByIDs :many
SELECT * FROM products WHERE status = 'listed' AND id = ANY(@ids::bigint[]);

-- name: UpsertProductFile :exec
INSERT INTO product_files (product_id, format, storage_key, bytes, sha256)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (product_id, format) DO UPDATE
SET storage_key = EXCLUDED.storage_key, bytes = EXCLUDED.bytes, sha256 = EXCLUDED.sha256;

-- name: ListProductFiles :many
SELECT * FROM product_files WHERE product_id = $1 ORDER BY format;

-- name: UpsertProductMedia :exec
INSERT INTO product_media (product_id, kind, position, storage_key)
VALUES ($1, $2, $3, $4)
ON CONFLICT (product_id, kind, position) DO UPDATE SET storage_key = EXCLUDED.storage_key;

-- name: ListProductMedia :many
SELECT * FROM product_media WHERE product_id = $1 ORDER BY kind, position;

-- name: GetProductMedia :one
SELECT * FROM product_media WHERE product_id = $1 AND kind = $2 AND position = $3;

-- name: GetEmbeddingHash :one
SELECT text_hash FROM product_embeddings WHERE product_id = $1;

-- name: UpsertEmbedding :exec
INSERT INTO product_embeddings (product_id, embedding, text_hash)
VALUES ($1, $2, $3)
ON CONFLICT (product_id) DO UPDATE SET embedding = EXCLUDED.embedding, text_hash = EXCLUDED.text_hash;

-- name: SearchListedProducts :many
SELECT p.*
FROM product_embeddings e
JOIN products p ON p.id = e.product_id
WHERE p.status = 'listed'
ORDER BY e.embedding <=> @query::vector
LIMIT @lim;
