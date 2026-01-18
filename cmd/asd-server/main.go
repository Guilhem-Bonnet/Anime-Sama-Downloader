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
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/buildinfo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	addr := flag.String("addr", envOr("ASD_ADDR", "127.0.0.1:8080"), "Adresse d'écoute (ex: 127.0.0.1:8080)")
	flag.Parse()

	logger := zerolog.New(os.Stdout).With().Timestamp().Str("app", "asd-server").Logger()
	log.Logger = logger

	logger.Info().Interface("build", buildinfo.Current()).Msg("starting")

	srv := httpapi.NewServer(logger)
	httpServer := &http.Server{
		Addr:              *addr,
		Handler:           srv.Router(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	shutdownCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

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
