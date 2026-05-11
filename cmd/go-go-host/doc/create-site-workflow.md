---
Title: "Create an Organization and Site"
Slug: "create-site-workflow"
Short: "Use the CLI to create org and site records before deploying bundles."
Topics:
  - go-go-host
  - sites
  - organizations
Commands:
  - org
  - site
Flags:
  - org-id
  - slug
  - name
IsTopLevel: false
IsTemplate: false
ShowPerDefault: true
SectionType: Tutorial
---

This tutorial covers the setup steps before uploading a deployment bundle. A site belongs to an organization and receives a generated primary host used by the runtime router.

## Create an organization

Create an organization as the current user. The creator becomes an org owner.

```bash
go-go-host org create --slug demo --name "Demo Org" --dev-user alice --output json
```

The plural alias also works:

```bash
go-go-host orgs list --dev-user alice --output table
```

## Create a site

Use the organization ID from the previous command to create a site.

```bash
go-go-host site create \
  --org-id org_123 \
  --slug hello \
  --name "Hello Site" \
  --dev-user alice \
  --output json
```

The site row contains `primary_host`, which is the Host header used when requesting the activated runtime.

## Troubleshooting

| Problem | Cause | Solution |
| --- | --- | --- |
| Site creation is forbidden | The current user is not an owner or developer for the org | Use the right `--dev-user` or membership |
| Slug validation fails | Slugs must be lowercase letters, numbers, and hyphens | Rename the org/site with a simple DNS-style slug |
| Host-header requests return 404 | No deployment is active or the Host header is wrong | Activate a deployment and use the site's `primary_host` |

## See Also

- `login-and-config`
- `deploy-workflow`
