#!/usr/bin/env sh
set -eu

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"

ADDR="${ASD_ADDR:-127.0.0.1:8080}"
DB="${ASD_DB_PATH:-asd.db}"

cd "$ROOT_DIR"
echo "[DEV] Backend (Go) sur http://$ADDR" >&2
exec go run ./cmd/asd-server -addr "$ADDR" -db "$DB"
