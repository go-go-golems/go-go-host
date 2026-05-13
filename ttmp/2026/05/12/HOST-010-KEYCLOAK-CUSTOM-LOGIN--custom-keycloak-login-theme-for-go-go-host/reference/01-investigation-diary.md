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
RelatedFiles: []
ExternalSources: []
Summary: ""
LastUpdated: 2026-05-12T21:41:07.877250084-04:00
WhatFor: ""
WhenToUse: ""
---

# Investigation Diary

## Goal

<!-- What is the purpose of this reference document? -->

## Context

<!-- Provide background context needed to use this reference -->

## Quick Reference

<!-- Provide copy/paste-ready content, API contracts, or quick-look tables -->

## Usage Examples

<!-- Show how to use this reference in practice -->

## Related

<!-- Link to related documents or resources -->

## 2026-05-12 — OS1 Keycloak login theme implementation

### Research
- Keycloak themes use FreeMarker templates + CSS + theme.properties
- Directory structure: `themes/<name>/login/{theme.properties, resources/css/, resources/img/, login.ftl, footer.ftl}`
- `theme.properties`: `parent=keycloak`, `import=common/keycloak`, `styles=css/login.css css/os1-overrides.css`
- Keycloak v26 login page uses PatternFly v4 CSS classes (`pf-c-form-control`, `pf-c-button pf-m-primary`, etc.)
- Theme caching must be disabled/cleared for dev iteration: `rm -rf /opt/keycloak/data/tmp/kc-gzip-cache`
- `loginTheme` must be set on the realm (via Admin API or realm import JSON)
- For production: package as JAR with `META-INF/keycloak-themes.json`, deploy to `providers/` dir

### Implementation
- Created `deployments/dev/keycloak/themes/go-go-host/login/` with:
  - `theme.properties` — extends `keycloak` parent, adds `os1-overrides.css`
  - `login.ftl` — custom FreeMarker template: social providers rendered ABOVE local login form, with "or" divider
  - `footer.ftl` — links to go-go-host and GitHub
  - `resources/css/os1-overrides.css` — pure monochrome OS1 overrides
- Updated `docker-compose.yaml` to mount theme directory
- Updated `realm-go-go-host.json` to set `loginTheme: go-go-host`
- Set theme via Admin API (realm import may not apply loginTheme)

### Design decisions
- Pure monochrome (black #111 on white #fff) — no color accents anywhere
- Social providers are the primary login method → rendered first, large buttons
- Local login (username/password) below a divider — secondary
- OS1 title bar with horizontal pinstripes and centered "go-go-host" label
- Compact font scale (11-12px), uppercase labels, border-box inputs
- Black Sign In button with 2px border + 4px box-shadow, press animation on click

### Remaining for production
- Package theme as JAR for the production Keycloak deployment
- Set loginTheme on production realm via Admin API or Terraform
- Add GitHub OIDC identity provider to the production realm (currently dev only has username/password)
- Upload screenshots to ticket
