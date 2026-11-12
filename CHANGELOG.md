# Changelog

## 2026-11-12 — Etapy A a B plánu 01: repo a srovnání s infra pravidly

- Git repo (`main`), `CLAUDE.md`, `.claude/settings.json` a gitignore podle DevConventions.
- Z docs a `.env.example` odstraněné domovské cesty, názvy svazků a LAN adresy (placeholdery).
- Docker na Macu už není potřeba: `market_dev` a `market_test` v `market-db` na JODĚ (port 5438 jen
  na LAN), `make deploy-db`; e2e test odmítne smazat databázi bez koncovky `_test`.
- Compose bez publikovaných portů: tunnel `market-nas` → Caddy → market-web; `config.yml` tunelu
  se generuje na JODĚ z UUID (`make tunnel-config`), deploy zapisuje do `~/deploy/events.jsonl`.
- Komentáře v kódu přeložené do angličtiny (kromě nasazené migrace).

## 2026-10-25 — Fáze 1: MVP katalog + free tier (lokálně ověřeno, nenasazeno)

- `api/` — Go market-api: schéma (products, files, media, embeddings, customers, sessions, orders, licenses, download_tokens, payment_events), sqlc + golang-migrate.
- Ingest z ugc-backend: `packed` a verdikt ≠ FAIL → draft; kopie souborů (FBX, textura, GLB, náhledy); archivace kusů, které z packed zmizí, a návrat do draftu.
- Katalog: filtry nad JSONB (port ze shopu), facety počítané z vylistovaných produktů, detail s `<model-viewer>`.
- Magic-link login (Resend / log), admin podle `ADMIN_EMAILS`, admin stránka (cena, free/paid, listing, ingest, objednávky).
- `PaymentProvider` abstrakce + Free provider; jednotný checkout → fulfillment → licence (personal) → ZIP s LICENSE.txt + README.txt.
- Chat asistent (SSE + pgvector) přepsaný do Go, mock režim bez LLM.
- `deploy/` — compose (pgvector, api, web, cloudflared) + Makefile `deploy-nas`.
- Ověřeno proti **reálnému ugc-backend na JODA** (12 kusů, 4 s): katalog → detail s 3D → košík → magic link → free checkout → „My assets“ → ZIP se skutečnými soubory. Screenshoty v `docs/screenshots/phase1-*.jpg`.
- Oproti plánu: Qdrant nahrazen pgvectorem; turntably/OBJ/RBXM v ugc-backend neexistují (náhled = koncept PNG + preview GLB); facet „style“ vypuštěn (obsahuje fragmenty promptu); názvy se odvozují z promptu.

## 2026-10-25 — Fáze 0: průzkum

- `docs/SHOP_REUSE.md` — audit shop stacku (`Dev/Shopify`): SvelteKit + FastAPI + pgvector, ne Go + Qdrant; co přenést, co psát nově, doporučení stacku.
- `docs/UGC_API.md` — API ugc-backend (:8095): endpointy, stavy jobů (`packed` = k listingu), dostupné formáty (chybí OBJ, RBXM, turntably), návrh integrace a změn v ugc-backend.
- `docs/ROBLOX_COMPLIANCE.md` — prodej vlastních GLB/FBX mimo Roblox OK; „Roblox“ nesmí být v názvu obchodu; RBXM opt-in; Creator Store jako 2. kanál (USD, 100 % net); OAuth upload technicky možný.
- Plán (`Prompts/00-ROBLOXMARKET_PLAN.md`) doplněn o sekci „Aktualizace po Fázi 0“ a odpovědi na otevřené otázky.
- Produkční kód zatím žádný (čeká na potvrzení stacku a názvu).
