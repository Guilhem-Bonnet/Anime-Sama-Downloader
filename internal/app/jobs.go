package app

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/ports"
	"github.com/rs/xid"
)

type JobService struct {
	repo ports.JobRepository
	bus  ports.EventBus
}

func NewJobService(repo ports.JobRepository, bus ports.EventBus) *JobService {
	return &JobService{repo: repo, bus: bus}
}

type CreateJobRequest struct {
	Type   string          `json:"type"`
	Params json.RawMessage `json:"params,omitempty"`
}

type JobDTO struct {
	ID        string          `json:"id"`
	Type      string          `json:"type"`
	State     domain.JobState `json:"state"`
	Progress  float64         `json:"progress"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
	Params    json.RawMessage `json:"params,omitempty"`
}

func toDTO(j domain.Job) JobDTO {
	return JobDTO{
		ID:        j.ID,
		Type:      j.Type,
		State:     j.State,
		Progress:  j.Progress,
		CreatedAt: j.CreatedAt,
		UpdatedAt: j.UpdatedAt,
		Params:    json.RawMessage(j.ParamsJSON),
	}
}

func (s *JobService) Create(ctx context.Context, req CreateJobRequest) (JobDTO, error) {
	now := time.Now().UTC()
	job := domain.Job{
		ID:         xid.New().String(),
		Type:       req.Type,
		State:      domain.JobQueued,
		Progress:   0,
		CreatedAt:  now,
		UpdatedAt:  now,
		ParamsJSON: []byte(req.Params),
	}
	created, err := s.repo.Create(ctx, job)
	if err != nil {
		return JobDTO{}, err
	}
	s.publish("job.created", created)
	return toDTO(created), nil
}

func (s *JobService) Get(ctx context.Context, id string) (JobDTO, error) {
	job, err := s.repo.Get(ctx, id)
	if err != nil {
		return JobDTO{}, err
	}
	return toDTO(job), nil
}

func (s *JobService) List(ctx context.Context, limit int) ([]JobDTO, error) {
	jobs, err := s.repo.List(ctx, limit)
	if err != nil {
		return nil, err
	}
	out := make([]JobDTO, 0, len(jobs))
	for _, j := range jobs {
		out = append(out, toDTO(j))
	}
	return out, nil
}

func (s *JobService) Cancel(ctx context.Context, id string) (JobDTO, error) {
	// V1: on autorise l'annulation depuis queued/running/muxing.
	// On essaie en cascade.
	for _, expected := range []domain.JobState{domain.JobQueued, domain.JobRunning, domain.JobMuxing} {
		updated, err := s.repo.UpdateState(ctx, id, expected, domain.JobCanceled)
		if err == nil {
			s.publish("job.canceled", updated)
			return toDTO(updated), nil
		}
	}
	// fallback: renvoyer l'état actuel
	job, err := s.repo.Get(ctx, id)
	if err != nil {
		return JobDTO{}, err
	}
	return toDTO(job), nil
}

func (s *JobService) publish(topic string, job domain.Job) {
	if s.bus == nil {
		return
	}
	b, err := json.Marshal(toDTO(job))
	if err != nil {
		return
	}
	s.bus.Publish(topic, b)
}
