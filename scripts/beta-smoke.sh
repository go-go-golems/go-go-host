#!/usr/bin/env bash
set -euo pipefail

API_URL="${GO_GO_HOST_BETA_API_URL:-https://hosting.yolo.scapegoat.dev}"
SITE_HOST="${GO_GO_HOST_BETA_SITE_HOST:-hello.hosting.yolo.scapegoat.dev}"

need() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "missing required command: $1" >&2
    exit 2
  fi
}

need curl
need jq

section() { printf '\n== %s ==\n' "$*"; }

section "control-plane health"
curl -fsSI "$API_URL/healthz" | sed -n '1,8p'
curl -fsSI "$API_URL/readyz" | sed -n '1,8p'

section "browser config"
CONFIG=$(curl -fsS "$API_URL/api/v1/config")
echo "$CONFIG" | jq '{publicBaseUrl, baseDomain, devAuth, oidc}'
BASE_DOMAIN=$(echo "$CONFIG" | jq -r '.baseDomain')
PUBLIC_BASE_URL=$(echo "$CONFIG" | jq -r '.publicBaseUrl')
if [[ "$BASE_DOMAIN" != "hosting.yolo.scapegoat.dev" ]]; then
  echo "unexpected baseDomain: $BASE_DOMAIN" >&2
  exit 1
fi
if [[ "$PUBLIC_BASE_URL" != "$API_URL" ]]; then
  echo "unexpected publicBaseUrl: $PUBLIC_BASE_URL" >&2
  exit 1
fi

section "demo site root"
curl -fsSI "https://$SITE_HOST/" | sed -n '1,10p'
HOME=$(curl -fsS "https://$SITE_HOST/")
if ! grep -q "Hello from go-go-host beta" <<<"$HOME"; then
  echo "demo homepage did not contain expected marker" >&2
  exit 1
fi

section "demo platform context"
PLATFORM=$(curl -fsS "https://$SITE_HOST/platform")
echo "$PLATFORM" | jq .
HOST=$(echo "$PLATFORM" | jq -r '.host')
if [[ "$HOST" != "$SITE_HOST" ]]; then
  echo "unexpected platform host: $HOST" >&2
  exit 1
fi

section "demo db stats"
curl -fsS "https://$SITE_HOST/db" | jq '{overLimit, stats: {totalBytes: .stats.totalBytes, softMaxBytes: .stats.softMaxBytes, hardMaxBytes: .stats.hardMaxBytes}}'

section "demo asset"
curl -fsSI "https://$SITE_HOST/assets/style.css" | sed -n '1,10p'

section "ok"
echo "beta smoke passed for $API_URL and https://$SITE_HOST"
