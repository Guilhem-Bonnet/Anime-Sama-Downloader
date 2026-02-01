package httpapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/ports"
)

// MockAnimeSearchService is a mock implementation of ports.AnimeSearch for testing
type MockAnimeSearchService struct {
	results []domain.AnimeSearchResult
	err     error
}

func (m *MockAnimeSearchService) Search(ctx context.Context, query string) ([]domain.AnimeSearchResult, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.results, nil
}

func (m *MockAnimeSearchService) SearchWithFilters(ctx context.Context, filters ports.SearchFilters) ([]domain.AnimeSearchResult, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.results, nil
}

func testSearchHandler(t *testing.T, searchService ports.AnimeSearch) http.Handler {
	return http.HandlerFunc(NewSearchHandler(searchService).Search)
}

// TestSearchHandler_ValidQuery returns search results
func TestSearchHandler_ValidQuery(t *testing.T) {
	mockService := &MockAnimeSearchService{
		results: []domain.AnimeSearchResult{
			{
				ID:           "1",
				Title:        "Naruto",
				ThumbnailURL: "https://example.com/naruto.jpg",
				Year:         2002,
				Status:       "completed",
				EpisodeCount: 220,
			},
		},
	}

	handler := testSearchHandler(t, mockService)
	req := httptest.NewRequest("GET", "/search?q=naruto", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}

	var results []SearchResponse
	err := json.NewDecoder(w.Body).Decode(&results)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}

	if results[0].Title != "Naruto" {
		t.Errorf("expected title 'Naruto', got %q", results[0].Title)
	}
}

// TestSearchHandler_EmptyQuery returns empty array
func TestSearchHandler_EmptyQuery(t *testing.T) {
	mockService := &MockAnimeSearchService{
		results: []domain.AnimeSearchResult{},
	}

	handler := testSearchHandler(t, mockService)
	req := httptest.NewRequest("GET", "/search?q=", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var results []SearchResponse
	json.NewDecoder(w.Body).Decode(&results)

	if results == nil || len(results) != 0 {
		t.Errorf("expected empty array, got %v", results)
	}
}

// TestSearchHandler_MissingQueryParam returns empty array
func TestSearchHandler_MissingQueryParam(t *testing.T) {
	mockService := &MockAnimeSearchService{
		results: []domain.AnimeSearchResult{},
	}

	handler := testSearchHandler(t, mockService)
	req := httptest.NewRequest("GET", "/search", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var results []SearchResponse
	json.NewDecoder(w.Body).Decode(&results)

	if results == nil || len(results) != 0 {
		t.Errorf("expected empty array, got %v", results)
	}
}

// TestSearchHandler_ResponseFormat verifies JSON response format
func TestSearchHandler_ResponseFormat(t *testing.T) {
	mockService := &MockAnimeSearchService{
		results: []domain.AnimeSearchResult{
			{
				ID:           "1",
				Title:        "Naruto",
				ThumbnailURL: "https://example.com/naruto.jpg",
				Year:         2002,
				Status:       "completed",
				EpisodeCount: 220,
			},
		},
	}

	handler := testSearchHandler(t, mockService)
	req := httptest.NewRequest("GET", "/search?q=naruto", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	var results []SearchResponse
	json.NewDecoder(w.Body).Decode(&results)

	if len(results) == 0 {
		t.Fatal("expected at least one result")
	}

	// Check all fields are present
	result := results[0]
	if result.ID == "" || result.Title == "" || result.ThumbnailURL == "" || result.Status == "" {
		t.Error("missing required fields in response")
	}

	if result.Year == 0 || result.EpisodeCount == 0 {
		t.Error("expected numeric fields to be non-zero")
	}
}

// TestSearchHandler_MultipleResults returns multiple results
func TestSearchHandler_MultipleResults(t *testing.T) {
	mockService := &MockAnimeSearchService{
		results: []domain.AnimeSearchResult{
			{
				ID:           "1",
				Title:        "Naruto",
				ThumbnailURL: "https://example.com/naruto.jpg",
				Year:         2002,
				Status:       "completed",
				EpisodeCount: 220,
			},
			{
				ID:           "2",
				Title:        "Naruto Shippuden",
				ThumbnailURL: "https://example.com/naruto-shippuden.jpg",
				Year:         2007,
				Status:       "completed",
				EpisodeCount: 500,
			},
		},
	}

	handler := testSearchHandler(t, mockService)
	req := httptest.NewRequest("GET", "/search?q=naruto", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	var results []SearchResponse
	json.NewDecoder(w.Body).Decode(&results)

	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}

// TestSearchHandler_ServiceError handles service errors gracefully
func TestSearchHandler_ServiceError(t *testing.T) {
	mockService := &MockAnimeSearchService{
		err: fmt.Errorf("search error"),
	}

	handler := testSearchHandler(t, mockService)
	req := httptest.NewRequest("GET", "/search?q=naruto", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}

	var errResp map[string]string
	json.NewDecoder(w.Body).Decode(&errResp)

	if errResp["error"] != "search error" {
		t.Errorf("expected error message, got %v", errResp)
	}
}

// Mock FileListService for testing
type MockFileListService struct {
	fileList *domain.FileList
	err      error
}

func (m *MockFileListService) GetFileList(ctx context.Context, animeID string) (*domain.FileList, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.fileList, nil
}

func (m *MockFileListService) GetFilesByAnimeTitle(ctx context.Context, title string) (*domain.FileList, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.fileList, nil
}

// TestFileListHandler_GetFiles_Success tests successful file list retrieval
func TestFileListHandler_GetFiles_Success(t *testing.T) {
	mockService := &MockFileListService{
		fileList: &domain.FileList{
			AnimeID: "1",
			Files: []domain.File{
				{
					ID:       "1-ep1",
					Name:     "Naruto - Episode 1",
					Path:     "/downloads/Naruto/Episode_01.mkv",
					Size:     350000000,
					Duration: 1400,
					Type:     "video/x-matroska",
				},
				{
					ID:       "1-ep2",
					Name:     "Naruto - Episode 2",
					Path:     "/downloads/Naruto/Episode_02.mkv",
					Size:     360000000,
					Duration: 1460,
					Type:     "video/x-matroska",
				},
			},
		},
	}

	handler := NewFileListHandler(mockService)
	req := httptest.NewRequest("GET", "/api/v1/anime/1/files", nil)
	w := httptest.NewRecorder()

	// Simulate chi URL params
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
		URLParams: chi.RouteParams{Keys: []string{"animeId"}, Values: []string{"1"}},
	}))

	handler.GetFiles(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp FileListResponse
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.AnimeID != "1" {
		t.Errorf("expected anime ID '1', got '%s'", resp.AnimeID)
	}

	if len(resp.Files) != 2 {
		t.Errorf("expected 2 files, got %d", len(resp.Files))
	}

	if resp.Count != 2 {
		t.Errorf("expected count 2, got %d", resp.Count)
	}
}

// TestFileListHandler_GetFiles_NotFound tests 404 handling
func TestFileListHandler_GetFiles_NotFound(t *testing.T) {
	mockService := &MockFileListService{
		err: fmt.Errorf("anime not found"),
	}

	handler := NewFileListHandler(mockService)
	req := httptest.NewRequest("GET", "/api/v1/anime/999/files", nil)
	w := httptest.NewRecorder()

	// Simulate chi URL params
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
		URLParams: chi.RouteParams{Keys: []string{"animeId"}, Values: []string{"999"}},
	}))

	handler.GetFiles(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

// TestFileListHandler_GetFiles_NoAnimeId tests missing anime ID
func TestFileListHandler_GetFiles_NoAnimeId(t *testing.T) {
	mockService := &MockFileListService{}

	handler := NewFileListHandler(mockService)
	req := httptest.NewRequest("GET", "/api/v1/anime//files", nil)
	w := httptest.NewRecorder()

	// Simulate chi URL params with empty animeId
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
		URLParams: chi.RouteParams{Keys: []string{"animeId"}, Values: []string{""}},
	}))

	handler.GetFiles(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}
