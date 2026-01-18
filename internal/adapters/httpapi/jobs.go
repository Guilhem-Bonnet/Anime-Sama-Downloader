package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/app"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/httpjson"
	"github.com/go-chi/chi/v5"
)

type JobsHandler struct {
	jobs *app.JobService
}

func NewJobsHandler(jobs *app.JobService) *JobsHandler {
	return &JobsHandler{jobs: jobs}
}

func (h *JobsHandler) Routes(r chi.Router) {
	r.Route("/jobs", func(r chi.Router) {
		r.Post("/", h.create)
		r.Get("/", h.list)
		r.Get("/{id}", h.get)
		r.Post("/{id}/cancel", h.cancel)
	})
}

func (h *JobsHandler) create(w http.ResponseWriter, r *http.Request) {
	var req app.CreateJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if req.Type == "" {
		httpjson.WriteError(w, http.StatusBadRequest, "missing type")
		return
	}

	job, err := h.jobs.Create(r.Context(), req)
	if err != nil {
		httpjson.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpjson.Write(w, http.StatusCreated, job)
}

func (h *JobsHandler) list(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	jobs, err := h.jobs.List(r.Context(), limit)
	if err != nil {
		httpjson.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpjson.Write(w, http.StatusOK, jobs)
}

func (h *JobsHandler) get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	job, err := h.jobs.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, app.ErrNotFound) {
			httpjson.WriteError(w, http.StatusNotFound, "not found")
			return
		}
		httpjson.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpjson.Write(w, http.StatusOK, job)
}

func (h *JobsHandler) cancel(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	job, err := h.jobs.Cancel(r.Context(), id)
	if err != nil {
		if errors.Is(err, app.ErrNotFound) {
			httpjson.WriteError(w, http.StatusNotFound, "not found")
			return
		}
		httpjson.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpjson.Write(w, http.StatusOK, job)
}
