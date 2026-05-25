# Changelog

## 2026-05-12

- Initial workspace created


## 2026-05-12

Enabled Glazed AppName-based environment-variable parsing for go-go-host and go-go-host-agent CLI commands; verified with --print-parsed-fields and cmd package tests.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/cmd/go-go-host-agent/cmds/support.go — Sets AppName GO_GO_HOST_AGENT instead of overriding the middleware chain.
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/cmd/go-go-host/cmds/support.go — Sets AppName GO_GO_HOST instead of overriding the middleware chain.


## 2026-05-12

Ran full repository test suite after CLI env parsing change; go test ./... passes.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/cmd/go-go-host-agent/cmds/support.go — Agent CLI env parser change included in full test run.
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/cmd/go-go-host/cmds/support.go — Human CLI env parser change included in full test run.


## 2026-05-12

Added a ticket script that reproduces the Glazed required-field/env ordering issue: optional env fields parse, required env-only fields fail with Field required-name is required.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/glazed/pkg/cmds/fields/cobra.go — Code path implicated by the repro.
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/ttmp/2026/05/12/HOST-005-E2E-FIXES--go-go-host-end-to-end-bringup-and-smaller-fixes/scripts/01-reproduce-glazed-required-env.sh — Executable repro script.


## 2026-05-12

Found existing upstream Glazed issue #556 for required-field env/config validation ordering and added a detailed go-go-host repro comment with the ticket script output and likely code path.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/glazed/pkg/cmds/fields/cobra.go — Likely early required-check code path referenced upstream.
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/ttmp/2026/05/12/HOST-005-E2E-FIXES--go-go-host-end-to-end-bringup-and-smaller-fixes/scripts/01-reproduce-glazed-required-env.sh — Repro referenced from the upstream issue comment.

