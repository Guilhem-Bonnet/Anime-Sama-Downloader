package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
)

const settingsKey = "default"

type SettingsRepository struct {
	db *sql.DB
}

func NewSettingsRepository(db *sql.DB) *SettingsRepository {
	return &SettingsRepository{db: db}
}

func (r *SettingsRepository) Get(ctx context.Context) (domain.Settings, error) {
	var b []byte
	err := r.db.QueryRowContext(ctx, `SELECT value_json FROM settings WHERE key = ?`, settingsKey).Scan(&b)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Pas encore initialisé → valeurs par défaut.
			return domain.DefaultSettings(), nil
		}
		return domain.Settings{}, err
	}
	var s domain.Settings
	if err := json.Unmarshal(b, &s); err != nil {
		// Si corrompu : fallback safe.
		return domain.DefaultSettings(), nil
	}
	return s, nil
}

func (r *SettingsRepository) Put(ctx context.Context, settings domain.Settings) (domain.Settings, error) {
	b, err := json.Marshal(settings)
	if err != nil {
		return domain.Settings{}, err
	}
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO settings(key, value_json, updated_at)
		VALUES(?, ?, ?)
		ON CONFLICT(key) DO UPDATE SET value_json = excluded.value_json, updated_at = excluded.updated_at
	`, settingsKey, b, time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		return domain.Settings{}, err
	}
	return r.Get(ctx)
}
