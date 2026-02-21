package ports

import (
	"context"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
)

// AnimeDetailService provides anime detail information.
type AnimeDetailService interface {
	// GetDetail returns full details for a specific anime by ID.
	// Returns error if anime not found or context cancelled.
	GetDetail(ctx context.Context, id string) (domain.AnimeDetail, error)
}
