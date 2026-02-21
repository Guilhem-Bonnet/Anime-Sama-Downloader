-- +migrate Up

ALTER TABLE jobs ADD COLUMN file_list_json BLOB;

-- +migrate Down

-- SQLite ne supporte pas DROP COLUMN avant 3.35.0
-- La colonne restera mais sera inutilisée si on rollback
