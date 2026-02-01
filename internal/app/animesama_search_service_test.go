package app

import (
	"context"
	"testing"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
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
		},
		{
			ID:           "2",
			Title:        "Naruto Shippuden",
			ThumbnailURL: "https://example.com/naruto-shippuden.jpg",
			Year:         2007,
			Status:       "completed",
			EpisodeCount: 500,
		},
		{
			ID:           "3",
			Title:        "One Piece",
			ThumbnailURL: "https://example.com/one-piece.jpg",
			Year:         1999,
			Status:       "ongoing",
			EpisodeCount: 1050,
		},
		{
			ID:           "4",
			Title:        "Demon Slayer",
			ThumbnailURL: "https://example.com/demon-slayer.jpg",
			Year:         2019,
			Status:       "ongoing",
			EpisodeCount: 50,
		},
		{
			ID:           "5",
			Title:        "Death Note",
			ThumbnailURL: "https://example.com/death-note.jpg",
			Year:         2006,
			Status:       "completed",
			EpisodeCount: 37,
		},
		{
			ID:           "6",
			Title:        "Attack on Titan",
			ThumbnailURL: "https://example.com/aot.jpg",
			Year:         2013,
			Status:       "completed",
			EpisodeCount: 94,
		},
		{
			ID:           "7",
			Title:        "My Hero Academia",
			ThumbnailURL: "https://example.com/mha.jpg",
			Year:         2016,
			Status:       "ongoing",
			EpisodeCount: 120,
		},
		{
			ID:           "8",
			Title:        "Fullmetal Alchemist Brotherhood",
			ThumbnailURL: "https://example.com/fma.jpg",
			Year:         2009,
			Status:       "completed",
			EpisodeCount: 64,
		},
		{
			ID:           "9",
			Title:        "Steins;Gate",
			ThumbnailURL: "https://example.com/steins-gate.jpg",
			Year:         2011,
			Status:       "completed",
			EpisodeCount: 24,
		},
		{
			ID:           "10",
			Title:        "Code Geass",
			ThumbnailURL: "https://example.com/code-geass.jpg",
			Year:         2006,
			Status:       "completed",
			EpisodeCount: 50,
		},
		{
			ID:           "11",
			Title:        "Bleach",
			ThumbnailURL: "https://example.com/bleach.jpg",
			Year:         2004,
			Status:       "completed",
			EpisodeCount: 366,
		},
		{
			ID:           "12",
			Title:        "Dragon Ball Z",
			ThumbnailURL: "https://example.com/dbz.jpg",
			Year:         1989,
			Status:       "completed",
			EpisodeCount: 291,
		},
		{
			ID:           "13",
			Title:        "Fairy Tail",
			ThumbnailURL: "https://example.com/fairy-tail.jpg",
			Year:         2009,
			Status:       "completed",
			EpisodeCount: 328,
		},
		{
			ID:           "14",
			Title:        "Tokyo Ghoul",
			ThumbnailURL: "https://example.com/tokyo-ghoul.jpg",
			Year:         2014,
			Status:       "completed",
			EpisodeCount: 48,
		},
		{
			ID:           "15",
			Title:        "Jujutsu Kaisen",
			ThumbnailURL: "https://example.com/jujutsu-kaisen.jpg",
			Year:         2020,
			Status:       "ongoing",
			EpisodeCount: 60,
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
			ThumbnailURL: "https://example.com/anime.jpg",
			Year:         2020,
			Status:       "ongoing",
			EpisodeCount: 12,
		})
	}

	service := NewAnimeSamaSearchService(catalogue)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.Search(context.Background(), "anime")
	}
}
