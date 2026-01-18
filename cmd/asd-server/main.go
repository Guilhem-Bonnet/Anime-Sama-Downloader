package main

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/adapters/httpapi"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/adapters/memorybus"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/adapters/sqlite"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/app"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/buildinfo"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	def := config.Default()
	addr := flag.String("addr", def.Addr, "Adresse d'écoute (ex: 127.0.0.1:8080)")
	dbPath := flag.String("db", def.DBPath, "Chemin SQLite (ex: asd.db)")
	flag.Parse()

	logger := zerolog.New(os.Stdout).With().Timestamp().Str("app", "asd-server").Logger()
	log.Logger = logger

	logger.Info().Interface("build", buildinfo.Current()).Str("db", *dbPath).Msg("starting")

	ctx := context.Background()
	db, err := sqlite.Open(ctx, *dbPath)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to open db")
	}
	defer func() { _ = db.Close() }()

	bus := memorybus.New()
	jobsRepo := sqlite.NewJobsRepository(db.SQL)
	jobsSvc := app.NewJobService(jobsRepo, bus)
	settingsRepo := sqlite.NewSettingsRepository(db.SQL)
	settingsSvc := app.NewSettingsService(settingsRepo)

	srv := httpapi.NewServer(logger, jobsSvc, settingsSvc, bus)
	httpServer := &http.Server{
		Addr:              *addr,
		Handler:           srv.Router(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	shutdownCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Worker V1: exécute les jobs "queued" en tâche de fond.
	workers := 1
	opts := app.DefaultWorkerOptions()
	if s, err := settingsSvc.Get(ctx); err == nil {
		if s.MaxWorkers > 0 {
			workers = s.MaxWorkers
		}
		if s.Destination != "" {
			opts.Destination = s.Destination
		}
	}
	app.RunWorkers(shutdownCtx, logger, jobsRepo, bus, workers, opts)
	logger.Info().Int("workers", workers).Msg("workers started")

	go func() {
		logger.Info().Str("addr", *addr).Msg("listening")
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error().Err(err).Msg("http server crashed")
			stop()
		}
	}()

	<-shutdownCtx.Done()
	logger.Info().Msg("shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = httpServer.Shutdown(ctx)
	logger.Info().Msg("bye")
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
