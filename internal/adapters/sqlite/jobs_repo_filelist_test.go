package sqlite

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
)

// TestJobsRepository_FileListJSON_Store tests storing file list metadata with jobs
func TestJobsRepository_FileListJSON_Store(t *testing.T) {
	db := setupJobsTestDB(t)
	defer db.Close()

	repo := NewJobsRepository(db)
	ctx := context.Background()

	// Create a job with file list metadata
	fileList := domain.FileList{
		AnimeID: "1",
		Files: []domain.File{
			{
				ID:       "1-ep1",
				Name:     "Episode 1",
				Path:     "/downloads/ep1.mkv",
				Size:     350000000,
				Duration: 1400,
				Type:     "video/x-matroska",
			},
		},
	}
	fileListJSON, _ := json.Marshal(fileList)

	job := domain.Job{
		ID:         "job-123",
		Type:       "download",
		State:      domain.JobQueued,
		Progress:   0,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
		ParamsJSON: []byte(`{"anime_id":"1"}`),
	}

	// Create the job
	created, err := repo.Create(ctx, job)
	if err != nil {
		t.Fatalf("failed to create job: %v", err)
	}

	if created.ID != "job-123" {
		t.Errorf("expected job ID 'job-123', got '%s'", created.ID)
	}

	// Now update with file list metadata
	err = repo.UpdateFileList(ctx, "job-123", fileListJSON)
	if err != nil {
		t.Fatalf("failed to update file list: %v", err)
	}

	// Retrieve and verify
	retrieved, err := repo.Get(ctx, "job-123")
	if err != nil {
		t.Fatalf("failed to get job: %v", err)
	}

	if len(retrieved.FileListJSON) == 0 {
		t.Error("expected file list metadata to be stored")
	}

	// Verify we can unmarshal the file list
	var retrievedFileList domain.FileList
	err = json.Unmarshal(retrieved.FileListJSON, &retrievedFileList)
	if err != nil {
		t.Fatalf("failed to unmarshal file list: %v", err)
	}

	if retrievedFileList.AnimeID != "1" {
		t.Errorf("expected anime ID '1', got '%s'", retrievedFileList.AnimeID)
	}

	if len(retrievedFileList.Files) != 1 {
		t.Errorf("expected 1 file, got %d", len(retrievedFileList.Files))
	}
}

// TestJobsRepository_FileListJSON_Optional tests that file list is optional
func TestJobsRepository_FileListJSON_Optional(t *testing.T) {
	db := setupJobsTestDB(t)
	defer db.Close()

	repo := NewJobsRepository(db)
	ctx := context.Background()

	// Create a job without file list
	job := domain.Job{
		ID:         "job-456",
		Type:       "download",
		State:      domain.JobQueued,
		Progress:   0,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
		ParamsJSON: []byte(`{"anime_id":"2"}`),
	}

	created, err := repo.Create(ctx, job)
	if err != nil {
		t.Fatalf("failed to create job: %v", err)
	}

	// Verify file list is empty (not an error, just optional)
	if created.FileListJSON != nil && len(created.FileListJSON) > 0 {
		t.Error("expected empty file list for new job")
	}
}

// TestJobsRepository_LoadUnfinishedJobs_WithFileList tests that LoadUnfinishedJobs includes file list metadata
func TestJobsRepository_LoadUnfinishedJobs_WithFileList(t *testing.T) {
	db := setupJobsTestDB(t)
	defer db.Close()

	repo := NewJobsRepository(db)
	ctx := context.Background()

	// Create two unfinished jobs, one with file list
	fileList := domain.FileList{
		AnimeID: "1",
		Files: []domain.File{
			{ID: "1-ep1", Name: "Episode 1", Path: "/ep1.mkv", Size: 350000000, Duration: 1400, Type: "video/x-matroska"},
		},
	}
	fileListJSON, _ := json.Marshal(fileList)

	job1 := domain.Job{
		ID:         "job-unfinished-1",
		Type:       "download",
		State:      domain.JobQueued,
		Progress:   0,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
		ParamsJSON: []byte(`{"anime_id":"1"}`),
	}

	job2 := domain.Job{
		ID:         "job-unfinished-2",
		Type:       "download",
		State:      domain.JobRunning,
		Progress:   50,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
		ParamsJSON: []byte(`{"anime_id":"2"}`),
	}

	repo.Create(ctx, job1)
	repo.Create(ctx, job2)

	// Add file list to second job
	repo.UpdateFileList(ctx, "job-unfinished-2", fileListJSON)

	// Load unfinished jobs
	jobs, err := repo.LoadUnfinishedJobs(ctx)
	if err != nil {
		t.Fatalf("failed to load unfinished jobs: %v", err)
	}

	if len(jobs) < 2 {
		t.Errorf("expected at least 2 unfinished jobs, got %d", len(jobs))
	}

	// Find job2 in results
	var job2Result *domain.Job
	for i := range jobs {
		if jobs[i].ID == "job-unfinished-2" {
			job2Result = &jobs[i]
			break
		}
	}

	if job2Result == nil {
		t.Error("expected to find job-unfinished-2 in results")
		return
	}

	// Verify file list is included
	if len(job2Result.FileListJSON) == 0 {
		t.Error("expected file list metadata to be loaded")
	}

	var loadedFileList domain.FileList
	json.Unmarshal(job2Result.FileListJSON, &loadedFileList)
	if loadedFileList.AnimeID != "1" {
		t.Errorf("expected anime ID '1', got '%s'", loadedFileList.AnimeID)
	}
}

// TestJobsRepository_FileListJSON_ClearOnUpdate tests that file list can be cleared
func TestJobsRepository_FileListJSON_ClearOnUpdate(t *testing.T) {
	db := setupJobsTestDB(t)
	defer db.Close()

	repo := NewJobsRepository(db)
	ctx := context.Background()

	fileList := domain.FileList{
		AnimeID: "1",
		Files: []domain.File{
			{ID: "1-ep1", Name: "Episode 1", Path: "/ep1.mkv", Size: 350000000, Duration: 1400, Type: "video/x-matroska"},
		},
	}
	fileListJSON, _ := json.Marshal(fileList)

	job := domain.Job{
		ID:         "job-update-list",
		Type:       "download",
		State:      domain.JobQueued,
		Progress:   0,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
		ParamsJSON: []byte(`{"anime_id":"1"}`),
	}

	repo.Create(ctx, job)
	repo.UpdateFileList(ctx, "job-update-list", fileListJSON)

	// Verify it was stored
	retrieved1, _ := repo.Get(ctx, "job-update-list")
	if len(retrieved1.FileListJSON) == 0 {
		t.Error("file list should be stored initially")
	}

	// Clear the file list by passing nil
	repo.UpdateFileList(ctx, "job-update-list", nil)

	// Verify it was cleared
	retrieved2, _ := repo.Get(ctx, "job-update-list")
	if retrieved2.FileListJSON != nil && len(retrieved2.FileListJSON) > 0 {
		t.Error("file list should be cleared")
	}
}
