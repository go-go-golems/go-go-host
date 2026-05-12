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

