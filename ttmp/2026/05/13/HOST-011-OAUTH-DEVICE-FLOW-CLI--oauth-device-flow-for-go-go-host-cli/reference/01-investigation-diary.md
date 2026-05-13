---
Title: Investigation Diary
Ticket: HOST-011-OAUTH-DEVICE-FLOW-CLI
Status: active
Topics:
    - go-go-host
    - cli
    - oauth
    - oidc
    - keycloak
    - security
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: ../../../../../../../../../../code/wesen/terraform/keycloak/apps/go-go-host/envs/k3s-beta/main.tf
      Note: Production Keycloak client Terraform to enable device grant
    - Path: cmd/go-go-host/cmds/cli_config.go
      Note: CLI config schema and token persistence path
    - Path: cmd/go-go-host/cmds/login.go
      Note: Current CLI login command to extend with device flow
    - Path: cmd/go-go-host/cmds/support.go
      Note: Shared HTTP helpers for bearer/dev auth and future refresh integration
    - Path: deployments/dev/keycloak/realm-go-go-host.json
      Note: Local Keycloak client config to enable device grant
    - Path: internal/httpapi/handler.go
      Note: Public config endpoint and protected API route wiring
    - Path: internal/httpapi/oidc.go
      Note: Backend OIDC bearer verifier that must accept CLI client tokens
    - Path: sources/00-sources-readme.md
      Note: Index of captured protocol, Keycloak, Terraform, and live-production sources
    - Path: sources/05-keycloak-oidc-layers.md
      Note: Official Keycloak OIDC endpoint and grant-type reference
    - Path: sources/live-device-endpoint-disabled-response.json
      Note: Live negative check showing current client does not yet allow device flow
    - Path: sources/live-go-go-host-openid-configuration.json
      Note: Live production OIDC discovery document
    - Path: sources/01-rfc8628-oauth-device-authorization-grant.md
      Note: Primary OAuth Device Authorization Grant specification
ExternalSources:
    - https://datatracker.ietf.org/doc/html/rfc8628
    - https://www.keycloak.org/securing-apps/oidc-layers
    - https://github.com/keycloak/keycloak-community/blob/main/design/oauth2-device-authorization-grant.md
Summary: Chronological diary for HOST-011 OAuth Device Flow CLI research and design.
LastUpdated: 2026-05-13T08:21:03.52249013-04:00
WhatFor: Resume or review the OAuth Device Flow CLI research/design work.
WhenToUse: Use when implementing `go-go-host login` device flow, updating Keycloak client settings, or validating the ticket docs.
---


# Diary

## Goal

This diary records the investigation and documentation work for adding OAuth 2.0 Device Authorization Grant login to the `go-go-host` human CLI. It captures the ticket setup, source collection, live Keycloak observations, and design decisions needed to write the intern-facing implementation guide.

## Step 1: Create ticket and capture primary OAuth/Keycloak sources

This step created a dedicated docmgr ticket for the OAuth Device Flow CLI work and started the evidence collection before writing design conclusions. The goal was to make the ticket self-contained: an intern should be able to read the ticket, inspect the saved sources, and understand both the protocol and the current `go-go-host` codebase constraints.

I captured the primary protocol and platform references into the ticket's `sources/` folder. RFC 8628 is the normative source for the device flow itself. Keycloak's OIDC endpoint documentation explains how Keycloak exposes the flow. The Keycloak community design note and Terraform provider docs explain the Keycloak-specific client setting that must be enabled before the current CLI can use the flow.

### Prompt Context

**User prompt (verbatim):**

```text
Ok, et's create a new ticket to add oauth device flow to the CLI. I don't know much about device flow, so you'll haev to explain. You can use web search.

Create  a detailed analysis / design / implementation guide that is for a new intern, explaining all the parts of the system needed to understand what it is, with prose paragraphs and bullet point sand pseudocode and diagrams and api references and file references. It should be very clear and technical. Store in the ticket and the nupload to remarkable.


Keep a diary as you work.
```

**Assistant interpretation:** Create a new docmgr ticket, research OAuth Device Authorization Grant and Keycloak support, map the existing CLI/auth code, write a detailed intern-facing design/implementation guide, keep a chronological diary, and upload the final bundle to reMarkable.

**Inferred user intent:** Prepare a high-quality implementation plan for adding real browser-based human CLI login without requiring the user to paste tokens or enter passwords into the CLI.

### What I did
- Created ticket `HOST-011-OAUTH-DEVICE-FLOW-CLI` with topics `go-go-host,cli,oauth,oidc,keycloak,security`.
- Created the primary design doc:
  - `design-doc/01-oauth-device-flow-cli-analysis-design-and-implementation-guide.md`
- Created this investigation diary:
  - `reference/01-investigation-diary.md`
- Added starter tasks for research, codebase mapping, design writing, validation, and reMarkable upload.
- Saved primary external sources into `sources/`:
  - `01-rfc8628-oauth-device-authorization-grant.md`
  - `02-rfc8414-oauth-authorization-server-metadata.md`
  - `03-rfc7009-oauth-token-revocation.md`
  - `04-rfc8252-oauth-native-apps.md`
  - `05-keycloak-oidc-layers.md`
  - `06-keycloak-community-device-grant-design.md`
  - `07-terraform-provider-keycloak-openid-client-defuddle.md`
  - `08-terraform-provider-keycloak-openid-client.md`
- Saved live production observations into `sources/`:
  - `live-go-go-host-openid-configuration.json`
  - `live-device-endpoint-disabled-response.json`
- Added `sources/00-sources-readme.md` explaining what each source is for.

### Why
- The design document must be evidence-backed and useful to a new intern, so source material should be preserved next to the ticket rather than only referenced by URL.
- OAuth Device Flow has specific polling and error-handling rules; preserving RFC 8628 avoids relying on memory or secondary summaries.
- Keycloak support depends on both realm metadata and per-client settings, so the ticket needs Keycloak-specific evidence as well as the generic OAuth spec.

### What worked
- `docmgr ticket create-ticket` created the workspace successfully under `ttmp/2026/05/13/HOST-011-OAUTH-DEVICE-FLOW-CLI--oauth-device-flow-for-go-go-host-cli`.
- `defuddle parse ... --md` successfully captured the primary RFC and Keycloak pages.
- The live production OIDC discovery document includes `device_authorization_endpoint` and lists `urn:ietf:params:oauth:grant-type:device_code` in `grant_types_supported`.
- A direct POST to the live device authorization endpoint produced useful negative evidence: the endpoint exists, but `go-go-host-dashboard` is not yet allowed to initiate the flow.

### What didn't work
- Running `curl -fsS` against the live device authorization endpoint failed with exit code 22 because Keycloak correctly returned HTTP 400 for the disabled client. I reran without `-f` so the JSON error response could be saved:
  - `{"error":"unauthorized_client","error_description":"Client is not allowed to initiate OAuth 2.0 Device Authorization Grant. The flow is disabled for the client."}`

### What I learned
- Production Keycloak already advertises the device authorization endpoint at the realm metadata level.
- The current production client still needs `oauth2_device_authorization_grant_enabled = true` before the CLI can start a device flow.
- The Terraform provider supports this setting on `keycloak_openid_client` as `oauth2_device_authorization_grant_enabled`.

### What was tricky to build
- The live realm metadata can make it look like device flow is ready because `device_authorization_endpoint` exists. That is not sufficient. Keycloak also enforces a per-client setting, and the current client returns `unauthorized_client` until that setting is enabled.
- The design guide should explicitly distinguish realm capability, client capability, and CLI implementation.

### What warrants a second pair of eyes
- Confirm whether `go-go-host-dashboard` should be reused for CLI device flow or whether Terraform should create a separate public client such as `go-go-host-cli`.
- Confirm whether refresh tokens should be persisted in the current YAML CLI config or moved to a dedicated token store/keyring in the implementation phase.

### What should be done in the future
- Finish mapping the existing CLI auth/config code and backend token verifier.
- Write the full design/implementation guide.
- Run `docmgr doctor` and upload the final bundle to reMarkable.

### Code review instructions
- Review the captured sources through `sources/00-sources-readme.md` first.
- Review the live negative response before enabling Terraform settings; it defines the current production gap.
- Validate future implementation by checking both OIDC discovery and a successful device authorization request.

### Technical details
- Ticket path:
  - `ttmp/2026/05/13/HOST-011-OAUTH-DEVICE-FLOW-CLI--oauth-device-flow-for-go-go-host-cli`
- Live discovery endpoint:
  - `https://auth.yolo.scapegoat.dev/realms/go-go-host/.well-known/openid-configuration`
- Live device authorization endpoint:
  - `https://auth.yolo.scapegoat.dev/realms/go-go-host/protocol/openid-connect/auth/device`
- Keycloak per-client Terraform setting:
  - `oauth2_device_authorization_grant_enabled = true`

---

## Step 2: Map CLI/auth architecture and write the implementation guide

This step turned the saved source material and repository evidence into the primary intern-facing design document. The guide now explains the protocol first, then maps the current `go-go-host` code paths, then gives a phased implementation plan with API shapes, pseudocode, Terraform/client settings, tests, and security notes.

The main design recommendation is to create a separate public Keycloak client named `go-go-host-cli` with Device Authorization Grant enabled, then extend the backend verifier to accept both dashboard and CLI client IDs. This keeps dashboard PKCE login and CLI device login separate while allowing both token types to call the same API.

### Prompt Context

**User prompt (verbatim):** (same as Step 1, plus follow-up: "open with md-view view when done" and "writ large files in smaller chunks to avoid large writes. continue")

**Assistant interpretation:** Continue the ticket research/documentation workflow, write the large design document in smaller chunks, and open the finished markdown with `md-view view` after validation/upload.

**Inferred user intent:** Produce a complete, readable design guide without risking large write failures, and make the result easy to review locally.

### What I did
- Inspected current CLI auth/config files:
  - `cmd/go-go-host/cmds/login.go`
  - `cmd/go-go-host/cmds/cli_config.go`
  - `cmd/go-go-host/cmds/support.go`
  - `cmd/go-go-host/cmds/org.go`
  - `cmd/go-go-host/cmds/site.go`
  - `cmd/go-go-host/cmds/deployments.go`
  - `cmd/go-go-host/doc/login-and-config.md`
- Inspected backend OIDC/config routing:
  - `internal/httpapi/oidc.go`
  - `internal/httpapi/auth.go`
  - `internal/httpapi/handler.go`
  - `internal/config/config.go`
- Inspected local and production Keycloak configuration references:
  - `deployments/dev/keycloak/realm-go-go-host.json`
  - `/home/manuel/code/wesen/terraform/keycloak/apps/go-go-host/envs/k3s-beta/main.tf`
- Wrote the primary design guide in multiple append chunks:
  - `design-doc/01-oauth-device-flow-cli-analysis-design-and-implementation-guide.md`

### Why
- The design guide needed to explain OAuth Device Flow from first principles and connect each concept to actual implementation files.
- Writing in chunks kept edits manageable and avoided a single very large write.

### What worked
- Repository inspection confirmed that most CLI commands already share a small set of HTTP helpers, so device-flow token refresh can be integrated centrally.
- Backend OIDC verification already accepts access tokens by issuer/signature/expiry and then enforces `aud` or `azp`, which is a good foundation for accepting a CLI client ID.
- The design doc now includes diagrams, API contracts, config schemas, pseudocode, implementation phases, tests, alternatives, and file-level references.

### What didn't work
- One attempted file read targeted `cmd/go-go-host/cmds/deploy.go`, which does not exist. Deployment commands are in `cmd/go-go-host/cmds/deployments.go`.

### What I learned
- The CLI login command already anticipated a future browser OAuth flow in its help text.
- The strongest implementation path is not to replace existing dev/manual modes, but to add device flow as the default production mode when neither `--dev-user` nor `--bearer-token` is supplied.

### What was tricky to build
- The backend currently has a single `OIDCClientID`, but the clean client model uses separate dashboard and CLI clients. The design resolves this by adding an accepted-client list while keeping the existing field for dashboard config/backwards compatibility.
- The CLI has many commands, but they mostly call shared helpers. The implementation guide emphasizes central refresh-aware auth resolution to avoid duplicated token logic.

### What warrants a second pair of eyes
- Review the recommendation to add `go-go-host-cli` instead of reusing `go-go-host-dashboard`.
- Review the proposed YAML token storage approach before implementation; it is acceptable for beta but not as strong as OS keychain storage.
- Review the backend accepted-client-ID design so it does not accidentally loosen issuer/signature/expiry verification.

### What should be done in the future
- Implement the guide in focused code PRs.
- Consider OS keychain-backed token storage after the beta implementation works.
- Consider a custom Keycloak device verification theme page if the default Keycloak device page does not match the OS1 login theme enough.

### Code review instructions
- Start with the design doc executive summary and implementation phases.
- Cross-check the file references in the current-state architecture section.
- Validate the live production gap with `sources/live-device-endpoint-disabled-response.json`.

### Technical details
- Primary design doc:
  - `ttmp/2026/05/13/HOST-011-OAUTH-DEVICE-FLOW-CLI--oauth-device-flow-for-go-go-host-cli/design-doc/01-oauth-device-flow-cli-analysis-design-and-implementation-guide.md`
- Existing CLI config file mode:
  - `0600` in `saveCLIConfig`
- Recommended Keycloak client:
  - `go-go-host-cli`
- Recommended grant:
  - `urn:ietf:params:oauth:grant-type:device_code`

---

## Step 3: Validate the ticket and upload the guide to reMarkable

This step cleaned up docmgr validation issues, uploaded the final guide bundle to reMarkable, and recorded the delivery result. The first upload attempt exposed a LaTeX/Pandoc issue caused by literal `\n` sequences in the verbatim user prompt inside the diary. I replaced that inline escaped prompt with a fenced `text` block containing the actual line breaks, then the bundle upload succeeded.

I also fixed docmgr hygiene problems by adding missing vocabulary topics and making the saved Markdown source captures docmgr-friendly. The source files now have numeric prefixes and frontmatter, and `docmgr doctor` reports that all checks pass.

### Prompt Context

**User prompt (verbatim):** (same as Step 2)

**Assistant interpretation:** Finish validation/delivery, upload to reMarkable, and open the final Markdown locally for review.

**Inferred user intent:** Get a complete ticket deliverable that is both stored in the repo and available on the reMarkable for reading.

### What I did
- Added docmgr vocabulary entries for `cli`, `oauth`, `oidc`, and `keycloak`.
- Renamed Markdown source captures with numeric prefixes so docmgr no longer reports missing numeric prefixes.
- Added frontmatter to Markdown source captures so docmgr can validate the whole ticket workspace.
- Fixed the diary prompt block to avoid literal `\n` sequences that broke Pandoc/LaTeX.
- Ran `docmgr doctor --ticket HOST-011-OAUTH-DEVICE-FLOW-CLI --stale-after 30` successfully.
- Uploaded a reMarkable bundle containing:
  - the primary design guide,
  - the investigation diary,
  - the source index.

### Why
- The ticket should be clean and maintainable before handoff.
- The reMarkable upload should include the guide and enough context to understand source provenance without uploading every long RFC capture.

### What worked
- `docmgr doctor` eventually reported `All checks passed`.
- `remarquee upload bundle` succeeded and uploaded:
  - `HOST-011_OAuth_Device_Flow_CLI_Guide.pdf`
  - remote directory `/ai/2026/05/13/HOST-011-OAUTH-DEVICE-FLOW-CLI`

### What didn't work
- First upload failed with:
  - `pandoc failed: Error producing PDF.`
  - `! Undefined control sequence.`
  - `l.2004 ...ev to explain. You can use web search.\n`
- Root cause: the diary stored the original prompt as a single inline string with literal `\n` escape sequences. Pandoc/LaTeX interpreted `\n` as a LaTeX command.
- Fix: changed the prompt to a fenced `text` block with real line breaks.

### What I learned
- For diary prompt contexts that contain multiple paragraphs, fenced text blocks are safer for PDF generation than inline quoted strings containing escaped newlines.
- Docmgr scans Markdown files in `sources/`, so source captures should either have frontmatter or be stored with a convention that doctor accepts.

### What was tricky to build
- The source capture folder contains third-party Markdown that was not originally docmgr-formatted. Adding minimal frontmatter and numeric prefixes preserved the sources while keeping ticket validation clean.
- The upload failure was not in the design guide; it came from the diary. Bundled uploads require every included Markdown file to be Pandoc-safe.

### What warrants a second pair of eyes
- Confirm whether the reMarkable bundle should include only guide/diary/source index, or whether a second source-archive PDF should include the full captured RFCs.
- Review the docmgr vocabulary additions; they are broadly useful for future auth tickets.

### What should be done in the future
- If uploading source captures, bundle them separately because RFC captures are long and can make the main implementation guide harder to read.

### Code review instructions
- Check `sources/00-sources-readme.md` for the source inventory.
- Check the final `docmgr doctor` result.
- Check `/tmp/host011-remarquee-upload.log` if upload troubleshooting is needed.

### Technical details
- Validation command:
  - `docmgr doctor --ticket HOST-011-OAUTH-DEVICE-FLOW-CLI --stale-after 30`
- Upload command:
  - `remarquee upload bundle ... --name "HOST-011 OAuth Device Flow CLI Guide" --remote-dir "/ai/2026/05/13/HOST-011-OAUTH-DEVICE-FLOW-CLI" --toc-depth 2 --non-interactive`
- Upload result:
  - `OK: uploaded HOST-011_OAuth_Device_Flow_CLI_Guide.pdf -> /ai/2026/05/13/HOST-011-OAUTH-DEVICE-FLOW-CLI`

---

## Step 4: Add backend OIDC accepted-client support

This step implemented the backend half of the separate CLI-client design. The API can now publish a device-flow client ID in `/api/v1/config`, and its OIDC verifier can accept tokens whose `aud` or `azp` matches any configured accepted client ID rather than only the dashboard client.

This is the prerequisite for a clean `go-go-host-cli` Keycloak client. Without it, a token issued to a separate CLI client would verify cryptographically but fail the local audience/authorized-party check.

### Prompt Context

**User prompt (verbatim):** "create detailed tasks in the ticket (if necessary) and implement task by task, keeping a detailed diary, and committing at appropriate intervals."

**Assistant interpretation:** Add implementation-level tasks to HOST-011, then implement the design incrementally with tests, diary updates, and commits at logical checkpoints.

**Inferred user intent:** Move from design into real implementation while preserving reviewable history and ticket documentation.

### What I did
- Added detailed implementation tasks to `tasks.md`.
- Extended `internal/config.Config` with:
  - `OIDCDeviceClientID`
  - `OIDCAcceptedClientIDs`
- Defaulted `OIDCDeviceClientID` to `go-go-host-cli`.
- Defaulted accepted client IDs to `oidcClientId` plus `oidcDeviceClientId` when no explicit list is configured.
- Updated `/api/v1/config` to include `deviceClientId` inside the `oidc` object.
- Updated the frontend `OIDCConfig` TypeScript type to tolerate `deviceClientId`.
- Added `tokenMatchesAnyClient` and changed OIDC authentication to accept any configured client while preserving issuer/signature/expiry verification.
- Added/updated tests for config defaults, explicit accepted-client preservation, `/api/v1/config`, and multi-client token matching.

### Why
- The recommended architecture uses a separate public Keycloak client for CLI device flow.
- The existing API verifier only accepted one configured client ID, which would reject CLI-issued tokens.
- Publishing `deviceClientId` through `/api/v1/config` lets the CLI discover the correct client ID without hard-coding production values.

### What worked
- `go test ./internal/config ./internal/httpapi -count=1` passed.
- The existing verifier structure already used `SkipClientIDCheck` plus local `aud`/`azp` checks, which made the multi-client change small and explicit.

### What didn't work
- My first attempted edit to `internal/config/config_test.go` failed because the exact replacement block did not match the file. I re-read the file and applied a smaller targeted replacement.

### What I learned
- The backend had already been designed to accept Keycloak access tokens by checking `azp`, so accepting a second public client is a natural extension rather than a large auth rewrite.
- The config response is the right bootstrap point for CLI login because it is already public and already communicates OIDC settings to the browser dashboard.

### What was tricky to build
- The defaulting logic must not overwrite an explicit `oidcAcceptedClientIds` list. The implementation only constructs a default list when the list is empty.
- Empty client IDs must not match accidentally, so `tokenMatchesClient` now rejects an empty configured client ID.

### What warrants a second pair of eyes
- Review whether defaulting accepted clients to include `go-go-host-cli` is acceptable before the Keycloak client exists in every environment.
- Review the error message for rejected tokens; it now includes the accepted client list.

### What should be done in the future
- Wire production config values so `oidcDeviceClientId` and `oidcAcceptedClientIds` are explicit in beta deployment configuration.
- Add the actual `go-go-host-cli` Keycloak client in local realm JSON and Terraform.

### Code review instructions
- Start in `internal/config/config.go` for config shape/defaults.
- Then review `internal/httpapi/oidc.go` for multi-client matching.
- Validate with `go test ./internal/config ./internal/httpapi -count=1`.

### Technical details
- Test command:
  - `go test ./internal/config ./internal/httpapi -count=1`
- Important symbols:
  - `Config.OIDCDeviceClientID`
  - `Config.OIDCAcceptedClientIDs`
  - `tokenMatchesAnyClient`

---

## Step 5: Enable Keycloak CLI device-flow client

This step enabled the Keycloak side of OAuth Device Authorization Grant for the CLI. Local development now imports a `go-go-host-cli` public client with device flow enabled, and production Terraform now manages the same public client in the hosted `go-go-host` realm.

After applying Terraform, the live production device authorization endpoint changed from returning `unauthorized_client` for the relevant CLI client to returning a real `device_code`, `user_code`, verification URI, expiry, and polling interval. That proves the Keycloak-side prerequisite for CLI login is now in place.

### Prompt Context

**User prompt (verbatim):** (same as Step 4)

**Assistant interpretation:** Continue implementing HOST-011 task by task and commit at the Keycloak/config milestone.

**Inferred user intent:** Make device flow actually available in local and production auth infrastructure before writing the CLI polling code.

### What I did
- Added a local `go-go-host-cli` public client to `deployments/dev/keycloak/realm-go-go-host.json` with:
  - `oauth2.device.authorization.grant.enabled = true`
  - `oauth2.device.code.lifespan = 600`
  - `oauth2.device.polling.interval = 5`
- Added production Terraform resource `keycloak_openid_client.cli` with device flow enabled.
- Added Terraform variables for `cli_client_id`, `cli_device_code_lifespan`, and `cli_device_polling_interval`.
- Added Terraform output `cli_client_id`.
- Updated K3s GitOps app config to set:
  - `oidcDeviceClientId: go-go-host-cli`
  - `oidcAcceptedClientIds: [go-go-host-dashboard, go-go-host-cli]`
- Ran `terraform validate`.
- Ran Terraform plan/apply for production Keycloak.
- Verified a final Terraform plan returned `No changes`.
- Verified the live device authorization endpoint now returns a device authorization response for `client_id=go-go-host-cli`.

### Why
- The CLI cannot start device flow until the Keycloak client explicitly allows the grant.
- Production app config must eventually tell the API that CLI-issued tokens are accepted.
- Local realm import should match production behavior so implementation can be tested under `devctl`.

### What worked
- Terraform planned exactly one addition: `keycloak_openid_client.cli`.
- Terraform apply completed with `1 added, 0 changed, 0 destroyed`.
- Final `terraform plan -detailed-exitcode` reported no changes.
- Live endpoint check returned:
  - `has_device_code: true`
  - `verification_uri: https://auth.yolo.scapegoat.dev/realms/go-go-host/device`
  - `expires_in: 600`
  - `interval: 5`

### What didn't work
- N/A for this step. The Terraform provider supported the device-flow fields directly.

### What I learned
- Keycloak realm metadata advertising device flow is not enough; the per-client setting is decisive.
- Once `go-go-host-cli` was created with device flow enabled, the same endpoint immediately returned a valid device authorization response.

### What was tricky to build
- This step spans three repositories: app repo for local realm JSON and diary, Terraform repo for durable Keycloak realm/client state, and K3s GitOps for runtime API config. The commits are intentionally separate by ownership boundary.

### What warrants a second pair of eyes
- Confirm `standard_flow_enabled = false` is appropriate for the CLI client. It narrows the client to device flow instead of also allowing browser redirect code flow.
- Confirm the K3s config is applied only after an app image containing `OIDCAcceptedClientIDs` support is deployed; older binaries should ignore unknown YAML fields, but rollout ordering should still be considered.

### What should be done in the future
- Implement the CLI login flow against `go-go-host-cli`.
- After app image deployment, verify `/api/v1/config` returns `deviceClientId` in production.

### Code review instructions
- Review Terraform commit `7da4671` in `/home/manuel/code/wesen/terraform`.
- Review GitOps commit `99bff08` in `/home/manuel/code/wesen/2026-03-27--hetzner-k3s`.
- Review local realm JSON in this app repo.
- Validate production with a POST to the device authorization endpoint using `client_id=go-go-host-cli`.

### Technical details
- Terraform apply result:
  - `Apply complete! Resources: 1 added, 0 changed, 0 destroyed.`
- Terraform commit:
  - `7da4671 keycloak/go-go-host: add CLI device flow client`
- GitOps commit:
  - `99bff08 go-go-host: configure CLI OIDC accepted client`
- Live verification command:
  - `curl -sS -X POST -d 'client_id=go-go-host-cli' -d 'scope=openid profile email' https://auth.yolo.scapegoat.dev/realms/go-go-host/protocol/openid-connect/auth/device`

---

## Step 6: Implement CLI device login, refresh-aware auth, logout, and docs

This step implemented the human CLI side of OAuth Device Authorization Grant. `go-go-host login` now preserves existing dev/manual-token modes, but when neither `--dev-user` nor `--bearer-token` is provided it starts device flow: it reads public server OIDC config, discovers Keycloak endpoints, requests a device code, prints the verification URL and user code, polls the token endpoint according to RFC 8628 behavior, stores structured OIDC tokens, and validates the token with `/api/v1/me`.

The shared CLI config and HTTP helper layer now understands structured OIDC sessions. Commands continue to call the same helper functions, but a stored OIDC session can be refreshed before use when the access token is close to expiry. A new `go-go-host logout` command best-effort revokes the refresh token and clears local auth state.

### Prompt Context

**User prompt (verbatim):** (same as Step 4)

**Assistant interpretation:** Continue implementing the next HOST-011 tasks and commit after a coherent CLI-auth milestone.

**Inferred user intent:** Move beyond Keycloak setup into a usable CLI authentication implementation with tests and documentation.

### What I did
- Added `cmd/go-go-host/cmds/oidc_device.go` with helpers for:
  - public config bootstrap,
  - OIDC discovery,
  - device authorization request,
  - device-code token polling,
  - refresh-token grant,
  - token revocation,
  - scope parsing,
  - token-to-session conversion.
- Extended `CLIConfig` with structured `CLIOIDCSession` token storage.
- Updated `resolveCLISettings` to refresh OIDC access tokens when they are close to expiry and a refresh token is available.
- Updated `go-go-host login`:
  - existing `--dev-user` mode remains,
  - existing `--bearer-token` mode remains,
  - default production mode now runs OAuth Device Authorization Grant.
- Added `go-go-host logout` and registered it on the root command.
- Updated `cmd/go-go-host/doc/login-and-config.md` to document device flow, logout, troubleshooting, dev auth, and manual bearer fallback.
- Added tests for device polling behavior and scope parsing.

### Why
- Human CLI login needs to use the same centralized Keycloak login policies as the dashboard without asking for passwords or pasted tokens.
- Refresh-aware auth belongs in the shared config/helper layer so existing commands do not duplicate token logic.
- Logout should clear local tokens and revoke the refresh token when Keycloak supports revocation.

### What worked
- `go test ./cmd/go-go-host ./cmd/go-go-host/cmds ./internal/config ./internal/httpapi -count=1` passed.
- Full `go test ./... -count=1` passed.
- Polling tests cover `authorization_pending`, `slow_down`, success, denial, and scope parsing.

### What didn't work
- N/A for compilation and tests. Interactive production device login was not completed in this step because it requires browser approval by a human; the live device authorization endpoint itself was already verified in Step 5.

### What I learned
- The existing command structure made this change smaller than expected because all user commands already pass through shared auth helpers.
- Keeping dev-user and manual-bearer modes intact is important for local tests and smoke/debug workflows.

### What was tricky to build
- Polling needed to be testable without real five-second sleeps. I added `pollDeviceTokenWithSleeper` so tests can inject a no-op sleeper and inspect requested intervals.
- Refresh must preserve an existing refresh token if the provider returns only a new access token. The helper keeps the old refresh token when the token response omits a rotated one.
- Logout should not fail local cleanup just because network revocation fails. It records `revoke_error` but still clears local auth state.

### What warrants a second pair of eyes
- Review whether the YAML config should be the long-term token store, or whether a follow-up should move OIDC refresh tokens into the OS keychain.
- Review whether `go-go-host login` should attempt to open the browser automatically. The first implementation prints the URL only, which is safer for SSH/headless terminals.
- Review timeout behavior: the polling deadline uses `expires_in`; command cancellation relies on context cancellation from the CLI process.

### What should be done in the future
- Run a real browser-approved production login smoke once the new app binary with accepted-client support is deployed.
- Consider adding `--timeout` and `--no-open` flags if browser-opening support is added later.

### Code review instructions
- Start with `cmd/go-go-host/cmds/login.go` for user-visible login behavior.
- Review `cmd/go-go-host/cmds/oidc_device.go` for protocol correctness.
- Review `cmd/go-go-host/cmds/cli_config.go` for token storage and refresh behavior.
- Review `cmd/go-go-host/cmds/logout.go` for cleanup/revocation semantics.
- Validate with `go test ./... -count=1`.

### Technical details
- New files:
  - `cmd/go-go-host/cmds/oidc_device.go`
  - `cmd/go-go-host/cmds/oidc_device_test.go`
  - `cmd/go-go-host/cmds/logout.go`
- Updated docs:
  - `cmd/go-go-host/doc/login-and-config.md`
- Full validation:
  - `go test ./... -count=1`
