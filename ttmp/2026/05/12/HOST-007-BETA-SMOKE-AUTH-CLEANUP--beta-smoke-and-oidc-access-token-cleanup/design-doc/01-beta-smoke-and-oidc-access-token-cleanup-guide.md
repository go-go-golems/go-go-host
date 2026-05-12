---
Title: Beta smoke and OIDC access-token cleanup guide
Ticket: HOST-007-BETA-SMOKE-AUTH-CLEANUP
Status: active
Topics:
    - go-go-host
    - hosting
    - security
    - deployments
    - platform-admin
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: ../../../../../../../../../../code/wesen/2026-03-27--hetzner-k3s/gitops/kustomize/go-go-host/certificate.yaml
      Note: wildcard TLS certificate for generated site hosts
    - Path: ../../../../../../../../../../code/wesen/2026-03-27--hetzner-k3s/gitops/kustomize/go-go-host/deployment.yaml
      Note: live K3s image pin for rollout
    - Path: ../../../../../../../../../../code/wesen/2026-03-27--hetzner-k3s/gitops/kustomize/go-go-host/ingress.yaml
      Note: wildcard host ingress for dashboard and generated sites
    - Path: examples/hello-beta
      Note: durable source for the live hello beta demo site
    - Path: internal/httpapi/auth.go
      Note: auth middleware selects dev auth versus OIDC
    - Path: internal/httpapi/oidc.go
      Note: OIDC bearer-token verification and access-token/client matching
    - Path: internal/httpapi/oidc_bootstrap_test.go
      Note: unit coverage for token client matching and platform-admin bootstrap
    - Path: scripts/beta-smoke.sh
      Note: repeatable public beta smoke script
    - Path: web/admin/src/auth/oidc.ts
      Note: PKCE login token storage and bearer-token helper
    - Path: web/admin/src/services/goGoHostApi.ts
      Note: frontend API client attaches bearer tokens
ExternalSources: []
Summary: Intern-facing implementation guide for making the go-go-host beta smoke repeatable and cleaning up OIDC API bearer-token semantics.
LastUpdated: 2026-05-12T17:35:00-04:00
WhatFor: Use this guide to understand what the beta deployment proved, why access-token cleanup is needed, and how the smoke tooling should work.
WhenToUse: When maintaining the go-go-host beta deployment, onboarding engineers to auth/deploy flows, or extending CLI/device-flow support.
---


# Beta smoke and OIDC access-token cleanup guide

## Executive summary

The go-go-host beta now has a real public deployment at:

```text
https://hosting.yolo.scapegoat.dev
```

and a real public hosted demo app at:

```text
https://hello.hosting.yolo.scapegoat.dev
```

The HOST-006 work proved that the platform can run on the K3s cluster with Argo CD, Keycloak OIDC, Vault-managed secrets, shared Postgres, wildcard DNS, wildcard TLS, Traefik Ingress, and Goja runtime host routing. The next problem is not "can it run?" The next problem is making this beta repeatable and correcting one auth boundary before we build CLI/device-flow work on top of it.

During the first live demo-site deployment, API calls using the Keycloak access token failed with:

```text
verify id token: oidc: expected audience "go-go-host-dashboard" got []
```

API calls using the ID token succeeded. That means the current backend behaves as if dashboard API bearer tokens are ID tokens. This is good enough for the browser proof, but it is the wrong default for future API clients and OAuth Device Flow. APIs should accept access tokens, and browser code should prefer sending access tokens.

This ticket does three things:

1. preserve the live demo app source in the repo,
2. add a repeatable public beta smoke script,
3. update OIDC bearer-token validation so the API accepts Keycloak access tokens issued to the dashboard client while preserving issuer/signature/expiry checks and platform-admin bootstrap rules.

## Problem statement

The live beta deployment has moved through several maturity stages:

```text
local dev-auth MVP
  -> local Keycloak browser OIDC
  -> K3s Argo CD deployment
  -> wildcard DNS/TLS/Ingress
  -> public demo site
```

That is enough to demonstrate the product, but it leaves three operational gaps.

First, the public beta smoke was manual. We created an org, created a site, uploaded a temporary bundle from `/tmp`, activated it, then used curl to verify public endpoints. Those commands proved the platform, but they are not a durable runbook.

Second, the demo app source was not initially stored in the repository. The live deployment stored an immutable bundle and unpacked files under the beta pod's data directory, but a new engineer should not have to retrieve source from a PVC to understand what the demo site does.

Third, the OIDC API token behavior was discovered to be semantically wrong. The dashboard token exchange stores both `id_token` and `access_token`, but the frontend bearer helper returned the ID token. When we manually tried to use the access token against `/api/v1/me`, the backend rejected it because `github.com/coreos/go-oidc/v3/oidc` verified the access token as if it were an ID token with audience equal to the dashboard client ID.

## Current-state architecture

### Public beta topology

The live beta topology is:

```text
Browser
  -> https://hosting.yolo.scapegoat.dev
    -> Traefik Ingress
      -> Service go-go-host
        -> go-go-hostd pod
          -> HTTP API and embedded dashboard
          -> runtime supervisor
          -> Postgres control-plane DB
          -> PVC for bundles, unpacked deployments, and per-site SQLite DBs
```

Generated sites use host-based routing:

```text
https://<site-slug>.hosting.yolo.scapegoat.dev
```

The wildcard DNS and wildcard TLS path is now in place:

```text
*.hosting.yolo.scapegoat.dev
  -> 91.98.46.169
  -> Traefik wildcard Ingress
  -> go-go-host service
  -> runtime supervisor chooses runtime by Host header
```

The live demo site currently uses:

```text
https://hello.hosting.yolo.scapegoat.dev
```

### Site host generation

When a user creates a site, `SiteService.CreateSite` computes `primary_host` from the site slug and configured base domain:

```go
host := slug
baseDomain := strings.Trim(s.baseDomain, ".")
if baseDomain != "" && baseDomain != "localhost" {
    host = slug + "." + baseDomain
} else if baseDomain == "localhost" {
    host = slug + ".localhost"
}
```

For beta:

```yaml
baseDomain: hosting.yolo.scapegoat.dev
```

so slug `hello` becomes:

```text
hello.hosting.yolo.scapegoat.dev
```

### Runtime activation and host routing

Deployment activation builds a runtime spec with the site's primary host plus verified custom domains:

```go
hosts := []string{site.PrimaryHost}
if domains, err := s.store.ListVerifiedSiteDomains(ctx, site.ID); err == nil {
    for _, domain := range domains {
        hosts = append(hosts, domain.Hostname)
    }
}
```

The runtime supervisor stores active runtimes in a `byHost` map. Public requests that are not handled by `/api`, `/app`, `/admin`, `/healthz`, or `/readyz` fall through to the runtime supervisor. The supervisor normalizes `r.Host` and looks up the active site runtime.

This means the generated host is not just a display string. It is the lookup key for serving traffic.

### Current OIDC behavior

The dashboard uses OAuth Authorization Code + PKCE. The helper stores the token response in local storage:

```ts
{
  idToken: tokens.id_token,
  accessToken: tokens.access_token,
  refreshToken: tokens.refresh_token,
  expiresAt: ...
}
```

Before this ticket, the bearer helper returned:

```ts
return getStoredTokens()?.idToken;
```

The API receives bearer tokens in `Authorization: Bearer ...` and calls `oidcAuthenticator.authenticate`. Before this ticket, the backend used:

```go
provider.Verifier(&oidc.Config{ClientID: clientID})
```

and emitted errors named `verify id token` / `decode id token claims`.

That works for ID tokens whose `aud` is the dashboard client. It rejected the observed Keycloak access token because that access token had:

```json
{
  "typ": "Bearer",
  "azp": "go-go-host-dashboard",
  "realm_access": { "roles": ["go-go-host-admin"] },
  "email": "wesen@ruinwesen.com"
}
```

but did not have `aud: "go-go-host-dashboard"`.

## Gap analysis

### Gap 1: API bearer token should be access-token friendly

A browser client can technically send an ID token to its backend, but it is not the right abstraction for future API clients. OAuth access tokens are designed for API authorization. Device Flow, CLI login, and scripts should all be able to call the API with an access token.

The beta API must therefore accept Keycloak access tokens that are:

- issued by the configured issuer,
- cryptographically valid,
- unexpired,
- associated with the configured client either by `aud` or by `azp`,
- carrying the subject/email/role/group claims needed for user upsert and platform-admin bootstrap.

### Gap 2: Smoke test should be repeatable

The public demo site now proves the hosting flow, but we need a command that can be run after deployments, restarts, cert renewal, DNS changes, or image bumps.

The first smoke script should be read-only and safe:

- hit health and readiness,
- inspect `/api/v1/config`,
- fetch the demo site root,
- fetch `/platform`,
- fetch `/db`,
- fetch static assets.

A later authenticated mode can create/redeploy the site, but that should wait until access-token semantics are fixed live.

### Gap 3: Demo app source should be durable

The live bundle was originally created in `/tmp`. The source belongs in the repo so a new engineer can understand and reproduce the site.

The durable source location is:

```text
examples/hello-beta/
```

This directory should be treated as a fixture and tutorial app, not just a one-off artifact.

## Proposed solution

## 1. Bearer-token validation model

Use the OIDC provider to validate signature, issuer, and expiry for any bearer token. Do not disable cryptographic validation. Then enforce that the token is for the configured dashboard/API client using either `aud` or `azp`.

The practical Keycloak behavior we observed is:

- ID token: `aud = go-go-host-dashboard`
- access token: `azp = go-go-host-dashboard`, may not include the dashboard client in `aud`

Therefore the backend should accept either:

```text
aud contains go-go-host-dashboard
```

or:

```text
azp == go-go-host-dashboard
```

The implementation uses `SkipClientIDCheck: true` only for the library's built-in audience check. It still verifies the provider signature, issuer, and expiry, then performs a local audience/authorized-party check.

Pseudocode:

```go
token := bearerToken(r)
verified, err := verifier.Verify(ctx, token) // issuer/signature/expiry
claims := decodeClaims(verified)
if !tokenMatchesClient(clientID, verified.Audience, claims) {
    reject
}
user := upsertUser(issuer, verified.Subject, claims.Email, displayName(claims))
if shouldBootstrapPlatformAdmin(cfg, verified.Subject, claims) {
    ensurePlatformAdmin(user)
}
```

## 2. Frontend bearer-token preference

The dashboard should send the access token when present:

```ts
export function bearerToken(): string | undefined {
  const tokens = getStoredTokens();
  return tokens?.accessToken || tokens?.idToken;
}
```

Fallback to ID token keeps older/partial token responses usable, but the normal Keycloak PKCE path should send the access token.

## 3. Read-only beta smoke script

Add:

```text
scripts/beta-smoke.sh
```

Default targets:

```text
GO_GO_HOST_BETA_API_URL=https://hosting.yolo.scapegoat.dev
GO_GO_HOST_BETA_SITE_HOST=hello.hosting.yolo.scapegoat.dev
```

Checks:

```text
/healthz -> HTTP 200
/readyz -> HTTP 200
/api/v1/config -> expected publicBaseUrl/baseDomain/OIDC config
https://hello.hosting.yolo.scapegoat.dev/ -> contains expected marker
/platform -> host equals hello.hosting.yolo.scapegoat.dev
/db -> overLimit false and quota stats present
/assets/style.css -> HTTP 200 text/css
```

This script is intentionally unauthenticated so it can be run by operators without browser token extraction.

## 4. Durable demo app

Add:

```text
examples/hello-beta/go-go-host.json
examples/hello-beta/scripts/app.js
examples/hello-beta/assets/style.css
examples/hello-beta/README.md
```

The app should exercise enough of the platform to catch meaningful regressions:

- UI rendering,
- request/platform context,
- per-site SQLite writes,
- DB guard stats,
- static assets.

## Implementation phases

### Phase 1: Document and preserve current state

- Create HOST-007.
- Write this guide.
- Add a diary step explaining the observed token failure and the live demo resource IDs.
- Add detailed tasks.

### Phase 2: Code cleanup

- Update backend OIDC bearer validation in `internal/httpapi/oidc.go`.
- Add unit coverage in `internal/httpapi/oidc_bootstrap_test.go`.
- Update frontend bearer helper in `web/admin/src/auth/oidc.ts`.
- Rebuild embedded dashboard assets.

### Phase 3: Smoke/demo fixture

- Add `examples/hello-beta`.
- Add `scripts/beta-smoke.sh`.
- Validate against the live beta deployment.

### Phase 4: Deploy the auth fix

- Build and push a new image from the current commit.
- Update K3s GitOps image pin.
- Let Argo roll out the image.
- Verify dashboard login and live API access-token calls.

### Phase 5: Optional authenticated smoke

After the live backend accepts access tokens, add an optional smoke mode that can:

- reuse or create org `beta-demo`,
- reuse or create site `hello`,
- package `examples/hello-beta`,
- upload and activate a deployment,
- verify the public site.

## Testing and validation strategy

Local/repo validation:

```bash
go test ./...
pnpm --dir web/admin build
BUILD_WEB_LOCAL=1 go run ./cmd/build-web
scripts/beta-smoke.sh
```

Live validation after image rollout:

```bash
curl -fsS https://hosting.yolo.scapegoat.dev/api/v1/config | jq .
curl -fsS https://hello.hosting.yolo.scapegoat.dev/platform | jq .
```

Browser validation:

1. Open `https://hosting.yolo.scapegoat.dev/admin`.
2. Login through Keycloak/GitHub.
3. Confirm admin dashboard loads.
4. Confirm org/site/deployment views still load after the frontend switches to access-token bearer auth.

Access-token validation:

1. Extract a fresh browser access token or use future CLI auth tooling.
2. Call `/api/v1/me` with `Authorization: Bearer <access-token>`.
3. Expect HTTP 200 and `platformAdmin: true` for `wesen@ruinwesen.com`.

## Risks and review points

### Risk: `azp` acceptance can be too broad if issuer/client boundaries are loose

We only accept tokens from the configured issuer and with `azp` matching the configured client ID. That is appropriate for this beta. If multiple API audiences are introduced later, add an explicit API audience and Keycloak audience mapper rather than relying on dashboard-client `azp`.

### Risk: email-based platform-admin bootstrap is convenient but broad

The live config grants platform admin by email and role. That is useful while GitHub IdP/account linking settles, but production may prefer role- or subject-only bootstrap.

### Risk: read-only smoke can pass while authenticated flows are broken

The first `scripts/beta-smoke.sh` is intentionally unauthenticated. It proves the live site and public control-plane health, not login or upload. Authenticated smoke should be added after access-token rollout.

### Risk: live image lag

The repository can contain fixes before the K3s image pin is updated. Always check the GitOps image tag before assuming live behavior matches branch-tip source.

## Key file reference

App auth and API:

- `internal/httpapi/oidc.go` — OIDC bearer verification, claim decoding, user upsert, platform-admin bootstrap.
- `internal/httpapi/auth.go` — auth middleware selection between dev auth and OIDC.
- `internal/httpapi/oidc_bootstrap_test.go` — platform-admin and token client-matching unit tests.
- `web/admin/src/auth/oidc.ts` — PKCE login, token storage, bearer token helper, logout.
- `web/admin/src/services/goGoHostApi.ts` — RTK Query API client and upload bearer headers.

Demo and smoke:

- `examples/hello-beta/` — durable source for the live demo app.
- `scripts/beta-smoke.sh` — public beta smoke script.

Deployment and docs:

- `Dockerfile` — image build toolchain.
- `.github/workflows/publish-image.yaml` — GHCR publish path.
- `/home/manuel/code/wesen/2026-03-27--hetzner-k3s/gitops/kustomize/go-go-host/deployment.yaml` — live image pin.
- `/home/manuel/code/wesen/2026-03-27--hetzner-k3s/gitops/kustomize/go-go-host/ingress.yaml` — wildcard host routing.
- `/home/manuel/code/wesen/2026-03-27--hetzner-k3s/gitops/kustomize/go-go-host/certificate.yaml` — wildcard TLS certificate.
- `/home/manuel/code/wesen/2026-03-27--hetzner-k3s/gitops/kustomize/platform-cert-issuer/clusterissuer-dns01-digitalocean.yaml` — DNS-01 issuer.

## Glossary

**ID token**: OIDC token intended for the browser/client to learn who authenticated. It often has `aud` equal to the client ID.

**Access token**: OAuth token intended for APIs. In Keycloak it may carry roles and `azp`, and may require an audience mapper if the API wants a custom `aud`.

**azp / authorized party**: Claim identifying the client to which the token was issued. For the beta dashboard this is `go-go-host-dashboard`.

**Primary host**: Generated hostname stored on the site row and used by the runtime supervisor for host-based routing.

**Wildcard TLS**: Certificate covering `*.hosting.yolo.scapegoat.dev`, issued through DNS-01.
