package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/app"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/httpjson"
	"github.com/go-chi/chi/v5"
)

type SettingsHandler struct {
	settings *app.SettingsService
	onPut    func(domain.Settings)
}

func NewSettingsHandler(settings *app.SettingsService, onPut func(domain.Settings)) *SettingsHandler {
	return &SettingsHandler{settings: settings, onPut: onPut}
}

func (h *SettingsHandler) Routes(r chi.Router) {
	r.Get("/settings", h.get)
	r.Put("/settings", h.put)
	// Variante avec slash final (utile selon reverse-proxy / clients).
	r.Get("/settings/", h.get)
	r.Put("/settings/", h.put)
}

func (h *SettingsHandler) get(w http.ResponseWriter, r *http.Request) {
	s, err := h.settings.Get(r.Context())
	if err != nil {
		httpjson.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpjson.Write(w, http.StatusOK, s)
}

func (h *SettingsHandler) put(w http.ResponseWriter, r *http.Request) {
	var s domain.Settings
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	updated, err := h.settings.Put(r.Context(), s)
	if err != nil {
		httpjson.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if h.onPut != nil {
		h.onPut(updated)
	}
	httpjson.Write(w, http.StatusOK, updated)
}
