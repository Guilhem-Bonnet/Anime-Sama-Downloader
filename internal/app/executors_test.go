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
