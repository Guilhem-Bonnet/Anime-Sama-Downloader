package ports

import (
	"context"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
)

// IEpisodeResolver defines the interface for resolving episodes from anime sources.
// This abstraction allows dependency injection and enables mocking for unit tests.
type IEpisodeResolver interface {
	// Resolve fetches and resolves episodes for a given anime base URL.
	// It selects the best available player and returns the episode URLs.
	//
	// Parameters:
	//   ctx: context for cancellation and timeouts
	//   baseURL: the anime-sama base URL (e.g., "https://anime-sama.si/anime/...")
	//   preferredPlayer: player name preference (e.g., "betaseries", "uqload")
	//
	// Returns:
	//   domain.ResolvedEpisodes: selected player, episode URLs, max episode count
	//   error: if fetching or parsing fails
	Resolve(ctx context.Context, baseURL, preferredPlayer string) (domain.ResolvedEpisodes, error)
}
