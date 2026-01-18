# Intégration Jellyfin / Plex (Docker & naming “pro”)

Ce projet peut s’intégrer à un serveur média (Jellyfin/Plex) de deux façons :

1) **Un nommage compatible media-server** (scan fiable) : `Season 01` + `S01E01`
2) **Un refresh automatique** (déclenché après les téléchargements, avec debounce)

---

## 1) Nommage compatible Jellyfin/Plex

### Activer le mode “media”

- Via variables d’environnement :

```bash
ASD_OUTPUT_NAMING_MODE=media
```

- Ou via `config.ini` :

```ini
[OUTPUT]
naming_mode = media
```

### Layout généré (mode `media`)

- Dossiers :

```text
<root>/<Series>/Season 01/
```

- Fichier :

```text
<Series> - S01E01.mp4
```

### Option : séparer VF/VOSTFR dans des dossiers distincts

Si tu télécharges plusieurs langues et veux éviter de mélanger :

- Env :

```bash
ASD_MEDIA_SEPARATE_LANG=1
```

- Ou `config.ini` :

```ini
[OUTPUT]
media_separate_lang = true
```

Résultat :

```text
<root>/<Series> [VOSTFR]/Season 01/...
<root>/<Series> [VF]/Season 01/...
```

---

## 2) Refresh automatique Jellyfin

### Variables nécessaires

```bash
ASD_JELLYFIN_URL=http://jellyfin:8096
ASD_JELLYFIN_API_KEY=xxxxxxxxxxxxxxxx
```

Optionnel (debounce) :

```bash
ASD_MEDIA_REFRESH_DEBOUNCE_SECONDS=20
```

Notes :
- Le refresh est **best-effort** (si l’API est KO, ça ne casse pas les downloads).
- Par défaut, le refresh s’auto-active si Jellyfin/Plex est configuré.
- Tu peux forcer : `ASD_MEDIA_REFRESH_ENABLED=1` (ou `0` pour désactiver).

### Tester Jellyfin (URL + clé) avec `curl`

Depuis l’hôte (si Jellyfin expose `8096:8096`) :

```bash
curl -fsS "http://localhost:8096/System/Info/Public" | head
```

Tester la clé API (doit répondre HTTP 204/200 selon version) :

```bash
curl -fsS -X POST "http://localhost:8096/Library/Refresh" \
  -H "X-Emby-Token: $ASD_JELLYFIN_API_KEY" \
  -o /dev/null
echo "OK"
```

### Obtenir une API key Jellyfin

- Dashboard Jellyfin → **Advanced / API Keys** (selon version) → créer une clé.

### Checklist “premier démarrage” (Jellyfin + Docker)

1) Ouvre Jellyfin : `http://localhost:8096`.
2) Termine le wizard (compte admin).
3) Crée une clé API : **Dashboard → Advanced → API Keys**.
4) Ajoute une librairie **Séries** :
  - Type : Séries
  - Dossier : `/media`
5) Lance un scan (ou attends le refresh auto après un download).

Si tu ne vois pas les fichiers :
- Vérifie que le volume `./videos` (ou `ASD_HOST_DOWNLOAD_ROOT`) est bien monté dans Jellyfin sur `/media`.
- Vérifie que tu utilises `ASD_OUTPUT_NAMING_MODE=media` pour un scan “propre”.

---

## 3) Refresh automatique Plex

### Variables nécessaires

```bash
ASD_PLEX_URL=http://plex:32400
ASD_PLEX_TOKEN=xxxxxxxxxxxx
ASD_PLEX_SECTION_ID=1
```

- `ASD_PLEX_SECTION_ID` = l’ID de la section “Library” (Films/Séries) que tu veux rafraîchir.

Notes Plex :
- Le refresh utilise une requête HTTP vers `/library/sections/<id>/refresh`.
- Le token Plex se récupère via la doc Plex (compte/token) ; une fois défini, garde-le dans `.env` (non committé).

### Tester Plex (URL + token) avec `curl`

Vérifier que Plex répond :

```bash
curl -fsS "http://localhost:32400/identity" | head
```

Déclencher un refresh (retour souvent vide, mais doit être HTTP 200) :

```bash
curl -fsS "http://localhost:32400/library/sections/$ASD_PLEX_SECTION_ID/refresh?X-Plex-Token=$ASD_PLEX_TOKEN" \
  -o /dev/null
echo "OK"
```

---

## 4) Exemple Docker Compose (Downloader + Jellyfin)

### Option A (recommandé) : utiliser le `docker-compose.yml` du repo

Le repo inclut un service Jellyfin optionnel via un **profile**.

1) Crée ton `.env` :

```bash
cp .env.example .env
```

2) Renseigne au minimum :

```dotenv
ASD_HOST_DOWNLOAD_ROOT=./videos
ASD_OUTPUT_NAMING_MODE=media
ASD_JELLYFIN_API_KEY=...
```

3) Lance le stack :

```bash
docker compose --profile media up --build
```

- UI downloader (dev) : http://localhost:5173
- Jellyfin : http://localhost:8096

Dans Jellyfin : ajoute une librairie **Séries** pointant sur `/media`.

### Option B : compose minimal “standalone”

Exemple minimal avec volume partagé (adapte les chemins) :

```yaml
services:
  asd:
    image: ghcr.io/<ton-user>/<ton-repo>:latest
    environment:
      - ASD_DOWNLOAD_ROOT=/data/videos
      - ASD_ALLOWED_DEST_PREFIXES=/data/videos
      - ASD_OUTPUT_NAMING_MODE=media
      - ASD_JELLYFIN_URL=http://jellyfin:8096
      - ASD_JELLYFIN_API_KEY=${JELLYFIN_API_KEY}
      - ASD_MEDIA_REFRESH_DEBOUNCE_SECONDS=20
    volumes:
      - ./videos:/data/videos

  jellyfin:
    image: jellyfin/jellyfin:latest
    ports:
      - "8096:8096"
    volumes:
      - ./jellyfin-config:/config
      - ./jellyfin-cache:/cache
      - ./videos:/media
```

Côté Jellyfin :
- Ajouter une librairie **Séries** pointant sur `/media`.

### Option C : prod (backend sert la SPA build) + Jellyfin

Si tu utilises le mode prod (un seul service `app`), tu peux lancer Jellyfin via le même profile :

```bash
docker compose -f docker-compose.prod.yml --profile media up --build
```

- App : http://localhost:8000
- Jellyfin : http://localhost:8096

---

## 5) Conseils scan

- Évite les sous-dossiers “langue” sous `Season 01` (ça casse souvent le scan). Préfère `ASD_MEDIA_SEPARATE_LANG=1` si besoin.
- Si tu changes de mode de nommage, un **rescan** (ou suppression/ajout de librairie) peut être nécessaire.
