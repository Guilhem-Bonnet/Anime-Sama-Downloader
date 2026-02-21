package app

import (
	"context"
	"strings"
)

type AnimeSamaCandidatesResolver interface {
	ResolveCandidates(ctx context.Context, titles []string, maxCandidates int) ([]AnimeSamaCandidate, error)
}

type AniListTitleSearcher interface {
	SearchAnimeTitles(ctx context.Context, query string, limit int) ([]string, error)
}

type AnimeSamaHybridResolver struct {
	Catalogue AnimeSamaCandidatesResolver
	AniList   AniListTitleSearcher

	// How many AniList results we fetch when falling back.
	AniListLimit int
}

func NewAnimeSamaHybridResolver(catalogue AnimeSamaCandidatesResolver, anilist AniListTitleSearcher) *AnimeSamaHybridResolver {
	return &AnimeSamaHybridResolver{Catalogue: catalogue, AniList: anilist, AniListLimit: 5}
}

func (r *AnimeSamaHybridResolver) ResolveCandidates(ctx context.Context, titles []string, maxCandidates int) ([]AnimeSamaCandidate, error) {
	if r == nil || r.Catalogue == nil {
		return nil, nil
	}

	// First try: direct slug probing from given titles.
	cands, err := r.Catalogue.ResolveCandidates(ctx, titles, maxCandidates)
	if err != nil {
		return nil, err
	}
	if len(cands) > 0 {
		return cands, nil
	}

	// Fallback: use AniList search to enrich title variants (romaji/english/native/synonyms).
	if r.AniList == nil || len(titles) == 0 {
		return cands, nil
	}
	q := strings.TrimSpace(titles[0])
	if q == "" {
		return cands, nil
	}

	limit := r.AniListLimit
	if limit <= 0 {
		limit = 5
	}
	alts, err := r.AniList.SearchAnimeTitles(ctx, q, limit)
	if err != nil {
		// Best-effort fallback: if AniList is down, keep original result.
		return cands, nil
	}
	if len(alts) == 0 {
		return cands, nil
	}

	enriched := mergeTitles(titles, alts)
	return r.Catalogue.ResolveCandidates(ctx, enriched, maxCandidates)
}

func mergeTitles(base []string, extra []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(base)+len(extra))
	add := func(v string) {
		v = strings.TrimSpace(v)
		if v == "" {
			return
		}
		k := strings.ToLower(v)
		if _, ok := seen[k]; ok {
			return
		}
		seen[k] = struct{}{}
		out = append(out, v)
	}
	for _, v := range base {
		add(v)
	}
	for _, v := range extra {
		add(v)
		if len(out) >= 12 {
			break
		}
	}
	return out
}
