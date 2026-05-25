# go-go-host

`go-go-host` is a Go-based hosting platform for small server-side JavaScript sites executed with Goja. It provides a control plane for organizations, sites, deployments, deployment agents, audit, settings, domains, quotas, capabilities, and runtime operations. It also provides a hosted runtime that routes public traffic by host name to the active deployment for each site.

The repository contains the daemon, human CLI, agent CLI, embedded React dashboard, Postgres control-plane store, deployment validator, and hosted JavaScript runtime.

## Status

The project is an active v1 implementation. The core platform loop is present:

- `go-go-hostd` runs the daemon, applies Postgres migrations, serves the HTTP API, serves the embedded dashboard, and dispatches hosted-site traffic.
- `go-go-host` is the human/operator CLI.
- `go-go-host-agent` is the machine deployment CLI for signed agent workflows.
- The backend includes users, organizations, memberships, sites, deployments, agents, agent keys, deploy runs, audit events, site settings, custom-domain placeholders, quotas, capabilities, maintenance APIs, and runtime status.
- The runtime validates deployment bundles, creates per-site Goja runtimes, wires explicit host capabilities, uses per-site SQLite, and swaps live traffic only after runtime health checks.
- The dashboard is a React/Vite/RTK Query/Storybook application served under `/app/*` and `/admin/*`.
- Developer and agent documentation is embedded in the CLI and exposed through the dashboard docs area.

The platform is suitable for local development and beta-oriented implementation work. Production hardening items such as domain/TLS automation, stronger isolation, secrets management, backup/restore automation, and full observability remain active design areas.

## Repository layout

| Path | Purpose |
|---|---|
| `cmd/go-go-hostd` | Daemon entrypoint. Loads config, opens Postgres, applies migrations, restores active runtimes, and serves HTTP. |
| `cmd/go-go-host` | Human CLI for status, login/config, orgs, sites, deployments, agents, audit, and maintenance workflows. |
| `cmd/go-go-host-agent` | Agent CLI for key generation, enrollment, signed deploy runs, and bundle upload. |
| `internal/httpapi` | HTTP routes, auth middleware, request/response DTOs, admin APIs, agent APIs, docs API, and fallback routing. |
| `internal/control` | Product services and invariants for orgs, sites, deployments, agents, audit, maintenance, and runtime orchestration. |
| `internal/store` | Postgres store, embedded migrations, generated sqlc queries, and store wrappers. |
| `internal/deploy` | Bundle archive reading, manifest validation, safe capability policy, path policy, and artifact preparation. |
| `internal/runtime` | Per-site Goja runtime construction, supervisor, activation, restart, stop, runtime status, and host dispatch. |
| `internal/sitejs` | JavaScript-facing modules and HTTP bridge: `express`, `ui.dsl`, scoped database access, DB guard, sessions, and request/response helpers. |
| `internal/webadmin` | Embedded Vite dashboard serving and SPA fallback. |
| `web/admin` | React dashboard with RTK Query, React Router, MSW, Storybook, and OS1 styling. |
| `docs` | Stable contributor, architecture, and runbook documentation. |
| `ttmp` | Ticket workspaces, design docs, diaries, runbooks, and investigation artifacts. |

## Architecture overview

```text
Dashboard / CLI / Agent
  -> internal/httpapi
      -> internal/control
          -> internal/store
          -> internal/deploy
          -> internal/runtime.Supervisor
              -> internal/runtime.SiteRuntime
                  -> internal/sitejs modules
```

Public hosted traffic follows a separate path:

```text
Incoming HTTP request
  -> host fallback routing
  -> runtime supervisor lookup by Host header
  -> active SiteRuntime
  -> JavaScript route handler
```

Authorization and product invariants belong in `internal/control`. HTTP handlers and CLI commands should be thin adapters. Dashboard code should call the API through RTK Query and must not be the only enforcement point for server-side rules.

## Run locally

Start the daemon with the default development config:

```bash
go run ./cmd/go-go-hostd --config configs/dev.yaml
```

Check health and readiness:

```bash
curl -fsS http://127.0.0.1:8080/healthz
curl -fsS http://127.0.0.1:8080/readyz | jq .
```

Check the CLIs:

```bash
go run ./cmd/go-go-host status --api-url http://127.0.0.1:8080 --output table
go run ./cmd/go-go-host-agent status --api-url http://127.0.0.1:8080 --output json
```

Open the embedded dashboard:

- User dashboard: <http://127.0.0.1:8080/app>
- Platform admin dashboard: <http://127.0.0.1:8080/admin>

For the detailed local workflow, see [`docs/runbooks/local-development.md`](docs/runbooks/local-development.md).

## Local Postgres and Keycloak

Start local services:

```bash
docker compose -f deployments/dev/docker-compose.yaml up -d
```

Useful endpoints:

- Postgres: `postgres://go_go_host:go_go_host_dev@127.0.0.1:55432/go_go_host?sslmode=disable`
- Keycloak admin: `http://127.0.0.1:18080` with `admin` / `admin`

Run store/control tests against Postgres:

```bash
export GO_GO_HOST_TEST_DATABASE_URL='postgres://go_go_host:go_go_host_dev@127.0.0.1:55432/go_go_host?sslmode=disable'
go test ./internal/store ./internal/control -count=1
```

Available config files:

| Config | Purpose |
|---|---|
| `configs/dev.yaml` | Default development config. |
| `configs/dev.keycloak.yaml` | Keycloak-oriented local config. |
| `configs/dev.postgres-keycloak.yaml` | Local Postgres and Keycloak config. |
| `configs/production.example.yaml` | Production-shaped example config. |

## Dashboard development

Install and run the Vite dev server:

```bash
cd web/admin
pnpm install
pnpm dev
```

Build the dashboard:

```bash
pnpm build
```

Run Storybook:

```bash
pnpm storybook
```

Build Storybook:

```bash
pnpm storybook:build
```

Build and embed dashboard assets for the Go binary:

```bash
go run ./cmd/build-web
go test ./internal/webadmin
go build ./...
```

## Build and test

Common targets:

```bash
make test
make build
make web-build
make storybook-build
```

Direct commands:

```bash
go test ./...
go build ./...
```

For the validation matrix by change type, see [`docs/contributing/testing-and-validation.md`](docs/contributing/testing-and-validation.md).

## Deployment bundle model

A hosted site is deployed as an archive containing a `go-go-host.json` manifest, JavaScript scripts, and optional static assets. The validator checks archive paths, manifest fields, requested capabilities, quotas, and smoke behavior before a deployment can be activated.

Example bundle source:

- [`examples/hello-beta`](examples/hello-beta)

Core manifest fields include:

- `scriptsDir`
- `assetsDir`
- `entrypoint`
- `smokePath`
- `capabilities`
- `allowedPaths`
- `channel`

Deployment activation is separate from upload. Activation builds a new runtime and runs a health check before swapping live traffic.

## Hosted JavaScript capabilities

Hosted Goja sites receive explicit host-mediated capabilities only. The safe default model includes:

- `express` for route registration.
- `ui.dsl` for escaped HTML rendering.
- scoped `database` / `db` backed by the site's own SQLite database.
- limited `time` / `timer` support.
- static asset serving from the active deployment.
- DB guard visibility for per-site SQLite quota behavior.

Unrestricted `fs` is not a default hosted capability. `exec` must not be exposed to hosted v1 sites.

## API and dashboard docs

The daemon exposes human, admin, agent, docs, and runtime APIs under `/api/v1`. The route reference is maintained in [`docs/architecture/api-surface.md`](docs/architecture/api-surface.md).

The dashboard includes a docs area under org routes:

```text
/app/orgs/:orgId/docs
/app/orgs/:orgId/docs/:slug
```

The docs API serves embedded documentation from both the human CLI and agent CLI:

```text
GET /api/v1/docs
GET /api/v1/docs/{slug}
```

## CLI model

All human CLI verbs should follow the Glazed command pattern:

- command structs embed `*cmds.CommandDescription`;
- settings structs use `glazed` tags;
- flags and arguments are declared with `fields.New`;
- commands include Glazed output and command-settings sections;
- `RunIntoGlazeProcessor` emits stable rows through `types.NewRow`.

Agent CLI commands should preserve the signed-request model and must not use human credentials for machine deployments.

## Contributing

Start with:

- [`CONTRIBUTING.md`](CONTRIBUTING.md)
- [`docs/contributing/README.md`](docs/contributing/README.md)

The contributor docs cover:

- architecture and subsystem boundaries;
- backend service rules;
- runtime and deployment rules;
- dashboard workflow and OS1 styling;
- local development;
- testing and validation;
- docmgr ticket and diary workflow.

Do not introduce new runtime capabilities, auth semantics, deployment activation behavior, or dashboard design-system changes without reading the relevant contributor guide and adding tests.
