package app

import (
	"context"
	"errors"
	"time"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/ports"
	"github.com/rs/zerolog"
)

type WorkerOptions struct {
	PollInterval time.Duration
	StepInterval time.Duration
	Steps        int
	Destination  string
	// DestinationFunc, si défini, permet de résoudre la destination à l'exécution (ex: depuis les settings).
	DestinationFunc func(ctx context.Context) (string, error)
}

func DefaultWorkerOptions() WorkerOptions {
	return WorkerOptions{
		PollInterval: 750 * time.Millisecond,
		StepInterval: 400 * time.Millisecond,
		Steps:        10,
		Destination:  "videos",
	}
}

type Worker struct {
	logger zerolog.Logger
	repo   ports.JobRepository
	bus    ports.EventBus
	opts   WorkerOptions
	execs  ExecutorRegistry
}

func NewWorker(logger zerolog.Logger, repo ports.JobRepository, bus ports.EventBus, opts WorkerOptions) *Worker {
	if opts.PollInterval <= 0 {
		opts.PollInterval = DefaultWorkerOptions().PollInterval
	}
	if opts.StepInterval <= 0 {
		opts.StepInterval = DefaultWorkerOptions().StepInterval
	}
	if opts.Steps <= 0 {
		opts.Steps = DefaultWorkerOptions().Steps
	}
	if opts.Destination == "" {
		opts.Destination = DefaultWorkerOptions().Destination
	}
	return &Worker{logger: logger, repo: repo, bus: bus, opts: opts, execs: DefaultExecutorRegistry()}
}

func RunWorkers(ctx context.Context, logger zerolog.Logger, repo ports.JobRepository, bus ports.EventBus, count int, opts WorkerOptions) {
	if count <= 0 {
		count = 1
	}
	for i := 0; i < count; i++ {
		w := NewWorker(logger.With().Int("worker", i+1).Logger(), repo, bus, opts)
		go w.Run(ctx)
	}
}

func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.opts.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			job, err := w.repo.ClaimNextQueued(ctx)
			if err != nil {
				// Adapter-specific: on traite tout "not found" comme "rien à faire".
				if errors.Is(err, ErrNotFound) {
					continue
				}
				w.logger.Error().Err(err).Msg("claim next job failed")
				continue
			}

			w.execute(ctx, job)
		}
	}
}

func (w *Worker) execute(ctx context.Context, job domain.Job) {
	w.logger.Info().Str("job_id", job.ID).Str("type", job.Type).Msg("job claimed")
	PublishJobEvent(w.bus, "job.started", job)

	isCanceled := func() (bool, error) {
		current, err := w.repo.Get(ctx, job.ID)
		if err != nil {
			return false, err
		}
		return current.State == domain.JobCanceled, nil
	}

	updateProgress := func(progress float64) error {
		updated, err := w.repo.UpdateProgress(ctx, job.ID, progress)
		if err != nil {
			return err
		}
		PublishJobEvent(w.bus, "job.progress", updated)
		return nil
	}

	updateResult := func(resultJSON []byte) error {
		updated, err := w.repo.UpdateResult(ctx, job.ID, resultJSON)
		if err != nil {
			return err
		}
		PublishJobEvent(w.bus, "job.result", updated)
		return nil
	}

	exec := w.execs.Get(job.Type)
	destination := w.opts.Destination
	if w.opts.DestinationFunc != nil {
		if d, err := w.opts.DestinationFunc(ctx); err == nil && d != "" {
			destination = d
		}
	}
	err := exec.Execute(ctx, job, ExecEnv{
		UpdateProgress: updateProgress,
		UpdateResult:   updateResult,
		IsCanceled:     isCanceled,
		StepInterval:   w.opts.StepInterval,
		Steps:          w.opts.Steps,
		Destination:    destination,
	})
	if err != nil {
		w.logger.Error().Err(err).Str("job_id", job.ID).Msg("executor failed")
		code := "executor_error"
		message := err.Error()
		var coded *CodedError
		if errors.As(err, &coded) {
			if coded.Code != "" {
				code = coded.Code
			}
			if coded.Message != "" {
				message = coded.Message
			}
		}
		_, _ = w.repo.UpdateError(ctx, job.ID, code, message)
		failed, err2 := w.repo.UpdateState(ctx, job.ID, domain.JobRunning, domain.JobFailed)
		if err2 == nil {
			PublishJobEvent(w.bus, "job.failed", failed)
		}
		return
	}

	canceled, err := isCanceled()
	if err != nil {
		w.logger.Error().Err(err).Str("job_id", job.ID).Msg("failed to reload job")
		return
	}
	if canceled {
		w.logger.Info().Str("job_id", job.ID).Msg("job canceled")
		return
	}

	// Terminer: respecter running -> muxing -> completed.
	phase, err := w.repo.UpdateState(ctx, job.ID, domain.JobRunning, domain.JobMuxing)
	if err != nil {
		w.logger.Warn().Err(err).Str("job_id", job.ID).Msg("failed to mark job muxing")
		return
	}
	PublishJobEvent(w.bus, "job.muxing", phase)

	finished, err := w.repo.UpdateState(ctx, job.ID, domain.JobMuxing, domain.JobCompleted)
	if err != nil {
		w.logger.Warn().Err(err).Str("job_id", job.ID).Msg("failed to mark job completed")
		return
	}
	finished, _ = w.repo.UpdateProgress(ctx, job.ID, 1)
	PublishJobEvent(w.bus, "job.completed", finished)
}
