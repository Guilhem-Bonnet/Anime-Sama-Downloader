# 🧭 Migration (ancien Python → serveur Go)

Ce repo a été réécrit autour d’un **serveur Go** (API + jobs + scheduler) avec une **UI web**.

## 1) Ancien code Python

L’implémentation Python historique (CLI/TUI/FastAPI) n’est plus présente sur cette branche.
Si tu en as besoin, récupère-la via l’historique git (tags/commits antérieurs).

## 2) Nouveau point d’entrée

Le workflow principal passe par :

- serveur : `go run ./cmd/asd-server` (UI + API)
- UI : http://127.0.0.1:8080

## 3) Domaine Anime‑Sama

Le serveur travaille avec des URLs `anime-sama.si` et normalise certains domaines connus (ex: `anime-sama.tv` → `anime-sama.si`).

## 4) Docker

Deux variables utiles :

- `ASD_HOST_DOWNLOAD_ROOT` → monté sur `/data/videos` (vidéos)
- `ASD_HOST_DATA_ROOT` → monté sur `/data` (DB `/data/asd.db`)

```bash
cp .env.example .env
docker compose -f docker-compose.prod.yml up --build
```

Accès : http://localhost:8080
