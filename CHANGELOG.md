# Changelog (Go rewrite)

Ce dépôt a été réécrit autour d’un serveur Go (API + jobs + abonnements + scheduler) et d’une UI web (React/Vite) servie en SPA.

Les notes détaillées et exemples de commandes de l’implémentation Python historique sont disponibles dans l’historique git (tags/commits antérieurs).

## Janvier 2026

- Backend Go : jobs persistés, worker pool, scheduler abonnements, SSE (`/api/v1/events`)
- API : settings, jobs, subscriptions, endpoints AniList + import auto, résolution Anime‑Sama (`/api/v1/animesama/resolve`)
- UI : servie depuis `webapp/dist` (fallback SPA) + client `/api/v1/*`
- Docker : image Go + SPA build, ports par défaut sur `8080`
