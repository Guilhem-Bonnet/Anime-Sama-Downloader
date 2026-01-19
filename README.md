# Anime-Sama Downloader

[![Go 1.22+](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)

Téléchargeur d’épisodes depuis **anime-sama.si** piloté par un **serveur Go** (API + jobs + scheduler) et une **UI React** servie en SPA.

**Fork de :** https://github.com/sertrafurr/Anime-Sama-Downloader

## Crédits

- Projet original : https://github.com/sertrafurr/Anime-Sama-Downloader (fork, non affilié)
- Recherche / planning : https://anilist.co

## Points clés

- Abonnements (clé = base URL saison/langue) + scheduler de checks
- Jobs persistés + SSE (`/api/v1/events`) + worker pool
- AniList (token optionnel) : viewer, airing schedule, watchlist + import auto → abonnements
- HLS/M3U8 via `ffmpeg` si disponible

## Démarrage rapide (local)

### 1) Lancer le serveur

```bash
go run ./cmd/asd-server
```

Par défaut :
- UI/API : http://127.0.0.1:8080
- DB SQLite : `asd.db`

### 2) (Optionnel) Builder l’UI pour qu’elle soit servie par le serveur

```bash
npm -C webapp ci
npm -C webapp run build
```

Le serveur sert automatiquement `webapp/dist` si présent.

### 3) Ouvrir

- UI : http://127.0.0.1:8080
- OpenAPI : http://127.0.0.1:8080/api/v1/openapi.json

Voir aussi : [QUICK_START.md](QUICK_START.md)

## Dev UI (Vite)

```bash
./scripts/dev-backend.sh
./scripts/dev-frontend.sh
```

- SPA : http://127.0.0.1:5173 (proxy vers `/api/*`)
- Backend : http://127.0.0.1:8080

## Docker

### Dev (backend + frontend)

```bash
docker compose up --build
```

Accès : http://localhost:5173

### Prod (image unique, serveur Go sert la SPA build)

```bash
docker compose -f docker-compose.prod.yml up --build
```

Accès : http://localhost:8080

### Volumes (Docker)

- Vidéos : monter l’hôte sur `/data/videos` via `ASD_HOST_DOWNLOAD_ROOT`
- DB : monter l’hôte sur `/data` via `ASD_HOST_DATA_ROOT` (fichier `/data/asd.db`)

```bash
cp .env.example .env
# édite .env puis relance docker compose
```

## Configuration

La config se fait via l’UI (onglet **settings**) ou l’API : `PUT /api/v1/settings`.

Petit CLI fourni (optionnel) :

```bash
go run ./cmd/asd settings get
go run ./cmd/asd settings set --destination /chemin/vers/videos --max-concurrent-downloads 6
```

Variables d’environnement serveur :

```bash
ASD_ADDR=127.0.0.1:8080
ASD_DB_PATH=asd.db
ASD_WEB_DIST=webapp/dist
```

## Notes sur l’ancien code Python

L’implémentation Python historique a été retirée de cette branche. Pour la retrouver, utilise l’historique git (tags/commits antérieurs).

## Notes légales

Outil à but éducatif. Respecte le droit d’auteur et les CGU du site.
