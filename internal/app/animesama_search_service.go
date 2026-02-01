package app

import (
	"context"
	"sort"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
)

// AnimeSamaSearchService implements anime search functionality
type AnimeSamaSearchService struct {
	catalogue []domain.AnimeSearchResult
}

// NewAnimeSamaSearchService creates a new search service with the given anime catalogue
func NewAnimeSamaSearchService(catalogue []domain.AnimeSearchResult) *AnimeSamaSearchService {
	return &AnimeSamaSearchService{
		catalogue: catalogue,
	}
}

// Search performs an anime search using the provided query string.
// Results are ranked by relevance (exact match first, then partial matches).
// Returns max 50 results. Empty query returns empty slice.
func (s *AnimeSamaSearchService) Search(ctx context.Context, query string) ([]domain.AnimeSearchResult, error) {
	// Handle context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Normalize query: trim whitespace, lowercase, unicode normalization
	normalizedQuery := s.normalizeQuery(query)
	if normalizedQuery == "" {
		return []domain.AnimeSearchResult{}, nil
	}

	// Score all results
	var scored []domain.SearchResultWithScore
	for _, anime := range s.catalogue {
		score := s.calculateScore(normalizedQuery, anime)
		if score > 0 {
			scored = append(scored, domain.SearchResultWithScore{
				Result: anime,
				Score:  score,
			})
		}
	}

	// Sort by score (descending)
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].Score > scored[j].Score
	})

	// Limit to 50 results
	limit := 50
	if len(scored) < limit {
		limit = len(scored)
	}

	results := make([]domain.AnimeSearchResult, limit)
	for i := 0; i < limit; i++ {
		results[i] = scored[i].Result
	}

	return results, nil
}

// calculateScore assigns a relevance score to an anime based on the query.
// Exact match = 1000+, partial match varies based on position.
func (s *AnimeSamaSearchService) calculateScore(normalizedQuery string, anime domain.AnimeSearchResult) float64 {
	normalizedTitle := s.normalizeQuery(anime.Title)

	// Exact match (full title match)
	if normalizedTitle == normalizedQuery {
		return 1000.0 + 1.0/(float64(len(normalizedTitle))+1)
	}

	// Partial match (substring)
	if idx := strings.Contains(normalizedTitle, normalizedQuery); idx {
		pos := strings.Index(normalizedTitle, normalizedQuery)
		// Early matches score higher: position 0 = 100 points, decreases with position
		return 100.0 - float64(pos)
	}

	return 0
}

// normalizeQuery normalizes the query string for comparison.
// Applies: lowercase, trim whitespace, unicode normalization (NFD), accent removal
func (s *AnimeSamaSearchService) normalizeQuery(query string) string {
	// Trim whitespace
	query = strings.TrimSpace(query)

	// Lowercase
	query = strings.ToLower(query)

	// Unicode normalization: NFD + remove accents
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)))
	normalized, _, _ := transform.String(t, query)

	return normalized
}
