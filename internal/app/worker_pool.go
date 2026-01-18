package app

import (
	"context"
	"sync"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/ports"
	"github.com/rs/zerolog"
)

// WorkerPool gère un pool de workers ajustable à chaud.
// Les workers sont arrêtés via cancel() sur leur contexte.
//
// SetCount() peut être appelé plusieurs fois et est thread-safe.
type WorkerPool struct {
	parent context.Context

	logger zerolog.Logger
	repo   ports.JobRepository
	bus    ports.EventBus
	opts   WorkerOptions

	mu      sync.Mutex
	cancels []context.CancelFunc
	wg      sync.WaitGroup
}

func NewWorkerPool(parent context.Context, logger zerolog.Logger, repo ports.JobRepository, bus ports.EventBus, opts WorkerOptions) *WorkerPool {
	if parent == nil {
		parent = context.Background()
	}
	return &WorkerPool{parent: parent, logger: logger, repo: repo, bus: bus, opts: opts}
}

func (p *WorkerPool) Count() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.cancels)
}

func (p *WorkerPool) SetCount(n int) {
	if n <= 0 {
		n = 1
	}

	p.mu.Lock()
	current := len(p.cancels)

	if n == current {
		p.mu.Unlock()
		return
	}

	if n > current {
		for i := current; i < n; i++ {
			ctx, cancel := context.WithCancel(p.parent)
			p.cancels = append(p.cancels, cancel)
			idx := i
			p.wg.Add(1)
			go func() {
				defer p.wg.Done()
				w := NewWorker(p.logger.With().Int("worker", idx+1).Logger(), p.repo, p.bus, p.opts)
				w.Run(ctx)
			}()
		}
		p.mu.Unlock()
		return
	}

	// n < current : stoppe les derniers workers
	toStop := append([]context.CancelFunc(nil), p.cancels[n:]...)
	p.cancels = p.cancels[:n]
	p.mu.Unlock()

	for _, cancel := range toStop {
		cancel()
	}
}

func (p *WorkerPool) Close() {
	p.mu.Lock()
	toStop := append([]context.CancelFunc(nil), p.cancels...)
	p.cancels = nil
	p.mu.Unlock()

	for _, cancel := range toStop {
		cancel()
	}
	p.wg.Wait()
}
