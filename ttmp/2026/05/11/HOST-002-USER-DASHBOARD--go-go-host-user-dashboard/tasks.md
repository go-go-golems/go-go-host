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

- [ ] Add React Router and define routes for `/app`.
- [ ] Create `AppShell` organism with top bar, org selector, and content slot.
- [ ] Create `OrgSidebar` organism.
- [ ] Create `RequireSession` guard.
- [ ] Create `RequireOrgAccess` guard.
- [ ] Create `OrgRedirectOrOnboarding` behavior:
  - [ ] no orgs -> onboarding,
  - [ ] one org -> org sites,
  - [ ] many orgs -> last selected or picker.
- [ ] Preserve selected org in URL.
- [ ] Add Storybook stories for shell states:
  - [ ] loading session,
  - [ ] unauthenticated/error,
  - [ ] no orgs,
  - [ ] one org,
  - [ ] many orgs,
  - [ ] dev auth banner.

Exit criteria:

- [ ] `/app` can route to an org sites page in Storybook/MSW.
- [ ] Route guards have visible loading/error/empty states.
- [ ] Shell uses os-core layout/theme primitives or documented local equivalents.

## Phase 6: Organization onboarding and site list/create

Goal: deliver the first useful dashboard workflow: create org/site and view sites.

- [ ] Create `pages/NoOrgsPage` with create-org form.
- [ ] Create `pages/SitesPage`.
- [ ] Create `pages/CreateSitePage`.
- [ ] Create `organisms/SitesTable`.
- [ ] Add runtime badge fan-out or placeholder strategy for site rows.
- [ ] Create site form validation for slug/name.
- [ ] Wire `POST /api/v1/orgs`.
- [ ] Wire `GET /api/v1/orgs/{org_id}/sites`.
- [ ] Wire `POST /api/v1/orgs/{org_id}/sites`.
- [ ] Add Storybook page stories:
  - [ ] no org onboarding empty/success/error,
  - [ ] sites empty,
  - [ ] sites populated,
  - [ ] sites load error,
  - [ ] create site valid/invalid/forbidden.

Exit criteria:

- [ ] User can create org and site from the UI against dev daemon.
- [ ] Site list renders real API data.
- [ ] Page stories cover the full workflow without a daemon.

## Phase 7: Site overview and runtime details

Goal: make each site inspectable.

- [ ] Create `SiteLayout`.
- [ ] Create `SiteHeader` organism.
- [ ] Create `SiteTabs` organism.
- [ ] Create `pages/SiteOverviewPage`.
- [ ] Create `pages/RuntimePage`.
- [ ] Create `organisms/RuntimeStatusPanel`.
- [ ] Wire `GET /api/v1/sites/{site_id}/runtime`.
- [ ] Link deployment ID in runtime panel to deployment detail when present.
- [ ] Add copy/open affordance for `primaryHost`.
- [ ] Add Storybook stories for runtime ready/stopped/failed/loading/forbidden.

Exit criteria:

- [ ] Site overview shows host, active deployment ID, runtime badge, counters, and last error.
- [ ] Runtime page can refresh runtime status.
- [ ] Runtime states are easy to review in Storybook.

## Phase 8: Deployment upload, list, detail, activate, rollback

Goal: expose the complete deployment loop in the dashboard.

- [ ] Create `pages/DeploymentsPage`.
- [ ] Create `pages/DeploymentDetailPage`.
- [ ] Create `organisms/DeploymentTimeline`.
- [ ] Create `organisms/DeploymentUploadPanel`.
- [ ] Create `organisms/ValidationReportPanel`.
- [ ] Wire `POST /api/v1/sites/{site_id}/deployments` multipart upload.
- [ ] Wire `GET /api/v1/sites/{site_id}/deployments`.
- [ ] Wire `GET /api/v1/deployments/{deployment_id}`.
- [ ] Wire `POST /api/v1/deployments/{deployment_id}/activate`.
- [ ] Wire `POST /api/v1/sites/{site_id}/rollback`.
- [ ] Add safe JSON parsing helpers for `manifestJson` and `validationJson`.
- [ ] Add confirmation dialogs for activate and rollback.
- [ ] Add Storybook stories:
  - [ ] upload idle,
  - [ ] upload progress,
  - [ ] validation success,
  - [ ] validation rejected,
  - [ ] active deployment,
  - [ ] superseded deployment,
  - [ ] rejected deployment detail,
  - [ ] activation error,
  - [ ] rollback confirmation.

Exit criteria:

- [ ] User can upload a valid bundle and read validation output.
- [ ] User can activate a deployment.
- [ ] User can roll back to a previous deployment.
- [ ] Rejected validation reports are presented as normal user-facing output, not uncaught errors.

## Phase 9: Agents page and audit page

Goal: expose the initial agent/audit APIs without overpromising future deploy-run features.

- [ ] Create `pages/AgentsPage`.
- [ ] Create `organisms/AgentsTable`.
- [ ] Wire `GET /api/v1/orgs/{org_id}/agents`.
- [ ] Wire `POST /api/v1/orgs/{org_id}/agents`.
- [ ] Wire `POST /api/v1/orgs/{org_id}/agents/{agent_id}/revoke`.
- [ ] Add explicit preview notice for keys/grants/enrollment not yet implemented.
- [ ] Create `pages/AuditPage`.
- [ ] Create `organisms/AuditTimeline`.
- [ ] Wire `GET /api/v1/orgs/{org_id}/audit` with query filters.
- [ ] Preserve audit filters in URL query params.
- [ ] Add Storybook stories:
  - [ ] agents empty/populated,
  - [ ] create agent success/error,
  - [ ] revoke confirm/error,
  - [ ] audit populated,
  - [ ] audit filtered empty,
  - [ ] audit selected metadata.

Exit criteria:

- [ ] Agents can be listed and created from dashboard.
- [ ] Agent revoke is confirmation-gated.
- [ ] Audit list filters by action/actor/resource.

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
