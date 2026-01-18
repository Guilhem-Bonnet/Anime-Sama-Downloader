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
