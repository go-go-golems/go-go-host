---
Title: "Getting Started with go-go-host"
Slug: "getting-started"
Short: "Start the daemon and check it from the Glazed CLI."
Topics:
  - go-go-host
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

Start the phase-0 daemon with a local development config:

```bash
go run ./cmd/go-go-hostd --config configs/dev.yaml
```

In another terminal, check daemon health through the Glazed CLI:

```bash
go run ./cmd/go-go-host status --api-url http://127.0.0.1:8080 --output table
go run ./cmd/go-go-host status --api-url http://127.0.0.1:8080 --output json
```

The CLI emits structured rows so scripts can select table, JSON, YAML, or other Glazed outputs.
