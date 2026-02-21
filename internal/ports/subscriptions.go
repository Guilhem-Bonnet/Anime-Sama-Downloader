package ports

import (
	"context"
	"time"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, sub domain.Subscription) (domain.Subscription, error)
	Get(ctx context.Context, id string) (domain.Subscription, error)
	List(ctx context.Context, limit int) ([]domain.Subscription, error)
	Update(ctx context.Context, sub domain.Subscription) (domain.Subscription, error)
	Delete(ctx context.Context, id string) error
	Due(ctx context.Context, now time.Time, limit int) ([]domain.Subscription, error)
	// MarkDownloadedEpisodeMax met à jour lastDownloadedEpisode de façon atomique:
	// lastDownloadedEpisode = max(lastDownloadedEpisode, episode).
	MarkDownloadedEpisodeMax(ctx context.Context, id string, episode int) (domain.Subscription, error)
}
