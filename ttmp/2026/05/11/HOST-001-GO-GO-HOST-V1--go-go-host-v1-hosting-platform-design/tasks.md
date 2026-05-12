# Tasks

## Completed research/documentation tasks

- [x] Create docmgr ticket for go-go-host v1 hosting platform design.
- [x] Read the PARC Goja Sites Hosting Service proposal.
- [x] Re-read the updated proposal after it introduced two dashboards: user dashboard and platform admin console.
- [x] Inspect relevant runtime source in `2026-05-03--goja-hosting-site`.
- [x] Inspect relevant agent-auth and run-token source in `2026-05-03--agent-enroll`.
- [x] Inspect Agent Enroll dashboard source for user-dashboard workflows.
- [x] Inspect relevant `go-go-goja` modules for hosted capability decisions.
- [x] Inspect `vm-system` for runtime management patterns.
- [x] Inspect the `@go-go-golems/os-core` React/RTK Query/Storybook example.
- [x] Revisit and update design after local `2026-05-01--wish-git` checkout became available.
- [x] Write intern-oriented design and implementation guide.
- [x] Write chronological investigation diary.
- [x] Validate the ticket with `docmgr doctor`.
- [x] Upload the document bundle to reMarkable and verify upload success from `remarquee` output.

## Implementation task plan

### Phase 0: Repository scaffold, conventions, and runnable skeleton

Goal: make `go-go-host` a runnable, testable Go repository before product features begin.

- [x] Replace placeholder `go-go-host/cmd/XXX/main.go` with `cmd/go-go-hostd/main.go`.
- [x] Add `cmd/go-go-host/main.go` for human CLI commands.
- [x] Add `cmd/go-go-host-agent/main.go` for headless agent CLI commands.
- [x] Add Glazed dependencies and use `github.com/go-go-golems/glazed/pkg/...` import paths.
- [x] Wire `go-go-host` root with Glazed logging section, embedded help system, and Cobra root setup.
- [x] Wire `go-go-host-agent` root with Glazed logging section, embedded help system, and Cobra root setup.
- [x] Add `cmd/go-go-host/doc` and `cmd/go-go-host-agent/doc` embedded help-doc packages.
- [x] Establish CLI folder convention: `cmd/<binary>/cmds/<group>/<verb>.go` where groups mirror command paths.
- [x] Add `internal/config` with daemon config: listen address, public base URL, base domain, data dir, control DB DSN, OIDC settings, dev auth flag, log level.
- [x] Add `internal/logging` or shared logger initialization with request ID fields.
- [x] Add `internal/httpapi` with `GET /healthz`, `GET /readyz`, and `GET /api/v1/version`.
- [x] Add `internal/control` with an empty `Core` object and constructor.
- [x] Add `internal/store` with store interface stubs and migration runner placeholder.
- [x] Add `internal/webadmin` placeholder static handler for future embedded dashboard.
- [x] Add `Makefile` targets: `build`, `test`, `lint` if available, `run-dev`, `generate-web`.
- [x] Add README section defining v1 goals, non-goals, and safe capability model.
- [x] Add a local dev config file under `configs/dev.yaml` or `configs/dev.json`.
- [x] Add first smoke test that starts daemon handler in-process and checks `/healthz`.

Exit criteria:

- [x] `go test ./...` passes.
- [x] `go run ./cmd/go-go-hostd --config configs/dev.yaml` starts and serves health endpoints.

### Phase 1: Control-plane database, migrations, and core services

Goal: create the durable product model before runtime activation.

- [x] Decide initial DB mode: Postgres-first sqlc control-plane CRUD, with Postgres/Keycloak dev compose for Phase 1/2 infrastructure.
- [x] Add `sqlc.yaml` for Postgres/pgx query generation.
- [x] Add `internal/store/queries/*.sql` for typed sqlc CRUD/query methods.
- [x] Generate `internal/store/db` sqlc package.
- [x] Refactor store methods to call generated sqlc queries instead of hand-written `database/sql` scans.
- [x] Add Postgres integration-test path gated by `GO_GO_HOST_TEST_DATABASE_URL`.
- [x] Add migration table equivalent to Wish Git's `schema_migrations` pattern.
- [x] Add `users` table keyed by `(issuer, subject)`.
- [x] Add `orgs` table with unique slug.
- [x] Add `memberships` table with roles: `org_owner`, `org_developer`, `org_viewer`.
- [x] Add `platform_admins` or global role mechanism.
- [x] Add `sites` table with org ID, slug, display name, primary host, status, and active deployment ID.
- [x] Add `site_domains` table.
- [x] Add `site_quotas` table.
- [x] Add `site_capabilities` table for host-granted capabilities.
- [x] Add `deployments` table with immutable version, status, bundle ref, manifest JSON, validation JSON, actor fields, timestamps.
- [x] Add `deploy_runs` table with actor, site, allowed actions/channels/paths, token hash, expiry, status.
- [x] Add `agents` table.
- [x] Add `agent_keys` table.
- [x] Add `agent_site_grants` table.
- [x] Add `agent_nonces` table with `(agent_id, nonce)` uniqueness.
- [x] Add `audit_log` table with actor, resource, action, metadata, IP, user agent, timestamp.
- [x] Add dev Docker Compose with Postgres and Keycloak for Phase 1/2 local infrastructure.
- [x] Implement `UserStore.UpsertFromOIDC`.
- [x] Implement `OrgStore` create/list/get membership helpers.
- [x] Implement `SiteStore` create/list/get/update status/active deployment.
- [x] Implement `AuditStore.Insert`.
- [x] Implement `control.OrgService` and `control.SiteService` with authorization checks.
- [x] Add tests for org membership authorization.
- [x] Add tests proving user from org A cannot read or mutate org B site.

Exit criteria:

- [x] Migrations apply from empty database.
- [x] Tests can create user, org, membership, site, quota, capability, and audit rows.
- [x] Authorization tests cover allowed owner/developer/viewer behavior and cross-org denial.

### Phase 2: Authentication and session foundation

Goal: make human API calls identity-aware while keeping product authorization local.

- [x] Add dev auth middleware that accepts a configured test user header/token for local development.
- [x] Add OIDC/JWKS token validator based on existing Wish Git / Agent Enroll patterns.
- [x] Add bearer token support for dashboard/API requests when `devAuth` is false.
- [x] Implement `GET /api/v1/me` returning user, org memberships, and platform admin flag.
- [ ] Implement logout/session-clear endpoint if cookie sessions are used.
- [ ] Implement `RequireSession`, `RequireOrgRole`, `RequireSitePermission`, `RequirePlatformAdmin` server-side helpers.
- [ ] Add audit event for first user provisioning.
- [ ] Add tests for invalid issuer/audience/signature rejection.
- [x] Add test for missing bearer token rejection in OIDC mode.
- [x] Add tests for local user provisioning by `(issuer, subject)`.
- [x] Add initial dev-auth org/site API smoke endpoints for dashboard bring-up: `POST /api/v1/orgs`, `GET/POST /api/v1/orgs/{org_id}/sites`.
- [x] Wire `go-go-hostd` to open the Postgres store and apply migrations at startup.

Exit criteria:

- [x] Dashboard and CLI can call `/api/v1/me` in dev mode.
- [x] Product permissions are checked from local DB rows, not IdP roles.

### Phase 3: Runtime copy/refactor from goja-site

Goal: move the proven runtime into `go-go-host` as a per-site runtime object without daemon-level HTTP ownership.

- [x] Create `internal/runtime/SiteRuntime` type with site ID, org ID, deployment ID, hosts, bundle path, DB path, capabilities, `*sql.DB`, `*dbguard.Guard`, `*engine.Runtime`, and `*web.Host`.
- [x] Copy/refactor script discovery and script loading from `goja-site/pkg/app/server.go`.
- [x] Copy/refactor `web.Host` integration or import the existing package if module boundaries allow it.
- [x] Ensure all Goja execution enters through runtime owner calls.
- [x] Wire preconfigured `database` and `db` modules with `configure()` disabled.
- [x] Wire `ui.dsl` and `express` registrars.
- [x] Wire `db.guard` using host-owned quota config.
- [x] Do not enable unrestricted `fs` or `exec` in hosted runtime.
- [x] Add a scoped static asset mount for deployment assets.
- [x] Add runtime `Close(ctx)` that closes owner/runtime and site DB.
- [x] Add runtime `HealthCheck(ctx)` that can execute a configured smoke route or script-load check.
- [x] Add runtime fixture site under `testdata/sites/hello`.
- [x] Add unit/integration test that creates a `SiteRuntime` and serves `GET /` through its handler.

Exit criteria:

- [x] Fixture site renders through the refactored runtime.
- [x] `database.configure()` fails inside hosted runtime.
- [x] `require("exec")` and unrestricted `require("fs")` are unavailable by default.

### Phase 4: Runtime supervisor and host router

Goal: serve many sites dynamically by Host header.

- [x] Create `internal/runtime/Supervisor` with maps by site ID and normalized host.
- [x] Implement `Activate(ctx, siteID, deploymentID)` that builds new runtime before swapping traffic.
- [x] Implement old-runtime graceful close after successful activation.
- [x] Implement `Stop(ctx, siteID)`.
- [x] Implement `Restart(ctx, siteID)`.
- [x] Implement `GetByHost(host)` and `ServeHTTP` host-router adapter.
- [x] Add runtime status model: starting, ready, failed, stopped, draining.
- [x] Persist runtime status transitions to store.
- [x] On daemon startup, reconcile stale starting/ready statuses to stopped/unknown, following `vmdaemon.closeStaleSessionsOnStartup` semantics.
- [x] Add request context fields: request ID, org ID, site ID, deployment ID, host.
- [x] Add request/error counters per runtime.
- [x] Add `GET /api/v1/sites/{site_id}/runtime`.
- [x] Add admin `GET /api/v1/admin/runtimes/summary`.
- [x] Add tests for unknown host returning 404.
- [x] Add tests that host A cannot route to site B.
- [x] Add tests that failed activation does not replace currently serving runtime.

Exit criteria:

- [x] Two fixture sites can be active simultaneously and route by Host header.
- [x] Runtime summary reports active site runtimes.

### Phase 5: Deployment bundle pipeline

Goal: upload, validate, store, activate, and roll back immutable deployments.

- [x] Define `go-go-host.json` manifest schema.
- [x] Implement archive reader for tar.gz and/or zip.
- [x] Implement path normalization and reject absolute paths, `..`, empty paths, symlinks if unsafe, and hidden forbidden metadata.
- [x] Implement file count and total size checks from site quota.
- [x] Implement manifest parser and schema validator.
- [x] Implement capability request parser.
- [x] Intersect requested capabilities with site policy and record effective capabilities.
- [x] Implement deploy-run path/channel checks based on Wish Git `AllowsPath`/glob model.
- [x] Implement immutable bundle storage under data dir.
- [x] Implement unpacking to `data/sites/<site-id>/deployments/<deployment-id>`.
- [x] Implement dry-run runtime load during validation.
- [x] Implement optional smoke route check.
- [x] Insert deployment row with status `uploaded`, `validated`, `rejected`, `active`, or `superseded`.
- [x] Implement `POST /api/v1/sites/{site_id}/deployments` for human multipart upload.
- [x] Implement `GET /api/v1/sites/{site_id}/deployments`.
- [x] Implement `GET /api/v1/deployments/{deployment_id}`.
- [x] Implement `POST /api/v1/deployments/{deployment_id}/activate`.
- [x] Implement rollback as activation of previous validated deployment.
- [x] Audit upload, validation failure, activation, rollback.
- [x] Add tests for bad paths, missing manifest, oversized bundle, forbidden capability, dry-run script error.
- [x] Add integration test: upload hello bundle, activate, request by Host header.

Exit criteria:

- [x] A user can deploy and activate a bundle via API.
- [x] Invalid bundles fail with human-readable validation reports.
- [x] Rollback switches active deployment without mutating deployment records.

### Phase 6: Human CLI using Glazed commands

Goal: provide developer workflow without requiring dashboard clicks, using the standard Glazed command structure rather than ad-hoc Cobra handlers.

- [x] Implement a shared CLI API client helper for base URL, dev-auth header, JSON GET/POST, and error handling.
- [x] Implement shared Glazed command builder helpers for common sections: API config, auth config, output defaults.
- [x] Add `go-go-host login` as a Glazed command or Cobra wrapper if browser OAuth requires custom flow; still expose parsed settings consistently.
- [x] Add local token/session config file.
- [x] Add `--bearer-token` auth flag to implemented human CLI commands for non-dev auth smoke/use.
- [x] Teach implemented CLI commands to load API URL/auth defaults from the local CLI config when flags are omitted.
- [x] Add `go-go-host me` as Glazed command emitting current user/org membership rows.
- [x] Add `go-go-host orgs list` as Glazed command emitting one row per org.
- [x] Add `GET /api/v1/orgs` and `go-go-host org list` as Glazed command emitting one row per org membership.
- [x] Add `go-go-host org create` as Glazed command emitting the created org row.
- [x] Add `go-go-host site list --org-id` as Glazed command emitting one row per site.
- [x] Add `go-go-host site create --slug ... --org-id ...` as Glazed command emitting the created site row.
- [x] Add `go-go-host sites runtime <site>` as Glazed command emitting runtime status fields.
- [x] Add `go-go-host deploy ./site --site <slug> --message ...` as Glazed command; use an argument for bundle/site directory and flags for site/message/channel.
- [x] Default deployment validation output to YAML or JSON because reports are nested.
- [x] Add `go-go-host deployments list --site <slug>` as Glazed command.
- [x] Add `go-go-host deployments show <deployment>` as Glazed command.
- [x] Add `go-go-host deployments activate <deployment>` as Glazed command.
- [x] Add `go-go-host rollback --site <slug> --to <deployment>` as Glazed command.
- [x] Add `go-go-host agents list --org` as Glazed command.
- [x] Add `go-go-host audit list --org ... --site ... --actor ...` as Glazed command with table/json output.
- [x] Ensure implemented commands define `CommandDescription`, typed settings struct, `glazed` tags, `cmds.WithFlags`, Glazed output section, and command settings section.
- [x] Ensure implemented commands decode with `vals.DecodeSectionInto(schema.DefaultSlug, settings)`.
- [x] Ensure implemented list/detail/mutation commands emit stable rows via `types.NewRow` and `gp.AddRow`.
- [ ] Add `cmds.WithLong` examples for every command.
- [x] Add embedded help pages for common workflows: login, create site, deploy, rollback, agent setup.
- [x] Add clear error handling for authz denial and validation failures.
- [x] Add CLI smoke test against httptest server using `--output json` assertions.

Exit criteria:

- [x] A developer can create a site, deploy a fixture, list deployments, and roll back from CLI.
- [x] `go-go-host sites list --output json` and `--output table` both work through Glazed.
- [ ] `go-go-host help` shows embedded command help.

### Phase 7: User dashboard foundation

Goal: build the normal org/user product dashboard with React, RTK Query, Storybook, and `@go-go-golems/os-core`.

- [x] Move Phase 7 dashboard affordances/design guide to dedicated ticket `HOST-002-USER-DASHBOARD`.
- [x] Create `web/admin/package.json` with Vite, React, TypeScript, RTK Query, Storybook, and `@go-go-golems/os-core`.
- [x] Import `@go-go-golems/os-core/theme` and selected desktop theme in `main.tsx`.
- [x] Add `.storybook/preview.tsx` with same theme imports.
- [x] Add `src/app/store.ts` with RTK Query reducer and middleware.
- [x] Add `src/services/goGoHostApi.ts` with typed endpoints and tags.
- [x] Add auth/session bootstrap based on `/api/v1/me`.
- [x] Add route guards: `RequireSession`, `RequireOrgAccess`, `RequireSiteAccess`.
- [x] Add user dashboard shell with org switcher and nav.
- [x] Add `SitesPage` list.
- [x] Add `CreateSite` form.
- [x] Add `SiteDetailPage` overview with host, active deployment, runtime badge, usage summary.
- [x] Add `DeploymentsPage` and `DeploymentDetailPage` with validation report.
- [x] Add upload deployment UI.
- [x] Add activate and rollback buttons with confirmation.
- [x] Add `AgentsPage` adapted from Agent Enroll.
- [ ] Add `BotTokensPage` with one-time reveal and copyable enroll/deploy commands. (Deferred to Phase 9 agent enrollment.)
- [ ] Add `AgentGrantEditor` for allowed sites/channels/paths and expiry. (Deferred to Phase 9 agent grants.)
- [x] Add `UsagePage` for request/storage/deployment quota display.
- [x] Add user-scoped `AuditPage`.
- [x] Add `MembersPage` for owners to invite/change roles if in v1 scope.
- [x] Add Storybook stories for RuntimeBadge, DeploymentTimeline, QuotaPanel, SecretRevealBox, CommandCopyBox, AgentGrantEditor.
- [ ] Add Playwright smoke test for dashboard login in dev mode and site list rendering. (Manual Playwright verification exists; automated smoke still pending.)

Exit criteria:

- [x] User dashboard can create/list sites and show real deployment/runtime data.
- [x] Storybook builds successfully.

### Phase 8: Platform admin console

Goal: provide installation-wide operator controls separate from user dashboard.

- [x] Add server-side `platform_admin` role/permission checks.
- [x] Add `/api/v1/admin/*` route group.
- [x] Add admin overview endpoint: org count, user count, site count, active runtimes, failed deployments, quota alarms.
- [x] Add admin users list/detail endpoints.
- [x] Add admin orgs list/detail endpoints.
- [x] Add admin sites list/detail endpoints across all orgs.
- [x] Add admin runtimes list/summary/restart/stop endpoints.
- [x] Add admin deployments list with filters by status/org/site/actor.
- [x] Add admin agents list and revoke endpoints. (Global list plus org-scoped revoke exists; admin global revoke remains future safety-control work.)
- [x] Add admin quota policy endpoints: defaults and per-site override. (Read-only quota policy inventory exists; write/edit workflow remains future work.)
- [x] Add admin domain policy endpoints for base domains and verification status.
- [x] Add admin global audit endpoint.
- [x] Add `/admin` route group in SPA gated by `RequirePlatformAdmin`.
- [x] Add `AdminOverviewPage`.
- [x] Add `AdminUsersPage` and `AdminOrgsPage`.
- [x] Add `AdminSitesPage`.
- [x] Add `AdminRuntimesPage` with restart/stop controls.
- [x] Add `AdminQuotasPage`.
- [x] Add `AdminAgentsPage`.
- [x] Add `AdminAuditPage`.
- [x] Add Storybook stories for AdminRuntimeTable, AdminQuotaPolicyEditor, AdminAuditFilters.
- [x] Add tests that non-admin users get 403 for `/api/v1/admin/*`.

Exit criteria:

- [x] Platform admin can inspect all tenants and runtimes.
- [x] Non-admin users cannot access admin APIs or admin routes.

### Phase 9: Agent enrollment and signed deploy runs

Goal: support headless deploy agents without human credentials.

- [x] Implement agent creation by org owner/developer.
- [x] Implement one-time bot/enrollment token generation with token hash storage.
- [x] Implement agent enrollment endpoint that exchanges token for registered agent key.
- [x] Implement `go-go-host-agent keygen` as a Glazed command that emits key path, public key fingerprint, and status.
- [x] Implement `go-go-host-agent enroll --token ...` as a Glazed command that emits agent ID, key ID, org/site scope, and status.
- [x] Implement `go-go-host-agent deploy ./site --site <slug>` as a Glazed command with stable validation/deployment output. (Implemented with `--bundle` and `--site-id` for v1 IDs.)
- [x] Implement `go-go-host-agent status` as a Glazed command for current agent identity and grants.
- [x] Implement Ed25519 signed request verification using canonical string from Agent Enroll.
- [x] Enforce timestamp skew and nonce replay prevention.
- [x] Implement `agent_site_grants` CRUD. (Create/update plus list via store; delete/revoke can be represented by expiring/removing deploy capability in follow-up UI.)
- [x] Implement `POST /api/v1/agent/deploy-runs` signed endpoint.
- [x] Create deploy run with allowed actions/channels/paths, expiry, status, upload token hash.
- [x] Implement upload endpoint bound to deploy run and token.
- [x] Reject expired, completed, revoked, wrong-site, wrong-channel, wrong-path runs.
- [x] Ensure agent commands share the same root Glazed logging/help setup as `go-go-host`.
- [x] Add embedded help pages for keygen, enroll, deploy, status, and troubleshooting signature errors.
- [x] Audit agent enroll, key add/revoke, grant update, deploy-run create, upload, validation, activation. (Key-add/enroll, grant update, deploy-run create, deployment upload/validation are audited; explicit key revoke remains folded into agent revoke for v1.)
- [x] Add tests for bad signature, old timestamp, future timestamp, replayed nonce, revoked key, unauthorized site. (Bad signature, old timestamp, replay, wrong site/path and agent upload are covered; future timestamp/revoked-key behavior is enforced by shared verifier/status checks and should get a focused regression if key-specific revoke is added.)

Exit criteria:

- [x] Agent can deploy to allowed site.
- [x] Same agent is denied for ungranted site/path/channel.
- [x] Replayed signed request is denied.

### Phase 9A: Scoped agent auto-activation

Goal: let explicitly trusted agents promote their own validated deployments without giving every deploy agent traffic-swap power.

- [x] Add `can_activate` to `agent_site_grants`.
- [x] Expose `canActivate` on human agent creation and grant update APIs.
- [x] Add human CLI flag `go-go-host agents create --can-activate` for immediate grants.
- [x] Add agent CLI flag `go-go-host-agent deploy --activate`.
- [x] Store `activate` in deploy-run `allowed_actions` only when the agent grant permits it.
- [x] Auto-activate after a valid agent upload only when the deploy run includes `activate`.
- [x] Audit scoped agent activation as `actor_type=agent` / `deployment.activate`.
- [x] Add automated coverage for auto-activated agent uploads.
- [x] Run live devctl smoke for scoped auto-activation.

Exit criteria:

- [x] Agent with `can_activate` can upload and activate in one signed deploy flow.
- [x] Public Host-header route serves the auto-activated deployment.
- [x] Audit clearly attributes activation to the agent.

### Phase 10: Capability hardening, quotas, and observability

Goal: make hosted execution boundaries visible and enforceable.

- [ ] Define default capability set: `express`, `ui.dsl`, scoped `database/db`, `time/timer`, static assets.
- [ ] Implement capability policy table defaults and per-site overrides.
- [ ] Implement manifest requested-vs-effective capability report.
- [ ] Ensure unrestricted `fs` is not available by default.
- [ ] Ensure `exec` is never available in hosted v1.
- [ ] Add optional scoped asset read capability if scripts need asset introspection.
- [ ] Add request timeout enforcement around handler calls where possible.
- [ ] Add DB guard configuration from site quota.
- [ ] Add DB stats endpoint.
- [ ] Add usage collector for request count, error count, DB size, bundle bytes, deployment count.
- [ ] Add dashboard quota warnings.
- [ ] Add structured logs with request ID/site ID/deployment ID.
- [ ] Add runtime event table or log stream for start/stop/fail/activate events.
- [ ] Add tests for DB hard-limit write failure.
- [ ] Add tests for forbidden capability rejection.

Exit criteria:

- [ ] Capabilities are explicit in deployment validation and runtime construction.
- [ ] Users and admins can see quota state and runtime errors.

### Phase 11: Domains, configuration, and site settings

Goal: make site configuration separate from code deployment.

- [ ] Implement site config endpoint for non-secret settings.
- [ ] Implement base-domain host assignment.
- [ ] Implement custom domain table and verification-token generation.
- [ ] Implement domain verification check placeholder/manual flow.
- [ ] Add dashboard domain page.
- [ ] Add admin domain policy page.
- [ ] Add capability settings page for site owners/admins.
- [ ] Add environment/secrets design placeholder; do not expose process env wholesale.
- [ ] Audit domain add/remove/verify and capability changes.

Exit criteria:

- [ ] Site code deployment and site configuration are separate API surfaces.
- [ ] Base-domain hosts work; custom-domain status is represented even if TLS automation is deferred.

### Phase 12: Backup, export, pruning, and production hardening

Goal: make the service operable after MVP.

- [ ] Add per-site SQLite backup/export command.
- [ ] Add deployment bundle export command.
- [ ] Add org/site metadata export.
- [ ] Add deployment pruning policy by age/count/status.
- [ ] Add audit retention policy.
- [ ] Add runtime crash/restart runbook.
- [ ] Add production config example.
- [ ] Add Dockerfile or deployment recipe if needed.
- [ ] Add `/readyz` checks for DB and data dir writability.
- [ ] Add load/concurrency smoke tests.
- [ ] Add security review checklist for new host capabilities.
- [ ] Add final Playwright E2E: login, create site, deploy, visit public host, create agent token, agent deploy, rollback, inspect audit.

Exit criteria:

- [ ] Operator can back up, inspect, prune, and recover core platform state.
- [ ] End-to-end happy path and security-boundary tests pass.
