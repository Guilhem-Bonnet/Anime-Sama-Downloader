package ports

import (
	"context"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
)

type JobRepository interface {
	Create(ctx context.Context, job domain.Job) (domain.Job, error)
	Get(ctx context.Context, id string) (domain.Job, error)
	List(ctx context.Context, limit int) ([]domain.Job, error)
	// ClaimNextQueued passe le plus vieux job "queued" à l'état "running" et le renvoie.
	// Renvoie ErrNotFound (adapter-specific) s'il n'y a aucun job à exécuter.
	ClaimNextQueued(ctx context.Context) (domain.Job, error)
	UpdateProgress(ctx context.Context, id string, progress float64) (domain.Job, error)
	UpdateResult(ctx context.Context, id string, resultJSON []byte) (domain.Job, error)
	UpdateError(ctx context.Context, id string, code string, message string) (domain.Job, error)
	UpdateState(ctx context.Context, id string, expected domain.JobState, next domain.JobState) (domain.Job, error)
}

type EventBus interface {
	Publish(topic string, payload []byte)
	Subscribe() (ch <-chan Event, cancel func())
}

type Event struct {
	Topic   string
	Payload []byte
}
