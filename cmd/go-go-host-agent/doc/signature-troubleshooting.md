---
Title: "Troubleshooting agent signatures"
Slug: "agent-signature-troubleshooting"
Short: "Common causes for signed agent request failures."
Topics:
  - go-go-host-agent
  - troubleshooting
Commands:
  - deploy
  - enroll
Flags:
  - config
IsTopLevel: false
IsTemplate: false
ShowPerDefault: true
SectionType: GeneralTopic
---

Agent deploy requests are intentionally strict. Check these items when a deploy run is denied:

- The config file has `agentId`, `keyId`, and `privateKey` from a successful `go-go-host-agent enroll`.
- The machine clock is within five minutes of the daemon clock.
- Each signed request uses a fresh nonce. Reusing the same request or replaying captured headers is rejected.
- The agent and key are still active and have not been revoked.
- The deploy run's `siteId`, `channel`, and logical `path` match the agent site grant.
- The upload token from deploy-run creation is used only with the matching upload endpoint.

Use JSON output for easier debugging:

```bash
go-go-host-agent status --config ./agent.json --output json
go-go-host-agent deploy --config ./agent.json --bundle ./site.tar.gz --output json
```
