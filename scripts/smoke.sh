#!/usr/bin/env sh
set -eu

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"

mode="${1:-go}"

if [ "$mode" = "go" ]; then
  echo "[SMOKE] go test ./..." >&2
  cd "$ROOT_DIR"
  exec go test ./...
fi

if [ "$mode" = "docker" ]; then
  echo "[SMOKE] Docker dev + ping /api/v1/health" >&2
  cd "$ROOT_DIR"

  if ! command -v docker >/dev/null 2>&1; then
    echo "[ERREUR] docker est requis pour ce mode." >&2
    exit 1
  fi
  if ! docker compose version >/dev/null 2>&1; then
    echo "[ERREUR] docker compose est requis (plugin Compose v2)." >&2
    exit 1
  fi

  docker compose up --build -d
  cleanup() {
    docker compose down >/dev/null 2>&1 || true
  }
  trap cleanup EXIT INT TERM

  # wait a bit
  i=0
  while [ $i -lt 30 ]; do
    if curl -fsS http://127.0.0.1:8080/api/v1/health >/dev/null 2>&1; then
      echo "[OK] /api/v1/health" >&2
      exit 0
    fi
    i=$((i+1))
    sleep 1
  done

  echo "[ERREUR] /api/v1/health ne répond pas" >&2
  docker compose ps || true
  exit 1
fi

echo "Usage: $0 [go|docker]" >&2
echo "  go     : lance go test ./..." >&2
echo "  docker : lance docker compose (dev) et ping /api/v1/health" >&2
exit 2
