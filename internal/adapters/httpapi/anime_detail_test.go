package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/app"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
)

func TestAnimeDetailHandler_ValidID(t *testing.T) {
	svc := app.NewMockAnimeDetailService()
	handler := NewAnimeDetailHandler(svc)

	router := chi.NewRouter()
	handler.Routes(router)

	req := httptest.NewRequest(http.MethodGet, "/anime/mushishi", nil)
	rec := httptest.NewRecorder()

	// Add chi URL param
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "mushishi")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	var anime domain.AnimeDetail
	if err := json.Unmarshal(rec.Body.Bytes(), &anime); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if anime.ID != "mushishi" {
		t.Errorf("Expected ID 'mushishi', got '%s'", anime.ID)
	}

	if anime.Title != "Mushishi" {
		t.Errorf("Expected title 'Mushishi', got '%s'", anime.Title)
	}

	if anime.EpisodeCount != 26 {
		t.Errorf("Expected 26 episodes, got %d", anime.EpisodeCount)
	}

	if len(anime.Seasons) != 1 {
		t.Errorf("Expected 1 season, got %d", len(anime.Seasons))
	}
}

func TestAnimeDetailHandler_InvalidID(t *testing.T) {
	svc := app.NewMockAnimeDetailService()
	handler := NewAnimeDetailHandler(svc)

	router := chi.NewRouter()
	handler.Routes(router)

	req := httptest.NewRequest(http.MethodGet, "/anime/invalid-id", nil)
	rec := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "invalid-id")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", rec.Code)
	}

	var errResp map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("Failed to decode error response: %v", err)
	}

	if errResp["error"] == "" {
		t.Error("Expected error message in response")
	}
}

func TestAnimeDetailHandler_EmptyID(t *testing.T) {
	svc := app.NewMockAnimeDetailService()
	handler := NewAnimeDetailHandler(svc)

	router := chi.NewRouter()
	handler.Routes(router)

	req := httptest.NewRequest(http.MethodGet, "/anime/", nil)
	rec := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", rec.Code)
	}
}

func TestAnimeDetailHandler_ResponseFormat(t *testing.T) {
	svc := app.NewMockAnimeDetailService()
	handler := NewAnimeDetailHandler(svc)

	router := chi.NewRouter()
	handler.Routes(router)

	req := httptest.NewRequest(http.MethodGet, "/anime/mushishi", nil)
	rec := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "mushishi")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	router.ServeHTTP(rec, req)

	// Check Content-Type header
	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}

	// Validate JSON structure
	var anime domain.AnimeDetail
	if err := json.Unmarshal(rec.Body.Bytes(), &anime); err != nil {
		t.Fatalf("Response is not valid JSON: %v", err)
	}

	// Check required fields
	if anime.ID == "" {
		t.Error("Missing 'id' field")
	}
	if anime.Title == "" {
		t.Error("Missing 'title' field")
	}
	if anime.ThumbnailURL == "" {
		t.Error("Missing 'thumbnail_url' field")
	}
	if anime.Synopsis == "" {
		t.Error("Missing 'synopsis' field")
	}
	if anime.Year == 0 {
		t.Error("Missing 'year' field")
	}
	if anime.Status == "" {
		t.Error("Missing 'status' field")
	}
	if len(anime.Genres) == 0 {
		t.Error("Missing 'genres' field")
	}
	if anime.EpisodeCount == 0 {
		t.Error("Missing 'episode_count' field")
	}
	if len(anime.Seasons) == 0 {
		t.Error("Missing 'seasons' field")
	}
}

func TestAnimeDetailHandler_MultipleSeasons(t *testing.T) {
	svc := app.NewMockAnimeDetailService()
	handler := NewAnimeDetailHandler(svc)

	router := chi.NewRouter()
	handler.Routes(router)

	// Test with samurai-champloo which has 1 season with 5 episodes
	req := httptest.NewRequest(http.MethodGet, "/anime/samurai-champloo", nil)
	rec := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "samurai-champloo")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	var anime domain.AnimeDetail
	if err := json.Unmarshal(rec.Body.Bytes(), &anime); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if anime.ID != "samurai-champloo" {
		t.Errorf("Expected ID 'samurai-champloo', got '%s'", anime.ID)
	}

	if len(anime.Seasons) != 1 {
		t.Errorf("Expected 1 season, got %d", len(anime.Seasons))
	}

	// Verify season structure
	season1 := anime.Seasons[0]
	if season1.Number != 1 {
		t.Errorf("Expected season 1, got %d", season1.Number)
	}
	if len(season1.Episodes) != 5 {
		t.Errorf("Expected 5 episodes in season 1, got %d", len(season1.Episodes))
	}
}

func TestAnimeDetailHandler_SingleSeason(t *testing.T) {
	svc := app.NewMockAnimeDetailService()
	handler := NewAnimeDetailHandler(svc)

	router := chi.NewRouter()
	handler.Routes(router)

	req := httptest.NewRequest(http.MethodGet, "/anime/dororo", nil)
	rec := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "dororo")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	var anime domain.AnimeDetail
	if err := json.Unmarshal(rec.Body.Bytes(), &anime); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if anime.ID != "dororo" {
		t.Errorf("Expected ID 'dororo', got '%s'", anime.ID)
	}

	if len(anime.Seasons) != 1 {
		t.Errorf("Expected 1 season, got %d", len(anime.Seasons))
	}

	season := anime.Seasons[0]
	if season.Number != 1 {
		t.Errorf("Expected season number 1, got %d", season.Number)
	}
	if season.Name != "Saison 1" {
		t.Errorf("Expected season name 'Saison 1', got '%s'", season.Name)
	}
	if len(season.Episodes) != 4 {
		t.Errorf("Expected 4 episodes (mock fixture), got %d", len(season.Episodes))
	}

	// Verify first episode
	ep1 := season.Episodes[0]
	if ep1.Number != 1 {
		t.Errorf("Expected episode number 1, got %d", ep1.Number)
	}
	if ep1.SeasonNumber != 1 {
		t.Errorf("Expected season_number 1, got %d", ep1.SeasonNumber)
	}
	if ep1.URL == "" {
		t.Error("Episode should have a URL")
	}
}

func TestAnimeDetailHandler_OngoingStatus(t *testing.T) {
	svc := app.NewMockAnimeDetailService()
	handler := NewAnimeDetailHandler(svc)

	router := chi.NewRouter()
	handler.Routes(router)

	// Test with natsume-yuujinchou which is ongoing
	req := httptest.NewRequest(http.MethodGet, "/anime/natsume-yuujinchou", nil)
	rec := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "natsume-yuujinchou")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	var anime domain.AnimeDetail
	if err := json.Unmarshal(rec.Body.Bytes(), &anime); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if anime.Status != "ongoing" {
		t.Errorf("Expected status 'ongoing', got '%s'", anime.Status)
	}

	if anime.EpisodeCount != 13 {
		t.Errorf("Expected 13 episodes, got %d", anime.EpisodeCount)
	}
}
