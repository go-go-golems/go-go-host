# Backend service guidelines

This document defines how to change the Go backend. It covers HTTP handlers, control services, store access, authorization, audit, and tests.

## Layering

Backend changes should follow this dependency direction:

```text
internal/httpapi
  -> internal/control
      -> internal/store
      -> internal/deploy
      -> internal/runtime
```

Responsibilities:

| Layer | Owns | Does not own |
|---|---|---|
| `internal/httpapi` | Routes, auth middleware, request decoding, response DTOs, status codes. | Product authorization rules as the only enforcement point. |
| `internal/control` | Product invariants, permissions, orchestration, audit, service methods. | Raw SQL or HTTP-specific DTO concerns. |
| `internal/store` | Migrations, sqlc queries, transactions, store wrappers, database model conversion. | HTTP request handling or runtime activation decisions. |
| `internal/deploy` | Bundle validation, manifest parsing, archive policy, capability policy. | Actor permissions. |
| `internal/runtime` | Runtime lifecycle and host dispatch. | Membership checks or dashboard concerns. |

## Adding an HTTP endpoint

Use this sequence:

1. Define the product behavior in an `internal/control` service.
2. Add or update store methods if durable state is required.
3. Add the HTTP handler in `internal/httpapi`.
4. Register the route in `internal/httpapi/handler.go`.
5. Add integration tests for allowed and forbidden cases.
6. Update dashboard/CLI/docs if the endpoint is user-facing.

Handler shape:

```go
func handleUpdateThing(core *control.Core) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        p, err := requirePrincipal(r)
        if err != nil {
            writeError(w, http.StatusUnauthorized, err.Error())
            return
        }

        var req updateThingRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            writeError(w, http.StatusBadRequest, "invalid JSON body")
            return
        }

        out, err := core.Things.Update(r.Context(), p.User.ID, req)
        if errors.Is(err, control.ErrPermissionDenied) {
            writeError(w, http.StatusForbidden, err.Error())
            return
        }
        if err != nil {
            writeError(w, http.StatusInternalServerError, err.Error())
            return
        }

        writeJSON(w, http.StatusOK, thingToDTO(out))
    }
}
```

Keep handlers small. If a handler becomes the place where permissions, database updates, runtime calls, and audit are coordinated, move that logic into `internal/control`.

## Authorization

Authorization belongs in control services. UI state and CLI prechecks can improve usability, but they are not enforcement.

Common patterns:

- Viewing org/site resources should check membership or platform-admin status as appropriate.
- Deployment upload should require deploy-capable human role or a validated agent deploy run.
- Activation should be treated separately from upload.
- Platform-admin APIs should not rely on route naming alone; check platform-admin status in the handler or service path used by that API.
- Agent actions must validate active agent, active key, signature, timestamp, nonce, grant, and token status.

## Audit events

Write audit events for security-relevant mutations and operationally important actions:

- Agent creation, enrollment, revocation, key use/revocation, grant changes.
- Deploy-run creation and upload-token failures where useful.
- Deployment upload, validation failure, activation, rollback, prune, export where applicable.
- Site config, domain, capability, and maintenance changes.
- Platform-admin actions.

Audit writes should be close to the service method that performs the mutation.

## Store and sqlc workflow

For schema or query changes:

1. Add a new migration in `internal/store/migrations`.
2. Add or update SQL in `internal/store/queries`.
3. Regenerate sqlc output using the repository's established workflow.
4. Add store wrapper methods in `internal/store`.
5. Add tests using a live Postgres DSN when behavior depends on the database.

Do not add raw SQL to HTTP handlers. If a query is needed by product logic, expose it through the store layer.

## Error handling

Use stable error mapping:

| Error | HTTP status |
|---|---|
| Missing/invalid authentication | `401 Unauthorized` |
| Authenticated but not allowed | `403 Forbidden` |
| Invalid user input | `400 Bad Request` |
| Missing resource | `404 Not Found` when the API intentionally exposes that distinction |
| Internal dependency failure | `500 Internal Server Error` or a more specific service status where appropriate |

Avoid leaking implementation details for security-sensitive failures. Agent signature failures, token failures, and permission failures should not reveal more than necessary.

## Tests

For a backend feature, include tests at the layer that owns the behavior:

- Control-service tests for product rules.
- HTTP integration tests for routes, auth, JSON, and status mapping.
- Store tests for migrations and query semantics.
- Runtime/deploy tests when the behavior crosses into deployment or Goja runtime.

Minimum commands:

```bash
go test ./...
go build ./...
```

For store-backed behavior:

```bash
export GO_GO_HOST_TEST_DATABASE_URL='postgres://go_go_host:go_go_host_dev@127.0.0.1:55432/go_go_host?sslmode=disable'
go test ./internal/store ./internal/control -count=1
```
