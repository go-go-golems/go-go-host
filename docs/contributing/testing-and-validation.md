# Testing and validation

Use this document to choose the validation commands for a change. Do not rely on one broad command when the change affects a subsystem that has its own failure modes.

## Baseline commands

Run these for most Go changes:

```bash
go test ./...
go build ./...
```

Run a targeted test while iterating:

```bash
go test ./internal/control -run TestName -count=1
```

Use `-count=1` when you need to avoid cached results.

## Validation matrix

| Change type | Required validation |
|---|---|
| Go backend logic | `go test ./...`; targeted package tests; `go build ./...` |
| HTTP API | Relevant `internal/httpapi` tests; success and forbidden cases; bad-input case where applicable |
| Store/migration/sqlc | Start Postgres; set `GO_GO_HOST_TEST_DATABASE_URL`; run `go test ./internal/store ./internal/control -count=1` |
| Deployment validation | `go test ./internal/deploy ./internal/control ./internal/httpapi -count=1`; accepted and rejected bundle cases |
| Runtime or hosted JS module | `go test ./internal/runtime ./internal/sitejs/... -count=1`; load, health-check, and close behavior |
| Agent signing/deploy runs | `go test ./internal/control ./internal/httpapi -run 'Agent|DeployRun|Signed' -count=1`; replay/expired/forbidden cases |
| Dashboard TypeScript/API state | `cd web/admin && pnpm build` |
| Dashboard visual/component work | `cd web/admin && pnpm storybook:build`; Storybook stories for loading/error/empty/populated states |
| Embedded dashboard | `go run ./cmd/build-web`; `go test ./internal/webadmin`; `go build ./...` |
| OIDC/browser auth | Local Keycloak/dev stack plus `make oidc-e2e` when touching browser auth |
| Documentation ticket | `docmgr doctor --ticket <TICKET-ID> --stale-after 30` |

## Postgres-backed tests

Use live Postgres validation when changing migrations, sqlc queries, store wrappers, or control services that depend on database behavior.

```bash
docker compose -f deployments/dev/docker-compose.yaml up -d
export GO_GO_HOST_TEST_DATABASE_URL='postgres://go_go_host:go_go_host_dev@127.0.0.1:55432/go_go_host?sslmode=disable'
go test ./internal/store ./internal/control -count=1
```

If the command fails, include the failure in the handoff. Do not silently replace it with `go test ./...` if the changed code needs a real database.

## Dashboard validation

For dashboard changes, validate TypeScript and Storybook separately:

```bash
cd web/admin
pnpm build
pnpm storybook:build
```

A page or component change should normally include Storybook coverage for:

- Loading state.
- Empty state.
- Populated state.
- Error state.
- Permission-denied or disabled-action state when applicable.

Visual changes should be checked against `docs/contributing/playbooks/os1-admin-dashboard-ui-work-guidelines.md`.

## Deployment and runtime validation

Deployment and runtime changes must test both acceptance and rejection. Common cases:

- Missing `go-go-host.json` is rejected.
- Invalid archive paths are rejected.
- Unknown or disabled capabilities are rejected.
- Valid bundles unpack and dry-run successfully.
- Failed runtime health checks do not replace live traffic.
- Runtime close paths release resources.

Runtime changes should include tests that exercise `NewSiteRuntime`, `HealthCheck`, and `Supervisor.Activate` when relevant.

## Handoff format

Include a validation block in the final handoff or commit notes:

```text
Validation:
- go test ./... ✅
- go build ./... ✅
- cd web/admin && pnpm build ✅
- pnpm storybook:build not run: frontend not changed
- docmgr doctor --ticket HOST-123 --stale-after 30 ✅
```

If a command was skipped, say why. If a command failed, include the exact command and error summary.
