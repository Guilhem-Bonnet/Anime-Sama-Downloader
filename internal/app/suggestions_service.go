package app

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/ports"
)

// SuggestionsServiceImpl implements the SuggestionsService
type SuggestionsServiceImpl struct {
	catalogue      []domain.AnimeSearchResult
	mu             sync.RWMutex
	searchHistory  []searchRecord
	maxHistorySize int
	trendWindow    time.Duration
}

type searchRecord struct {
	query     string
	timestamp time.Time
}

// NewSuggestionsService creates a new suggestions service
func NewSuggestionsService(catalogue []domain.AnimeSearchResult) *SuggestionsServiceImpl {
	return &SuggestionsServiceImpl{
		catalogue:      catalogue,
		searchHistory:  []searchRecord{},
		maxHistorySize: 1000,
		trendWindow:    24 * time.Hour,
	}
}

// GetSuggestions returns categorized suggestions based on query
func (s *SuggestionsServiceImpl) GetSuggestions(ctx context.Context, query string, limit int) ([]ports.Suggestion, error) {
	if limit <= 0 {
		limit = 10
	}

	query = strings.TrimSpace(query)
	if query == "" {
		// Return trending suggestions if no query
		return s.getTrendingSuggestions(limit), nil
	}

	suggestions := []ports.Suggestion{}

	// 1. Recent searches matching query
	recentMatches := s.getRecentMatches(query, limit/3)
	suggestions = append(suggestions, recentMatches...)

	// 2. Popular anime matching query
	if len(suggestions) < limit {
		popularMatches := s.getPopularMatches(query, limit-len(suggestions))
		suggestions = append(suggestions, popularMatches...)
	}

	// 3. Genre-based suggestions if query matches a genre
	if len(suggestions) < limit {
		genreMatches := s.getGenreMatches(query, limit-len(suggestions))
		suggestions = append(suggestions, genreMatches...)
	}

	// Limit results
	if len(suggestions) > limit {
		suggestions = suggestions[:limit]
	}

	return suggestions, nil
}

// GetTrendingSuggestions returns currently trending searches
func (s *SuggestionsServiceImpl) GetTrendingSuggestions(ctx context.Context, limit int) ([]ports.Suggestion, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.getTrendingSuggestions(limit), nil
}

// TrackSearch records a search for trending calculation
func (s *SuggestionsServiceImpl) TrackSearch(query string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	query = strings.TrimSpace(query)
	if query == "" {
		return nil
	}

	s.searchHistory = append(s.searchHistory, searchRecord{
		query:     query,
		timestamp: time.Now(),
	})

	// Keep history size under control
	if len(s.searchHistory) > s.maxHistorySize {
		s.searchHistory = s.searchHistory[1:]
	}

	return nil
}

// getRecentMatches returns recent searches matching the query
func (s *SuggestionsServiceImpl) getRecentMatches(query string, limit int) []ports.Suggestion {
	s.mu.RLock()
	defer s.mu.RUnlock()

	suggestions := []ports.Suggestion{}
	seen := make(map[string]bool)
	normalizedQuery := normalizeForMatching(query)

	// Iterate in reverse (newest first)
	for i := len(s.searchHistory) - 1; i >= 0 && len(suggestions) < limit; i-- {
		record := s.searchHistory[i]
		if seen[record.query] {
			continue
		}

		if strings.Contains(normalizeForMatching(record.query), normalizedQuery) {
			seen[record.query] = true
			suggestions = append(suggestions, ports.Suggestion{
				Query:    record.query,
				Category: ports.SuggestionCategoryRecent,
				Score:    calculateRecencyScore(record.timestamp),
			})
		}
	}

	return suggestions
}

// getPopularMatches returns anime matching the query, sorted by episode count (popularity proxy)
func (s *SuggestionsServiceImpl) getPopularMatches(query string, limit int) []ports.Suggestion {
	s.mu.RLock()
	defer s.mu.RUnlock()

	suggestions := []ports.Suggestion{}
	normalizedQuery := normalizeForMatching(query)

	// Find anime matching query
	for _, anime := range s.catalogue {
		if len(suggestions) >= limit {
			break
		}

		if strings.Contains(normalizeForMatching(anime.Title), normalizedQuery) {
			suggestions = append(suggestions, ports.Suggestion{
				Query:    anime.Title,
				Category: ports.SuggestionCategoryPopular,
				Score:    calculatePopularityScore(anime.EpisodeCount),
				Metadata: map[string]interface{}{
					"anime_id": anime.ID,
					"episodes": anime.EpisodeCount,
					"year":     anime.Year,
				},
			})
		}
	}

	return suggestions
}

// getGenreMatches returns suggestions for anime by matching genres
func (s *SuggestionsServiceImpl) getGenreMatches(query string, limit int) []ports.Suggestion {
	s.mu.RLock()
	defer s.mu.RUnlock()

	suggestions := []ports.Suggestion{}
	normalizedQuery := normalizeForMatching(query)
	genreCounts := make(map[string]int)

	// Find genre matches
	for _, anime := range s.catalogue {
		for _, genre := range anime.Genres {
			if strings.Contains(normalizeForMatching(genre), normalizedQuery) {
				genreCounts[genre]++
			}
		}
	}

	// Create suggestions from genres
	for genre, count := range genreCounts {
		if len(suggestions) >= limit {
			break
		}
		suggestions = append(suggestions, ports.Suggestion{
			Query:    "genre:" + genre,
			Category: ports.SuggestionCategoryGenre,
			Score:    float64(count) / 100.0, // Normalize count as score
			Metadata: map[string]interface{}{
				"anime_count": count,
			},
		})
	}

	return suggestions
}

// getTrendingSuggestions returns trending searches from recent history
func (s *SuggestionsServiceImpl) getTrendingSuggestions(limit int) []ports.Suggestion {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Count searches in trend window
	trendingCounts := make(map[string]int)
	now := time.Now()

	for _, record := range s.searchHistory {
		if now.Sub(record.timestamp) <= s.trendWindow {
			trendingCounts[record.query]++
		}
	}

	// Sort by count and create suggestions
	suggestions := []ports.Suggestion{}
	for query, count := range trendingCounts {
		suggestions = append(suggestions, ports.Suggestion{
			Query:    query,
			Category: ports.SuggestionCategoryTrending,
			Score:    float64(count),
		})
	}

	// Sort by score (simple bubble sort for small dataset)
	for i := 0; i < len(suggestions); i++ {
		for j := i + 1; j < len(suggestions); j++ {
			if suggestions[j].Score > suggestions[i].Score {
				suggestions[i], suggestions[j] = suggestions[j], suggestions[i]
			}
		}
	}

	if len(suggestions) > limit {
		suggestions = suggestions[:limit]
	}

	return suggestions
}

// Helper functions

func normalizeForMatching(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func calculateRecencyScore(timestamp time.Time) float64 {
	// More recent = higher score
	// Decay over 7 days
	daysAgo := time.Since(timestamp).Hours() / 24
	return (7 - daysAgo) / 7 * 100
}

func calculatePopularityScore(episodeCount int) float64 {
	// More episodes = more popular (proxy)
	// Normalize to 0-100 range (assuming max ~500 episodes)
	return float64(episodeCount) / 500 * 100
}
