package app

import (
	"context"
	"sync"
)

// DynamicLimiter limite le nombre d'opérations concurrentes.
// Le plafond peut être modifié à chaud via SetLimit.
//
// Le pattern est volontairement simple et ne dépend pas de packages externes.
// Acquire respecte le contexte.
type DynamicLimiter struct {
	mu       sync.Mutex
	limit    int
	inFlight int
	notify   chan struct{}
}

func NewDynamicLimiter(limit int) *DynamicLimiter {
	if limit <= 0 {
		limit = 1
	}
	return &DynamicLimiter{limit: limit, notify: make(chan struct{})}
}

func (l *DynamicLimiter) Limit() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.limit
}

func (l *DynamicLimiter) InFlight() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.inFlight
}

func (l *DynamicLimiter) SetLimit(limit int) {
	if limit <= 0 {
		limit = 1
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	if l.limit == limit {
		return
	}
	l.limit = limit
	l.signalLocked()
}

func (l *DynamicLimiter) Acquire(ctx context.Context) error {
	for {
		l.mu.Lock()
		limit := l.limit
		if limit <= 0 {
			limit = 1
		}
		if l.inFlight < limit {
			l.inFlight++
			l.mu.Unlock()
			return nil
		}
		ch := l.notify
		l.mu.Unlock()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ch:
		}
	}
}

func (l *DynamicLimiter) Release() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.inFlight > 0 {
		l.inFlight--
	}
	l.signalLocked()
}

func (l *DynamicLimiter) signalLocked() {
	// Réveille tous les waiters en fermant le channel et en recréant.
	// C'est OK même si aucun waiter n'écoute.
	close(l.notify)
	l.notify = make(chan struct{})
}
