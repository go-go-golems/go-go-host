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
