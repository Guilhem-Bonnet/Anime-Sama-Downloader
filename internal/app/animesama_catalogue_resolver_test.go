package app

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAnimeSamaCatalogueResolver_ResolveCandidates(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only one slug exists in this fake server.
		if r.URL.Path == "/catalogue/solo-leveling/" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	r := NewAnimeSamaCatalogueResolver().WithBaseURL(ts.URL)
	cands, err := r.ResolveCandidates(context.Background(), []string{"Solo Leveling"}, 3)
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if len(cands) != 1 {
		t.Fatalf("expected 1 candidate, got %d", len(cands))
	}
	if cands[0].Slug != "solo-leveling" {
		t.Fatalf("unexpected slug: %q", cands[0].Slug)
	}
}
