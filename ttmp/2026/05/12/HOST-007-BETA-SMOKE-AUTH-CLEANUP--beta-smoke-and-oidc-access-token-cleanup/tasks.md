---
Ticket: HOST-007-BETA-SMOKE-AUTH-CLEANUP
Title: Tasks
Status: active
Topics:
    - go-go-host
    - hosting
    - security
    - deployments
    - platform-admin
DocType: reference
Intent: long-term
---

# Tasks

## Documentation and ticket setup

- [x] Create HOST-007 ticket workspace.
- [x] Create primary intern-facing guide document.
- [x] Create chronological investigation diary.
- [x] Backfill current beta findings from HOST-006: live demo site, wildcard TLS, and ID-token/access-token discovery.
- [x] Add detailed implementation task list.
- [ ] Relate all key files to the guide and diary.
- [ ] Run `docmgr doctor`.
- [ ] Upload updated bundle to reMarkable.

## P0: OIDC access-token cleanup

- [x] Confirm current behavior: backend accepts ID token but rejects Keycloak access token with `expected audience "go-go-host-dashboard" got []`.
- [x] Change frontend bearer token helper to prefer `accessToken` and fall back to `idToken` only when no access token is available.
- [x] Change backend verifier to validate OIDC bearer token signature, issuer, and expiry while accepting either matching `aud` or matching `azp`/authorized party.
- [x] Add unit coverage for audience/authorized-party token matching.
- [x] Keep platform-admin bootstrap behavior based on subject/email/roles/groups/client roles.
- [ ] Rebuild and publish a new Docker image containing the auth cleanup.
- [ ] Update K3s GitOps image pin and roll out the new image.
- [ ] Verify live API calls with the browser access token after rollout.
- [ ] Verify dashboard/GitHub login still reaches admin and org/site pages.

## P0: Repeatable beta smoke

- [x] Add durable `examples/hello-beta` bundle source matching the live demo site.
- [x] Add `scripts/beta-smoke.sh` for control-plane and live demo-site smoke.
- [x] Validate `scripts/beta-smoke.sh` against `https://hosting.yolo.scapegoat.dev` and `https://hello.hosting.yolo.scapegoat.dev` before image rollout.
- [x] Discover that pod restart/image rollout drops in-memory active runtime registrations.
- [x] Add daemon startup restoration for deployments whose database status is `active`.
- [ ] Validate `scripts/beta-smoke.sh` again after deploying the startup-restore image.
- [ ] Extend smoke script with optional authenticated create/upload/activate mode after access-token rollout.
- [ ] Document required environment variables for authenticated mode.

## P1: Demo-site lifecycle

- [x] Record live demo resources: org `beta-demo`, site `hello`, deployment `dep_181c0489-b037-4732-b7b3-3cc99bf4ea52`.
- [x] Preserve source for the demo app under `examples/hello-beta`.
- [ ] Add a make target for packaging the demo bundle.
- [ ] Decide whether the live demo site should be treated as permanent beta fixture data or recreated by smoke automation.

## P1: CLI/device-flow preparation

- [ ] Document that API clients should use access tokens, not ID tokens, once rollout is complete.
- [ ] Sketch future CLI OAuth Device Flow: `go-go-host login`, token cache, and authenticated API calls.
- [ ] Keep deploy-agent identity separate: enrollment token plus Ed25519 signed requests remains the machine-auth model.

## Validation checklist

- [x] `go test ./internal/httpapi ./internal/config ./internal/store ./internal/control ./internal/runtime`
- [x] `pnpm --dir web/admin build`
- [x] `BUILD_WEB_LOCAL=1 go run ./cmd/build-web`
- [x] `scripts/beta-smoke.sh`
- [x] `go test ./...`
- [x] New access-token image `sha-23b66ec` deployed to K3s.
- [x] Live access-token API verification passes against `/api/v1/me`.
- [ ] New startup-restore image deployed to K3s.
- [ ] Live beta smoke passes after image rollout without manual reactivation.
