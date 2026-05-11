# go-go-host

`go-go-host` is the first v1 implementation of a Goja sites hosting platform. It will host small server-side JavaScript sites in Go, route public traffic by host name, store immutable deployments, support human and agent deployment workflows, and provide separate user and platform-admin dashboards.

## Current phase

This repository is in Phase 0 scaffold work:

- `go-go-hostd` starts the daemon and serves health/version endpoints.
- `go-go-host` is the human CLI and uses Glazed command structure.
- `go-go-host-agent` is the headless agent CLI and uses Glazed command structure.
- Control-plane services, runtime supervision, deployments, dashboards, and agent enrollment will be added in later phases.

## Run locally

```bash
go run ./cmd/go-go-hostd --config configs/dev.yaml
```

In another shell:

```bash
go run ./cmd/go-go-host status --api-url http://127.0.0.1:8080 --output table
go run ./cmd/go-go-host-agent status --api-url http://127.0.0.1:8080 --output json
```

## Build and test

```bash
make test
make build
```

## Optional local Postgres and Keycloak

Phase 1 introduces the control-plane schema, so Postgres belongs in Phase 1 development infrastructure. Keycloak is exercised in Phase 2 authentication work, but it is included in the dev compose file now so the local stack is ready before auth wiring begins.

```bash
docker compose -f deployments/dev/docker-compose.yaml up -d
```

Useful local endpoints:

- Postgres: `postgres://go_go_host:go_go_host_dev@127.0.0.1:55432/go_go_host?sslmode=disable`
- Keycloak admin: `http://127.0.0.1:18080` with `admin` / `admin`

The store layer is Postgres-first and generated with sqlc. Postgres integration tests run when `GO_GO_HOST_TEST_DATABASE_URL` is set:

```bash
export GO_GO_HOST_TEST_DATABASE_URL='postgres://go_go_host:go_go_host_dev@127.0.0.1:55432/go_go_host?sslmode=disable'
go test ./internal/store ./internal/control
```

The daemon still defaults to `configs/dev.yaml` during the Phase 1 skeleton. Use `configs/dev.postgres-keycloak.yaml` when wiring Postgres and OIDC in the next steps.

## Capability model

Hosted Goja sites should receive explicit host-mediated capabilities only. The safe default target is:

- `express` for route registration;
- `ui.dsl` for escaped HTML values;
- scoped `database` / `db` backed by the site's own SQLite database;
- limited `time` / `timer` support;
- static asset serving from the active deployment.

Unrestricted `fs` is not a default hosted capability. `exec` must not be exposed to hosted v1 sites.

## Dashboard model

The product has two dashboard surfaces:

- `/app/*`: user dashboard for organization users and developers.
- `/admin/*`: platform admin console for installation operators.

Both will be implemented as a React/Vite/RTK Query/Storybook frontend using `@go-go-golems/os-core`.

## CLI model

All CLI verbs should be Glazed commands:

- command structs embed `*cmds.CommandDescription`;
- settings structs use `glazed` tags;
- flags/arguments are declared with `fields.New`;
- commands include Glazed output and command-settings sections;
- `RunIntoGlazeProcessor` emits stable rows through `types.NewRow`.
