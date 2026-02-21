package app

import (
	"context"
	"os"
	"strings"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/ports"
)

type SettingsService struct {
	repo ports.SettingsRepository
}

func NewSettingsService(repo ports.SettingsRepository) *SettingsService {
	return &SettingsService{repo: repo}
}

func (s *SettingsService) Get(ctx context.Context) (domain.Settings, error) {
	settings, err := s.repo.Get(ctx)
	if err != nil {
		return domain.Settings{}, err
	}

	// Auto-migration: en Docker, on monte généralement /data/videos.
	// Si la destination est restée au défaut relatif "videos", on bascule vers /data/videos
	// pour coller à la config docker-compose (et éviter des fichiers perdus dans /app/videos).
	if dest, ok := autoDefaultDestination(settings.Destination); ok {
		settings.Destination = dest
		// Best-effort: on persiste pour que l'UI affiche la vraie destination.
		if updated, putErr := s.repo.Put(ctx, settings); putErr == nil {
			settings = updated
		}
	}

	return settings, nil
}

func autoDefaultDestination(current string) (string, bool) {
	cur := strings.TrimSpace(current)
	if cur != "" && cur != domain.DefaultSettings().Destination {
		return "", false
	}

	// Override explicite.
	if env := strings.TrimSpace(os.Getenv("ASD_DEFAULT_DESTINATION")); env != "" {
		return env, true
	}

	// Heuristique Docker: /data/videos existe généralement (volume mount).
	if _, err := os.Stat("/data/videos"); err == nil {
		return "/data/videos", true
	}

	return "", false
}

func (s *SettingsService) Put(ctx context.Context, settings domain.Settings) (domain.Settings, error) {
	// Validation légère V1.
	if settings.Destination == "" {
		settings.Destination = domain.DefaultSettings().Destination
	}
	if settings.OutputNamingMode == "" {
		settings.OutputNamingMode = domain.DefaultSettings().OutputNamingMode
	}
	if settings.MaxWorkers <= 0 {
		settings.MaxWorkers = domain.DefaultSettings().MaxWorkers
	}
	if settings.MaxConcurrentDownloads <= 0 {
		settings.MaxConcurrentDownloads = domain.DefaultSettings().MaxConcurrentDownloads
	}
	return s.repo.Put(ctx, settings)
}
