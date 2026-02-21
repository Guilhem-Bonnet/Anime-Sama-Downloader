package app

import (
	"context"
	"fmt"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
)

// EpisodeResolver handles fetching and parsing anime episodes from anime-sama.
// It extracts the episode resolution logic from SubscriptionService to follow
// Single Responsibility Principle.
type EpisodeResolver struct {
	// Optionnel: injecter HTTP client si besoin custom (ex: timeouts, rate limiting)
}

// NewEpisodeResolver creates a new episode resolver.
func NewEpisodeResolver() *EpisodeResolver {
	return &EpisodeResolver{}
}

// Resolve fetches episodes from anime-sama baseURL and resolves player selection.
//
// Flow:
// 1. Fetch episodes.js from baseURL
// 2. Parse players map
// 3. Select best player (fallback if preferred unavailable)
// 4. Return URLs + max episode count
func (r *EpisodeResolver) Resolve(ctx context.Context, baseURL, preferredPlayer string) (domain.ResolvedEpisodes, error) {
	// Step 1: Fetch episodes.js
	jsText, err := FetchEpisodesJS(ctx, baseURL)
	if err != nil {
		return domain.ResolvedEpisodes{}, fmt.Errorf("fetch episodes.js: %w", err)
	}

	// Step 2: Parse episodes
	eps, err := ParseEpisodesJS(jsText)
	if err != nil {
		return domain.ResolvedEpisodes{}, fmt.Errorf("parse episodes.js: %w", err)
	}

	// Step 3: Select player using existing helper
	selected, urls := selectPlayer(preferredPlayer, eps.Players)

	// Step 4: Calculate max episode
	maxAvail := MaxAvailableEpisode(urls)
	if maxAvail < 0 {
		maxAvail = 0
	}

	return domain.ResolvedEpisodes{
		SelectedPlayer: selected,
		URLs:           urls,
		MaxEpisode:     maxAvail,
	}, nil
}
