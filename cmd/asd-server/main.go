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
	
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})
	
	mux.HandleFunc("GET /api/search", func(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"results":[]}`))
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
