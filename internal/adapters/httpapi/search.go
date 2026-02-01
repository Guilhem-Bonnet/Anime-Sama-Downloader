package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

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
	Query  string   `json:"q"`
	Genres []string `json:"genres"`      // Filter by genres (e.g., ["Action", "Adventure"])
	Status string   `json:"status"`      // Filter by status: "ongoing", "completed", "planning", "" (all)
	YearMin int     `json:"year_min"`    // Minimum year (e.g., 2020), 0 = no filter
	YearMax int     `json:"year_max"`    // Maximum year (e.g., 2023), 0 = no filter
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

// Search handles GET /api/v1/search?q={query}&genres={genre1,genre2}&status={status}&year_min={year}&year_max={year}
func (h *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	query := r.URL.Query().Get("q")
	
	// Parse filters
	filters := ports.SearchFilters{
		Query: query,
	}
	
	// Parse genres (comma-separated)
	if genresParam := r.URL.Query().Get("genres"); genresParam != "" {
		filters.Genres = parseCommaSeparated(genresParam)
	}
	
	// Parse status
	filters.Status = r.URL.Query().Get("status")
	
	// Parse year range
	if yearMinStr := r.URL.Query().Get("year_min"); yearMinStr != "" {
		if yearMin, err := strconv.Atoi(yearMinStr); err == nil {
			filters.YearMin = yearMin
		}
	}
	if yearMaxStr := r.URL.Query().Get("year_max"); yearMaxStr != "" {
		if yearMax, err := strconv.Atoi(yearMaxStr); err == nil {
			filters.YearMax = yearMax
		}
	}

	// Call search service with filters
	results, err := h.searchService.SearchWithFilters(r.Context(), filters)
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

// parseCommaSeparated parses a comma-separated string into a slice of trimmed strings
func parseCommaSeparated(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// FileListHandler handles file listing requests
type FileListHandler struct {
	fileListService ports.FileListService
}

// NewFileListHandler creates a new file list handler
func NewFileListHandler(fileListService ports.FileListService) *FileListHandler {
	return &FileListHandler{
		fileListService: fileListService,
	}
}

// FileResponse represents a single file in the HTTP response
type FileResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	Size     int64  `json:"size"`
	Duration int    `json:"duration"`
	Type     string `json:"type"`
}

// FileListResponse represents the complete file list response
type FileListResponse struct {
	AnimeID string          `json:"anime_id"`
	Files   []FileResponse  `json:"files"`
	Count   int             `json:"count"`
}

// GetFiles handles GET /api/v1/anime/{animeId}/files
func (h *FileListHandler) GetFiles(w http.ResponseWriter, r *http.Request) {
	animeID := chi.URLParam(r, "animeId")

	if animeID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "anime ID required"})
		return
	}

	fileList, err := h.fileListService.GetFileList(r.Context(), animeID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// Map domain files to HTTP response
	files := make([]FileResponse, len(fileList.Files))
	for i, file := range fileList.Files {
		files[i] = FileResponse{
			ID:       file.ID,
			Name:     file.Name,
			Path:     file.Path,
			Size:     file.Size,
			Duration: file.Duration,
			Type:     file.Type,
		}
	}

	response := FileListResponse{
		AnimeID: fileList.AnimeID,
		Files:   files,
		Count:   len(files),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// RegisterSearchRoutes registers search routes in the chi router
func RegisterSearchRoutes(r chi.Router, searchService ports.AnimeSearch, fileListService ports.FileListService) {
	handler := NewSearchHandler(searchService)
	r.Get("/search", handler.Search)

	// Register autocomplete route
	autocompleteHandler := NewAutocompleteHandler(searchService)
	r.Get("/search/autocomplete", autocompleteHandler.handleAutocomplete)

	// Register file listing route if service provided
	if fileListService != nil {
		fileListHandler := NewFileListHandler(fileListService)
		r.Get("/anime/{animeId}/files", fileListHandler.GetFiles)
	}
}
