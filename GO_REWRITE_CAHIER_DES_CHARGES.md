# Cahier des charges — Réécriture Go (Anime-Sama-Downloader)

Date : 2026-01-18  
Branche : `go-rewrite`

## 1) Objectif
Recréer l’application sous forme d’un **service** ("app service") en Go, avec :
- une **API HTTP** pour piloter toutes les commandes (téléchargements, recherche, configuration, suivi)
- une **Web UI** (front statique servi par le service)
- une distribution **"lite"** sans Web UI (client CLI, ou serveur headless sans assets)

Contraintes clés :
- **builds / releases propres** (un binaire par OS, publication GitHub Releases)
- **parité fonctionnelle** avec les features actuelles (naming Jellyfin/Plex, refresh, etc.)
- architecture maintenable (Clean/Hexagonal) + bonnes pratiques (tests, CI, observabilité)

## 2) Périmètre fonctionnel (ligne rouge)
### 2.1 Fonctionnalités “must-have” (V1)
- **Recherche** (équivalent actuel) : requête texte → résultats normalisés
- **Téléchargement** : créer un job de téléchargement d’épisode(s)
- **Suivi** : progression, logs, état des jobs en temps réel
- **Annulation** : annuler un job en cours + arrêt propre des workers
- **Nommage** :
  - mode `legacy`
  - mode `media-server` (Jellyfin/Plex-friendly) : `Series/Season 01/Series - S01E01.ext`
  - option séparation par langue si activée
- **Refresh media servers** (best-effort, jamais bloquant) :
  - Jellyfin : refresh bibliothèque
  - Plex : refresh section(s)
- **Persistance** :
  - jobs persistés (historique + reprise après restart)
  - settings persistés
- **Web UI** :
  - config (destinations, tokens, base URL, limites)
  - écran jobs (liste, détail, progression)
  - bouton “Tester API Jellyfin/Plex”

### 2.2 Hors-périmètre initial (post-V1)
- Multi-user / auth avancée (OAuth) (V1.1+)
- Plugins third-party (V2)
- Scheduler “calendrier” (V1.1+)

## 3) Livrables
- `asd-server` : service HTTP (API + events + UI statique)
- `asd` : client CLI (pilotage via API) + option standalone si souhaité
- Documentation :
  - installation
  - configuration
  - API (OpenAPI)
  - guide Jellyfin/Plex
- Pipeline CI/CD : build multi-OS, checksums, release auto

## 4) Architecture cible (Clean / Hexagonal)
### 4.1 Couches
1) **Domain** (pur)
- Entités : `Series`, `Episode`, `DownloadJob`, `JobState`, `NamingPlan`
- Règles : transitions d’état, validation des options

2) **Application** (use-cases)
- `SearchUseCase`
- `CreateJobUseCase`
- `CancelJobUseCase`
- `GetJobUseCase`, `ListJobsUseCase`
- `RefreshMediaUseCase`, `TestMediaApisUseCase`
- `UpdateSettingsUseCase`

3) **Ports** (interfaces)
- `SourceProvider` (Anime-Sama aujourd’hui, extensible)
- `Downloader` (HLS/TS, MP4)
- `Transcoder/Muxer` (ffmpeg)
- `JobRepository` (SQLite)
- `SettingsRepository` (SQLite)
- `MediaServerClient` (Jellyfin, Plex)
- `EventBus` (SSE)

4) **Adapters**
- HTTP API (Gin/Echo/Chi)
- SSE events
- SQLite (sqlc + migrations)
- Jellyfin/Plex clients
- wrapper ffmpeg

### 4.2 Design patterns à appliquer (ligne rouge)
- **Strategy** : naming mode, provider, downloader
- **State Machine** : états de job (voir §5)
- **Observer/EventBus** : push events vers UI/CLI
- **Repository** : DB isolée
- **Dependency Injection** : construction explicite + interfaces
- **Retry Policy** : backoff/jitter, timeouts
- **Circuit Breaker (soft)** : sur Jellyfin/Plex (optionnel V1)

## 5) Modèle Jobs (state machine)
États minimaux :
- `queued` → `running` → `muxing` → `completed`
- `queued|running|muxing` → `canceled`
- `queued|running|muxing` → `failed`

Règles :
- transitions atomiques (persistées)
- logs associés au job (stream + storage)

## 6) API HTTP (v1)
Base : `/api/v1`

### 6.1 Health / Meta
- `GET /health`
- `GET /version`

### 6.2 Search
- `GET /search?q=...` (ou `POST /search`)

### 6.3 Jobs
- `POST /jobs` : crée un job (body = sélection épisodes + options)
- `GET /jobs` : liste paginée + filtre par état
- `GET /jobs/{id}` : détail
- `POST /jobs/{id}/cancel`
- `GET /jobs/{id}/logs` (optionnel si events)

### 6.4 Events
- `GET /events` (SSE) : progression + logs + changements d’état

### 6.5 Settings
- `GET /settings`
- `PUT /settings`

### 6.6 Media servers
- `POST /media/test`
- `POST /media/refresh`

### 6.7 OpenAPI
- `/openapi.json` généré (ligne rouge)

## 7) Stockage
- SQLite unique (par défaut) : `asd.db`
- Tables : `jobs`, `job_events`/`job_logs`, `settings`, `migrations`
- Migrations versionnées

## 8) Observabilité & robustesse
- Logs structurés (JSON) + niveaux
- Correlation ID (request → job)
- Timeouts par défaut (HTTP fetch, media servers)
- Concurrence configurable : `max_workers`, `max_concurrent_downloads`
- Arrêt propre : drain queue, stop workers, flush refresh

## 9) Packaging / releases
- `goreleaser` : Linux/macOS/Windows + checksums
- Artifacts : `asd-server`, `asd`
- Option : embed UI (assets `dist/`) dans `asd-server`

## 10) Jalons (proposition)
### M0 — Socle repo (1-2 jours)
- Structure Go (modules, lint, tests)
- Serveur HTTP minimal + health
- SQLite + migrations

### M1 — Jobs + events (3-7 jours)
- CRUD jobs
- workers + queue
- SSE

### M2 — Download pipeline (1-2 semaines)
- provider Anime-Sama
- téléchargement + muxing
- naming modes

### M3 — Jellyfin/Plex (2-4 jours)
- refresh debounce
- test API endpoints

### M4 — Web UI intégrée (3-7 jours)
- UI settings + jobs
- build frontend + embed

### M5 — Parité & release (1 semaine)
- tests end-to-end minimal
- goreleaser + release v0.1

## 11) Critères d’acceptation (V1)
- Un job de batch episodes peut : se lancer, se suivre, s’annuler, se reprendre après restart
- Sorties conformes aux modes de nommage, chemins safe
- Refresh Jellyfin/Plex déclenché après succès, sans casser le téléchargement si indisponible
- Web UI : config + suivi jobs + test API
- Release : binaires multi-OS + checksums, install simple

## 12) Questions ouvertes (à trancher)
- Mode “lite” :
  - A) `asd` = client CLI uniquement (reco)
  - B) `asd` = standalone executor (plus lourd mais autonome)
- Transport events : SSE vs WebSocket (SSE recommandé pour simplicité)
- Front : garder React/Vite (probable), ou UI plus simple

