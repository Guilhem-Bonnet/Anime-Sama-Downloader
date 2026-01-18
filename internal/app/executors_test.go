package app

import (
	"context"
	"testing"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
)

func TestDownloadStubExecutor_InvalidURL(t *testing.T) {
	ex := DownloadStubExecutor{}

	job := domain.Job{
		ID:         "job1",
		Type:       "download",
		ParamsJSON: []byte(`{"url":"not a url"}`),
	}

	var last float64
	err := ex.Execute(context.Background(), job, ExecEnv{
		UpdateProgress: func(p float64) error { last = p; return nil },
		IsCanceled:     func() (bool, error) { return false, nil },
		StepInterval:   0,
		Steps:          0,
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if last != 0 {
		t.Fatalf("expected no progress updates, got %v", last)
	}
}

func TestDownloadStubExecutor_ImmediateSuccess(t *testing.T) {
	ex := DownloadStubExecutor{}

	job := domain.Job{
		ID:         "job2",
		Type:       "download",
		ParamsJSON: []byte(`{"url":"https://example.com/video.m3u8"}`),
	}

	var last float64
	err := ex.Execute(context.Background(), job, ExecEnv{
		UpdateProgress: func(p float64) error { last = p; return nil },
		IsCanceled:     func() (bool, error) { return false, nil },
		StepInterval:   0,
		Steps:          0,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if last != 1 {
		t.Fatalf("expected progress 1, got %v", last)
	}
}
