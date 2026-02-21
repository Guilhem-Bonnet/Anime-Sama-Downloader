package app

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
)

func TestAniListService_Viewer_RequiresToken(t *testing.T) {
	svc := NewAniListService(func(ctx context.Context) (domain.Settings, error) {
		return domain.Settings{}, nil
	})
	_, err := svc.Viewer(context.Background())
	if err != ErrAniListNotConfigured {
		t.Fatalf("expected ErrAniListNotConfigured, got %v", err)
	}
}

func TestAniListService_Viewer_SendsBearerToken(t *testing.T) {
	var gotAuth string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"Viewer":{"id":123,"name":"Guilhem"}}}`))
	}))
	defer ts.Close()

	svc := NewAniListService(func(ctx context.Context) (domain.Settings, error) {
		return domain.Settings{AniListToken: "tok"}, nil
	}).WithEndpoint(ts.URL)

	viewer, err := svc.Viewer(context.Background())
	if err != nil {
		t.Fatalf("viewer: %v", err)
	}
	if viewer.ID != 123 || viewer.Name != "Guilhem" {
		t.Fatalf("unexpected viewer: %+v", viewer)
	}
	if !strings.HasPrefix(gotAuth, "Bearer ") {
		t.Fatalf("expected Bearer auth, got %q", gotAuth)
	}
}

func TestAniListService_AiringSchedule_WorksWithoutToken(t *testing.T) {
	var gotAuth string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"Page":{"airingSchedules":[]}}}`))
	}))
	defer ts.Close()

	svc := NewAniListService(func(ctx context.Context) (domain.Settings, error) {
		return domain.Settings{}, nil
	}).WithEndpoint(ts.URL)

	_, err := svc.AiringSchedule(context.Background(), time.Now().Add(-time.Hour), time.Now().Add(time.Hour), 10)
	if err != nil {
		t.Fatalf("airing: %v", err)
	}
	if gotAuth != "" {
		t.Fatalf("expected no Authorization header, got %q", gotAuth)
	}
}
