package app

import (
	"context"
	"strings"
	"testing"
)

type fakeCatalogueResolver struct {
	calls [][]string
	// map key is first title (lowercase) -> candidates
	results map[string][]AnimeSamaCandidate
}

func (f *fakeCatalogueResolver) ResolveCandidates(ctx context.Context, titles []string, maxCandidates int) ([]AnimeSamaCandidate, error) {
	f.calls = append(f.calls, append([]string(nil), titles...))
	if len(titles) == 0 {
		return nil, nil
	}
	for _, t := range titles {
		key := strings.ToLower(strings.TrimSpace(t))
		if got, ok := f.results[key]; ok {
			return got, nil
		}
	}
	return nil, nil
}

type fakeAniListSearcher struct {
	alts []string
	err  error
}

func (f *fakeAniListSearcher) SearchAnimeTitles(ctx context.Context, query string, limit int) ([]string, error) {
	if f.err != nil {
		return nil, f.err
	}
	return append([]string(nil), f.alts...), nil
}

func TestAnimeSamaHybridResolver_FallbackToAniListWhenNoCandidates(t *testing.T) {
	ctx := context.Background()

	cat := &fakeCatalogueResolver{results: map[string][]AnimeSamaCandidate{
		"jigokuraku": {{CatalogueURL: "https://anime-sama.si/catalogue/jigokuraku/", Slug: "jigokuraku", MatchedTitle: "Jigokuraku", Score: 1.0}},
	}}
	ali := &fakeAniListSearcher{alts: []string{"Jigokuraku"}}

	r := NewAnimeSamaHybridResolver(cat, ali)
	cands, err := r.ResolveCandidates(ctx, []string{"Hell's Paradise"}, 3)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(cands) != 1 || cands[0].Slug != "jigokuraku" {
		t.Fatalf("unexpected candidates: %+v", cands)
	}
	if len(cat.calls) < 2 {
		t.Fatalf("expected 2 catalogue calls (direct + fallback), got %d", len(cat.calls))
	}
}

func TestAnimeSamaHybridResolver_NoFallbackWhenDirectCandidates(t *testing.T) {
	ctx := context.Background()

	cat := &fakeCatalogueResolver{results: map[string][]AnimeSamaCandidate{
		"hell's paradise": {{CatalogueURL: "https://anime-sama.si/catalogue/hell-s-paradise/", Slug: "hell-s-paradise", MatchedTitle: "Hell's Paradise", Score: 1.0}},
	}}
	ali := &fakeAniListSearcher{alts: []string{"Jigokuraku"}}

	r := NewAnimeSamaHybridResolver(cat, ali)
	cands, err := r.ResolveCandidates(ctx, []string{"Hell's Paradise"}, 3)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(cands) != 1 || cands[0].Slug != "hell-s-paradise" {
		t.Fatalf("unexpected candidates: %+v", cands)
	}
	if len(cat.calls) != 1 {
		t.Fatalf("expected 1 catalogue call, got %d", len(cat.calls))
	}
}
