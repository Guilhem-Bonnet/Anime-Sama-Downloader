package app

import (
	"context"
	"testing"
)

type stubAniList struct {
	entries []AniListWatchlistEntry
	err     error
}

func (s stubAniList) Watchlist(ctx context.Context, statuses []string) ([]AniListWatchlistEntry, error) {
	return s.entries, s.err
}

type stubResolver struct {
	cands []AnimeSamaCandidate
	err   error
}

func (s stubResolver) ResolveCandidates(ctx context.Context, titles []string, maxCandidates int) ([]AnimeSamaCandidate, error) {
	return s.cands, s.err
}

type stubSubs struct {
	created []SubscriptionDTO
	list    []SubscriptionDTO
	err     error
}

func (s *stubSubs) List(ctx context.Context, limit int) ([]SubscriptionDTO, error) {
	return s.list, nil
}

func (s *stubSubs) Create(ctx context.Context, baseURL, label, player string) (SubscriptionDTO, error) {
	if s.err != nil {
		return SubscriptionDTO{}, s.err
	}
	dto := SubscriptionDTO{ID: "sub1", BaseURL: baseURL, Label: label, Player: player}
	s.created = append(s.created, dto)
	return dto, nil
}

func TestAniListImport_AutoConfirm_SkipsAmbiguousAndCreatesConfident(t *testing.T) {
	ctx := context.Background()

	entry := AniListWatchlistEntry{}
	entry.Status = "CURRENT"
	entry.Progress = 0
	entry.Media.ID = 123
	entry.Media.Synonyms = []string{}
	entry.Media.Title.Romaji = "Test Anime"
	entry.Media.Title.English = "Test Anime"

	anilist := stubAniList{entries: []AniListWatchlistEntry{entry}}

	// Two high-score candidates => ambiguous => skipped.
	resolverAmbig := stubResolver{cands: []AnimeSamaCandidate{
		{CatalogueURL: "https://anime-sama.si/catalogue/test-anime/", Slug: "test-anime", MatchedTitle: "Test Anime", Score: 0.99},
		{CatalogueURL: "https://anime-sama.si/catalogue/test-anime-2/", Slug: "test-anime-2", MatchedTitle: "Test Anime", Score: 0.98},
	}}

	subs := &stubSubs{}

	svcAmbig := &AniListImportService{anilist: anilist, resolver: resolverAmbig, subs: subs}
	resAmbig, err := svcAmbig.AutoConfirm(ctx, AniListImportAutoRequest{Season: 1, Lang: "vostfr", MaxCandidates: 3, MinScore: 0.95})
	if err != nil {
		t.Fatalf("autoConfirm(ambig): %v", err)
	}
	if len(resAmbig.Created) != 0 {
		t.Fatalf("expected 0 created (ambig), got %d", len(resAmbig.Created))
	}
	if len(resAmbig.Skipped) != 1 {
		t.Fatalf("expected 1 skipped (ambig), got %d", len(resAmbig.Skipped))
	}

	// Now test the happy path (single confident candidate) with subs stub.
	resolverOK := stubResolver{cands: []AnimeSamaCandidate{{
		CatalogueURL: "https://anime-sama.si/catalogue/test-anime/", Slug: "test-anime", MatchedTitle: "Test Anime", Score: 1.0,
	}}}

	svcOK := &AniListImportService{anilist: anilist, resolver: resolverOK, subs: subs}
	res, err := svcOK.AutoConfirm(ctx, AniListImportAutoRequest{Season: 1, Lang: "vostfr", MaxCandidates: 3, MinScore: 0.95})
	if err != nil {
		t.Fatalf("autoConfirm: %v", err)
	}
	if len(res.Created) != 1 {
		t.Fatalf("expected 1 created, got %d", len(res.Created))
	}
	if len(res.Skipped) != 0 {
		t.Fatalf("expected 0 skipped, got %d", len(res.Skipped))
	}
}
