---
Title: "Agent Setup Preview"
Slug: "agent-setup"
Short: "Understand the planned agent deployment workflow and current CLI status."
Topics:
  - go-go-host
  - agents
  - deployments
Commands:
  - go-go-host-agent
  - deploy
Flags:
  - dev-user
IsTopLevel: false
IsTemplate: false
ShowPerDefault: true
SectionType: GeneralTopic
---

This page explains the intended agent workflow for go-go-host. The v1 control-plane foundations already include agent tables, but the human CLI agent commands are intentionally deferred until the agent APIs are implemented.

## Current status

Use human deployment commands for now:

```bash
go-go-host deploy --site-id site_123 --path ./bundle.tar.gz --dev-user alice
go-go-host deployments activate --deployment-id dep_123 --dev-user alice
```

The separate `go-go-host-agent` binary exists as a scaffold for future non-human deployment flows.

## Planned flow

The planned agent workflow will let an org owner create an agent, grant scoped site/channel/path permissions, and hand the agent a one-time enrollment command. Agents will upload bundles through constrained deploy runs rather than broad user sessions.

The deployment validator already contains path and channel policy hooks so the later agent API can reuse the same checks.

## Troubleshooting

| Problem | Cause | Solution |
| --- | --- | --- |
| `go-go-host agents list` is not available | Agent HTTP APIs are not implemented yet | Use human deployment commands in v1 until Phase 9 |
| Agent binary has limited commands | It is currently a scaffold | Track the agent API phase before depending on it |
| Need path/channel restrictions today | Bundle validation supports policy hooks, but no user-facing grant editor exists yet | Keep deployments manual or add policy plumbing in the next phase |

## See Also

- `deploy-workflow`
- `rollback-workflow`
