package app

import (
	"context"
	"log/slog"
	"time"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
)

// DownloadService manages anime episode downloads.
type DownloadService struct {
	downloadRepo domain.IDownloadRepository
	eventBus     domain.IEventBus
	logger       *slog.Logger
}

// NewDownloadService creates a new DownloadService.
func NewDownloadService(downloadRepo domain.IDownloadRepository, eventBus domain.IEventBus, logger *slog.Logger) *DownloadService {
	return &DownloadService{
		downloadRepo: downloadRepo,
		eventBus:     eventBus,
		logger:       logger,
	}
}

// StartDownload enqueues a new download job.
func (ds *DownloadService) StartDownload(ctx context.Context, downloadID, animeID string, episodeNum int) (*domain.Download, error) {
	// Create download record
	download := &domain.Download{
		DownloadID: downloadID,
		AnimeID:    animeID,
		EpisodeNum: episodeNum,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// TODO: Save to repository

	// Emit event
	ds.eventBus.Emit(domain.EventDownloadQueued, map[string]interface{}{
		"download_id": downloadID,
		"anime_id":    animeID,
	})

	ds.logger.Info("download enqueued", slog.String("download_id", downloadID), slog.String("anime_id", animeID))
	return download, nil
}
