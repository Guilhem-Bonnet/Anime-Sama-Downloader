package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/adapters/memorybus"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/app"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/buildinfo"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("starting anime-sama downloader MVP")

	// Check system dependencies
	ffmpegStatus := buildinfo.CheckFFmpeg()
	if !ffmpegStatus.Available {
		logger.Warn("ffmpeg not found - HLS/M3U8 downloads will not work",
			"dependency", "ffmpeg",
			"message", ffmpegStatus.Message)
	} else {
		logger.Info("ffmpeg available",
			"dependency", "ffmpeg",
			"version", ffmpegStatus.Version)
	}

	eventBus := memorybus.NewEventBus()
	_ = app.NewSearchService(nil, logger)
	_ = app.NewDownloadService(nil, eventBus, logger)
	jobWorker := app.NewJobWorker(nil, eventBus, logger)

	go func() {
		ctx := context.Background()
		jobWorker.Start(ctx)
	}()

	mux := http.NewServeMux()

	// CORS middleware
	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	}

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Search endpoint with mock data
	mux.HandleFunc("GET /api/search", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		mockData := `{
			"results": [
				{"anime_id":"mushishi","title":"Mushishi","episodes":26,"source":"AnimeSama","image_url":"/assets/cover-placeholder.svg"},
				{"anime_id":"mononoke","title":"Mononoke","episodes":12,"source":"AnimeSama","image_url":"/assets/cover-placeholder.svg"},
				{"anime_id":"natsume-yuujinchou","title":"Natsume Yuujinchou","episodes":13,"source":"AnimeSama","image_url":"/assets/cover-placeholder.svg"},
				{"anime_id":"samurai-champloo","title":"Samurai Champloo","episodes":26,"source":"AnimeSama","image_url":"/assets/cover-placeholder.svg"},
				{"anime_id":"dororo","title":"Dororo","episodes":24,"source":"AnimeSama","image_url":"/assets/cover-placeholder.svg"},
				{"anime_id":"spice-and-wolf","title":"Spice and Wolf","episodes":13,"source":"AnimeSama","image_url":"/assets/cover-placeholder.svg"}
			]
		}`
		w.Write([]byte(mockData))
	})

	// Downloads list
	mux.HandleFunc("GET /api/downloads", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		mockData := `{
			"downloads": [
				{"download_id":"dl-001","anime_id":"attack-on-titan","episode_number":1,"status":"completed","progress":100},
				{"download_id":"dl-002","anime_id":"demon-slayer","episode_number":1,"status":"running","progress":65},
				{"download_id":"dl-003","anime_id":"jujutsu-kaisen","episode_number":1,"status":"pending","progress":0}
			]
		}`
		w.Write([]byte(mockData))
	})

	// Create download
	mux.HandleFunc("POST /api/downloads", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			AnimeID       string `json:"anime_id"`
			EpisodeNumber int    `json:"episode_number"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		resp := fmt.Sprintf(`{"download_id":"dl-%d","status":"queued"}`, time.Now().Unix())
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(resp))
	})

	addr := ":8000"
	httpServer := &http.Server{
		Addr:    addr,
		Handler: corsMiddleware(mux),
	}

	logger.Info(fmt.Sprintf("listening on %s", addr))

	// Graceful shutdown handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		logger.Info("received shutdown signal, starting graceful shutdown",
			"signal", sig.String())

		// Use 30s timeout for graceful shutdown (allows in-flight requests to complete)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		logger.Info("shutting down HTTP server...")
		if err := httpServer.Shutdown(ctx); err != nil {
			logger.Error("HTTP server shutdown error", "error", err.Error())
		}

		logger.Info("graceful shutdown complete")
	}()

	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("server error", "error", err.Error())
		os.Exit(1)
	}
}
