package ports

import (
	"context"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
)

type JobRepository interface {
	Create(ctx context.Context, job domain.Job) (domain.Job, error)
	Get(ctx context.Context, id string) (domain.Job, error)
	List(ctx context.Context, limit int) ([]domain.Job, error)
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
