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
- [x] Rebuild and publish a new Docker image containing the auth cleanup.
- [x] Update K3s GitOps image pin and roll out the new image.
- [x] Verify live API calls with the browser access token after rollout.
- [x] Verify dashboard/GitHub login still reaches admin and org/site pages.

## P0: Repeatable beta smoke

- [x] Add durable `examples/hello-beta` bundle source matching the live demo site.
- [x] Add `scripts/beta-smoke.sh` for control-plane and live demo-site smoke.
- [x] Validate `scripts/beta-smoke.sh` against `https://hosting.yolo.scapegoat.dev` and `https://hello.hosting.yolo.scapegoat.dev` before image rollout.
- [x] Discover that pod restart/image rollout drops in-memory active runtime registrations.
- [x] Add daemon startup restoration for deployments whose database status is `active`.
- [x] Validate `scripts/beta-smoke.sh` again after deploying the startup-restore image.
- [x] Add `scripts/beta-agent-smoke.sh` authenticated agent publishing smoke.
- [x] Document required environment variables in the script help/error output.
- [x] Validate `scripts/beta-agent-smoke.sh` live with an access token; deployment `dep_353a2977-6f57-4602-b6bc-eb94754a664a` became active and temporary agent was revoked.
- [ ] Add a non-agent authenticated human upload smoke mode if still needed.

## P1: Demo-site lifecycle

- [x] Record live demo resources: org `beta-demo`, site `hello`, initial deployment `dep_181c0489-b037-4732-b7b3-3cc99bf4ea52`.
- [x] Preserve source for the demo app under `examples/hello-beta`.
- [x] Link `/assets/style.css` from the demo `/` page using `ui.link({ rel: "stylesheet", href: "/assets/style.css" })`.
- [x] Redeploy the styled demo as user deployment `dep_728e1491-30b9-49c0-b435-bbc0eb224a61`.
- [x] Redeploy the demo through agent publishing as deployment `dep_aba73759-dc63-47c4-9a32-ade076330a1a`.
- [ ] Add a make target for packaging the demo bundle.
- [ ] Decide whether the live demo site should be treated as permanent beta fixture data or recreated by smoke automation.

## P1: Agent publishing smoke

- [x] Create a beta deploy agent with an intentionally narrow `bundles/**` grant and observe upload rejection for archive paths.
- [x] Create a beta deploy agent with grant path `**`, enroll it with `go-go-host-agent`, deploy `examples/hello-beta`, and auto-activate it.
- [x] Verify the agent-published site at `https://hello.hosting.yolo.scapegoat.dev/`.
- [x] Revoke both temporary smoke-test agents after the test.
- [ ] Clarify/document the meaning of agent grant `path`: it currently constrains archive entry paths, not only the logical upload path.

## P0: Agent bundle-path semantics fix

- [x] Define `allowedBundlePaths` as the agent grant policy for logical deployment artifact paths, for example `bundles/**`, `bundles/previews/**`, or `bundles/releases/**`.
- [x] Define `bundlePath` as the deploy-run request field that carries the logical artifact path; it should match `allowedBundlePaths` and should not be confused with the uploaded tar/zip's internal file names.
- [x] Rename the operator CLI grant flag from `go-go-host agents create --path` to `--bundle-path`, keeping `--path` as a deprecated compatibility alias.
- [x] Rename the agent CLI deploy flag from `go-go-host-agent deploy --path` to `--bundle-path`, keeping `--path` as a deprecated compatibility alias.
- [x] Add API request aliases so `allowedBundlePaths` is preferred while existing `allowedPaths` requests continue to work during beta migration.
- [x] Add API request aliases so `bundlePath` is preferred while existing deploy-run `path` requests continue to work during beta migration.
- [x] Stop passing deploy-run allowed paths into `deploy.ValidateAndStore` as archive-entry `AllowedPaths`; agent grant paths now authorize the logical bundle path only.
- [x] Keep the regular bundle/archive safety validator unchanged for traversal, absolute paths, unsafe symlinks, manifest path validation, size limits, and capability policy.
- [x] Update response DTOs and docs to label the persisted DB field `allowed_paths` as logical bundle paths until a future DB migration renames it.
- [x] Add/adjust tests so an agent with `allowedBundlePaths: ["bundles/**"]` can deploy a normal archive containing `go-go-host.json`, `scripts/app.js`, and `assets/style.css` using `--bundle-path bundles/hello-beta.tar.gz`.
- [x] Add/adjust tests so the same agent is denied when requesting `--bundle-path private/hello-beta.tar.gz`.
- [x] Add CLI/help docs explaining that `--bundle` is the local file path and `--bundle-path` is the logical artifact path checked against the grant.
- [x] Re-run live agent publishing smoke with `allowedBundlePaths: ["bundles/**"]` and `go-go-host-agent deploy --bundle-path bundles/hello-beta-agent-smoke.tar.gz`.

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
- [x] New startup-restore image `sha-f137ff9` deployed to K3s.
- [x] Live beta smoke passes after image rollout without manual reactivation.
