package ports

import (
	"context"
	"encoding/json"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
)

// JobCreator defines the interface for creating download jobs.
// It abstracts the job creation mechanism to allow dependency injection
// and easier testing (mock implementations).
type JobCreator interface {
	Create(ctx context.Context, req JobCreationRequest) (domain.Job, error)
}

// JobCreationRequest is the request structure for creating jobs.
type JobCreationRequest struct {
	Type   string          `json:"type"`
	Params json.RawMessage `json:"params,omitempty"`
}
