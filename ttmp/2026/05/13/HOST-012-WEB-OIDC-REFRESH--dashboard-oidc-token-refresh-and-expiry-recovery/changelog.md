# Changelog

## 2026-05-13

- Initial workspace created


## 2026-05-13

Step 1: created ticket, implementation guide, task list, and diary for dashboard OIDC refresh

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/ttmp/2026/05/13/HOST-012-WEB-OIDC-REFRESH--dashboard-oidc-token-refresh-and-expiry-recovery/design-doc/01-dashboard-oidc-refresh-issue-and-implementation-guide.md — Issue explanation and implementation plan
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/auth/oidc.ts — Current browser token storage and target refresh helper location
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/services/goGoHostApi.ts — Current API bearer attachment and target retry wrapper


## 2026-05-13

Step 2: implemented dashboard OIDC refresh, RTK Query 401 retry, upload retry, and focused Vitest coverage

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/auth/oidc.test.ts — Refresh helper tests
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/auth/oidc.ts — Token metadata
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/services/goGoHostApi.test.ts — RTK Query refresh and retry tests
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/services/goGoHostApi.ts — Async auth base query


## 2026-05-13

Step 3: rebuilt embedded dashboard assets so OIDC refresh ships in the Go binary

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/internal/webadmin/dist/assets/index-BHlgXeht.js — Generated dashboard bundle containing OIDC refresh implementation
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/internal/webadmin/dist/index.html — Embedded dashboard entrypoint updated to new JS bundle


## 2026-05-13

Step 4: deployed HOST-012 dashboard refresh image sha-6c833cb to beta and verified production serves index-BHlgXeht.js

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/ttmp/2026/05/13/HOST-012-WEB-OIDC-REFRESH--dashboard-oidc-token-refresh-and-expiry-recovery/reference/01-diary.md — Deployment diary and validation notes

