package httpapi

import (
	"net/http"
	"time"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/buildinfo"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/httpjson"
	"github.com/rs/zerolog/hlog"
)

const defaultRequestTimeout = 30 * time.Second

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	httpjson.Write(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleVersion(w http.ResponseWriter, r *http.Request) {
	httpjson.Write(w, http.StatusOK, buildinfo.Current())
}

func accessLogFn(r *http.Request, status, size int, duration time.Duration) {
	logger := hlog.FromRequest(r)
	logger.Info().
		Int("status", status).
		Int("size", size).
		Dur("duration", duration).
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Msg("http")
}
