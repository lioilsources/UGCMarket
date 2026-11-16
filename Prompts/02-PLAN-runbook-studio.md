# 02 — Návod pro zákazníky: model z marketu na postavě v Roblox Studiu

Navazuje na `01-PLAN-dalsi-kroky.md`. Stav ověřen 16. 11. 2026. Pořadí etap je závazné: 0 → 4.

## Cíl

Zákazník si stáhne ZIP a neví, co s ním. Potřebuje návod, podle kterého dojde od souboru
k modelu nasazenému na postavě v Roblox Studiu. Čtenář je často začátečník (13+), nezná naši
pipeline a anglicky čte průměrně.

Hotovo je, když člověk, který postup nezná, podle návodu nasadí na postavu helmu i batoh
a nemusí se na nic ptát.

## Výchozí stav (ověřeno 16. 11.)

- **Market běží** na `market.ol1n.com`. Na produkci je stav větve `chore/cloudflared-2026.11`
  (`main` + nová verze tunelu), která čeká na merge. `phase2-payments` není mergnutá ani nasazená.
- **Vylistováno 5 kusů**, z toho 4 zdarma (kabuto helma, dva batohy, kabelka).
- **ZIP** (`api/internal/shop/download.go`): `{slug}.fbx`, `{slug}_tex.png`, `{slug}.glb`,
  `LICENSE.txt`, `README.txt` (šablona `api/internal/shop/templates/readme.txt`).
- **FBX:** do 4000 trojúhelníků, měřítko ve studs, +Z nahoru, −Y dopředu, jeden materiál,
  textura není vložená. Produkt má v `attributes` pole `category`, `roblox_attachment`
  a `bbox_studs` (odpovídá limitům Robloxu pro měřítko „Normal“).
- **Studio na M2:** Roblox Studio 0.737 ze září (při spuštění se aktualizuje), obsahuje binárku
  `StudioMCP`. Terminál má povolené nahrávání obrazovky, takže `screencapture` funguje.
  Nástroje pro ovládání plochy (computer use) v plánovací session nebyly; Claude in Chrome
  desktopové aplikace neovládá.
- **Návod zatím neexistuje.** Detail produktu slibuje „README.txt (import guide for Roblox
  Studio)“, README má k importu čtyři řádky a cestu v menu, která už podle dokumentace neplatí.

## Co vím ze souborů (ve Studiu zatím nevyzkoušeno)

Tohle jsou nejpravděpodobnější místa, kde zákazník uvízne. Ověř je jako první.

1. **Textura se nejspíš sama nenačte.** FBX z linky `trellis` odkazuje na `model_tex.png`
   a na absolutní cestu build stroje; FBX z linky `sf3d` (např. box batoh) nemá odkaz žádný.
   V ZIPu se textura jmenuje `{slug}_tex.png`. Zákazník ji tedy asi musí nahrát ručně
   (Asset Manager → moderace → `TextureID`).
2. **Bod uchycení se nejspíš neimportuje.** FBX obsahuje prázdný objekt `<Attachment>_Att`
   (např. `HatAttachment_Att`). Specifikace Robloxu říká, že import bodů uchycení z cizích
   nástrojů u rigid accessories není podporovaný; správný bod vytvoří Accessory Fitting Tool.
3. **Meč není accessory.** Kategorie `sword` má `RightGripAttachment`, což není typ accessory
   pro Marketplace. Patří k němu jiný postup (Tool). Žádný meč teď není vylistovaný.
4. **Import nahrává na účet.** Importer má ve výchozím stavu zapnuté „Upload to Roblox“:
   mesh i textura se nahrají na účet přihlášeného uživatele a projdou moderací Robloxu.

## Co má vzniknout

1. **Stránka `/guide`** na webu, anglicky (jazyk obchodu), čitelná i na telefonu.
2. **Screenshoty z Roblox Studia na M2** ve `web/static/guide/`.
3. **Odkazy na návod:** u tlačítka Download v `/me`, v detailu produktu (ideálně rovnou na
   oddíl podle kategorie, např. `/guide#back`), v patičce a v `README.txt` v balíčku.
4. **`README.txt`:** opravený postup importu podle skutečnosti, krátký „Quick start“ a odkaz.
5. **`docs/STUDIO_RUNBOOK.md`** česky, pro nás: co se ve Studiu skutečně děje, verze Studia
   a datum, jak přefotit screenshoty, nalezené vady assetů.
6. Zápis do `CHANGELOG.md` a screenshot hotové stránky v `docs/screenshots/`.
7. **Tentýž návod na `olin.now`** v sekci marketu (Olinovo zadání z 16. 11., viz Etapa 3b).

## Etapa 0 — prostředí (než začneš psát)

- [ ] Zjisti, čím budeš Studio ovládat. Zkus `ToolSearch` na computer use. Když není,
      kliká Olin podle tvého seznamu kroků a ty fotíš a ověřuješ. Oprávnění systému neměň.
- [ ] Studio MCP (po Olinově souhlasu, viz R7): `execute_luau`, `search_game_tree`,
      `inspect_instance`, `start_stop_play`, `screen_capture` (jen viewport), `save_place`.
      Lokální soubor importovat neumí, dialog Import musí někdo odklikat.
- [ ] Vyzkoušej `screencapture` na okno Studia a ověř, že na snímku je okno, ne jen plocha.
- [ ] Přidej `run/` do `.gitignore` (zatím tam není) a stáhni si do něj dva balíčky z produkce:
      **Ornate Samurai Kabuto Helmet** (helma, linka `trellis`) a **Cute Structured Box
      Backpack** (záda, linka `sf3d`). Pokrývají obě linky i dva body uchycení.
      Na produkci existuje testovací zákazník `prod-check` (smyšlená adresa na example.com),
      který obě licence má. Resend není zapojený, přihlašovací odkaz je v logu `market-api`.
- [ ] Kus „Neon Samurai Helmet“ (`smoke-001`) nepoužívej, je to testovací koule.

## Etapa 1 — projít postup naostro

Projdi celý postup s oběma kusy a zapisuj do `docs/STUDIO_RUNBOOK.md`, co se stalo.
Dokumentace Robloxu je vodítko, pravda je Studio na M2.

| Krok | Podle dokumentace | Co ověřit |
|---|---|---|
| Import | **File ⟩ Import** nebo **Home ⟩ Import**, pak **Import** | Stojí model rovně, čelem správně, ve správné velikosti (`bbox_studs`)? Jaká nastavení musí zákazník změnit? |
| Textura | sama přes `TextureID`, jinak Asset Manager ⟩ Import, po moderaci ID do `TextureID` | Načte se sama? Načte se, když se soubor vedle FBX jmenuje `model_tex.png`? Jak dlouho trvá moderace a co uživatel vidí? |
| Bod uchycení | z cizích nástrojů se neimportuje | Co udělá importer s `<Attachment>_Att`? Vadí, když tam zůstane? |
| Accessory | **Avatar ⟩ Accessory** → **Part** → **Next** → typ **Accessory** + umístění → **Next** → **Generate MeshPart Accessory** | Sedí model na postavě bez ručního posouvání? Je novější nástroj Avatar Setup jednodušší? |
| Zkouška | Rig Builder + přetažení Accessory pod postavu, nebo Avatar Setup ⟩ Test in Experience | Která cesta je pro začátečníka nejkratší? |
| Použití ve hře | `Humanoid:AddAccessory` nebo `HumanoidDescription` | Minimální skript, který nasadí accessory každému hráči; ověř v playtestu. |

- [ ] Publikaci na Marketplace nezkoušej (poplatky, ověření totožnosti, nevratné).
- [ ] Když najdeš vadu assetu (textura, orientace, měřítko, uchycení), zastav se a řekni
      to Olinovi. Malá oprava v marketu (např. název textury v ZIPu) do rozsahu patří.
      Oprava v UGC Factory je jiné repo: jen návrh Olinovi.

## Etapa 2 — text návodu

Osnova (anglicky, názvy oddílů uprav podle toho, co vyšlo v Etapě 1):

1. What's in the ZIP: který soubor je k čemu, FBX pro Studio, GLB pro jiné nástroje.
2. Before you start: Studio, přihlášení, upozornění, že import nahrává model a texturu
   na jejich účet a prochází moderací Robloxu.
3. Import the model.
4. Apply the texture.
5. Turn it into an accessory, s tabulkou kategorie → typ a umístění:
   hat a helmet → Hat, hair → Hair, back → Back, waist → Waist, neck, shoulder, front, face.
6. Try it on a character.
7. Use it in your experience, s krátkým skriptem.
8. Troubleshooting: model bez textury, špatná velikost, sedí křivě, zamítnutá moderace
   (odkaz na `/licenses`).
9. Publishing: jeden odstavec a odkaz na dokumentaci Robloxu, žádné sliby.

Pravidla textu:

- Číslované kroky, jedna akce na krok. Názvy tlačítek a menu tučně a přesně tak, jak jsou
  ve Studiu. Krátké věty, žádný žargon bez vysvětlení.
- Nepiš nic, co jsi sám ve Studiu neprošel.
- Formulace podle `docs/ROBLOX_COMPLIANCE.md` (oddíl „Dopad na plán“, body 1 a 2): modely
  „compatible with Roblox Studio“, ne „Roblox items“ ani „UGC items for your avatar“. Nic
  o Robuxech. Patička s prohlášením o nezávislosti už na webu je.

## Etapa 3 — web

- [ ] `web/src/routes/guide/+page.svelte`, statická stránka bez API. Obrázky s rozměry,
      `loading="lazy"`, popisky `alt`, kliknutím zvětšit.
- [ ] Odkazy podle „Co má vzniknout“, bod 3. V `README.txt` adresu skládej z `PUBLIC_URL`.
- [ ] Změny v `p/[slug]/+page.svelte` a `me/+page.svelte` drž malé: `phase2-payments`
      tytéž soubory mění a bude se mergovat.
- [ ] Gate `make test`. E2e test už ZIP rozbaluje (`api/internal/httpapi/e2e_test.go`):
      přidej do něj kontrolu, že `README.txt` obsahuje odkaz na návod.

## Etapa 3b — návod na olin.now

Olin chce návod nasadit i na `olin.now` do sekce marketu. Stav ověřen 16. 11.:

- Repo `lioilsources/ol1n.now` je **veřejné**, statický web na GitHub Pages. Lokální checkout
  leží vedle tohoto repa. Než začneš, přečti jeho `CLAUDE.md` (oddíly „Adding a New App“,
  „App subpages“ a „Deployment“).
- **Sekce marketu tam zatím není.** V `apps/` je jen `ugcfactory`. Založ ji vedle něj jako
  `ugcmarket` (R10) podle „Adding a New App“: `apps/ugcmarket/meta.md` a `icon.png` 512×512
  (bez ikony build tiše použije písmeno).
- **Návod jako podstránka sekce:** `apps/ugcmarket/pages/guide.html` (ručně psaný fragment těla
  stránky; podstránky jsou jen v základním jazyce). Screenshoty jako WebP.
- **Jeden zdroj obsahu.** Text a obrázky vznikají v tomto repu a do `ol1n.now` se přenášejí
  skriptem v `bin/`, stejně jako tamní `make import-*` berou data ze sousedních checkoutů.
  Neudržuj dvě ručně psané kopie.
- **Veřejný název je „Ol1n Market“**, ne „RobloxMarket“ (pravidla značky, viz `CLAUDE.md`).
  Slug `ugcmarket` je jen adresa sekce, název obchodu nemění.
- Sekce odkazuje na `market.ol1n.com` a obráceně. Pro text i screenshoty platí stejná
  pravidla jako na marketu.
- **Git a nasazení:** lokální checkout je na cizí rozdělané větvi s necommitnutými soubory.
  Nesahej na ni, pracuj ve worktree z `origin/main`. Samostatný PR v `ol1n.now`;
  `make deploy` staví z `origin/main`, takže nasazení čeká na Olinovo „mergni“ a „nasaď“.

## Etapa 4 — ověření a předání

- [ ] Nech návod přečíst subagenta, který zná jen text a screenshoty, a oprav nejasné kroky.
- [ ] Olin projde návod s jiným kusem (např. Mini Flap Shoulder Bag) bez nápovědy.
- [ ] CHANGELOG, screenshot stránky. Nasazení je zveřejnění: až na Olinovo „nasaď“.
      Platí zvlášť pro market a zvlášť pro olin.now.

## Screenshoty — pravidla

- Skutečné Studio na M2, anglické rozhraní, jedno téma po celý návod. Do `STUDIO_RUNBOOK.md`
  zapiš verzi Studia a datum.
- Celé okno přes `screencapture`, samotný viewport může přes MCP `screen_capture`.
- Na snímku nesmí být: logo Robloxu ani ikona aplikace, uživatelské jméno a avatar účtu
  (pravý horní roh Studia), domovské cesty v dialogu pro výběr souboru, jiná okna.
  Ořízni na panel, o který v kroku jde.
- Postava je figurína z Rig Builderu, ne avatar účtu.
- Prvek, na který se kliká, zvýrazni rámečkem. Ořez a zvýraznění dělej skriptem v `bin/`
  se seznamem výřezů, aby šly snímky po změně Studia přefotit stejně.
- Surové snímky nech v `run/` a necommituj. Do repa jdou jen ořezané:
  `web/static/guide/NN-nazev.webp`, šířka do 1600 px, zhruba do 200 kB.

## Rozhodnutí vlastníka

- **R6:** návod jen anglicky (doporučuji, obchod je anglicky), nebo i česky?
- **R7:** smí zkušební importy jít na Olinův účet Roblox? Každý import nahraje mesh a texturu
  na účet a projde moderací. Druhá možnost je vedlejší účet. A smí se Claude Code připojit
  ke Studiu přes MCP? (Trvalá konfigurace; klient pak čte a mění otevřené místo.)
- **R8:** publikace na Marketplace jen jako odkaz na dokumentaci Robloxu (doporučuji), nebo
  vlastní oddíl?
- **R9:** meč: vylistovat a popsat postup přes Tool teď, nebo později?
- **R10 (rozhodnuto 16. 11. 2026):** slug sekce na olin.now je **`ugcmarket`**.
  K druhé půlce (návod i na marketu jako `/guide`, nebo jen odkaz na olin.now) se Olin
  nevyjádřil. Platí proto doporučení: návod je na obou místech z jednoho zdroje.

## Kroky Olina

1. Rozhodnout R6–R9 (R10 je rozhodnuté) a merge větví `chore/cloudflared-2026.11`
   a `docs/plans-02-03`.
2. Na M2 spustit Studio, nechat ho aktualizovat, přihlásit se (heslo agent nezadává),
   přepnout rozhraní do angličtiny a založit prázdné místo (Baseplate) uložené lokálně.
3. Pro MCP: **Assistant ⟩ … ⟩ Manage MCP Servers ⟩ Enable Studio as MCP server**,
   pak **Quick connect ⟩ Claude Code**.
4. Být asi půl hodiny u M2 kvůli klikání v dialozích, pokud session nemá computer use.
5. Na konci říct „nasaď“.

## Git

- Větev `feat/studio-guide`. Základ je `main`, pokud už obsahuje `docs/plans-02-03`.
  Jinak vycházej z `docs/plans-02-03`: stojí na `chore/cloudflared-2026.11`, tedy na stavu
  produkce (nasazení ze samotného `main` by vrátilo starou verzi tunelu), a jsou v ní plány.
- Na `phase2-payments` nesahej.
- Plány 02 a 03 jsou commitnuté ve větvi `docs/plans-02-03`.

## Pravidla pro práci

- Realita vítězí: když se plán rozchází se Studiem, oprav plán a řekni to Olinovi.
- Heslo k Robloxu agent nezadává. Nic se nepublikuje na Marketplace, nic se nekupuje.
- Na JODĚ jen deklarované deploy akce, žádné sudo. SPARK není potřeba.

## Dokumentace Robloxu (přípona `.md` vrací čistý text)

- Importer: https://create.roblox.com/docs/studio/importer
- Import rigid accessory: https://create.roblox.com/docs/art/accessories/creating-rigid/importing
- Accessory Fitting Tool: https://create.roblox.com/docs/art/accessories/creating-rigid/converting
- Specifikace (limity, body uchycení): https://create.roblox.com/docs/art/accessories/specifications
- Avatar Setup: https://create.roblox.com/docs/avatar-setup
- Studio MCP: https://create.roblox.com/docs/studio/mcp
