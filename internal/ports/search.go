package ports

import (
	"context"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
)

// AnimeSearch defines the interface for anime search functionality.
type AnimeSearch interface {
	// Search performs an anime search using the provided query string.
	// Returns a sorted list of results (most relevant first), max 50 results.
	// Empty query returns empty slice (no error).
	Search(ctx context.Context, query string) ([]domain.AnimeSearchResult, error)
}
