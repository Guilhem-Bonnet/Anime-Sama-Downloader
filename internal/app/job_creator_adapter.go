package app

import (
	"context"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/ports"
)

// JobServiceAdapter adapts JobService to implement ports.JobCreator interface.
// This allows SubscriptionService to depend on the JobCreator interface
// instead of the concrete JobService implementation.
type JobServiceAdapter struct {
	service *JobService
}

// NewJobServiceAdapter creates an adapter for JobService.
func NewJobServiceAdapter(service *JobService) *JobServiceAdapter {
	if service == nil {
		return nil
	}
	return &JobServiceAdapter{service: service}
}

// Create creates a new job and returns the domain.Job.
// It implements ports.JobCreator interface.
func (a *JobServiceAdapter) Create(ctx context.Context, req ports.JobCreationRequest) (domain.Job, error) {
	internalReq := CreateJobRequest{
		Type:   req.Type,
		Params: req.Params,
	}
	dto, err := a.service.Create(ctx, internalReq)
	if err != nil {
		return domain.Job{}, err
	}
	// Convert JobDTO back to domain.Job
	return domain.Job{
		ID:           dto.ID,
		Type:         dto.Type,
		State:        dto.State,
		Progress:     dto.Progress,
		CreatedAt:    dto.CreatedAt,
		UpdatedAt:    dto.UpdatedAt,
		ParamsJSON:   []byte(dto.Params),
		ResultJSON:   []byte(dto.Result),
		ErrorCode:    dto.ErrorCode,
		ErrorMessage: dto.Error,
	}, nil
}
