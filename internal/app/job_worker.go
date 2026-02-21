package app

import (
	"context"
	"log/slog"
	"time"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
)

// JobWorker processes queued jobs in the background.
type JobWorker struct {
	jobRepo  domain.IJobRepository
	eventBus domain.IEventBus
	logger   *slog.Logger
}

// NewJobWorker creates a new JobWorker instance.
func NewJobWorker(jobRepo domain.IJobRepository, eventBus domain.IEventBus, logger *slog.Logger) *JobWorker {
	return &JobWorker{
		jobRepo:  jobRepo,
		eventBus: eventBus,
		logger:   logger,
	}
}

// Start begins processing jobs in a background loop.
func (jw *JobWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			jw.logger.Info("job worker stopping")
			return
		case <-ticker.C:
			jw.processJobs(ctx)
		}
	}
}

// processJobs fetches and executes pending jobs.
func (jw *JobWorker) processJobs(ctx context.Context) {
	// TODO: Fetch pending jobs from repository
	// For now, simulate job processing
}
