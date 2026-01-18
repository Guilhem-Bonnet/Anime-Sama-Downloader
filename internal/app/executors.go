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
		return fmt.Errorf("missing params.url")
	}
	u, err := url.Parse(p.URL)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return fmt.Errorf("invalid params.url")
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("unsupported url scheme")
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
			return err
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
		return err
	}

	tmpPath := dstPath + ".part"
	out, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer func() {
		_ = out.Close()
	}()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.URL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "asd-server")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		_ = os.Remove(tmpPath)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("http error: %s", resp.Status)
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
				return werr
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
			return rerr
		}
	}

	if err := out.Close(); err != nil {
		_ = os.Remove(tmpPath)
		return err
	}
	if err := os.Rename(tmpPath, dstPath); err != nil {
		_ = os.Remove(tmpPath)
		return err
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
