---
title: "Implementation Plan — Anime-Sama Downloader v1.0"
date: "31 janvier 2026"
version: "1.0"
status: "Ready for Development"
timeline: "Week 1-3"
---

# Implementation Plan — Anime-Sama Downloader v1.0

_Detailed task breakdown following the Architecture Decision Document_

---

## 📊 Project Timeline

**Week 1 (Days 1-5):** MVP Foundation + Prototype Validation
**Week 2-3 (Days 6-15):** Features + Polish + Testing
**Post-MVP:** Monitoring, docs, optimizations

---

## 🏗️ Phase 1: MVP Foundation (Week 1, Days 1-5)

### Day 1: Project Setup & Database

#### Task 1.1: Backend Project Structure
- [ ] `go mod init anime-sama/asd` (if fresh)
- [ ] Create `internal/domain/` package structure
- [ ] Create `internal/app/` package structure
- [ ] Create `internal/adapters/` subpackages
- [ ] Create `internal/config/` package
- [ ] Create `cmd/asd-server/` entry point
- **Depends on:** None
- **Time estimate:** 30 min
- **Files to create:** ~15 .go files (empty with package declarations)

#### Task 1.2: Frontend Project Structure
- [ ] `npm install` (if fresh or update)
- [ ] Create `src/components/ui/` directory
- [ ] Create `src/components/custom/` directory
- [ ] Create `src/components/layouts/` directory
- [ ] Create `src/components/features/` directory
- [ ] Create `src/stores/` directory
- [ ] Create `src/hooks/` directory
- [ ] Create `src/utils/` directory
- [ ] Create `src/types/` directory
- [ ] Create `src/styles/` directory
- **Depends on:** None
- **Time estimate:** 20 min
- **Files to create:** ~10 directories

#### Task 1.3: SQLite Database Setup
- [ ] Create `internal/adapters/sqlite/db.go` with connection logic
- [ ] Create `internal/adapters/sqlite/migrations.go` with migration runner
- [ ] Create `migrations/001_init.sql` with schema:
  - [ ] `downloads` table
  - [ ] `jobs` table
  - [ ] `settings` table
  - [ ] `subscriptions` table (post-MVP placeholder)
- [ ] Create `internal/adapters/sqlite/migrations/001_init.sql`
- [ ] Test connection in main.go
- **Depends on:** Task 1.1
- **Time estimate:** 1 hour
- **Schema spec:**
  ```sql
  CREATE TABLE downloads (
    download_id TEXT PRIMARY KEY,
    job_id TEXT,
    anime_id TEXT,
    episode_number INTEGER,
    metadata JSON,
    created_at TEXT,
    FOREIGN KEY(job_id) REFERENCES jobs(job_id)
  );
  
  CREATE TABLE jobs (
    job_id TEXT PRIMARY KEY,
    status TEXT,
    progress_percent INTEGER,
    error_message TEXT,
    created_at TEXT,
    updated_at TEXT
  );
  
  CREATE TABLE settings (
    key TEXT PRIMARY KEY,
    value TEXT,
    updated_at TEXT
  );
  ```

#### Task 1.4: Docker Setup
- [ ] Verify `Dockerfile` exists (from project)
- [ ] Verify `docker-compose.yml` exists
- [ ] Test `docker compose up` builds successfully
- [ ] Expose ports 8000 (backend) and 5173 (frontend)
- **Depends on:** Task 1.1, 1.2
- **Time estimate:** 20 min

**Day 1 Goal:** ✅ Project structure + database ready

---

### Day 2: Domain Layer & Entities

#### Task 2.1: Domain Entities
- [ ] `internal/domain/download.go` — Download entity struct
  - Fields: ID, JobID, AnimeID, EpisodeNum, Metadata, CreatedAt
  - Methods: IsCompleted(), GetFilePath()
- [ ] `internal/domain/job.go` — Job entity struct
  - Fields: ID, Status, Progress, Error, CreatedAt, UpdatedAt
  - Methods: IsRunning(), IsCompleted(), HasError()
- [ ] `internal/domain/errors.go` — Custom error types
  - ErrorCode enum
  - AppError struct
- **Depends on:** Task 1.1
- **Time estimate:** 1 hour
- **Pattern:** Follow Go naming (PascalCase types, camelCase methods)

#### Task 2.2: Domain Interfaces
- [ ] `internal/domain/resolver.go` — IResolver interface
  ```go
  type IResolver interface {
    Resolve(ctx context.Context, query string) (*SearchResult, error)
  }
  ```
- [ ] `internal/domain/queue.go` — IJobQueue interface (if separate)
  ```go
  type IJobQueue interface {
    Enqueue(job *Job) error
    Dequeue() (*Job, error)
    UpdateStatus(jobID string, status string) error
  }
  ```
- **Depends on:** Task 2.1
- **Time estimate:** 30 min

#### Task 2.3: Mocks for Testing
- [ ] `internal/domain/mocks.go` — Mock implementations
  - MockResolver
  - MockJobQueue
- [ ] Set up testify/mock for assertions
- **Depends on:** Task 2.2
- **Time estimate:** 30 min

**Day 2 Goal:** ✅ Domain layer complete + mockable

---

### Day 3: Application Layer (Services)

#### Task 3.1: SearchService
- [ ] `internal/app/search_service.go`
  ```go
  type SearchService struct {
    resolvers []domain.IResolver
  }
  
  func (s *SearchService) SearchAnime(ctx, query string) (*domain.SearchResult, error)
  ```
- [ ] Implement multi-source search (AnimeSama → MangaDex fallback)
- [ ] Add error handling + logging
- [ ] Unit tests: `search_service_test.go`
  - Test successful search
  - Test fallback on primary failure
  - Test error handling
- **Depends on:** Task 2.2
- **Time estimate:** 2 hours

#### Task 3.2: DownloadService
- [ ] `internal/app/download_service.go`
  ```go
  type DownloadService struct {
    queue domain.IJobQueue
  }
  
  func (s *DownloadService) StartDownload(ctx, downloadID string) error
  ```
- [ ] Enqueue logic
- [ ] Validation
- [ ] Logging
- [ ] Unit tests: `download_service_test.go`
  - Test successful enqueue
  - Test validation
  - Test duplicate prevention
- **Depends on:** Task 2.2, Task 1.3
- **Time estimate:** 1.5 hours

#### Task 3.3: JobWorker
- [ ] `internal/app/job_worker.go`
  ```go
  type JobWorker struct {
    queue domain.IJobQueue
    bus domain.IEventBus
  }
  
  func (w *JobWorker) Start(ctx context.Context) error
  ```
- [ ] Goroutine loop to process jobs
- [ ] Update progress
- [ ] Emit events
- [ ] Basic tests (mock queue/bus)
- **Depends on:** Task 2.2, Task 3.4 (EventBus)
- **Time estimate:** 2 hours

#### Task 3.4: SettingsService
- [ ] `internal/app/settings_service.go`
  ```go
  type SettingsService struct {
    repo domain.ISettingsRepository
  }
  
  func (s *SettingsService) Get(key string) (string, error)
  func (s *SettingsService) Set(key, value string) error
  ```
- [ ] CRUD operations
- [ ] Caching (optional)
- [ ] Unit tests
- **Depends on:** Task 2.2, Task 1.3
- **Time estimate:** 1 hour

**Day 3 Goal:** ✅ All core services implemented + tested

---

### Day 4: Adapters (HTTP API + Event Bus)

#### Task 4.1: HTTP API Handlers
- [ ] `internal/adapters/httpapi/handlers.go` — Setup HTTP server
  - Router setup (chi or mux)
  - Middleware (CORS, logging, error handling)
- [ ] `internal/adapters/httpapi/download_handler.go`
  - `GET /api/downloads` — List
  - `POST /api/downloads` — Create
  - `GET /api/downloads/{id}` — Get one
  - `DELETE /api/downloads/{id}` — Cancel (future)
- [ ] `internal/adapters/httpapi/job_handler.go`
  - `GET /api/downloads/{id}/progress` (SSE) — Stream
- [ ] `internal/adapters/httpapi/settings_handler.go`
  - `GET /api/settings/{key}`
  - `PUT /api/settings/{key}`
- [ ] Error response formatter
- [ ] Tests for each endpoint
- **Depends on:** Task 3.1, 3.2, 3.3, 3.4
- **Time estimate:** 3 hours

#### Task 4.2: SQLite Repository Adapter
- [ ] `internal/adapters/sqlite/download_repo.go`
  - Implement domain.IDownloadRepository
  - CRUD operations
- [ ] `internal/adapters/sqlite/job_repo.go`
  - Implement domain.IJobRepository
  - CRUD operations
- [ ] `internal/adapters/sqlite/settings_repo.go`
  - Implement domain.ISettingsRepository
- [ ] Database tests (use test fixtures)
- **Depends on:** Task 1.3, Task 2.2
- **Time estimate:** 2 hours

#### Task 4.3: Event Bus
- [ ] `internal/adapters/memorybus/event_bus.go`
  ```go
  type EventBus struct {
    subscribers map[string][]Handler
    mu sync.RWMutex
  }
  
  func (b *EventBus) Subscribe(event string, handler Handler)
  func (b *EventBus) Emit(event string, payload interface{})
  ```
- [ ] Implement domain.IEventBus interface
- [ ] Thread-safe pub/sub
- [ ] Tests
- **Depends on:** Task 2.2
- **Time estimate:** 1 hour

#### Task 4.4: Main Entry Point
- [ ] `cmd/asd-server/main.go`
  - Initialize config
  - Setup database
  - Initialize services
  - Start HTTP server + job worker
  - Graceful shutdown
- [ ] Integration test: server starts + responds to ping
- **Depends on:** Task 1.1, 1.3, 4.1, 4.2, 4.3
- **Time estimate:** 1 hour

**Day 4 Goal:** ✅ Backend API fully functional + event bus working

---

### Day 5: Frontend (React) + Real-time Integration

#### Task 5.1: Zustand Stores
- [ ] `webapp/src/stores/ui.store.ts`
  - State: mode (simple/expert), currentView, modals
  - Actions: setMode, setView, openModal, closeModal
  - Persist mode to localStorage
- [ ] `webapp/src/stores/jobs.store.ts`
  - State: downloads array, selectedDownloadId, progress
  - Actions: addDownload, updateProgress, removeDownload
  - Subscribe to SSE events
- [ ] `webapp/src/stores/search.store.ts`
  - State: query, results, filters, isLoading
  - Actions: search, setFilters, clearResults
- [ ] Tests for store actions
- **Depends on:** Task 1.2
- **Time estimate:** 2 hours

#### Task 5.2: Custom React Components
- [ ] `webapp/src/components/custom/StatusBadge.tsx`
  - Display job status (queued, running, completed, failed)
  - Color-coded badges
- [ ] `webapp/src/components/custom/DownloadProgress.tsx`
  - Progress bar + percent
  - ETA display
  - Live update from SSE
- [ ] `webapp/src/components/custom/ModeToggle.tsx`
  - Simple ↔ Expert switch
  - Persist selection
- [ ] `webapp/src/components/custom/RuleCard.tsx` (future, placeholder)
- [ ] `webapp/src/components/custom/FormStepper.tsx` (future, placeholder)
- [ ] `webapp/src/components/custom/LogViewer.tsx` (future, placeholder)
- [ ] Each component: `.tsx` + `.test.tsx`
- **Depends on:** Task 1.2, Task 5.1
- **Time estimate:** 3 hours

#### Task 5.3: API Integration Utils
- [ ] `webapp/src/utils/api.ts`
  - Base fetch wrapper with error handling
  - API_BASE_URL configuration
  - JSON field mapping (snake_case ↔ camelCase)
- [ ] `webapp/src/types/api.ts`
  - Download interface
  - Job interface
  - API response types
- [ ] Error handling utils
- **Depends on:** Task 1.2
- **Time estimate:** 1 hour

#### Task 5.4: SSE Hook
- [ ] `webapp/src/hooks/useSSE.ts`
  ```typescript
  export const useSSE = (jobId: string, onProgress: (percent: number) => void) => {
    useEffect(() => {
      const es = new EventSource(`/api/downloads/${jobId}/progress`)
      es.onmessage = (e) => onProgress(JSON.parse(e.data).percent)
      return () => es.close()
    }, [jobId])
  }
  ```
- [ ] Error handling on disconnect
- [ ] Cleanup on unmount
- [ ] Test: mock EventSource
- **Depends on:** Task 1.2
- **Time estimate:** 1 hour

#### Task 5.5: Main Pages
- [ ] `webapp/src/App.tsx` — Main layout
  - Header + navigation
  - Render by currentView (search, downloads, settings)
- [ ] `webapp/src/components/features/search/SearchBar.tsx`
  - Input + debounce
  - Call useSearchStore.search()
- [ ] `webapp/src/components/features/downloads/DownloadList.tsx`
  - List all downloads
  - Use useJobsStore
  - Display StatusBadge + DownloadProgress
- [ ] `webapp/src/components/layouts/ErrorBoundary.tsx`
  - Catch React errors
  - Display user-friendly message
- [ ] App.tsx tests
- **Depends on:** Task 5.1, 5.2, 5.3, 5.4
- **Time estimate:** 2 hours

#### Task 5.6: Tailwind + Sakura Night
- [ ] Verify `tailwind.config.ts` has Sakura Night tokens
  - Colors: #0A0E1A (base), #D946EF (magenta), #06B6D4 (cyan)
  - Typography: Inter (body), Noto Serif JP (titles)
  - Spacing scale
- [ ] `webapp/src/styles/colors.css` — CSS variables
- [ ] `webapp/src/styles/globals.css` — Global styles
- **Depends on:** Task 1.2
- **Time estimate:** 30 min

**Day 5 Goal:** ✅ Frontend + backend integrated, real-time working

---

### Week 1 Checkpoint: MVP Foundation Complete ✅

**Running:**
```bash
# Terminal 1: Backend
cd /home/guilhem/Anime-Sama-Downloader
go run cmd/asd-server/main.go

# Terminal 2: Frontend
cd /home/guilhem/Anime-Sama-Downloader/webapp
npm run dev
```

**Test:**
- [ ] Open http://localhost:5173
- [ ] Click "Search" → API call to backend
- [ ] Click "Download" → Job queued + SSE progress visible
- [ ] Mode toggle works + persists
- [ ] All 5 custom components render
- [ ] No console errors

**Prototype User Testing Ready:**
- [ ] Figma/Framer prototypes created (parallel track)
- [ ] User testing sessions scheduled (Alex + Maya personas)
- [ ] Feedback collection method defined

---

## 🚀 Phase 2: Features + Polish (Week 2-3, Days 6-15)

### Week 2: Features Implementation

#### Task 6: Multi-source Search Reliability
- [ ] Implement AnimeSama resolver (`internal/adapters/animesama/`)
- [ ] Implement MangaDex resolver (`internal/adapters/mangadex/`)
- [ ] Metadata matching algorithm
- [ ] Retry logic (exponential backoff)
- [ ] Cache layer (optional)
- [ ] Tests: 10+ test cases per resolver
- **Time estimate:** 4 hours
- **Depends on:** Day 4 API complete

#### Task 7: Job Worker Improvements
- [ ] Implement 7-layer fallback strategy
- [ ] Proper error recovery
- [ ] Job persistence across restarts
- [ ] Rate limiting (if needed)
- [ ] Logging per job step
- **Time estimate:** 3 hours

#### Task 8: Settings & Configuration
- [ ] Jellyfin integration setup
- [ ] AniList token storage (encrypted)
- [ ] Settings UI page in React
- [ ] Save/load settings
- **Time estimate:** 3 hours

#### Task 9: Error Handling & Validation
- [ ] Input validation on all endpoints
- [ ] User-friendly error messages
- [ ] Error boundary in React
- [ ] Retry UI prompts
- **Time estimate:** 2 hours

### Week 3: Testing + Polish

#### Task 10: Testing & Coverage
- [ ] Unit tests: target 70%+ coverage
  - Go: `go test ./...`
  - React: `npm test`
- [ ] Integration tests (E2E scenarios)
- [ ] Manual testing checklist
- [ ] Performance testing (Lighthouse)
- **Time estimate:** 6 hours

#### Task 11: UI Polish
- [ ] Animations (transitions)
- [ ] Loading skeletons
- [ ] Empty states
- [ ] Error states
- [ ] Dark mode refinements
- [ ] Responsive design fixes
- **Time estimate:** 4 hours

#### Task 12: Documentation
- [ ] API documentation (OpenAPI)
- [ ] Component Storybook (optional)
- [ ] Development setup guide
- [ ] Architecture review doc
- **Time estimate:** 2 hours

---

## ✅ Task Dependencies Graph

```
Day 1:
  1.1 → 1.2 → 1.3 → 1.4

Day 2:
  1.1 → 2.1 → 2.2 → 2.3

Day 3:
  2.2 → 3.1, 3.2, 3.3, 3.4

Day 4:
  3.1-3.4, 1.3 → 4.1, 4.2, 4.3
  All → 4.4 (main.go)

Day 5:
  1.2 → 5.1 → 5.5
  1.2 → 5.2, 5.3, 5.4
  All → 5.5 (App.tsx)
  1.2 → 5.6 (Tailwind)

Week 2:
  Day 5 ✅ → 6, 7, 8, 9

Week 3:
  Week 2 ✅ → 10, 11, 12
```

---

## 📋 Daily Checklist Template

**Each day, track:**
- [ ] All tasks completed?
- [ ] Tests passing?
- [ ] No critical bugs?
- [ ] Code compiles/runs?
- [ ] Commits pushed?

---

## 🎯 Success Criteria (Week 1)

- [ ] Backend API responding to requests
- [ ] Frontend loading without errors
- [ ] Real-time SSE streaming working
- [ ] Mode toggle persistent
- [ ] Database persisting data
- [ ] 3+ manual test cases passing
- [ ] Prototype ready for user testing

---

## 📞 Questions During Implementation?

Refer back to **architecture.md** for:
- Pattern questions → Section: "Implementation Patterns"
- Structure questions → Section: "Project Structure"
- Decision questions → Section: "ADRs"

---

**Status:** Ready to start Task 1.1
**Next:** Begin backend setup
