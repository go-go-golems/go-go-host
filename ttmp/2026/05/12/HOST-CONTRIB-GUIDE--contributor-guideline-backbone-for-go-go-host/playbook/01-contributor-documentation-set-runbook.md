---
Title: Contributor documentation set runbook
Ticket: HOST-CONTRIB-GUIDE
Status: active
Topics:
    - go-go-host
    - contributor-guidelines
    - onboarding
    - documentation
DocType: playbook
Intent: long-term
Owners: []
RelatedFiles:
    - Path: AGENT.md
      Note: Current compact agent-facing rules to reconcile with repository-facing docs
    - Path: README.md
      Note: |-
        Should link to the future contributor documentation entrypoint
        Future link target for the contributing docs entrypoint
    - Path: docs/contributing/playbooks/os1-admin-dashboard-ui-work-guidelines.md
      Note: |-
        Existing stable playbook that the frontend guide should link rather than duplicate
        Stable dashboard playbook that future frontend docs should link rather than duplicate
ExternalSources: []
Summary: Step-by-step runbook for implementing the proposed go-go-host contributor documentation backbone.
LastUpdated: 2026-05-12T21:07:00-04:00
WhatFor: Use this as the execution checklist when turning the HOST-CONTRIB-GUIDE design into repository-facing docs.
WhenToUse: Before creating or reorganizing contribution, architecture, local-development, testing, or docmgr workflow documents.
---


# Contributor documentation set runbook

## Purpose

This playbook turns the design document in this ticket into concrete repository-facing documentation. It is intentionally command-oriented so a future human or coding agent can implement the documentation set without rediscovering the analysis.

## Environment assumptions

Work from the repository root:

```bash
cd /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host
```

Required tools:

- Go toolchain used by the repository.
- `pnpm` for `web/admin`.
- `docmgr` for ticket docs.
- `devctl` and Docker for local stack verification when updating runbooks.

## Phase 0: read before editing

Read these files first:

```bash
# Existing stable guidance
$EDITOR AGENT.md README.md docs/contributing/playbooks/os1-admin-dashboard-ui-work-guidelines.md

# Ticket design package
$EDITOR ttmp/2026/05/12/HOST-CONTRIB-GUIDE--contributor-guideline-backbone-for-go-go-host/design-doc/01-contributor-guideline-backbone-design.md
$EDITOR ttmp/2026/05/12/HOST-CONTRIB-GUIDE--contributor-guideline-backbone-for-go-go-host/reference/01-investigation-diary.md
```

Exit criteria:

- You understand that `ttmp` is working memory and `docs` is stable contributor guidance.
- You know the backend layering rule: `httpapi -> control -> store/runtime`.
- You know the dashboard playbook already exists and must not be duplicated.

## Phase 1: create the stable skeleton

Create directories:

```bash
mkdir -p docs/contributing docs/architecture docs/runbooks
```

Create or update these files:

```text
docs/contributing/README.md
docs/contributing/architecture-map.md
docs/contributing/backend-service-guidelines.md
docs/contributing/runtime-and-deployment-guidelines.md
docs/contributing/frontend-dashboard-guidelines.md
docs/contributing/testing-and-validation.md
docs/contributing/docmgr-and-ticket-workflow.md
docs/runbooks/local-development.md
docs/architecture/api-surface.md
docs/architecture/data-model.md
```

Also consider adding a tiny top-level `CONTRIBUTING.md`:

```markdown
# Contributing

Start with [docs/contributing/README.md](docs/contributing/README.md).
```

Update `README.md` with a contribution section:

```markdown
## Contributing

For architecture, local development, validation commands, and lane-specific guidance, see `docs/contributing/README.md`.
```

Exit criteria:

- A contributor can start at README and find the correct guide.
- Dashboard guidance links to the existing OS1 playbook.
- No stable doc is only a placeholder; each should contain at least purpose, core rules, file references, and validation commands.

## Phase 2: fill the docs from evidence

Use this source map:

| Stable doc | Primary evidence |
|---|---|
| `architecture-map.md` | `internal/httpapi/handler.go`, `internal/control/core.go`, `internal/runtime/supervisor.go`, `web/admin/src/app/routes.tsx` |
| `backend-service-guidelines.md` | `internal/control/*.go`, `internal/httpapi/*.go`, `internal/store/*.go`, `internal/store/queries/*.sql` |
| `runtime-and-deployment-guidelines.md` | `internal/deploy/bundle.go`, `internal/control/deployments.go`, `internal/runtime/runtime.go`, `internal/sitejs/web/host.go` |
| `frontend-dashboard-guidelines.md` | `docs/contributing/playbooks/os1-admin-dashboard-ui-work-guidelines.md`, `web/admin/src/services/goGoHostApi.ts`, `web/admin/src/services/msw`, `web/admin/src/pages/*/*.stories.tsx` |
| `testing-and-validation.md` | `Makefile`, `web/admin/package.json`, `.github/workflows/*`, existing tests |
| `docmgr-and-ticket-workflow.md` | this ticket, existing `ttmp/2026/05/*` tickets, docmgr conventions |
| `local-development.md` | `Makefile`, `.devctl.yaml`, `plugins/go-go-host-devctl.py`, `deployments/dev/docker-compose.yaml`, `configs/*.yaml` |
| `api-surface.md` | route table in `internal/httpapi/handler.go` and RTK endpoints in `web/admin/src/services/goGoHostApi.ts` |
| `data-model.md` | `internal/store/migrations/*.sql`, `internal/store/queries/*.sql`, `internal/store/*.go` |

Use line references while drafting. Example:

```bash
nl -ba internal/httpapi/handler.go | sed -n '17,110p'
nl -ba internal/control/core.go | sed -n '9,37p'
nl -ba internal/deploy/bundle.go | sed -n '19,180p'
```

Exit criteria:

- Every major claim points to source or existing docs.
- Every guide has at least one “Do not” section for dangerous shortcuts.
- Every guide has a validation checklist.

## Phase 3: validate commands

Run baseline validation:

```bash
go test ./...
go build ./...
```

For docs that mention frontend workflow, also run:

```bash
cd web/admin
pnpm build
pnpm storybook:build
cd ../..
```

For docs that mention embedded dashboard workflow, run:

```bash
go run ./cmd/build-web
go test ./internal/webadmin
```

For docs that mention docmgr workflow, run:

```bash
docmgr doctor --ticket HOST-CONTRIB-GUIDE --stale-after 30
```

Exit criteria:

- Commands in docs have been tested or clearly marked as requiring a service dependency.
- Any failing command is documented with the exact failure and either fixed or listed as an open issue.

## Phase 4: update docmgr bookkeeping

Relate important stable docs to this ticket as they are created:

```bash
docmgr doc relate \
  --doc ttmp/2026/05/12/HOST-CONTRIB-GUIDE--contributor-guideline-backbone-for-go-go-host/design-doc/01-contributor-guideline-backbone-design.md \
  --file-note "/home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/docs/contributing/README.md:Stable contribution entrypoint created from this design"
```

Update tasks:

```bash
docmgr task check --ticket HOST-CONTRIB-GUIDE --id 1,2,3,4
```

Update changelog:

```bash
docmgr changelog update --ticket HOST-CONTRIB-GUIDE \
  --entry "Implemented stable contributor documentation skeleton and linked it from README." \
  --file-note "/home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/docs/contributing/README.md:Contributor documentation entrypoint"
```

Exit criteria:

- `docmgr doctor --ticket HOST-CONTRIB-GUIDE --stale-after 30` passes.
- Ticket tasks and changelog reflect the implemented docs.

## Failure modes

### README and source disagree

If README says the system is still in an earlier phase, update README or explicitly link to fresher architecture docs. Do not copy stale README claims into new docs.

### Dashboard docs drift from OS1 playbook

If a new frontend guide starts restating visual rules, replace the duplicated text with a link to `docs/contributing/playbooks/os1-admin-dashboard-ui-work-guidelines.md` and keep only workflow-specific additions.

### API surface table becomes stale

Until route extraction is automated, add an explicit maintenance note to `docs/architecture/api-surface.md`: every change to `internal/httpapi/handler.go` that adds/removes a public route must update the API surface doc.

### Tests require services

If a validation command requires Postgres, Keycloak, or devctl, document prerequisites and provide a smaller local fallback where appropriate.

## Final handoff checklist

- [ ] `docs/contributing/README.md` exists and is linked from `README.md`.
- [ ] Architecture map has diagrams and file references.
- [ ] Backend guide documents `httpapi -> control -> store/runtime`.
- [ ] Runtime/deployment guide documents manifest, capabilities, validation, and activation.
- [ ] Frontend guide links to the OS1 dashboard playbook.
- [ ] Testing guide has a matrix by contribution lane.
- [ ] Local development runbook has copy/paste commands.
- [ ] API surface and data model docs exist.
- [ ] `docmgr doctor --ticket HOST-CONTRIB-GUIDE --stale-after 30` passes.
