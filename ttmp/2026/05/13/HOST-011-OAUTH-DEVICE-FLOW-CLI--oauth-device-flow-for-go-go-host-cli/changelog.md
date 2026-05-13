# Changelog

## 2026-05-13

- Initial workspace created


## 2026-05-13

Step 1: created ticket and saved OAuth Device Flow, Keycloak, Terraform, and live production source material under sources/

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/ttmp/2026/05/13/HOST-011-OAUTH-DEVICE-FLOW-CLI--oauth-device-flow-for-go-go-host-cli/reference/01-investigation-diary.md — Diary records ticket setup and source capture
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/ttmp/2026/05/13/HOST-011-OAUTH-DEVICE-FLOW-CLI--oauth-device-flow-for-go-go-host-cli/sources/00-sources-readme.md — Index of captured sources for HOST-011


## 2026-05-13

Step 2: mapped CLI/auth architecture and wrote the OAuth Device Flow intern implementation guide

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/ — Primary intern-facing analysis
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/ttmp/2026/05/13/HOST-011-OAUTH-DEVICE-FLOW-CLI--oauth-device-flow-for-go-go-host-cli/reference/01-investigation-diary.md — Diary records architecture mapping and guide writing


## 2026-05-13

Step 3: validated HOST-011 docs and uploaded the OAuth Device Flow CLI guide bundle to reMarkable

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/ttmp/2026/05/13/HOST-011-OAUTH-DEVICE-FLOW-CLI--oauth-device-flow-for-go-go-host-cli/reference/01-investigation-diary.md — Diary records validation
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/ttmp/2026/05/13/HOST-011-OAUTH-DEVICE-FLOW-CLI--oauth-device-flow-for-go-go-host-cli/sources/00-sources-readme.md — Validated source index included in uploaded bundle


## 2026-05-13

Step 4: implemented backend OIDC accepted-client support and device client config discovery

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/internal/config/config.go — Added OIDC device client and accepted client IDs
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/internal/httpapi/handler.go — Publishes deviceClientId in /api/v1/config
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/internal/httpapi/oidc.go — Accepts tokens for any configured OIDC client by aud/azp
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/services/types.ts — Frontend type updated for deviceClientId

