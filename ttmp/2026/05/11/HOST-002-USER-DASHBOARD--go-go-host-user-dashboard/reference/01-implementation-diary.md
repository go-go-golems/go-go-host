---
Title: Implementation Diary
Ticket: HOST-002-USER-DASHBOARD
Status: active
Topics:
    - go-go-host
    - frontend
    - storybook
    - rtk-query
    - dashboard
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: "Chronological diary for the go-go-host user dashboard ticket."
LastUpdated: 2026-05-11T19:15:00-04:00
WhatFor: "Track dashboard design and implementation decisions, commands, validation, failures, and review notes."
WhenToUse: "Append after every dashboard design or implementation slice."
---

# Implementation Diary

## Step 1: Split user dashboard work into its own ticket

The user dashboard design was originally created in the backend/platform ticket `HOST-001-GO-GO-HOST-V1`. The user asked to create a separate ticket for dashboard implementation, move the design guide over, and add a phased task list so the dashboard can be built incrementally rather than all at once.

### Prompt Context

**User prompt (verbatim):** "Let's create a separate ticket for the user dashboard, and move the design guide over, and add a detailed list of tasks, split in phases, to build teh dashboard page in increasing phases, so that we don't do everything at once."

**Assistant interpretation:** Create a dedicated docmgr ticket, move/copy the dashboard design guide into it, replace the original with a pointer, and create a detailed phased dashboard implementation plan.

**Inferred user intent:** Decouple frontend/dashboard tracking from the backend platform ticket and make the implementation sequence reviewable in smaller increments.

### What I did

- Created ticket `HOST-002-USER-DASHBOARD`.
- Added design doc `design-doc/01-user-dashboard-affordances-page-designs-and-component-system-guide.md`.
- Moved the Phase 7 dashboard design guide content into the new ticket and updated the frontmatter ticket metadata.
- Replaced the old `HOST-001` copy with a short archived pointer to the new ticket.
- Added a detailed phased task list in `tasks.md`.
- Added this implementation diary.

### Why

The dashboard is large enough to deserve its own ticket. It spans React, Storybook, MSW, RTK Query, Go embed serving, Playwright, and a component system. Keeping it in the backend v1 ticket would make planning and review harder.

### What worked

- `docmgr ticket create-ticket` created the new workspace cleanly.
- `docmgr doc add` created the new design doc and diary skeletons.
- The previous design guide could be moved without losing its technical detail.

### What didn't work

- N/A.

### What was tricky

- This was a move rather than a pure copy. To avoid two active source-of-truth docs, the old ticket now keeps only a historical pointer to the new ticket.

### What warrants review

- Confirm ticket ID `HOST-002-USER-DASHBOARD` is acceptable.
- Review task phases before implementation begins to ensure scope and order match product priorities.

### Next steps

- Relate the new dashboard design doc to backend API files.
- Run `docmgr doctor` for the new ticket.
- Upload the new ticket bundle to reMarkable if requested or after validation.

## Step 2: Scaffolded Vite/React dashboard, Storybook, MSW, and devctl services

I started implementation for `HOST-002-USER-DASHBOARD` by building the minimal frontend scaffold and making Storybook runnable through devctl so progress can be watched live.

### Prompt Context

**User prompt (verbatim):** "go ahead, commit at appropriate intervals, keep a detailed diary" followed by "Create a devctl setup to run storybook as well" and "then run storybook, so i can track your progress."

**Assistant interpretation:** Begin Phase 1/2 implementation, add devctl supervision for Storybook, run Storybook in the background, validate build/storybook, update docs, and commit a focused slice.

**Inferred user intent:** The user wants an observable frontend development loop, not just static code changes.

### What I did

- Created `web/admin` Vite/React/TypeScript scaffold.
- Added Storybook 8 with React/Vite.
- Added RTK Query store and initial `goGoHostApi` service.
- Added MSW fixtures and handlers for current backend routes.
- Added Storybook MSW setup and mock Redux provider.
- Added a local `@go-go-golems/os-core` shim package under `web/admin/vendor/os-core` because the real private package is currently inaccessible to npm in this environment.
- Added the first atom:
  - `components/atoms/StatusPill`
- Added the first page story:
  - `pages/AppBootstrapPage`
- Added Makefile targets for install/dev/build/storybook/storybook-build.
- Added a devctl plugin and `.devctl.yaml` with services:
  - `storybook` on `http://127.0.0.1:6007`
  - `web-admin` on `http://127.0.0.1:5173`
- Ran Storybook through devctl.

### Why

This gives the frontend work a working inner loop immediately. The first Storybook page already exercises RTK Query + MSW + fake store, which is the pattern every later page should follow.

### What worked

- `pnpm install` completed successfully.
- `make web-build` passes.
- `make storybook-build` passes.
- `devctl plugins list` and `devctl plan` work.
- `devctl up --force` started both services.
- Storybook is reachable at `http://127.0.0.1:6007`.
- Vite dev server is reachable at `http://127.0.0.1:5173`.

### What didn't work

- `npm view @go-go-golems/os-core version` failed with a GitHub Packages 403. I added a local package shim named `@go-go-golems/os-core` that exports `theme.css` and `themes/desktop.css`. This keeps imports and theming shape compatible while allowing local Storybook to run. The shim should be replaced with the real package once credentials/package access are available.
- Storybook initially tried port `6006`, but another node process was already listening there. I changed this project to use `6007` and updated the devctl health check accordingly.
- The first `pnpm build` failed because the stories imported `@storybook/react` types but only `@storybook/react-vite` was installed. Adding `@storybook/react` fixed TypeScript.

### What I learned

- Storybook 8 builds cleanly with Vite 6 and React 19 in this repo.
- The devctl plugin protocol is enough for a repo-local Storybook/Vite supervision setup with no custom daemon code.

### What was tricky

- Keeping the os-core dependency requirement while avoiding a blocked private package install required a local shim. I intentionally kept the import paths as `@go-go-golems/os-core/theme.css` so replacing the shim later should be mechanical.

### What warrants review

- Confirm whether using Storybook port `6007` is acceptable given the existing process on `6006`.
- Confirm whether the local os-core shim should be committed short-term or replaced immediately with a proper private registry setup.
- Review whether devctl should launch only Storybook or both Storybook and Vite; I included both so dashboard and component work can be watched separately.

### Validation

Commands run:

```bash
pnpm install
make web-build
make storybook-build
devctl plugins list
devctl plan
devctl up --force
devctl status --tail-lines 20
curl -I http://127.0.0.1:6007
curl -I http://127.0.0.1:5173
go test ./...
```

### Current running services

- Storybook: `http://127.0.0.1:6007`
- Vite dashboard dev server: `http://127.0.0.1:5173`

Use:

```bash
cd go-go-host
devctl status --tail-lines 20
devctl logs --service storybook
devctl logs --service web-admin
```

## Step 3: Fixed MSW worker registration and completed initial atom set

I fixed the Storybook MSW registration error by installing the MSW service worker into `web/admin/public` and serving that directory through Storybook. Then I continued Phase 3 by adding the initial atom set with stories.

### Prompt Context

**User prompt (verbatim):** "[MSW] Failed to register the Service Worker: on storybook" followed by "continue"

**Assistant interpretation:** First fix the visible Storybook/MSW runtime issue, then continue implementing dashboard components while keeping Storybook live.

**Inferred user intent:** Keep the live Storybook environment healthy and continue building reviewable component increments.

### What I did

- Ran `pnpm exec msw init public --save` in `web/admin`.
- Added `web/admin/public/mockServiceWorker.js`.
- Added `staticDirs: ['../public']` to Storybook config so the worker is served from Storybook.
- Restarted Storybook through devctl.
- Added atom components and stories:
  - `RuntimeStatusDot`
  - `RoleBadge`
  - `CopyButton`
  - `EmptyState`
  - `ErrorCallout`
  - `LoadingBlock`
  - `Timestamp`
  - `CodeBlock`
  - `JsonTree`
- Added `components/atoms/index.ts` barrel exports.
- Added the missing `Org` TypeScript type.
- Expanded fixtures to cover ready/stopped/failed runtime states and active/superseded/rejected deployments.
- Updated Phase 2/3 task checkboxes.

### Why

MSW page stories need the service worker to be available at runtime; otherwise page stories fail even if `storybook:build` succeeds. The atom set establishes the low-level UI vocabulary required before molecules and page layouts.

### What worked

- `make web-build` passes.
- `make storybook-build` passes.
- Storybook was restarted and is alive at `http://127.0.0.1:6007`.
- Vite dev server remains alive at `http://127.0.0.1:5173`.

### What didn't work

- The MSW worker was initially missing because Storybook does not automatically serve the worker unless it exists in a static directory. Adding `public/mockServiceWorker.js` and Storybook `staticDirs` fixed it.

### What I learned

- For this Storybook/MSW setup, `public` must be treated as part of the Storybook static asset contract, not just a Vite app convention.

### What was tricky

- Keeping atoms small and non-data-fetching while still making stories useful required story variants for states rather than embedding API behavior.

### What warrants review

- `CopyButton` relies on `navigator.clipboard`; if Storybook browser permissions are restrictive, we may want an injectable `copyText` helper for tests.
- `Timestamp` uses browser locale rendering; if deterministic screenshots become important, we should add a formatter abstraction.

### Validation

Commands run:

```bash
cd web/admin && pnpm exec msw init public --save
make web-build
make storybook-build
devctl restart storybook
devctl status --tail-lines 10
```

## Step 4: Switched from local os-core shim to real macOS1 theme imports

The user asked why the dashboard was not using the retro macOS1 look from go-go-os-core. I had initially added a local os-core shim because npm access to the private/package-registry version returned 403, and then added `macos1-bridge.css` as an adapter from the scaffold's generic dashboard classes to HyperCard/macOS1 CSS variables. That was only a temporary bridge, not the final intended architecture.

### What I changed

- Pointed `@go-go-golems/os-core` at the local go-go-os-frontend package checkout using `link:/home/manuel/workspaces/2026-05-11/npm-packages-go-go-os/go-go-os-frontend/packages/os-core`.
- Replaced shim imports with real os-core imports:
  - `@go-go-golems/os-core/theme`
  - `@go-go-golems/os-core/desktop-theme-macos1`
- Wrapped the app and Storybook stories in the required HyperCard scope manually:
  - `data-widget="hypercard"`
  - `className="theme-macos1"`
- Kept `macos1-bridge.css` as a compatibility adapter so the dashboard scaffold's current classes inherit macOS1 variables while we migrate components toward os-core primitives and `data-part` selectors.
- Updated `CopyButton` to expose `data-part="btn"` so os-core button styling can apply even without importing the TS `Btn` component.

### Why the bridge exists

The bridge exists because the first dashboard scaffold used local class names like `.dashboard-panel`, `.site-card`, and `.status-pill`, while the macOS1 theme in os-core is scoped to `[data-widget="hypercard"].theme-macos1` and primarily styles os-core `data-part` primitives. The bridge maps our early scaffold classes onto the os-core macOS1 tokens. It should shrink over time as we replace class-only widgets with true os-core primitives or matching `data-part` markup.

### What did not work

Using the local package as a `file:` dependency caused pnpm to pack only part of the source tree, so TypeScript could not resolve several os-core internal exports. Switching to a `link:` dependency fixed the missing files because node_modules now points at the full local checkout.

### Validation

Commands run:

```bash
cd web/admin && pnpm install
make web-build
make storybook-build
devctl restart storybook
devctl status --tail-lines 5
```

Storybook is alive at `http://127.0.0.1:6007`.

## Step 5: Added first molecule set for dashboard workflows

I continued Phase 4 by building the first prop-driven molecule components on top of the atom set and the macOS1-themed Storybook environment.

### Prompt Context

**User prompt (verbatim):** "continue"

**Assistant interpretation:** Continue implementation after fixing the macOS1/os-core setup, keeping Storybook live and adding the next layer of reusable UI components.

**Inferred user intent:** Move from atoms toward dashboard workflow components while preserving Storybook coverage.

### What I did

- Added molecule components and stories:
  - `RuntimeBadge`
  - `SiteHostCopy`
  - `DeploymentStatusPill`
  - `ValidationSummary`
  - `MetricCard`
  - `ManifestSummary`
  - `AgentStatusBadge`
  - `OrgSwitcher`
- Added `components/molecules/index.ts` barrel exports.
- Kept molecules prop-driven and free of RTK Query data fetching.
- Used existing atoms where appropriate:
  - `RuntimeStatusDot`
  - `StatusPill`
  - `CopyButton`
  - `ErrorCallout`
  - `JsonTree`
  - `RoleBadge`
- Validated build and Storybook build.
- Restarted Storybook through devctl so the new stories are visible.

### Why

Molecules are the next layer needed before building organisms and pages. These components represent reusable product concepts: runtime state, deployment status, validation output, manifest display, metrics, site host copy actions, agent status, and org switching.

### What worked

- `make web-build` passes.
- `make storybook-build` passes.
- Storybook restarted and remains alive on `http://127.0.0.1:6007`.

### What didn't work

- N/A in this slice.

### What I learned

- The macOS1 bridge is sufficient for early molecules, but upcoming organisms should use more os-core `data-part` conventions directly to reduce custom bridge styling.

### What was tricky

- Keeping `ValidationSummary` useful without overfitting to current backend JSON required treating validation errors/warnings/capabilities as optional and displayable independently.

### What warrants review

- `SiteHostCopy` currently exposes copy-host and copy-curl actions. Later we may add an explicit open-public-site action once host/public URL behavior is finalized.
- `OrgSwitcher` is prop-only for now; route integration belongs in the shell/routing phase.

### Validation

Commands run:

```bash
make web-build
make storybook-build
devctl restart storybook
devctl status --tail-lines 5
```

## Step 6: Added first organism set for shell, tables, timelines, and panels

I continued from molecules into the first organism layer. These components are still Storybook-first and mostly prop-driven, but they now represent larger dashboard regions that can be composed into pages.

### Prompt Context

**User prompt (verbatim):** "continue"

**Assistant interpretation:** Continue building the dashboard incrementally, moving from molecules toward organisms while keeping Storybook live and validated.

### What I did

- Added organism components and stories:
  - `AppShell`
  - `OrgSidebar`
  - `SitesTable`
  - `RuntimeStatusPanel`
  - `DeploymentTimeline`
  - `ValidationReportPanel`
  - `AgentsTable`
  - `AuditTimeline`
- Added `components/organisms/index.ts` barrel exports.
- Reused atoms and molecules instead of duplicating status/copy/metric/JSON UI.
- Kept organisms mostly prop-driven so page data fetching can remain in page components.
- Rebuilt and restarted Storybook.

### Why

Organisms are the next step before real pages. They turn the atom/molecule vocabulary into dashboard sections: shell chrome, site lists, runtime panels, deployment history, validation reports, agent tables, and audit timelines.

### What worked

- `make web-build` passes.
- `make storybook-build` passes.
- Storybook restarted and remains alive at `http://127.0.0.1:6007`.

### What didn't work

- N/A.

### What was tricky

- `ValidationReportPanel` needs to parse `manifestJson` and `validationJson`, which are currently strings in the API DTO. I added safe parsing fallback behavior so malformed JSON becomes visible UI state rather than a crash.

### What warrants review

- `AuditTimeline` and `DeploymentTimeline` use simple responsive grids. They may need mobile-specific treatment once pages are wired.
- `AppShell` is currently a static shell organism. Route integration and real navigation state still belong to the shell/routing phase.

### Validation

Commands run:

```bash
make web-build
make storybook-build
devctl restart storybook
devctl status --tail-lines 5
```

## Step 7: Added initial routing guards and sites pages

I moved from component stories into the first
I moved from component stories into the first real application-routing slice. The dashboard now has `/app` routing, session/org guards, an organization layout, and the first sites pages backed by RTK Query and MSW stories.

### Prompt Context

**User prompt (verbatim):** "continue"

**Assistant interpretation:** Continue the Storybook-first dashboard implementation by wiring Phase 5 shell/routing and beginning Phase 6 page work.

### What I did

- Added router setup:
  - `src/app/routes.tsx`
  - `src/app/routing/guards.tsx`
  - `src/app/routing/OrgRedirectOrOnboarding.tsx`
  - `src/app/routing/OrgLayout.tsx`
- Replaced the bootstrap-only `App` with a `RouterProvider` app.
- Added session/org guards:
  - `RequireSession`
  - `RequireOrgAccess`
- Added organization redirect/onboarding behavior:
  - no memberships -> `NoOrgsPage`
  - first membership -> org sites page
- Added initial pages and stories:
  - `NoOrgsPage`
  - `SitesPage`
  - `CreateSitePage`
- Wired `SitesPage` to `useListSitesQuery` for `GET /api/v1/orgs/{org_id}/sites`.
- Added MSW-backed page story states for populated, empty, and load error site lists.

### Why

The component catalog is now large enough to start validating real dashboard composition. Routing and guards let later pages share consistent shell behavior instead of each page reinventing loading, access-denied, and organization selection flows.

### What worked

- `make web-build` passes.
- `make storybook-build` passes.
- Storybook was restarted and remains alive through `devctl`.

### What didn't work

- N/A.

### What was tricky

- Page stories need their own `MemoryRouter` route context so hooks like `useParams` and `useNavigate` work in isolation.
- The sites page currently uses fixture runtime status data because runtime-by-site aggregation is not yet exposed as a page API; this should be replaced when the page detail/runtime APIs are wired.

### What warrants review

- `CreateSitePage` is only a form shell; mutation wiring and validation still need the next slice.
- Agents and audit routes are placeholders in the shared org layout and should be replaced by real pages next.

### Validation

Commands run:

```bash
make web-build
make storybook-build
devctl restart storybook
devctl status --tail-lines 5
```

## Step 8: Wired create-site form mutation and interaction stories

I finished the first useful site-management workflow slice by turning `CreateSitePage` from a form shell into a validated mutation-backed page.

### Prompt Context

**User prompt (verbatim):** "continue" and reminder to keep a diary and commit at appropriate intervals.

**Assistant interpretation:** Continue the `HOST-002` page work, keep the implementation diary current, validate the frontend, and commit the focused slice.

### What I did

- Added `CreateSiteRequest` to dashboard service types.
- Added RTK Query `createSite` mutation for `POST /api/v1/orgs/{org_id}/sites`.
- Invalidated the org site list cache after successful create.
- Added a shared `apiErrorMessage()` helper for RTK Query errors.
- Added MSW `POST /api/v1/orgs/:orgId/sites` handler.
- Reworked `CreateSitePage` into a real form:
  - slug and name state,
  - DNS-safe lowercase slug validation,
  - name presence/length validation,
  - preview host using `/api/v1/config` base domain,
  - API error rendering,
  - success navigation to the future site detail route.
- Added Storybook interaction stories for:
  - invalid slug,
  - successful create,
  - forbidden create.
- Added `@storybook/test` as a dev dependency so interaction stories can use `userEvent`, `within`, and `expect` with type support.

### Why

The dashboard needed a complete create-site loop before moving to site detail/deployment pages. This validates that form state, client validation, RTK Query mutations, MSW, and Storybook interactions all work together.

### What worked

- `make web-build` passes.
- `make storybook-build` passes.
- Storybook restarted successfully through `devctl`.

### What didn't work

- Initial Storybook interaction imports failed because `@storybook/test` was not installed. I added it to `devDependencies` and restored the documented import path.
- The first `apiErrorMessage` implementation assumed every `FetchBaseQueryError` had an `error` property; TypeScript correctly flagged that only some variants do. I fixed this with discriminated `in` checks.

### What was tricky

- The create success story needs a placeholder route for `/app/orgs/:orgId/sites/:siteId` so navigation can be asserted inside Storybook without implementing the full site detail page yet.

### What warrants review

- The slug regex is intentionally conservative and DNS-ish. If backend slug policy differs, align both sides in a later pass.
- Success navigation currently points at a not-yet-real site detail route; that route should be implemented next.

### Validation

Commands run:

```bash
make web-build
make storybook-build
devctl restart storybook
```

## Step 9: Completed practical Phase 6-8 dashboard workflow for real testing

The user asked to push through phases 6, 7, and 8 so there is a real dashboard workflow to test against the dev daemon. I implemented the missing organization/site creation pieces, site overview/runtime pages, and the deployment upload/list/detail/activate/rollback loop.

### Prompt Context

**User prompt (verbatim):** "do phase 6 7 8 because then i'll have something to test for real."

**Assistant interpretation:** Prioritize end-to-end usability over polishing every subcomponent. The important path is: create org, create site, inspect runtime/deployments, upload a bundle, activate, and rollback.

### What I did

Phase 6:

- Added `createOrg` RTK Query mutation for `POST /api/v1/orgs`.
- Reworked `NoOrgsPage` into a real create-organization form with validation, API error rendering, and success navigation into create-site.
- Added MSW `POST /api/v1/orgs` handler.
- Replaced fixture-only site runtime badges with runtime fan-out queries from the real `GET /api/v1/sites/{site_id}/runtime` endpoint.

Phase 7:

- Added `SiteHeader` and `SiteTabs` organisms.
- Added `SiteLayout` route wrapper.
- Added `SiteOverviewPage` that composes runtime status, deployment timeline, active deployment validation report, host copy affordance, and raw site DTO debug output.
- Added `RuntimePage` with runtime refresh and 10-second polling.
- Wired nested routes under `/app/orgs/:orgId/sites/:siteId`.

Phase 8:

- Added `getDeployment`, `uploadDeployment`, `activateDeployment`, and `rollbackDeployment` RTK Query endpoints.
- Implemented multipart bundle upload with `FormData`.
- Treated validation-rejected upload responses that include `{ deployment, report, manifest }` as displayable results even when the HTTP status is `400`.
- Added `DeploymentUploadPanel` organism with manifest and validation output.
- Added `DeploymentsPage` with upload panel, deployment timeline, and rollback confirmation.
- Added `DeploymentDetailPage` with deployment status, validation report, and activation confirmation.
- Added safe JSON parsing helpers for deployment `manifestJson` and `validationJson`.
- Expanded MSW handlers for deployment upload, get, activate, and rollback.
- Added Storybook stories for site overview, runtime page, upload states, deployment list states, and deployment detail states.

### Why

This creates a coherent browser workflow that can be exercised against the real daemon instead of only isolated components:

1. Open `/app`.
2. Create an organization if none exists.
3. Create a site.
4. Open the site overview.
5. Inspect runtime/deployments.
6. Upload a bundle.
7. Inspect validation output.
8. Activate or roll back deployments.

### What worked

- `make web-build` passes.
- `make storybook-build` passes.
- The Storybook/MSW workflow can exercise these pages without a backend.

### What didn't work

- N/A in this slice.

### What was tricky

- Upload validation failures are intentionally not treated as fatal UI errors when the backend returns a deployment plus validation report. This lets rejected bundles appear as normal reviewable deployment records.
- The site list runtime badge fan-out uses per-site child query components so the parent can keep `SitesTable` simple while still getting live status.

### What warrants review

- Activate and rollback currently use `window.confirm`; this satisfies confirmation gating, but should be replaced by the planned themed `ConfirmActionDialog` molecule/organism.
- Upload progress is currently represented by `isLoading`, not byte-level progress.
- The runtime panel shows deployment ID but does not yet render it as a direct link to deployment detail.

### Validation

Commands run:

```bash
make web-build
make storybook-build
```

## Step 10: Updated devctl to run the full dashboard test stack

The user asked whether devctl runs the dashboard and what to test. I updated the devctl plugin so `devctl up` now starts the complete local stack needed for real browser testing.

### Prompt Context

**User prompt (verbatim):** "ok, can you update devctl to run it (if you haven't already), and tell me what I should test."

**Assistant interpretation:** Make devctl start Postgres, the Go daemon/API, the Vite dashboard, and Storybook together, then provide a concrete manual test checklist.

### What I did

- Updated `plugins/go-go-host-devctl.py` launch plan to include four services:
  - `postgres` using `docker compose -f deployments/dev/docker-compose.yaml up postgres`,
  - `go-go-hostd` using `go run ./cmd/go-go-hostd --config configs/dev.yaml`,
  - `web-admin` using `pnpm dev`,
  - `storybook` using `pnpm storybook`.
- Added service config metadata for:
  - Postgres on `127.0.0.1:55432`,
  - daemon/API on `127.0.0.1:8080`,
  - dashboard on `127.0.0.1:5173`,
  - Storybook on `127.0.0.1:6007`.
- Made the daemon service wait for Postgres TCP readiness before starting.
- Added Docker validation to devctl plugin validation.
- Restarted the dev environment with `devctl up --force --skip-validate`.

### What worked

- `devctl plan` shows all four services.
- `devctl up --force --skip-validate` started all four services.
- `devctl status --tail-lines 20` reports all four services alive.
- `curl http://127.0.0.1:8080/healthz` returns `{"status":"ok"}`.
- `curl http://127.0.0.1:8080/api/v1/me` works with dev auth.
- `curl -I http://127.0.0.1:5173/app` returns the Vite app HTML.

### What didn't work

- I used `--skip-validate` when restarting because Docker availability and current environment were already known; normal future `devctl up --force` can run validation.

### What warrants review

- The Postgres devctl service attaches to `docker compose up postgres`; stopping devctl stops the attached process, while Docker may keep or recreate the container depending on Compose signal behavior. If this is annoying, a later refinement can add explicit prepare/cleanup ops.

### Validation

Commands run:

```bash
devctl plan
devctl up --force --skip-validate
devctl status --tail-lines 20
curl -fsS http://127.0.0.1:8080/healthz
curl -fsS -H 'X-Go-Go-Host-User: alice' http://127.0.0.1:8080/api/v1/me
curl -fsSI http://127.0.0.1:5173/app
```

## Step 11: Fixed tar bundles created with `tar -C dir -czf bundle .`

The user uploaded `/tmp/go-go-host-test-bundle.tar.gz` and the dashboard showed a validation report with `missing go-go-host.json manifest`, even though the bundle was created from a directory containing the manifest.

### Prompt Context

**User prompt (verbatim):** "I uploaded /tmp/go-go-host-test-bundle.tar.gz and I got: { \"bytes\": 708, \"files\": 3, \"valid\": false, \"errors\": [ \"missing go-go-host.json manifest\" ] }"

**Assistant interpretation:** Investigate the deployment validator because the provided bundle recipe should have produced a valid manifest.

### Root cause

The suggested bundle command was:

```bash
tar -C /tmp/go-go-host-test-bundle -czf /tmp/go-go-host-test-bundle.tar.gz .
```

That creates tar entries such as `./go-go-host.json` and `./scripts/app.js`. The validator accepted the paths after cleaning them, but it looked for the manifest using the raw archive name exactly equal to `go-go-host.json`. Therefore `./go-go-host.json` was counted as a file but not recognized as the manifest.

### Fix

- Added `normalizeArchiveName()` in `internal/deploy/bundle.go`.
- `readTarGz()` and `readZip()` now normalize file names with `filepath.Clean`/slash conversion when building archive entries.
- Added `TestValidateAndStoreAcceptsDotSlashManifest` to prove `./go-go-host.json` is accepted.

### Validation

Commands run:

```bash
go test ./internal/deploy ./internal/control ./internal/httpapi
make web-build
devctl restart go-go-hostd
devctl status --tail-lines 8
```

The daemon has been restarted, so the same uploaded bundle should now validate correctly.

## Step 12: Added Phase 9 agents and audit pages

After confirming the upload fix worked, I continued implementation work by exposing the initial agent and audit APIs in the dashboard.

### Prompt Context

**User prompt (verbatim):** "ok, it works. Continue with implementation work"

**Assistant interpretation:** Continue after the fixed real deployment workflow, moving to the next dashboard phase: agent and audit pages.

### What I did

- Added service types for agent create/revoke requests.
- Added RTK Query endpoints:
  - `createAgent` -> `POST /api/v1/orgs/{org_id}/agents`,
  - `revokeAgent` -> `POST /api/v1/orgs/{org_id}/agents/{agent_id}/revoke`,
  - expanded audit filters for action, actor type, actor ID, and resource ID.
- Added MSW handlers for agent create, agent revoke, and filtered audit list.
- Added `AgentsPage`:
  - preview notice that enrollment keys/grants/deploy-run tokens are pending,
  - create-agent form,
  - agents table,
  - confirmation-gated revoke,
  - empty/error states.
- Added `AuditPage`:
  - action and actor-type filters,
  - filters preserved in URL query params,
  - audit timeline with selected metadata display,
  - empty/error states.
- Wired `/app/orgs/:orgId/agents` and `/app/orgs/:orgId/audit` to the real pages.
- Added Storybook page stories for populated, empty/error-ish, create success/error, selected metadata, and URL-filtered states.

### Why

The backend already exposes initial agent and audit APIs. Adding the pages makes the dashboard useful for inspecting who/what changed deployments and for creating/revoking initial automation identities without pretending that the future enrollment/grants flow is complete.

### What worked

- `make web-build` passes.
- `make storybook-build` passes.

### What didn't work

- N/A.

### What was tricky

- Audit filters use backend snake_case query names (`actor_type`, `actor_id`, `resource_id`) while the frontend type uses camelCase. The RTK Query endpoint maps the currently used filters through request params; the page preserves visible filters in URL query params.

### What warrants review

- Revoke uses `window.confirm`; this should later become the themed confirmation dialog.
- `AgentsTable` shows revoke actions, but real enrollment key/grant management still needs backend APIs before we can go deeper.

### Validation

Commands run:

```bash
make web-build
make storybook-build
```

## Step 13: Implemented Phase 10 usage/members placeholders and Phase 11 embedded SPA serving

The user confirmed the deployment workflow works and asked to continue with phases 10 and 11. I added the remaining product-placeholder pages and then wired the built React dashboard into the Go daemon under `/app`.

### Prompt Context

**User prompt (verbatim):** "ok, works, continue. do phase 10 and 11"

**Assistant interpretation:** Finish the usage/members surfaces and make the dashboard available from `go-go-hostd` itself, not only from the Vite dev server.

### What I did

Phase 10:

- Added `QuotaPanel` organism with API-pending quota explanation and runtime-derived counters.
- Added `MembersTable` organism based on `/api/v1/me` memberships and roles.
- Added `UsagePage`:
  - lists sites,
  - fans out runtime status requests,
  - aggregates request/error counters,
  - shows per-site runtime counters.
- Added `MembersPage`:
  - shows membership roles from `/api/v1/me`,
  - explicitly notes that membership mutation APIs are pending.
- Wired org sidebar active state for usage/members.
- Wired routes:
  - `/app/orgs/:orgId/usage`,
  - `/app/orgs/:orgId/members`.
- Added Storybook stories for quota counters, empty usage, current memberships, role variants, and empty memberships.

Phase 11:

- Set Vite production `base` to `/app/` so embedded assets are requested as `/app/assets/...`.
- Added `make web-embed`, which builds `web/admin` and copies `web/admin/dist` into `internal/webadmin/dist`.
- Replaced the placeholder webadmin handler with an embedded `go:embed` SPA handler.
- The handler serves assets directly and falls back nested app routes to `index.html`.
- Updated HTTP routing so `/app`, `/app/`, and nested `/app/...` routes are served by the embedded dashboard while `/api/*` continues to route to JSON APIs.
- Kept `/admin` separate by routing it through the same compatibility handler for now; it can be split to a separate admin bundle later.
- Added tests:
  - `internal/webadmin/handler_test.go`,
  - `internal/httpapi/webadmin_integration_test.go`.
- Added `.gitignore` exceptions so `internal/webadmin/dist` is committed while generated root/web dist directories remain ignored.

### What worked

- `make web-build` passes.
- `make storybook-build` passes.
- `make web-embed` succeeds.
- `go test ./...` passes.
- `curl -I http://127.0.0.1:8080/app` returns embedded dashboard HTML.
- `curl -I http://127.0.0.1:8080/app/orgs/org_123/sites` returns embedded dashboard HTML.
- Devctl services were restarted successfully.

### What didn't work

- The first webadmin test saw HTTP 301 redirects because `http.FileServer` redirects `/index.html` to `./`. I fixed the handler to serve the embedded index bytes directly for SPA fallback routes.

### What warrants review

- `/admin` is still not a distinct platform admin bundle; it currently uses the compatibility webadmin handler. The route is separate, but the product surface remains future work.
- `internal/webadmin/dist` is committed as the production embedded snapshot. Whenever frontend changes should ship in `go-go-hostd`, run `make web-embed` before committing.

### Validation

Commands run:

```bash
make web-build
make storybook-build
go test ./internal/webadmin ./internal/httpapi ./internal/deploy
make web-embed
go test ./...
devctl restart go-go-hostd
devctl restart web-admin
devctl restart storybook
curl -fsSI http://127.0.0.1:8080/app
curl -fsSI http://127.0.0.1:8080/app/orgs/org_123/sites
```
