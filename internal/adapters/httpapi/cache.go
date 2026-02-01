package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/app"
)

// CacheHandler handles catalogue cache operations
type CacheHandler struct {
	cache *app.CatalogueCache
}

// NewCacheHandler creates a new cache handler
func NewCacheHandler(cache *app.CatalogueCache) *CacheHandler {
	return &CacheHandler{cache: cache}
}

// RefreshCache manually triggers a cache refresh
func (h *CacheHandler) RefreshCache(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	_, err := h.cache.Refresh(ctx)
	if err != nil {
		http.Error(w, "Failed to refresh cache: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Cache refreshed successfully",
		"stats":   h.cache.Stats(),
	})
}

// CacheStats returns cache statistics
func (h *CacheHandler) CacheStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(h.cache.Stats())
}
