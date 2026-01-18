-- +migrate Up

CREATE TABLE IF NOT EXISTS schema_migrations (
  version INTEGER PRIMARY KEY,
  applied_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS jobs (
  id TEXT PRIMARY KEY,
  type TEXT NOT NULL,
  state TEXT NOT NULL,
  progress REAL NOT NULL DEFAULT 0,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL,
  params_json BLOB,
  result_json BLOB,
  error_code TEXT,
  error_message TEXT
);

CREATE INDEX IF NOT EXISTS idx_jobs_updated_at ON jobs(updated_at);

-- +migrate Down
DROP TABLE IF EXISTS jobs;
DROP TABLE IF EXISTS schema_migrations;
