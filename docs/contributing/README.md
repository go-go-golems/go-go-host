# Contributing to go-go-host

This guide explains how to work on `go-go-host` without breaking the codebase's layering, security model, runtime model, or dashboard conventions. It is intended for both human contributors and coding agents.

## Required starting points

Before making a non-trivial change, identify the contribution area and read the relevant guide.

| Area | Read first | Primary files |
|---|---|---|
| Backend API or product behavior | [`backend-service-guidelines.md`](backend-service-guidelines.md) | `internal/httpapi`, `internal/control`, `internal/store` |
| Deployment validation or hosted runtime | [`runtime-and-deployment-guidelines.md`](runtime-and-deployment-guidelines.md) | `internal/deploy`, `internal/runtime`, `internal/sitejs` |
| Dashboard UI or frontend API state | [`frontend-dashboard-guidelines.md`](frontend-dashboard-guidelines.md) | `web/admin/src`, `docs/contributing/playbooks` |
| Tests, builds, and handoff validation | [`testing-and-validation.md`](testing-and-validation.md) | `Makefile`, `web/admin/package.json`, `.github/workflows` |
| Local development and services | [`../runbooks/local-development.md`](../runbooks/local-development.md) | `configs`, `deployments/dev`, `.devctl.yaml` |
| Ticket docs, diaries, and research reports | [`docmgr-and-ticket-workflow.md`](docmgr-and-ticket-workflow.md) | `ttmp` |
| System structure | [`architecture-map.md`](architecture-map.md) | `cmd`, `internal`, `web/admin` |

## Repository structure

| Path | Purpose |
|---|---|
| `cmd/go-go-hostd` | Daemon entrypoint. Loads config, opens the store, applies migrations, constructs `control.Core`, and serves HTTP. |
| `cmd/go-go-host` | Human CLI. Commands should be Glazed commands and should call the API instead of duplicating server rules. |
| `cmd/go-go-host-agent` | Machine deployment CLI. Handles key generation, enrollment, signing, deploy-run creation, and bundle upload. |
| `internal/httpapi` | HTTP transport. Owns route registration, request decoding, response DTOs, auth middleware, and HTTP status mapping. |
| `internal/control` | Product logic. Owns authorization, service invariants, audit events, deployment workflows, and orchestration. |
| `internal/store` | Persistence. Owns migrations, generated sqlc queries, store wrappers, and database model conversion. |
| `internal/deploy` | Bundle handling. Owns archive validation, manifest parsing, path policy, capability policy, and bundle storage preparation. |
| `internal/runtime` | Hosted runtime lifecycle. Owns `SiteRuntime`, `Supervisor`, activation, restart, stop, host dispatch, and runtime status. |
| `internal/sitejs` | JavaScript-facing host modules and HTTP bridge. Owns `express`, `ui.dsl`, database guard exposure, request/response DTOs, and session helpers. |
| `internal/webadmin` | Embedded dashboard file serving. Owns `go:embed` of built Vite assets and SPA fallback behavior. |
| `web/admin` | React dashboard. Owns routes, pages, components, RTK Query API state, MSW fixtures, and Storybook stories. |
| `docs` | Stable contributor, architecture, and runbook documentation. |
| `ttmp` | Ticket workspaces, diaries, design docs, screenshots, and temporary research. |

## Layering rules

The normal backend dependency direction is:

```text
HTTP handler / CLI adapter
    -> control service
        -> store wrapper / runtime / deploy package
            -> database / filesystem / Goja runtime
```

Follow these rules:

- Put authorization and product invariants in `internal/control`, not only in HTTP handlers, CLI commands, or React components.
- Use `internal/httpapi` for transport concerns: JSON decoding, request path variables, auth middleware, response DTOs, and status codes.
- Use `internal/store` for persistence. Do not add ad-hoc SQL to handlers or frontend-facing code.
- Use `internal/deploy` for bundle validation and manifest policy. Do not duplicate bundle policy in the UI or CLI as the only enforcement point.
- Use `internal/runtime` and `internal/sitejs` for hosted JavaScript execution. Do not expose new host capabilities without explicit policy and tests.
- Use `web/admin/src/services/goGoHostApi.ts` for dashboard API calls. Do not scatter raw `fetch` calls through pages except for special cases such as multipart upload where the API client deliberately wraps the behavior.

## Standard change workflow

For a normal feature, use this sequence:

1. Identify the contribution area and read the relevant guide.
2. Inspect existing code in the same layer before creating new patterns.
3. Implement the server-side invariant first if the feature changes product behavior.
4. Add tests at the layer that owns the invariant.
5. Add transport, CLI, or dashboard integration after the invariant is tested.
6. Run the validation commands for the contribution area.
7. Update stable docs or a ticket diary when the change introduces a new workflow, invariant, or debugging lesson.
8. Commit focused changes at logical checkpoints.

## Stop and ask before changing

Get explicit review before making changes in these areas:

- Authentication, OIDC verification, dev auth, platform-admin bootstrap, or session behavior.
- Agent signature verification, nonce handling, key status, deploy-run tokens, or grant checks.
- Deployment activation, rollback, active deployment status transitions, or supervisor traffic swap behavior.
- Bundle path validation, safe capability policy, or hosted JavaScript module exposure.
- Runtime isolation, request timeouts, database guard enforcement, or filesystem access.
- Schema migrations that rewrite existing data or change core entity relationships.
- Dashboard visual system changes that diverge from the OS1 playbook.
- Backwards-compatibility adapters or shims that were not explicitly requested.

## Minimum validation

Run the validation for the changed area. Common commands:

```bash
go test ./...
go build ./...
```

For frontend changes:

```bash
cd web/admin
pnpm build
pnpm storybook:build
```

For docmgr ticket work:

```bash
docmgr doctor --ticket <TICKET-ID> --stale-after 30
```

See [`testing-and-validation.md`](testing-and-validation.md) for the full matrix.

## Documentation expectations

Use stable docs and ticket docs for different purposes:

- Put durable contributor guidance in `docs`.
- Put investigation logs, design alternatives, implementation diaries, and screenshots in `ttmp` ticket workspaces.
- Promote a ticket-local lesson into `docs` when it becomes a repeated workflow or a rule future contributors must follow.

Do not leave important operational knowledge only in chat transcripts or commit messages.
