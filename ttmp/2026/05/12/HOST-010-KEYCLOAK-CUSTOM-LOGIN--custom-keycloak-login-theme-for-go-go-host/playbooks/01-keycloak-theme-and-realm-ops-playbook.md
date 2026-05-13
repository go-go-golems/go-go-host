---
Title: Keycloak Theme and Realm Operations Playbook
Ticket: HOST-010-KEYCLOAK-CUSTOM-LOGIN
Status: active
Topics:
  - keycloak
  - gitops
  - terraform
  - auth
DocType: playbook
Intent: long-term
Summary: "How go-go-host Keycloak theme and realm configuration were deployed, and what the long-term ownership model should be."
---

# Keycloak Theme and Realm Operations Playbook

## Current production deployment status

The go-go-host Keycloak login customization is currently deployed with a hybrid model:

1. **Theme artifact and Keycloak pod wiring are GitOps-managed.**
   - Repo: `/home/manuel/code/wesen/2026-03-27--hetzner-k3s`
   - Commit: `7ec5a75` — `keycloak: add go-go-host OS1 login theme and volume mount`
   - Files:
     - `gitops/kustomize/keycloak/keycloak-theme-configmap.yaml`
     - `gitops/kustomize/keycloak/deployment.yaml`
     - `gitops/kustomize/keycloak/kustomization.yaml`
   - Mechanism: a small Keycloak theme JAR is stored in a Kubernetes ConfigMap and mounted into the Keycloak pod at:
     - `/opt/keycloak/providers/go-go-host-keycloak-theme.jar`

2. **Realm state was updated manually via Keycloak Admin/kcadm.**
   - Production realm: `go-go-host`
   - Setting applied:
     - `loginTheme=go-go-host`
   - Verified with:
     - `kcadm.sh get realms/go-go-host | grep loginTheme`
   - This part is **not yet Terraform-managed** and should be treated as an operational/manual mutation until migrated.

3. **GitHub identity provider exists in production realm state.**
   - Verified via `kcadm.sh get realms/go-go-host/identity-provider/instances`
   - Provider alias: `github`
   - Provider ID: `github`
   - Enabled: `true`

## What is documented today

Detailed implementation notes are in:

- `reference/01-investigation-diary.md`
  - Step 1: research
  - Step 2: local dev theme implementation
  - Step 3: GitHub IdP verification and production theme deployment

Screenshots are in:

- `sources/screenshots/host-010-prod-login.png`

GitOps deployment is captured in the K3s repo commit:

- `/home/manuel/code/wesen/2026-03-27--hetzner-k3s`
- Commit `7ec5a75`

## Recommended long-term ownership model

Use a clear split of responsibilities:

### GitOps owns Kubernetes runtime state

GitOps should manage:

- Keycloak Deployment/StatefulSet/Service/Ingress
- Mounted theme/provider artifacts
- Keycloak image version
- Resource requests/limits
- Health probes
- Vault/ExternalSecret references
- Init containers or mounted provider JARs

For this theme specifically, GitOps can either:

1. **Keep the ConfigMap-mounted JAR** for small/simple theme artifacts.
   - Pros: simple, visible in GitOps, no image build pipeline needed.
   - Cons: ConfigMap size limit, binary blob in Git, not ideal for larger assets.

2. **Prefer a custom Keycloak image for long-term production.**
   - Build an image extending `quay.io/keycloak/keycloak:<version>`.
   - Copy theme JAR into `/opt/keycloak/providers/` or copy theme files into `/opt/keycloak/themes/`.
   - Push versioned image tag.
   - GitOps only changes `image:`.
   - This is cleaner once themes/providers grow or become release artifacts.

Recommended long-term default: **custom Keycloak image** for versioned provider/theme artifacts, deployed by GitOps.

### Terraform should own realm configuration

Terraform should manage durable Keycloak realm state:

- Realm `go-go-host`
- Login theme selection: `login_theme = "go-go-host"`
- OIDC clients (dashboard, CLI, etc.)
- Redirect URIs and web origins
- Identity providers (GitHub)
- Identity-provider mappers
- First broker login flow selection
- Required actions and realm policies

This avoids manual drift where the pod has the theme but the realm does not select it.

### Vault should own secrets

Secrets should not be committed to GitOps or Terraform source files.

GitHub IdP client secret should live in Vault or another secret manager. Be careful with Terraform: even `sensitive = true` values usually still land in Terraform state. If Terraform manages the IdP, the Terraform backend must be treated as secret-bearing infrastructure.

Acceptable approaches:

1. Terraform manages Keycloak IdP and reads secret from a secure backend.
2. A Kubernetes bootstrap job reads the secret from Vault and applies IdP config via `kcadm.sh`.
3. A separate secret-sync process writes the IdP secret into Keycloak, while Terraform manages non-secret metadata.

## Practical migration path from current state

### Phase 1: Keep current GitOps theme JAR

Already done:

- Keycloak pod has theme JAR mounted through GitOps.
- Realm `go-go-host` has `loginTheme=go-go-host` manually applied.
- GitHub IdP is enabled.

### Phase 2: Make realm state durable

Choose one:

- Terraform Keycloak provider, if this repo already owns Keycloak realm state.
- A GitOps-managed bootstrap Job if the team prefers Keycloak realm operations to run in-cluster.

Minimum durable state:

- `loginTheme=go-go-host`
- GitHub IdP alias `github`
- GitHub client ID and secret source
- First broker login flow
- Trust email / store token decisions

### Phase 3: Move theme artifact to a custom image

If the theme grows beyond tiny CSS/FTL files, build a custom image:

```Dockerfile
FROM quay.io/keycloak/keycloak:26.1.0
COPY go-go-host-keycloak-theme.jar /opt/keycloak/providers/go-go-host-keycloak-theme.jar
```

Then GitOps owns only:

```yaml
image: ghcr.io/go-go-golems/keycloak-go-go-host:<tag>
```

## Validation commands

Check Argo CD status:

```bash
kubectl get application keycloak -n argocd -o jsonpath='{.status.sync.status} {.status.health.status}'
```

Check mounted provider JAR:

```bash
kubectl exec -n keycloak deploy/keycloak -- ls -l /opt/keycloak/providers/go-go-host-keycloak-theme.jar
```

Check production realm theme:

```bash
kubectl exec -n keycloak deploy/keycloak -- /opt/keycloak/bin/kcadm.sh get realms/go-go-host | grep loginTheme
```

Check production GitHub IdP:

```bash
kubectl exec -n keycloak deploy/keycloak -- /opt/keycloak/bin/kcadm.sh get realms/go-go-host/identity-provider/instances
```

Open production login page:

```text
https://auth.yolo.scapegoat.dev/realms/go-go-host/protocol/openid-connect/auth?client_id=go-go-host-dashboard&redirect_uri=https%3A%2F%2Fhosting.yolo.scapegoat.dev%2Fapp%2Fauth%2Fcallback&response_type=code&scope=openid&code_challenge=E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM&code_challenge_method=S256
```
