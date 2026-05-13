# Testing and validation

Validation is not a ritual at the end of a change. It is how you discover whether you changed the layer you thought you changed. A backend permission change that only passes a React smoke test is not validated. A dashboard page that only passes `go test` is not validated. Each contribution lane has its own evidence.

The purpose of this guide is to make that evidence explicit. When you hand off a change, a reviewer should know what was tested, what was not tested, and what still depends on local services such as Postgres, Keycloak, or a browser.

## The baseline

For ordinary Go changes, start with the baseline:

```bash
go test ./...
go build ./...
```

This catches compile errors, package-level unit tests, and integration tests that do not need external services. It does not prove that Postgres migrations work against a live database, that browser auth works, or that Storybook states render correctly.

## Validation matrix

| Change type | Required validation | Why this evidence matters |
|---|---|---|
| Go backend logic | `go test ./...`; targeted package test with `-run`; `go build ./...` | Product invariants live in Go services and must be executable without the dashboard. |
| HTTP API | Relevant `internal/httpapi` tests; forbidden and allowed auth cases | The route may decode correctly but still expose the wrong resource or status code. |
| Store or migration | Start Postgres; set `GO_GO_HOST_TEST_DATABASE_URL`; run `go test ./internal/store ./internal/control` | sqlc code and migrations need a real database to prove schema behavior. |
| Deployment validation | `go test ./internal/deploy ./internal/control ./internal/httpapi`; accepted and rejected bundle cases | Bundle validation is a security boundary, not a UI convenience. |
| Runtime or hosted JS module | `go test ./internal/runtime ./internal/sitejs/...`; runtime smoke fixtures | Goja runtimes are stateful and resource-owning; load, health, and close paths all matter. |
| Dashboard API state | `cd web/admin && pnpm build`; Storybook/MSW stories | TypeScript and RTK Query tags catch frontend contract mistakes. |
| Dashboard visual work | Storybook states plus screenshot review | The OS1 visual system depends on layout, not just types. |
| Embedded dashboard | `go run ./cmd/build-web`; `go test ./internal/webadmin`; `go build ./...` | The production binary serves embedded Vite assets, not the dev server. |
| OIDC/auth browser flow | dev stack plus `make oidc-e2e` when applicable | Token acquisition and callback handling need a browser-shaped test. |
| Documentation ticket | `docmgr doctor --ticket <TICKET-ID> --stale-after 30` | Ticket docs should remain searchable, related, and vocabulary-valid. |

## How to test a backend API change

A backend API change usually crosses handler, service, store, and DTO code. Test the deepest invariant first, then the transport.

```text
control service test
  proves product rule
HTTP integration test
  proves route, auth, JSON, and status mapping
CLI/frontend test or story
  proves caller behavior if user-facing
```

A useful API test usually has at least three cases:

- The actor can perform the operation when they have the right role.
- The actor receives `403` or an equivalent permission error when they lack the role.
- Bad input produces a stable client error instead of an internal server error.

## How to test a dashboard change

Dashboard validation is a state exercise. A page that only renders the happy path is unfinished because the real network has loading, empty, error, denied, and stale states.

For each page or major component, prefer Storybook stories for:

- Loading state.
- Empty state.
- Populated state.
- Error state.
- Permission-denied or unavailable-action state when relevant.

Run:

```bash
cd web/admin
pnpm build
pnpm storybook:build
```

If the change is visual, take screenshots from Storybook or the app. Compare them against the OS1 dashboard playbook rather than against generic SaaS dashboard habits.

## How to test a deployment/runtime change

Deployment and runtime tests should prove both rejection and success. A validator that only accepts good input may still let bad input through.

Useful cases include:

- Missing `go-go-host.json` is rejected.
- Unknown or disabled capability is rejected.
- Unsafe paths such as `../outside` are rejected.
- A valid bundle dry-runs and smoke-checks successfully.
- A runtime that fails health check does not replace live traffic.
- Closing or replacing a runtime releases resources.

## Recording validation in handoff

Every handoff should include a short validation block:

```text
Validation:
- go test ./... ✅
- cd web/admin && pnpm build ✅
- pnpm storybook:build not run: frontend not changed
- docmgr doctor --ticket HOST-123 --stale-after 30 ✅
```

Do not hide skipped validation. A skipped command is not a failure if it is named and justified. It becomes a problem when reviewers have to infer what was not checked.
