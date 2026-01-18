package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/ports"
)

type JobsRepository struct {
	db *sql.DB
}

func NewJobsRepository(db *sql.DB) *JobsRepository {
	return &JobsRepository{db: db}
}

func (r *JobsRepository) Create(ctx context.Context, job domain.Job) (domain.Job, error) {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO jobs(id, type, state, progress, created_at, updated_at, params_json, result_json, error_code, error_message)
		VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, job.ID, job.Type, string(job.State), job.Progress,
		job.CreatedAt.Format(time.RFC3339), job.UpdatedAt.Format(time.RFC3339), job.ParamsJSON, job.ResultJSON, job.ErrorCode, job.ErrorMessage)
	if err != nil {
		return domain.Job{}, err
	}
	return r.Get(ctx, job.ID)
}

func (r *JobsRepository) Get(ctx context.Context, id string) (domain.Job, error) {
	var j domain.Job
	var createdAt, updatedAt string
	err := r.db.QueryRowContext(ctx, `
		SELECT id, type, state, progress, created_at, updated_at, params_json, result_json, error_code, error_message
		FROM jobs WHERE id = ?
	`, id).Scan(&j.ID, &j.Type, &j.State, &j.Progress, &createdAt, &updatedAt, &j.ParamsJSON, &j.ResultJSON, &j.ErrorCode, &j.ErrorMessage)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Job{}, ports.ErrNotFound
		}
		return domain.Job{}, err
	}
	j.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	j.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
	return j, nil
}

func (r *JobsRepository) List(ctx context.Context, limit int) ([]domain.Job, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, type, state, progress, created_at, updated_at, params_json, result_json, error_code, error_message
		FROM jobs ORDER BY updated_at DESC LIMIT ?
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []domain.Job{}
	for rows.Next() {
		var j domain.Job
		var createdAt, updatedAt string
		if err := rows.Scan(&j.ID, &j.Type, &j.State, &j.Progress, &createdAt, &updatedAt, &j.ParamsJSON, &j.ResultJSON, &j.ErrorCode, &j.ErrorMessage); err != nil {
			return nil, err
		}
		j.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		j.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
		out = append(out, j)
	}
	return out, rows.Err()
}

func (r *JobsRepository) ClaimNextQueued(ctx context.Context) (domain.Job, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.Job{}, err
	}
	defer func() { _ = tx.Rollback() }()

	var id string
	err = tx.QueryRowContext(ctx, `
		SELECT id
		FROM jobs
		WHERE state = ?
		ORDER BY created_at ASC
		LIMIT 1
	`, string(domain.JobQueued)).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Job{}, ports.ErrNotFound
		}
		return domain.Job{}, err
	}

	res, err := tx.ExecContext(ctx, `
		UPDATE jobs
		SET state = ?, updated_at = ?
		WHERE id = ? AND state = ?
	`, string(domain.JobRunning), time.Now().UTC().Format(time.RFC3339), id, string(domain.JobQueued))
	if err != nil {
		return domain.Job{}, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.Job{}, ports.ErrNotFound
	}
	if err := tx.Commit(); err != nil {
		return domain.Job{}, err
	}
	return r.Get(ctx, id)
}

func (r *JobsRepository) UpdateProgress(ctx context.Context, id string, progress float64) (domain.Job, error) {
	res, err := r.db.ExecContext(ctx, `
		UPDATE jobs
		SET progress = ?, updated_at = ?
		WHERE id = ?
	`, progress, time.Now().UTC().Format(time.RFC3339), id)
	if err != nil {
		return domain.Job{}, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.Job{}, ports.ErrNotFound
	}
	return r.Get(ctx, id)
}

func (r *JobsRepository) UpdateResult(ctx context.Context, id string, resultJSON []byte) (domain.Job, error) {
	res, err := r.db.ExecContext(ctx, `
		UPDATE jobs
		SET result_json = ?, updated_at = ?
		WHERE id = ?
	`, resultJSON, time.Now().UTC().Format(time.RFC3339), id)
	if err != nil {
		return domain.Job{}, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.Job{}, ports.ErrNotFound
	}
	return r.Get(ctx, id)
}

func (r *JobsRepository) UpdateError(ctx context.Context, id string, code string, message string) (domain.Job, error) {
	res, err := r.db.ExecContext(ctx, `
		UPDATE jobs
		SET error_code = ?, error_message = ?, updated_at = ?
		WHERE id = ?
	`, code, message, time.Now().UTC().Format(time.RFC3339), id)
	if err != nil {
		return domain.Job{}, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.Job{}, ports.ErrNotFound
	}
	return r.Get(ctx, id)
}

func (r *JobsRepository) UpdateState(ctx context.Context, id string, expected domain.JobState, next domain.JobState) (domain.Job, error) {
	if !domain.CanTransition(expected, next) {
		return domain.Job{}, domain.ErrInvalidTransition
	}
	res, err := r.db.ExecContext(ctx, `
		UPDATE jobs
		SET state = ?, updated_at = ?
		WHERE id = ? AND state = ?
	`, string(next), time.Now().UTC().Format(time.RFC3339), id, string(expected))
	if err != nil {
		return domain.Job{}, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.Job{}, ports.ErrNotFound
	}
	return r.Get(ctx, id)
}
