---
Title: Investigation diary
Ticket: HOST-007-BETA-SMOKE-AUTH-CLEANUP
Status: active
Topics:
    - go-go-host
    - hosting
    - security
    - deployments
    - platform-admin
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: examples/hello-beta
      Note: demo app fixture added during HOST-007
    - Path: internal/httpapi/oidc.go
      Note: implemented OIDC bearer access-token matching
    - Path: scripts/beta-smoke.sh
      Note: live beta smoke script added during HOST-007
    - Path: web/admin/src/auth/oidc.ts
      Note: frontend now prefers access tokens
ExternalSources: []
Summary: Chronological diary for beta smoke and OIDC access-token cleanup work.
LastUpdated: 2026-05-12T17:35:00-04:00
WhatFor: Use this to resume or audit HOST-007 implementation.
WhenToUse: When changing auth token semantics, beta smoke scripts, or demo-site fixtures.
---


# Diary

## Goal

Capture the work needed to turn the successful go-go-host beta demo into repeatable smoke tooling and to fix the OIDC bearer-token semantics discovered while scripting the demo deployment.

## Step 1: Ticket setup and evidence-backed guide

I created HOST-007 after the live demo site proved that the platform can serve a real Goja app at `https://hello.hosting.yolo.scapegoat.dev`. The guide explains the current beta architecture, why the ID-token/access-token distinction matters, and how the next implementation work should be structured for a new intern.

The main technical finding is that the public beta path works, but the API authentication boundary needs cleanup before future CLI/device-flow work. Browser login stores both ID and access tokens; the previous frontend sent the ID token; the backend verified all bearer tokens as ID tokens. That was enough for the dashboard, but not correct for API clients.

### Prompt Context

**User prompt (verbatim):** "Createa new ticket.

Create a guide document that explains everything going on for a new intern and upload to remarkable, so that we can keep track of what work was needed and why and how it was going to be solved, in addition to the diary. Then add detailed tasks to the ticket, and tackle it."

**Assistant interpretation:** Create a new docmgr ticket for beta smoke/auth cleanup, write a comprehensive intern-facing guide and diary, upload to reMarkable, define tasks, and begin implementing the highest-priority items.

**Inferred user intent:** Preserve the hard-won deployment knowledge, convert the next phase into a structured ticket, and make immediate progress instead of only planning.

### What I did

- Created ticket `HOST-007-BETA-SMOKE-AUTH-CLEANUP`.
- Created the design guide `design-doc/01-beta-smoke-and-oidc-access-token-cleanup-guide.md`.
- Created this investigation diary.
- Wrote detailed tasks in `tasks.md`.
- Reviewed current auth and deployment files:
  - `internal/httpapi/auth.go`
  - `internal/httpapi/oidc.go`
  - `web/admin/src/auth/oidc.ts`
  - `web/admin/src/services/goGoHostApi.ts`
  - `configs/production.example.yaml`
  - `deployments/dev/keycloak/realm-go-go-host.json`
- Used the live demo deployment evidence from HOST-006 as starting context.

### Why

- The live beta deployment now has enough moving parts that an intern needs a map before changing it.
- OIDC access-token cleanup is a prerequisite for CLI/device-flow work.
- A repeatable smoke script prevents future infrastructure or image changes from silently breaking the public demo path.

### What worked

- The repository already had clear separation between OIDC auth, frontend PKCE storage, RTK Query bearer injection, runtime host routing, and deployment activation.
- The live demo resource IDs and public URLs provided concrete acceptance criteria for the smoke script.

### What didn't work

- The first scripted API attempt with the Keycloak access token failed during HOST-006 with:

```text
{"error":"verify id token: oidc: expected audience \"go-go-host-dashboard\" got []"}
```

That failure became the central technical driver for HOST-007.

### What I learned

- Keycloak access tokens for the public dashboard client can identify the client through `azp` even when `aud` does not include the dashboard client ID.
- The backend naming around `idToken` made the previous behavior explicit: it really was validating bearer tokens as ID tokens.

### What was tricky to build

The tricky part is preserving security while accepting access tokens. Disabling the library's client ID check alone would be too broad. The safer beta approach is to keep issuer/signature/expiry validation through the OIDC provider, then add a local check requiring either `aud` or `azp` to match the configured client.

### What warrants a second pair of eyes

- Whether accepting `azp == go-go-host-dashboard` is sufficient for beta, or whether Keycloak should instead add a dedicated API audience mapper and the backend should require that audience.
- Whether platform-admin bootstrap should continue accepting email matches after GitHub IdP is stable.

### What should be done in the future

- Add authenticated smoke mode after the live image contains the access-token fix.
- Add CLI OAuth Device Flow design/implementation as a separate ticket.

### Code review instructions

- Start with `internal/httpapi/oidc.go` and `web/admin/src/auth/oidc.ts`.
- Confirm the verifier still validates issuer/signature/expiry.
- Confirm `scripts/beta-smoke.sh` is read-only and safe to run repeatedly.
- Validate with `go test ./...`, `pnpm --dir web/admin build`, and `scripts/beta-smoke.sh`.

### Technical details

Representative commands:

```bash
docmgr ticket create-ticket --ticket HOST-007-BETA-SMOKE-AUTH-CLEANUP --title "Beta smoke and OIDC access-token cleanup" --topics go-go-host,hosting,security,deployments,platform-admin
rg -n "idToken|accessToken|Verifier|SkipClientIDCheck|bearerToken" internal web/admin/src
scripts/beta-smoke.sh
```

## Step 2: Implement access-token preference, backend token matching, demo fixture, and smoke script

I implemented the first code changes for HOST-007. The frontend now prefers the access token for API calls, while the backend can validate a Keycloak bearer token that identifies the dashboard client through either `aud` or `azp`. I also preserved the live demo app source under `examples/hello-beta` and added a read-only beta smoke script.

This step does not yet mean the live cluster has the new auth behavior. The K3s deployment still needs a new image build, image push, and GitOps image pin update before the live API accepts access tokens in the way the source branch now does.

### Prompt Context

**User prompt (verbatim):** (same as Step 1)

**Assistant interpretation:** Begin implementation after creating the ticket and guide.

**Inferred user intent:** Make tangible progress on the highest-priority cleanup while keeping documentation synchronized.

**Commit (code):** pending at time of diary entry.

### What I did

- Updated `web/admin/src/auth/oidc.ts` so `bearerToken()` returns `accessToken || idToken`.
- Updated `internal/httpapi/oidc.go`:
  - renamed error context from ID-token-specific language to OIDC bearer-token language,
  - added `aud` and `azp` claims,
  - set the OIDC verifier to `SkipClientIDCheck: true` while preserving provider validation,
  - added `tokenMatchesClient` to require `aud` or `azp` to match the configured client ID.
- Added unit coverage for `tokenMatchesClient` in `internal/httpapi/oidc_bootstrap_test.go`.
- Rebuilt embedded dashboard assets with `BUILD_WEB_LOCAL=1 go run ./cmd/build-web`.
- Added durable demo app source under `examples/hello-beta`.
- Added `scripts/beta-smoke.sh` and ran it against the live beta deployment.

### Why

- Future CLI/device-flow work needs access-token semantics.
- The dashboard should not keep sending ID tokens when an access token is present.
- The live demo app should be reproducible from source, not only from a deployed bundle.
- The smoke script gives operators a fast no-auth health check for both dashboard/API and the generated-site path.

### What worked

Validation passed:

```bash
go test ./internal/httpapi ./internal/config ./internal/store ./internal/control ./internal/runtime
pnpm --dir web/admin build
BUILD_WEB_LOCAL=1 go run ./cmd/build-web
scripts/beta-smoke.sh
go test ./...
```

`scripts/beta-smoke.sh` verified:

```text
https://hosting.yolo.scapegoat.dev/healthz
https://hosting.yolo.scapegoat.dev/readyz
https://hosting.yolo.scapegoat.dev/api/v1/config
https://hello.hosting.yolo.scapegoat.dev/
https://hello.hosting.yolo.scapegoat.dev/platform
https://hello.hosting.yolo.scapegoat.dev/db
https://hello.hosting.yolo.scapegoat.dev/assets/style.css
```

### What didn't work

No new command failure occurred in this step. The known live failure remains the old image rejecting access tokens; this will persist until a new image is deployed.

### What I learned

- A read-only smoke script is immediately useful even before authenticated smoke exists.
- The demo app is a good fixture because it touches HTML rendering, platform context, SQLite state, DB guard stats, and static assets.

### What was tricky to build

The auth cleanup has an important subtlety: `SkipClientIDCheck` sounds dangerous if read in isolation. In this implementation it is paired with explicit local matching against `aud` and `azp`. This preserves a client binding while allowing Keycloak's access-token shape.

### What warrants a second pair of eyes

- Review the claim matching policy in `tokenMatchesClient`.
- Decide if we should add a Keycloak audience mapper and require `aud` instead of allowing `azp` long term.

### What should be done in the future

- Build/push a new image and update K3s GitOps.
- Verify live access-token API calls after rollout.
- Add authenticated smoke mode.

### Code review instructions

Review:

```text
internal/httpapi/oidc.go
internal/httpapi/oidc_bootstrap_test.go
web/admin/src/auth/oidc.ts
examples/hello-beta/
scripts/beta-smoke.sh
```

Validate:

```bash
go test ./...
pnpm --dir web/admin build
BUILD_WEB_LOCAL=1 go run ./cmd/build-web
scripts/beta-smoke.sh
```

### Technical details

The source for the public demo site is now located at:

```text
examples/hello-beta/
```

The live deployment serving `https://hello.hosting.yolo.scapegoat.dev` is deployment:

```text
dep_181c0489-b037-4732-b7b3-3cc99bf4ea52
```

The code currently serving live traffic came from an equivalent bundle first built under `/tmp`; `examples/hello-beta` preserves that app source for future rebuilds and redeploys.

## Step 3: Deploy access-token fix and discover restart/runtime restoration gap

I built and pushed the first HOST-007 image, updated the K3s GitOps image pin, and let Argo roll out the new pod. The access-token fix worked live: `/api/v1/me` accepted the browser access token and returned the expected user, memberships, and platform-admin status. The rollout also exposed a separate runtime lifecycle bug: active deployments are stored in the database, but the in-memory runtime supervisor did not automatically restore them after a daemon restart.

This discovery is exactly why the read-only smoke script is useful. After the rollout, the control-plane health checks passed, but `https://hello.hosting.yolo.scapegoat.dev/` returned 404 because the new process had not reactivated the site's active deployment into the supervisor host map.

### Prompt Context

**User prompt (verbatim):** "where is the code for hello.hosting.yolo.scapegoat.dev ? btw ?

then continue"

**Assistant interpretation:** Answer that the demo source is now in `examples/hello-beta`, then continue implementing and deploying the HOST-007 work.

**Inferred user intent:** Keep momentum, make the beta cleanup live, and ensure the demo site remains understandable and reproducible.

**Commit (code):** pending for startup-restore fix at time of diary entry.

### What I did

- Committed and pushed the first HOST-007 source changes as `23b66ec Add beta smoke and OIDC access token support`.
- Built and pushed:

```text
ghcr.io/go-go-golems/go-go-host:sha-23b66ec
```

- Updated K3s GitOps deployment image from `sha-4187ea3` to `sha-23b66ec` and pushed K3s commit:

```text
5779ad3 Bump go-go-host beta image
```

- Refreshed Argo and waited for rollout.
- Verified live access-token API auth:

```bash
curl -fsS https://hosting.yolo.scapegoat.dev/api/v1/me \
  -H "Authorization: Bearer <browser access token>"
```

returned:

```json
{
  "email": "wesen@ruinwesen.com",
  "platformAdmin": true,
  "memberships": ["beta-demo", "wesen"]
}
```

- Ran `scripts/beta-smoke.sh` after rollout and found the demo site root returned 404.
- Implemented startup active-runtime restoration:
  - `DeploymentService.RestoreActiveRuntimes(ctx)` lists active deployments and activates their runtime specs into the supervisor.
  - `cmd/go-go-hostd/main.go` calls restoration after stale runtime-status reconciliation and before starting the HTTP server.

### Why

- The access-token fix needed to be proven live, not just locally.
- The public demo site should survive normal pod restarts and image rollouts. Requiring manual reactivation after every daemon restart is not acceptable for beta.

### What worked

- The access-token fix worked in the live cluster after deploying `sha-23b66ec`.
- Argo rolled out the image successfully.
- The smoke script immediately caught the runtime restoration issue.
- The codebase already had enough activation logic to reuse for restoration.

### What didn't work

`scripts/beta-smoke.sh` failed after image rollout with:

```text
== demo site root ==
curl: (22) The requested URL returned error: 404
HTTP/2 404
```

The root cause was not TLS, DNS, or Ingress. `/healthz`, `/readyz`, and `/api/v1/config` passed. The failure was inside go-go-host: the supervisor's host map was empty after process restart.

### What I learned

- Runtime status reconciliation on daemon startup marked stale runtime statuses stopped, but there was no corresponding step to reload database-active deployments into memory.
- Activation is both a database operation and an in-memory supervisor operation. The database remembered `active_deployment_id`; the new process needed to reconstruct the supervisor state.

### What was tricky to build

The restoration path needed to avoid inventing a second activation implementation. I added a method that lists active deployments and calls the existing internal `activate` helper with a system actor. This reuses manifest parsing, host list construction, quota loading, runtime construction, and supervisor activation.

### What warrants a second pair of eyes

- `RestoreActiveRuntimes` currently uses `ListAdminDeployments` with status `active` and limit `1000`. That is fine for beta but should become a purpose-built store query before broad production use.
- Startup restoration currently returns the first error and fails daemon startup. Review whether beta should instead log per-site failures and continue serving other active sites.

### What should be done in the future

- Add unit/integration coverage for restart restoration.
- Consider making runtime status reconciliation and active restoration a single explicit startup phase with structured logs.

### Code review instructions

Review:

```text
internal/control/deployments.go
cmd/go-go-hostd/main.go
```

Validate:

```bash
go test ./internal/control ./cmd/go-go-hostd ./internal/httpapi
go test ./...
```

After deploying the next image, validate:

```bash
scripts/beta-smoke.sh
```

### Technical details

The durable source for the demo site is:

```text
examples/hello-beta/
```

The live site URL remains:

```text
https://hello.hosting.yolo.scapegoat.dev/
```

## Step 4: Deploy startup restoration and confirm beta smoke passes after rollout

I built and deployed the startup-restoration fix so active sites survive a normal image rollout. This completed the immediate HOST-007 implementation loop: access tokens work live, the demo app source is durable, the smoke script exists, and the live demo site remains reachable after pod restart.

### Prompt Context

**User prompt (verbatim):** (same as Step 3)

**Assistant interpretation:** Continue from the runtime restoration discovery and make the live beta healthy again.

**Inferred user intent:** Finish the implementation, not just document or partially fix it.

**Commit (code):** `f137ff9` — "Restore active runtimes on daemon startup"

### What I did

- Built and pushed:

```text
ghcr.io/go-go-golems/go-go-host:sha-f137ff9
```

- Updated K3s GitOps deployment image and pushed K3s commit:

```text
13ac467 Deploy go-go-host startup runtime restore
```

- Refreshed Argo and waited for rollout.
- Confirmed Argo and Deployment state:

```text
go-go-host Synced Healthy
deployment image ghcr.io/go-go-golems/go-go-host:sha-f137ff9
```

- Re-ran `scripts/beta-smoke.sh` after rollout.
- Re-verified live access-token API calls against `/api/v1/me`.

### Why

The first HOST-007 image proved access-token support but broke the public demo site until manual reactivation. A beta hosting service must restore active sites on daemon restart, because pod restarts, image rollouts, node drains, and crashes are normal operational events.

### What worked

After deploying `sha-f137ff9`, `scripts/beta-smoke.sh` passed:

```text
beta smoke passed for https://hosting.yolo.scapegoat.dev and https://hello.hosting.yolo.scapegoat.dev
```

The demo site returned HTTP 200 after the pod restart without manual reactivation.

The browser access token also worked live:

```json
{
  "email": "wesen@ruinwesen.com",
  "platformAdmin": true,
  "memberships": ["beta-demo", "wesen"]
}
```

### What didn't work

No new failure occurred after deploying `sha-f137ff9`. The earlier `sha-23b66ec` rollout failure is now fixed by startup restoration.

### What I learned

- The read-only smoke script is valuable as a post-rollout gate, not just a one-off status check.
- Access-token semantics and runtime restoration are tightly connected operationally: both need to work before a scripted authenticated deploy smoke is trustworthy.

### What was tricky to build

The main sharp edge was Argo status timing. Immediately after pushing image bumps, the Deployment image and Argo summary can briefly disagree while refresh/sync is in progress. I waited for both Argo `Synced Healthy` and the Deployment container image to show the expected SHA before treating rollout as complete.

### What warrants a second pair of eyes

- Startup restoration failure policy: currently one failed active deployment can fail daemon startup. That may be desirable for beta visibility, but production may prefer best-effort restoration with per-site degraded status.
- The store query for active runtime restoration should eventually be purpose-built rather than reusing admin inventory listing.

### What should be done in the future

- Add automated tests for startup restoration.
- Add authenticated smoke mode that creates/redeploys the demo site using an access token.
- Add a Makefile target for packaging `examples/hello-beta`.

### Code review instructions

Review commits:

```text
23b66ec Add beta smoke and OIDC access token support
f137ff9 Restore active runtimes on daemon startup
```

Review K3s commits:

```text
5779ad3 Bump go-go-host beta image
13ac467 Deploy go-go-host startup runtime restore
```

Validate live:

```bash
scripts/beta-smoke.sh
curl -fsS https://hosting.yolo.scapegoat.dev/api/v1/me -H "Authorization: Bearer <access-token>" | jq .
```
