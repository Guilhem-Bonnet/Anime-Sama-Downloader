#!/usr/bin/env sh
set -eu

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
cd "$ROOT_DIR/webapp"

if [ ! -d node_modules ]; then
  npm ci
fi

exec npm run dev
