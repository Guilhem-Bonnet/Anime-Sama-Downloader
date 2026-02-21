package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/app"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
)

// Tests
func TestAniListHandler_ViewerDisabled(t *testing.T) {
	handler := NewAniListHandler(nil)

	httpReq := httptest.NewRequest(http.MethodGet, "/anilist/viewer", nil)
	rr := httptest.NewRecorder()

	router := chi.NewRouter()
	handler.Routes(router)
	router.ServeHTTP(rr, httpReq)

	if rr.Code != http.StatusNotImplemented {
		t.Errorf("expected status %d, got %d", http.StatusNotImplemented, rr.Code)
	}
}

func TestAniListHandler_AiringDisabled(t *testing.T) {
	handler := NewAniListHandler(nil)

	httpReq := httptest.NewRequest(http.MethodGet, "/anilist/airing", nil)
	rr := httptest.NewRecorder()

	router := chi.NewRouter()
	handler.Routes(router)
	router.ServeHTTP(rr, httpReq)

	if rr.Code != http.StatusNotImplemented {
		t.Errorf("expected status %d, got %d", http.StatusNotImplemented, rr.Code)
	}
}

func TestAniListHandler_WatchlistDisabled(t *testing.T) {
	handler := NewAniListHandler(nil)

	httpReq := httptest.NewRequest(http.MethodGet, "/anilist/watchlist", nil)
	rr := httptest.NewRecorder()

	router := chi.NewRouter()
	handler.Routes(router)
	router.ServeHTTP(rr, httpReq)

	if rr.Code != http.StatusNotImplemented {
		t.Errorf("expected status %d, got %d", http.StatusNotImplemented, rr.Code)
	}
}

func TestAniListHandler_ViewerNotConfigured(t *testing.T) {
	svc := app.NewAniListService(func(ctx context.Context) (domain.Settings, error) {
		return domain.Settings{AniListToken: ""}, nil
	})
	handler := NewAniListHandler(svc)

	httpReq := httptest.NewRequest(http.MethodGet, "/anilist/viewer", nil)
	rr := httptest.NewRecorder()

	router := chi.NewRouter()
	handler.Routes(router)
	router.ServeHTTP(rr, httpReq)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}
