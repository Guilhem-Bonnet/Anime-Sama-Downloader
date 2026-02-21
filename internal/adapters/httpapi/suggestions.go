package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/ports"
)

// SuggestionsHandler handles search suggestions
type SuggestionsHandler struct {
	suggestionsService ports.SuggestionsService
}

// NewSuggestionsHandler creates a new suggestions handler
func NewSuggestionsHandler(suggestionsService ports.SuggestionsService) *SuggestionsHandler {
	return &SuggestionsHandler{suggestionsService: suggestionsService}
}

// GetSuggestions returns search suggestions based on query prefix
// GET /api/v1/suggestions?q=naruto&limit=10
func (h *SuggestionsHandler) GetSuggestions(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	limitStr := r.URL.Query().Get("limit")

	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	suggestions, err := h.suggestionsService.GetSuggestions(r.Context(), query, limit)
	if err != nil {
		http.Error(w, "Failed to get suggestions: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"suggestions": suggestions,
		"count":       len(suggestions),
	})
}

// GetTrendingSuggestions returns currently trending searches
// GET /api/v1/suggestions/trending?limit=10
func (h *SuggestionsHandler) GetTrendingSuggestions(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")

	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	suggestions, err := h.suggestionsService.GetTrendingSuggestions(r.Context(), limit)
	if err != nil {
		http.Error(w, "Failed to get trending suggestions: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"suggestions": suggestions,
		"count":       len(suggestions),
	})
}
