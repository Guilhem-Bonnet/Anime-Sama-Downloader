package domain

import "context"

// IDownloadRepository defines data access for downloads.
type IDownloadRepository interface {
	Create(ctx context.Context, download *Download) error
	GetByID(ctx context.Context, downloadID string) (*Download, error)
	List(ctx context.Context) ([]*Download, error)
	Update(ctx context.Context, download *Download) error
	Delete(ctx context.Context, downloadID string) error
}

// IJobRepository defines data access for jobs.
type IJobRepository interface {
	Create(ctx context.Context, job *Job) error
	GetByID(ctx context.Context, jobID string) (*Job, error)
	List(ctx context.Context) ([]*Job, error)
	Update(ctx context.Context, job *Job) error
	Delete(ctx context.Context, jobID string) error
}

// ISettingsRepository defines data access for settings.
type ISettingsRepository interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string) error
	Delete(ctx context.Context, key string) error
}
