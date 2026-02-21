---
stepsCompleted:
  - step-01-init
  - step-02-discovery
  - step-03-success
  - step-04-journeys
  - step-05-domain
  - step-06-innovation
  - step-07-project-type
  - step-08-scoping
  - step-09-functional
  - step-10-nonfunctional
  - step-11-polish
  - step-12-complete
inputDocuments:
  - _bmad-output/planning-artifacts/00-PROJECT-BRIEF.md
  - _bmad-output/planning-artifacts/01-PERSONAS-AND-JOURNEYS.md
  - _bmad-output/planning-artifacts/02-DESIGN-SYSTEM-SAKURA-NIGHT.md
  - _bmad-output/planning-artifacts/03-TECHNICAL-ARCHITECTURE.md
  - _bmad-output/planning-artifacts/04-BACKLOG-SPRINTS.md
  - _bmad-output/planning-artifacts/05-PM-EXECUTION-GUIDE.md
  - _bmad-output/planning-artifacts/README.md
workflowType: 'prd'
documentCounts:
  briefCount: 1
  researchCount: 0
  brainstormingCount: 0
  projectDocsCount: 6
classification:
  projectType: Web app + API publique + CLI (version allégée)
  domain: Gestion de contenus médias / téléchargement d’anime (général, non régulé)
  complexity: Medium
  projectContext: brownfield
---

# Product Requirements Document - Anime-Sama-Downloader

**Author:** Guilhem
**Date:** 31 janvier 2026

## Success Criteria

### User Success

- Télécharger un épisode en ≤ 4 clics et ≤ 30s jusqu’au job “en cours”.
- Recherche compréhensible et sans ambiguïté : saisons/épisodes/langues explicitement clarifiés.
- Taux d’échec perçu par l’utilisateur < 5% sur 50 actions.

### Business Success

- Usage privé (serveurs personnels).
- ≤ 3 utilisateurs actifs.
- Une seule queue partagée pour tous les utilisateurs.

### Technical Success

- API P95 < 300ms (hors scraping).
- Taux d’échec des jobs < 2%.
- Coverage tests ≥ 70%.
- Lighthouse ≥ 90.
- Uptime local ≥ 99%.

### Measurable Outcomes

- 90% des téléchargements démarrent en ≤ 30s après action.
- 0 confusion signalée sur saisons/langues/épisodes (tests utilisateurs).
- Max 3 utilisateurs actifs sans dégradation perceptible.

## Product Scope

### MVP - Minimum Viable Product

- UI/UX refonte
- Recherche fiable
- Téléchargement fiable
- Jobs + queue partagée
- Settings essentiels

### Growth Features (Post-MVP)

- Subscriptions
- Calendar
- Notifications
- UI avancée (animations, polish)

### Vision (Future)

- Multi-users (limité à 3)
- AniList
- Jellyfin
- API publique + CLI

## User Journeys

### 1) Alex (Primary — Happy Path)

**Opening** : Alex ouvre l’app le soir. Il veut le dernier épisode sans réfléchir.
**Rising** : Il voit un dashboard clair, tape le titre, filtre VOSTFR, saison/épisode évidents.
**Climax** : Il clique “Télécharger” et voit la queue démarrer immédiatement.
**Resolution** : Progression visible, notification claire, fichier bien nommé. Il se sent en confiance.

### 2) Alex (Primary — Edge Case)

**Opening** : Alex cherche un anime avec plusieurs saisons et versions.
**Rising** : Il hésite entre VF/VOSTFR, saisons multiples. L’UI explicite les labels et propose une prévisualisation.
**Climax** : Il corrige facilement un choix erroné sans perdre sa sélection.
**Resolution** : Il comprend enfin la différence saisons/épisodes/langues et termine le téléchargement.

### 3) Maya (Admin/Ops)

**Opening** : Maya gère un serveur privé partagé (≤ 3 users).
**Rising** : Elle configure settings globaux, vérifie la queue partagée et les limites.
**Climax** : Elle crée des abonnements et déclenche des jobs groupés.
**Resolution** : La queue reste stable, les jobs s’enchaînent proprement, elle garde la maîtrise.

### 4) API/CLI User (Lite)

**Opening** : Un utilisateur veut aller vite sans UI.
**Rising** : Il lance un CLI minimal avec recherche + téléchargement direct.
**Climax** : Il récupère l’épisode via une seule commande, feedback clair.
**Resolution** : Il obtient le résultat sans friction, et peut automatiser.

### Journey Requirements Summary

- Recherche explicite (saisons/épisodes/langues)
- Téléchargement déclenché en ≤ 30s
- Queue partagée claire et observable
- CLI rapide avec feedback compréhensible
- UX de correction sans perte (edge cases)

## Domain-Specific Requirements

### Compliance & Regulatory

- Usage privé recommandé (serveurs personnels), pas d’exposition publique.
- Avertissement d’usage dans README + UI Settings (optionnel).

### Technical Constraints

- Limite de queue : 10 jobs simultanés, surplus en attente.
- Rétention logs : à définir (ex : 7/30 jours), rotation configurable.

### Integration Requirements

- Dépendance scraping anime-sama (site + source GitHub).
- ffmpeg requis pour muxing.

### Risk Mitigations

- Changement HTML anime-sama → tests + fallback + alerte.
- Étudier récupération directe des sources si besoin (charge inconnue).
- Rate limit strict pour éviter ban.

## Web App Specific Requirements

### Project-Type Overview

Single Page Application (SPA) React avec API REST. Navigateurs modernes supportés : Firefox, Chrome, Safari. Usage privé, pas de SEO. CLI allégé en complément.

### Technical Architecture Considerations

- **Architecture** : SPA (Single Page Application)
- **Navigateurs supportés** : Firefox, Chrome, Safari (versions modernes)
- **Real-time** : SSE pour stream jobs/progress
- **Accessibilité** : WCAG AA si aligné avec design system
- **SEO** : Non applicable (usage privé)

### API & Integration Requirements

**Versioning :**
- URL path : `/api/v1`
- Header optionnel : `X-API-Version`

**Rate Limits :**
- Public : 100 req/min/IP
- Local : 300 req/min/IP
- Burst allowance : 30 requêtes

**Endpoints clés (Web + CLI) :**
- `GET /api/v1/search` — Recherche animes
- `GET/POST /api/v1/subscriptions` — Gestion abonnements
- `GET /api/v1/jobs` — Status queue + jobs
- `GET/PUT /api/v1/settings` — Configuration
- `POST /api/v1/animesama/resolve` — Résolution anime-sama
- `GET /api/v1/events` — SSE stream (jobs progress)

### CLI & Command-Line Interface

- Endpoints accessibles via CLI (interface allégée)
- Même rate limiting appliqué
- Feedback clair (JSON ou text)

## Project Scoping & Phased Development

### MVP Strategy & Philosophy

**MVP Approach :** Problème-résolution minimal viable. Focus sur fiabilité recherche + téléchargement pour valider product-market fit avec utilisateurs privés.

**Resource Requirements (Équipe A) :** 7-8 personnes
- 1 Tech Lead (architecture + reviews)
- 2 Backend devs (Go, core services)
- 2 Frontend devs (React, UI)
- 1 UX/UI designer (Sakura Night implementation)
- 1 QA engineer (tests + E2E)
- 1 DevOps/Infrastructure (optionnel)

**Timeline :** 10-13 semaines (2,5-3 mois)

### MVP Feature Set (Phase 1 — 2-3 semaines)

**Core User Journeys Supported :**
- Alex (Happy Path) : téléchargement fiable en ≤ 4 clics
- Alex (Edge Case) : gestion des recherches ambiguës
- Maya (Admin) : config settings + queue partagée

**Must-Have Capabilities :**
- UI/UX refonte (design system "Sakura Night")
- Recherche anime fiable (saisons/langues explicites)
- Téléchargement fiable (queue ≤ 10 jobs parallèles)
- Job tracking + progress bar SSE
- Settings configuration (basiques)
- Cleanup code + tests ≥ 70%

### Post-MVP Features

**Phase 2 (Growth — 3-4 semaines après MVP) :**
- Subscriptions (auto-check recurring)
- Calendar (airing schedule)
- Notifications (completion, errors)
- Animations + UI polish

**Phase 3 (Vision — post-Phase 2) :**
- Multi-users auth (≤ 3 users)
- AniList integration (auto-import watchlist)
- Jellyfin integration (sync + naming)
- API publique + CLI lite

### Risk Mitigation Strategy

**Technical Risks :**
| Risque | Probabilité | Mitigation |
|--------|-------------|-----------|
| Changement domaine anime-sama | Fréquent | Source GitHub comme fallback, monitoring, tests d'intégration |
| Fermeture anime-sama | Possible | Documenter fallback, étudier récupération directe sources (Phase 3) |
| Changement structure site HTML | Fréquent | Abstraction scraper, rate limiting strict, alertes |

**Resource Risks :**
- Équipe réduite ? Repasser à Scénario B (4-5 personnes), délai +30%.
- Designer indisponible ? Utiliser design system existant, focus sur implémentation.

**Domain Risks :**
- Anime-sama bans scraping ? Rate limit + fallbacks + GitHub source minimisent.

### MVP Success Criteria

- **User** : Télécharger 1 épisode en ≤ 4 clics, ≤ 30s → job en cours
- **Tech** : API P95 < 300ms, job failure < 2%, tests ≥ 70%, Lighthouse ≥ 90
- **Business** : Support ≤ 3 users concurrents, queue partagée stable

## Functional Requirements (MVP + Growth + Vision)

### Capability 1: Content Discovery & Search

**FR1: User can search anime by title with autocomplete**
- When user types in search field, autocomplete suggestions appear (max 10 results)
- Suggestions filtered from cached anime database
- <100ms response time
- Phase: MVP

**FR2: User can filter search results by status (ongoing, completed, planning)**
- Filter controls appear below search field
- Multiple filters can be combined (AND logic)
- Persist selected filters in session storage
- Phase: MVP

**FR3: User can view detailed anime metadata page**
- Display: title, synopsis, episode count, air dates, genres, source URL
- Fetch additional metadata from anime-sama.si for episodes list
- Load time <300ms P95
- Phase: MVP

**FR4: System maintains local cache of anime metadata**
- Cache refreshed every 24 hours or on user request
- Fallback to cache if anime-sama.si unavailable
- Cache stored in SQLite database
- Phase: MVP

**FR5: User can sort search results by title, date added, popularity**
- Sort options available as dropdown in search results header
- Default sort: by title (A-Z)
- Persist sort preference in settings
- Phase: MVP

**FR6: System implements anime-sama.si scraper with abstraction layer**
- Scraper module decoupled from UI/API
- Support for HTML structure version detection
- Fallback to GitHub mirror if anime-sama.si unavailable
- Include retry logic (3 attempts, exponential backoff)
- Phase: MVP

### Capability 2: Download & Queue Management

**FR7: User can add anime to download queue with episode range selection**
- User selects start/end episode numbers or "all episodes"
- Queue entry created with initial status "queued"
- Max 10 concurrent downloads enforced (surplus queued)
- Phase: MVP

**FR8: User can view real-time download progress via Server-Sent Events (SSE)**
- Open SSE connection to `/api/v1/jobs/{jobId}/progress`
- Stream updates: current episode, bytes downloaded, ETA, download speed
- Update frequency: every 500ms
- Connection auto-recovers on disconnect (3-second retry)
- Phase: MVP

**FR9: User can pause/resume/cancel individual downloads**
- Pause button stops ffmpeg process, preserves partial files
- Resume button restarts from last checkpoint
- Cancel button removes job from queue and cleans up files
- Changes reflected in UI within 1 second
- Phase: MVP

**FR10: User can reorder queued jobs (drag-and-drop)**
- Reorder affects execution order for queued jobs
- Not applicable to running jobs
- Persist new order to database
- Phase: Growth

**FR11: System stores downloaded episodes locally with HLS encoding**
- Episode files encoded to H.264 video + AAC audio via ffmpeg
- HLS playlist (.m3u8) + 10-second segments (.ts) generated
- Stored in `data/videos/{animeId}/` directory structure
- Re-encode skipped if HLS already exists (checksum validation)
- Phase: MVP

**FR12: User can view list of downloaded anime with episode counts**
- Library view shows: anime title, cover image, episode count, last updated date
- Click to view downloaded episodes for specific anime
- Delete button removes local files and database entries
- Phase: MVP

### Capability 3: Settings & Configuration

**FR13: User can configure download quality (720p, 1080p, original)**
- Quality option stored in settings (SQLite `settings` table)
- Affects ffmpeg encoding parameters
- Default: 720p
- Phase: MVP

**FR14: User can configure download destination path**
- Settings page allows custom path input with file picker
- Validate path is writable before saving
- Fallback to `./data/videos/` if custom path unavailable
- Phase: MVP

**FR15: User can enable/disable HLS encoding**
- Toggle in settings (default: enabled)
- If disabled, store original anime files without re-encoding
- Saves ~70% encoding time but increases storage
- Phase: MVP

**FR16: User can configure notification preferences (email, in-app, none)**
- Notification settings persist in database
- Phase: Growth

**FR17: User can import/export settings as JSON file**
- Export button downloads `settings.json`
- Import button accepts JSON file upload with validation
- Backup/restore full configuration
- Phase: Growth

**FR18: System persists all user settings to SQLite database**
- Settings table schema: `id`, `key`, `value`, `updated_at`
- Settings cached in-memory with 5-minute refresh TTL
- Changes immediately written to database
- Phase: MVP

### Capability 4: Jobs & Monitoring

**FR19: System maintains persistent job queue in SQLite**
- Jobs table: `id`, `animeId`, `status`, `episodes`, `createdAt`, `startedAt`, `completedAt`, `errorMessage`
- Statuses: `queued`, `downloading`, `encoding`, `completed`, `failed`, `paused`
- Restart unfinished jobs on application restart
- Phase: MVP

**FR20: User can view job history with timestamps**
- History page shows all completed/failed jobs (scrollable)
- Display: anime title, episodes downloaded, duration, status, timestamp
- Filter by status (completed, failed) and date range
- Phase: Growth

**FR21: System logs all operations to rotating file**
- Log file: `logs/asd-server.log` (rotate daily, keep 7 days)
- Log levels: ERROR, WARN, INFO, DEBUG
- Include: timestamp, level, component, message
- Phase: MVP

**FR22: User can view system health metrics (uptime, job success rate, resource usage)**
- Dashboard shows: total jobs, success count, failure count, current download speed
- Calculated from job history and real-time stats
- Auto-refresh every 5 seconds
- Phase: Growth

**FR23: System monitors anime-sama.si availability and alerts user**
- Health check runs every 30 minutes (ping anime-sama.si)
- If unavailable for 2+ consecutive checks, alert appears in UI
- Log to system logs
- Phase: MVP

**FR24: System implements automatic retry logic for failed jobs**
- Retry failed jobs up to 3 times with exponential backoff (1min, 5min, 15min)
- Mark as permanently failed after 3 retries
- Log failure reason (network timeout, encoding error, metadata unavailable, etc.)
- Phase: MVP

### Capability 5: Multi-User Support

**FR25: System supports multiple user accounts with role-based access (admin, user)**
- Admin: can manage users, view all jobs, modify global settings
- User: can only view/manage own jobs, modify own settings
- Phase: Vision (Phase 3)

**FR26: User can share download queue with other users**
- Admin creates shared queue accessible to multiple users
- Shared job progress visible to all queue members
- Downloaded episodes accessible to all users
- Phase: Vision

**FR27: System maintains per-user download history and preferences**
- Each user has isolated settings and job history
- Shared downloads appear in all users' libraries
- Phase: Vision

**FR28: User can set download speed limits per user (throttle bandwidth)**
- Admin sets global limit; users can set personal sub-limit
- Applied via ffmpeg `-b:v` and `-b:a` parameters
- Phase: Vision

**FR29: System enforces per-user download quota (max episodes/week)**
- Admin configurable quota (e.g., 100 episodes/week)
- Track used quota, reject if exceeded, reset weekly
- Phase: Vision

### Capability 6: API & CLI

**FR30: System exposes RESTful API at `/api/v1/`**
- Endpoints: GET/POST jobs, GET progress, GET anime search, GET settings
- Authentication: API key (bearer token) or session cookie for local usage
- Rate limiting: 100 req/min/IP (public), 300 req/min/IP (local)
- Response format: JSON with consistent error structure
- Phase: Growth (Phase 2) for basic endpoints; Phase 3 for full API

**FR31: API implements versioning via URL path and optional header**
- Primary: `/api/v1/` (URL path version)
- Optional: `X-API-Version: 1` header for client preference
- Support up to 3 major versions concurrently
- Phase: Growth

**FR32: System provides CLI tool for headless usage**
- CLI commands: `asd search <title>`, `asd download <id> [episodes]`, `asd status`, `asd settings`
- Distributed as standalone binary (Go) or shell wrapper
- Support scripting and cron job integration
- Phase: Vision

**FR33: API endpoints implement pagination for list results**
- Query params: `limit` (default 20, max 100), `offset` (default 0)
- Response includes `total`, `limit`, `offset`, `items[]`
- Phase: Growth

**FR34: API returns consistent error responses with HTTP status codes**
- 400: Bad request (invalid params)
- 401: Unauthorized (invalid API key)
- 404: Not found
- 429: Rate limit exceeded
- 500: Server error
- Error body: `{ "error": "message", "code": "ERROR_CODE" }`
- Phase: Growth

**FR35: System logs all API requests to audit trail**
- Log: timestamp, method, path, status code, response time, IP, API key (hashed)
- Keep 30 days of audit logs
- Phase: Growth

### Capability 7: Data Persistence

**FR36: System uses SQLite for all persistent data**
- Database file: `data/asd.db`
- Tables: `anime`, `episodes`, `jobs`, `settings`, `users` (Phase 3)
- Automatic migrations on schema changes
- Phase: MVP

**FR37: System backs up database daily to timestamped file**
- Backup path: `data/backups/asd-{YYYY-MM-DD}.db`
- Keep last 30 days of backups
- Automated cleanup of old backups
- Phase: MVP

**FR38: System validates database integrity on startup**
- Run `PRAGMA integrity_check` on application start
- If integrity issues detected, restore from latest backup and alert admin
- Log to system logs
- Phase: MVP

**FR39: System exports user data as JSON/CSV for backup/migration**
- Export includes: anime metadata, job history, settings, downloaded files list
- Format: JSON (primary) or CSV (for spreadsheet view)
- Phase: Growth

### Capability 8: Subscriptions & Calendar

**FR40: User can view calendar of anime airing dates**
- Calendar displays all currently airing anime with air date/time
- Grouped by day of week (e.g., all Monday releases, all Thursday releases)
- Shows episode number and air time (e.g., "15:00 JST")
- Sortable by: day of week, popularity, genre
- Search/filter by genre, language (VOSTFR, VF)
- Display covers/thumbnails for each anime
- Phase: Growth (Phase 2)

**FR41: User can view detailed anime information and subscribe**
- Click anime in calendar → open detail modal/page
- Display: title, cover, synopsis, episode count, genres, air time, current status
- "Subscribe" button triggers subscription flow
- Link to search page if user wants to download existing episodes
- Phase: Growth (Phase 2)

**FR42: User can select subscription scope (all episodes, specific season, future only)**
- Subscription modal offers three options:
  - **All available**: Download all existing episodes + future releases
  - **Specific season**: Select which season(s) to download + future releases of selected season
  - **Future releases only**: Only new episodes as they air (no backlog)
- For "All available" and "Specific season": offer option to download immediately or on-demand
- Confirm selection before subscribing
- Phase: Growth (Phase 2)

**FR43: System automatically downloads new episodes when aired**
- New episode released → system detects via anime-sama.si scraper
- Auto-download triggered within 30 minutes of air time (configurable)
- Episode added to download queue with priority based on subscription date
- User receives notification (in-app + optional email) when download starts/completes
- Download occurs even if app is minimized (background job)
- If download fails, retry logic applies (same as FR24)
- Phase: Growth (Phase 2)

**FR44: User can manage subscriptions (list, modify, unsubscribe)**
- "My Subscriptions" page shows all active subscriptions
- Display: anime title, cover, subscription type (all/season/future), next air date, last downloaded episode
- Modify button: change scope (all → future only, specific season → different season, etc.)
- Unsubscribe button: remove subscription, optionally keep downloaded episodes
- Pause subscription: temporarily stop auto-downloads without unsubscribing
- View subscription history: completed, paused, unsubscribed (keep for 30 days)
- Phase: Growth (Phase 2)

---

**Functional Requirements Summary:**
- **Total FRs:** 44 organized in 8 capabilities
- **MVP (Phase 1):** 24 FRs (search, download, queue, settings, jobs, persistence)
- **Growth (Phase 2):** 15 additional FRs (API, notifications, history, export, subscriptions, calendar)
- **Vision (Phase 3):** 5 additional FRs (multi-user, advanced CLI)
- **Binding Contract:** No feature outside these FRs will be implemented
- **Dependencies:** All downstream UX design, architecture, and development scoped by these FRs


## Non-Functional Requirements

### Performance

**NFR1: API Latency**
- All API responses must be < 300ms P95 (excluding anime-sama.si scraping calls)
- Measured from request receipt to response completion
- Baseline: 95th percentile of response times across all GET/POST endpoints

**NFR2: Search Autocomplete Performance**
- Autocomplete suggestions must appear within < 100ms from last keystroke
- Filtered results delivered from cached anime database
- Maximum 10 suggestions per query

**NFR3: Server-Sent Events (SSE) Streaming**
- Job progress updates sent every 500ms ± 100ms
- Connection stability: auto-recover within 3 seconds on network disconnect
- Message delivery rate: >99% for non-network-failure scenarios
- Support up to 10 concurrent SSE connections

**NFR4: ffmpeg HLS Encoding Performance**
- H.264 video encoding at real-time speed or faster (1.0x play speed minimum)
- For 24-30 minute episodes: complete encoding within 30-40 minutes wall-clock time
- Quality: 720p H.264 + AAC audio (default), configurable to 1080p or original
- Acceptable degradation: -5% performance if system CPU load >80%

**NFR5: Metadata Cache Hit Rate**
- Anime metadata cache hit rate >90% for repeat searches
- Cache refresh every 24 hours or on manual user request
- Cache size: <500MB for typical anime database (10,000+ titles)

### Reliability

**NFR6: Job Queue Persistence**
- All jobs in queue persisted to SQLite database
- Jobs surviving application restart, resume from checkpoint
- State preserved: job ID, anime ID, episode range, progress, status
- Recovery time on restart: <5 seconds for queue restoration

**NFR7: Job Failure Recovery**
- Failed jobs automatically retried up to 3 times
- Retry backoff: 1 minute → 5 minutes → 15 minutes
- After 3 retries, mark job as permanently failed and log reason
- Failure reasons logged: network timeout, encoding error, metadata unavailable, ffmpeg crash, disk space full
- Success rate target: >98% on first attempt or within retry window

**NFR8: anime-sama.si Availability Fallback**
- Health check every 30 minutes (ping anime-sama.si homepage)
- If unreachable for 2+ consecutive checks (60 minutes total), alert user
- Fallback to GitHub mirror source for scraping
- User-facing error message: "anime-sama.si temporarily unavailable, using fallback source"
- Failover transparent to user; no action required

**NFR9: Database Integrity Validation**
- PRAGMA integrity_check executed on application startup
- If corruption detected, automatically restore from latest backup
- Alert admin/user if restore occurs (log file + optional UI notification)
- Backups retained: 30 days of daily backups

### Security

**NFR10: Data Encryption**
- In Transit: HTTPS/TLS 1.3+ for all API communication (Phase 3 public API)
- At Rest: Downloaded anime files stored in user's local filesystem (user responsible for encryption)
- Database: Stored unencrypted in SQLite (local access only, appropriate for private usage)
- Threat Model: Assume local network is trusted; public internet threats mitigated by HTTPS

**NFR11: API Authentication & Rate Limiting**
- Local Usage: Session cookies (optional; no auth required for localhost)
- Public API (Phase 3): Bearer token authentication (API keys)
- Rate Limits: 100 requests/min/IP (public), 300 requests/min/IP (local), 30-request burst allowance
- Throttle response: HTTP 429 with Retry-After header

**NFR12: Input Validation & SQL Injection Prevention**
- All user inputs validated: search terms, file paths, episode ranges, settings values
- SQL queries use parameterized statements (prevent injection)
- Search terms: max 200 characters, alphanumeric + spaces only
- Episode ranges: 1-9999, validated as integers
- File paths: validated against whitelist directory, no `../` sequences allowed

**NFR13: Sensitive Data Handling**
- Never log: passwords, API keys, auth tokens, user identities
- Audit trail: API requests logged with method/path/status, IP, API key hashed (SHA-256)
- Log retention: 30 days of audit logs
- Error messages: Don't expose internal paths, stack traces, or database schema

### Accessibility

**NFR14: WCAG AA Compliance**
- All interactive UI elements must be WCAG AA compliant
- Color contrast ratio ≥ 4.5:1 for normal text, ≥ 3:1 for large text
- Images must have alt text describing content/function
- Form fields must have associated labels
- Target: 100% WCAG AA compliance for critical user flows (search, download, queue)

**NFR15: Semantic HTML & ARIA Labels**
- Semantic HTML5 tags used (nav, main, article, aside, etc.)
- ARIA labels for dynamic content (job progress, status updates)
- Role attributes for custom components (role="button", role="tab")
- Focus indicators visible on all interactive elements

**NFR16: Keyboard Navigation**
- All UI interactive elements navigable via keyboard (Tab, Enter, Escape)
- Modal dialogs: focus trapped, Escape closes
- Dropdowns: arrow keys to navigate options
- Buttons: Enter or Space to activate
- Search field: Enter submits search, Escape clears focus

---

**Non-Functional Requirements Summary:**
- **Total NFRs:** 16 (Performance: 5, Reliability: 4, Security: 4, Accessibility: 3)
- **Scalability & Compliance:** Not included (out of scope for this product)
- **Testability:** All NFRs are measurable and have specific success criteria
- **Binding:** NFRs become quality acceptance criteria for QA, dev review, and product sign-off

