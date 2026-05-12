# Changelog

## 2026-05-12

- Initial workspace created


## 2026-05-12

Created HOST-007 intern guide and diary, implemented access-token preference/client matching, added examples/hello-beta, and added scripts/beta-smoke.sh.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/examples/hello-beta — demo fixture
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/internal/httpapi/oidc.go — OIDC bearer-token cleanup
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/scripts/beta-smoke.sh — public beta smoke script
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/ttmp/2026/05/12/HOST-007-BETA-SMOKE-AUTH-CLEANUP--beta-smoke-and-oidc-access-token-cleanup/design-doc/01-beta-smoke-and-oidc-access-token-cleanup-guide.md — intern-facing HOST-007 guide


## 2026-05-12

Deployed access-token image sha-23b66ec, verified live access-token API auth, discovered demo-site 404 after rollout, and added daemon startup restoration for active deployments.

### Related Files

- /home/manuel/code/wesen/2026-03-27--hetzner-k3s/gitops/kustomize/go-go-host/deployment.yaml — deployed sha-23b66ec during HOST-007
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/cmd/go-go-hostd/main.go — daemon startup calls active runtime restoration
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/internal/control/deployments.go — startup restoration for active deployments


## 2026-05-12

Built and deployed sha-f137ff9 with startup active-runtime restoration; Argo is healthy, live access-token auth works, and scripts/beta-smoke.sh passes after rollout.

### Related Files

- /home/manuel/code/wesen/2026-03-27--hetzner-k3s/gitops/kustomize/go-go-host/deployment.yaml — live image now sha-f137ff9
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/internal/control/deployments.go — deployed startup active runtime restoration
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/scripts/beta-smoke.sh — post-rollout beta smoke passed


## 2026-05-12

Linked the demo CSS from the root page, redeployed the demo, tested live signed agent publishing, documented the bundles/** grant-path rejection, and revoked temporary smoke agents.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/examples/hello-beta/scripts/app.js — demo page now links stylesheet
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/ttmp/2026/05/12/HOST-007-BETA-SMOKE-AUTH-CLEANUP--beta-smoke-and-oidc-access-token-cleanup/reference/01-investigation-diary.md — agent publishing smoke details


## 2026-05-12

Added P0 tasks to rename/fix agent path semantics around allowedBundlePaths, bundlePath, and --bundle-path so grants authorize logical bundle paths rather than archive entries.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/ttmp/2026/05/12/HOST-007-BETA-SMOKE-AUTH-CLEANUP--beta-smoke-and-oidc-access-token-cleanup/reference/01-investigation-diary.md — bundle path naming decision
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/ttmp/2026/05/12/HOST-007-BETA-SMOKE-AUTH-CLEANUP--beta-smoke-and-oidc-access-token-cleanup/tasks.md — agent bundle-path semantics tasks


## 2026-05-12

Implemented preferred agent bundle-path semantics locally: allowedBundlePaths/bundlePath API aliases, --bundle-path CLI flags, and removal of agent grant path checks from archive-entry validation.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/cmd/go-go-host-agent/cmds/deploy.go — agent CLI --bundle-path flag
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/cmd/go-go-host/cmds/agents.go — operator CLI --bundle-path grant flag
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/internal/httpapi/agents_audit.go — allowedBundlePaths and bundlePath API aliases
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/internal/httpapi/deployments.go — agent uploads no longer pass logical bundle paths to archive validator

