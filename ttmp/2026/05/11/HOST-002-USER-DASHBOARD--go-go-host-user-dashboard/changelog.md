# Changelog

## 2026-05-11

- Initial workspace created


## 2026-05-11

Created dedicated user dashboard ticket, moved Phase 7 design guide into it, replaced old copy with pointer, and added detailed phased task plan.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/ttmp/2026/05/11/HOST-002-USER-DASHBOARD--go-go-host-user-dashboard/tasks.md — Phased dashboard implementation plan


## 2026-05-11

Scaffolded web/admin with Vite React TypeScript, RTK Query, MSW Storybook, os-core shim, Makefile targets, devctl Storybook/Vite services, and first StatusPill/AppBootstrap stories.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/package.json — Dashboard frontend scaffold


## 2026-05-11

Fixed Storybook MSW worker registration and added initial atom components with stories: RuntimeStatusDot, RoleBadge, CopyButton, EmptyState, ErrorCallout, LoadingBlock, Timestamp, CodeBlock, and JsonTree.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/public/mockServiceWorker.js — MSW worker served by Storybook
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/components/atoms/index.ts — Atom component barrel exports


## 2026-05-11

Switched dashboard from temporary os-core shim to local linked go-go-os-core macOS1 theme imports and documented the macos1 bridge as a temporary class-to-token adapter.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/app/macos1-bridge.css — Temporary bridge from dashboard scaffold classes to os-core macOS1 tokens


## 2026-05-11

Added first dashboard molecule set with Storybook stories: RuntimeBadge, SiteHostCopy, DeploymentStatusPill, ValidationSummary, MetricCard, ManifestSummary, AgentStatusBadge, and OrgSwitcher.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/components/molecules/index.ts — Dashboard molecule barrel exports


## 2026-05-11

Added first dashboard organism set with Storybook stories: AppShell, OrgSidebar, SitesTable, RuntimeStatusPanel, DeploymentTimeline, ValidationReportPanel, AgentsTable, and AuditTimeline.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/components/organisms/index.ts — Dashboard organism barrel exports


## 2026-05-11

Added initial /app routing with session/org guards plus NoOrgsPage, SitesPage, and CreateSitePage Storybook-backed page shells.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/app/routes.tsx — Initial dashboard route tree


## 2026-05-11

Wired CreateSitePage to RTK Query create-site mutation with validation, MSW POST handler, and Storybook interaction stories for invalid, success, and forbidden states.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/pages/CreateSitePage/CreateSitePage.tsx — Create-site page form and mutation wiring


## 2026-05-11

Completed practical Phase 6-8 dashboard workflow: org/site create, site layout/runtime overview, deployment upload/list/detail/activate/rollback pages, MSW handlers, and Storybook states.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/app/routes.tsx — Nested site and deployment route tree


## 2026-05-11

Updated devctl to run the full local dashboard test stack: Postgres, go-go-hostd, Vite web-admin, and Storybook.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/plugins/go-go-host-devctl.py — Full dashboard devctl launch plan


## 2026-05-11

Fixed deployment bundle validation for tar archives whose entries are prefixed with ./, allowing tar -C bundle-dir -czf bundle.tar.gz . to find go-go-host.json.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/internal/deploy/bundle.go — Normalize archive entry names before manifest detection


## 2026-05-11

Added Phase 9 dashboard pages for agents and audit, including create/revoke agent mutations, audit filters in URL query params, MSW handlers, and Storybook states.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/pages/AgentsPage/AgentsPage.tsx — Agents page API wiring


## 2026-05-11

Implemented Phase 10 usage/members pages and Phase 11 embedded dashboard serving from go-go-hostd under /app with build copy flow and route tests.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/internal/webadmin/handler.go — Embedded SPA handler for /app


## 2026-05-11

Replaced manual dashboard embedding with a Dagger-backed cmd/build-web pipeline, pinned pnpm, rebuilt embedded assets, and verified /app CSS assets are served by go-go-hostd.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/cmd/build-web/main.go — Dagger build-web pipeline for embedded dashboard


## 2026-05-11

Fixed embedded dashboard styling mismatch by making macOS1 bridge CSS self-contained, removing unstable os-core side-effect imports, simplifying Dagger build-web, and validating embedded/Storybook pages with Playwright.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/app/macos1-bridge.css — Self-contained macOS1 token fallback


## 2026-05-11

Restored @go-go-golems/os-core as a published npmjs dependency, rebuilt embedded dashboard with Dagger, and verified embedded styling with Playwright.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/package.json — Uses @go-go-golems/os-core from npmjs

