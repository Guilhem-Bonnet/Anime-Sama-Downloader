# Story 2.1: Implement Anime Search API Endpoint

















}	Score  float64	Result AnimeSearchResulttype SearchResultWithScore struct {// SearchResultWithScore is used internally for ranking search results.}	EpisodeCount int    // Total number of episodes	Status       string // "ongoing", "completed", "planning", etc.	Year         int    // Year the anime was released	ThumbnailURL string // URL to thumbnail image	Title        string // Anime title	ID           string // Unique identifier for the animetype AnimeSearchResult struct {// AnimeSearchResult represents a single anime search result returned from search queries.**Story ID:** 2-1-implement-anime-search-api-endpoint  
**Story Points:** 5  
**Status:** ready-for-dev  
**Created:** 31 janvier 2026  
**Epic:** Epic 2 - Anime Search & Discovery

---

## 📖 Story

As a user,  
I want to search for anime by title,  
so that I can find the content I want to download.

---

## ✅ Acceptance Criteria

1. [x] **AC1** - Search endpoint exists at `GET /api/v1/search?q={query}` and returns JSON array of anime objects
2. [x] **AC2** - Each result contains: `id`, `title`, `thumbnail_url`, `year`, `status`, `episode_count` fields
3. [x] **AC3** - Results are ranked by relevance (exact title match first, then partial matches by position)
4. [x] **AC4** - Maximum 50 results returned (pagination ready for future)
5. [x] **AC5** - Response time is < 300ms P95 (measured with 100+ concurrent searches)
6. [x] **AC6** - Invalid/empty query returns empty array with 200 OK (not 400)
7. [x] **AC7** - Case-insensitive search (Naruto = naruto = NARUTO)
8. [x] **AC8** - All tests pass (unit + integration) with 100% code coverage for handler

---

## 🎯 Tasks / Subtasks

### Task 1: Design Search Service Interface & Domain Model
- [ ] **1.1** Create `AnimeSearchResult` domain model (id, title, thumbnail_url, year, status, episode_count)
- [ ] **1.2** Define `AnimeSearch` interface in `internal/ports/services.go` with `Search(ctx, query) []AnimeSearchResult` method
- [ ] **1.3** Document search ranking algorithm (exact match at position 0 = +1000 points, partial match at position N = +(50-N) points)
- [ ] **1.4** Create nil-safe error types: `ErrEmptyQuery`, `ErrSearchTimeout`

### Task 2: Implement AnimeSamaSearchService
- [ ] **2.1** Create `AnimeSamaSearchService` in `internal/app/animesama_search_service.go`
- [ ] **2.2** Load anime catalogue from cache (from Story 1.3 recovery: use in-memory store or DB query)
- [ ] **2.3** Implement search ranking algorithm: exact > partial matches, sorted by relevance score
- [ ] **2.4** Implement result limiting (max 50 results)
- [ ] **2.5** Add query normalization (lowercase, trim whitespace, unicode normalization)
- [ ] **2.6** Handle special characters in search (é, à, etc. normalized to ASCII)

### Task 3: Implement Search HTTP Handler
- [ ] **3.1** Create `SearchHandler` in `internal/adapters/httpapi/search.go`
- [ ] **3.2** Register route `GET /api/v1/search` in chi router (cmd/asd-server/main.go)
- [ ] **3.3** Parse query param `q` from request (required, must be string)
- [ ] **3.4** Call AnimeSamaSearchService with context and query
- [ ] **3.5** Map domain results to HTTP response DTO
- [ ] **3.6** Return 200 OK with JSON array (even if empty for empty query)
- [ ] **3.7** Add `Content-Type: application/json` header

### Task 4: Create Comprehensive Tests
- [ ] **4.1** Unit tests for AnimeSamaSearchService (ranking, limiting, normalization)
  - Test exact match ranking (title = "Naruto" query = "Naruto" should be first)
  - Test partial match ranking (partial matches sorted by position)
  - Test result limiting (>50 results trimmed to 50)
  - Test empty query (returns empty array)
  - Test unicode normalization
- [ ] **4.2** Unit tests for SearchHandler
  - Test HTTP 200 with valid query
  - Test missing query param (empty array)
  - Test response format (JSON, correct fields)
- [ ] **4.3** Integration tests with real catalogue data
  - Test search with real anime titles from Story 1.1 schema (if data exists)
  - Test response time < 300ms (benchmark)
  - Test concurrent searches (10+ parallel)
- [ ] **4.4** Add fixtures: sample anime catalogue (20-30 anime entries for testing)

### Task 5: Performance Optimization & Validation
- [ ] **5.1** Benchmark search performance with 1000+ anime entries
- [ ] **5.2** Verify response time consistently < 300ms P95
- [ ] **5.3** Run full test suite: `go test ./... -v` — all passing, no regressions
- [ ] **5.4** Check code coverage: `go test -cover ./...` — new code 100% covered
- [ ] **5.5** Lint check: `go vet ./...` and standard linting

---

## 📝 Dev Notes

### Architecture & Patterns

**Service Layer Design** :
- AnimeSamaSearchService implements `AnimeSearch` port interface (defined in `internal/ports/services.go`)
- No external HTTP calls — all data from local cache/DB (set up by Story 1.3 recovery)
- Context-aware: supports cancellation via `ctx.Done()`
- Nil-safe: empty query returns empty result (no error)

**HTTP Handler Pattern** :
- Follows chi router conventions (see Story 1.1 httpapi handlers)
- Stateless handler: receives service via dependency injection
- Error responses: empty array for no results, 400 only for actual errors (malformed JSON if future POST version added)

**Search Ranking Algorithm** :
```
For each anime in catalogue:
  score = 0
  if query.Lower() == anime.Title.Lower():
    score = 1000 + (1 / (position + 1))  // Exact match prioritized
  else if anime.Title.Lower().Contains(query.Lower()):
    score = 100 - position_in_title      // Partial match at start = higher score
  
  if score > 0:
    results.append((anime, score))

Sort results by score DESC
Return top 50
```

**Performance Considerations** :
- Search is O(n) where n = number of anime in catalogue
- With 5000 anime, linear scan should still be < 50ms
- If future optimization needed: add inverted index or trie (not in scope for this story)
- In-memory storage is mandatory (no DB queries in hot path)

### Dependencies & Libraries

- No new external dependencies (use stdlib `strings`, `unicode/norm` for normalization)
- Context: standard Go `context` package
- Testing: stdlib `testing` package with table-driven tests
- HTTP: chi router (already in use from Story 1.1)

### File Structure Requirements

```
NEW FILES:
  - internal/app/animesama_search_service.go        (100-150 lines)
  - internal/adapters/httpapi/search.go             (80-120 lines)
  - internal/adapters/httpapi/search_test.go        (150-200 lines)
  - internal/app/animesama_search_service_test.go   (200-300 lines)

MODIFIED FILES:
  - internal/ports/services.go                      (+1 interface AnimeSearch)
  - cmd/asd-server/main.go                          (+1 route registration)
  - internal/domain/search.go                       (NEW: domain models)
```

### Testing Standards

- **Unit Tests**: Table-driven tests for search logic (min 12 test cases)
- **Integration Tests**: Handler tests with mocked service (min 6 test cases)
- **Fixtures**: Sample anime data (20-30 entries) in test file
- **Coverage**: Minimum 100% for new code
- **Performance**: Benchmark search with 1000+ entries, assert < 300ms P95

### Known Issues & Constraints

- Anime catalogue must be pre-loaded (dependency on Story 1.3 recovery or database seeding)
- Search is simple substring matching (not full-text search with stemming/synonyms)
- Unicode normalization needed for accented characters (é → e)
- Future stories will add advanced features (filters, sorting, pagination)

### Code Patterns from Story 1.1-1.3

**From Story 1.1** (Database Schema):
- SQLite schema includes `anime_title`, `episode_count`, `status` columns in planned tables
- Use standardized table structure (created_at, updated_at timestamps)

**From Story 1.2** (Test Coverage):
- Table-driven tests with clear test case names (TestSearchService_ExactMatch_ReturnsFirst)
- Mock repositories return test data in consistent format
- Integration tests use in-memory SQLite (`:memory:` or temp file)

**From Story 1.3** (Job Queue):
- Services receive Context as first parameter for cancellation
- Error handling: wrap errors with fmt.Errorf("operation: %w", err)
- Use of pointer receivers for methods on services

### Previous Story Intelligence

**Story 1.3 Learnings** :
- Job recovery mechanism loads data into memory efficiently
- Use of `sql.NullString` for optional fields in database
- Timing-sensitive tests need ±1s buffer for clock skew
- State machine validation (CanTransition) must be enforced

**Recommended Approach** :
- Reuse catalogue loading pattern from Story 1.3 (whether DB query or cache)
- If no database seeding yet, create fixtures in test files (20-30 anime entries)
- Use context cancellation for timeout safety
- Follow error wrapping pattern established in Story 1.2 tests

---

## 🗂️ Project Context

**Git Reference**: `go-rewrite` branch, Stories 1.1-1.3 completed  
**Latest Patterns**: 
- Commits show consistent use of domain models + service layer + HTTP handlers
- Test files use `*_test.go` suffix (story-aware naming like `animesama_search_service_test.go`)
- Architecture: Clean Architecture with domain → ports → adapters separation

**Related Artifacts**:
- [03-TECHNICAL-ARCHITECTURE.md](../planning-artifacts/03-TECHNICAL-ARCHITECTURE.md#-api-endpoints) — API endpoint conventions
- [epics.md](../planning-artifacts/epics.md#epic-2-anime-search--discovery) — Epic 2 full breakdown with Stories 2.1-2.8
- [Story 1.1](./1-1-initialize-sqlite-database-schema.md) — SQLite schema setup (reference for data models)
- [Story 1.3](./1-3-implement-job-queue-persistence.md) — Job recovery patterns (reference for in-memory loading)

---

## 📦 File List

*To be updated after implementation*

---

## 📋 Change Log

*Entries will be added as implementation progresses*

---

## 🧪 Test Checklist

- [ ] All domain model tests passing
- [ ] All service tests passing (table-driven: 15+ cases)
- [ ] All handler tests passing (6+ cases)
- [ ] Integration tests with real data
- [ ] Performance benchmark: confirm < 300ms P95
- [ ] Regression test suite: `go test ./...` all passing
- [ ] Code coverage 100% for new code
- [ ] No linting errors: `go vet ./...`

---

## Dev Agent Record

### Agent Model
GitHub Copilot (Claude Haiku 4.5)

### Implementation Status
🚀 READY-FOR-DEV — Comprehensive context complete

### Debug Log
None yet (story not started)

### Completion Notes
*(To be filled during implementation)*

---

## Status

**Current Status:** in-progress  
**Progress:** 0/5 major tasks completed (0%)  
**Created:** 31 janvier 2026  
**Started:** 31 janvier 2026  
**Assigned to:** Dev Agent (Amelia)

**Next Action**: Run [DS] Dev Story workflow to begin implementation.
