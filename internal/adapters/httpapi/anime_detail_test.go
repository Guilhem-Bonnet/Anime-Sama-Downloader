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

	req := httptest.NewRequest(http.MethodGet, "/anime/naruto", nil)
	rec := httptest.NewRecorder()

	// Add chi URL param
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "naruto")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	var anime domain.AnimeDetail
	if err := json.Unmarshal(rec.Body.Bytes(), &anime); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if anime.ID != "naruto" {
		t.Errorf("Expected ID 'naruto', got '%s'", anime.ID)
	}

	if anime.Title != "Naruto" {
		t.Errorf("Expected title 'Naruto', got '%s'", anime.Title)
	}

	if anime.EpisodeCount != 220 {
		t.Errorf("Expected 220 episodes, got %d", anime.EpisodeCount)
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

	req := httptest.NewRequest(http.MethodGet, "/anime/naruto", nil)
	rec := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "naruto")
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

	// Test with Naruto Shippuden which has 2 seasons
	req := httptest.NewRequest(http.MethodGet, "/anime/naruto-shippuden", nil)
	rec := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "naruto-shippuden")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	var anime domain.AnimeDetail
	if err := json.Unmarshal(rec.Body.Bytes(), &anime); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if anime.ID != "naruto-shippuden" {
		t.Errorf("Expected ID 'naruto-shippuden', got '%s'", anime.ID)
	}

	if len(anime.Seasons) != 2 {
		t.Errorf("Expected 2 seasons, got %d", len(anime.Seasons))
	}

	// Verify season structure
	season1 := anime.Seasons[0]
	if season1.Number != 1 {
		t.Errorf("Expected season 1, got %d", season1.Number)
	}
	if len(season1.Episodes) == 0 {
		t.Error("Season 1 should have episodes")
	}

	season2 := anime.Seasons[1]
	if season2.Number != 2 {
		t.Errorf("Expected season 2, got %d", season2.Number)
	}
	if len(season2.Episodes) == 0 {
		t.Error("Season 2 should have episodes")
	}
}

func TestAnimeDetailHandler_SingleSeason(t *testing.T) {
	svc := app.NewMockAnimeDetailService()
	handler := NewAnimeDetailHandler(svc)

	router := chi.NewRouter()
	handler.Routes(router)

	req := httptest.NewRequest(http.MethodGet, "/anime/attack-on-titan", nil)
	rec := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "attack-on-titan")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	var anime domain.AnimeDetail
	if err := json.Unmarshal(rec.Body.Bytes(), &anime); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if anime.ID != "attack-on-titan" {
		t.Errorf("Expected ID 'attack-on-titan', got '%s'", anime.ID)
	}

	if len(anime.Seasons) != 1 {
		t.Errorf("Expected 1 season, got %d", len(anime.Seasons))
	}

	season := anime.Seasons[0]
	if season.Number != 1 {
		t.Errorf("Expected season number 1, got %d", season.Number)
	}
	if season.Name != "Season 1" {
		t.Errorf("Expected season name 'Season 1', got '%s'", season.Name)
	}
	if len(season.Episodes) != 2 {
		t.Errorf("Expected 2 episodes (mock fixture), got %d", len(season.Episodes))
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

	// Test with One Piece which is ongoing
	req := httptest.NewRequest(http.MethodGet, "/anime/one-piece", nil)
	rec := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "one-piece")
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

	if anime.EpisodeCount != 1100 {
		t.Errorf("Expected 1100 episodes, got %d", anime.EpisodeCount)
	}
}
