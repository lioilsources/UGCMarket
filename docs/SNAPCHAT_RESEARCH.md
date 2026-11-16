# Snapchat – rešerše: doplňky pro avatara, generátor a obchod

> **Toto není právní ani daňová rada.** Rešerše veřejných zdrojů k 2026-11-16.
> Většinu stránek jsem četl přes automatický výtah nebo úryvky z vyhledávání, ne celé v originále.
> Před rozhodnutím, které stojí peníze nebo čas, ověř citace v originále.
>
> Legenda: **[P]** = primární (dokumentace, podmínky a nápověda Snapu),
> **[S]** = sekundární (média, blogy, marketing partnerů). Vše navštíveno **2026-11-16**.

## Shrnutí

Snapchat **není** otevřená ekonomika doplňků pro avatara jako Roblox. Oblečení pro Bitmoji prodává
Snap sám a značky, se kterými má smlouvu. Nezávislý tvůrce tam vlastní kusy nahrát ani prodat nemůže.
Otevřená část Snapchatu jsou **Lenses** (AR efekty z Lens Studia). Naše modely se do nich technicky
hodí, ale výplaty za Lenses nejsou dostupné tvůrcům z Česka.

| | Roblox | Snapchat |
|---|---|---|
| Kdo smí prodávat doplňky pro avatara | kdokoli ověřený (Marketplace) | Snap a smluvní značky |
| Otevřená cesta pro tvůrce | Marketplace, Creator Store, vlastní hry | Lenses (AR), Lens Games |
| Výdělek tvůrce | prodej položek | výplaty za úspěšné Lenses, jen ve 42 zemích bez ČR |
| Co kupující potřebuje od nás | FBX do 4000 trojúhelníků | model do Lens do 8 MB, do 100 000 trojúhelníků |
| Vlastní AI generátor platformy | ano: text → 3D model ve Studiu | ano: text → 3D, text → Lens, rekvizity pro Bitmoji |

## Zdroje

| ID | Zdroj | Typ |
|---|---|---|
| BF | Exclusive Snapchat Bitmoji Fashion, https://support.bitmoji.com/hc/en-us/articles/20278438888084-Exclusive-Snapchat-Bitmoji-Fashion | [P] |
| SCP | How do I get exclusive Bitmoji fashion on Snapchat?, https://help.snapchat.com/hc/en-us/articles/40679147635348-How-do-I-get-exclusive-Bitmoji-fashion-on-Snapchat | [P] |
| SUITE | Bitmoji Suite, https://developers.snap.com/lens-studio/features/bitmoji-suite/overview a `/outfit` | [P] |
| B3D | Bitmoji 3D, https://developers.snap.com/lens-studio/features/bitmoji-avatar/bitmoji-3d | [P] |
| BAV | Bitmoji Avatar overview, https://developers.snap.com/lens-studio/features/bitmoji-avatar/overview | [P] |
| SUBM | Lens submission guidelines, https://developers.snap.com/lens-studio/publishing/submitting/submission-guidelines | [P] |
| LCR | Lens Creator Rewards, https://developers.snap.com/lens-studio/monetization/lens-creator-rewards | [P] |
| LPT | Snap Lens+ Payout Program Terms (účinnost 2026-01-09), https://www.snap.com/terms/lens-payout | [P] |
| LSLA | Lens Studio License Agreement (účinnost 2024-06-17), https://www.snap.com/terms/lens-studio-license-agreement | [P] |
| LST | Lens Studio Terms (účinnost 2024-09-16), https://www.snap.com/terms/lens-studio-terms | [P] |
| ALIB | Asset Library publishing guide, https://developers.snap.com/lens-studio/assets-pipeline/asset-library/asset-library-publishing-guide | [P] |
| DG | Digital Goods for Lenses (dokumentace Lens Studio 4.55.1), https://developers.snap.com/lens-studio/4.55.1/references/guides/lens-features/digital-goods/digital-goods-for-lenses | [P], stará verze |
| MCP | Developer mode with Claude Code, https://developers.snap.com/lens-studio/features/lens-studio-ai/developer-mode/with-claude-code | [P] |
| T3D | Snap AR blog o text-to-3D (2026-05-04), https://ar.snap.com/blog/text-to-3d-ai-ar-platform | [P] |
| EASY | Easy Lens, https://developers.snap.com/lens-studio/features/lens-studio-ai/easy-lens/overview | [P] |
| HEAD | Head attached 3D objects, https://developers.snap.com/lens-studio/features/ar-tracking/face/head-attached-3d-objects | [P] |
| BODY | Object Tracking 3D, https://developers.snap.com/lens-studio/features/ar-tracking/body/object-tracking-3d | [P] |
| BRAND | Snap brand guidelines, https://www.snap.com/brand-guidelines (pravidla jsou v PDF ze září 2025, **nečteno**) | [P] |
| DRESSX | DRESSX o Bitmoji Avatar Store a kolekci Diesel, https://dressx.com/news/what-is-snapchat-s-bitmoji-avatar-store-and-how-can-brands-benefit-from-it | [S], marketing |
| SMT | Social Media Today o Bitmoji Suite, https://www.socialmediatoday.com/news/snapchat-adds-gaming-bitmoji-creation-elements-lens-studio/749737/ | [S] |
| CK | Engadget o Commerce Kit, https://www.engadget.com/social-media/snap-is-bringing-shopping-features-to-its-ar-glasses-170000985.html | [S] |
| NI | Net Influencer o konci aplikace Bitmoji, https://www.netinfluencer.com/snap-to-shut-down-standalone-bitmoji-app-consolidate-features-into-snapchat/ | [S] |

---

## 1. Oblečení pro Bitmoji: smí ho vyrábět a prodávat nezávislý tvůrce?

**Verdikt: NE** (stav k dnešku; program pro nezávislé tvůrce jsem nenašel).

- Uživatel si oblečení pro Bitmoji bere v editoru přímo ve Snapchatu. Část kusů se odemyká za
  Snap Tokens, část je vyhrazená předplatitelům Snapchat+ [P BF, SCP]. Kusy získané za tokeny
  fungují jen uvnitř Snapchatu [P BF].
- Placené kolekce dodávají značky ve spolupráci se Snapem. DRESSX se uvádí jako technický partner,
  který pro značky digitální kusy vyrábí (např. Diesel, podzim 2025) [S DRESSX].
- Samostatná aplikace Bitmoji skončila 30. 6. 2026, vše je teď ve Snapchatu [S NI].
- Veřejnou specifikaci oblečení pro Bitmoji (kostra, formát, limity) ani přihlášku pro tvůrce
  jsem nenašel. Cesta dovnitř vede přes obchodní dohodu se Snapem nebo přes partnera jako DRESSX.

## 2. Bitmoji v Lens Studiu

**Verdikt: jen uvnitř Lens a jen nekomerčně.**

- Bitmoji Suite umí tři věci: obléknout Bitmoji do kusů **z knihovny Snapu**, přidat animaci
  a **vygenerovat rekvizity z textu** ve stylu Bitmoji [P SUITE]. Vlastní model oblečení se
  podle dokumentace nahrát nedá.
- Komponenta Bitmoji 3D vloží do scény avatara uživatele nebo kamaráda [P B3D]. Rekvizity jde
  k avatarovi připojit (stránka „Adding Bitmoji Props“ [P BAV]; její obsah se mi nepodařilo načíst).
- Dokumentace výslovně říká: *„Bitmoji avatars cannot be used for commercial purposes, including
  Ads.“* [P B3D]. Jestli se za komerční účel počítá i Lens přihlášený do výplat, není jasné.

## 3. Naše doplňky v Lens (AR na uživateli)

**Verdikt: ANO, technicky to sedí dobře.**

- Lens Studio umí připnout 3D objekt na hlavu (Head Binding: střed hlavy, čelo, brada a další
  body) [P HEAD] a na klouby těla nebo ruky (Object Tracking 3D, 3D Hand Tracking) [P BODY].
  Helma, klobouk a vlasy patří na hlavu, batoh a křídla na trup, meč do ruky.
- Limity při odeslání Lens [P SUBM]: zabalená Lens nejvýš **8 MB**, textury nejvýš 2048 × 2048,
  mesh nejvýš **100 000 trojúhelníků** (60 000 s kostrou), nejvýš 150 MB paměti.
- Náš FBX (do 4000 trojúhelníků, textura 1024 px, zhruba 0,8 MB) se vejde s rezervou.
  Surový GLB (kolem 11 MB) se nevejde. Pro AR by se hodil střední export, viz „Dopad na plán“.
- Lens Studio na M2 zatím není nainstalované. Nic z tohoto oddílu jsem nezkoušel.

## 4. Výdělek z Lenses

**Verdikt: NE pro tvůrce z Česka.**

- Program Lens Creator Rewards má dvě větve [P LCR]: **Top Performer Payouts** (odměna pro
  Lenses, které se do 90 dnů od odeslání dostanou mezi nejúspěšnější) a **Lens+ Payouts**
  (od 16. 10. 2025; podíl z předplatného za Lenses označené jako exkluzivní, jen pro schválené
  tvůrce).
- Podmínka: věk 18+ a sídlo v jedné ze **42 zemí**. V seznamu je Německo, Rakousko nebo Itálie,
  **Česko v něm není** [P LCR]. Podmínky Lens+ navíc chtějí, aby tvůrce Lens odeslal fyzicky
  z oprávněné země [P LPT].
- Lens+ je jen na přihlášku, výplata od 100 USD, vzorec je neveřejný a Snap může program kdykoli
  změnit nebo ukončit [P LPT].
- Prodej digitálního zboží uvnitř Lens za tokeny dokumentace označuje za ukončený [P DG, stará
  verze]. Nový Commerce Kit je jen pro vybrané vývojáře a míří na brýle Specs [S CK].
- Sekundární zdroje uvádějí zhruba 1 milion USD měsíčně rozdělený mezi víc než tisíc tvůrců.
  To je v průměru nízké stovky dolarů na tvůrce, a neověřené.

## 5. Prodej modelů pro Lenses na Ol1n Marketu

**Verdikt: právně spíš ANO, poptávka nejasná.**

- Podmínky Lens Studia neřeší, odkud tvůrce modely má. Tvůrce ale Snapu ručí, že má k obsahu
  všechna práva, a dává Snapu širokou sublicencovatelnou licenci k odeslané Lens [P LST].
  Naše licence by tedy musela výslovně dovolit použití v Lens a udělení této licence Snapu,
  stejně jako už to dělá u Robloxu.
- Americká verze licence Lens Studia říká, že nástroj je určený jen pro původní díla, která
  uživatel sám vytvořil a vlastní [P LSLA]. Ve verzi pro zbytek světa tato věta není. Pro
  zákazníky z USA je to nejasné místo.
- Placené tržiště modelů pro Lens Studio jsem nenašel, oficiální ani neoficiální. Asset Library
  v Lens Studiu je zdarma; tvůrce do ní může kus poslat ke schválení a zveřejní se pod jeho
  účtem [P ALIB].
- Tvůrci Lenses mají zdarma knihovnu modelů a generátor z textu přímo v nástroji [P T3D].
  Důvod platit za hotový model je tedy slabší než u Robloxu.

## 6. Generátor: co už nabízí Snap sám

**Verdikt: obecný generátor nemá proti Snapu výhodu.**

- Lens Studio generuje 3D model z textu a textury na existující model [P T3D].
- Easy Lens postaví celou Lens z textového zadání, jako aplikace, na webu i uvnitř Lens Studia [P EASY].
- Bitmoji Suite generuje rekvizity ve stylu Bitmoji [P SUITE].
- Naše výhoda může být jen v tom, co Snap nedává: ucelené kolekce v jednom stylu, kusy už
  usazené na hlavu nebo tělo a stejný kus pro víc platforem najednou.

## 7. Automatizace

**Verdikt: částečně ANO, bez plně bezobslužného režimu.**

- Lens Studio má vlastní MCP server a oficiální návod na připojení Claude Code. Agent umí číst
  a měnit scénu, přidávat objekty a komponenty, psát skripty a instalovat kusy z Asset Library [P MCP].
- Server patří k jednomu oknu Lens Studia a nástroj musí běžet. Režim bez okna ani příkaz pro
  sestavení z příkazové řádky jsem nenašel.
- Jestli MCP umí importovat lokální soubor a odeslat Lens, dokumentace neříká. Odeslání Lens
  je zveřejnění a dělá se přes účet tvůrce.

## 8. Značka a formulace

- V Lens nesmí být odkazy na „Snap“, „Snapchat“, „Spectacles“, „Bitmoji“, logo ducha ani
  rozhraní Snapchatu, a nic, co naznačuje partnerství se Snapem [P SUBM].
- Zakázané jsou QR kódy, výzvy typu „sleduj“ a „sdílej“ a pokusy dostat uživatele na cizí web [P SUBM].
  Lens tedy nemůže vést zákazníky na náš obchod.
- Licence Lens Studia zakazuje používat jméno a loga Snapu bez písemného souhlasu [P LSLA].
- Pravidla pro popisné užití („compatible with Lens Studio“) jsou v PDF, které jsem nečetl [P BRAND].

---

## Dopad na plán

1. **Oficiální oblečení pro Bitmoji nedělat.** Je to uzavřený obchod. Jediný krok je dotaz na
   Snap nebo na partnera typu DRESSX. Olin 16. 11. rozhodl dotaz poslat, koncept je
   v `docs/SNAPCHAT_INQUIRY.md`.
2. **Na výplaty z Lenses nespoléhat.** Z Česka nejsou dostupné a obcházet podmínku sídla nejde.
3. **Lenses dávají smysl jako výloha a měření zájmu,** ne jako příjem: bezplatná Lens s naší
   helmou ukáže, které kusy lidi baví. Zákazníky do obchodu ale přivést nesmí.
4. **Prodej modelů pro Lenses je levné rozšíření obchodu,** protože modely už existují. Poptávka
   je ale neověřená, takže nejdřív malý pokus, ne nová produktová řada.
5. **Do UGC Factory by přibyl druhý profil exportu** vedle `nas/spec/roblox_spec.json`: víc
   trojúhelníků, metry místo studs, počátek v bodě uchycení. Je to jiné repo, jen návrh.
6. **Licence:** doplnit povolení použít model v Lens a udělit Snapu licenci podle jeho podmínek.
   Texty licencí jsou pořád DRAFT a čekají na právní kontrolu.

## Nejistoty

- **Snap Tokens:** nápověda Bitmoji s nimi dál počítá, stará dokumentace Lens Studia mluví
  o jejich ukončení. Jak se oblečení pro Bitmoji v roce 2026 platí, jsem jednoznačně neověřil.
- **Program pro tvůrce oblečení Bitmoji** mohl vzniknout a já ho nenašel. Lens Fest 2026 jsem
  ve výsledcích neviděl.
- **Seznam 42 zemí** je z výtahu stránky, ne z originálu. Ověřit přímo před rozhodnutím.
- **Rekvizity pro Bitmoji:** technické požadavky a body uchycení neznám, stránka se nenačetla.
- **Body uchycení na těle a ruce** znám jen z úryvků, nezkoušel jsem je.
- **Pravidla značky** pro popisné užití názvů Snapu jsem nečetl (PDF).
- **Podmínky generativních nástrojů Snapu** jsem nečetl.
- **Poptávka** po placených modelech pro Lenses: žádný důkaz pro ani proti.
- **Daně a platby** (výplaty ze zahraničí, DPH) jsou mimo rozsah.
