package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
)

// AniListDetailService implements ports.AnimeDetailService using the AniList GraphQL API.
// IDs with the "al-" prefix (e.g. "al-187901") are resolved via AniList.
// Other IDs are delegated to a fallback service (if provided).
type AniListDetailService struct {
	endpoint string
	client   *http.Client
	fallback interface {
		GetDetail(ctx context.Context, id string) (domain.AnimeDetail, error)
	}
}

// NewAniListDetailService creates a new detail service backed by AniList.
// fallback may be nil; if so, non-AniList IDs will return "not found".
func NewAniListDetailService(fallback interface {
	GetDetail(ctx context.Context, id string) (domain.AnimeDetail, error)
}) *AniListDetailService {
	return &AniListDetailService{
		endpoint: "https://graphql.anilist.co",
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
		fallback: fallback,
	}
}

// aniListDetailMedia represents a single media entry from an AniList detail query.
type aniListDetailMedia struct {
	ID    int `json:"id"`
	Title struct {
		Romaji  string `json:"romaji"`
		English string `json:"english"`
	} `json:"title"`
	CoverImage struct {
		Large  string `json:"large"`
		Medium string `json:"medium"`
	} `json:"coverImage"`
	Description  *string  `json:"description"`
	SeasonYear   *int     `json:"seasonYear"`
	StartDate    struct {
		Year *int `json:"year"`
	} `json:"startDate"`
	Status       string   `json:"status"`
	Episodes     *int     `json:"episodes"`
	Genres       []string `json:"genres"`
	AverageScore *int     `json:"averageScore"`
}

type aniListDetailData struct {
	Media aniListDetailMedia `json:"Media"`
}

const aniListDetailQuery = `query($id:Int){
	Media(id:$id, type:ANIME){
		id
		title { romaji english }
		coverImage { large medium }
		description(asHtml:false)
		seasonYear
		startDate { year }
		status
		episodes
		genres
		averageScore
	}
}`

// GetDetail returns anime details by ID.
func (s *AniListDetailService) GetDetail(ctx context.Context, id string) (domain.AnimeDetail, error) {
	select {
	case <-ctx.Done():
		return domain.AnimeDetail{}, ctx.Err()
	default:
	}

	// Handle AniList IDs (al-{numericID})
	if strings.HasPrefix(id, "al-") {
		numStr := strings.TrimPrefix(id, "al-")
		aniListID, err := strconv.Atoi(numStr)
		if err != nil {
			return domain.AnimeDetail{}, fmt.Errorf("invalid anilist id: %s", id)
		}
		return s.fetchFromAniList(ctx, id, aniListID)
	}

	// Delegate non-AniList IDs to fallback
	if s.fallback != nil {
		return s.fallback.GetDetail(ctx, id)
	}

	return domain.AnimeDetail{}, fmt.Errorf("anime not found: %s", id)
}

// fetchFromAniList queries AniList for anime details by numeric ID.
func (s *AniListDetailService) fetchFromAniList(ctx context.Context, originalID string, aniListID int) (domain.AnimeDetail, error) {
	reqBody := aniListGraphQLRequest{
		Query: aniListDetailQuery,
		Variables: map[string]any{
			"id": aniListID,
		},
	}

	b, err := json.Marshal(reqBody)
	if err != nil {
		return domain.AnimeDetail{}, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, s.endpoint, bytes.NewReader(b))
	if err != nil {
		return domain.AnimeDetail{}, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("User-Agent", "asd-server")

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return domain.AnimeDetail{}, fmt.Errorf("anilist detail: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return domain.AnimeDetail{}, fmt.Errorf("anilist http error: %s", resp.Status)
	}

	var out aniListGraphQLResponse[aniListDetailData]
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return domain.AnimeDetail{}, err
	}
	if len(out.Errors) > 0 {
		return domain.AnimeDetail{}, fmt.Errorf("anilist: %s", out.Errors[0].Message)
	}

	return s.toAnimeDetail(originalID, out.Data.Media), nil
}

// toAnimeDetail converts an AniList media response to our domain model.
func (s *AniListDetailService) toAnimeDetail(id string, m aniListDetailMedia) domain.AnimeDetail {
	title := m.Title.English
	if title == "" {
		title = m.Title.Romaji
	}

	thumbnail := m.CoverImage.Large
	if thumbnail == "" {
		thumbnail = m.CoverImage.Medium
	}

	synopsis := ""
	if m.Description != nil {
		synopsis = stripHTMLTags(*m.Description)
	}

	year := 0
	if m.SeasonYear != nil {
		year = *m.SeasonYear
	} else if m.StartDate.Year != nil {
		year = *m.StartDate.Year
	}

	episodes := 0
	if m.Episodes != nil {
		episodes = *m.Episodes
	}

	// Build placeholder seasons/episodes (AniList doesn't provide episode-level detail for free)
	seasons := []domain.Season{}
	if episodes > 0 {
		eps := make([]domain.Episode, episodes)
		for i := 0; i < episodes; i++ {
			eps[i] = domain.Episode{
				Number:       i + 1,
				Title:        fmt.Sprintf("Épisode %d", i+1),
				SeasonNumber: 1,
				URL:          "", // download URL resolved later via anime-sama
			}
		}
		seasons = append(seasons, domain.Season{
			Number:   1,
			Name:     "Saison 1",
			Episodes: eps,
		})
	}

	return domain.AnimeDetail{
		ID:           id,
		Title:        title,
		ThumbnailURL: thumbnail,
		Synopsis:     synopsis,
		Year:         year,
		Status:       mapAniListStatus(m.Status),
		Genres:       m.Genres,
		EpisodeCount: episodes,
		Seasons:      seasons,
	}
}

// stripHTMLTags removes basic HTML tags from AniList description text.
func stripHTMLTags(s string) string {
	var result strings.Builder
	inTag := false
	for _, r := range s {
		switch {
		case r == '<':
			inTag = true
		case r == '>':
			inTag = false
		case !inTag:
			result.WriteRune(r)
		}
	}
	// Normalize <br> to line breaks (AniList uses <br>)
	out := result.String()
	out = strings.ReplaceAll(out, "\n\n\n", "\n\n")
	return strings.TrimSpace(out)
}
