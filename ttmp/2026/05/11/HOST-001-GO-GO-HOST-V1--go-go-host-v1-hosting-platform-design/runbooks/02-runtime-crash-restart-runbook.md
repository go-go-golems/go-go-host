---
Title: Runtime Crash and Restart Runbook
Ticket: HOST-001-GO-GO-HOST-V1
DocType: playbook
Topics: [go-go-host, vm-runtime, operations]
Status: active
Intent: operational
---

# Runtime Crash and Restart Runbook

Use this when a hosted site returns errors, disappears from Host-header routing, or shows `failed`/`stopped` runtime state.

## 1. Triage

```bash
curl -fsS http://127.0.0.1:8080/readyz | jq .
go run ./cmd/go-go-host status --api-url http://127.0.0.1:8080
go run ./cmd/go-go-host audit list --dev-user dev-user --org-id ORG_ID --limit 50 --output json | jq '.[] | select(.resourceId=="SITE_ID" or .resourceType=="deployment")'
```

Check dashboard pages:

- `/app/orgs/{orgId}/sites/{siteId}/runtime`
- `/app/orgs/{orgId}/sites/{siteId}/deployments`
- `/admin/runtimes`

## 2. Restart one runtime

```bash
curl -X POST http://127.0.0.1:8080/api/v1/admin/runtimes/SITE_ID/restart \
  -H 'X-Go-Go-Host-User: dev-user'
```

Expected result: runtime status returns `ready` or the error explains the activation/load failure.

## 3. Stop a noisy runtime

```bash
curl -X POST http://127.0.0.1:8080/api/v1/admin/runtimes/SITE_ID/stop \
  -H 'X-Go-Go-Host-User: dev-user'
```

## 4. Roll back

```bash
go run ./cmd/go-go-host rollback --site-id SITE_ID --dev-user OWNER --output json
```

## 5. Export evidence before destructive work

```bash
go run ./cmd/go-go-host maintenance export metadata --site-id SITE_ID --dev-user OWNER -o /tmp/site-metadata.json
go run ./cmd/go-go-host maintenance export db --site-id SITE_ID --dev-user OWNER -o /tmp/site.sqlite
```

## 6. Escalation checklist

- Preserve deployment bundle SHA256 and deployment ID.
- Preserve runtime `lastError` and audit events.
- Do not edit site SQLite files while a runtime is running unless the runtime is stopped.
- If a bundle is malicious or corrupt, revoke agent keys and prune only after exporting evidence.
