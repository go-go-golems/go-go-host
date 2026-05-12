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

