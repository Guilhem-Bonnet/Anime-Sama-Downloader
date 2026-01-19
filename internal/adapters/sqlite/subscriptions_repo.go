package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/ports"
)

type SubscriptionsRepository struct {
	db *sql.DB
}

func NewSubscriptionsRepository(db *sql.DB) *SubscriptionsRepository {
	return &SubscriptionsRepository{db: db}
}

func (r *SubscriptionsRepository) Create(ctx context.Context, sub domain.Subscription) (domain.Subscription, error) {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO subscriptions(
			id, base_url, label, player,
			last_scheduled_episode, last_downloaded_episode, last_available_episode,
			next_check_at, last_checked_at,
			created_at, updated_at
		) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		sub.ID, sub.BaseURL, sub.Label, sub.Player,
		sub.LastScheduledEpisode, sub.LastDownloadedEpisode, sub.LastAvailableEpisode,
		sub.NextCheckAt.UTC().Format(time.RFC3339), sub.LastCheckedAt.UTC().Format(time.RFC3339),
		sub.CreatedAt.UTC().Format(time.RFC3339), sub.UpdatedAt.UTC().Format(time.RFC3339),
	)
	if err != nil {
		// modernc.org/sqlite retourne souvent une erreur texte du type:
		// "constraint failed: UNIQUE constraint failed: subscriptions.base_url (2067)"
		msg := strings.ToLower(err.Error())
		if strings.Contains(msg, "unique constraint failed") && strings.Contains(msg, "subscriptions.base_url") {
			return domain.Subscription{}, ports.ErrConflict
		}
		return domain.Subscription{}, err
	}
	return r.Get(ctx, sub.ID)
}

func (r *SubscriptionsRepository) Get(ctx context.Context, id string) (domain.Subscription, error) {
	var sub domain.Subscription
	var nextCheck, lastChecked, created, updated string
	err := r.db.QueryRowContext(ctx, `
		SELECT id, base_url, label, player,
			last_scheduled_episode, last_downloaded_episode, last_available_episode,
			next_check_at, last_checked_at,
			created_at, updated_at
		FROM subscriptions
		WHERE id = ?
	`, id).Scan(
		&sub.ID, &sub.BaseURL, &sub.Label, &sub.Player,
		&sub.LastScheduledEpisode, &sub.LastDownloadedEpisode, &sub.LastAvailableEpisode,
		&nextCheck, &lastChecked,
		&created, &updated,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Subscription{}, ports.ErrNotFound
		}
		return domain.Subscription{}, err
	}
	if t, err := time.Parse(time.RFC3339, nextCheck); err == nil {
		sub.NextCheckAt = t
	}
	if t, err := time.Parse(time.RFC3339, lastChecked); err == nil {
		sub.LastCheckedAt = t
	}
	if t, err := time.Parse(time.RFC3339, created); err == nil {
		sub.CreatedAt = t
	}
	if t, err := time.Parse(time.RFC3339, updated); err == nil {
		sub.UpdatedAt = t
	}
	return sub, nil
}

func (r *SubscriptionsRepository) List(ctx context.Context, limit int) ([]domain.Subscription, error) {
	q := `
		SELECT id FROM subscriptions
		ORDER BY updated_at DESC
	`
	args := []any{}
	if limit > 0 {
		q += ` LIMIT ?`
		args = append(args, limit)
	}

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := make([]string, 0)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	out := make([]domain.Subscription, 0, len(ids))
	for _, id := range ids {
		sub, err := r.Get(ctx, id)
		if err != nil {
			return nil, err
		}
		out = append(out, sub)
	}
	return out, nil
}

func (r *SubscriptionsRepository) Update(ctx context.Context, sub domain.Subscription) (domain.Subscription, error) {
	_, err := r.db.ExecContext(ctx, `
		UPDATE subscriptions
		SET base_url = ?, label = ?, player = ?,
			last_scheduled_episode = ?, last_downloaded_episode = ?, last_available_episode = ?,
			next_check_at = ?, last_checked_at = ?,
			updated_at = ?
		WHERE id = ?
	`,
		sub.BaseURL, sub.Label, sub.Player,
		sub.LastScheduledEpisode, sub.LastDownloadedEpisode, sub.LastAvailableEpisode,
		sub.NextCheckAt.UTC().Format(time.RFC3339), sub.LastCheckedAt.UTC().Format(time.RFC3339),
		sub.UpdatedAt.UTC().Format(time.RFC3339),
		sub.ID,
	)
	if err != nil {
		msg := strings.ToLower(err.Error())
		if strings.Contains(msg, "unique constraint failed") && strings.Contains(msg, "subscriptions.base_url") {
			return domain.Subscription{}, ports.ErrConflict
		}
		return domain.Subscription{}, err
	}
	return r.Get(ctx, sub.ID)
}

func (r *SubscriptionsRepository) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM subscriptions WHERE id = ?`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ports.ErrNotFound
	}
	return nil
}

func (r *SubscriptionsRepository) Due(ctx context.Context, now time.Time, limit int) ([]domain.Subscription, error) {
	q := `
		SELECT id FROM subscriptions
		WHERE next_check_at <= ?
		ORDER BY next_check_at ASC
	`
	args := []any{now.UTC().Format(time.RFC3339)}
	if limit > 0 {
		q += ` LIMIT ?`
		args = append(args, limit)
	}
	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	out := make([]domain.Subscription, 0, len(ids))
	for _, id := range ids {
		sub, err := r.Get(ctx, id)
		if err != nil {
			return nil, err
		}
		out = append(out, sub)
	}
	return out, nil
}

func (r *SubscriptionsRepository) MarkDownloadedEpisodeMax(ctx context.Context, id string, episode int) (domain.Subscription, error) {
	if episode <= 0 {
		return r.Get(ctx, id)
	}

	res, err := r.db.ExecContext(ctx, `
		UPDATE subscriptions
		SET last_downloaded_episode = CASE
			WHEN ? > last_downloaded_episode THEN ?
			ELSE last_downloaded_episode
		END,
		updated_at = ?
		WHERE id = ?
	`, episode, episode, time.Now().UTC().Format(time.RFC3339), id)
	if err != nil {
		return domain.Subscription{}, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.Subscription{}, ports.ErrNotFound
	}
	return r.Get(ctx, id)
}
