-- name: CreateLoginToken :exec
INSERT INTO login_tokens (token_hash, email, expires_at) VALUES ($1, $2, $3);

-- name: ConsumeLoginToken :one
UPDATE login_tokens SET used_at = now()
WHERE token_hash = $1 AND used_at IS NULL AND expires_at > now()
RETURNING email;

-- name: UpsertCustomer :one
INSERT INTO customers (email) VALUES ($1)
ON CONFLICT ((lower(email))) DO UPDATE SET email = customers.email
RETURNING *;

-- name: CreateSession :exec
INSERT INTO sessions (token_hash, customer_id, expires_at) VALUES ($1, $2, $3);

-- name: GetSessionCustomer :one
SELECT c.*
FROM sessions s
JOIN customers c ON c.id = s.customer_id
WHERE s.token_hash = $1 AND s.expires_at > now();

-- name: DeleteSession :exec
DELETE FROM sessions WHERE token_hash = $1;

-- name: PurgeExpiredTokens :exec
WITH a AS (DELETE FROM login_tokens WHERE login_tokens.expires_at < now() - interval '1 day'),
     b AS (DELETE FROM download_tokens WHERE download_tokens.expires_at < now() - interval '1 day')
DELETE FROM sessions WHERE sessions.expires_at < now();
