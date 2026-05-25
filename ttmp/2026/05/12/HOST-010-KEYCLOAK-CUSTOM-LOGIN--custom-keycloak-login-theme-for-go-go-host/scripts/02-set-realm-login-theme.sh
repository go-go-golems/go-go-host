#!/usr/bin/env bash
# 02-set-realm-login-theme.sh
# Set the go-go-host login theme on the go-go-host realm via the Keycloak Admin API.
#
# Usage:
#   bash scripts/02-set-realm-login-theme.sh

set -euo pipefail

KC_URL="http://127.0.0.1:18080"
KC_ADMIN="admin"
KC_ADMIN_PASS="admin"
REALM="go-go-host"
THEME="go-go-host"

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

echo "==> Setting loginTheme=$THEME on realm $REALM..."
curl -sf -X PUT "$KC_URL/admin/realms/$REALM" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"loginTheme\": \"$THEME\"}" > /dev/null

echo "==> Verifying..."
CURRENT_THEME=$(curl -sf "$KC_URL/admin/realms/$REALM" \
  -H "Authorization: Bearer $TOKEN" | jq -r '.loginTheme')

if [ "$CURRENT_THEME" = "$THEME" ]; then
  echo "==> SUCCESS: Realm $REALM loginTheme=$THEME"
else
  echo "==> WARNING: Expected loginTheme=$THEME, got loginTheme=$CURRENT_THEME"
fi
