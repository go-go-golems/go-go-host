# Changelog

## 2026-05-12

- Initial workspace created


## 2026-05-12

Created production readiness guide with prioritized beta launch backlog, evidence-backed architecture analysis, API sketches, pseudocode, implementation phases, tests, risks, and file references.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/ttmp/2026/05/12/HOST-006-PRODUCTION-READINESS--go-go-host-production-readiness-and-beta-launch-plan/design-doc/01-go-go-host-production-readiness-and-beta-launch-implementation-guide.md — primary deliverable
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/ttmp/2026/05/12/HOST-006-PRODUCTION-READINESS--go-go-host-production-readiness-and-beta-launch-plan/reference/01-investigation-diary.md — chronological investigation record
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/ttmp/2026/05/12/HOST-006-PRODUCTION-READINESS--go-go-host-production-readiness-and-beta-launch-plan/tasks.md — prioritized implementation backlog


## 2026-05-12

Implemented Phase 1 local Keycloak/OIDC foundations: realm import, dev Keycloak config, devctl Keycloak startup, OIDC config API, frontend PKCE login/callback/logout, bearer-token API calls, OIDC platform-admin bootstrap, gated Playwright smoke, and embedded dashboard rebuild.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/configs/dev.keycloak.yaml — local OIDC daemon config
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/deployments/dev/keycloak/realm-go-go-host.json — local Keycloak realm/client/users/role import
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/internal/httpapi/oidc.go — OIDC platform-admin bootstrap from subject/email/roles
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/plugins/go-go-host-devctl.py — devctl starts Keycloak and uses OIDC config
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/scripts/oidc-login-playwright.mjs — gated OIDC browser smoke
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/auth/oidc.ts — dashboard PKCE/token/logout helper


## 2026-05-12

Verified Phase 1 OIDC browser flow live with built-in Playwright tooling: Keycloak redirect, platform-admin login, admin dashboard access, logout, Alice login, and non-admin `/admin` denial.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/ttmp/2026/05/12/HOST-006-PRODUCTION-READINESS--go-go-host-production-readiness-and-beta-launch-plan/reference/01-investigation-diary.md — live Playwright verification notes
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/ttmp/2026/05/12/HOST-006-PRODUCTION-READINESS--go-go-host-production-readiness-and-beta-launch-plan/tasks.md — Phase 1 live smoke task updates

