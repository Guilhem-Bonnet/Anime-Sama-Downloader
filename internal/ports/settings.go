package ports

import (
	"context"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
)

type SettingsRepository interface {
	Get(ctx context.Context) (domain.Settings, error)
	Put(ctx context.Context, settings domain.Settings) (domain.Settings, error)
}
