# Tasks

This task list intentionally splits the user dashboard into increasing phases. Do not implement all pages at once. Each phase should produce working, reviewed, Storybook-covered increments.

## Phase 0: Ticket and design setup

Goal: make the dashboard work trackable independently from the backend/platform ticket.

- [x] Create dedicated ticket `HOST-002-USER-DASHBOARD`.
- [x] Move dashboard design guide from `HOST-001-GO-GO-HOST-V1` into this ticket.
- [x] Add this phased task plan.
- [x] Relate initial dashboard implementation files as they are created.
- [x] Keep `reference/01-implementation-diary.md` updated after every implementation slice.
- [ ] Upload major design/task updates to reMarkable when they materially change.

Exit criteria:

- [x] Design guide lives in this ticket.
- [x] Dashboard work can proceed without mixing frontend tasks into the backend platform ticket.

## Phase 1: Frontend scaffold and build plumbing

Goal: create the smallest Vite/React/Storybook project that can build and be served later.

- [x] Create `web/admin/package.json` with React, Vite, TypeScript, RTK Query, Storybook, MSW, and `@go-go-golems/os-core` dependencies.
- [x] Create `web/admin/pnpm-lock.yaml` or align with workspace package manager if one exists.
- [x] Create `web/admin/tsconfig.json` and `tsconfig.node.json`.
- [x] Create `web/admin/vite.config.ts` with `/api` dev proxy to `go-go-hostd`.
- [x] Create `web/admin/index.html`.
- [x] Create `web/admin/src/main.tsx`.
- [x] Create `web/admin/src/app/App.tsx` with a minimal placeholder route.
- [x] Create `web/admin/.storybook/main.ts`.
- [x] Create `web/admin/.storybook/preview.tsx`.
- [x] Import `@go-go-golems/os-core` theme CSS or document exact package exports if import names differ.
- [x] Add root Makefile targets:
  - [x] `make web-install`
  - [x] `make web-dev`
  - [x] `make web-build`
  - [x] `make storybook`
  - [x] `make storybook-build`
- [x] Add local validation via Makefile targets for `pnpm build` and `pnpm storybook:build`.

Exit criteria:

- [x] `cd web/admin && pnpm build` passes.
- [x] `cd web/admin && pnpm storybook:build` passes.
- [x] At least one placeholder story renders.

## Phase 2: API types, RTK Query, MSW fixtures, and fake store

Goal: establish data contracts before building pages.

- [x] Create `web/admin/src/services/types.ts`.
- [ ] Add TypeScript types for:
  - [x] `ConfigResponse`
  - [x] `MeResponse`
  - [x] `Membership`
  - [x] `Org`
  - [x] `Site`
  - [x] `Deployment`
  - [x] `ValidationReport`
  - [x] `RuntimeStatus`
  - [x] `Agent`
  - [x] `AuditEvent`
- [x] Create `web/admin/src/services/goGoHostApi.ts` with RTK Query endpoints for current backend APIs.
- [x] Add API tags for cache invalidation: `Me`, `Org`, `Site`, `Deployment`, `Runtime`, `Agent`, `Audit`, `Config`.
- [x] Create `web/admin/src/app/store.ts`.
- [x] Create `web/admin/src/app/providers/AppProviders.tsx`.
- [x] Create `web/admin/src/app/providers/MockAppProviders.tsx` for Storybook.
- [x] Create `web/admin/src/services/msw/fixtures.ts` with realistic org/site/deployment/runtime/agent/audit data.
- [x] Create `web/admin/src/services/msw/handlers.ts` matching Go routes.
- [x] Create `web/admin/src/services/msw/browser.ts` and `server.ts` if needed by Storybook/tests.
- [x] Configure Storybook MSW addon.
- [x] Add story-level handler override examples for loading/error/empty states.

Exit criteria:

- [x] Page stories can fetch `/api/v1/me` through MSW.
- [x] Fixtures cover ready/stopped/failed runtime states and validated/rejected/active deployments.
- [x] No Storybook story requires a live daemon.

## Phase 3: Design tokens, atoms, and primitive stories

Goal: build the visual vocabulary before page composition.

- [x] Create `components/atoms/RuntimeStatusDot` with story states: ready, failed, stopped, starting, draining.
- [x] Create `components/atoms/StatusPill` with story states for site/deployment/agent statuses.
- [x] Create `components/atoms/RoleBadge` with owner/developer/viewer stories.
- [x] Create `components/atoms/CopyButton` with idle/copied/error stories.
- [x] Create `components/atoms/EmptyState` with and without action.
- [x] Create `components/atoms/ErrorCallout` with auth/network/validation examples.
- [x] Create `components/atoms/LoadingBlock` with small/large variants.
- [x] Create `components/atoms/Timestamp` with absolute/relative/empty variants.
- [x] Create `components/atoms/CodeBlock` for shell and JSON examples.
- [x] Create `components/atoms/JsonTree` for manifest and validation report examples.
- [x] Ensure each implemented atom directory contains:
  - [x] `Component.tsx`
  - [x] `Component.stories.tsx`
  - [x] `index.ts`
- [x] Use os-core tokens/CSS variables instead of hardcoded colors wherever possible.

Exit criteria:

- [ ] All atom stories render in Storybook.
- [ ] Atom stories include non-happy states where relevant.
- [ ] Atoms do not import RTK Query.

## Phase 4: Molecules and interaction widgets

Goal: build reusable dashboard controls and panels from atoms.

- [x] Create `components/molecules/RuntimeBadge`.
- [x] Create `components/molecules/OrgSwitcher`.
- [x] Create `components/molecules/SiteHostCopy`.
- [x] Create `components/molecules/DeploymentStatusPill`.
- [x] Create `components/molecules/ValidationSummary`.
- [ ] Create `components/molecules/FileDropZone`.
- [ ] Create `components/molecules/ConfirmActionDialog`.
- [ ] Create `components/molecules/FilterToolbar`.
- [x] Create `components/molecules/MetricCard`.
- [x] Create `components/molecules/ManifestSummary`.
- [ ] Create `components/molecules/AuditEventRow`.
- [x] Create `components/molecules/AgentStatusBadge`.
- [ ] Add interaction tests for `FileDropZone` and `ConfirmActionDialog` if test tooling is available.

Exit criteria:

- [x] Every implemented molecule has Storybook stories.
- [ ] Stories cover loading/error/empty/permission states where relevant.
- [x] Implemented molecules remain prop-driven and do not own route-level data fetching.

## Phase 5: Shell, routing, and session/org guards

Goal: make `/app` navigable with real route structure and mocked data.

- [x] Add React Router and define routes for `/app`.
- [x] Create `AppShell` organism with top bar, org selector, and content slot.
- [x] Create `OrgSidebar` organism.
- [x] Create `RequireSession` guard.
- [x] Create `RequireOrgAccess` guard.
- [x] Create `OrgRedirectOrOnboarding` behavior:
  - [x] no orgs -> onboarding,
  - [x] one org -> org sites,
  - [ ] many orgs -> last selected or picker.
- [x] Preserve selected org in URL.
- [x] Add initial Storybook stories for shell states:
  - [ ] loading session,
  - [ ] unauthenticated/error,
  - [x] no orgs,
  - [x] one org,
  - [x] many orgs,
  - [x] dev auth banner.

Exit criteria:

- [x] `/app` can route to an org sites page in Storybook/MSW.
- [x] Route guards have visible loading/error/empty states.
- [x] Shell uses os-core macOS1 theme scope and documented local bridge equivalents.

## Phase 6: Organization onboarding and site list/create

Goal: deliver the first useful dashboard workflow: create org/site and view sites.

- [x] Create `pages/NoOrgsPage` onboarding shell with create-org affordance placeholder.
- [x] Create `pages/SitesPage`.
- [x] Create `pages/CreateSitePage` form shell.
- [x] Create `organisms/SitesTable`.
- [x] Add runtime badge fan-out or placeholder strategy for site rows.
- [x] Create site form validation for slug/name.
- [x] Wire `POST /api/v1/orgs`.
- [x] Wire `GET /api/v1/orgs/{org_id}/sites`.
- [x] Wire `POST /api/v1/orgs/{org_id}/sites`.
- [x] Add Storybook page stories:
  - [x] no org onboarding empty state,
  - [x] sites empty,
  - [x] sites populated,
  - [x] sites load error,
  - [x] create site valid/invalid/forbidden.

Exit criteria:

- [x] User can create org and site from the UI against dev daemon.
- [x] Site list renders real API data.
- [x] Page stories cover the full workflow without a daemon.

## Phase 7: Site overview and runtime details

Goal: make each site inspectable.

- [x] Create `SiteLayout`.
- [x] Create `SiteHeader` organism.
- [x] Create `SiteTabs` organism.
- [x] Create `pages/SiteOverviewPage`.
- [x] Create `pages/RuntimePage`.
- [x] Create `organisms/RuntimeStatusPanel`.
- [x] Wire `GET /api/v1/sites/{site_id}/runtime`.
- [ ] Link deployment ID in runtime panel to deployment detail when present.
- [x] Add copy/open affordance for `primaryHost`.
- [x] Add Storybook stories for runtime ready/stopped/failed/forbidden.

Exit criteria:

- [x] Site overview shows host, active deployment ID, runtime badge, counters, and last error.
- [x] Runtime page can refresh runtime status.
- [x] Runtime states are easy to review in Storybook.

## Phase 8: Deployment upload, list, detail, activate, rollback

Goal: expose the complete deployment loop in the dashboard.

- [x] Create `pages/DeploymentsPage`.
- [x] Create `pages/DeploymentDetailPage`.
- [x] Create `organisms/DeploymentTimeline`.
- [x] Create `organisms/DeploymentUploadPanel`.
- [x] Create `organisms/ValidationReportPanel`.
- [x] Wire `POST /api/v1/sites/{site_id}/deployments` multipart upload.
- [x] Wire `GET /api/v1/sites/{site_id}/deployments`.
- [x] Wire `GET /api/v1/deployments/{deployment_id}`.
- [x] Wire `POST /api/v1/deployments/{deployment_id}/activate`.
- [x] Wire `POST /api/v1/sites/{site_id}/rollback`.
- [x] Add safe JSON parsing helpers for `manifestJson` and `validationJson`.
- [x] Add confirmation dialogs for activate and rollback.
- [ ] Add Storybook stories:
  - [x] upload idle,
  - [ ] upload progress,
  - [x] validation success,
  - [x] validation rejected,
  - [x] active deployment,
  - [x] superseded deployment,
  - [x] rejected deployment detail,
  - [x] activation error,
  - [x] rollback confirmation.

Exit criteria:

- [x] User can upload a valid bundle and read validation output.
- [x] User can activate a deployment.
- [x] User can roll back to a previous deployment.
- [x] Rejected validation reports are presented as normal user-facing output, not uncaught errors.

## Phase 9: Agents page and audit page

Goal: expose the initial agent/audit APIs without overpromising future deploy-run features.

- [x] Create `pages/AgentsPage`.
- [x] Create `organisms/AgentsTable`.
- [x] Wire `GET /api/v1/orgs/{org_id}/agents`.
- [x] Wire `POST /api/v1/orgs/{org_id}/agents`.
- [x] Wire `POST /api/v1/orgs/{org_id}/agents/{agent_id}/revoke`.
- [x] Add explicit preview notice for keys/grants/enrollment not yet implemented.
- [x] Create `pages/AuditPage`.
- [x] Create `organisms/AuditTimeline`.
- [x] Wire `GET /api/v1/orgs/{org_id}/audit` with query filters.
- [x] Preserve audit filters in URL query params.
- [ ] Add Storybook stories:
  - [x] agents empty/populated,
  - [x] create agent success/error,
  - [x] revoke confirm/error,
  - [x] audit populated,
  - [x] audit filtered empty,
  - [x] audit selected metadata.

Exit criteria:

- [x] Agents can be listed and created from dashboard.
- [x] Agent revoke is confirmation-gated.
- [x] Audit list filters by action/actor/resource.

## Phase 10: Usage and members placeholder pages

Goal: reserve product surfaces for near-future backend features while showing useful current data.

- [ ] Create `pages/UsagePage`.
- [ ] Create `organisms/QuotaPanel` with API-pending state.
- [ ] Show runtime counters from `GET /api/v1/sites/{site_id}/runtime`.
- [ ] Create `pages/MembersPage`.
- [ ] Create `organisms/MembersTable` based on `/api/v1/me` memberships.
- [ ] Add API-pending callouts for membership mutation and quota endpoints.
- [ ] Add Storybook stories for quota pending and role variants.

Exit criteria:

- [ ] Users can see request/error counters in usage.
- [ ] Users understand which usage/member features are pending backend APIs.

## Phase 11: Backend embedding and SPA serving

Goal: serve the built dashboard from the Go daemon.

- [ ] Add production build copy/embed flow for `web/admin/dist`.
- [ ] Replace `internal/webadmin.NewPlaceholderHandler()` with embedded SPA handler.
- [ ] Preserve `/api/*` route behavior.
- [ ] Serve `/app/` and nested `/app/orgs/...` routes as SPA index.
- [ ] Keep `/admin/` separate, either placeholder or future admin bundle.
- [ ] Add Go tests:
  - [ ] `/app/` returns dashboard index,
  - [ ] `/app/orgs/org_123/sites` returns dashboard index,
  - [ ] `/api/v1/version` still returns JSON,
  - [ ] unknown public host fallback still reaches runtime supervisor.

Exit criteria:

- [ ] `go-go-hostd` serves the dashboard from `/app/` in production mode.
- [ ] SPA routing and API routing do not conflict.

## Phase 12: End-to-end validation and release polish

Goal: make dashboard safe to hand to users.

- [ ] Add Playwright smoke test for dev dashboard login/session bootstrap.
- [ ] Add Playwright smoke test for site list rendering.
- [ ] Add Playwright smoke test for deployment upload using a fixture bundle if feasible.
- [ ] Add Storybook accessibility checks if tooling is available.
- [ ] Add `README` section for dashboard development.
- [ ] Add help doc or ticket note explaining dashboard dev workflow.
- [ ] Run full validation:
  - [ ] `go test ./...`
  - [ ] `pnpm build`
  - [ ] `pnpm storybook:build`
  - [ ] Playwright smoke tests.

Exit criteria:

- [ ] Dashboard is usable for create-site, deploy, activate, runtime inspect, agents list/create, and audit list.
- [ ] Storybook is the authoritative component/page review surface.
- [ ] CI/local validation covers both Go and web assets.
