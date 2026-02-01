package domain

// AnimeSearchResult represents a single anime search result.
type AnimeSearchResult struct {
	ID           string   // Unique identifier for the anime
	Title        string   // Anime title
	ThumbnailURL string   // URL to thumbnail image
	Year         int      // Year the anime was released
	Status       string   // "ongoing", "completed", "planning", etc.
	EpisodeCount int      // Total number of episodes
	Genres       []string // List of genres (e.g., ["Action", "Adventure"])
}

// SearchResultWithScore is used internally for ranking search results.
type SearchResultWithScore struct {
	Result AnimeSearchResult
	Score  float64
}
