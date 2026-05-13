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
