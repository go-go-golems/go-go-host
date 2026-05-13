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

