package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
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
}

func NewAnimeSamaHandler(resolver AnimeSamaResolver) *AnimeSamaHandler {
	return &AnimeSamaHandler{resolver: resolver}
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
	})
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
