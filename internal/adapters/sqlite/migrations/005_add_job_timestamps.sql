-- +migrate Up

ALTER TABLE jobs ADD COLUMN started_at TEXT;
ALTER TABLE jobs ADD COLUMN completed_at TEXT;

CREATE INDEX IF NOT EXISTS idx_jobs_state ON jobs(state);

-- +migrate Down

-- SQLite ne supporte pas DROP COLUMN avant 3.35.0
-- Les colonnes resteront mais seront inutilisées si on rollback
