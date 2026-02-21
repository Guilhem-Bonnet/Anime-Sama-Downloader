package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/ports"
)

// RecommendationsHandler handles anime recommendations
type RecommendationsHandler struct {
	recommendationsService ports.RecommendationsService
}

// NewRecommendationsHandler creates a new recommendations handler
func NewRecommendationsHandler(recommendationsService ports.RecommendationsService) *RecommendationsHandler {
	return &RecommendationsHandler{recommendationsService: recommendationsService}
}

// Routes registers recommendations endpoints.
func (h *RecommendationsHandler) Routes(r chi.Router) {
	r.Route("/recommendations", func(r chi.Router) {
		r.Get("/similar", h.GetSimilarAnime)
		r.Get("/query", h.GetRecommendationsByQuery)
		r.Get("/genres", h.GetRecommendationsForGenres)
	})
}

// GetSimilarAnime returns anime similar to the given anime
// GET /api/v1/recommendations/similar?anime_id=1&limit=10
func (h *RecommendationsHandler) GetSimilarAnime(w http.ResponseWriter, r *http.Request) {
	animeID := r.URL.Query().Get("anime_id")
	limitStr := r.URL.Query().Get("limit")

	if animeID == "" {
		http.Error(w, "anime_id parameter is required", http.StatusBadRequest)
		return
	}

	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	recommendations, err := h.recommendationsService.GetSimilarAnime(animeID, limit)
	if err != nil {
		http.Error(w, "Failed to get recommendations: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"recommendations": recommendations,
		"count":           len(recommendations),
	})
}

// GetRecommendationsByQuery returns recommendations based on search query
// GET /api/v1/recommendations/query?q=naruto&limit=10
func (h *RecommendationsHandler) GetRecommendationsByQuery(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	limitStr := r.URL.Query().Get("limit")

	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	recommendations, err := h.recommendationsService.GetRecommendationsByQuery(query, limit)
	if err != nil {
		http.Error(w, "Failed to get recommendations: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"recommendations": recommendations,
		"count":           len(recommendations),
	})
}

// GetRecommendationsForGenres returns recommendations for multiple genres
// GET /api/v1/recommendations/genres?genres=Action,Adventure&limit=10
func (h *RecommendationsHandler) GetRecommendationsForGenres(w http.ResponseWriter, r *http.Request) {
	genresStr := r.URL.Query().Get("genres")
	limitStr := r.URL.Query().Get("limit")

	if genresStr == "" {
		http.Error(w, "genres parameter is required", http.StatusBadRequest)
		return
	}

	// Parse comma-separated genres
	genres := strings.Split(genresStr, ",")
	for i, g := range genres {
		genres[i] = strings.TrimSpace(g)
	}

	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	recommendations, err := h.recommendationsService.GetRecommendationsForGenres(genres, limit)
	if err != nil {
		http.Error(w, "Failed to get recommendations: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"recommendations": recommendations,
		"count":           len(recommendations),
	})
}
