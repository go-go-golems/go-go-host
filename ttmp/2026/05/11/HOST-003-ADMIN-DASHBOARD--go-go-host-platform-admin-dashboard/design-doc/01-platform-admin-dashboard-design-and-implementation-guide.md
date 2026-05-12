---
Title: Platform admin dashboard design and implementation guide
Ticket: HOST-003-ADMIN-DASHBOARD
Status: active
Topics:
    - dashboard
    - frontend
    - go-go-host
    - rtk-query
    - storybook
    - platform-admin
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: "Design and phased implementation guide for the go-go-host platform admin dashboard under /admin."
LastUpdated: 2026-05-11T23:20:00-04:00
WhatFor: "Use this when implementing admin-side dashboard routes, APIs, Storybook stories, and tests."
WhenToUse: "When working on HOST-003 admin dashboard features or platform-admin control plane views."
---

# Platform admin dashboard design and implementation guide

## Executive summary

`HOST-003-ADMIN-DASHBOARD` adds the platform operator dashboard for go-go-host. It should be built exactly like the user-facing dashboard from `HOST-002`: Storybook-first, RTK Query for all API calls, MSW fixtures for page states, embedded Vite SPA delivery, and macOS1/HyperCard visual language from `@go-go-golems/os-core`.

The user dashboard answers “what can an org member do with their sites?” The admin dashboard answers “what is happening across the whole host?” It lives under `/admin/*`, is gated by platform-admin identity, and provides global visibility into tenants, sites, runtimes, deployments, agents, audit events, quotas, capabilities, and domains.

## Problem statement

`/admin` currently shares the embedded SPA handler but has no real admin route tree or operator workflows. Backend platform-admin checks exist for runtime summary, and `/api/v1/me` exposes `platformAdmin`, but there is no coherent admin dashboard surface.

Operators need a safe control room that can:

- inspect all orgs/sites/runtimes without joining every org,
- identify failed runtimes and bad deployments quickly,
- inspect global audit activity,
- later stop/restart runtimes and manage policy safely,
- preserve tenant boundaries for normal users.

## Proposed solution

Build a dedicated admin dashboard route tree in the existing React app:

```text
/admin
/admin/overview
/admin/runtimes
/admin/orgs
/admin/users
/admin/sites
/admin/deployments
/admin/deployments/:deploymentId
/admin/agents
/admin/audit
/admin/quotas
/admin/capabilities
/admin/domains
```

Start with the smallest useful MVP:

1. `RequirePlatformAdmin` guard using `/api/v1/me`.
2. `AdminLayout` and `AdminSidebar` using the same `AppShell` primitives.
3. `GET /api/v1/admin/runtimes/summary` RTK Query integration.
4. `AdminOverviewPage` with summary counters.
5. `AdminRuntimesPage` and `AdminRuntimeTable` with Storybook stories.

Then expand backend inventory APIs and corresponding pages phase by phase.

## Design decisions

### Storybook-first implementation

Every admin component/page should have stories before or alongside route wiring. Stories should cover loading, empty, healthy, degraded, forbidden, and malformed-data cases using MSW or props.

### One SPA, two scopes

The Vite app serves both `/app/*` and `/admin/*`. This avoids duplicating build pipelines and keeps shared atoms/molecules/theme code in one package. User routes stay org-scoped; admin routes stay platform-scoped.

### Platform-admin guard at both layers

The browser guard improves UX, but every `/api/v1/admin/*` endpoint must independently check `core.Store.IsPlatformAdmin`. Non-admin users should get a friendly dashboard denial and a backend 403.

### Read-only first, operations later

The first admin dashboard slice is observability only. Destructive controls such as runtime stop/restart must come after confirmation dialogs, audit logging, and backend tests.

### Keep os-core canonical

`@go-go-golems/os-core` owns theme tokens. The local bridge CSS should only adapt dashboard-specific class names to the token contract.

## API surface plan

Existing:

- `GET /api/v1/me` includes `platformAdmin`.
- `GET /api/v1/admin/runtimes/summary` returns supervisor summary and is platform-admin gated.

Planned:

- `GET /api/v1/admin/orgs`
- `GET /api/v1/admin/users`
- `GET /api/v1/admin/sites`
- `GET /api/v1/admin/deployments`
- `GET /api/v1/admin/audit`
- `GET /api/v1/admin/agents`
- `POST /api/v1/admin/runtimes/{site_id}/restart`
- `POST /api/v1/admin/runtimes/{site_id}/stop`
- quota/capability/domain policy endpoints after backend tables stabilize.

## Component plan

Atoms/molecules should reuse HOST-002 components where possible:

- `StatusPill`
- `RuntimeStatusDot`
- `Timestamp`
- `CodeBlock`
- `JsonTree`
- `ErrorCallout`
- `LoadingBlock`
- `MetricCard`
- `RuntimeBadge`

New admin-specific organisms:

- `AdminSidebar`
- `AdminRuntimeTable`
- `AdminInventoryTable`
- `AdminAuditFilters`
- `AdminRuntimeActions`
- `AdminQuotaPolicyEditor`

New pages:

- `AdminOverviewPage`
- `AdminRuntimesPage`
- `AdminOrgsPage`
- `AdminUsersPage`
- `AdminSitesPage`
- `AdminDeploymentsPage`
- `AdminDeploymentDetailPage`
- `AdminAgentsPage`
- `AdminAuditPage`
- `AdminQuotasPage`
- `AdminCapabilitiesPage`
- `AdminDomainsPage`

## Alternatives considered

### Separate admin SPA

Rejected for v1 because it would duplicate the Vite, Storybook, Dagger, embed, theme, and RTK setup. A single SPA with route-level separation is simpler.

### Backend-rendered admin pages

Rejected because the user dashboard already established React/RTK/Storybook as the product surface, and admin workflows need the same tables, filters, and interaction patterns.

### Build operations before inventory

Rejected because read-only visibility is safer and gives immediate operator value without adding stop/restart risk.

## Implementation plan

See `tasks.md` for detailed phase tasks. The working order is:

1. Ticket/docs/tasks.
2. Admin shell/routing/guard.
3. Runtime summary data model and pages.
4. Backend inventory APIs.
5. Inventory pages.
6. Runtime operations with confirmations and audit.
7. Policy pages.
8. E2E validation.

## Validation plan

Run after each slice:

```bash
go test ./...
make web-build
make storybook-build
go run ./cmd/build-web
docmgr doctor --ticket HOST-003-ADMIN-DASHBOARD --stale-after 30
```

Browser smoke:

```text
http://127.0.0.1:8080/admin
http://127.0.0.1:8080/admin/overview
http://127.0.0.1:8080/admin/runtimes
http://127.0.0.1:6007/?path=/story/admin-pages-adminoverviewpage--with-runtimes
```
