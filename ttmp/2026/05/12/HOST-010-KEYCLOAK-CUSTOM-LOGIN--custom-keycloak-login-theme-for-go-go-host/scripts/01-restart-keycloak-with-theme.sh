#!/usr/bin/env bash
# 01-restart-keycloak-with-theme.sh
# Restart the dev Keycloak container with the go-go-host login theme mounted.
#
# Prerequisites:
#   - Docker Compose stack is running (devctl up or docker compose up)
#   - The theme files exist at deployments/dev/keycloak/themes/go-go-host/
#
# Usage:
#   bash scripts/01-restart-keycloak-with-theme.sh

set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../../../.." && pwd)"
COMPOSE_DIR="$REPO_ROOT/deployments/dev"

echo "==> Restarting Keycloak with go-go-host login theme..."
cd "$COMPOSE_DIR"
docker compose restart keycloak
echo "==> Waiting for Keycloak to become ready..."
for i in $(seq 1 30); do
  if curl -sf http://127.0.0.1:18080/health/ready > /dev/null 2>&1; then
    echo "==> Keycloak is ready at http://127.0.0.1:18080"
    exit 0
  fi
  sleep 2
done
echo "==> WARNING: Keycloak did not report ready within 60s, but may still be starting."
echo "    Try: curl -sf http://127.0.0.1:18080/health/ready"
