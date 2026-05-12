---
Title: Implementation diary
Ticket: HOST-003-ADMIN-DASHBOARD
Status: active
Topics:
    - dashboard
    - frontend
    - go-go-host
    - rtk-query
    - storybook
    - platform-admin
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: "Chronological implementation diary for the platform admin dashboard."
LastUpdated: 2026-05-11T23:20:00-04:00
WhatFor: "Use this to resume admin dashboard work with context, validations, and known issues."
WhenToUse: "Before continuing HOST-003 implementation or reviewing admin dashboard changes."
---

# Implementation diary

## Step 1: Ticket setup and phase plan

The user asked to return to the first platform ticket, create a new ticket for the admin-side dashboard, approach it exactly like the user-facing dashboard, add detailed phased tasks, and begin work.

### Prompt Context

**User prompt (verbatim):** "Let's go back to the first ticket, and now create a new ticket to tackle the admin side dashboard, and approach it exactly the same way as the user facing dashboard. Then add detailed tasks split in phases in the new ticket, and then work on it"

**Assistant interpretation:** Split admin dashboard execution out of `HOST-001-GO-GO-HOST-V1` into a new dedicated ticket, mirror the HOST-002 Storybook-first frontend approach, then start with admin shell/routing/runtime-summary MVP.

### What I did

- Created ticket `HOST-003-ADMIN-DASHBOARD` titled `go-go-host platform admin dashboard`.
- Added this implementation diary.
- Added `Platform admin dashboard design and implementation guide`.
- Replaced the generated placeholder task file with a detailed 10-phase plan.

### Initial approach

The first implementation slice should stay read-only and use already-existing backend support:

- `/api/v1/me` for session and `platformAdmin` guard.
- `/api/v1/admin/runtimes/summary` for runtime inventory.
- Existing embedded SPA handler for `/admin/*`.
- Existing macOS1 theme package and dashboard components.

### Next implementation slice

- Add `RequirePlatformAdmin`.
- Add `AdminLayout` and `AdminSidebar`.
- Add `AdminOverviewPage`, `AdminRuntimesPage`, and `AdminRuntimeTable`.
- Add RTK Query endpoint and MSW fixtures.
- Add Storybook stories.
- Validate with Go tests, web build, Storybook build, Dagger build, docmgr doctor.

## Step 2: Admin dashboard shell and runtime summary MVP

I started the admin dashboard implementation using the same approach as the user dashboard: RTK Query APIs, reusable organisms, Storybook stories, MSW fixtures, embedded Vite SPA build, and Playwright browser verification.

### What changed

- Added `RequirePlatformAdmin` guard backed by `/api/v1/me.platformAdmin`.
- Added `/admin` route tree:
  - `/admin` redirects to `/admin/overview`.
  - `/admin/overview` renders the first platform overview page.
  - `/admin/runtimes` renders runtime inventory.
- Added `AdminLayout` using the existing `AppShell` plus an admin sidebar.
- Added `AdminSidebar` organism and stories.
- Added `AdminRuntimeTable` organism and stories.
- Added `AdminOverviewPage` with active site, runtime, host, request, and failed-runtime counters.
- Added `AdminRuntimesPage` with refresh and runtime detail rows.
- Added `AdminRuntimeSummary` TypeScript contract.
- Added RTK Query endpoint `useGetAdminRuntimeSummaryQuery` for `GET /api/v1/admin/runtimes/summary`.
- Added MSW fixture and handler for admin runtime summary.
- Extended embedded SPA integration test coverage to `/admin`, `/admin/overview`, and `/admin/runtimes`.
- Rebuilt embedded assets with Dagger.

### Validation

Commands run:

```bash
make web-build
go test ./...
make storybook-build
go run ./cmd/build-web
devctl restart go-go-hostd
devctl restart web-admin
devctl restart storybook
curl -fsS -o /tmp/admin-index.html -w '%{http_code}\n' http://127.0.0.1:8080/admin
curl -fsS -o /tmp/admin-runtimes.html -w '%{http_code}\n' http://127.0.0.1:8080/admin/runtimes
```

Results:

- `make web-build`: passed.
- `go test ./...`: passed.
- `make storybook-build`: passed.
- `go run ./cmd/build-web`: passed and exported embedded dist.
- `/admin` and `/admin/runtimes` returned HTTP 200 from `go-go-hostd`.

### Browser verification

Playwright checks:

- `http://127.0.0.1:8080/admin` renders the non-admin denial state for the current dev user.
- `http://127.0.0.1:6007/?path=/story/admin-pages-adminoverviewpage--with-runtimes` renders the admin overview story with runtime summary data.

Screenshots:

- `admin-dashboard-denied-mvp.png`
- `storybook-admin-overview-mvp.png`

### Issues and follow-ups

- Local dev user is not seeded as a platform admin, so embedded `/admin` currently verifies the denial path. Add a dev runbook or seed command for platform-admin browser verification.
- Storybook still emits noisy MSW warnings for unhandled static/module requests during browser inspection; production Storybook build passes.
- Admin shell denied state is functional but should get a dedicated Storybook guard story and perhaps a more centered panel treatment.

## Step 3: Admin inventory APIs and first inventory pages

I continued from the admin runtime MVP into the next platform inventory slice.

### What changed

Backend/store:

- Added sqlc admin inventory queries in `internal/store/queries/admin.sql`.
- Generated `internal/store/db/admin.sql.go`.
- Added store wrappers in `internal/store/admin.go` for:
  - org inventory,
  - user inventory,
  - site inventory with runtime status,
  - deployment inventory with org/site/status filters and limit.
- Added platform-admin-gated HTTP endpoints:
  - `GET /api/v1/admin/orgs`,
  - `GET /api/v1/admin/users`,
  - `GET /api/v1/admin/sites`,
  - `GET /api/v1/admin/deployments`.
- Reused the same `requirePlatformAdmin` helper for runtime summary and inventory endpoints.

Frontend:

- Added TypeScript contracts for admin org/user/site/deployment rows.
- Added RTK Query endpoints:
  - `useListAdminOrgsQuery`,
  - `useListAdminUsersQuery`,
  - `useListAdminSitesQuery`,
  - `useListAdminDeploymentsQuery`.
- Added MSW fixtures and handlers for the inventory endpoints.
- Added routes/pages:
  - `/admin/orgs` → `AdminOrgsPage`,
  - `/admin/users` → `AdminUsersPage`,
  - `/admin/sites` → `AdminSitesPage`,
  - `/admin/deployments` → `AdminDeploymentsPage`.
- Added Storybook stories for populated/empty/forbidden or filtered states.

### Validation

Commands run:

```bash
sqlc generate
make web-build
go test ./...
make storybook-build
go run ./cmd/build-web
```

Results:

- sqlc generation succeeded.
- TypeScript/Vite build passed.
- Go tests passed.
- Storybook production build passed.
- Dagger build-web exported embedded admin assets.

### Browser verification

Playwright checked:

- `http://127.0.0.1:6007/?path=/story/admin-pages-adminsitespage--populated`

Screenshot:

- `storybook-admin-sites-inventory.png`

### Follow-ups

- Add backend integration tests with a seeded platform admin to prove inventory endpoints return all tenants.
- Add a dev runbook or seed command so `/admin` can be inspected against the real embedded daemon as a platform admin.
- Add global admin audit and agent inventory endpoints next.

## Step 4: Dev platform-admin seeding and embedded admin verification

I added a local-dev way to seed platform-admin users so the real embedded `/admin` dashboard can be exercised without manually inserting rows into Postgres.

### What changed

- Added `devPlatformAdminSubjects` to daemon config.
- Default config and `configs/dev.yaml` now include `dev-user`, so the normal browser dev identity becomes a platform admin in local development.
- Updated dev auth to call `AddPlatformAdmin` after `UpsertUserFromOIDC` when the dev subject is configured as a platform admin.
- Added integration coverage for admin inventory endpoints:
  - non-admin dev users get `403`,
  - configured platform-admin dev users can query tenant org/site inventory.

### Validation

Commands run:

```bash
go test ./...
devctl restart go-go-hostd
curl -fsS http://127.0.0.1:8080/api/v1/me | jq '{email:.user.email, platformAdmin}'
curl -fsS http://127.0.0.1:8080/api/v1/admin/orgs | jq 'length'
```

Results:

- Go tests passed.
- `/api/v1/me` now reports `platformAdmin: true` for `dev-user` under `configs/dev.yaml`.
- `/api/v1/admin/orgs` returns inventory for the local dev database.

### Browser verification

Playwright checked embedded admin pages against the real daemon:

- `http://127.0.0.1:8080/admin/overview`
- `http://127.0.0.1:8080/admin/sites`

Screenshots:

- `embedded-admin-overview-dev-admin.png`
- `embedded-admin-sites-dev-admin.png`

### Notes

This is dev-auth-only behavior. Production OIDC users still need explicit `platform_admins` rows or a later admin bootstrap workflow.

## Step 5: Global admin agents and audit

I added the next operator inventory surfaces: global agents and global audit.

### What changed

Backend/store:

- Extended `internal/store/queries/admin.sql` with:
  - `ListAdminAgents`,
  - `ListAdminAuditEvents`.
- Regenerated sqlc output.
- Added store wrappers:
  - `ListAdminAgents`,
  - `ListAdminAuditEvents`.
- Added platform-admin-gated endpoints:
  - `GET /api/v1/admin/agents`,
  - `GET /api/v1/admin/audit`.

Frontend:

- Added `AdminAgent` TypeScript contract.
- Added RTK Query endpoints:
  - `useListAdminAgentsQuery`,
  - `useListAdminAuditQuery`.
- Added MSW fixtures/handlers for global agents and global audit.
- Added pages/routes:
  - `/admin/agents`,
  - `/admin/audit`.
- Added Storybook stories for both pages.

### Validation

Commands run:

```bash
sqlc generate
make web-build
go test ./...
make storybook-build
go run ./cmd/build-web
devctl restart go-go-hostd
curl -fsS http://127.0.0.1:8080/api/v1/admin/audit | jq 'length'
curl -fsS http://127.0.0.1:8080/api/v1/admin/agents | jq 'length'
```

Results:

- Web build passed.
- Go tests passed.
- Storybook build passed.
- Dagger embedded build passed.
- Embedded daemon returned global audit and agent rows for the seeded dev platform admin.

### Browser verification

Playwright checked embedded pages:

- `http://127.0.0.1:8080/admin/audit`
- `http://127.0.0.1:8080/admin/agents`

Screenshots:

- `embedded-admin-audit.png`
- `embedded-admin-agents.png`

### Follow-ups

- Add global agent revoke controls with confirmation once admin operation safety controls are in place.
- Add deployment detail under `/admin/deployments/:deploymentId` next.
- Add richer audit filters for actor ID/resource ID/time range.
