# API surface

This document summarizes the HTTP API surface registered by `internal/httpapi/handler.go`. Update it when adding, removing, or changing public routes.

## Route categories

| Category | Auth model | Route prefix |
|---|---|---|
| Health/config | Public or deployment-local | `/healthz`, `/readyz`, `/api/v1/version`, `/api/v1/config` |
| User dashboard API | Human authentication | `/api/v1/me`, `/api/v1/orgs`, `/api/v1/sites`, `/api/v1/deployments` |
| Platform admin API | Human authentication plus platform-admin authorization | `/api/v1/admin` |
| Agent API | Signed agent requests and upload tokens | `/api/v1/agent` |
| Embedded dashboard | Browser static/SPAs | `/app`, `/admin` |

## Health and config

| Method | Path | Purpose |
|---|---|---|
| `GET` | `/healthz` | Basic process health. |
| `GET` | `/readyz` | Readiness checks such as DB and data directory. |
| `GET` | `/api/v1/version` | Version metadata. |
| `GET` | `/api/v1/config` | Browser-facing configuration. |

## User/org/site API

| Method | Path | Purpose |
|---|---|---|
| `GET` | `/api/v1/me` | Current user, memberships, and platform-admin flag. |
| `GET` | `/api/v1/orgs` | List organizations visible to the user. |
| `POST` | `/api/v1/orgs` | Create an organization. |
| `GET` | `/api/v1/orgs/{org_id}/sites` | List sites in an organization. |
| `POST` | `/api/v1/orgs/{org_id}/sites` | Create a site. |
| `GET` | `/api/v1/orgs/{org_id}/agents` | List deployment agents for an org. |
| `POST` | `/api/v1/orgs/{org_id}/agents` | Create an agent and optional enrollment token/grant. |
| `POST` | `/api/v1/orgs/{org_id}/agents/{agent_id}/revoke` | Revoke an agent. |
| `POST` | `/api/v1/orgs/{org_id}/agents/{agent_id}/enrollment-token` | Create a new enrollment token. |
| `GET` | `/api/v1/orgs/{org_id}/agents/{agent_id}/keys` | List agent keys. |
| `POST` | `/api/v1/orgs/{org_id}/agents/{agent_id}/keys/{key_id}/revoke` | Revoke an agent key. |
| `POST` | `/api/v1/orgs/{org_id}/agents/{agent_id}/grants` | Upsert an agent site grant. |
| `GET` | `/api/v1/orgs/{org_id}/audit` | List org audit events. |

## Agent API

| Method | Path | Purpose |
|---|---|---|
| `POST` | `/api/v1/agent/enroll` | Exchange enrollment token and public key for an agent key. |
| `POST` | `/api/v1/agent/deploy-runs` | Create a signed deploy run and upload token. |
| `POST` | `/api/v1/agent/deploy-runs/{run_id}/upload` | Upload a bundle using a deploy-run upload token. |

## Site runtime, settings, and deployments

| Method | Path | Purpose |
|---|---|---|
| `GET` | `/api/v1/sites/{site_id}/runtime` | Get runtime status for a site. |
| `GET` | `/api/v1/sites/{site_id}/db/stats` | Get per-site DB guard stats. |
| `GET` | `/api/v1/sites/{site_id}/config` | List site config values. |
| `PUT` | `/api/v1/sites/{site_id}/config` | Upsert a site config value. |
| `DELETE` | `/api/v1/sites/{site_id}/config` | Delete a site config value. |
| `GET` | `/api/v1/sites/{site_id}/capabilities` | List site capability settings. |
| `PUT` | `/api/v1/sites/{site_id}/capabilities` | Upsert a site capability setting. |
| `GET` | `/api/v1/sites/{site_id}/domains` | List site domains. |
| `POST` | `/api/v1/sites/{site_id}/domains` | Add a site domain. |
| `POST` | `/api/v1/sites/{site_id}/domains/{domain_id}/verify` | Verify a site domain. |
| `DELETE` | `/api/v1/sites/{site_id}/domains/{domain_id}` | Delete a site domain. |
| `GET` | `/api/v1/sites/{site_id}/environment` | Return environment/secrets placeholder information. |
| `POST` | `/api/v1/sites/{site_id}/deployments` | Upload a deployment bundle as a human user. |
| `GET` | `/api/v1/sites/{site_id}/deployments` | List site deployments. |
| `POST` | `/api/v1/sites/{site_id}/rollback` | Roll back to a previous deployment. |
| `GET` | `/api/v1/sites/{site_id}/export/metadata` | Export site metadata. |
| `GET` | `/api/v1/sites/{site_id}/export/db` | Export site SQLite database. |
| `POST` | `/api/v1/sites/{site_id}/deployments/prune` | Prune old deployment artifacts. |
| `GET` | `/api/v1/deployments/{deployment_id}` | Get deployment details. |
| `GET` | `/api/v1/deployments/{deployment_id}/bundle` | Export a deployment bundle. |
| `POST` | `/api/v1/deployments/{deployment_id}/activate` | Activate a deployment. |

## Platform admin API

| Method | Path | Purpose |
|---|---|---|
| `GET` | `/api/v1/admin/runtimes/summary` | Runtime summary across the platform. |
| `POST` | `/api/v1/admin/runtimes/{site_id}/restart` | Restart a site runtime. |
| `POST` | `/api/v1/admin/runtimes/{site_id}/stop` | Stop a site runtime. |
| `GET` | `/api/v1/admin/orgs` | List org inventory. |
| `GET` | `/api/v1/admin/users` | List user inventory. |
| `GET` | `/api/v1/admin/sites` | List site inventory. |
| `GET` | `/api/v1/admin/deployments` | List deployments across orgs/sites. |
| `GET` | `/api/v1/admin/deployments/{deployment_id}` | Get admin deployment details. |
| `GET` | `/api/v1/admin/agents` | List agents. |
| `GET` | `/api/v1/admin/audit` | List audit events. |
| `GET` | `/api/v1/admin/quotas` | List quota state. |
| `GET` | `/api/v1/admin/capabilities` | List capability state. |
| `GET` | `/api/v1/admin/domains` | List domain state. |
| `POST` | `/api/v1/admin/audit/retention` | Apply audit retention. |

## Maintenance rule

When adding a route:

1. Register it in `internal/httpapi/handler.go`.
2. Add or update handler tests.
3. Add or update RTK Query endpoints if the dashboard uses it.
4. Update this document.
5. Confirm the auth category is correct.
