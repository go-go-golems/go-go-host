---
Ticket: HOST-006-PRODUCTION-READINESS
Title: Tasks
Status: active
Topics:
    - go-go-host
    - hosting
    - security
    - deployments
    - agents
    - platform-admin
DocType: reference
Intent: long-term
---

# Tasks

## Documentation delivery

- [x] Create HOST-006 production readiness ticket.
- [x] Create primary design/implementation guide.
- [x] Create investigation diary.
- [x] Gather evidence from auth, frontend, devctl, runtime, settings, maintenance, docs, packaging, and E2E files.
- [x] Write production-readiness priority plan in order of necessity.
- [x] Add implementation phases, API sketches, pseudocode, test strategy, and file reference map.
- [x] Run `docmgr doctor`.
- [x] Upload final bundle to reMarkable.

## P0: Real local/staging auth

- [ ] Add Keycloak realm import with public dashboard client, redirect URIs, web origins, and seed users.
- [ ] Add `configs/dev.keycloak.yaml` with `devAuth: false`.
- [ ] Update devctl to start Keycloak and Keycloak Postgres for the OIDC profile.
- [ ] Expose OIDC frontend config from `/api/v1/config`.
- [ ] Add dashboard OAuth Authorization Code + PKCE login, callback, logout, and token refresh.
- [ ] Attach bearer tokens in RTK Query `prepareHeaders`.
- [ ] Add backend OIDC claim mapping tests for issuer, audience, expiry, signature, email, subject, and display name.
- [ ] Add Playwright login/logout/user-isolation E2E.

## P0: Platform-admin bootstrap

- [ ] Add config fields for OIDC admin subjects, emails, and/or roles.
- [ ] Add claim parsing for realm/client roles if Keycloak role bootstrap is chosen.
- [ ] Add admin bootstrap audit event.
- [ ] Add tests proving admin user can access `/admin` and normal users get `403`.
- [ ] Document beta operator bootstrap procedure.

## P0: Release/deploy pipeline

- [ ] Add image build/push GitHub Action.
- [ ] Add staging deployment recipe with health gate on `/readyz`.
- [ ] Decide and document migration policy.
- [ ] Add release image tags by version and commit SHA.
- [ ] Add rollback procedure for image and data preservation.

## P1: Domains/TLS

- [ ] Implement DNS TXT/CNAME verification checks.
- [ ] Add domain verification detail API and dashboard copy.
- [ ] Add fake DNS resolver tests.
- [ ] Write edge/TLS routing recipe for wildcard base domain and custom domains.
- [ ] Add domain recheck/expiry policy.

## P1: Runtime isolation

- [ ] Add Goja interrupt support on request timeout.
- [ ] Add per-site concurrency limiter.
- [ ] Review or restrict app-level `db.guard.configure`.
- [ ] Add crash-loop protection and restart policy.
- [ ] Decide subprocess/container isolation threshold for broader beta.
- [ ] Add security tests for CPU loops, DB hard limit, denied capabilities, and runtime panic containment.

## P1: Secrets and external app readiness

- [ ] Design encrypted site secrets table and key-management strategy.
- [ ] Add secrets API that never returns plaintext from list operations.
- [ ] Add JS `secrets` module with scoped runtime access.
- [ ] Add dashboard secrets UI with create/update/delete but no value display after save.
- [ ] Add redaction and audit tests.

## P1: Backup/restore and observability

- [ ] Add scheduled Postgres backups.
- [ ] Add per-site SQLite backup/snapshot workflow.
- [ ] Add bundle/object-store backup target.
- [ ] Add restore CLI/API and restore drill.
- [ ] Add Prometheus or OpenTelemetry metrics.
- [ ] Add alert rules and dashboards for auth, deploy, runtime, quota, disk, and agent-security failures.

## P2: Beta onboarding/support

- [ ] Add invitation/member-management workflow.
- [ ] Publish CLI docs as web docs or docs site.
- [ ] Add acceptable-use/beta terms and abuse-response contact.
- [ ] Add support runbooks and feedback-capture workflow.
- [ ] Add scheduled full E2E in CI/nightly.
