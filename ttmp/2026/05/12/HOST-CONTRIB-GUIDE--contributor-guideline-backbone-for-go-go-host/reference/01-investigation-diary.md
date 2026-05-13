---
Title: Investigation diary
Ticket: HOST-CONTRIB-GUIDE
Status: active
Topics:
    - go-go-host
    - contributor-guidelines
    - onboarding
    - documentation
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: ttmp/2026/05/12/HOST-CONTRIB-GUIDE--contributor-guideline-backbone-for-go-go-host/design-doc/01-contributor-guideline-backbone-design.md
      Note: Primary design deliverable produced during the investigation
    - Path: ttmp/2026/05/12/HOST-CONTRIB-GUIDE--contributor-guideline-backbone-for-go-go-host/playbook/01-contributor-documentation-set-runbook.md
      Note: Implementation runbook produced during the investigation
ExternalSources: []
Summary: Chronological diary for the contributor guideline backbone investigation.
LastUpdated: 2026-05-12T20:55:00-04:00
WhatFor: Use this to continue or audit the investigation that produced the contributor documentation proposal.
WhenToUse: Read before revising the proposed contribution guidelines, onboarding docs, or documentation map.
---


# Diary

## Goal

Capture the repository investigation used to design a strong contribution-guideline backbone for `go-go-host`, including architecture mapping, existing docs, comparable project conventions, and ticket delivery steps.

## Step 1: Create the documentation ticket and orient on the repository

I created a new docmgr ticket for the contributor-guideline work and immediately added the primary design document plus this diary. The intent was to keep the research trail and final recommendation in a ticket workspace instead of producing an ad-hoc markdown file.

### Prompt Context

**User prompt (verbatim):** "Create a new docmgr to put together great contribution guidelines for go-go-host. We are trying to put together the best documents for contributing / guidelines for developing for go-go-host/. We already have one for the dashboard. Look around (since we also haev a whole set of other projects of a similar style), look also in ttmp/ folders which often have runbooks or similar. 

Keep a diary as you look around.

CCreate  a detailed analysis / design / implementation guide that is for a new intern, explaining all the parts of the system needed to understand what it is, with prose paragraphs and bullet point sand pseudocode and diagrams and api references and file references. It should be very clear and technical. Store in the ticket and the nupload to remarkable.

This should be a suggestion for a set of guideline and documentation for people contributing to get them started, to ensure proper guidelines are followed, to make it easy to have a unified codebase as more work is done. A lot of humans and agents are going to work on this and they need a solid backbone of documentation."

**Assistant interpretation:** Create a new docmgr ticket containing an intern-ready contributor-guideline design package, inspect existing docs and similar projects, maintain a diary, validate the ticket, and upload the result to reMarkable.

**Inferred user intent:** Establish a durable documentation backbone so future human and agent contributors can work consistently across backend, frontend, runtime, deployment, dashboard, and operations areas.

### What I did
- Ran `docmgr status --summary-only` from the workspace root to confirm the docs root and ticket inventory.
- Created ticket `HOST-CONTRIB-GUIDE` with topics `go-go-host,contributor-guidelines,onboarding,documentation`.
- Added the design document `design-doc/01-contributor-guideline-backbone-design.md`.
- Added the diary `reference/01-investigation-diary.md`.
- Added four ticket tasks covering architecture mapping, comparable-project inspection, guide writing, validation, and upload.

### Why
- The user explicitly requested a new docmgr workspace and a diary.
- A ticket workspace makes the recommendation searchable through docmgr and keeps the future documentation implementation plan close to the evidence.

### What worked
- `docmgr status --summary-only` identified the active docs root as `go-go-host/ttmp`.
- Ticket and document creation succeeded.

### What didn't work
- I ran `find ttmp/HOST-CONTRIB-GUIDE -maxdepth 3 -type f` after ticket creation, but docmgr stores tickets under dated paths, so the command failed with `find: ‘ttmp/HOST-CONTRIB-GUIDE’: No such file or directory`.

### What I learned
- This repository already has a substantial `ttmp/2026/05/...` history, including v1 architecture docs, dashboard guides, hardening docs, production-readiness docs, and runbooks.
- The new ticket path is `ttmp/2026/05/12/HOST-CONTRIB-GUIDE--contributor-guideline-backbone-for-go-go-host`.

### What was tricky to build
- The main trick was scoping: the request is not to write a narrow `CONTRIBUTING.md`, but to propose a whole contributor documentation system that can support many future humans and agents.

### What warrants a second pair of eyes
- Confirm that `HOST-CONTRIB-GUIDE` is the desired ticket ID and that the recommended document names match the repository’s long-term naming preferences.

### What should be done in the future
- Convert the proposed document set into committed repository-facing files after review.

### Code review instructions
- Start with the ticket index and this diary.
- Validate with `docmgr doctor --ticket HOST-CONTRIB-GUIDE --stale-after 30` after all docs are written.

### Technical details
- Commands run:
  - `docmgr status --summary-only`
  - `docmgr ticket create-ticket --ticket HOST-CONTRIB-GUIDE --title "Contributor guideline backbone for go-go-host" --topics go-go-host,contributor-guidelines,onboarding,documentation`
  - `docmgr doc add --ticket HOST-CONTRIB-GUIDE --doc-type design-doc --title "Contributor guideline backbone design"`
  - `docmgr doc add --ticket HOST-CONTRIB-GUIDE --doc-type reference --title "Investigation diary"`
  - `docmgr task add --ticket HOST-CONTRIB-GUIDE ...`

## Step 2: Inspect existing repository docs, dashboard guidelines, and runbooks

I mapped the local documentation inventory before drawing conclusions. The repository already has a promoted dashboard playbook under `docs/contributing`, and several ticket-local runbooks/design docs under `ttmp` that should shape any contributor-guideline package.

### Prompt Context

**User prompt (verbatim):** (same as Step 1)

**Assistant interpretation:** Use current repository docs and prior ticket docs as evidence for the proposed contributor guideline set.

**Inferred user intent:** Reuse existing good practices instead of inventing new documentation in isolation.

### What I did
- Listed repository docs with `find docs -maxdepth 5 -type f`.
- Listed ticket docs with `find ttmp -maxdepth 6 -type f -name '*.md'`.
- Read `AGENT.md`, `README.md`, `docs/contributing/playbooks/os1-admin-dashboard-ui-work-guidelines.md`, the HOST-003 local admin dashboard runbook, and the HOST-008 OS1 dashboard playbook.

### Why
- The user mentioned an existing dashboard guideline and asked to look in `ttmp` for runbooks.
- These documents provide the repository’s current voice: evidence-based, command-oriented, and intern-readable.

### What worked
- The dashboard guideline is already a strong model for a promoted contributor playbook.
- The local admin dashboard runbook has a concise command/verification format that should be copied for operational guides.

### What didn't work
- N/A.

### What I learned
- `AGENT.md` contains high-value contributor rules but is too agent-specific and generic to serve as the full human onboarding guide.
- `README.md` still describes “Phase 0 scaffold work,” while the source has advanced far beyond that into OIDC, agents, deployments, admin dashboards, and runtime operations.
- The repository already has a durable dashboard playbook at `docs/contributing/playbooks/os1-admin-dashboard-ui-work-guidelines.md`; the future documentation set should not duplicate it, but should index it from a larger contribution guide.

### What was tricky to build
- Some documents are repository-facing (`docs/contributing/...`) while others are ticket-local (`ttmp/...`). The proposed guideline backbone needs a rule for promotion: ticket docs are research and implementation memory; stable docs belong under `docs/contributing` or `docs/architecture`.

### What warrants a second pair of eyes
- Decide whether repository-facing docs should live under `docs/contributing/` only, or split into `docs/contributing`, `docs/architecture`, `docs/runbooks`, and `docs/developer`.

### What should be done in the future
- Promote stable guidance from `ttmp` into repository docs once reviewed.

### Code review instructions
- Compare the final recommendation against `docs/contributing/playbooks/os1-admin-dashboard-ui-work-guidelines.md` to ensure the dashboard-specific guidance remains authoritative.

### Technical details
- Key files inspected:
  - `AGENT.md`
  - `README.md`
  - `docs/contributing/playbooks/os1-admin-dashboard-ui-work-guidelines.md`
  - `ttmp/2026/05/11/HOST-003-ADMIN-DASHBOARD--go-go-host-platform-admin-dashboard/playbooks/01-local-admin-dashboard-runbook.md`
  - `ttmp/2026/05/12/HOST-008-OS1-DASHBOARD-VISUAL-POLISH--os1-dashboard-visual-polish-and-widget-reuse/playbook/01-os1-admin-dashboard-ui-work-guidelines.md`

## Step 3: Map architecture and comparable project conventions

I then inspected the source tree and neighboring projects to understand the actual system contributors need to work in. The design recommendation is anchored in backend, runtime, deployment, frontend, CLI, web-embed, and development-workflow files rather than only in prose docs.

### Prompt Context

**User prompt (verbatim):** (same as Step 1)

**Assistant interpretation:** Produce a file-backed architecture map that a new intern can use as the contributor onboarding spine.

**Inferred user intent:** Prevent future contributors from making changes in the wrong layer or bypassing existing conventions.

### What I did
- Listed source directories under `cmd`, `internal`, and `web/admin/src`.
- Counted package/file distribution to identify core subsystems.
- Inspected line-number evidence for `internal/httpapi/handler.go`, `internal/control/core.go`, `internal/deploy/bundle.go`, `internal/runtime/runtime.go`, `internal/runtime/supervisor.go`, `internal/sitejs/web/host.go`, `internal/store/store.go`, `web/admin/src/app/routes.tsx`, `web/admin/src/services/goGoHostApi.ts`, `Makefile`, `.devctl.yaml`, and `web/admin/package.json`.
- Read comparable `AGENT.md` files from `go-go-goja` and `glazed`.

### Why
- Contribution guidelines must be grounded in the actual layering and tools.
- Neighboring projects show the shared go-go-golems style: Go/Cobra/Glazed, context-aware APIs, gofmt/tests/lint, no stray modules, and explicit web build/embed patterns.

### What worked
- `rg -n` produced useful line references for the final design doc.
- The source tree has clear boundaries: `cmd`, `internal/config`, `internal/httpapi`, `internal/control`, `internal/store`, `internal/deploy`, `internal/runtime`, `internal/sitejs`, `internal/webadmin`, and `web/admin`.

### What didn't work
- A broad `nl -ba` command produced too much output and was truncated. I switched to targeted reads and `rg -n` evidence.

### What I learned
- The current `Core` object is the key backend boundary: transports should go through control services, not directly into store/runtime.
- HTTP route registration in `internal/httpapi/handler.go` shows which APIs are authenticated user APIs and which are signed agent APIs.
- Deployment validation performs path/capability checks, stores immutable artifacts, then dry-runs a runtime before marking a deployment validated.
- Runtime activation builds and health-checks a new runtime before swapping traffic.
- The dashboard uses React 19, Vite 6, RTK Query, Storybook, MSW, and OS-core.

### What was tricky to build
- The repository’s `README.md` and current source are out of sync. The contributor guide should explicitly warn new contributors to treat source and recent `ttmp` docs as fresher until README is updated.

### What warrants a second pair of eyes
- Verify if `README.md` should be refreshed as part of implementing the contributor documentation set.

### What should be done in the future
- Add a generated or maintained architecture map so this evidence does not have to be rediscovered each time.

### Code review instructions
- Start code review of the final design doc at the architecture map and evidence table.
- Spot-check cited source lines before promotion to repository-facing docs.

### Technical details
- Representative commands:
  - `find cmd internal web/admin/src -maxdepth 3 -type f | sort`
  - `find internal cmd -type f -name '*.go' | awk -F/ '{print $1"/"$2}' | sort | uniq -c | sort -nr`
  - `rg -n "func NewHandler|type Core|SafeCapabilities|createApi|createBrowserRouter" cmd internal web/admin/src README.md AGENT.md docs/contributing -S`

## Step 4: Write deliverables and prepare validation/upload

I wrote the primary design document and a companion runbook. The design document explains the system and proposes a concrete documentation set; the runbook is the copy/paste checklist for implementing that documentation set in the repository.

### Prompt Context

**User prompt (verbatim):** (same as Step 1)

**Assistant interpretation:** Produce final ticket documents, relate important files, run docmgr validation, and upload the bundle to reMarkable.

**Inferred user intent:** Make the output durable and immediately useful for future contributors.

### What I did
- Added a playbook document named `Contributor documentation set runbook`.
- Wrote an intern-ready design guide covering system purpose, architecture, contribution lanes, documentation set proposal, API/file references, pseudocode, diagrams, implementation plan, testing, and review checklists.
- Wrote this diary entry before validation and upload.

### Why
- A design doc is good for reasoning and onboarding; a playbook is better for execution.
- Separating the two helps future agents implement the docs without re-reading the entire analysis every time.

### What worked
- The repository has enough existing docs and code evidence to propose a strong documentation backbone without guessing.

### What didn't work
- N/A at this step.

### What I learned
- The best documentation backbone should have both stable repository docs and ticket-local research docs. Stable docs should be concise, reviewed, and linked from README; ticket docs should preserve implementation history and investigation details.

### What was tricky to build
- The guide needed to be both a system explainer and a meta-guide for future documentation. I handled this by making the first half explain the system and the second half define the proposed documentation set and implementation plan.

### What warrants a second pair of eyes
- The proposed document taxonomy and promotion criteria should be reviewed by maintainers before bulk-creating repository-facing docs.

### What should be done in the future
- Implement Phase 1 of the documentation set: contributor overview, architecture map, backend/service guide, frontend guide index, runtime/deployment guide, local development runbook, and testing matrix.

### Code review instructions
- Review `design-doc/01-contributor-guideline-backbone-design.md` first.
- Review `playbook/01-contributor-documentation-set-runbook.md` as the implementation checklist.
- Validate with `docmgr doctor --ticket HOST-CONTRIB-GUIDE --stale-after 30`.

### Technical details
- Files created/updated in this step:
  - `ttmp/2026/05/12/HOST-CONTRIB-GUIDE--contributor-guideline-backbone-for-go-go-host/design-doc/01-contributor-guideline-backbone-design.md`
  - `ttmp/2026/05/12/HOST-CONTRIB-GUIDE--contributor-guideline-backbone-for-go-go-host/playbook/01-contributor-documentation-set-runbook.md`
  - `ttmp/2026/05/12/HOST-CONTRIB-GUIDE--contributor-guideline-backbone-for-go-go-host/reference/01-investigation-diary.md`
