---
Title: Investigation diary
Ticket: HOST-005-E2E-FIXES
Status: active
Topics:
    - go-go-host
    - hosting
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: glazed/pkg/cmds/fields/cobra.go
      Note: Required check happens while gathering Cobra flags based on cmd.Flags().Changed before env values are applied.
    - Path: go-go-host/cmd/go-go-host-agent/cmds/support.go
      Note: Headless agent CLI Glazed parser configuration now enables GO_GO_HOST_AGENT_* env overrides.
    - Path: go-go-host/cmd/go-go-host/cmds/support.go
      Note: Human CLI Glazed parser configuration now enables GO_GO_HOST_* env overrides.
    - Path: go-go-host/ttmp/2026/05/12/HOST-005-E2E-FIXES--go-go-host-end-to-end-bringup-and-smaller-fixes/scripts/01-reproduce-glazed-required-env.sh
      Note: Standalone repro confirming required Glazed fields fail before env middleware can satisfy them.
ExternalSources: []
Summary: Chronological notes for the end-to-end go-go-host bringup ticket and smaller CLI fixes.
LastUpdated: 2026-05-12T16:50:00-04:00
WhatFor: Record what was investigated and changed while making go-go-host work end to end.
WhenToUse: Resume HOST-005 work, validate CLI/env behavior, or review smaller follow-up issues.
---



# Investigation diary

## Goal

Bring go-go-host to a working end-to-end state and use this ticket to track smaller correctness and ergonomics fixes discovered along the way.

## Context

The first reported issue was Glazed CLI environment-variable parsing: `go-go-host me --print-parsed-fields` did not pick up `GO_GO_HOST_DEV_USER=...`, so local dev workflows had to pass `--dev-user` explicitly. The same expectation applies to agent-related CLI commands.

## Diary

### 2026-05-12 — Ticket creation and Glazed env parsing

- Created `HOST-005-E2E-FIXES` with tasks for human CLI env support, agent CLI env support, and the larger end-to-end smoke path.
- Confirmed the human CLI's `BuildGlazedCobraCommand` supplied `MiddlewaresFunc: glazedcli.CobraCommandDefaultMiddlewares`.
- Checked Glazed's `CobraParserConfig`: a custom `MiddlewaresFunc` replaces the built-in parser chain, so `AppName`-based env loading never runs unless the caller re-adds env middleware manually.
- Changed the human CLI builder to leave `MiddlewaresFunc` unset and set `AppName: "GO_GO_HOST"` so fields such as `dev-user` map to `GO_GO_HOST_DEV_USER`.
- Changed the headless agent CLI builder to leave `MiddlewaresFunc` unset and set `AppName: "GO_GO_HOST_AGENT"` so fields such as `api-url` map to `GO_GO_HOST_AGENT_API_URL`.
- Validation:
  - `GO_GO_HOST_DEV_USER=alice go run ./cmd/go-go-host me --print-parsed-fields` now shows `default.dev-user` sourced from env with `env_key: GO_GO_HOST_DEV_USER`.
  - `GO_GO_HOST_DEV_USER=alice go run ./cmd/go-go-host agents list --org-id org_123 --print-parsed-fields` shows the same env provenance for the `agents` verb.
  - `GO_GO_HOST_AGENT_API_URL=http://example.invalid go run ./cmd/go-go-host-agent status --print-parsed-fields` shows `default.api-url` sourced from env with `env_key: GO_GO_HOST_AGENT_API_URL`.
  - `go test ./cmd/go-go-host/... ./cmd/go-go-host-agent/...` passes.
  - `go test ./...` passes.

### 2026-05-12 — Reproduced required-field env issue in Glazed

- Added `scripts/01-reproduce-glazed-required-env.sh` to create a temporary Go module against the local `glazed` checkout.
- The repro confirms that AppName env parsing itself works for an optional field (`REQ_ENV_TEST_OPTIONAL_NAME=from-env`).
- The same setup fails when the field is marked required (`fields.WithRequired(true)`) and only provided by env: `Field required-name is required`.
- The likely root cause is in Glazed's Cobra source: `GatherFlagsFromCobraCommand` checks `cmd.Flags().Changed(flagName)` and returns the required-field error before later env middleware can merge `REQ_ENV_TEST_REQUIRED_NAME`.
- Found existing upstream Glazed issue [go-go-golems/glazed#556](https://github.com/go-go-golems/glazed/issues/556), which already tracks this exact bug.
- Added a detailed go-go-host-specific reproduction comment: https://github.com/go-go-golems/glazed/issues/556#issuecomment-4432888793

## Quick Reference

Use env variables instead of repeating common flags:

```bash
GO_GO_HOST_DEV_USER=alice go-go-host me
GO_GO_HOST_DEV_USER=alice go-go-host agents list --org-id org_123
GO_GO_HOST_AGENT_API_URL=http://127.0.0.1:8080 go-go-host-agent status
```

For diagnostics, append `--print-parsed-fields` and confirm the field log contains `source: env` and the expected `env_key`.

## Related

- Human CLI command builder: `cmd/go-go-host/cmds/support.go`
- Agent CLI command builder: `cmd/go-go-host-agent/cmds/support.go`
