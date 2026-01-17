#!/usr/bin/env sh
set -eu

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
cd "$ROOT_DIR/webapp"

HOST="${ASD_WEBAPP_HOST:-127.0.0.1}"
PORT="${ASD_WEBAPP_PORT:-5173}"

if [ ! -d node_modules ]; then
  echo "[DEV] Installation des dépendances frontend (npm ci)..." >&2
  npm ci
fi

echo "[DEV] Frontend (Vite) sur http://$HOST:$PORT" >&2
exec npm run dev -- --host "$HOST" --port "$PORT"
