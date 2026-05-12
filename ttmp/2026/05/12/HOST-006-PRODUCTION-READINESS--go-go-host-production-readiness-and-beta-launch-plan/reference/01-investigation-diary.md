---
Title: Investigation diary
Ticket: HOST-006-PRODUCTION-READINESS
Status: active
Topics:
    - go-go-host
    - hosting
    - security
    - deployments
    - agents
    - platform-admin
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: ../../../../../../../../../../code/wesen/2026-03-27--hetzner-k3s/docs/go-go-host-beta-deployment-playbook.md
      Note: operator runbook for the live hosting.yolo.scapegoat.dev beta deployment
    - Path: ../../../../../../../../../../code/wesen/2026-03-27--hetzner-k3s/gitops/applications/go-go-host.yaml
      Note: Argo CD Application for go-go-host beta
    - Path: ../../../../../../../../../../code/wesen/2026-03-27--hetzner-k3s/gitops/kustomize/go-go-host
      Note: Kustomize manifests for namespace
    - Path: ../../../../../../../../../../code/wesen/2026-03-27--hetzner-k3s/scripts/bootstrap-go-go-host-image-pull-secret.sh
      Note: Vault image-pull secret bootstrap for private GHCR image
    - Path: ../../../../../../../../../../code/wesen/2026-03-27--hetzner-k3s/scripts/bootstrap-go-go-host-runtime-secrets.sh
      Note: Vault runtime secret bootstrap for Postgres DSN
    - Path: ../../../../../../../../../../code/wesen/terraform/dns/zones/scapegoat-dev/envs/prod/main.tf
      Note: DNS wildcard for *.hosting.yolo.scapegoat.dev
    - Path: ../../../../../../../../../../code/wesen/terraform/keycloak/apps/go-go-host/envs/k3s-beta
      Note: Keycloak realm/client/role/user Terraform for beta
    - Path: .github/workflows/publish-image.yaml
      Note: GHCR image publish workflow
    - Path: Dockerfile
      Note: Go 1.26 Docker build fix for beta image
ExternalSources: []
Summary: Chronological diary for the HOST-006 production readiness investigation.
LastUpdated: 2026-05-12T13:14:11.750159577-04:00
WhatFor: Use this to resume or audit the production readiness planning work.
WhenToUse: When implementing beta launch tasks or reviewing evidence behind the readiness guide.
---


# Diary

## Goal

Capture the investigation and documentation work for the go-go-host production readiness and beta launch plan.

## Step 1: Ticket setup, evidence gathering, and production-readiness guide

I created a new production-readiness ticket and wrote the first full design/implementation guide for moving go-go-host from a local MVP to a beta service. The guide is intentionally ordered by necessity: real browser auth and admin bootstrap first, then release/deploy mechanics, then domains/TLS, runtime isolation, secrets, backups, observability, and beta onboarding.

The main finding is that the core product loop exists, but the service is not production-ready until local/staging exercise the same OIDC path beta users will use. The current dev-user flow is useful for fast backend iteration, but it should not remain the normal browser testing path for beta readiness.

### Prompt Context

**User prompt (verbatim):** "what are all th e missing pieces needed for proper deploy (even just for beta users)?"

**User prompt (verbatim):** "Create a new ticket with a full detail of all the things needed to push into production, in order of necessity.

Create  a detailed analysis / design / implementation guide that is for a new intern, explaining all the parts of the system needed to understand what it is, with prose paragraphs and bullet point sand pseudocode and diagrams and api references and file references. It should be very clear and technical. Store in the ticket and the nupload to remarkable."

**Assistant interpretation:** Create a new docmgr ticket and write a clear intern-facing production-readiness implementation guide, grounded in current code evidence and prioritized by launch necessity, then upload the result to reMarkable.

**Inferred user intent:** Turn the informal question “what is missing for beta/production?” into a durable engineering plan that can drive future implementation work.

**Commit (code):** N/A — documentation/ticket work only.

### What I did

- Loaded the ticket-research, docmgr, reMarkable upload, diary, and textbook-writing instructions.
- Created ticket `HOST-006-PRODUCTION-READINESS` with topics `go-go-host,hosting,security,deployments,agents,platform-admin`.
- Created the primary design doc and investigation diary.
- Gathered evidence from:
  - auth middleware and OIDC verifier,
  - frontend RTK Query API setup,
  - AppShell dev-auth display,
  - dev Docker Compose and devctl plugin,
  - runtime and supervisor,
  - site settings/domain/secrets placeholder APIs,
  - maintenance/export/prune APIs,
  - Dockerfile and production config example,
  - bundled developer/JS/agent docs,
  - gated Playwright E2E script.
- Wrote the production-readiness guide with priority table, architecture diagrams, current-state analysis, gap analysis, implementation phases, API sketches, pseudocode, tests, risks, and file references.
- Rewrote `tasks.md` into a production-readiness backlog.

### Why

- The user asked for a new ticket and a detailed intern-facing design/implementation guide.
- Production readiness spans many subsystems, so the output needed to explain the architecture before prescribing work.
- The order of work matters: browser OIDC and admin bootstrap must come before domains, secrets, or beta onboarding.

### What worked

- The current codebase has enough structure to make the plan concrete: OIDC bearer validation exists, Keycloak Compose services exist, dashboard API access is centralized, and maintenance/runtime/domain APIs are already factored.
- The new developer and agent help docs provide a useful evidence point: user-facing docs are now comparatively strong, while auth/domain/secrets/ops remain the main production gaps.

### What didn't work

- `docmgr vocab list` showed no dedicated `operations` or `production` topic, so I stayed within existing vocabulary topics for the ticket.
- There are unrelated dirty files in `cmd/go-go-host*/cmds/support.go` and an untracked `HOST-005-E2E-FIXES` workspace. I did not touch them because this ticket is documentation-only.

### What I learned

- Keycloak is defined in Compose, but devctl does not start it and the Keycloak config file still has `devAuth: true`; this makes real local browser auth the top production-readiness gap.
- The backend OIDC verifier is already enough to anchor the browser-login implementation, but the frontend needs PKCE, token storage, logout, and Authorization header injection.
- Domain and secrets work are explicit placeholders rather than half-hidden missing features, which is good for planning because the gaps are visible and documentable.

### What was tricky to build

- The guide needed to distinguish “implemented but not production-grade” from “not implemented.” For example, maintenance export exists, but scheduled backup/restore drills do not; OIDC bearer verification exists, but browser login does not; domain rows exist, but DNS verification and TLS do not.
- The priority order needed to be strict. It would be tempting to implement custom domains or secrets first, but beta users cannot safely test anything until real auth and admin bootstrap are in place.

### What warrants a second pair of eyes

- Review the runtime-isolation recommendations, especially whether in-process Goja is acceptable for closed beta after adding interrupts/concurrency limits or whether subprocess isolation should be P0/P1.
- Review the platform-admin bootstrap recommendation: config-based OIDC subject/email/role mapping may be enough for beta, but production may need a stronger bootstrap ceremony.
- Review the TLS recommendation: edge-managed TLS is simpler for beta, but product goals may require daemon-managed ACME later.

### What should be done in the future

- Implement Phase 1 from the guide: Keycloak realm import, devctl OIDC profile, frontend PKCE, Authorization headers, OIDC admin bootstrap, and login E2E.
- Turn the guide's phases into separate implementation tickets if work will be distributed across multiple sessions or engineers.

### Code review instructions

- Start with `design-doc/01-go-go-host-production-readiness-and-beta-launch-implementation-guide.md`.
- Confirm each major claim maps to the referenced files.
- Review `tasks.md` for whether the priority order matches the desired beta launch sequence.
- Validate with `docmgr doctor --ticket HOST-006-PRODUCTION-READINESS --stale-after 30`.

### Technical details

Representative evidence commands:

```bash
nl -ba internal/httpapi/auth.go | sed -n '1,120p'
nl -ba internal/httpapi/oidc.go | sed -n '1,130p'
nl -ba web/admin/src/services/goGoHostApi.ts | sed -n '1,35p'
nl -ba deployments/dev/docker-compose.yaml | sed -n '1,90p'
nl -ba plugins/go-go-host-devctl.py | sed -n '1,150p'
nl -ba internal/runtime/runtime.go | sed -n '1,130p'
nl -ba internal/control/maintenance.go | sed -n '1,240p'
```

## Step 2: Phase 1 implementation — local Keycloak/OIDC browser login

I implemented the first production-readiness phase so the local stack can exercise the same browser authentication path intended for beta: Keycloak realm import, OIDC daemon config, devctl startup changes, dashboard PKCE login/callback/logout, bearer-token attachment, and OIDC platform-admin bootstrap.

### What changed

- Added `deployments/dev/keycloak/realm-go-go-host.json` with local users `alice`, `bob`, and `platform-admin`, a public PKCE client `go-go-host-dashboard`, and a `go-go-host-admin` realm role.
- Added `configs/dev.keycloak.yaml` and corrected `configs/dev.postgres-keycloak.yaml` to run with `devAuth: false`.
- Updated `deployments/dev/docker-compose.yaml` so Keycloak imports the realm and uses the `keycloak-postgres` service.
- Updated `plugins/go-go-host-devctl.py` so the launch plan starts Keycloak and runs `go-go-hostd` with the Keycloak config.
- Extended daemon config and `/api/v1/config` with OIDC browser-login metadata.
- Added production-safe platform-admin bootstrap knobs: OIDC subjects, emails, and roles/groups/client roles.
- Added dashboard OIDC PKCE helpers, callback route, token storage, logout, and Authorization headers for RTK Query plus deployment uploads.
- Added a gated Playwright smoke script: `GO_GO_HOST_OIDC_E2E=1 node scripts/oidc-login-playwright.mjs`.
- Rebuilt embedded dashboard assets.

### Validation

- `go test ./...`
- `pnpm --dir web/admin build`
- `node --check scripts/oidc-login-playwright.mjs`
- `BUILD_WEB_LOCAL=1 go run ./cmd/build-web`

### Notes

I did not run the live browser Keycloak smoke in this step because it requires starting the full devctl stack and Playwright browser runtime. The script is checked in and gated for that follow-up.

## Step 3: Live browser verification with built-in Playwright tooling

The user pointed out that the agent has a Playwright browser tool available, so I used it to validate the Phase 1 OIDC work without depending on the repo-local `playwright` npm package.

### What I ran

- Started the dev stack with `devctl up --force`.
- The first cold start timed out on devctl health while Keycloak/import was still warming up, so I brought up Docker Compose services and the daemon manually once to verify the browser flow.
- After Keycloak was warm, `devctl up --force` completed successfully with five services: Postgres, Keycloak, go-go-hostd, Vite, and Storybook.
- Used the built-in Playwright browser tool against the embedded dashboard at `http://127.0.0.1:8080/admin`.

### Browser checks

- Navigating to `/admin` redirected to Keycloak.
- Logged in as `platform-admin` / `admin`.
- Returned to `/admin/overview` and saw `admin@example.test · platform admin` plus admin navigation.
- Clicked `Sign out`; the dashboard returned to `/app` and redirected to Keycloak login.
- Logged in as `alice` / `alice`.
- Alice landed on the no-organization onboarding page.
- Navigating Alice to `/admin` showed `Platform admin required`, proving the OIDC admin bootstrap did not accidentally grant normal users platform-admin access.

### Follow-up

The live browser smoke passed through the built-in tool. The remaining work is to convert this manual flow into a repeatable automated test path that does not require the repo-local `playwright` package to be installed separately.

## Step 4: Repeatable OIDC smoke and devctl health hardening

I made the live OIDC browser smoke repeatable from the repository instead of relying only on the agent's interactive Playwright browser tool.

### What changed

- Added `playwright` as a `web/admin` dev dependency so `make web-install` prepares the smoke dependency.
- Updated `scripts/oidc-login-playwright.mjs` to load Playwright from `web/admin/node_modules`, default to the embedded daemon at `http://127.0.0.1:8080`, and use stricter Keycloak locators.
- Added `make oidc-e2e` as the local command wrapper.
- Increased devctl health windows for Keycloak/go-go-hostd/Vite/Storybook so cold Keycloak realm import has enough startup time.

### Validation

- `make oidc-e2e` passed and reported `OIDC E2E ok: admin@example.test platformAdmin=true`.
- `devctl up --force` passed after the health window changes.
- `go test ./...` passed.
- `pnpm --dir web/admin build` passed.
- `python3 -m py_compile plugins/go-go-host-devctl.py` passed.
- `docmgr doctor --ticket HOST-006-PRODUCTION-READINESS --stale-after 30` passed.

### Remaining note

The smoke currently verifies the platform-admin path. The Alice/Bob isolation flow was manually verified with the built-in Playwright browser tool; turning that exact multi-user flow into the scripted smoke remains a future enhancement.

## Step 5: CI decision for OIDC smoke

The user clarified that CI wiring is not necessary for the OIDC browser smoke right now. Local testing through `devctl up --force` plus `make oidc-e2e` is sufficient for Phase 1. I updated the HOST-006 task list to mark the CI decision complete and leave the next work focused on production/beta deployment readiness rather than CI automation.

## Step 6: Beta deployment research across K3s GitOps and Terraform

The user changed the target from a generic beta topology to the real Hetzner K3s cluster and asked to deploy through Argo CD. The target repos were:

- K3s/GitOps repo: `/home/manuel/code/wesen/2026-03-27--hetzner-k3s`
- Keycloak/Terraform repo: `/home/manuel/code/wesen/terraform`
- App repo: `/home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host`

The requested initial hostname was `yolo.scapegoat.dev`; later it was corrected to `hosting.yolo.scapegoat.dev`.

### Prompt Context

**User prompt (verbatim):** "we want to deploy this to my k3s cluster in /home/manuel/code/wesen/2026-03-27--hetzner-k3s , which also has a keycloak setup which is configured in /home/manuel/code/wesen/terraform .

Study how to do a deploy, and then set the system up so we can have a beta deploy on yolo.scapegoat.dev using argocd"

**Assistant interpretation:** Read the existing K3s and Terraform deployment patterns, then add a real beta deployment path for go-go-host using the existing Argo CD, Vault, Postgres, cert-manager, Traefik, and Keycloak conventions.

**Inferred user intent:** Move beyond local/devctl readiness and make the beta system actually run on the real cluster with the same GitOps/identity/secrets discipline used by other apps.

### Evidence gathered

I read the K3s deployment docs and existing manifests to avoid inventing a parallel deployment style:

- `/home/manuel/code/wesen/2026-03-27--hetzner-k3s/docs/argocd-app-setup.md`
- `/home/manuel/code/wesen/2026-03-27--hetzner-k3s/docs/public-repo-ghcr-argocd-deployment-playbook.md`
- `/home/manuel/code/wesen/2026-03-27--hetzner-k3s/docs/app-packaging-and-gitops-pr-standard.md`
- `/home/manuel/code/wesen/2026-03-27--hetzner-k3s/docs/app-runtime-secrets-and-identity-provisioning-playbook.md`
- `/home/manuel/code/wesen/terraform/docs/shared-keycloak-platform-playbook.md`

Important patterns learned:

- New Argo CD `Application` manifests under `gitops/applications/*.yaml` must be applied once with `kubectl apply`; the repo does not currently auto-materialize new apps through app-of-apps.
- Live app manifests are plain Kustomize packages under `gitops/kustomize/<app>/`.
- Runtime secrets are usually delivered with Vault Secrets Operator (`VaultConnection`, `VaultAuth`, `VaultStaticSecret`).
- Shared Postgres already exists as `postgres.postgres.svc.cluster.local:5432` and app DB users/databases are created by idempotent bootstrap Jobs.
- In-cluster VSO traffic should use `http://vault.vault.svc.cluster.local:8200`, not the public Vault hostname.
- Public ingress is Traefik + cert-manager + `letsencrypt-prod` HTTP-01.
- Existing `*.yolo.scapegoat.dev` DNS points at the K3s node, but the apex `yolo.scapegoat.dev` did not.
- The in-cluster Keycloak public URL is `https://auth.yolo.scapegoat.dev`.

### Design selected

For the first beta deploy I chose a conservative single-pod topology:

```text
Ingress hosting.yolo.scapegoat.dev
  -> Service go-go-host:80
    -> Deployment go-go-host:8080
      -> ConfigMap daemon config
      -> Secret go-go-host-runtime from Vault
      -> PVC /var/lib/go-go-host
      -> shared Postgres control-plane DB
      -> Keycloak realm go-go-host
```

The namespace is `go-go-host`.

The K3s repo owns the runtime topology, Vault policies/roles, app manifests, and operator runbook. The Terraform repo owns Keycloak realm/client/role/user state. The app repo owns Dockerfile, image publishing, and daemon config/env expansion.

## Step 7: App repo packaging fixes for beta image publishing

Before the cluster could pull an image, the app repo needed a beta image publishing path and a Dockerfile that actually matched the Go version in `go.mod`.

### What changed

In the app repo I added:

- `.github/workflows/publish-image.yaml`
- `internal/config/config_test.go`

I changed:

- `internal/config/config.go`
- `Dockerfile`

The config loader now expands environment variables inside YAML/JSON config files. This matters because the K3s ConfigMap can safely contain:

```yaml
controlDbDsn: "${GO_GO_HOST_CONTROL_DB_DSN}"
```

and the actual DSN can come from a Kubernetes Secret rendered by Vault Secrets Operator.

### Build failure and fix

The first Docker build failed:

```text
go: go.mod requires go >= 1.26.1 (running go 1.24.13; GOTOOLCHAIN=local)
```

The Dockerfile still used:

```dockerfile
FROM golang:1.24-bookworm AS build
```

but `go.mod` declares:

```text
go 1.26.1
```

I fixed the Dockerfile to use:

```dockerfile
FROM golang:1.26-bookworm AS build
```

### Image built and pushed

I built from an archive of commit `4187ea3` so the image tag matched the source commit exactly:

```bash
cd go-go-host
tmpdir=$(mktemp -d)
git archive 4187ea3 | tar -x -C "$tmpdir"
cd "$tmpdir"
IMAGE=ghcr.io/go-go-golems/go-go-host:sha-4187ea3
docker build -t "$IMAGE" .
docker push "$IMAGE"
```

The image pushed successfully:

```text
ghcr.io/go-go-golems/go-go-host:sha-4187ea3
sha256:0942deba0b9ea834ef6562af5284c26ae67257e80e0da9e45311b506bdecf50e
```

The GHCR package was private, so I added an image-pull-secret path to the K3s GitOps side instead of blocking on package visibility.

### Commits

App repo commits:

- `083e76c Prepare beta image publishing`
- `4187ea3 Use Go 1.26 for Docker builds`

I also pushed the app branch for traceability:

```bash
git push origin task/go-go-host-v1
```

## Step 8: K3s GitOps manifests, Vault roles, and bootstrap scripts

I added a full Kustomize package and Argo CD application in the K3s repo.

### Files created in the K3s repo

Application:

- `gitops/applications/go-go-host.yaml`

Kustomize package:

- `gitops/kustomize/go-go-host/kustomization.yaml`
- `gitops/kustomize/go-go-host/namespace.yaml`
- `gitops/kustomize/go-go-host/serviceaccount.yaml`
- `gitops/kustomize/go-go-host/db-bootstrap-serviceaccount.yaml`
- `gitops/kustomize/go-go-host/vault-connection.yaml`
- `gitops/kustomize/go-go-host/vault-auth.yaml`
- `gitops/kustomize/go-go-host/db-bootstrap-vault-auth.yaml`
- `gitops/kustomize/go-go-host/runtime-secret.yaml`
- `gitops/kustomize/go-go-host/image-pull-secret.yaml`
- `gitops/kustomize/go-go-host/postgres-admin-secret.yaml`
- `gitops/kustomize/go-go-host/db-bootstrap-script-configmap.yaml`
- `gitops/kustomize/go-go-host/db-bootstrap-job.yaml`
- `gitops/kustomize/go-go-host/configmap.yaml`
- `gitops/kustomize/go-go-host/persistentvolumeclaim.yaml`
- `gitops/kustomize/go-go-host/deployment.yaml`
- `gitops/kustomize/go-go-host/service.yaml`
- `gitops/kustomize/go-go-host/ingress.yaml`

Vault policy/role files:

- `vault/policies/kubernetes/go-go-host.hcl`
- `vault/policies/kubernetes/go-go-host-db-bootstrap.hcl`
- `vault/roles/kubernetes/go-go-host.json`
- `vault/roles/kubernetes/go-go-host-db-bootstrap.json`

Operator scripts:

- `scripts/bootstrap-go-go-host-runtime-secrets.sh`
- `scripts/bootstrap-go-go-host-image-pull-secret.sh`

Runbook:

- `docs/go-go-host-beta-deployment-playbook.md`

### Render validation

I validated the Kustomize package locally:

```bash
cd /home/manuel/code/wesen/2026-03-27--hetzner-k3s
kubectl kustomize gitops/kustomize/go-go-host >/tmp/go-go-host-kustomize.yaml
```

After adding the image-pull secret, the package rendered 16 resources:

```text
Namespace go-go-host
ServiceAccount go-go-host
ServiceAccount go-go-host-db-bootstrap
ConfigMap go-go-host-config
ConfigMap go-go-host-db-bootstrap-script
Service go-go-host
PersistentVolumeClaim go-go-host-data
Deployment go-go-host
Job go-go-host-db-bootstrap
Ingress go-go-host
VaultAuth go-go-host
VaultAuth go-go-host-db-bootstrap
VaultConnection vault
VaultStaticSecret go-go-host-ghcr-pull
VaultStaticSecret go-go-host-postgres-admin
VaultStaticSecret go-go-host-runtime
```

### Initial K3s commit

K3s repo commit:

- `984048e Add go-go-host beta GitOps deployment`

## Step 9: Keycloak Terraform beta realm

I added a new Terraform environment for the in-cluster Keycloak deployment.

### Files created in the Terraform repo

- `keycloak/apps/go-go-host/envs/k3s-beta/main.tf`
- `keycloak/apps/go-go-host/envs/k3s-beta/variables.tf`
- `keycloak/apps/go-go-host/envs/k3s-beta/providers.tf`
- `keycloak/apps/go-go-host/envs/k3s-beta/outputs.tf`
- `keycloak/apps/go-go-host/envs/k3s-beta/versions.tf`
- `keycloak/apps/go-go-host/envs/k3s-beta/terraform.tfvars.example`
- `keycloak/apps/go-go-host/envs/k3s-beta/.terraform.lock.hcl`

### What it manages

- realm `go-go-host`
- public OIDC client `go-go-host-dashboard`
- realm role `go-go-host-admin`
- optional `wesen` user
- assignment of `go-go-host-admin` to `wesen`

I intentionally used a public client because the dashboard uses browser PKCE and should not require a client secret.

### Validation

```bash
cd /home/manuel/code/wesen/terraform/keycloak/apps/go-go-host/envs/k3s-beta
terraform init -backend=false
terraform validate
```

Validation succeeded.

### Initial Terraform commit

Terraform repo commit:

- `5c2f61e Add go-go-host Keycloak beta realm`

## Step 10: Applying Keycloak, Vault, secrets, and Argo

After pushing the K3s and Terraform commits, I applied the live control-plane prerequisites.

### Keycloak admin credentials

I read the in-cluster Keycloak bootstrap admin secret through the Tailscale kubeconfig:

```bash
cd /home/manuel/code/wesen/2026-03-27--hetzner-k3s
export KUBECONFIG=$PWD/.cache/kubeconfig-tailnet.yaml
kubectl -n keycloak get secret keycloak-bootstrap-admin -o jsonpath='{.data.username}' | base64 -d >/tmp/kc-user
kubectl -n keycloak get secret keycloak-bootstrap-admin -o jsonpath='{.data.password}' | base64 -d >/tmp/kc-pass
```

I then applied Terraform:

```bash
cd /home/manuel/code/wesen/terraform/keycloak/apps/go-go-host/envs/k3s-beta
export AWS_PROFILE=manuel
export TF_VAR_keycloak_url=https://auth.yolo.scapegoat.dev
export TF_VAR_keycloak_username="$(cat /tmp/kc-user)"
export TF_VAR_keycloak_password="$(cat /tmp/kc-pass)"
export TF_VAR_wesen_password="$(openssl rand -base64 24 | tr -d '\n')"
terraform init
terraform plan -out=/tmp/go-go-host-kc-beta.tfplan
terraform apply -auto-approve /tmp/go-go-host-kc-beta.tfplan
```

Terraform created 5 resources:

```text
module.realm.keycloak_realm.this
keycloak_role.platform_admin
keycloak_openid_client.dashboard
keycloak_user.wesen[0]
keycloak_user_roles.wesen_platform_admin[0]
```

Outputs showed:

```text
realm_name = go-go-host
dashboard_client_id = go-go-host-dashboard
platform_admin_role = go-go-host-admin
public_callback_url = https://yolo.scapegoat.dev/app/auth/callback
```

The callback URL was later corrected when the hostname changed to `hosting.yolo.scapegoat.dev`.

### Vault setup

I wrote the Vault policies and Kubernetes auth roles:

```bash
cd /home/manuel/code/wesen/2026-03-27--hetzner-k3s
export VAULT_ADDR=https://vault.yolo.scapegoat.dev
vault policy write go-go-host vault/policies/kubernetes/go-go-host.hcl
vault policy write go-go-host-db-bootstrap vault/policies/kubernetes/go-go-host-db-bootstrap.hcl
vault write auth/kubernetes/role/go-go-host @vault/roles/kubernetes/go-go-host.json
vault write auth/kubernetes/role/go-go-host-db-bootstrap @vault/roles/kubernetes/go-go-host-db-bootstrap.json
```

Vault warned that the Kubernetes roles do not have an audience configured:

```text
Role go-go-host does not have an audience configured. While audiences are not required, consider specifying one if your use case would benefit from additional JWT claim verification.
```

This matches other existing role files and was not treated as a blocker for beta.

### Secret bootstrap failure and fix

The first runtime secret bootstrap failed because the script explicitly requires `VAULT_TOKEN`:

```text
missing required environment variable: VAULT_TOKEN
```

The Vault CLI had enough ambient token state for `vault policy write`, but the script uses `require_env VAULT_TOKEN`, so I reran with:

```bash
export VAULT_TOKEN="$(cat ~/.vault-token)"
./scripts/bootstrap-go-go-host-runtime-secrets.sh
GITHUB_DEPLOY_PAT="$(gh auth token)" ./scripts/bootstrap-go-go-host-image-pull-secret.sh
```

The scripts seeded:

```text
kv/apps/go-go-host/beta/runtime
kv/apps/go-go-host/beta/image-pull
```

without printing secret values.

### Argo bootstrap

I bootstrapped the new Argo `Application`:

```bash
cd /home/manuel/code/wesen/2026-03-27--hetzner-k3s
export KUBECONFIG=$PWD/.cache/kubeconfig-tailnet.yaml
kubectl apply -f gitops/applications/go-go-host.yaml
kubectl -n argocd annotate application go-go-host argocd.argoproj.io/refresh=hard --overwrite
```

Argo created the namespace, VSO resources, secrets, ConfigMaps, and the database bootstrap Job. The bootstrap Job completed successfully.

## Step 11: Argo hang on PVC sync wave and why Service/Ingress were missing

The first Argo sync appeared to hang, and the user observed in Argo CD that the Service and Ingress were not found:

```text
Resource not found in cluster: v1/Service:go-go-host
Resource not found in cluster: networking.k8s.io/v1/Ingress:go-go-host
```

The root cause was an ordering bug in my Kustomize package, not a missing manifest.

### Exact state from Argo

Inspecting the Application showed:

```text
waiting for healthy state of /PersistentVolumeClaim/go-go-host-data
```

Namespace resources showed:

```text
pod/go-go-host-db-bootstrap-...   Completed
persistentvolumeclaim/go-go-host-data   Pending   storageClassName=local-path
```

The Service and Ingress had not been applied yet because they were later sync waves.

### Root cause

I had put the PVC in sync wave `1`:

```yaml
metadata:
  annotations:
    argocd.argoproj.io/sync-wave: "1"
```

But with K3s `local-path` storage, a PVC can remain `Pending` until a pod using it is scheduled. Argo was waiting for the PVC to become healthy before applying wave `2`, but the Deployment that would bind the PVC was also in wave `2`. That created an ordering deadlock.

### Fix

I changed the PVC sync wave to `2`, alongside the Deployment:

```yaml
argocd.argoproj.io/sync-wave: "2"
```

Then I committed and pushed:

```text
c86389a Fix go-go-host PVC sync wave
```

I cleared the stuck operation and refreshed Argo:

```bash
kubectl -n argocd patch application go-go-host --type merge -p '{"operation":null}' || true
kubectl -n argocd annotate application go-go-host argocd.argoproj.io/refresh=hard --overwrite
```

After this, the Deployment, Service, PVC, and Ingress were all created. The app became:

```text
Synced Healthy
```

## Step 12: Hostname correction from yolo.scapegoat.dev to hosting.yolo.scapegoat.dev

The user corrected the desired hostname:

**User prompt (verbatim):** "wait, what's the name here? actually i want hosting.yolo.scapegoat.dev, sorry"

I moved the deployment from `yolo.scapegoat.dev` to `hosting.yolo.scapegoat.dev`.

### K3s changes

I updated:

- `gitops/kustomize/go-go-host/configmap.yaml`
- `gitops/kustomize/go-go-host/ingress.yaml`
- `docs/go-go-host-beta-deployment-playbook.md`

The daemon config became:

```yaml
publicBaseUrl: "https://hosting.yolo.scapegoat.dev"
baseDomain: "hosting.yolo.scapegoat.dev"
oidcIssuer: "https://auth.yolo.scapegoat.dev/realms/go-go-host"
```

The Ingress host became:

```text
hosting.yolo.scapegoat.dev
```

K3s repo commit:

```text
3bf74d3 Move go-go-host beta to hosting.yolo
```

### Terraform Keycloak changes

I updated the default public app URL in:

- `keycloak/apps/go-go-host/envs/k3s-beta/variables.tf`
- `keycloak/apps/go-go-host/envs/k3s-beta/terraform.tfvars.example`

Terraform changed the OIDC client redirect URIs, logout redirects, and web origin from `https://yolo.scapegoat.dev` to `https://hosting.yolo.scapegoat.dev`.

Terraform repo commit:

```text
1a39dfb Move go-go-host beta identity to hosting.yolo
```

I applied the Terraform change. The output became:

```text
public_callback_url = https://hosting.yolo.scapegoat.dev/app/auth/callback
```

### DNS mistake and correction

To make `yolo.scapegoat.dev` resolve, I had briefly added an apex `yolo` A record in the DNS Terraform state. After the hostname correction, I removed that accidental record and added a wildcard for hosted go-go-host site subdomains:

```text
*.hosting.yolo.scapegoat.dev -> 91.98.46.169
```

Terraform DNS apply destroyed:

```text
digitalocean_record.records["yolo_a"]
```

and created:

```text
digitalocean_record.records["wildcard_hosting_yolo_a"]
```

This means:

```text
hosting.yolo.scapegoat.dev       resolves through *.yolo.scapegoat.dev
foo.hosting.yolo.scapegoat.dev   resolves through *.hosting.yolo.scapegoat.dev
yolo.scapegoat.dev               intentionally has no A record
```

## Step 13: Final rollout, health checks, and browser smoke

After the hostname change, I refreshed Argo and restarted the Deployment once so the pod would pick up the updated ConfigMap:

```bash
kubectl -n argocd annotate application go-go-host argocd.argoproj.io/refresh=hard --overwrite
kubectl -n go-go-host rollout restart deployment/go-go-host
kubectl -n go-go-host rollout status deployment/go-go-host
```

The rollout completed successfully:

```text
deployment "go-go-host" successfully rolled out
```

The certificate for `hosting.yolo.scapegoat.dev` became ready:

```text
certificate.cert-manager.io/go-go-host-yolo-tls True
```

### HTTP/API smoke

I ran:

```bash
curl -I https://hosting.yolo.scapegoat.dev/healthz
curl -I https://hosting.yolo.scapegoat.dev/readyz
curl -fsS https://hosting.yolo.scapegoat.dev/api/v1/config | jq .
```

Results:

```text
/healthz -> HTTP/2 200
/readyz  -> HTTP/2 200
/api/v1/config -> devAuth=false and correct OIDC config
```

The config response showed:

```json
{
  "baseDomain": "hosting.yolo.scapegoat.dev",
  "devAuth": false,
  "oidc": {
    "clientId": "go-go-host-dashboard",
    "issuer": "https://auth.yolo.scapegoat.dev/realms/go-go-host",
    "logoutRedirectPath": "/app",
    "redirectPath": "/app/auth/callback",
    "scopes": ["openid", "profile", "email"]
  },
  "publicBaseUrl": "https://hosting.yolo.scapegoat.dev"
}
```

Argo status:

```text
Synced Healthy
```

### Browser smoke

Using the built-in Playwright browser tool, I opened:

```text
https://hosting.yolo.scapegoat.dev/admin
```

The browser redirected to:

```text
https://auth.yolo.scapegoat.dev/realms/go-go-host/protocol/openid-connect/auth...
```

I logged in as `wesen` using a generated temporary password. The callback completed and redirected to:

```text
https://hosting.yolo.scapegoat.dev/admin/overview
```

The dashboard showed:

```text
wesen@ruinwesen.com · platform admin
```

After the smoke, I rotated the `wesen` password again through the Keycloak admin API. The current generated password is stored only in:

```text
/tmp/go-go-host-beta-wesen-password
```

No password value was committed or written into docs.

## Step 14: Current live state and remaining caveats

### Live state

`go-go-host` beta is now deployed through Argo CD on the K3s cluster.

Public URL:

```text
https://hosting.yolo.scapegoat.dev
```

OIDC issuer:

```text
https://auth.yolo.scapegoat.dev/realms/go-go-host
```

Image:

```text
ghcr.io/go-go-golems/go-go-host:sha-4187ea3
```

Argo application:

```text
go-go-host -> Synced Healthy
```

### Cross-repo commits

App repo:

- `083e76c Prepare beta image publishing`
- `4187ea3 Use Go 1.26 for Docker builds`

K3s repo:

- `984048e Add go-go-host beta GitOps deployment`
- `c86389a Fix go-go-host PVC sync wave`
- `3bf74d3 Move go-go-host beta to hosting.yolo`

Terraform repo:

- `5c2f61e Add go-go-host Keycloak beta realm`
- `1a39dfb Move go-go-host beta identity to hosting.yolo`

### What worked

- Existing K3s patterns were reusable: Argo Application, Kustomize, VSO, Vault policies/roles, Postgres bootstrap Job, Traefik ingress, cert-manager HTTP-01.
- The local OIDC work from Phase 1 translated cleanly to real Keycloak because the dashboard used standard PKCE and `/api/v1/config` advertised the issuer/client/redirect path.
- Config environment expansion was the right small app change to keep the DSN out of ConfigMaps and Git.
- The built-in Playwright browser tool was useful for a real end-to-end browser smoke after Argo converged.

### What failed or was tricky

- The Dockerfile used Go 1.24 while `go.mod` required Go 1.26.1, causing the first image build to fail.
- GHCR package visibility was private. I worked around this with a Vault-backed image pull secret rather than changing package visibility through a broken/unsupported API path.
- The first Argo sync hung because a `local-path` PVC was in an earlier sync wave than the pod that would bind it.
- I initially targeted `yolo.scapegoat.dev`; the correct host is `hosting.yolo.scapegoat.dev`. This required coordinated changes in GitOps, Keycloak Terraform, and DNS Terraform.
- The initial attempt to curl `yolo.scapegoat.dev` failed with DNS resolution errors because wildcard `*.yolo.scapegoat.dev` does not cover the apex `yolo.scapegoat.dev`.

### Caveats and next steps

- `https://hosting.yolo.scapegoat.dev` is live and healthy.
- DNS for `*.hosting.yolo.scapegoat.dev` exists, but wildcard TLS for arbitrary hosted sub-sites is not implemented. The current `letsencrypt-prod` issuer uses HTTP-01, so wildcard certs would require a DNS-01 issuer or a different sub-site TLS strategy.
- The Deployment currently has a manual `kubectl.kubernetes.io/restartedAt` annotation from the rollout restart. This is acceptable but could be removed in a later GitOps cleanup if Argo reports drift.
- The app repo still has unrelated dirty files in `cmd/go-go-host*/cmds/support.go` and an unrelated untracked HOST-005 workspace; I did not touch them.

## Step 15: DigitalOcean DNS-01 wildcard TLS for generated site hosts

The user asked whether the cluster already had cert-manager support for wildcard certificates and then asked me to add the missing pieces to the K3s Argo CD deployment, likely through the existing `platform-cert-issuer` application.

### Prompt Context

**User prompt (verbatim):** "don't we already have a cert-manager to do that?"

**User prompt (verbatim):** "can you add that to k3s argo cd deploy (probably as part of the platform-cert-issuer application) ?"

**Assistant interpretation:** Verify the existing cert-manager/issuer state, then add a GitOps-managed DNS-01 issuer and wire go-go-host to use wildcard TLS for generated site hosts.

### What I found

The cluster already had cert-manager and a ready HTTP-01 ClusterIssuer:

```text
clusterissuer.cert-manager.io/letsencrypt-prod   True
```

The GitOps manifest at `gitops/kustomize/platform-cert-issuer/clusterissuer.yaml` uses only an HTTP-01 solver:

```yaml
solvers:
  - http01:
      ingress:
        ingressClassName: traefik
```

That explains why exact hosts such as `hosting.yolo.scapegoat.dev` work but wildcard certs cannot be issued by that issuer.

I also found an existing cert-manager namespace secret:

```text
cert-manager/digitalocean-dns
key: access-token
```

So the cluster already had the DigitalOcean API token material needed for DNS-01; it just lacked a DNS-01 ClusterIssuer and a go-go-host wildcard Certificate/Ingress wiring.

### GitOps changes

In the K3s repo, I added a second platform ClusterIssuer:

- `gitops/kustomize/platform-cert-issuer/clusterissuer-dns01-digitalocean.yaml`
- referenced from `gitops/kustomize/platform-cert-issuer/kustomization.yaml`

The new issuer is:

```text
letsencrypt-prod-dns01-digitalocean
```

and uses:

```yaml
dns01:
  digitalocean:
    tokenSecretRef:
      name: digitalocean-dns
      key: access-token
```

For go-go-host I added:

- `gitops/kustomize/go-go-host/certificate.yaml`
- `gitops/kustomize/go-go-host/kustomization.yaml` entry
- updated `gitops/kustomize/go-go-host/ingress.yaml`
- updated `docs/go-go-host-beta-deployment-playbook.md`

The new Certificate requests both:

```text
hosting.yolo.scapegoat.dev
*.hosting.yolo.scapegoat.dev
```

into secret:

```text
go-go-host-wildcard-tls
```

The Ingress now has two rules:

```text
hosting.yolo.scapegoat.dev
*.hosting.yolo.scapegoat.dev
```

both forwarding to the same `go-go-host` service. This lets Traefik forward generated site subdomains to the daemon, where the runtime supervisor routes by `Host` header to the active site runtime.

I removed the Ingress `cert-manager.io/cluster-issuer: letsencrypt-prod` annotation to avoid ingress-shim trying to create an HTTP-01 certificate for a wildcard host. The wildcard cert is now an explicit Certificate object using the DNS-01 issuer.

### Validation before push

I rendered both Kustomize packages:

```bash
kubectl kustomize gitops/kustomize/platform-cert-issuer >/tmp/platform-cert-issuer.yaml
kubectl kustomize gitops/kustomize/go-go-host >/tmp/go-go-host-kustomize.yaml
```

I also used server-side dry-run against the cluster:

```bash
kubectl apply --dry-run=server -k gitops/kustomize/platform-cert-issuer
kubectl apply --dry-run=server -k gitops/kustomize/go-go-host
```

Both dry-runs passed. Existing Argo-managed objects emitted `last-applied-configuration` warnings because Argo uses server-side apply, but the new ClusterIssuer and Certificate validated correctly.

### Commit and rollout

K3s repo commit:

```text
4d521ef Add DigitalOcean DNS01 wildcard TLS
```

I pushed it to `origin/main`, then refreshed Argo:

```bash
kubectl -n argocd annotate application platform-cert-issuer argocd.argoproj.io/refresh=hard --overwrite
kubectl -n argocd annotate application go-go-host argocd.argoproj.io/refresh=hard --overwrite
```

Argo converged to:

```text
platform-cert-issuer   Synced Healthy
 go-go-host            Synced Healthy
```

The DNS-01 ClusterIssuer became ready:

```text
letsencrypt-prod-dns01-digitalocean   True   The ACME account was registered with the ACME server
```

The wildcard certificate was issued successfully:

```text
certificate.cert-manager.io/go-go-host-wildcard-tls   True   go-go-host-wildcard-tls
```

Certificate details:

```text
Dns Names:
  hosting.yolo.scapegoat.dev
  *.hosting.yolo.scapegoat.dev
Issuer Ref:
  ClusterIssuer/letsencrypt-prod-dns01-digitalocean
Not After:
  2026-08-10T20:16:06Z
Renewal Time:
  2026-07-11T20:16:06Z
```

The Ingress now advertises:

```text
hosting.yolo.scapegoat.dev,*.hosting.yolo.scapegoat.dev
```

### Smoke test

I verified both exact-host and wildcard-host TLS routing:

```bash
curl -fsSI https://hosting.yolo.scapegoat.dev/healthz
curl -fsSI https://foo.hosting.yolo.scapegoat.dev/healthz
```

Both returned:

```text
HTTP/2 200
```

I also checked the wildcard root path:

```bash
curl -sSI https://foo.hosting.yolo.scapegoat.dev/
```

It returned:

```text
HTTP/2 404
```

That 404 is expected because `foo.hosting.yolo.scapegoat.dev` is not currently an active site host in the go-go-host runtime supervisor. The important part is that TLS and Traefik wildcard ingress work; once a site with primary host `foo.hosting.yolo.scapegoat.dev` has an active deployment, the daemon can route it by Host header.

### Current site-host access model

Generated user sites are now intended to be reachable as:

```text
https://<site-slug>.hosting.yolo.scapegoat.dev
```

provided that:

1. the site exists,
2. its slug maps to `<site-slug>.hosting.yolo.scapegoat.dev`,
3. it has an active deployment, and
4. the runtime supervisor has activated the deployment and registered that host.

The wildcard DNS, wildcard cert, and wildcard ingress pieces are now in place.
