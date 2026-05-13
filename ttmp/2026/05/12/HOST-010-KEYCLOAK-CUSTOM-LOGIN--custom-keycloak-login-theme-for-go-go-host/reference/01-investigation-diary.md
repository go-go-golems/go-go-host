---
Title: Investigation Diary
Ticket: HOST-010-KEYCLOAK-CUSTOM-LOGIN
Status: active
Topics:
    - keycloak
    - auth
    - theming
    - devops
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: ../../../../../../../../../../code/wesen/2026-03-27--hetzner-k3s/gitops/kustomize/keycloak/deployment.yaml
      Note: Production theme JAR mount
    - Path: ../../../../../../../../../../code/wesen/2026-03-27--hetzner-k3s/gitops/kustomize/keycloak/keycloak-theme-configmap.yaml
      Note: Production Keycloak theme ConfigMap
    - Path: ../../../../../../../../../../code/wesen/terraform/keycloak/apps/go-go-host/envs/k3s-beta/main.tf
      Note: Vault-backed GitHub OAuth Terraform data source
    - Path: ../../../../../../../../../../code/wesen/terraform/keycloak/apps/go-go-host/envs/k3s-beta/variables.tf
      Note: Vault/fallback credential source variables
    - Path: deployments/dev/docker-compose.yaml
      Note: |-
        Added theme volume mount (Step 2)
        Dev Keycloak theme volume mount
    - Path: deployments/dev/keycloak/realm-go-go-host.json
      Note: Added loginTheme field (Step 2)
    - Path: deployments/dev/keycloak/themes/go-go-host/login/footer.ftl
      Note: Custom footer with go-go-host and GitHub links (Step 2)
    - Path: deployments/dev/keycloak/themes/go-go-host/login/login.ftl
      Note: Custom FreeMarker template with social providers above local login (Step 2)
    - Path: deployments/dev/keycloak/themes/go-go-host/login/resources/css/os1-overrides.css
      Note: Pure monochrome OS1 CSS overrides (Step 2)
    - Path: deployments/dev/keycloak/themes/go-go-host/login/theme.properties
      Note: Theme config extending keycloak parent (Step 2)
    - Path: docs/deployment.md
      Note: Operator docs for Vault-backed production IdP credentials
    - Path: ttmp/2026/05/12/HOST-010-KEYCLOAK-CUSTOM-LOGIN--custom-keycloak-login-theme-for-go-go-host/scripts/03-add-github-idp.sh
      Note: Repeatable local GitHub IdP setup script
ExternalSources: []
Summary: Step-by-step diary of implementing a custom OS1-styled Keycloak login theme for go-go-host
LastUpdated: 2026-05-12T21:50:00-04:00
WhatFor: Record the research, implementation, and deployment steps for the Keycloak login theme customization
WhenToUse: Resume work on HOST-010, debug theme issues, or deploy to production
---




# Diary

## Goal

Implement a custom Keycloak login theme for the go-go-host platform that matches the OS1/classic Mac aesthetic of the admin dashboard, with monochrome styling, social providers rendered prominently above local login, and an OS1-style title bar.

## Step 1: Research Keycloak theme customization

Researched Keycloak 26's theme system to understand the directory structure, FreeMarker templating, CSS override strategy, and deployment options (volume mount vs JAR). Read the official Keycloak docs at keycloak.org/ui-customization/themes and inspected the base `login.ftl` and `template.ftl` from the Keycloak themes JAR.

**Prompt Context**

**User prompt (verbatim):** "Can we make a custom landing page for the keycloak authentication? If you need to search the web, use `surf kagi search --query ...`. Create a new docmgr ticket and do the research / read what you need, so that we can first tweak the local docker compose keycloak, and then the realm for hosting on the production site. Keep a diary as you work, store scripts in scripts/ folder of the ticket."

**Assistant interpretation:** Create a custom Keycloak login page theme styled like the go-go-host dashboard, with proper research, a docmgr ticket, diary, and scripts for both local dev and production deployment.

**Inferred user intent:** Replace the default Keycloak login page with a branded OS1-styled page that matches the rest of the go-go-host dashboard, starting with dev and then deploying to production.

### What I did
- Created HOST-010 docmgr ticket with design doc, diary, tasks, and scripts directory
- Added 5 tasks: research, create theme, docker-compose mount, verify, production deploy
- Searched web via `surf kagi search` for Keycloak 26 theme structure
- Read official Keycloak docs via `defuddle parse` from keycloak.org/ui-customization/themes
- Extracted the base `login.ftl` from the Keycloak themes JAR inside the running container
- Examined the dev docker-compose.yaml and realm-go-go-host.json for current Keycloak setup
- Examined the production Keycloak deployment.yaml in the K3s GitOps repo

### Why
- Need to understand Keycloak's theme SPI before building: directory layout, theme.properties format, FreeMarker template structure, CSS class names in v26, and caching behavior
- Need to know the dev vs production deployment differences before designing the mounting strategy

### What worked
- `surf kagi search` returned relevant results including the official docs and a 2026 blog post
- `defuddle parse` extracted clean article content from the official Keycloak docs
- `docker cp` + `jar tf`/`jar xf` successfully extracted the base login.ftl from the running Keycloak container's theme JAR at `/opt/keycloak/lib/lib/main/org.keycloak.keycloak-themes-26.0.8.jar`

### What didn't work
- `defuddle` as a subcommand of `surf` doesn't exist; `defuddle` is a separate CLI tool that takes `parse <url>`
- `jar` command isn't available inside the Keycloak container; had to `docker cp` the JAR out and extract locally
- Keycloak theme JAR path differs from docs: it's at `lib/lib/main/` not `lib/`

### What I learned
- Keycloak v26 login pages use PatternFly v4 CSS classes: `.pf-c-form-control`, `.pf-c-button.pf-m-primary`, `.pf-c-input-group`, `.card-pf`, `.login-pf-page`
- Theme structure: `themes/<name>/login/{theme.properties, login.ftl, footer.ftl, resources/css/, resources/img/, messages/}`
- `theme.properties` supports `parent=keycloak`, `import=common/keycloak`, `styles=css/login.css css/my-override.css`
- For production: package as JAR with `META-INF/keycloak-themes.json` → deploy to `providers/`
- For dev: mount theme directory into `/opt/keycloak/themes/` via docker compose volume
- `loginTheme` must be set on the realm via Admin API (`PUT /admin/realms/{realm}`) — realm import JSON may not reliably apply it
- Theme cache must be cleared for dev: `rm -rf /opt/keycloak/data/tmp/kc-gzip-cache`

### What was tricky to build
- Finding the actual DOM class names used by Keycloak v26's login page — the base FreeMarker template uses `properties.kcXxxClass` variables that resolve from the `keycloak` theme's `theme.properties`, making it hard to predict exact CSS class names without inspecting the rendered page
- Keycloak aggressively caches themes; dev iteration requires clearing the gzip cache directory after every CSS change

### What warrants a second pair of eyes
- The choice to extend `keycloak` parent (not `base`) means we inherit all PatternFly CSS but only override what we target — any Keycloak upgrade that changes class names will silently break our styling
- The `login.ftl` override copies the entire base template; upgrades to Keycloak's login.ftl won't be picked up automatically

### What should be done in the future
- Consider a regression test that screenshots the login page after Keycloak upgrades
- When Keycloak upgrades, diff the base login.ftl against our override to detect template drift

### Code review instructions
- Review `deployments/dev/keycloak/themes/go-go-host/login/theme.properties` to verify parent/import are correct
- Inspect the rendered login page DOM to confirm CSS class selectors match

### Technical details
- Keycloak theme JAR path in container: `/opt/keycloak/lib/lib/main/org.keycloak.keycloak-themes-26.0.8.jar`
- Base login.ftl location in JAR: `theme/base/login/login.ftl`
- Base template.ftl location: `theme/base/login/template.ftl`
- Admin API to set login theme: `PUT /admin/realms/go-go-host` with body `{"loginTheme": "go-go-host"}`
- Theme cache: `/opt/keycloak/data/tmp/kc-gzip-cache`

---

## Step 2: Create OS1-styled login theme and mount in dev Keycloak

Created the custom `go-go-host` login theme with pure monochrome OS1 styling, mounted it into the dev Keycloak container, and verified it renders correctly. The initial version had teal accents which were later removed per user feedback.

**Commit (code):** f38feb4 — "Add OS1-styled Keycloak login theme (HOST-010)"

### What I did
- Created `deployments/dev/keycloak/themes/go-go-host/login/` with:
  - `theme.properties` — extends `keycloak` parent, imports `common/keycloak`, adds `css/login.css css/os1-overrides.css`
  - `login.ftl` — custom FreeMarker template: social providers rendered ABOVE local login form with "or" divider
  - `footer.ftl` — links to go-go-host and GitHub
  - `resources/css/os1-overrides.css` — OS1-styled CSS overrides
- Updated `deployments/dev/docker-compose.yaml` to add volume mount: `./keycloak/themes/go-go-host:/opt/keycloak/themes/go-go-host:ro`
- Updated `deployments/dev/keycloak/realm-go-go-host.json` to add `"loginTheme": "go-go-host"`
- Set theme via Admin API (realm import didn't reliably apply loginTheme)
- Created scripts in `ttmp/.../scripts/`:
  - `01-restart-keycloak-with-theme.sh` — restarts dev Keycloak and waits for readiness
  - `02-set-realm-login-theme.sh` — sets loginTheme on the go-go-host realm via Admin REST API

### Why
- The dev Keycloak needs the theme files available inside the container; volume mount is simplest for dev iteration
- Social providers (GitHub OIDC) should be the primary login method for most users, so they need to be above the username/password form
- Pure monochrome matches the OS1 dashboard aesthetic — no colored accents on the login page

### What worked
- Volume mount in docker-compose.yaml took effect immediately after container recreation
- `devctl restart keycloak` successfully recreated the container with the new volume mount
- Admin API `PUT /admin/realms/go-go-host` with `{"loginTheme": "go-go-host"}` worked immediately
- `#kc-header-wrapper::after` with `content: 'go-go-host'` renders the title text on top of the striped background
- `font-size: 0` on `#kc-header-wrapper` hides the raw realm display name text while keeping the `::after` pseudo-element visible
- Keycloak's `footer.ftl` macro `<#macro content>` successfully renders custom footer content

### What didn't work
- Initial CSS selectors didn't match Keycloak v26 DOM: used `.login-pf body` and `input#kc-login` but the actual classes are `.pf-c-form-control`, `.pf-c-button.pf-m-primary`, `.card-pf`, `.login-pf-page`
- First attempt had teal/green accents on button, links, and focus outlines — user explicitly said "no greenish accent either"
- First attempt hid `#kc-header` entirely (`display: none`) which removed the title bar; had to change to `display: block`
- `#kc-header-wrapper * { display: none }` to hide realm name also hid the `::after` pseudo-element in some cases; switched to `font-size: 0` trick
- Footer links (`#kc-login-footer-links a`) were still blue because the CSS selectors didn't cover all link states; fixed with `a, a:visited, a:link { color: #111 !important }`
- `devctl restart keycloak` sends a signal to the existing container — the volume mount change required container recreation (which `docker compose restart` does in this case)

### What I learned
- Keycloak v26 login DOM structure: `.login-pf-page > #kc-header > #kc-header-wrapper` (realm name), `.card-pf` (login card), `.login-pf-header > #kc-page-title` (page title), `#kc-content > #kc-content-wrapper > #kc-form > #kc-form-wrapper > form` (form)
- The `login.ftl` template has named sections (`header`, `form`, `socialProviders`, `info`) that `template.ftl` renders in a specific order; overriding `login.ftl` lets us reorder (social providers first)
- Keycloak's theme cache at `/opt/keycloak/data/tmp/kc-gzip-cache` must be cleared (`rm -rf`) for CSS changes to take effect in dev mode
- `!important` is necessary for virtually every CSS override because PatternFly and Keycloak's own stylesheets use high-specificity selectors

### What was tricky to build
- Getting the OS1 title bar to show "go-go-host" instead of the realm display name "go-go-host local dev": the realm name is a raw text node inside `#kc-header-wrapper`, not wrapped in a child element. The trick: set `font-size: 0` on the wrapper to hide the text node, then use `::after { content: 'go-go-host'; font-size: 11px }` to render the desired title.
- Reordering social providers above the form: the base `template.ftl` renders the `socialProviders` section after the form. To put social first, we had to override `login.ftl` entirely and move the social provider rendering into the `form` section above the local login form.
- The realm import JSON's `loginTheme` field was set to `null` after import; had to set it separately via the Admin REST API.

### What warrants a second pair of eyes
- The `login.ftl` override copies the entire base template; if Keycloak's base template changes in a future version (e.g., new form fields, CSRF tokens, session handling), our override won't pick up those changes. This is an accepted risk documented in Keycloak's own theme guide.
- The `!important` on every CSS rule is fragile; a future PatternFly upgrade that adds more specific selectors could silently override our styles.

### What should be done in the future
- Package theme as JAR for production deployment: create `META-INF/keycloak-themes.json`, build JAR, add as init container or volume to the production Keycloak deployment
- Set `loginTheme` on the production realm via the Keycloak Admin API or Terraform
- Add GitHub OIDC identity provider to the production realm so the social provider buttons actually appear
- Consider a CI check that screenshots the login page and compares against a reference image after Keycloak upgrades
- When Keycloak is upgraded, diff the base `login.ftl` against our override to detect template drift

### Code review instructions
- Start in `deployments/dev/keycloak/themes/go-go-host/login/login.ftl` — verify social providers section is above the form, and the divider renders
- Check `deployments/dev/keycloak/themes/go-go-host/login/resources/css/os1-overrides.css` — verify all `!important` overrides target the correct PatternFly v4 classes
- Open `http://127.0.0.1:18080/realms/go-go-host/protocol/openid-connect/auth?client_id=go-go-host-dashboard&redirect_uri=http://127.0.0.1:5173/app/auth/callback&response_type=code&scope=openid&code_challenge=E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM&code_challenge_method=S256` to view the login page
- Verify: monochrome only, striped title bar with "go-go-host", black Sign In button, no colored accents

### Technical details
- Theme directory: `/home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/deployments/dev/keycloak/themes/go-go-host/login/`
- Keycloak container name: `go-go-host-keycloak`
- Theme mount: `./keycloak/themes/go-go-host:/opt/keycloak/themes/go-go-host:ro`
- Admin API token: `curl -sf http://127.0.0.1:18080/realms/master/protocol/openid-connect/token -d "client_id=admin-cli" -d "username=admin" -d "password=admin" -d "grant_type=password"`
- Set theme: `curl -X PUT http://127.0.0.1:18080/admin/realms/go-go-host -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -d '{"loginTheme": "go-go-host"}'`
- Clear theme cache: `docker exec go-go-host-keycloak rm -rf /opt/keycloak/data/tmp/kc-gzip-cache`
- Production Keycloak image: `quay.io/keycloak/keycloak:26.1.0` at `auth.yolo.scapegoat.dev`
- Production K3s deployment: `/home/manuel/code/wesen/2026-03-27--hetzner-k3s/gitops/kustomize/keycloak/deployment.yaml`

---

## Step 3: Verify GitHub IdP and deploy the theme to production Keycloak

This step completed the end-to-end path: the local realm has a GitHub identity provider so the social-login button renders above the username/password form, and the production Keycloak deployment now loads the go-go-host theme JAR from GitOps. Production was verified by opening the real auth URL for `hosting.yolo.scapegoat.dev` and capturing a screenshot.

The last visual polish was removing the excessive top whitespace inside the card body above the GitHub button. The root cause was combined padding from `#kc-content` and `#kc-social-providers`; after forcing `#kc-content` top padding to zero and removing social-provider top padding, the gap dropped from roughly 30px to about 10–12px while keeping the layout readable.

**Prompt Context**

**User prompt (verbatim):** "continue. Also the title bar lines should stop at the boundaries of the login window, they spread the whole width right  now"

**Assistant interpretation:** Continue the Keycloak IdP/theme work and fix the striped OS1 title bar so it is constrained to the login window, not the entire page width.

**Inferred user intent:** Make the custom Keycloak login page feel like one coherent OS1 window before deploying and verifying it.

**User prompt (verbatim):** "this is good, it was the drop shadow. Let's continue and setup idP, the feedback is not always perfect from pinocchio"

**Assistant interpretation:** Accept the current title bar/card alignment, continue with GitHub IdP setup, and don't over-optimize based on visual model feedback.

**Inferred user intent:** Prioritize functional IdP setup and real-page verification over chasing minor screenshot critique.

**User prompt (verbatim):** "there's whitespace above idp github and the top of the window. it's the margin-top of kc-content-wrapper and then the padding of kc-content or margin higher up."

**Assistant interpretation:** Remove the internal vertical whitespace above the GitHub social login button inside the Keycloak card body.

**Inferred user intent:** Tighten the login layout so social login is prominent and visually attached to the OS1 window body.

### What I did
- Fixed link hover behavior so footer links no longer invert into a full black box; hover now removes underline and uses a simple black outline.
- Added/used a GitHub IdP setup script: `scripts/03-add-github-idp.sh`.
- Added the local GitHub IdP via Keycloak Admin API using `GITHUB_CLIENT_ID` and `GITHUB_CLIENT_SECRET` from `.envrc`.
- Verified local realm IdPs via Admin API: `["github"]`.
- Fixed the title bar width by constraining `#kc-header` / `.login-pf-page-header` to `500px` centered.
- Fixed excessive whitespace above the GitHub button by setting `#kc-content { padding: 0 1rem 0.75rem }` and `#kc-social-providers { padding-top: 0; margin-top: 0 }`.
- Repackaged the theme as `go-go-host-keycloak-theme.jar` with `META-INF/keycloak-themes.json`.
- Created a production Kubernetes ConfigMap from the JAR: `keycloak-theme-go-go-host`.
- Updated production Keycloak GitOps deployment to mount the JAR into `/opt/keycloak/providers/go-go-host-keycloak-theme.jar`.
- Added `keycloak-theme-configmap.yaml` to the production Keycloak kustomization.
- Triggered Argo CD sync for `keycloak`; verified `Synced Healthy` and new pod ready.
- Set production realm `go-go-host` loginTheme to `go-go-host` using `kcadm.sh` inside the Keycloak pod.
- Verified production `go-go-host` realm has GitHub IdP enabled and `loginTheme: go-go-host`.
- Captured production screenshot: `sources/screenshots/host-010-prod-login.png`.

### Why
- The login theme should work in both local dev and production, and social login should be the obvious primary path.
- Production Keycloak needs a packaged theme deployment mechanism; a ConfigMap-mounted JAR in `/opt/keycloak/providers/` is simple and GitOps-friendly for this small theme.
- Realm UI configuration lives in Keycloak's database, so `loginTheme` had to be set via Admin API/kcadm after the provider was available.

### What worked
- GitHub IdP appeared on the local login page above the local login form with an `OR` divider.
- Production Keycloak accepted the theme JAR mounted in `/opt/keycloak/providers/` and rendered the custom OS1 theme after rollout.
- Argo CD synced the production Keycloak deployment successfully; app status: `Synced Healthy`.
- Production login page at `https://auth.yolo.scapegoat.dev/realms/go-go-host/protocol/openid-connect/auth?...` shows the OS1 theme and GitHub social-login button.

### What didn't work
- Playwright MCP browser target was closed, so direct browser tools failed with: `Error: browserBackend.callTool: Target page, context or browser has been closed`.
- Workaround: used the repo-local `playwright` package under `web/admin` with a Node script to render screenshots headlessly.
- Initial Admin API check returned `[]` for identity providers because the user-created IdP was not yet present in the local Keycloak realm; after running the API setup, it returned `["github"]`.
- Pinocchio feedback over-emphasized a perceived title-bar/card mismatch; user clarified the visual was acceptable and the apparent issue was mostly the drop shadow.

### What I learned
- Keycloak production can load small custom theme JARs directly from a ConfigMap mounted into `/opt/keycloak/providers/`.
- The GitHub IdP renders in the login template only when the realm actually has an enabled identity provider instance; CSS/template work alone is not enough to see the social section.
- For this Keycloak login DOM, the main contributors to top spacing above social providers are `#kc-content` padding and `#kc-social-providers` padding.

### What was tricky to build
- The title bar and card are siblings (`#kc-header.login-pf-page-header` and `.card-pf`), not a single wrapper. To visually unify them, CSS has to constrain and align both separately rather than relying on shared parent geometry.
- Deploying a Keycloak theme declaratively is only half the story: the provider JAR can be GitOps-managed, but the realm selection (`loginTheme`) is stored in Keycloak's database and must be set via Admin API/kcadm or Terraform.
- The local environment's `.envrc` exports GitHub secrets automatically; commands emitted direnv noise in output, so scripts and diary should avoid printing secret values.

### What warrants a second pair of eyes
- Mounting the JAR from a ConfigMap into `/opt/keycloak/providers/` is simple, but confirm Keycloak provider loading from ConfigMap-backed files remains acceptable for the production lifecycle.
- Confirm that realm config (`loginTheme`, GitHub IdP settings) is made durable by the intended IaC source (Terraform or a bootstrap job), not only a manual Admin API mutation.
- Check whether GitHub OAuth app callback URLs are correct for both local and production; GitHub OAuth Apps usually have a single callback URL, so local/prod may need separate apps.

### What should be done in the future
- Move production realm configuration into Terraform or a repeatable bootstrap script: `loginTheme=go-go-host`, GitHub IdP alias/client credentials, first-broker-login flow settings.
- Consider adding a Keycloak smoke script that checks the production auth page for `social-github` and `go-go-host` theme assets.
- Store production GitHub client secret through Vault/Terraform, never in Git.

### Code review instructions
- Review GitOps commit `7ec5a75` in `/home/manuel/code/wesen/2026-03-27--hetzner-k3s`: it adds the theme ConfigMap and mounts the JAR in Keycloak.
- Review `deployments/dev/keycloak/themes/go-go-host/login/resources/css/os1-overrides.css` for the final spacing and title-bar constraints.
- Validate local: open the local OIDC auth URL and check that the GitHub button is above the form with `OR` divider.
- Validate production: open the production OIDC auth URL and confirm the custom OS1 theme renders.

### Technical details
- Local GitHub callback URL: `http://127.0.0.1:18080/realms/go-go-host/broker/github/endpoint`
- Production GitHub callback URL: `https://auth.yolo.scapegoat.dev/realms/go-go-host/broker/github/endpoint`
- Production Keycloak pod verification:
  - `kubectl get application keycloak -n argocd -o jsonpath='{.status.sync.status} {.status.health.status}'`
  - `kubectl exec -n keycloak <pod> -- /opt/keycloak/bin/kcadm.sh get realms/go-go-host | grep loginTheme`
- Theme JAR mount path: `/opt/keycloak/providers/go-go-host-keycloak-theme.jar`
- Production screenshot: `ttmp/2026/05/12/HOST-010-KEYCLOAK-CUSTOM-LOGIN--custom-keycloak-login-theme-for-go-go-host/sources/screenshots/host-010-prod-login.png`

---

## Step 4: Move Keycloak realm configuration into Terraform and document the deployment pipeline

This step made the Keycloak realm configuration durable by moving the go-go-host login theme selection and GitHub identity provider into the existing Terraform environment for the beta realm. It also added a top-level deployment document in the app repo that explains the whole release path: local dev, embedded dashboard build, Docker image, GitOps rollout, Keycloak runtime, Terraform-managed realm state, DNS, and secrets.

The important ownership split is now explicit: GitOps owns Kubernetes runtime artifacts such as the Keycloak pod and mounted theme JAR, while Terraform owns Keycloak realm state such as `loginTheme` and GitHub IdP configuration. The existing manual GitHub IdP was imported into Terraform state before applying, and a follow-up plan verified no drift.

**Prompt Context**

**User prompt (verbatim):** "do 2. move realm config to terraform. Write a document in docs/deployment.md about how this is all set up and works (not justk eycloak, the entire deployment pipeline)"

**Assistant interpretation:** Make the previously manual production Keycloak realm settings durable in Terraform, then write a broad deployment pipeline document for go-go-host.

**Inferred user intent:** Remove manual drift from production auth setup and create a stable operator-facing reference for future deployments.

### What I did
- Updated Terraform module `keycloak/modules/realm-base` with optional `login_theme` support.
- Updated Terraform env `keycloak/apps/go-go-host/envs/k3s-beta` to pass `login_theme = "go-go-host"` into the realm module.
- Added Terraform-managed `keycloak_oidc_github_identity_provider.github` for the production `go-go-host` realm.
- Added GitHub IdP input variables and outputs to the go-go-host k3s-beta Terraform env.
- Added `ignore_changes = [federated_identity]` to the managed `wesen` user so Terraform does not remove runtime account links created by Keycloak first-broker-login.
- Imported existing production GitHub IdP into Terraform state with ID `go-go-host/github`.
- Ran `terraform apply`; it updated the GitHub IdP metadata/settings.
- Ran a final `terraform plan -detailed-exitcode`; it reported `No changes`.
- Added `docs/deployment.md` in the app repo documenting the full deployment pipeline.

### Why
- Production realm settings were previously manual Keycloak DB state. If the realm was recreated or drifted, the pod would still have the theme JAR but the realm might stop selecting it.
- GitHub IdP configuration contains operational policy decisions (`trust_email`, `sync_mode`, visibility, display name) that should be versioned and reviewed.
- The project needed a single deployment overview that includes more than Keycloak: app build, dashboard embedding, image tagging, GitOps, Argo CD, Terraform, DNS, and secrets.

### What worked
- `terraform init` and `terraform validate` succeeded for `keycloak/apps/go-go-host/envs/k3s-beta`.
- `terraform import 'keycloak_oidc_github_identity_provider.github[0]' go-go-host/github` succeeded.
- The first full plan revealed a real risk: Terraform would remove an existing `federated_identity` block from the `wesen` user.
- Adding `ignore_changes = [federated_identity]` prevented Terraform from removing runtime GitHub account linking.
- `terraform apply /tmp/go-go-host-kc.plan` completed successfully with `0 added, 1 changed, 0 destroyed`.
- Final `terraform plan -detailed-exitcode` reported `No changes`.

### What didn't work
- Initial `terraform state list` failed because AWS credentials for the S3 backend were not selected:
  - `Error: No valid credential sources found`
  - `failed to refresh cached credentials, no EC2 IMDS role found`
- Fix: run Terraform with `AWS_PROFILE=manuel`.
- The first full plan would have removed the user's existing GitHub federated identity link. This was not acceptable because those links are runtime account state created by users logging in.

### What I learned
- The existing Terraform repo already had a `go-go-host/envs/k3s-beta` environment; we did not need to create a new Terraform stack.
- The Keycloak provider imports a GitHub IdP using ID format `realm/alias`, e.g. `go-go-host/github`.
- Keycloak user resources can observe runtime `federated_identity` links; without lifecycle ignore settings, Terraform may try to remove user-linked IdP state.

### What was tricky to build
- Separating realm-owned configuration from runtime user state: the GitHub IdP belongs in Terraform, but individual users' federated identity links are created dynamically by the broker login flow and should not be declared by Terraform.
- Terraform variables for secrets had to be provided from environment/non-committed sources. The GitHub client secret is now a Terraform input, which means the remote state backend must be treated as secret-bearing infrastructure.

### What warrants a second pair of eyes
- Confirm that `ignore_changes = [federated_identity]` is the right long-term policy for managed bootstrap users in Keycloak.
- Confirm whether Terraform should manage the `wesen` bootstrap user at all after GitHub SSO is working.
- Confirm whether GitHub OAuth credentials should be sourced by Terraform directly or applied by a Vault-backed bootstrap job to avoid secrets in Terraform state.

### What should be done in the future
- Move GitHub OAuth credentials into the intended long-term secret source, preferably Vault-backed automation rather than ad hoc environment variables.
- Consider replacing ConfigMap-mounted theme JARs with a custom Keycloak image if theme assets grow.
- Add a smoke check that verifies the production OIDC login page contains the custom theme and `social-github` link.

### Code review instructions
- Terraform review starts in `/home/manuel/code/wesen/terraform/keycloak/apps/go-go-host/envs/k3s-beta/main.tf`.
- Module review starts in `/home/manuel/code/wesen/terraform/keycloak/modules/realm-base/main.tf`.
- Deployment docs review starts in `docs/deployment.md` in the app repo.
- Validate Terraform with:
  - `terraform validate`
  - `terraform plan -detailed-exitcode` from `keycloak/apps/go-go-host/envs/k3s-beta` with the required `TF_VAR_*` values.

### Technical details
- Terraform import command used:
  - `terraform import 'keycloak_oidc_github_identity_provider.github[0]' go-go-host/github`
- Apply result:
  - `Apply complete! Resources: 0 added, 1 changed, 0 destroyed.`
- Final plan result:
  - `No changes. Your infrastructure matches the configuration.`
- New app documentation:
  - `docs/deployment.md`

---

## Step 5: Correct production GitHub OAuth credentials

This step corrected the production GitHub identity-provider credentials after noticing that the first Terraform apply used the generic `GITHUB_CLIENT_ID` and `GITHUB_CLIENT_SECRET` variables. Those variables correspond to the local OAuth app. Production must use the live OAuth app credentials exposed as `GITHUB_LIVE_CLIENT_ID` and `GITHUB_LIVE_CLIENT_SECRET`.

The Terraform plan with the live variables showed exactly one in-place change: the GitHub identity provider's client ID and client secret. Applying that plan updated production Keycloak, and a final `terraform plan -detailed-exitcode` reported no changes.

**Prompt Context**

**User prompt (verbatim):** "what idp settings did you use for github, i think you might have used the ones for local host and not the ones for the live site. The live site are in .envrc now with +GITHUB_LIVE_CLIENT_SECRET and +GITHUB_LIVE_CLIENT_ID"

**Assistant interpretation:** Check which environment variables were used for the production GitHub IdP and correct production if the local OAuth app credentials were applied.

**Inferred user intent:** Ensure the production Keycloak GitHub broker uses the OAuth app whose callback URL points at `auth.yolo.scapegoat.dev`, not the localhost OAuth app.

### What I did
- Confirmed the previous production Terraform command used `TF_VAR_github_client_id="$GITHUB_CLIENT_ID"` and `TF_VAR_github_client_secret="$GITHUB_CLIENT_SECRET"`.
- Re-ran the production Terraform plan using:
  - `TF_VAR_github_client_id="$GITHUB_LIVE_CLIENT_ID"`
  - `TF_VAR_github_client_secret="$GITHUB_LIVE_CLIENT_SECRET"`
- Verified the plan only changed `keycloak_oidc_github_identity_provider.github[0]` client ID/secret in-place.
- Applied the plan.
- Ran final `terraform plan -detailed-exitcode`; it reported `No changes`.
- Updated `docs/deployment.md` to explicitly distinguish local GitHub vars from production live vars.

### Why
- GitHub OAuth Apps normally have callback URLs tied to a particular environment. The localhost app and production app must not be confused.
- If production Keycloak uses the local client ID/secret, GitHub login can fail because the callback URL registered in GitHub does not match production's broker endpoint.

### What worked
- Terraform detected only the intended credential change.
- Apply completed with `0 added, 1 changed, 0 destroyed`.
- Final plan reported no drift.

### What didn't work
- The earlier production apply used the non-live variables. This was a process error caused by ambiguous environment variable names.

### What I learned
- The `.envrc` now exports both local and live GitHub OAuth credentials. Production Terraform must map `GITHUB_LIVE_*` into `TF_VAR_github_*`.
- The documentation must spell this out because Terraform variable names are environment-neutral while shell variable names distinguish local vs live.

### What was tricky to build
- The Terraform resource marks `client_id` and `client_secret` as sensitive, so the plan correctly hides values. That protects secrets but also means the operator must be careful about which shell variables are mapped into Terraform.

### What warrants a second pair of eyes
- Confirm the live GitHub OAuth app callback URL is exactly `https://auth.yolo.scapegoat.dev/realms/go-go-host/broker/github/endpoint`.
- Confirm the local GitHub OAuth app callback URL remains `http://127.0.0.1:18080/realms/go-go-host/broker/github/endpoint`.

### What should be done in the future
- Consider using explicit Terraform variable names such as `github_live_client_id` in production envs, or a wrapper script, to make accidental local/live credential swaps harder.
- Move GitHub OAuth credentials into a durable secret source and make the mapping explicit there.

### Code review instructions
- Review `docs/deployment.md` for the live/local variable guidance.
- Validate production Terraform with the live variables and confirm `terraform plan -detailed-exitcode` returns no changes.

### Technical details
- Correct production mapping:
  - `export TF_VAR_github_client_id="$GITHUB_LIVE_CLIENT_ID"`
  - `export TF_VAR_github_client_secret="$GITHUB_LIVE_CLIENT_SECRET"`
- Apply result:
  - `Apply complete! Resources: 0 added, 1 changed, 0 destroyed.`
- Final plan result:
  - `No changes. Your infrastructure matches the configuration.`

---

## Step 6: Move production GitHub OAuth input to Vault-backed Terraform data

This step made the production Keycloak GitHub identity provider easier to operate safely by moving the normal credential source from shell variables into Vault KV. Terraform still manages the Keycloak IdP resource, but operators no longer need to map `GITHUB_LIVE_CLIENT_ID` and `GITHUB_LIVE_CLIENT_SECRET` into Terraform variables for normal production applies.

The production live OAuth credentials were seeded into Vault at `kv/apps/go-go-host/prod/github-oauth` with keys `client_id` and `client_secret`. The Keycloak Terraform stack now reads that secret through the Vault provider and uses it as the source for `keycloak_oidc_github_identity_provider.github[0]`. A final plan, with no `TF_VAR_github_client_id` or `TF_VAR_github_client_secret` set, reported no changes.

**Prompt Context**

**User prompt (verbatim):** "go ahead"

**Assistant interpretation:** Proceed with the previously suggested follow-up: move production GitHub OAuth credential handling away from ad hoc shell variables and into a durable secret flow.

**Inferred user intent:** Make the production Keycloak IdP deployment less error-prone and avoid repeating the local/live credential mix-up.

### What I did
- Seeded Vault KV v2 secret `kv/apps/go-go-host/prod/github-oauth` using the live GitHub OAuth variables from `.envrc`:
  - `client_id="$GITHUB_LIVE_CLIENT_ID"`
  - `client_secret="$GITHUB_LIVE_CLIENT_SECRET"`
- Added the HashiCorp Vault provider to the production go-go-host Keycloak Terraform environment.
- Added a Vault data source for `apps/go-go-host/prod/github-oauth`.
- Changed the GitHub IdP resource to use local values derived from Vault by default.
- Kept `github_client_id` / `github_client_secret` as nullable fallback inputs when `github_oauth_use_vault = false`.
- Updated `terraform.tfvars.example` to document the Vault-backed default and emergency fallback mode.
- Updated `docs/deployment.md` to document normal Vault-backed Terraform operation and secret rotation.
- Ran `terraform fmt`, `terraform init`, and a production `terraform plan -detailed-exitcode` with Vault enabled and no GitHub `TF_VAR_*` credentials set.

### Why
- The previous shell-variable flow required operators to remember that `GITHUB_CLIENT_ID` was local and `GITHUB_LIVE_CLIENT_ID` was production.
- Storing the production OAuth app credentials in Vault gives the production stack a single durable source for the live IdP credentials.
- This reduces the chance of applying localhost OAuth credentials to production Keycloak.

### What worked
- The existing OIDC Vault token in `~/.vault-token` had enough permission to write and read the new secret path.
- Terraform read `data.vault_kv_secret_v2.github_oauth[0]` successfully.
- The final production Terraform plan reported `No changes. Your infrastructure matches the configuration.`

### What didn't work
- N/A. The implementation and validation worked on the first pass after adding the Vault provider and data source.

### What I learned
- Using Vault as the Terraform input source fixes the operator input problem, but it does not make Terraform state non-sensitive. The Keycloak provider still stores sensitive resource state in the encrypted remote backend.
- The deployment docs need to state both facts: use Vault for the input source, and treat the S3 Terraform state backend as secret-bearing infrastructure.

### What was tricky to build
- The fallback mode needed to remain available without requiring Vault for every possible one-off apply. The Terraform variables now default to Vault-backed credentials, while `github_client_id` and `github_client_secret` are nullable and only used when `github_oauth_use_vault = false`.
- Terraform conditionals must keep references to the counted Vault data source guarded so fallback mode can still work.

### What warrants a second pair of eyes
- Confirm the chosen Vault path, `kv/apps/go-go-host/prod/github-oauth`, matches the broader naming convention for production app secrets.
- Confirm the operator policy for this path is appropriate; current validation used an existing OIDC token.
- Confirm the Terraform state backend controls are acceptable for sensitive Keycloak provider state.

### What should be done in the future
- Consider moving Keycloak admin credentials into a similarly durable flow, while preserving a safe break-glass path.
- Consider a small wrapper script for production applies that exports `AWS_PROFILE`, `VAULT_ADDR`, and Keycloak admin credentials consistently.

### Code review instructions
- Start in `/home/manuel/code/wesen/terraform/keycloak/apps/go-go-host/envs/k3s-beta/main.tf` and review the Vault data source plus the `local.github_client_*` values.
- Review `/home/manuel/code/wesen/terraform/keycloak/apps/go-go-host/envs/k3s-beta/variables.tf` for defaults and fallback behavior.
- Review `docs/deployment.md` for the operator-facing flow.
- Validate with:
  - `terraform fmt`
  - `terraform init`
  - `terraform plan -detailed-exitcode` with `VAULT_TOKEN` set and without `TF_VAR_github_client_id` / `TF_VAR_github_client_secret`.

### Technical details
- Vault address: `https://vault.yolo.scapegoat.dev`
- Vault KV mount: `kv`
- Vault secret path: `apps/go-go-host/prod/github-oauth`
- Secret keys:
  - `client_id`
  - `client_secret`
- Terraform validation result:
  - `data.vault_kv_secret_v2.github_oauth[0]: Read complete`
  - `No changes. Your infrastructure matches the configuration.`
