---
Title: "Rollback a Site"
Slug: "rollback-workflow"
Short: "Activate the previous validated deployment for a site."
Topics:
  - go-go-host
  - rollback
  - deployments
Commands:
  - rollback
  - deployments
Flags:
  - site-id
  - deployment-id
IsTopLevel: false
IsTemplate: false
ShowPerDefault: true
SectionType: Tutorial
---

This tutorial covers rollback for a site whose current deployment should be replaced by the previous validated or superseded deployment.

## Inspect deployments

List deployments for the site before rolling back. This shows which deployment is active and which previous versions are available.

```bash
go-go-host deployments list --site-id site_123 --dev-user alice --output table
```

Use JSON or YAML when you want the full manifest and validation report fields.

```bash
go-go-host deployments show --deployment-id dep_123 --dev-user alice --output yaml
```

## Roll back

Rollback activates the previous validated/superseded deployment. It does not mutate bundle contents, so every deployment remains immutable and auditable.

```bash
go-go-host rollback --site-id site_123 --dev-user alice
```

The command returns the deployment row that became active.

## Verify runtime state

After rollback, inspect runtime status and request the public host.

```bash
go-go-host site runtime --site-id site_123 --dev-user alice --output table
curl -H 'Host: hello.localhost' http://127.0.0.1:8080/
```

## Troubleshooting

| Problem | Cause | Solution |
| --- | --- | --- |
| Rollback fails with no rows | There is no previous validated/superseded deployment | Upload and validate another deployment before relying on rollback |
| Rollback activates but traffic does not change | The public request is using the wrong Host header | Use the site's `primary_host` value from `site list` |
| Runtime status reports failed | The previous deployment no longer loads successfully | Inspect daemon logs and deploy a fixed bundle |

## See Also

- `deploy-workflow`
- `getting-started`
