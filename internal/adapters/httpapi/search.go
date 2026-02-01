package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/ports"
)

// SearchHandler handles anime search requests
type SearchHandler struct {
	searchService ports.AnimeSearch
}

// NewSearchHandler creates a new search handler
func NewSearchHandler(searchService ports.AnimeSearch) *SearchHandler {
	return &SearchHandler{
		searchService: searchService,
	}
}

// SearchRequest represents the search query parameters
type SearchRequest struct {
	Query string `json:"q"`
}

// SearchResponse represents a single search result in the HTTP response
type SearchResponse struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	ThumbnailURL string `json:"thumbnail_url"`
	Year         int    `json:"year"`
	Status       string `json:"status"`
	EpisodeCount int    `json:"episode_count"`
}

// Search handles GET /api/v1/search?q={query}
func (h *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	// Parse query parameter
	query := r.URL.Query().Get("q")

	// Call search service with context
	results, err := h.searchService.Search(r.Context(), query)
	if err != nil {
		// Context cancelled or other error
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// Map domain results to HTTP response
	response := make([]SearchResponse, len(results))
	for i, result := range results {
		response[i] = SearchResponse{
			ID:           result.ID,
			Title:        result.Title,
			ThumbnailURL: result.ThumbnailURL,
			Year:         result.Year,
			Status:       result.Status,
			EpisodeCount: result.EpisodeCount,
		}
	}

	// Return results
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// RegisterSearchRoutes registers search routes in the chi router
func RegisterSearchRoutes(r chi.Router, searchService ports.AnimeSearch) {
	handler := NewSearchHandler(searchService)
	r.Get("/search", handler.Search)
	
	// Register autocomplete route
	autocompleteHandler := NewAutocompleteHandler(searchService)
	r.Get("/search/autocomplete", autocompleteHandler.handleAutocomplete)
}
