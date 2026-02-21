package sqlite_test

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/adapters/sqlite"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
)

// TestJobRecovery_ApplicationRestart tests the complete job recovery scenario.
func TestJobRecovery_ApplicationRestart(t *testing.T) {
	// Create a temporary database file
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	ctx := context.Background()

	// === PHASE 1: Initial application session ===
	db1, err := sqlite.Open(ctx, dbPath)
	if err != nil {
		t.Fatalf("Failed to open database (session 1): %v", err)
	}

	repo1 := sqlite.NewJobsRepository(db1.SQL)

	// Create jobs
	now := time.Now().UTC()
	jobs := []domain.Job{
		{ID: "job1", Type: "download", State: domain.JobQueued, CreatedAt: now.Add(-3 * time.Minute), UpdatedAt: now.Add(-3 * time.Minute)},
		{ID: "job2", Type: "download", State: domain.JobQueued, CreatedAt: now.Add(-2 * time.Minute), UpdatedAt: now.Add(-2 * time.Minute)},
		{ID: "job3", Type: "download", State: domain.JobCompleted, CreatedAt: now.Add(-5 * time.Minute), UpdatedAt: now.Add(-5 * time.Minute)},
	}

	for _, j := range jobs {
		_, err := repo1.Create(ctx, j)
		if err != nil {
			t.Fatalf("Failed to create job %s: %v", j.ID, err)
		}
	}

	// Start job2 (simulate it was in progress)
	_, err = repo1.UpdateState(ctx, "job2", domain.JobQueued, domain.JobRunning)
	if err != nil {
		t.Fatalf("Failed to start job2: %v", err)
	}

	// Close database (simulate application shutdown)
	if err := db1.Close(); err != nil {
		t.Fatalf("Failed to close database: %v", err)
	}

	// === PHASE 2: Application restart ===
	db2, err := sqlite.Open(ctx, dbPath)
	if err != nil {
		t.Fatalf("Failed to open database (session 2): %v", err)
	}
	defer db2.Close()

	repo2 := sqlite.NewJobsRepository(db2.SQL)

	// Measure recovery time
	startRecovery := time.Now()
	unfinished, err := repo2.LoadUnfinishedJobs(ctx)
	recoveryDuration := time.Since(startRecovery)

	if err != nil {
		t.Fatalf("Failed to load unfinished jobs: %v", err)
	}

	// Verify recovery time (AC5: within 5 seconds)
	if recoveryDuration > 5*time.Second {
		t.Errorf("Recovery took too long: %v (expected < 5s)", recoveryDuration)
	}

	// Verify correct jobs were recovered
	if len(unfinished) != 2 {
		t.Fatalf("Expected 2 unfinished jobs, got %d", len(unfinished))
	}

	// Verify job order (FIFO: job1, job2)
	expectedOrder := []string{"job1", "job2"}
	for i, expectedID := range expectedOrder {
		if unfinished[i].ID != expectedID {
			t.Errorf("Expected job at position %d to be %s, got %s", i, expectedID, unfinished[i].ID)
		}
	}

	// Verify job1 is still queued
	if unfinished[0].ID == "job1" && unfinished[0].State != domain.JobQueued {
		t.Errorf("Job1 should be queued, got %v", unfinished[0].State)
	}

	// Verify job2 is still running
	if unfinished[1].ID == "job2" && unfinished[1].State != domain.JobRunning {
		t.Errorf("Job2 should be running, got %v", unfinished[1].State)
	}

	// Verify job2 has StartedAt set
	if unfinished[1].ID == "job2" && unfinished[1].StartedAt == nil {
		t.Error("Job2 should have StartedAt set from before restart")
	}
}

// TestJobRecovery_MultipleRestarts tests recovery across multiple restarts.
func TestJobRecovery_MultipleRestarts(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	ctx := context.Background()

	// Session 1: Create job
	db1, _ := sqlite.Open(ctx, dbPath)
	repo1 := sqlite.NewJobsRepository(db1.SQL)

	job := domain.Job{
		ID:        "persistent-job",
		Type:      "download",
		State:     domain.JobQueued,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	repo1.Create(ctx, job)
	db1.Close()

	// Session 2: Start job
	db2, _ := sqlite.Open(ctx, dbPath)
	repo2 := sqlite.NewJobsRepository(db2.SQL)
	repo2.UpdateState(ctx, "persistent-job", domain.JobQueued, domain.JobRunning)
	db2.Close()

	// Session 3: Recover and verify
	db3, _ := sqlite.Open(ctx, dbPath)
	defer db3.Close()
	repo3 := sqlite.NewJobsRepository(db3.SQL)

	unfinished, err := repo3.LoadUnfinishedJobs(ctx)
	if err != nil {
		t.Fatalf("Failed to recover: %v", err)
	}

	if len(unfinished) != 1 {
		t.Fatalf("Expected 1 unfinished job, got %d", len(unfinished))
	}

	if unfinished[0].State != domain.JobRunning {
		t.Errorf("Job should be running, got %v", unfinished[0].State)
	}

	if unfinished[0].StartedAt == nil {
		t.Error("Job should have StartedAt preserved across restarts")
	}
}

// TestJobRecovery_PerformanceWithManyJobs tests recovery performance with many jobs.
func TestJobRecovery_PerformanceWithManyJobs(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	ctx := context.Background()

	// Create database and populate with many jobs
	db, err := sqlite.Open(ctx, dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	repo := sqlite.NewJobsRepository(db.SQL)

	// Create 100 unfinished jobs and 50 completed jobs
	now := time.Now().UTC()
	for i := 0; i < 100; i++ {
		job := domain.Job{
			ID:        fmt.Sprintf("job-%d", i),
			Type:      "download",
			State:     domain.JobQueued,
			CreatedAt: now.Add(-time.Duration(i) * time.Minute),
			UpdatedAt: now.Add(-time.Duration(i) * time.Minute),
		}
		repo.Create(ctx, job)
	}

	for i := 100; i < 150; i++ {
		job := domain.Job{
			ID:        fmt.Sprintf("job-%d", i),
			Type:      "download",
			State:     domain.JobCompleted,
			CreatedAt: now.Add(-time.Duration(i) * time.Minute),
			UpdatedAt: now.Add(-time.Duration(i) * time.Minute),
		}
		repo.Create(ctx, job)
	}

	// Close and reopen (simulate restart)
	db.Close()
	db, _ = sqlite.Open(ctx, dbPath)
	defer db.Close()
	repo = sqlite.NewJobsRepository(db.SQL)

	// Measure recovery time
	start := time.Now()
	unfinished, err := repo.LoadUnfinishedJobs(ctx)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("Failed to load unfinished jobs: %v", err)
	}

	// Should recover only unfinished jobs
	if len(unfinished) != 100 {
		t.Errorf("Expected 100 unfinished jobs, got %d", len(unfinished))
	}

	// Performance check (should complete within 5 seconds)
	if duration > 5*time.Second {
		t.Errorf("Recovery took too long: %v (expected < 5s)", duration)
	}

	t.Logf("Recovery of %d jobs completed in %v", len(unfinished), duration)
}

// TestJobRecovery_EmptyDatabase tests recovery with no jobs.
func TestJobRecovery_EmptyDatabase(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	ctx := context.Background()

	db, err := sqlite.Open(ctx, dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	repo := sqlite.NewJobsRepository(db.SQL)

	unfinished, err := repo.LoadUnfinishedJobs(ctx)
	if err != nil {
		t.Fatalf("Failed to load unfinished jobs: %v", err)
	}

	if len(unfinished) != 0 {
		t.Errorf("Expected 0 unfinished jobs in empty database, got %d", len(unfinished))
	}
}

// TestJobRecovery_OnlyCompletedJobs tests recovery when all jobs are completed.
func TestJobRecovery_OnlyCompletedJobs(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	ctx := context.Background()

	db, err := sqlite.Open(ctx, dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	repo := sqlite.NewJobsRepository(db.SQL)

	// Create only completed jobs
	now := time.Now().UTC()
	for i := 0; i < 10; i++ {
		job := domain.Job{
			ID:        fmt.Sprintf("completed-job-%d", i),
			Type:      "download",
			State:     domain.JobCompleted,
			CreatedAt: now,
			UpdatedAt: now,
		}
		repo.Create(ctx, job)
	}

	// Close and reopen
	db.Close()
	db, _ = sqlite.Open(ctx, dbPath)
	defer db.Close()
	repo = sqlite.NewJobsRepository(db.SQL)

	unfinished, err := repo.LoadUnfinishedJobs(ctx)
	if err != nil {
		t.Fatalf("Failed to load unfinished jobs: %v", err)
	}

	if len(unfinished) != 0 {
		t.Errorf("Expected 0 unfinished jobs when all are completed, got %d", len(unfinished))
	}
}
