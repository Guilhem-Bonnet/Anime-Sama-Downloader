# Story 2.2: Implement Search Autocomplete

**Story ID:** 2-2-implement-search-autocomplete  
**Story Points:** 3  
**Status:** ready-for-dev  
**Created:** 31 janvier 2025  
**Epic:** Epic 2 - Anime Search & Discovery

---

## 📖 Story

As a user,  
I want autocomplete suggestions as I type in the search field,  
so that I can quickly find anime titles without typing the full name.

---

## ✅ Acceptance Criteria

1. [ ] **AC1** - Autocomplete suggestions appear after typing at least 2 characters in the search field
2. [ ] **AC2** - Suggestions are debounced (300ms wait after last keystroke before API call)
3. [ ] **AC3** - Maximum 10 suggestions are displayed below the search field
4. [ ] **AC4** - Each suggestion shows: anime title + thumbnail + year
5. [ ] **AC5** - Clicking a suggestion navigates to anime detail page (or triggers full search)
6. [ ] **AC6** - Suggestions update within 150ms of typing pause (P95)
7. [ ] **AC7** - Keyboard navigation works (↑↓ to select, Enter to choose, Esc to close)
8. [ ] **AC8** - All tests pass (backend + frontend) with coverage for autocomplete logic

---

## 🎯 Tasks / Subtasks

### Task 1: Create Autocomplete Backend Endpoint
- [ ] **1.1** Create `GET /api/v1/search/autocomplete?q={query}` endpoint in `internal/adapters/httpapi/search.go`
- [ ] **1.2** Reuse `AnimeSamaSearchService` from Story 2-1 with limit=10 (no need to create new service)
- [ ] **1.3** Return lightweight DTO: `{id, title, thumbnail_url, year}` (no status/episode_count to reduce payload)
- [ ] **1.4** Add early return if query length < 2 characters (return empty array)
- [ ] **1.5** Ensure response time < 150ms (same in-memory search as Story 2-1)

### Task 2: Implement Debounce Logic in Frontend
- [ ] **2.1** Create custom React hook `useDebounce(value, delay)` or use existing library (lodash.debounce)
- [ ] **2.2** Wrap search input onChange handler with 300ms debounce
- [ ] **2.3** Cancel pending API calls if new input arrives (AbortController)
- [ ] **2.4** Show loading spinner during API call (debounce complete → API pending)
- [ ] **2.5** Clear suggestions when input is cleared or < 2 characters

### Task 3: Create AutocompleteSuggestions Component
- [ ] **3.1** Create `webapp/src/components/search/AutocompleteSuggestions.tsx` component
- [ ] **3.2** Display suggestions as dropdown list positioned below search input
- [ ] **3.3** Each suggestion item shows: thumbnail (48x48), title (bold), year (gray text)
- [ ] **3.4** Implement click handler: navigate to `/anime/{id}` or trigger full search
- [ ] **3.5** Add "No results found" message when suggestions array is empty (but query >= 2 chars)
- [ ] **3.6** Add subtle animation (fade-in) when suggestions appear
- [ ] **3.7** Close suggestions when clicking outside (useClickOutside hook)

### Task 4: Implement Keyboard Navigation
- [ ] **4.1** Track selected suggestion index in component state
- [ ] **4.2** Handle ArrowDown: increment index (wrap to 0 at end)
- [ ] **4.3** Handle ArrowUp: decrement index (wrap to last at start)
- [ ] **4.4** Handle Enter: trigger navigation/search with selected suggestion
- [ ] **4.5** Handle Escape: close suggestions dropdown
- [ ] **4.6** Highlight selected suggestion with background color change
- [ ] **4.7** Scroll suggestion into view if outside visible area

### Task 5: Testing & Performance Validation
- [ ] **5.1** Backend tests: autocomplete handler (6+ tests - valid query, short query, empty, performance)
- [ ] **5.2** Frontend tests: AutocompleteSuggestions component (render, click, keyboard nav)
- [ ] **5.3** Integration test: debounce behavior (ensure API not called on every keystroke)
- [ ] **5.4** Benchmark autocomplete endpoint: confirm < 150ms P95
- [ ] **5.5** Visual regression test: suggestions dropdown styling matches design system
- [ ] **5.6** All tests passing, zero regressions

---

## 📝 Dev Notes

### Architecture Overview

**Backend** (Go):
- Reuse `AnimeSamaSearchService.Search()` from Story 2-1
- New handler: `AutocompleteHandler` (or extend `SearchHandler`)
- New route: `GET /api/v1/search/autocomplete`
- DTO: Lightweight response (omit status/episode_count to reduce payload)

**Frontend** (React/TypeScript):
- New component: `AutocompleteSuggestions.tsx`
- Custom hook: `useDebounce.ts` (or lodash.debounce)
- API call: `GET /api/v1/search/autocomplete?q={query}`
- State management: Local component state (no Zustand needed for this)

### Debounce Implementation

**Why 300ms?**
- Balance between responsiveness and reducing API load
- User typing speed: 300ms is comfortable pause
- Prevents API spam (10 keystrokes = 1 API call instead of 10)

**Example implementation:**
```typescript
const useDebounce = (value: string, delay: number) => {
  const [debouncedValue, setDebouncedValue] = useState(value);
  
  useEffect(() => {
    const timer = setTimeout(() => setDebouncedValue(value), delay);
    return () => clearTimeout(timer);
  }, [value, delay]);
  
  return debouncedValue;
};

// Usage in SearchInput component:
const [query, setQuery] = useState('');
const debouncedQuery = useDebounce(query, 300);

useEffect(() => {
  if (debouncedQuery.length >= 2) {
    fetchAutocompleteSuggestions(debouncedQuery);
  }
}, [debouncedQuery]);
```

### API Response Format

**Request:**
```
GET /api/v1/search/autocomplete?q=naruto
```

**Response (200 OK):**
```json
[
  {
    "id": "anime-123",
    "title": "Naruto",
    "thumbnail_url": "https://cdn.anime-sama.si/naruto.jpg",
    "year": 2002
  },
  {
    "id": "anime-456",
    "title": "Naruto Shippuden",
    "thumbnail_url": "https://cdn.anime-sama.si/naruto-shippuden.jpg",
    "year": 2007
  }
]
```

**Empty query (< 2 chars):**
```json
[]
```

### Keyboard Navigation Pattern

Standard autocomplete UX:
- **ArrowDown**: Move to next suggestion (highlight)
- **ArrowUp**: Move to previous suggestion
- **Enter**: Select highlighted suggestion (navigate/search)
- **Escape**: Close suggestions dropdown
- **Tab**: Close suggestions (optional: select first?)

**State management:**
```typescript
const [selectedIndex, setSelectedIndex] = useState(-1); // -1 = no selection

const handleKeyDown = (e: KeyboardEvent) => {
  switch (e.key) {
    case 'ArrowDown':
      setSelectedIndex((prev) => (prev + 1) % suggestions.length);
      break;
    case 'ArrowUp':
      setSelectedIndex((prev) => (prev - 1 + suggestions.length) % suggestions.length);
      break;
    case 'Enter':
      if (selectedIndex >= 0) navigate(`/anime/${suggestions[selectedIndex].id}`);
      break;
    case 'Escape':
      setSuggestions([]);
      break;
  }
};
```

### Styling Guidelines

Follow **Sakura Night Design System** from [02-DESIGN-SYSTEM-SAKURA-NIGHT.md](../planning-artifacts/02-DESIGN-SYSTEM-SAKURA-NIGHT.md):

- Background: `gray-800` (`#1F2937`)
- Selected item: `pink-500/10` background with `pink-500` left border
- Text: `gray-100` (title), `gray-400` (year)
- Thumbnail: Rounded corners (`rounded-md`)
- Dropdown: `shadow-lg` with smooth `transition-all duration-200`

### Performance Considerations

- Autocomplete uses same in-memory search as Story 2-1 (1.29ms per search)
- Limit to 10 results (vs 50 for full search) → even faster
- Debounce prevents API spam (critical for UX)
- AbortController cancels pending requests (avoid race conditions)

### Dependencies & Libraries

**Backend**:
- No new dependencies (reuse Story 2-1 service)
- chi router (existing)

**Frontend**:
- React 18+ (existing)
- TypeScript (existing)
- Optional: `lodash.debounce` (or custom hook)
- Optional: `react-use` (for useClickOutside)

### File Structure

```
NEW FILES:
  - internal/adapters/httpapi/autocomplete.go              (80-100 lines)
  - internal/adapters/httpapi/autocomplete_test.go         (120-150 lines)
  - webapp/src/components/search/AutocompleteSuggestions.tsx  (150-200 lines)
  - webapp/src/hooks/useDebounce.ts                        (15-20 lines)
  - webapp/src/components/search/AutocompleteSuggestions.test.tsx (100-120 lines)

MODIFIED FILES:
  - internal/adapters/httpapi/router.go                    (add autocomplete route)
  - webapp/src/components/ui/Input.tsx (or SearchInput.tsx)  (integrate autocomplete)
```

### Testing Strategy

**Backend Tests** (Go):
- `TestAutocompleteHandler_ValidQuery` - returns 10 results max
- `TestAutocompleteHandler_ShortQuery` - query < 2 chars returns []
- `TestAutocompleteHandler_EmptyQuery` - empty string returns []
- `TestAutocompleteHandler_ResponseFormat` - DTO has id, title, thumbnail_url, year only
- `TestAutocompleteHandler_Performance` - response < 150ms
- `TestAutocompleteHandler_CaseInsensitive` - Naruto = naruto

**Frontend Tests** (React Testing Library):
- `renders suggestions list` - component displays suggestions
- `clicking suggestion navigates` - onClick handler called
- `keyboard navigation works` - ArrowDown/Up/Enter/Esc
- `debounce delays API call` - not called immediately
- `shows no results message` - when suggestions.length === 0

### Known Issues & Constraints

- Autocomplete uses same catalogue as full search (no separate index)
- No fuzzy matching (exact/partial substring only, same as Story 2-1)
- No caching of autocomplete results (could add React Query later)
- Keyboard navigation requires focus on input (standard behavior)

### Code Patterns from Story 2-1

**From Story 2-1** (Anime Search API):
- Reuse `AnimeSamaSearchService.Search()` method
- Same ranking algorithm (exact > partial matches)
- Same unicode normalization (é → e)
- Same error handling (empty query returns [])

**Recommended Approach**:
- Create thin wrapper handler `AutocompleteHandler` that calls `searchService.Search(ctx, query)` and slices to 10 results
- Frontend: Use AbortController to cancel pending requests on new input
- Testing: Focus on debounce behavior and keyboard navigation (critical UX)

---

## 🗂️ Project Context

**Git Reference**: `go-rewrite` branch, Stories 1.1-1.3, 2.1 completed  
**Dependencies**: Story 2-1 (Implement Anime Search API) — MUST be complete (reuses AnimeSamaSearchService)  
**Latest Patterns**: 
- Story 2-1 established AnimeSamaSearchService with 1.29ms search performance
- Clean Architecture: domain → ports → adapters separation
- Frontend uses React 18 + TypeScript + Vite

**Related Artifacts**:
- [Story 2-1](./2-1-implement-anime-search-api-endpoint.md) — Anime Search API (dependency)
- [02-DESIGN-SYSTEM-SAKURA-NIGHT.md](../planning-artifacts/02-DESIGN-SYSTEM-SAKURA-NIGHT.md) — UI styling guidelines
- [03-TECHNICAL-ARCHITECTURE.md](../planning-artifacts/03-TECHNICAL-ARCHITECTURE.md#-api-endpoints) — API conventions
- [epics.md](../planning-artifacts/epics.md#epic-2-anime-search--discovery) — Epic 2 stories

---

## 📦 File List

*To be updated after implementation*

---

## 📋 Change Log

*Entries will be added as implementation progresses*

---

## 🧪 Test Checklist

- [ ] Backend autocomplete handler tests passing (6+ tests)
- [ ] Frontend AutocompleteSuggestions component tests passing (5+ tests)
- [ ] Debounce behavior validated (API not called on every keystroke)
- [ ] Keyboard navigation tested (ArrowDown/Up/Enter/Esc)
- [ ] Performance benchmark: autocomplete < 150ms P95
- [ ] Visual regression: suggestions dropdown matches design
- [ ] Full test suite: `go test ./...` + `npm test` all passing
- [ ] Zero regressions (197 existing Go tests still pass)

---

## Dev Agent Record

### Agent Model
GitHub Copilot (Claude Sonnet 4.5)

### Implementation Status
🚀 READY-FOR-DEV — Comprehensive context complete, Story 2-1 dependency satisfied

### Debug Log
None yet (story not started)

### Completion Notes
*(To be filled during implementation)*

---

## Status

**Current Status:** ready-for-dev  
**Progress:** 0/5 major tasks completed (0%)  
**Created:** 31 janvier 2025  
**Assigned to:** Dev Agent (Amelia)

**Next Action**: Run implementation workflow to begin Task 1 (backend endpoint creation).
