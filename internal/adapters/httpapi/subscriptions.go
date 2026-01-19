package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/app"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/httpjson"
	"github.com/go-chi/chi/v5"
)

type SubscriptionsHandler struct {
	subs *app.SubscriptionService
}

func NewSubscriptionsHandler(subs *app.SubscriptionService) *SubscriptionsHandler {
	return &SubscriptionsHandler{subs: subs}
}

func (h *SubscriptionsHandler) Routes(r chi.Router) {
	r.Route("/subscriptions", func(r chi.Router) {
		r.Post("/", h.create)
		r.Get("/", h.list)
		r.Post("/sync-all", h.syncAll)
		r.Get("/{id}", h.get)
		r.Put("/{id}", h.update)
		r.Delete("/{id}", h.delete)
		r.Post("/{id}/sync", h.sync)
	})
}

type createSubscriptionRequest struct {
	BaseURL string `json:"baseUrl"`
	Label   string `json:"label"`
	Player  string `json:"player,omitempty"`
}

func (h *SubscriptionsHandler) create(w http.ResponseWriter, r *http.Request) {
	var req createSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	sub, err := h.subs.Create(r.Context(), req.BaseURL, req.Label, req.Player)
	if err != nil {
		if errors.Is(err, app.ErrConflict) {
			httpjson.WriteError(w, http.StatusConflict, "subscription already exists")
			return
		}
		httpjson.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	httpjson.Write(w, http.StatusCreated, sub)
}

func (h *SubscriptionsHandler) list(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	subs, err := h.subs.List(r.Context(), limit)
	if err != nil {
		httpjson.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpjson.Write(w, http.StatusOK, subs)
}

func (h *SubscriptionsHandler) get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	sub, err := h.subs.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, app.ErrNotFound) {
			httpjson.WriteError(w, http.StatusNotFound, "not found")
			return
		}
		httpjson.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpjson.Write(w, http.StatusOK, sub)
}

func (h *SubscriptionsHandler) update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var dto app.SubscriptionDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	dto.ID = id
	updated, err := h.subs.Update(r.Context(), dto)
	if err != nil {
		if errors.Is(err, app.ErrNotFound) {
			httpjson.WriteError(w, http.StatusNotFound, "not found")
			return
		}
		httpjson.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	httpjson.Write(w, http.StatusOK, updated)
}

func (h *SubscriptionsHandler) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.subs.Delete(r.Context(), id); err != nil {
		if errors.Is(err, app.ErrNotFound) {
			httpjson.WriteError(w, http.StatusNotFound, "not found")
			return
		}
		httpjson.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *SubscriptionsHandler) sync(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	enqueue := true
	if v := r.URL.Query().Get("enqueue"); v == "0" || v == "false" {
		enqueue = false
	}
	res, err := h.subs.SyncOnce(r.Context(), id, enqueue)
	if err != nil {
		if errors.Is(err, app.ErrNotFound) {
			httpjson.WriteError(w, http.StatusNotFound, "not found")
			return
		}
		httpjson.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	httpjson.Write(w, http.StatusOK, res)
}

type syncAllResponse struct {
	Results []app.SyncResult `json:"results"`
	Errors  []struct {
		ID    string `json:"id"`
		Error string `json:"error"`
	} `json:"errors"`
}

func (h *SubscriptionsHandler) syncAll(w http.ResponseWriter, r *http.Request) {
	enqueue := true
	if v := r.URL.Query().Get("enqueue"); v == "0" || v == "false" {
		enqueue = false
	}
	dueOnly := false
	if v := r.URL.Query().Get("dueOnly"); v == "1" || v == "true" {
		dueOnly = true
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	subs, err := h.subs.List(r.Context(), limit)
	if err != nil {
		httpjson.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	res := syncAllResponse{Results: []app.SyncResult{}, Errors: []struct {
		ID    string `json:"id"`
		Error string `json:"error"`
	}{} }

	now := time.Now().UTC()
	for _, sub := range subs {
		if dueOnly && sub.NextCheckAt.After(now) {
			continue
		}
		rr, err := h.subs.SyncOnce(r.Context(), sub.ID, enqueue)
		if err != nil {
			res.Errors = append(res.Errors, struct {
				ID    string `json:"id"`
				Error string `json:"error"`
			}{ID: sub.ID, Error: err.Error()})
			continue
		}
		res.Results = append(res.Results, rr)
	}

	httpjson.Write(w, http.StatusOK, res)
}
