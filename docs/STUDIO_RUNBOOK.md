# Roblox Studio — runbook pro návod `/guide`

Pro nás, ne pro zákazníky. Co se ve Studiu skutečně děje, jak přefotit screenshoty
a nalezené vady assetů. Plán: `Prompts/02-PLAN-runbook-studio.md`.

## Stav

| | |
|---|---|
| Verze Studia | zatím neprošlo (na M2 je 0.737, při spuštění se aktualizuje) |
| Datum průchodu | — |
| Ověřené oddíly návodu | jen „What's in the ZIP“ (popisuje náš balíček, ne Studio) |

Stav k 2026-11-16: Studio na M2 neběželo, MCP server Studia nebyl připojený a Docker na JODĚ
byl zastavený (market vracel 530), takže se nestáhly ani testovací ZIPy. Hotová je kostra:
stránka, odkazy, README, skript na screenshoty.

## Jak je návod postavený

- **Jeden zdroj:** `web/src/lib/guide/guide.js`. Stránka `/guide` i budoucí export na olin.now
  (Etapa 3b) renderují tentýž obsah přes `web/src/lib/guide/render.js`.
- Oddíl má `verified: true`, až když prošel Studiem na M2. Neověřené oddíly se ukážou jen
  ve `vite dev` (štítek „draft“ a TODO body), na produkci nejsou. Dokud nejsou ověřené všechny,
  produkce ukazuje upozornění, že návod vzniká.
- Kotvy pro produkty: detail produktu odkazuje na `/guide#cat-<category>`. Řádky tabulky
  kategorií v oddílu `accessory` musí mít `rowIds: ['cat-hat', …]`.
- `README.txt` v balíčku má odkaz `{PUBLIC_URL}/guide`; e2e test to kontroluje.

## Průchod ve Studiu (Etapa 1)

Kusy: **Ornate Samurai Kabuto Helmet** (linka trellis, hlava) a **Cute Structured Box Backpack**
(linka sf3d, záda). Stáhnout jako testovací zákazník `prod-check` do `run/zips/`
(přihlašovací odkaz je v logu `market-api`, Resend není zapojený).

Zapisuj sem ke každému kroku: přesné názvy tlačítek, co se stalo, kolik to trvalo, screenshot.

| Krok | Výsledek |
|---|---|
| Import (cesta v menu, nastavení, orientace, velikost proti `bbox_studs`) | — |
| Co importer udělá s `<Attachment>_Att` | — |
| Textura: načte se sama? Jinak nahrání PNG a nastavení na MeshPart, doba moderace | — |
| Accessory Fitting Tool vs. Avatar Setup | — |
| Zkouška na figuríně z Rig Builderu | — |
| Skript `Humanoid:AddAccessory` v playtestu | — |

## Screenshoty

1. Celé okno Studia: `screencapture -x -l <windowid> run/guide-raw/<název>.png`
   (ID okna: `osascript -e 'tell app "System Events" to …'` nebo `screencapture -iW` ručně).
   Samotný viewport jde i přes MCP `screen_capture`.
2. Do `docs/guide-shots.json` zapiš výřez, rámečky kolem klikaného prvku a rozmazání
   všeho soukromého (jméno a avatar účtu vpravo nahoře, domovské cesty v dialogu).
3. `bin/guide-shots` vyrobí `web/static/guide/NN-nazev.webp` (do 1600 px, cíl ~200 kB)
   a rozměry do `web/src/lib/guide/shots.json`.
4. Surové snímky zůstávají v `run/` a necommitují se.

Po aktualizaci Studia stačí přefotit surové snímky se stejnými názvy a pustit skript znovu;
když se posune rozhraní, opravit souřadnice v `guide-shots.json`.

## Vady assetů

Zatím žádné ověřené. Podezření z plánu (ověřit jako první): textura se nenačte sama, bod
uchycení se neimportuje.

## Co má udělat Olin na M2

1. Spustit Roblox Studio, nechat ho aktualizovat a přihlásit se (heslo agent nezadává).
2. Přepnout rozhraní do angličtiny a založit prázdné místo **Baseplate**, uložit lokálně.
3. Zapnout MCP server Studia: **Assistant ⟩ … ⟩ Manage MCP Servers ⟩ Enable Studio as MCP
   server**, pak **Quick connect ⟩ Claude Code**. Session Claude Code pak restartovat, aby
   nástroje Studia viděla.
4. Být asi půl hodiny u M2: dialog Import a výběr souboru agent neodklikne, kliká Olin podle
   seznamu kroků a agent fotí a ověřuje.
