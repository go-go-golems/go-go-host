#!/usr/bin/env bash
set -euo pipefail

API_URL="${GO_GO_HOST_BETA_API_URL:-https://hosting.yolo.scapegoat.dev}"
ORG_ID="${GO_GO_HOST_BETA_ORG_ID:-org_36cc42ac-d5d7-441a-809d-6fefb7e3c761}"
SITE_ID="${GO_GO_HOST_BETA_SITE_ID:-site_0fcba219-8bc9-412f-a0e0-41a4066c7a21}"
SITE_HOST="${GO_GO_HOST_BETA_SITE_HOST:-hello.hosting.yolo.scapegoat.dev}"
CHANNEL="${GO_GO_HOST_BETA_CHANNEL:-default}"
BUNDLE_PATH="${GO_GO_HOST_BETA_BUNDLE_PATH:-bundles/hello-beta-agent-smoke.tar.gz}"
AGENT_GRANT_PATH="${GO_GO_HOST_BETA_AGENT_GRANT_PATH:-bundles/**}"
BUNDLE_SOURCE_DIR="${GO_GO_HOST_BETA_BUNDLE_SOURCE_DIR:-examples/hello-beta}"
BEARER_TOKEN="${GO_GO_HOST_BETA_BEARER_TOKEN:-}"
KEEP_AGENT="${GO_GO_HOST_BETA_KEEP_AGENT:-0}"

need() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "missing required command: $1" >&2
    exit 2
  fi
}

section() { printf '\n== %s ==\n' "$*"; }

need curl
need jq
need tar

if [[ -z "$BEARER_TOKEN" ]]; then
  cat >&2 <<'MSG'
missing GO_GO_HOST_BETA_BEARER_TOKEN

Set it to a Keycloak/go-go-host access token for a user that owns the target org/site.
Example extraction during manual browser testing:
  JSON.parse(localStorage.getItem('go-go-host.oidc.tokens')).accessToken
MSG
  exit 2
fi

if [[ ! -d "$BUNDLE_SOURCE_DIR" ]]; then
  echo "missing bundle source directory: $BUNDLE_SOURCE_DIR" >&2
  exit 2
fi

TMPDIR=$(mktemp -d)
AGENT_CONFIG="$TMPDIR/agent.json"
BUNDLE="$TMPDIR/hello-beta.tar.gz"
AGENT_ID=""
cleanup() {
  if [[ -n "$AGENT_ID" && "$KEEP_AGENT" != "1" ]]; then
    curl -fsS -X POST "$API_URL/api/v1/orgs/$ORG_ID/agents/$AGENT_ID/revoke" \
      -H "Authorization: Bearer $BEARER_TOKEN" >/dev/null || true
  fi
  rm -rf "$TMPDIR"
}
trap cleanup EXIT

section "package demo bundle"
tar -C "$BUNDLE_SOURCE_DIR" -czf "$BUNDLE" .
tar -tzf "$BUNDLE" | sed -n '1,20p'

section "create scoped agent"
AGENT_JSON=$(go run ./cmd/go-go-host agents create \
  --api-url "$API_URL" \
  --bearer-token "$BEARER_TOKEN" \
  --org-id "$ORG_ID" \
  --name "beta-agent-smoke-$(date +%Y%m%d%H%M%S)" \
  --site-id "$SITE_ID" \
  --channel "$CHANNEL" \
  --bundle-path "$AGENT_GRANT_PATH" \
  --can-activate \
  --output json)
echo "$AGENT_JSON" | jq .
AGENT_ID=$(echo "$AGENT_JSON" | jq -r '.[0].id')
ENROLLMENT_TOKEN=$(echo "$AGENT_JSON" | jq -r '.[0].enrollment_token')
if [[ -z "$AGENT_ID" || "$AGENT_ID" == "null" || -z "$ENROLLMENT_TOKEN" || "$ENROLLMENT_TOKEN" == "null" ]]; then
  echo "agent create response did not contain id/enrollment_token" >&2
  exit 1
fi

section "keygen and enroll"
go run ./cmd/go-go-host-agent keygen --config "$AGENT_CONFIG" --api-url "$API_URL" --output json | jq .
go run ./cmd/go-go-host-agent enroll --config "$AGENT_CONFIG" --api-url "$API_URL" --token "$ENROLLMENT_TOKEN" --output json | jq .

section "deploy and activate via signed agent"
DEPLOY_JSON=$(go run ./cmd/go-go-host-agent deploy \
  --config "$AGENT_CONFIG" \
  --bundle "$BUNDLE" \
  --site-id "$SITE_ID" \
  --channel "$CHANNEL" \
  --bundle-path "$BUNDLE_PATH" \
  --activate \
  --output json)
echo "$DEPLOY_JSON" | jq .
DEPLOYMENT_ID=$(echo "$DEPLOY_JSON" | jq -r '.[0].deployment_id')
VALID=$(echo "$DEPLOY_JSON" | jq -r '.[0].valid')
ACTIVATED=$(echo "$DEPLOY_JSON" | jq -r '.[0].activated')
if [[ "$VALID" != "true" || "$ACTIVATED" != "true" ]]; then
  echo "agent deployment did not validate and activate" >&2
  exit 1
fi

section "verify public hosted site"
PLATFORM=$(curl -fsS "https://$SITE_HOST/platform")
echo "$PLATFORM" | jq .
LIVE_DEPLOYMENT_ID=$(echo "$PLATFORM" | jq -r '.deploymentId')
if [[ "$LIVE_DEPLOYMENT_ID" != "$DEPLOYMENT_ID" ]]; then
  echo "expected live deployment $DEPLOYMENT_ID, got $LIVE_DEPLOYMENT_ID" >&2
  exit 1
fi
curl -fsS "https://$SITE_HOST/" | grep -q "Hello from go-go-host beta"
curl -fsSI "https://$SITE_HOST/assets/style.css" | sed -n '1,10p'

section "revoke smoke agent"
if [[ "$KEEP_AGENT" == "1" ]]; then
  echo "keeping smoke agent $AGENT_ID because GO_GO_HOST_BETA_KEEP_AGENT=1"
else
  curl -fsS -X POST "$API_URL/api/v1/orgs/$ORG_ID/agents/$AGENT_ID/revoke" \
    -H "Authorization: Bearer $BEARER_TOKEN" | jq .
  AGENT_ID=""
fi

section "ok"
echo "agent beta smoke passed; deployment $DEPLOYMENT_ID is live at https://$SITE_HOST"
