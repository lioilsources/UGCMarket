# ROBLOXMARKET_PLAN.md — handoff pro Opuse

> **Projekt:** Roblox Market — prodej 3D modelů z UGC Factory (za peníze, crypto, nebo zdarma)
> **Doména:** `market.ol1n.com` (návrh)
> **Stav:** koncept / MVP ke zkoušce — cíl je ověřit, jestli lidi koupí AI-generované Roblox assety přímo od Olina
> **Autor plánu:** Claude (chat) pro Olina, 2026-10-25
> **Určeno pro:** Claude Code (Opus) — Director session + sub-sessions přes agent hub

---

## ⚠️ Aktualizace po Fázi 0 (2026-10-25) — realita vs. plán

Podklady: `docs/SHOP_REUSE.md`, `docs/UGC_API.md`, `docs/ROBLOX_COMPLIANCE.md`. Kde se níže uvedené liší od zbytku plánu, **platí tahle sekce**.

**Rozpory s plánem:**
1. **Shop stack není Go + Qdrant.** `shop.ol1n.com` (repo `Dev/Shopify`, „ChatShop“) = SvelteKit 2 PWA + Python FastAPI „brain“ + alembic + Postgres 17 s pgvector (bge-m3, HNSW) a JSONB/GIN. Žádné sqlc/golang-migrate. Je to scaffold + 1 iterace (5 commitů), **nenasazený**, bez košíku na serveru, objednávek, plateb, loginu a e-mailů. Přenositelné: SSE RAG chat, auto-filtry z atributů, JSONB filtr → SQL, chat-nad-galerií UI.
2. **ugc-backend nemá stavy `reviewed`/`listed`.** Joby: `new → approved → converting → converted (PASS/WARN/FAIL) → packed`. „Připraveno k listingu“ = `status=packed` a `verdict != FAIL`. Market drží vlastní `draft → listed`.
3. **Formáty:** existuje GLB (high-poly, 2k textury), `preview.glb` (512px, pro viewer), FBX (≤4k tris) + PNG textura, zip. **OBJ a RBXM se nevyrábějí, turntable rendery nejsou automatické** (`render_views.py` jen ručně). Prodávaný GLB a FBX se liší geometrií → v MVP prodávat FBX + preview z GLB, OBJ/GLB z konvertovaného modelu doplnit v ugc-backend.
4. **Značka:** podle Roblox brand guidelines se jméno „Roblox“ nesmí vázat na prodej produktů bez licence → **veřejný název „RobloxMarket“ / Roblox logo nelze**. Repo se může jmenovat interně jakkoli, doména `market.ol1n.com` je OK; „Roblox“ jen popisně („compatible with Roblox Studio“) + disclaimer „not affiliated“.
5. **Licence modelů pipeline:** SF3D draft modely = Stability Community License; Illustrious/TRELLIS.2 výstupy ověřit pro komerční prodej (otevřené).
6. Fantasy characters (`/v1/fc/`) mimo MVP (Mixamo animace = licenční problém pro přeprodej).

**Odpovědi Olina na sekci 9 (2026-10-25):** jednotná cena + free pack od začátku · crypto = **BTC/Lightning stačí** → self-hosted BTCPay bez USDT pluginů · login = e-mail magic link · licence = personal + commercial (bez exkluzivity). Otevřené zůstávají: refundy při zamítnutí moderací (návrh v `ROBLOX_COMPLIANCE.md`: oprava nebo refund do 14 dní), výše ceny.

**Rozhodnutí Olina (2026-10-25):**
- **Stack: Go backend podle plánu** — `market-api` v Go (sqlc + golang-migrate + Postgres JSONB/GIN), frontend = SvelteKit ze shopu (jen UI, volá market-api), chat/RAG přepsat do Go. Vektory: pgvector v Postgresu (shop už ho má, jedna DB méně) — Qdrant jen pokud by pgvector nestačil. LLM/embeddingy přes AiStack LiteLLM na Sparku (`shop-embed` = bge-m3), bez LLM mock režim jako ve shopu.
- **Veřejný název: „Ol1n Market“** (`market.ol1n.com`). „Roblox“ jen popisně + disclaimer „not affiliated with Roblox Corporation“.
- Nový samostatný repo, ne monorepo se shopem.

---

## 0. TL;DR pro Opuse

1. Stavíš **market.ol1n.com** = Shopify-like obchod, ale produkty jsou 3D modely (GLB/FBX/RBXM) vyrobené v **UGC Factory** (`ugc-backend` na JODA :8095 + Blender worker + TRELLIS na Sparku).
2. **Znovu použij stack z `shop.ol1n.com`** (Go backend, Postgres JSONB/GIN, Qdrant RAG, SSE chat, ingest state machine). Nezakládej nový stack — forkni/parametrizuj existující, ať se koncept dá otestovat za pár dní.
3. Tři platební cesty: **Stripe** (karta), **crypto** (self-hosted BTCPay Server — sedí do Olinovy infra), **free** (jen účet + stažení).
4. MVP = katalog + detail s 3D náhledem + košík + jedna platba + licence + download link. Nic víc.
5. **Nejdřív ověř Roblox pravidla** (sekce 7) — rozhoduje o tom, co přesně prodáváme (soubory vs. publikaci do Roblox účtu kupujícího).

---

## 1. Cíl a koncept

- Olin má pipeline (UGC Factory), která generuje texturované 3D assety pro Roblox (z obrázku/promptu → GLB → Blender cleanup → Roblox-ready).
- Dnes nemá kanál, kde by si je lidi mohli **koupit přímo od něj**. Roblox Creator Store je jeden kanál (Robux), ale Olin chce **vlastní market** s reálnými penězi / crypto / free tierem — a chce použít stejný stack jako u `shop.ol1n.com`, aby se koncept dal rychle ověřit.
- Hypotéza k otestování: *"Roblox devs koupí hotový, stylově konzistentní 3D asset pack za pár dolarů/USDT, když je kvalita OK a dodání okamžité."*

**Co market NENÍ (zatím):** marketplace třetích stran, aukce, UGC avatar katalog uvnitř Robloxu, subscription. Jen Olinovy assety.

---

## 2. Co už existuje — a co z toho použít

| Komponenta | Kde | Použít jak |
|---|---|---|
| `shop.ol1n.com` | repo shop e-shopu (Opus: najdi na JODA / GitHub `lioilsources`) | **Základ marketu.** Katalog, košík, checkout, chat-asistent, ingest state machine. |
| Postgres JSONB/GIN | shop stack | Produkt = `products` řádek s JSONB `attributes` (poly count, formáty, styl, tagy, Roblox category). GIN index pro filtry. |
| Qdrant RAG + SSE chat | shop stack | "Najdi mi low-poly středověký meč v cartoon stylu" → embedding search přes katalog + obrázky. |
| Go ingest pipeline (state machine) | shop stack | Místo Shopify feedu ingestuje **z `ugc-backend` API** — stav `generated → reviewed → listed`. |
| `ugc-backend` (Go, JODA :8095) | UGC Factory | Zdroj pravdy pro assety: metadata, soubory, náhledy. Market ho volá, nekopíruje logiku. |
| Blender headless worker (JODA) | UGC Factory | Generuje **náhledy** (turntable PNG/WebP, GLB pro web viewer) a **exporty** (FBX, RBXM, OBJ). |
| TRELLIS / ComfyUI (Spark) | UGC Factory | Výroba assetů — mimo scope marketu, jen konzument výstupů. |
| `ugc_studio` (Flutter) | UGC Factory | Admin: tlačítko "Publikovat na market" + nastavení ceny/licence. |
| Caddy + Cloudflare Tunnel + Docker Compose na JODA | infra | Hosting `market.ol1n.com` stejně jako ostatní `*.ol1n.com`. |
| sqlc + golang-migrate | infra pattern | Migrace a dotazy. |

**Opus, první úkol:** projdi repo shopu a napiš do `docs/SHOP_REUSE.md`, co je přenositelné 1:1, co se parametrizuje (konfig "vertical = 3d-assets"), a co se píše nově. Teprve pak kóduj.

---

## 3. Architektura (cílový stav MVP)

```
Kupující (web)  ──►  market-web (frontend ze shop stacku + 3D viewer)
                          │ REST/SSE
                          ▼
                    market-api (Go, JODA, nový port např. :8096)
                     ├─ catalog / cart / orders / licenses / downloads
                     ├─ chat (SSE) ── Qdrant (embeddingy katalogu)
                     ├─ payments: Stripe | BTCPay | free
                     └─ ingest: pull z ugc-backend :8095
                          │
          ┌───────────────┼──────────────────┐
          ▼               ▼                  ▼
   Postgres (JSONB)   MinIO / FS na JODA   ugc-backend :8095
   products, orders,  soubory assetů        (metadata, Blender
   licenses, dl_tokens (signed URL dl)       exporty, náhledy)
          
Platby:  Stripe (webhook) ──┐
         BTCPay Server ─────┼──► market-api /webhooks/* ──► order.paid ──► vydat licenci + download token
         free (bez platby) ─┘
```

Zásady:
- **Soubory assetů nejdou přes Cloudflare Tunnel veřejně napřímo** — stahování přes `market-api` se signed, časově omezenými tokeny (nebo MinIO presigned URL). Free tier taky přes token (kvůli statistice + rate limitu).
- **Idempotentní webhooky** (Stripe `event.id`, BTCPay `invoiceId`) — zapisuj do tabulky `payment_events`.
- `ugc-backend` zůstává jediným místem, kde se assety vyrábějí a verzují; market drží jen `asset_ref` + snapshot metadat v JSONB.

---

## 4. Datový model (návrh, sqlc + migrate)

```sql
products        (id, slug, title, description, price_cents, currency, tier ENUM('free','paid'),
                 asset_ref TEXT, attributes JSONB, status ENUM('draft','listed','archived'),
                 created_at, updated_at)                      -- GIN na attributes
product_files   (id, product_id, format ENUM('glb','fbx','obj','rbxm','zip'), storage_key, bytes, sha256)
product_media   (id, product_id, kind ENUM('thumb','turntable','viewer_glb'), storage_key)
customers       (id, email, roblox_user_id NULL, wallet_hint NULL, created_at)
orders          (id, customer_id, total_cents, currency, status ENUM('pending','paid','failed','refunded'),
                 provider ENUM('stripe','btcpay','free'), provider_ref, created_at, paid_at)
order_items     (order_id, product_id, price_cents)
licenses        (id, order_id, product_id, customer_id, license_type ENUM('personal','commercial'),
                 license_key, issued_at)                      -- PDF/text licence generovaná z šablony
download_tokens (token, license_id, expires_at, max_uses, uses)
payment_events  (provider, event_id PK, payload JSONB, processed_at)
```

`attributes` JSONB příklad:
```json
{"style":"cartoon","polycount":1840,"formats":["glb","fbx","rbxm"],"roblox_category":"Prop",
 "tags":["medieval","sword","weapon"],"bundle_of":null,"source_pipeline":"trellis-v2"}
```

---

## 5. Platby — tři cesty

| Cesta | Nástroj | Proč | Poznámky |
|---|---|---|---|
| **Karta / fiat** | Stripe Checkout (hosted) + webhook | Nejmenší práce, PCI mimo nás. Olin má OSVČ + živnost → Stripe CZ účet. | Test mode nejdřív. Ceny v USD (cílovka = Roblox devs, globální). |
| **Crypto** | **BTCPay Server** (self-hosted, Docker na JODA) | Sedí do self-hosted filozofie, žádné KYC třetí strany, BTC/Lightning + volitelně USDT přes plugin. | Alternativa s menší prací: Coinbase Commerce / NOWPayments (hosted, ale závislost). Opus: zjisti aktuální stav BTCPay pluginů pro USDT/ETH, Olin chce primárně stablecoin. |
| **Free** | žádný provider | Lead magnet + test funnelu. | Vyžaduje e-mail (magic link login), vydá licenci `personal`, počítá se do statistik. |

Checkout flow je **jeden**: `order.create → vyber provider → provider.start → webhook/paid → fulfill`. Free = provider, který hned vrátí `paid`. Díky tomu je kód pro licence a download společný.

Daně: digitální zboží v EU = OSS/VAT MOSS téma. **MVP: prodávej jen "jako OSVČ bez DPH" pokud pod limitem, ale ukládej zemi kupujícího** (Stripe to vrátí, u crypta se ptej). Opus pořeš flag `tax_country`, ne výpočet daně. *(Olin: ověř si u účetního; toto není daňová rada.)*

---

## 6. Dodání (fulfillment) a licence

- Po `paid`: vygeneruj `license` (šablona `personal` / `commercial`, markdown → PDF přes existující tooling nebo jen text), vytvoř `download_token`, pošli e-mail (Resend / SMTP přes Cloudflare — zjisti, co shop už používá).
- Stránka **"Moje assety"** = seznam licencí + opakované stažení (token s obnovou).
- Formáty v balíčku: **GLB** (univerzální), **FBX** (Studio import), **RBXM** (přímý import do Roblox Studia) + `LICENSE.txt` + `README` (jak importovat, kolize, LOD).
- **Volitelně později:** "Publikovat do mého Roblox účtu" přes **Roblox Open Cloud API + OAuth 2.0** (upload assetu jménem kupujícího). Je to killer feature pro UX, ale podmíněná sekcí 7 — do MVP **ne**.

---

## 7. ⚠️ Roblox pravidla — ověřit PŘED stavbou checkoutu

Tohle rozhoduje, co se smí prodávat mimo Roblox. Opus: udělej rešerši (WebSearch/WebFetch na aktuální Roblox Terms of Use, Creator Store policies, Open Cloud docs) a zapiš do `docs/ROBLOX_COMPLIANCE.md` s datem a odkazy. Otázky:

1. Smí se **3D modely (vlastní IP, generické soubory GLB/FBX)** prodávat mimo Roblox za fiat/crypto? (Očekávání: ano — je to prodej vlastního digitálního díla, stejně jako na Sketchfab/Gumroad. Rizikové je jen prodávat "Roblox položky/Robux" mimo platformu.)
2. Je problém dodávat **`.rbxm`** (Roblox-specifický formát) off-platform? Nebo radši jen GLB/FBX + návod?
3. **Creator Store** — jaké jsou aktuální podmínky pro placené modely (Robux vs. USD), jaké jsou fee a požadavky na verifikaci účtu. Může to být **druhý kanál** vedle marketu (stejný asset, dvě místa).
4. **Open Cloud Assets API s OAuth** — lze dnes uploadovat asset do účtu jiného uživatele s jeho souhlasem? Jaké scopes, limity, moderace.
5. **AI-generovaný obsah** — má Roblox požadavek na označení / omezení pro assety z TRELLIS/ComfyUI? Jak je to s IP (Olin generuje z vlastních promptů/obrázků).
6. Moderace: každý upload na Roblox prochází moderací — pokud by asset neprošel, kupující má problém → do licence/README dát disclaimer + politiku refundu.

Dokud není bod 1–2 jasný, **MVP dodává jen GLB/FBX/OBJ** (platformově neutrální) a RBXM je opt-in.

---

## 8. Fáze

### Fáze 0 — Průzkum (1 session)
- [x] `docs/SHOP_REUSE.md` — co ze shop stacku použít (sekce 2)
- [x] `docs/ROBLOX_COMPLIANCE.md` — odpovědi na sekci 7
- [x] Zjistit API `ugc-backend` :8095 (endpointy pro list assetů, soubory, náhledy) — napsat `docs/UGC_API.md`
- [x] Rozhodnout (nový repo): nový repo `lioilsources/RobloxMarket` vs. monorepo se shopem s `VERTICAL=3d`. **Doporučení: nový repo, shop jako go module / zkopírované balíčky** — ať se shop nerozbije.

### Fáze 1 — MVP katalog + free tier (2–3 sessions)
- [~] Repo, Docker Compose na JODA (`deploy/market-stack.yaml`, cloudflared místo Caddy jako UGC), Postgres + migrace — **compose hotový, nenasazeno** (chybí tunnel `market-nas` + DNS + `.env`)
- [x] Ingest z `ugc-backend` → `products` (status `draft`)
- [x] Admin: web `/admin` (ne `ugc_studio`)  — nastavit cenu/tier/licenci, `listed`
- [x] Web: katalog s filtry (JSONB), detail s **3D viewerem** (`<model-viewer>` nebo three.js GLB), turntable náhledy — *turntably v ugc neexistují; náhled = koncept PNG + preview GLB*
- [x] Magic-link login (e-mail), **free download** přes token, licence `personal` (Resend klíč + doména zatím chybí → e-maily jen v logu)
- [x] Chat asistent (SSE + **pgvector**) nad katalogem, zatím mock režim (LLM na SPARKu přes broker) — přenést ze shopu, nový embedding job

### Fáze 2 — Placení (1–2 sessions)
- [ ] Stripe Checkout + webhook + `payment_events`
- [ ] BTCPay Server v Docker Compose na JODA + webhook (nebo hosted alternativa dle Fáze 0)
- [ ] Jednotný fulfillment (sekce 6), e-maily, stránka "Moje assety"
- [ ] Refund flow (ruční, přes admin)

### Fáze 3 — Roblox integrace a růst (po validaci)
- [ ] Paralelní listing na Creator Store (pokud dává smysl)
- [ ] "Publikovat do mého Robloxu" přes Open Cloud OAuth
- [ ] Bundles / packy, slevové kódy, metriky do Discord botu (Ol1nBot)
- [ ] Zapojit promo přes OpenClaw promo bot

### Definice hotovo pro "koncept vyzkoušen"
Olin dokáže: vygenerovat asset v UGC Factory → kliknout "publikovat" → asset je na market.ol1n.com s náhledem → cizí člověk ho koupí kartou NEBO crypto NEBO stáhne zdarma → dostane soubory + licenci → Olin vidí objednávku. Celé bez ručního zásahu.

---

## 9. Otevřené otázky pro Olina (Opus se ptej, když nejsou zodpovězené)

1. **Ceny:** jednotná cena za asset (např. $2–5) nebo tiery podle složitosti? Nějaký free pack jako lead magnet hned od začátku?
2. **Licence:** jen `personal` + `commercial`, nebo i exkluzivní prodej (asset se po koupi stáhne z marketu)?
3. **Crypto:** BTC/Lightning stačí, nebo nutně USDT/USDC (→ ovlivní volbu BTCPay vs. hosted)?
4. **Účet:** stačí e-mail magic link, nebo chceš i "Login with Roblox" (OAuth) hned v MVP?
5. **Kde běží frontend shopu** — Flutter web, nebo něco jiného? (Opus zjistí sám, ale Olin může ušetřit čas.)
6. **Refundy a moderace:** co když kupujícímu Roblox asset zamítne — refund ano/ne?

---

## 10. Instrukce pro Opuse — jak začít

1. Přečti tento soubor celý. Pak **Fáze 0** — nepiš produkční kód, dokud nejsou tři `docs/*.md` hotové.
2. Používej Director → sub-sessions: jedna na rešerši Roblox pravidel, jedna na audit shop repa, jedna na průzkum `ugc-backend` API. Výsledky sloučit do `docs/`.
3. Stack bez výjimky: Go + sqlc + golang-migrate + Postgres (JSONB/GIN) + Qdrant + Docker Compose na JODA + Caddy + Cloudflare Tunnel. Frontend = to, co má shop.
4. Platební kód piš s abstrakcí `PaymentProvider` (Stripe / BTCPay / Free) od prvního commitu, i když MVP spustí jen Free.
5. Každou fázi zakonči krátkým `CHANGELOG` záznamem + screenshotem funkčního flow (Olin to chce vidět, ne číst).
6. Když narazíš na rozpor mezi tímto plánem a realitou kódu/pravidel Robloxu — **realita vyhrává**, plán aktualizuj a oznam Olinovi.
