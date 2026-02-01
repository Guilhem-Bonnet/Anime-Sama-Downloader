# Story 2.3: Implement Anime Detail View

**Story ID:** 2-3-implement-anime-detail-view  
**Story Points:** 5  
**Status:** ready-for-dev  
**Created:** 31 janvier 2025  
**Epic:** Epic 2 - Anime Search & Discovery

---

## 📖 Story

As a user,  
I want to view detailed information about an anime,  
so that I can decide whether to download it.

---

## ✅ Acceptance Criteria

1. [ ] **AC1** - Detail page loads when clicking anime from search results or entering /anime/:id URL
2. [ ] **AC2** - Page displays: title, thumbnail, synopsis, year, status, genre tags, episode count
3. [ ] **AC3** - "Download" button is visible and accessible
4. [ ] **AC4** - Episode list is displayed with season/episode numbers (e.g., "S01E12")
5. [ ] **AC5** - If anime has multiple seasons, seasons are displayed as tabs (Season 1, Season 2, etc.)
6. [ ] **AC6** - Page loads within 200ms (P95)
7. [ ] **AC7** - Loading state shown while fetching data
8. [ ] **AC8** - Error state shown if anime not found (404) or API error

---

## 🎯 Tasks / Subtasks

### Task 1: Create Anime Detail Domain Model & Service
- [ ] **1.1** Create `AnimeDetail` domain model extending AnimeSearchResult with: synopsis, genres, seasons, episodes
- [ ] **1.2** Define episode structure: `Episode { number, title, url, season_number }`
- [ ] **1.3** Define season structure: `Season { number, name, episodes[] }`
- [ ] **1.4** Create `AnimeDetailService` interface in ports with `GetDetail(ctx, id) (AnimeDetail, error)` method
- [ ] **1.5** Implement mock service returning test data (real scraping to be implemented in later stories)

### Task 2: Implement Anime Detail Backend Endpoint
- [ ] **2.1** Create `AnimeDetailHandler` in `internal/adapters/httpapi/anime_detail.go`
- [ ] **2.2** Register route `GET /api/v1/anime/:id` in chi router
- [ ] **2.3** Parse anime ID from URL parameter
- [ ] **2.4** Call AnimeDetailService.GetDetail(ctx, id)
- [ ] **2.5** Return 404 if anime not found
- [ ] **2.6** Return 200 with JSON detail object: `{id, title, thumbnail_url, synopsis, year, status, genres[], episode_count, seasons[]}`
- [ ] **2.7** Add Content-Type: application/json header

### Task 3: Create Frontend Anime Detail Page
- [ ] **3.1** Create `webapp/src/pages/AnimeDetailPage.tsx` component
- [ ] **3.2** Use React Router to extract :id param from URL
- [ ] **3.3** Fetch anime detail from `GET /api/v1/anime/:id` on mount
- [ ] **3.4** Show loading spinner while fetching (reuse existing loader component)
- [ ] **3.5** Show error message if 404 or API error (with retry button)
- [ ] **3.6** Display anime info: large thumbnail (300x400), title (H1), year, status badge
- [ ] **3.7** Display synopsis in scrollable text block (max-height 300px)
- [ ] **3.8** Display genre tags as pills (same style as status badges)

### Task 4: Implement Episode List with Season Tabs
- [ ] **4.1** Create `EpisodeList` component accepting seasons[] prop
- [ ] **4.2** If anime has 1 season, display episodes directly (no tabs)
- [ ] **4.3** If anime has multiple seasons, render tab navigation (Season 1, Season 2, etc.)
- [ ] **4.4** Display episodes as grid: 8 columns x N rows (reuse .epgrid from styles.css)
- [ ] **4.5** Each episode shows: episode number, title (tooltip on hover)
- [ ] **4.6** Add checkboxes for episode selection (prepare for download feature)
- [ ] **4.7** "Select All" / "Deselect All" buttons above episode grid

### Task 5: Implement Download Button & Action
- [ ] **5.1** Create "Download Selected Episodes" button (prominent, Sakura Night primary style)
- [ ] **5.2** Button disabled if no episodes selected
- [ ] **5.3** On click, show confirmation modal: "Download 5 episodes of Naruto?"
- [ ] **5.4** Modal has "Cancel" and "Confirm" buttons
- [ ] **5.5** On confirm, create download job via POST /api/v1/jobs (from Story 3.1, placeholder for now)
- [ ] **5.6** Show success toast: "Download added to queue"
- [ ] **5.7** Optionally navigate to /queue page after toast

### Task 6: Testing & Performance Validation
- [ ] **6.1** Backend tests: AnimeDetailHandler (6+ tests - valid ID, invalid ID, 404, format)
- [ ] **6.2** Frontend tests: AnimeDetailPage component (render, loading, error, episode selection)
- [ ] **6.3** Integration test: click anime in search → detail page loads
- [ ] **6.4** Performance test: detail page loads within 200ms
- [ ] **6.5** Visual regression: detail page matches Sakura Night design
- [ ] **6.6** All tests passing, zero regressions

---

## 📝 Dev Notes

### Architecture Overview

**Backend** (Go):
- New domain model: `AnimeDetail` (extends AnimeSearchResult)
- New service: `AnimeDetailService` interface (mock implementation for now)
- New handler: `AnimeDetailHandler` for GET /api/v1/anime/:id
- Route: `GET /api/v1/anime/:id` returns full anime details

**Frontend** (React/TypeScript):
- New page: `AnimeDetailPage.tsx`
- New component: `EpisodeList.tsx` (with season tabs)
- New component: `DownloadConfirmModal.tsx` (optional)
- Routing: `/anime/:id` route added to React Router

### Domain Model Design

```go
// internal/domain/anime_detail.go
type AnimeDetail struct {
    ID           string   `json:"id"`
    Title        string   `json:"title"`
    ThumbnailURL string   `json:"thumbnail_url"`
    Synopsis     string   `json:"synopsis"`
    Year         int      `json:"year"`
    Status       string   `json:"status"`      // "ongoing", "completed"
    Genres       []string `json:"genres"`      // ["Action", "Adventure", "Shonen"]
    EpisodeCount int      `json:"episode_count"`
    Seasons      []Season `json:"seasons"`
}

type Season struct {
    Number   int       `json:"number"`   // 1, 2, 3...
    Name     string    `json:"name"`     // "Season 1", "Part 2", etc.
    Episodes []Episode `json:"episodes"`
}

type Episode struct {
    Number       int    `json:"number"`        // Episode number within season
    Title        string `json:"title"`         // Episode title (optional)
    SeasonNumber int    `json:"season_number"` // Which season this belongs to
    URL          string `json:"url"`           // Download URL (for later use)
}
```

### API Response Format

**Request:**
```
GET /api/v1/anime/naruto-shippuden
```

**Response (200 OK):**
```json
{
  "id": "naruto-shippuden",
  "title": "Naruto Shippuden",
  "thumbnail_url": "https://cdn.anime-sama.si/naruto-shippuden.jpg",
  "synopsis": "Naruto Uzumaki is back after two years of training...",
  "year": 2007,
  "status": "completed",
  "genres": ["Action", "Adventure", "Shonen"],
  "episode_count": 500,
  "seasons": [
    {
      "number": 1,
      "name": "Season 1",
      "episodes": [
        {"number": 1, "title": "Homecoming", "season_number": 1, "url": "..."},
        {"number": 2, "title": "The Akatsuki Makes Its Move", "season_number": 1, "url": "..."}
      ]
    }
  ]
}
```

**Response (404 Not Found):**
```json
{
  "error": "Anime not found",
  "id": "invalid-id"
}
```

### Frontend Component Structure

```tsx
// AnimeDetailPage.tsx
function AnimeDetailPage() {
  const { id } = useParams();
  const [anime, setAnime] = useState<AnimeDetail | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedEpisodes, setSelectedEpisodes] = useState<number[]>([]);

  useEffect(() => {
    fetchAnimeDetail(id);
  }, [id]);

  return (
    <div className="container">
      {loading && <Loader />}
      {error && <ErrorMessage message={error} />}
      {anime && (
        <>
          <AnimeHeader anime={anime} />
          <EpisodeList 
            seasons={anime.seasons}
            selectedEpisodes={selectedEpisodes}
            onSelectionChange={setSelectedEpisodes}
          />
          <DownloadButton 
            disabled={selectedEpisodes.length === 0}
            onClick={() => handleDownload(selectedEpisodes)}
          />
        </>
      )}
    </div>
  );
}
```

### Episode Selection UI

Use existing `.epgrid` and `.ep` styles from `styles.css`:
```css
.epgrid { display: grid; grid-template-columns: repeat(8, minmax(0, 1fr)); gap: 8px; }
.ep { display: flex; align-items: center; gap: 6px; padding: 8px 10px; border: 1px solid var(--border); border-radius: 12px; background: rgba(0,0,0,.18); }
.ep input { accent-color: var(--accent); }
```

### Season Tabs Pattern

Standard tab UI pattern:
- Tabs above episode list (horizontal nav)
- Active tab highlighted with pink accent border
- Click tab to switch seasons
- Episode grid updates reactively

```tsx
<div className="season-tabs">
  {seasons.map(season => (
    <button 
      key={season.number}
      className={`tab ${activeSeason === season.number ? 'active' : ''}`}
      onClick={() => setActiveSeason(season.number)}
    >
      {season.name}
    </button>
  ))}
</div>
<EpisodeGrid episodes={seasons[activeSeason - 1].episodes} />
```

### Mock Service Implementation

For this story, create mock service that returns hardcoded data:
```go
// internal/app/mock_anime_detail_service.go
type MockAnimeDetailService struct {
    fixtures map[string]domain.AnimeDetail
}

func (s *MockAnimeDetailService) GetDetail(ctx context.Context, id string) (domain.AnimeDetail, error) {
    detail, ok := s.fixtures[id]
    if !ok {
        return domain.AnimeDetail{}, fmt.Errorf("anime not found: %s", id)
    }
    return detail, nil
}
```

Real scraping service will be implemented in later stories (Story 2.6 or Epic 3).

### Performance Considerations

- Detail page loads within 200ms (backend mock returns immediately)
- Episode list rendering: virtualize if >500 episodes (use react-window)
- Season tabs: lazy load episodes on tab switch (optional optimization)
- Thumbnail loading: use lazy loading + placeholder

### Dependencies & Libraries

**Backend**:
- chi router (existing) for URL params
- No new dependencies

**Frontend**:
- React Router (existing) for :id param
- Optional: react-window for virtualized episode list (only if performance issue)
- Optional: react-tooltip for episode title tooltips

### File Structure

```
NEW FILES:
  - internal/domain/anime_detail.go                        (60-80 lines)
  - internal/ports/anime_detail.go                         (10-15 lines)
  - internal/app/mock_anime_detail_service.go              (100-150 lines)
  - internal/adapters/httpapi/anime_detail.go              (100-120 lines)
  - internal/adapters/httpapi/anime_detail_test.go         (150-200 lines)
  - webapp/src/pages/AnimeDetailPage.tsx                   (250-350 lines)
  - webapp/src/components/anime/EpisodeList.tsx            (150-200 lines)
  - webapp/src/components/anime/DownloadConfirmModal.tsx   (80-100 lines)

MODIFIED FILES:
  - internal/adapters/httpapi/router.go                    (register anime detail route)
  - cmd/asd-server/main.go                                 (initialize mock detail service)
  - webapp/src/App.tsx                                     (add /anime/:id route)
```

### Testing Strategy

**Backend Tests** (Go):
- `TestAnimeDetailHandler_ValidID` - returns 200 with full detail
- `TestAnimeDetailHandler_InvalidID` - returns 404
- `TestAnimeDetailHandler_ResponseFormat` - validates JSON structure
- `TestAnimeDetailHandler_MultipleSeasons` - validates seasons array
- `TestAnimeDetailHandler_SingleSeason` - validates single season
- `TestAnimeDetailHandler_ContextCancellation` - respects ctx.Done()

**Frontend Tests** (React Testing Library):
- `renders loading state initially` - shows spinner
- `renders anime detail after fetch` - displays title, synopsis, etc.
- `renders error state on 404` - shows error message
- `episode selection works` - checkbox toggle functional
- `download button disabled when no episodes selected` - button disabled state
- `season tabs switch correctly` - clicking tab updates episode list

### Known Issues & Constraints

- Mock service returns hardcoded data (real scraping in future stories)
- Download button creates placeholder job (real job creation in Story 3.1)
- Episode URLs are mock data (real URLs scraped in Epic 3)
- No episode filtering (by watched/unwatched, to be added in Epic 5)

### Code Patterns from Stories 2-1 & 2-2

**From Story 2-1** (Anime Search):
- Reuse domain model patterns (ID, Title, ThumbnailURL, Year, Status)
- Reuse error handling (404 returns JSON error object)
- Reuse handler pattern (NewHandler, Routes method)

**From Story 2-2** (Autocomplete):
- Reuse AbortController pattern for cancelling fetches
- Reuse loading/error state management
- Reuse Sakura Night styling patterns

**Recommended Approach**:
- Start with backend: create domain models + mock service + handler + tests
- Then frontend: create page + episode list + routing
- Finally: styling + performance validation

---

## 🗂️ Project Context

**Git Reference**: `go-rewrite` branch, Stories 1.1-1.3, 2.1-2.2 completed  
**Dependencies**: Stories 2-1 (Search API) and 2-2 (Autocomplete) MUST be complete  
**Latest Patterns**: 
- Story 2-1: AnimeSamaSearchService with 1.29ms performance
- Story 2-2: AutocompleteSuggestions with debounce + keyboard nav
- Clean Architecture: domain → ports → adapters
- Sakura Night design system for all UI components

**Related Artifacts**:
- [Story 2-1](./2-1-implement-anime-search-api-endpoint.md) — Search API (provides anime IDs)
- [Story 2-2](./2-2-implement-search-autocomplete.md) — Autocomplete (navigation source)
- [02-DESIGN-SYSTEM-SAKURA-NIGHT.md](../planning-artifacts/02-DESIGN-SYSTEM-SAKURA-NIGHT.md) — UI styling
- [03-TECHNICAL-ARCHITECTURE.md](../planning-artifacts/03-TECHNICAL-ARCHITECTURE.md) — Backend patterns
- [epics.md](../planning-artifacts/epics.md#epic-2-anime-search--discovery) — Epic 2 roadmap

---

## 📦 File List

*To be updated after implementation*

---

## 📋 Change Log

*Entries will be added as implementation progresses*

---

## 🧪 Test Checklist

- [ ] Backend anime detail handler tests passing (6+ tests)
- [ ] Frontend AnimeDetailPage component tests passing (5+ tests)
- [ ] Episode selection functionality validated
- [ ] Season tabs switching validated (if multiple seasons)
- [ ] Performance test: page loads < 200ms
- [ ] Visual regression: matches Sakura Night design
- [ ] Full test suite: `go test ./...` + `npm test` all passing
- [ ] Zero regressions (203 existing tests still pass)

---

## Dev Agent Record

### Agent Model
GitHub Copilot (Claude Sonnet 4.5)

### Implementation Status
🚀 READY-FOR-DEV — Comprehensive context complete, Stories 2-1 & 2-2 dependencies satisfied

### Debug Log
None yet (story not started)

### Completion Notes
*(To be filled during implementation)*

---

## Status

**Current Status:** ready-for-dev  
**Progress:** 0/6 major tasks completed (0%)  
**Created:** 31 janvier 2025  
**Assigned to:** Dev Agent (Amelia)

**Next Action**: Run implementation workflow to begin Task 1 (domain models & service creation).
