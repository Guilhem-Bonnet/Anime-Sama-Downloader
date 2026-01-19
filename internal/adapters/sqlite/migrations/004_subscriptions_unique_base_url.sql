-- +migrate Up

-- Best-effort deduplication before enforcing uniqueness.
-- Keep the earliest rowid for each base_url.
DELETE FROM subscriptions
WHERE rowid NOT IN (
  SELECT MIN(rowid) FROM subscriptions GROUP BY base_url
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_subscriptions_base_url ON subscriptions(base_url);

-- +migrate Down

DROP INDEX IF EXISTS ux_subscriptions_base_url;
