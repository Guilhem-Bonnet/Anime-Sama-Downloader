package app

import (
	"context"
	"testing"
	"time"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/ports"
)

func testCatalogueForSuggestions() []domain.AnimeSearchResult {
	return []domain.AnimeSearchResult{
		{
			ID:           "1",
			Title:        "Naruto",
			Year:         2002,
			Status:       "completed",
			EpisodeCount: 220,
			Genres:       []string{"Action", "Adventure", "Shonen"},
		},
		{
			ID:           "2",
			Title:        "Naruto Shippuden",
			Year:         2007,
			Status:       "completed",
			EpisodeCount: 500,
			Genres:       []string{"Action", "Adventure", "Shonen"},
		},
		{
			ID:           "3",
			Title:        "One Piece",
			Year:         1999,
			Status:       "ongoing",
			EpisodeCount: 1050,
			Genres:       []string{"Action", "Adventure", "Comedy"},
		},
		{
			ID:           "4",
			Title:        "Dragon Ball Z",
			Year:         1989,
			Status:       "completed",
			EpisodeCount: 291,
			Genres:       []string{"Action", "Adventure", "Shonen"},
		},
		{
			ID:           "5",
			Title:        "Attack on Titan",
			Year:         2013,
			Status:       "completed",
			EpisodeCount: 94,
			Genres:       []string{"Action", "Drama", "Dark Fantasy"},
		},
	}
}

func TestSuggestionsService_GetSuggestions_Empty(t *testing.T) {
	service := NewSuggestionsService(testCatalogueForSuggestions())
	ctx := context.Background()

	suggestions, err := service.GetSuggestions(ctx, "", 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Empty query should return trending (which is empty initially)
	if len(suggestions) != 0 {
		t.Errorf("expected 0 suggestions for empty query, got %d", len(suggestions))
	}
}

func TestSuggestionsService_GetSuggestions_PopularMatches(t *testing.T) {
	service := NewSuggestionsService(testCatalogueForSuggestions())
	ctx := context.Background()

	suggestions, err := service.GetSuggestions(ctx, "naruto", 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(suggestions) == 0 {
		t.Fatal("expected suggestions for 'naruto'")
	}

	// Should find both Naruto and Naruto Shippuden
	hasNaruto := false
	hasNarutoShippuden := false
	for _, s := range suggestions {
		if s.Query == "Naruto" {
			hasNaruto = true
		}
		if s.Query == "Naruto Shippuden" {
			hasNarutoShippuden = true
		}
	}

	if !hasNaruto {
		t.Error("expected 'Naruto' in suggestions")
	}
	if !hasNarutoShippuden {
		t.Error("expected 'Naruto Shippuden' in suggestions")
	}
}

func TestSuggestionsService_GetSuggestions_WithRecentSearches(t *testing.T) {
	service := NewSuggestionsService(testCatalogueForSuggestions())
	ctx := context.Background()

	// Track some searches
	service.TrackSearch("Naruto")
	service.TrackSearch("One Piece")
	service.TrackSearch("Naruto")

	suggestions, err := service.GetSuggestions(ctx, "naruto", 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// First suggestion should be from recent history (higher score)
	if len(suggestions) > 0 && suggestions[0].Category != ports.SuggestionCategoryRecent {
		t.Errorf("expected first suggestion to be recent, got %s", suggestions[0].Category)
	}
}

func TestSuggestionsService_TrackSearch(t *testing.T) {
	service := NewSuggestionsService(testCatalogueForSuggestions())

	// Track searches
	service.TrackSearch("Naruto")
	service.TrackSearch("One Piece")

	// Get trending
	trending, err := service.GetTrendingSuggestions(context.Background(), 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(trending) == 0 {
		t.Fatal("expected trending suggestions after tracking searches")
	}

	if len(trending) != 2 {
		t.Errorf("expected 2 trending, got %d", len(trending))
	}
}

func TestSuggestionsService_TrackSearch_Deduplication(t *testing.T) {
	service := NewSuggestionsService(testCatalogueForSuggestions())

	// Track same search multiple times
	service.TrackSearch("Naruto")
	service.TrackSearch("Naruto")
	service.TrackSearch("Naruto")

	trending, err := service.GetTrendingSuggestions(context.Background(), 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should have 1 entry but with score of 3
	if len(trending) != 1 {
		t.Errorf("expected 1 trending entry, got %d", len(trending))
	}

	if trending[0].Score != 3 {
		t.Errorf("expected score 3 (3 searches), got %v", trending[0].Score)
	}
}

func TestSuggestionsService_GetGenreMatches(t *testing.T) {
	service := NewSuggestionsService(testCatalogueForSuggestions())
	ctx := context.Background()

	suggestions, err := service.GetSuggestions(ctx, "action", 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should find "Action" genre suggestion
	hasActionGenre := false
	for _, s := range suggestions {
		if s.Query == "genre:Action" {
			hasActionGenre = true
			if s.Category != ports.SuggestionCategoryGenre {
				t.Errorf("expected genre category, got %s", s.Category)
			}
		}
	}

	if !hasActionGenre {
		t.Error("expected 'genre:Action' in suggestions")
	}
}

func TestSuggestionsService_CaseInsensitive(t *testing.T) {
	service := NewSuggestionsService(testCatalogueForSuggestions())
	ctx := context.Background()

	tests := []string{"NARUTO", "naruto", "Naruto", "NaRuTo"}

	for _, query := range tests {
		suggestions, err := service.GetSuggestions(ctx, query, 10)
		if err != nil {
			t.Fatalf("unexpected error for query %q: %v", query, err)
		}

		if len(suggestions) == 0 {
			t.Errorf("expected suggestions for query %q", query)
		}
	}
}

func TestSuggestionsService_Limit(t *testing.T) {
	service := NewSuggestionsService(testCatalogueForSuggestions())
	ctx := context.Background()

	// Track many searches
	for i := 0; i < 20; i++ {
		service.TrackSearch("test" + string(rune(i)))
	}

	suggestions, err := service.GetSuggestions(ctx, "test", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(suggestions) > 5 {
		t.Errorf("expected max 5 suggestions, got %d", len(suggestions))
	}
}

func TestSuggestionsService_TrendWindow(t *testing.T) {
	service := NewSuggestionsService(testCatalogueForSuggestions())
	ctx := context.Background()

	// Track a search
	service.TrackSearch("Old Search")

	// Manually add an old search to history (before trend window)
	service.mu.Lock()
	service.searchHistory[0].timestamp = time.Now().Add(-48 * time.Hour)
	service.mu.Unlock()

	trending, err := service.GetTrendingSuggestions(ctx, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Old search should not be in trending (outside 24h window)
	for _, s := range trending {
		if s.Query == "Old Search" {
			t.Error("expected old search to be outside trend window")
		}
	}
}
