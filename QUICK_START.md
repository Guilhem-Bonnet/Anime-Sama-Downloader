# 🚀 Quick Start (Go server + UI + Docker)

Ce guide donne les commandes “prêtes à copier-coller”. Pour le détail : [README.md](README.md).

## 1) Prérequis

- Go 1.22+
- Node 20+ (uniquement si tu rebuild l’UI)
- `ffmpeg` (recommandé si tu télécharges du HLS/M3U8)

## 2) Lancer le serveur (local)

```bash
go run ./cmd/asd-server
```

Ouvre : http://127.0.0.1:8080

Changer l’adresse / la DB :

```bash
ASD_ADDR=127.0.0.1:8099 ASD_DB_PATH=./asd.db go run ./cmd/asd-server
```

## 3) Builder l’UI (pour qu’elle soit servie par le serveur)

```bash
npm -C webapp ci
npm -C webapp run build
```

Ensuite, le serveur sert automatiquement `webapp/dist`.

## 4) Dev UI (Vite)

Backend :

```bash
./scripts/dev-backend.sh
```

Frontend :

```bash
./scripts/dev-frontend.sh
```

Ouvre :
- http://127.0.0.1:5173 (SPA, proxy vers le backend)
- http://127.0.0.1:8080 (backend)

## 5) Docker

### Dev

```bash
docker compose up --build
```

Accès : http://localhost:5173

### Prod

```bash
docker compose -f docker-compose.prod.yml up --build
```

Accès : http://localhost:8080

### Volumes (Docker)

- Vidéos : `/data/videos` (hôte → variable `ASD_HOST_DOWNLOAD_ROOT`)
- DB : `/data/asd.db` (hôte → variable `ASD_HOST_DATA_ROOT`)

```bash
cp .env.example .env
# éditer .env puis relancer docker compose
```

---

## ⚡ Astuces

1. Configure la destination via l’UI (**settings**) ou via l’API `PUT /api/v1/settings`
2. Mets ton token AniList dans `settings.anilistToken` pour activer viewer/watchlist/import
3. La spec OpenAPI est disponible sur `/api/v1/openapi.json`
