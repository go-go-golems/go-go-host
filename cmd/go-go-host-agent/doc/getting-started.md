---
Title: "Getting Started with go-go-host-agent"
Slug: "agent-getting-started"
Short: "Check agent CLI wiring before enrollment features land."
Topics:
  - go-go-host-agent
  - getting-started
Commands:
  - status
Flags:
  - api-url
IsTopLevel: true
IsTemplate: false
ShowPerDefault: true
SectionType: Tutorial
---

The phase-0 agent CLI is wired with Glazed help, logging, and output sections. It can check daemon status before signed enrollment and deployment flows are implemented.

```bash
go run ./cmd/go-go-host-agent status --api-url http://127.0.0.1:8080 --output json
```
