---
Title: "Login and CLI Configuration"
Slug: "login-and-config"
Short: "Store API URL and OAuth/device-flow auth defaults for go-go-host CLI commands."
Topics:
  - go-go-host
  - login
  - configuration
  - oauth
  - oidc
Commands:
  - login
  - logout
  - me
Flags:
  - api-url
  - dev-user
  - bearer-token
  - client-id
  - scopes
IsTopLevel: true
IsTemplate: false
ShowPerDefault: true
SectionType: Tutorial
---

This tutorial covers CLI authentication and local configuration. The production login path uses OAuth 2.0 Device Authorization Grant through Keycloak. The CLI prints a browser URL and a short user code, waits while you approve the request in Keycloak, then stores the returned tokens in the local CLI config file.

The config file is stored at the OS user config path by default:

```text
~/.config/go-go-host/config.yaml
```

Set `GO_GO_HOST_CLI_CONFIG` to use a different file for tests or isolated environments.

## Production login with device flow

Run `login` with the production API URL and no manual auth flags:

```bash
go-go-host login --api-url https://hosting.yolo.scapegoat.dev
```

The CLI reads `/api/v1/config`, discovers Keycloak's OIDC metadata, starts Device Authorization Grant with the configured CLI client, and prints instructions like:

```text
Open this URL in your browser:
  https://auth.yolo.scapegoat.dev/realms/go-go-host/device?user_code=WDJB-MJHT

Enter this code if prompted:
  WDJB-MJHT

Waiting for browser authorization...
```

Open the URL, log in through Keycloak, and confirm the code. After approval, the CLI stores the access token, refresh token, issuer, client ID, scopes, and expiry in the config file. Future commands use the stored token automatically:

```bash
go-go-host me --output table
go-go-host org list
go-go-host site list --org-id org_...
```

If the access token is close to expiry and a refresh token is available, the CLI refreshes it before making API requests.

## Override client ID or scopes

The server config should normally publish the correct device-flow client ID. Use `--client-id` only for testing a different Keycloak client:

```bash
go-go-host login \
  --api-url https://hosting.yolo.scapegoat.dev \
  --client-id go-go-host-cli
```

Override scopes with a comma- or space-separated list:

```bash
go-go-host login \
  --api-url https://hosting.yolo.scapegoat.dev \
  --scopes "openid profile email"
```

## Logout

Use `logout` to clear local tokens:

```bash
go-go-host logout
```

If the current session has a refresh token and Keycloak publishes a revocation endpoint, logout best-effort revokes the refresh token before clearing local state. Local state is cleared even if revocation fails.

## Configure dev auth

In local dev mode, the daemon accepts `X-Go-Go-Host-User`. Store that value with `login`:

```bash
go-go-host login --api-url http://127.0.0.1:8080 --dev-user alice
```

Then verify the session context:

```bash
go-go-host me --output table
```

## Configure manual bearer auth

For one-off non-dev auth smoke tests, store a bearer token directly:

```bash
go-go-host login --api-url https://host.example.com --bearer-token "$TOKEN"
```

Commands also accept `--bearer-token` directly, which overrides the stored value for that invocation. This mode is retained for debugging, but production human login should use device flow.

## Troubleshooting

| Problem | Cause | Solution |
| --- | --- | --- |
| `login` says the server did not publish OIDC configuration | The daemon is running with `devAuth: true` or OIDC config is missing | Use `--dev-user` locally, or configure `oidcIssuer` / `oidcClientId` / `oidcDeviceClientId` for non-dev |
| Device authorization returns `unauthorized_client` | The Keycloak client has Device Authorization Grant disabled | Enable `oauth2_device_authorization_grant_enabled` for `go-go-host-cli` in Terraform/local realm config |
| `me` returns unauthorized after login | The API does not accept the CLI client ID in `aud` or `azp` | Add `go-go-host-cli` to `oidcAcceptedClientIds` and redeploy the API config |
| Commands still use the wrong API URL | A local config file is taking precedence | Pass `--api-url` explicitly or set `GO_GO_HOST_CLI_CONFIG` to a test config path |
| Tests affect your normal config | The default config lives in the OS user config directory | Set `GO_GO_HOST_CLI_CONFIG=$(mktemp)` for smoke tests |

## See Also

- `getting-started`
- `deploy-workflow`
