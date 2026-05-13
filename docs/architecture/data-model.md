# Data model

This document summarizes the main persistent entities in the `go-go-host` control plane. The authoritative schema is in `internal/store/migrations`, and generated queries come from `internal/store/queries`.

## Persistence stack

```text
internal/store/migrations/*.sql
  -> database schema
internal/store/queries/*.sql
  -> sqlc-generated query methods
internal/store/*.go
  -> store wrappers and domain model helpers
internal/control/*.go
  -> product behavior using store wrappers
```

Do not bypass this stack from HTTP handlers or dashboard code.

## Main entities

| Entity | Purpose |
|---|---|
| `users` | Human identities from dev auth or OIDC. |
| `orgs` | Tenant boundary for sites, agents, memberships, and audit. |
| `memberships` | User role in an organization. Roles drive view/deploy/owner permissions. |
| `platform_admins` | Installation-level operator access. |
| `sites` | Hosted application record. Owns primary host, active deployment, status, and org. |
| `site_domains` | Additional hostnames for a site and verification status. |
| `site_quotas` | Bundle size, DB size, request timeout, and related per-site limits. |
| `site_capabilities` | Per-site capability policy. |
| `deployments` | Immutable deployment versions, artifact paths, manifest JSON, validation JSON, status, and activation time. |
| `deploy_runs` | Short-lived agent upload authorization records. |
| `agents` | Machine identities scoped to an organization. |
| `agent_keys` | Public keys for agent signed requests. |
| `agent_site_grants` | Agent permissions for specific sites, channels, paths, and activation rights. |
| `agent_nonces` | Replay protection for signed agent requests. |
| `audit_log` | Security and operational event history. |

## Core relationships

```mermaid
erDiagram
  users ||--o{ memberships : has
  orgs ||--o{ memberships : contains
  orgs ||--o{ sites : owns
  orgs ||--o{ agents : owns
  orgs ||--o{ audit_log : records
  sites ||--o{ deployments : has
  sites ||--o{ site_domains : has
  sites ||--o{ site_capabilities : has
  sites ||--o{ site_quotas : has
  agents ||--o{ agent_keys : has
  agents ||--o{ agent_site_grants : has
  agents ||--o{ deploy_runs : creates
  sites ||--o{ agent_site_grants : grants
  sites ||--o{ deploy_runs : targets
```

## Status fields

Status fields are part of the product contract. Treat changes to status values as compatibility-affecting changes.

Common status areas:

- Site status: active or lifecycle-related site availability.
- Deployment status: uploaded, validated, rejected, active, superseded, and related states.
- Agent status: active or revoked.
- Agent key status: active or revoked.
- Deploy-run status: pending, uploading, completed, expired, revoked, or similar lifecycle states.
- Domain verification status.

When adding a status value, update:

1. Migration or constants if needed.
2. Store model conversion.
3. Control-service transitions.
4. API DTOs and frontend types.
5. Dashboard rendering and filters.
6. Tests for allowed and forbidden transitions.

## Migration workflow

For schema changes:

1. Add a new numbered migration under `internal/store/migrations`.
2. Keep migrations forward-only and deterministic.
3. Add or update SQL queries under `internal/store/queries`.
4. Regenerate sqlc output using the repository's configured sqlc workflow.
5. Add store wrapper methods.
6. Add tests with a live Postgres DSN.

Validation:

```bash
docker compose -f deployments/dev/docker-compose.yaml up -d
export GO_GO_HOST_TEST_DATABASE_URL='postgres://go_go_host:go_go_host_dev@127.0.0.1:55432/go_go_host?sslmode=disable'
go test ./internal/store ./internal/control -count=1
```

## Store wrapper rules

Store wrappers should:

- Accept `context.Context`.
- Return domain structs or clear primitive values.
- Hide sqlc-specific details from control services where practical.
- Keep transaction boundaries explicit.
- Avoid HTTP status concepts.
- Avoid dashboard-specific DTO fields.

Control services should not know more SQL than necessary. If a service needs a new data access pattern, add a store method.

## Audit model

Audit rows should capture actor, action, resource type, resource ID, org, and timestamp. Use audit for security-relevant and operationally important mutations.

Examples:

- `deployment.upload`
- `deployment.validation_failed`
- `deployment.activate`
- `deployment.rollback`
- `agent.enroll`
- `agent.grant.upsert`
- `agent.signature.invalid`
- `agent.upload_token.invalid`

Keep action names stable. Dashboards and operational docs may depend on them.
