package sqlite

import (
	"context"
	"testing"
	"time"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
)

// TestUpdateState_SetsStartedAtOnRunning tests that started_at is set when transitioning to running.
func TestUpdateState_SetsStartedAtOnRunning(t *testing.T) {
	ctx := context.Background()
	db, err := Open(ctx, ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	repo := NewJobsRepository(db.SQL)

	// Create a job in queued state
	job := domain.Job{
		ID:        "test-job-1",
		Type:      "download",
		State:     domain.JobQueued,
		Progress:  0.0,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	created, err := repo.Create(ctx, job)
	if err != nil {
		t.Fatalf("Failed to create job: %v", err)
	}

	// Verify started_at is initially nil
	if created.StartedAt != nil {
		t.Errorf("StartedAt should be nil for queued job, got %v", created.StartedAt)
	}

	// Transition to running
	updated, err := repo.UpdateState(ctx, job.ID, domain.JobQueued, domain.JobRunning)
	if err != nil {
		t.Fatalf("Failed to update state: %v", err)
	}

	// Verify started_at is now set
	if updated.StartedAt == nil {
		t.Error("StartedAt should be set when transitioning to running")
	}

	// Verify state changed
	if updated.State != domain.JobRunning {
		t.Errorf("Expected state to be running, got %v", updated.State)
	}
}

// TestUpdateState_SetsCompletedAtOnTerminal tests that completed_at is set on terminal states.
func TestUpdateState_SetsCompletedAtOnTerminal(t *testing.T) {
	ctx := context.Background()
	db, err := Open(ctx, ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	repo := NewJobsRepository(db.SQL)

	// Create and transition job to running
	job := domain.Job{
		ID:        "test-job-2",
		Type:      "download",
		State:     domain.JobQueued,
		Progress:  0.0,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	_, err = repo.Create(ctx, job)
	if err != nil {
		t.Fatalf("Failed to create job: %v", err)
	}

	running, err := repo.UpdateState(ctx, job.ID, domain.JobQueued, domain.JobRunning)
	if err != nil {
		t.Fatalf("Failed to transition to running: %v", err)
	}

	// Verify completed_at is initially nil
	if running.CompletedAt != nil {
		t.Errorf("CompletedAt should be nil for running job, got %v", running.CompletedAt)
	}

	// Transition to muxing
	muxing, err := repo.UpdateState(ctx, job.ID, domain.JobRunning, domain.JobMuxing)
	if err != nil {
		t.Fatalf("Failed to transition to muxing: %v", err)
	}

	// Verify completed_at is still nil (muxing is not terminal)
	if muxing.CompletedAt != nil {
		t.Errorf("CompletedAt should be nil for muxing job, got %v", muxing.CompletedAt)
	}

	// Transition to completed (terminal state)
	completed, err := repo.UpdateState(ctx, job.ID, domain.JobMuxing, domain.JobCompleted)
	if err != nil {
		t.Fatalf("Failed to transition to completed: %v", err)
	}

	// Verify completed_at is now set
	if completed.CompletedAt == nil {
		t.Error("CompletedAt should be set when transitioning to completed")
	}

	// Verify state changed
	if completed.State != domain.JobCompleted {
		t.Errorf("Expected state to be completed, got %v", completed.State)
	}
}

// TestUpdateState_SetsCompletedAtOnFailed tests that completed_at is set on failed state.
func TestUpdateState_SetsCompletedAtOnFailed(t *testing.T) {
	ctx := context.Background()
	db, err := Open(ctx, ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	repo := NewJobsRepository(db.SQL)

	job := domain.Job{
		ID:        "test-job-3",
		Type:      "download",
		State:     domain.JobQueued,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	_, err = repo.Create(ctx, job)
	if err != nil {
		t.Fatalf("Failed to create job: %v", err)
	}

	// Transition directly to failed
	failed, err := repo.UpdateState(ctx, job.ID, domain.JobQueued, domain.JobFailed)
	if err != nil {
		t.Fatalf("Failed to transition to failed: %v", err)
	}

	// Verify completed_at is set
	if failed.CompletedAt == nil {
		t.Error("CompletedAt should be set when transitioning to failed")
	}

	if failed.State != domain.JobFailed {
		t.Errorf("Expected state to be failed, got %v", failed.State)
	}
}

// TestLoadUnfinishedJobs tests loading of queued and running jobs.
func TestLoadUnfinishedJobs(t *testing.T) {
	ctx := context.Background()
	db, err := Open(ctx, ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	repo := NewJobsRepository(db.SQL)

	// Create jobs in various states
	now := time.Now().UTC()

	jobs := []domain.Job{
		{ID: "job1", Type: "download", State: domain.JobQueued, CreatedAt: now.Add(-3 * time.Minute), UpdatedAt: now.Add(-3 * time.Minute)},
		{ID: "job2", Type: "download", State: domain.JobRunning, CreatedAt: now.Add(-2 * time.Minute), UpdatedAt: now.Add(-2 * time.Minute)},
		{ID: "job3", Type: "download", State: domain.JobCompleted, CreatedAt: now.Add(-1 * time.Minute), UpdatedAt: now.Add(-1 * time.Minute)},
		{ID: "job4", Type: "download", State: domain.JobQueued, CreatedAt: now, UpdatedAt: now},
		{ID: "job5", Type: "download", State: domain.JobFailed, CreatedAt: now.Add(-5 * time.Minute), UpdatedAt: now.Add(-5 * time.Minute)},
	}

	for _, j := range jobs {
		_, err := repo.Create(ctx, j)
		if err != nil {
			t.Fatalf("Failed to create job %s: %v", j.ID, err)
		}
	}

	// Load unfinished jobs
	unfinished, err := repo.LoadUnfinishedJobs(ctx)
	if err != nil {
		t.Fatalf("Failed to load unfinished jobs: %v", err)
	}

	// Should only return queued and running jobs
	if len(unfinished) != 3 {
		t.Errorf("Expected 3 unfinished jobs, got %d", len(unfinished))
	}

	// Verify order (should be sorted by created_at ASC)
	expectedOrder := []string{"job1", "job2", "job4"}
	for i, expectedID := range expectedOrder {
		if unfinished[i].ID != expectedID {
			t.Errorf("Expected job at position %d to be %s, got %s", i, expectedID, unfinished[i].ID)
		}
	}

	// Verify states
	for _, j := range unfinished {
		if j.State != domain.JobQueued && j.State != domain.JobRunning {
			t.Errorf("Unfinished job %s has unexpected state: %v", j.ID, j.State)
		}
	}
}

// TestLoadUnfinishedJobs_Empty tests loading when no unfinished jobs exist.
func TestLoadUnfinishedJobs_Empty(t *testing.T) {
	ctx := context.Background()
	db, err := Open(ctx, ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	repo := NewJobsRepository(db.SQL)

	// Create only completed jobs
	now := time.Now().UTC()
	jobs := []domain.Job{
		{ID: "job1", Type: "download", State: domain.JobCompleted, CreatedAt: now, UpdatedAt: now},
		{ID: "job2", Type: "download", State: domain.JobFailed, CreatedAt: now, UpdatedAt: now},
	}

	for _, j := range jobs {
		_, err := repo.Create(ctx, j)
		if err != nil {
			t.Fatalf("Failed to create job %s: %v", j.ID, err)
		}
	}

	// Load unfinished jobs
	unfinished, err := repo.LoadUnfinishedJobs(ctx)
	if err != nil {
		t.Fatalf("Failed to load unfinished jobs: %v", err)
	}

	// Should return empty slice
	if len(unfinished) != 0 {
		t.Errorf("Expected 0 unfinished jobs, got %d", len(unfinished))
	}
}

// TestStateTransitionTimestamps_FullLifecycle tests full job lifecycle with timestamps.
func TestStateTransitionTimestamps_FullLifecycle(t *testing.T) {
	ctx := context.Background()
	db, err := Open(ctx, ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	repo := NewJobsRepository(db.SQL)

	// Create job
	job := domain.Job{
		ID:        "lifecycle-job",
		Type:      "download",
		State:     domain.JobQueued,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	created, err := repo.Create(ctx, job)
	if err != nil {
		t.Fatalf("Failed to create job: %v", err)
	}

	// Verify initial state
	if created.StartedAt != nil {
		t.Error("StartedAt should be nil for queued job")
	}
	if created.CompletedAt != nil {
		t.Error("CompletedAt should be nil for queued job")
	}

	// Transition: queued -> running
	running, err := repo.UpdateState(ctx, job.ID, domain.JobQueued, domain.JobRunning)
	if err != nil {
		t.Fatalf("Failed to transition to running: %v", err)
	}
	if running.StartedAt == nil {
		t.Error("StartedAt should be set after transitioning to running")
	}
	if running.CompletedAt != nil {
		t.Error("CompletedAt should still be nil for running job")
	}

	// Transition: running -> muxing
	muxing, err := repo.UpdateState(ctx, job.ID, domain.JobRunning, domain.JobMuxing)
	if err != nil {
		t.Fatalf("Failed to transition to muxing: %v", err)
	}
	if muxing.StartedAt == nil {
		t.Error("StartedAt should persist after transitioning to muxing")
	}
	if muxing.CompletedAt != nil {
		t.Error("CompletedAt should still be nil for muxing job")
	}

	// Transition: muxing -> completed
	completed, err := repo.UpdateState(ctx, job.ID, domain.JobMuxing, domain.JobCompleted)
	if err != nil {
		t.Fatalf("Failed to transition to completed: %v", err)
	}
	if completed.StartedAt == nil {
		t.Error("StartedAt should persist in completed job")
	}
	if completed.CompletedAt == nil {
		t.Error("CompletedAt should be set after transitioning to completed")
	}

	// Verify timestamps are reasonable
	if completed.CompletedAt.Before(*completed.StartedAt) {
		t.Error("CompletedAt should be after StartedAt")
	}
	if completed.StartedAt.Before(completed.CreatedAt) {
		t.Error("StartedAt should be after CreatedAt")
	}
}
