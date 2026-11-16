# 03 — Snapchat: generátor a obchod s doplňky

Podklad: `docs/SNAPCHAT_RESEARCH.md` (stav 16. 11. 2026). Přečti ho celý, než začneš.
Navazuje na plány 01 a 02; nic z nich neruší.

## Závěr rešerše ve třech větách

Oblečení pro Bitmoji je uzavřený obchod Snapu a smluvních značek, nezávislý tvůrce tam nic
nenahraje ani neprodá. Otevřená část Snapchatu jsou Lenses (AR efekty); naše modely se do nich
technicky hodí, ale výplaty za Lenses nejsou pro tvůrce z Česka. Zbývá levné rozšíření obchodu
o modely připravené pro Lens Studio a Lenses jako výloha, obojí s neověřenou poptávkou.

## Doporučení

Nestavět novou produktovou řadu. Udělat **časově omezený pokus** (etapy 0 a 1, dohromady
dva až tři dny práce) a podle výsledku se rozhodnout. Pokus nezačínat dřív, než poběží placený
prodej na marketu (plán 01, etapa D na produkci a rozhodnutí R5) a návod (plán 02). Důvod:
Snapchat nepřinese příjem dřív než hotový prodej, který už je skoro hotový.

**Olin 16. 11. 2026: pokus se udělá** (S1). K pořadí se nevyjádřil, proto platí doporučení
výše, dokud neřekne jinak. Etapa 0 Olina nepotřebuje a může běžet kdykoli.

## Co nedělat

- **Oficiální oblečení pro Bitmoji.** Není jak. Jediný krok je dotaz na partnera a na Snap:
  koncept je v `docs/SNAPCHAT_INQUIRY.md`, odesílá Olin (S3).
- **Stavět na výplatách z Lenses.** Česko není mezi 42 zeměmi. Podmínku sídla neobcházet.
- **Lens, která vede lidi do obchodu.** Pravidla zakazují odkazy, QR kódy a výzvy k akci.
- **Slova „Snapchat“, „Snap“, „Bitmoji“ a logo ducha** v Lens, v názvu produktu a v názvu
  obchodu. Popisné užití na webu až po přečtení pravidel značky (etapa 0).
- **Bitmoji v čemkoli komerčním.** Dokumentace to výslovně zakazuje.

## Etapa 0 — dočíst, co rešerše nechala otevřené (půl dne, bez Olina)

- [ ] Pravidla značky Snapu (PDF ze září 2025): smíme psát „compatible with Lens Studio“?
- [ ] Seznam oprávněných zemí pro výplaty přímo v originále, včetně data stránky.
- [ ] Stránka „Adding Bitmoji Props“ v prohlížeči (výtah se nenačetl): body uchycení, požadavky.
- [ ] Jak se v roce 2026 platí oblečení pro Bitmoji (tokeny, předplatné) a jestli vznikl
      program pro tvůrce. Hledej oznámení z Lens Festu 2026 a newsroom Snapu.
- [ ] Podmínky generativních nástrojů Snapu.
- [ ] Poptávka: existují placené modely nebo šablony pro Lens Studio (Gumroad, Etsy, vlastní
      weby tvůrců)? Kolik stojí a prodávají se? Stačí deset nálezů, nebo zjištění, že nejsou.
- [ ] Zapsat odpověď na dotaz z `docs/SNAPCHAT_INQUIRY.md`, pokud už přišla.
- [ ] Opravit `docs/SNAPCHAT_RESEARCH.md` podle zjištění. Když se změní závěr, řekni to Olinovi.

## Etapa 1 — technický pokus na M2 (jeden až dva dny)

Cíl: zjistit, kolik ruční práce stojí jeden kus v Lens a jestli to jde automatizovat.
Nic se nezveřejňuje.

- [ ] Lens Studio na M2 není. Instalaci a přihlášení dělá Olin svým osobním účtem Snapchat
      (S2, viz Kroky Olina). Heslo agent nezadává.
- [ ] Vzít dva kusy z marketu: kabuto helmu (hlava) a box batoh (záda). FBX i GLB.
- [ ] Helma na hlavu přes Head Binding, batoh na trup přes sledování těla. Zapsat: sedí
      měřítko a orientace? Kolik minut ruční úpravy? Jak velká je výsledná Lens proti limitu 8 MB?
- [ ] Vyzkoušet surový GLB: vejde se, když se textura zmenší? Jaký počet trojúhelníků vypadá
      na telefonu dobře a nezpomaluje?
- [ ] Připojit Claude Code k Lens Studiu přes jeho MCP server (po Olinově souhlasu, je to trvalá
      konfigurace). Zjistit, co jde bez klikání: import souboru, usazení, náhled, export.
- [ ] Vyzkoušet generátor z textu přímo v Lens Studiu na stejném zadání („ornate samurai kabuto
      helmet with neon accents“) a porovnat s naším kusem. To je poctivý test, jestli máme co nabídnout.
- [ ] Výstup: `docs/SNAPCHAT_TECH.md` (česky) se screenshoty náhledu a tabulkou času na kus.

## Etapa 2 — rozhodnutí (Olin)

Podle etap 0 a 1 jedna ze tří možností:

- **Stop.** Snapchat odložit, rešerši nechat v `docs/` a vrátit se k ní po změně pravidel.
- **Jen obchod** (etapa 3).
- **Obchod a výloha** (etapy 3 a 4).

## Etapa 3 — modely pro Lens Studio v obchodě (jen po „ano“)

- [ ] **Profil exportu** v UGC Factory vedle `nas/spec/roblox_spec.json`: víc trojúhelníků
      (hodnotu určí etapa 1), metry místo studs, počátek v bodě uchycení, textura do 2048 px.
      Je to jiné repo: nejdřív návrh Olinovi, pak samostatný PR tam.
- [ ] **Market:** soubory produktu podle cíle (Roblox Studio, Lens Studio), filtr v katalogu,
      v ZIPu druhý soubor s návodem. Kus bez exportu pro Lens se pro Lens nenabízí.
- [ ] **Licence:** doplnit povolení použít model v Lens a udělit Snapu licenci, kterou jeho
      podmínky chtějí. Texty jsou DRAFT; bez právní kontroly (R5) se neprodává.
- [ ] **Návod** pro Lens Studio stejným postupem jako plán 02 (projít naostro, screenshoty z M2).
- [ ] **Měření:** nejdřív pět až deset kusů zdarma a počítat stažení podle cíle. Placený
      prodej až když je o ně zájem.

## Etapa 4 — Lenses jako výloha (jen po „ano“)

- [ ] Šablona Lens s jedním místem pro model, usazení podle kategorie.
- [ ] Z vylistovaných kusů připravit první tři Lenses. Bez loga, odkazů a výzev.
- [ ] Odeslání Lens je zveřejnění pod Olinovým účtem: jen na jeho výslovný pokyn pro každou
      dávku. Účet a heslo agent nezadává.
- [ ] Počty použití z portálu My Lenses zapisovat k produktům v adminu jako signál zájmu.
- [ ] Bezplatný kus do Asset Library Lens Studia je možná cesta k viditelnosti; rozhodnout
      až podle etapy 0 (pravidla pro odeslání a licence, kterou tím Snap dostane).

## Rozhodnutí vlastníka

Rozhodnuto 16. 11. 2026:

- **S1: pokus se udělá** (etapy 0 a 1).
- **S2: osobní účet Snapchat.** Lenses z etapy 4 by tedy vyšly pod Olinovým osobním účtem;
  před etapou 4 to ještě jednou potvrdit.
- **S3: dotaz poslat.** Koncept je v `docs/SNAPCHAT_INQUIRY.md`, odesílá Olin.

Otevřené:

- **S4:** smí vzniknout návrh druhého profilu exportu v UGC Factory? Potřeba až pro etapu 3.

## Kroky Olina (pro etapu 1)

1. Odeslat dotaz z `docs/SNAPCHAT_INQUIRY.md` (nejdřív DRESSX) a dát vědět, co odpověděli.
2. Na M2 nainstalovat Lens Studio a přihlásit se osobním účtem Snapchat.
3. V Lens Studiu zapnout MCP server (AI Assistant ⟩ AI Model Context Protocol ⟩ Configure Server)
   a dovolit připojení Claude Code.
4. Mít po ruce telefon se Snapchatem na náhled Lens.

## Pravidla pro práci

- Realita vítězí: co neprošlo Lens Studiem na M2, není ověřené. Pravidla Snapu se mění,
  před každou etapou znovu ověř body, o které se opíráš.
- Nic se neodesílá ke schválení ani nezveřejňuje bez Olinova pokynu. Nic se nekupuje.
- UGC Factory je jiné repo: jen návrhy a samostatný PR po souhlasu.
- Formulace na webu podle rešerše, oddíl 8. Platí i pravidla z `docs/ROBLOX_COMPLIANCE.md`.
- Každá etapa končí zápisem do `CHANGELOG.md`.
