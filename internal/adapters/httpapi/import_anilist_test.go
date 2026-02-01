package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/app"
	"github.com/go-chi/chi/v5"
)

func TestAniListImportHandler_PreviewDisabled(t *testing.T) {
	handler := NewAniListImportHandler(nil)

	reqBody := app.AniListImportPreviewRequest{Statuses: []string{"CURRENT"}}
	body, _ := json.Marshal(reqBody)
	httpReq := httptest.NewRequest(http.MethodPost, "/import/anilist/preview", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	router := chi.NewRouter()
	handler.Routes(router)
	router.ServeHTTP(rr, httpReq)

	if rr.Code != http.StatusNotImplemented {
		t.Errorf("expected status %d, got %d", http.StatusNotImplemented, rr.Code)
	}
}

func TestAniListImportHandler_AutoDisabled(t *testing.T) {
	handler := NewAniListImportHandler(nil)

	reqBody := app.AniListImportAutoRequest{Statuses: []string{"CURRENT"}}
	body, _ := json.Marshal(reqBody)
	httpReq := httptest.NewRequest(http.MethodPost, "/import/anilist/auto", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	router := chi.NewRouter()
	handler.Routes(router)
	router.ServeHTTP(rr, httpReq)

	if rr.Code != http.StatusNotImplemented {
		t.Errorf("expected status %d, got %d", http.StatusNotImplemented, rr.Code)
	}
}

func TestAniListImportHandler_ConfirmDisabled(t *testing.T) {
	handler := NewAniListImportHandler(nil)

	reqBody := app.AniListImportConfirmRequest{Items: []app.AniListImportConfirmItem{}}
	body, _ := json.Marshal(reqBody)
	httpReq := httptest.NewRequest(http.MethodPost, "/import/anilist/confirm", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	router := chi.NewRouter()
	handler.Routes(router)
	router.ServeHTTP(rr, httpReq)

	if rr.Code != http.StatusNotImplemented {
		t.Errorf("expected status %d, got %d", http.StatusNotImplemented, rr.Code)
	}
}

func TestAniListImportHandler_PreviewInvalidJSON(t *testing.T) {
	handler := NewAniListImportHandler(nil)

	httpReq := httptest.NewRequest(http.MethodPost, "/import/anilist/preview", bytes.NewReader([]byte("bad json")))
	rr := httptest.NewRecorder()

	router := chi.NewRouter()
	handler.Routes(router)
	router.ServeHTTP(rr, httpReq)

	if rr.Code != http.StatusNotImplemented {
		t.Errorf("expected status %d, got %d", http.StatusNotImplemented, rr.Code)
	}
}
