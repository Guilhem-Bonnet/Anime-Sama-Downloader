package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/ports"
)

// AniListSearchService implements ports.AnimeSearch using the AniList GraphQL API.
// No authentication token is required for public search queries.
type AniListSearchService struct {
	endpoint string
	client   *http.Client
}

// NewAniListSearchService creates a new search service backed by AniList.
func NewAniListSearchService() *AniListSearchService {
	return &AniListSearchService{
		endpoint: "https://graphql.anilist.co",
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// aniListSearchMediaEntry represents a single media entry from AniList search.
type aniListSearchMediaEntry struct {
	ID         int `json:"id"`
	Title      struct {
		Romaji  string `json:"romaji"`
		English string `json:"english"`
	} `json:"title"`
	CoverImage struct {
		Large  string `json:"large"`
		Medium string `json:"medium"`
	} `json:"coverImage"`
	SeasonYear   *int     `json:"seasonYear"`
	StartDate    struct {
		Year *int `json:"year"`
	} `json:"startDate"`
	Status       string   `json:"status"`
	Episodes     *int     `json:"episodes"`
	Genres       []string `json:"genres"`
	AverageScore *int     `json:"averageScore"`
}

type aniListSearchPageData struct {
	Page struct {
		Media []aniListSearchMediaEntry `json:"media"`
	} `json:"Page"`
}

// Search performs an anime search via AniList GraphQL.
func (s *AniListSearchService) Search(ctx context.Context, query string) ([]domain.AnimeSearchResult, error) {
	return s.SearchWithFilters(ctx, ports.SearchFilters{Query: query})
}

// SearchWithFilters performs a filtered anime search via AniList GraphQL.
func (s *AniListSearchService) SearchWithFilters(ctx context.Context, filters ports.SearchFilters) ([]domain.AnimeSearchResult, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	query := strings.TrimSpace(filters.Query)

	// Build dynamic GraphQL query based on filters
	gqlQuery, variables := s.buildQuery(query, filters)

	reqBody := aniListGraphQLRequest{
		Query:     gqlQuery,
		Variables: variables,
	}

	var out aniListGraphQLResponse[aniListSearchPageData]
	if err := s.doRequest(ctx, reqBody, &out); err != nil {
		return nil, fmt.Errorf("anilist search: %w", err)
	}
	if len(out.Errors) > 0 {
		return nil, fmt.Errorf("anilist search: %s", out.Errors[0].Message)
	}

	results := make([]domain.AnimeSearchResult, 0, len(out.Data.Page.Media))
	for _, m := range out.Data.Page.Media {
		results = append(results, s.toSearchResult(m))
	}

	return results, nil
}

// buildQuery constructs the GraphQL query and variables depending on filters.
func (s *AniListSearchService) buildQuery(query string, filters ports.SearchFilters) (string, map[string]any) {
	vars := map[string]any{
		"page":    1,
		"perPage": 25,
		"type":    "ANIME",
	}

	// Build argument parts
	argDefs := []string{"$page:Int", "$perPage:Int", "$type:MediaType"}
	mediaArgs := []string{"type:$type", "sort:SEARCH_MATCH"}

	if query != "" {
		argDefs = append(argDefs, "$search:String")
		mediaArgs = append(mediaArgs, "search:$search")
		vars["search"] = query
	} else {
		// When no query, sort by popularity
		mediaArgs = []string{"type:$type", "sort:POPULARITY_DESC"}
	}

	if len(filters.Genres) > 0 {
		argDefs = append(argDefs, "$genres:[String]")
		mediaArgs = append(mediaArgs, "genre_in:$genres")
		vars["genres"] = filters.Genres
	}

	if filters.Status != "" {
		aniListStatus := mapStatusToAniList(filters.Status)
		if aniListStatus != "" {
			argDefs = append(argDefs, "$status:MediaStatus")
			mediaArgs = append(mediaArgs, "status:$status")
			vars["status"] = aniListStatus
		}
	}

	if filters.YearMin > 0 {
		argDefs = append(argDefs, "$yearGreater:FuzzyDateInt")
		mediaArgs = append(mediaArgs, "startDate_greater:$yearGreater")
		vars["yearGreater"] = filters.YearMin*10000 // AniList uses YYYYMMDD format
	}

	if filters.YearMax > 0 {
		argDefs = append(argDefs, "$yearLesser:FuzzyDateInt")
		mediaArgs = append(mediaArgs, "startDate_lesser:$yearLesser")
		vars["yearLesser"] = (filters.YearMax+1)*10000 - 1 // end of year
	}

	gql := fmt.Sprintf(`query(%s){
		Page(page:$page, perPage:$perPage){
			media(%s){
				id
				title { romaji english }
				coverImage { large medium }
				seasonYear
				startDate { year }
				status
				episodes
				genres
				averageScore
			}
		}
	}`, strings.Join(argDefs, ","), strings.Join(mediaArgs, ","))

	return gql, vars
}

// mapStatusToAniList converts our status string to AniList MediaStatus enum.
func mapStatusToAniList(status string) string {
	switch strings.ToLower(status) {
	case "ongoing", "releasing":
		return "RELEASING"
	case "completed", "finished":
		return "FINISHED"
	case "planning", "not_yet_released":
		return "NOT_YET_RELEASED"
	case "cancelled":
		return "CANCELLED"
	case "hiatus":
		return "HIATUS"
	default:
		return ""
	}
}

// mapAniListStatus converts AniList MediaStatus to our status string.
func mapAniListStatus(status string) string {
	switch status {
	case "RELEASING":
		return "ongoing"
	case "FINISHED":
		return "completed"
	case "NOT_YET_RELEASED":
		return "planning"
	case "CANCELLED":
		return "cancelled"
	case "HIATUS":
		return "hiatus"
	default:
		return strings.ToLower(status)
	}
}

// toSearchResult converts an AniList media entry to our domain model.
func (s *AniListSearchService) toSearchResult(m aniListSearchMediaEntry) domain.AnimeSearchResult {
	// Prefer English title, fallback to Romaji
	title := m.Title.English
	if title == "" {
		title = m.Title.Romaji
	}

	// Thumbnail: prefer large, fallback to medium
	thumbnail := m.CoverImage.Large
	if thumbnail == "" {
		thumbnail = m.CoverImage.Medium
	}

	// Year: prefer seasonYear, fallback to startDate.year
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

	// Use AniList ID as string identifier (prefixed to avoid collision)
	id := fmt.Sprintf("al-%d", m.ID)

	return domain.AnimeSearchResult{
		ID:           id,
		Title:        title,
		ThumbnailURL: thumbnail,
		Year:         year,
		Status:       mapAniListStatus(m.Status),
		EpisodeCount: episodes,
		Genres:       m.Genres,
	}
}

// doRequest executes a GraphQL request against AniList (no auth required for search).
func (s *AniListSearchService) doRequest(ctx context.Context, req aniListGraphQLRequest, out any) error {
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

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("anilist http error: %s", resp.Status)
	}

	return json.NewDecoder(resp.Body).Decode(out)
}
