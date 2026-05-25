# Changelog

## 2026-05-11

- Initial workspace created


## 2026-05-11

Created admin dashboard ticket, design guide, diary, and 10-phase Storybook-first implementation plan.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/ttmp/2026/05/11/HOST-003-ADMIN-DASHBOARD--go-go-host-platform-admin-dashboard/design-doc/01-platform-admin-dashboard-design-and-implementation-guide.md — Admin dashboard design guide
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/ttmp/2026/05/11/HOST-003-ADMIN-DASHBOARD--go-go-host-platform-admin-dashboard/tasks.md — Detailed phased admin dashboard task plan


## 2026-05-11

Implemented the first admin dashboard MVP: platform-admin guard, /admin routes, admin shell/sidebar, runtime summary RTK Query endpoint, overview/runtimes pages, Storybook stories, MSW fixtures, and embedded /admin route coverage.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/internal/httpapi/webadmin_integration_test.go — Embedded /admin SPA route coverage
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/app/routes.tsx — Admin route tree
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/components/organisms/AdminRuntimeTable/AdminRuntimeTable.tsx — Admin runtime inventory table
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/pages/AdminOverviewPage/AdminOverviewPage.tsx — Admin overview MVP


## 2026-05-11

Added platform-admin inventory APIs and initial admin org/user/site/deployment pages with RTK Query, MSW fixtures, and Storybook coverage.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/internal/httpapi/admin_inventory.go — platform-admin inventory handlers
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/internal/store/queries/admin.sql — sqlc admin inventory queries
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/pages/AdminDeploymentsPage/AdminDeploymentsPage.tsx — admin deployment inventory page
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/pages/AdminSitesPage/AdminSitesPage.tsx — admin site inventory page


## 2026-05-11

Added dev platform-admin subject seeding, backend admin inventory access tests, and verified embedded /admin pages as a seeded dev platform admin.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/internal/config/config.go — devPlatformAdminSubjects config
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/internal/httpapi/admin_inventory_integration_test.go — admin inventory authorization coverage
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/internal/httpapi/auth.go — dev auth auto-seeds configured platform admins


## 2026-05-11

Added global admin agents and audit endpoints/pages, including sqlc queries, RTK Query hooks, MSW fixtures, Storybook stories, embedded build, and Playwright verification.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/internal/httpapi/admin_inventory.go — admin agents and audit handlers
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/internal/store/queries/admin.sql — global agents and audit queries
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/pages/AdminAgentsPage/AdminAgentsPage.tsx — global agents page
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/pages/AdminAuditPage/AdminAuditPage.tsx — global audit page


## 2026-05-11

Added admin deployment detail endpoint and page with manifest, validation, actor, bundle, org/site metadata, Storybook stories, embedded build, and Playwright verification.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/internal/httpapi/admin_inventory.go — admin deployment detail handler
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/internal/store/queries/admin.sql — GetAdminDeployment query
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/pages/AdminDeploymentDetailPage/AdminDeploymentDetailPage.tsx — admin deployment detail page


## 2026-05-11

Added platform-admin runtime restart/stop endpoints with audit logging, themed confirmation dialog, admin runtime action buttons, RTK Query mutations, MSW handlers, and fixed API 404 fallback flushing.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/internal/httpapi/fallback.go — API 404 fallback flush fix
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/internal/httpapi/runtime.go — admin runtime operation endpoints
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/components/molecules/ConfirmActionDialog/ConfirmActionDialog.tsx — themed confirmation dialog
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/pages/AdminRuntimesPage/AdminRuntimesPage.tsx — runtime operation UI wiring


## 2026-05-11

Completed admin dashboard phases 7-9 with read-only quota/capability/domain APIs and pages, Storybook coverage, favicon, responsive table overflow, and dialog keyboard polish.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/internal/httpapi/admin_inventory.go — admin policy endpoint handlers
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/internal/store/queries/admin.sql — admin quota/capability/domain queries
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/pages/AdminCapabilitiesPage/AdminCapabilitiesPage.tsx — capability policy page
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/pages/AdminDomainsPage/AdminDomainsPage.tsx — domain policy page
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/pages/AdminQuotasPage/AdminQuotasPage.tsx — read-only quota page

