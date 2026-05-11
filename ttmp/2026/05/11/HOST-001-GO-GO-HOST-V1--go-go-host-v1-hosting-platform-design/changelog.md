# Changelog

## 2026-05-11

- Initial workspace created


## 2026-05-11

Created evidence-backed go-go-host v1 intern design guide and diary covering runtime refactor, deployment model, dashboard, agent deploys, vm-system reuse, and implementation phases.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/ttmp/2026/05/11/HOST-001-GO-GO-HOST-V1--go-go-host-v1-hosting-platform-design/design-doc/01-go-go-host-v1-hosting-platform-intern-design-and-implementation-guide.md — Primary design deliverable


## 2026-05-11

Validated ticket with docmgr doctor and uploaded bundled PDF to reMarkable at /ai/2026/05/11/HOST-001-GO-GO-HOST-V1.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/ttmp/2026/05/11/HOST-001-GO-GO-HOST-V1--go-go-host-v1-hosting-platform-design/tasks.md — Updated completion checklist after validation and upload


## 2026-05-11

Updated design guide with concrete local Wish Git analysis: scoped run policy, pre-receive boundary validation, delegated run API, SSH certificate future path, and schema references.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/2026-05-01--wish-git/internal/policy/authorize.go — Wish Git policy source now referenced locally


## 2026-05-11

Re-uploaded updated go-go-host v1 design bundle with Wish Git additions to reMarkable at /ai/2026/05/11/HOST-001-GO-GO-HOST-V1.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/ttmp/2026/05/11/HOST-001-GO-GO-HOST-V1--go-go-host-v1-hosting-platform-design/design-doc/01-go-go-host-v1-hosting-platform-intern-design-and-implementation-guide.md — Updated design bundle source


## 2026-05-11

Updated research for two dashboard surfaces and replaced tasks with detailed phased implementation backlog.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/ttmp/2026/05/11/HOST-001-GO-GO-HOST-V1--go-go-host-v1-hosting-platform-design/tasks.md — Detailed phased implementation task plan


## 2026-05-11

Uploaded two-dashboard updated design copy and standalone implementation tasks PDF to reMarkable.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/ttmp/2026/05/11/HOST-001-GO-GO-HOST-V1--go-go-host-v1-hosting-platform-design/tasks.md — Standalone tasks PDF source


## 2026-05-11

Updated design and implementation task plan to require Glazed command structure for go-go-host and go-go-host-agent CLIs.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/ttmp/2026/05/11/HOST-001-GO-GO-HOST-V1--go-go-host-v1-hosting-platform-design/tasks.md — Glazed CLI tasks added to Phase 0 Phase 6 and Phase 9


## 2026-05-11

Uploaded Glazed-updated design copy and standalone tasks PDF to reMarkable.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/ttmp/2026/05/11/HOST-001-GO-GO-HOST-V1--go-go-host-v1-hosting-platform-design/design-doc/01-go-go-host-v1-hosting-platform-intern-design-and-implementation-guide.md — Glazed-updated design PDF source


## 2026-05-11

Implemented Phase 0 scaffold: daemon, health/version API, Glazed human CLI, Glazed agent CLI, config/control/store/webadmin placeholders, dev config, README, Makefile, and tests.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/cmd/go-go-host/cmds/status.go — First Glazed human CLI command


## 2026-05-11

Resolved workspace dependency issue by removing stale goja-site module-local replace and adding version-specific go.work replace for go-go-goja v0.0.0.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go.work — Central workspace replacement for go-go-goja v0.0.0


## 2026-05-11

Implemented Phase 1 control-plane schema, migration runner, initial store methods, org/site services, authorization tests, and dev Postgres/Keycloak compose stack.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/internal/store/migrations/001_initial_schema.sql — Initial control-plane schema


## 2026-05-11

Converted Phase 1 persistence to Postgres sqlc: added sqlc config, Postgres migrations, query files, generated db package, pgx store wrapper, advisory migration lock, and Postgres integration tests.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/sqlc.yaml — sqlc configuration for Postgres pgx generated store


## 2026-05-11

Wired daemon to Postgres store/migrations and started Phase 2 dev-auth API with /me plus org/site endpoints.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/internal/httpapi/auth.go — Dev auth middleware and local user provisioning


## 2026-05-11

Added OIDC bearer-token auth path and initial Glazed CLI commands for me, org create, site create, and site list; updated task checklist.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/cmd/go-go-host/cmds/me.go — New Glazed me command


## 2026-05-11

Added org membership listing API and Glazed org list command, plus bearer-token flags for implemented CLI commands.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/cmd/go-go-host/cmds/org.go — Glazed org list/create commands


## 2026-05-11

Added Glazed login command and local CLI config persistence for API URL, dev user, and bearer token; commands now load saved defaults.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/cmd/go-go-host/cmds/login.go — Local CLI login/config command


## 2026-05-11

Started Phase 3 runtime refactor: copied site JS support packages, added SiteRuntime, fixture hosted site, and tests for render/configure/fs/exec behavior.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/internal/runtime/runtime.go — Phase 3 SiteRuntime implementation


## 2026-05-11

Completed Phase 3 runtime health check and test, preparing SiteRuntime for supervisor activation gates.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/internal/runtime/runtime.go — HealthCheck implementation


## 2026-05-11

Started Phase 4 runtime supervisor with Host-header routing, activation swap, stop, status summary, and failure-preserves-current-runtime tests.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/internal/runtime/supervisor.go — Runtime supervisor and host router


## 2026-05-11

Extended Phase 4 supervisor with restart specs, request/error counters, control-core wiring, site runtime status API, and admin runtime summary API.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/internal/httpapi/runtime.go — Runtime status API handlers

