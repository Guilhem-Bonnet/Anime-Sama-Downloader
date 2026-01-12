#!/usr/bin/env sh
set -eu

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"

if [ ! -x "$ROOT_DIR/venv/bin/python" ]; then
  echo "[ERROR] venv introuvable. Crée-la via: python3 -m venv venv && ./venv/bin/pip install -r requirements.txt" >&2
  exit 1
fi

cd "$ROOT_DIR"
exec "$ROOT_DIR/venv/bin/python" -m uvicorn utils.ui.web.app:create_app --factory --host 127.0.0.1 --port 8000 --reload
