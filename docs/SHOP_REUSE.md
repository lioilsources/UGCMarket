# SHOP_REUSE.md — audit `shop.ol1n.com` (ChatShop) pro RobloxMarket

> Datum: 2026-10-25 · Zdroj: repo `lioilsources/Shopify`, read-only audit
> Plán: `Prompts/00-ROBLOXMARKET_PLAN.md` (sekce 2, 4, 10)

## TL;DR

- Shop **není** Go + Qdrant + sqlc. Je to **SvelteKit 2 / Svelte 5 PWA** (`apps/web`) + **Python FastAPI „brain“**
  (`services/brain`, SQLAlchemy 2 + alembic) + **Postgres 17 s pgvector** (embeddingy přímo v DB, žádný Qdrant).
- Je to **scaffold + 1 iterace** (5 commitů, 2026-07-22 a 2026-08-28; ~1 300 řádků vlastního kódu). README to říká
  samo: *„Není tam platební brána, objednávkový tok ani reálný katalog — je to funkční jádro, ne e-shop na klíč.“*
- **Chybí:** checkout, objednávky, platby, auth/login, e-maily, admin auth, detail produktu, stahování. Košík je jen
  `localStorage`. **Nasazený není** — `infra/nas/Caddyfile` má `shop.example.com`, iterace 3 „deploy na NAS“ je otevřená.
- Skutečně cenné k převzetí: **SSE RAG chat**, **dynamické filtry z JSONB (filter_schemas + merge)**, **JSONB→SQL
  překladač filtrů**, UI layout chat 20 % / galerie 80 %, napojení na LiteLLM gateway na Sparku (`shop`, `shop-embed`).
- Obchodní jádro marketu (orders, payments, licenses, download tokens, magic link, e-mail) se **píše nově** v každé variantě.

---

## 1. Realita vs. plán (sekce 2 plánu)

| Předpoklad plánu | Realita v repu | Soubory |
|---|---|---|
| Go backend | **Python 3.12 FastAPI** (`brain`) + SvelteKit server endpointy jako BFF (čtou DB přímo přes `postgres.js`). Žádný Go kód. | `services/brain/brain/main.py`, `apps/web/src/routes/api/*/+server.ts` |
| Postgres JSONB/GIN | **Ano.** `products.attributes JSONB` + GIN index `ix_products_attributes`. | `services/brain/alembic/versions/0001_initial.py` |
| Qdrant RAG | **Ne.** pgvector, tabulka `product_embeddings vector(1024)` + HNSW (cosine). Qdrant nikde. | `0001_initial.py`, `brain/rag.py` |
| SSE chat | **Ano.** FastAPI `StreamingResponse`, eventy `products` → `token`* → `done`; SvelteKit proxy. | `brain/chat.py`, `apps/web/src/routes/api/chat/+server.ts`, `lib/components/Chat.svelte` |
| Go ingest state machine | **Ne state machine.** Jeden synchronní POST `/ingest`: lookup (stub) → LLM extrakce → merge filtrů → embedding → insert. Žádné stavy, fronta ani retry. | `brain/ingest.py`, `brain/filters.py` |
| sqlc + golang-migrate | **alembic** (1 migrace) + SQLAlchemy ORM v Pythonu; v TS raw SQL přes `postgres` tagged templates. | `services/brain/alembic/`, `apps/web/src/lib/server/db.ts` |
| Katalog | Ano (galerie + kategorie + filtry). Bez detailu produktu, bez slugů, bez stránkování. | `routes/+page.svelte`, `api/products/+server.ts` |
| Košík | Jen klientský (Svelte 5 `$state` + `localStorage`), bez stránky košíku. | `apps/web/src/lib/cart.svelte.ts` |
| Checkout / platby / objednávky | **Neexistuje.** | — |
| Auth / login | **Neexistuje** (ani admin; `/admin/ingest` je veřejný). | — |
| E-maily | **Neexistuje** (žádný Resend/SMTP). | — |
| Caddy + Cloudflare Tunnel + Compose na JODA | Compose + Caddy (Let's Encrypt, porty 80/443) **připravené, nenasazené**. Cloudflare Tunnel **ne**. | `infra/nas/docker-compose.prod.yml`, `infra/nas/Caddyfile` |
| Frontend „Flutter web?“ (otázka 9.5) | **SvelteKit PWA** (adapter-node, Tailwind 4, service worker, manifest). | `apps/web/` |

---

## 2. Inventář komponent

Legenda zralosti: **funguje** (běží lokálně v mock režimu dle `Prompts/01-PLAN-next-iterations.md`), **test** (pytest), **stub**.

| # | Komponenta | Co dělá | Soubory | Zralost | Klasifikace |
|---|---|---|---|---|---|
| 1 | SSE chat (backend) | Embedding dotazu → top-5 pgvector → prompt s kontextem → stream z LLM; mock odpověď bez LLM | `brain/chat.py`, `brain/rag.py`, `brain/llm.py` | funguje (mock i Spark) | **parametrizovat**: system prompt (CZ→EN, „asistent 3D asset marketu“), karta produktu (`price_czk`→`price_cents`+`currency`, `thumb_url`, `tier`), filtrovat jen `status='listed'` |
| 2 | Embeddings klient | OpenAI-compatible `/embeddings` (LiteLLM `shop-embed` = bge-m3, 1024 dim) nebo deterministický hash-mock | `brain/embeddings.py` | funguje, **test** | **1:1** (konfig `EMBEDDINGS_MODEL`) |
| 3 | LLM klient | `AsyncOpenAI` stream/complete přes LiteLLM gateway | `brain/llm.py` | funguje | **1:1** |
| 4 | Filter schema merge | Z atributů odvodí druh filtru (`bool/range/set/text`), rozšiřuje výčty a meze | `brain/filters.py`, `tests/test_filters.py` | funguje, **test** | **parametrizovat**: `PRICE_ATTR=price_czk`→`price_cents`/USD, případně zafixovat schéma pro vertical `3d-assets` (viz §3) |
| 5 | Ingest | Synchronní pipeline, `_lookup()` je **stub** (vrací `raw_text` nebo název) | `brain/ingest.py` | stub lookup, zbytek funguje | **psát nově** (pull z ugc-backend, stavy) — znovu použít jen `_extract`/merge/embed (§5) |
| 6 | Seed / re-embed | Demo data (trička, boty, HDD); `--reembed` přepočte vektory | `brain/seed.py` | funguje | **parametrizovat** (demo assety), `reembed()` 1:1 |
| 7 | Migrace + ORM | 4 tabulky, alembic | `alembic/versions/0001_initial.py`, `brain/models.py` | funguje | **parametrizovat** + rozšířit o tabulky obchodu |
| 8 | BFF products | JSONB filtry → SQL (`?|`, `->>`, ILIKE, numeric range) | `apps/web/src/routes/api/products/+server.ts` | funguje, netestováno | **parametrizovat** (cena, `status`, stránkování) — logika 1:1 |
| 9 | BFF categories | kategorie + `json_agg` filter schémat | `routes/api/categories/+server.ts` | funguje | **1:1** |
| 10 | BFF chat/ingest proxy | SSE/JSON proxy na brain | `routes/api/chat/+server.ts`, `routes/api/ingest/+server.ts` | funguje | **1:1** (chat); ingest proxy nahradit |
| 11 | UI: Chat | SSE parser přes `fetch` + `TextDecoderStream`, karty doporučení | `lib/components/Chat.svelte` | funguje | **parametrizovat** (texty, karta → odkaz na detail) |
| 12 | UI: Filters | Dynamický sidebar podle `kind` | `lib/components/Filters.svelte` | funguje | **1:1** |
| 13 | UI: ProductCard / galerie | Karta s obrázkem, cenou, „Do košíku“; layout chat/galerie | `lib/components/ProductCard.svelte`, `routes/+page.svelte`, `+layout.svelte` | funguje | **parametrizovat** (badge free/paid, formáty, hover turntable) |
| 14 | Košík | `localStorage`, qty | `lib/cart.svelte.ts` | funguje (jen klient) | **parametrizovat**: digitální zboží ⇒ bez `qty`, dedupe; serverový košík se nepíše, jde rovnou do `orders` |
| 15 | Admin ingest form | Ruční vložení položky | `routes/admin/ingest/+page.svelte` | funguje, bez auth | **psát nově** (admin listing: cena/tier/licence/`listed`) |
| 16 | PWA | manifest + service worker (cache-first build assetů) | `src/service-worker.ts`, `static/manifest.webmanifest` | funguje | **1:1** (pozor: necachovat GLB/download) |
| 17 | Deploy | compose prod (caddy/web/brain/db), Dockerfiles | `infra/nas/*`, `apps/web/Dockerfile`, `services/brain/Dockerfile` | nenasazeno | **parametrizovat** (Cloudflare Tunnel místo LE, viz §7) |
| 18 | Spark LLM compose | záložní vLLM+TEI | `infra/spark/docker-compose.llm.yml` | zastaralé | **nepoužít** — Spark jen přes `sparkctl`/AiStack LiteLLM |
| — | Detail produktu + 3D viewer | — | — | chybí | **psát nově** |
| — | Customers / magic link / session | — | — | chybí | **psát nově** |
| — | Orders, `PaymentProvider` (Free/Stripe/BTCPay), `payment_events` | — | — | chybí | **psát nově** |
| — | Licence, `download_tokens`, signed download, „Moje assety“ | — | — | chybí | **psát nově** |
| — | E-mail (magic link, potvrzení) | — | — | chybí | **psát nově** |

Odhad: ze shopu se reálně přenese ~600–800 řádků (chat, filtry, galerie, BFF). Obchodní část marketu bude 2–3× větší.

---

## 3. Datový model shopu a mapování na plán (sekce 4)

Shop (`0001_initial.py`):

```
categories         id, slug UNIQUE, name
filter_schemas     id, category_id→categories, attr_key, label, kind(text|range|set|bool), config JSONB
                   UNIQUE(category_id, attr_key)
products           id, category_id, name, description, price_czk NUMERIC(12,2), image_url,
                   attributes JSONB  (GIN ix_products_attributes), created_at
product_embeddings product_id PK→products ON DELETE CASCADE, embedding vector(1024)  (HNSW cosine)
```

Atributy jsou volné JSONB (`{"size":["S","M"],"color":"černá","capacity_tb":4,"nas_ready":true}`); filtry se z nich
generují automaticky (`filters.merge_attributes`). Pozn.: filtr `set` přes `attributes->>key = ANY(..) OR ?|` GIN
částečně využije; `range` přes `(attributes->>k)::numeric` GIN nevyužije (pro tisíce položek nevadí).

| Plán (sekce 4) | Shop | Akce |
|---|---|---|
| `products.id, slug, title, description` | `id, name, description` (bez slugu) | přidat `slug`, `name`→`title` |
| `price_cents, currency` | `price_czk NUMERIC` | změnit na `price_cents INT` + `currency` (USD) |
| `tier`, `status`, `asset_ref`, `updated_at` | chybí | přidat (`status` draft/listed/archived, `asset_ref` = ugc job id) |
| `attributes JSONB` + GIN | **shodné** | 1:1; klíče `style, polycount, formats[], roblox_category, tags[], source_pipeline` |
| — | `categories` + `filter_schemas` | **ponechat** — v plánu chybí, ale je to jádro dynamických filtrů. Kategorie = `roblox_category` (Prop, Weapon, Accessory, Character, …) |
| — | `product_embeddings vector(1024)` | ponechat (nahrazuje Qdrant); volitelně `image_embedding` (§4) |
| `product_files`, `product_media` | `image_url` (1 obrázek) | **nové** tabulky |
| `customers, orders, order_items, licenses, download_tokens, payment_events` | chybí | **nové** |

Doporučení pro vertical `3d-assets`: **nenechat filtry „růst“ z LLM**, ale seedovat pevné `filter_schemas`
(`style` set, `polycount` range, `formats` set, `tier` set, `price` range, `tags` set) — data z ugc-backend jsou
strukturovaná, takže auto-merge zůstane jen pro nové tagy.

---

## 4. Chat / RAG

**Jak funguje dnes** (`brain/chat.py`, `brain/rag.py`):
1. poslední user zpráva → `embed()` → `ORDER BY embedding <=> query LIMIT 5` (pgvector cosine),
2. SSE event `products` (karty) jde **hned na začátku**, pak systémový prompt + `product_context()` (id, název, cena,
   kategorie, atributy, 200 znaků popisu) + celá historie → `chat.completions(stream=True)`,
3. eventy `token` a `done`. Model nevrací vybraná id — karty = top-5 z vyhledávání (README tvrdí jinak; realita je tato).

**LLM endpoint:** AiStack LiteLLM gateway na Sparku `http://SPARK_HOST:8080/v1`,
model `shop` = `nvidia/Qwen3-32B-FP4` přes vLLM (:8015), embeddingy `shop-embed` = `BAAI/bge-m3` (vLLM pooling, :8016,
1024 dim). Modul žije v AiStacku (`deploy/docker-compose.shopify.yaml`), viz `Prompts/01-PLAN-next-iterations.md`.
Pozn.: `embeddings.py` posílá `encoding_format: "float"` (fix pro LiteLLM, commit `65035d5`).

**Mock režim:** prázdné `LLM_BASE_URL` → šablonová odpověď; prázdné `EMBEDDINGS_BASE_URL` → hash-vektor slov
(`mock_embed`). Po přepnutí na reálný model `python -m brain.seed --reembed`.

**Co je potřeba pro 3D katalog:**
- Text k embeddingu: `title + description + style + roblox_category + tags + prompt z UGC Factory` (prompt je nejlepší
  popis assetu, ugc-backend ho drží v `jobs.prompt`). Popisy EN (cílovka globální), bge-m3 je multilingvální.
- Filtrovat RAG na `status='listed'`; přidat **hybrid** (pgvector + `filters` z chatu, např. LLM vytáhne `{"style":"cartoon"}`
  a aplikuje JSONB filtr) — v shopu není, plán iterace 4 to zmiňuje.
- Chat má umět i „změnit galerii“ (pitch to slibuje, kód ne) — event `filters` v SSE, Chat.svelte ho propíše do galerie.
- **Image embeddings: v shopu nejsou.** Pro „najdi podobný“ / dotaz obrázkem: CLIP/SigLIP nad turntable náhledy
  (`vector(768)` ve sloupci `image_embedding`), služba na Sparku přes `sparkctl` (ne trvale). Do MVP **ne**, text stačí.
- Spark je sdílený: chat za běhu volá LiteLLM `shop`, který nemusí být nahoře → mock fallback je nutnost, ne luxus.
  Dávkový re-embed katalogu pouštět přes `sparkctl exec <projekt> shop-embed …`.

---

## 5. Ingest pipeline

**Dnes** (`brain/ingest.py`): `POST /ingest {category_slug, name, price_czk, image_url, url?, raw_text?, attributes}` →
`_lookup` (stub, vrací `raw_text|name`) → `_extract` (LLM vrátí JSON `{attributes, description}`, při chybě `{}`) →
`Product` insert → `merge_attributes` → embedding → commit. **Žádné stavy**, žádná idempotence, žádný update existující
položky (každý POST = nový produkt).

**Adaptace na ugc-backend** (kontrakt viz `docs/UGC_API.md`; ugc-api je Go/SQLite, `GET /items`, `GET /jobs/{id}`,
`/jobs/{id}/preview|glb`, `GET /packed/{id}/download`, SSE `GET /events`; přes `ugc.ol1n.com` je za CF Access, z JODA LAN
`http://NAS_HOST:8095` / `http://ugc-api:8095` bez tunelu):

```
ugc packed item ──pull──▶ draft ──(media+files staženy, sha256)──▶ ready ──(admin: cena/tier/licence)──▶ listed ──▶ archived
                        └─ failed (retry)
```
- Stavy = `products.status` (+ `ingest_state` / `ingest_error` pro technické kroky). Upsert podle `asset_ref` (= ugc job id).
- Spouštění: cron/poll `GET /items` (vrací celý katalog bez filtru — filtrovat podle `state` na straně marketu,
  případně do ugc-api doplnit `?state=`) + volitelně poslech SSE `/events`; ruční tlačítko v adminu.
- `_lookup` nahradit čtením metadat z ugc (`name`, `description_en`, `category`, `tags`, `style`, `prompt`, `collection`).
- `_extract` (LLM) použít jen volitelně na doplnění EN popisu / tagů z promptu; atributy jsou strukturované.
- `polycount` počítá Blender worker (`blender_scripts/convert.py` — `loop_triangles`, výsledek ve `jobs.report`;
  ověřit formát v `docs/UGC_API.md`), `formats` = co je v `packed/{id}`; ne LLM.
- Soubory: GLB/FBX/zip z `packed/{id}` zkopírovat do vlastního storage marketu (`product_files`, sha256) — market
  nesmí servírovat soubory přímo z ugc-api (je za Access a není to jeho role). Viewer GLB a turntable → `product_media`.
- `merge_attributes` + `embed` + `reembed` 1:1.

---

## 6. Frontend

**Stack:** SvelteKit 2.16, Svelte 5 (runes `$state/$derived/$effect/$props`), TypeScript, Tailwind 4 (`@tailwindcss/vite`),
adapter-node (port 3000), `postgres` 3.4 pro BFF SQL, PWA. Jediná stránka `/` + `/admin/ingest`; UI česky, ceny v Kč,
emoji v hlavičce. Žádné testy, žádný i18n, žádný routing na detail.

**Integrace `<model-viewer>`** se dotkne:
- nová route `routes/p/[slug]/+page.svelte` + `+page.server.ts` (detail: viewer, turntable, formáty, polycount, licence, cena, CTA),
- `ProductCard.svelte` + `Chat.svelte` — odkaz na detail místo přímého „Do košíku“; hover turntable WebP,
- `app.html` nebo komponenta `ModelViewer.svelte` načítající `@google/model-viewer` (npm, dynamic import v `onMount` —
  SSR ho neumí) — UGC Factory už `model-viewer.min.js` používá (`UGCFactory/nas/viewer/`, `characters_api.go:411`),
- `types.ts` (`Product` + `media`, `files`, `tier`, `price_cents`), `api/products` (vrací `thumb_url`, `viewer_glb_url`),
- `service-worker.ts` — vyloučit `*.glb` a `/download/*` z cache,
- web GLB servírovat z marketu (`/media/{id}.glb`, s `Cache-Control`), ne z ugc-api.

Nové stránky navíc: košík/checkout, login (magic link), „Moje assety“, admin listing.

---

## 7. Deployment

**Shop:** `infra/nas/docker-compose.prod.yml` (projekt `chatshop-prod`): `caddy:2` (80/443, Let's Encrypt) → `web:3000`;
`brain:8000` (při startu `alembic upgrade head`), `db` = `pgvector/pgvector:pg17` (volume `pgdata`, init
`infra/init/01-extensions.sql` = `CREATE EXTENSION vector`). Interní docker síť, ven jen Caddy. Caddyfile má
placeholder `shop.example.com`. **Nikde nasazeno** (iterace 3 otevřená, žádný cloudflared, žádný záznam pro shop.ol1n.com
v repu). Dev: `infra/docker-compose.yml` (Postgres na 5432).

**Vzor, který na JODA skutečně běží** (UGC Factory, `UGCFactory/nas/ugc-stack.yaml` + `cloudflared/config.yml`):
explicitní `name:` projektu, kontejner `cloudflared` s dedikovaným tunelem v compose, ingress `hostname → http://svc:port`,
LAN port publikovaný (`8095:8095`), CF Access pro interní věci, data přes `${DATA_DIR}` (snap docker vidí jen `/home`, `/media`).

**Pro market:** převzít shop compose, ale **místo Caddy+LE použít Cloudflare Tunnel** (vlastní tunel `market-nas`,
`market.ol1n.com → http://market-web:3000`, **bez** Access pro veřejný web; admin cesty za Access nebo vlastní auth),
LAN port např. `8096`, Postgres bez publikovaného portu, `DATA_DIR` pro soubory assetů (ZFS dataset projektu), zálohy `pgdata` + assetů.
Caddy je zbytečný, pokud jde vše přes tunel (Cloudflare dělá TLS). Brain/API volá ugc-api po docker síti nebo `http://NAS_HOST:8095`.

---

## 8. Doporučení stacku („realita vyhrává“)

| | (a) Shop stack as-is: SvelteKit + FastAPI + pgvector | (b) Go backend dle plánu + SvelteKit frontend | (c) Hybrid: SvelteKit (BFF dělá obchod) + FastAPI jen pro AI |
|---|---|---|---|
| Co se převezme | vše z §2 téměř 1:1 | frontend + SQL logiku filtrů; chat/RAG/filters přepsat do Go | frontend + BFF 1:1, brain 1:1 (chat, embed, merge) |
| Obchod (orders, Stripe, licence, tokens) | psát v Pythonu (Stripe SDK ok, BTCPay přes HTTP) | psát v Go (`stripe-go`, vzory z ugc-api: stdlib mux, SSE `sse.go`) | psát v SvelteKit server routes (TS, `stripe` npm, form actions, cookies/session nativně) |
| Plusy | nejrychlejší start, jeden ORM/migrace (alembic), AI kód hotový | odpovídá plánu a ugc-api; jeden statický binár; silné typy | každý jazyk dělá, co umí; žádný přepis hotového; auth/session/checkout přirozeně ve webové vrstvě |
| Mínusy | obchod ve třech vrstvách (TS BFF + Python) ⇒ dvojí DB přístup jako dnes | přepis ~500 řádků funkčního Pythonu bez přidané hodnoty; plán počítal s reuse, ne přepisem; ugc-api navíc **není** sqlc/Postgres (SQLite, ruční SQL), takže „Go pattern se sqlc“ neexistuje ani tam | dvě runtime (Node + Python), migrace musí vlastnit jedna strana |
| Odhad MVP (Fáze 1+2) | ~6–8 sessions | ~9–12 sessions | ~6–8 sessions |

**Doporučení: (c) hybrid, prakticky = (a) s jasným rozdělením rolí.**
- **SvelteKit (apps/web)** vlastní veřejný web, detail s `<model-viewer>`, košík, checkout, magic-link session, „Moje assety“,
  webhooky (`/webhooks/stripe`, `/webhooks/btcpay`) a signed download (`/download/[token]` streamuje soubor z `DATA_DIR`).
  Abstrakce `PaymentProvider` (Free/Stripe/BTCPay) v `src/lib/server/payments/` od prvního commitu.
- **FastAPI brain** zůstává AI službou: chat (SSE), embed/re-embed, ingest z ugc-backend (poll + stavy), merge filtrů.
- **Migrace vlastní alembic** (jeden zdroj pravdy o schématu; TS jen dotazuje). Postgres + pgvector, **bez Qdrantu**.
- **Deploy:** compose na JODA + vlastní Cloudflare Tunnel podle `ugc-stack.yaml`; LLM/embeddingy přes LiteLLM na Sparku
  s mock fallbackem; dávky přes `sparkctl`.
- Plán (sekce 0 bod 2, sekce 2, sekce 10 bod 3) aktualizovat: „Go + sqlc + golang-migrate + Qdrant“ → „SvelteKit + FastAPI
  + alembic + pgvector“ a oznámit Olinovi. Go se vyplatí až ve chvíli, kdy by market měl splynout s ugc-api — to plán
  výslovně nechce („market ho volá, nekopíruje logiku“).

Nový repo `RobloxMarket` (ne fork s `VERTICAL`): zkopírovat `apps/web`, `services/brain`, `infra/` ze shopu jako výchozí
commit s odkazem na zdrojový hash `608415e`; shop zůstane netknutý (nemá uživatele ani nasazení, divergence nevadí).
