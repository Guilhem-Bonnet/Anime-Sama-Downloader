package app

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestDynamicLimiter_AcquireRelease(t *testing.T) {
	l := NewDynamicLimiter(1)

	ctx := context.Background()
	if err := l.Acquire(ctx); err != nil {
		t.Fatalf("Acquire: %v", err)
	}

	acquired := make(chan struct{})
	go func() {
		_ = l.Acquire(ctx)
		close(acquired)
	}()

	select {
	case <-acquired:
		t.Fatalf("second acquire should block")
	case <-time.After(50 * time.Millisecond):
	}

	l.Release()
	select {
	case <-acquired:
	case <-time.After(250 * time.Millisecond):
		t.Fatalf("second acquire should have proceeded")
	}

	l.Release()
}

func TestDynamicLimiter_SetLimitWakesWaiters(t *testing.T) {
	l := NewDynamicLimiter(1)
	ctx := context.Background()

	if err := l.Acquire(ctx); err != nil {
		t.Fatalf("Acquire: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	done := make(chan struct{})
	go func() {
		defer wg.Done()
		_ = l.Acquire(ctx)
		close(done)
	}()

	// Toujours bloqué tant que limit=1.
	select {
	case <-done:
		t.Fatalf("acquire should block")
	case <-time.After(50 * time.Millisecond):
	}

	// En augmentant le plafond, le waiter doit passer.
	l.SetLimit(2)
	select {
	case <-done:
	case <-time.After(250 * time.Millisecond):
		t.Fatalf("waiter should have been woken by SetLimit")
	}

	l.Release()
	l.Release()
	wg.Wait()
}

func TestDynamicLimiter_AcquireHonorsContext(t *testing.T) {
	l := NewDynamicLimiter(1)
	if err := l.Acquire(context.Background()); err != nil {
		t.Fatalf("Acquire: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	start := time.Now()
	err := l.Acquire(ctx)
	if err == nil {
		t.Fatalf("expected error")
	}
	if time.Since(start) < 40*time.Millisecond {
		t.Fatalf("expected acquire to wait for context timeout")
	}

	l.Release()
}

// TestDynamicLimiter_LimitMethod tests the Limit() getter
func TestDynamicLimiter_LimitMethod(t *testing.T) {
	tests := []struct {
		initial  int
		expected int
	}{
		{5, 5},
		{1, 1},
		{100, 100},
		{0, 1},  // defaults to 1
		{-5, 1}, // defaults to 1
	}

	for _, tt := range tests {
		l := NewDynamicLimiter(tt.initial)
		if l.Limit() != tt.expected {
			t.Fatalf("NewDynamicLimiter(%d).Limit() = %d, expected %d", tt.initial, l.Limit(), tt.expected)
		}
	}
}

// TestDynamicLimiter_InFlightMethod tests the InFlight() getter
func TestDynamicLimiter_InFlightMethod(t *testing.T) {
	l := NewDynamicLimiter(5)

	// Initially empty
	if l.InFlight() != 0 {
		t.Fatalf("initial InFlight() should be 0, got %d", l.InFlight())
	}

	// After one acquire
	ctx := context.Background()
	l.Acquire(ctx)
	if l.InFlight() != 1 {
		t.Fatalf("InFlight() after 1 Acquire should be 1, got %d", l.InFlight())
	}

	// After two acquires
	l.Acquire(ctx)
	if l.InFlight() != 2 {
		t.Fatalf("InFlight() after 2 Acquires should be 2, got %d", l.InFlight())
	}

	// After one release
	l.Release()
	if l.InFlight() != 1 {
		t.Fatalf("InFlight() after 1 Release should be 1, got %d", l.InFlight())
	}

	// After all releases
	l.Release()
	if l.InFlight() != 0 {
		t.Fatalf("InFlight() after all Releases should be 0, got %d", l.InFlight())
	}
}

// TestDynamicLimiter_SetLimitZeroDefaultsToOne tests SetLimit(0) defaults to 1
func TestDynamicLimiter_SetLimitZeroDefaultsToOne(t *testing.T) {
	l := NewDynamicLimiter(5)
	l.SetLimit(0)
	if l.Limit() != 1 {
		t.Fatalf("SetLimit(0) should default to 1, got %d", l.Limit())
	}
}

// TestDynamicLimiter_SetLimitNegativeDefaultsToOne tests SetLimit(-n) defaults to 1
func TestDynamicLimiter_SetLimitNegativeDefaultsToOne(t *testing.T) {
	l := NewDynamicLimiter(5)
	l.SetLimit(-10)
	if l.Limit() != 1 {
		t.Fatalf("SetLimit(-10) should default to 1, got %d", l.Limit())
	}
}

// TestDynamicLimiter_SetLimitIncrease tests increasing limit wakes waiters
func TestDynamicLimiter_SetLimitIncrease(t *testing.T) {
	l := NewDynamicLimiter(1)
	ctx := context.Background()

	l.Acquire(ctx)
	if l.Limit() != 1 || l.InFlight() != 1 {
		t.Fatalf("expected limit 1, in-flight 1")
	}

	// Now try to acquire second (should block due to limit)
	acquired := make(chan struct{})
	go func() {
		l.Acquire(ctx)
		close(acquired)
	}()

	select {
	case <-acquired:
		t.Fatalf("second acquire should block")
	case <-time.After(50 * time.Millisecond):
	}

	// Increase limit
	l.SetLimit(2)
	if l.Limit() != 2 {
		t.Fatalf("expected limit 2, got %d", l.Limit())
	}

	// Now second acquire should proceed
	select {
	case <-acquired:
	case <-time.After(250 * time.Millisecond):
		t.Fatalf("second acquire should have proceeded after SetLimit")
	}

	if l.InFlight() != 2 {
		t.Fatalf("expected 2 in-flight, got %d", l.InFlight())
	}

	l.Release()
	l.Release()
}

// TestDynamicLimiter_MultipleSetLimit tests multiple SetLimit calls
func TestDynamicLimiter_MultipleSetLimit(t *testing.T) {
	l := NewDynamicLimiter(1)

	l.SetLimit(3)
	if l.Limit() != 3 {
		t.Fatalf("expected 3, got %d", l.Limit())
	}

	l.SetLimit(5)
	if l.Limit() != 5 {
		t.Fatalf("expected 5, got %d", l.Limit())
	}

	l.SetLimit(2)
	if l.Limit() != 2 {
		t.Fatalf("expected 2, got %d", l.Limit())
	}
}
