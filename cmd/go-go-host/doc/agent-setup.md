---
Title: "Agent Setup"
Slug: "agent-setup"
Short: "Create deployment agents, grants, and enrollment tokens for machine deploys."
Topics:
  - go-go-host
  - agents
  - deployments
Commands:
  - agents
  - audit
Flags:
  - dev-user
  - org-id
  - site-id
IsTopLevel: false
IsTemplate: false
ShowPerDefault: true
SectionType: GeneralTopic
---

Agent deployment is implemented in v1. A human operator creates the agent and site grant with `go-go-host agents create`; the machine then uses `go-go-host-agent keygen`, `enroll`, and `deploy`.

For the complete operator guide, run:

```bash
go-go-host help agent-guide
```

For the machine-side guide, run:

```bash
go-go-host-agent help agent-guide
go-go-host-agent help agent-keygen-enroll-deploy
go-go-host-agent help agent-signature-troubleshooting
```

The shortest local flow is:

```bash
go-go-host agents create \
  --api-url http://127.0.0.1:8080 \
  --dev-user alice \
  --org-id ORG_ID \
  --name ci-agent \
  --site-id SITE_ID \
  --channel default \
  --bundle-path '**' \
  --can-activate \
  --output json
```

Then on the machine:

```bash
go-go-host-agent keygen --config ./agent.json --api-url http://127.0.0.1:8080
go-go-host-agent enroll --config ./agent.json --api-url http://127.0.0.1:8080 --token ENROLLMENT_TOKEN
go-go-host-agent deploy --config ./agent.json --bundle ./site.tar.gz --site-id SITE_ID --channel default --bundle-path bundles/site.tar.gz --activate
```
