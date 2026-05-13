# Changelog

## 2026-05-12

- Initial workspace created


## 2026-05-12

Implemented OS1 monochrome Keycloak login theme with social providers above local login

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/deployments/dev/keycloak/themes/go-go-host/login/login.ftl — Custom FreeMarker template with social providers first
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/deployments/dev/keycloak/themes/go-go-host/login/resources/css/os1-overrides.css — Pure monochrome OS1 CSS overrides
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/deployments/dev/keycloak/themes/go-go-host/login/theme.properties — Theme config extending keycloak parent


## 2026-05-12

Step 3: verified GitHub IdP locally and deployed OS1 login theme to production Keycloak (repo commit c57d79b, gitops commit 7ec5a75)

### Related Files

- /home/manuel/code/wesen/2026-03-27--hetzner-k3s/gitops/kustomize/keycloak/deployment.yaml — Production Keycloak provider JAR mount
- /home/manuel/code/wesen/2026-03-27--hetzner-k3s/gitops/kustomize/keycloak/keycloak-theme-configmap.yaml — Production theme JAR ConfigMap
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/deployments/dev/keycloak/themes/go-go-host/login/resources/css/os1-overrides.css — Final spacing/titlebar/link hover fixes


## 2026-05-12

Added playbook documenting current GitOps/manual Keycloak deployment split and recommended long-term GitOps/Terraform ownership model

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/ttmp/2026/05/12/HOST-010-KEYCLOAK-CUSTOM-LOGIN--custom-keycloak-login-theme-for-go-go-host/playbooks/01-keycloak-theme-and-realm-ops-playbook.md — Documents current deployment and long-term operating model


## 2026-05-12

Step 4: moved go-go-host Keycloak loginTheme and GitHub IdP into Terraform; added full deployment pipeline docs

### Related Files

- /home/manuel/code/wesen/terraform/keycloak/apps/go-go-host/envs/k3s-beta/main.tf — Terraform owns login theme and GitHub IdP
- /home/manuel/code/wesen/terraform/keycloak/modules/realm-base/main.tf — Realm module supports login_theme
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/docs/deployment.md — Top-level go-go-host deployment pipeline documentation


## 2026-05-12

Step 5: corrected production GitHub IdP to use GITHUB_LIVE_CLIENT_ID/GITHUB_LIVE_CLIENT_SECRET and documented live/local variable mapping

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/docs/deployment.md — Documents live vs local GitHub OAuth credential mapping for production Terraform

