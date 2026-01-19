package app

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
)

var ErrAniListNotConfigured = errors.New("anilist not configured")

type AniListService struct {
	settings func(ctx context.Context) (domain.Settings, error)
	endpoint string
	client   *http.Client
}

func NewAniListService(settingsGetter func(ctx context.Context) (domain.Settings, error)) *AniListService {
	return &AniListService{
		settings: settingsGetter,
		endpoint: "https://graphql.anilist.co",
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (s *AniListService) WithEndpoint(endpoint string) *AniListService {
	if strings.TrimSpace(endpoint) != "" {
		s.endpoint = strings.TrimSpace(endpoint)
	}
	return s
}

type aniListGraphQLRequest struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables,omitempty"`
}

type aniListGraphQLError struct {
	Message string `json:"message"`
}

type aniListGraphQLResponse[T any] struct {
	Data   T                   `json:"data"`
	Errors []aniListGraphQLError `json:"errors,omitempty"`
}

type AniListViewer struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type viewerData struct {
	Viewer AniListViewer `json:"Viewer"`
}

func (s *AniListService) Viewer(ctx context.Context) (AniListViewer, error) {
	if s == nil || s.settings == nil {
		return AniListViewer{}, ErrAniListNotConfigured
	}
	st, err := s.settings(ctx)
	if err != nil {
		return AniListViewer{}, err
	}
	token := strings.TrimSpace(st.AniListToken)
	if token == "" {
		return AniListViewer{}, ErrAniListNotConfigured
	}

	req := aniListGraphQLRequest{Query: `query { Viewer { id name } }`}
	var out aniListGraphQLResponse[viewerData]
	if err := s.do(ctx, token, req, &out); err != nil {
		return AniListViewer{}, err
	}
	if len(out.Errors) > 0 {
		return AniListViewer{}, errors.New(out.Errors[0].Message)
	}
	return out.Data.Viewer, nil
}

type AniListAiringScheduleEntry struct {
	ID       int `json:"id"`
	AiringAt int `json:"airingAt"`
	Episode  int `json:"episode"`
	Media    struct {
		ID    int `json:"id"`
		Title struct {
			Romaji  string `json:"romaji"`
			English string `json:"english"`
			Native  string `json:"native"`
		} `json:"title"`
	} `json:"media"`
}

type airingPageData struct {
	Page struct {
		AiringSchedules []AniListAiringScheduleEntry `json:"airingSchedules"`
	} `json:"Page"`
}

func (s *AniListService) AiringSchedule(ctx context.Context, from, to time.Time, limit int) ([]AniListAiringScheduleEntry, error) {
	if limit <= 0 {
		limit = 50
	}
	start := int(from.UTC().Unix())
	end := int(to.UTC().Unix())

	req := aniListGraphQLRequest{
		Query: `query($start:Int,$end:Int,$page:Int,$perPage:Int){
			Page(page:$page, perPage:$perPage){
				airingSchedules(airingAt_greater:$start, airingAt_lesser:$end, sort: AIRING_AT){
					id airingAt episode
					media{ id title{ romaji english native } }
				}
			}
		}`,
		Variables: map[string]any{"start": start, "end": end, "page": 1, "perPage": limit},
	}

	var out aniListGraphQLResponse[airingPageData]
	if err := s.do(ctx, "", req, &out); err != nil {
		return nil, err
	}
	if len(out.Errors) > 0 {
		return nil, errors.New(out.Errors[0].Message)
	}
	return out.Data.Page.AiringSchedules, nil
}

type AniListWatchlistEntry struct {
	Status   string `json:"status"`
	Progress int    `json:"progress"`
	Media    struct {
		ID       int      `json:"id"`
		Synonyms []string `json:"synonyms"`
		Title    struct {
			Romaji  string `json:"romaji"`
			English string `json:"english"`
			Native  string `json:"native"`
		} `json:"title"`
	} `json:"media"`
}

type watchlistData struct {
	MediaListCollection struct {
		Lists []struct {
			Entries []AniListWatchlistEntry `json:"entries"`
		} `json:"lists"`
	} `json:"MediaListCollection"`
}

func (s *AniListService) Watchlist(ctx context.Context, statuses []string) ([]AniListWatchlistEntry, error) {
	if s == nil || s.settings == nil {
		return nil, ErrAniListNotConfigured
	}
	st, err := s.settings(ctx)
	if err != nil {
		return nil, err
	}
	token := strings.TrimSpace(st.AniListToken)
	if token == "" {
		return nil, ErrAniListNotConfigured
	}

	viewer, err := s.Viewer(ctx)
	if err != nil {
		return nil, err
	}
	if len(statuses) == 0 {
		statuses = []string{"CURRENT", "PLANNING"}
	}

	req := aniListGraphQLRequest{
		Query: `query($userId:Int,$statusIn:[MediaListStatus]){
			MediaListCollection(userId:$userId, type: ANIME, status_in:$statusIn){
				lists{
					entries{ status progress media{ id synonyms title{ romaji english native } } }
				}
			}
		}`,
		Variables: map[string]any{"userId": viewer.ID, "statusIn": statuses},
	}

	var out aniListGraphQLResponse[watchlistData]
	if err := s.do(ctx, token, req, &out); err != nil {
		return nil, err
	}
	if len(out.Errors) > 0 {
		return nil, errors.New(out.Errors[0].Message)
	}

	flat := make([]AniListWatchlistEntry, 0)
	for _, l := range out.Data.MediaListCollection.Lists {
		for _, e := range l.Entries {
			flat = append(flat, e)
		}
	}
	return flat, nil
}

func (s *AniListService) do(ctx context.Context, token string, req aniListGraphQLRequest, out any) error {
	b, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, s.endpoint, bytes.NewReader(b))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("User-Agent", "asd-server")
	if strings.TrimSpace(token) != "" {
		httpReq.Header.Set("Authorization", "Bearer "+strings.TrimSpace(token))
	}

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		// AniList tends to return JSON, but we keep it simple.
		return errors.New("anilist http error: " + resp.Status)
	}

	return json.NewDecoder(resp.Body).Decode(out)
}
