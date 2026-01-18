package sqlite

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/ports"
)

func TestJobsRepository_ClaimNextQueued(t *testing.T) {
	ctx := context.Background()
	db, err := Open(ctx, ":memory:")
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	repo := NewJobsRepository(db.SQL)

	// Aucun job -> not found
	if _, err := repo.ClaimNextQueued(ctx); err == nil || !errors.Is(err, ports.ErrNotFound) {
		t.Fatalf("expected ErrNotFound when no queued jobs, got %v", err)
	}

	now := time.Now().UTC()
	_, err = repo.Create(ctx, domain.Job{
		ID:        "job1",
		Type:      "noop",
		State:     domain.JobQueued,
		Progress:  0,
		CreatedAt: now.Add(-2 * time.Minute),
		UpdatedAt: now.Add(-2 * time.Minute),
	})
	if err != nil {
		t.Fatalf("Create(job1): %v", err)
	}
	_, err = repo.Create(ctx, domain.Job{
		ID:        "job2",
		Type:      "noop",
		State:     domain.JobQueued,
		Progress:  0,
		CreatedAt: now.Add(-1 * time.Minute),
		UpdatedAt: now.Add(-1 * time.Minute),
	})
	if err != nil {
		t.Fatalf("Create(job2): %v", err)
	}

	claimed, err := repo.ClaimNextQueued(ctx)
	if err != nil {
		t.Fatalf("ClaimNextQueued: %v", err)
	}
	if claimed.ID != "job1" {
		t.Fatalf("expected to claim oldest (job1), got %q", claimed.ID)
	}
	if claimed.State != domain.JobRunning {
		t.Fatalf("expected claimed state running, got %q", claimed.State)
	}

	updated, err := repo.UpdateProgress(ctx, claimed.ID, 0.5)
	if err != nil {
		t.Fatalf("UpdateProgress: %v", err)
	}
	if updated.Progress != 0.5 {
		t.Fatalf("expected progress=0.5, got %v", updated.Progress)
	}
}
