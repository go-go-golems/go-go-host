# go-go-host deployment pipeline

This document describes how `go-go-host` is built, configured, and deployed from local development through the beta Kubernetes environment.

The short version:

- **Application code and embedded dashboard assets** live in this repo.
- **Container images** are built from this repo and pushed to GHCR.
- **Kubernetes runtime state** is owned by the K3s GitOps repo and reconciled by Argo CD.
- **Keycloak realm configuration and DNS** are owned by Terraform.
- **Runtime secrets** come from Vault/Kubernetes secret sync, not from committed Kubernetes Secrets.

## Repositories and ownership

| Concern | Owner | Path |
|---|---|---|
| Application source, Go daemon/CLI/agent, web/admin dashboard, dev compose, docs | App repo | `/home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host` |
| K3s Kubernetes manifests, Argo CD Applications, Keycloak Deployment, go-go-host Deployment | GitOps repo | `/home/manuel/code/wesen/2026-03-27--hetzner-k3s` |
| Keycloak realm/client/IdP state, DNS, other shared infrastructure | Terraform repo | `/home/manuel/code/wesen/terraform` |

## Local development

Use `devctl` for the standard local stack.

```bash
cd /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host
devctl up --force
```

The devctl plugin starts:

- Postgres on `127.0.0.1:55432`
- Keycloak on `127.0.0.1:18080`
- go-go-host daemon on `127.0.0.1:8080`
- Vite dashboard on `127.0.0.1:5173`
- Storybook on `127.0.0.1:6007`

Local Keycloak fixture files live under:

```text
deployments/dev/keycloak/
```

The local Keycloak realm import is:

```text
deployments/dev/keycloak/realm-go-go-host.json
```

The local OS1 login theme is mounted from:

```text
deployments/dev/keycloak/themes/go-go-host/login/
```

The dev compose mount is declared in:

```text
deployments/dev/docker-compose.yaml
```

Keycloak admin in local dev:

```text
URL:      http://127.0.0.1:18080/admin/master/console/
Username: admin
Password: admin
```

## Frontend build and embedding

The production Go binary serves embedded dashboard assets from `internal/webadmin/dist`.

Build and embed the dashboard with:

```bash
cd /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host
BUILD_WEB_LOCAL=1 go run ./cmd/build-web
```

This runs the Vite build in `web/admin` and copies the result into:

```text
internal/webadmin/dist
```

Run this before building the production container whenever dashboard code changed.

Useful validation:

```bash
cd /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host
pnpm --dir web/admin build
go test ./internal/webadmin
go build ./cmd/go-go-hostd
```

## Application image build and deploy

The production image is built from this repo's `Dockerfile` and pushed to GHCR.

The image tag convention used during beta is:

```text
ghcr.io/go-go-golems/go-go-host:sha-<short-git-sha>
```

Example:

```bash
cd /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host
SHORT_SHA=$(git rev-parse --short HEAD)
docker build -t ghcr.io/go-go-golems/go-go-host:sha-${SHORT_SHA} .
docker push ghcr.io/go-go-golems/go-go-host:sha-${SHORT_SHA}
```

After pushing the image, update the K3s GitOps repo:

```text
/home/manuel/code/wesen/2026-03-27--hetzner-k3s/gitops/kustomize/go-go-host/deployment.yaml
```

Change the container image:

```yaml
image: ghcr.io/go-go-golems/go-go-host:sha-<short-git-sha>
```

Commit and push the GitOps repo:

```bash
cd /home/manuel/code/wesen/2026-03-27--hetzner-k3s
git add gitops/kustomize/go-go-host/deployment.yaml
git commit -m "go-go-host: update image to sha-<short-git-sha>"
git push origin main
```

Argo CD reconciles the deployment. You can trigger a sync manually:

```bash
kubectl patch application go-go-host -n argocd \
  --type merge \
  -p '{"operation":{"sync":{"revision":"HEAD"}}}'
```

Validate rollout:

```bash
kubectl get application go-go-host -n argocd -o jsonpath='{.status.sync.status} {.status.health.status}'
kubectl get pods -n go-go-host
kubectl get pod -n go-go-host <pod> -o jsonpath='{.spec.containers[0].image}'
```

Public beta URLs:

```text
Control plane: https://hosting.yolo.scapegoat.dev
Demo site:     https://hello.hosting.yolo.scapegoat.dev
```

## Kubernetes/GitOps runtime state

The K3s GitOps repo owns Kubernetes resources for the beta control plane:

```text
/home/manuel/code/wesen/2026-03-27--hetzner-k3s/gitops/kustomize/go-go-host/
```

Important files:

| File | Purpose |
|---|---|
| `deployment.yaml` | go-go-host daemon image, env, probes, pod settings |
| `service.yaml` | Cluster service |
| `ingress.yaml` | public routing for `hosting.yolo.scapegoat.dev` and hosted wildcard traffic |
| `certificate.yaml` | TLS certificate resources |
| `configmap.yaml` | non-secret app config |
| `vault-auth.yaml`, `vault-connection.yaml` | Vault integration |
| `runtime-secret.yaml` | runtime secret sync target/reference |
| `db-bootstrap-*` | database bootstrap job and service account |

Argo CD Application:

```text
/home/manuel/code/wesen/2026-03-27--hetzner-k3s/gitops/applications/go-go-host.yaml
```

## Keycloak deployment vs realm configuration

There are two separate layers:

1. **Keycloak server runtime** — Kubernetes/GitOps.
2. **Keycloak realm state** — Terraform.

Do not conflate them.

### Keycloak server runtime: GitOps

Production Keycloak Kubernetes resources live in:

```text
/home/manuel/code/wesen/2026-03-27--hetzner-k3s/gitops/kustomize/keycloak/
```

The custom go-go-host login theme is currently deployed as a small Keycloak theme JAR stored in a ConfigMap:

```text
gitops/kustomize/keycloak/keycloak-theme-configmap.yaml
```

The JAR is mounted into the Keycloak pod at:

```text
/opt/keycloak/providers/go-go-host-keycloak-theme.jar
```

The mount is declared in:

```text
gitops/kustomize/keycloak/deployment.yaml
```

This part is GitOps-managed and reconciled by Argo CD.

### Keycloak realm state: Terraform

The go-go-host beta Keycloak realm is managed in the Terraform repo:

```text
/home/manuel/code/wesen/terraform/keycloak/apps/go-go-host/envs/k3s-beta/
```

This environment owns:

- realm `go-go-host`
- realm display name
- selected login theme: `go-go-host`
- dashboard OIDC client `go-go-host-dashboard`
- redirect URIs and web origins
- platform admin role
- optional bootstrap/admin user
- GitHub identity provider

Apply from that directory with Keycloak admin credentials and GitHub OAuth credentials provided through environment variables or another non-committed secret source.

Typical operator flow:

```bash
cd /home/manuel/code/wesen/terraform/keycloak/apps/go-go-host/envs/k3s-beta
export AWS_PROFILE=manuel
export TF_VAR_keycloak_url=https://auth.yolo.scapegoat.dev
export TF_VAR_keycloak_username=...
export TF_VAR_keycloak_password=...
# Use the LIVE GitHub OAuth App credentials for production.
# The local variables GITHUB_CLIENT_ID / GITHUB_CLIENT_SECRET are for localhost.
export TF_VAR_github_client_id="$GITHUB_LIVE_CLIENT_ID"
export TF_VAR_github_client_secret="$GITHUB_LIVE_CLIENT_SECRET"
export TF_VAR_wesen_password=...
terraform init
terraform plan
terraform apply
```

The Terraform backend uses shared S3 remote state:

```text
bucket: go-go-golems-tf-state
key:    keycloak/apps/go-go-host/k3s-beta/terraform.tfstate
region: us-east-1
```

### GitHub identity provider

The GitHub IdP is represented in Terraform with `keycloak_oidc_github_identity_provider`.

Callback URLs:

```text
Local:      http://127.0.0.1:18080/realms/go-go-host/broker/github/endpoint
Production: https://auth.yolo.scapegoat.dev/realms/go-go-host/broker/github/endpoint
```

For GitHub OAuth Apps, local and production often need separate apps because an OAuth App typically has a single callback URL. In this repo's shell environment the naming convention is:

```bash
# local Keycloak / localhost callback
GITHUB_CLIENT_ID=...
GITHUB_CLIENT_SECRET=...

# production Keycloak / auth.yolo.scapegoat.dev callback
GITHUB_LIVE_CLIENT_ID=...
GITHUB_LIVE_CLIENT_SECRET=...
```

When running production Terraform for `keycloak/apps/go-go-host/envs/k3s-beta`, map the live variables into Terraform:

```bash
export TF_VAR_github_client_id="$GITHUB_LIVE_CLIENT_ID"
export TF_VAR_github_client_secret="$GITHUB_LIVE_CLIENT_SECRET"
```

Secrets must not be committed. The GitHub client secret should be sourced from Vault or the operator environment. If Terraform manages the IdP secret, treat the Terraform state backend as secret-bearing infrastructure.

## DNS and TLS

DNS is managed by Terraform in the shared Terraform repo, not by this app repo.

Relevant repository:

```text
/home/manuel/code/wesen/terraform
```

Public beta hosts currently include:

```text
hosting.yolo.scapegoat.dev
*.hosting.yolo.scapegoat.dev
auth.yolo.scapegoat.dev
```

TLS and ingress resources for K3s are GitOps-managed in the K3s repo. DNS records are Terraform-managed.

## Secrets

Runtime secrets should not be committed to Kubernetes manifests.

The production model is:

- durable secret material lives in Vault or another approved secret backend,
- Kubernetes gets only synced/runtime Secret objects,
- GitOps commits references, Vault auth, and wiring,
- Terraform variables carrying secrets must come from environment variables or a secure backend.

For go-go-host, examples include:

- database runtime credentials,
- Keycloak admin credentials,
- GitHub OAuth client secret,
- deploy-agent enrollment/signing material.

## Standard deployment checklist

Use this checklist for normal app releases:

1. Build and embed dashboard assets if frontend changed.
2. Run Go and frontend validation.
3. Commit app repo changes.
4. Build and push `ghcr.io/go-go-golems/go-go-host:sha-<short-sha>`.
5. Update GitOps `go-go-host/deployment.yaml` image tag.
6. Commit and push GitOps repo.
7. Trigger/wait for Argo CD sync.
8. Verify pod image, health, control-plane URL, and demo hosted site.
9. If auth/realm config changed, run Terraform from `keycloak/apps/go-go-host/envs/k3s-beta`.
10. Record validation and any manual operations in the relevant docmgr ticket diary.

## Standard Keycloak theme/realm checklist

Use this checklist when changing auth UI or IdP behavior:

1. Edit theme files in this repo under `deployments/dev/keycloak/themes/go-go-host/login/`.
2. Validate locally with dev Keycloak.
3. Package theme as a JAR.
4. Update GitOps Keycloak theme artifact or custom Keycloak image.
5. Push GitOps and sync Keycloak.
6. Manage realm selection and IdPs through Terraform in `keycloak/apps/go-go-host/envs/k3s-beta`.
7. Open the production OIDC auth URL and verify the page renders.

## Validation URLs

Production auth page (requires no login to view the login screen):

```text
https://auth.yolo.scapegoat.dev/realms/go-go-host/protocol/openid-connect/auth?client_id=go-go-host-dashboard&redirect_uri=https%3A%2F%2Fhosting.yolo.scapegoat.dev%2Fapp%2Fauth%2Fcallback&response_type=code&scope=openid&code_challenge=E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM&code_challenge_method=S256
```

Production docs/API health checks:

```bash
curl -I https://hosting.yolo.scapegoat.dev/
curl -I https://hello.hosting.yolo.scapegoat.dev/
```
