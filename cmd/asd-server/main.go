package main

import (
"context"
"fmt"
"log/slog"
"net/http"
"os"
"os/signal"
"syscall"
"time"

"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/adapters/memorybus"
"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/app"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("starting anime-sama downloader MVP")

	eventBus := memorybus.NewEventBus()
	_ = app.NewSearchService(nil, logger)
	_ = app.NewDownloadService(nil, eventBus, logger)
	jobWorker := app.NewJobWorker(nil, eventBus, logger)

	go func() {
		ctx := context.Background()
		jobWorker.Start(ctx)
	}()

	mux := http.NewServeMux()
	
	// CORS middleware for development
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
	
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})
	
	mux.HandleFunc("GET /api/search", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		mockData := `{
			"results": [
				{
					"anime_id": "attack-on-titan",
		w.Header().Set("Content-Type", "application/json")
		mockData := `{
			"downloads": [
				{
					"download_id": "dl-001",
					"anime_id": "attack-on-titan",
					"episode_number": 1,
					"status": "completed",
					"progress": 100
				},
				{
					"download_id": "dl-002",
					"anime_id": "demon-slayer",
					"episode_number": 1,
					"status": "running",
					"progress": 65
				},
				{
					"download_id": "dl-003",
					"anime_id": "jujutsu-kaisen",
					"episode_number": 1,
					"status": "pending",
					"progress": 0
				}
			]
		}`
		w.Write([]byte(mockData
					"source": "AnimeSama",
					"image_url": "https://cdn.myanimelist.net/images/anime/10/47347.jpg"
				},
				{
					"anime_id": "demon-slayer",
					"title": "Demon Slayer: Kimetsu no Yaiba",
					"episodes": 26,
					"source": "AnimeSama",
					"image_url": "https://cdn.myanimelist.net/images/anime/1286/99889.jpg"
				},
				{
					"anime_id": "jujutsu-kaisen",
					"titlecorsMiddleware(mux)"Jujutsu Kaisen",
					"episodes": 24,
					"source": "AnimeSama",
					"image_url": "https://cdn.myanimelist.net/images/anime/1171/109222.jpg"
				},
				{
					"anime_id": "my-hero-academia",
					"title": "My Hero Academia",
					"episodes": 113,
					"source": "AnimeSama",
					"image_url": "https://cdn.myanimelist.net/images/anime/10/78745.jpg"
				},
				{
					"anime_id": "one-piece",
					"title": "One Piece",
					"episodes": 1000,
					"source": "AnimeSama",
					"image_url": "https://cdn.myanimelist.net/images/anime/6/73245.jpg"
				},
				{
					"anime_id": "naruto",
					"title": "Naruto Shippuden",
					"episodes": 500,
					"source": "AnimeSama",
					"image_url": "https://cdn.myanimelist.net/images/anime/5/17407.jpg"
				}
			]
		}`
		w.Write([]byte(mockData))
	})
	
	mux.HandleFunc("GET /api/downloads", func(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"downloads":[]}`))
	})

	addr := ":8000"
	httpServer := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	logger.Info(fmt.Sprintf("listening on %s", addr))

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		httpServer.Shutdown(ctx)
	}()

	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("server error", "error", err.Error())
		os.Exit(1)
	}
}
