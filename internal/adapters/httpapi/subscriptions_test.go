package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/adapters/memorybus"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/app"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
	"github.com/go-chi/chi/v5"
)

// Test stubs for SubscriptionRepository
type stubSubscriptionRepo struct {
	createFn func(ctx context.Context, sub domain.Subscription) (domain.Subscription, error)
	getFn    func(ctx context.Context, id string) (domain.Subscription, error)
	listFn   func(ctx context.Context, limit int) ([]domain.Subscription, error)
	updateFn func(ctx context.Context, sub domain.Subscription) (domain.Subscription, error)
	deleteFn func(ctx context.Context, id string) error
	dueFn    func(ctx context.Context, now time.Time, limit int) ([]domain.Subscription, error)
	markFn   func(ctx context.Context, id string, episode int) (domain.Subscription, error)
}

func (r *stubSubscriptionRepo) Create(ctx context.Context, sub domain.Subscription) (domain.Subscription, error) {
	if r.createFn != nil {
		return r.createFn(ctx, sub)
	}
	return sub, nil
}

func (r *stubSubscriptionRepo) Get(ctx context.Context, id string) (domain.Subscription, error) {
	if r.getFn != nil {
		return r.getFn(ctx, id)
	}
	return domain.Subscription{}, nil
}

func (r *stubSubscriptionRepo) List(ctx context.Context, limit int) ([]domain.Subscription, error) {
	if r.listFn != nil {
		return r.listFn(ctx, limit)
	}
	return []domain.Subscription{}, nil
}

func (r *stubSubscriptionRepo) Update(ctx context.Context, sub domain.Subscription) (domain.Subscription, error) {
	if r.updateFn != nil {
		return r.updateFn(ctx, sub)
	}
	return sub, nil
}

func (r *stubSubscriptionRepo) Delete(ctx context.Context, id string) error {
	if r.deleteFn != nil {
		return r.deleteFn(ctx, id)
	}
	return nil
}

func (r *stubSubscriptionRepo) Due(ctx context.Context, now time.Time, limit int) ([]domain.Subscription, error) {
	if r.dueFn != nil {
		return r.dueFn(ctx, now, limit)
	}
	return []domain.Subscription{}, nil
}

func (r *stubSubscriptionRepo) MarkDownloadedEpisodeMax(ctx context.Context, id string, episode int) (domain.Subscription, error) {
	if r.markFn != nil {
		return r.markFn(ctx, id, episode)
	}
	return domain.Subscription{}, nil
}

// Tests
func TestSubscriptionsHandler_Create_Success(t *testing.T) {
	repo := &stubSubscriptionRepo{
		createFn: func(ctx context.Context, sub domain.Subscription) (domain.Subscription, error) {
			sub.ID = "sub-1"
			return sub, nil
		},
	}
	bus := memorybus.New()
	handler := NewSubscriptionsHandler(app.NewSubscriptionService(repo, nil, nil, bus))

	reqBody := createSubscriptionRequest{
		BaseURL: "https://example.com",
		Label:   "Test Anime",
	}
	body, _ := json.Marshal(reqBody)
	httpReq := httptest.NewRequest(http.MethodPost, "/subscriptions", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	router := chi.NewRouter()
	handler.Routes(router)
	router.ServeHTTP(rr, httpReq)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}
}

func TestSubscriptionsHandler_List_Success(t *testing.T) {
	repo := &stubSubscriptionRepo{
		listFn: func(ctx context.Context, limit int) ([]domain.Subscription, error) {
			return []domain.Subscription{
				{ID: "sub-1", BaseURL: "https://example.com", Label: "Test 1"},
				{ID: "sub-2", BaseURL: "https://example2.com", Label: "Test 2"},
			}, nil
		},
	}
	bus := memorybus.New()
	handler := NewSubscriptionsHandler(app.NewSubscriptionService(repo, nil, nil, bus))

	httpReq := httptest.NewRequest(http.MethodGet, "/subscriptions", nil)
	rr := httptest.NewRecorder()

	router := chi.NewRouter()
	handler.Routes(router)
	router.ServeHTTP(rr, httpReq)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestSubscriptionsHandler_Get_Success(t *testing.T) {
	repo := &stubSubscriptionRepo{
		getFn: func(ctx context.Context, id string) (domain.Subscription, error) {
			return domain.Subscription{ID: "sub-1", BaseURL: "https://example.com", Label: "Test"}, nil
		},
	}
	bus := memorybus.New()
	handler := NewSubscriptionsHandler(app.NewSubscriptionService(repo, nil, nil, bus))

	httpReq := httptest.NewRequest(http.MethodGet, "/subscriptions/sub-1", nil)
	rr := httptest.NewRecorder()

	router := chi.NewRouter()
	handler.Routes(router)
	router.ServeHTTP(rr, httpReq)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestSubscriptionsHandler_Delete_Success(t *testing.T) {
	repo := &stubSubscriptionRepo{
		deleteFn: func(ctx context.Context, id string) error {
			return nil
		},
	}
	bus := memorybus.New()
	handler := NewSubscriptionsHandler(app.NewSubscriptionService(repo, nil, nil, bus))

	httpReq := httptest.NewRequest(http.MethodDelete, "/subscriptions/sub-1", nil)
	rr := httptest.NewRecorder()

	router := chi.NewRouter()
	handler.Routes(router)
	router.ServeHTTP(rr, httpReq)

	if rr.Code != http.StatusNoContent {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, rr.Code)
	}
}

func TestSubscriptionsHandler_Create_InvalidJSON(t *testing.T) {
	repo := &stubSubscriptionRepo{}
	bus := memorybus.New()
	handler := NewSubscriptionsHandler(app.NewSubscriptionService(repo, nil, nil, bus))

	httpReq := httptest.NewRequest(http.MethodPost, "/subscriptions", bytes.NewReader([]byte("invalid")))
	rr := httptest.NewRecorder()

	router := chi.NewRouter()
	handler.Routes(router)
	router.ServeHTTP(rr, httpReq)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}
