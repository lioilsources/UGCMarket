# Roblox compliance – RobloxMarket (market.ol1n.com)

> **Toto není právní ani daňová rada.** Jde o rešerši veřejných pravidel Roblox k 2026-10-25.
> Pravidla se mění (ToU má novou verzi s **Effective Date 1. 11. 2026**, poslední úprava stránky 2026-10-07).
> Před spuštěním placeného prodeje doporučuji krátkou konzultaci s právníkem (IP / spotřebitelské právo EU).
>
> Legenda zdrojů: **[P]** = primární (Roblox ToU, Community Standards, Creator Hub docs, help center),
> **[S]** = sekundární (DevForum příspěvky uživatelů, blogy, média). Všechny odkazy navštíveny **2026-10-25**.

## Zdroje

| ID | Zdroj | Typ |
|---|---|---|
| ToU | Roblox Terms of Use (User + Creator Terms), https://en.help.roblox.com/hc/en-us/articles/115004647846-Roblox-Terms-of-Use (stránka upd. 2026-10-07, Effective 2026-11-01) | [P] |
| CS | Roblox Community Standards, https://about.roblox.com/community-standards | [P] |
| LOGO | Roblox Name and Logo – Community Usage Guidelines, https://en.help.roblox.com/hc/en-us/articles/115001708126 (upd. 2026-10-07) | [P] |
| SELL | Sell on the Creator Store, https://create.roblox.com/docs/production/sell-on-creator-store | [P] |
| CSTORE | Creator Store, https://create.roblox.com/docs/production/creator-store | [P] |
| ASSETS | Open Cloud – Usage guide for assets, https://create.roblox.com/docs/cloud/guides/usage-assets | [P] |
| OAREG | OAuth 2.0 app registration and review, https://create.roblox.com/docs/cloud/auth/oauth2-registration | [P] |
| OADEV | OAuth 2.0 app implementation, https://create.roblox.com/docs/cloud/auth/oauth2-develop | [P] |
| TPAP | Creator Third-Party App Policy, https://en.help.roblox.com/hc/en-us/articles/37924211313044 (upd. 2026-10-08) | [P] |
| RL | Open Cloud rate limits, https://create.roblox.com/docs/cloud/reference/rate-limits | [P] |
| GENAI | Generative AI in experiences, https://create.roblox.com/docs/generative-AI | [P] |
| AIDS | AI data sharing, https://create.roblox.com/docs/ai-data-sharing | [P] |
| DF-OAUTH | DevForum oznámení OAuth 2.0 (2023-06-01), https://devforum.roblox.com/t/2401354 | [P/S] – oficiální post Roblox staff, ale starší |
| DF-MOD | DevForum vlákna o chybné moderaci, např. https://devforum.roblox.com/t/account-banned-for-sexual-content-for-uploading-a-cat-texture/4101908 | [S] |

---

## 1. Prodej vlastních 3D modelů mimo Roblox za fiat/krypto

**Verdikt: ANO** (pro obecné soubory GLB/FBX/OBJ s vlastní IP) – s výhradou názvu a brandingu (viz níže: **NE** pro „Roblox“ v názvu obchodu).

**Co je zakázáno (ToU/CS):**
- ToU (User Terms, Robux & Virtual Content): *„You may not use, acquire, or distribute Robux or Virtual Content except through the Services …“*; Roblox neuznává third-party služby, které prodávají Robux nebo Virtual Content. **Virtual Content** = položky získané *na* platformě (avatar items, game passes, přístup do experience) [P ToU].
- Prodej/nákup účtů a přístupových údajů je zakázán [P ToU, Account Responsibilities].
- CS → Roblox Economy > Fraud: zakázané jsou *„off-platform, secondary market transactions for account access or virtual content from games or Marketplace“* a *„Using third-party services to buy, sell, trade, or give away Robux“* [P CS].
- CS → Promotional Offers: *„Developers may not independently use off-platform services or products to sell on-platform items“* bez smlouvy s Roblox [P CS].
- CS → Directing Users Off-Platform: na platformě (v experience, popisech assetů, chatu) nesmíš odkazovat na externí weby kromě funkce Social Links [P CS]. → **Odkaz na market.ol1n.com nesmí být v popisu assetu / v experience.**

**Proč je náš model OK:** prodáváme **soubory a licenci** k vlastnímu dílu, ne Robux, účty ani on-platform položky. ToU výslovně říká, že Creator si ponechává autorská práva k UGC (Creator Terms, *Ownership of UGC*) [P ToU]. Nic v ToU nezakazuje prodávat vlastní 3D díla mimo Roblox.

**Marketing / trademark – zde je reálné riziko:**
- LOGO: *„The use of the Roblox name … badge, logo … may not be tied to the development, distribution or sale of any products or services … without a Roblox licensing agreement.“* Logo a „Tilt“ nejsou povoleny vůbec (kromě badge „Now on Roblox“ pro propagaci experience). Nesmí se naznačovat partnerství/endorsement [P LOGO].
- ToU Creator Terms: zákaz „Roblox“, „Blox“ a podobných jmen v názvech experiences, username, display names, community names; Roblox Trademark License platí jen *uvnitř* experiences [P ToU]. Creator nesmí používat své UGC tak, aby naznačoval afilaci s Roblox [P ToU, *Limitations on Use of UGC*].
- **Závěr:** název „RobloxMarket“ a jakékoli logo Roblox na webu jsou v rozporu s LOGO guidelines. Faktické zmínky („modely optimalizované pro Roblox Studio“) jsou běžné a obecně chráněné nominativním užitím, ale Roblox guideline je striktnější → **nejasné**, držet minimum.

## 2. Doručování `.rbxm` mimo Roblox; licence k obsahu ze Studia

**Verdikt: nejasné, spíš ANO (OK)** pro .rbxm obsahující výhradně vlastní geometrii/textury; doporučení **GLB/FBX jako default, .rbxm jen opt-in**.

- Formát .rbxm sám o sobě ToU nezakazuje distribuovat; Roblox ho sám přijímá jako vstup Open Cloud uploadu (Model: `.fbx, .gltf, .glb, .rbxm, .rbxmx`) [P ASSETS]. Pozor: *„.rbxm or .rbxmx files edited outside of Roblox Studio might not upload or function“* (u Animation) [P ASSETS].
- Studio License: Studio smíš používat *„solely for the purpose of creating … and Publishing Creator’s UGC on the Services“* [P ToU, Creator Terms – Roblox Studio]. Export .rbxm a jeho prodej jako souboru je **šedá zóna** – výslovně zakázáno není, ale účel licence je publikace na Roblox. Kupující .rbxm stejně použije jen ve Studiu, tj. v souladu s účelem.
- **Templates a Roblox IP:** Templates a prvky Studia patří Roblox, nesmí tvořit součást prodávaného souboru [P ToU]. Licensed Music se nesmí exportovat mimo Services [P ToU]. → .rbxm nesmí obsahovat Roblox templates, toolbox assety cizích autorů, licencovanou hudbu, ani skripty závislé na `require(assetId)`/`LoadAssetAsync`.
- **Licence pro Roblox:** obsah *uložený* lokálně ve Studiu Roblox licenci nedává; licence vzniká **publikací** na Services. UGC License: Creator dává Roblox *„perpetual, worldwide, non-exclusive, royalty-free right and license (with the right to sublicense …)“* k hostování, úpravám, distribuci a *„training, development, and use of machine learning … including inputs and outputs within Studio“* [P ToU]. Platí pro cokoli, co Olin **nebo kupující** nahraje na Roblox. Exkluzivní právo Roblox se týká aplikace UGC na Classic Avatar [P ToU].
- Důsledek pro LICENSE: kupující při uploadu na Roblox zaručuje, že *„own[s] or otherwise control[s] all the rights“* a uděluje Roblox výše uvedenou sublicencovatelnou licenci [P ToU]. **Naše licence proto musí kupujícímu výslovně dovolit nahrát model na Roblox a udělit Roblox licenci vyžadovanou jeho ToU.**

## 3. Creator Store jako druhý kanál

**Verdikt: ANO, je to možné; exkluzivita nebyla nalezena (nejasné, ale nic ji nevyžaduje).**

Aktuální stav [P SELL, CSTORE]:
- **Ceny jen v USD** (*„You can only price assets on the Creator Store in USD“*); Robux pricing už neexistuje.
- Prodat lze jen **Models a Plugins**. Models **$2.99–$49.99**, Plugins $4.99–$249.99. Free distribuce: Model, Plugin, MeshPart, Decal, Audio.
- Poplatky: *„100% of net proceeds“* – strhávají se jen daně a payment processing fees. Peníze jdou do **30denního escrow**, výplata měsíčně (minima dle země).
- **Seller account přes Stripe** (jeden Stripe účet na Roblox účet, země nelze změnit, nedostupné např. BR/CN/IN/RU – ČR není na seznamu výjimek, ale ověřit při onboardingu). 18+ (nebo 13–17 se souhlasem rodiče), **2-Step verification**. Jen assety vlastněné **individuálním účtem** (ne group).
- Distribuce vyžaduje účet starší 2 dnů, ověřený **age check nebo government ID** (telefon nestačí). Limity / 30 dní: verified 200 models/meshes/images, unverified 10.
- Refund/chargeback i moderace již prodaného assetu (bez úspěšného odvolání do 5 dnů) → částka se strhne prodejci.
- AI data sharing: placené assety po 2024-07-10 mají sdílení pro AI trénink **zapnuté defaultně** (lze vypnout), **free assety vždy sdílené** [P AIDS, CSTORE].
- Exkluzivita: SELL ani CSTORE ji nezmiňují; UGC License je **non-exclusive** [P ToU]. Prodej stejného modelu na vlastním webu je tedy slučitelný.
- Pozor: v popisu Creator Store assetu **nesmí být odkaz na vlastní web** (CS – Directing Users Off-Platform) [P CS].

## 4. Open Cloud Assets API přes OAuth 2.0 (upload do účtu kupujícího)

**Verdikt: ANO technicky (s jeho souhlasem), ale s review a betou.**

- Scopes: `asset:read` + `asset:write` (*„Create and update assets on the user’s behalf“*), identita `openid:read` (+ `profile:read`) [P ASSETS, TPAP, DF-OAUTH]. Asset patří uživateli v `creationContext` požadavku [P ASSETS].
- Kategorie aplikace: **Creation & Productivity Tools** (povoluje `asset:read`, `asset:write`); všechny scopes musí spadat do jedné kategorie, jen minimum nutné, *„Roblox does not approve scopes for planned or future functionality“* [P TPAP, OAREG].
- Vývojář musí být **ID verified** (pro uživatele aplikace to neplatí, stačí 13+) [P/S DF-OAUTH, OAuth overview].
- **Private mode = max 10 uživatelů**; pro veřejné nasazení review: popis, thumbnail, Entry Link, HTTPS URL na Privacy Policy a ToS, **demo video ≤ 1 min** celého flow, zdůvodnění každého scope. Po schválení nelze zpět do private; každá změna public app = nové review [P OAREG].
- Tokeny: access token 15 min, refresh token uložit bezpečně, PKCE doporučeno [P OADEV].
- Upload: 1 asset / volání, max **20 MB**; Model z `.fbx/.gltf/.glb/.rbxm/.rbxmx` se importuje jako Model s MeshParts (balíček v Inventory > My Packages); update obsahu jen pro .fbx. Assets API je **Beta** [P ASSETS].
- Rate limits: OAuth limity jsou **per access token** [P RL]; konkrétní číslo pro `POST /assets/v1/assets` (vyhledávání uvádí 120/min) jsem v primárním zdroji neověřil [S].
- Moderace: výsledek v `moderationResult` operace (`MODERATION_STATE_APPROVED` …) → polling přes `GET /assets/v1/operations/{id}` [P ASSETS].
- Data policy: nesmí se prodávat data z Roblox API ani je použít k tréninku AI; při ztrátě přístupu data smazat [P TPAP].
- **Riziko:** CS zakazuje *„use off-platform services … to sell on-platform items“* [P CS]. Placený „upload do Roblox účtu“ by mohl být vnímán jako prodej on-platform položky. **Doporučení:** prodávat jen soubor+licenci, upload nabídnout jako bezplatný nástroj pro vlastníka licence (stejně by mohl soubor nahrát ručně) – **nejasné**, review může dopadnout různě.

## 5. AI-generovaný obsah

**Verdikt: nejasné, spíš ANO (bez zvláštního omezení).**

- V ToU, CS ani Creator Store docs jsem **nenašel povinnost označovat** assety vytvořené generativní AI ani zákaz takových assetů. ToU sekce „AI Features“ se týká jen AI nástrojů *Roblox* (uživatel si ponechává práva k Prompts/Outputs) [P ToU]. GENAI vyžaduje transparentnost u generativní AI *běžící v experiences*, ne u statických assetů [P GENAI].
- Platí obecné: uploader zaručuje, že vlastní/kontroluje práva a obsah neporušuje IP třetích stran; Roblox může proaktivně odstranit podezřele porušující assety [P ToU, CS – IP Violations].
- IP rizika pipeline: (a) licence vah Illustrious/SDXL derivátů a TRELLIS.2 – ověřit, že komerční použití výstupů je povoleno (zejm. některé checkpointy/LoRA mají non-commercial licence); (b) prompty nesmí cílit na známé postavy/značky (Roblox i jiné franšízy); (c) autorskoprávní ochrana čistě AI výstupu je v EU/US nejistá → licence má být smluvní (licenční podmínky), ne jen spoléhat na copyright; (d) ToU zakazuje přistupovat ke Services kvůli použití Virtual Content pro trénink/použití ML [P ToU] → **nepoužívat Roblox assety jako vstupy/reference pipeline.**
- Doporučení: dobrovolně uvádět „AI-assisted (generated + manual cleanup)“ – transparentní vůči zákazníkům a v EU v souladu s duchem AI Act (čl. 50 se primárně týká deepfakes/textu, ne nutně 3D assetů – **nejasné**).

## 6. Moderace a reklamace

**Verdikt: ANO, každý upload prochází moderací; výsledek nelze předem garantovat.**

- Creator Store: human + automated moderace proaktivně i reaktivně; porušení → moderace assetu **a účtu** [P CSTORE]. Odvolání možné; DevForum hlásí false positives (např. UV-unwrapped textury označené jako sexual content) a rychlá zamítnutí odvolání [S DF-MOD].
- Open Cloud vrací stav moderace v operaci; zamítnutý asset se nenačte [P ASSETS].
- Opakovaně zamítané uploady mohou vést k postihu účtu kupujícího → zákazník musí vědět, že riziko uploadu nese on, a my musíme pre-screenovat (žádné texty, loga, nahota, zbraně v realistickém kontextu, extremistické symboly).

**Návrh textu do LICENSE/README (EN, pro zákazníky):**

```
ROBLOX DISCLAIMER
This product is an independent 3D asset (GLB/FBX/OBJ, optional RBXM) created by
[Olin / OSVČ name, IČO]. It is not affiliated with, endorsed by, or sponsored by
Roblox Corporation. "Roblox" is a trademark of Roblox Corporation, used here only
to describe compatibility.

You may upload this asset to your own Roblox account and use it in your
experiences. By uploading, you grant Roblox the license required by the Roblox
Terms of Use; this license permits that sublicense.

Every Roblox upload is subject to Roblox moderation, which we do not control.
We pre-screen all assets against the Roblox Community Standards, but we cannot
guarantee approval. Do not upload content to an account you cannot risk.

This asset was created with AI-assisted tools (image/3D generation) followed by
manual cleanup.

REFUNDS
If Roblox moderation rejects an unmodified asset within 14 days of purchase,
send us the asset/operation ID and the rejection notice and we will, at our
choice, provide a fixed version or a full refund. Modified assets, rejections
caused by your own changes, and account-level actions by Roblox are excluded.
Statutory consumer rights (including EU withdrawal rights where applicable)
are not affected.
```

Pozn.: u digitálního obsahu v EU zaniká právo na odstoupení jen při výslovném souhlasu spotřebitele se zahájením plnění a potvrzení ztráty práva (checkout checkbox) – **ověřit s právníkem**.

---

## Dopad na plán (MVP)

1. **Přejmenovat projekt / doménu v UI.** „RobloxMarket“ je v rozporu s LOGO guidelines (jméno Roblox spojené s prodejem produktů). Interní kódové jméno může zůstat, veřejně např. „ol1n Market“, „Blockforge Assets“ (bez „Blox“!). Žádné Roblox logo, Tilt, screenshoty UI Roblox s logem, ani „Now on Roblox“ badge.
2. **Formulace k použití:** „3D modely kompatibilní s Roblox Studio“, „optimalizováno pro import do Roblox Studia (3D Importer)“, disclaimer o neafilaci v patičce.
   **Vyhnout se:** „Roblox items“, „Roblox store/shop“, „buy Roblox assets“, „official“, „UGC items for your avatar“, cokoli o Robux („cheaper than Robux“, platba Robuxy, gift cards), „limiteds“, prodej účtů, „we add the item to your Roblox inventory“.
3. **Formáty:** default **GLB + FBX** (+ OBJ volitelně), textury PNG; **.rbxm jen opt-in** a jen s čistě vlastní geometrií (žádné templates, cizí toolbox assety, skripty). Držet se limitů importu (soubor ≤ 20 MB pro Open Cloud).
4. **Nepřijímat Robux** ani nic, co by připomínalo směnu Robux ↔ fiat/krypto. Krypto platby za soubory jsou z pohledu Roblox irelevantní (daňově/AML řešit zvlášť).
5. **Creator Store jako 2. kanál dává smysl** pro discoverability (Toolbox) a důvěru: USD ceny $2.99–$49.99, 100 % net, Stripe. Kroky: ID verifikace, 2SV, Stripe onboarding (ČR), asset na individuálním účtu, **vypnout AI data sharing** u placených assetů (free jsou sdílené vždy). V popisech Creator Store **žádné odkazy** na web. Free „teaser“ modely na Creator Store = AI trénink Roblox nevyhnutelný.
6. **OAuth upload až v2:** registrace app v kategorii Creation & Productivity, jen `openid:read profile:read asset:read asset:write`; MVP běží v private mode (≤10 testerů), pro public připravit privacy policy, ToS, demo video. Upload nabídnout zdarma po nákupu licence, ne jako placenou službu „doručení do Roblox“.
7. **Pre-moderace v pipeline:** automatická kontrola textur (NSFW klasifikátor, OCR na text/loga), ručně zkontrolovat před publikací; logovat prompty a verze modelů (doklad o původu pro IP spory).
8. **LICENSE:** výslovně povolit upload na Roblox a sublicenci dle Roblox ToU; zakázat přeprodej samotných souborů; uvést AI-assisted původ; reklamační postup výše.

## Nejistoty

- **Plné znění ToU** – stránka vrací WebFetch 403, text získán přes Zendesk Help Center API (stejný obsah, ale zobrazená verze je ta s Effective Date 2026-11-01; předchozí verze neporovnána).
- **Studio License vs. export .rbxm** k prodeji mimo platformu – výslovně neřešeno; jen výklad.
- **Placený OAuth upload** vs. CS „off-platform sale of on-platform items“ – výklad nejasný, závisí na review Roblox.
- **Rate limit** `POST /assets/v1/assets` pro OAuth – konkrétní číslo (120/min) jen ze sekundárního výpisu vyhledávání; primární stránka RL čísla neuvádí.
- **Exkluzivita Creator Store** – nenalezena ani výslovně vyloučena; plné „Creator Store seller terms“ / Stripe Connected Account Agreement jsem nečetl.
- **Dostupnost Stripe seller onboardingu v ČR** – nevylučuje se, ale neověřeno.
- **AI labeling** – žádné pravidlo nenalezeno; může přibýt (Roblox často mění AI politiku). EU AI Act čl. 50 vs. 3D assety nejasné.
- **Licence modelů pipeline** (Illustrious checkpoint, LoRA, TRELLIS.2) pro komerční výstupy – mimo rozsah, nutno ověřit zvlášť.
- **Moderace a odvolání** – fakta o průběhu odvolání jen ze sekundárních zdrojů (DevForum).
- Daňové a spotřebitelské aspekty (DPH OSS u digitálních služeb do EU, krypto platby u OSVČ) – mimo rozsah, nejde o radu.
