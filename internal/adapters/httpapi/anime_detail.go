package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/ports"
)

// AnimeDetailHandler handles anime detail requests.
type AnimeDetailHandler struct {
	detailService ports.AnimeDetailService
}

// NewAnimeDetailHandler creates a new anime detail handler.
func NewAnimeDetailHandler(detailService ports.AnimeDetailService) *AnimeDetailHandler {
	return &AnimeDetailHandler{
		detailService: detailService,
	}
}

// handleAnimeDetail handles GET /api/v1/anime/:id
func (h *AnimeDetailHandler) handleAnimeDetail(w http.ResponseWriter, r *http.Request) {
	// Extract anime ID from URL parameter
	id := chi.URLParam(r, "id")
	if id == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Missing anime ID"})
		return
	}

	// Fetch anime detail from service
	detail, err := h.detailService.GetDetail(r.Context(), id)
	if err != nil {
		// Assume "not found" if error contains "not found"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Anime not found",
			"id":    id,
		})
		return
	}

	// Return anime detail
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(detail)
}

// Routes registers anime detail routes.
func (h *AnimeDetailHandler) Routes(r chi.Router) {
	r.Get("/anime/{id}", h.handleAnimeDetail)
}
