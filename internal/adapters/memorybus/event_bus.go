// Package memorybus implements an in-memory event bus.
package memorybus

import (
	"sync"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
)

// EventBus is an in-memory publish/subscribe event bus.
type EventBus struct {
	subscribers map[string][]domain.EventHandler
	mu          sync.RWMutex
}

// NewEventBus creates a new EventBus.
func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[string][]domain.EventHandler),
	}
}

// Subscribe registers a handler for an event.
func (b *EventBus) Subscribe(event string, handler domain.EventHandler) func() {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Create a wrapper that allows tracking
	index := len(b.subscribers[event])
	b.subscribers[event] = append(b.subscribers[event], handler)

	// Return unsubscribe function
	return func() {
		b.mu.Lock()
		defer b.mu.Unlock()

		handlers := b.subscribers[event]
		if index < len(handlers) {
			// Remove handler from slice by index
			b.subscribers[event] = append(handlers[:index], handlers[index+1:]...)
		}
	}
}

// Emit publishes an event to all subscribers.
func (b *EventBus) Emit(event string, payload interface{}) {
	b.mu.RLock()
	handlers := b.subscribers[event]
	b.mu.RUnlock()

	// Call handlers in goroutines to avoid blocking
	for _, handler := range handlers {
		go handler(payload)
	}
}
