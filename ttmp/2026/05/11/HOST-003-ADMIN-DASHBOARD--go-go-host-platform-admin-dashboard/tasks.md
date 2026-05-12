# Tasks

## Phase 0 — Ticket setup and scope

- [x] Create `HOST-003-ADMIN-DASHBOARD` ticket workspace.
- [x] Add design guide and implementation diary.
- [x] Split admin dashboard work out of `HOST-001-GO-GO-HOST-V1` while keeping backend/platform dependencies visible.
- [x] Relate core backend and frontend files to the ticket index/design doc.
- [x] Keep `docmgr doctor --ticket HOST-003-ADMIN-DASHBOARD --stale-after 30` green after each focused slice.

## Phase 1 — Admin shell, routing, and access guards

- [x] Serve the same embedded Vite SPA for `/admin` and nested `/admin/*` routes.
- [x] Add `RequirePlatformAdmin` route guard backed by `/api/v1/me.platformAdmin`.
- [x] Add an admin-specific shell/sidebar that is visually consistent with the user dashboard but clearly labeled as platform scope.
- [x] Add `/admin` redirect to `/admin/overview`.
- [x] Add admin denial state for authenticated non-platform-admin users.
- [ ] Add Storybook stories for admin shell and denied/loading states.
- [x] Add integration tests that `/admin`, `/admin/overview`, and nested admin routes serve the embedded SPA index.

## Phase 2 — Admin overview and runtime inventory MVP

- [x] Add RTK Query endpoint for `GET /api/v1/admin/runtimes/summary`.
- [x] Define TypeScript contracts for admin runtime summary responses.
- [x] Add `AdminOverviewPage` with active runtime/site/host/error counters.
- [x] Add `AdminRuntimesPage` with runtime rows, status pills, host list, deployment IDs, request/error counters, and last error summaries.
- [x] Add `AdminRuntimeTable` organism with Storybook stories for empty, healthy, degraded, and denied states.
- [x] Poll runtime summary on admin pages without hammering the API.
- [x] Verify non-admin API response renders a friendly platform-admin-required callout.

## Phase 3 — Platform inventory APIs

- [x] Add admin-only API endpoint for all orgs with counts and creation timestamps.
- [x] Add admin-only API endpoint for all sites with org, host, status, active deployment, and runtime status.
- [x] Add admin-only API endpoint for all deployments with org/site filters, deployment status filters, and limit.
- [x] Add admin-only API endpoint for global agents with org/site/grant context.
- [x] Add admin-only global audit endpoint with filters for actor, org, resource, action, and time range.
- [x] Add backend tests proving platform admins can query all tenants.
- [x] Add backend tests proving normal users receive 403 for every `/api/v1/admin/*` inventory endpoint.

## Phase 4 — Admin orgs/users/sites pages

- [x] Add `AdminOrgsPage` with org list, membership count, site count, deployment count, and quick links.
- [x] Add `AdminUsersPage` with known users, platform-admin marker, membership summary, and last activity placeholder.
- [x] Add `AdminSitesPage` with global site inventory and filters by org/status/runtime state.
- [x] Add Storybook stories and MSW fixtures for each inventory page.
- [x] Add empty/error/loading states for every page.

## Phase 5 — Admin deployments and audit pages

- [x] Add `AdminDeploymentsPage` with status/org/site/actor filters.
- [x] Add `AdminDeploymentDetailPage` for manifest, validation report, activation timeline, actor, and bundle metadata.
- [x] Add `AdminAuditPage` for global audit with URL-backed filters.
- [ ] Add Storybook interaction stories for filtering and bad JSON metadata.
- [x] Link runtime deployment IDs to the admin deployment detail route.

## Phase 6 — Runtime operations and safety controls

- [x] Design admin runtime stop/restart API with explicit audit logging.
- [x] Add backend stop/restart endpoints gated to platform admins.
- [x] Add themed confirmation dialog for destructive runtime actions.
- [x] Add `AdminRuntimeActions` component and stories for success/failure/disabled states.
- [x] Ensure operations update runtime summary cache and audit logs.
- [x] Add tests for forbidden runtime operation by non-admin users.

## Phase 7 — Quotas, capabilities, domains, and policy pages

- [ ] Add admin quota policy API design for defaults and per-site overrides.
- [ ] Add `AdminQuotasPage` with read-only current usage first, then editable policy forms.
- [ ] Add admin capability policy page showing requested vs effective capabilities.
- [ ] Add admin domain policy page for base domains, verification status, and custom-domain placeholders.
- [ ] Add Storybook stories for quota/capability/domain edge cases.

## Phase 8 — Admin agents and enrollment oversight

- [x] Add `AdminAgentsPage` with global agent list and revoke affordance.
- [ ] Add enrollment-key visibility once backend enrollment exists.
- [ ] Add grants inspection once `agent_site_grants` exists.
- [ ] Add audit trail links from each agent row.
- [x] Add Storybook stories for active/revoked/stale agents.

## Phase 9 — Visual polish, accessibility, and responsiveness

- [ ] Match the macOS1/HyperCard style established by the user dashboard.
- [ ] Keep os-core as the canonical theme package; admin CSS should only add app-specific mappings.
- [ ] Verify keyboard navigation for sidebar, filters, tables, and confirmation dialogs.
- [ ] Add responsive states for small laptop widths and narrow panes.
- [ ] Add missing favicon or suppress known favicon 404 noise.

## Phase 10 — E2E validation and release readiness

- [ ] Add Playwright smoke for non-admin denial under `/admin`.
- [x] Add Playwright smoke for admin overview/runtime inventory with seeded platform admin.
- [x] Add embedded SPA tests for `/admin/*` fallback behavior.
- [ ] Add devctl runbook for creating a local platform admin user.
- [ ] Run `go test ./...`, `make web-build`, `make storybook-build`, and `go run ./cmd/build-web`.
- [ ] Upload the final admin dashboard design/task bundle to reMarkable if requested.
