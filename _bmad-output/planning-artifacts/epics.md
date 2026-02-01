---
stepsCompleted:
  - step-01-validate-prerequisites
  - step-02-design-epics
  - step-03-create-stories
inputDocuments:
  - _bmad-output/planning-artifacts/prd.md
status: complete
totalEpics: 8
totalStories: 49
---

# Anime-Sama-Downloader - Epic Breakdown

## Overview

This document provides the complete epic and story breakdown for Anime-Sama-Downloader, decomposing the requirements from the PRD into implementable stories.

## Requirements Inventory

### Functional Requirements

**Capability 1: Content Discovery & Search**

- **FR1:** User can search anime by title with autocomplete (Phase: MVP)
- **FR2:** User can filter search results by status (ongoing, completed, planning) (Phase: MVP)
- **FR3:** User can view detailed anime metadata page (Phase: MVP)
- **FR4:** System maintains local cache of anime metadata (Phase: MVP)
- **FR5:** User can sort search results by title, date added, popularity (Phase: MVP)
- **FR6:** System implements anime-sama.si scraper with abstraction layer (Phase: MVP)

**Capability 2: Download & Queue Management**

- **FR7:** User can add anime to download queue with episode range selection (Phase: MVP)
- **FR8:** User can view real-time download progress via Server-Sent Events (SSE) (Phase: MVP)
- **FR9:** User can pause/resume/cancel individual downloads (Phase: MVP)
- **FR10:** User can reorder queued jobs (drag-and-drop) (Phase: Growth)
- **FR11:** System stores downloaded episodes locally with HLS encoding (Phase: MVP)
- **FR12:** User can view list of downloaded anime with episode counts (Phase: MVP)

**Capability 3: Settings & Configuration**

- **FR13:** User can configure download quality (720p, 1080p, original) (Phase: MVP)
- **FR14:** User can configure download destination path (Phase: MVP)
- **FR15:** User can enable/disable HLS encoding (Phase: MVP)
- **FR16:** User can configure notification preferences (email, in-app, none) (Phase: Growth)
- **FR17:** User can import/export settings as JSON file (Phase: Growth)
- **FR18:** System persists all user settings to SQLite database (Phase: MVP)

**Capability 4: Jobs & Monitoring**

- **FR19:** System maintains persistent job queue in SQLite (Phase: MVP)
- **FR20:** User can view job history with timestamps (Phase: Growth)
- **FR21:** System logs all operations to rotating file (Phase: MVP)
- **FR22:** User can view system health metrics (uptime, job success rate, resource usage) (Phase: Growth)
- **FR23:** System monitors anime-sama.si availability and alerts user (Phase: MVP)
- **FR24:** System implements automatic retry logic for failed jobs (Phase: MVP)

**Capability 5: Multi-User Support**

- **FR25:** System supports multiple user accounts with role-based access (admin, user) (Phase: Vision/Phase 3)
- **FR26:** User can share download queue with other users (Phase: Vision)
- **FR27:** System maintains per-user download history and preferences (Phase: Vision)
- **FR28:** User can set download speed limits per user (throttle bandwidth) (Phase: Vision)
- **FR29:** System enforces per-user download quota (max episodes/week) (Phase: Vision)

**Capability 6: API & CLI**

- **FR30:** System exposes RESTful API at `/api/v1/` (Phase: Growth for basic endpoints; Phase 3 for full API)
- **FR31:** API implements versioning via URL path and optional header (Phase: Growth)
- **FR32:** System provides CLI tool for headless usage (Phase: Vision)
- **FR33:** API endpoints implement pagination for list results (Phase: Growth)
- **FR34:** API returns consistent error responses with HTTP status codes (Phase: Growth)
- **FR35:** System logs all API requests to audit trail (Phase: Growth)

**Capability 7: Data Persistence**

- **FR36:** System uses SQLite for all persistent data (Phase: MVP)
- **FR37:** System backs up database daily to timestamped file (Phase: MVP)
- **FR38:** System validates database integrity on startup (Phase: MVP)
- **FR39:** System exports user data as JSON/CSV for backup/migration (Phase: Growth)

**Capability 8: Subscriptions & Calendar**

- **FR40:** User can view calendar of anime airing dates (Phase: Growth/Phase 2)
- **FR41:** User can view detailed anime information and subscribe (Phase: Growth/Phase 2)
- **FR42:** User can select subscription scope (all episodes, specific season, future only) (Phase: Growth/Phase 2)
- **FR43:** System automatically downloads new episodes when aired (Phase: Growth/Phase 2)
- **FR44:** User can manage subscriptions (list, modify, unsubscribe) (Phase: Growth/Phase 2)

**Total: 44 Functional Requirements**
- MVP (Phase 1): 24 FRs
- Growth (Phase 2): 15 FRs
- Vision (Phase 3): 5 FRs

### Non-Functional Requirements

**Performance**

- **NFR1:** API Latency — All API responses < 300ms P95 (excluding anime-sama.si scraping calls)
- **NFR2:** Search Autocomplete Performance — Suggestions appear within < 100ms from last keystroke
- **NFR3:** Server-Sent Events (SSE) Streaming — Updates every 500ms ± 100ms, auto-recover within 3s on disconnect
- **NFR4:** ffmpeg HLS Encoding Performance — Real-time speed or faster (1.0x play speed minimum), 24-30 min episodes in 30-40 min wall-clock
- **NFR5:** Metadata Cache Hit Rate — >90% for repeat searches

**Reliability**

- **NFR6:** Job Queue Persistence — All jobs persisted to SQLite, survive restart, <5s recovery time
- **NFR7:** Job Failure Recovery — Failed jobs retried 3x with backoff (1/5/15 min), >98% success rate target
- **NFR8:** anime-sama.si Availability Fallback — Health check every 30 min, failover to GitHub mirror if down 2+ checks
- **NFR9:** Database Integrity Validation — PRAGMA integrity_check on startup, auto-restore from backup if corruption

**Security**

- **NFR10:** Data Encryption — HTTPS/TLS 1.3+ for public API (Phase 3), local storage unencrypted
- **NFR11:** API Authentication & Rate Limiting — Bearer tokens (Phase 3), 100 req/min/IP (public), 300 req/min/IP (local)
- **NFR12:** Input Validation & SQL Injection Prevention — All inputs validated, parameterized queries, search terms max 200 chars
- **NFR13:** Sensitive Data Handling — Never log passwords/API keys/tokens, audit trail with hashed keys, 30 days retention

**Accessibility**

- **NFR14:** WCAG AA Compliance — Contrast ratio ≥ 4.5:1 for normal text, alt text for images, labeled form fields
- **NFR15:** Semantic HTML & ARIA Labels — HTML5 semantic tags, ARIA labels for dynamic content, role attributes
- **NFR16:** Keyboard Navigation — All elements navigable via Tab/Enter/Escape, focus management in modals

**Total: 16 Non-Functional Requirements**
- Performance: 5 NFRs
- Reliability: 4 NFRs
- Security: 4 NFRs
- Accessibility: 3 NFRs

### Additional Requirements

**Technical Architecture Context (from PRD):**

- SPA React frontend with Zustand state management
- Go backend with Clean Architecture
- SQLite database (data/asd.db)
- ffmpeg for HLS encoding (H.264 + AAC)
- Server-Sent Events (SSE) for real-time progress
- API versioning via URL path (/api/v1)
- Rate limiting: 100 req/min/IP (public), 300 req/min/IP (local)

**Domain-Specific Requirements (from PRD):**

- Scraping dependency: anime-sama.si (site + GitHub mirror fallback)
- Job queue limit: Max 10 concurrent downloads
- Log retention: 7-30 days, rotating daily
- Usage recommendation: Private servers, ≤3 concurrent users
- Browser support: Firefox, Chrome, Safari (modern versions)

**Design System (from PRD):**

- Sakura Night theme (dark mode, magenta/cyan accents)
- WCAG AA accessibility target
- Design tokens CSS approach
- Component library structure

### FR Coverage Map

**FR1** → Epic 2 (Search autocomplete)
**FR2** → Epic 2 (Filter results)
**FR3** → Epic 2 (Metadata page)
**FR4** → Epic 2 (Cache)
**FR5** → Epic 2 (Sort results)
**FR6** → Epic 2 (Scraper)
**FR7** → Epic 3 (Add to queue)
**FR8** → Epic 3 (SSE progress)
**FR9** → Epic 3 (Pause/resume/cancel)
**FR10** → Epic 5 (Reorder queue)
**FR11** → Epic 3 (HLS encoding)
**FR12** → Epic 3 (Downloaded library)
**FR13** → Epic 4 (Quality config)
**FR14** → Epic 4 (Download path)
**FR15** → Epic 4 (HLS toggle)
**FR16** → Epic 4 (Notifications)
**FR17** → Epic 4 (Import/export settings)
**FR18** → Epic 1 (Settings persistence)
**FR19** → Epic 1 (Job queue persistence)
**FR20** → Epic 5 (Job history)
**FR21** → Epic 1 (Logging)
**FR22** → Epic 5 (Health metrics)
**FR23** → Epic 2 (Site monitoring)
**FR24** → Epic 3 (Retry logic)
**FR25** → Epic 8 (Role-based access)
**FR26** → Epic 8 (Shared queue)
**FR27** → Epic 8 (Per-user history)
**FR28** → Epic 8 (Bandwidth limits)
**FR29** → Epic 8 (Quotas)
**FR30** → Epic 7 (REST API)
**FR31** → Epic 7 (API versioning)
**FR32** → Epic 7 (CLI tool)
**FR33** → Epic 7 (Pagination)
**FR34** → Epic 7 (Error responses)
**FR35** → Epic 7 (Audit trail)
**FR36** → Epic 1 (SQLite)
**FR37** → Epic 1 (Backups)
**FR38** → Epic 1 (DB integrity)
**FR39** → Epic 7 (Export data)
**FR40** → Epic 6 (Calendar view)
**FR41** → Epic 6 (Subscribe flow)
**FR42** → Epic 6 (Subscription scope)
**FR43** → Epic 6 (Auto-download)
**FR44** → Epic 6 (Manage subscriptions)

**Coverage: 44/44 FRs (100%)**

## Epic List

### Epic 1: Project Foundation & Infrastructure
Installer l'infrastructure de base pour permettre le développement de toutes les fonctionnalités. L'application peut démarrer, persister des données, et logger les opérations.

**FRs covered:** FR18, FR19, FR21, FR36, FR37, FR38, NFR6, NFR9, NFR12

### Epic 2: Anime Search & Discovery
Les utilisateurs peuvent chercher, filtrer, et découvrir des anime avec métadonnées complètes. Trouver n'importe quel anime en <30s avec des informations claires.

**FRs covered:** FR1, FR2, FR3, FR4, FR5, FR6, FR23, NFR1, NFR2, NFR5, NFR8

### Epic 3: Download Queue & Job Management
Les utilisateurs peuvent télécharger des anime avec suivi en temps réel et contrôle de la queue. Télécharger des épisodes en ≤4 clics, voir la progression, et gérer la queue.

**FRs covered:** FR7, FR8, FR9, FR11, FR12, FR24, NFR3, NFR4, NFR7

### Epic 4: Settings & Configuration
Les utilisateurs peuvent personnaliser les comportements de téléchargement et notifications. Configurer la qualité, le chemin de téléchargement, et les notifications selon les préférences.

**FRs covered:** FR13, FR14, FR15, FR16, FR17

### Epic 5: Advanced Queue Features
Les utilisateurs peuvent optimiser leur queue avec réorganisation et historique détaillé. Réorganiser la queue par priorité et consulter l'historique complet des téléchargements.

**FRs covered:** FR10, FR20, FR22

### Epic 6: Subscriptions & Calendar Discovery
Les utilisateurs peuvent découvrir des anime via un calendrier et s'abonner pour téléchargement automatique. Découvrir facilement les nouveaux anime du moment et s'abonner pour recevoir automatiquement les nouveaux épisodes.

**FRs covered:** FR40, FR41, FR42, FR43, FR44

### Epic 7: Public API & CLI Access
Les développeurs et power users peuvent accéder aux fonctionnalités via API REST et CLI. Automatiser des téléchargements via scripts ou applications tierces.

**FRs covered:** FR30, FR31, FR32, FR33, FR34, FR35, FR39, NFR10, NFR11, NFR13

### Epic 8: Multi-User Support (Vision)
Administrateurs peuvent gérer plusieurs utilisateurs avec rôles et quotas. Partager le serveur avec la famille (≤3 users) avec contrôle des accès et quotas.

**FRs covered:** FR25, FR26, FR27, FR28, FR29

## Epic 1: Project Foundation & Infrastructure

Installer l'infrastructure de base pour permettre le développement de toutes les fonctionnalités. L'application peut démarrer, persister des données, et logger les opérations.

### Story 1.1: Initialize SQLite Database Schema

As a developer,
I want to initialize the SQLite database with core tables,
So that the application can persist settings and job queue data.

**Acceptance Criteria:**

**Given** the application is starting for the first time
**When** the database initialization runs
**Then** the following tables are created: `settings`, `jobs`
**And** the `settings` table has schema: `id`, `key`, `value`, `updated_at`
**And** the `jobs` table has schema: `id`, `animeId`, `status`, `episodes`, `createdAt`, `startedAt`, `completedAt`, `errorMessage`
**And** PRAGMA integrity_check passes with "ok" status
**And** default settings are inserted (download_path, quality, hls_enabled)

### Story 1.2: Implement Application Logging System

As a developer,
I want a structured logging system with file rotation,
So that operations can be traced and debugged.

**Acceptance Criteria:**

**Given** the application is running
**When** any operation occurs (API call, job state change, error)
**Then** a log entry is written to `logs/asd-server.log`
**And** log entry includes: timestamp, level (ERROR/WARN/INFO/DEBUG), component, message
**And** log file rotates daily with pattern `asd-server-YYYY-MM-DD.log`
**And** old log files beyond 7 days are automatically deleted
**And** log level is configurable via environment variable

### Story 1.3: Implement Job Queue Persistence

As a system,
I want to persist job queue state to SQLite,
So that jobs survive application restarts.

**Acceptance Criteria:**

**Given** a job is added to the queue
**When** the job state changes (queued → downloading → completed)
**Then** the job row in `jobs` table is updated with new status and timestamps
**And** `startedAt` is set when job transitions to "downloading"
**And** `completedAt` is set when job transitions to "completed" or "failed"
**And** on application restart, unfinished jobs (status = "queued" or "downloading") are reloaded
**And** reloaded jobs resume from their last checkpoint within 5 seconds

### Story 1.4: Implement Automatic Database Backups

As a system administrator,
I want automatic daily backups of the database,
So that data is protected against corruption.

**Acceptance Criteria:**

**Given** the application has been running for 24 hours
**When** the backup scheduler triggers
**Then** a backup file is created at `data/backups/asd-YYYY-MM-DD.db`
**And** the backup is a complete copy of `data/asd.db`
**And** backups older than 30 days are automatically deleted
**And** backup operation does not block active database writes
**And** backup completion is logged with success/failure status

### Story 1.5: Implement Database Integrity Validation on Startup

As a system,
I want to validate database integrity on startup,
So that corrupted databases are detected and restored.

**Acceptance Criteria:**

**Given** the application is starting
**When** database integrity check runs via PRAGMA integrity_check
**Then** if result is "ok", application continues normally
**And** if corruption is detected, the latest backup from `data/backups/` is restored
**And** if restore occurs, an alert is logged (ERROR level) with timestamp and reason
**And** if no backups exist and corruption detected, application fails with clear error message
**And** integrity check completes within 3 seconds for databases < 1GB

### Story 1.6: Implement Input Validation Utilities

As a developer,
I want reusable input validation functions,
So that user inputs are sanitized consistently.

**Acceptance Criteria:**

**Given** user input is received (search term, file path, episode range)
**When** validation function is called
**Then** search terms are validated: max 200 characters, alphanumeric + spaces only
**And** episode ranges are validated: integers between 1-9999
**And** file paths are validated: no `../` sequences, must be within allowed directory
**And** validation errors return clear messages (e.g., "Invalid episode range: must be 1-9999")
**And** all SQL queries use parameterized statements (no string concatenation)


## Epic 2: Anime Search & Discovery

Permettre aux utilisateurs de rechercher et découvrir des animés. Ils peuvent trouver exactement ce qu'ils cherchent rapidement.

### Story 2.1: Implement Anime Search API Endpoint

As a user,
I want to search for anime by title,
So that I can find the content I want to download.

**Acceptance Criteria:**

**Given** I send a GET request to `/api/v1/search?q=naruto`
**When** the search is processed
**Then** I receive a JSON array of matching anime
**And** each result contains: `id`, `title`, `thumbnail_url`, `year`, `status`, `episode_count`
**And** results are ranked by relevance (exact match > partial match)
**And** maximum 50 results are returned
**And** response time is < 300ms (P95)

### Story 2.2: Implement Search Autocomplete

As a user,
I want autocomplete suggestions as I type,
So that I can quickly find anime titles.

**Acceptance Criteria:**

**Given** I type at least 2 characters in the search field
**When** 300ms have elapsed since last keystroke
**Then** suggestions appear below the search field
**And** maximum 10 suggestions are shown
**And** each suggestion shows: title + thumbnail + year
**And** clicking a suggestion navigates to anime detail page
**And** suggestions update within 150ms of typing pause

### Story 2.3: Implement Anime Detail View

As a user,
I want to view detailed information about an anime,
So that I can decide whether to download it.

**Acceptance Criteria:**

**Given** I click on an anime from search results
**When** the detail page loads
**Then** I see: title, thumbnail, synopsis, year, status, genre tags, episode count
**And** a "Download" button is visible
**And** episode list is displayed with season/episode numbers
**And** page loads within 200ms
**And** if anime has multiple seasons, seasons are displayed as tabs

### Story 2.4: Implement Search Result Filters

As a user,
I want to filter search results by genre, year, and status,
So that I can narrow down options.

**Acceptance Criteria:**

**Given** I have search results displayed
**When** I apply filters (genre: Action, year: 2020-2023, status: Complete)
**Then** results update to show only matching anime
**And** filter selections persist during the session
**And** filter chips are displayed showing active filters
**And** I can remove filters individually by clicking on chips
**And** filter application completes within 100ms

### Story 2.5: Implement Recent Searches History

As a user,
I want to see my recent searches,
So that I can quickly re-search without retyping.

**Acceptance Criteria:**

**Given** I have performed searches in the past
**When** I click on the search field
**Then** my last 10 searches are displayed as clickable items
**And** clicking a recent search executes that search
**And** recent searches are stored in localStorage
**And** I can clear recent searches with a "Clear History" button
**And** searches are deduplicated (same query only appears once)

### Story 2.6: Implement Anime Catalogue Cache

As a system,
I want to cache anime-sama.si catalogue in memory,
So that searches are fast without repeated scraping.

**Acceptance Criteria:**

**Given** the application starts
**When** first search request arrives
**Then** full anime catalogue is fetched from anime-sama.si (or GitHub mirror if site down)
**And** catalogue is stored in memory with TTL of 6 hours
**And** subsequent searches use in-memory cache
**And** cache refresh happens automatically every 6 hours
**And** manual cache refresh endpoint exists at `/api/v1/cache/refresh`

### Story 2.7: Implement Fallback to GitHub Mirror

As a system,
I want to fallback to GitHub anime-sama-mirror if site is down,
So that searches still work during outages.

**Acceptance Criteria:**

**Given** anime-sama.si is unreachable (timeout > 5s or HTTP 5xx)
**When** catalogue fetch is attempted
**Then** system automatically tries GitHub mirror at `github.com/Anime-Sama/anime-sama-mirror`
**And** fallback happens transparently (no user action required)
**And** a warning is logged indicating fallback was used
**And** next catalogue refresh re-attempts primary site first
**And** fallback response time is < 10s

### Story 2.8: Implement Search Result Pagination

As a user,
I want to paginate through search results,
So that I can browse large result sets efficiently.

**Acceptance Criteria:**

**Given** search returns > 50 results
**When** I view the search results page
**Then** first 50 results are displayed
**And** pagination controls appear at bottom (Previous, Page 1/N, Next)
**And** clicking Next loads results 51-100
**And** URL updates to reflect current page (e.g., `?page=2`)
**And** page navigation completes within 100ms


## Epic 3: Download Queue & Job Management

Permettre aux utilisateurs de télécharger des épisodes et suivre la progression. Ils peuvent gérer une queue de téléchargements avec visibilité en temps réel.

### Story 3.1: Implement Download Job Creation

As a user,
I want to start a download by selecting episodes,
So that content is saved to my local storage.

**Acceptance Criteria:**

**Given** I'm viewing an anime detail page
**When** I select episodes (e.g., "1-12" or "1,3,5") and click "Download"
**Then** a job is created with status "queued"
**And** job appears immediately in the queue view
**And** a success notification appears: "Download added to queue"
**And** job includes: anime_id, title, episodes, quality (from settings), hls_enabled
**And** if episodes are invalid (e.g., "1-999" for 24-episode anime), an error is shown

### Story 3.2: Implement Download Queue View

As a user,
I want to see all my download jobs in a queue,
So that I can monitor progress.

**Acceptance Criteria:**

**Given** I navigate to the queue page
**When** the page loads
**Then** I see a list of all jobs (queued, downloading, completed, failed)
**And** each job shows: anime title, episodes, status, progress %, start time, estimated time remaining
**And** jobs are grouped by status: Downloading (top), Queued, Completed (bottom)
**And** page loads within 200ms
**And** progress updates automatically every 2 seconds via SSE

### Story 3.3: Implement Real-Time Download Progress with SSE

As a user,
I want to see download progress update in real-time,
So that I know when content will be ready.

**Acceptance Criteria:**

**Given** a job is actively downloading
**When** I'm viewing the queue page
**Then** progress bar updates every 2 seconds without page refresh
**And** SSE connection is established to `/api/v1/jobs/stream`
**And** each SSE event contains: job_id, status, progress_percentage, current_episode
**And** if SSE connection drops, UI shows "Connection lost" and retries every 5s
**And** bandwidth usage is displayed (e.g., "2.5 MB/s")

### Story 3.4: Implement Download Executor with Concurrency Limit

As a system,
I want to limit concurrent downloads to 10,
So that system resources are not overwhelmed.

**Acceptance Criteria:**

**Given** 15 jobs are in the queue
**When** downloads begin
**Then** exactly 10 jobs transition to "downloading" status
**And** remaining 5 jobs stay "queued"
**And** as a downloading job completes, the next queued job starts automatically
**And** each job runs in a separate goroutine
**And** job state transitions are persisted to SQLite immediately

### Story 3.5: Implement Job Cancellation

As a user,
I want to cancel an active download,
So that I can stop unwanted jobs.

**Acceptance Criteria:**

**Given** a job is downloading or queued
**When** I click "Cancel" button on the job
**Then** job status changes to "cancelled"
**And** if job was downloading, download process stops within 3 seconds
**And** partial files are kept in `data/videos/{anime_id}/partial/`
**And** job is removed from active executors
**And** next queued job starts if executor slot becomes available
**And** a notification confirms cancellation

### Story 3.6: Implement Job Retry

As a user,
I want to retry failed jobs,
So that I can recover from temporary errors.

**Acceptance Criteria:**

**Given** a job has failed (status = "failed")
**When** I click "Retry" button
**Then** job status changes to "queued"
**And** job is re-added to the queue with same parameters
**And** download attempts from the beginning (not resuming partial files)
**And** a notification appears: "Job re-queued"
**And** error message from previous attempt is cleared

### Story 3.7: Implement Download Folder Organization

As a system,
I want to organize downloaded files by anime,
So that content is easy to locate.

**Acceptance Criteria:**

**Given** a job is downloading episodes
**When** an episode completes
**Then** file is saved to `{download_path}/{anime_title_sanitized}/S{season}E{episode}.mp4`
**And** anime title is sanitized (no special chars, max 100 chars)
**And** if HLS is enabled, file extension is `.m3u8` with chunks in `chunks/` subfolder
**And** completed files have read permissions set to 644
**And** folder structure is created automatically if it doesn't exist

### Story 3.8: Implement Job Status Filters

As a user,
I want to filter jobs by status,
So that I can focus on specific job types.

**Acceptance Criteria:**

**Given** I'm viewing the queue page
**When** I select a status filter (All, Downloading, Queued, Completed, Failed)
**Then** job list updates to show only matching jobs
**And** filter selection persists during session (localStorage)
**And** job count badges show: "Downloading (3)", "Queued (7)", etc.
**And** filter application completes within 50ms
**And** default filter is "All"

### Story 3.9: Implement Download Completion Notifications

As a user,
I want to be notified when downloads complete,
So that I know content is ready.

**Acceptance Criteria:**

**Given** a job completes successfully
**When** all episodes finish downloading
**Then** a browser notification appears: "{Anime Title} download complete"
**And** notification includes anime thumbnail
**And** clicking notification navigates to completed jobs view
**And** in-app toast notification also appears for 5 seconds
**And** notification permission is requested on first download


## Epic 4: Settings & Configuration

Permettre aux utilisateurs de personnaliser l'application. Ils peuvent configurer le comportement selon leurs préférences.

### Story 4.1: Implement Settings Page UI

As a user,
I want a settings page to configure the application,
So that I can customize behavior to my preferences.

**Acceptance Criteria:**

**Given** I navigate to the settings page
**When** the page loads
**Then** I see organized sections: Download, Video, Interface, Advanced
**And** each setting has a label, description, and input field
**And** current values are pre-populated from database
**And** a "Save" button is visible at the bottom
**And** page loads within 200ms

### Story 4.2: Implement Download Path Configuration

As a user,
I want to set my download folder path,
So that files are saved where I want them.

**Acceptance Criteria:**

**Given** I'm on the settings page
**When** I change the "Download Path" field to `/home/user/Videos`
**And** click "Save"
**Then** the value is persisted to `settings` table with key `download_path`
**And** a success toast appears: "Settings saved"
**And** new downloads use the updated path
**And** if path is invalid (doesn't exist/no permissions), an error is shown
**And** default value is `data/videos`

### Story 4.3: Implement Video Quality Selection

As a user,
I want to choose video quality (480p, 720p, 1080p),
So that I can balance quality and file size.

**Acceptance Criteria:**

**Given** I'm on the settings page
**When** I select "720p" from the quality dropdown
**And** click "Save"
**Then** the value is persisted with key `video_quality`
**And** new downloads request 720p sources from anime-sama.si
**And** if 720p is unavailable for a specific anime, highest available quality is used
**And** job view displays selected quality
**And** default value is "1080p"

### Story 4.4: Implement HLS Streaming Toggle

As a user,
I want to enable/disable HLS streaming,
So that I can choose between direct downloads and adaptive streaming.

**Acceptance Criteria:**

**Given** I'm on the settings page
**When** I toggle "Enable HLS Streaming" to ON
**And** click "Save"
**Then** the value is persisted with key `hls_enabled`
**And** new downloads generate `.m3u8` playlists instead of single `.mp4` files
**And** ffmpeg is used to segment videos into HLS chunks
**And** HLS files are organized: `{anime}/playlist.m3u8` + `chunks/segment-*.ts`
**And** default value is OFF (false)

### Story 4.5: Implement Concurrent Downloads Limit

As a user,
I want to set max concurrent downloads (1-20),
So that I can control system resource usage.

**Acceptance Criteria:**

**Given** I'm on the settings page
**When** I change "Max Concurrent Downloads" to 5
**And** click "Save"
**Then** the value is persisted with key `max_concurrent_downloads`
**And** download executor respects new limit immediately (within 5s)
**And** if I set 3 and 10 jobs are running, excess jobs continue until completion (no interruption)
**And** future jobs respect new limit
**And** default value is 10, min is 1, max is 20

### Story 4.6: Implement Settings Validation

As a system,
I want to validate all settings before persisting,
So that invalid configurations are rejected.

**Acceptance Criteria:**

**Given** I'm on the settings page
**When** I enter invalid values (e.g., download path with `../`, concurrent limit = 0)
**And** click "Save"
**Then** form submission is blocked
**And** error messages appear below invalid fields
**And** valid fields are not highlighted as errors
**And** no values are persisted until all validations pass
**And** validation rules match input validation utilities from Story 1.6


## Epic 5: Advanced Queue Features

Ajouter des fonctionnalités avancées pour gérer la queue de téléchargements. Les utilisateurs peuvent prioriser et organiser leurs jobs.

### Story 5.1: Implement Job Priority Reordering

As a user,
I want to reorder queued jobs by dragging,
So that important downloads are processed first.

**Acceptance Criteria:**

**Given** multiple jobs are queued
**When** I drag a job to a different position in the queue
**And** drop it in the new position
**Then** job order is updated immediately in the UI
**And** job priorities are persisted to database with numeric `priority` field
**And** download executor processes jobs by priority (highest first)
**And** drag handles are visible on queued jobs only (not on downloading/completed)
**And** priority changes take effect within 2 seconds

### Story 5.2: Implement Bulk Job Operations

As a user,
I want to select multiple jobs and perform bulk actions,
So that I can manage many jobs efficiently.

**Acceptance Criteria:**

**Given** I'm viewing the queue page
**When** I select checkboxes on multiple jobs
**Then** a bulk actions toolbar appears: "Cancel Selected", "Retry Selected", "Delete Selected"
**And** clicking "Cancel Selected" cancels all selected jobs in one operation
**And** a confirmation dialog appears for destructive actions (Delete)
**And** bulk operation completion is confirmed with toast: "3 jobs cancelled"
**And** "Select All" checkbox is available in the table header

### Story 5.3: Implement Job Details Sidebar

As a user,
I want to click on a job to see full details,
So that I can inspect logs and metadata.

**Acceptance Criteria:**

**Given** I'm viewing the queue page
**When** I click on a job row
**Then** a sidebar opens on the right side (300px width)
**And** sidebar shows: anime title, thumbnail, episodes, quality, status, timestamps, error messages (if failed)
**And** for active jobs, real-time logs are streamed in a scrollable panel
**And** sidebar includes "Close" button or can be closed by clicking outside
**And** sidebar opening/closing has smooth animation (200ms)

### Story 5.4: Implement Auto-Retry for Failed Jobs

As a system,
I want to automatically retry failed jobs up to 3 times,
So that transient errors are handled gracefully.

**Acceptance Criteria:**

**Given** a job fails due to network error
**When** failure is detected
**Then** job is automatically re-queued if retry count < 3
**And** retry attempt increments (`retry_count` field in jobs table)
**And** exponential backoff is applied: 1 min, 5 min, 15 min delays
**And** after 3 failed attempts, job status becomes "failed" permanently
**And** job details show: "Retry 2/3" badge
**And** auto-retry can be disabled in settings


## Epic 6: Subscriptions & Calendar Discovery

Permettre aux utilisateurs de s'abonner à des animés en cours. Les épisodes sont téléchargés automatiquement à leur sortie.

### Story 6.1: Implement Anime Calendar View

As a user,
I want to see a calendar of anime release dates,
So that I can discover new episodes.

**Acceptance Criteria:**

**Given** I navigate to the calendar page
**When** the page loads
**Then** I see a monthly calendar grid with today's date highlighted
**And** each date shows anime titles with new episodes releasing that day
**And** anime entries show: thumbnail, title, episode number, release time
**And** clicking an anime entry opens its detail page
**And** I can navigate between months with Previous/Next buttons
**And** calendar data is fetched from anime-sama.si release schedule

### Story 6.2: Implement Subscribe Button on Anime Detail Page

As a user,
I want to click "Subscribe" on an anime detail page,
So that new episodes are tracked.

**Acceptance Criteria:**

**Given** I'm viewing an anime detail page
**When** I click the "Subscribe" button
**Then** button changes to "Subscribed" with checkmark icon
**And** a subscription scope modal appears with options: "All Future Episodes", "Current Season Only", "Custom Range"
**And** selecting "All Future Episodes" creates subscription with scope = "all"
**And** subscription is persisted to `subscriptions` table with: anime_id, user_id, scope, created_at, last_checked
**And** subscribed anime appears in "My Subscriptions" page
**And** button state persists on page refresh

### Story 6.3: Implement Subscription Scope Selection

As a user,
I want to choose subscription scope (all episodes, season, range),
So that I only auto-download what I want.

**Acceptance Criteria:**

**Given** I'm subscribing to an anime
**When** the subscription modal appears
**Then** I see three options: "All Future Episodes", "Current Season Only", "Custom Episode Range (e.g., 1-24)"
**And** selecting "Current Season Only" creates subscription with scope = "season:{season_number}"
**And** selecting "Custom Range" shows two input fields for start/end episodes
**And** invalid ranges (e.g., "50-99" for 24-episode anime) show error: "Invalid range for this anime"
**And** scope can be edited later from "My Subscriptions" page

### Story 6.4: Implement Subscription Manager Page

As a user,
I want to view all my subscriptions in one place,
So that I can manage them.

**Acceptance Criteria:**

**Given** I navigate to "My Subscriptions" page
**When** the page loads
**Then** I see a list of subscribed anime with: thumbnail, title, scope, last checked time, status (active/paused)
**And** each subscription has "Edit Scope", "Pause", "Unsubscribe" buttons
**And** clicking "Pause" stops auto-downloads but keeps subscription
**And** clicking "Unsubscribe" removes subscription after confirmation
**And** clicking "Edit Scope" reopens the scope modal
**And** page shows empty state with "No subscriptions" message if list is empty

### Story 6.5: Implement Auto-Download Scheduler for Subscriptions

As a system,
I want to check for new episodes every 30 minutes,
So that subscribed content is downloaded automatically.

**Acceptance Criteria:**

**Given** subscriptions exist in the database
**When** scheduler runs (every 30 minutes)
**Then** anime-sama.si is queried for new episodes for each subscribed anime
**And** if new episode matches subscription scope, a download job is created automatically
**And** job is added with status "queued" and tagged with `subscription_id`
**And** user receives notification: "{Anime Title} Episode {N} added to queue"
**And** last_checked timestamp is updated for each subscription
**And** scheduler respects paused subscriptions (skips them)

### Story 6.6: Implement Subscription Table in Database

As a developer,
I want a subscriptions table in SQLite,
So that subscription data is persisted.

**Acceptance Criteria:**

**Given** the application starts
**When** database initialization runs
**Then** a `subscriptions` table is created with schema:
  - `id` (PRIMARY KEY)
  - `user_id` (INTEGER, for multi-user support, defaults to 1)
  - `anime_id` (TEXT, anime identifier from anime-sama.si)
  - `scope` (TEXT, e.g., "all", "season:2", "range:1-24")
  - `status` (TEXT, "active" or "paused")
  - `created_at` (DATETIME)
  - `last_checked` (DATETIME)
**And** UNIQUE constraint exists on (user_id, anime_id)
**And** table indexes are created on: user_id, status, last_checked

### Story 6.7: Implement Subscription Notifications

As a user,
I want to be notified when new episodes are auto-downloaded,
So that I know content is available.

**Acceptance Criteria:**

**Given** a subscription auto-downloads a new episode
**When** the job is created by the scheduler
**Then** I receive a browser notification: "{Anime Title} Episode {N} downloading"
**And** notification includes anime thumbnail
**And** clicking notification navigates to the queue page
**And** in-app badge on "Queue" menu item shows count of new subscription jobs
**And** notifications can be disabled in settings (key: `notify_subscriptions`)


## Epic 7: Public API & CLI Access

Exposer une API REST publique et un CLI pour intégrations externes. Les utilisateurs peuvent automatiser workflows via scripts.

### Story 7.1: Implement API Authentication with API Keys

As a developer,
I want to generate API keys for authentication,
So that external clients can access the API securely.

**Acceptance Criteria:**

**Given** I'm logged into the web app
**When** I navigate to Settings > API Keys
**Then** I see a button "Generate New API Key"
**And** clicking generates a new key (UUID format: `ask_xxxxxxxxxxxxxxxx`)
**And** key is stored in `api_keys` table with: id, key, user_id, created_at, last_used
**And** key is displayed once with warning: "Save this key, it won't be shown again"
**And** API requests include header: `Authorization: Bearer ask_xxxxxxxxxxxxxxxx`
**And** invalid keys return HTTP 401 with message: "Invalid API key"

### Story 7.2: Implement API Rate Limiting per Key

As a system,
I want to rate limit API requests per key,
So that abuse is prevented.

**Acceptance Criteria:**

**Given** an API key exists
**When** requests are made using that key
**Then** rate limit is 300 requests per minute per key
**And** excess requests return HTTP 429 with headers: `X-RateLimit-Limit: 300`, `X-RateLimit-Remaining: 0`, `Retry-After: 42`
**And** rate limit resets every 60 seconds (sliding window)
**And** rate limits are tracked in-memory (Redis-like structure)
**And** burst allowance is 30 requests

### Story 7.3: Implement RESTful API Endpoints

As a developer,
I want comprehensive REST endpoints,
So that I can build integrations.

**Acceptance Criteria:**

**Given** I have a valid API key
**When** I send requests to documented endpoints
**Then** the following routes are available:
  - `GET /api/v1/search?q=naruto` - Search anime
  - `GET /api/v1/anime/:id` - Get anime details
  - `POST /api/v1/jobs` - Create download job (body: {anime_id, episodes})
  - `GET /api/v1/jobs` - List all jobs
  - `GET /api/v1/jobs/:id` - Get job details
  - `DELETE /api/v1/jobs/:id` - Cancel/delete job
  - `GET /api/v1/subscriptions` - List subscriptions
  - `POST /api/v1/subscriptions` - Create subscription
  - `DELETE /api/v1/subscriptions/:id` - Delete subscription
**And** all responses use JSON format
**And** all endpoints respect authentication and rate limits

### Story 7.4: Implement API Versioning

As a system architect,
I want API versioning in URLs and headers,
So that breaking changes don't affect existing clients.

**Acceptance Criteria:**

**Given** the API is evolving
**When** new versions are released
**Then** version is specified in URL: `/api/v1/`, `/api/v2/`
**And** clients can optionally use header: `X-API-Version: 2`
**And** if no version specified, latest version is used
**And** deprecated versions return warning header: `X-API-Deprecated: true`
**And** API documentation specifies version for each endpoint

### Story 7.5: Implement OpenAPI Specification

As a developer,
I want an OpenAPI (Swagger) spec,
So that I can generate clients and explore the API.

**Acceptance Criteria:**

**Given** the API is running
**When** I access `/api/v1/openapi.json`
**Then** a valid OpenAPI 3.0 spec is returned
**And** spec includes: all endpoints, request/response schemas, authentication requirements
**And** Swagger UI is available at `/api/docs` for interactive exploration
**And** spec is auto-generated from Go code annotations
**And** spec includes examples for each endpoint

### Story 7.6: Implement CLI Tool

As a user,
I want a command-line tool to interact with the downloader,
So that I can script workflows.

**Acceptance Criteria:**

**Given** the CLI binary `asd` is installed
**When** I run CLI commands
**Then** the following commands work:
  - `asd search "naruto"` - Search anime
  - `asd download <anime_id> --episodes 1-12` - Create job
  - `asd queue list` - List jobs
  - `asd queue cancel <job_id>` - Cancel job
  - `asd subscribe <anime_id> --scope all` - Subscribe
  - `asd config set download_path /path/to/folder` - Change settings
**And** CLI reads API key from env var `ASD_API_KEY` or config file `~/.asd/config.yaml`
**And** CLI outputs JSON when `--json` flag is used
**And** CLI has `--help` for each command

### Story 7.7: Implement Webhook Notifications

As a developer,
I want to register webhooks for events,
So that my application is notified of changes.

**Acceptance Criteria:**

**Given** I have a valid API key
**When** I register a webhook via `POST /api/v1/webhooks` with payload: `{url, events: ["job.completed", "subscription.new_episode"]}`
**Then** webhook is stored in `webhooks` table
**And** when a job completes, a POST request is sent to webhook URL with payload: `{event: "job.completed", job_id, anime_title}`
**And** webhook requests include signature header: `X-Webhook-Signature: sha256(body + secret)`
**And** failed webhooks are retried 3 times with exponential backoff
**And** webhooks can be deleted via `DELETE /api/v1/webhooks/:id`

### Story 7.8: Implement API Documentation Portal

As a developer,
I want comprehensive API documentation,
So that I can integrate without reading code.

**Acceptance Criteria:**

**Given** the API is running
**When** I access `/docs`
**Then** I see a documentation portal with:
  - Getting started guide (authentication, first request)
  - Endpoint reference (all routes with examples)
  - SDKs section (links to Go, Python, JavaScript clients)
  - Rate limiting details
  - Error codes reference
  - Webhook guide with signature verification example
**And** documentation is searchable
**And** code examples are available in multiple languages


## Epic 8: Multi-User Support

Permettre plusieurs utilisateurs d'utiliser l'application avec comptes séparés. Chaque utilisateur a ses propres queues, subscriptions, et settings.

### Story 8.1: Implement User Table and Authentication Schema

As a developer,
I want a users table in the database,
So that user accounts can be stored.

**Acceptance Criteria:**

**Given** the application starts
**When** database initialization runs
**Then** a `users` table is created with schema:
  - `id` (PRIMARY KEY, INTEGER)
  - `username` (TEXT, UNIQUE, NOT NULL)
  - `email` (TEXT, UNIQUE, NOT NULL)
  - `password_hash` (TEXT, NOT NULL)
  - `created_at` (DATETIME)
  - `last_login` (DATETIME)
**And** a default admin user is created: username="admin", password="changeme"
**And** passwords are hashed using bcrypt with cost factor 12
**And** indexes exist on: username, email

### Story 8.2: Implement User Registration and Login UI

As a user,
I want to create an account and log in,
So that I can use the application with my own data.

**Acceptance Criteria:**

**Given** I'm on the login page
**When** I enter credentials and click "Login"
**Then** if valid, I'm redirected to the dashboard with session cookie set
**And** if invalid, error message appears: "Invalid username or password"
**And** I can click "Create Account" to access registration form
**And** registration requires: username (3-20 chars), email (valid format), password (8+ chars)
**And** duplicate usernames/emails show error: "Username already exists"
**And** session expires after 7 days of inactivity

### Story 8.3: Implement Session Management with JWT

As a system,
I want to use JWT tokens for session management,
So that authentication is stateless and scalable.

**Acceptance Criteria:**

**Given** a user logs in successfully
**When** authentication succeeds
**Then** a JWT token is generated with payload: {user_id, username, exp: 7 days}
**And** token is signed with secret from env var `JWT_SECRET`
**And** token is set as HTTP-only cookie: `session_token`
**And** all authenticated API requests verify JWT signature
**And** expired tokens return HTTP 401 with message: "Session expired"
**And** token refresh endpoint exists: `POST /api/v1/auth/refresh`

### Story 8.4: Implement User-Scoped Data Isolation

As a system,
I want all user data isolated by user_id,
So that users only see their own data.

**Acceptance Criteria:**

**Given** multiple users exist
**When** a user queries jobs, subscriptions, or settings
**Then** only rows with matching `user_id` are returned
**And** all database queries include `WHERE user_id = ?` clause
**And** attempting to access another user's data returns HTTP 403: "Forbidden"
**And** database migrations add `user_id` column to: jobs, subscriptions, api_keys tables
**And** foreign key constraints enforce referential integrity: `FOREIGN KEY (user_id) REFERENCES users(id)`

### Story 8.5: Implement User Profile Settings

As a user,
I want to edit my profile (email, password),
So that I can keep my account updated.

**Acceptance Criteria:**

**Given** I'm logged in
**When** I navigate to Profile settings
**Then** I see my current: username (read-only), email (editable), "Change Password" button
**And** clicking "Change Password" shows form: Current Password, New Password, Confirm New Password
**And** changing password requires correct current password
**And** new password must meet requirements: 8+ chars, 1 uppercase, 1 number
**And** successful changes show toast: "Profile updated"
**And** email changes send confirmation email to new address

### Story 8.6: Implement Admin User Management

As an admin,
I want to manage user accounts,
So that I can add/remove users and reset passwords.

**Acceptance Criteria:**

**Given** I'm logged in as admin (role = "admin")
**When** I navigate to Admin > Users
**Then** I see a table of all users: username, email, created_at, last_login, role
**And** I can click "Create User" to add a new account
**And** I can click "Delete" on any user (with confirmation dialog)
**And** I can click "Reset Password" to generate temporary password
**And** I can change user roles: "user" or "admin"
**And** admin role is stored in `users.role` column (TEXT, default "user")
**And** non-admin users cannot access `/admin/*` routes (HTTP 403)

### Story 8.7: Implement User Statistics Dashboard

As a user,
I want to see statistics about my usage,
So that I can track activity.

**Acceptance Criteria:**

**Given** I'm logged in
**When** I view the dashboard homepage
**Then** I see statistics cards:
  - Total downloads (count of completed jobs)
  - Active subscriptions (count)
  - Total storage used (sum of file sizes)
  - Downloads this month (count)
**And** statistics are calculated from database queries with `WHERE user_id = ?`
**And** storage calculation scans `data/videos/{user_id}/` directory
**And** statistics update every time page is visited (no caching)

