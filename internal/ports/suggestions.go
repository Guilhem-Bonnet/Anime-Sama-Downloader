package ports

import "context"

// SuggestionCategory represents the category of a suggestion
type SuggestionCategory string

const (
	SuggestionCategoryRecent   SuggestionCategory = "recent"
	SuggestionCategoryPopular  SuggestionCategory = "popular"
	SuggestionCategoryTrending SuggestionCategory = "trending"
	SuggestionCategoryGenre    SuggestionCategory = "genre"
)

// Suggestion represents a single search suggestion
type Suggestion struct {
	Query    string                 `json:"query"`
	Category SuggestionCategory     `json:"category"`
	Score    float64                `json:"score"`
	Metadata map[string]interface{} `json:"metadata,omitempty"` // e.g., anime count, genre info
}

// SuggestionsService defines the interface for search suggestions
type SuggestionsService interface {
	// GetSuggestions returns search suggestions based on query prefix
	// Returns categorized suggestions (recent, popular, trending, genre-based)
	GetSuggestions(ctx context.Context, query string, limit int) ([]Suggestion, error)

	// GetTrendingSuggestions returns currently trending search terms
	GetTrendingSuggestions(ctx context.Context, limit int) ([]Suggestion, error)

	// TrackSearch records a search query for trending/popularity calculation
	TrackSearch(query string) error
}
