---
Title: Local admin dashboard runbook
Ticket: HOST-003-ADMIN-DASHBOARD
Status: active
Topics:
  - dashboard
  - go-go-host
  - platform-admin
DocType: playbook
Intent: operational
Summary: "How to run and verify the platform-admin dashboard locally with devctl."
LastUpdated: 2026-05-11T23:35:00-04:00
---

# Local admin dashboard runbook

## Start stack

```bash
devctl up --force
```

The dev config seeds `dev-user` as a platform admin via:

```yaml
devPlatformAdminSubjects:
  - dev-user
```

## Verify identity

```bash
curl -fsS http://127.0.0.1:8080/api/v1/me | jq '{email:.user.email, platformAdmin}'
```

Expected:

```json
{
  "email": "dev-user@dev.local",
  "platformAdmin": true
}
```

## Verify admin APIs

```bash
for p in orgs users sites deployments agents audit quotas capabilities domains; do
  echo "$p"
  curl -fsS "http://127.0.0.1:8080/api/v1/admin/$p" | jq 'length'
done
```

## Browser URLs

- Embedded admin overview: <http://127.0.0.1:8080/admin/overview>
- Runtime operations: <http://127.0.0.1:8080/admin/runtimes>
- Inventory: <http://127.0.0.1:8080/admin/sites>
- Quotas: <http://127.0.0.1:8080/admin/quotas>
- Capabilities: <http://127.0.0.1:8080/admin/capabilities>
- Domains: <http://127.0.0.1:8080/admin/domains>

## Validation before committing

```bash
go test ./...
make web-build
make storybook-build
go run ./cmd/build-web
docmgr doctor --ticket HOST-003-ADMIN-DASHBOARD --stale-after 30
```
