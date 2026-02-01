package sqlite

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/ports"
)

// Helper to create in-memory test database with jobs schema
func setupJobsTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	// Create jobs table
	_, err = db.Exec(`
		CREATE TABLE jobs(
			id TEXT PRIMARY KEY,
			type TEXT NOT NULL,
			state TEXT NOT NULL,
			progress REAL NOT NULL DEFAULT 0,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			started_at TEXT,
			completed_at TEXT,
			params_json BLOB,
			result_json BLOB,
			error_code TEXT,
			error_message TEXT
		)
	`)
	if err != nil {
		t.Fatalf("failed to create jobs table: %v", err)
	}

	return db
}

// TestUpdateState_QueuedToRunning_SetsStartedAt verifies startedAt is set
// when transitioning from queued to running
func TestUpdateState_QueuedToRunning_SetsStartedAt(t *testing.T) {
	db := setupJobsTestDB(t)
	defer db.Close()
	repo := NewJobsRepository(db)
	ctx := context.Background()

	// Create a queued job
	now := time.Now().UTC()
	job := domain.Job{
		ID:        "job1",
		Type:      "download",
		State:     domain.JobQueued,
		Progress:  0,
		CreatedAt: now,
		UpdatedAt: now,
	}

	createdJob, err := repo.Create(ctx, job)
	if err != nil {
		t.Fatalf("failed to create job: %v", err)
	}

	if createdJob.StartedAt != nil {
		t.Errorf("expected nil StartedAt, got %v", createdJob.StartedAt)
	}

	// Transition to running
	updated, err := repo.UpdateState(ctx, "job1", domain.JobQueued, domain.JobRunning)

	if err != nil {
		t.Fatalf("failed to update state: %v", err)
	}

	if updated.State != domain.JobRunning {
		t.Errorf("expected state Running, got %v", updated.State)
	}

	if updated.StartedAt == nil {
		t.Error("expected StartedAt to be set")
	}

	if updated.CompletedAt != nil {
		t.Errorf("expected nil CompletedAt, got %v", updated.CompletedAt)
	}
}

// TestUpdateState_RunningToMuxingToCompleted verifies full state sequence
// with proper state machine transitions
func TestUpdateState_RunningToMuxingToCompleted_SetsCompletedAt(t *testing.T) {
	db := setupJobsTestDB(t)
	defer db.Close()
	repo := NewJobsRepository(db)
	ctx := context.Background()

	now := time.Now().UTC()
	job := domain.Job{
		ID:        "job2",
		Type:      "download",
		State:     domain.JobQueued,
		Progress:  0,
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, err := repo.Create(ctx, job)
	if err != nil {
		t.Fatalf("failed to create job: %v", err)
	}

	// Transition: Queued → Running
	_, err = repo.UpdateState(ctx, "job2", domain.JobQueued, domain.JobRunning)
	if err != nil {
		t.Fatalf("failed to transition to running: %v", err)
	}

	// Transition: Running → Muxing
	_, err = repo.UpdateState(ctx, "job2", domain.JobRunning, domain.JobMuxing)
	if err != nil {
		t.Fatalf("failed to transition to muxing: %v", err)
	}

	// Transition: Muxing → Completed
	completed, err := repo.UpdateState(ctx, "job2", domain.JobMuxing, domain.JobCompleted)

	if err != nil {
		t.Fatalf("failed to transition to completed: %v", err)
	}

	if completed.State != domain.JobCompleted {
		t.Errorf("expected state Completed, got %v", completed.State)
	}

	if completed.CompletedAt == nil {
		t.Error("expected CompletedAt to be set")
	}
}

// TestUpdateState_RunningToFailed_SetsCompletedAt verifies completedAt
// is set when transitioning to failed (terminal state)
func TestUpdateState_RunningToFailed_SetsCompletedAt(t *testing.T) {
	db := setupJobsTestDB(t)
	defer db.Close()
	repo := NewJobsRepository(db)
	ctx := context.Background()

	now := time.Now().UTC()
	job := domain.Job{
		ID:        "job3",
		Type:      "download",
		State:     domain.JobQueued,
		Progress:  0,
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, _ = repo.Create(ctx, job)
	_, _ = repo.UpdateState(ctx, "job3", domain.JobQueued, domain.JobRunning)

	failed, err := repo.UpdateState(ctx, "job3", domain.JobRunning, domain.JobFailed)
	if err != nil {
		t.Fatalf("failed to transition to failed: %v", err)
	}

	if failed.State != domain.JobFailed {
		t.Errorf("expected state Failed, got %v", failed.State)
	}

	if failed.CompletedAt == nil {
		t.Error("expected CompletedAt to be set on failed transition")
	}
}

// TestUpdateState_InvalidTransition returns ErrInvalidTransition
// for disallowed state transitions
func TestUpdateState_InvalidTransition(t *testing.T) {
	db := setupJobsTestDB(t)
	defer db.Close()
	repo := NewJobsRepository(db)
	ctx := context.Background()

	now := time.Now().UTC()
	job := domain.Job{
		ID:        "job4",
		Type:      "download",
		State:     domain.JobQueued,
		Progress:  0,
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, _ = repo.Create(ctx, job)

	// Try invalid transition: Queued -> Muxing (skip Running, not allowed)
	_, err := repo.UpdateState(ctx, "job4", domain.JobQueued, domain.JobMuxing)
	if err != domain.ErrInvalidTransition {
		t.Errorf("expected ErrInvalidTransition, got %v", err)
	}
}

// TestUpdateState_ExpectedStateMismatch returns ErrNotFound
// when expected state doesn't match current state (optimistic locking)
func TestUpdateState_ExpectedStateMismatch(t *testing.T) {
	db := setupJobsTestDB(t)
	defer db.Close()
	repo := NewJobsRepository(db)
	ctx := context.Background()

	now := time.Now().UTC()
	job := domain.Job{
		ID:        "job5",
		Type:      "download",
		State:     domain.JobQueued,
		Progress:  0,
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, _ = repo.Create(ctx, job)
	_, _ = repo.UpdateState(ctx, "job5", domain.JobQueued, domain.JobRunning)

	// Try to update with wrong expected state (should fail like optimistic locking)
	_, err := repo.UpdateState(ctx, "job5", domain.JobQueued, domain.JobFailed)
	if err != ports.ErrNotFound {
		t.Errorf("expected ErrNotFound (optimistic lock), got %v", err)
	}
}

// TestLoadUnfinishedJobs returns queued and running jobs, excludes completed
func TestLoadUnfinishedJobs_ReturnsQueuedAndRunning(t *testing.T) {
	db := setupJobsTestDB(t)
	defer db.Close()
	repo := NewJobsRepository(db)
	ctx := context.Background()

	now := time.Now().UTC()

	// Create 5 jobs with different states
	jobs := []domain.Job{
		{ID: "queued1", Type: "download", State: domain.JobQueued, Progress: 0, CreatedAt: now, UpdatedAt: now},
		{ID: "queued2", Type: "download", State: domain.JobQueued, Progress: 0, CreatedAt: now.Add(1 * time.Second), UpdatedAt: now},
		{ID: "running1", Type: "download", State: domain.JobRunning, Progress: 50, CreatedAt: now.Add(2 * time.Second), UpdatedAt: now},
		{ID: "completed1", Type: "download", State: domain.JobCompleted, Progress: 100, CreatedAt: now.Add(3 * time.Second), UpdatedAt: now},
		{ID: "failed1", Type: "download", State: domain.JobFailed, Progress: 25, CreatedAt: now.Add(4 * time.Second), UpdatedAt: now},
	}

	for _, j := range jobs {
		_, _ = repo.Create(ctx, j)
	}

	unfinished, err := repo.LoadUnfinishedJobs(ctx)
	if err != nil {
		t.Fatalf("failed to load unfinished jobs: %v", err)
	}

	if len(unfinished) != 3 {
		t.Errorf("expected 3 unfinished jobs, got %d", len(unfinished))
	}

	// Verify it contains only queued and running, in FIFO order
	expectedIDs := []string{"queued1", "queued2", "running1"}
	for i, expected := range expectedIDs {
		if i >= len(unfinished) || unfinished[i].ID != expected {
			t.Errorf("expected job %s at index %d, got %v", expected, i, unfinished)
		}
	}
}

// TestLoadUnfinishedJobs_EmptyListWhenAllCompleted returns empty when no unfinished
func TestLoadUnfinishedJobs_EmptyListWhenAllCompleted(t *testing.T) {
	db := setupJobsTestDB(t)
	defer db.Close()
	repo := NewJobsRepository(db)
	ctx := context.Background()

	now := time.Now().UTC()

	// Create only terminal jobs
	jobs := []domain.Job{
		{ID: "completed1", Type: "download", State: domain.JobCompleted, Progress: 100, CreatedAt: now, UpdatedAt: now},
		{ID: "failed1", Type: "download", State: domain.JobFailed, Progress: 25, CreatedAt: now.Add(1 * time.Second), UpdatedAt: now},
		{ID: "canceled1", Type: "download", State: domain.JobCanceled, Progress: 10, CreatedAt: now.Add(2 * time.Second), UpdatedAt: now},
	}

	for _, j := range jobs {
		_, _ = repo.Create(ctx, j)
	}

	unfinished, err := repo.LoadUnfinishedJobs(ctx)
	if err != nil {
		t.Fatalf("failed to load unfinished jobs: %v", err)
	}

	if len(unfinished) != 0 {
		t.Errorf("expected 0 unfinished jobs, got %d", len(unfinished))
	}
}

// TestLoadUnfinishedJobs_OrderByCreatedAsc verifies FIFO order (oldest first)
func TestLoadUnfinishedJobs_OrderByCreatedAsc(t *testing.T) {
	db := setupJobsTestDB(t)
	defer db.Close()
	repo := NewJobsRepository(db)
	ctx := context.Background()

	now := time.Now().UTC()

	// Create in reverse order of createdAt
	_, _ = repo.Create(ctx, domain.Job{ID: "job3", Type: "download", State: domain.JobQueued, Progress: 0,
		CreatedAt: now.Add(2 * time.Second), UpdatedAt: now})
	_, _ = repo.Create(ctx, domain.Job{ID: "job1", Type: "download", State: domain.JobQueued, Progress: 0,
		CreatedAt: now, UpdatedAt: now})
	_, _ = repo.Create(ctx, domain.Job{ID: "job2", Type: "download", State: domain.JobQueued, Progress: 0,
		CreatedAt: now.Add(1 * time.Second), UpdatedAt: now})

	unfinished, err := repo.LoadUnfinishedJobs(ctx)
	if err != nil {
		t.Fatalf("failed to load unfinished jobs: %v", err)
	}

	expectedOrder := []string{"job1", "job2", "job3"}
	for i, expected := range expectedOrder {
		if i >= len(unfinished) || unfinished[i].ID != expected {
			t.Errorf("expected job %s at index %d, got %s", expected, i, unfinished[i].ID)
		}
	}
}

// TestUpdateState_QueuedToCanceled_SetsCompletedAt verifies cancellation
// from queued state sets completedAt
func TestUpdateState_QueuedToCanceled_SetsCompletedAt(t *testing.T) {
	db := setupJobsTestDB(t)
	defer db.Close()
	repo := NewJobsRepository(db)
	ctx := context.Background()

	now := time.Now().UTC()
	job := domain.Job{
		ID:        "job7",
		Type:      "download",
		State:     domain.JobQueued,
		Progress:  0,
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, _ = repo.Create(ctx, job)

	canceled, err := repo.UpdateState(ctx, "job7", domain.JobQueued, domain.JobCanceled)
	if err != nil {
		t.Fatalf("failed to cancel job: %v", err)
	}

	if canceled.CompletedAt == nil {
		t.Error("expected CompletedAt to be set on queued->canceled transition")
	}
}

// TestUpdateState_ProgressPreservedDuringStateChange verifies progress
// is preserved when state changes
func TestUpdateState_ProgressPreserved(t *testing.T) {
	db := setupJobsTestDB(t)
	defer db.Close()
	repo := NewJobsRepository(db)
	ctx := context.Background()

	now := time.Now().UTC()
	job := domain.Job{
		ID:        "job8",
		Type:      "download",
		State:     domain.JobQueued,
		Progress:  0,
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, _ = repo.Create(ctx, job)
	_, _ = repo.UpdateState(ctx, "job8", domain.JobQueued, domain.JobRunning)

	// Update progress
	withProgress, err := repo.UpdateProgress(ctx, "job8", 75.5)
	if err != nil {
		t.Fatalf("failed to update progress: %v", err)
	}

	if withProgress.Progress != 75.5 {
		t.Errorf("expected progress 75.5, got %v", withProgress.Progress)
	}

	// State change should preserve progress
	updated, err := repo.UpdateState(ctx, "job8", domain.JobRunning, domain.JobMuxing)
	if err != nil {
		t.Fatalf("failed to transition to muxing: %v", err)
	}

	if updated.Progress != 75.5 {
		t.Errorf("expected progress preserved at 75.5, got %v", updated.Progress)
	}
}
