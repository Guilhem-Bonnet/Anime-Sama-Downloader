package httpapi

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/app"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/httpjson"
	"github.com/go-chi/chi/v5"
)

type AniListHandler struct {
	svc *app.AniListService
}

func NewAniListHandler(svc *app.AniListService) *AniListHandler {
	return &AniListHandler{svc: svc}
}

func (h *AniListHandler) Routes(r chi.Router) {
	r.Route("/anilist", func(r chi.Router) {
		r.Get("/viewer", h.viewer)
		r.Get("/airing", h.airing)
		r.Get("/watchlist", h.watchlist)
	})
}

func (h *AniListHandler) viewer(w http.ResponseWriter, r *http.Request) {
	if h.svc == nil {
		httpjson.WriteError(w, http.StatusNotImplemented, "anilist disabled")
		return
	}
	viewer, err := h.svc.Viewer(r.Context())
	if err != nil {
		if err == app.ErrAniListNotConfigured {
			httpjson.WriteError(w, http.StatusBadRequest, "anilist not configured (set settings.anilistToken)")
			return
		}
		httpjson.WriteError(w, http.StatusBadGateway, err.Error())
		return
	}
	httpjson.Write(w, http.StatusOK, viewer)
}

func (h *AniListHandler) airing(w http.ResponseWriter, r *http.Request) {
	if h.svc == nil {
		httpjson.WriteError(w, http.StatusNotImplemented, "anilist disabled")
		return
	}
	days, _ := strconv.Atoi(r.URL.Query().Get("days"))
	if days <= 0 {
		days = 7
	}
	if days > 30 {
		days = 30
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	now := time.Now().UTC()
	entries, err := h.svc.AiringSchedule(r.Context(), now.Add(-15*time.Minute), now.Add(time.Duration(days)*24*time.Hour), limit)
	if err != nil {
		httpjson.WriteError(w, http.StatusBadGateway, err.Error())
		return
	}
	httpjson.Write(w, http.StatusOK, entries)
}

func (h *AniListHandler) watchlist(w http.ResponseWriter, r *http.Request) {
	if h.svc == nil {
		httpjson.WriteError(w, http.StatusNotImplemented, "anilist disabled")
		return
	}
	statuses := r.URL.Query()["status"]
	entries, err := h.svc.Watchlist(r.Context(), statuses)
	if err != nil {
		if err == app.ErrAniListNotConfigured {
			httpjson.WriteError(w, http.StatusBadRequest, "anilist not configured (set settings.anilistToken)")
			return
		}
		httpjson.WriteError(w, http.StatusBadGateway, err.Error())
		return
	}
	httpjson.Write(w, http.StatusOK, entries)
}
