package memorybus

import (
	"sync"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/ports"
)

type Bus struct {
	mu    sync.Mutex
	subs  map[chan ports.Event]struct{}
	alive bool
}

func New() *Bus {
	return &Bus{subs: make(map[chan ports.Event]struct{}), alive: true}
}

func (b *Bus) Publish(topic string, payload []byte) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if !b.alive {
		return
	}
	evt := ports.Event{Topic: topic, Payload: payload}
	for ch := range b.subs {
		select {
		case ch <- evt:
		default:
			// drop si le client est trop lent
		}
	}
}

func (b *Bus) Subscribe() (<-chan ports.Event, func()) {
	ch := make(chan ports.Event, 64)
	b.mu.Lock()
	if !b.alive {
		close(ch)
		b.mu.Unlock()
		return ch, func() {}
	}
	b.subs[ch] = struct{}{}
	b.mu.Unlock()

	cancel := func() {
		b.mu.Lock()
		if _, ok := b.subs[ch]; ok {
			delete(b.subs, ch)
			close(ch)
		}
		b.mu.Unlock()
	}

	return ch, cancel
}
