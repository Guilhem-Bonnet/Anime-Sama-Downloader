package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

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
		ID:                    s.ID,
		BaseURL:               s.BaseURL,
		Label:                 s.Label,
		Player:                s.Player,
		LastScheduledEpisode:  s.LastScheduledEpisode,
		LastDownloadedEpisode: s.LastDownloadedEpisode,
		LastAvailableEpisode:  s.LastAvailableEpisode,
		NextCheckAt:           s.NextCheckAt,
		LastCheckedAt:         s.LastCheckedAt,
		CreatedAt:             s.CreatedAt,
		UpdatedAt:             s.UpdatedAt,
	}
}

func (s *SubscriptionService) Create(ctx context.Context, baseURL, label, player string) (SubscriptionDTO, error) {
	baseURL = strings.TrimSpace(baseURL)
	label = strings.TrimSpace(label)
	player = strings.TrimSpace(player)
	if baseURL == "" {
		return SubscriptionDTO{}, errors.New("missing baseUrl")
	}
	if player == "" {
		player = "auto"
	}
	canon, err := CanonicalizeAnimeSamaBaseURL(baseURL)
	if err != nil {
		return SubscriptionDTO{}, err
	}
	if label == "" {
		label = defaultLabelForBaseURL(canon)
		if label == "" {
			label = "Anime"
		}
	}

	now := time.Now().UTC()
	sub := domain.Subscription{
		ID:                    xid.New().String(),
		BaseURL:               canon,
		Label:                 label,
		Player:                player,
		LastScheduledEpisode:  0,
		LastDownloadedEpisode: 0,
		LastAvailableEpisode:  0,
		NextCheckAt:           now,
		LastCheckedAt:         time.Time{},
		CreatedAt:             now,
		UpdatedAt:             now,
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
	Subscription        SubscriptionDTO `json:"subscription"`
	SelectedPlayer      string          `json:"selectedPlayer"`
	MaxAvailableEpisode int             `json:"maxAvailableEpisode"`
	EnqueuedEpisodes    []int           `json:"enqueuedEpisodes"`
	EnqueuedJobIDs      []string        `json:"enqueuedJobIds"`
	Message             string          `json:"message,omitempty"`
}

type EpisodeStatus struct {
	Episode    int  `json:"episode"`
	Available  bool `json:"available"`
	Scheduled  bool `json:"scheduled"`
	Downloaded bool `json:"downloaded"`
}

type EpisodesResponse struct {
	Subscription        SubscriptionDTO `json:"subscription"`
	SelectedPlayer      string          `json:"selectedPlayer"`
	MaxAvailableEpisode int             `json:"maxAvailableEpisode"`
	Episodes            []EpisodeStatus `json:"episodes"`
}

type EnqueueSkippedEpisode struct {
	Episode int    `json:"episode"`
	Reason  string `json:"reason"`
}

type EnqueueEpisodesResponse struct {
	Subscription     SubscriptionDTO         `json:"subscription"`
	SelectedPlayer   string                  `json:"selectedPlayer"`
	EnqueuedEpisodes []int                   `json:"enqueuedEpisodes"`
	EnqueuedJobIDs   []string                `json:"enqueuedJobIds"`
	Skipped          []EnqueueSkippedEpisode `json:"skipped"`
}

func (s *SubscriptionService) Episodes(ctx context.Context, id string) (EpisodesResponse, error) {
	sub, err := s.repo.Get(ctx, id)
	if err != nil {
		return EpisodesResponse{}, err
	}

	jsText, err := FetchEpisodesJS(ctx, sub.BaseURL)
	if err != nil {
		return EpisodesResponse{}, err
	}
	eps, err := ParseEpisodesJS(jsText)
	if err != nil {
		return EpisodesResponse{}, err
	}

	selected, urls := selectPlayer(sub.Player, eps.Players)
	maxAvail := MaxAvailableEpisode(urls)
	if maxAvail < 0 {
		maxAvail = 0
	}

	out := make([]EpisodeStatus, 0, maxAvail)
	for ep := 1; ep <= maxAvail; ep++ {
		available := false
		if ep-1 >= 0 && ep-1 < len(urls) {
			available = strings.TrimSpace(urls[ep-1]) != ""
		}
		out = append(out, EpisodeStatus{
			Episode:    ep,
			Available:  available,
			Scheduled:  sub.LastScheduledEpisode >= ep,
			Downloaded: sub.LastDownloadedEpisode >= ep,
		})
	}

	return EpisodesResponse{
		Subscription:        toSubscriptionDTO(sub),
		SelectedPlayer:      selected,
		MaxAvailableEpisode: maxAvail,
		Episodes:            out,
	}, nil
}

func (s *SubscriptionService) EnqueueEpisodes(ctx context.Context, id string, episodes []int) (EnqueueEpisodesResponse, error) {
	if s.jobs == nil {
		return EnqueueEpisodesResponse{}, errors.New("job service not configured")
	}
	if len(episodes) == 0 {
		return EnqueueEpisodesResponse{}, errors.New("missing episodes")
	}

	// Normalize input.
	seen := map[int]struct{}{}
	norm := make([]int, 0, len(episodes))
	for _, ep := range episodes {
		if ep <= 0 {
			continue
		}
		if _, ok := seen[ep]; ok {
			continue
		}
		seen[ep] = struct{}{}
		norm = append(norm, ep)
	}
	if len(norm) == 0 {
		return EnqueueEpisodesResponse{}, errors.New("no valid episodes")
	}
	sort.Ints(norm)

	sub, err := s.repo.Get(ctx, id)
	if err != nil {
		return EnqueueEpisodesResponse{}, err
	}

	jsText, err := FetchEpisodesJS(ctx, sub.BaseURL)
	if err != nil {
		return EnqueueEpisodesResponse{}, err
	}
	eps, err := ParseEpisodesJS(jsText)
	if err != nil {
		return EnqueueEpisodesResponse{}, err
	}

	selected, urls := selectPlayer(sub.Player, eps.Players)

	enqueuedEpisodes := []int{}
	enqueuedJobIDs := []string{}
	skipped := []EnqueueSkippedEpisode{}
	maxEp := MaxAvailableEpisode(urls)
	if maxEp < 0 {
		maxEp = 0
	}

	for _, ep := range norm {
		if ep > maxEp {
			skipped = append(skipped, EnqueueSkippedEpisode{Episode: ep, Reason: "not available"})
			continue
		}
		if ep-1 < 0 || ep-1 >= len(urls) {
			skipped = append(skipped, EnqueueSkippedEpisode{Episode: ep, Reason: "missing url"})
			continue
		}
		u := strings.TrimSpace(urls[ep-1])
		if u == "" {
			skipped = append(skipped, EnqueueSkippedEpisode{Episode: ep, Reason: "missing url"})
			continue
		}

		params := map[string]any{
			"url":            u,
			"path":           filepath.ToSlash(filepath.Join("subscriptions", sub.ID, fmt.Sprintf("%s-ep-%02d.mp4", safeLabel(sub.Label), ep))),
			"filename":       "",
			"subscriptionId": sub.ID,
			"episode":        ep,
			"source":         "anime-sama",
		}
		b, _ := json.Marshal(params)
		created, err := s.jobs.Create(ctx, CreateJobRequest{Type: "download", Params: b})
		if err != nil {
			skipped = append(skipped, EnqueueSkippedEpisode{Episode: ep, Reason: err.Error()})
			continue
		}
		enqueuedEpisodes = append(enqueuedEpisodes, ep)
		enqueuedJobIDs = append(enqueuedJobIDs, created.ID)
		if ep > sub.LastScheduledEpisode {
			sub.LastScheduledEpisode = ep
		}
	}

	sub.UpdatedAt = time.Now().UTC()
	updated, err := s.repo.Update(ctx, sub)
	if err != nil {
		return EnqueueEpisodesResponse{}, err
	}
	s.publish("subscription.updated", updated)

	return EnqueueEpisodesResponse{
		Subscription:     toSubscriptionDTO(updated),
		SelectedPlayer:   selected,
		EnqueuedEpisodes: enqueuedEpisodes,
		EnqueuedJobIDs:   enqueuedJobIDs,
		Skipped:          skipped,
	}, nil
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
				"url":            u,
				"path":           filepath.ToSlash(filepath.Join("subscriptions", sub.ID, fmt.Sprintf("%s-ep-%02d.mp4", safeLabel(sub.Label), ep))),
				"filename":       "",
				"subscriptionId": sub.ID,
				"episode":        ep,
				"source":         "anime-sama",
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
		Subscription:        toSubscriptionDTO(updated),
		SelectedPlayer:      selected,
		MaxAvailableEpisode: maxAvail,
		EnqueuedEpisodes:    enqueuedEpisodes,
		EnqueuedJobIDs:      enqueuedJobIDs,
		Message:             "note: episodes.js urls are host/embed urls; full video extraction pipeline is not implemented yet",
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

// SafeLabel expose le nettoyage des labels pour usage externe (ex: téléchargements ad-hoc).
func SafeLabel(s string) string {
	return safeLabel(s)
}

func selectPlayer(preferred string, players map[string][]string) (string, []string) {
	selected := preferred
	if selected == "" || strings.EqualFold(selected, "auto") {
		selected = BestPlayer(players)
		if selected == "auto" {
			selected = ""
		}
	}
	urls := players[selected]
	if len(urls) == 0 {
		selected = BestPlayer(players)
		urls = players[selected]
	}
	return selected, urls
}

func defaultLabelForBaseURL(canonBaseURL string) string {
	u, err := url.Parse(strings.TrimSpace(canonBaseURL))
	if err != nil {
		return ""
	}
	segs := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(segs) == 0 {
		return ""
	}

	// Typical pattern: /catalogue/<slug>/saisonN/<lang>/
	slug := ""
	for i := 0; i < len(segs)-1; i++ {
		if segs[i] == "catalogue" {
			slug = segs[i+1]
			break
		}
	}
	if slug == "" {
		// Fallback: best-effort guess (second segment is often a slug).
		if len(segs) >= 2 {
			slug = segs[1]
		} else {
			slug = segs[0]
		}
	}
	title := prettifySlug(slug)
	if title == "" {
		return ""
	}

	season := 0
	lang := ""
	for _, seg := range segs {
		if season == 0 {
			if n := parseSeasonSegment(seg); n > 0 {
				season = n
				continue
			}
		}
		if lang == "" {
			if l := normalizeLangSegment(seg); l != "" {
				lang = l
				continue
			}
		}
	}

	qual := []string{}
	if season > 0 {
		qual = append(qual, fmt.Sprintf("S%d", season))
	}
	if lang != "" {
		qual = append(qual, strings.ToUpper(lang))
	}
	if len(qual) == 0 {
		return title
	}
	return fmt.Sprintf("%s (%s)", title, strings.Join(qual, " "))
}

// DefaultLabelForBaseURL génère un label lisible à partir d'une baseUrl canonique.
func DefaultLabelForBaseURL(canonBaseURL string) string {
	return defaultLabelForBaseURL(canonBaseURL)
}

func parseSeasonSegment(seg string) int {
	seg = strings.ToLower(strings.TrimSpace(seg))
	if !strings.HasPrefix(seg, "saison") {
		return 0
	}
	n, err := strconv.Atoi(strings.TrimPrefix(seg, "saison"))
	if err != nil || n <= 0 {
		return 0
	}
	return n
}

func normalizeLangSegment(seg string) string {
	seg = strings.ToLower(strings.TrimSpace(seg))
	switch seg {
	case "vostfr", "vf", "vo", "vosten", "vost":
		return seg
	default:
		return ""
	}
}

func prettifySlug(slug string) string {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return ""
	}
	slug = strings.NewReplacer("-", " ", "_", " ").Replace(slug)
	fields := strings.Fields(slug)
	if len(fields) == 0 {
		return ""
	}
	for i, w := range fields {
		fields[i] = titleWord(w)
	}
	return strings.Join(fields, " ")
}

func titleWord(w string) string {
	if w == "" {
		return ""
	}
	r := []rune(w)
	r[0] = unicode.ToUpper(r[0])
	for i := 1; i < len(r); i++ {
		r[i] = unicode.ToLower(r[i])
	}
	return string(r)
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
