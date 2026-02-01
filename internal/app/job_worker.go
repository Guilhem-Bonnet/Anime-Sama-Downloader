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

// executeJob runs a single job and updates its progress.
func (jw *JobWorker) executeJob(ctx context.Context, job *domain.Job) {
	jw.eventBus.Emit(domain.EventJobStarted, map[string]interface{}{
		"job_id": job.ID,
	})

	// Simulate work with progress updates
	for progress := 10; progress <= 100; progress += 10 {
		select {
		case <-ctx.Done():
			return
		case <-time.After(100 * time.Millisecond):
			// Update job progress
			jw.eventBus.Emit(domain.EventJobProgress, map[string]interface{}{
				"job_id":   job.ID,
				"progress": progress,
			})
		}
	}

	// Mark as completed
	jw.eventBus.Emit(domain.EventJobCompleted, map[string]interface{}{
		"job_id": job.ID,
	})
}
