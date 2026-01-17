#!/usr/bin/env sh
set -eu

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
PY="$ROOT_DIR/venv/bin/python"

HOST="${ASD_WEB_HOST:-127.0.0.1}"
PORT="${ASD_WEB_PORT:-8000}"

if [ ! -x "$PY" ]; then
  echo "[ERREUR] venv introuvable." >&2
  echo "        Crée-la via: python3 -m venv venv && ./venv/bin/pip install -r requirements.txt" >&2
  exit 1
fi

cd "$ROOT_DIR"
echo "[DEV] Backend (FastAPI) sur http://$HOST:$PORT" >&2
exec "$PY" -m uvicorn utils.ui.web.app:create_app --factory --host "$HOST" --port "$PORT" --reload
