package app

import (
	"context"
	"os"
	"strconv"
	"sync"
	"time"
)

// RateLimiter implements a token bucket rate limiter for HTTP requests.
// It limits requests per second to avoid getting banned by external services.
type RateLimiter struct {
	mu         sync.Mutex
	rps        float64        // requests per second
	tokens     float64        // current available tokens
	maxTokens  float64        // maximum tokens (burst)
	lastUpdate time.Time      // last token update time
	waitCount  int            // number of requests waiting
}

// NewRateLimiter creates a new rate limiter with the specified requests per second.
// If rps <= 0, defaults to 1.0 (1 request per second).
func NewRateLimiter(rps float64) *RateLimiter {
	if rps <= 0 {
		rps = 1.0
	}
	return &RateLimiter{
		rps:        rps,
		tokens:     rps, // start with full bucket
		maxTokens:  rps, // allow burst equal to RPS
		lastUpdate: time.Now(),
	}
}

// NewRateLimiterFromEnv creates a rate limiter from RATE_LIMIT_RPS env var.
// Defaults to 1.0 if not set or invalid.
func NewRateLimiterFromEnv() *RateLimiter {
	rps := 1.0
	if val := os.Getenv("RATE_LIMIT_RPS"); val != "" {
		if parsed, err := strconv.ParseFloat(val, 64); err == nil && parsed > 0 {
			rps = parsed
		}
	}
	return NewRateLimiter(rps)
}

// RPS returns the current rate limit in requests per second.
func (r *RateLimiter) RPS() float64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.rps
}

// SetRPS updates the rate limit. Takes effect immediately.
func (r *RateLimiter) SetRPS(rps float64) {
	if rps <= 0 {
		rps = 1.0
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rps = rps
	r.maxTokens = rps
}

// WaitCount returns the number of goroutines currently waiting for a token.
func (r *RateLimiter) WaitCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.waitCount
}

// Wait blocks until a token is available or context is cancelled.
// Returns nil on success, context error on cancellation.
func (r *RateLimiter) Wait(ctx context.Context) error {
	r.mu.Lock()
	r.waitCount++
	r.mu.Unlock()

	defer func() {
		r.mu.Lock()
		r.waitCount--
		r.mu.Unlock()
	}()

	for {
		r.mu.Lock()
		r.refillTokens()

		if r.tokens >= 1.0 {
			r.tokens -= 1.0
			r.mu.Unlock()
			return nil
		}

		// Calculate wait time for next token
		waitDuration := time.Duration(float64(time.Second) / r.rps)
		r.mu.Unlock()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(waitDuration):
			// Try again
		}
	}
}

// TryAcquire attempts to acquire a token without waiting.
// Returns true if token was acquired, false if rate limited.
func (r *RateLimiter) TryAcquire() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.refillTokens()

	if r.tokens >= 1.0 {
		r.tokens -= 1.0
		return true
	}
	return false
}

// refillTokens adds tokens based on elapsed time. Must be called with lock held.
func (r *RateLimiter) refillTokens() {
	now := time.Now()
	elapsed := now.Sub(r.lastUpdate).Seconds()
	r.lastUpdate = now

	// Add tokens based on elapsed time
	r.tokens += elapsed * r.rps
	if r.tokens > r.maxTokens {
		r.tokens = r.maxTokens
	}
}
