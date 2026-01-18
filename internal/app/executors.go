package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/url"
	"time"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
)

type JobExecutor interface {
	Execute(ctx context.Context, job domain.Job, env ExecEnv) error
}

type ExecEnv struct {
	UpdateProgress func(progress float64) error
	IsCanceled     func() (bool, error)
	StepInterval   time.Duration
	Steps          int
}

type ExecutorRegistry struct {
	byType   map[string]JobExecutor
	fallback JobExecutor
}

func (r ExecutorRegistry) Get(jobType string) JobExecutor {
	if r.byType != nil {
		if ex, ok := r.byType[jobType]; ok {
			return ex
		}
	}
	return r.fallback
}

func DefaultExecutorRegistry() ExecutorRegistry {
	return ExecutorRegistry{
		byType: map[string]JobExecutor{
			"noop":     NoopExecutor{},
			"sleep":    SleepExecutor{},
			"download": DownloadStubExecutor{},
		},
		fallback: DefaultExecutor{},
	}
}

type NoopExecutor struct{}

func (NoopExecutor) Execute(ctx context.Context, job domain.Job, env ExecEnv) error {
	canceled, err := env.IsCanceled()
	if err != nil {
		return err
	}
	if canceled {
		return nil
	}
	return env.UpdateProgress(1)
}

type SleepExecutor struct{}

type sleepParams struct {
	Duration   string `json:"duration"`
	DurationMs int64  `json:"durationMs"`
	Seconds    int64  `json:"seconds"`
}

func (SleepExecutor) Execute(ctx context.Context, job domain.Job, env ExecEnv) error {
	dur := time.Second
	p := sleepParams{}
	if len(job.ParamsJSON) > 0 {
		_ = json.Unmarshal(job.ParamsJSON, &p)
	}
	if p.Duration != "" {
		if d, err := time.ParseDuration(p.Duration); err == nil {
			dur = d
		}
	} else if p.DurationMs > 0 {
		dur = time.Duration(p.DurationMs) * time.Millisecond
	} else if p.Seconds > 0 {
		dur = time.Duration(p.Seconds) * time.Second
	}
	if dur <= 0 {
		return env.UpdateProgress(1)
	}

	step := env.StepInterval
	if step <= 0 {
		step = 200 * time.Millisecond
	}

	start := time.Now()
	ticker := time.NewTicker(step)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			canceled, err := env.IsCanceled()
			if err != nil {
				return err
			}
			if canceled {
				return nil
			}

			elapsed := time.Since(start)
			progress := float64(elapsed) / float64(dur)
			progress = math.Max(0, math.Min(1, progress))
			if err := env.UpdateProgress(progress); err != nil {
				return err
			}
			if progress >= 1 {
				return nil
			}
		}
	}
}

type DefaultExecutor struct{}

func (DefaultExecutor) Execute(ctx context.Context, job domain.Job, env ExecEnv) error {
	steps := env.Steps
	if steps <= 0 {
		steps = DefaultWorkerOptions().Steps
	}
	step := env.StepInterval
	if step <= 0 {
		step = DefaultWorkerOptions().StepInterval
	}

	for i := 1; i <= steps; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(step):
		}

		canceled, err := env.IsCanceled()
		if err != nil {
			return err
		}
		if canceled {
			return nil
		}

		progress := float64(i) / float64(steps)
		progress = math.Max(0, math.Min(1, progress))
		if err := env.UpdateProgress(progress); err != nil {
			return err
		}
	}
	return nil
}

type DownloadStubExecutor struct{}

type downloadParams struct {
	URL string `json:"url"`
}

func (DownloadStubExecutor) Execute(ctx context.Context, job domain.Job, env ExecEnv) error {
	p := downloadParams{}
	if len(job.ParamsJSON) > 0 {
		_ = json.Unmarshal(job.ParamsJSON, &p)
	}
	if p.URL == "" {
		return fmt.Errorf("missing params.url")
	}
	u, err := url.Parse(p.URL)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return fmt.Errorf("invalid params.url")
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("unsupported url scheme")
	}

	canceled, err := env.IsCanceled()
	if err != nil {
		return err
	}
	if canceled {
		return nil
	}

	// Mode "stub" : si StepInterval est à 0, on termine immédiatement (utile pour tests).
	if env.StepInterval <= 0 {
		return env.UpdateProgress(1)
	}

	steps := 3
	step := env.StepInterval
	if step > 200*time.Millisecond {
		step = 200 * time.Millisecond
	}

	for i := 1; i <= steps; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(step):
		}

		canceled, err := env.IsCanceled()
		if err != nil {
			return err
		}
		if canceled {
			return nil
		}

		progress := float64(i) / float64(steps)
		progress = math.Max(0, math.Min(1, progress))
		if err := env.UpdateProgress(progress); err != nil {
			return err
		}
	}
	return nil
}

var ErrExecutorFailed = errors.New("executor failed")
