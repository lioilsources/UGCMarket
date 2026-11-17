# Roblox Studio — runbook pro návod `/guide`

Pro nás, ne pro zákazníky. Co se ve Studiu skutečně děje, jak přefotit screenshoty
a nalezené vady assetů. Plán: `Prompts/02-PLAN-runbook-studio.md`.

## Stav

| | |
|---|---|
| Verze Studia | 0.747.0.7471481 (MacBook) |
| Datum průchodu | 2026-11-17 |
| Ověřené oddíly návodu | všechny kromě Publishing (jen odkaz na dokumentaci Robloxu, R8). Helma (trellis): import, textura v importu, accessory Hat, zkouška, skript v playtestu. Batoh (sf3d): textura po importu, měřítko, accessory Back, zkouška na figuríně |

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
- Export na olin.now: `bin/export-olinnow ../ol1n.now` zapíše `apps/ugcmarket/pages/guide.html`
  (čeština, výchozí jazyk webu) a `guide.<lang>.html` pro dalších 44 jazyků sekce, snímky do
  `static/apps/ugcmarket/guide/`. Ručně tam needitovat. Styly `.ugc-guide` jsou v
  `assets/css/store.css` tamního repa.
- **Překlady:** `web/src/lib/guide/i18n/<lang>.json`, klíče z `en.json` (`i18n.js`). Po změně
  `guide.js` export odmítne běžet a vypíše změněné klíče: `bin/export-olinnow --update-en`, pak
  ty klíče přeložit ve všech jazycích. Tučné názvy z UI Studia, `kód` a URL zůstávají anglicky.
  Market `/guide` je dál jen anglicky. Překlady dělali agenti 2026-11-17, rodilý mluvčí je nečetl.

## Průchod ve Studiu (Etapa 1)

Kusy: **Ornate Samurai Kabuto Helmet** (linka trellis, hlava) a **Cute Structured Box Backpack**
(linka sf3d, záda). Stáhnout jako testovací zákazník `prod-check` do `run/zips/`
(přihlašovací odkaz je v logu `market-api`, Resend není zapojený).

Zapisuj sem ke každému kroku: přesné názvy tlačítek, co se stalo, kolik to trvalo, screenshot.

Průchod 2026-11-17 na MacBooku (MacBookPro18,3), Roblox Studio **0.747.0.7471481**, anglické
rozhraní, tmavé téma, místo `guide-test.rbxl` z šablony Baseplate. Klikal agent přes CGEvent
(terminál má oprávnění Zpřístupnění), stav ověřoval přes Studio MCP. Helma prošla celá. Batoh
ze sf3d má vadu měřítka (viz Vady assetů); prošel s ručním zmenšením a ruční texturou, obojí
je v návodu jako postup a v Troubleshooting.

| Krok | Výsledek |
|---|---|
| Import (cesta v menu, nastavení, orientace, velikost proti `bbox_studs`) | **Home ⟩ Import** (ikona krychle se šipkou mezi Script a Material) → systémový dialog „Open 3D File (FBX, OBJ, GLB, GLTF)“ → okno **Import Preview**. Pozor: **Terrain Editor** má vlastní tlačítko **Import** (výšková mapa), Olin na něj napoprvé klikl. Výchozí nastavení stačí: Upload to Roblox ✓, Add to Workspace ✓, Set Pivot to Scene Origin ✓, Anchored ✗, File Transform World Forward = Front, World Up = Top. Helma stojí rovně, čelem k −Z (Front). Velikost helmy 1.86 × 1.81 × 1.87 (metadata říkají [1.87, 2.5, 1.87], viz vady). Po importu se dole otevře **Import Queue** se zelenou fajfkou. Výsledek: `Model` › `world` › `MeshPart Handle`. |
| Co importer udělá s `<Attachment>_Att` | Ve stromu importu se vůbec neobjeví, do místa se nedostane. Nevadí. Bod uchycení vytvoří Accessory Fitting Tool. |
| Textura: načte se sama? | **Trellis (helma):** ne. Řádek `base_color_texture` je odškrtnutý s chybou „Could not read the color map texture file associated with this object.“ Oprava přímo v importu: kliknout na řádek `base_color_texture` (pod `world ⟩ Handle`), vpravo u **Color File Path** ikona složky → vybrat `{slug}_tex.png` ze ZIPu. Řádek se zaškrtne, náhled je barevný, po importu má `Handle` vyplněné `TextureID`. Žádný Asset Manager, žádné čekání na moderaci v rozhraní. **Sf3d (batoh):** FBX na texturu neodkazuje, v importu žádný řádek s texturou není; `TextureID` je po importu prázdné. **Sf3d, textura po importu (ověřeno):** **Home ⟩ Assets** otevře **Asset Manager** → tlačítko **Import** v Asset Manageru → vybrat `{slug}_tex.png` → v **Import Queue** **Start Import** → po zelené fajfce je obrázek v Asset Manageru, číslo ve sloupci **ID** (kontextové menu Copy se neotevřelo, ID se opisuje) → v Exploreru vybrat `Handle` → Properties, filtr `texture` → **TextureContent** = `rbxassetid://<ID>` (Olin psal ručně, funguje hned, žádné čekání na moderaci v rozhraní). |
| Měřítko (sf3d) | Vybrat horní `Model` v Exploreru → Properties, filtr `scale` → **Scale** (sekce Pivot) = `0.3` → batoh 2.40 × 2.07 × 1.46, `TextureID` zůstane. Vlastnost a pole v UI ověřené (Olinův snímek); hodnotu nastavil agent přes MCP (`Model:ScaleTo`), protože jeho psaní do pole se netrefilo (otevřelo přejmenování v Exploreru). |
| Accessory Fitting Tool vs. Avatar Setup | **Avatar ⟩ Accessory** otevře panel **Accessory Fitting Tool**. Kroky: **Select** (klik na model ve viewportu vybere rovnou `Handle`) → **Continue** → **Asset Type**: **Accessory** + typ v rozbalovacím seznamu (abecedně: Back, Eyebrows, Eyelashes, Face, Front, Hair, Hat, …), **Select Body Type** nechat Classic Blocky → do scény přibude figurína **Mannequin**, helma po výběru **Hat** sama skočí na hlavu → **Continue** → krok úprav (modrý kvádr = povolený prostor) → **Generate MeshPart Accessory** dole v panelu (v úzkém panelu je popisek oříznutý; pomohlo zavřít Toolbox a Import Queue). Potvrzení „Handle has been successfully generated into Workspace!“ je samostatné okno, které se schová za hlavní okno Studia. Výsledek: `Accessory` jménem **Handle** (podle partu!), `AccessoryType = Unknown`, uvnitř `MeshPart Handle` s texturou a nový `Attachment HatAttachment`. Avatar Setup nezkoušen. **Batoh:** výběr přes klik do pole **Part** a pak na model ve viewportu (bez kliknutí do pole se vybere celý Model a pole zůstane prázdné) → **Accessory** + **Back** → batoh skočí figuríně na záda → **Continue** → **Generate MeshPart Accessory** → OK. Výsledek: `Accessory` **Handle**, `BodyBackAttachment`, textura i velikost zachované. |
| Zkouška na figuríně z Rig Builderu | **Avatar ⟩ Character** → **Generate Rig** (R15, Masculine, seznam postav) → **Block Avatar (2012)** → v `Workspace` přibude `Rig`. Accessory přetáhnout v Exploreru na `Rig`: helma sedí na hlavě (0.6 studu nad středem hlavy, `AccessoryRigidConstraint`), čelem dopředu, bez ručního posouvání. Batoh: přesunut do `Rig` přes MCP (přetažení v Exploreru ověřeno u helmy), sedí na zádech (0.5 studu za středem trupu, `AccessoryRigidConstraint`). Past: v `Workspace` po Fitting Toolu zůstane i původní importovaný `Model` se skoro stejným jménem; Olin napoprvé přetáhl ten. Do návodu: původní model smazat a accessory přejmenovat. Přejmenování přes F2 na Macu nefungovalo (bez fn), přes pole Name v Properties se agentovi nepodařilo psát; přejmenováno přes MCP, v UI neověřeno. |
| Skript `Humanoid:AddAccessory` v playtestu | Accessory v `ServerStorage` jako `KabutoHelmet`, `Script` v `ServerScriptService` (text skriptu je v `web/src/lib/guide/guide.js`, oddíl `game`). Playtest (R15, avatar účtu): hráč dostal helmu na `HatAttachment` s `AccessoryRigidConstraint`, konzole bez chyb. `AccessoryType = Unknown` nevadí. |

## Screenshoty

Surové snímky ze 2026-11-17 jsou v `run/guide-raw/` (01–22), výřezy v `docs/guide-shots.json`.
Na snímcích Fitting Toolu je batoh, protože vznikly v pozdějším rozložení panelů.
Neověřeno: Back u kusu z linky trellis (`mini-flap-shoulder-bag-in-quilted-lambskin`), prod-check
na něj nemá licenci a objednávka by vznikla na produkci.

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

Ověřeno 2026-11-17 ve Studiu 0.747:

1. **Sf3d: měřítko (blokující).** `cute-structured-box-backpack-in-woven-straw` se importuje
   jako 8.0 × 6.9 × 4.9 studu, tedy zhruba čtyřikrát větší než postava (trup ~2 study).
   Produkt nemá `bbox_studs` (`None`). Týká se obou vylistovaných kusů z linky sf3d
   (`cute-structured-box-backpack-in-woven-straw`, `cute-classic-school-backpack-in-soft-suede`).
   Kusy z linky trellis `bbox_studs` mají. Oprava patří do UGC Factory (škálování na studs
   v exportu sf3d), tady jen návrh. Snímek `run/guide-raw/90-defect-backpack-scale.jpg`.
2. **Sf3d: FBX bez odkazu na texturu.** Importér nenabídne řádek textury, `TextureID` zůstane
   prázdné. Zákazník musí texturu nastavit až po importu (postup ověřen, je v návodu).
   Rozhodnutí 2026-11-17: sf3d kusy se nestahují do draftu, vada se jen zapisuje.
3. **Trellis: odkaz na texturu vede na cestu build stroje** (`/data/converted/…/model_tex.png`).
   Importér hlásí chybu, ale texturu jde vybrat ručně přímo v okně importu (viz tabulka).
   Kdyby FBX odkazoval relativně na `{slug}_tex.png`, načetla by se nejspíš sama (neověřeno;
   oprava v UGC Factory).
4. **Trellis: `bbox_studs` nesedí s FBX.** Helma má v metadatech [1.87, 2.5, 1.87], ve Studiu
   1.86 × 1.81 × 1.87. Zákazníkovi nevadí, ale filtr podle velikosti v katalogu lže.
5. ZIP se rozbalí do další složky se stejným jménem (`slug/slug/…`). Drobnost do návodu.

## Co má udělat Olin na M2

1. Spustit Roblox Studio, nechat ho aktualizovat a přihlásit se (heslo agent nezadává).
2. Přepnout rozhraní do angličtiny a založit prázdné místo **Baseplate**, uložit lokálně.
3. Zapnout MCP server Studia: **Assistant ⟩ … ⟩ Manage MCP Servers ⟩ Enable Studio as MCP
   server**, pak **Quick connect ⟩ Claude Code**. Session Claude Code pak restartovat, aby
   nástroje Studia viděla.
4. Být asi půl hodiny u M2: dialog Import a výběr souboru agent neodklikne, kliká Olin podle
   seznamu kroků a agent fotí a ověřuje.
