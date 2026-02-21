package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// aniListHTTPClient wraps an *http.Client with retry-on-429 logic and
// an in-memory response cache suitable for AniList GraphQL queries.
type aniListHTTPClient struct {
	client   *http.Client
	endpoint string

	// Simple TTL cache: key = SHA of request body → cached response bytes
	mu    sync.RWMutex
	cache map[string]cacheEntry
}

type cacheEntry struct {
	data      []byte
	expiresAt time.Time
}

const (
	aniListCacheTTL     = 5 * time.Minute
	aniListMaxRetries   = 3
	aniListBaseBackoff  = 2 * time.Second
	aniListCacheMaxSize = 200
)

// NewAniListHTTPClient creates a shared HTTP client with retry and cache.
func NewAniListHTTPClient(endpoint string) *aniListHTTPClient {
	return &aniListHTTPClient{
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
		endpoint: endpoint,
		cache:    make(map[string]cacheEntry),
	}
}

// do executes a GraphQL request with retry on 429 and caching.
func (c *aniListHTTPClient) do(ctx context.Context, req aniListGraphQLRequest, out any) error {
	b, err := json.Marshal(req)
	if err != nil {
		return err
	}

	cacheKey := string(b)

	// Check cache first
	if cached, ok := c.getFromCache(cacheKey); ok {
		return json.Unmarshal(cached, out)
	}

	// Execute with retry
	var respBody []byte
	for attempt := 0; attempt <= aniListMaxRetries; attempt++ {
		if attempt > 0 {
			backoff := aniListBaseBackoff * time.Duration(1<<(attempt-1))
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
			}
		}

		respBody, err = c.executeRequest(ctx, b)
		if err == nil {
			break
		}

		// Only retry on 429 (rate limit)
		if !isRateLimitError(err) {
			return err
		}
	}
	if err != nil {
		return err
	}

	// Store in cache
	c.putInCache(cacheKey, respBody)

	return json.Unmarshal(respBody, out)
}

func (c *aniListHTTPClient) executeRequest(ctx context.Context, body []byte) ([]byte, error) {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("User-Agent", "asd-server")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 429 {
		// Parse Retry-After header if present
		retryAfter := aniListBaseBackoff
		if ra := resp.Header.Get("Retry-After"); ra != "" {
			if secs, parseErr := strconv.Atoi(ra); parseErr == nil {
				retryAfter = time.Duration(secs) * time.Second
			}
		}
		return nil, &rateLimitError{retryAfter: retryAfter}
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("anilist http error: %s", resp.Status)
	}

	return respBody, nil
}

func (c *aniListHTTPClient) getFromCache(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.cache[key]
	if !ok || time.Now().After(entry.expiresAt) {
		return nil, false
	}
	return entry.data, true
}

func (c *aniListHTTPClient) putInCache(key string, data []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Evict expired entries if cache is full
	if len(c.cache) >= aniListCacheMaxSize {
		now := time.Now()
		for k, v := range c.cache {
			if now.After(v.expiresAt) {
				delete(c.cache, k)
			}
		}
		// If still full, clear oldest half
		if len(c.cache) >= aniListCacheMaxSize {
			count := 0
			for k := range c.cache {
				delete(c.cache, k)
				count++
				if count >= aniListCacheMaxSize/2 {
					break
				}
			}
		}
	}

	c.cache[key] = cacheEntry{
		data:      data,
		expiresAt: time.Now().Add(aniListCacheTTL),
	}
}

// rateLimitError signals a 429 response from AniList.
type rateLimitError struct {
	retryAfter time.Duration
}

func (e *rateLimitError) Error() string {
	return fmt.Sprintf("anilist rate limited (retry after %s)", e.retryAfter)
}

func isRateLimitError(err error) bool {
	_, ok := err.(*rateLimitError)
	return ok
}
