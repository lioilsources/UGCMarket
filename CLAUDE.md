# RobloxMarket (Ol1n Market)

<!-- Max 80 řádků. Co se nevejde, jde do docs/ a sem jeden odkaz. Veřejné repo: žádné LAN IP, domovské cesty, ID. -->

## Účel
Obchod `market.ol1n.com` s 3D modely z UGC Factory (FBX, GLB, textura) za kartu, BTC nebo zdarma.
Veřejný název „Ol1n Market“: slovo Roblox v názvu obchodu zakazují pravidla značky.

## Build a test
```
cp .env.example .env   # DB na JODĚ, heslo z deploy/.env na JODĚ
make test              # gate: go vet + Go testy (e2e proti market_test) + svelte-check
make api / make web    # dev: API :8196, web :5173
make deploy-db         # jen Postgres na JODĚ
make deploy-nas        # celý stack = zveřejnění, jen na Olinovo „nasaď“
```

## Mapa
- `api/` Go market-api: katalog, ingest z ugc-backend, magic link, checkout, licence, stahování, chat.
- `api/internal/db/` sqlc (generované, needitovat) + migrace v `migrations/`.
- `web/` SvelteKit PWA, `/api/*` proxy na market-api.
- `deploy/` compose pro JODU: db, market-api, market-web, Caddy, cloudflared.
- `docs/` rešerše Fáze 0, screenshoty. `Prompts/` plány.
- `bin/` nástroje (`guide-shots`), `run/` běhová data (gitignore).

## Invarianty
- Platby jen přes `PaymentProvider`/`WebhookProvider` (`api/internal/shop/payments.go`).
  Webhooky idempotentní přes `payment_events`; fulfillment musí jít pustit opakovaně.
- E2e test maže schéma public, proto běží jen proti DB s koncovkou `_test`.
- Ingest z ugc-backend jen čte (GET); soubory se kopují, neodkazují. Nový kus je vždy draft.
- Download token se nikdy neloguje (`/api/download/…` je v logu redigovaný).
- Docker na Macu neběží: dev/test DB jsou `market_dev` a `market_test` v `market-db` na JODĚ.
- Na JODĚ žádné porty kromě Postgresu na LAN adrese (`DB_BIND`); web jen přes tunnel a Caddy.
- `.env`, `deploy/.env`, `deploy/cloudflared/{config.yml,credentials.json}` nejsou v gitu.
- Texty licencí v `api/internal/shop/templates/` jsou DRAFT; placený prodej až po právní kontrole.

## Odchylky
- `jodadeploy` zatím neexistuje: deploy přes `make deploy-nas`, který zapisuje do `events.jsonl`.
- Odchylka od lang.md: komentáře v migraci `0001_init` zůstávají česky (nasazená migrace se needituje).

## Kam dál
- `Prompts/01-PLAN-dalsi-kroky.md` aktuální plán, `Prompts/00-ROBLOXMARKET_PLAN.md` původní zadání.
- `Prompts/02-PLAN-runbook-studio.md` návod pro zákazníky, stav v `docs/STUDIO_RUNBOOK.md`.
- `Prompts/03-PLAN-snapchat.md` pokus se Snapchatem, podklad `docs/SNAPCHAT_RESEARCH.md`.
- `README.md` architektura a API, `CHANGELOG.md`, `docs/ROBLOX_COMPLIANCE.md`, `docs/UGC_API.md`.
