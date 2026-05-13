# Local development runbook

Local development should resemble the product enough to catch real integration mistakes, but it should also be fast enough that contributors use it every day. This runbook gives you the common loops: Go-only, daemon plus CLI, dashboard dev server, Storybook, embedded dashboard, and Postgres/Keycloak-backed development.

Run commands from the repository root unless a section says otherwise:

```bash
cd /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host
```

## 1. Fast Go loop

Use this loop when you are changing pure Go code or writing package tests:

```bash
go test ./...
go build ./...
```

To run one package or one test:

```bash
go test ./internal/control -run TestName -count=1
```

The `-count=1` flag disables test caching. Use it when you are iterating on behavior and want to know what just happened.

## 2. Run the daemon and CLIs

Start the daemon:

```bash
go run ./cmd/go-go-hostd --config configs/dev.yaml
```

In another shell, ask the human and agent CLIs for status:

```bash
go run ./cmd/go-go-host status --api-url http://127.0.0.1:8080 --output table
go run ./cmd/go-go-host-agent status --api-url http://127.0.0.1:8080 --output json
```

When running long-lived processes, prefer tmux so the server can be inspected and killed cleanly:

```bash
tmux new -s go-go-host-dev
# inside tmux
go run ./cmd/go-go-hostd --config configs/dev.yaml
```

If a process is stuck on a port, kill by port rather than guessing the PID:

```bash
lsof-who -p 8080 -k
```

## 3. Postgres and Keycloak services

Start the development services:

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

Use this path when you touch migrations, sqlc queries, memberships, agents, deployments, audit, or any control-plane state.

## 4. devctl loop

The repository has a devctl plugin wired through `.devctl.yaml`. Use it when you want the project-local orchestration path rather than manual commands:

```bash
devctl up --force
```

Then verify the daemon:

```bash
curl -fsS http://127.0.0.1:8080/healthz
curl -fsS http://127.0.0.1:8080/readyz | jq .
```

If you are working on platform-admin behavior in dev auth mode, verify identity:

```bash
curl -fsS http://127.0.0.1:8080/api/v1/me | jq '{email:.user.email, platformAdmin}'
```

## 5. Dashboard dev server

The dashboard lives in `web/admin` and is a Vite application.

Install dependencies if needed:

```bash
cd web/admin
pnpm install
```

Run the dev server:

```bash
pnpm dev
```

Open the dashboard through Vite while the daemon is running. The Vite config proxies API calls to the backend in development.

Before handing off dashboard code, build it:

```bash
pnpm build
```

## 6. Storybook loop

Storybook is the safest place to develop deterministic UI states because MSW can provide known data without requiring a live backend.

```bash
cd web/admin
pnpm storybook
```

Build Storybook before handoff when page/component stories changed:

```bash
pnpm storybook:build
```

A useful page story includes loading, error, empty, and populated states. The OS1 visual system is easier to preserve when those states are visible side by side.

## 7. Embedded dashboard loop

Production serves the dashboard from embedded Vite assets, not from the Vite dev server. When you change build, routing, base paths, or `internal/webadmin`, validate the embedded path:

```bash
go run ./cmd/build-web
go test ./internal/webadmin
go build ./...
```

If the browser works in Vite but not in the daemon, suspect base paths, embedded `dist`, or SPA fallback behavior.

## 8. Common failure modes

| Symptom | Likely cause | First check |
|---|---|---|
| `readyz` reports DB failure | Postgres is not running or DSN is wrong | `docker compose -f deployments/dev/docker-compose.yaml ps` |
| Dashboard route works in Vite but not embedded | SPA fallback or base path mismatch | `go run ./cmd/build-web` and `go test ./internal/webadmin` |
| Store tests skip or fail unexpectedly | `GO_GO_HOST_TEST_DATABASE_URL` is unset or points to a stale DB | Print the env var and restart compose services |
| Port already in use | Old daemon or Vite process is still running | `lsof-who -p 8080 -k` or `lsof-who -p 5173 -k` |
| Storybook state differs from app state | MSW fixture drifted from API DTO | Compare `web/admin/src/services/types.ts`, handlers, and backend DTOs |

## 9. Before you stop

A local session should end with enough evidence that the next person knows what happened. For code changes, record commands and outcomes in the handoff. For ticket work, update the diary and run:

```bash
docmgr doctor --ticket <TICKET-ID> --stale-after 30
```
