#!/usr/bin/env sh
set -eu

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
PY="$ROOT_DIR/venv/bin/python"

if [ ! -x "$PY" ]; then
  echo "[ERROR] venv introuvable. Crée-la via: python3 -m venv venv && ./venv/bin/pip install -r requirements.txt" >&2
  exit 1
fi

mode="${1:-py}"

if [ "$mode" = "py" ]; then
  echo "[SMOKE] py_compile"
  exec "$PY" -m py_compile \
    main.py \
    utils/config.py \
    utils/output_paths.py \
    utils/search.py \
    utils/tui.py \
    utils/ui/web/app.py
fi

if [ "$mode" = "docker" ]; then
  echo "[SMOKE] docker compose up (dev) + ping /api/health"
  cd "$ROOT_DIR"
  docker compose up --build -d
  cleanup() {
    docker compose down >/dev/null 2>&1 || true
  }
  trap cleanup EXIT INT TERM

  # wait a bit
  i=0
  while [ $i -lt 30 ]; do
    if curl -fsS http://127.0.0.1:8000/api/health >/dev/null 2>&1; then
      echo "[OK] /api/health"
      exit 0
    fi
    i=$((i+1))
    sleep 1
  done

  echo "[ERROR] /api/health ne répond pas" >&2
  docker compose ps || true
  exit 1
fi

echo "Usage: $0 [py|docker]" >&2
exit 2
