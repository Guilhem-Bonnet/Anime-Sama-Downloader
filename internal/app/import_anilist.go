package app

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

type AniListImportPreviewRequest struct {
	Statuses      []string `json:"statuses"`
	Season        int      `json:"season"`
	Lang          string   `json:"lang"`
	MaxCandidates int      `json:"maxCandidates"`
}

type AniListImportPreviewItem struct {
	AniListMediaID int                 `json:"anilistMediaId"`
	Title          string              `json:"title"`
	Titles         map[string]string   `json:"titles"`
	Synonyms       []string            `json:"synonyms"`
	Candidates     []AniListImportCand `json:"candidates"`
}

type AniListImportCand struct {
	CatalogueURL string  `json:"catalogueUrl"`
	BaseURL      string  `json:"baseUrl"`
	Slug         string  `json:"slug"`
	MatchedTitle string  `json:"matchedTitle"`
	Score        float64 `json:"score"`
}

type AniListImportPreviewResponse struct {
	Items []AniListImportPreviewItem `json:"items"`
}

type AniListImportConfirmRequest struct {
	Items []AniListImportConfirmItem `json:"items"`
}

type AniListImportConfirmItem struct {
	BaseURL string `json:"baseUrl"`
	Label   string `json:"label"`
	Player  string `json:"player,omitempty"`
}

type AniListImportConfirmResponse struct {
	Created []SubscriptionDTO        `json:"created"`
	Errors  []AniListImportItemError `json:"errors"`
}

type AniListImportAutoRequest struct {
	Statuses      []string `json:"statuses"`
	Season        int      `json:"season"`
	Lang          string   `json:"lang"`
	MaxCandidates int      `json:"maxCandidates"`
	MinScore      float64  `json:"minScore"`
}

type AniListImportAutoSkipped struct {
	AniListMediaID int     `json:"anilistMediaId"`
	Title          string  `json:"title"`
	Reason         string  `json:"reason"`
	BaseURL        string  `json:"baseUrl,omitempty"`
	TopScore       float64 `json:"topScore,omitempty"`
}

type AniListImportAutoResponse struct {
	Created []SubscriptionDTO          `json:"created"`
	Skipped []AniListImportAutoSkipped `json:"skipped"`
	Errors  []AniListImportItemError   `json:"errors"`
}

type AniListImportItemError struct {
	BaseURL string `json:"baseUrl"`
	Error   string `json:"error"`
}

type AniListImportService struct {
	anilist  interface {
		Watchlist(ctx context.Context, statuses []string) ([]AniListWatchlistEntry, error)
	}
	resolver interface {
		ResolveCandidates(ctx context.Context, titles []string, maxCandidates int) ([]AnimeSamaCandidate, error)
	}
	subs interface {
		Create(ctx context.Context, baseURL, label, player string) (SubscriptionDTO, error)
		List(ctx context.Context, limit int) ([]SubscriptionDTO, error)
	}
}

func NewAniListImportService(anilist *AniListService, resolver *AnimeSamaCatalogueResolver, subs *SubscriptionService) *AniListImportService {
	return &AniListImportService{anilist: anilist, resolver: resolver, subs: subs}
}

func (s *AniListImportService) Preview(ctx context.Context, req AniListImportPreviewRequest) (AniListImportPreviewResponse, error) {
	if s == nil || s.anilist == nil || s.resolver == nil {
		return AniListImportPreviewResponse{}, fmt.Errorf("import service not configured")
	}
	season := req.Season
	if season <= 0 {
		season = 1
	}
	lang := strings.TrimSpace(strings.ToLower(req.Lang))
	if lang == "" {
		lang = "vostfr"
	}

	entries, err := s.anilist.Watchlist(ctx, req.Statuses)
	if err != nil {
		return AniListImportPreviewResponse{}, err
	}

	items := make([]AniListImportPreviewItem, 0, len(entries))
	for _, e := range entries {
		titles := []string{}
		if e.Media.Title.Romaji != "" {
			titles = append(titles, e.Media.Title.Romaji)
		}
		if e.Media.Title.English != "" {
			titles = append(titles, e.Media.Title.English)
		}
		if e.Media.Title.Native != "" {
			titles = append(titles, e.Media.Title.Native)
		}
		for _, syn := range e.Media.Synonyms {
			if strings.TrimSpace(syn) != "" {
				titles = append(titles, syn)
			}
		}

		cands, _ := s.resolver.ResolveCandidates(ctx, titles, req.MaxCandidates)
		mapped := make([]AniListImportCand, 0, len(cands))
		for _, c := range cands {
			base := strings.TrimRight(strings.TrimSpace(c.CatalogueURL), "/") + fmt.Sprintf("/saison%d/%s/", season, lang)
			mapped = append(mapped, AniListImportCand{
				CatalogueURL: c.CatalogueURL,
				BaseURL:      base,
				Slug:         c.Slug,
				MatchedTitle: c.MatchedTitle,
				Score:        c.Score,
			})
		}

		display := firstNonEmpty(e.Media.Title.English, e.Media.Title.Romaji, e.Media.Title.Native)
		items = append(items, AniListImportPreviewItem{
			AniListMediaID: e.Media.ID,
			Title:          display,
			Titles: map[string]string{
				"romaji":  e.Media.Title.Romaji,
				"english": e.Media.Title.English,
				"native":  e.Media.Title.Native,
			},
			Synonyms:   e.Media.Synonyms,
			Candidates: mapped,
		})
	}

	return AniListImportPreviewResponse{Items: items}, nil
}

func (s *AniListImportService) Confirm(ctx context.Context, req AniListImportConfirmRequest) (AniListImportConfirmResponse, error) {
	if s == nil || s.subs == nil {
		return AniListImportConfirmResponse{}, fmt.Errorf("subscription service not configured")
	}

	created := []SubscriptionDTO{}
	errorsOut := []AniListImportItemError{}
	for _, it := range req.Items {
		label := strings.TrimSpace(it.Label)
		if label == "" {
			label = "Anime"
		}
		sub, err := s.subs.Create(ctx, it.BaseURL, label, it.Player)
		if err != nil {
			if errors.Is(err, ErrConflict) {
				errorsOut = append(errorsOut, AniListImportItemError{BaseURL: it.BaseURL, Error: "already subscribed"})
				continue
			}
			errorsOut = append(errorsOut, AniListImportItemError{BaseURL: it.BaseURL, Error: err.Error()})
			continue
		}
		created = append(created, sub)
	}
	return AniListImportConfirmResponse{Created: created, Errors: errorsOut}, nil
}

func (s *AniListImportService) AutoConfirm(ctx context.Context, req AniListImportAutoRequest) (AniListImportAutoResponse, error) {
	if s == nil || s.subs == nil {
		return AniListImportAutoResponse{}, fmt.Errorf("subscription service not configured")
	}

	minScore := req.MinScore
	if minScore <= 0 {
		minScore = 0.95
	}

	prev, err := s.Preview(ctx, AniListImportPreviewRequest{
		Statuses:      req.Statuses,
		Season:        req.Season,
		Lang:          req.Lang,
		MaxCandidates: req.MaxCandidates,
	})
	if err != nil {
		return AniListImportAutoResponse{}, err
	}

	existing, _ := s.subs.List(ctx, 10000)
	seenBase := map[string]struct{}{}
	for _, sub := range existing {
		canon, err := CanonicalizeAnimeSamaBaseURL(sub.BaseURL)
		if err != nil {
			canon = strings.TrimSpace(sub.BaseURL)
		}
		if canon != "" {
			seenBase[canon] = struct{}{}
		}
	}

	created := []SubscriptionDTO{}
	skipped := []AniListImportAutoSkipped{}
	errorsOut := []AniListImportItemError{}

	for _, it := range prev.Items {
		if len(it.Candidates) == 0 {
			skipped = append(skipped, AniListImportAutoSkipped{AniListMediaID: it.AniListMediaID, Title: it.Title, Reason: "no candidates"})
			continue
		}

		best := it.Candidates[0]
		secondScore := 0.0
		if len(it.Candidates) > 1 {
			secondScore = it.Candidates[1].Score
		}
		if best.Score < minScore {
			skipped = append(skipped, AniListImportAutoSkipped{AniListMediaID: it.AniListMediaID, Title: it.Title, Reason: "low confidence", BaseURL: best.BaseURL, TopScore: best.Score})
			continue
		}
		if len(it.Candidates) > 1 && secondScore >= minScore {
			skipped = append(skipped, AniListImportAutoSkipped{AniListMediaID: it.AniListMediaID, Title: it.Title, Reason: "ambiguous (multiple high-score candidates)", BaseURL: best.BaseURL, TopScore: best.Score})
			continue
		}

		canonBase, err := CanonicalizeAnimeSamaBaseURL(best.BaseURL)
		if err != nil {
			canonBase = strings.TrimSpace(best.BaseURL)
		}
		if canonBase != "" {
			if _, ok := seenBase[canonBase]; ok {
				skipped = append(skipped, AniListImportAutoSkipped{AniListMediaID: it.AniListMediaID, Title: it.Title, Reason: "already subscribed", BaseURL: canonBase, TopScore: best.Score})
				continue
			}
		}

		label := strings.TrimSpace(it.Title)
		if label == "" {
			label = "Anime"
		}
		sub, err := s.subs.Create(ctx, best.BaseURL, label, "auto")
		if err != nil {
			errorsOut = append(errorsOut, AniListImportItemError{BaseURL: best.BaseURL, Error: err.Error()})
			continue
		}
		created = append(created, sub)
		if canon, err := CanonicalizeAnimeSamaBaseURL(sub.BaseURL); err == nil && canon != "" {
			seenBase[canon] = struct{}{}
		} else if strings.TrimSpace(sub.BaseURL) != "" {
			seenBase[strings.TrimSpace(sub.BaseURL)] = struct{}{}
		}
	}

	return AniListImportAutoResponse{Created: created, Skipped: skipped, Errors: errorsOut}, nil
}

func firstNonEmpty(vs ...string) string {
	for _, v := range vs {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}
