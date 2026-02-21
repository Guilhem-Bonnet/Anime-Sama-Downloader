package ports

import (
	"context"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
)

// SearchFilters represents search filtering options.
type SearchFilters struct {
	Query   string   // Search query string
	Genres  []string // Filter by genres (empty = no filter)
	Status  string   // Filter by status: "ongoing", "completed", "planning", "" (all)
	YearMin int      // Minimum year (0 = no filter)
	YearMax int      // Maximum year (0 = no filter)
}

// AnimeSearch defines the interface for anime search functionality.
type AnimeSearch interface {
	// Search performs an anime search using the provided query string.
	// Returns a sorted list of results (most relevant first), max 50 results.
	// Empty query returns empty slice (no error).
	Search(ctx context.Context, query string) ([]domain.AnimeSearchResult, error)

	// SearchWithFilters performs a filtered anime search.
	// Applies genre, status, and year filters to the results.
	// Returns filtered and sorted results, max 50 results.
	SearchWithFilters(ctx context.Context, filters SearchFilters) ([]domain.AnimeSearchResult, error)
}
