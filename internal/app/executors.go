package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
)

type JobExecutor interface {
	Execute(ctx context.Context, job domain.Job, env ExecEnv) error
}

type ExecEnv struct {
	UpdateProgress func(progress float64) error
	UpdateResult   func(resultJSON []byte) error
	IsCanceled     func() (bool, error)
	StepInterval   time.Duration
	Steps          int
	Destination    string
	CreateJob      func(jobType string, paramsJSON []byte) (domain.Job, error)
	GetJob         func(jobID string) (domain.Job, error)
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
			"download": DownloadExecutor{},
			"spawn":    SpawnExecutor{},
			"wait":     WaitExecutor{},
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

type DownloadExecutor struct{}

type downloadParams struct {
	URL      string `json:"url"`
	Filename string `json:"filename,omitempty"`
	Path     string `json:"path,omitempty"`
}

type downloadResult struct {
	URL         string `json:"url"`
	Path        string `json:"path"`
	Bytes       int64  `json:"bytes"`
	ContentType string `json:"contentType,omitempty"`
}

func (DownloadExecutor) Execute(ctx context.Context, job domain.Job, env ExecEnv) error {
	p := downloadParams{}
	if len(job.ParamsJSON) > 0 {
		_ = json.Unmarshal(job.ParamsJSON, &p)
	}
	if p.URL == "" {
		return &CodedError{Code: "invalid_params", Message: "missing params.url"}
	}
	u, err := url.Parse(p.URL)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return &CodedError{Code: "invalid_params", Message: "invalid params.url"}
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return &CodedError{Code: "invalid_params", Message: "unsupported url scheme"}
	}

	baseDir := strings.TrimSpace(env.Destination)
	if baseDir == "" {
		baseDir = domain.DefaultSettings().Destination
		if baseDir == "" {
			baseDir = "videos"
		}
	}

	filename := strings.TrimSpace(p.Filename)
	if filename == "" {
		filename = path.Base(u.Path)
		if filename == "/" || filename == "." || filename == "" {
			filename = job.ID
		}
	}
	filename = sanitizeFilename(filename)
	if filename == "" {
		filename = job.ID
	}

	var dstPath string
	if strings.TrimSpace(p.Path) != "" {
		dstPath, err = safeJoin(baseDir, p.Path)
		if err != nil {
			return &CodedError{Code: "invalid_params", Message: err.Error()}
		}
		// Si path ressemble à un dossier, ajouter le filename.
		if strings.HasSuffix(p.Path, "/") || strings.HasSuffix(p.Path, string(os.PathSeparator)) {
			dstPath = filepath.Join(dstPath, filename)
		} else {
			// Si path pointe vers un fichier sans extension, on garde tel quel.
			if fi := filepath.Base(dstPath); fi == "." || fi == string(os.PathSeparator) {
				dstPath = filepath.Join(dstPath, filename)
			}
		}
	} else {
		dstPath = filepath.Join(baseDir, filename)
	}

	canceled, err := env.IsCanceled()
	if err != nil {
		return err
	}
	if canceled {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(dstPath), 0o755); err != nil {
		return &CodedError{Code: "io_error", Message: "failed to create destination directory", Err: err}
	}

	tmpPath := dstPath + ".part"
	out, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return &CodedError{Code: "io_error", Message: "failed to create temp file", Err: err}
	}
	defer func() {
		_ = out.Close()
	}()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.URL, nil)
	if err != nil {
		return &CodedError{Code: "invalid_params", Message: "failed to build http request", Err: err}
	}
	req.Header.Set("User-Agent", "asd-server")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		_ = os.Remove(tmpPath)
		return &CodedError{Code: "network_error", Message: "http request failed", Err: err}
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		_ = os.Remove(tmpPath)
		return &CodedError{Code: "http_status", Message: fmt.Sprintf("http error: %s", resp.Status)}
	}

	total := resp.ContentLength
	buf := make([]byte, 128*1024)
	var downloaded int64
	lastUpdate := time.Now()

	for {
		canceled, err := env.IsCanceled()
		if err != nil {
			_ = os.Remove(tmpPath)
			return err
		}
		if canceled {
			_ = os.Remove(tmpPath)
			return nil
		}

		n, rerr := resp.Body.Read(buf)
		if n > 0 {
			if _, werr := out.Write(buf[:n]); werr != nil {
				_ = os.Remove(tmpPath)
				return &CodedError{Code: "io_error", Message: "failed to write temp file", Err: werr}
			}
			downloaded += int64(n)
		}

		now := time.Now()
		if total > 0 && now.Sub(lastUpdate) >= 250*time.Millisecond {
			progress := float64(downloaded) / float64(total)
			progress = math.Max(0, math.Min(0.999, progress))
			_ = env.UpdateProgress(progress)
			lastUpdate = now
		}

		if rerr != nil {
			if errors.Is(rerr, io.EOF) {
				break
			}
			_ = os.Remove(tmpPath)
			return &CodedError{Code: "io_error", Message: "failed while reading http response", Err: rerr}
		}
	}

	if err := out.Close(); err != nil {
		_ = os.Remove(tmpPath)
		return &CodedError{Code: "io_error", Message: "failed to close temp file", Err: err}
	}
	if err := os.Rename(tmpPath, dstPath); err != nil {
		_ = os.Remove(tmpPath)
		return &CodedError{Code: "io_error", Message: "failed to move temp file into place", Err: err}
	}

	// Résultat.
	if env.UpdateResult != nil {
		res := downloadResult{
			URL:         p.URL,
			Path:        dstPath,
			Bytes:       downloaded,
			ContentType: resp.Header.Get("Content-Type"),
		}
		if b, err := json.Marshal(res); err == nil {
			_ = env.UpdateResult(b)
		}
	}
	_ = env.UpdateProgress(1)
	return nil
}

func sanitizeFilename(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ReplaceAll(name, "\\", "_")
	name = strings.ReplaceAll(name, "/", "_")
	name = strings.ReplaceAll(name, "\x00", "")
	return name
}

func safeJoin(baseDir, rel string) (string, error) {
	if filepath.IsAbs(rel) {
		return "", fmt.Errorf("params.path must be relative")
	}
	clean := filepath.Clean(rel)
	if clean == "." {
		return baseDir, nil
	}
	if clean == ".." || strings.HasPrefix(clean, ".."+string(os.PathSeparator)) {
		return "", fmt.Errorf("invalid params.path")
	}
	out := filepath.Join(baseDir, clean)
	relToBase, err := filepath.Rel(baseDir, out)
	if err != nil {
		return "", fmt.Errorf("invalid params.path")
	}
	if relToBase == ".." || strings.HasPrefix(relToBase, ".."+string(os.PathSeparator)) {
		return "", fmt.Errorf("invalid params.path")
	}
	return out, nil
}

var ErrExecutorFailed = errors.New("executor failed")

type SpawnExecutor struct{}

type spawnJobSpec struct {
	Type   string          `json:"type"`
	Params json.RawMessage `json:"params,omitempty"`
}

type spawnParams struct {
	Jobs []spawnJobSpec `json:"jobs"`
}

type spawnResult struct {
	JobIDs []string `json:"jobIds"`
}

func (SpawnExecutor) Execute(ctx context.Context, job domain.Job, env ExecEnv) error {
	if env.CreateJob == nil {
		return &CodedError{Code: "executor_error", Message: "missing env.CreateJob"}
	}

	p := spawnParams{}
	if len(job.ParamsJSON) > 0 {
		_ = json.Unmarshal(job.ParamsJSON, &p)
	}
	if len(p.Jobs) == 0 {
		return &CodedError{Code: "invalid_params", Message: "missing params.jobs"}
	}

	ids := make([]string, 0, len(p.Jobs))
	for _, spec := range p.Jobs {
		if strings.TrimSpace(spec.Type) == "" {
			return &CodedError{Code: "invalid_params", Message: "spawn job missing type"}
		}

		canceled, err := env.IsCanceled()
		if err != nil {
			return err
		}
		if canceled {
			return nil
		}

		created, err := env.CreateJob(spec.Type, []byte(spec.Params))
		if err != nil {
			return &CodedError{Code: "executor_error", Message: "failed to create child job", Err: err}
		}
		ids = append(ids, created.ID)
	}

	if env.UpdateResult != nil {
		res := spawnResult{JobIDs: ids}
		if b, err := json.Marshal(res); err == nil {
			_ = env.UpdateResult(b)
		}
	}
	_ = env.UpdateProgress(1)
	return nil
}

type WaitExecutor struct{}

type waitParams struct {
	JobIDs       []string `json:"jobIds"`
	FailOnFailed *bool    `json:"failOnFailed,omitempty"`
	TimeoutMs    int64    `json:"timeoutMs,omitempty"`
	PollMs       int64    `json:"pollMs,omitempty"`
}

type waitChildSummary struct {
	ID        string          `json:"id"`
	State     domain.JobState `json:"state"`
	Progress  float64         `json:"progress"`
	ErrorCode string          `json:"errorCode,omitempty"`
	Error     string          `json:"error,omitempty"`
}

type waitResult struct {
	JobIDs      []string           `json:"jobIds"`
	Total       int                `json:"total"`
	Done        int                `json:"done"`
	Children    []waitChildSummary `json:"children"`
	CompletedAt time.Time          `json:"completedAt"`
}

func (WaitExecutor) Execute(ctx context.Context, job domain.Job, env ExecEnv) error {
	if env.GetJob == nil {
		return &CodedError{Code: "executor_error", Message: "missing env.GetJob"}
	}

	p := waitParams{}
	if len(job.ParamsJSON) > 0 {
		_ = json.Unmarshal(job.ParamsJSON, &p)
	}
	if len(p.JobIDs) == 0 {
		return &CodedError{Code: "invalid_params", Message: "missing params.jobIds"}
	}

	ids := make([]string, 0, len(p.JobIDs))
	seen := map[string]struct{}{}
	for _, id := range p.JobIDs {
		id = strings.TrimSpace(id)
		if id == "" {
			return &CodedError{Code: "invalid_params", Message: "jobIds must be non-empty"}
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return &CodedError{Code: "invalid_params", Message: "missing params.jobIds"}
	}

	failOnFailed := true
	if p.FailOnFailed != nil {
		failOnFailed = *p.FailOnFailed
	}

	poll := 300 * time.Millisecond
	if p.PollMs > 0 {
		poll = time.Duration(p.PollMs) * time.Millisecond
	}
	if poll <= 0 {
		poll = 300 * time.Millisecond
	}

	var deadline time.Time
	if p.TimeoutMs > 0 {
		deadline = time.Now().Add(time.Duration(p.TimeoutMs) * time.Millisecond)
	}

	ticker := time.NewTicker(poll)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}

		canceled, err := env.IsCanceled()
		if err != nil {
			return err
		}
		if canceled {
			return nil
		}

		if !deadline.IsZero() && time.Now().After(deadline) {
			return &CodedError{Code: "timeout", Message: "wait timeout"}
		}

		summaries := make([]waitChildSummary, 0, len(ids))
		done := 0
		for _, id := range ids {
			child, err := env.GetJob(id)
			if err != nil {
				if errors.Is(err, ErrNotFound) {
					return &CodedError{Code: "not_found", Message: fmt.Sprintf("child job not found: %s", id)}
				}
				return err
			}
			s := waitChildSummary{
				ID:        child.ID,
				State:     child.State,
				Progress:  child.Progress,
				ErrorCode: child.ErrorCode,
				Error:     child.ErrorMessage,
			}
			summaries = append(summaries, s)

			if child.State.IsTerminal() {
				done++
				if failOnFailed && (child.State == domain.JobFailed || child.State == domain.JobCanceled) {
					return &CodedError{Code: "child_failed", Message: fmt.Sprintf("child job failed: %s", child.ID)}
				}
			}
		}

		if env.UpdateProgress != nil {
			progress := float64(done) / float64(len(ids))
			progress = math.Max(0, math.Min(0.999, progress))
			_ = env.UpdateProgress(progress)
		}

		if done == len(ids) {
			if env.UpdateResult != nil {
				res := waitResult{
					JobIDs:      ids,
					Total:       len(ids),
					Done:        done,
					Children:    summaries,
					CompletedAt: time.Now().UTC(),
				}
				if b, err := json.Marshal(res); err == nil {
					_ = env.UpdateResult(b)
				}
			}
			_ = env.UpdateProgress(1)
			return nil
		}
	}
}
