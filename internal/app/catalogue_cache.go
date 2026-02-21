package app

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
)

// CatalogueCache manages an in-memory cache of anime catalogue with TTL
type CatalogueCache struct {
	mu              sync.RWMutex
	catalogue       []domain.AnimeSearchResult
	lastRefreshTime time.Time
	ttl             time.Duration
	fetcher         func(ctx context.Context) ([]domain.AnimeSearchResult, error)
	stopRefresh     chan struct{}
	stopOnce        sync.Once
	refreshInterval time.Duration
}

// NewCatalogueCache creates a new catalogue cache with the given TTL and fetcher function
func NewCatalogueCache(ttl time.Duration, refreshInterval time.Duration, fetcher func(ctx context.Context) ([]domain.AnimeSearchResult, error)) *CatalogueCache {
	if refreshInterval <= 0 {
		refreshInterval = 30 * time.Minute
	}
	cache := &CatalogueCache{
		catalogue:       []domain.AnimeSearchResult{},
		ttl:             ttl,
		fetcher:         fetcher,
		stopRefresh:     make(chan struct{}),
		refreshInterval: refreshInterval,
	}

	// Start auto-refresh goroutine
	go cache.autoRefresh()

	return cache
}

// GetCatalogue returns the cached catalogue or fetches it if expired
func (c *CatalogueCache) GetCatalogue(ctx context.Context) ([]domain.AnimeSearchResult, error) {
	c.mu.RLock()
	if c.isFresh() && len(c.catalogue) > 0 {
		defer c.mu.RUnlock()
		return c.catalogue, nil
	}
	c.mu.RUnlock()

	// Cache is stale or empty, refresh it
	return c.Refresh(ctx)
}

// Refresh fetches fresh data and updates the cache
func (c *CatalogueCache) Refresh(ctx context.Context) ([]domain.AnimeSearchResult, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Double-check: another goroutine may have refreshed while we waited for the lock.
	if c.isFresh() && len(c.catalogue) > 0 {
		return c.catalogue, nil
	}

	// Fetch fresh data
	catalogue, err := c.fetcher(ctx)
	if err != nil {
		// If we have stale data, return it as fallback
		if len(c.catalogue) > 0 {
			log.Printf("Warning: Failed to refresh catalogue, using stale data: %v", err)
			return c.catalogue, nil
		}
		return nil, err
	}

	// Update cache
	c.catalogue = catalogue
	c.lastRefreshTime = time.Now()
	log.Printf("Catalogue cache refreshed: %d items loaded", len(catalogue))

	return c.catalogue, nil
}

// isFresh checks if the cache is still fresh (within TTL)
func (c *CatalogueCache) isFresh() bool {
	if c.lastRefreshTime.IsZero() {
		return false
	}
	return time.Since(c.lastRefreshTime) < c.ttl
}

// autoRefresh periodically refreshes the cache in the background
func (c *CatalogueCache) autoRefresh() {
	ticker := time.NewTicker(c.refreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ctx := context.Background()
			_, err := c.Refresh(ctx)
			if err != nil {
				log.Printf("Auto-refresh failed: %v", err)
			}
		case <-c.stopRefresh:
			return
		}
	}
}

// Stop stops the auto-refresh goroutine. Safe to call multiple times.
func (c *CatalogueCache) Stop() {
	c.stopOnce.Do(func() { close(c.stopRefresh) })
}

// Stats returns cache statistics
func (c *CatalogueCache) Stats() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var age time.Duration
	if !c.lastRefreshTime.IsZero() {
		age = time.Since(c.lastRefreshTime)
	}

	return map[string]interface{}{
		"items":           len(c.catalogue),
		"last_refresh":    c.lastRefreshTime,
		"cache_age_secs":  age.Seconds(),
		"ttl_secs":        c.ttl.Seconds(),
		"is_fresh":        c.isFresh(),
		"refresh_in_secs": (c.ttl - age).Seconds(),
	}
}
