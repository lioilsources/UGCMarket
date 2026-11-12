# 01 — Plán dalších kroků (pro navazující session)

Navazuje na `00-ROBLOXMARKET_PLAN.md` (sekce „Aktualizace po Fázi 0“ platí přednostně).
Stav ověřen 12. 11. 2026. Pořadí etap je závazné: A → B → C → D → E.

## Výchozí stav

- **Hotovo lokálně:** Fáze 0 (tři dokumenty v `docs/`) a Fáze 1 (Go `api/`, SvelteKit `web/`,
  free flow ověřený na skutečných datech z UGC Factory, screenshoty v `docs/screenshots/`).
  Podrobnosti v `CHANGELOG.md` a `README.md`.
- **Není git repozitář.** Nic není commitnuté ani na GitHubu.
- **JODA:** v deploy adresáři projektu leží kód z 27. 10., `.env` (heslo k DB, admin e-mail,
  `DATA_DIR` na budoucí dataset `pool-market`) a sestavené image `market-api`, `market-web`.
  Nic neběží. Dataset `pool/market` neexistuje, tunel `market-nas` neexistuje,
  `market.ol1n.com` se nepřekládá.
- **Platby:** jen `Free` provider. Rozhraní `PaymentProvider`/`WebhookProvider`, obecný
  webhook handler a tabulka `payment_events` existují.
- **Chat:** mock režim (hash embeddingy, šablonová odpověď), bez LLM.

## Co se od Fáze 1 změnilo v pravidlech (kód to zatím porušuje)

Zdroj: `~/.claude/shared/infra.md`, `infra-architektura.md`, `ports.md`, `public-repo.md`.

1. **Na M2 neběží Docker.** `make dev-db` a e2e test počítají s lokálními kontejnery.
   Vývojová i testovací DB musí být na JODĚ (`market_dev`, `market_test`), adresa v `.env`.
2. **JODA: žádné porty na hostu, vlastní tunel, Caddy.** Compose teď publikuje `8096` a `8097`.
3. **Deploy přes `jodadeploy`** a zápis do `events.jsonl`, ne vlastní `make deploy-nas`.
   Ověř, jestli nástroj už existuje (`which jodadeploy`); když ne, nech Makefile a zapiš odchylku.
4. **Veřejné repo:** v `deploy/.env.example`, `docs/SHOP_REUSE.md` a `docs/UGC_API.md` jsou
   domovské cesty, názvy svazků nebo LAN adresy. Chybí `CLAUDE.md` repa.
5. **Interaktivní chat potřebuje fallback bez GPU.** Mock ho splňuje, ale musí zůstat i po
   zapnutí LLM (SPARK má okna a lease, chat nesmí spadnout).

## Etapa A — repo a hygiena (bez Olina)

- [x] Skill `bootstrap-repo`: `CLAUDE.md` (do 80 řádků; gate = `make test`), settings, gitignore.
- [x] Vyčistit soukromé údaje z bodu 4 (cesty a adresy nahradit placeholdery, fakta nechat
      v `infra.md`). Ověřit agentem `repo-auditor`.
- [x] `git init`, první commit se soubory přidanými jmenovitě. `.env` a `credentials.json` ne.
- [ ] Založení repa na GitHubu **až po rozhodnutí R1**.

## Etapa B — srovnání s infra pravidly

- [x] Dev a test DB na JODĚ: compose služba nebo druhá databáze v produkčním Postgresu
      projektu, port jen na LAN adrese. `Makefile` čte `DATABASE_URL` z `.env`, cíl `dev-db` zrušit.
      E2e test maže schéma, proto smí běžet jen proti `market_test` (pojistka na jméno DB v testu).
- [x] Compose bez publikovaných portů; web za Caddy podle vzoru ostatních stacků na JODĚ,
      tunel míří na Caddy. `make status` přepsat na kontrolu přes veřejnou adresu.
- [x] Deploy podle bodu 3. Zapsat projekt do `ports.md` (PR do DevConventions).
- [x] `make test` musí projít z M2 bez Dockeru.

Stav 12. 11.: hotovo. `market-db` běží na JODĚ (port 5438 jen na LAN, data dočasně v `data/`
deploy adresáře), `make test` prošel z M2 bez Dockeru. `jodadeploy` neexistuje → `make deploy-nas`
se zápisem do `events.jsonl` (odchylka v `CLAUDE.md`). Porty: PR DevConventions #5.

## Etapa C — nasazení (čeká na Olina, viz „Kroky Olina“)

- [ ] Po vytvoření tunelu: `make deploy-nas` (sám najde UUID, připraví `config.yml` a credentials), kontrola čtyř kontejnerů a `https://market.ol1n.com`.
- [ ] Po vytvoření datasetu: zastavit `db`, přesunout `data/` do datasetu, přepnout `DATA_DIR`, start.
- [ ] Přihlásit se jako admin, spustit ingest, zkontrolovat duplicitní samurajskou helmu
      (dva joby v UGC) a `smoke-001`. Listing a free pack nastaví Olin, nebo je potvrdí.
- [ ] Projít free flow na produkci, uložit screenshot, zapsat do `CHANGELOG.md`.
- [ ] Po dodání klíče Resend: ověřit doručení magic linku na skutečnou adresu.

Nasazení je zveřejnění. Proveď ho, až Olin řekne „nasaď“.

## Etapa D — Fáze 2: platby

- [ ] **Stripe:** Checkout Session, webhook s ověřením podpisu, idempotence přes
      `payment_events`, testovací režim s testovacími klíči. Vrácení peněz (`refunded` →
      zneplatnit licenci a download tokeny).
- [ ] **BTCPay (BTC/Lightning):** provider + webhook. Kde poběží, je rozhodnutí R3.
- [ ] Komerční licence jako druhá cena položky (R2), `license_type` už ve schématu je.
- [ ] Checkout: zaškrtnutí souhlasu se ztrátou práva na odstoupení (EU, digitální obsah),
      uložit čas souhlasu k objednávce; `tax_country` vyplnit z platební brány.
- [ ] Reklamace kvůli moderaci Robloxu do 14 dnů (slib v `README.txt` balíčku): admin akce.
- [ ] E2e test obou providerů s falešným webhookem; screenshot placeného flow; CHANGELOG.

## Etapa E — Fáze 3: Roblox integrace

Podle `00-ROBLOXMARKET_PLAN.md` a `docs/ROBLOX_COMPLIANCE.md` (RBXM jen opt-in, Creator Store
jako druhý kanál). Začít až po ostrém provozu plateb. Před začátkem znovu ověřit pravidla Robloxu.

## Návrhy změn v ugc-backend (jiné repo, jen návrh Olinovi)

Příznak „publikovat na market“, turntable rendery, OBJ, GLB exportovaný z převedeného modelu.
Detail v `docs/UGC_API.md`.

## Rozhodnutí vlastníka

- **R1:** repo `lioilsources/RobloxMarket` veřejné, nebo soukromé? (Obsahuje návrhy licencí a obchod.)
- **R2:** cena komerční licence a potvrzení jednotné ceny 3 USD.
- **R3:** BTCPay na JODĚ s prořezaným uzlem (JODA má málo paměti), nebo jiné řešení?
- **R4:** LLM pro chat: trvale přes gateway na SPARKu s fallbackem na mock, nebo nechat mock?
      Zapnutí jen přes spark-broker.
- **R5:** právní kontrola textů licencí (teď DRAFT) a ověření licencí SF3D a Illustrious
      pro komerční prodej. Bez toho nespouštět placený prodej.

## Kroky Olina (blokují etapu C)

1. Na JODĚ vytvořit ZFS dataset projektu a pod ním dataset pro Postgres (vzor v `infra.md`,
   oddíl JODA; chce sudo).
2. Na JODĚ založit tunel `market-nas` a DNS záznam `market.ol1n.com` (agentovi to blokuje
   auto mode jako změnu DNS).
3. Resend: účet, ověřená doména, API klíč do `.env` na JODĚ.

## Pravidla pro práci

- Realita vítězí: když se plán rozchází se skutečností, oprav plán a řekni to Olinovi.
- Každá etapa končí zápisem do `CHANGELOG.md` a screenshotem funkčního flow.
- Skutečné platební klíče ani hesla agent nezadává; ostré klíče do `.env` dává Olin.
- SPARK jen přes `sparkctl`. Na JODĚ jen deklarované deploy akce, žádné sudo.
