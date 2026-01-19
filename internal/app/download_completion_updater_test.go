package app

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/ports"
	"github.com/rs/zerolog"
)

type memSubsRepo struct {
	mu   sync.Mutex
	byID map[string]domain.Subscription
}

func newMemSubsRepo() *memSubsRepo {
	return &memSubsRepo{byID: map[string]domain.Subscription{}}
}

func (r *memSubsRepo) Create(ctx context.Context, sub domain.Subscription) (domain.Subscription, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.byID[sub.ID] = sub
	return sub, nil
}

func (r *memSubsRepo) Get(ctx context.Context, id string) (domain.Subscription, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	sub, ok := r.byID[id]
	if !ok {
		return domain.Subscription{}, ports.ErrNotFound
	}
	return sub, nil
}

func (r *memSubsRepo) List(ctx context.Context, limit int) ([]domain.Subscription, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]domain.Subscription, 0, len(r.byID))
	for _, sub := range r.byID {
		out = append(out, sub)
	}
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out, nil
}

func (r *memSubsRepo) Update(ctx context.Context, sub domain.Subscription) (domain.Subscription, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.byID[sub.ID]; !ok {
		return domain.Subscription{}, ports.ErrNotFound
	}
	r.byID[sub.ID] = sub
	return sub, nil
}

func (r *memSubsRepo) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.byID[id]; !ok {
		return ports.ErrNotFound
	}
	delete(r.byID, id)
	return nil
}

func (r *memSubsRepo) Due(ctx context.Context, now time.Time, limit int) ([]domain.Subscription, error) {
	return nil, nil
}

func (r *memSubsRepo) MarkDownloadedEpisodeMax(ctx context.Context, id string, episode int) (domain.Subscription, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	sub, ok := r.byID[id]
	if !ok {
		return domain.Subscription{}, ports.ErrNotFound
	}
	if episode > sub.LastDownloadedEpisode {
		sub.LastDownloadedEpisode = episode
		r.byID[id] = sub
	}
	return sub, nil
}

func TestDownloadCompletionUpdater_MarksLastDownloadedEpisodeMax(t *testing.T) {
	repo := newMemSubsRepo()
	if _, err := repo.Create(context.Background(), domain.Subscription{ID: "sub1"}); err != nil {
		t.Fatalf("create subscription: %v", err)
	}

	u := NewDownloadCompletionUpdater(zerolog.Nop(), nil, repo)

	payload := JobDTO{
		ID:    "job1",
		Type:  "download",
		State: domain.JobCompleted,
		Params: func() json.RawMessage {
			b, _ := json.Marshal(map[string]any{"subscriptionId": "sub1", "episode": 3})
			return b
		}(),
	}
	b, _ := json.Marshal(payload)

	u.handleEvent(context.Background(), ports.Event{Topic: "job.completed", Payload: b})
	sub, _ := repo.Get(context.Background(), "sub1")
	if sub.LastDownloadedEpisode != 3 {
		t.Fatalf("expected lastDownloadedEpisode=3, got %d", sub.LastDownloadedEpisode)
	}

	// Lower episode should not reduce value.
	payload.Params = func() json.RawMessage {
		b, _ := json.Marshal(map[string]any{"subscriptionId": "sub1", "episode": 2})
		return b
	}()
	b, _ = json.Marshal(payload)
	u.handleEvent(context.Background(), ports.Event{Topic: "job.completed", Payload: b})
	sub, _ = repo.Get(context.Background(), "sub1")
	if sub.LastDownloadedEpisode != 3 {
		t.Fatalf("expected lastDownloadedEpisode to stay 3, got %d", sub.LastDownloadedEpisode)
	}
}
