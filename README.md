# Anime-Sama Downloader

[![Python 3.9+](https://img.shields.io/badge/Python-3.9+-3776AB?logo=python&logoColor=white)](https://www.python.org/)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)

Téléchargeur d’épisodes depuis **anime-sama.si** avec plusieurs interfaces (CLI / TUI / Web) et une file de téléchargements avec parallélisme.

**Fork de :** https://github.com/sertrafurr/Anime-Sama-Downloader

## Points clés

- Recherche par nom (AniList) + résolution automatique d’URL
- Téléchargements en parallèle (jusqu’à 10)
- Support MP4 + HLS/M3U8 (segments `.ts`) + conversion MP4 (FFmpeg/MoviePy)
- Interface Web (FastAPI + React) et TUI (Textual)
- Configuration simple : `config.ini` + variables d’environnement

## Démarrage rapide

### Installation (Python)

```bash
python3 -m pip install -r requirements.txt
```

Optionnel mais recommandé : `ffmpeg` (conversion plus rapide).

### Utilisation (CLI)

```bash
# Mode interactif
python main.py

# Recherche + téléchargement
python main.py -s "demon slayer" -e 1-5 --quick

# URL directe
python main.py -u "https://anime-sama.si/catalogue/sword-art-online/saison1/vostfr/" -e 1-10 -t

# TUI (Textual)
python main.py --tui
```

### Interface Web (dev local)

```bash
./scripts/dev-backend.sh
./scripts/dev-frontend.sh
```

- SPA : http://127.0.0.1:5173
- Backend : http://127.0.0.1:8000

Voir aussi : [QUICK_START.md](QUICK_START.md)

## Docker

### Dev (backend + frontend)

```bash
docker compose up --build
```

Accès : http://localhost:5173

### Prod (backend sert la SPA build)

```bash
docker compose -f docker-compose.prod.yml up --build
```

Accès : http://localhost:8000

### Sortie vidéos en Docker (important)

- **Dans le conteneur**, la sortie est **toujours** `/data/videos`.
- **Sur l’hôte**, ce dossier est monté via : `ASD_HOST_DOWNLOAD_ROOT`.

Pour personnaliser le dossier sur l’hôte :

```bash
cp .env.example .env
# édite .env
```

Dans l’interface Web en Docker : la “destination” est un **sous-dossier relatif** sous `/data/videos` (pour éviter d’écrire dans le FS interne du conteneur).

## Configuration

Priorité (du plus fort au plus faible) : **CLI > variables d’env > `config.ini`**.

### `config.ini`

```bash
cp config.ini.example config.ini
```

### Variables d’environnement utiles

```bash
ASD_CONFIG=./config.ini
ASD_DOWNLOAD_ROOT=./videos
ASD_MAX_CONCURRENT_DOWNLOADS=10

# Domaine (si ça bouge)
ASD_SITE_BASE_URL=https://anime-sama.si

# Web
ASD_WEB_HOST=0.0.0.0
ASD_WEB_PORT=8000

# Frontend (Vite) en dev local
ASD_WEBAPP_HOST=127.0.0.1
ASD_WEBAPP_PORT=5173

# Docker: dossier hôte monté sur /data/videos
ASD_HOST_DOWNLOAD_ROOT=/chemin/absolu/sur/hote
```

## Documentation

- [QUICK_START.md](QUICK_START.md) : commandes prêtes à l’emploi
- [SEARCH_GUIDE.md](SEARCH_GUIDE.md) : recherche, résolution d’URL, conseils
- [MIGRATION.md](MIGRATION.md) : notes de migration (domaines, Docker, nouveautés)
- [CHANGELOG.md](CHANGELOG.md) : historique des changements

## Notes légales

Outil à but éducatif. Respecte le droit d’auteur et les CGU du site.
