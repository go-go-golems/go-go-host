---
Title: "Deploy a Site Bundle"
Slug: "deploy-workflow"
Short: "Create a site, upload a bundle, activate it, and inspect runtime status."
Topics:
  - go-go-host
  - deployments
  - workflow
Commands:
  - deploy
  - deployments
  - site
Flags:
  - site-id
  - path
  - dev-user
IsTopLevel: true
IsTemplate: false
ShowPerDefault: true
SectionType: Tutorial
---

This tutorial covers the normal local developer deployment loop. It assumes the daemon is running in dev-auth mode and that you have already created an organization.

## Create a site

Create a site in an organization. The command returns a stable `site_id` and a primary host that will route to the runtime after activation.

```bash
go-go-host site create \
  --org-id org_123 \
  --slug hello \
  --name "Hello Site" \
  --dev-user alice \
  --output json
```

## Prepare a bundle

A deployment bundle is a `.tar.gz` or `.zip` archive containing a `go-go-host.json` manifest and site files. A minimal manifest looks like this:

```json
{
  "scriptsDir": "scripts",
  "assetsDir": "assets",
  "smokePath": "/",
  "capabilities": ["time", "timer"]
}
```

The validator rejects absolute paths, parent traversal, unsafe links, oversized bundles, and forbidden capabilities before a deployment can become active.

## Upload and validate

Upload the bundle to create a deployment record. The upload performs manifest validation and a dry-run runtime load.

```bash
go-go-host deploy \
  --site-id site_123 \
  --path ./hello.tar.gz \
  --message "initial deploy" \
  --dev-user alice \
  --output yaml
```

Use YAML or JSON output when reviewing validation reports, because reports contain nested fields.

## Activate and verify

Activate a validated deployment, then inspect runtime state.

```bash
go-go-host deployments activate --deployment-id dep_123 --dev-user alice
go-go-host site runtime --site-id site_123 --dev-user alice --output table
```

After activation, requests with the site's Host header route through the supervisor to the hosted Goja runtime.

## Troubleshooting

| Problem | Cause | Solution |
| --- | --- | --- |
| Upload returns a validation report | The archive failed manifest, path, quota, capability, or dry-run checks | Re-run with `--output yaml` and fix the reported errors |
| Activation fails | The deployment is not in an activatable status or the runtime failed to load | Inspect `deployments show` and server logs |
| Runtime status is stopped after daemon restart | Runtime state is in-process and stale records are reconciled on startup | Reactivate the desired deployment |

## See Also

- `getting-started`
- `rollback-workflow`
- `agent-setup`
