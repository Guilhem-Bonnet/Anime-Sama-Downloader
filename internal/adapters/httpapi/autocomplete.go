package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/ports"
)

// AutocompleteHandler handles autocomplete search requests.
type AutocompleteHandler struct {
	searchService ports.AnimeSearch
}

// NewAutocompleteHandler creates a new autocomplete handler.
func NewAutocompleteHandler(searchService ports.AnimeSearch) *AutocompleteHandler {
	return &AutocompleteHandler{
		searchService: searchService,
	}
}

// AutocompleteSuggestion is a lightweight DTO for autocomplete results.
type AutocompleteSuggestion struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	ThumbnailURL string `json:"thumbnail_url"`
	Year         int    `json:"year"`
}

// handleAutocomplete handles GET /api/v1/search/autocomplete?q={query}
func (h *AutocompleteHandler) handleAutocomplete(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	// Early return for queries < 2 characters
	if len(query) < 2 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]AutocompleteSuggestion{})
		return
	}

	// Use search service from Story 2-1
	results, err := h.searchService.Search(r.Context(), query)
	if err != nil {
		// Log error but return empty array (graceful degradation)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]AutocompleteSuggestion{})
		return
	}

	// Limit to 10 results for autocomplete
	limit := 10
	if len(results) > limit {
		results = results[:limit]
	}

	// Convert to lightweight DTO (omit status and episode_count)
	suggestions := make([]AutocompleteSuggestion, len(results))
	for i, result := range results {
		suggestions[i] = AutocompleteSuggestion{
			ID:           result.ID,
			Title:        result.Title,
			ThumbnailURL: result.ThumbnailURL,
			Year:         result.Year,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(suggestions)
}

// Routes registers autocomplete routes on the provided router.
func (h *AutocompleteHandler) Routes(r http.Handler) {
	// Note: This will be called from router.go with chi.Router
	// For now, we'll register directly in router.go
}
