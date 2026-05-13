#!/usr/bin/env bash
# 03-add-github-idp.sh
# Add GitHub as an OIDC identity provider to the go-go-host realm.
#
# Prerequisites:
#   - Keycloak running at http://127.0.0.1:18080
#   - GitHub OAuth App created at https://github.com/settings/developers
#   - Environment variables GITHUB_CLIENT_ID and GITHUB_CLIENT_SECRET set
#
# Usage:
#   GITHUB_CLIENT_ID=xxx GITHUB_CLIENT_SECRET=yyy bash scripts/03-add-github-idp.sh
#
# If no credentials are provided, creates a placeholder IdP for styling testing.

set -euo pipefail

KC_URL="http://127.0.0.1:18080"
KC_ADMIN="admin"
KC_ADMIN_PASS="admin"
REALM="go-go-host"

CLIENT_ID="${GITHUB_CLIENT_ID:-placeholder-github-client-id}"
CLIENT_SECRET="${GITHUB_CLIENT_SECRET:-placeholder-github-client-secret}"

echo "==> Getting admin token..."
TOKEN=$(curl -sf "$KC_URL/realms/master/protocol/openid-connect/token" \
  -d "client_id=admin-cli" \
  -d "username=$KC_ADMIN" \
  -d "password=$KC_ADMIN_PASS" \
  -d "grant_type=password" | jq -r '.access_token')

if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ]; then
  echo "ERROR: Failed to get admin token. Is Keycloak running at $KC_URL?"
  exit 1
fi

echo "==> Adding GitHub identity provider to realm $REALM..."
curl -sf -X POST "$KC_URL/admin/realms/$REALM/identity-provider/instances" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"alias\": \"github\",
    \"providerId\": \"github\",
    \"enabled\": true,
    \"firstBrokerLoginFlowAlias\": \"first broker login\",
    \"config\": {
      \"clientId\": \"$CLIENT_ID\",
      \"clientSecret\": \"$CLIENT_SECRET\",
      \"useJwksUrl\": \"true\",
      \"syncMode\": \"IMPORT\",
      \"trustEmail\": \"true\",
      \"storeToken\": \"true\",
      \"addReadTokenRoleOnCreate\": \"false\",
      \"acceptsPromptNoneForwardFromClient\": \"false\"
    }
  }" > /dev/null

echo "==> Verifying..."
IDPS=$(curl -sf "$KC_URL/admin/realms/$REALM/identity-provider/instances" \
  -H "Authorization: Bearer $TOKEN" | jq -r '.[].alias')

if echo "$IDPS" | grep -q "github"; then
  echo "==> SUCCESS: GitHub IdP added to realm $REALM"
  echo "    Callback URL: $KC_URL/realms/$REALM/broker/github/endpoint"
  echo ""
  echo "    Set your GitHub OAuth App callback to:"
  echo "    $KC_URL/realms/$REALM/broker/github/endpoint"
else
  echo "==> WARNING: GitHub IdP not found. Available: $IDPS"
fi
