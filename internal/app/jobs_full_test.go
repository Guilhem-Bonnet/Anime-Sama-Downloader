package app

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
)

// MockJobRepository for testing
type mockJobRepository struct {
	createFn          func(ctx context.Context, job domain.Job) (domain.Job, error)
	getFn             func(ctx context.Context, id string) (domain.Job, error)
	listFn            func(ctx context.Context, limit int) ([]domain.Job, error)
	claimNextQueuedFn func(ctx context.Context) (domain.Job, error)
	updateProgressFn  func(ctx context.Context, id string, progress float64) (domain.Job, error)
	updateResultFn    func(ctx context.Context, id string, resultJSON []byte) (domain.Job, error)
	updateErrorFn     func(ctx context.Context, id string, code string, message string) (domain.Job, error)
	updateStateFn     func(ctx context.Context, id string, expected domain.JobState, next domain.JobState) (domain.Job, error)
}

func (m *mockJobRepository) Create(ctx context.Context, job domain.Job) (domain.Job, error) {
	if m.createFn != nil {
		return m.createFn(ctx, job)
	}
	return job, nil
}

func (m *mockJobRepository) Get(ctx context.Context, id string) (domain.Job, error) {
	if m.getFn != nil {
		return m.getFn(ctx, id)
	}
	return domain.Job{}, nil
}

func (m *mockJobRepository) List(ctx context.Context, limit int) ([]domain.Job, error) {
	if m.listFn != nil {
		return m.listFn(ctx, limit)
	}
	return []domain.Job{}, nil
}

func (m *mockJobRepository) ClaimNextQueued(ctx context.Context) (domain.Job, error) {
	if m.claimNextQueuedFn != nil {
		return m.claimNextQueuedFn(ctx)
	}
	return domain.Job{}, errors.New("no queued job")
}

func (m *mockJobRepository) UpdateProgress(ctx context.Context, id string, progress float64) (domain.Job, error) {
	if m.updateProgressFn != nil {
		return m.updateProgressFn(ctx, id, progress)
	}
	return domain.Job{}, nil
}

func (m *mockJobRepository) UpdateResult(ctx context.Context, id string, resultJSON []byte) (domain.Job, error) {
	if m.updateResultFn != nil {
		return m.updateResultFn(ctx, id, resultJSON)
	}
	return domain.Job{}, nil
}

func (m *mockJobRepository) UpdateError(ctx context.Context, id string, code string, message string) (domain.Job, error) {
	if m.updateErrorFn != nil {
		return m.updateErrorFn(ctx, id, code, message)
	}
	return domain.Job{}, nil
}

func (m *mockJobRepository) UpdateState(ctx context.Context, id string, expected domain.JobState, next domain.JobState) (domain.Job, error) {
	if m.updateStateFn != nil {
		return m.updateStateFn(ctx, id, expected, next)
	}
	return domain.Job{}, nil
}

// Test cases for JobService.Create()
func TestJobService_Create_Success(t *testing.T) {
	params := CreateJobRequest{
		Type:   "download",
		Params: json.RawMessage(`{"url": "https://example.com/video.mp4"}`),
	}

	repo := &mockJobRepository{
		createFn: func(ctx context.Context, job domain.Job) (domain.Job, error) {
			return job, nil
		},
	}

	bus := &mockEventBus{}
	service := NewJobService(repo, bus)
	ctx := context.Background()

	dto, err := service.Create(ctx, params)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dto.ID == "" {
		t.Error("expected ID to be set")
	}
	if dto.Type != "download" {
		t.Errorf("expected type 'download', got %q", dto.Type)
	}
	if dto.State != domain.JobQueued {
		t.Errorf("expected initial state 'queued', got %v", dto.State)
	}
	if dto.Progress != 0 {
		t.Errorf("expected initial progress 0, got %v", dto.Progress)
	}
}

func TestJobService_Create_WithEmptyParams(t *testing.T) {
	params := CreateJobRequest{
		Type:   "cleanup",
		Params: json.RawMessage{},
	}

	repo := &mockJobRepository{
		createFn: func(ctx context.Context, job domain.Job) (domain.Job, error) {
			return job, nil
		},
	}

	service := NewJobService(repo, nil)
	ctx := context.Background()

	dto, err := service.Create(ctx, params)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dto.Type != "cleanup" {
		t.Errorf("expected type 'cleanup', got %q", dto.Type)
	}
}

func TestJobService_Create_RepositoryError(t *testing.T) {
	params := CreateJobRequest{
		Type:   "download",
		Params: json.RawMessage(`{}`),
	}

	repo := &mockJobRepository{
		createFn: func(ctx context.Context, job domain.Job) (domain.Job, error) {
			return domain.Job{}, errors.New("database error")
		},
	}

	service := NewJobService(repo, nil)
	ctx := context.Background()

	_, err := service.Create(ctx, params)

	if err == nil {
		t.Fatal("expected error from repository")
	}
}

// Test cases for JobService.Get()
func TestJobService_Get_Success(t *testing.T) {
	job := domain.Job{
		ID:    "job-1",
		Type:  "download",
		State: domain.JobQueued,
	}

	repo := &mockJobRepository{
		getFn: func(ctx context.Context, id string) (domain.Job, error) {
			if id == "job-1" {
				return job, nil
			}
			return domain.Job{}, errors.New("not found")
		},
	}

	service := NewJobService(repo, nil)
	ctx := context.Background()

	result, err := service.Get(ctx, "job-1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "job-1" {
		t.Errorf("expected ID 'job-1', got %q", result.ID)
	}
	if result.State != domain.JobQueued {
		t.Errorf("expected state 'queued', got %v", result.State)
	}
}

func TestJobService_Get_NotFound(t *testing.T) {
	repo := &mockJobRepository{
		getFn: func(ctx context.Context, id string) (domain.Job, error) {
			return domain.Job{}, errors.New("not found")
		},
	}

	service := NewJobService(repo, nil)
	ctx := context.Background()

	_, err := service.Get(ctx, "non-existent")

	if err == nil {
		t.Fatal("expected error for non-existent job")
	}
}

// Test cases for JobService.List()
func TestJobService_List_Success(t *testing.T) {
	jobs := []domain.Job{
		{ID: "job-1", Type: "download", State: domain.JobQueued},
		{ID: "job-2", Type: "download", State: domain.JobQueued},
	}

	repo := &mockJobRepository{
		listFn: func(ctx context.Context, limit int) ([]domain.Job, error) {
			if limit < len(jobs) {
				return jobs[:limit], nil
			}
			return jobs, nil
		},
	}

	service := NewJobService(repo, nil)
	ctx := context.Background()

	result, err := service.List(ctx, 10)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 jobs, got %d", len(result))
	}
	if result[0].ID != "job-1" {
		t.Errorf("expected first job ID 'job-1', got %q", result[0].ID)
	}
}

func TestJobService_List_Empty(t *testing.T) {
	repo := &mockJobRepository{
		listFn: func(ctx context.Context, limit int) ([]domain.Job, error) {
			return []domain.Job{}, nil
		},
	}

	service := NewJobService(repo, nil)
	ctx := context.Background()

	result, err := service.List(ctx, 10)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected 0 jobs, got %d", len(result))
	}
}

func TestJobService_List_RespectLimit(t *testing.T) {
	jobs := []domain.Job{
		{ID: "job-1", Type: "download"},
		{ID: "job-2", Type: "download"},
		{ID: "job-3", Type: "download"},
	}

	repo := &mockJobRepository{
		listFn: func(ctx context.Context, limit int) ([]domain.Job, error) {
			if limit < len(jobs) {
				return jobs[:limit], nil
			}
			return jobs, nil
		},
	}

	service := NewJobService(repo, nil)
	ctx := context.Background()

	result, err := service.List(ctx, 2)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 jobs due to limit, got %d", len(result))
	}
}

// Test cases for JobService.UpdateProgress()
func TestJobService_Cancel_Success(t *testing.T) {
	updatedJob := domain.Job{
		ID:       "job-1",
		Type:     "download",
		State:    domain.JobCanceled,
		Progress: 0.75,
	}

	repo := &mockJobRepository{
		updateStateFn: func(ctx context.Context, id string, expected domain.JobState, next domain.JobState) (domain.Job, error) {
			if id == "job-1" && next == domain.JobCanceled {
				return updatedJob, nil
			}
			return domain.Job{}, errors.New("not found")
		},
	}

	service := NewJobService(repo, nil)
	ctx := context.Background()

	result, err := service.Cancel(ctx, "job-1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.State != domain.JobCanceled {
		t.Errorf("expected state 'canceled', got %v", result.State)
	}
}

func TestJobService_Cancel_AlreadyRunning(t *testing.T) {
	stateBeforeCancel := domain.JobRunning
	canceledJob := domain.Job{
		ID:    "job-1",
		Type:  "download",
		State: domain.JobCanceled,
	}

	repo := &mockJobRepository{
		updateStateFn: func(ctx context.Context, id string, expected domain.JobState, next domain.JobState) (domain.Job, error) {
			// Simuler: queued check échoue, running check réussit
			if expected == domain.JobRunning && next == domain.JobCanceled {
				return canceledJob, nil
			}
			return domain.Job{}, errors.New("invalid transition")
		},
		getFn: func(ctx context.Context, id string) (domain.Job, error) {
			return domain.Job{ID: "job-1", State: stateBeforeCancel}, nil
		},
	}

	service := NewJobService(repo, nil)
	ctx := context.Background()

	result, err := service.Cancel(ctx, "job-1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.State != domain.JobCanceled {
		t.Errorf("expected state to be canceled, got %v", result.State)
	}
}

// Test EventBus integration
func TestJobService_PublishesCreatedEvent(t *testing.T) {
	eventPublished := false
	var publishedTopic string

	repo := &mockJobRepository{
		createFn: func(ctx context.Context, job domain.Job) (domain.Job, error) {
			return job, nil
		},
	}

	bus := &mockEventBus{
		publishFn: func(topic string, data []byte) {
			eventPublished = true
			publishedTopic = topic
		},
	}

	service := NewJobService(repo, bus)
	ctx := context.Background()

	params := CreateJobRequest{
		Type:   "download",
		Params: json.RawMessage(`{}`),
	}

	service.Create(ctx, params)

	if !eventPublished {
		t.Error("expected event to be published")
	}
	if publishedTopic != "job.created" {
		t.Errorf("expected topic 'job.created', got %q", publishedTopic)
	}
}

// Test JobDTO conversion
func TestJobDTO_Conversion(t *testing.T) {
	now := time.Now().UTC()
	job := domain.Job{
		ID:           "job-1",
		Type:         "download",
		State:        domain.JobQueued,
		Progress:     0.0,
		CreatedAt:    now,
		UpdatedAt:    now,
		ParamsJSON:   []byte(`{"url": "https://example.com"}`),
		ResultJSON:   []byte{},
		ErrorCode:    "",
		ErrorMessage: "",
	}

	dto := ToJobDTO(job)

	if dto.ID != job.ID {
		t.Errorf("ID mismatch: expected %q, got %q", job.ID, dto.ID)
	}
	if dto.Type != job.Type {
		t.Errorf("Type mismatch: expected %q, got %q", job.Type, dto.Type)
	}
	if dto.State != job.State {
		t.Errorf("State mismatch: expected %v, got %v", job.State, dto.State)
	}
	if dto.Progress != job.Progress {
		t.Errorf("Progress mismatch: expected %v, got %v", job.Progress, dto.Progress)
	}
}

// Test JobService initialization
func TestJobService_NewJobService(t *testing.T) {
	repo := &mockJobRepository{}
	bus := &mockEventBus{}
	service := NewJobService(repo, bus)

	if service == nil {
		t.Fatal("expected service to be initialized")
	}
}

// Test with nil EventBus (should not panic)
func TestJobService_Create_WithNilEventBus(t *testing.T) {
	repo := &mockJobRepository{
		createFn: func(ctx context.Context, job domain.Job) (domain.Job, error) {
			return job, nil
		},
	}

	service := NewJobService(repo, nil)
	ctx := context.Background()

	params := CreateJobRequest{
		Type:   "download",
		Params: json.RawMessage(`{}`),
	}

	_, err := service.Create(ctx, params)

	if err != nil {
		t.Fatalf("unexpected error with nil EventBus: %v", err)
	}
}
