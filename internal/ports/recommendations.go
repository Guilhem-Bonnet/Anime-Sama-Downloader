package ports

// RecommendationScore represents how well an anime matches recommendation criteria
type RecommendationScore struct {
	AnimeID string  `json:"anime_id"`
	Title   string  `json:"title"`
	Score   float64 `json:"score"`
	Reason  string  `json:"reason"` // e.g., "2 genre matches", "Same status"
}

// RecommendationsService defines the interface for anime recommendations
type RecommendationsService interface {
	// GetSimilarAnime returns anime similar to the given anime by genres and status
	GetSimilarAnime(animeID string, limit int) ([]RecommendationScore, error)

	// GetRecommendationsByQuery returns anime recommended based on search query
	// Uses genres and status from matching anime to find similar ones
	GetRecommendationsByQuery(query string, limit int) ([]RecommendationScore, error)

	// GetRecommendationsForGenres returns anime matching multiple genres
	GetRecommendationsForGenres(genres []string, limit int) ([]RecommendationScore, error)
}
