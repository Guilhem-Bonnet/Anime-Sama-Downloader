package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/app"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
	"github.com/go-chi/chi/v5"
)

// Mock JobRepository for tests
type testJobRepository struct {
	createFn        func(ctx context.Context, job domain.Job) (domain.Job, error)
	getFn           func(ctx context.Context, id string) (domain.Job, error)
	listFn          func(ctx context.Context, limit int) ([]domain.Job, error)
	claimNextQueuedFn func(ctx context.Context) (domain.Job, error)
	updateProgressFn func(ctx context.Context, id string, progress float64) (domain.Job, error)
	updateResultFn  func(ctx context.Context, id string, resultJSON []byte) (domain.Job, error)
	updateErrorFn   func(ctx context.Context, id string, code string, message string) (domain.Job, error)
	updateStateFn   func(ctx context.Context, id string, expected domain.JobState, next domain.JobState) (domain.Job, error)
}

func (m *testJobRepository) Create(ctx context.Context, job domain.Job) (domain.Job, error) {
	if m.createFn != nil {
		return m.createFn(ctx, job)
	}
	return job, nil
}

func (m *testJobRepository) Get(ctx context.Context, id string) (domain.Job, error) {
	if m.getFn != nil {
		return m.getFn(ctx, id)
	}
	return domain.Job{}, nil
}

func (m *testJobRepository) List(ctx context.Context, limit int) ([]domain.Job, error) {
	if m.listFn != nil {
		return m.listFn(ctx, limit)
	}
	return []domain.Job{}, nil
}

func (m *testJobRepository) ClaimNextQueued(ctx context.Context) (domain.Job, error) {
	if m.claimNextQueuedFn != nil {
		return m.claimNextQueuedFn(ctx)
	}
	return domain.Job{}, errors.New("no queued job")
}

func (m *testJobRepository) UpdateProgress(ctx context.Context, id string, progress float64) (domain.Job, error) {
	if m.updateProgressFn != nil {
		return m.updateProgressFn(ctx, id, progress)
	}
	return domain.Job{}, nil
}

func (m *testJobRepository) UpdateResult(ctx context.Context, id string, resultJSON []byte) (domain.Job, error) {
	if m.updateResultFn != nil {
		return m.updateResultFn(ctx, id, resultJSON)
	}
	return domain.Job{}, nil
}

func (m *testJobRepository) UpdateError(ctx context.Context, id string, code string, message string) (domain.Job, error) {
	if m.updateErrorFn != nil {
		return m.updateErrorFn(ctx, id, code, message)
	}
	return domain.Job{}, nil
}

func (m *testJobRepository) UpdateState(ctx context.Context, id string, expected domain.JobState, next domain.JobState) (domain.Job, error) {
	if m.updateStateFn != nil {
		return m.updateStateFn(ctx, id, expected, next)
	}
	return domain.Job{}, nil
}

// Test POST /jobs - Create job
func TestJobsHandler_Create_Success(t *testing.T) {
	jobRepo := &testJobRepository{
		createFn: func(ctx context.Context, job domain.Job) (domain.Job, error) {
			return job, nil
		},
	}
	jobService := app.NewJobService(jobRepo, nil)
	handler := NewJobsHandler(jobService)
	req := app.CreateJobRequest{
		Type:   "download",
		Params: json.RawMessage(`{"url": "https://example.com"}`),
	}

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/jobs", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.create(w, httpReq)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var response app.JobDTO
	json.Unmarshal(w.Body.Bytes(), &response)

	if response.ID == "" {
		t.Error("expected job ID to be set")
	}
	if response.Type != "download" {
		t.Errorf("expected job type 'download', got %q", response.Type)
	}
}

func TestJobsHandler_Create_InvalidJSON(t *testing.T) {
	jobRepo := &testJobRepository{}
	handler := NewJobsHandler(app.NewJobService(jobRepo, nil))

	httpReq := httptest.NewRequest("POST", "/jobs", bytes.NewBuffer([]byte("invalid json")))
	w := httptest.NewRecorder()

	handler.create(w, httpReq)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestJobsHandler_Create_MissingType(t *testing.T) {
	jobRepo := &testJobRepository{}
	handler := NewJobsHandler(app.NewJobService(jobRepo, nil))

	req := app.CreateJobRequest{
		Type: "", // missing type
	}

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/jobs", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.create(w, httpReq)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// Test GET /jobs - List jobs
func TestJobsHandler_List_Success(t *testing.T) {
	expectedJobs := []domain.Job{
		{ID: "job-1", Type: "download", State: domain.JobQueued},
		{ID: "job-2", Type: "download", State: domain.JobQueued},
	}

	jobRepo := &testJobRepository{
		listFn: func(ctx context.Context, limit int) ([]domain.Job, error) {
			return expectedJobs, nil
		},
	}

	handler := NewJobsHandler(app.NewJobService(jobRepo, nil))

	httpReq := httptest.NewRequest("GET", "/jobs?limit=10", nil)
	w := httptest.NewRecorder()

	handler.list(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response []app.JobDTO
	json.Unmarshal(w.Body.Bytes(), &response)

	if len(response) != 2 {
		t.Errorf("expected 2 jobs, got %d", len(response))
	}
}

func TestJobsHandler_List_Empty(t *testing.T) {
	jobRepo := &testJobRepository{
		listFn: func(ctx context.Context, limit int) ([]domain.Job, error) {
			return []domain.Job{}, nil
		},
	}

	handler := NewJobsHandler(app.NewJobService(jobRepo, nil))

	httpReq := httptest.NewRequest("GET", "/jobs", nil)
	w := httptest.NewRecorder()

	handler.list(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response []app.JobDTO
	json.Unmarshal(w.Body.Bytes(), &response)

	if len(response) != 0 {
		t.Errorf("expected 0 jobs, got %d", len(response))
	}
}

// Test GET /jobs/{id} - Get job
func TestJobsHandler_Get_Success(t *testing.T) {
	expectedJob := domain.Job{
		ID:    "job-1",
		Type:  "download",
		State: domain.JobQueued,
	}

	jobRepo := &testJobRepository{
		getFn: func(ctx context.Context, id string) (domain.Job, error) {
			if id == "job-1" {
				return expectedJob, nil
			}
			return domain.Job{}, app.ErrNotFound
		},
	}

	handler := NewJobsHandler(app.NewJobService(jobRepo, nil))

	httpReq := httptest.NewRequest("GET", "/jobs/job-1", nil)
	w := httptest.NewRecorder()

	// Simulate chi routing
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "job-1")
	httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), chi.RouteCtxKey, rctx))

	handler.get(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response app.JobDTO
	json.Unmarshal(w.Body.Bytes(), &response)

	if response.ID != "job-1" {
		t.Errorf("expected job ID 'job-1', got %q", response.ID)
	}
}

func TestJobsHandler_Get_NotFound(t *testing.T) {
	jobRepo := &testJobRepository{
		getFn: func(ctx context.Context, id string) (domain.Job, error) {
			return domain.Job{}, app.ErrNotFound
		},
	}

	handler := NewJobsHandler(app.NewJobService(jobRepo, nil))

	httpReq := httptest.NewRequest("GET", "/jobs/non-existent", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "non-existent")
	httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), chi.RouteCtxKey, rctx))

	handler.get(w, httpReq)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

// Test POST /jobs/{id}/cancel - Cancel job
func TestJobsHandler_Cancel_Success(t *testing.T) {
	canceledJob := domain.Job{
		ID:    "job-1",
		Type:  "download",
		State: domain.JobCanceled,
	}

	jobRepo := &testJobRepository{
		getFn: func(ctx context.Context, id string) (domain.Job, error) {
			if id == "job-1" {
				return canceledJob, nil
			}
			return domain.Job{}, app.ErrNotFound
		},
		updateStateFn: func(ctx context.Context, id string, expected domain.JobState, next domain.JobState) (domain.Job, error) {
			if id == "job-1" && next == domain.JobCanceled {
				return canceledJob, nil
			}
			return domain.Job{}, app.ErrNotFound
		},
	}

	handler := NewJobsHandler(app.NewJobService(jobRepo, nil))

	httpReq := httptest.NewRequest("POST", "/jobs/job-1/cancel", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "job-1")
	httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), chi.RouteCtxKey, rctx))

	handler.cancel(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response app.JobDTO
	json.Unmarshal(w.Body.Bytes(), &response)

	if response.State != domain.JobCanceled {
		t.Errorf("expected state 'canceled', got %v", response.State)
	}
}

func TestJobsHandler_Cancel_NotFound(t *testing.T) {
	jobRepo := &testJobRepository{
		getFn: func(ctx context.Context, id string) (domain.Job, error) {
			return domain.Job{}, app.ErrNotFound
		},
		updateStateFn: func(ctx context.Context, id string, expected domain.JobState, next domain.JobState) (domain.Job, error) {
			return domain.Job{}, app.ErrNotFound
		},
	}

	handler := NewJobsHandler(app.NewJobService(jobRepo, nil))

	httpReq := httptest.NewRequest("POST", "/jobs/non-existent/cancel", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "non-existent")
	httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), chi.RouteCtxKey, rctx))

	handler.cancel(w, httpReq)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

// Mock JobService for testing
type mockJobService struct {
	createFn func(ctx context.Context, req app.CreateJobRequest) (app.JobDTO, error)
	getFn    func(ctx context.Context, id string) (app.JobDTO, error)
	listFn   func(ctx context.Context, limit int) ([]app.JobDTO, error)
	cancelFn func(ctx context.Context, id string) (app.JobDTO, error)
}

func (m *mockJobService) Create(ctx context.Context, req app.CreateJobRequest) (app.JobDTO, error) {
	if m.createFn != nil {
		return m.createFn(ctx, req)
	}
	return app.JobDTO{}, nil
}

func (m *mockJobService) Get(ctx context.Context, id string) (app.JobDTO, error) {
	if m.getFn != nil {
		return m.getFn(ctx, id)
	}
	return app.JobDTO{}, nil
}

func (m *mockJobService) List(ctx context.Context, limit int) ([]app.JobDTO, error) {
	if m.listFn != nil {
		return m.listFn(ctx, limit)
	}
	return []app.JobDTO{}, nil
}

func (m *mockJobService) Cancel(ctx context.Context, id string) (app.JobDTO, error) {
	if m.cancelFn != nil {
		return m.cancelFn(ctx, id)
	}
	return app.JobDTO{}, nil
}
