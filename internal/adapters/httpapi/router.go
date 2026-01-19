package httpapi

import (
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/app"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/ports"
)

type Server struct {
	logger   zerolog.Logger
	jobs     *app.JobService
	settings *app.SettingsService
	subs     *app.SubscriptionService
	anilist  *app.AniListService
	importer *app.AniListImportService
	resolver AnimeSamaResolver
	bus      ports.EventBus
	// downloadLimiter est optionnel et permet d'appliquer maxConcurrentDownloads à chaud.
	downloadLimiter *app.DynamicLimiter
	// onSettingsUpdated est optionnel (ex: ajuster maxWorkers).
	onSettingsUpdated func(domain.Settings)
}

func NewServer(logger zerolog.Logger, jobs *app.JobService, settings *app.SettingsService, subs *app.SubscriptionService, anilist *app.AniListService, importer *app.AniListImportService, resolver AnimeSamaResolver, bus ports.EventBus, downloadLimiter *app.DynamicLimiter, onSettingsUpdated func(domain.Settings)) *Server {
	return &Server{logger: logger, jobs: jobs, settings: settings, subs: subs, anilist: anilist, importer: importer, resolver: resolver, bus: bus, downloadLimiter: downloadLimiter, onSettingsUpdated: onSettingsUpdated}
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
		if s.subs != nil {
			NewSubscriptionsHandler(s.subs).Routes(r)
		}
		if s.anilist != nil {
			NewAniListHandler(s.anilist).Routes(r)
		}
		if s.importer != nil {
			NewAniListImportHandler(s.importer).Routes(r)
		}
		if s.resolver != nil {
			NewAnimeSamaHandler(s.resolver).Routes(r)
		}
		if s.settings != nil {
			NewSettingsHandler(s.settings, func(updated domain.Settings) {
				if s.downloadLimiter == nil {
					// noop
				} else if updated.MaxConcurrentDownloads > 0 {
					s.downloadLimiter.SetLimit(updated.MaxConcurrentDownloads)
				}
				if s.onSettingsUpdated != nil {
					s.onSettingsUpdated(updated)
				}
			}).Routes(r)
		}
	})

	// Compat: certaines UIs/clients historiques écoutent /api/events.
	r.Get("/api/events", s.handleEvents)

	// UI (best-effort): sert webapp/dist si présent (SPA fallback vers index.html).
	if h := newSPAHandlerFromEnvOrDefault(); h != nil {
		r.Mount("/", h)
	}

	return r
}

func newSPAHandlerFromEnvOrDefault() http.Handler {
	distDir := strings.TrimSpace(os.Getenv("ASD_WEB_DIST"))
	if distDir == "" {
		distDir = "webapp/dist"
	}
	if st, err := os.Stat(distDir); err != nil || !st.IsDir() {
		return nil
	}

	fs := http.FileServer(http.Dir(distDir))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ne jamais servir de fichiers pour les routes API.
		if strings.HasPrefix(r.URL.Path, "/api/") || strings.HasPrefix(r.URL.Path, "/api/v1/") {
			http.NotFound(w, r)
			return
		}

		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			http.NotFound(w, r)
			return
		}

		clean := path.Clean(r.URL.Path)
		if clean == "/" {
			fs.ServeHTTP(w, r)
			return
		}

		reqPath := strings.TrimPrefix(clean, "/")
		candidate := filepath.Join(distDir, filepath.FromSlash(reqPath))
		if st, err := os.Stat(candidate); err == nil && !st.IsDir() {
			fs.ServeHTTP(w, r)
			return
		}

		// Fallback SPA: servir index.html.
		r2 := *r
		r2.URL = newCopyURL(r.URL)
		r2.URL.Path = "/"
		fs.ServeHTTP(w, &r2)
	})
}

func newCopyURL(u *url.URL) *url.URL {
	if u == nil {
		return &url.URL{}
	}
	copy := *u
	return &copy
}
