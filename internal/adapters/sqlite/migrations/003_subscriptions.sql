-- +migrate Up

CREATE TABLE IF NOT EXISTS subscriptions (
  id TEXT PRIMARY KEY,
  base_url TEXT NOT NULL,
  label TEXT NOT NULL,
  player TEXT NOT NULL,
  last_scheduled_episode INTEGER NOT NULL DEFAULT 0,
  last_downloaded_episode INTEGER NOT NULL DEFAULT 0,
  last_available_episode INTEGER NOT NULL DEFAULT 0,
  next_check_at TEXT NOT NULL,
  last_checked_at TEXT NOT NULL,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_subscriptions_next_check_at ON subscriptions(next_check_at);

-- +migrate Down

DROP TABLE IF EXISTS subscriptions;
