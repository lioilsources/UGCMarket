# Ol1n Market 🧊

Obchod s 3D modely z UGC Factory — `market.ol1n.com`. Hotové kusy (stav `packed`
v ugc-backend) se natáhnou do katalogu jako draft, admin je vylistuje (cena / free),
kupující je stáhne jako ZIP (FBX + textura + GLB + LICENSE + README).
Veřejně **„Ol1n Market“** — jméno Roblox se v názvu obchodu použít nesmí (viz
`docs/ROBLOX_COMPLIANCE.md`). Interní/repo jméno RobloxMarket zůstává.

```
prohlížeč ──► market-web (SvelteKit, :3000)  ── /api/* proxy ──►  market-api (Go, :8096)
                                                                   ├─ Postgres 17 + pgvector (JSONB/GIN, embeddingy)
                                                                   ├─ soubory /data (kopie z UGC)
                                                                   ├─ ingest ◄── ugc-api :8095 (síť ugcfactory_ugc)
                                                                   ├─ chat SSE ── LiteLLM na SPARKu (nebo mock)
                                                                   └─ platby: PaymentProvider (Free; Stripe/BTCPay = Fáze 2)
```

## Struktura

| Adresář | Co to je |
|---|---|
| `api/` | Go `market-api`: katalog, ingest z UGC, magic-link login, checkout, licence, stahování, chat. sqlc + golang-migrate (migrace v binárce, pouští se při startu) |
| `web/` | SvelteKit PWA převzatá ze shopu (`Dev/Shopify/apps/web`): chat nahoře, galerie s filtry, detail s `<model-viewer>`, košík, „My assets“, admin |
| `deploy/` | compose stack pro JODA + Cloudflare tunnel `market-nas` |
| `docs/` | Fáze 0 rešerše (`SHOP_REUSE`, `UGC_API`, `ROBLOX_COMPLIANCE`) + screenshoty |
| `Prompts/` | plán (`00-ROBLOXMARKET_PLAN.md`, sekce „Aktualizace po Fázi 0“ platí přednostně) |

## Lokální vývoj

Docker na Macu neběží. Vývojová a testovací databáze (`market_dev`, `market_test`) jsou
v Postgresu na JODĚ, port 5438 jen na LAN adrese.

```bash
cp .env.example .env           # adresa JODY a heslo z deploy/.env na JODĚ
make deploy-db                 # jednou: Postgres na JODĚ (nic veřejného)
make api                       # market-api :8196 proti market_dev, ingest z UGC_URL (jen GET)
make web                       # http://localhost:5173
make test                      # go vet + Go testy (e2e maže market_test) + svelte-check
```

## API (market-api)

| Metoda | Cesta | Kdo |
|---|---|---|
| GET | `/api/products?q=&tier=&filters={json}` · `/api/products/{slug}` · `/api/catalog/facets` | veřejné |
| GET | `/api/media/{id}/thumb` · `/api/media/{id}/viewer.glb` | veřejné (draft jen admin) |
| POST | `/api/chat` (SSE: `products`, `token`, `done`) | veřejné |
| POST | `/api/auth/login` · GET `/api/auth/verify` · POST `/api/auth/logout` | veřejné |
| GET/POST | `/api/me`, `/api/me/licenses`, `/api/me/licenses/{id}/download`, `/api/checkout` | přihlášený |
| GET | `/api/download/{token}` (24 h, 5 použití) | držitel tokenu |
| POST | `/api/webhooks/{provider}` (idempotentní přes `payment_events`) | platební brána |
| GET/PATCH/POST | `/api/admin/products[/{id}]`, `/api/admin/ingest`, `/api/admin/orders` | admin |

## Nasazení (JODA)

```
tunnel market-nas ─► Caddy :80 ─► market-web :3000 ─► market-api :8096 ─► market-db
```
Na hostu není publikovaný žádný port kromě Postgresu na LAN adrese (`DB_BIND`).

Jednou ručně (Olin):
1. ZFS dataset `pool/market` a pod ním `pool/market/postgres` (recordsize 16K), mountpoint
   pod domovským adresářem, aby ho viděl snap docker. Do té doby `DATA_DIR` v `deploy/.env`
   na JODĚ ukazuje do `data/` v deploy adresáři projektu.
2. `cloudflared tunnel create market-nas` a `cloudflared tunnel route dns market-nas market.ol1n.com`.
3. Resend: ověřená doména a `RESEND_API_KEY` v `deploy/.env` na JODĚ.

Pak:
```bash
make deploy-nas                # rsync, config tunelu z UUID, compose up, zápis do events.jsonl
make status
```
