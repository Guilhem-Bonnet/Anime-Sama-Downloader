package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/app"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/httpjson"
	"github.com/go-chi/chi/v5"
)

type AniListImportHandler struct {
	svc *app.AniListImportService
}

func NewAniListImportHandler(svc *app.AniListImportService) *AniListImportHandler {
	return &AniListImportHandler{svc: svc}
}

func (h *AniListImportHandler) Routes(r chi.Router) {
	r.Route("/import/anilist", func(r chi.Router) {
		r.Post("/preview", h.preview)
		r.Post("/auto", h.auto)
		r.Post("/confirm", h.confirm)
	})
}

func (h *AniListImportHandler) preview(w http.ResponseWriter, r *http.Request) {
	if h.svc == nil {
		httpjson.WriteError(w, http.StatusNotImplemented, "import disabled")
		return
	}
	var req app.AniListImportPreviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	res, err := h.svc.Preview(r.Context(), req)
	if err != nil {
		if err == app.ErrAniListNotConfigured {
			httpjson.WriteError(w, http.StatusBadRequest, "anilist not configured (set settings.anilistToken)")
			return
		}
		httpjson.WriteError(w, http.StatusBadGateway, err.Error())
		return
	}
	httpjson.Write(w, http.StatusOK, res)
}

func (h *AniListImportHandler) confirm(w http.ResponseWriter, r *http.Request) {
	if h.svc == nil {
		httpjson.WriteError(w, http.StatusNotImplemented, "import disabled")
		return
	}
	var req app.AniListImportConfirmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	res, err := h.svc.Confirm(r.Context(), req)
	if err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	httpjson.Write(w, http.StatusOK, res)
}

func (h *AniListImportHandler) auto(w http.ResponseWriter, r *http.Request) {
	if h.svc == nil {
		httpjson.WriteError(w, http.StatusNotImplemented, "import disabled")
		return
	}
	var req app.AniListImportAutoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	res, err := h.svc.AutoConfirm(r.Context(), req)
	if err != nil {
		if err == app.ErrAniListNotConfigured {
			httpjson.WriteError(w, http.StatusBadRequest, "anilist not configured (set settings.anilistToken)")
			return
		}
		httpjson.WriteError(w, http.StatusBadGateway, err.Error())
		return
	}
	httpjson.Write(w, http.StatusOK, res)
}
