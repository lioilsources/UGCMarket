-- Ol1n Market: katalog 3D assetů z UGC Factory + objednávky, licence, stahování.
-- Enumy jsou TEXT + CHECK (snáz se rozšiřují migrací než Postgres ENUM).

CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE products (
    id                BIGSERIAL PRIMARY KEY,
    slug              TEXT NOT NULL UNIQUE,
    title             TEXT NOT NULL,
    description       TEXT NOT NULL DEFAULT '',
    price_cents       INTEGER NOT NULL DEFAULT 0 CHECK (price_cents >= 0),
    currency          TEXT NOT NULL DEFAULT 'USD',
    tier              TEXT NOT NULL DEFAULT 'paid' CHECK (tier IN ('free', 'paid')),
    -- id jobu v ugc-backend; ugc zůstává zdrojem pravdy, market drží snapshot
    asset_ref         TEXT NOT NULL UNIQUE,
    attributes        JSONB NOT NULL DEFAULT '{}',
    status            TEXT NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'listed', 'archived')),
    -- updated_at jobu/itemu v ugc v okamžiku posledního ingestu (detekce změn)
    source_updated_at TIMESTAMPTZ,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX products_attributes_gin ON products USING GIN (attributes);
CREATE INDEX products_status_idx ON products (status);

CREATE TABLE product_files (
    id          BIGSERIAL PRIMARY KEY,
    product_id  BIGINT NOT NULL REFERENCES products (id) ON DELETE CASCADE,
    format      TEXT NOT NULL CHECK (format IN ('glb', 'fbx', 'obj', 'rbxm', 'png', 'zip')),
    storage_key TEXT NOT NULL,
    bytes       BIGINT NOT NULL,
    sha256      TEXT NOT NULL,
    UNIQUE (product_id, format)
);

CREATE TABLE product_media (
    id          BIGSERIAL PRIMARY KEY,
    product_id  BIGINT NOT NULL REFERENCES products (id) ON DELETE CASCADE,
    kind        TEXT NOT NULL CHECK (kind IN ('thumb', 'turntable', 'viewer_glb')),
    position    INTEGER NOT NULL DEFAULT 0,
    storage_key TEXT NOT NULL,
    UNIQUE (product_id, kind, position)
);

CREATE TABLE product_embeddings (
    product_id BIGINT PRIMARY KEY REFERENCES products (id) ON DELETE CASCADE,
    embedding  vector(1024) NOT NULL,
    -- sha256 textu, ze kterého embedding vznikl — přepočítá se jen při změně
    text_hash  TEXT NOT NULL
);
CREATE INDEX product_embeddings_hnsw ON product_embeddings USING hnsw (embedding vector_cosine_ops);

CREATE TABLE customers (
    id             BIGSERIAL PRIMARY KEY,
    email          TEXT NOT NULL,
    roblox_user_id TEXT,
    wallet_hint    TEXT,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX customers_email_idx ON customers (lower(email));

-- Tokeny (magic link, session, download) se ukládají jen jako sha256 hash.
CREATE TABLE login_tokens (
    token_hash TEXT PRIMARY KEY,
    email      TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    used_at    TIMESTAMPTZ
);

CREATE TABLE sessions (
    token_hash  TEXT PRIMARY KEY,
    customer_id BIGINT NOT NULL REFERENCES customers (id) ON DELETE CASCADE,
    expires_at  TIMESTAMPTZ NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE orders (
    id           BIGSERIAL PRIMARY KEY,
    customer_id  BIGINT NOT NULL REFERENCES customers (id),
    total_cents  INTEGER NOT NULL CHECK (total_cents >= 0),
    currency     TEXT NOT NULL DEFAULT 'USD',
    status       TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'paid', 'failed', 'refunded')),
    provider     TEXT NOT NULL CHECK (provider IN ('stripe', 'btcpay', 'free')),
    provider_ref TEXT,
    -- země kupujícího pro OSS/DPH evidenci; daň se zatím nepočítá
    tax_country  TEXT,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    paid_at      TIMESTAMPTZ
);
CREATE INDEX orders_customer_idx ON orders (customer_id);

CREATE TABLE order_items (
    order_id     BIGINT NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
    product_id   BIGINT NOT NULL REFERENCES products (id),
    price_cents  INTEGER NOT NULL CHECK (price_cents >= 0),
    license_type TEXT NOT NULL DEFAULT 'personal' CHECK (license_type IN ('personal', 'commercial')),
    PRIMARY KEY (order_id, product_id)
);

CREATE TABLE licenses (
    id           BIGSERIAL PRIMARY KEY,
    order_id     BIGINT NOT NULL REFERENCES orders (id),
    product_id   BIGINT NOT NULL REFERENCES products (id),
    customer_id  BIGINT NOT NULL REFERENCES customers (id),
    license_type TEXT NOT NULL CHECK (license_type IN ('personal', 'commercial')),
    license_key  TEXT NOT NULL UNIQUE,
    issued_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (order_id, product_id)
);
CREATE INDEX licenses_customer_idx ON licenses (customer_id);

CREATE TABLE download_tokens (
    token_hash TEXT PRIMARY KEY,
    license_id BIGINT NOT NULL REFERENCES licenses (id) ON DELETE CASCADE,
    expires_at TIMESTAMPTZ NOT NULL,
    max_uses   INTEGER NOT NULL,
    uses       INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Idempotence webhooků: (provider, event_id) se zpracuje jen jednou.
CREATE TABLE payment_events (
    provider     TEXT NOT NULL,
    event_id     TEXT NOT NULL,
    payload      JSONB NOT NULL,
    processed_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (provider, event_id)
);
