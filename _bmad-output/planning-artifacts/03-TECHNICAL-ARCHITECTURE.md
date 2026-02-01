# 🏗️ ARCHITECTURE TECHNIQUE

**Projet** : Anime-Sama Downloader v1.0  
**Date** : 31 janvier 2026  
**Auteur** : Winston (Architect)

---

## 📊 VUE D'ENSEMBLE

Architecture **Clean Architecture** avec séparation claire entre Backend (Go) et Frontend (React). Communication via API REST + SSE pour real-time updates.

```
┌────────────────────────────────────────────────┐
│              FRONTEND (React SPA)              │
│  - Pages (Dashboard, Subs, Calendar, Jobs)    │
│  - Components (UI, Anime, Layout)             │
│  - State Management (Zustand)                 │
│  - API Client + SSE                           │
└───────────────┬────────────────────────────────┘
                │ HTTP REST + WebSocket (SSE)
┌───────────────▼────────────────────────────────┐
│           BACKEND (Go HTTP Server)             │
│  ┌──────────────────────────────────────────┐ │
│  │        API Layer (Chi Router)            │ │
│  │  /api/v1/* + /auth/* + /ws + OpenAPI     │ │
│  └────────────┬─────────────────────────────┘ │
│               │                                │
│  ┌────────────▼─────────────────────────────┐ │
│  │      Application Services Layer          │ │
│  │  - SubscriptionService                   │ │
│  │  - JobService (Worker Pool + Queue)      │ │
│  │  - AniListSyncService                    │ │
│  │  - AnimeSamaScraperService               │ │
│  │  - JellyfinWebhookService                │ │
│  │  - UserService (multi-users)             │ │
│  │  - SettingsService                       │ │
│  └────────────┬─────────────────────────────┘ │
│               │                                │
│  ┌────────────▼─────────────────────────────┐ │
│  │         Domain Layer (Entities)          │ │
│  │  Subscription, Job, User, Settings       │ │
│  └────────────┬─────────────────────────────┘ │
│               │                                │
│  ┌────────────▼─────────────────────────────┐ │
│  │       Ports (Interfaces)                 │ │
│  │  JobRepo, SubRepo, UserRepo, EventBus    │ │
│  └────────────┬─────────────────────────────┘ │
│               │                                │
│  ┌────────────▼─────────────────────────────┐ │
│  │           Adapters                       │ │
│  │  - SQLiteAdapter (repos)                 │ │
│  │  - AnimeSamaClient (HTTP scraper)        │ │
│  │  - AniListClient (GraphQL)               │ │
│  │  - JellyfinClient (REST API)             │ │
│  │  - EventBusAdapter (SSE)                 │ │
│  └──────────────────────────────────────────┘ │
└────────────────────────────────────────────────┘
                │
┌───────────────▼────────────────────────────────┐
│          EXTERNAL SYSTEMS                      │
│  - anime-sama.si (scraping)                    │
│  - AniList GraphQL API                         │
│  - Jellyfin/Plex REST API                      │
│  - ffmpeg (HLS muxing)                         │
└────────────────────────────────────────────────┘
```

---

## 🔧 BACKEND (Go)

### Stack technique

- **Language** : Go 1.22+
- **Router** : `chi/v5` (léger, performant)
- **Database** : SQLite (`modernc.org/sqlite`)
- **Logging** : `zerolog` (structured logging)
- **Testing** : `testing` stdlib + table-driven tests
- **OpenAPI** : Schema v3 pour documentation API

### Structure de dossiers

```
cmd/
├── asd/              # CLI (legacy, optionnel)
└── asd-server/       # Serveur HTTP principal
    └── main.go

internal/
├── domain/           # Entités métier (Subscription, Job, User...)
│   ├── subscription.go
│   ├── job.go
│   ├── user.go
│   └── settings.go
│
├── ports/            # Interfaces (contrats)
│   ├── repositories.go   # JobRepo, SubRepo, UserRepo...
│   ├── services.go       # External APIs (AniList, Jellyfin)
│   └── eventbus.go       # EventBus interface
│
├── app/              # Application services (business logic)
│   ├── subscription_service.go
│   ├── job_service.go
│   ├── anilist_sync_service.go
│   ├── animesama_scraper_service.go
│   ├── jellyfin_webhook_service.go
│   ├── user_service.go
│   └── worker_pool.go
│
├── adapters/         # Implémentations concrètes
│   ├── sqlite/          # Persistence SQLite
│   │   ├── subscription_repo.go
│   │   ├── job_repo.go
│   │   └── user_repo.go
│   ├── httpapi/         # Handlers HTTP
│   │   ├── subscriptions.go
│   │   ├── jobs.go
│   │   ├── auth.go
│   │   └── openapi.go
│   ├── animesama/       # Client HTTP anime-sama
│   │   └── scraper.go
│   ├── anilist/         # Client GraphQL AniList
│   │   └── client.go
│   ├── jellyfin/        # Client REST Jellyfin
│   │   └── client.go
│   └── memorybus/       # EventBus in-memory + SSE
│       └── eventbus.go
│
├── config/           # Configuration
│   └── config.go
│
└── buildinfo/        # Version, build info
    └── buildinfo.go
```

### Services principaux

#### SubscriptionService

**Responsabilités** :
- CRUD abonnements
- Auto-labeling depuis baseURL
- Sync avec AniList watchlist
- Calcul `nextCheckAt` (scheduler)

**Méthodes clés** :
```go
type SubscriptionService interface {
    Create(ctx context.Context, baseURL, player string) (*Subscription, error)
    Update(ctx context.Context, id string, updates map[string]interface{}) error
    Delete(ctx context.Context, id string) error
    List(ctx context.Context, limit int) ([]Subscription, error)
    GetEpisodes(ctx context.Context, baseURL string) ([]Episode, error)
    EnqueueEpisodes(ctx context.Context, subID string, episodeNumbers []int) error
    SyncFromAniList(ctx context.Context, userID string) ([]Subscription, error)
}
```

#### JobService

**Responsabilités** :
- CRUD jobs (queued, running, completed, failed)
- Worker pool (N workers concurrents)
- SSE events (`job.created`, `job.progress`, `job.completed`)
- Retry logic (max 3 attempts)

**Méthodes clés** :
```go
type JobService interface {
    Create(ctx context.Context, typ string, params map[string]interface{}) (*Job, error)
    Cancel(ctx context.Context, jobID string) error
    List(ctx context.Context, filters map[string]interface{}) ([]Job, error)
    Get(ctx context.Context, jobID string) (*Job, error)
    UpdateProgress(ctx context.Context, jobID string, progress float64) error
}
```

#### AniListSyncService

**Responsabilités** :
- OAuth2 flow (token exchange)
- Fetch viewer watchlist
- Match anime → anime-sama candidates
- Auto-create subscriptions

**Flow** :
1. User clicks "Connect AniList"
2. OAuth redirect → get token
3. Fetch watchlist (`CURRENT`, `PLANNING`)
4. Pour chaque anime : resolve candidates anime-sama
5. User valide matches → create subscriptions

#### JellyfinWebhookService (NEW)

**Responsabilités** :
- Webhook POST vers Jellyfin après download
- Trigger library scan
- Update metadata (AniList/AniDB)

**Configuration** :
```yaml
jellyfin:
  enabled: true
  url: http://jellyfin:8096
  api_key: YOUR_API_KEY
  library_id: anime_library
```

---

## 🎨 FRONTEND (React)

### Stack technique

- **Framework** : React 18 + TypeScript
- **Bundler** : Vite 5
- **Router** : React Router v6
- **State** : Zustand (lightweight, simple)
- **Styling** : CSS Modules + Design tokens
- **API** : Axios + React Query (caching)
- **Icons** : Lucide React
- **Animations** : Framer Motion

### Structure de dossiers

```
webapp/
├── src/
│   ├── main.tsx                    # Entry point
│   │
│   ├── pages/                      # Route-based pages
│   │   ├── Dashboard.tsx           # Home (nouveautés + abonnements)
│   │   ├── Search.tsx              # Recherche anime
│   │   ├── Subscriptions.tsx       # Liste abonnements
│   │   ├── Calendar.tsx            # Calendrier sorties
│   │   ├── Jobs.tsx                # Queue téléchargements
│   │   └── Settings.tsx            # Configuration
│   │
│   ├── components/                 # Atomic Design
│   │   ├── ui/                     # Primitives (Button, Input, Card...)
│   │   │   ├── Button.tsx
│   │   │   ├── Input.tsx
│   │   │   ├── Card.tsx
│   │   │   ├── Modal.tsx
│   │   │   ├── Badge.tsx
│   │   │   ├── ProgressBar.tsx
│   │   │   └── Toast.tsx
│   │   │
│   │   ├── anime/                  # Anime-specific components
│   │   │   ├── AnimeCard.tsx
│   │   │   ├── EpisodeList.tsx
│   │   │   ├── EpisodeSelector.tsx
│   │   │   └── SubscriptionCard.tsx
│   │   │
│   │   └── layout/                 # Layout components
│   │       ├── Header.tsx
│   │       ├── Sidebar.tsx
│   │       ├── TabNav.tsx
│   │       └── PageContainer.tsx
│   │
│   ├── hooks/                      # Custom React hooks
│   │   ├── useSubscriptions.ts     # Fetch + cache subscriptions
│   │   ├── useJobs.ts              # Fetch + cache jobs
│   │   ├── useSSE.ts               # SSE real-time connection
│   │   ├── useSearch.ts            # Anime search
│   │   └── useAuth.ts              # Authentication
│   │
│   ├── stores/                     # Zustand stores
│   │   ├── subscriptionsStore.ts
│   │   ├── jobsStore.ts
│   │   ├── authStore.ts
│   │   └── settingsStore.ts
│   │
│   ├── lib/                        # Utils & helpers
│   │   ├── api.ts                  # API client (axios)
│   │   ├── sse.ts                  # SSE client
│   │   ├── theme.ts                # Theme utils
│   │   └── utils.ts                # Misc helpers
│   │
│   └── styles/
│       ├── tokens.css              # Design tokens (variables)
│       ├── globals.css             # Reset + base styles
│       ├── components.css          # UI components styles
│       └── animations.css          # Keyframes
│
├── index.html
├── vite.config.ts
├── tsconfig.json
└── package.json
```

### State Management (Zustand)

```typescript
// subscriptionsStore.ts
import { create } from 'zustand';

interface SubscriptionsStore {
  subscriptions: Subscription[];
  loading: boolean;
  error: string | null;
  
  fetch: () => Promise<void>;
  create: (baseURL: string) => Promise<void>;
  delete: (id: string) => Promise<void>;
  sync: (id: string) => Promise<void>;
}

export const useSubscriptionsStore = create<SubscriptionsStore>((set, get) => ({
  subscriptions: [],
  loading: false,
  error: null,
  
  fetch: async () => {
    set({ loading: true });
    try {
      const data = await api.listSubscriptions();
      set({ subscriptions: data, loading: false });
    } catch (err) {
      set({ error: err.message, loading: false });
    }
  },
  
  // ... autres méthodes
}));
```

### SSE Real-time Updates

```typescript
// useSSE.ts
export function useSSE() {
  const jobsStore = useJobsStore();
  
  useEffect(() => {
    const eventSource = new EventSource('/api/v1/events');
    
    eventSource.addEventListener('job.progress', (e) => {
      const data = JSON.parse(e.data);
      jobsStore.updateJob(data.jobId, { progress: data.progress });
    });
    
    eventSource.addEventListener('job.completed', (e) => {
      const data = JSON.parse(e.data);
      jobsStore.updateJob(data.jobId, { state: 'completed' });
      toast.success(`Episode ${data.episode} téléchargé !`);
    });
    
    return () => eventSource.close();
  }, []);
}
```

---

## 🔌 API ENDPOINTS

### Subscriptions

```
GET    /api/v1/subscriptions              # List all
POST   /api/v1/subscriptions              # Create
GET    /api/v1/subscriptions/:id          # Get one
PUT    /api/v1/subscriptions/:id          # Update
DELETE /api/v1/subscriptions/:id          # Delete
POST   /api/v1/subscriptions/:id/sync     # Sync check
POST   /api/v1/subscriptions/:id/enqueue  # Enqueue episodes
GET    /api/v1/subscriptions/:id/episodes # List episodes
```

### Jobs

```
GET    /api/v1/jobs                       # List all
POST   /api/v1/jobs                       # Create (enqueue)
GET    /api/v1/jobs/:id                   # Get one
DELETE /api/v1/jobs/:id                   # Cancel
POST   /api/v1/jobs/sync-all              # Sync all subscriptions
```

### AniList

```
GET    /api/v1/anilist/viewer             # Get viewer info
GET    /api/v1/anilist/airing             # Airing schedule
POST   /api/v1/anilist/import-auto        # Auto-import watchlist
```

### Anime-Sama

```
POST   /api/v1/animesama/resolve          # Resolve anime title → candidates
POST   /api/v1/animesama/scan             # Scan options (seasons, langs)
GET    /api/v1/animesama/episodes         # List episodes for baseURL
POST   /api/v1/animesama/enqueue          # Enqueue download
```

### Settings

```
GET    /api/v1/settings                   # Get settings
PUT    /api/v1/settings                   # Update settings
```

### Auth (NEW)

```
POST   /api/v1/auth/login                 # Login (JWT)
POST   /api/v1/auth/logout                # Logout
POST   /api/v1/auth/refresh               # Refresh token
GET    /api/v1/auth/me                    # Current user
```

### SSE

```
GET    /api/v1/events                     # SSE stream
```

---

## 🗄️ DATABASE SCHEMA (SQLite)

```sql
-- Users (NEW)
CREATE TABLE users (
  id TEXT PRIMARY KEY,
  username TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  role TEXT NOT NULL DEFAULT 'user', -- 'admin' | 'user'
  anilist_token TEXT,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
);

-- Subscriptions
CREATE TABLE subscriptions (
  id TEXT PRIMARY KEY,
  user_id TEXT REFERENCES users(id),
  base_url TEXT NOT NULL,
  label TEXT NOT NULL,
  player TEXT NOT NULL,
  last_scheduled_episode INTEGER NOT NULL DEFAULT 0,
  last_downloaded_episode INTEGER NOT NULL DEFAULT 0,
  last_available_episode INTEGER NOT NULL DEFAULT 0,
  next_check_at DATETIME NOT NULL,
  last_checked_at DATETIME,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
);

-- Jobs
CREATE TABLE jobs (
  id TEXT PRIMARY KEY,
  user_id TEXT REFERENCES users(id),
  type TEXT NOT NULL, -- 'download' | 'spawn' | 'sync'
  state TEXT NOT NULL, -- 'queued' | 'running' | 'completed' | 'failed' | 'canceled'
  progress REAL NOT NULL DEFAULT 0.0,
  params TEXT, -- JSON
  result TEXT, -- JSON
  error TEXT,
  attempts INTEGER NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  started_at DATETIME,
  finished_at DATETIME
);

-- Settings (global)
CREATE TABLE settings (
  key TEXT PRIMARY KEY,
  value TEXT NOT NULL,
  updated_at DATETIME NOT NULL
);
```

---

## 🔐 SÉCURITÉ

### Authentication (JWT)

- Tokens JWT signés avec secret (env var `JWT_SECRET`)
- Expiration : 1 heure (access token), 30 jours (refresh token)
- Stored in httpOnly cookies (XSS protection)

### Authorization

- Role-based : `admin` | `user`
- Admins : Full access
- Users : Own subscriptions + jobs only

### Input Validation

- Tous les inputs validés côté backend
- Protection path traversal (downloads)
- Sanitization HTML (logs, titles)

### Rate Limiting

- Anime-sama scraping : Max 1 req/sec (éviter ban)
- API publique : Max 100 req/min par IP
- AniList : Respect rate limits (90 points/min)

---

## 📊 PERFORMANCE

### Backend

- Worker pool : 10 workers max (configurable)
- Database : SQLite avec WAL mode (concurrent reads)
- Caching : In-memory cache pour episodes list (5 min TTL)
- Graceful shutdown : Context cancellation + drain workers

### Frontend

- Code splitting : React.lazy() pour pages
- Image optimization : WebP + lazy loading
- Bundle size : < 300KB (gzipped)
- Lighthouse score target : ≥ 90

---

## 🔧 CONFIGURATION

### Environment Variables

```bash
# Server
ASD_HOST=127.0.0.1
ASD_PORT=8080
ASD_DB_PATH=asd.db

# Paths
ASD_DESTINATION=/data/videos
ASD_MAX_CONCURRENT_DOWNLOADS=10

# Auth
JWT_SECRET=your_secret_key

# Jellyfin (optionnel)
JELLYFIN_ENABLED=true
JELLYFIN_URL=http://jellyfin:8096
JELLYFIN_API_KEY=your_api_key
JELLYFIN_LIBRARY_ID=anime_library

# AniList (optionnel)
ANILIST_CLIENT_ID=your_client_id
ANILIST_CLIENT_SECRET=your_client_secret
```

---

## 🧪 TESTS

### Backend

- Unit tests : 70%+ coverage
- Integration tests : API endpoints + database
- Table-driven tests pour business logic

```bash
go test ./... -cover -race
```

### Frontend

- Unit tests : Components (React Testing Library)
- E2E tests : User flows (Playwright)

```bash
npm test
npm run test:e2e
```

---

## 🚀 DÉPLOIEMENT

### Docker Compose (Production)

```yaml
version: '3.8'
services:
  asd:
    build: .
    ports:
      - "8080:8080"
    volumes:
      - ./data:/data
      - ./videos:/videos
    environment:
      - ASD_DESTINATION=/videos
      - JELLYFIN_URL=http://jellyfin:8096
```

### CI/CD (GitHub Actions)

- Lint (golangci-lint, ESLint)
- Tests (Go + React)
- Build Docker image
- Deploy (optionnel)

---

**🎯 Prochaines étapes** : Implémenter les nouveaux services (UserService, JellyfinWebhookService) dans Sprint 3.
