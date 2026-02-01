package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
)

// mockAutocompleteSearchService implements ports.AnimeSearch for testing
type mockAutocompleteSearchService struct {
	results []domain.AnimeSearchResult
	err     error
}

func (m *mockAutocompleteSearchService) Search(ctx context.Context, query string) ([]domain.AnimeSearchResult, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.results, nil
}

func TestAutocompleteHandler_ValidQuery(t *testing.T) {
	// Arrange: Create mock service with 15 results
	results := make([]domain.AnimeSearchResult, 15)
	for i := 0; i < 15; i++ {
		results[i] = domain.AnimeSearchResult{
			ID:           "anime-" + string(rune('1'+i)),
			Title:        "Test Anime " + string(rune('A'+i)),
			ThumbnailURL: "https://example.com/thumb" + string(rune('1'+i)) + ".jpg",
			Year:         2020 + i,
			Status:       "completed",
			EpisodeCount: 12,
		}
	}

	mockService := &mockAutocompleteSearchService{results: results}
	handler := NewAutocompleteHandler(mockService)

	// Act: Make request with valid query
	req := httptest.NewRequest(http.MethodGet, "/api/v1/search/autocomplete?q=test", nil)
	w := httptest.NewRecorder()
	handler.handleAutocomplete(w, req)

	// Assert: Status 200 and max 10 results
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var suggestions []AutocompleteSuggestion
	if err := json.NewDecoder(w.Body).Decode(&suggestions); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(suggestions) != 10 {
		t.Errorf("Expected 10 suggestions (limit), got %d", len(suggestions))
	}

	// Verify first suggestion has all required fields
	if suggestions[0].ID == "" || suggestions[0].Title == "" || suggestions[0].ThumbnailURL == "" || suggestions[0].Year == 0 {
		t.Errorf("First suggestion missing required fields: %+v", suggestions[0])
	}
}

func TestAutocompleteHandler_ShortQuery(t *testing.T) {
	// Arrange
	mockService := &mockAutocompleteSearchService{
		results: []domain.AnimeSearchResult{
			{ID: "anime-1", Title: "Naruto", ThumbnailURL: "thumb.jpg", Year: 2002},
		},
	}
	handler := NewAutocompleteHandler(mockService)

	// Act: Query with 1 character (< 2)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/search/autocomplete?q=n", nil)
	w := httptest.NewRecorder()
	handler.handleAutocomplete(w, req)

	// Assert: Returns empty array without calling service
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var suggestions []AutocompleteSuggestion
	if err := json.NewDecoder(w.Body).Decode(&suggestions); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(suggestions) != 0 {
		t.Errorf("Expected empty array for short query, got %d suggestions", len(suggestions))
	}
}

func TestAutocompleteHandler_EmptyQuery(t *testing.T) {
	// Arrange
	mockService := &mockAutocompleteSearchService{
		results: []domain.AnimeSearchResult{},
	}
	handler := NewAutocompleteHandler(mockService)

	// Act: Empty query parameter
	req := httptest.NewRequest(http.MethodGet, "/api/v1/search/autocomplete?q=", nil)
	w := httptest.NewRecorder()
	handler.handleAutocomplete(w, req)

	// Assert: Returns empty array
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var suggestions []AutocompleteSuggestion
	if err := json.NewDecoder(w.Body).Decode(&suggestions); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(suggestions) != 0 {
		t.Errorf("Expected empty array, got %d suggestions", len(suggestions))
	}
}

func TestAutocompleteHandler_ResponseFormat(t *testing.T) {
	// Arrange: Mock service returns full domain objects
	mockService := &mockAutocompleteSearchService{
		results: []domain.AnimeSearchResult{
			{
				ID:           "anime-1",
				Title:        "Naruto",
				ThumbnailURL: "https://cdn.example.com/naruto.jpg",
				Year:         2002,
				Status:       "completed", // Should NOT be in autocomplete response
				EpisodeCount: 220,         // Should NOT be in autocomplete response
			},
		},
	}
	handler := NewAutocompleteHandler(mockService)

	// Act
	req := httptest.NewRequest(http.MethodGet, "/api/v1/search/autocomplete?q=naruto", nil)
	w := httptest.NewRecorder()
	handler.handleAutocomplete(w, req)

	// Assert: Response only contains id, title, thumbnail_url, year
	var suggestions []AutocompleteSuggestion
	if err := json.NewDecoder(w.Body).Decode(&suggestions); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(suggestions) != 1 {
		t.Fatalf("Expected 1 suggestion, got %d", len(suggestions))
	}

	s := suggestions[0]
	if s.ID != "anime-1" || s.Title != "Naruto" || s.ThumbnailURL != "https://cdn.example.com/naruto.jpg" || s.Year != 2002 {
		t.Errorf("Suggestion has incorrect fields: %+v", s)
	}

	// Verify Content-Type header
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}
}

func TestAutocompleteHandler_ServiceError(t *testing.T) {
	// Arrange: Mock service returns error
	mockService := &mockAutocompleteSearchService{
		err: context.DeadlineExceeded,
	}
	handler := NewAutocompleteHandler(mockService)

	// Act
	req := httptest.NewRequest(http.MethodGet, "/api/v1/search/autocomplete?q=test", nil)
	w := httptest.NewRecorder()
	handler.handleAutocomplete(w, req)

	// Assert: Graceful degradation (empty array, 200 OK)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 (graceful degradation), got %d", w.Code)
	}

	var suggestions []AutocompleteSuggestion
	if err := json.NewDecoder(w.Body).Decode(&suggestions); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(suggestions) != 0 {
		t.Errorf("Expected empty array on error (graceful degradation), got %d suggestions", len(suggestions))
	}
}

func TestAutocompleteHandler_CaseInsensitive(t *testing.T) {
	// Arrange
	mockService := &mockAutocompleteSearchService{
		results: []domain.AnimeSearchResult{
			{ID: "anime-1", Title: "Naruto", ThumbnailURL: "thumb.jpg", Year: 2002},
		},
	}
	handler := NewAutocompleteHandler(mockService)

	// Act: Query with uppercase (service handles normalization)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/search/autocomplete?q=NARUTO", nil)
	w := httptest.NewRecorder()
	handler.handleAutocomplete(w, req)

	// Assert: Service is called (case sensitivity handled by AnimeSamaSearchService)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var suggestions []AutocompleteSuggestion
	if err := json.NewDecoder(w.Body).Decode(&suggestions); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(suggestions) != 1 {
		t.Errorf("Expected 1 suggestion (case insensitive), got %d", len(suggestions))
	}
}
