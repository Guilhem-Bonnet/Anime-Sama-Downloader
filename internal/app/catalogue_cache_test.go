package app

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
)

// mockFetcher returns a mock catalogue
func mockFetcher(ctx context.Context) ([]domain.AnimeSearchResult, error) {
	return []domain.AnimeSearchResult{
		{ID: "1", Title: "Naruto", Year: 2002, Status: "completed", Genres: []string{"Action"}},
		{ID: "2", Title: "One Piece", Year: 1999, Status: "ongoing", Genres: []string{"Adventure"}},
	}, nil
}

// errorFetcher always returns an error
func errorFetcher(ctx context.Context) ([]domain.AnimeSearchResult, error) {
	return nil, errors.New("fetch error")
}

func TestCatalogueCache_GetCatalogue_InitialLoad(t *testing.T) {
	cache := NewCatalogueCache(6*time.Hour, 6*time.Hour, mockFetcher)
	defer cache.Stop()

	ctx := context.Background()
	catalogue, err := cache.GetCatalogue(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(catalogue) != 2 {
		t.Errorf("expected 2 items, got %d", len(catalogue))
	}

	if catalogue[0].Title != "Naruto" {
		t.Errorf("expected first item 'Naruto', got %q", catalogue[0].Title)
	}
}

func TestCatalogueCache_GetCatalogue_CacheHit(t *testing.T) {
	fetchCount := 0
	countingFetcher := func(ctx context.Context) ([]domain.AnimeSearchResult, error) {
		fetchCount++
		return mockFetcher(ctx)
	}

	cache := NewCatalogueCache(6*time.Hour, 24*time.Hour, countingFetcher)
	defer cache.Stop()

	ctx := context.Background()

	// First call: cache miss (fetch)
	_, err := cache.GetCatalogue(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if fetchCount != 1 {
		t.Errorf("expected 1 fetch, got %d", fetchCount)
	}

	// Second call: cache hit (no fetch)
	_, err = cache.GetCatalogue(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if fetchCount != 1 {
		t.Errorf("expected still 1 fetch (cache hit), got %d", fetchCount)
	}
}

func TestCatalogueCache_GetCatalogue_CacheMiss_Expired(t *testing.T) {
	fetchCount := 0
	countingFetcher := func(ctx context.Context) ([]domain.AnimeSearchResult, error) {
		fetchCount++
		return mockFetcher(ctx)
	}

	// Set TTL to 100ms
	cache := NewCatalogueCache(100*time.Millisecond, 24*time.Hour, countingFetcher)
	defer cache.Stop()

	ctx := context.Background()

	// First call: cache miss (fetch)
	_, err := cache.GetCatalogue(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if fetchCount != 1 {
		t.Errorf("expected 1 fetch, got %d", fetchCount)
	}

	// Wait for cache to expire
	time.Sleep(150 * time.Millisecond)

	// Second call: cache expired (fetch again)
	_, err = cache.GetCatalogue(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if fetchCount != 2 {
		t.Errorf("expected 2 fetches (cache expired), got %d", fetchCount)
	}
}

func TestCatalogueCache_Refresh_ManualTrigger(t *testing.T) {
	cache := NewCatalogueCache(6*time.Hour, 24*time.Hour, mockFetcher)
	defer cache.Stop()

	ctx := context.Background()

	// Manual refresh
	catalogue, err := cache.Refresh(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(catalogue) != 2 {
		t.Errorf("expected 2 items, got %d", len(catalogue))
	}

	// Verify stats show fresh cache
	stats := cache.Stats()
	if !stats["is_fresh"].(bool) {
		t.Error("expected cache to be fresh after manual refresh")
	}
}

func TestCatalogueCache_ErrorHandling_FallbackToStale(t *testing.T) {
	// Start with working fetcher
	fetcher := mockFetcher
	cache := NewCatalogueCache(6*time.Hour, 24*time.Hour, func(ctx context.Context) ([]domain.AnimeSearchResult, error) {
		return fetcher(ctx)
	})
	defer cache.Stop()

	ctx := context.Background()

	// First load: success
	_, err := cache.GetCatalogue(ctx)
	if err != nil {
		t.Fatalf("unexpected error on first load: %v", err)
	}

	// Switch to error fetcher
	fetcher = errorFetcher

	// Manual refresh with error: should return stale data
	catalogue, err := cache.Refresh(ctx)
	if err != nil {
		t.Fatalf("expected no error (fallback to stale), got: %v", err)
	}

	if len(catalogue) != 2 {
		t.Errorf("expected stale data (2 items), got %d", len(catalogue))
	}
}

func TestCatalogueCache_Stats(t *testing.T) {
	cache := NewCatalogueCache(6*time.Hour, 24*time.Hour, mockFetcher)
	defer cache.Stop()

	ctx := context.Background()

	// Initial stats (no data)
	stats := cache.Stats()
	if stats["items"].(int) != 0 {
		t.Errorf("expected 0 items initially, got %d", stats["items"])
	}

	// Load cache
	_, err := cache.GetCatalogue(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Stats after loading
	stats = cache.Stats()
	if stats["items"].(int) != 2 {
		t.Errorf("expected 2 items after loading, got %d", stats["items"])
	}

	if !stats["is_fresh"].(bool) {
		t.Error("expected cache to be fresh after loading")
	}

	if stats["cache_age_secs"].(float64) < 0 {
		t.Error("expected non-negative cache age")
	}
}

func TestCatalogueCache_AutoRefresh(t *testing.T) {
	fetchCount := 0
	countingFetcher := func(ctx context.Context) ([]domain.AnimeSearchResult, error) {
		fetchCount++
		return mockFetcher(ctx)
	}

	// Set auto-refresh to 200ms
	cache := NewCatalogueCache(6*time.Hour, 200*time.Millisecond, countingFetcher)
	defer cache.Stop()

	ctx := context.Background()

	// Initial load
	_, err := cache.GetCatalogue(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if fetchCount != 1 {
		t.Errorf("expected 1 fetch initially, got %d", fetchCount)
	}

	// Wait for auto-refresh to trigger
	time.Sleep(300 * time.Millisecond)

	// Fetch count should have increased due to auto-refresh
	if fetchCount < 2 {
		t.Errorf("expected at least 2 fetches (auto-refresh), got %d", fetchCount)
	}
}
