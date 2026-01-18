package httpapi

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/adapters/sqlite"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/app"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
	"github.com/go-chi/chi/v5"
)

func TestSettingsHandler_PutUpdatesDownloadLimiter(t *testing.T) {
	ctx := context.Background()
	db, err := sqlite.Open(ctx, ":memory:")
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	repo := sqlite.NewSettingsRepository(db.SQL)
	svc := app.NewSettingsService(repo)
	lim := app.NewDynamicLimiter(1)

	h := NewSettingsHandler(svc, func(updated domain.Settings) {
		lim.SetLimit(updated.MaxConcurrentDownloads)
	})

	r := chi.NewRouter()
	h.Routes(r)

	body := []byte(`{"destination":"videos","outputNamingMode":"legacy","separateLang":false,"maxWorkers":2,"maxConcurrentDownloads":2}`)
	req := httptest.NewRequest(http.MethodPut, "/settings", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status: want %d, got %d", http.StatusOK, rr.Code)
	}
	if lim.Limit() != 2 {
		t.Fatalf("limiter limit: want %d, got %d", 2, lim.Limit())
	}
}
