package httpapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/app"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/ports"
)

type Server struct {
	logger   zerolog.Logger
	jobs     *app.JobService
	settings *app.SettingsService
	bus      ports.EventBus
}

func NewServer(logger zerolog.Logger, jobs *app.JobService, settings *app.SettingsService, bus ports.EventBus) *Server {
	return &Server{logger: logger, jobs: jobs, settings: settings, bus: bus}
}

func (s *Server) Router() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(defaultRequestTimeout))
	r.Use(hlog.NewHandler(s.logger))
	r.Use(hlog.RequestIDHandler("request_id", "Request-Id"))
	r.Use(hlog.RemoteAddrHandler("remote_ip"))
	r.Use(hlog.UserAgentHandler("user_agent"))
	r.Use(hlog.AccessHandler(accessLogFn))

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", s.handleHealth)
		r.Get("/version", s.handleVersion)
		r.Get("/openapi.json", s.handleOpenAPI)
		r.Get("/events", s.handleEvents)

		if s.jobs != nil {
			NewJobsHandler(s.jobs).Routes(r)
		}
		if s.settings != nil {
			NewSettingsHandler(s.settings).Routes(r)
		}
	})

	return r
}
