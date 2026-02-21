package app

import (
	"context"
	"testing"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/ports"
)

// testCatalogue returns a sample anime catalogue for testing
func testCatalogue() []domain.AnimeSearchResult {
	return []domain.AnimeSearchResult{
		{
			ID:           "1",
			Title:        "Naruto",
			ThumbnailURL: "https://example.com/naruto.jpg",
			Year:         2002,
			Status:       "completed",
			EpisodeCount: 220,
			Genres:       []string{"Action", "Adventure", "Shonen"},
		},
		{
			ID:           "2",
			Title:        "Naruto Shippuden",
			ThumbnailURL: "https://example.com/naruto-shippuden.jpg",
			Year:         2007,
			Status:       "completed",
			EpisodeCount: 500,
			Genres:       []string{"Action", "Adventure", "Shonen"},
		},
		{
			ID:           "3",
			Title:        "One Piece",
			ThumbnailURL: "https://example.com/one-piece.jpg",
			Year:         1999,
			Status:       "ongoing",
			EpisodeCount: 1050,
			Genres:       []string{"Action", "Adventure", "Comedy"},
		},
		{
			ID:           "4",
			Title:        "Demon Slayer",
			ThumbnailURL: "https://example.com/demon-slayer.jpg",
			Year:         2019,
			Status:       "ongoing",
			EpisodeCount: 50,
			Genres:       []string{"Action", "Supernatural", "Shonen"},
		},
		{
			ID:           "5",
			Title:        "Death Note",
			ThumbnailURL: "https://example.com/death-note.jpg",
			Year:         2006,
			Status:       "completed",
			EpisodeCount: 37,
			Genres:       []string{"Mystery", "Thriller", "Psychological"},
		},
		{
			ID:           "6",
			Title:        "Attack on Titan",
			ThumbnailURL: "https://example.com/aot.jpg",
			Year:         2013,
			Status:       "completed",
			EpisodeCount: 94,
			Genres:       []string{"Action", "Drama", "Dark Fantasy"},
		},
		{
			ID:           "7",
			Title:        "My Hero Academia",
			ThumbnailURL: "https://example.com/mha.jpg",
			Year:         2016,
			Status:       "ongoing",
			EpisodeCount: 120,
			Genres:       []string{"Action", "Superhero", "Shonen"},
		},
		{
			ID:           "8",
			Title:        "Fullmetal Alchemist Brotherhood",
			ThumbnailURL: "https://example.com/fma.jpg",
			Year:         2009,
			Status:       "completed",
			EpisodeCount: 64,
			Genres:       []string{"Action", "Adventure", "Fantasy"},
		},
		{
			ID:           "9",
			Title:        "Steins;Gate",
			ThumbnailURL: "https://example.com/steins-gate.jpg",
			Year:         2011,
			Status:       "completed",
			EpisodeCount: 24,
			Genres:       []string{"Sci-Fi", "Thriller", "Drama"},
		},
		{
			ID:           "10",
			Title:        "Code Geass",
			ThumbnailURL: "https://example.com/code-geass.jpg",
			Year:         2006,
			Status:       "completed",
			EpisodeCount: 50,
			Genres:       []string{"Mecha", "Sci-Fi", "Drama"},
		},
		{
			ID:           "11",
			Title:        "Bleach",
			ThumbnailURL: "https://example.com/bleach.jpg",
			Year:         2004,
			Status:       "completed",
			EpisodeCount: 366,
			Genres:       []string{"Action", "Supernatural", "Shonen"},
		},
		{
			ID:           "12",
			Title:        "Dragon Ball Z",
			ThumbnailURL: "https://example.com/dbz.jpg",
			Year:         1989,
			Status:       "completed",
			EpisodeCount: 291,
			Genres:       []string{"Action", "Adventure", "Shonen"},
		},
		{
			ID:           "13",
			Title:        "Fairy Tail",
			ThumbnailURL: "https://example.com/fairy-tail.jpg",
			Year:         2009,
			Status:       "completed",
			EpisodeCount: 328,
			Genres:       []string{"Action", "Adventure", "Fantasy"},
		},
		{
			ID:           "14",
			Title:        "Tokyo Ghoul",
			ThumbnailURL: "https://example.com/tokyo-ghoul.jpg",
			Year:         2014,
			Status:       "completed",
			EpisodeCount: 48,
			Genres:       []string{"Action", "Horror", "Psychological"},
		},
		{
			ID:           "15",
			Title:        "Jujutsu Kaisen",
			ThumbnailURL: "https://example.com/jujutsu-kaisen.jpg",
			Year:         2020,
			Status:       "ongoing",
			EpisodeCount: 60,
			Genres:       []string{"Action", "Supernatural", "Shonen"},
		},
	}
}

// TestAnimeSamaSearchService_ExactMatch verifies exact title matches rank first
func TestAnimeSamaSearchService_ExactMatch(t *testing.T) {
	service := NewAnimeSamaSearchService(testCatalogue())

	results, err := service.Search(context.Background(), "Naruto")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) == 0 {
		t.Fatal("expected results, got none")
	}

	// First result should be exact match "Naruto"
	if results[0].Title != "Naruto" {
		t.Errorf("expected first result 'Naruto', got %q", results[0].Title)
	}
}

// TestAnimeSamaSearchService_PartialMatch verifies partial matches are found
func TestAnimeSamaSearchService_PartialMatch(t *testing.T) {
	service := NewAnimeSamaSearchService(testCatalogue())

	results, err := service.Search(context.Background(), "Naruto")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should find both "Naruto" and "Naruto Shippuden"
	if len(results) < 2 {
		t.Fatalf("expected at least 2 results, got %d", len(results))
	}

	if results[1].Title != "Naruto Shippuden" {
		t.Errorf("expected second result 'Naruto Shippuden', got %q", results[1].Title)
	}
}

// TestAnimeSamaSearchService_CaseInsensitive verifies search is case-insensitive
func TestAnimeSamaSearchService_CaseInsensitive(t *testing.T) {
	service := NewAnimeSamaSearchService(testCatalogue())

	tests := []struct {
		query    string
		expected string
	}{
		{"naruto", "Naruto"},
		{"NARUTO", "Naruto"},
		{"NaRuTo", "Naruto"},
	}

	for _, tt := range tests {
		results, err := service.Search(context.Background(), tt.query)
		if err != nil {
			t.Fatalf("unexpected error for query %q: %v", tt.query, err)
		}

		if len(results) == 0 {
			t.Fatalf("expected results for query %q", tt.query)
		}

		if results[0].Title != tt.expected {
			t.Errorf("query %q: expected %q, got %q", tt.query, tt.expected, results[0].Title)
		}
	}
}

// TestAnimeSamaSearchService_EmptyQuery returns empty results
func TestAnimeSamaSearchService_EmptyQuery(t *testing.T) {
	service := NewAnimeSamaSearchService(testCatalogue())

	results, err := service.Search(context.Background(), "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("expected empty results for empty query, got %d", len(results))
	}
}

// TestAnimeSamaSearchService_MaxResults verifies max 50 results returned
func TestAnimeSamaSearchService_MaxResults(t *testing.T) {
	// Create catalogue with 100+ anime
	catalogue := []domain.AnimeSearchResult{}
	for i := 0; i < 100; i++ {
		catalogue = append(catalogue, domain.AnimeSearchResult{
			ID:           string(rune(i)),
			Title:        "anime",
			ThumbnailURL: "https://example.com/anime.jpg",
			Year:         2020,
			Status:       "ongoing",
			EpisodeCount: 12,
		})
	}

	service := NewAnimeSamaSearchService(catalogue)
	results, err := service.Search(context.Background(), "anime")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) > 50 {
		t.Errorf("expected max 50 results, got %d", len(results))
	}
}

// TestAnimeSamaSearchService_ContextCancellation verifies context cancellation is respected
func TestAnimeSamaSearchService_ContextCancellation(t *testing.T) {
	service := NewAnimeSamaSearchService(testCatalogue())

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := service.Search(ctx, "Naruto")
	if err == nil {
		t.Error("expected context error, got nil")
	}
}

// TestAnimeSamaSearchService_UnicodeNormalization verifies accented characters are handled
func TestAnimeSamaSearchService_UnicodeNormalization(t *testing.T) {
	catalogue := []domain.AnimeSearchResult{
		{
			ID:           "1",
			Title:        "Café Paradise",
			ThumbnailURL: "https://example.com/cafe.jpg",
			Year:         2020,
			Status:       "ongoing",
			EpisodeCount: 12,
		},
	}

	service := NewAnimeSamaSearchService(catalogue)

	// Search with and without accents should find the same result
	results1, _ := service.Search(context.Background(), "Café")
	results2, _ := service.Search(context.Background(), "Cafe")

	if len(results1) == 0 || len(results2) == 0 {
		t.Fatal("expected results for both accented and non-accented queries")
	}

	if results1[0].ID != results2[0].ID {
		t.Error("expected same result for accented and non-accented queries")
	}
}

// TestAnimeSamaSearchService_Whitespace verifies whitespace is trimmed
func TestAnimeSamaSearchService_Whitespace(t *testing.T) {
	service := NewAnimeSamaSearchService(testCatalogue())

	results1, _ := service.Search(context.Background(), "Naruto")
	results2, _ := service.Search(context.Background(), "  Naruto  ")
	results3, _ := service.Search(context.Background(), "\tNaruto\t")

	if len(results1) == 0 || len(results2) == 0 || len(results3) == 0 {
		t.Fatal("expected results for all whitespace variations")
	}

	if results1[0].ID != results2[0].ID || results1[0].ID != results3[0].ID {
		t.Error("expected same result for queries with different whitespace")
	}
}

// TestAnimeSamaSearchService_NoMatch returns empty results
func TestAnimeSamaSearchService_NoMatch(t *testing.T) {
	service := NewAnimeSamaSearchService(testCatalogue())

	results, err := service.Search(context.Background(), "NonExistentAnime12345")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("expected empty results for non-existent query, got %d", len(results))
	}
}

// TestAnimeSamaSearchService_RankingConsistency verifies ranking is consistent
func TestAnimeSamaSearchService_RankingConsistency(t *testing.T) {
	service := NewAnimeSamaSearchService(testCatalogue())

	// Run multiple times to ensure consistent ranking
	var previousIDs []string
	for i := 0; i < 3; i++ {
		results, _ := service.Search(context.Background(), "Dragon")
		if i == 0 {
			for _, r := range results {
				previousIDs = append(previousIDs, r.ID)
			}
		} else {
			for j, r := range results {
				if j < len(previousIDs) && r.ID != previousIDs[j] {
					t.Errorf("iteration %d: ranking not consistent at position %d", i, j)
				}
			}
		}
	}
}

// BenchmarkAnimeSamaSearchService_LargeDataset measures performance with 1000+ anime
func BenchmarkAnimeSamaSearchService_LargeDataset(b *testing.B) {
	// Create catalogue with 1000 anime
	catalogue := []domain.AnimeSearchResult{}
	for i := 0; i < 1000; i++ {
		catalogue = append(catalogue, domain.AnimeSearchResult{
			ID:           string(rune(i)),
			Title:        "anime title number " + string(rune(i)),
			ThumbnailURL: "https://example.com/img.jpg",
			Year:         2000 + (i % 24), // Year range 2000-2023
			Status:       "completed",
			EpisodeCount: 12,
			Genres:       []string{"Action"},
		})
	}

	service := NewAnimeSamaSearchService(catalogue)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = service.Search(context.Background(), "anime")
	}
}

// TestSearchWithFilters_GenreFilter verifies genre filtering
func TestSearchWithFilters_GenreFilter(t *testing.T) {
	service := NewAnimeSamaSearchService(testCatalogue())
	ctx := context.Background()

	tests := []struct {
		name           string
		genres         []string
		expectedTitles []string
	}{
		{
			name:           "filter by Action",
			genres:         []string{"Action"},
			expectedTitles: []string{"Naruto", "Naruto Shippuden", "One Piece", "Demon Slayer", "Attack on Titan", "My Hero Academia", "Fullmetal Alchemist Brotherhood", "Bleach", "Dragon Ball Z", "Fairy Tail", "Tokyo Ghoul", "Jujutsu Kaisen"},
		},
		{
			name:           "filter by Mystery",
			genres:         []string{"Mystery"},
			expectedTitles: []string{"Death Note"},
		},
		{
			name:           "filter by Shonen",
			genres:         []string{"Shonen"},
			expectedTitles: []string{"Naruto", "Naruto Shippuden", "Demon Slayer", "My Hero Academia", "Bleach", "Dragon Ball Z", "Jujutsu Kaisen"},
		},
		{
			name:           "filter by multiple genres (OR logic)",
			genres:         []string{"Mecha", "Superhero"},
			expectedTitles: []string{"Code Geass", "My Hero Academia"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filters := ports.SearchFilters{
				Query:  "", // No text filter, just genre
				Genres: tt.genres,
			}

			results, err := service.SearchWithFilters(ctx, filters)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(results) != len(tt.expectedTitles) {
				t.Errorf("expected %d results, got %d", len(tt.expectedTitles), len(results))
			}

			// Verify all expected titles are in results (order doesn't matter for genre filter only)
			resultTitles := make(map[string]bool)
			for _, r := range results {
				resultTitles[r.Title] = true
			}

			for _, expected := range tt.expectedTitles {
				if !resultTitles[expected] {
					t.Errorf("expected result %q not found", expected)
				}
			}
		})
	}
}

// TestSearchWithFilters_StatusFilter verifies status filtering
func TestSearchWithFilters_StatusFilter(t *testing.T) {
	service := NewAnimeSamaSearchService(testCatalogue())
	ctx := context.Background()

	tests := []struct {
		name          string
		status        string
		expectedCount int
	}{
		{
			name:          "filter completed anime",
			status:        "completed",
			expectedCount: 11,
		},
		{
			name:          "filter ongoing anime",
			status:        "ongoing",
			expectedCount: 4,
		},
		{
			name:          "case insensitive",
			status:        "COMPLETED",
			expectedCount: 11,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filters := ports.SearchFilters{
				Query:  "", // No text filter
				Status: tt.status,
			}

			results, err := service.SearchWithFilters(ctx, filters)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(results) != tt.expectedCount {
				t.Errorf("expected %d results, got %d", tt.expectedCount, len(results))
			}
		})
	}
}

// TestSearchWithFilters_YearRangeFilter verifies year range filtering
func TestSearchWithFilters_YearRangeFilter(t *testing.T) {
	service := NewAnimeSamaSearchService(testCatalogue())
	ctx := context.Background()

	tests := []struct {
		name           string
		yearMin        int
		yearMax        int
		expectedTitles []string
	}{
		{
			name:           "filter 2010 and later",
			yearMin:        2010,
			yearMax:        0,
			expectedTitles: []string{"Steins;Gate", "Attack on Titan", "Tokyo Ghoul", "My Hero Academia", "Demon Slayer", "Jujutsu Kaisen"},
		},
		{
			name:           "filter 2000s only",
			yearMin:        2000,
			yearMax:        2009,
			expectedTitles: []string{"Naruto", "Bleach", "Code Geass", "Death Note", "Naruto Shippuden", "Fullmetal Alchemist Brotherhood", "Fairy Tail"},
		},
		{
			name:           "filter before 2000",
			yearMin:        0,
			yearMax:        1999,
			expectedTitles: []string{"Dragon Ball Z", "One Piece"},
		},
		{
			name:           "filter exact year 2020",
			yearMin:        2020,
			yearMax:        2020,
			expectedTitles: []string{"Jujutsu Kaisen"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filters := ports.SearchFilters{
				Query:   "",
				YearMin: tt.yearMin,
				YearMax: tt.yearMax,
			}

			results, err := service.SearchWithFilters(ctx, filters)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(results) != len(tt.expectedTitles) {
				t.Errorf("expected %d results, got %d", len(tt.expectedTitles), len(results))
			}

			resultTitles := make(map[string]bool)
			for _, r := range results {
				resultTitles[r.Title] = true
			}

			for _, expected := range tt.expectedTitles {
				if !resultTitles[expected] {
					t.Errorf("expected result %q not found", expected)
				}
			}
		})
	}
}

// TestSearchWithFilters_CombinedFilters verifies multiple filters work together
func TestSearchWithFilters_CombinedFilters(t *testing.T) {
	service := NewAnimeSamaSearchService(testCatalogue())
	ctx := context.Background()

	// Query + Genre + Status + Year
	filters := ports.SearchFilters{
		Query:   "naruto",
		Genres:  []string{"Action"},
		Status:  "completed",
		YearMin: 2000,
		YearMax: 2010,
	}

	results, err := service.SearchWithFilters(ctx, filters)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should find:
	// - "Naruto" (2002, Action, completed) - YES
	// - "Naruto Shippuden" (2007, Action, completed) - YES
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}

	if results[0].Title != "Naruto" {
		t.Errorf("expected first result 'Naruto', got %q", results[0].Title)
	}
}

// TestSearchWithFilters_NoMatches verifies no results when filters don't match
func TestSearchWithFilters_NoMatches(t *testing.T) {
	service := NewAnimeSamaSearchService(testCatalogue())
	ctx := context.Background()

	filters := ports.SearchFilters{
		Query:   "naruto",
		Genres:  []string{"Mecha"}, // Naruto is not Mecha
		Status:  "completed",
		YearMin: 2000,
		YearMax: 2010,
	}

	results, err := service.SearchWithFilters(ctx, filters)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

// TestSearchWithFilters_EmptyFilters behaves like regular Search
func TestSearchWithFilters_EmptyFilters(t *testing.T) {
	service := NewAnimeSamaSearchService(testCatalogue())
	ctx := context.Background()

	filters := ports.SearchFilters{
		Query: "naruto",
		// All other filters empty
	}

	resultsFiltered, err := service.SearchWithFilters(ctx, filters)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resultsRegular, err := service.Search(ctx, "naruto")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resultsFiltered) != len(resultsRegular) {
		t.Errorf("empty filters should match regular search: %d vs %d", len(resultsFiltered), len(resultsRegular))
	}
}
