package app

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
)

func TestDownloadExecutor_InvalidURL(t *testing.T) {
	ex := DownloadExecutor{}

	job := domain.Job{
		ID:         "job1",
		Type:       "download",
		ParamsJSON: []byte(`{"url":"not a url"}`),
	}

	var last float64
	err := ex.Execute(context.Background(), job, ExecEnv{
		UpdateProgress: func(p float64) error { last = p; return nil },
		UpdateResult:   func([]byte) error { return nil },
		IsCanceled:     func() (bool, error) { return false, nil },
		Destination:    t.TempDir(),
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	var coded *CodedError
	if !errors.As(err, &coded) || coded.Code != "invalid_params" {
		t.Fatalf("expected invalid_params coded error, got %T (%v)", err, err)
	}
	if last != 0 {
		t.Fatalf("expected no progress updates, got %v", last)
	}
}

func TestDownloadExecutor_HTTPStatusErrorIsCoded(t *testing.T) {
	ex := DownloadExecutor{}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "nope", http.StatusNotFound)
	}))
	defer srv.Close()

	job := domain.Job{
		ID:         "job404",
		Type:       "download",
		ParamsJSON: []byte(`{"url":"` + srv.URL + `"}`),
	}

	err := ex.Execute(context.Background(), job, ExecEnv{
		UpdateProgress: func(float64) error { return nil },
		UpdateResult:   func([]byte) error { return nil },
		IsCanceled:     func() (bool, error) { return false, nil },
		Destination:    t.TempDir(),
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	var coded *CodedError
	if !errors.As(err, &coded) || coded.Code != "http_status" {
		t.Fatalf("expected http_status coded error, got %T (%v)", err, err)
	}
}

func TestDownloadExecutor_HTMLResponseIsNotMedia(t *testing.T) {
	ex := DownloadExecutor{}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte("<html><body>embed</body></html>"))
	}))
	defer srv.Close()

	job := domain.Job{
		ID:         "jobhtml",
		Type:       "download",
		ParamsJSON: []byte(`{"url":"` + srv.URL + `"}`),
	}

	err := ex.Execute(context.Background(), job, ExecEnv{
		UpdateProgress: func(float64) error { return nil },
		UpdateResult:   func([]byte) error { return nil },
		IsCanceled:     func() (bool, error) { return false, nil },
		Destination:    t.TempDir(),
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	var coded *CodedError
	if !errors.As(err, &coded) || coded.Code != "not_media" {
		t.Fatalf("expected not_media coded error, got %T (%v)", err, err)
	}
}

func TestDownloadExecutor_DownloadsToFileAndSetsResult(t *testing.T) {
	ex := DownloadExecutor{}

	payload := []byte("hello world")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Length", strconv.Itoa(len(payload)))
		_, _ = w.Write(payload)
	}))
	defer srv.Close()

	dest := t.TempDir()

	job := domain.Job{
		ID:         "job2",
		Type:       "download",
		ParamsJSON: []byte(`{"url":"` + srv.URL + `","filename":"file.bin"}`),
	}

	var last float64
	var resultJSON []byte
	err := ex.Execute(context.Background(), job, ExecEnv{
		UpdateProgress: func(p float64) error { last = p; return nil },
		UpdateResult:   func(b []byte) error { resultJSON = append([]byte(nil), b...); return nil },
		IsCanceled:     func() (bool, error) { return false, nil },
		Destination:    dest,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if last != 1 {
		t.Fatalf("expected progress 1, got %v", last)
	}
	if len(resultJSON) == 0 {
		t.Fatalf("expected result JSON to be set")
	}
	var res downloadResult
	if err := json.Unmarshal(resultJSON, &res); err != nil {
		t.Fatalf("invalid result JSON: %v", err)
	}
	if res.Path == "" {
		t.Fatalf("expected result path")
	}
	if res.Bytes != int64(len(payload)) {
		t.Fatalf("expected %d bytes, got %d", len(payload), res.Bytes)
	}
	if _, err := os.Stat(res.Path); err != nil {
		t.Fatalf("expected file to exist: %v", err)
	}
	if filepath.Dir(res.Path) != dest {
		t.Fatalf("expected file to be under destination")
	}
	data, err := os.ReadFile(res.Path)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	if string(data) != string(payload) {
		t.Fatalf("unexpected file content")
	}
}

func TestDownloadExecutor_PathTraversalRejected(t *testing.T) {
	ex := DownloadExecutor{}

	job := domain.Job{
		ID:         "job3",
		Type:       "download",
		ParamsJSON: []byte(`{"url":"https://example.com","path":"../escape"}`),
	}

	err := ex.Execute(context.Background(), job, ExecEnv{
		UpdateProgress: func(float64) error { return nil },
		UpdateResult:   func([]byte) error { return nil },
		IsCanceled:     func() (bool, error) { return false, nil },
		Destination:    t.TempDir(),
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	var coded *CodedError
	if !errors.As(err, &coded) || coded.Code != "invalid_params" {
		t.Fatalf("expected invalid_params coded error, got %T (%v)", err, err)
	}
}

func TestSpawnExecutor_MissingJobs(t *testing.T) {
	ex := SpawnExecutor{}

	job := domain.Job{ID: "spawn1", Type: "spawn", ParamsJSON: []byte(`{"jobs":[]}`)}
	err := ex.Execute(context.Background(), job, ExecEnv{
		UpdateProgress: func(float64) error { return nil },
		UpdateResult:   func([]byte) error { return nil },
		IsCanceled:     func() (bool, error) { return false, nil },
		CreateJob:      func(string, []byte) (domain.Job, error) { return domain.Job{}, nil },
	})
	if err == nil {
		t.Fatalf("expected error")
	}
	var coded *CodedError
	if !errors.As(err, &coded) || coded.Code != "invalid_params" {
		t.Fatalf("expected invalid_params coded error, got %T (%v)", err, err)
	}
}

func TestSpawnExecutor_CreatesJobsAndSetsResult(t *testing.T) {
	ex := SpawnExecutor{}

	job := domain.Job{ID: "spawn2", Type: "spawn", ParamsJSON: []byte(`{"jobs":[{"type":"noop"},{"type":"sleep","params":{"seconds":1}}]}`)}

	created := make([]domain.Job, 0, 2)
	var resultJSON []byte
	var last float64

	err := ex.Execute(context.Background(), job, ExecEnv{
		UpdateProgress: func(p float64) error { last = p; return nil },
		UpdateResult:   func(b []byte) error { resultJSON = append([]byte(nil), b...); return nil },
		IsCanceled:     func() (bool, error) { return false, nil },
		CreateJob: func(jobType string, paramsJSON []byte) (domain.Job, error) {
			j := domain.Job{ID: jobType + "-id", Type: jobType, ParamsJSON: paramsJSON, CreatedAt: time.Now(), UpdatedAt: time.Now()}
			created = append(created, j)
			return j, nil
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if last != 1 {
		t.Fatalf("expected progress 1, got %v", last)
	}
	if len(created) != 2 {
		t.Fatalf("expected 2 created jobs, got %d", len(created))
	}
	if len(resultJSON) == 0 {
		t.Fatalf("expected result JSON")
	}
	var res struct {
		JobIDs []string `json:"jobIds"`
	}
	if err := json.Unmarshal(resultJSON, &res); err != nil {
		t.Fatalf("invalid result JSON: %v", err)
	}
	if len(res.JobIDs) != 2 {
		t.Fatalf("expected 2 result jobIds, got %d", len(res.JobIDs))
	}
}

func TestWaitExecutor_MissingJobIDs(t *testing.T) {
	ex := WaitExecutor{}
	job := domain.Job{ID: "wait1", Type: "wait", ParamsJSON: []byte(`{"jobIds":[]}`)}

	err := ex.Execute(context.Background(), job, ExecEnv{
		UpdateProgress: func(float64) error { return nil },
		UpdateResult:   func([]byte) error { return nil },
		IsCanceled:     func() (bool, error) { return false, nil },
		GetJob:         func(string) (domain.Job, error) { return domain.Job{}, nil },
	})
	if err == nil {
		t.Fatalf("expected error")
	}
	var coded *CodedError
	if !errors.As(err, &coded) || coded.Code != "invalid_params" {
		t.Fatalf("expected invalid_params coded error, got %T (%v)", err, err)
	}
}

func TestWaitExecutor_WaitsAndSetsResult(t *testing.T) {
	ex := WaitExecutor{}
	job := domain.Job{ID: "wait2", Type: "wait", ParamsJSON: []byte(`{"jobIds":["a","b"],"pollMs":1,"timeoutMs":200}`)}

	var calls int
	var last float64
	var resultJSON []byte

	err := ex.Execute(context.Background(), job, ExecEnv{
		UpdateProgress: func(p float64) error { last = p; return nil },
		UpdateResult:   func(b []byte) error { resultJSON = append([]byte(nil), b...); return nil },
		IsCanceled:     func() (bool, error) { return false, nil },
		GetJob: func(id string) (domain.Job, error) {
			calls++
			state := domain.JobRunning
			if calls >= 5 {
				state = domain.JobCompleted
			}
			return domain.Job{ID: id, Type: "noop", State: state, Progress: 1, UpdatedAt: time.Now()}, nil
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if last != 1 {
		t.Fatalf("expected progress 1, got %v", last)
	}
	if len(resultJSON) == 0 {
		t.Fatalf("expected result JSON")
	}
	var res struct {
		Total    int `json:"total"`
		Done     int `json:"done"`
		Children []struct {
			ID    string `json:"id"`
			State string `json:"state"`
		} `json:"children"`
	}
	if err := json.Unmarshal(resultJSON, &res); err != nil {
		t.Fatalf("invalid result JSON: %v", err)
	}
	if res.Total != 2 || res.Done != 2 {
		t.Fatalf("expected total=2 done=2, got total=%d done=%d", res.Total, res.Done)
	}
	if len(res.Children) != 2 {
		t.Fatalf("expected 2 children, got %d", len(res.Children))
	}
}

func TestWaitExecutor_FailsOnChildFailedByDefault(t *testing.T) {
	ex := WaitExecutor{}
	job := domain.Job{ID: "wait3", Type: "wait", ParamsJSON: []byte(`{"jobIds":["ok","bad"],"pollMs":1,"timeoutMs":200}`)}

	err := ex.Execute(context.Background(), job, ExecEnv{
		UpdateProgress: func(float64) error { return nil },
		UpdateResult:   func([]byte) error { return nil },
		IsCanceled:     func() (bool, error) { return false, nil },
		GetJob: func(id string) (domain.Job, error) {
			state := domain.JobCompleted
			if id == "bad" {
				state = domain.JobFailed
			}
			return domain.Job{ID: id, Type: "noop", State: state, Progress: 1, UpdatedAt: time.Now()}, nil
		},
	})
	if err == nil {
		t.Fatalf("expected error")
	}
	var coded *CodedError
	if !errors.As(err, &coded) || coded.Code != "child_failed" {
		t.Fatalf("expected child_failed coded error, got %T (%v)", err, err)
	}
}

// TestNoopExecutor_Execute_Success tests NoopExecutor completes successfully
func TestNoopExecutor_Execute_Success(t *testing.T) {
	ex := NoopExecutor{}
	var progress []float64

	err := ex.Execute(context.Background(), domain.Job{}, ExecEnv{
		UpdateProgress: func(p float64) error {
			progress = append(progress, p)
			return nil
		},
		UpdateResult: func([]byte) error { return nil },
		IsCanceled:   func() (bool, error) { return false, nil },
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(progress) != 1 || progress[0] != 1.0 {
		t.Fatalf("expected progress [1.0], got %v", progress)
	}
}

// TestNoopExecutor_Execute_Canceled tests NoopExecutor respects cancellation
func TestNoopExecutor_Execute_Canceled(t *testing.T) {
	ex := NoopExecutor{}
	var progress []float64

	err := ex.Execute(context.Background(), domain.Job{}, ExecEnv{
		UpdateProgress: func(p float64) error {
			progress = append(progress, p)
			return nil
		},
		UpdateResult: func([]byte) error { return nil },
		IsCanceled:   func() (bool, error) { return true, nil },
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(progress) > 0 {
		t.Fatalf("expected no progress when canceled, got %v", progress)
	}
}

// TestNoopExecutor_Execute_CancelError tests NoopExecutor propagates cancel check errors
func TestNoopExecutor_Execute_CancelError(t *testing.T) {
	ex := NoopExecutor{}

	err := ex.Execute(context.Background(), domain.Job{}, ExecEnv{
		UpdateProgress: func(p float64) error { return nil },
		UpdateResult:   func([]byte) error { return nil },
		IsCanceled:     func() (bool, error) { return false, errors.New("cancel failed") },
	})

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

// TestSleepExecutor_Execute_ZeroDuration tests SleepExecutor with zero duration (defaults to 1 second)
func TestSleepExecutor_Execute_ZeroDuration(t *testing.T) {
	ex := SleepExecutor{}
	var progress []float64

	job := domain.Job{ParamsJSON: []byte(`{"durationMs": 0}`)}
	start := time.Now()
	err := ex.Execute(context.Background(), job, ExecEnv{
		UpdateProgress: func(p float64) error {
			progress = append(progress, p)
			return nil
		},
		UpdateResult: func([]byte) error { return nil },
		IsCanceled:   func() (bool, error) { return false, nil },
		StepInterval: 200 * time.Millisecond,
	})
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	// When durationMs is 0, defaults to 1 second
	if elapsed < 900*time.Millisecond {
		t.Fatalf("expected at least 900ms, got %v", elapsed)
	}
	// Should have multiple updates
	if len(progress) < 2 {
		t.Fatalf("expected multiple progress updates for 1 second sleep, got %d", len(progress))
	}
	// Final progress should be 1.0
	if progress[len(progress)-1] != 1.0 {
		t.Fatalf("expected final progress 1.0, got %v", progress[len(progress)-1])
	}
}

// TestSleepExecutor_Execute_WithDuration tests SleepExecutor sleeps and updates progress
func TestSleepExecutor_Execute_WithDuration(t *testing.T) {
	ex := SleepExecutor{}
	var progress []float64

	job := domain.Job{ParamsJSON: []byte(`{"durationMs": 50}`)}
	start := time.Now()
	err := ex.Execute(context.Background(), job, ExecEnv{
		UpdateProgress: func(p float64) error {
			progress = append(progress, p)
			return nil
		},
		UpdateResult: func([]byte) error { return nil },
		IsCanceled:   func() (bool, error) { return false, nil },
		StepInterval: 10 * time.Millisecond,
	})
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if elapsed < 50*time.Millisecond {
		t.Fatalf("expected at least 50ms, got %v", elapsed)
	}
	if len(progress) < 2 {
		t.Fatalf("expected multiple progress updates, got %d", len(progress))
	}
	if progress[len(progress)-1] != 1.0 {
		t.Fatalf("expected final progress 1.0, got %v", progress[len(progress)-1])
	}
}

// TestSleepExecutor_Execute_Canceled tests SleepExecutor respects cancellation
func TestSleepExecutor_Execute_Canceled(t *testing.T) {
	ex := SleepExecutor{}
	cancelAt := 0
	callCount := 0

	job := domain.Job{ParamsJSON: []byte(`{"durationMs": 500}`)}
	start := time.Now()
	err := ex.Execute(context.Background(), job, ExecEnv{
		UpdateProgress: func(p float64) error { return nil },
		UpdateResult:   func([]byte) error { return nil },
		IsCanceled: func() (bool, error) {
			callCount++
			return callCount > cancelAt, nil
		},
		StepInterval: 10 * time.Millisecond,
	})
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if elapsed > 200*time.Millisecond {
		t.Fatalf("expected quick return due to cancellation, got %v", elapsed)
	}
}

// TestSleepExecutor_Execute_DurationString tests duration string parameter
func TestSleepExecutor_Execute_DurationString(t *testing.T) {
	ex := SleepExecutor{}
	var progress []float64

	job := domain.Job{ParamsJSON: []byte(`{"duration": "30ms"}`)}
	start := time.Now()
	err := ex.Execute(context.Background(), job, ExecEnv{
		UpdateProgress: func(p float64) error {
			progress = append(progress, p)
			return nil
		},
		UpdateResult: func([]byte) error { return nil },
		IsCanceled:   func() (bool, error) { return false, nil },
		StepInterval: 5 * time.Millisecond,
	})
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if elapsed < 30*time.Millisecond {
		t.Fatalf("expected at least 30ms, got %v", elapsed)
	}
	if len(progress) < 2 {
		t.Fatalf("expected multiple progress updates, got %d", len(progress))
	}
}

// TestSleepExecutor_Execute_SecondsParam tests seconds parameter
func TestSleepExecutor_Execute_SecondsParam(t *testing.T) {
	ex := SleepExecutor{}
	var progress []float64

	job := domain.Job{ParamsJSON: []byte(`{"seconds": 0.05}`)}
	start := time.Now()
	err := ex.Execute(context.Background(), job, ExecEnv{
		UpdateProgress: func(p float64) error {
			progress = append(progress, p)
			return nil
		},
		UpdateResult: func([]byte) error { return nil },
		IsCanceled:   func() (bool, error) { return false, nil },
		StepInterval: 5 * time.Millisecond,
	})
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if elapsed < 45*time.Millisecond {
		t.Fatalf("expected at least 45ms, got %v", elapsed)
	}
}

// TestDefaultExecutor_Execute_Success tests DefaultExecutor completes successfully
func TestDefaultExecutor_Execute_Success(t *testing.T) {
	ex := DefaultExecutor{}
	var progress []float64

	err := ex.Execute(context.Background(), domain.Job{}, ExecEnv{
		UpdateProgress: func(p float64) error {
			progress = append(progress, p)
			return nil
		},
		UpdateResult: func([]byte) error { return nil },
		IsCanceled:   func() (bool, error) { return false, nil },
		Steps:        5,
		StepInterval: 5 * time.Millisecond,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(progress) < 2 {
		t.Fatalf("expected multiple progress updates, got %d", len(progress))
	}
	if progress[len(progress)-1] != 1.0 {
		t.Fatalf("expected final progress 1.0, got %v", progress[len(progress)-1])
	}
}

// TestDefaultExecutor_Execute_Canceled tests DefaultExecutor respects cancellation
func TestDefaultExecutor_Execute_Canceled(t *testing.T) {
	ex := DefaultExecutor{}
	callCount := 0

	err := ex.Execute(context.Background(), domain.Job{}, ExecEnv{
		UpdateProgress: func(p float64) error { return nil },
		UpdateResult:   func([]byte) error { return nil },
		IsCanceled: func() (bool, error) {
			callCount++
			return callCount > 1, nil
		},
		Steps:        100,
		StepInterval: 5 * time.Millisecond,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

// TestDefaultExecutor_Execute_DefaultSteps tests DefaultExecutor uses default steps
func TestDefaultExecutor_Execute_DefaultSteps(t *testing.T) {
	ex := DefaultExecutor{}
	var progress []float64

	err := ex.Execute(context.Background(), domain.Job{}, ExecEnv{
		UpdateProgress: func(p float64) error {
			progress = append(progress, p)
			return nil
		},
		UpdateResult: func([]byte) error { return nil },
		IsCanceled:   func() (bool, error) { return false, nil },
		// No Steps or StepInterval set - should use defaults
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(progress) == 0 {
		t.Fatalf("expected progress updates with defaults")
	}
}

// TestExecutorRegistry_Get_TypeExists tests registry returns correct executor type
func TestExecutorRegistry_Get_TypeExists(t *testing.T) {
	reg := DefaultExecutorRegistry()

	ex := reg.Get("noop")
	if _, ok := ex.(NoopExecutor); !ok {
		t.Fatalf("expected NoopExecutor, got %T", ex)
	}

	ex = reg.Get("sleep")
	if _, ok := ex.(SleepExecutor); !ok {
		t.Fatalf("expected SleepExecutor, got %T", ex)
	}
}

// TestExecutorRegistry_Get_TypeNotFound tests registry returns fallback executor
func TestExecutorRegistry_Get_TypeNotFound(t *testing.T) {
	reg := DefaultExecutorRegistry()

	ex := reg.Get("unknown-executor-type")
	if _, ok := ex.(DefaultExecutor); !ok {
		t.Fatalf("expected DefaultExecutor fallback, got %T", ex)
	}
}

// TestExecutorRegistry_Get_AllBuiltInTypes tests all built-in executor types exist
func TestExecutorRegistry_Get_AllBuiltInTypes(t *testing.T) {
	reg := DefaultExecutorRegistry()

	types := []string{"noop", "sleep", "download", "spawn", "wait"}
	for _, jobType := range types {
		ex := reg.Get(jobType)
		if ex == nil {
			t.Fatalf("expected executor for type %q, got nil", jobType)
		}
	}
}

// TestExecutorRegistry_Get_EmptyRegistry tests Get with empty registry
func TestExecutorRegistry_Get_EmptyRegistry(t *testing.T) {
	reg := ExecutorRegistry{} // Empty - no types, no fallback

	ex := reg.Get("anything")
	if ex != nil {
		t.Fatalf("expected nil from empty registry, got %v", ex)
	}
}

// TestExecutorRegistry_Get_NilMapFallback tests registry with nil byType map uses fallback
func TestExecutorRegistry_Get_NilMapFallback(t *testing.T) {
	reg := ExecutorRegistry{
		byType:   nil,
		fallback: DefaultExecutor{},
	}

	ex := reg.Get("any-type")
	if _, ok := ex.(DefaultExecutor); !ok {
		t.Fatalf("expected DefaultExecutor fallback, got %T", ex)
	}
}

// TestSleepExecutor_Execute_ContextCanceled tests SleepExecutor respects context cancellation
func TestSleepExecutor_Execute_ContextCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	ex := SleepExecutor{}

	job := domain.Job{ParamsJSON: []byte(`{"durationMs": 500}`)}

	cancel()

	err := ex.Execute(ctx, job, ExecEnv{
		UpdateProgress: func(p float64) error { return nil },
		UpdateResult:   func([]byte) error { return nil },
		IsCanceled:     func() (bool, error) { return false, nil },
		StepInterval:   10 * time.Millisecond,
	})

	if err == nil {
		t.Fatalf("expected context canceled error, got nil")
	}
}

// TestDefaultExecutor_Execute_ContextCanceled tests DefaultExecutor respects context cancellation
func TestDefaultExecutor_Execute_ContextCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	ex := DefaultExecutor{}

	cancel()

	err := ex.Execute(ctx, domain.Job{}, ExecEnv{
		UpdateProgress: func(p float64) error { return nil },
		UpdateResult:   func([]byte) error { return nil },
		IsCanceled:     func() (bool, error) { return false, nil },
		Steps:          5,
	})

	if err == nil {
		t.Fatalf("expected context canceled error, got nil")
	}
}

// TestSleepExecutor_Execute_ProgressError tests SleepExecutor propagates progress update errors
func TestSleepExecutor_Execute_ProgressError(t *testing.T) {
	ex := SleepExecutor{}

	job := domain.Job{ParamsJSON: []byte(`{"durationMs": 50}`)}
	err := ex.Execute(context.Background(), job, ExecEnv{
		UpdateProgress: func(p float64) error { return errors.New("update failed") },
		UpdateResult:   func([]byte) error { return nil },
		IsCanceled:     func() (bool, error) { return false, nil },
		StepInterval:   10 * time.Millisecond,
	})

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

// TestDefaultExecutor_Execute_ProgressError tests DefaultExecutor propagates progress update errors
func TestDefaultExecutor_Execute_ProgressError(t *testing.T) {
	ex := DefaultExecutor{}

	err := ex.Execute(context.Background(), domain.Job{}, ExecEnv{
		UpdateProgress: func(p float64) error { return errors.New("update failed") },
		UpdateResult:   func([]byte) error { return nil },
		IsCanceled:     func() (bool, error) { return false, nil },
		Steps:          5,
	})

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

// TestDefaultExecutor_Execute_CancelError tests DefaultExecutor propagates cancel check errors
func TestDefaultExecutor_Execute_CancelError(t *testing.T) {
	ex := DefaultExecutor{}

	err := ex.Execute(context.Background(), domain.Job{}, ExecEnv{
		UpdateProgress: func(p float64) error { return nil },
		UpdateResult:   func([]byte) error { return nil },
		IsCanceled:     func() (bool, error) { return false, errors.New("cancel check failed") },
		Steps:          5,
	})

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

// TestSleepExecutor_Execute_InvalidDurationString tests invalid duration string
func TestSleepExecutor_Execute_InvalidDurationString(t *testing.T) {
	ex := SleepExecutor{}
	var progress []float64

	job := domain.Job{ParamsJSON: []byte(`{"duration": "invalid"}`)}
	start := time.Now()
	err := ex.Execute(context.Background(), job, ExecEnv{
		UpdateProgress: func(p float64) error {
			progress = append(progress, p)
			return nil
		},
		UpdateResult: func([]byte) error { return nil },
		IsCanceled:   func() (bool, error) { return false, nil },
		StepInterval: 50 * time.Millisecond,
	})
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	// Should fallback to default 1 second
	if elapsed < 900*time.Millisecond {
		t.Fatalf("expected at least 900ms (1 second default), got %v", elapsed)
	}
}

// TestSleepExecutor_Execute_EmptyParams tests empty params JSON
func TestSleepExecutor_Execute_EmptyParams(t *testing.T) {
	ex := SleepExecutor{}
	var progress []float64

	job := domain.Job{ParamsJSON: []byte{}}
	start := time.Now()
	err := ex.Execute(context.Background(), job, ExecEnv{
		UpdateProgress: func(p float64) error {
			progress = append(progress, p)
			return nil
		},
		UpdateResult: func([]byte) error { return nil },
		IsCanceled:   func() (bool, error) { return false, nil },
		StepInterval: 50 * time.Millisecond,
	})
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	// Should use default 1 second
	if elapsed < 900*time.Millisecond {
		t.Fatalf("expected at least 900ms, got %v", elapsed)
	}
}
