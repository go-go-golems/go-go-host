---
Title: "Agent keygen, enroll, and deploy"
Slug: "agent-keygen-enroll-deploy"
Short: "Create Ed25519 agent keys, enroll with a one-time token, and upload signed deployments."
Topics:
  - go-go-host-agent
  - deployment
Commands:
  - keygen
  - enroll
  - deploy
  - status
Flags:
  - config
  - token
  - bundle
  - site-id
IsTopLevel: true
IsTemplate: false
ShowPerDefault: true
SectionType: Tutorial
---

The agent CLI uses an Ed25519 key pair stored in a local config file. A human user first creates an agent enrollment token with `go-go-host agents create`. The agent exchanges that one-time token for a registered key, then signs deploy-run creation requests.

```bash
go-go-host-agent keygen --config ./agent.json

go-go-host-agent enroll \
  --config ./agent.json \
  --token enroll_...

go-go-host-agent status --config ./agent.json --output table

go-go-host-agent deploy \
  --config ./agent.json \
  --bundle ./site.tar.gz \
  --site-id site_123 \
  --channel default \
  --path bundles/site.tar.gz
```

Signed requests include `X-Go-Go-Agent-ID`, `X-Go-Go-Agent-Key-ID`, `X-Go-Go-Agent-Timestamp`, `X-Go-Go-Agent-Nonce`, and `X-Go-Go-Agent-Signature`. The server rejects bad signatures, old/future timestamps, replayed nonces, revoked agents/keys, and deploy runs outside the agent's site/channel/path grant.
