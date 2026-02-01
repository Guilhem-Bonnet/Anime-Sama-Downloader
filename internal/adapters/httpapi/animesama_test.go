package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/app"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
	"github.com/go-chi/chi/v5"
)

// Stub AnimeSamaResolver for testing
type stubAnimeSamaResolver struct {
	resolveFn func(ctx context.Context, titles []string, maxCandidates int) ([]app.AnimeSamaCandidate, error)
}

func (s *stubAnimeSamaResolver) ResolveCandidates(ctx context.Context, titles []string, maxCandidates int) ([]app.AnimeSamaCandidate, error) {
	if s.resolveFn != nil {
		return s.resolveFn(ctx, titles, maxCandidates)
	}
	return []app.AnimeSamaCandidate{}, nil
}

// Stub JobRepository for testing
type stubJobRepo struct {
	createFn func(ctx context.Context, job domain.Job) (domain.Job, error)
}

func (s *stubJobRepo) Create(ctx context.Context, job domain.Job) (domain.Job, error) {
	if s.createFn != nil {
		return s.createFn(ctx, job)
	}
	job.ID = "job-1"
	return job, nil
}

func (s *stubJobRepo) Get(ctx context.Context, id string) (domain.Job, error) {
	return domain.Job{}, nil
}

func (s *stubJobRepo) List(ctx context.Context, limit int) ([]domain.Job, error) {
	return []domain.Job{}, nil
}

func (s *stubJobRepo) ClaimNextQueued(ctx context.Context) (domain.Job, error) {
	return domain.Job{}, nil
}

func (s *stubJobRepo) UpdateProgress(ctx context.Context, id string, progress float64) (domain.Job, error) {
	return domain.Job{}, nil
}

func (s *stubJobRepo) UpdateResult(ctx context.Context, id string, resultJSON []byte) (domain.Job, error) {
	return domain.Job{}, nil
}

func (s *stubJobRepo) UpdateError(ctx context.Context, id string, code string, message string) (domain.Job, error) {
	return domain.Job{}, nil
}

func (s *stubJobRepo) UpdateState(ctx context.Context, id string, expected domain.JobState, next domain.JobState) (domain.Job, error) {
	return domain.Job{}, nil
}

// Tests
func TestAnimeSamaHandler_Resolve_Success(t *testing.T) {
	resolver := &stubAnimeSamaResolver{
		resolveFn: func(ctx context.Context, titles []string, maxCandidates int) ([]app.AnimeSamaCandidate, error) {
			return []app.AnimeSamaCandidate{
				{
					CatalogueURL: "https://anime-sama.si/catalogue/test",
					MatchedTitle: "Test Anime",
					Score:        95.0,
				},
			}, nil
		},
	}
	jobRepo := &stubJobRepo{}
	jobService := app.NewJobService(jobRepo, nil)
	handler := NewAnimeSamaHandler(resolver, jobService)

	reqBody := animeSamaResolveRequest{
		Titles:        []string{"Test Anime"},
		Season:        1,
		Lang:          "vostfr",
		MaxCandidates: 5,
	}
	body, _ := json.Marshal(reqBody)
	httpReq := httptest.NewRequest(http.MethodPost, "/animesama/resolve", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	router := chi.NewRouter()
	handler.Routes(router)
	router.ServeHTTP(rr, httpReq)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var response animeSamaResolveResponse
	json.Unmarshal(rr.Body.Bytes(), &response)

	if len(response.Candidates) == 0 {
		t.Errorf("expected candidates, got %d", len(response.Candidates))
	}
}

func TestAnimeSamaHandler_Resolve_InvalidJSON(t *testing.T) {
	resolver := &stubAnimeSamaResolver{}
	jobRepo := &stubJobRepo{}
	jobService := app.NewJobService(jobRepo, nil)
	handler := NewAnimeSamaHandler(resolver, jobService)

	httpReq := httptest.NewRequest(http.MethodPost, "/animesama/resolve", bytes.NewReader([]byte("invalid json")))
	rr := httptest.NewRecorder()

	router := chi.NewRouter()
	handler.Routes(router)
	router.ServeHTTP(rr, httpReq)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestAnimeSamaHandler_Episodes_InvalidJSON(t *testing.T) {
	resolver := &stubAnimeSamaResolver{}
	jobRepo := &stubJobRepo{}
	jobService := app.NewJobService(jobRepo, nil)
	handler := NewAnimeSamaHandler(resolver, jobService)

	httpReq := httptest.NewRequest(http.MethodPost, "/animesama/episodes", bytes.NewReader([]byte("invalid")))
	rr := httptest.NewRecorder()

	router := chi.NewRouter()
	handler.Routes(router)
	router.ServeHTTP(rr, httpReq)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestAnimeSamaHandler_Scan_InvalidJSON(t *testing.T) {
	resolver := &stubAnimeSamaResolver{}
	jobRepo := &stubJobRepo{}
	jobService := app.NewJobService(jobRepo, nil)
	handler := NewAnimeSamaHandler(resolver, jobService)

	httpReq := httptest.NewRequest(http.MethodPost, "/animesama/scan", bytes.NewReader([]byte("bad json")))
	rr := httptest.NewRecorder()

	router := chi.NewRouter()
	handler.Routes(router)
	router.ServeHTTP(rr, httpReq)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestAnimeSamaHandler_Enqueue_InvalidJSON(t *testing.T) {
	resolver := &stubAnimeSamaResolver{}
	jobRepo := &stubJobRepo{}
	jobService := app.NewJobService(jobRepo, nil)
	handler := NewAnimeSamaHandler(resolver, jobService)

	httpReq := httptest.NewRequest(http.MethodPost, "/animesama/enqueue", bytes.NewReader([]byte("bad")))
	rr := httptest.NewRecorder()

	router := chi.NewRouter()
	handler.Routes(router)
	router.ServeHTTP(rr, httpReq)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}
