---
Title: "Getting Started with go-go-host-agent"
Slug: "agent-getting-started"
Short: "Check agent CLI wiring, create a key, enroll, and deploy."
Topics:
  - go-go-host-agent
  - getting-started
Commands:
  - status
  - keygen
  - enroll
  - deploy
Flags:
  - api-url
  - config
IsTopLevel: true
IsTemplate: false
ShowPerDefault: true
SectionType: Tutorial
---

`go-go-host-agent` is the headless deploy CLI for CI workers and other machine identities. It uses a local Ed25519 key pair, a one-time enrollment token, signed deploy-run creation, and upload tokens scoped to a single deploy run.

Check daemon reachability:

```bash
go run ./cmd/go-go-host-agent status --api-url http://127.0.0.1:8080 --output json
```

Create local keys and enroll after a human operator gives you an `enroll_...` token:

```bash
go-go-host-agent keygen --config ./agent.json
go-go-host-agent enroll --config ./agent.json --token enroll_...
```

Deploy a prepared bundle to an allowed site/channel/path:

```bash
go-go-host-agent deploy \
  --config ./agent.json \
  --bundle ./site.tar.gz \
  --site-id site_123 \
  --channel default \
  --bundle-path bundles/site.tar.gz
```

Add `--activate` only when the human-created grant includes scoped auto-activation permission for that site.
