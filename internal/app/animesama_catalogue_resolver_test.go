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

// TestResolve_SubtitleStripping valide que "Hell's Paradise: Jigokuraku" → slot hells-paradise.
// Régresse le bug où le titre complet AniList (avec sous-titre après ":") empêchait
// de trouver le slug anime-sama qui ne contient que le titre principal.
func TestResolve_SubtitleStripping(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/catalogue/hells-paradise/" {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	cases := []struct {
		title    string
		wantSlug string
	}{
		{"Hell's Paradise: Jigokuraku", "hells-paradise"},
		{"Shingeki no Kyojin: The Final Season", "shingeki-no-kyojin"},
	}

	for _, tc := range cases {
		tsLocal := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// slug must end with the expected value
			if r.URL.Path == "/catalogue/"+tc.wantSlug+"/" {
				w.WriteHeader(http.StatusOK)
				return
			}
			w.WriteHeader(http.StatusNotFound)
		}))

		res := NewAnimeSamaCatalogueResolver().WithBaseURL(tsLocal.URL)
		cands, err := res.ResolveCandidates(context.Background(), []string{tc.title}, 3)
		tsLocal.Close()
		if err != nil {
			t.Errorf("%q: resolve error: %v", tc.title, err)
			continue
		}
		if len(cands) == 0 {
			t.Errorf("%q: expected candidate with slug %q, got none", tc.title, tc.wantSlug)
			continue
		}
		if cands[0].Slug != tc.wantSlug {
			t.Errorf("%q: want slug %q, got %q", tc.title, tc.wantSlug, cands[0].Slug)
		}
	}
	_ = ts // keep the first server alive (not actually used — test uses tsLocal)
}

