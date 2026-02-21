package app

import (
	"strings"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/ports"
)

// RecommendationsServiceImpl implements the RecommendationsService
type RecommendationsServiceImpl struct {
	catalogue []domain.AnimeSearchResult
}

// NewRecommendationsService creates a new recommendations service
func NewRecommendationsService(catalogue []domain.AnimeSearchResult) *RecommendationsServiceImpl {
	return &RecommendationsServiceImpl{catalogue: catalogue}
}

// GetSimilarAnime returns anime similar to the given anime
func (s *RecommendationsServiceImpl) GetSimilarAnime(animeID string, limit int) ([]ports.RecommendationScore, error) {
	if limit <= 0 {
		limit = 10
	}

	// Find the source anime
	var sourceAnime *domain.AnimeSearchResult
	for i := range s.catalogue {
		if s.catalogue[i].ID == animeID {
			sourceAnime = &s.catalogue[i]
			break
		}
	}

	if sourceAnime == nil {
		return []ports.RecommendationScore{}, nil
	}

	// Score all other anime based on similarity
	type scoredAnime struct {
		anime  domain.AnimeSearchResult
		score  float64
		reason string
	}

	scored := []scoredAnime{}

	for _, anime := range s.catalogue {
		if anime.ID == animeID {
			continue // Skip the source anime
		}

		score := 0.0
		reason := ""

		// Genre matches (1 point per matching genre, max 3)
		genreMatches := 0
		for _, sourceGenre := range sourceAnime.Genres {
			for _, animeGenre := range anime.Genres {
				if strings.EqualFold(sourceGenre, animeGenre) {
					genreMatches++
					break
				}
			}
		}

		if genreMatches > 3 {
			genreMatches = 3
		}
		score += float64(genreMatches) * 2.0

		// Status match (1 point)
		statusMatch := false
		if strings.EqualFold(sourceAnime.Status, anime.Status) {
			score += 1.0
			statusMatch = true
		}

		// Year proximity (0.5 points if within 2 years)
		yearDiff := sourceAnime.Year - anime.Year
		if yearDiff < 0 {
			yearDiff = -yearDiff
		}
		if yearDiff <= 2 {
			score += 0.5
		}

		// Popularity proximity (0.5 points if episode count is similar)
		episodeDiff := sourceAnime.EpisodeCount - anime.EpisodeCount
		if episodeDiff < 0 {
			episodeDiff = -episodeDiff
		}
		// Normalize to 0-100 range (assume max ~500 episodes)
		if float64(episodeDiff) <= 100 {
			score += 0.5
		}

		if score > 0 {
			reasonParts := []string{}
			if genreMatches > 0 {
				reasonParts = append(reasonParts, strings.TrimSpace(strings.Join([]string{string('0' + rune(genreMatches)), "genre"}, " ")))
			}
			if statusMatch {
				reasonParts = append(reasonParts, "status")
			}

			reason = strings.Join(reasonParts, ", ")

			scored = append(scored, scoredAnime{
				anime:  anime,
				score:  score,
				reason: reason,
			})
		}
	}

	// Sort by score (descending)
	for i := 0; i < len(scored); i++ {
		for j := i + 1; j < len(scored); j++ {
			if scored[j].score > scored[i].score {
				scored[i], scored[j] = scored[j], scored[i]
			}
		}
	}

	// Limit results
	if len(scored) > limit {
		scored = scored[:limit]
	}

	// Convert to RecommendationScore
	results := make([]ports.RecommendationScore, len(scored))
	for i, s := range scored {
		results[i] = ports.RecommendationScore{
			AnimeID: s.anime.ID,
			Title:   s.anime.Title,
			Score:   s.score,
			Reason:  s.reason,
		}
	}

	return results, nil
}

// GetRecommendationsByQuery returns recommendations based on search query
func (s *RecommendationsServiceImpl) GetRecommendationsByQuery(query string, limit int) ([]ports.RecommendationScore, error) {
	if limit <= 0 {
		limit = 10
	}

	query = strings.TrimSpace(query)
	if query == "" {
		return []ports.RecommendationScore{}, nil
	}

	// Find anime matching the query
	normalizedQuery := normalizeForMatching(query)
	matchingAnimes := []domain.AnimeSearchResult{}

	for _, anime := range s.catalogue {
		if strings.Contains(normalizeForMatching(anime.Title), normalizedQuery) {
			matchingAnimes = append(matchingAnimes, anime)
		}
	}

	if len(matchingAnimes) == 0 {
		return []ports.RecommendationScore{}, nil
	}

	// Collect all genres and statuses from matching anime
	genres := make(map[string]bool)
	statuses := make(map[string]bool)

	for _, anime := range matchingAnimes {
		for _, genre := range anime.Genres {
			genres[genre] = true
		}
		statuses[anime.Status] = true
	}

	// Score anime based on genre and status matches
	type scoredAnime struct {
		anime domain.AnimeSearchResult
		score float64
	}

	scored := []scoredAnime{}

	for _, anime := range s.catalogue {
		// Don't recommend if already in query results
		isInQuery := false
		for _, q := range matchingAnimes {
			if q.ID == anime.ID {
				isInQuery = true
				break
			}
		}
		if isInQuery {
			continue
		}

		score := 0.0

		// Count matching genres
		for _, animeGenre := range anime.Genres {
			if genres[animeGenre] {
				score += 2.0
			}
		}

		// Status match
		if statuses[anime.Status] {
			score += 1.0
		}

		if score > 0 {
			scored = append(scored, scoredAnime{anime, score})
		}
	}

	// Sort by score
	for i := 0; i < len(scored); i++ {
		for j := i + 1; j < len(scored); j++ {
			if scored[j].score > scored[i].score {
				scored[i], scored[j] = scored[j], scored[i]
			}
		}
	}

	// Limit results
	if len(scored) > limit {
		scored = scored[:limit]
	}

	// Convert to RecommendationScore
	results := make([]ports.RecommendationScore, len(scored))
	for i, s := range scored {
		results[i] = ports.RecommendationScore{
			AnimeID: s.anime.ID,
			Title:   s.anime.Title,
			Score:   s.score,
		}
	}

	return results, nil
}

// GetRecommendationsForGenres returns anime matching multiple genres
func (s *RecommendationsServiceImpl) GetRecommendationsForGenres(genres []string, limit int) ([]ports.RecommendationScore, error) {
	if limit <= 0 {
		limit = 10
	}

	if len(genres) == 0 {
		return []ports.RecommendationScore{}, nil
	}

	type scoredAnime struct {
		anime domain.AnimeSearchResult
		score float64
	}

	scored := []scoredAnime{}

	// Score anime by genre matches
	for _, anime := range s.catalogue {
		score := 0.0

		for _, animeGenre := range anime.Genres {
			for _, targetGenre := range genres {
				if strings.EqualFold(animeGenre, targetGenre) {
					score += 1.0
					break
				}
			}
		}

		if score > 0 {
			scored = append(scored, scoredAnime{anime, score})
		}
	}

	// Sort by score
	for i := 0; i < len(scored); i++ {
		for j := i + 1; j < len(scored); j++ {
			if scored[j].score > scored[i].score {
				scored[i], scored[j] = scored[j], scored[i]
			}
		}
	}

	// Limit results
	if len(scored) > limit {
		scored = scored[:limit]
	}

	// Convert to RecommendationScore
	results := make([]ports.RecommendationScore, len(scored))
	for i, s := range scored {
		results[i] = ports.RecommendationScore{
			AnimeID: s.anime.ID,
			Title:   s.anime.Title,
			Score:   s.score,
		}
	}

	return results, nil
}
