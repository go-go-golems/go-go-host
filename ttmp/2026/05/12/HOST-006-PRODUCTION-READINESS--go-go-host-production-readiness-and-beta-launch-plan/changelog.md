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


## 2026-05-12

Made OIDC smoke repeatable with a Playwright dev dependency, GO_GO_HOST_OIDC_E2E=1 node scripts/oidc-login-playwright.mjs, default embedded-dashboard target, and longer devctl health windows for cold Keycloak startup.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/Makefile — oidc-e2e target
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/plugins/go-go-host-devctl.py — longer health windows for Keycloak and dashboard services
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/scripts/oidc-login-playwright.mjs — repeatable OIDC browser smoke
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/package.json — Playwright dev dependency


## 2026-05-12

Closed Phase 1 CI-smoke follow-up by deciding local GO_GO_HOST_OIDC_E2E=1 node scripts/oidc-login-playwright.mjs
OIDC E2E ok: admin@example.test platformAdmin=true testing is sufficient for now.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/ttmp/2026/05/12/HOST-006-PRODUCTION-READINESS--go-go-host-production-readiness-and-beta-launch-plan/reference/01-investigation-diary.md — CI decision note
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/ttmp/2026/05/12/HOST-006-PRODUCTION-READINESS--go-go-host-production-readiness-and-beta-launch-plan/tasks.md — Phase 1 CI decision


## 2026-05-12

Backfilled detailed K3s/Argo/Keycloak/Vault/DNS beta deployment diary, including the PVC sync-wave hang, hostname correction to hosting.yolo.scapegoat.dev, smoke tests, commits, and caveats.

### Related Files

- /home/manuel/code/wesen/2026-03-27--hetzner-k3s/docs/go-go-host-beta-deployment-playbook.md — K3s beta deploy runbook
- /home/manuel/code/wesen/2026-03-27--hetzner-k3s/gitops/kustomize/go-go-host — go-go-host Argo/Kustomize package
- /home/manuel/code/wesen/terraform/dns/zones/scapegoat-dev/envs/prod/main.tf — hosting.yolo DNS wildcard
- /home/manuel/code/wesen/terraform/keycloak/apps/go-go-host/envs/k3s-beta — Keycloak beta realm Terraform
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/ttmp/2026/05/12/HOST-006-PRODUCTION-READINESS--go-go-host-production-readiness-and-beta-launch-plan/reference/01-investigation-diary.md — backfilled beta deployment diary


## 2026-05-12

Added K3s GitOps-managed DigitalOcean DNS-01 ClusterIssuer and go-go-host wildcard Certificate/Ingress so generated site hosts under *.hosting.yolo.scapegoat.dev can terminate TLS and reach the daemon.

### Related Files

- /home/manuel/code/wesen/2026-03-27--hetzner-k3s/gitops/kustomize/go-go-host/certificate.yaml — wildcard TLS Certificate for hosting.yolo.scapegoat.dev
- /home/manuel/code/wesen/2026-03-27--hetzner-k3s/gitops/kustomize/go-go-host/ingress.yaml — wildcard ingress rule for generated site hosts
- /home/manuel/code/wesen/2026-03-27--hetzner-k3s/gitops/kustomize/platform-cert-issuer/clusterissuer-dns01-digitalocean.yaml — DNS-01 ClusterIssuer for DigitalOcean


## 2026-05-12

Created and activated the first live demo site at https://hello.hosting.yolo.scapegoat.dev, validating wildcard DNS/TLS/Ingress, host routing, bundle upload, activation, assets, SQLite state, and platform context.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/ttmp/2026/05/12/HOST-006-PRODUCTION-READINESS--go-go-host-production-readiness-and-beta-launch-plan/reference/01-investigation-diary.md — live demo site smoke results

