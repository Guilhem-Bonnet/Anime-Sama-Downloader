package app

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/ports"
)

// MockSubscriptionRepository for testing
type mockSubscriptionRepository struct {
	createFn                   func(ctx context.Context, sub domain.Subscription) (domain.Subscription, error)
	getFn                      func(ctx context.Context, id string) (domain.Subscription, error)
	listFn                     func(ctx context.Context, limit int) ([]domain.Subscription, error)
	updateFn                   func(ctx context.Context, sub domain.Subscription) (domain.Subscription, error)
	deleteFn                   func(ctx context.Context, id string) error
	dueFn                      func(ctx context.Context, now time.Time, limit int) ([]domain.Subscription, error)
	markDownloadedEpisodeMaxFn func(ctx context.Context, id string, episode int) (domain.Subscription, error)
}

func (m *mockSubscriptionRepository) Create(ctx context.Context, sub domain.Subscription) (domain.Subscription, error) {
	if m.createFn != nil {
		return m.createFn(ctx, sub)
	}
	return sub, nil
}

func (m *mockSubscriptionRepository) Get(ctx context.Context, id string) (domain.Subscription, error) {
	if m.getFn != nil {
		return m.getFn(ctx, id)
	}
	return domain.Subscription{}, nil
}

func (m *mockSubscriptionRepository) List(ctx context.Context, limit int) ([]domain.Subscription, error) {
	if m.listFn != nil {
		return m.listFn(ctx, limit)
	}
	return []domain.Subscription{}, nil
}

func (m *mockSubscriptionRepository) Update(ctx context.Context, sub domain.Subscription) (domain.Subscription, error) {
	if m.updateFn != nil {
		return m.updateFn(ctx, sub)
	}
	return sub, nil
}

func (m *mockSubscriptionRepository) Delete(ctx context.Context, id string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}

func (m *mockSubscriptionRepository) Due(ctx context.Context, now time.Time, limit int) ([]domain.Subscription, error) {
	if m.dueFn != nil {
		return m.dueFn(ctx, now, limit)
	}
	return []domain.Subscription{}, nil
}

func (m *mockSubscriptionRepository) MarkDownloadedEpisodeMax(ctx context.Context, id string, episode int) (domain.Subscription, error) {
	if m.markDownloadedEpisodeMaxFn != nil {
		return m.markDownloadedEpisodeMaxFn(ctx, id, episode)
	}
	return domain.Subscription{}, nil
}

// MockEventBus for testing
type mockEventBus struct {
	publishFn   func(topic string, data []byte)
	subscribeFn func() (<-chan ports.Event, func())
}

func (m *mockEventBus) Publish(topic string, data []byte) {
	if m.publishFn != nil {
		m.publishFn(topic, data)
	}
}

func (m *mockEventBus) Subscribe() (<-chan ports.Event, func()) {
	if m.subscribeFn != nil {
		return m.subscribeFn()
	}
	ch := make(chan ports.Event)
	return ch, func() {}
}

// Test cases for SubscriptionService.Create()
func TestSubscriptionService_Create_Success(t *testing.T) {
	repo := &mockSubscriptionRepository{
		createFn: func(ctx context.Context, sub domain.Subscription) (domain.Subscription, error) {
			return sub, nil
		},
	}
	bus := &mockEventBus{}

	service := NewSubscriptionService(repo, nil, nil, bus)
	ctx := context.Background()

	dto, err := service.Create(ctx, "https://anime-sama.fr/catalogue/naruto", "Naruto", "auto")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dto.ID == "" {
		t.Error("expected ID to be set")
	}
	if dto.BaseURL == "" {
		t.Error("expected BaseURL to be set")
	}
	if dto.Label != "Naruto" {
		t.Errorf("expected label 'Naruto', got %q", dto.Label)
	}
	if dto.Player != "auto" {
		t.Errorf("expected player 'auto', got %q", dto.Player)
	}
}

func TestSubscriptionService_Create_MissingBaseURL(t *testing.T) {
	repo := &mockSubscriptionRepository{}
	service := NewSubscriptionService(repo, nil, nil, nil)
	ctx := context.Background()

	_, err := service.Create(ctx, "", "Naruto", "auto")
	if err == nil {
		t.Fatal("expected error for missing baseURL")
	}
	if !errors.Is(err, errors.New("missing baseUrl")) && err.Error() != "missing baseUrl" {
		t.Errorf("expected 'missing baseUrl' error, got %v", err)
	}
}

func TestSubscriptionService_Create_InvalidBaseURL(t *testing.T) {
	repo := &mockSubscriptionRepository{}
	service := NewSubscriptionService(repo, nil, nil, nil)
	ctx := context.Background()

	_, err := service.Create(ctx, "not a valid url", "", "auto")
	if err == nil {
		t.Fatal("expected error for invalid baseURL")
	}
}

func TestSubscriptionService_Create_DefaultPlayer(t *testing.T) {
	repo := &mockSubscriptionRepository{
		createFn: func(ctx context.Context, sub domain.Subscription) (domain.Subscription, error) {
			return sub, nil
		},
	}
	service := NewSubscriptionService(repo, nil, nil, nil)
	ctx := context.Background()

	dto, err := service.Create(ctx, "https://anime-sama.fr/catalogue/naruto", "Naruto", "")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dto.Player != "auto" {
		t.Errorf("expected default player 'auto', got %q", dto.Player)
	}
}

func TestSubscriptionService_Create_AutoLabel(t *testing.T) {
	repo := &mockSubscriptionRepository{
		createFn: func(ctx context.Context, sub domain.Subscription) (domain.Subscription, error) {
			return sub, nil
		},
	}
	service := NewSubscriptionService(repo, nil, nil, nil)
	ctx := context.Background()

	dto, err := service.Create(ctx, "https://anime-sama.fr/catalogue/naruto", "", "auto")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dto.Label == "" {
		t.Error("expected label to be auto-generated")
	}
}

// Test cases for SubscriptionService.Update()
func TestSubscriptionService_Update_Success(t *testing.T) {
	existing := domain.Subscription{
		ID:      "sub-1",
		BaseURL: "https://anime-sama.fr/catalogue/naruto",
		Label:   "Naruto",
		Player:  "auto",
	}

	repo := &mockSubscriptionRepository{
		getFn: func(ctx context.Context, id string) (domain.Subscription, error) {
			if id == "sub-1" {
				return existing, nil
			}
			return domain.Subscription{}, errors.New("not found")
		},
		updateFn: func(ctx context.Context, sub domain.Subscription) (domain.Subscription, error) {
			return sub, nil
		},
	}

	service := NewSubscriptionService(repo, nil, nil, nil)
	ctx := context.Background()

	updateDTO := SubscriptionDTO{
		ID:    "sub-1",
		Label: "Updated Label",
	}

	result, err := service.Update(ctx, updateDTO)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Label != "Updated Label" {
		t.Errorf("expected label 'Updated Label', got %q", result.Label)
	}
}

func TestSubscriptionService_Update_NotFound(t *testing.T) {
	repo := &mockSubscriptionRepository{
		getFn: func(ctx context.Context, id string) (domain.Subscription, error) {
			return domain.Subscription{}, errors.New("not found")
		},
	}

	service := NewSubscriptionService(repo, nil, nil, nil)
	ctx := context.Background()

	updateDTO := SubscriptionDTO{ID: "non-existent"}
	_, err := service.Update(ctx, updateDTO)

	if err == nil {
		t.Fatal("expected error for non-existent subscription")
	}
}

func TestSubscriptionService_Update_ManualEpisodeAdjustment(t *testing.T) {
	existing := domain.Subscription{
		ID:                    "sub-1",
		BaseURL:               "https://anime-sama.fr/catalogue/naruto",
		Label:                 "Naruto",
		LastDownloadedEpisode: 0,
	}

	repo := &mockSubscriptionRepository{
		getFn: func(ctx context.Context, id string) (domain.Subscription, error) {
			return existing, nil
		},
		updateFn: func(ctx context.Context, sub domain.Subscription) (domain.Subscription, error) {
			return sub, nil
		},
	}

	service := NewSubscriptionService(repo, nil, nil, nil)
	ctx := context.Background()

	updateDTO := SubscriptionDTO{
		ID:                    "sub-1",
		LastDownloadedEpisode: 25,
	}

	result, err := service.Update(ctx, updateDTO)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.LastDownloadedEpisode != 25 {
		t.Errorf("expected LastDownloadedEpisode 25, got %d", result.LastDownloadedEpisode)
	}
}

// Test cases for SubscriptionService.Delete()
func TestSubscriptionService_Delete_Success(t *testing.T) {
	deleteCalled := false
	repo := &mockSubscriptionRepository{
		deleteFn: func(ctx context.Context, id string) error {
			deleteCalled = true
			return nil
		},
	}

	bus := &mockEventBus{}
	service := NewSubscriptionService(repo, nil, nil, bus)
	ctx := context.Background()

	err := service.Delete(ctx, "sub-1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !deleteCalled {
		t.Error("expected Delete to be called on repository")
	}
}

func TestSubscriptionService_Delete_Error(t *testing.T) {
	repo := &mockSubscriptionRepository{
		deleteFn: func(ctx context.Context, id string) error {
			return errors.New("database error")
		},
	}

	service := NewSubscriptionService(repo, nil, nil, nil)
	ctx := context.Background()

	err := service.Delete(ctx, "sub-1")

	if err == nil {
		t.Fatal("expected error from Delete")
	}
}

// Test cases for SubscriptionService.List()
func TestSubscriptionService_List_Success(t *testing.T) {
	subs := []domain.Subscription{
		{ID: "sub-1", Label: "Naruto"},
		{ID: "sub-2", Label: "Bleach"},
	}

	repo := &mockSubscriptionRepository{
		listFn: func(ctx context.Context, limit int) ([]domain.Subscription, error) {
			return subs, nil
		},
	}

	service := NewSubscriptionService(repo, nil, nil, nil)
	ctx := context.Background()

	result, err := service.List(ctx, 10)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 subscriptions, got %d", len(result))
	}
	if result[0].Label != "Naruto" {
		t.Errorf("expected first label 'Naruto', got %q", result[0].Label)
	}
}

func TestSubscriptionService_List_Empty(t *testing.T) {
	repo := &mockSubscriptionRepository{
		listFn: func(ctx context.Context, limit int) ([]domain.Subscription, error) {
			return []domain.Subscription{}, nil
		},
	}

	service := NewSubscriptionService(repo, nil, nil, nil)
	ctx := context.Background()

	result, err := service.List(ctx, 10)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected 0 subscriptions, got %d", len(result))
	}
}

// Test cases for SubscriptionService.Get()
func TestSubscriptionService_Get_Success(t *testing.T) {
	sub := domain.Subscription{
		ID:    "sub-1",
		Label: "Naruto",
	}

	repo := &mockSubscriptionRepository{
		getFn: func(ctx context.Context, id string) (domain.Subscription, error) {
			if id == "sub-1" {
				return sub, nil
			}
			return domain.Subscription{}, errors.New("not found")
		},
	}

	service := NewSubscriptionService(repo, nil, nil, nil)
	ctx := context.Background()

	result, err := service.Get(ctx, "sub-1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Label != "Naruto" {
		t.Errorf("expected label 'Naruto', got %q", result.Label)
	}
}

func TestSubscriptionService_Get_NotFound(t *testing.T) {
	repo := &mockSubscriptionRepository{
		getFn: func(ctx context.Context, id string) (domain.Subscription, error) {
			return domain.Subscription{}, errors.New("not found")
		},
	}

	service := NewSubscriptionService(repo, nil, nil, nil)
	ctx := context.Background()

	_, err := service.Get(ctx, "non-existent")

	if err == nil {
		t.Fatal("expected error for non-existent subscription")
	}
}

// Test EventBus integration
func TestSubscriptionService_PublishesCreatedEvent(t *testing.T) {
	eventPublished := false
	var publishedTopic string

	repo := &mockSubscriptionRepository{
		createFn: func(ctx context.Context, sub domain.Subscription) (domain.Subscription, error) {
			return sub, nil
		},
	}

	bus := &mockEventBus{
		publishFn: func(topic string, data []byte) {
			eventPublished = true
			publishedTopic = topic
		},
	}

	service := NewSubscriptionService(repo, nil, nil, bus)
	ctx := context.Background()

	service.Create(ctx, "https://anime-sama.fr/catalogue/naruto", "Naruto", "auto")

	if !eventPublished {
		t.Error("expected event to be published")
	}
	if publishedTopic != "subscription.created" {
		t.Errorf("expected topic 'subscription.created', got %q", publishedTopic)
	}
}

// Test DTO conversion
func TestSubscriptionDTO_Conversion(t *testing.T) {
	now := time.Now().UTC()
	sub := domain.Subscription{
		ID:                    "sub-1",
		BaseURL:               "https://anime-sama.fr/catalogue/naruto",
		Label:                 "Naruto",
		Player:                "auto",
		LastScheduledEpisode:  10,
		LastDownloadedEpisode: 5,
		LastAvailableEpisode:  15,
		NextCheckAt:           now,
		LastCheckedAt:         now,
		CreatedAt:             now,
		UpdatedAt:             now,
	}

	dto := toSubscriptionDTO(sub)

	if dto.ID != sub.ID {
		t.Errorf("ID mismatch: expected %q, got %q", sub.ID, dto.ID)
	}
	if dto.Label != sub.Label {
		t.Errorf("Label mismatch: expected %q, got %q", sub.Label, dto.Label)
	}
	if dto.LastDownloadedEpisode != 5 {
		t.Errorf("LastDownloadedEpisode mismatch: expected 5, got %d", dto.LastDownloadedEpisode)
	}
}

// Test input trimming and validation
func TestSubscriptionService_Create_TrimsWhitespace(t *testing.T) {
	repo := &mockSubscriptionRepository{
		createFn: func(ctx context.Context, sub domain.Subscription) (domain.Subscription, error) {
			return sub, nil
		},
	}

	service := NewSubscriptionService(repo, nil, nil, nil)
	ctx := context.Background()

	dto, err := service.Create(ctx, "  https://anime-sama.fr/catalogue/naruto  ", "  Naruto  ", "  auto  ")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dto.Label != "Naruto" {
		t.Errorf("expected label to be trimmed, got %q", dto.Label)
	}
	if dto.Player != "auto" {
		t.Errorf("expected player to be trimmed, got %q", dto.Player)
	}
}
