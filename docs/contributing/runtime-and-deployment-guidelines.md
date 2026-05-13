# Runtime and deployment guidelines

This document defines how to work on deployment bundles, manifest validation, hosted JavaScript capabilities, runtime activation, and public request handling.

## Key files

| Area | Files |
|---|---|
| Bundle manifest and validation | `internal/deploy/bundle.go` |
| Deployment upload, dry-run, activation, rollback | `internal/control/deployments.go` |
| Runtime construction | `internal/runtime/runtime.go` |
| Runtime supervision and host dispatch | `internal/runtime/supervisor.go` |
| JavaScript HTTP bridge | `internal/sitejs/web` |
| UI DSL module | `internal/sitejs/uidsl` |
| Database guard | `internal/sitejs/dbguard` |

## Deployment lifecycle

A deployment moves through these steps:

```text
bundle upload
  -> archive/path/manifest/capability validation
  -> immutable deployment row and artifact paths
  -> unpack bundle
  -> dry-run SiteRuntime
  -> smoke health check
  -> validated or rejected deployment status
  -> optional activation
  -> live runtime traffic swap
```

Do not treat upload and activation as the same operation. Upload creates and validates a candidate deployment. Activation changes live traffic.

## Manifest contract

The bundle manifest is `go-go-host.json`. It describes the app entrypoint, script and asset locations, smoke path, requested capabilities, path policy, and channel. The manifest is parsed and validated in `internal/deploy/bundle.go`.

Rules:

- Required paths must be relative bundle paths.
- Paths must not escape the unpacked bundle directory.
- Requested capabilities must be allowed by platform/site policy.
- The server decides effective capabilities; the manifest only requests them.
- The smoke path should be stable and safe to call during validation.

## Capability policy

Hosted JavaScript must receive explicit host-mediated capabilities only. Safe default capabilities include route registration, UI rendering, scoped database access, timers where allowed, and static assets. Do not expose unrestricted filesystem, process execution, host environment, or arbitrary network access by default.

When adding a capability:

1. Define the capability name and threat model.
2. Add validation in `internal/deploy`.
3. Add site policy support if it can be enabled/disabled per site.
4. Wire the module or registrar in `internal/runtime` or `internal/sitejs`.
5. Add rejected and accepted tests.
6. Update JS API docs and dashboard capability/admin docs.

Do not add a module directly to every runtime without a policy decision.

## Runtime construction

`SiteRuntime` owns one active Goja runtime for one site deployment. Runtime construction should:

- Open the per-site SQLite database.
- Configure the database guard.
- Register only approved native modules and middleware.
- Register `express`, `ui.dsl`, and other host modules through explicit registrars.
- Mount static assets only when assets are configured and allowed.
- Load scripts through the runtime owner.
- Close all resources if construction fails.

The Goja runtime is not a general-purpose sandbox boundary. Treat every new host API as an expansion of what hosted code can do.

## Activation and supervisor behavior

`Supervisor.Activate` must preserve the live-traffic safety sequence:

```text
set status starting
build next runtime
run health check
lock supervisor maps
remove old host mappings
install new site and host mappings
set status ready
unlock
persist status
close old runtime asynchronously
```

Do not swap traffic before the new runtime is built and health-checked. Do not delete old mappings until the replacement is ready.

## Public request handling

Public hosted-site requests are routed by host name to the active runtime. The request path is:

```text
HTTP fallback
  -> Supervisor.GetByHost
  -> SiteRuntime.ServeHTTP
  -> sitejs/web.Host
  -> route registry or static asset mount
  -> Goja handler call
```

Public request handling should not perform control-plane membership checks. The control plane already decided which deployment is active.

## Required tests

Deployment/runtime changes should include both rejection and success cases.

Common required cases:

- Missing manifest is rejected.
- Invalid JSON manifest is rejected.
- Unsafe paths are rejected.
- Unknown or disabled capabilities are rejected.
- Valid bundle is unpacked and dry-run successfully.
- Runtime health-check failure marks the deployment invalid or prevents activation.
- Activation preserves old runtime when new runtime fails.
- Restart and stop behavior update status correctly.

Useful commands:

```bash
go test ./internal/deploy ./internal/runtime ./internal/sitejs/... -count=1
go test ./internal/control ./internal/httpapi -run 'Deployment|Runtime|Agent' -count=1
```

## Review checklist

Before merging runtime or deployment work, verify:

- The final enforcement point is server-side.
- The change does not rely on dashboard or CLI validation only.
- Capabilities are explicit and documented.
- Path validation cannot be bypassed through archive layout or manifest fields.
- Runtime construction and failure paths close resources.
- Activation does not swap traffic before a successful health check.
- Audit events exist for upload, validation failure, activation, and rollback where applicable.
