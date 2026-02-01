package app

import (
	"context"
	"testing"
	"time"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
)

type stubSubRepo struct {
	created domain.Subscription
}

func (s *stubSubRepo) Create(ctx context.Context, sub domain.Subscription) (domain.Subscription, error) {
	s.created = sub
	return sub, nil
}
func (s *stubSubRepo) Get(ctx context.Context, id string) (domain.Subscription, error) {
	panic("not used")
}
func (s *stubSubRepo) List(ctx context.Context, limit int) ([]domain.Subscription, error) {
	panic("not used")
}
func (s *stubSubRepo) Update(ctx context.Context, sub domain.Subscription) (domain.Subscription, error) {
	panic("not used")
}
func (s *stubSubRepo) Delete(ctx context.Context, id string) error {
	panic("not used")
}
func (s *stubSubRepo) Due(ctx context.Context, now time.Time, limit int) ([]domain.Subscription, error) {
	panic("not used")
}
func (s *stubSubRepo) MarkDownloadedEpisodeMax(ctx context.Context, id string, episode int) (domain.Subscription, error) {
	panic("not used")
}

func TestSubscriptionService_Create_AutoLabelFromBaseURL(t *testing.T) {
	repo := &stubSubRepo{}
	svc := NewSubscriptionService(repo, nil, nil)

	ctx := context.Background()
	dto, err := svc.Create(ctx, "https://anime-sama.si/catalogue/solo-leveling/saison1/vostfr/", "", "auto")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if dto.Label == "" {
		t.Fatalf("expected non-empty label")
	}
	if dto.Label != "Solo Leveling (S1 VOSTFR)" {
		t.Fatalf("unexpected label: %q", dto.Label)
	}
}

func TestSubscriptionService_Create_CustomLabel(t *testing.T) {
	repo := &stubSubRepo{}
	svc := NewSubscriptionService(repo, nil, nil)

	ctx := context.Background()
	customLabel := "My Custom Anime"
	dto, err := svc.Create(ctx, "https://anime-sama.si/catalogue/test/saison1/vostfr/", customLabel, "auto")
	if err != nil {
		t.Fatalf("Create with custom label: %v", err)
	}
	if dto.Label != customLabel {
		t.Fatalf("expected label %q, got %q", customLabel, dto.Label)
	}
}

func TestSubscriptionService_Get(t *testing.T) {
	sub := domain.Subscription{
		ID:      "sub-1",
		BaseURL: "https://anime-sama.si/catalogue/test/saison1/vostfr/",
		Label:   "Test Anime",
		Player:  "auto",
	}
	repo := &getSubRepo{sub: sub}
	svc := NewSubscriptionService(repo, nil, nil)

	ctx := context.Background()
	dto, err := svc.Get(ctx, "sub-1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if dto.ID != "sub-1" {
		t.Fatalf("expected ID sub-1, got %s", dto.ID)
	}
	if dto.Label != "Test Anime" {
		t.Fatalf("expected label 'Test Anime', got %s", dto.Label)
	}
}

func TestSubscriptionService_List(t *testing.T) {
	subs := []domain.Subscription{
		{ID: "sub-1", Label: "Anime 1", BaseURL: "url1", Player: "auto"},
		{ID: "sub-2", Label: "Anime 2", BaseURL: "url2", Player: "auto"},
	}
	repo := &listSubRepo{subs: subs}
	svc := NewSubscriptionService(repo, nil, nil)

	ctx := context.Background()
	dtos, err := svc.List(ctx, 10)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(dtos) != 2 {
		t.Fatalf("expected 2 subscriptions, got %d", len(dtos))
	}
	if dtos[0].ID != "sub-1" {
		t.Fatalf("expected first sub-1, got %s", dtos[0].ID)
	}
}

func TestSubscriptionService_Delete(t *testing.T) {
	repo := &deleteSubRepo{deleted: false}
	svc := NewSubscriptionService(repo, nil, nil)

	ctx := context.Background()
	err := svc.Delete(ctx, "sub-1")
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if !repo.deleted {
		t.Fatalf("Delete should have marked deleted=true")
	}
}

// Additional stub repos for extended tests
type getSubRepo struct {
	sub domain.Subscription
}

func (r *getSubRepo) Create(ctx context.Context, sub domain.Subscription) (domain.Subscription, error) {
	panic("not used")
}
func (r *getSubRepo) Get(ctx context.Context, id string) (domain.Subscription, error) {
	return r.sub, nil
}
func (r *getSubRepo) List(ctx context.Context, limit int) ([]domain.Subscription, error) {
	panic("not used")
}
func (r *getSubRepo) Update(ctx context.Context, sub domain.Subscription) (domain.Subscription, error) {
	panic("not used")
}
func (r *getSubRepo) Delete(ctx context.Context, id string) error {
	panic("not used")
}
func (r *getSubRepo) Due(ctx context.Context, now time.Time, limit int) ([]domain.Subscription, error) {
	panic("not used")
}
func (r *getSubRepo) MarkDownloadedEpisodeMax(ctx context.Context, id string, episode int) (domain.Subscription, error) {
	panic("not used")
}

type listSubRepo struct {
	subs []domain.Subscription
}

func (r *listSubRepo) Create(ctx context.Context, sub domain.Subscription) (domain.Subscription, error) {
	panic("not used")
}
func (r *listSubRepo) Get(ctx context.Context, id string) (domain.Subscription, error) {
	panic("not used")
}
func (r *listSubRepo) List(ctx context.Context, limit int) ([]domain.Subscription, error) {
	return r.subs, nil
}
func (r *listSubRepo) Update(ctx context.Context, sub domain.Subscription) (domain.Subscription, error) {
	panic("not used")
}
func (r *listSubRepo) Delete(ctx context.Context, id string) error {
	panic("not used")
}
func (r *listSubRepo) Due(ctx context.Context, now time.Time, limit int) ([]domain.Subscription, error) {
	panic("not used")
}
func (r *listSubRepo) MarkDownloadedEpisodeMax(ctx context.Context, id string, episode int) (domain.Subscription, error) {
	panic("not used")
}

type deleteSubRepo struct {
	deleted bool
}

func (r *deleteSubRepo) Create(ctx context.Context, sub domain.Subscription) (domain.Subscription, error) {
	panic("not used")
}
func (r *deleteSubRepo) Get(ctx context.Context, id string) (domain.Subscription, error) {
	panic("not used")
}
func (r *deleteSubRepo) List(ctx context.Context, limit int) ([]domain.Subscription, error) {
	panic("not used")
}
func (r *deleteSubRepo) Update(ctx context.Context, sub domain.Subscription) (domain.Subscription, error) {
	panic("not used")
}
func (r *deleteSubRepo) Delete(ctx context.Context, id string) error {
	r.deleted = true
	return nil
}
func (r *deleteSubRepo) Due(ctx context.Context, now time.Time, limit int) ([]domain.Subscription, error) {
	panic("not used")
}
func (r *deleteSubRepo) MarkDownloadedEpisodeMax(ctx context.Context, id string, episode int) (domain.Subscription, error) {
	panic("not used")
}
