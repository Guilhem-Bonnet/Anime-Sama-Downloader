package app

import (
	"context"

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
	return s.repo.Get(ctx)
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
