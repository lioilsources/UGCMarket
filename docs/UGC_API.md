# UGC Factory API (`ugc-backend` / `ugc-api`) — reference pro market

> Zdroj: repo UGCFactory (jen čteno, stav k 2026-10-25). Cesty níže jsou
> relativní k tomuto repu, odkazy `soubor:řádek` míří na handler. Shrnutí na konci (§7)
> říká, co plán marketu předpokládá a co v ugc-backendu reálně **není**.

## 1. Přehled: jak služba běží

| Věc | Hodnota | Zdroj |
|---|---|---|
| Proces | Go `ugc-api`, jeden binár, stdlib `net/http` ServeMux | `nas/main.go:44-117` |
| Port | `:8095` (`UGC_ADDR`), publikovaný na hostu JODA `8095:8095` | `nas/main.go:45`, `nas/ugc-stack.yaml:31-32` |
| Compose | projekt `ugcfactory`, soubor `nas/ugc-stack.yaml`; služby `ugc-api`, `ugc-blender` (konverze GLB→FBX), `ugc-fc` (fantasy characters), `cloudflared` | `nas/ugc-stack.yaml:13-80` |
| Síť | compose síť `ugc`; worker volá `http://ugc-api:8095` | `nas/ugc-stack.yaml:42,79` |
| Interní URL | `http://NAS_HOST:8095` (LAN), Spark pushuje na tutéž adresu | `nas/ugc-stack.yaml:8`, `README.md:61`, `spark/internal/ugc/nas.go:20` |
| Veřejná URL | `https://ugc.ol1n.com` přes dedikovaný tunnel `ugc-nas` | `nas/cloudflared/config.yml:6-19` |
| DB | **SQLite** (`modernc.org/sqlite`, WAL, `MaxOpenConns=1`), soubor `/data/ugc.db` | `nas/store.go:54-63`, `nas/main.go:52` |
| Data na disku | `/data` = `${DATA_DIR}` z `.env`, na JODA `/media/storage/ugc` (snap docker vidí jen `/home` a `/media`) | `nas/ugc-stack.yaml:28`, `README.md:219-220` |

**Auth:**
- **Veřejně:** Cloudflare Access vynucený přímo na konektoru tunelu (`access.required: true`, team `ol1n`,
  AUD `1654b518…e074`) — `nas/cloudflared/config.yml:15-19`. Access app `ugc` má policy
  `app-service-token` (jakýkoli platný service token) + `owner-sso` (`nas/README.md:180-183`).
  Strojový klient tedy posílá `CF-Access-Client-Id` + `CF-Access-Client-Secret` (stejný vzor jako
  ugc-api → Spark, `nas/main.go:481-483`).
- **Na LAN / v compose síti: API je úplně otevřené**, žádný token (`nas/main.go:5-6`).
- Jediný API klíč existuje jen pro prefix `/v1/fc/`: `FC_API_KEYS` → `X-API-Key` nebo
  `Authorization: Bearer` (`nas/characters_api.go:39-54`); nenastaveno = otevřeno.
  V `ugc-stack.yaml` se `FC_API_KEYS` ani nepředává (`nas/ugc-stack.yaml:21-26`).
- `/worker/*` nemá auth vůbec. Komentář tvrdí, že je tunel nepublikuje (`main.go:97`), ale ingress
  posílá na `ugc-api:8095` **všechny** cesty (`config.yml:10-11`) — chrání je tedy jen Access.
- CORS: `Access-Control-Allow-Origin: *` jen u `/jobs/{id}/preview` a `/jobs/{id}/glb` (`nas/main.go:339`).

**Layout souborů (`/data`)** — adresáře vytváří `nas/main.go:47-51`:

```
/data/ugc.db (+ -wal, -shm)          SQLite; při kopírování brát všechny tři (README.md:221-222)
/data/incoming/{id}/model.glb        surový mesh ze Sparku (TRELLIS 2048² tex ~12 MB / SF3D 1024²)
/data/incoming/{id}/preview.png      2D KONCEPT obrázek ze Sparku (ne render 3D!), volitelný
/data/incoming/{id}/preview.glb      tentýž GLB s texturami zmenšenými na 512 px (~2 MB)
/data/jobs/{id}.json                 zadání pro blender worker
/data/converted/{id}/model.fbx       Roblox-ready FBX (≤4000 tris, +Z up, −Y fwd, textura zvlášť)
/data/converted/{id}/model_tex.png   bake 1024², EMIT (bez světel)
/data/converted/{id}/report.json     report konverze (tris, bbox, watertight, …)
/data/converted/{id}/view_NNN.png    jen pokud někdo RUČNĚ spustil render_views.py
/data/packed/{id}/{name}.fbx         balíček pro Studio; {name} = první 4 slova promptu
/data/packed/{id}/{name}_tex.png
/data/packed/{id}/item.json          katalogová metadata (name, popisy, cena v Robux, tags, limited)
/data/rejected/{id}/                 přesunutý incoming/ po rejectu (README: 30 dní pro undo)
/data/characters/{id}/               fantasy characters (§5)
```

## 2. Endpointy

Všechny registrace jsou v `nas/main.go:69-101` a `nas/characters_api.go:56-78`. Chyby vždy
`{"error": "…"}` (`nas/main.go:172-174`). Žádná paginace kromě `limit`.

| Metoda | Cesta | Účel | Request → Response | Handler |
|---|---|---|---|---|
| GET | `/healthz` | ping DB | → `{"status":"ok"}` | `main.go:176` |
| POST | `/jobs` | push hotového meshe ze Sparku | multipart `glb` (povinný), `preview` (PNG, volitelný), `meta` (JSON: `id, prompt, category*, style, backend, seed, collection, symmetry`) → `201 Job`; stejné `id` = remesh: přepíše soubory, vrátí `200` existující Job | `main.go:185` |
| GET | `/jobs?status=&limit=` | seznam jobů, `created_at DESC`, limit default 200, max 1000; **filtr jen přesná shoda jednoho statusu** | → `[]Job` | `main.go:296`, `store.go:141` |
| GET | `/jobs/{id}` | detail | → `Job` | `main.go:324` |
| GET | `/jobs/{id}/preview` | koncept PNG z `incoming/` | → `image/png` | `main.go:343` |
| GET | `/jobs/{id}/glb[?full=1]` | `preview.glb` (512 px), s `full=1` originál `model.glb` | → `model/gltf-binary` | `main.go:347` |
| GET | `/viewer/{id}` | samostatná HTML stránka s `<model-viewer>` nad `/jobs/{id}/glb` | → HTML | `main.go:399` |
| GET | `/viewer/assets/*` | statika (`model-viewer.min.js`) | → soubory | `main.go:76-77` |
| POST | `/jobs/{id}/approve` | stage-aware: `new`/`failed`→`approved`; `converted`→pack (nebo `remeshing` u SF3D draftu, `202`) | → `Job` | `main.go:410` |
| POST | `/jobs/{id}/reject` | z `new/converted/failed/remeshing` → `rejected`, přesun do `rejected/` | → `Job` | `main.go:543` |
| POST | `/jobs/{id}/reconvert[?symmetry=radial\|none]` | z `converted/failed/remeshing/packed` zpět do fronty `approved` | → `Job` | `main.go:505` |
| POST | `/jobs/{id}/reroll` | nový roll stejného promptu na Sparku, job → `rerolled` | → `Job` | `main.go:563` |
| POST | `/generate` | proxy batch composeru na Spark ImageStudio | libovolný JSON → `202 {"status":"queued"}` | `main.go:589` |
| GET | `/pipeline/jobs` | proxy rozpracovaných jobů na Sparku (`/ugc/jobs`), při výpadku `[]` | → JSON ze Sparku | `main.go:609` |
| GET | `/items` | katalog zabalených kusů, **bez filtru a limitu** | → `[]Item` | `main.go:800`, `store.go:262` |
| PATCH | `/items/{id}` | částečná úprava (`name, description_cs, description_en, price_robux, collection, state, tags, limited`); **bez validace `state`** | JSON → `Item` | `main.go:809` |
| GET | `/packed/{id}/download` | zip všech souborů z `packed/{id}/` (FBX + tex PNG + item.json) | → `application/zip` | `main.go:847` |
| GET | `/events` | **SSE**, globální stream | viz níže | `main.go:87`, `sse.go:44` |
| GET | `/app/` | Flutter web build (`ugc_studio`) | — | `main.go:90-96` |
| POST | `/worker/claim` | worker si vezme 1 `approved` job (→ `converting`), jinak `204` | → `Job` | `main.go:662` |
| POST | `/worker/result/{id}` | `{verdict, error, report}` → `converted`/`failed` | → `Job` | `main.go:676` |

**SSE (`/events`)**: rámec `data: {"event": "<typ>", "data": <payload>}\n\n`, po připojení komentář
`: connected` (`sse.go:27,68,76`). Bez `id:`/`Last-Event-ID`, bez replaye, bez heartbeatu; pomalý
klient (buffer 16) zprávy **ztrácí** (`sse.go:37-40,54`). Události:
`job.created` (Job, `main.go:264`), `job.updated` (Job; approve/reject/reroll/reconvert/claim/result),
`item.created` (Item, po packu, `main.go:761`), `item.updated` (`main.go:842`), plus FC události
(globální stream dostává i topic-eventy postav, `sse.go:33-36`).

## 3. Datový model a stavy

Dvě tabulky (`nas/store.go:66-83`):

**`jobs` / `Job`** (`store.go:18-34`): `id` (`[A-Za-z0-9_-]{1,128}`, `main.go:283`), `status`, `prompt`,
`category`, `style`, `backend` (`sf3d`|`trellis`), `seed`, `collection`, `symmetry` (`""`|`radial`),
`verdict` (`PASS`/`WARN`/`FAIL`), `error`, `report` (JSON z `convert.py`), `meta` (syrový JSON ze
Sparku — navíc `reroll_of`, `checkpoint`, `created_at`, `spark/internal/ugc/pipeline.go:296-301`),
`created_at`, `updated_at`.

```
new ──approve──▶ approved ──claim──▶ converting ──result──▶ converted (verdict PASS/WARN/FAIL)
 │                  ▲                                │            │
 ├─reject─▶ rejected│                        error ─▶ failed ──approve──┘(retry)
 └─reroll─▶ rerolled (terminální)                     converted ──approve──▶ packed
                    └──── remeshing ◀── approve (jen backend=sf3d a Spark nastaven)
```
Přesné hodnoty `status`: `new, approved, converting, converted, failed, remeshing, packed,
rejected, rerolled` (`main.go:415-444,548,581`, `store.go:195-214,299`).

**`items` / `Item`** (`store.go:37-50`) — vzniká **až při packu** (`main.go:723-755`): `id` (= job id),
`name`, `description_cs/en`, `category`, `price_robux`, `tags` (JSON), `collection`, `limited` (JSON),
`state` — komentář říká `packed/uploaded/moderation/approved/onsale` (Roblox Marketplace workflow),
ale kód nastavuje jen `packed` (`main.go:723`, `store.go:223`) a Flutter appka `PATCH /items`
vůbec nevolá (grep `app/lib`). Reconvert+pack vrací `state` na `packed` (`main.go:728-729`).

**Triage = dvě síta**: 1) approve na `new` (pustí do konverze), 2) approve na `converted` = „Zabalit“
(pack). Reject je možný z obou.

**Mapování na market (`generated → reviewed → listed`):**

| Market | ugc-backend | Poznámka |
|---|---|---|
| generated | `new` (případně i `approved`/`converting`/`converted`) | jen surový GLB + koncept, nic Roblox-ready |
| reviewed | `converted` s `verdict ∈ {PASS, WARN}` | prošlo 1. sítem a konverzí, ale **ještě ne 2. sítem** |
| **ready to list** | **`status = packed`** (= existuje `items` řádek) | jediný stav, kdy majitel kus vědomě schválil podruhé a existuje finální FBX |
| listed | — | v ugc-backendu neexistuje, drží ho market (`products.status`) |

`verdict` u `packed` může být i `WARN` (např. symetrické kusy nejsou watertight, `README.md:173-176`);
`FAIL` kusy jdou zabalit taky (approve na `converted` verdikt nekontroluje, `main.go:424-440`) —
market by měl `FAIL` filtrovat.

## 4. Soubory k assetu

| Formát | Existuje? | Kde / jak | Pozn. |
|---|---|---|---|
| **GLB** | ANO (surový) | `incoming/{id}/model.glb`, HTTP `GET /jobs/{id}/glb?full=1` | high-poly výstup TRELLIS/SF3D (TRELLIS ~32k vertexů, 2048² textura, ~12 MB). **Není** to zkonvertovaný ≤4k-tris model; neškálovaný na Roblox studs. |
| GLB pro web viewer | ANO | `incoming/{id}/preview.glb`, `GET /jobs/{id}/glb` | `nas/glbpreview.go:20-24`: jen downsample PNG textur na 512 px, geometrie stejná; generuje se async po pushi + backfill při startu (`main.go:235-239,880`) |
| **FBX** | ANO | `converted/{id}/model.fbx`, `packed/{id}/{name}.fbx` | `blender_scripts/convert.py:444-454,507`; textura **neembedovaná** (`path_mode=STRIP`), ≤`max_tris` 4000, attachment empty `{Attachment}_Att` |
| Textura PNG | ANO | `model_tex.png` / `{name}_tex.png`, 1024² | `convert.py:497`, `spec/roblox_spec.json` `texture_max_px` |
| ZIP | ANO | `GET /packed/{id}/download` (FBX + PNG + item.json), streamuje se on-the-fly | `main.go:847-877` |
| **OBJ** | **NE** | nikde v kódu (grep `obj_export`) | |
| **RBXM** | **NE** | nikde v repu; do Robloxu se jde ručně přes Studio Accessory Fitting Tool (`README.md:77-79`) | |
| Turntable / render 3D | **NE automaticky** | `blender_scripts/render_views.py` existuje, ale **nikdo ho nevolá** — jen ruční `docker exec … --angles 0,180` (`README.md:178-182`); výstup `converted/{id}/view_000.png`, 720 px, EMIT osvětlení (`render_views.py:1-21`). Žádný HTTP endpoint. | |
| Náhled PNG | ANO, ale 2D | `incoming/{id}/preview.png` = koncept obrázek ze Sparku (`spark/.../pipeline.go:304`), ne render meshe | po rejectu se přesune → 404 |

**Metadata k dispozici:** z `Job`: `prompt`, `category` (`hat, hair, face, neck, shoulder, front,
back, waist, sword, helmet` — `nas/spec/roblox_spec.json`), `style`, `collection`, `backend`, `seed`,
`symmetry`, `verdict`, `meta.checkpoint`. Z `report` (`convert.py:512-533`): `tri_count`, `max_tris`,
`uv_ok`, `texture_size`, `bbox` (Blender units po škálování), `bbox_limit_studs`, `bbox_target_studs`,
`bbox_within_limit`, `watertight`, `open_edges`, `components`, `floater_faces_removed`,
`symmetry_seam_verts`, `material_count`, `attachment` (např. `HatAttachment`), `fbx`, `texture`,
`verdict_reasons` (`worker/worker.py:63-64`). Z `Item`: `name`, popisy, `tags`, `price_robux`.
Velikosti souborů ani hashe API nevrací.

## 5. Fantasy characters (`/v1/fc/`)

Druhá linka: 2D obrázek → rigovaný, animovaný humanoid (Mixamo kostra). Vlastní tabulky
`characters`, `character_steps`, `animations`, `exports` (`nas/characters.go:159-203`), stavy
`uploaded → preprocessed → meshed → cleaned → rigged → animated → exported → done` + `failed`
(`nas/README.md:57-60`). Endpointy (`nas/characters_api.go:60-73`): `POST/GET /v1/fc/characters`,
`GET/DELETE /v1/fc/characters/{id}`, `…/retry`, `…/animations`, `…/download?format=glb|fbx|zip`
(`characters_api.go:352`), `…/file/{artifact}`, `…/events` (SSE per postava, `:385`), `…/viewer`,
`…/export {target: roblox|luanti|user}`, `GET /v1/fc/exports/{id}`, `GET/PUT /v1/fc/animations`.
Soubory: `final_glb`, `final_fbx`, `preview_mp4` (turntable, default vypnutý — `FC_PREVIEW=thumb`),
`thumb_png`; `Character.tri_count` (`characters.go:19-38`). Auth `FC_API_KEYS`.

**Pro market:** pro MVP **ne** — README ji vede jako rozpracovanou (ComfyUI workflowy na Sparku
chybí, `README.md:42-44`), je určená pro vlastní hry (MountainsSimulator, DoggioWars) a animace
pocházejí z Mixamo (`seed_animlib.py --license mixamo-embedded`) — **přeprodej Mixamo klipů je
licenčně problematický**. Později jako samostatný typ produktu „rigged character“ (GLB/FBX/zip +
thumb už existují).

## 6. Integrace pro market

**Doporučený pull (bez změn v ugc-backendu, funguje hned):**
1. market-api běží na JODA a volá `http://ugc-api:8095` přes sdílenou docker síť (připojit
   market compose k externí síti `ugcfactory_ugc`) nebo `http://NAS_HOST:8095` — bez Access tokenu.
2. **Polling** každých ~5 min: `GET /jobs?status=packed&limit=1000` + `GET /items` (spárovat podle
   `id`). Upsert podle `asset_ref = job.id`; změna = `job.updated_at` / `item.updated_at` > snapshot.
   Nový kus → `products.status='draft'`, nikdy rovnou `listed`.
3. **SSE `/events` jen jako trigger** („něco se změnilo, pullni hned“) na `item.created`/`job.updated`
   — ne jako zdroj pravdy: žádný replay, ztrácí zprávy (`sse.go:37-40`). Webhook neexistuje.
4. Kus, který v ugc zmizí ze `packed` (reconvert → `approved`, reject), označit `archived`/skrýt.
   Pozor: reconvert z `packed` jde zpět do `approved` a pak znovu `packed` s novým FBX — při
   návratu přepočítat sha256 souborů (verzi produktu).

**Snapshot do `products.attributes` (JSONB):**
```json
{"ugc_id":"…","prompt":"…","category":"hat","roblox_attachment":"HatAttachment",
 "style":"…","collection":"Crowns","backend":"trellis","symmetry":"radial",
 "verdict":"WARN","verdict_reasons":["mesh neni watertight"],
 "polycount":3870,"texture_px":1024,"uv_ok":true,"watertight":false,
 "bbox_studs":[1.8,1.2,1.8],"materials":1,"formats":["glb","fbx"],
 "name":"…","tags":[…],"price_robux":0,"ugc_updated_at":"2026-…"}
```
(`polycount` z `report.tri_count` platí pro FBX; tris surového GLB backend nezná.)

**Soubory:** market-api na JODA → **read-only bind mount** `/media/storage/ugc:/ugc:ro` a kopie do
vlastního úložiště (MinIO/FS) při ingestu: `incoming/{id}/preview.glb` → `viewer_glb`,
`incoming/{id}/preview.png` → `thumb`, `packed/{id}/*` → `fbx` + textura, `incoming/{id}/model.glb` →
`glb`. **Kopírovat, ne odkazovat** — reject/reconvert soubory přesouvá a přepisuje. HTTP fallback:
`/jobs/{id}/glb`, `/jobs/{id}/glb?full=1`, `/jobs/{id}/preview`, `/packed/{id}/download`.
Zákazník nikdy nesmí dostat URL na ugc.ol1n.com (je za Access).

**Návrhy malých změn v ugc-backendu (jen návrhy, nic neimplementováno):**
1. **Flag „publikovat na market“**: sloupec `items.market_state` (`""|publish|published|withdrawn`)
   + tlačítko v `ugc_studio`; market pak pulluje `?market=publish`. Nepoužívat `items.state`, ten je
   určený pro Roblox Marketplace workflow.
2. **Filtr a inkrement**: `GET /items?state=&updated_since=` a `GET /jobs?status=a,b&updated_since=`
   (dnes jen jeden status, `store.go:144-147`; items bez filtru vůbec).
3. **Jeden „asset“ endpoint** `GET /assets/{id}` = Job + Item + seznam souborů s `bytes`/`sha256`/URL.
4. **Export GLB z konvertovaného modelu** (`bpy.ops.export_scene.gltf` vedle FBX v `convert.py:507`),
   aby prodávaný GLB odpovídal FBX (≤4k tris, Roblox scale) — dnes se GLB a FBX liší geometrií.
5. **OBJ export** (`bpy.ops.wm.obj_export` + MTL + PNG) ve stejném kroku.
6. **Turntable automaticky**: zavolat `render_views.py` ve workeru po konverzi (např. 8 úhlů, WebP)
   + endpoint `GET /jobs/{id}/views/{n}`. CPU render 720 px/16 samples by měl být levný (na rozdíl
   od FC Cycles turntable).
7. **Interní auth token** (`UGC_API_KEYS`, stejný mechanismus jako `fcAuth`) pro read endpointy,
   nebo aspoň Access service token pro market, pokud by běžel mimo JODA.
8. **Webhook / outbox**: tabulka událostí s rostoucím `seq` a `GET /events?since=seq`, aby šel
   spolehlivě dohnat výpadek (alternativa k pollingu).
9. **RBXM**: headless generování není triviální (Roblox nemá oficiální Linux tooling; existuje
   `rbx-dom`/Lune pro zápis `.rbxm` s `MeshPart`, ale mesh musí být nahraný na Roblox asset ID) —
   plán to stejně řadí jako opt-in, do MVP nedávat.

## 7. Mezery a neznámé

- **RBXM a OBJ neexistují**; plán (§ „Blender headless worker … exporty (FBX, RBXM, OBJ)“ a
  „turntable PNG/WebP“) předpokládá víc, než worker umí. MVP reálně: **GLB (surový) + FBX + PNG + zip**.
- **Turntable** neexistuje automaticky; jediný auto-náhled je 2D koncept + `preview.glb`.
- GLB (high-poly, originální měřítko) ≠ FBX (≤4k tris, Roblox studs) — rozhodnout, co se prodává.
- Stav `listed`/„publikovat na market“ v ugc-backendu není; `items.state` nikdo nemění.
- Bez auth na LAN; market na JODA dostane plný zápis (approve/reject/PATCH) — dobrovolně volat jen GET.
- Kusy po remeshi: druhý push nepřepíše `meta`/`prompt` v DB (`main.go:247-259`), `backend` se přepíše
  na `trellis` (`main.go:497`) — metadata jsou jinak stabilní.
- `rejected/` má „30 dní pak cron smaže“ (`main.go:552`) — cron v repu není; neověřeno na JODA.
- Item `name` je odvozené z promptu (`safeName`, `main.go:765`), anglické a krátké; popisy
  `description_cs/en` jsou prázdné, dokud je někdo nevyplní (appka to nedělá) → market potřebuje
  vlastní copywriting (LLM přes SPARK broker).
- Licence výstupů: SF3D = Stability Community License (do $1M), TRELLIS MIT (`README.md:95`,
  `docs/LICENSES.md`) — ověřit, zda SF3D kusy (hlavně draft, který nebyl remeshnutý) smí na prodej.
- Neověřeno živě (repo jen čteno): skutečný počet `packed` kusů, velikosti souborů, zda `DATA_DIR`
  je opravdu `/media/storage/ugc` na produkci a zda je `ugc-api` dosažitelné z jiného compose projektu.
