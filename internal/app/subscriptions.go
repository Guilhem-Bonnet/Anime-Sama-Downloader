package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/ports"
	"github.com/rs/xid"
)

type SubscriptionService struct {
	repo ports.SubscriptionRepository
	jobs *JobService
	bus  ports.EventBus
}

func NewSubscriptionService(repo ports.SubscriptionRepository, jobs *JobService, bus ports.EventBus) *SubscriptionService {
	return &SubscriptionService{repo: repo, jobs: jobs, bus: bus}
}

type SubscriptionDTO struct {
	ID string `json:"id"`

	BaseURL string `json:"baseUrl"`
	Label   string `json:"label"`
	Player  string `json:"player"`

	LastScheduledEpisode  int `json:"lastScheduledEpisode"`
	LastDownloadedEpisode int `json:"lastDownloadedEpisode"`
	LastAvailableEpisode  int `json:"lastAvailableEpisode"`

	NextCheckAt   time.Time `json:"nextCheckAt"`
	LastCheckedAt time.Time `json:"lastCheckedAt"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func toSubscriptionDTO(s domain.Subscription) SubscriptionDTO {
	return SubscriptionDTO{
		ID: s.ID,
		BaseURL: s.BaseURL,
		Label: s.Label,
		Player: s.Player,
		LastScheduledEpisode: s.LastScheduledEpisode,
		LastDownloadedEpisode: s.LastDownloadedEpisode,
		LastAvailableEpisode: s.LastAvailableEpisode,
		NextCheckAt: s.NextCheckAt,
		LastCheckedAt: s.LastCheckedAt,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

func (s *SubscriptionService) Create(ctx context.Context, baseURL, label, player string) (SubscriptionDTO, error) {
	baseURL = strings.TrimSpace(baseURL)
	label = strings.TrimSpace(label)
	player = strings.TrimSpace(player)
	if baseURL == "" {
		return SubscriptionDTO{}, errors.New("missing baseUrl")
	}
	if label == "" {
		return SubscriptionDTO{}, errors.New("missing label")
	}
	if player == "" {
		player = "auto"
	}
	canon, err := CanonicalizeAnimeSamaBaseURL(baseURL)
	if err != nil {
		return SubscriptionDTO{}, err
	}

	now := time.Now().UTC()
	sub := domain.Subscription{
		ID: xid.New().String(),
		BaseURL: canon,
		Label: label,
		Player: player,
		LastScheduledEpisode: 0,
		LastDownloadedEpisode: 0,
		LastAvailableEpisode: 0,
		NextCheckAt: now,
		LastCheckedAt: time.Time{},
		CreatedAt: now,
		UpdatedAt: now,
	}
	created, err := s.repo.Create(ctx, sub)
	if err != nil {
		return SubscriptionDTO{}, err
	}
	s.publish("subscription.created", created)
	return toSubscriptionDTO(created), nil
}

func (s *SubscriptionService) Get(ctx context.Context, id string) (SubscriptionDTO, error) {
	sub, err := s.repo.Get(ctx, id)
	if err != nil {
		return SubscriptionDTO{}, err
	}
	return toSubscriptionDTO(sub), nil
}

func (s *SubscriptionService) List(ctx context.Context, limit int) ([]SubscriptionDTO, error) {
	subs, err := s.repo.List(ctx, limit)
	if err != nil {
		return nil, err
	}
	out := make([]SubscriptionDTO, 0, len(subs))
	for _, sub := range subs {
		out = append(out, toSubscriptionDTO(sub))
	}
	return out, nil
}

func (s *SubscriptionService) Update(ctx context.Context, dto SubscriptionDTO) (SubscriptionDTO, error) {
	existing, err := s.repo.Get(ctx, dto.ID)
	if err != nil {
		return SubscriptionDTO{}, err
	}
	if strings.TrimSpace(dto.BaseURL) != "" {
		canon, err := CanonicalizeAnimeSamaBaseURL(dto.BaseURL)
		if err != nil {
			return SubscriptionDTO{}, err
		}
		existing.BaseURL = canon
	}
	if strings.TrimSpace(dto.Label) != "" {
		existing.Label = strings.TrimSpace(dto.Label)
	}
	if strings.TrimSpace(dto.Player) != "" {
		existing.Player = strings.TrimSpace(dto.Player)
	}
	// Allow manually adjusting lastDownloadedEpisode (useful for initial bootstrap).
	if dto.LastDownloadedEpisode >= 0 {
		existing.LastDownloadedEpisode = dto.LastDownloadedEpisode
	}
	if dto.LastScheduledEpisode >= 0 {
		existing.LastScheduledEpisode = dto.LastScheduledEpisode
	}
	existing.UpdatedAt = time.Now().UTC()
	updated, err := s.repo.Update(ctx, existing)
	if err != nil {
		return SubscriptionDTO{}, err
	}
	s.publish("subscription.updated", updated)
	return toSubscriptionDTO(updated), nil
}

func (s *SubscriptionService) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err == nil {
		s.publishRaw("subscription.deleted", map[string]any{"id": id})
	}
	return err
}

type SyncResult struct {
	Subscription SubscriptionDTO `json:"subscription"`
	SelectedPlayer string `json:"selectedPlayer"`
	MaxAvailableEpisode int `json:"maxAvailableEpisode"`
	EnqueuedEpisodes []int `json:"enqueuedEpisodes"`
	EnqueuedJobIDs []string `json:"enqueuedJobIds"`
	Message string `json:"message,omitempty"`
}

// SyncOnce fetches episodes.js, updates availability fields, and optionally enqueues download jobs
// for newly-available episodes. This is a best-effort MVP: episode URLs are host/embed URLs.
func (s *SubscriptionService) SyncOnce(ctx context.Context, id string, enqueue bool) (SyncResult, error) {
	sub, err := s.repo.Get(ctx, id)
	if err != nil {
		return SyncResult{}, err
	}

	jsText, err := FetchEpisodesJS(ctx, sub.BaseURL)
	if err != nil {
		sub.LastCheckedAt = time.Now().UTC()
		sub.NextCheckAt = time.Now().UTC().Add(30 * time.Minute)
		sub.UpdatedAt = time.Now().UTC()
		_, _ = s.repo.Update(ctx, sub)
		return SyncResult{}, err
	}

	eps, err := ParseEpisodesJS(jsText)
	if err != nil {
		return SyncResult{}, err
	}

	selected := sub.Player
	if selected == "" || strings.EqualFold(selected, "auto") {
		selected = BestPlayer(eps.Players)
		if selected == "auto" {
			selected = ""
		}
	}
	urls := eps.Players[selected]
	if len(urls) == 0 {
		// fallback to best
		selected = BestPlayer(eps.Players)
		urls = eps.Players[selected]
	}

	maxAvail := MaxAvailableEpisode(urls)
	if maxAvail < 0 {
		maxAvail = 0
	}

	now := time.Now().UTC()
	sub.LastAvailableEpisode = maxAvail
	sub.LastCheckedAt = now

	// Basic next-check policy: if new stuff is available beyond what we've scheduled, check more often.
	if sub.LastScheduledEpisode < maxAvail {
		sub.NextCheckAt = now.Add(10 * time.Minute)
	} else {
		sub.NextCheckAt = now.Add(2 * time.Hour)
	}

	enqueuedEpisodes := []int{}
	enqueuedJobIDs := []string{}
	if enqueue && s.jobs != nil && sub.LastScheduledEpisode < maxAvail {
		from := sub.LastScheduledEpisode + 1
		for ep := from; ep <= maxAvail; ep++ {
			if ep-1 < 0 || ep-1 >= len(urls) {
				continue
			}
			u := urls[ep-1]
			if strings.TrimSpace(u) == "" {
				continue
			}

			params := map[string]any{
				"url": u,
				"path": filepath.ToSlash(filepath.Join("subscriptions", sub.ID, fmt.Sprintf("%s-ep-%02d.mp4", safeLabel(sub.Label), ep))),
				"filename": "",
				"subscriptionId": sub.ID,
				"episode": ep,
				"source": "anime-sama",
			}
			b, _ := json.Marshal(params)
			created, err := s.jobs.Create(ctx, CreateJobRequest{Type: "download", Params: b})
			if err != nil {
				// stop on first enqueue error
				break
			}
			enqueuedEpisodes = append(enqueuedEpisodes, ep)
			enqueuedJobIDs = append(enqueuedJobIDs, created.ID)
			sub.LastScheduledEpisode = ep
		}
	}

	sub.UpdatedAt = time.Now().UTC()
	updated, err := s.repo.Update(ctx, sub)
	if err != nil {
		return SyncResult{}, err
	}
	s.publish("subscription.synced", updated)

	return SyncResult{
		Subscription: toSubscriptionDTO(updated),
		SelectedPlayer: selected,
		MaxAvailableEpisode: maxAvail,
		EnqueuedEpisodes: enqueuedEpisodes,
		EnqueuedJobIDs: enqueuedJobIDs,
		Message: "note: episodes.js urls are host/embed urls; full video extraction pipeline is not implemented yet",
	}, nil
}

func safeLabel(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "/", "-")
	s = strings.ReplaceAll(s, "\\", "-")
	s = strings.ReplaceAll(s, "\x00", "")
	if s == "" {
		return "anime"
	}
	return s
}

func (s *SubscriptionService) publish(topic string, sub domain.Subscription) {
	if s.bus == nil {
		return
	}
	b, err := json.Marshal(toSubscriptionDTO(sub))
	if err != nil {
		return
	}
	s.bus.Publish(topic, b)
}

func (s *SubscriptionService) publishRaw(topic string, v any) {
	if s.bus == nil {
		return
	}
	b, err := json.Marshal(v)
	if err != nil {
		return
	}
	s.bus.Publish(topic, b)
}
