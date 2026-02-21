package app

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
)

func TestAniListService_Watchlist_UsesTokenAndFlattensEntries(t *testing.T) {
	calls := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.Header().Set("Content-Type", "application/json")
		if calls == 1 {
			_, _ = w.Write([]byte(`{"data":{"Viewer":{"id":42,"name":"me"}}}`))
			return
		}
		_, _ = w.Write([]byte(`{"data":{"MediaListCollection":{"lists":[{"entries":[{"status":"CURRENT","progress":3,"media":{"id":123,"synonyms":["Solo"],"title":{"romaji":"Solo Leveling","english":"Solo Leveling","native":""}}}]}]}}}`))
	}))
	defer ts.Close()

	svc := NewAniListService(func(ctx context.Context) (domain.Settings, error) {
		return domain.Settings{AniListToken: "tok"}, nil
	}).WithEndpoint(ts.URL)

	entries, err := svc.Watchlist(context.Background(), []string{"CURRENT"})
	if err != nil {
		t.Fatalf("watchlist: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Media.ID != 123 {
		t.Fatalf("unexpected media id: %d", entries[0].Media.ID)
	}
}
