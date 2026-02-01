package app

import (
	"testing"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
)

func TestRecommendationsService_GetSimilarAnime(t *testing.T) {
	catalogue := testCatalogueForSuggestions()
	service := NewRecommendationsService(catalogue)

	// Get recommendations for Naruto (ID: "1")
	recommendations, err := service.GetSimilarAnime("1", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(recommendations) == 0 {
		t.Fatal("expected recommendations for Naruto")
	}

	// Should find Naruto Shippuden (same genres)
	hasNarutoShippuden := false
	for _, rec := range recommendations {
		if rec.Title == "Naruto Shippuden" {
			hasNarutoShippuden = true
			if rec.Score <= 0 {
				t.Error("expected positive score for Naruto Shippuden")
			}
			break
		}
	}

	if !hasNarutoShippuden {
		t.Error("expected Naruto Shippuden in recommendations")
	}
}

func TestRecommendationsService_GetSimilarAnime_NonExistent(t *testing.T) {
	catalogue := testCatalogueForSuggestions()
	service := NewRecommendationsService(catalogue)

	recommendations, err := service.GetSimilarAnime("999", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(recommendations) != 0 {
		t.Errorf("expected empty recommendations for non-existent anime, got %d", len(recommendations))
	}
}

func TestRecommendationsService_GetRecommendationsByQuery(t *testing.T) {
	catalogue := testCatalogueForSuggestions()
	service := NewRecommendationsService(catalogue)

	// Search for "Naruto"
	recommendations, err := service.GetRecommendationsByQuery("Naruto", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(recommendations) == 0 {
		t.Fatal("expected recommendations for Naruto query")
	}

	// Should find other Action/Adventure/Shonen anime
	hasActionAdventure := false
	for _, rec := range recommendations {
		if rec.Title == "One Piece" || rec.Title == "Dragon Ball Z" {
			hasActionAdventure = true
			break
		}
	}

	if !hasActionAdventure {
		t.Error("expected other Action/Adventure anime in recommendations")
	}
}

func TestRecommendationsService_GetRecommendationsByQuery_Empty(t *testing.T) {
	catalogue := testCatalogueForSuggestions()
	service := NewRecommendationsService(catalogue)

	recommendations, err := service.GetRecommendationsByQuery("", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(recommendations) != 0 {
		t.Errorf("expected empty recommendations for empty query, got %d", len(recommendations))
	}
}

func TestRecommendationsService_GetRecommendationsForGenres(t *testing.T) {
	catalogue := testCatalogueForSuggestions()
	service := NewRecommendationsService(catalogue)

	// Get recommendations for Action genre
	recommendations, err := service.GetRecommendationsForGenres([]string{"Action"}, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(recommendations) == 0 {
		t.Fatal("expected recommendations for Action genre")
	}

	// All recommendations should have Action genre
	for _, rec := range recommendations {
		// Find the anime to check genre
		found := false
		for _, anime := range catalogue {
			if anime.Title == rec.Title {
				hasAction := false
				for _, g := range anime.Genres {
					if g == "Action" {
						hasAction = true
						break
					}
				}
				if !hasAction {
					t.Errorf("expected %q to have Action genre", rec.Title)
				}
				found = true
				break
			}
		}
		if !found {
			t.Errorf("anime %q not found in catalogue", rec.Title)
		}
	}
}

func TestRecommendationsService_GetRecommendationsForGenres_Multiple(t *testing.T) {
	catalogue := testCatalogueForSuggestions()
	service := NewRecommendationsService(catalogue)

	// Get recommendations for Action + Adventure
	recommendations, err := service.GetRecommendationsForGenres([]string{"Action", "Adventure"}, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(recommendations) == 0 {
		t.Fatal("expected recommendations for Action + Adventure genres")
	}

	// Higher scored anime should have both genres
	if len(recommendations) > 0 {
		topAnime := recommendations[0]
		found := false
		for _, anime := range catalogue {
			if anime.Title == topAnime.Title {
				hasAction := false
				hasAdventure := false
				for _, g := range anime.Genres {
					if g == "Action" {
						hasAction = true
					}
					if g == "Adventure" {
						hasAdventure = true
					}
				}
				if !hasAction || !hasAdventure {
					t.Errorf("expected top recommendation to have both Action and Adventure genres")
				}
				found = true
				break
			}
		}
		if !found {
			t.Errorf("top anime %q not found in catalogue", topAnime.Title)
		}
	}
}

func TestRecommendationsService_ScoringLogic(t *testing.T) {
	catalogue := []domain.AnimeSearchResult{
		{
			ID:           "1",
			Title:        "Anime A",
			Year:         2020,
			Status:       "ongoing",
			EpisodeCount: 100,
			Genres:       []string{"Action", "Adventure"},
		},
		{
			ID:           "2",
			Title:        "Anime B",
			Year:         2020,
			Status:       "ongoing",
			EpisodeCount: 110,
			Genres:       []string{"Action", "Adventure"}, // All matches: highest score
		},
		{
			ID:           "3",
			Title:        "Anime C",
			Year:         2020,
			Status:       "completed",
			EpisodeCount: 100,
			Genres:       []string{"Action", "Comedy"}, // 1 genre match, different status
		},
		{
			ID:           "4",
			Title:        "Anime D",
			Year:         2018,
			Status:       "completed",
			EpisodeCount: 50,
			Genres:       []string{"Comedy"}, // No genre match
		},
	}

	service := NewRecommendationsService(catalogue)

	recommendations, err := service.GetSimilarAnime("1", 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Anime B should be first (all genres and status match)
	if len(recommendations) == 0 {
		t.Fatal("expected recommendations")
	}

	if recommendations[0].Title != "Anime B" {
		t.Errorf("expected Anime B to be first, got %s", recommendations[0].Title)
	}

	// Anime C should be second (1 genre match)
	if len(recommendations) < 2 {
		t.Fatal("expected at least 2 recommendations")
	}

	if recommendations[1].Title != "Anime C" {
		t.Errorf("expected Anime C to be second, got %s", recommendations[1].Title)
	}
}

func TestRecommendationsService_Limit(t *testing.T) {
	catalogue := testCatalogueForSuggestions()
	service := NewRecommendationsService(catalogue)

	recommendations, err := service.GetSimilarAnime("1", 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(recommendations) > 2 {
		t.Errorf("expected max 2 recommendations, got %d", len(recommendations))
	}
}
