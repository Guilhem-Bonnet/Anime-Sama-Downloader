-- +migrate Up

CREATE TABLE IF NOT EXISTS settings (
  key TEXT PRIMARY KEY,
  value_json BLOB NOT NULL,
  updated_at TEXT NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS settings;
