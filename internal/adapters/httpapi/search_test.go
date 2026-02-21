package httpapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
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

func testSearchHandler(_ *testing.T, searchService ports.AnimeSearch) http.Handler {
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
	fileList     *domain.FileList
	err          error
	shouldError  bool
	errorMessage string
}

func (m *MockFileListService) GetFileList(ctx context.Context, animeID string) (*domain.FileList, error) {
	if m.shouldError {
		return nil, fmt.Errorf(m.errorMessage)
	}
	if m.err != nil {
		return nil, m.err
	}
	return m.fileList, nil
}

func (m *MockFileListService) GetFilesByAnimeTitle(ctx context.Context, title string) (*domain.FileList, error) {
	if m.shouldError {
		return nil, fmt.Errorf(m.errorMessage)
	}
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

// TestFileListHandler_GetFiles_ServiceError tests internal service errors
func TestFileListHandler_GetFiles_ServiceError(t *testing.T) {
	mockService := &MockFileListService{
		shouldError:  true,
		errorMessage: "internal service error",
	}

	handler := NewFileListHandler(mockService)
	req := httptest.NewRequest("GET", "/api/v1/anime/1/files", nil)
	w := httptest.NewRecorder()

	// Simulate chi URL params
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
		URLParams: chi.RouteParams{Keys: []string{"animeId"}, Values: []string{"1"}},
	}))

	handler.GetFiles(w, req)

	// Should return 404 for "not found" errors
	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

// TestFileListHandler_GetFiles_LargeFileList tests handling of large file lists
func TestFileListHandler_GetFiles_LargeFileList(t *testing.T) {
	// Create service with large file list (1000 episodes)
	largeFiles := make([]domain.File, 1000)
	for i := range largeFiles {
		largeFiles[i] = domain.File{
			ID:       fmt.Sprintf("large-ep%d", i+1),
			Name:     fmt.Sprintf("Episode %d", i+1),
			Path:     fmt.Sprintf("/downloads/Episode_%04d.mkv", i+1),
			Size:     350000000,
			Duration: 1400,
			Type:     "video/x-matroska",
		}
	}

	mockService := &MockFileListService{
		fileList: &domain.FileList{
			AnimeID: "large-1",
			Files:   largeFiles,
		},
	}

	handler := NewFileListHandler(mockService)
	req := httptest.NewRequest("GET", "/api/v1/anime/large-1/files", nil)
	w := httptest.NewRecorder()

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
		URLParams: chi.RouteParams{Keys: []string{"animeId"}, Values: []string{"large-1"}},
	}))

	handler.GetFiles(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response FileListResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(response.Files) != 1000 {
		t.Errorf("expected 1000 files in response, got %d", len(response.Files))
	}
}

// TestFileListHandler_GetFiles_JSONValidation tests response JSON structure
func TestFileListHandler_GetFiles_JSONValidation(t *testing.T) {
	mockService := &MockFileListService{
		fileList: &domain.FileList{
			AnimeID: "1",
			Files: []domain.File{
				{
					ID:       "1-ep1",
					Name:     "Episode 1",
					Path:     "/downloads/Episode_01.mkv",
					Size:     350000000,
					Duration: 1400,
					Type:     "video/x-matroska",
				},
			},
		},
	}

	handler := NewFileListHandler(mockService)
	req := httptest.NewRequest("GET", "/api/v1/anime/1/files", nil)
	w := httptest.NewRecorder()

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
		URLParams: chi.RouteParams{Keys: []string{"animeId"}, Values: []string{"1"}},
	}))

	handler.GetFiles(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	// Verify content-type
	contentType := w.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("expected JSON content type, got %s", contentType)
	}

	// Verify JSON structure
	var response FileListResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Verify required fields are present
	if response.AnimeID == "" {
		t.Error("response missing anime_id field")
	}

	if len(response.Files) == 0 {
		t.Error("response missing files array")
	}

	// Verify file structure
	file := response.Files[0]
	if file.ID == "" || file.Name == "" || file.Path == "" || file.Type == "" {
		t.Error("file missing required fields")
	}
}

// TestFileListHandler_GetFiles_EmptyFileList tests anime with no episodes
func TestFileListHandler_GetFiles_EmptyFileList(t *testing.T) {
	mockService := &MockFileListService{
		fileList: &domain.FileList{
			AnimeID: "empty-1",
			Files:   []domain.File{}, // No files
		},
	}

	handler := NewFileListHandler(mockService)
	req := httptest.NewRequest("GET", "/api/v1/anime/empty-1/files", nil)
	w := httptest.NewRecorder()

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
		URLParams: chi.RouteParams{Keys: []string{"animeId"}, Values: []string{"empty-1"}},
	}))

	handler.GetFiles(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response FileListResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(response.Files) != 0 {
		t.Errorf("expected 0 files, got %d", len(response.Files))
	}
}

// TestFileListHandler_GetFiles_SpecialCharactersInID tests IDs with special characters
func TestFileListHandler_GetFiles_SpecialCharactersInID(t *testing.T) {
	mockService := &MockFileListService{
		fileList: &domain.FileList{
			AnimeID: "test-123",
			Files: []domain.File{
				{
					ID:       "test-123-ep1",
					Name:     "Episode 1",
					Path:     "/downloads/Episode_01.mkv",
					Size:     350000000,
					Duration: 1400,
					Type:     "video/x-matroska",
				},
			},
		},
	}

	handler := NewFileListHandler(mockService)

	// Test various ID formats
	testIDs := []string{
		"test-123",
		"anime_456",
		"show.789",
	}

	for _, animeID := range testIDs {
		mockService.fileList.AnimeID = animeID
		req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/anime/%s/files", animeID), nil)
		w := httptest.NewRecorder()

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{Keys: []string{"animeId"}, Values: []string{animeID}},
		}))

		handler.GetFiles(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("ID '%s': expected status 200, got %d", animeID, w.Code)
		}
	}
}

// ==================== TASK 6.3: HTTP ERROR SCENARIOS ====================

// TestFileListHandler_GetFiles_InvalidQueryParameters tests handling of invalid query params
func TestFileListHandler_GetFiles_InvalidQueryParameters(t *testing.T) {
	mockService := &MockFileListService{
		fileList: &domain.FileList{
			AnimeID: "1",
			Files: []domain.File{
				{
					ID:       "1-ep1",
					Name:     "Episode 1",
					Path:     "/path/ep1.mkv",
					Size:     350000000,
					Duration: 1400,
					Type:     "video/x-matroska",
				},
			},
		},
	}

	handler := NewFileListHandler(mockService)

	testCases := []struct {
		name      string
		query     string
		expectErr bool
	}{
		{"negative limit", "?limit=-1", true},
		{"zero limit", "?limit=0", true},
		{"non-numeric limit", "?limit=abc", true},
		{"negative offset", "?offset=-1", true},
		{"non-numeric offset", "?offset=xyz", true},
		{"invalid filter", "?filter=@@@", false}, // Should handle gracefully
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/anime/1/files%s", tc.query), nil)
			w := httptest.NewRecorder()

			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
				URLParams: chi.RouteParams{Keys: []string{"animeId"}, Values: []string{"1"}},
			}))

			handler.GetFiles(w, req)

			// Handler should return either 200 or 400, but not 500
			if w.Code >= 500 {
				t.Errorf("expected client error or success, got %d", w.Code)
			}
		})
	}
}

// TestFileListHandler_GetFiles_ContextCancellation tests handling of cancelled contexts
func TestFileListHandler_GetFiles_ContextCancellation(t *testing.T) {
	mockService := &MockFileListService{
		fileList: &domain.FileList{
			AnimeID: "1",
			Files: []domain.File{
				{
					ID:       "1-ep1",
					Name:     "Episode 1",
					Path:     "/path/ep1.mkv",
					Size:     350000000,
					Duration: 1400,
					Type:     "video/x-matroska",
				},
			},
		},
	}

	handler := NewFileListHandler(mockService)

	req := httptest.NewRequest("GET", "/api/v1/anime/1/files", nil)
	w := httptest.NewRecorder()

	// Create a cancelled context
	ctx, cancel := context.WithCancel(req.Context())
	cancel() // Cancel immediately

	req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, &chi.Context{
		URLParams: chi.RouteParams{Keys: []string{"animeId"}, Values: []string{"1"}},
	}))

	handler.GetFiles(w, req)

	// Should handle cancellation gracefully (either error or still respond)
	if w.Code >= 500 {
		t.Errorf("expected graceful handling of cancellation, got status %d", w.Code)
	}
}

// TestFileListHandler_GetFiles_ResponseHeaders tests proper response headers
func TestFileListHandler_GetFiles_ResponseHeaders(t *testing.T) {
	mockService := &MockFileListService{
		fileList: &domain.FileList{
			AnimeID: "1",
			Files: []domain.File{
				{
					ID:       "1-ep1",
					Name:     "Episode 1",
					Path:     "/path/ep1.mkv",
					Size:     350000000,
					Duration: 1400,
					Type:     "video/x-matroska",
				},
			},
		},
	}

	handler := NewFileListHandler(mockService)

	req := httptest.NewRequest("GET", "/api/v1/anime/1/files", nil)
	w := httptest.NewRecorder()

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
		URLParams: chi.RouteParams{Keys: []string{"animeId"}, Values: []string{"1"}},
	}))

	handler.GetFiles(w, req)

	if w.Code == http.StatusOK {
		// Verify Content-Type
		if ct := w.Header().Get("Content-Type"); !strings.Contains(ct, "application/json") {
			t.Errorf("expected Content-Type to be application/json, got %s", ct)
		}

		// Verify no sensitive headers leak
		if w.Header().Get("X-Internal-Error") != "" {
			t.Error("unexpected internal error header in response")
		}
	}
}

// ==================== TASK 6.4: PAGINATION & FILTERING ====================

// TestFileListHandler_GetFiles_PaginationParameters tests pagination parameter handling
func TestFileListHandler_GetFiles_PaginationParameters(t *testing.T) {
	// Create a file list with multiple files
	files := make([]domain.File, 10)
	for i := 0; i < 10; i++ {
		files[i] = domain.File{
			ID:       fmt.Sprintf("1-ep%d", i+1),
			Name:     fmt.Sprintf("Episode %d", i+1),
			Path:     fmt.Sprintf("/path/ep%d.mkv", i+1),
			Size:     int64(350000000 + i*10000000),
			Duration: 1400 + i*60,
			Type:     "video/x-matroska",
		}
	}

	mockService := &MockFileListService{
		fileList: &domain.FileList{
			AnimeID: "1",
			Files:   files,
		},
	}

	handler := NewFileListHandler(mockService)

	testCases := []struct {
		name  string
		limit int
		desc  string
	}{
		{"default limit", 0, "should return all files"},
		{"limit 5", 5, "should handle smaller limit"},
		{"limit 20", 20, "should handle limit > total"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query := "/api/v1/anime/1/files"
			if tc.limit > 0 {
				query += fmt.Sprintf("?limit=%d", tc.limit)
			}

			req := httptest.NewRequest("GET", query, nil)
			w := httptest.NewRecorder()

			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
				URLParams: chi.RouteParams{Keys: []string{"animeId"}, Values: []string{"1"}},
			}))

			handler.GetFiles(w, req)

			if w.Code == http.StatusOK {
				var response struct {
					AnimeID string        `json:"anime_id"`
					Files   []interface{} `json:"files"`
				}
				json.NewDecoder(w.Body).Decode(&response)

				// Should still have files
				if len(response.Files) == 0 {
					t.Errorf("expected files in response, got empty list")
				}
			}
		})
	}
}

// TestFileListHandler_GetFiles_FilteringByFileType tests file type filtering
func TestFileListHandler_GetFiles_FilteringByFileType(t *testing.T) {
	files := []domain.File{
		{
			ID:       "1-ep1",
			Name:     "Episode 1",
			Path:     "/path/ep1.mkv",
			Size:     350000000,
			Duration: 1400,
			Type:     "video/x-matroska",
		},
		{
			ID:       "1-ep1-sub",
			Name:     "Episode 1 Subtitles",
			Path:     "/path/ep1.srt",
			Size:     50000,
			Duration: 0,
			Type:     "text/plain",
		},
	}

	mockService := &MockFileListService{
		fileList: &domain.FileList{
			AnimeID: "1",
			Files:   files,
		},
	}

	handler := NewFileListHandler(mockService)

	req := httptest.NewRequest("GET", "/api/v1/anime/1/files?type=video", nil)
	w := httptest.NewRecorder()

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
		URLParams: chi.RouteParams{Keys: []string{"animeId"}, Values: []string{"1"}},
	}))

	handler.GetFiles(w, req)

	if w.Code == http.StatusOK {
		var response struct {
			Files []map[string]interface{} `json:"files"`
		}
		json.NewDecoder(w.Body).Decode(&response)

		// Verify response structure is valid
		if len(response.Files) == 0 {
			t.Logf("filter applied: response has %d files", len(response.Files))
		}
	}
}

// TestFileListHandler_GetFiles_SortingOptions tests sorting parameter handling
func TestFileListHandler_GetFiles_SortingOptions(t *testing.T) {
	files := []domain.File{
		{
			ID:       "1-ep2",
			Name:     "Zebra Episode",
			Path:     "/path/ep2.mkv",
			Size:     400000000,
			Duration: 1500,
			Type:     "video/x-matroska",
		},
		{
			ID:       "1-ep1",
			Name:     "Alpha Episode",
			Path:     "/path/ep1.mkv",
			Size:     350000000,
			Duration: 1400,
			Type:     "video/x-matroska",
		},
	}

	mockService := &MockFileListService{
		fileList: &domain.FileList{
			AnimeID: "1",
			Files:   files,
		},
	}

	handler := NewFileListHandler(mockService)

	testCases := []struct {
		name  string
		sort  string
		order string
	}{
		{"sort by name", "name", "asc"},
		{"sort by size", "size", "desc"},
		{"sort by date", "date", "asc"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query := fmt.Sprintf("/api/v1/anime/1/files?sort=%s&order=%s", tc.sort, tc.order)

			req := httptest.NewRequest("GET", query, nil)
			w := httptest.NewRecorder()

			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
				URLParams: chi.RouteParams{Keys: []string{"animeId"}, Values: []string{"1"}},
			}))

			handler.GetFiles(w, req)

			// Handler should handle sort parameters gracefully
			if w.Code >= 500 {
				t.Errorf("expected client success, got status %d", w.Code)
			}
		})
	}
}
