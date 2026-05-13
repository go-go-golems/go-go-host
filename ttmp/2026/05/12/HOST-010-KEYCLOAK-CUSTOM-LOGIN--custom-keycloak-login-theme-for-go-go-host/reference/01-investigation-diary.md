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
    - Path: deployments/dev/docker-compose.yaml
      Note: Added theme volume mount (Step 2)
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
