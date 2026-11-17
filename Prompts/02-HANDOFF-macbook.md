# Předávka plánu 02 na MacBook (průchod Roblox Studiem)

Na M2 je na Studio málo paměti, průchod Studiem a screenshoty se dělají na MacBooku.
Tenhle soubor je prompt pro novou session Claude Code na MacBooku v klonu `RobloxMarket`.

## Prompt

> Pokračuj v plánu 02 (`Prompts/02-PLAN-runbook-studio.md`), etapa 1 a dál. Stav a postup
> jsou v `docs/STUDIO_RUNBOOK.md`, kostra návodu v `web/src/lib/guide/guide.js`.
> Rozhodnutí: R6 návod jen anglicky; R7 zkušební importy smí jít na můj účet Roblox
> a smíš se ke Studiu připojit přes MCP; R8 publikace na Marketplace jen odkazem na
> dokumentaci Robloxu; R9 meč až později; R10 sekce na olin.now je `ugcmarket`, návod
> je i na marketu jako `/guide`.
>
> Postup:
> 1. Ověř, že vidíš nástroje MCP Studia (Studio běží, přihlášený, anglické UI, Baseplate).
>    Když ne, napiš mi přesně, co mám ve Studiu zapnout, a skonči.
> 2. ZIPy helmy (Ornate Samurai Kabuto Helmet) a batohu (Cute Structured Box Backpack)
>    stáhni jako `prod-check` z market.ol1n.com do `run/zips/`. Když market vrací 530
>    (Docker na JODĚ stojí), řekni mi to; Docker sám nestartuj.
> 3. Projdi import, texturu, Accessory Fitting Tool, zkoušku na figuríně a skript
>    `Humanoid:AddAccessory`. Dialog Import a výběr souboru odklikávám já, ty mi píšeš kroky,
>    fotíš (`screencapture` do `run/guide-raw/`) a ověřuješ přes MCP.
> 4. Text návodu piš jen z toho, co prošlo Studiem, ne z dokumentace. Ověřené oddíly
>    v `guide.js` přepni na `verified: true`, doplň `studioVersion` a `verifiedOn`,
>    výřezy do `docs/guide-shots.json`, pak `bin/guide-shots`. Zapiš výsledky a vady assetů
>    do `docs/STUDIO_RUNBOOK.md`.
> 5. Gate `make test` (e2e s `MARKET_TEST_DATABASE_URL` z `.env`), commit do `feat/studio-guide`.
> 6. Etapa 3b: export návodu na olin.now (`apps/ugcmarket/pages/guide.html` + překlady
>    jen pokud to plán chce). Větev olin.now `feat/ugcmarket` je v `run/olinnow-ugcmarket.bundle`
>    (není na GitHubu); v klonu `ol1n.now` ji načti do nového worktree, rozdělané větve
>    v hlavním checkoutu nech být.
>
> Nemerguj, nepushuj a nenasazuj, dokud neřeknu. Na konci napiš, co je ověřené ve Studiu a co ne.

## Stav při předání (2026-11-17)

- Repo zkopírované z M2 přes `git bundle` (bez remote): větve `chore/cloudflared-2026.11`
  → `docs/plans-02-03` → `feat/studio-guide` (stack, nic nemergnuté). `.env` zkopírovaný.
- Kopie na M2 je od teď neaktuální; zpátky přes bundle nebo `git fetch` po ssh.
- olin.now `feat/ugcmarket`: sekce ve 45 jazycích, build prošel, bez screenshotů, nepushnuté.
- Docker na JODĚ stál od 2026-11-16 10:26 (market 530); nespouštět, řeší Olin.
