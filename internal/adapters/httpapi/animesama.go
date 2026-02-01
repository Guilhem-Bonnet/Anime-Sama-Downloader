package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"fmt"
	"errors"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/app"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/httpjson"
	"github.com/go-chi/chi/v5"
)

type AnimeSamaResolver interface {
	ResolveCandidates(ctx context.Context, titles []string, maxCandidates int) ([]app.AnimeSamaCandidate, error)
}

type AnimeSamaHandler struct {
	resolver AnimeSamaResolver
	jobs     *app.JobService
}

func NewAnimeSamaHandler(resolver AnimeSamaResolver, jobs *app.JobService) *AnimeSamaHandler {
	return &AnimeSamaHandler{resolver: resolver, jobs: jobs}
}

type animeSamaResolveRequest struct {
	Titles        []string `json:"titles"`
	Season        int      `json:"season"`
	Lang          string   `json:"lang"`
	MaxCandidates int      `json:"maxCandidates"`
}

type animeSamaResolveCandidate struct {
	CatalogueURL string  `json:"catalogueUrl"`
	BaseURL      string  `json:"baseUrl"`
	Slug         string  `json:"slug"`
	MatchedTitle string  `json:"matchedTitle"`
	Score        float64 `json:"score"`
}

type animeSamaResolveResponse struct {
	Candidates []animeSamaResolveCandidate `json:"candidates"`
}

func (h *AnimeSamaHandler) Routes(r chi.Router) {
	r.Route("/animesama", func(r chi.Router) {
		r.Post("/resolve", h.resolve)
		r.Post("/scan", h.scan)
		r.Post("/episodes", h.episodes)
		r.Post("/enqueue", h.enqueue)
	})
}

type animeSamaEpisodesRequest struct {
	BaseURL string `json:"baseUrl"`
}

type animeSamaEpisodeStatus struct {
	Episode   int  `json:"episode"`
	Available bool `json:"available"`
}

type animeSamaEpisodesResponse struct {
	BaseURL             string                 `json:"baseUrl"`
	SelectedPlayer      string                 `json:"selectedPlayer"`
	MaxAvailableEpisode int                    `json:"maxAvailableEpisode"`
	Episodes            []animeSamaEpisodeStatus `json:"episodes"`
}

func selectPlayerAdHoc(players map[string][]string) (string, []string) {
	selected := app.BestPlayer(players)
	urls := players[selected]
	if len(urls) == 0 {
		// fallback best-effort: pick any non-empty key
		for k, v := range players {
			selected = k
			urls = v
			break
		}
	}
	return selected, urls
}

func (h *AnimeSamaHandler) episodes(w http.ResponseWriter, r *http.Request) {
	var req animeSamaEpisodesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	canon, err := app.CanonicalizeAnimeSamaBaseURL(req.BaseURL)
	if err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	jsText, err := app.FetchEpisodesJS(r.Context(), canon)
	if err != nil {
		httpjson.WriteError(w, http.StatusBadGateway, err.Error())
		return
	}
	eps, err := app.ParseEpisodesJS(jsText)
	if err != nil {
		httpjson.WriteError(w, http.StatusBadGateway, err.Error())
		return
	}

	selected, urls := selectPlayerAdHoc(eps.Players)
	maxEp := app.MaxAvailableEpisode(urls)
	if maxEp < 0 {
		maxEp = 0
	}

	out := make([]animeSamaEpisodeStatus, 0, maxEp)
	for ep := 1; ep <= maxEp; ep++ {
		available := false
		if ep-1 >= 0 && ep-1 < len(urls) {
			available = strings.TrimSpace(urls[ep-1]) != ""
		}
		out = append(out, animeSamaEpisodeStatus{Episode: ep, Available: available})
	}

	httpjson.Write(w, http.StatusOK, animeSamaEpisodesResponse{
		BaseURL:             canon,
		SelectedPlayer:      selected,
		MaxAvailableEpisode: maxEp,
		Episodes:            out,
	})
}

type animeSamaEnqueueRequest struct {
	BaseURL   string `json:"baseUrl"`
	Label     string `json:"label"`
	Episodes  []int  `json:"episodes"`
}

type animeSamaEnqueueSkippedEpisode struct {
	Episode int    `json:"episode"`
	Reason  string `json:"reason"`
}

type animeSamaEnqueueResponse struct {
	BaseURL          string                         `json:"baseUrl"`
	Label            string                         `json:"label"`
	SelectedPlayer   string                         `json:"selectedPlayer"`
	EnqueuedEpisodes []int                          `json:"enqueuedEpisodes"`
	EnqueuedJobIDs   []string                       `json:"enqueuedJobIds"`
	Skipped          []animeSamaEnqueueSkippedEpisode `json:"skipped"`
}

func (h *AnimeSamaHandler) enqueue(w http.ResponseWriter, r *http.Request) {
	if h == nil || h.jobs == nil {
		httpjson.WriteError(w, http.StatusNotImplemented, "job service not configured")
		return
	}

	var req animeSamaEnqueueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if len(req.Episodes) == 0 {
		httpjson.WriteError(w, http.StatusBadRequest, "missing episodes")
		return
	}
	canon, err := app.CanonicalizeAnimeSamaBaseURL(req.BaseURL)
	if err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	label := strings.TrimSpace(req.Label)
	if label == "" {
		label = app.DefaultLabelForBaseURL(canon)
		if strings.TrimSpace(label) == "" {
			label = "Anime"
		}
	}

	// Normalize input.
	seen := map[int]struct{}{}
	norm := make([]int, 0, len(req.Episodes))
	for _, ep := range req.Episodes {
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
		httpjson.WriteError(w, http.StatusBadRequest, "no valid episodes")
		return
	}
	sort.Ints(norm)

	jsText, err := app.FetchEpisodesJS(r.Context(), canon)
	if err != nil {
		httpjson.WriteError(w, http.StatusBadGateway, err.Error())
		return
	}
	eps, err := app.ParseEpisodesJS(jsText)
	if err != nil {
		httpjson.WriteError(w, http.StatusBadGateway, err.Error())
		return
	}

	selected, urls := selectPlayerAdHoc(eps.Players)
	maxEp := app.MaxAvailableEpisode(urls)
	if maxEp < 0 {
		maxEp = 0
	}

	enqueuedEpisodes := []int{}
	enqueuedJobIDs := []string{}
	skipped := []animeSamaEnqueueSkippedEpisode{}

	root := filepath.ToSlash(filepath.Join("adhoc", app.SafeLabel(label)))
	for _, ep := range norm {
		if ep > maxEp {
			skipped = append(skipped, animeSamaEnqueueSkippedEpisode{Episode: ep, Reason: "not available"})
			continue
		}
		if ep-1 < 0 || ep-1 >= len(urls) {
			skipped = append(skipped, animeSamaEnqueueSkippedEpisode{Episode: ep, Reason: "missing url"})
			continue
		}
		u := strings.TrimSpace(urls[ep-1])
		if u == "" {
			skipped = append(skipped, animeSamaEnqueueSkippedEpisode{Episode: ep, Reason: "missing url"})
			continue
		}

		params := map[string]any{
			"url":      u,
			"path":     filepath.ToSlash(filepath.Join(root, fmt.Sprintf("%s-ep-%02d.mp4", app.SafeLabel(label), ep))),
			"filename": "",
			"baseUrl":  canon,
			"label":    label,
			"episode":  ep,
			"source":   "anime-sama",
			"mode":     "adhoc",
		}
		b, _ := json.Marshal(params)
		created, err := h.jobs.Create(r.Context(), app.CreateJobRequest{Type: "download", Params: b})
		if err != nil {
			skipped = append(skipped, animeSamaEnqueueSkippedEpisode{Episode: ep, Reason: err.Error()})
			continue
		}
		enqueuedEpisodes = append(enqueuedEpisodes, ep)
		enqueuedJobIDs = append(enqueuedJobIDs, created.ID)
	}

	if len(enqueuedEpisodes) == 0 && len(skipped) > 0 {
		// Rien de planifié: on renvoie quand même 200 avec détail, mais si tout est vide => erreur claire.
		// Cas typique: saison vide.
	}

	if len(enqueuedEpisodes) == 0 && len(skipped) == 0 {
		httpjson.WriteError(w, http.StatusBadRequest, errors.New("nothing to enqueue").Error())
		return
	}

	httpjson.Write(w, http.StatusOK, animeSamaEnqueueResponse{
		BaseURL:          canon,
		Label:            label,
		SelectedPlayer:   selected,
		EnqueuedEpisodes: enqueuedEpisodes,
		EnqueuedJobIDs:   enqueuedJobIDs,
		Skipped:          skipped,
	})
}

type animeSamaScanRequest struct {
	CatalogueURL string   `json:"catalogueUrl"`
	MaxSeason    int      `json:"maxSeason"`
	Langs        []string `json:"langs"`
}

type animeSamaScanOption struct {
	BaseURL             string `json:"baseUrl"`
	Season              int    `json:"season"`
	Lang                string `json:"lang"`
	SelectedPlayer      string `json:"selectedPlayer"`
	MaxAvailableEpisode int    `json:"maxAvailableEpisode"`
}

type animeSamaScanResponse struct {
	Options []animeSamaScanOption `json:"options"`
}

func (h *AnimeSamaHandler) scan(w http.ResponseWriter, r *http.Request) {
	var req animeSamaScanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	catalogueURL := strings.TrimSpace(req.CatalogueURL)
	if catalogueURL == "" {
		httpjson.WriteError(w, http.StatusBadRequest, "missing catalogueUrl")
		return
	}
	if !strings.HasSuffix(catalogueURL, "/") {
		catalogueURL += "/"
	}

	maxSeason := req.MaxSeason
	if maxSeason <= 0 {
		maxSeason = 5
	}
	if maxSeason > 20 {
		maxSeason = 20
	}

	langs := make([]string, 0, 4)
	if len(req.Langs) == 0 {
		langs = []string{"vostfr", "vf"}
	} else {
		seen := map[string]struct{}{}
		for _, l := range req.Langs {
			l = strings.ToLower(strings.TrimSpace(l))
			if l == "" {
				continue
			}
			if _, ok := seen[l]; ok {
				continue
			}
			seen[l] = struct{}{}
			langs = append(langs, l)
			if len(langs) >= 6 {
				break
			}
		}
		if len(langs) == 0 {
			langs = []string{"vostfr", "vf"}
		}
	}

	baseRoot := strings.TrimRight(catalogueURL, "/")
	out := make([]animeSamaScanOption, 0, maxSeason*len(langs))

	for season := 1; season <= maxSeason; season++ {
		for _, lang := range langs {
			baseURL := fmt.Sprintf("%s/saison%d/%s/", baseRoot, season, lang)
			jsText, err := app.FetchEpisodesJS(r.Context(), baseURL)
			if err != nil {
				continue
			}
			eps, err := app.ParseEpisodesJS(jsText)
			if err != nil {
				continue
			}
			selected := app.BestPlayer(eps.Players)
			urls := eps.Players[selected]
			maxEp := app.MaxAvailableEpisode(urls)
			if maxEp <= 0 {
				continue
			}
			out = append(out, animeSamaScanOption{BaseURL: baseURL, Season: season, Lang: lang, SelectedPlayer: selected, MaxAvailableEpisode: maxEp})
			if len(out) >= 40 {
				break
			}
		}
		if len(out) >= 40 {
			break
		}
	}

	httpjson.Write(w, http.StatusOK, animeSamaScanResponse{Options: out})
}

func (h *AnimeSamaHandler) resolve(w http.ResponseWriter, r *http.Request) {
	if h == nil || h.resolver == nil {
		httpjson.WriteError(w, http.StatusNotImplemented, "resolver disabled")
		return
	}

	var req animeSamaResolveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if len(req.Titles) == 0 {
		httpjson.WriteError(w, http.StatusBadRequest, "missing titles")
		return
	}
	season := req.Season
	if season <= 0 {
		season = 1
	}
	lang := strings.TrimSpace(strings.ToLower(req.Lang))
	if lang == "" {
		lang = "vostfr"
	}
	maxC := req.MaxCandidates
	if maxC <= 0 {
		maxC = 5
	}
	if maxC > 10 {
		maxC = 10
	}

	cands, err := h.resolver.ResolveCandidates(r.Context(), req.Titles, maxC)
	if err != nil {
		httpjson.WriteError(w, http.StatusBadGateway, err.Error())
		return
	}

	out := make([]animeSamaResolveCandidate, 0, len(cands))
	for _, c := range cands {
		base := strings.TrimRight(strings.TrimSpace(c.CatalogueURL), "/") + "/" + "saison" + strconv.Itoa(season) + "/" + lang + "/"
		out = append(out, animeSamaResolveCandidate{
			CatalogueURL: c.CatalogueURL,
			BaseURL:      base,
			Slug:         c.Slug,
			MatchedTitle: c.MatchedTitle,
			Score:        c.Score,
		})
	}

	httpjson.Write(w, http.StatusOK, animeSamaResolveResponse{Candidates: out})
}
