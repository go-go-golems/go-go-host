# Local development runbook

This runbook lists the standard local development commands for `go-go-host`.

Run commands from the repository root unless noted otherwise:

```bash
cd /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host
```

## Build and test

```bash
go test ./...
go build ./...
```

Run one package or test:

```bash
go test ./internal/control -run TestName -count=1
```

## Run the daemon

```bash
go run ./cmd/go-go-hostd --config configs/dev.yaml
```

Check health from another shell:

```bash
curl -fsS http://127.0.0.1:8080/healthz
curl -fsS http://127.0.0.1:8080/readyz | jq .
```

Check the CLIs:

```bash
go run ./cmd/go-go-host status --api-url http://127.0.0.1:8080 --output table
go run ./cmd/go-go-host-agent status --api-url http://127.0.0.1:8080 --output json
```

Prefer tmux for long-running local servers:

```bash
tmux new -s go-go-host-dev
# inside tmux
go run ./cmd/go-go-hostd --config configs/dev.yaml
```

Kill a stuck process by port:

```bash
lsof-who -p 8080 -k
```

## Start Postgres and Keycloak

```bash
docker compose -f deployments/dev/docker-compose.yaml up -d
```

Useful endpoints:

- Postgres: `postgres://go_go_host:go_go_host_dev@127.0.0.1:55432/go_go_host?sslmode=disable`
- Keycloak admin: `http://127.0.0.1:18080` with `admin` / `admin`

Run Postgres-backed tests:

```bash
export GO_GO_HOST_TEST_DATABASE_URL='postgres://go_go_host:go_go_host_dev@127.0.0.1:55432/go_go_host?sslmode=disable'
go test ./internal/store ./internal/control -count=1
```

Use this path when changing migrations, sqlc queries, memberships, agents, deployments, audit, or store-backed control services.

## Run devctl stack

The repository has a devctl plugin in `.devctl.yaml`.

```bash
devctl up --force
```

Verify identity in dev auth mode:

```bash
curl -fsS http://127.0.0.1:8080/api/v1/me | jq '{email:.user.email, platformAdmin}'
```

## Run dashboard dev server

```bash
cd web/admin
pnpm install
pnpm dev
```

Build before handoff:

```bash
pnpm build
```

## Run Storybook

```bash
cd web/admin
pnpm storybook
```

Build Storybook before handoff when pages, components, fixtures, or stories changed:

```bash
pnpm storybook:build
```

## Validate embedded dashboard

Production uses embedded Vite assets served by the Go binary. Validate this path when changing dashboard build, routes, base paths, or `internal/webadmin`:

```bash
go run ./cmd/build-web
go test ./internal/webadmin
go build ./...
```

## Common issues

| Symptom | Check |
|---|---|
| Port already in use | `lsof-who -p 8080 -k` or `lsof-who -p 5173 -k` |
| `/readyz` reports DB failure | `docker compose -f deployments/dev/docker-compose.yaml ps` and the configured DSN |
| Store tests fail to connect | `GO_GO_HOST_TEST_DATABASE_URL` and Postgres container status |
| Dashboard works in Vite but not embedded | Run `go run ./cmd/build-web` and `go test ./internal/webadmin` |
| Storybook data differs from app data | Compare `web/admin/src/services/types.ts`, MSW handlers, and backend DTOs |

## End-of-session checklist

- Stop or detach long-running processes.
- Record validation commands and results in the handoff.
- Update the ticket diary if the work is ticket-based.
- Run `docmgr doctor --ticket <TICKET-ID> --stale-after 30` for docmgr work.
