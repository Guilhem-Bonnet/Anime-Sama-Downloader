package app

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestRateLimiter_NewRateLimiter(t *testing.T) {
	tests := []struct {
		name        string
		rps         float64
		expectedRPS float64
	}{
		{"positive rps", 2.0, 2.0},
		{"zero rps defaults to 1", 0, 1.0},
		{"negative rps defaults to 1", -5.0, 1.0},
		{"fractional rps", 0.5, 0.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limiter := NewRateLimiter(tt.rps)
			if limiter.RPS() != tt.expectedRPS {
				t.Errorf("expected RPS %f, got %f", tt.expectedRPS, limiter.RPS())
			}
		})
	}
}

func TestRateLimiter_SetRPS(t *testing.T) {
	limiter := NewRateLimiter(1.0)

	limiter.SetRPS(5.0)
	if limiter.RPS() != 5.0 {
		t.Errorf("expected RPS 5.0, got %f", limiter.RPS())
	}

	// Zero should default to 1
	limiter.SetRPS(0)
	if limiter.RPS() != 1.0 {
		t.Errorf("expected RPS 1.0 after setting 0, got %f", limiter.RPS())
	}
}

func TestRateLimiter_TryAcquire(t *testing.T) {
	limiter := NewRateLimiter(1.0) // 1 request per second

	// First request should succeed (bucket starts full)
	if !limiter.TryAcquire() {
		t.Error("first TryAcquire should succeed")
	}

	// Second immediate request should fail (no tokens)
	if limiter.TryAcquire() {
		t.Error("second immediate TryAcquire should fail")
	}

	// Wait for refill
	time.Sleep(1100 * time.Millisecond)

	// Now should succeed again
	if !limiter.TryAcquire() {
		t.Error("TryAcquire after refill should succeed")
	}
}

func TestRateLimiter_Wait_Success(t *testing.T) {
	limiter := NewRateLimiter(10.0) // 10 requests per second

	ctx := context.Background()
	start := time.Now()

	// Should acquire immediately (bucket full)
	err := limiter.Wait(ctx)
	if err != nil {
		t.Errorf("Wait should succeed: %v", err)
	}

	elapsed := time.Since(start)
	if elapsed > 50*time.Millisecond {
		t.Errorf("first Wait took too long: %v", elapsed)
	}
}

func TestRateLimiter_Wait_RateLimited(t *testing.T) {
	limiter := NewRateLimiter(2.0) // 2 requests per second

	ctx := context.Background()

	// Drain the bucket
	limiter.Wait(ctx)
	limiter.Wait(ctx)

	// Third request should wait
	start := time.Now()
	err := limiter.Wait(ctx)
	if err != nil {
		t.Errorf("Wait should succeed: %v", err)
	}

	elapsed := time.Since(start)
	if elapsed < 400*time.Millisecond {
		t.Errorf("Wait should have been rate limited, elapsed: %v", elapsed)
	}
}

func TestRateLimiter_Wait_ContextCancelled(t *testing.T) {
	limiter := NewRateLimiter(0.1) // Very slow: 1 request per 10 seconds

	// Drain the bucket
	limiter.TryAcquire()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := limiter.Wait(ctx)
	if err != context.DeadlineExceeded {
		t.Errorf("expected DeadlineExceeded, got: %v", err)
	}
}

func TestRateLimiter_WaitCount(t *testing.T) {
	limiter := NewRateLimiter(0.5) // Very slow

	// Drain the bucket
	limiter.TryAcquire()

	if limiter.WaitCount() != 0 {
		t.Error("WaitCount should be 0 initially")
	}

	// Start a waiting goroutine
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()
		limiter.Wait(ctx)
	}()

	// Give goroutine time to start waiting
	time.Sleep(50 * time.Millisecond)

	if limiter.WaitCount() != 1 {
		t.Errorf("expected WaitCount 1, got %d", limiter.WaitCount())
	}

	wg.Wait()

	if limiter.WaitCount() != 0 {
		t.Errorf("expected WaitCount 0 after goroutine exits, got %d", limiter.WaitCount())
	}
}

func TestRateLimiter_ConcurrentAccess(t *testing.T) {
	limiter := NewRateLimiter(100.0) // High rate for fast test

	var wg sync.WaitGroup
	var successCount atomic.Int32
	numGoroutines := 50

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			if err := limiter.Wait(ctx); err == nil {
				successCount.Add(1)
			}
		}()
	}

	wg.Wait()

	// All should succeed with high rate limit
	if successCount.Load() != int32(numGoroutines) {
		t.Errorf("expected %d successes, got %d", numGoroutines, successCount.Load())
	}
}

func TestRateLimiter_RespectRateLimit(t *testing.T) {
	limiter := NewRateLimiter(5.0) // 5 requests per second

	ctx := context.Background()
	numRequests := 10

	start := time.Now()
	for i := 0; i < numRequests; i++ {
		if err := limiter.Wait(ctx); err != nil {
			t.Fatalf("Wait failed: %v", err)
		}
	}
	elapsed := time.Since(start)

	// 10 requests at 5 RPS should take at least 800ms
	// (first 5 are burst from tokens, then 5 more at ~200ms each = ~1s)
	// We use a conservative minimum to avoid flaky tests
	expectedMin := 800 * time.Millisecond
	if elapsed < expectedMin {
		t.Errorf("rate limit not respected: %d requests in %v (expected >= %v)", numRequests, elapsed, expectedMin)
	}
}
