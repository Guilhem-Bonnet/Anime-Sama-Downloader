package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/adapters/httpapi"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/adapters/memorybus"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/adapters/sqlite"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/app"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/buildinfo"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
)

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	logger.Info().Msg("starting anime-sama downloader MVP")

	// Check system dependencies
	ffmpegStatus := buildinfo.CheckFFmpeg()
	if !ffmpegStatus.Available {
		logger.Warn().Str("dependency", "ffmpeg").Str("message", ffmpegStatus.Message).Msg("ffmpeg not found - HLS/M3U8 downloads will not work")
	} else {
		logger.Info().Str("dependency", "ffmpeg").Str("version", ffmpegStatus.Version).Msg("ffmpeg available")
	}

	addr := envOr("ASD_ADDR", "127.0.0.1:8000")
	dbPath := envOr("ASD_DB_PATH", "asd.db")

	ctx := context.Background()
	db, err := sqlite.Open(ctx, dbPath)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to open sqlite database")
	}
	defer func() {
		_ = db.Close()
	}()

	bus := memorybus.New()

	jobsRepo := sqlite.NewJobsRepository(db.SQL)
	settingsRepo := sqlite.NewSettingsRepository(db.SQL)
	subsRepo := sqlite.NewSubscriptionsRepository(db.SQL)

	settingsService := app.NewSettingsService(settingsRepo)
	jobService := app.NewJobService(jobsRepo, bus)
	jobCreator := app.NewJobServiceAdapter(jobService)
	subscriptionsService := app.NewSubscriptionService(subsRepo, jobCreator, nil, bus)

	currentSettings, err := settingsService.Get(ctx)
	if err != nil {
		logger.Warn().Err(err).Msg("failed to load settings, using defaults")
		currentSettings = domain.DefaultSettings()
	}

	downloadLimiter := app.NewDynamicLimiter(currentSettings.MaxConcurrentDownloads)
	workerOpts := app.DefaultWorkerOptions()
	workerOpts.DownloadLimiter = downloadLimiter
	workerOpts.DestinationFunc = func(ctx context.Context) (string, error) {
		s, err := settingsService.Get(ctx)
		if err != nil {
			return domain.DefaultSettings().Destination, err
		}
		return s.Destination, nil
	}
	workerOpts.MaxConcurrentDownloadsFunc = func(ctx context.Context) (int, error) {
		s, err := settingsService.Get(ctx)
		if err != nil {
			return domain.DefaultSettings().MaxConcurrentDownloads, err
		}
		return s.MaxConcurrentDownloads, nil
	}

	workerPool := app.NewWorkerPool(ctx, logger, jobsRepo, bus, workerOpts)
	workerPool.SetCount(currentSettings.MaxWorkers)
	defer workerPool.Close()

	resolver := app.NewAnimeSamaCatalogueResolver()
	anilistService := app.NewAniListService(settingsService.Get)
	importer := app.NewAniListImportService(anilistService, resolver, subscriptionsService)

	_, animeCatalogue := devCatalogue()
	aniListHTTP := app.NewAniListHTTPClient("https://graphql.anilist.co")
	searchService := app.NewAniListSearchService(aniListHTTP)
	fileListService := app.NewFileListService(animeCatalogue)
	mockDetail := app.NewMockAnimeDetailService()
	detailService := app.NewAniListDetailService(aniListHTTP, mockDetail)
	recommendationsService := app.NewRecommendationsService(nil) // TODO: feed from AniList

	server := httpapi.NewServer(
		logger,
		jobService,
		settingsService,
		subscriptionsService,
		anilistService,
		importer,
		resolver,
		bus,
		downloadLimiter,
		searchService,
		detailService,
		fileListService,
		recommendationsService,
		func(updated domain.Settings) {
			if updated.MaxWorkers > 0 {
				workerPool.SetCount(updated.MaxWorkers)
			}
		},
	)

	httpServer := &http.Server{
		Addr:              addr,
		Handler:           server.Router(),
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	logger.Info().Str("addr", addr).Msg("listening")

	// Graceful shutdown handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		logger.Info().Str("signal", sig.String()).Msg("received shutdown signal, starting graceful shutdown")

		// Use 30s timeout for graceful shutdown (allows in-flight requests to complete)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		logger.Info().Msg("shutting down HTTP server...")
		if err := httpServer.Shutdown(ctx); err != nil {
			logger.Error().Err(err).Msg("HTTP server shutdown error")
		}
		logger.Info().Msg("graceful shutdown complete")
	}()

	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error().Err(err).Msg("server error")
		os.Exit(1)
	}
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func devCatalogue() ([]domain.AnimeSearchResult, []*domain.Anime) {
	search := []domain.AnimeSearchResult{
		{
			ID:           "mushishi",
			Title:        "Mushishi",
			ThumbnailURL: "/assets/cover-placeholder.svg",
			Year:         2005,
			Status:       "completed",
			EpisodeCount: 26,
			Genres:       []string{"Drama", "Mystery", "Supernatural"},
		},
		{
			ID:           "mononoke",
			Title:        "Mononoke",
			ThumbnailURL: "/assets/cover-placeholder.svg",
			Year:         2007,
			Status:       "completed",
			EpisodeCount: 12,
			Genres:       []string{"Horror", "Mystery"},
		},
		{
			ID:           "natsume-yuujinchou",
			Title:        "Natsume Yuujinchou",
			ThumbnailURL: "/assets/cover-placeholder.svg",
			Year:         2008,
			Status:       "ongoing",
			EpisodeCount: 13,
			Genres:       []string{"Slice of Life", "Supernatural"},
		},
		{
			ID:           "samurai-champloo",
			Title:        "Samurai Champloo",
			ThumbnailURL: "/assets/cover-placeholder.svg",
			Year:         2004,
			Status:       "completed",
			EpisodeCount: 26,
			Genres:       []string{"Action", "Adventure"},
		},
		{
			ID:           "dororo",
			Title:        "Dororo",
			ThumbnailURL: "/assets/cover-placeholder.svg",
			Year:         2019,
			Status:       "completed",
			EpisodeCount: 24,
			Genres:       []string{"Action", "Drama"},
		},
		{
			ID:           "spice-and-wolf",
			Title:        "Spice and Wolf",
			ThumbnailURL: "/assets/cover-placeholder.svg",
			Year:         2008,
			Status:       "completed",
			EpisodeCount: 13,
			Genres:       []string{"Fantasy", "Romance"},
		},
	}

	anime := make([]*domain.Anime, 0, len(search))
	for _, item := range search {
		anime = append(anime, &domain.Anime{
			ID:           item.ID,
			Title:        item.Title,
			ThumbnailURL: item.ThumbnailURL,
			Year:         item.Year,
			Status:       item.Status,
			EpisodeCount: item.EpisodeCount,
			Genres:       item.Genres,
		})
	}
	return search, anime
}
