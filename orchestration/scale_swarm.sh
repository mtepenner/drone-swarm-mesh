#!/usr/bin/env bash
set -euo pipefail

COUNT="${1:-25}"
if ! [[ "$COUNT" =~ ^[0-9]+$ ]]; then
  echo "usage: $0 <replica-count>" >&2
  exit 1
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
docker compose -f "$SCRIPT_DIR/docker-compose.yaml" up -d --build --scale drone-agent="$COUNT"
