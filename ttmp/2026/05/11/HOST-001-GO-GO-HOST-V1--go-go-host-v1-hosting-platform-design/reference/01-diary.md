---
Title: Diary
Ticket: HOST-001-GO-GO-HOST-V1
Status: active
Topics:
    - goja
    - hosting
    - go-go-host
    - rtk-query
    - storybook
    - vm-runtime
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: ../../../../../../../../../../code/wesen/go-go-golems/go-go-parc/Research/Institute/Proposals/2026/05/11/PROPOSAL - Goja Sites Hosting Service.md
    - Path: ../../../../../../../2026-05-01--wish-git/internal/policy/authorize.go
      Note: New local Wish Git evidence recorded in Step 4
    - Path: Dockerfile
      Note: Phase 12 production image recipe
    - Path: cmd/go-go-host-agent/cmds/support.go
      Note: agent CLI signing
    - Path: cmd/go-go-host/cmds/maintenance.go
      Note: Phase 12 operator CLI
    - Path: internal/control/agent_runs.go
      Note: Phase 9 service logic for enrollment
    - Path: internal/control/deployments.go
      Note: Activation now includes verified custom domains
    - Path: internal/control/maintenance.go
      Note: Phase 12 export/prune/retention service
    - Path: internal/control/services.go
      Note: Phase 11 site config
    - Path: internal/httpapi/agent_signed_integration_test.go
      Note: Phase 9 security regression coverage
    - Path: internal/httpapi/agents_audit.go
      Note: agent enrollment and signed deploy-run HTTP endpoints
    - Path: internal/httpapi/site_settings.go
      Note: Phase 11 site settings API handlers and DTOs
    - Path: internal/store/queries/sites.sql
      Note: Phase 11 sqlc queries for config and domains
    - Path: scripts/final-e2e-playwright.mjs
      Note: Phase 12 gated final E2E script
    - Path: ttmp/2026/05/11/HOST-001-GO-GO-HOST-V1--go-go-host-v1-hosting-platform-design/design-doc/01-go-go-host-v1-hosting-platform-intern-design-and-implementation-guide.md
      Note: Primary design deliverable recorded in diary
    - Path: web/admin/src/pages/SiteSettingsPage/SiteSettingsPage.tsx
      Note: Phase 11 operator settings page
ExternalSources: []
Summary: Chronological diary for creating the go-go-host v1 design ticket and intern implementation guide.
LastUpdated: 2026-05-11T17:05:00-04:00
WhatFor: Use this to continue the investigation or review how the design was produced.
WhenToUse: When extending the go-go-host v1 design, implementing phases, or auditing source evidence.
---






# Diary

## Goal

Capture the investigation and documentation work for the first v1 design of `go-go-host`, including source evidence, design choices, dashboard direction, VM-runtime analysis, and delivery steps.

## Step 1: Ticket setup and proposal intake

I created a docmgr ticket for the first `go-go-host` v1 hosting-platform design and read the upstream PARC proposal. The proposal established the product goal: turn the existing trusted `goja-site` prototype into a multi-tenant platform with users, organizations, deployments, agents, audit, runtime state, and an admin console.

This step created the workspace that now contains the design guide and this diary, so future implementation work can be related to the same ticket rather than living as ad-hoc Markdown.

### Prompt Context

**User prompt (verbatim):** "Create a new docmgr to create the first v1 of go-go-host the hosting platform. Read /home/manuel/code/wesen/go-go-golems/go-go-parc/Research/Institute/Proposals/2026/05/11/PROPOSAL - Goja Sites Hosting Service.md

And then Create  a detailed analysis / design / implementation guide that is for a new intern, explaining all the parts of the system needed to understand what it is, with prose paragraphs and bullet point sand pseudocode and diagrams and api references and file references. It should be very clear and technical. Store in the ticket and the nupload to remarkable.


Read the relevant source in this workspace and all.

We want to copy/refactor the runtime into the new service.

We also want a phase for a proper user dashboaord.

We also want to use a react/rtk-query/storybook solution using https://www.npmjs.com/package/@go-go-golems/os-core for building the site (see /home/manuel/workspaces/2026-05-11/npm-packages-go-go-os/2026-05-11--npm-go-go-os-test for an example).

Furthermore, also analyze ./vm-system that might have useful pieces for vm runtime management as well.

Keep a diary as you work."

**Assistant interpretation:** Create a docmgr ticket and write an evidence-backed design/implementation guide for an intern to build go-go-host v1, including runtime refactor guidance, dashboard architecture, VM-system reuse analysis, and reMarkable delivery.

**Inferred user intent:** Establish the initial technical blueprint for building a real go-go-host hosting platform from existing prototypes while preserving a reviewable investigation trail.

**Commit (code):** N/A — documentation-only work in the ticket workspace.

### What I did

- Read the proposal at `/home/manuel/code/wesen/go-go-golems/go-go-parc/Research/Institute/Proposals/2026/05/11/PROPOSAL - Goja Sites Hosting Service.md`.
- Ran `docmgr status --summary-only` to confirm the docmgr root.
- Created ticket `HOST-001-GO-GO-HOST-V1` with topics `goja,hosting,go-go-host,rtk-query,storybook,vm-runtime`.
- Created the primary design document and diary document.

### Why

- The user explicitly requested a new docmgr ticket and a diary.
- Starting with the proposal prevents the design from drifting away from the product goals.

### What worked

- `docmgr ticket create-ticket` created the expected ticket workspace under `go-go-host/ttmp/2026/05/11/HOST-001-GO-GO-HOST-V1--go-go-host-v1-hosting-platform-design`.
- `docmgr doc add` created the design doc and diary stubs.

### What didn't work

- N/A for this step.

### What I learned

- The current workspace contains copies of the important source projects: `2026-05-03--goja-hosting-site`, `2026-05-03--agent-enroll`, `go-go-goja`, `go-go-host`, and `vm-system`.
- The static multi-site prototype is close to the requested product runtime but lacks the control-plane model.

### What was tricky to build

- The ticket had to be created from the workspace's docmgr config, whose root is `go-go-host/ttmp`, not the workspace root.

### What warrants a second pair of eyes

- The chosen ticket ID and title are intentionally broad. If a smaller implementation ticket is desired, split this later into runtime, deployment, dashboard, and agent sub-tickets.

### What should be done in the future

- Add implementation tickets for each phase once coding begins.

### Code review instructions

- Start with `go-go-host/ttmp/2026/05/11/HOST-001-GO-GO-HOST-V1--go-go-host-v1-hosting-platform-design/index.md`.
- Confirm the design and diary documents exist and have accurate frontmatter.

### Technical details

Commands run:

```bash
pwd && rg --files -g '!*node_modules*' -g '!*.png' -g '!*.jpg' | sed -n '1,200p' && echo '---docmgr---' && docmgr status --summary-only

docmgr ticket create-ticket --ticket HOST-001-GO-GO-HOST-V1 --title "go-go-host v1 hosting platform design" --topics goja,hosting,go-go-host,rtk-query,storybook,vm-runtime

docmgr doc add --ticket HOST-001-GO-GO-HOST-V1 --doc-type design-doc --title "go-go-host v1 hosting platform intern design and implementation guide"

docmgr doc add --ticket HOST-001-GO-GO-HOST-V1 --doc-type reference --title "Diary"
```

## Step 2: Evidence gathering across runtime, agent, dashboard, and vm-system sources

I mapped the source files that should shape v1. The important runtime evidence came from `goja-site`: the current `Server` owns the database, Goja runtime, route host, and HTTP server; `MultiServer` already dispatches by Host header; and `web.Host` already converts HTTP requests into Goja handler calls through a runtime owner.

I also inspected Agent Enroll for signed machine requests and run tokens, `go-go-goja` modules for capability boundaries, the npm OS-core example for the React/RTK Query/Storybook dashboard stack, and `vm-system` for runtime-management patterns such as core/ports separation, runtime registry, persistent runtime state, and startup reconciliation.

### Prompt Context

**User prompt (verbatim):** (same as Step 1)

**Assistant interpretation:** Gather concrete line-backed evidence from all relevant code before writing the design.

**Inferred user intent:** Ensure recommendations are grounded in the actual workspace instead of generic platform architecture advice.

**Commit (code):** N/A — investigation only.

### What I did

- Listed workspace files with `rg --files`.
- Counted and inspected key runtime files under `2026-05-03--goja-hosting-site`.
- Inspected Agent Enroll signed request and run-token files.
- Inspected `go-go-goja` `express`, `database`, `fs`, and `exec` modules.
- Inspected `vm-system` core, ports, daemon, store migrations, HTTP transport, and session manager.
- Inspected `/home/manuel/workspaces/2026-05-11/npm-packages-go-go-os/2026-05-11--npm-go-go-os-test` package scripts, theme imports, RTK Query store/API slice, and control-panel UI.

### Why

- The final guide needed file references, API sketches, diagrams, pseudocode, and actionable implementation phases.
- The user specifically asked to analyze `./vm-system` and the npm OS-core example.

### What worked

- The relevant source files were present in the workspace.
- `vm-system` provided several reusable architectural ideas even though its runtime execution model should not be copied wholesale.
- The npm example provided a clear dashboard stack: Vite, React, RTK Query, Storybook, `@go-go-golems/os-core`, and OS widget packages.

### What didn't work

- The workspace did not appear to include `2026-05-01--wish-git`, even though the proposal references it. I therefore did not use local Wish Git line references in the final design.
- One attempted read of `examples/03-rtk-query-control-panel/src/ControlPanelApp.tsx` failed because the real path is `examples/03-rtk-query-control-panel/src/features/ControlPanelApp/ControlPanelApp.tsx`.

### What I learned

- `goja-site` currently enables `fs`, `path`, `time`, and `timer` middleware for its prototype runtime. v1 should treat `fs` as too powerful for default hosted sites.
- The database module already supports the exact capability boundary needed for hosting: preconfigured DB with `configure()` disabled.
- Agent Enroll's signed request verifier is directly reusable for agent deploys.
- `vm-system`'s `Core`/ports split is a good template for keeping HTTP handlers thin and services testable.

### What was tricky to build

- The source tree contains multiple related projects, and some proposal references point outside the current workspace. I avoided claiming local evidence for absent files and focused on sources actually available here.

### What warrants a second pair of eyes

- Confirm whether Wish Git should be copied into this workspace or referenced only through the proposal in future design revisions.
- Confirm whether the first implementation should use Postgres immediately or SQLite-compatible migrations during prototyping.

### What should be done in the future

- When implementation begins, add exact tests for runtime activation, agent replay denial, and forbidden capability rejection.

### Code review instructions

- Review the evidence claims in the design doc against the referenced files.
- Pay special attention to the capability-security conclusions around `fs` and `exec`.

### Technical details

Representative commands run:

```bash
for f in \
2026-05-03--goja-hosting-site/pkg/app/server.go \
2026-05-03--goja-hosting-site/pkg/app/multi_server.go \
2026-05-03--goja-hosting-site/pkg/app/multi_config.go \
2026-05-03--goja-hosting-site/pkg/web/host.go \
2026-05-03--goja-hosting-site/pkg/web/route_registry.go \
2026-05-03--goja-hosting-site/pkg/dbguard/guard.go \
2026-05-03--agent-enroll/internal/agent/signature.go \
2026-05-03--agent-enroll/internal/runs/runs.go \
go-go-goja/modules/express/express.go \
go-go-goja/modules/database/database.go \
vm-system/pkg/vmcontrol/core.go \
vm-system/pkg/vmcontrol/runtime_registry.go \
vm-system/pkg/vmdaemon/app.go \
vm-system/pkg/vmtransport/http/server.go \
; do [ -f "$f" ] && echo "$f $(wc -l < "$f")"; done

find /home/manuel/workspaces/2026-05-11/npm-packages-go-go-os/2026-05-11--npm-go-go-os-test/examples -maxdepth 5 -type f | sort | sed -n '1,160p'
```

## Step 3: Design guide authoring

I wrote the primary design document as an intern-oriented implementation guide. It includes an executive summary, problem statement, mental model, current-source analysis, proposed backend architecture, runtime supervisor design, deployment model, API sketches, dashboard plan, implementation phases, testing strategy, design decisions, alternatives, risks, and a file reference map.

The document emphasizes copying/refactoring the existing runtime first, then adding product control-plane concepts around it. It also calls out security boundaries: preconfigured DB is acceptable, unrestricted `fs` should not be default, and `exec` must not be exposed in hosted v1.

### Prompt Context

**User prompt (verbatim):** (same as Step 1)

**Assistant interpretation:** Produce the actual detailed design deliverable in the docmgr ticket.

**Inferred user intent:** Give a new intern enough context to start implementing without repeatedly rediscovering the runtime, dashboard, deployment, and agent-auth architecture.

**Commit (code):** N/A — documentation-only work.

### What I did

- Wrote `design-doc/01-go-go-host-v1-hosting-platform-intern-design-and-implementation-guide.md`.
- Included Mermaid diagrams for the mental model, architecture, and deployment state machine.
- Added pseudocode for `Core`, `RuntimeSupervisor`, `SiteRuntime`, deployment validation, integration tests, and RTK Query API slices.
- Added API reference sketches for human, site, deployment, and agent endpoints.
- Added phased implementation plan with completion criteria.

### Why

- The user asked for prose paragraphs, bullet points, pseudocode, diagrams, API references, and file references.
- Intern-facing docs need both conceptual explanation and concrete starting points.

### What worked

- The existing source structure mapped cleanly into proposed packages: `internal/runtime`, `internal/control`, `internal/store`, `internal/httpapi`, `internal/deploy`, `internal/agentauth`, and `web/admin`.
- The React dashboard phase could be tied directly to the OS-core example's real store/API/story structure.

### What didn't work

- The final design is necessarily a design document, not an implementation. It does not create code scaffolding yet.

### What I learned

- The safest v1 path is to keep public site runtime serving and control-plane API serving in the same daemon initially, while defining interfaces that can later support process-backed runtimes.
- The dashboard should be part of the platform architecture, because it is the main way operators will understand deployments, agent grants, runtime health, and audit events.

### What was tricky to build

- Balancing detail against navigability required stable sections and explicit file maps. I used tables and pseudocode to keep the long design scannable.

### What warrants a second pair of eyes

- Review the data model before migrations are written.
- Review whether `site_capabilities` should be per site, per deployment, or both.
- Review the deployment upload API shape: multipart upload is simpler, signed upload URL is closer to production.

### What should be done in the future

- Convert phases into docmgr tasks or implementation tickets.
- Add a smaller architecture decision record for control DB choice.

### Code review instructions

- Start with the `Executive summary`, then read `Current-state source analysis`, then `Implementation phases`.
- Validate file references against the workspace before coding.

### Technical details

Primary doc path:

```text
go-go-host/ttmp/2026/05/11/HOST-001-GO-GO-HOST-V1--go-go-host-v1-hosting-platform-design/design-doc/01-go-go-host-v1-hosting-platform-intern-design-and-implementation-guide.md
```

## Step 4: Added local Wish Git evidence

After the Wish Git repository was added locally under `./2026-05-01--wish-git/`, I revisited the design guide and replaced the earlier proposal-level Wish Git treatment with concrete local source references. The update strengthens the deployment-policy section because Wish Git contains exactly the boundary-enforcement pattern go-go-host should emulate: create a scoped delegated run first, then check scope again at the final write boundary.

This did not change the v1 recommendation to start with bundle uploads. It made the future Git-deploy path more concrete and added a reusable policy model for deploy-run actions, channels, and bundle paths.

### Prompt Context

**User prompt (verbatim):** "I added the wish-git locally ( ./2026-05-01--wish-git/)

anything to update?"

**Assistant interpretation:** Inspect the newly available Wish Git checkout and update the existing ticket documentation if it changes or strengthens the design.

**Inferred user intent:** Keep the design evidence-backed now that a previously missing referenced project is available locally.

**Commit (code):** N/A — documentation-only update.

### What I did

- Listed files under `2026-05-01--wish-git`.
- Read `internal/policy/authorize.go`, `internal/policy/patterns.go`, `internal/githook/pre_receive.go`, `internal/api/agent_runs.go`, `internal/api/ssh_certs.go`, and `internal/store/migrations/001_initial_schema.sql`.
- Added a new design-guide section: `Existing Wish Git scoped-delegation and boundary-enforcement pattern`.
- Updated the `Git-only deployment` alternative with concrete local Wish Git references.
- Added Wish Git files to the design-guide file reference map.

### Why

- The first pass could not line-reference Wish Git because it was not present in the workspace.
- Wish Git provides concrete source evidence for scoped delegated runs, allowed action/ref/path policy, final-boundary validation, TTL caps, SSH certificate issuance, local OIDC user mapping, and audit schema.

### What worked

- The newly added repository contained the expected files.
- The policy code was compact enough to translate directly into go-go-host deploy-run pseudocode.
- The pre-receive hook provided a clear analogy for bundle validation.

### What didn't work

- N/A for this step.

### What I learned

- `AllowsPath` intentionally denies when the allowed path list is empty, which is a good safe default for deploy-run bundle scopes.
- Wish Git caps/defaults run TTL to at most 3600 seconds in the run creation handler; go-go-host should similarly cap deploy-run/upload-token lifetimes.
- The SSH certificate path is more relevant to a future Git remote deployment source than to bundle-upload v1.

### What was tricky to build

- The design had to preserve the earlier decision that Git is a later deployment source, while still extracting Wish Git's useful policy and boundary-validation patterns for bundle v1.

### What warrants a second pair of eyes

- Review whether go-go-host should reuse Wish Git's `/**` glob convention exactly, or define a stricter bundle-path matcher.
- Review whether deploy-run `AllowedChannels` should use the same glob helper as `AllowedPaths`.

### What should be done in the future

- Add implementation tests equivalent to Wish Git's policy/pre-receive tests for deploy-run path and channel scopes.
- Consider a later Git deployment ticket that starts from Wish Git's SSH certificate flow.

### Code review instructions

- Review `2026-05-01--wish-git/internal/policy/authorize.go` first.
- Then review `2026-05-01--wish-git/internal/githook/pre_receive.go` and compare it to the proposed bundle validator.
- Finally review the design guide's `Existing Wish Git scoped-delegation and boundary-enforcement pattern` section.

### Technical details

Commands run:

```bash
find 2026-05-01--wish-git -maxdepth 3 -type f | sort | sed -n '1,160p'

for f in \
  2026-05-01--wish-git/internal/policy/authorize.go \
  2026-05-01--wish-git/internal/policy/patterns.go \
  2026-05-01--wish-git/internal/githook/pre_receive.go \
  2026-05-01--wish-git/internal/api/agent_runs.go \
  2026-05-01--wish-git/internal/api/ssh_certs.go \
  2026-05-01--wish-git/internal/store/migrations.go; do
  echo "===== $f"
  nl -ba "$f" | sed -n '1,220p'
done

nl -ba 2026-05-01--wish-git/internal/store/migrations/001_initial_schema.sql | sed -n '1,240p'
```

## Step 5: Updated research for two dashboards and implementation task plan

The original PARC proposal was updated to distinguish a normal user dashboard from a platform admin console. I re-read the proposal, inspected the local Agent Enroll dashboard sources that the proposal now identifies as a useful user-dashboard reference, and updated the design guide so it no longer treats the dashboard as one generic admin surface.

I also replaced the short phase list in `tasks.md` with a detailed implementation checklist split into twelve phases. The checklist now includes backend scaffold, database/auth, runtime refactor, deployment pipeline, human CLI, user dashboard, platform admin console, agent deploys, capability hardening, domains/configuration, and production hardening.

### Prompt Context

**User prompt (verbatim):** "we updated the original PROPOSAL to have 2 dashboards: one user, one admin.

Read again, then update your research.

Then, create a detailed list of tasks for implementing the system, split into phases."

**Assistant interpretation:** Re-read the updated proposal, update the research/design to reflect separate user and platform-admin dashboards, and turn the implementation plan into a detailed phased task list.

**Inferred user intent:** Make the design more product-accurate before implementation starts, and provide an actionable intern-ready backlog.

**Commit (code):** N/A — documentation-only update.

### What I did

- Re-read `/home/manuel/code/wesen/go-go-golems/go-go-parc/Research/Institute/Proposals/2026/05/11/PROPOSAL - Goja Sites Hosting Service.md`.
- Inspected Agent Enroll dashboard files under `2026-05-03--agent-enroll/web/dashboard/src`.
- Updated the design guide's dashboard section to split `/app/*` user dashboard routes from `/admin/*` platform admin routes.
- Added user dashboard workflows: first login, site creation, deployment, rollback, bot token creation, agent grants, runtime inspection, usage, and audit.
- Added platform admin workflows: global overview, users/org support, runtime operations, quota policy, deployment moderation, domains, agent abuse response, global audit, and platform settings.
- Rewrote `tasks.md` as a detailed phased implementation plan with concrete checkboxes and exit criteria.

### Why

- A normal user dashboard and a platform admin console have different audiences, authorization scopes, API endpoints, and navigation models.
- Implementation tasks need to track that split explicitly so the team does not accidentally build only an operator tool and miss the product's user workflow.

### What worked

- Agent Enroll's dashboard provided concrete local patterns for session guards, org-scoped routing, RTK Query endpoints, one-time token reveal, copyable commands, agent lists, polling runs, and revoke confirmation.
- The implementation plan now has enough granularity to be used as a first project backlog.

### What didn't work

- N/A for this step.

### What I learned

- The user dashboard should likely ship before or alongside the platform admin console because it is how normal users create sites, deploy code, and manage agents.
- The admin console should use explicit `/api/v1/admin/*` endpoints and server-side `platform_admin` checks, not client-only hiding.

### What was tricky to build

- The design needed to keep one embedded SPA for implementation simplicity while making the product distinction between user and admin surfaces clear. I resolved this by proposing separate route groups, feature folders, role guards, and endpoint naming conventions.

### What warrants a second pair of eyes

- Review whether Phase 7 user dashboard should come before the human CLI, or whether the current ordering should remain CLI first for faster deployment testing.
- Review whether members/invitations belong in v1 user-dashboard scope or should be deferred.

### What should be done in the future

- Turn each phase in `tasks.md` into separate implementation tickets once coding begins.
- Add UI wireframes or Storybook-first component tasks for the user dashboard and admin console.

### Code review instructions

- Review the design guide section `Dashboard design`.
- Review `tasks.md` from `Phase 0` through `Phase 12` and confirm the ordering matches implementation priorities.
- Check the Agent Enroll dashboard references if adapting UI code or patterns.

### Technical details

Commands run:

```bash
find 2026-05-03--agent-enroll/web/dashboard/src -maxdepth 3 -type f | sort | sed -n '1,160p'

for f in \
  2026-05-03--agent-enroll/web/dashboard/src/App.tsx \
  2026-05-03--agent-enroll/web/dashboard/src/pages/AgentsPage.tsx \
  2026-05-03--agent-enroll/web/dashboard/src/components/organisms/EnrollmentTokenPanel.tsx \
  2026-05-03--agent-enroll/web/dashboard/src/components/molecules/SecretRevealBox.tsx \
  2026-05-03--agent-enroll/web/dashboard/src/components/molecules/CommandCopyBox.tsx \
  2026-05-03--agent-enroll/web/dashboard/src/api/kanbanApi.ts; do
  echo "===== $f"
  nl -ba "$f" | sed -n '1,180p'
done
```

## Step 6: Added Glazed command-structure requirement

The implementation plan now explicitly requires both command-line binaries to use the Glazed command structure. This affects `go-go-host` for human workflows and `go-go-host-agent` for headless agent workflows.

This update keeps the CLI implementation aligned with the go-go-golems ecosystem: commands should expose typed settings, consistent output formats, embedded help, logging, and structured rows instead of hand-written Cobra-only handlers.

### Prompt Context

**User prompt (verbatim):** "Alright, and use the glazed command structure for implementation."

**Assistant interpretation:** Update the design and task plan so CLI implementation uses Glazed commands and standard Glazed/Cobra root wiring.

**Inferred user intent:** Ensure the future implementation follows established go-go-golems CLI conventions rather than inventing a bespoke command architecture.

**Commit (code):** N/A — documentation-only update.

### What I did

- Loaded the `glazed-command-authoring` skill instructions.
- Added a design-guide section: `CLI command architecture: use Glazed commands`.
- Added root command requirements for Glazed logging and embedded help setup.
- Added command folder layout for `cmd/go-go-host` and `cmd/go-go-host-agent`.
- Added a Glazed command skeleton for `sites list`.
- Updated Phase 0 tasks to add Glazed dependencies, root wiring, doc packages, and CLI folder conventions.
- Reworked Phase 6 into `Human CLI using Glazed commands` with concrete command-authoring requirements.
- Updated Phase 9 so agent CLI commands are Glazed commands too.

### Why

- go-go-host commands produce structured resources: sites, deployments, agents, runs, audit events, runtime status, and validation reports.
- Glazed provides consistent `--output`, `--fields`, table/JSON/YAML output, command settings, logging, and help.

### What worked

- The existing phased task list had clear CLI phases, so the Glazed requirement fit cleanly into Phase 0, Phase 6, and Phase 9.

### What didn't work

- N/A for this step.

### What I learned

- Login may still need a small custom Cobra wrapper if browser OAuth has special control flow, but parsed settings and root wiring should remain consistent with Glazed conventions.

### What was tricky to build

- Some commands are mutations or uploads rather than pure tabular queries. The design resolves this by requiring mutation commands to emit stable result rows and validation commands to default to YAML/JSON when reports are nested.

### What warrants a second pair of eyes

- Review whether `go-go-host deploy ./site` should be a top-level Glazed command or live under a `deployments upload` group.
- Review output defaults for nested deployment validation reports.

### What should be done in the future

- Add embedded Glazed help pages for login, deploy, rollback, agent enrollment, and troubleshooting signed requests.

### Code review instructions

- Review the design section `CLI command architecture: use Glazed commands`.
- Review `tasks.md` Phase 0, Phase 6, and Phase 9.

### Technical details

Key Glazed requirements captured in the design:

```text
- command struct embeds *cmds.CommandDescription
- settings struct uses glazed tags
- constructor uses cmds.NewCommandDescription
- flags/arguments use fields.New
- command includes settings.NewGlazedSchema and cli.NewCommandSettingsSection
- RunIntoGlazeProcessor decodes vals.DecodeSectionInto(schema.DefaultSlug, settings)
- rows are emitted with types.NewRow and gp.AddRow
- root wires logging.AddLoggingSectionToRootCommand and help_cmd.SetupCobraRootCommand
```

## Step 7: Implemented Phase 0 scaffold with Glazed CLI roots

I implemented the Phase 0 scaffold in `go-go-host`. The repository now has a runnable daemon, two Glazed-wired CLI roots, placeholder control/store/webadmin packages, health/version APIs, local dev config, Makefile targets, and a smoke test for the HTTP handler.

A small workspace dependency issue surfaced while validating. The old module-local `replace` in `2026-05-03--goja-hosting-site/go.mod` pointed at a missing `../corporate-headquarters/go-go-goja` checkout. I removed that nested replace and centralized the local mapping in `go.work` with a version-specific replacement for `github.com/go-go-golems/go-go-goja v0.0.0 => ./go-go-goja`, because the workspace still has a `v0.0.0` requirement from `goja-hosting-site`.

### Prompt Context

**User prompt (verbatim):** "go ahead"

**Assistant interpretation:** Start implementation with Phase 0, following the just-added Glazed command structure requirement.

**Inferred user intent:** Move from planning into code, beginning with the scaffold and CLI architecture.

**Commit (code):** N/A — changes are uncommitted.

### What I did

- Replaced the placeholder `cmd/XXX/main.go` with `cmd/go-go-hostd/main.go`.
- Added `cmd/go-go-host` root and `status` Glazed command.
- Added `cmd/go-go-host-agent` root and `status` Glazed command.
- Added embedded help doc packages for both CLIs.
- Added `internal/config`, `internal/control`, `internal/store`, `internal/httpapi`, and `internal/webadmin`.
- Added `configs/dev.yaml`.
- Rewrote `README.md` for go-go-host v1 scope.
- Rewrote `Makefile` with build/test/run targets.
- Updated module path to `github.com/go-go-golems/go-go-host` and added Glazed/Cobra/YAML/logging dependencies.
- Removed stale nested replace from `2026-05-03--goja-hosting-site/go.mod`.
- Added version-specific workspace replace in `go.work` for `github.com/go-go-golems/go-go-goja v0.0.0`.

### Why

- Phase 0 needs a runnable skeleton before control-plane or runtime work starts.
- Glazed root wiring now exists from the beginning, avoiding a later CLI refactor.

### What worked

- `go test ./...` passes inside `go-go-host`.
- A local daemon starts from `configs/dev.yaml` and serves `/healthz`.
- `go-go-host status --output json` works against the daemon.
- `go-go-host-agent status --output json` works against the daemon.

### What didn't work

- First `go mod tidy` failed because I initially used `github.com/go-go-golems/glazed v0.0.0`, which is not a real revision.
- After switching to `v1.2.7`, tests still failed because a stale nested replace in `2026-05-03--goja-hosting-site/go.mod` pointed to a missing path.
- Removing that nested replace was not sufficient by itself because the workspace still had a `github.com/go-go-golems/go-go-goja v0.0.0` requirement. A version-specific `go.work` replace fixed it.

### What I learned

- In this workspace, `go.work use ./go-go-goja` alone did not stop the Go command from trying to resolve the `v0.0.0` requirement from another workspace module. The version-specific workspace replace is the cleanest central fix while avoiding stale module-local paths.

### What was tricky to build

- The CLI skeleton had to satisfy Glazed conventions while keeping Phase 0 small. I implemented only `status` for each CLI, but with the right command struct, settings tags, sections, decoding, and row output pattern.

### What warrants a second pair of eyes

- Review whether the `go.work` version-specific replace should stay or whether `2026-05-03--goja-hosting-site` should instead require a real tagged `go-go-goja` version.
- Review whether `go-go-hostd` should also use Glazed root wiring or remain plain Cobra because it is a daemon-only command.

### What should be done in the future

- Start Phase 1 migrations and control-plane services.
- Add real `go-go-host me`, `orgs`, and `sites` Glazed command groups after auth/store exist.

### Code review instructions

- Start with `go-go-host/cmd/go-go-host/root.go` and `go-go-host/cmd/go-go-host/cmds/status.go` for the Glazed pattern.
- Review `go-go-host/cmd/go-go-host-agent/root.go` and `cmds/status.go` for the agent CLI mirror.
- Review `go-go-host/cmd/go-go-hostd/main.go` and `internal/httpapi/handler.go` for the daemon skeleton.
- Validate with `cd go-go-host && go test ./...`.

### Technical details

Commands run:

```bash
cd go-go-host && go mod tidy
cd go-go-host && go fmt ./... && go test ./...
cd go-go-host && go test ./...

cd go-go-host
(go run ./cmd/go-go-hostd --config configs/dev.yaml > /tmp/go-go-hostd.log 2>&1 & echo $! > /tmp/go-go-hostd.pid)
curl -fsS http://127.0.0.1:8080/healthz
go run ./cmd/go-go-host status --api-url http://127.0.0.1:8080 --output json
go run ./cmd/go-go-host-agent status --api-url http://127.0.0.1:8080 --output json
kill $(cat /tmp/go-go-hostd.pid)
```

## Step 8: Implemented Phase 1 control-plane schema and service skeleton

I continued with Phase 1 after the workspace `go-go-goja` version issue was fixed. The implementation now has an embedded control-plane schema, migration runner, initial store methods, org/site services with authorization checks, and tests for migrations, membership authorization, cross-org denial, and user upsert behavior.

I also added a local Docker Compose stack for development infrastructure. Postgres belongs naturally in Phase 1 because this phase introduces the control-plane database. Keycloak is primarily Phase 2 authentication work, but it is included in the compose file now so the local Postgres + Keycloak stack is ready before OIDC validation is wired.

### Prompt Context

**User prompt (verbatim):** "ok i fixed the 0.0.0 package reference, we should be good now i think.

Continue with phase 1. Also, in which phase should we introduce a keycloak docker compose setup (and postgresql in there too, i guess)"

**Assistant interpretation:** Proceed into Phase 1 implementation and clarify where Postgres/Keycloak Docker Compose belongs in the phased plan.

**Inferred user intent:** Move the project from skeleton to durable product state while preparing the dev infrastructure for database and auth work.

**Commit (code):** N/A — changes are uncommitted.

### What I did

- Added embedded SQL migrations under `go-go-host/internal/store/migrations`.
- Added schema tables for users, orgs, memberships, platform admins, sites, domains, quotas, capabilities, deployments, deploy runs, agents, keys, grants, nonces, and audit log.
- Added migration runner and store constructor in `internal/store`.
- Added store models and methods for user upsert, org membership, site creation/list/update, quotas, and audit events.
- Added `control.OrgService` and `control.SiteService` with role checks.
- Added tests for migrations and core rows in `internal/store`.
- Added tests for owner/developer/viewer behavior and cross-org denial in `internal/control`.
- Added `deployments/dev/docker-compose.yaml` with Postgres and Keycloak.
- Added `configs/dev.postgres-keycloak.yaml` for the future Postgres/OIDC local stack.
- Updated README with where Postgres and Keycloak fit in the phases.
- Marked Phase 1 checklist items complete in `tasks.md`.

### Why

- Phase 1 needs durable product rows before deployment/runtime work begins.
- Postgres is the intended control-plane database, so local infrastructure should appear in Phase 1.
- Keycloak is needed in Phase 2, but adding it to the same compose stack now avoids an auth-infra surprise later.

### What worked

- `cd go-go-host && go test ./...` passes.
- Store tests apply migrations from an empty SQLite database and create user/org/site/quota/audit rows.
- Control tests prove an owner can create a site, a viewer cannot create a site, and a user outside an org cannot list its sites.

### What didn't work

- N/A in this step.

### What I learned

- Keeping Phase 1 SQLite-compatible made the store tests fast and local. The compose stack still prepares for Postgres wiring in the next iteration.
- The SQL uses a small number of SQLite conveniences, so a later Postgres store pass should either introduce driver-specific SQL or switch to a query builder/pgx-style implementation.

### What was tricky to build

- The first schema needed to include future tables for deployments and agents while only implementing store methods for users, orgs, sites, quotas, and audit. I kept unimplemented tables present so future phases do not need a conceptual schema reset.

### What warrants a second pair of eyes

- Review whether to keep SQLite-compatible store tests while adding Postgres integration tests, or move fully to Postgres once the compose stack is used in CI/local dev.
- Review the Keycloak compose setup before Phase 2 realm/client import automation is added.

### What should be done in the future

- Add Postgres driver support and run migrations against the compose Postgres service.
- Add Keycloak realm/client bootstrap in Phase 2.
- Add HTTP API endpoints for org/site CRUD now that the service methods exist.

### Code review instructions

- Review `go-go-host/internal/store/migrations/001_initial_schema.sql` first.
- Then review `go-go-host/internal/store/store.go`, `users.go`, `orgs.go`, `sites.go`, and `audit.go`.
- Then review `go-go-host/internal/control/services.go` and tests.
- Validate with `cd go-go-host && go test ./...`.

### Technical details

Commands run:

```bash
cd go-go-host && go mod tidy && go fmt ./... && go test ./...
cd go-go-host && go test ./...
```

## Step 9: Converted Phase 1 store to Postgres sqlc

I added sqlc to the control-plane persistence layer and converted the Phase 1 store from hand-written `database/sql` scans to generated Postgres/pgx queries. The schema is now Postgres-shaped, query files live under `internal/store/queries`, and generated code lives under `internal/store/db`.

I also changed the tests so ordinary `go test ./...` still passes without a local database, while real Postgres integration tests run when `GO_GO_HOST_TEST_DATABASE_URL` is set. I validated the integration path against the dev Docker Compose Postgres service.

### Prompt Context

**User prompt (verbatim):** "alright, let's add the tasksa and do  postgresql sqlc."

**Assistant interpretation:** Add explicit implementation tasks for sqlc/Postgres and perform the refactor now.

**Inferred user intent:** Avoid growing a hand-written CRUD layer and switch the control-plane persistence to the intended Postgres/sqlc architecture early.

**Commit (code):** N/A — changes are uncommitted.

### What I did

- Added `go-go-host/sqlc.yaml`.
- Converted `internal/store/migrations/001_initial_schema.sql` to Postgres syntax and Postgres types.
- Added sqlc query files under `internal/store/queries`.
- Generated `internal/store/db` with `sqlc generate`.
- Refactored `internal/store` to use `pgxpool` and generated `db.Queries`.
- Added migration advisory locking to avoid concurrent migration races in parallel package tests.
- Updated integration tests to use `GO_GO_HOST_TEST_DATABASE_URL` and skip when it is absent.
- Validated integration tests against the Docker Compose Postgres service.
- Updated `tasks.md` with sqlc/Postgres tasks and marked them complete.
- Updated README with sqlc/Postgres test instructions.

### Why

- The control-plane CRUD is relational and stable enough for sqlc.
- Postgres is the intended production control-plane database.
- Adding sqlc now prevents the hand-written scanning layer from expanding during upcoming API/dashboard/agent phases.

### What worked

- `sqlc generate` completed successfully.
- `cd go-go-host && go test ./...` passes without a database.
- With Docker Compose Postgres running, `GO_GO_HOST_TEST_DATABASE_URL=... go test ./internal/store ./internal/control` passes.

### What didn't work

- The first Postgres integration run hit a migration race when `./internal/store` and `./internal/control` applied migrations concurrently against the same database. I fixed this with a Postgres advisory lock around migration application.

### What I learned

- Package-level integration tests can run concurrently against the same database, so migration application must be concurrency-safe even in development.
- Keeping integration tests gated by `GO_GO_HOST_TEST_DATABASE_URL` preserves fast local unit test behavior while allowing real Postgres validation.

### What was tricky to build

- sqlc with pgx emits `pgtype.Timestamptz` for timestamps and generated table row types under the `db` package. I kept public store models separate and added conversion helpers so control services do not depend directly on generated pgx types.

### What warrants a second pair of eyes

- Review whether we want generated `db` types to propagate into services or keep the current wrapper/conversion boundary.
- Review schema choices such as text IDs with prefixes versus UUID columns.
- Review whether test databases should be isolated per package instead of sharing one database with unique test slugs.

### What should be done in the future

- Wire `go-go-hostd` to open the Postgres store from `ControlDBDSN` and apply migrations at startup.
- Add HTTP API endpoints using the `control` services.
- Add Keycloak realm/client bootstrap in Phase 2.

### Code review instructions

- Review `go-go-host/sqlc.yaml` and `go-go-host/internal/store/queries/*.sql` first.
- Review generated code under `go-go-host/internal/store/db` only for shape, not line-by-line style.
- Review `go-go-host/internal/store/store.go` for migration locking and pgxpool setup.
- Validate with:
  - `cd go-go-host && go test ./...`
  - optional Postgres: `GO_GO_HOST_TEST_DATABASE_URL='postgres://go_go_host:go_go_host_dev@127.0.0.1:55432/go_go_host?sslmode=disable' go test ./internal/store ./internal/control`

### Technical details

Commands run:

```bash
cd go-go-host && sqlc generate
cd go-go-host && go mod tidy && go fmt ./... && go test ./...
cd go-go-host && docker compose -f deployments/dev/docker-compose.yaml up -d postgres
GO_GO_HOST_TEST_DATABASE_URL='postgres://go_go_host:go_go_host_dev@127.0.0.1:55432/go_go_host?sslmode=disable' go test ./internal/store ./internal/control
cd go-go-host && docker compose -f deployments/dev/docker-compose.yaml down
```

## Step 10: Wired daemon to Postgres store and started Phase 2 dev-auth API

I continued from the Postgres/sqlc store work by wiring the daemon to open the Postgres store, apply migrations at startup, and pass the store-backed core into the HTTP API. I also started Phase 2 with a dev-auth middleware and the first identity-aware control-plane endpoints.

The new API can provision a local dev user from request headers, return `/api/v1/me`, create organizations, create sites under an organization, and list sites for an organization. Authorization still comes from local membership rows, not identity-provider roles, which preserves the product authorization model.

### Prompt Context

**User prompt (verbatim):** "Let's continue"

**Assistant interpretation:** Continue the next implementation step after Postgres/sqlc, focusing on daemon store wiring and early Phase 2 auth/API foundations.

**Inferred user intent:** Keep moving the scaffold toward a usable control plane that the dashboard and CLI can call.

**Commit (code):** N/A — changes are uncommitted.

### What I did

- Updated `go-go-hostd` to open `store.Open(ctx, cfg.ControlDBDSN)`, apply migrations, close the store on shutdown, and construct `control.NewCoreWithStore`.
- Changed default/dev config to use the local compose Postgres DSN.
- Added dev auth middleware using `X-Go-Go-Host-User`, `X-Go-Go-Host-Email`, and `X-Go-Go-Host-Name` headers, with a default `dev-user` fallback.
- Added `GET /api/v1/me`.
- Added `POST /api/v1/orgs`.
- Added `GET /api/v1/orgs/{org_id}/sites`.
- Added `POST /api/v1/orgs/{org_id}/sites`.
- Added lowercase JSON DTOs for org/site API responses.
- Added a Postgres-gated HTTP integration test for dev-auth org/site flow.
- Updated Phase 2 task status.

### Why

- The daemon previously used an in-memory/no-store core, which was no longer enough after moving to Postgres/sqlc.
- `/api/v1/me` is the first dashboard dependency for both `/app/*` and `/admin/*`.
- Basic org/site endpoints let the dashboard begin real RTK Query integration before deployment/runtime phases exist.

### What worked

- `cd go-go-host && go test ./...` passes without a database.
- With compose Postgres, `GO_GO_HOST_TEST_DATABASE_URL=... go test ./internal/store ./internal/control ./internal/httpapi` passes.
- Manual daemon smoke on an alternate port successfully called `/api/v1/me`, created an org, created a site, and listed sites.

### What didn't work

- The first manual smoke reused occupied ports from earlier daemon processes. I killed the leftover listeners and reran on another port.
- One smoke attempt reused an existing org slug in a persistent dev volume and returned an error response without an `id`; rerunning with a fresh volume and unique slug validated the flow.

### What I learned

- `go run` can leave the compiled child process listening after killing the wrapper PID in shell smoke tests. Future manual smokes should track/kill the actual listener or use `timeout` plus a unique port.
- API responses should use explicit DTOs; returning store structs produced capitalized JSON field names.

### What was tricky to build

- The HTTP routing uses Go 1.22 `ServeMux` method/path patterns inside a dev-auth sub-mux. The outer mux must mount both exact and subtree prefixes carefully so authenticated API routes reach the inner method-aware mux.

### What warrants a second pair of eyes

- Dev auth currently provisions `dev-user` by default when no header is present. That is convenient locally but must never be enabled in production config.
- The API route surface is intentionally minimal; review before dashboard code depends on names permanently.

### What should be done in the future

- Add real OIDC/JWKS validation against Keycloak.
- Add platform-admin checks for `/api/v1/admin/*`.
- Add CLI commands for `me`, `org create`, `sites list`, and `sites create` using Glazed.

### Code review instructions

- Review `cmd/go-go-hostd/main.go` for store startup/shutdown wiring.
- Review `internal/httpapi/auth.go`, `api.go`, `dto.go`, and `handler.go` for dev-auth route behavior.
- Validate with `go test ./...` and, optionally, compose-backed integration tests.

### Technical details

Commands run:

```bash
cd go-go-host && go fmt ./... && go test ./...
GO_GO_HOST_TEST_DATABASE_URL='postgres://go_go_host:go_go_host_dev@127.0.0.1:55432/go_go_host?sslmode=disable' go test ./internal/store ./internal/control ./internal/httpapi
```

## Step 11: Added OIDC bearer path and initial Glazed control-plane CLI commands

I continued both tracks requested after the Phase 2 start: real OIDC/JWKS bearer-token plumbing and Glazed CLI commands for the first control-plane APIs. The HTTP auth layer now supports dev-auth locally and an OIDC verifier path when `devAuth` is false, while the human CLI can call `/me`, create an organization, create a site, and list sites.

The CLI commands are intentionally small but follow the Glazed structure: `CommandDescription`, typed settings with `glazed` tags, shared sections, `vals.DecodeSectionInto`, and stable row output through `types.NewRow` and `gp.AddRow`.

### Prompt Context

**User prompt (verbatim):** "continue, both, check off tasks as you go."

**Assistant interpretation:** Continue both previously suggested next steps: add Glazed CLI commands and wire the real OIDC/JWKS bearer-token path, while updating the task checklist as work completes.

**Inferred user intent:** Keep implementation moving across backend auth/API and developer workflow, and keep the docmgr task plan accurate.

**Commit (code):** N/A — changes are uncommitted.

### What I did

- Added `internal/httpapi/oidc.go` using `github.com/coreos/go-oidc/v3/oidc`.
- Refactored auth middleware to use dev-auth when enabled and OIDC bearer-token verification when disabled.
- Added a missing-bearer-token test for OIDC mode.
- Added CLI JSON request helpers with dev-auth header support.
- Added `go-go-host me` Glazed command.
- Added `go-go-host org create` Glazed command.
- Added `go-go-host site create` and `go-go-host site list` Glazed commands.
- Wired the new CLI commands into the root command.
- Ran ordinary tests and compose-backed integration/smoke tests.
- Checked off completed Phase 2 and Phase 6 tasks.

### Why

- The dashboard and CLI both need a real identity endpoint before deeper product workflows can be useful.
- OIDC belongs in Phase 2, and the CLI needs to exercise the same control-plane APIs the dashboard will use.

### What worked

- `cd go-go-host && go test ./...` passes.
- Compose-backed integration tests for store/control/httpapi pass with `GO_GO_HOST_TEST_DATABASE_URL`.
- Manual CLI smoke against a running daemon passed for `me`, `org create`, `site create`, and `site list` with JSON output.

### What didn't work

- N/A for code after the fixes. Earlier port/listener leftovers were cleaned up.

### What I learned

- The Glazed command helpers are enough to add CRUD commands quickly once the API endpoints exist.
- OIDC verification can be lazy-initialized on the first bearer request, which avoids daemon startup depending on Keycloak availability.

### What was tricky to build

- The OIDC path should not force discovery during dev-auth mode or daemon startup. Lazy verifier initialization keeps dev mode fast and avoids failing local runs when Keycloak is not up.
- The CLI JSON output shape from Glazed is a list of rows, so smoke scripts need to read the first row for IDs.

### What warrants a second pair of eyes

- The OIDC test currently covers missing bearer rejection, not a full signed-token JWKS verification matrix. Add issuer/audience/signature tests when a test JWKS helper or Keycloak test container is in place.
- Command names are currently singular (`org`, `site`) while older task text used plural (`orgs`, `sites`). Decide before documenting CLI public API.

### What should be done in the future

- Add `go-go-host org list`.
- Add token storage/login flow once browser OAuth is implemented.
- Add full OIDC/JWKS tests for invalid issuer, audience, and signature.
- Add platform-admin route guards.

### Code review instructions

- Review `internal/httpapi/auth.go` and `internal/httpapi/oidc.go` for auth behavior.
- Review `cmd/go-go-host/cmds/me.go`, `org.go`, `site.go`, and `support.go` for Glazed command patterns.
- Validate with `go test ./...` and optional compose-backed CLI smoke.

### Technical details

Commands run:

```bash
cd go-go-host && go get github.com/coreos/go-oidc/v3/oidc@latest && go mod tidy
cd go-go-host && go fmt ./... && go test ./...
GO_GO_HOST_TEST_DATABASE_URL='postgres://go_go_host:go_go_host_dev@127.0.0.1:55432/go_go_host?sslmode=disable' go test ./internal/store ./internal/control ./internal/httpapi
go run ./cmd/go-go-host me --api-url http://127.0.0.1:18083 --dev-user cli-smoke --output json
go run ./cmd/go-go-host org create --api-url http://127.0.0.1:18083 --dev-user cli-smoke --slug cli-smoke-org --name 'CLI Smoke Org' --output json
go run ./cmd/go-go-host site create --api-url http://127.0.0.1:18083 --dev-user cli-smoke --org-id "$ORG" --slug cli-smoke-site --name 'CLI Smoke Site' --output json
go run ./cmd/go-go-host site list --api-url http://127.0.0.1:18083 --dev-user cli-smoke --org-id "$ORG" --output json
```

## Step 12: Added org listing and bearer-token CLI plumbing

I continued the CLI/API foundation by adding an organization listing endpoint and command. The API now exposes `GET /api/v1/orgs` for the current user's organization memberships, and the CLI has `go-go-host org list` to emit those memberships as Glazed rows.

I also added a `--bearer-token` flag to implemented human CLI commands. This does not replace the future browser login/token storage workflow, but it lets the same commands exercise the non-dev bearer-token path once Keycloak tokens are available.

### Prompt Context

**User prompt (verbatim):** "continue."

**Assistant interpretation:** Continue incremental Phase 2 and CLI work, focusing on still-open low-level tasks before moving to runtime or dashboard.

**Inferred user intent:** Keep building usable control-plane API/CLI surface area while maintaining tests and checklist accuracy.

**Commit (code):** N/A — changes are uncommitted.

### What I did

- Added `GET /api/v1/orgs` returning the current user's organization memberships.
- Added `go-go-host org list` as a Glazed command.
- Added `--bearer-token` to implemented human CLI commands.
- Updated CLI helper methods to send either dev-auth headers or bearer tokens.
- Ran normal tests, Postgres-backed integration tests, and a compose-backed CLI smoke for `org create` and `org list`.
- Updated `tasks.md`, diary, and changelog.

### Why

- A dashboard needs org membership listing before org-scoped pages can work.
- The CLI needed a way to exercise non-dev auth before a full browser login flow exists.

### What worked

- `go test ./...` passes.
- Postgres-backed integration tests pass.
- `go-go-host org create --output json` and `go-go-host org list --output json` worked against a live daemon.

### What didn't work

- N/A in this step after the small compile fix in the CLI helper argument list.

### What I learned

- Returning memberships rather than bare organizations is more useful for `/api/v1/orgs`, because the dashboard needs each org's role immediately.

### What was tricky to build

- The CLI helper signature changed to include both dev user and bearer token, so all helper callers needed to be kept consistent.

### What warrants a second pair of eyes

- Decide whether public CLI command names should be singular (`org`, `site`) or plural (`orgs`, `sites`) before writing final user docs.

### What should be done in the future

- Add browser OAuth/login and local token storage so users do not paste `--bearer-token` manually.
- Add admin-only API surfaces and platform-admin guards.

### Code review instructions

- Review `internal/httpapi/admin.go` for `GET /api/v1/orgs` response shape.
- Review `cmd/go-go-host/cmds/org.go` and `support.go` for CLI auth behavior.
- Validate with ordinary and Postgres-backed tests.

### Technical details

Commands run:

```bash
cd go-go-host && go fmt ./... && go test ./...
GO_GO_HOST_TEST_DATABASE_URL='postgres://go_go_host:go_go_host_dev@127.0.0.1:55432/go_go_host?sslmode=disable' go test ./internal/store ./internal/control ./internal/httpapi
go run ./cmd/go-go-host org create --api-url http://127.0.0.1:18084 --dev-user cli-smoke-2 --slug cli-smoke-org-2 --name 'CLI Smoke Org 2' --output json
go run ./cmd/go-go-host org list --api-url http://127.0.0.1:18084 --dev-user cli-smoke-2 --output json
```

## Step 13: Added CLI login/config persistence

I added the first local CLI configuration workflow. `go-go-host login` now stores an API URL plus either a dev user or bearer token in a YAML config file. Implemented CLI commands load that file when flags are omitted, so developers can run `go-go-host me`, `go-go-host org create`, and related commands without repeating `--api-url` and auth flags each time.

This is still a bridge command, not a browser OAuth flow. It creates the persistence location and command shape that a future PKCE/browser login can reuse.

### Prompt Context

**User prompt (verbatim):** "continue."

**Assistant interpretation:** Continue closing Phase 2/CLI gaps, especially the local token/session config task.

**Inferred user intent:** Make the current CLI workflows more realistic and reduce repeated local-dev flags while preserving the path to real OAuth later.

**Commit (code):** N/A — changes are uncommitted.

### What I did

- Added `cmd/go-go-host/cmds/cli_config.go`.
- Added `cmd/go-go-host/cmds/login.go`.
- Wired `go-go-host login` into the root command.
- Updated `me`, `org`, and `site` commands to resolve API URL/dev user/bearer token from saved config unless flags override them.
- Added `GO_GO_HOST_CLI_CONFIG` override for tests/smoke scripts.
- Ran tests and a compose-backed CLI config smoke.
- Updated `tasks.md`, diary, and changelog.

### Why

- Phase 6 called for a local token/session config file.
- Requiring `--api-url` and `--dev-user` on every command is too clumsy for the dashboard/CLI bring-up loop.

### What worked

- `go-go-host login --api-url ... --dev-user ... --output json` wrote a config file.
- Follow-up commands used the saved config successfully without explicit `--api-url` or `--dev-user`.
- `go test ./...` passes.

### What didn't work

- N/A in this step.

### What I learned

- A simple YAML config file is enough for Phase 2 dev auth and bearer-token smoke tests. Real browser OAuth can later replace how the bearer token gets into this file.

### What was tricky to build

- Glazed defaults still provide a default API URL, so the config resolver needs to treat the static default as overridable by saved config unless the user explicitly provides a different API URL. This is acceptable for now but may warrant a cleaner flag/config precedence model later.

### What warrants a second pair of eyes

- Review config file permissions and whether bearer tokens should be stored in plaintext YAML for local development only.
- Decide whether `GO_GO_HOST_CLI_CONFIG` should remain a documented test override.

### What should be done in the future

- Implement browser OAuth/PKCE in `go-go-host login` and store refresh/session material safely.
- Add `go-go-host logout` or config-clear command.

### Code review instructions

- Review `cmd/go-go-host/cmds/cli_config.go` and `login.go`.
- Validate with `GO_GO_HOST_CLI_CONFIG=$(mktemp -d)/config.yaml go-go-host login ...` followed by `go-go-host me`.

### Technical details

Commands run:

```bash
cd go-go-host && go fmt ./... && go test ./...
GO_GO_HOST_CLI_CONFIG=$CFG go run ./cmd/go-go-host login --api-url http://127.0.0.1:18085 --dev-user config-smoke --output json
GO_GO_HOST_CLI_CONFIG=$CFG go run ./cmd/go-go-host me --output json
GO_GO_HOST_CLI_CONFIG=$CFG go run ./cmd/go-go-host org create --slug config-smoke-org --name 'Config Smoke Org' --output json
```

## Step 14: Started Phase 3 by refactoring the hosted Goja runtime

I started Phase 3 after committing the completed control-plane/auth/CLI foundation. This step copied the reusable `goja-site` web, UI DSL, and DB guard components into `go-go-host` and introduced `internal/runtime.SiteRuntime`, which is the per-site runtime object the supervisor will manage later.

The fixture runtime now loads JavaScript from disk, registers Express-style routes, renders `ui.dsl` HTML, uses a per-site SQLite database through preconfigured `database` and `db` modules, and serves static assets under `/assets`. The tests prove the first important safety properties: the fixture renders, `database.configure()` is disabled, and unrestricted `fs`/`exec` are not available by default.

### Prompt Context

**User prompt (verbatim):** "continue committing at appropriate intervals, and keeping a frequent diary"

**Assistant interpretation:** Commit the accumulated work before starting Phase 3, then continue implementation in smaller reviewable chunks and update the diary frequently.

**Inferred user intent:** Keep the repository history usable and keep a written trail as runtime refactoring begins.

**Commit (code):** `cc20654cf0e787c9d22f54bec39cec0e82ae3a03` — "Scaffold control plane with Postgres sqlc and CLI"

### What I did

- Committed the previous scaffold/control-plane/sqlc/auth/CLI work.
- Copied reusable runtime support packages from `goja-site` into `internal/sitejs`:
  - `internal/sitejs/web`
  - `internal/sitejs/uidsl`
  - `internal/sitejs/dbguard`
- Added `internal/runtime/runtime.go` with `Spec`, `CapabilitySet`, and `SiteRuntime`.
- Added a fixture hosted site under `internal/runtime/testdata/sites/hello`.
- Added runtime tests for rendering, disabled database reconfiguration, and unavailable `fs`/`exec` modules.
- Ran `go mod tidy`, `go fmt ./...`, and `go test ./...`.
- Updated Phase 3 task checkboxes.

### Why

- The platform needs a per-site runtime object before supervisor, deployment activation, or Host-header routing can be implemented.
- Copy/refactoring the runtime support packages keeps `go-go-host` independent enough to evolve the hosted runtime policy without changing the original prototype project.

### What worked

- The fixture route renders `Hello from go-go-host` through the refactored runtime.
- The hosted database module rejects `configure()` because it is preconfigured by Go.
- `require("exec")` and `require("fs")` fail by default.
- `go test ./...` passes.

### What didn't work

- N/A after adjusting the runtime test to call Goja through the real owner with a `*goja.Runtime` callback.

### What I learned

- The `goja-site` runtime support packages copy cleanly because non-test files do not import the original module path.
- The safest default runtime middleware list is very small: `path`, plus `time/timer` when enabled. Leaving out `fs` and `exec` gives us the expected hosted safety baseline.

### What was tricky to build

- The test for unavailable modules needed to exercise CommonJS `require()` inside the runtime owner, not by directly inspecting Go objects outside the owner. This keeps the test aligned with hosted script behavior.

### What warrants a second pair of eyes

- Review whether copying `internal/sitejs/*` is the preferred long-term strategy or whether these should become a shared module later.
- Review the default capability policy before exposing user-authored deployments.

### What should be done in the future

- Add `HealthCheck(ctx)` for smoke route validation.
- Add runtime supervisor and Host-header routing in Phase 4.
- Eventually add a scoped asset/filesystem capability instead of unrestricted `fs`.

### Code review instructions

- Start with `internal/runtime/runtime.go`.
- Then review copied support packages only for package-path changes and hosted-policy assumptions.
- Validate with `cd go-go-host && go test ./internal/runtime ./...`.

### Technical details

Commands run:

```bash
cd go-go-host && go test ./...
git add Makefile README.md go.mod go.sum sqlc.yaml configs deployments internal cmd ttmp
git commit -m "Scaffold control plane with Postgres sqlc and CLI"
mkdir -p internal/sitejs/web internal/sitejs/uidsl internal/sitejs/dbguard
cp ../2026-05-03--goja-hosting-site/pkg/web/*.go internal/sitejs/web/
cp ../2026-05-03--goja-hosting-site/pkg/uidsl/*.go internal/sitejs/uidsl/
cp ../2026-05-03--goja-hosting-site/pkg/dbguard/*.go internal/sitejs/dbguard/
rm -f internal/sitejs/web/*_test.go internal/sitejs/uidsl/*_test.go internal/sitejs/dbguard/*_test.go
go mod tidy && go fmt ./... && go test ./...
```

## Step 15: Completed the Phase 3 runtime health check

I filled the remaining Phase 3 runtime gap by adding `SiteRuntime.HealthCheck(ctx)`. The method performs an in-process GET request against a configured health path, defaulting to `/`, and treats non-2xx/3xx responses as unhealthy.

This gives the future supervisor a simple activation gate: build a new runtime, call `HealthCheck`, then swap traffic only if the new runtime is serving successfully.

### Prompt Context

**User prompt (verbatim):** (same as Step 14)

**Assistant interpretation:** Continue Phase 3 in a small commit-sized increment and keep the diary up to date.

**Inferred user intent:** Finish the per-site runtime object before moving to supervisor/host routing.

**Commit (code):** `f69b5a394de90c3867fc8d8cd4b19ef2f4b49551` — "Add hosted Goja site runtime"

### What I did

- Added `HealthPath` to runtime `Spec`.
- Added `SiteRuntime.HealthCheck(ctx)`.
- Added `TestRuntimeHealthCheck`.
- Marked the health-check Phase 3 task complete.

### Why

- Phase 4 activation needs a health check before traffic is swapped to a new runtime.
- The design explicitly calls for a smoke-route check or script-load check.

### What worked

- `go test ./internal/runtime ./...` passes.
- The fixture runtime health check succeeds against `/`.

### What didn't work

- I briefly started a custom response recorder, then simplified to the standard library `httptest.NewRecorder`, which is clearer and already correct.

### What I learned

- The health check can stay fully in-process for now; no socket or daemon listener is needed to validate a `SiteRuntime` before supervisor activation.

### What was tricky to build

- The health-check request needs a valid URL, but host does not matter yet because this is inside one runtime handler. The future supervisor will own host-level routing.

### What warrants a second pair of eyes

- Review whether the default health path should remain `/` or come from deployment manifest validation in Phase 5.

### What should be done in the future

- Feed manifest smoke-test path into `RuntimeSpec.HealthPath`.
- Use `HealthCheck` in `Supervisor.Activate`.

### Code review instructions

- Review `internal/runtime/runtime.go` around `HealthCheck`.
- Validate with `go test ./internal/runtime`.

### Technical details

Commands run:

```bash
go fmt ./internal/runtime && go test ./internal/runtime ./...
```

## Step 16: Started Phase 4 with runtime supervisor and Host-header routing

I started Phase 4 by adding an in-memory runtime supervisor. It owns maps by site ID and normalized host, activates runtimes by building and health-checking a new `SiteRuntime` before swapping traffic, and exposes an `http.Handler` that dispatches public requests by Host header.

The tests cover the important activation invariants: two fixture sites can serve concurrently by host, unknown hosts return 404, failed activation does not replace the currently serving runtime, and stopped runtimes are removed from routing.

### Prompt Context

**User prompt (verbatim):** "continue"

**Assistant interpretation:** Move from Phase 3 runtime object into Phase 4 supervisor/host routing and continue updating tasks and diary.

**Inferred user intent:** Build toward the first real hosted-site smoke loop.

**Commit (code):** N/A — current supervisor changes are uncommitted at this diary entry.

### What I did

- Added `internal/runtime/supervisor.go`.
- Added `internal/runtime/supervisor_test.go`.
- Added runtime status model: `starting`, `ready`, `failed`, `stopped`, `draining`.
- Added host normalization and supervisor summary.
- Implemented `Activate`, `Stop`, `GetByHost`, `Status`, `Summary`, and `ServeHTTP`.
- Updated runtime fixture helper reuse between runtime and supervisor tests.
- Updated Phase 4 task checkboxes.

### Why

- The product needs a dynamic Host-header router before deployments can activate live sites.
- Activation must build and health-check a new runtime before swapping, so bad deployments cannot take down a good runtime.

### What worked

- `go test ./internal/runtime ./...` passes.
- Two fixture sites route independently by Host header.
- A bad activation records failed status but leaves the previous runtime serving.

### What didn't work

- N/A in this step.

### What I learned

- The supervisor can remain store-independent for this step. Persistence of status transitions can be layered in after the in-memory invariants are correct.

### What was tricky to build

- Failed activation status and currently serving runtime are separate concepts: after a failed activation, the site should report the failed attempt but continue serving the previous runtime. Tests now encode that invariant.

### What warrants a second pair of eyes

- Review whether `Status(siteID)` should report the current serving runtime status or the last activation attempt. It currently reports the last activation attempt, which is useful for deployment diagnostics but may need separate fields later.

### What should be done in the future

- Add `Restart(ctx, siteID)` by retaining the last successful `Spec`.
- Persist runtime status transitions to the control DB.
- Add public request context and counters.
- Wire supervisor into the daemon after deployment records exist.

### Code review instructions

- Review `internal/runtime/supervisor.go` for activation swap ordering.
- Review `internal/runtime/supervisor_test.go` for host routing and failure invariants.
- Validate with `go test ./internal/runtime`.

### Technical details

Commands run:

```bash
go fmt ./internal/runtime && go test ./internal/runtime ./...
```

## Step 17: Added supervisor restart, counters, and runtime API endpoints

I continued Phase 4 by extending the supervisor beyond basic Host-header routing. It now retains the last successful runtime spec for restart, tracks request/error counters, and is reachable through the control-plane core so HTTP API handlers can report runtime status.

This step also added the first runtime-related API endpoints: org members can query a site's runtime status, and platform admins can query the global runtime summary. These endpoints are still backed by in-memory supervisor state; persisted runtime transitions remain a later Phase 4 item.

### Prompt Context

**User prompt (verbatim):** "go ahead"

**Assistant interpretation:** Continue implementing the remaining Phase 4 supervisor/API pieces.

**Inferred user intent:** Move toward a daemon-wired runtime surface that can be inspected from API/dashboard/CLI.

**Commit (code):** N/A — current changes are uncommitted at this diary entry.

### What I did

- Added `Supervisor.Restart(ctx, siteID)` using retained successful specs.
- Added per-site request and error counters in runtime status.
- Added `statusResponseWriter` to observe final response status codes.
- Added `Supervisor` to `control.Core`.
- Added `GET /api/v1/sites/{site_id}/runtime`.
- Added `GET /api/v1/admin/runtimes/summary` with platform-admin check.
- Updated Phase 4 task checkboxes.
- Ran `go test ./...`.

### Why

- Restart and counters are needed for operational runtime management.
- Dashboard and CLI need runtime status endpoints before they can show useful runtime health.

### What worked

- `go test ./...` passes.
- Existing supervisor tests still pass, and new restart/counter tests cover the added behavior.

### What didn't work

- N/A in this step.

### What I learned

- It is useful to keep supervisor in `control.Core` now, even before it is fully backed by deployment records. That gives HTTP handlers a stable dependency path.

### What was tricky to build

- Runtime status after a failed activation currently represents the last activation attempt while the previous runtime may still serve traffic. Counters remain associated by site ID, so they carry over across replacement attempts.

### What warrants a second pair of eyes

- Review whether admin runtime summary should be available only to `platform_admin` rows or also to org owners filtered by org.
- Review whether request counters should reset on restart or persist per site until process restart. They currently persist per site in memory.

### What should be done in the future

- Persist runtime status transitions to the control DB.
- Add request context fields and structured runtime logs.
- Wire deployment activation to `Supervisor.Activate`.

### Code review instructions

- Review `internal/runtime/supervisor.go` and `internal/httpapi/runtime.go`.
- Validate with `go test ./internal/runtime ./internal/httpapi ./...`.

### Technical details

Commands run:

```bash
go fmt ./... && go test ./...
```

## Step 18: Finished Phase 4 persistence, reconciliation, and daemon fallback routing

I completed the remaining Phase 4 runtime-supervisor tasks. Runtime status transitions are now persisted to the control DB through a supervisor recorder, daemon startup reconciles stale runtime statuses to stopped, public requests receive platform context fields before entering JavaScript, and the main HTTP handler falls back to the supervisor for public hosted-site traffic.

This does not yet provide deployment-driven activation, but the runtime and supervisor layer is now complete enough for the next phase to create deployment records and call `Supervisor.Activate`.

### Prompt Context

**User prompt (verbatim):** "do full phase 4, don't stop before you are done. Commit at appropriate intervals and keep a diary"

**Assistant interpretation:** Complete all remaining Phase 4 tasks, commit in coherent chunks, and keep the ticket diary updated.

**Inferred user intent:** Finish the runtime supervisor layer before moving to deployment pipeline work.

**Commit (code):** N/A — current persistence/fallback changes are uncommitted at this diary entry.

### What I did

- Added `runtime_status` migration and sqlc queries.
- Regenerated sqlc code.
- Added store methods for runtime status upsert/list/get and stale reconciliation.
- Added a control-layer runtime status recorder and passed it into the supervisor.
- Persisted supervisor status transitions and request/error counter updates.
- Added daemon startup stale-runtime reconciliation.
- Added platform context support in copied `web` package and injected request ID/org/site/deployment/host before hosted JavaScript runs.
- Added HTTP fallback routing from API mux to supervisor for public site traffic.
- Added tests for runtime status persistence/reconciliation and HTTP fallback routing.
- Updated Phase 4 checklist to complete.

### Why

- Phase 4 requires supervisor state to be visible and recoverable at the control-plane level.
- Public hosted-site traffic must reach the supervisor after API/app/admin routes are considered.
- JavaScript handlers need safe platform context for debugging and request correlation.

### What worked

- `go test ./...` passes.
- Compose-backed Postgres integration tests pass for store/control/httpapi.
- HTTP fallback routing test proves public Host-header traffic can reach an active supervisor runtime.

### What didn't work

- N/A after iterative implementation.

### What I learned

- Keeping runtime persistence behind a small `StatusRecorder` interface avoids coupling `internal/runtime` directly to the control-plane store package.
- Fallback routing can be implemented without changing the public API mux shape, but future production routing may want a more explicit top-level router.

### What was tricky to build

- Request counters are updated after serving a response, while status persistence is asynchronous for request counts. This avoids blocking public traffic on control DB writes but means counters are eventually consistent.

### What warrants a second pair of eyes

- Review fallback routing semantics: API 404s with bodies should not be converted into hosted-site requests, while unknown public paths should fall through to supervisor.
- Review whether runtime counter persistence should be batched rather than written after every request.

### What should be done in the future

- Phase 5 should create deployment records and activate them via `Supervisor.Activate`.
- Add structured logs for public requests with platform context.
- Consider batching runtime metrics writes.

### Code review instructions

- Review `internal/runtime/supervisor.go`, `internal/store/runtime_status.go`, `internal/httpapi/fallback.go`, and `internal/sitejs/web/platform.go`.
- Validate with `go test ./...` and compose-backed integration tests.

### Technical details

Commands run:

```bash
sqlc generate
go fmt ./... && go test ./...
GO_GO_HOST_TEST_DATABASE_URL='postgres://go_go_host:go_go_host_dev@127.0.0.1:55432/go_go_host?sslmode=disable' go test ./internal/store ./internal/control ./internal/httpapi
```

## Step 19: Implemented Phase 5 deployment bundle pipeline

I implemented the first full deployment pipeline slice: bundle validation, immutable storage/unpack, deployment records, upload/list/detail/activate/rollback APIs, dry-run runtime validation, and an end-to-end deploy-and-serve integration test.

### Prompt Context

**User prompt (verbatim):** "do phase 5, keep a diary, commit at appropriat eintervals"

**Assistant interpretation:** Complete Phase 5, keep the ticket diary current, validate the work, and commit a coherent implementation slice.

**Inferred user intent:** Move from a manually activatable runtime supervisor to an API-driven bundle deployment workflow.

**Commit (code):** N/A — this entry was written before committing the Phase 5 slice.

### What I did

- Added `internal/deploy` bundle validation and storage helpers.
- Defined the initial `go-go-host.json` manifest shape.
- Implemented tar.gz and zip archive readers.
- Rejected unsafe paths: absolute paths, parent traversal, unsafe components, symlinks/hardlinks, and hidden metadata components.
- Added quota-oriented file/byte validation.
- Parsed requested capabilities and rejected capabilities outside the site policy.
- Added deploy path/channel policy checks for the current bundle validator.
- Stored bundle archives under the configured data directory and unpacked deployments under `data/sites/<site-id>/deployments/<deployment-id>`.
- Added dry-run runtime load and optional smoke route validation before a deployment becomes `validated`.
- Added deployment store queries/methods for create/list/get/status/activation/previous deployment.
- Added `DeploymentService` for upload, list, get, activate, and rollback.
- Added API endpoints:
  - `POST /api/v1/sites/{site_id}/deployments`
  - `GET /api/v1/sites/{site_id}/deployments`
  - `GET /api/v1/deployments/{deployment_id}`
  - `POST /api/v1/deployments/{deployment_id}/activate`
  - `POST /api/v1/sites/{site_id}/rollback`
- Added tests for invalid paths, missing manifest, oversized bundle, forbidden capabilities, and an integration deploy/activate/serve flow.
- Updated the Phase 5 checklist to complete.

### Why

Phase 5 is the handoff point from runtime mechanics to a deployable product workflow. The control plane now has enough deployment records and validation behavior for a user to upload an immutable site bundle, activate it, and serve it by Host header.

### What worked

- `go test ./...` passes.
- Compose-backed Postgres integration tests pass for store/control/httpapi.
- The integration test uploads a hello bundle, activates it, and requests the deployed site through the public Host-header fallback path.

### What didn't work

- The first integration bundle used a prototype `require('web')` style that is not the current runtime API, so dry-run validation correctly rejected it with `Invalid module`. I changed the bundle fixture to use the copied runtime's supported `require('express')` API.
- The public fallback router originally avoided fallback when the stdlib mux wrote a 404 body. Public `/` requests therefore produced an empty response instead of reaching the supervisor. I tightened the fallback rule to forward 404s for non-control paths while preserving `/api`, `/app`, `/admin`, and health routes.

### What I learned

- Dry-run runtime validation is useful immediately: it caught an invalid script before the deployment could be activated.
- Deployment activation should remain a separate explicit step after upload/validation, because it gives clients a clear validation report before traffic changes.

### What was tricky to build

- Deployment creation needs a stable deployment ID before archive/unpack paths are known. I currently create an uploaded placeholder row, validate/store the bundle, then update the row with artifacts and validation status.
- The current rollback implementation activates the previous validated/superseded deployment rather than mutating bundle contents.

### What warrants a second pair of eyes

- Review whether creating a placeholder deployment before validation is the desired long-term transaction shape.
- Review whether request-counter persistence should be batched before serious traffic.
- Review whether `go-go-host.json` should support richer routes, build metadata, and explicit capability configuration.

### What should be done in the future

- Add CLI commands for deploy/list/show/activate/rollback in Phase 6.
- Add dashboard surfaces for validation reports in Phase 7.
- Add deployment artifact garbage collection and stricter bundle content type sniffing.
- Add a richer site capability policy table rather than the current safe default capability map.

### Code review instructions

- Start with `internal/deploy/bundle.go` for validation behavior.
- Then review `internal/control/deployments.go` for service orchestration.
- Then review `internal/httpapi/deployments.go` and `internal/httpapi/deployments_integration_test.go` for the API contract and smoke path.
- Validate with `go test ./...` and the compose-backed Postgres integration test command.

### Technical details

Commands run:

```bash
go fmt ./... && go test ./...
GO_GO_HOST_TEST_DATABASE_URL='postgres://go_go_host:go_go_host_dev@127.0.0.1:55432/go_go_host?sslmode=disable' go test ./internal/store ./internal/control ./internal/httpapi
```

## Step 20: Added Phase 6 deployment/runtime CLI commands

I continued into Phase 6 and added the first developer workflow CLI slice on top of the Phase 5 deployment APIs.

### Prompt Context

**User prompt (verbatim):** "continue"

**Assistant interpretation:** Continue from completed Phase 5 into the next checklist items, starting Phase 6 CLI workflow support.

**Inferred user intent:** Make the deployment pipeline usable from the human `go-go-host` CLI, then commit a focused slice.

**Commit (code):** N/A — this entry was written before committing the CLI slice.

### What I did

- Added `go-go-host deploy` as a Glazed command for multipart bundle upload.
- Added `go-go-host deployments list`.
- Added `go-go-host deployments show`.
- Added `go-go-host deployments activate`.
- Added `go-go-host rollback`.
- Added `go-go-host site runtime` and `go-go-host sites runtime` alias support.
- Added `orgs` and `sites` aliases for the existing singular command groups.
- Added multipart upload support to the shared CLI HTTP helper.
- Improved CLI HTTP error handling so non-2xx responses include API validation/authz response bodies.
- Added `WithLong` examples for the new deployment/runtime commands.
- Updated Phase 6 checklist items for the implemented CLI deployment workflow.

### Why

Phase 5 made deployments possible through the API, but users still needed raw HTTP clients. This slice makes the core developer loop available through Glazed commands: create/list sites using existing commands, upload a bundle, list/show deployments, activate a deployment, inspect runtime status, and roll back.

### What worked

- `go test ./...` passes after adding the new CLI commands.
- The commands reuse the existing config resolution logic for API URL, dev user, and bearer token.

### What didn't work

- I initially tried to merge nested validation/manifest fields into a Glazed row using a nonexistent `Row.Merge` helper. I changed the code to copy ordered-map pairs explicitly.

### What I learned

- Glazed rows are ordered maps, so adding extra nested fields requires explicit `Set` calls or constructing a row with all pairs up front.

### What was tricky to build

- Multipart upload needed a separate helper rather than forcing the existing JSON helper to handle file bodies.
- Nested validation reports are emitted as row fields for JSON/YAML output, while table output remains best for the stable scalar deployment columns.

### What warrants a second pair of eyes

- Review CLI flag names before public docs harden. The current commands use `--site-id` and `--deployment-id`; older checklist wording still mentions slug/positional forms.
- Decide whether `go-go-host deploy` should accept a directory and build a bundle locally, or keep requiring an already-created `.tar.gz`/`.zip` archive for v1.

### What should be done in the future

- Add CLI smoke tests against an httptest server.
- Add embedded help workflow pages for deploy and rollback.
- Add agent/audit CLI commands after those APIs exist.

### Code review instructions

- Review `cmd/go-go-host/cmds/deployments.go` and `cmd/go-go-host/cmds/support.go` first.
- Then review `cmd/go-go-host/cmds/site.go` and `cmd/go-go-host/root.go` for command wiring.
- Validate with `go test ./...`.

### Technical details

Commands run:

```bash
go fmt ./cmd/go-go-host ./cmd/go-go-host/cmds
go test ./...
```

## Step 21: Added embedded CLI workflow help pages

I added embedded Glazed help pages for the first user-facing CLI workflows: login/config, organization/site creation, deployment, rollback, and the current agent setup preview.

### Prompt Context

**User prompt (verbatim):** "continue"

**Assistant interpretation:** Continue Phase 6 after the initial deployment CLI command slice and address the embedded help workflow checklist item.

**Inferred user intent:** Keep moving through the developer CLI experience instead of stopping after command wiring.

**Commit (code):** N/A — this entry was written before committing the help-page slice.

### What I did

- Added `login-and-config.md`.
- Added `create-site-workflow.md`.
- Added `deploy-workflow.md`.
- Added `rollback-workflow.md`.
- Added `agent-setup.md` as a preview/status page for the future agent workflow.
- Marked the embedded help workflow checklist item complete.

### Why

The CLI now has enough commands for an end-to-end developer workflow, but discoverability matters. Embedded help pages make the workflow available through `go-go-host help` without requiring external docs.

### What worked

- The existing `go:embed *.md` and help system wiring already picked up additional Markdown pages.
- `go test ./...` passes with the new embedded docs.

### What didn't work

- N/A.

### What I learned

- The current root command was already using the standard Glazed help-system integration, so adding workflow docs only required Markdown files with proper frontmatter.

### What was tricky to build

- The agent workflow page needed to be explicit that the agent CLI/API is still a preview, so users do not expect non-existent `agents list` behavior yet.

### What warrants a second pair of eyes

- Review command names in the docs before public release. The current v1 commands prefer IDs (`--site-id`, `--deployment-id`) over slug-based positional shortcuts.

### What should be done in the future

- Add real agent and audit help pages when those APIs and commands exist.
- Add a CLI smoke test that verifies help-page lookup and JSON output for common commands.

### Code review instructions

- Review files under `cmd/go-go-host/doc/*.md`.
- Validate with `go test ./...` and `go-go-host help deploy-workflow` once the binary is built.

## Step 22: Added CLI smoke tests plus audit and agent APIs/commands

I added smoke coverage for the CLI HTTP helpers, then implemented the initial audit and agent API surface with matching Glazed CLI commands.

### Prompt Context

**User prompt (verbatim):** "go ahead with smoke tests, then implement the audit/agent APIs and their CLI commands"

**Assistant interpretation:** First add CLI smoke tests for the existing workflow, then add the missing backend APIs and CLI commands for audit and agents.

**Inferred user intent:** Close the remaining Phase 6 gaps that were blocked on missing APIs.

**Commit (code):** N/A — this entry was written before committing this slice.

### What I did

- Added CLI HTTP smoke tests using `httptest`:
  - JSON auth request helper.
  - Multipart bundle upload helper.
  - Error-body propagation for failed API responses.
- Added sqlc agent queries and regenerated store code.
- Added audit filtered query support.
- Added store models/methods for agents and agent site grants.
- Added control services for agents and audit.
- Added HTTP APIs:
  - `GET /api/v1/orgs/{org_id}/agents`
  - `POST /api/v1/orgs/{org_id}/agents`
  - `POST /api/v1/orgs/{org_id}/agents/{agent_id}/revoke`
  - `GET /api/v1/orgs/{org_id}/audit`
- Added integration test covering agent create/list and audit filtering.
- Added CLI commands:
  - `go-go-host agents list`
  - `go-go-host agents create`
  - `go-go-host audit list`
- Updated Phase 6 checklist entries for agents, audit, and CLI smoke tests.

### Why

The CLI checklist still had explicit agent and audit commands, but the server did not expose enough API to support them. This slice adds the minimal useful control-plane API and CLI workflow while preserving the larger future agent enrollment/deploy-run work for later phases.

### What worked

- `go test ./...` passes.
- Compose-backed Postgres integration tests pass for store/control/httpapi.
- The agent creation API emits an audit event that is visible through the audit list API.

### What didn't work

- sqlc generated positional names (`Column2`, etc.) for the first filtered audit query because the SQL used raw positional parameters. I adapted the store wrapper to those generated names rather than spending time on a cosmetic query rewrite.

### What I learned

- The existing schema had enough agent tables to support a useful initial list/create/revoke API without a migration.
- Audit list filtering can stay simple for v1: org-scoped access plus resource/actor/action filters covers the immediate CLI need.

### What was tricky to build

- Avoiding over-scoping agents: the current API intentionally creates and lists agent records only. Full enrollment keys, grants, and deploy-run token issuance should remain a later focused phase.

### What warrants a second pair of eyes

- Review agent API authorization. Currently org owners and developers can create/revoke agents, while viewers can list.
- Review audit filters and whether the CLI should expose `--site-id` as an alias for resource filtering once resource typing is more consistent.

### What should be done in the future

- Add full agent enrollment keys and grant editor APIs.
- Add deploy-run APIs that use agent grants for path/channel restrictions.
- Add command-level smoke tests that execute Cobra/Glazed output modes end-to-end, beyond HTTP helper smoke tests.

### Code review instructions

- Review `internal/control/agents.go`, `internal/httpapi/agents_audit.go`, `cmd/go-go-host/cmds/agents.go`, and `cmd/go-go-host/cmds/audit.go`.
- Validate with `go test ./...` and the compose-backed integration command.

## Step 23: Wrote Phase 7 dashboard affordances and component-system design

I created a new Phase 7 dashboard design document for the `/app` user/org developer dashboard. The document is written as an intern-facing implementation guide and includes product affordances, backend evidence, page descriptions, ASCII screenshots, component taxonomy, Storybook/MSW requirements, RTK Query endpoint sketches, testing strategy, and phased implementation guidance.

### Prompt Context

**User prompt (verbatim):** "Ok, for phase 7, we want a new design document which lists all the affordances the site should provide, and then a set of page descriptions with ascii screenshots, and then a list of the widgets for the component system (atomm/molecules/organisms, building on top of go-go-os-core ) with one directory per widget, with a storybook stories page as well. We always want to have storybook stories (using msw / fake store when relevant), for the pages too. Create  a detailed analysis / design / implementation guide that is for a new intern, explaining all the parts of the system needed to understand what it is, with prose paragraphs and bullet point sand pseudocode and diagrams and api references and file references. It should be very clear and technical. Store in the ticket and the nupload to remarkable."

**Assistant interpretation:** Produce documentation rather than code: create a new ticket design doc for Phase 7, make it comprehensive and implementation-ready, update ticket bookkeeping, and upload the bundle to reMarkable.

**Inferred user intent:** Align the dashboard build before implementation so future frontend work has clear affordances, page contracts, component boundaries, and Storybook expectations.

**Commit (code):** N/A — documentation-only work is uncommitted at this entry.

### What I did

- Added a new design document under the existing HOST-001 ticket:
  - `design-doc/02-phase-7-user-dashboard-affordances-page-designs-and-component-system-guide.md`
- Grounded the design in current backend files and APIs.
- Listed dashboard affordances for global shell, orgs, sites, deployments, runtime, agents, audit, usage, and members.
- Defined `/app` route structure.
- Wrote page descriptions with ASCII screenshots.
- Defined atoms, molecules, organisms, and pages with one directory per widget.
- Required Storybook stories for every widget and page.
- Specified MSW/fake-store Storybook setup.
- Sketched RTK Query API endpoints and TypeScript types.
- Added implementation phases and testing strategy.
- Marked the Phase 7 design-guide checklist item complete.

### Why

The backend and CLI now expose enough functionality to start the dashboard, but the frontend needs a shared blueprint before implementation. The design document turns the user's requirements into page-level and component-level contracts.

### What worked

- The existing ticket structure made it easy to add a second design doc.
- The current backend routes and schema provide enough evidence to define API references and UI behavior.

### What didn't work

- N/A.

### What I learned

- The current `/app` and `/admin` routes are explicitly placeholders, which cleanly separates Phase 7 user dashboard work from future platform admin console work.

### What was tricky to build

- The design needed to include future-looking components such as quota and agent grants without pretending those APIs are complete. I called those out as API-pending or preview states.

### What warrants a second pair of eyes

- Review the route naming and whether dashboard URLs should use org/site IDs only or introduce slug lookup APIs later.
- Review the component taxonomy before implementation to avoid over- or under-splitting widgets.

### What should be done in the future

- Implement Phase 7 scaffold under `web/admin`.
- Verify exact `@go-go-golems/os-core` package exports before coding theme imports.
- Add Storybook CI once scaffold exists.

### Code review instructions

- Start with the new Phase 7 design doc.
- Compare API references against `internal/httpapi/handler.go`, `internal/httpapi/deployments.go`, `internal/httpapi/runtime.go`, and `internal/httpapi/agents_audit.go`.
- Confirm that every requested page/widget/storybook requirement is represented before implementation begins.

## Step 18: Split platform admin dashboard into HOST-003

The admin-side dashboard work that was originally tracked as future `/admin` work in this platform ticket is now split into a dedicated ticket: `HOST-003-ADMIN-DASHBOARD`.

The new ticket should follow the same implementation style as `HOST-002-USER-DASHBOARD`: Storybook-first React pages, MSW fixtures, RTK Query API integration, embedded Vite SPA delivery, macOS1 styling via os-core, and frequent docmgr task/changelog/diary updates.

Initial HOST-003 scope starts with read-only platform admin observability: route guard, admin shell, overview, and runtime summary using the existing `/api/v1/admin/runtimes/summary` endpoint.

## Step 19: Reconciled HOST-001 Phase 7/8 task state after dashboard tickets

After completing the dedicated user dashboard (`HOST-002`) and platform admin dashboard (`HOST-003`) slices, I returned to the original platform ticket and reconciled Phase 7 and Phase 8 tasks.

Phase 7 is now mostly complete from the HOST-002 work: Vite/React/RTK Query/Storybook scaffold, os-core theme imports, user/org/site routes, deployment workflow, agents/audit/usage/members pages, embedded SPA serving, and Storybook build validation are all done. I left the bot token/enrollment and agent grant editor items unchecked because those belong with the still-open Phase 9 agent enrollment/grants work, not the completed user dashboard MVP. I also left the automated Playwright smoke test unchecked because we have manual Playwright verification but not a committed automated E2E test.

Phase 8 is now checked off for the v1 admin dashboard scope delivered through HOST-003: platform-admin guard, `/api/v1/admin/*`, global inventory endpoints, runtime restart/stop, admin pages, Storybook coverage, non-admin 403 tests, seeded dev admin verification, and read-only policy pages. I annotated quota/admin-agent items where the current delivered scope is read-only/listing and future write workflows remain separate follow-up work.

## Step 20: Agent enrollment schema, signed deploy-run backend, and agent CLI

Phase 9 needed to turn the placeholder agent CLI and existing agent inventory table into an end-to-end machine deploy path. I implemented a compact v1 flow: a human creates an agent and receives a one-time enrollment token, the machine generates an Ed25519 key, enrollment registers that public key, and future deploy-run creation requests are signed with timestamp and nonce replay protection.

The implementation intentionally separates signed deploy-run creation from bundle upload. Signed requests authorize intent and produce a short-lived upload token; multipart bundle upload is then bound to that deploy run and token. This keeps large multipart bodies out of the signature path while still preventing arbitrary uploads.

### Prompt Context

**User prompt (verbatim):** "Do phase 9, in one go, keep a detailed diary, commit at appropriat eintervals."

**Assistant interpretation:** Complete HOST-001 Phase 9 agent enrollment and signed deploy-run work in a single sustained implementation pass, writing diary entries and committing focused slices.

**Inferred user intent:** Finish the remaining headless deployment-agent MVP so CI-style agents can deploy without human browser credentials, with enough documentation and commit history to review safely.

**Commit (code):** 326581d5da38cd793a734642d4a8a28b02581b46 — "Add signed agent enrollment and deploy runs"

### What I did

- Added migration `004_agent_enrollment_runs.sql` for hashed one-time enrollment tokens and supporting indexes.
- Extended `internal/store/queries/agents.sql` and regenerated sqlc output for:
  - enrollment tokens,
  - agent keys,
  - nonce insertion,
  - deploy-run creation/fetch/finish,
  - last-seen updates.
- Added store models/wrappers for `AgentEnrollmentToken`, `AgentKey`, and `DeployRun`.
- Added `internal/control/agent_runs.go` for:
  - enrollment token generation and hashing,
  - immediate optional site grants during agent creation,
  - enrollment token exchange for registered public keys,
  - canonical Ed25519 signature verification,
  - five-minute timestamp skew checks,
  - nonce replay prevention,
  - grant checks for site/channel/path,
  - deploy-run creation and upload-token validation.
- Extended deployment upload to support `actorType=agent`/`actorID=<agent>` after the signed deploy-run/token gate.
- Added HTTP endpoints:
  - `POST /api/v1/agent/enroll`,
  - `POST /api/v1/agent/deploy-runs`,
  - `POST /api/v1/agent/deploy-runs/{run_id}/upload`,
  - `POST /api/v1/orgs/{org_id}/agents/{agent_id}/grants`.
- Updated human `go-go-host agents create` output to include the one-time enrollment token and optional immediate grant flags.
- Implemented `go-go-host-agent` Glazed commands:
  - `keygen`,
  - `enroll`,
  - `deploy`,
  - enhanced `status`.

### Why

- Agents need machine credentials that are independent from human OIDC/dev-auth sessions.
- Enrollment tokens let humans create/approve a machine identity without copying private keys into the control plane.
- Ed25519 signatures, timestamps, and nonces give the platform a minimal replay-resistant request-auth layer for headless agents.
- Separating deploy-run authorization from upload keeps authorization JSON small and deterministic while still binding the later bundle upload to a short-lived server-side run.

### What worked

- `sqlc generate` produced clean typed accessors after adding the new queries.
- `go test ./...` passed after the first compile fix.
- The CLI root already had Glazed logging/help setup, so the new agent commands could reuse the existing command builder.
- `go run ./cmd/go-go-host-agent keygen --config /tmp/go-go-host-agent-test.json --output json` produced a local config and JSON row as expected.

### What didn't work

- First test run failed because sqlc named the deploy-run create parameter `ActorID`, not `AgentID`:

```text
internal/store/agents.go:141:108: unknown field AgentID in struct literal of type db.CreateDeployRunParams
```

I fixed `Store.CreateDeployRun` to populate `ActorID: input.AgentID` because the SQL inserts the same value into both `actor_id` and `agent_id`.

- The first `go-go-host-agent enroll` response struct used a combined anonymous field tag that would not decode `orgId` into `OrgID`. I fixed it in a follow-up commit with explicit JSON tags for `id`, `orgId`, `status`, and `name`.

### What I learned

- The existing schema already had most Phase 9 tables (`agent_keys`, `agent_site_grants`, `agent_nonces`, `deploy_runs`), so the missing durable piece was enrollment-token state.
- `net/http` route patterns can expose unauthenticated agent endpoints alongside the existing user-authenticated API by mounting `/api/v1/agent/` directly to the API mux rather than through dev/OIDC middleware.
- Signing JSON deploy-run creation is much simpler and safer than trying to include multipart upload bodies in the agent signature flow.

### What was tricky to build

- The canonical string must match exactly between CLI and server. I centralized the canonical format in `control.AgentCanonicalString()` and reused it from the CLI to avoid drift.
- Request body reading had ordering hazards: the server must hash the exact JSON body used for signature verification and then decode it. I fixed this by reading the body once in the deploy-run handler, decoding from a `bytes.Reader`, and passing the original bytes to verification.
- Grant path matching needed to be useful without becoming a full policy language. The v1 matcher supports empty/all, exact paths, glob-style `filepath.Match`, and simple `prefix/**` semantics.

### What warrants a second pair of eyes

- The signature canonicalization contract is security-critical; review `control.AgentCanonicalString`, CLI `signRequest`, and HTTP `verifyAgentRequest` together.
- Enrollment tokens are returned once in API/CLI output. Review logs and UI surfaces before production use to ensure they are not accidentally persisted beyond operator intent.
- Upload tokens are bearer tokens bound to deploy runs. Review expiry/status transitions and confirm all upload failure paths mark runs as rejected where desired.
- There is no dedicated per-key revoke endpoint yet; v1 revocation is handled at the agent status level. If per-key rotation becomes a requirement, add explicit key revoke APIs/tests.

### What should be done in the future

- Add a first-class grant editor in the dashboard instead of only optional immediate grant fields and the grant API.
- Add production-grade token redaction guidance for CLI output/logging.
- Add key rotation and per-key revoke workflows.
- Add automated CLI E2E against a live devctl stack, not only httptest-backed API coverage and command smoke.

### Code review instructions

- Start with `internal/control/agent_runs.go`; review token generation, signature verification, nonce insertion, grant checks, and deploy-run creation.
- Then review `internal/httpapi/agents_audit.go` and `internal/httpapi/deployments.go` for transport boundaries and upload-token handling.
- Review `cmd/go-go-host-agent/cmds/support.go`, `keygen.go`, `enroll.go`, and `deploy.go` for the CLI-side canonical signing flow and local config persistence.
- Validate with:

```bash
sqlc generate
go test ./...
go run ./cmd/go-go-host-agent --help
go run ./cmd/go-go-host-agent keygen --config /tmp/go-go-host-agent-test.json --output json
```

### Technical details

Canonical signature string:

```text
METHOD\nPATH?QUERY\nSHA256_HEX_BODY\nRFC3339_TIMESTAMP\nNONCE
```

Signed request headers:

```text
X-Go-Go-Agent-ID
X-Go-Go-Agent-Key-ID
X-Go-Go-Agent-Timestamp
X-Go-Go-Agent-Nonce
X-Go-Go-Agent-Signature
```

The upload endpoint does not require a signature; it requires the short-lived `X-Go-Go-Upload-Token` generated by the signed deploy-run endpoint.

## Step 21: Agent CLI help pages and enroll decoding fix

After the main Phase 9 commit, I added operator-facing help pages for the new agent commands and fixed a CLI response decoding bug before moving on. This was a documentation-heavy slice, but it also caught a real usability issue: the enrolled org ID would not have been persisted correctly because the JSON tag did not match `orgId`.

### Prompt Context

**User prompt (verbatim):** (same as Step 20)

**Assistant interpretation:** Continue completing Phase 9 with documentation and focused commits after the main backend/CLI implementation.

**Inferred user intent:** Make the new agent workflow discoverable from the CLI and avoid leaving known decode bugs in committed code.

**Commit (docs):** eafc6f5794a39dd943f2d26c0fa213ca72ae598b — "Document agent enrollment workflow"

**Commit (code):** 8411ef74c57c5ac7d6a1b5aa581c6f83c6ddcaf9 — "Fix agent enroll response decoding"

### What I did

- Rewrote `cmd/go-go-host-agent/doc/getting-started.md` to describe keygen/enroll/deploy instead of the old placeholder status-only flow.
- Added `keygen-enroll-deploy.md` with the canonical operator workflow.
- Added `signature-troubleshooting.md` with common failure modes: clock skew, replay, revoked credentials, and grant mismatches.
- Fixed `cmd/go-go-host-agent/cmds/enroll.go` so the response struct decodes `orgId` into local agent config.

### Why

- Phase 9 adds a multi-step workflow; without embedded help pages, users would need to infer too much from flags.
- The CLI config needs `orgId`/`siteId` persisted correctly for `status` output and future grant-aware workflows.

### What worked

- The docs are picked up by the existing `//go:embed *.md` loader in `cmd/go-go-host-agent/doc/doc.go`.
- `go test ./cmd/go-go-host-agent/...` passed after adding the docs.
- Full `go test ./...` passed after the decode fix.

### What didn't work

- No command failure in this step; the issue was found by inspection before a user hit it.

### What I learned

- Keeping focused commits helped: the code-heavy signed-agent commit, docs commit, decode-fix commit, and test-coverage commit each tell a distinct review story.

### What was tricky to build

- The help pages need to describe enough of the security model without over-promising production key rotation or dashboard grant editing. I phrased revoke/key-rotation items as troubleshooting/future operational concerns.

### What warrants a second pair of eyes

- Review whether help text should warn more strongly that enrollment tokens are one-time secrets and should not be pasted into shared logs.

### What should be done in the future

- Add a longer production runbook after production platform-admin bootstrap and key rotation are implemented.

### Code review instructions

- Review `cmd/go-go-host-agent/doc/*.md` with the new commands open.
- Run `go run ./cmd/go-go-host-agent help agent-keygen-enroll-deploy` once help output formatting is being reviewed interactively.

### Technical details

The agent docs now cover:

- key generation,
- one-time token enrollment,
- signed deploy-run creation,
- upload-token-bound bundle upload,
- common signature failure reasons.

## Step 22: Signed agent deploy-run security coverage and Phase 9 task reconciliation

I added integration coverage for the security boundaries that make the signed-agent flow meaningful: successful deploy-run creation and upload for an allowed grant, nonce replay denial, timestamp skew denial, bad signature denial, and denial for ungranted site/path access. Then I reconciled the HOST-001 Phase 9 checklist with the delivered v1 scope.

### Prompt Context

**User prompt (verbatim):** (same as Step 20)

**Assistant interpretation:** Finish Phase 9 by proving the happy path and key denial cases, then update ticket bookkeeping.

**Inferred user intent:** Avoid a superficially complete agent workflow that lacks security regression tests or ticket traceability.

**Commit (code):** 530a5d97f8a50147b7623675e6acdb5e447af757 — "Cover agent deploy run upload security"

### What I did

- Added `internal/httpapi/agent_signed_integration_test.go`.
- Covered:
  - human creates agent with immediate site/channel/path grant,
  - agent enrolls with Ed25519 public key,
  - signed deploy-run creation succeeds,
  - upload-token-bound bundle upload succeeds,
  - replayed nonce is denied,
  - old timestamp is denied,
  - bad signature is denied,
  - wrong path is denied,
  - wrong site is denied.
- Updated `HOST-001` Phase 9 tasks and exit criteria to checked with notes where v1 scope differs from future richer workflows.

### Why

- Signed request systems are easy to get wrong silently; the regression test exercises the verifier through HTTP instead of only unit-level helpers.
- The ticket checklist needed to reflect the implemented behavior while preserving caveats around per-key revoke and richer grant editing.

### What worked

- `go test ./...` passes.
- `go run ./cmd/go-go-host-agent --help` works.
- `go run ./cmd/go-go-host-agent keygen --config /tmp/go-go-host-agent-test.json --output json` emits a Glazed JSON row.

### What didn't work

- Integration tests that require Postgres still skip unless `GO_GO_HOST_TEST_DATABASE_URL` is set, following existing project behavior. The new tests are wired into the same integration harness.

### What I learned

- The existing deployment test helper `writeHelloBundle` made it straightforward to prove the upload half of the deploy-run flow without adding more fixture files.
- The current status-code mapping treats replay and old timestamp as bad request while permission failures become forbidden. This is acceptable for v1, but API clients should treat both as non-retryable without fresh signing state.

### What was tricky to build

- The first version of the signed-flow test only created a deploy run. That proved signing/grants but not the upload-token endpoint. I added `uploadAgentBundleViaAPI` so the test now satisfies the Phase 9 exit criterion that an agent can deploy to an allowed site.

### What warrants a second pair of eyes

- Confirm that the integration test assertions cover the most important security boundaries. Future timestamp and explicit key-revoke tests should be added when per-key revoke is implemented.
- Confirm whether replay denial should be reported as `403` instead of `400`; current behavior exposes the duplicate nonce as a validation error.

### What should be done in the future

- Add a live devctl CLI E2E that runs `go-go-host agents create`, `go-go-host-agent keygen`, `enroll`, and `deploy` against a real daemon.
- Add dashboard UI for creating grants and displaying the one-time enrollment token more safely.
- Add key rotation and explicit key revoke.

### Code review instructions

- Review `internal/httpapi/agent_signed_integration_test.go` as the executable description of the Phase 9 security contract.
- Run:

```bash
go test ./...
```

- If a Postgres test DB is available, set `GO_GO_HOST_TEST_DATABASE_URL` and rerun to execute the integration path instead of skipping it.

### Technical details

The test uses Ed25519 directly and signs the same canonical string used by the CLI. This avoids treating the CLI as a test dependency while still exercising the HTTP verifier and deploy-run/upload endpoints.

## Step 23: Live devctl agent deployment smoke test

I ran the Phase 9 workflow against the real devctl stack rather than only httptest and command-level smoke. This live pass was valuable because it exposed two path-policy bugs that unit tests had not caught: deploy runs were storing the requested logical path as the upload validation policy, and the deployment archive validator did not treat `**` as an allow-all nested-path pattern.

After fixing both issues and restarting `go-go-hostd`, the live workflow succeeded: a dev-auth user created an org/site/agent, the agent generated a key, enrolled with a one-time token, created a signed deploy run, uploaded a bundle, the deployment validated as `created_by_type=agent`, a human activated it, and the public Host-header route served the uploaded site.

### Prompt Context

**User prompt (verbatim):** "run a live test against devctl"

**Assistant interpretation:** Start/restart the local devctl stack and execute a real CLI/API agent enrollment and deployment workflow against `http://127.0.0.1:8080`.

**Inferred user intent:** Confirm that Phase 9 works outside tests, catch integration gaps, and report exactly what was smoke-tested.

**Commit (code):** 9122ef2044048e7d757a7a0dbbd7fb07a6b5195d — "Preserve grant paths for agent uploads"

**Commit (code):** ce1192f7c22415d7918318f5dba22b86f3f618f3 — "Allow double-star deployment path grants"

**Commit (test):** 2210e0dc7841be08151e025cadebc963c4719012 — "Test double-star deployment path grants"

### What I did

- Ran `devctl up --force` and confirmed Postgres, `go-go-hostd`, Vite, and Storybook were alive.
- Ran a full live CLI flow using `go run`:
  - `go-go-host org create`,
  - `go-go-host site create`,
  - `go-go-host agents create --site-id ... --path '**'`,
  - `go-go-host-agent keygen`,
  - `go-go-host-agent enroll`,
  - `go-go-host-agent status`,
  - `go-go-host-agent deploy`,
  - `go-go-host deployments activate`,
  - `curl -H 'Host: live-agent-site-...localhost' http://127.0.0.1:8080/`,
  - `go-go-host audit list`.
- Fixed deploy-run creation to preserve grant archive paths for upload validation instead of narrowing validation to the requested logical deploy path.
- Fixed deployment bundle path matching so `**` allows nested archive paths.
- Added a regression test for `AllowedPaths: []string{"**"}` against nested archive entries.

### Why

- The previous test suite proved the signed HTTP path, but not a live daemon compiled/restarted under devctl with real CLI binaries and the existing `/tmp/go-go-host-test-bundle.tar.gz`.
- Agent grants use path policies twice: once for signed deploy-run authorization and once for archive-entry validation. The live smoke made it clear those two concepts must not be conflated.

### What worked

Final successful live output included:

```text
"status": "validated",
"valid": true,
"created_by_type": "agent"
```

Activation returned:

```text
"status": "active",
"created_by_type": "agent",
"created_by_id": "agt_3908903e-6307-45ff-b43b-735a04103da5"
```

The public route returned the uploaded HTML:

```html
<title>Hello from dashboard upload</title>
<h1>Hello from go-go-host</h1>
```

Audit showed the expected sequence:

```text
org.create
site.create
agent.create
agent.grant.upsert
agent.enrollment_token.create
agent.enroll
deploy_run.create   actor_type=agent
deployment.upload   actor_type=agent
deployment.activate actor_type=user
```

### What didn't work

First live deploy failed with rejected validation because the deploy run stored the requested logical path (`/tmp/go-go-host-test-bundle.tar.gz`) as the archive validation allow-list:

```text
path "assets/style.css" is not allowed by deployment policy
path "scripts/app.js" is not allowed by deployment policy
path "go-go-host.json" is not allowed by deployment policy
```

I fixed that by storing `grant.AllowedPaths` in the deploy run for upload validation.

Second live deploy still failed for nested paths because the deploy package treated `**` like a `filepath.Match` pattern, which does not match slashes:

```text
path "assets/style.css" is not allowed by deployment policy
path "scripts/app.js" is not allowed by deployment policy
```

I fixed `internal/deploy/bundle.go` so `pattern == "**"` explicitly allows all archive paths and added `TestValidateAndStoreAllowsDoubleStarPolicy`.

A separate CLI usage mistake happened during activation: I initially passed the deployment ID positionally, but the Glazed command expects `--deployment-id`. The corrected command worked:

```bash
go run ./cmd/go-go-host deployments activate \
  --api-url http://127.0.0.1:8080 \
  --dev-user live-agent-1778594396 \
  --deployment-id dep_4710ba96-8150-4df8-aebb-4eb7ea807ab1 \
  --output json
```

### What I learned

- Live devctl smoke was necessary; the httptest integration used a narrower grant path (`bundles/**`) and did not expose the deploy-package `**` mismatch.
- The deployment policy should use grant archive path allow-lists, while the signed deploy-run request path should be treated as a requested logical path for authorization only.
- The current agent workflow uploads and validates but does not auto-activate; activation is still a human/user operation in this smoke.

### What was tricky to build

- There were two distinct path languages in play: agent deploy-run request paths and archive-entry validation paths. The symptom looked like a permission failure, but the root cause was that the wrong path value was being carried from signed deploy-run creation into bundle validation.
- `filepath.Match("**", "assets/style.css")` does not mean recursive allow-all. I made `**` an explicit sentinel in the deploy validator to match the grant helper behavior.

### What warrants a second pair of eyes

- Review whether the deploy-run response should include both requested logical path and effective archive path policy to make this distinction visible.
- Review whether `**` should be accepted as a production policy shorthand or replaced with a clearer explicit value such as `all`.
- Review the current human activation step: if CI agents should auto-activate in some channels, that should be a deliberate grant flag and audited separately.

### What should be done in the future

- Add a scripted live E2E under `ttmp/.../scripts` or a Makefile target that repeats this devctl smoke.
- Add CLI help examples showing `--deployment-id` for activation and `--path bundles/app.tar.gz` for agent deploy-run requests.
- Add a focused test for a signed deploy run with grant `**` and nested archive entries if the integration harness is guaranteed to run under CI with Postgres.

### Code review instructions

- Review `internal/control/agent_runs.go` around `CreateDeployRun`: the effective upload validation policy now comes from `grant.AllowedPaths`.
- Review `internal/deploy/bundle.go` and `internal/deploy/bundle_test.go` for `**` path policy semantics.
- Reproduce with:

```bash
devctl up --force
# create org/site/agent, keygen/enroll/deploy as shown in this diary step
curl -H 'Host: <site>.localhost' http://127.0.0.1:8080/
```

### Technical details

Successful live IDs from the final run:

```text
user: live-agent-1778594396
org: org_d03dc993-65fb-4a2b-bd43-7d5593741828
site: site_dc002068-d404-4a85-aa4d-fae323959865
host: live-agent-site-1778594396.localhost
agent: agt_3908903e-6307-45ff-b43b-735a04103da5
key: ak_aaa7d034-a01e-4f7b-a75b-9e2a4c87e05a
deploy run: dr_75846d83-0ffc-46a8-8f4f-c81d34f4c133
deployment: dep_4710ba96-8150-4df8-aebb-4eb7ea807ab1
```

## Step 24: Scoped agent auto-activation

The user asked whether human activation was done in the UI, then asked to add scoped agent auto-activation. I added a new Phase 9A rather than silently expanding Phase 9, because this changes the trust model: some agents can now promote validated deployments to live traffic, but only when a human-created site grant explicitly includes `can_activate`.

The resulting flow keeps the safe default from Phase 9: normal agents still upload and validate only. Auto-activation requires both grant permission (`canActivate`) and an explicit agent request (`go-go-host-agent deploy --activate`). When both are present, the upload endpoint validates the bundle, activates the deployment as the agent, returns `activated: true`, and audit records `deployment.activate` with `actor_type=agent`.

### Prompt Context

**User prompt (verbatim):** "add scoped agent auto-activation. Add new phase and new tasks and then build it"

**Assistant interpretation:** Add a new HOST-001 phase for scoped agent activation, update tasks first, then implement backend, CLI, tests, docs, and live validation.

**Inferred user intent:** Allow CI/trusted agents to deploy directly to live traffic in controlled scopes instead of requiring a human promotion step every time.

**Commit (code):** 248494d649cc32501f219bfd075b3b4847740224 — "Add scoped agent auto activation"

**Commit (docs/tasks):** fb4a7cbd8baa8aec8f648773073d30522b43d6b8 — "Document scoped agent auto activation"

### What I did

- Added migration `005_agent_auto_activate.sql` with `agent_site_grants.can_activate`.
- Updated sqlc queries/generated code and store models for `can_activate`.
- Extended agent creation/grant APIs with `canActivate`.
- Added human CLI flag:
  - `go-go-host agents create --can-activate`.
- Added agent CLI flag:
  - `go-go-host-agent deploy --activate`.
- Updated signed deploy-run creation so `activate` is appended to `allowed_actions` only when the grant has `CanActivate`.
- Added `DeploymentService.ActivateAsAgent()` and shared activation internals so user and agent activations use the same runtime swap/DB update path but produce different audit actors.
- Updated agent upload handling so a valid upload with deploy-run action `activate` auto-activates before returning.
- Updated the signed-agent integration test so the happy path requests auto-activation and asserts the upload response is `active` / `activated: true`.
- Added Phase 9A tasks and exit criteria to HOST-001.
- Updated agent help docs to mention `--activate` and the required grant.

### Why

- Some CI agents should be able to continuously deploy to a known site/channel/path without a human click after every validation.
- This should not be a global agent capability. It belongs on the same site grant that already scopes deploy rights.
- Requiring `--activate` keeps auto-promotion explicit per deploy-run even for agents that are allowed to do it.

### What worked

- `go test ./...` passed after the schema/model/API/CLI changes.
- Live devctl smoke succeeded with:

```text
"activated": true
"status": "active"
```

- The public Host-header route served the auto-activated deployment immediately after agent deploy.
- Audit showed:

```text
deployment.activate actor_type=agent
```

### What didn't work

- No compile failures after the main implementation beyond the expected sqlc row-shape adjustment in store wrappers: adding `can_activate` changed sqlc return types from the table model to generated row structs for explicit SELECT/RETURNING lists. I fixed this by adding separate `agentGrantFromListRow` and `agentGrantFromUpsertRow` helpers.

### What I learned

- The Phase 9 deploy-run model already anticipated multiple actions via `allowed_actions`, so auto-activation fit naturally as an additional action rather than a separate deploy-run table.
- Keeping activation in `DeploymentService` avoided duplicating runtime supervisor swap logic in HTTP handlers.

### What was tricky to build

- The tricky part was preserving authorization layering. The signed deploy-run endpoint decides whether the run may include `activate`; the upload endpoint only executes activation if the already-created deploy run contains that action. This avoids trusting an upload-time request parameter.
- User activation and agent activation share most mechanics but must differ in authorization and audit attribution. I split the public methods (`Activate`, `ActivateAsAgent`) and moved common work to a private `activate` helper.

### What warrants a second pair of eyes

- Review `ActivateAsAgent` to ensure grant expiry/status/org checks are sufficient for production.
- Review whether `can_activate` should be constrained by channel in a stronger way than the deploy-run's existing channel grant check.
- Review whether auto-activation should be limited to deployments created by the same agent. The current upload path naturally activates the deployment it just created; the service method is more general and relies on agent/site grant.

### What should be done in the future

- Add dashboard UI for `canActivate` grants with strong warnings.
- Add an admin/user audit filter preset for agent activations.
- Add key rotation/revoke before treating auto-activation as production-ready for high-risk sites.

### Code review instructions

- Review `internal/control/agent_runs.go` for deploy-run `activate` authorization.
- Review `internal/control/deployments.go` for `ActivateAsAgent` and shared `activate` behavior.
- Review `internal/httpapi/deployments.go` for upload-time auto-activation behavior.
- Review `cmd/go-go-host-agent/cmds/deploy.go` and `cmd/go-go-host/cmds/agents.go` for the new CLI flags.
- Validate with:

```bash
go test ./...
devctl restart go-go-hostd
# create org/site/agent with --can-activate, enroll, deploy with --activate
```

### Technical details

Successful live scoped auto-activation smoke:

```text
org: org_7bca711c-991a-4bc1-959a-9a3ac728e7f0
site: site_61c023b3-ad2a-4e5b-888d-36752f2df7c6
agent: agt_d3855b5a-44cc-4279-bef3-50af4b77446e
key: ak_24993377-21e9-4a71-bc57-42f4d77aeac7
deploy run: dr_5762f6f6-e1b5-4ece-b894-f27834816597
deployment: dep_0790020c-9c3d-4cc2-a918-48ad9e3b2eff
```

Agent deploy output:

```json
{
  "activated": true,
  "status": "active",
  "valid": true
}
```

## Step 25: Phase 10A immediate agent hardening slice

I added a new Phase 10A to capture the hardening items that make the just-built agent auto-activation workflow easier and safer to operate immediately. Then I implemented the first slice: `canActivate` is now visible in the user dashboard, backend grant creation enforces owner-only auto-activation grants, agent keys can be listed with fingerprints and last-used timestamps, and operators can revoke individual signing keys without revoking the whole agent.

This is deliberately not the whole hardening guide. Security-failure audit events, upload token one-time transitions, and bundle SHA256 traceability remain open in Phase 10A. The committed slice focuses on operator control over the highest-risk thing we just added: agents that can activate live traffic.

### Prompt Context

**User prompt (verbatim):** "add these as phase 10A, with the same detailed tickets. Then get to work on phase 10."

**Assistant interpretation:** Add a detailed HOST-001 Phase 10A for the immediate hardening recommendations from HOST-004, then begin implementation rather than leaving it as planning only.

**Inferred user intent:** Turn the hardening guide into actionable project tasks and immediately make the agent/auto-activation system safer to use.

**Commit (code):** f31621773be37f26612d923782a85a61d2f1ecb9 — "Add agent key inventory and revoke controls"

### What I did

- Added migration `006_agent_key_inventory.sql` with `agent_keys.last_used_at`.
- Added key list/revoke sqlc queries and store wrappers.
- Added backend APIs:
  - `GET /api/v1/orgs/{org_id}/agents/{agent_id}/keys`,
  - `POST /api/v1/orgs/{org_id}/agents/{agent_id}/keys/{key_id}/revoke`.
- Added key fingerprints derived from SHA256 of the public key.
- Touched key `last_used_at` after successful signed request verification.
- Added owner-only enforcement when a human creates/updates a grant with `canActivate`.
- Added dashboard pieces:
  - `AgentKeysTable`,
  - key list/revoke actions on `AgentsPage`,
  - `canActivate` create-agent warning checkbox,
  - one-time enrollment token display from create response.
- Extended signed-agent integration coverage to list keys, assert last-used/fingerprint, revoke key, and prove revoked keys are denied.

### What worked

- `go test ./...` passes.
- `pnpm --dir web/admin build` passes.
- Storybook has an `AgentKeysTable` story.

### What didn't work

- The first web build failed because the new story imported `Meta`/`StoryObj` from `@storybook/react-vite`; existing project stories import from `@storybook/react`. I changed the import to match the repo convention.
- Adding `last_used_at` changed sqlc row shapes for key queries. I fixed the generated-query alignment by selecting/returning `last_used_at` consistently.

### What warrants a second pair of eyes

- The dashboard create-agent form currently exposes `canActivate` without selecting an immediate site grant. The backend only creates an activation-capable grant when a `siteId` is provided, so a fuller grant editor/detail page should follow.
- Key revoke uses developer/owner deploy role today. Decide whether key revoke should be owner-only when the agent has activation grants.
- The fingerprint format is short (`SHA256:` plus first 8 bytes of digest). Review whether to display a longer fingerprint.

### What should be done in the future

- Finish Phase 10A open items: security-failure audit events, upload token one-time semantics, bundle SHA256.
- Add a proper AgentDetailPage with grant editor instead of expanding the list page further.
- Add key rotation/replacement-key success flow.

### Code review instructions

- Start with `internal/control/agents.go`, `internal/control/agent_runs.go`, and `internal/httpapi/agents_audit.go`.
- Review `web/admin/src/pages/AgentsPage/AgentsPage.tsx` and `AgentKeysTable` for operator UX.
- Validate with:

```bash
go test ./...
pnpm --dir web/admin build
```

## Step 26: Completed Phase 10 and Phase 10A hardening pass

I completed the Phase 10/10A pass as an implementation slice rather than only a planning update. The work hardened the recently-added agent auto-activation path, made deployment artifacts more traceable, connected site capability policy to validation, configured runtime DB quota enforcement from site quotas, added a DB stats API, recorded runtime events, and rebuilt the embedded dashboard.

This pass also exposed one operational detail during live smoke: the live auto-activation flow succeeded, DB stats were reachable, and the public host served the deployment, but the CLI did not initially expose `bundleSha256` in `deployments show`. I added the CLI row field after seeing that the API field was present but the CLI output dropped it.

### Prompt Context

**User prompt (verbatim):** "do phase 10 and phase 10a in their entirety, commit at appropriate intervals, keep a detailed diary."

**Assistant interpretation:** Finish both the original Phase 10 hardening/observability checklist and the newly added Phase 10A immediate agent-hardening checklist, committing in reviewable chunks and recording the implementation journey.

**Inferred user intent:** Move the platform from “feature works” to “feature is safer, observable, and easier to operate,” especially around agent auto-activation.

**Commits (code/docs):**

- f31621773be37f26612d923782a85a61d2f1ecb9 — "Add agent key inventory and revoke controls"
- a23680bcc2e97f69e889f88720aeb227641b1ac1 — "Add Phase 10 deployment hardening foundations"
- 0fa578223d1b612e05beb43e19d9f8f7e27aeda1 — "Surface bundle hashes and quota warnings"
- 2b2e0e7ab67606688f24a2257be6daadaf6107d8 — "Apply site capability policy during validation"
- 2df31cb1b96d26041faf5d2644dd0f2d6f508b6d — "Add agent key rotation tokens"
- b917a61dabade64fc35c3445cbfe740852bb4173 — "Rebuild embedded dashboard for hardening UI"
- 3f7287ef7c0501fbfd3964468254555db6474651 — "Complete Phase 10 hardening checklist"

### What I did

Phase 10A agent hardening:

- Added key inventory and revoke APIs.
- Added key fingerprints and `last_used_at` tracking.
- Added key rotation token endpoint for existing agents.
- Added dashboard key table, key revoke action, and rotation-token action.
- Added owner-only enforcement for `canActivate` grant creation/update.
- Extended signed-agent integration coverage for revoked-key denial, replacement-key success, and second-upload denial.

Phase 10 capability/quota/deployment hardening:

- Expanded safe capability policy to include `express`, `ui.dsl`, `database`, `db`, timers, static assets, and sqlite compatibility.
- Added default site capability rows on site creation.
- Wired site capability policy into deployment validation so requested/effective capability reports use stored policy.
- Added request timeout enforcement through the runtime supervisor using site quota `request_timeout_ms`.
- Configured DB guard soft/hard byte limits from site quota during dry-run and runtime activation.
- Added DB stats endpoint: `GET /api/v1/sites/{site_id}/db/stats`.
- Added runtime events table and status-recorder insertion for runtime status changes.
- Added `bundle_sha256` to deployments, API DTOs, CLI rows, and dashboard deployment details.
- Added dashboard quota/runtime error warning text.
- Added a runtime hard-limit test and kept existing forbidden-capability tests.
- Rebuilt embedded dashboard assets with `go run ./cmd/build-web`.

### Why

- Agent auto-activation makes key compromise and grant mistakes more serious, so operators need key visibility, key revoke, and key rotation.
- Deployment validation should use stored site policy, not just hard-coded defaults, otherwise policy pages are observability-only and not enforcement.
- Bundle hashes are the simplest artifact integrity handle: they make it possible to identify what was uploaded and running.
- Runtime quota/DB stats visibility closes the loop between site quota rows, runtime behavior, and dashboard/operator diagnostics.

### What worked

- `go test ./...` passed after each focused slice.
- `pnpm --dir web/admin build` passed after dashboard changes.
- `make storybook-build` passed with the new AgentKeys story.
- `go run ./cmd/build-web` rebuilt embedded assets successfully through Dagger.
- A live devctl smoke after the hardening foundations verified:
  - agent auto-activation still works,
  - public Host-header routing serves the uploaded site,
  - `/api/v1/sites/{site_id}/db/stats` returns DB byte/quota stats.

### What didn't work

- The first story build after adding `AgentKeysTable` failed because I imported Storybook types from `@storybook/react-vite` instead of the repo convention `@storybook/react`. I fixed the import.
- Adding `last_used_at` changed sqlc row shapes for agent-key queries; I fixed the query return lists so create/get/list all include the new column.
- After the live Phase 10 smoke, `go-go-host deployments show` did not display `bundle_sha256` because the CLI DTO/row mapping did not include the field. I added `BundleSHA256` to the CLI deployment DTO and row output.

### What I learned

- The current schema already had useful policy tables (`site_capabilities`, `site_quotas`), but the runtime/deploy path needed explicit wiring to make them enforcement mechanisms.
- `deploy_runs.allowed_actions` was a good extension point for auto-activation and later hardening; it now also works naturally with one-time upload semantics.
- Bundle hash visibility needs to be end-to-end: storing it in DB is not enough if API/CLI/dashboard omit it.

### What was tricky to build

- The agent upload state transition needed to happen after token hash validation but before accepting a bundle. I kept validation in `ValidateUploadToken`, then atomically moved the run from `pending` to `uploading` with a store method before returning control to the upload handler.
- Capability policy had two layers: default safe capabilities in code and stored per-site capability rows. New sites now get default rows, while existing sites with no rows still fall back to code defaults to avoid breaking older dev data.
- Runtime request timeout uses `http.TimeoutHandler`, which can limit request handling time at the HTTP layer. It is not a perfect Goja interruption mechanism, but it is an immediately useful guard until deeper VM interruption controls are designed.

### What warrants a second pair of eyes

- Review the security-failure audit events. They intentionally avoid logging raw signatures/tokens, but the exact action names and metadata should be checked before production.
- Review `ActivateAsAgent` and `canActivate` owner-only enforcement together to confirm the trust model is consistent.
- Review DB hard-limit behavior. The test proves an extremely low hard limit fails fixture writes; production thresholds and error UX should still be tuned.
- Review runtime event volume. The status recorder currently records status changes and request-counter updates can also persist status frequently; event retention/pruning should be part of Phase 12.

### What should be done in the future

- Add a dedicated AgentDetailPage instead of continuing to grow `AgentsPage`.
- Add a richer security-events view instead of relying only on generic audit filtering.
- Add deploy-run expiry cleanup and nonce retention cleanup under Phase 12 production hardening.
- Replace `http.TimeoutHandler` with deeper Goja interrupt support if hosted scripts can still consume CPU after timeout responses.

### Code review instructions

Review in this order:

1. `internal/control/agent_runs.go` — signed request audit, deploy-run state transitions, upload token checks.
2. `internal/control/agents.go` — key list/revoke/rotation token behavior.
3. `internal/control/deployments.go` — capability policy lookup, DB quota runtime spec, agent activation.
4. `internal/runtime/runtime.go` and `internal/runtime/supervisor.go` — DB guard and request timeout wiring.
5. `internal/store/queries/*.sql` and migrations `006`/`007` — schema/query changes.
6. `web/admin/src/pages/AgentsPage/AgentsPage.tsx` and `AgentKeysTable` — operator UI.
7. `cmd/go-go-host/cmds/deployments.go` — bundle hash CLI output.

Validate with:

```bash
go test ./...
pnpm --dir web/admin build
make storybook-build
go run ./cmd/build-web
docmgr doctor --ticket HOST-001-GO-GO-HOST-V1 --stale-after 30
```

### Technical details

New/changed APIs:

```http
GET  /api/v1/orgs/{org_id}/agents/{agent_id}/keys
POST /api/v1/orgs/{org_id}/agents/{agent_id}/keys/{key_id}/revoke
POST /api/v1/orgs/{org_id}/agents/{agent_id}/enrollment-token
GET  /api/v1/sites/{site_id}/db/stats
```

New schema changes:

```sql
ALTER TABLE agent_keys ADD COLUMN IF NOT EXISTS last_used_at TIMESTAMPTZ;
ALTER TABLE deployments ADD COLUMN IF NOT EXISTS bundle_sha256 TEXT NOT NULL DEFAULT '';
CREATE TABLE IF NOT EXISTS runtime_events (...);
```

New deploy-run upload state:

```text
pending -> uploading -> completed
pending -> uploading -> rejected
```

## Step 27: Completed Phase 11 site settings and domain management

I implemented Phase 11 as a full site-settings slice. Sites now have a control-plane API for non-secret config, custom domains with verification tokens, capability policy edits, and an explicit environment/secrets placeholder that documents what is intentionally not exposed to hosted code.

The dashboard now has a site-level Settings page that brings these pieces together in one operator workflow. This keeps code deployment separate from site configuration while preserving the earlier safety boundary: no process environment passthrough, no plaintext secret API, and no unrestricted filesystem access.

### Prompt Context

**User prompt (verbatim):** "Continue with phase 11. Make sure to read the diary and docmgr skills. Keep a diary as you work, commit at appropriate intervals"

**Assistant interpretation:** Resume from the current HOST-001 diary/task state, read the diary/docmgr operating instructions, implement all Phase 11 tasks, and commit focused backend/frontend/docs slices.

**Inferred user intent:** Finish the site settings/domains/configuration milestone with traceable code changes and enough diary detail for later review or continuation.

**Commit (code):** 21d2c172d05b4ca9abaedea452414089bc1d8f62 — "Add site settings and domain APIs"

**Commit (frontend):** 2c99c2875bed875617ce0bf8f1b9d20f5cbe9cf8 — "Add site settings dashboard page"

**Commit (polish):** 92ffade — "Polish pending domain DTOs"

**Commit (embedded assets):** a7ffffef97f847b8b26eb1f75202603fd0e2de9d — "Rebuild embedded dashboard for site settings"

### What I did

- Read `/home/manuel/.pi/agent/skills/diary/SKILL.md` and `/home/manuel/.pi/agent/skills/docmgr/SKILL.md` before coding.
- Re-read the HOST-001 diary and Phase 11 task list.
- Added migration `008_phase11_site_settings_domains.sql` for `site_config` and domain indexes.
- Added sqlc queries and store wrappers for:
  - site config list/upsert/delete,
  - site domain create/list/list-verified/get/verify/delete,
  - capability upsert wrapper.
- Added `SiteService` methods for config, domains, capabilities, and environment-placeholder authorization.
- Added HTTP APIs:
  - `GET/PUT/DELETE /api/v1/sites/{site_id}/config`,
  - `GET/PUT /api/v1/sites/{site_id}/capabilities`,
  - `GET/POST /api/v1/sites/{site_id}/domains`,
  - `POST /api/v1/sites/{site_id}/domains/{domain_id}/verify`,
  - `DELETE /api/v1/sites/{site_id}/domains/{domain_id}`,
  - `GET /api/v1/sites/{site_id}/environment`.
- Added audit events for config upsert/delete, capability update, domain add/verify/delete.
- Included verified custom domains in runtime activation host specs.
- Added integration coverage for config, capability toggles, domain add/verify, environment placeholder, and audit visibility.
- Added RTK Query hooks/types/MSW fixtures for site settings.
- Added `SiteSettingsPage` with Storybook coverage and wired it into the site tabs/routes.
- Rebuilt embedded dashboard assets via `go run ./cmd/build-web`.
- Marked Phase 11 tasks and exit criteria complete.

### Why

- Phase 11’s goal was to separate code deployment from site configuration.
- Domain state needed a first-class API rather than only admin read-only inventory.
- Capability policy was already enforced in Phase 10, but site owners needed a UI/API to inspect and adjust it.
- Secrets/environment needed an explicit placeholder so v1 does not accidentally grow unsafe process-env exposure.

### What worked

- `sqlc generate` handled the new site queries cleanly.
- The existing `site_domains` table already matched the required custom-domain lifecycle, so Phase 11 only needed query/service/API/UI wiring plus one supporting config table.
- The dashboard Settings page built successfully with Vite and Storybook.
- The embedded dashboard pipeline produced new static assets successfully through Dagger.

### What didn't work

- The first backend test compile failed because I added a duplicate `createTestOrgViaAPI` helper in `site_settings_integration_test.go`; the package already had one in `deployments_integration_test.go`.

Exact failure:

```text
# github.com/go-go-golems/go-go-host/internal/httpapi [github.com/go-go-golems/go-go-host/internal/httpapi.test]
internal/httpapi/site_settings_integration_test.go:128:6: createTestOrgViaAPI redeclared in this block
	internal/httpapi/deployments_integration_test.go:62:6: other declaration of createTestOrgViaAPI
FAIL	github.com/go-go-golems/go-go-host/internal/httpapi [build failed]
FAIL
```

I removed the duplicate helper and reused the existing test helper.

### What I learned

- Base-domain primary host assignment was already implemented in `SiteService.CreateSite`, so Phase 11 mostly needed custom-domain lifecycle and dashboard surfacing.
- Runtime activation is the natural point to include verified custom domains because it already builds the host list for the supervisor.
- Keeping the environment/secrets endpoint read-only and descriptive is a useful guardrail: it makes the deferred security decision visible instead of leaving a vague gap.

### What was tricky to build

- DELETE with a JSON body is supported by the current RTK Query/fetch path, but it is less ergonomic than a path-key route. I kept the backend simple for config deletion (`DELETE /config` with `{key}`) because config keys may contain dots and would otherwise require careful path escaping.
- Verified domain hosts are included on activation, not hot-added to an already-running runtime. This avoids implicit runtime mutation from a domain verification endpoint, but it means operators need a future restart/reactivation affordance if they expect a verified domain to route immediately.
- Capability edits are owner-only, while config/domain edits allow owners and developers. That is intentional because capabilities change runtime security boundaries.

### What warrants a second pair of eyes

- Review `internal/control/services.go` authorization choices: developers can edit non-secret config/domains, but only owners can toggle capabilities.
- Review hostname validation; it intentionally requires at least one dot and lowercases/trims a trailing dot.
- Review whether domain verification should remain a manual placeholder in v1 or perform DNS TXT/CNAME checks before marking verified.
- Review whether config deletion should become `DELETE /api/v1/sites/{site_id}/config/{key}` with URL-escaped keys.

### What should be done in the future

- Add DNS verification checks for `_go-go-host.<hostname>` TXT or a documented CNAME target.
- Hot-add verified domains to running runtimes or prompt the operator to restart/reactivate.
- Add encrypted secret storage and injection design before exposing any secret/runtime environment API.
- Add an admin write policy page for domain allow/deny rules if platform-level domain governance becomes necessary.

### Code review instructions

Start with backend control flow:

1. `internal/control/services.go` — site config, domain, capability, and environment-placeholder service logic.
2. `internal/httpapi/site_settings.go` — HTTP DTOs and handlers.
3. `internal/store/queries/sites.sql` and `internal/store/sites.go` — persistence API.
4. `internal/control/deployments.go` — verified custom domains added to activation hosts.
5. `web/admin/src/pages/SiteSettingsPage/SiteSettingsPage.tsx` — operator UI.
6. `web/admin/src/services/goGoHostApi.ts` — RTK Query surface.

Validate with:

```bash
go test ./...
pnpm --dir web/admin build
make storybook-build
go run ./cmd/build-web
docmgr doctor --ticket HOST-001-GO-GO-HOST-V1 --stale-after 30
```

### Technical details

New API summary:

```http
GET    /api/v1/sites/{site_id}/config
PUT    /api/v1/sites/{site_id}/config
DELETE /api/v1/sites/{site_id}/config
GET    /api/v1/sites/{site_id}/capabilities
PUT    /api/v1/sites/{site_id}/capabilities
GET    /api/v1/sites/{site_id}/domains
POST   /api/v1/sites/{site_id}/domains
POST   /api/v1/sites/{site_id}/domains/{domain_id}/verify
DELETE /api/v1/sites/{site_id}/domains/{domain_id}
GET    /api/v1/sites/{site_id}/environment
```

New audit actions:

```text
site.config.upsert
site.config.delete
site.capability.update
site.domain.add
site.domain.verify
site.domain.delete
```

## Step 28: Completed Phase 12 maintenance and production hardening

I continued from Phase 11 into the final Phase 12 operability work. The platform now has operator-facing export and pruning APIs, CLI commands for maintenance workflows, stronger readiness checks, production deployment scaffolding, and runbooks for runtime crashes and capability review.

This step focused on recoverability and safe operations rather than new product features. Operators can export site metadata, SQLite data, and immutable deployment bundles before pruning or incident response, and they have a documented path for runtime restart/rollback and future capability review.

### Prompt Context

**User prompt (verbatim):** "continue"

**Assistant interpretation:** Continue with the next unchecked platform phase, which is Phase 12 backup/export/pruning/production hardening, while preserving the existing diary/commit discipline.

**Inferred user intent:** Carry the HOST-001 implementation through the remaining operability checklist without stopping for another planning prompt.

**Commit (code):** fdfe91022b5b44daca5b29babd7b1faf6fe659b0 — "Add maintenance export and pruning APIs"

**Commit (CLI):** 5b8be383ebab7fe83a1ef307779b20f7e95f3051 — "Add maintenance CLI commands"

**Commit (docs/deploy):** b62f0036b36a66929eb217bf4f6033bc9dc2ed3e — "Add production hardening runbooks"

**Commit (tests/script):** dcbf75c4643493747185be9935f65e4d911e9666 — "Add production smoke coverage"

**Commit (Docker polish):** fb9a43c — "Fix Dockerfile for sqlite CGO runtime"

**Commit (tasks/e2e script):** 35c0c29 — "Complete Phase 12 checklist"

### What I did

- Added maintenance sqlc queries and store wrappers for deployment pruning and audit retention.
- Added `MaintenanceService` with:
  - site metadata export,
  - site SQLite DB export path resolution,
  - deployment bundle export path resolution,
  - deployment prune workflow,
  - platform-admin audit retention workflow.
- Added HTTP APIs:
  - `GET /api/v1/sites/{site_id}/export/metadata`,
  - `GET /api/v1/sites/{site_id}/export/db`,
  - `GET /api/v1/deployments/{deployment_id}/bundle`,
  - `POST /api/v1/sites/{site_id}/deployments/prune`,
  - `POST /api/v1/admin/audit/retention`.
- Changed `/readyz` from a static response to DB ping plus data-dir writeability checks.
- Added integration coverage for readiness, metadata export, and prune dry-run.
- Added `go-go-host maintenance` CLI commands:
  - `maintenance export metadata`,
  - `maintenance export db`,
  - `maintenance export bundle`,
  - `maintenance prune-deployments`,
  - `maintenance audit-retention`.
- Added runtime concurrency/load smoke coverage in `internal/runtime/supervisor_test.go`.
- Added `scripts/final-e2e-playwright.mjs`, gated by `GO_GO_HOST_E2E=1`, to exercise browser rendering plus API deploy/activate/rollback/audit and signed-agent deploy flow against a live devctl stack.
- Added runbooks:
  - `runbooks/02-runtime-crash-restart-runbook.md`,
  - `runbooks/03-host-capability-security-review-checklist.md`.
- Added `configs/production.example.yaml`.
- Added a Dockerfile and fixed it to use CGO because the runtime depends on `github.com/mattn/go-sqlite3`.
- Marked Phase 12 tasks and exit criteria complete.

### Why

- Phase 12 is about making the MVP operable after development: backup/export, recovery, pruning, readiness, production configuration, and safety checklists.
- Export endpoints and CLI commands give operators a concrete evidence-preservation path before destructive actions.
- Pruning and audit retention need explicit APIs because otherwise operators must edit database rows and filesystem artifacts by hand.
- `/readyz` needs to prove dependencies are usable, not just that the HTTP process is alive.

### What worked

- The maintenance API integrated cleanly with existing site authorization helpers and deployment/store models.
- Metadata export was straightforward because Phase 11 had already made config/capabilities/domains first-class store surfaces.
- The runtime supervisor handled 400 concurrent requests in the new smoke test and reported the expected request count.
- `go test ./...`, `pnpm --dir web/admin build`, and `docmgr doctor` passed after the Phase 12 implementation.

### What didn't work

- The initial Dockerfile used `CGO_ENABLED=0` and a static distroless base. That is wrong for this repo because the hosted runtime uses `github.com/mattn/go-sqlite3`, which requires CGO. I changed the Dockerfile to build with `CGO_ENABLED=1` and run on `debian:bookworm-slim` with certificates and a non-root user.

### What I learned

- The Phase 10/11 work made Phase 12 simpler: bundle hashes, site config, domain rows, and capability policy all export naturally as metadata.
- Maintenance endpoints need path containment checks because deployment bundle refs and SQLite paths are ultimately filesystem paths.
- The operator CLI does not need to understand the internal data-dir layout; it can download from authenticated API endpoints.

### What was tricky to build

- `deployment.bundle_ref` is stored as a filesystem path after validation. Export code therefore checks that the resolved path stays inside `dataDir` before serving it.
- Pruning must never delete the currently active deployment, even if its status or timestamp would otherwise match. The store query excludes the active deployment ID.
- Audit retention is platform-admin-only and records its own audit event after deletion, so the retention event itself remains visible.
- The Playwright E2E is intentionally gated by `GO_GO_HOST_E2E=1` because it needs a running devctl stack, a test bundle, and the Playwright package in the operator environment.

### What warrants a second pair of eyes

- Review prune policy semantics: default statuses are `rejected,superseded`; active deployment is excluded; `keepLatest` preserves the newest matching old deployments.
- Review audit-retention behavior and whether dry-run should count candidate rows instead of only returning a descriptive no-delete response.
- Review Dockerfile runtime dependencies if the final deployment target is distroless or Alpine; sqlite CGO makes static builds non-trivial.
- Review export endpoints for filename/content-type expectations and whether large files need streaming-specific limits.

### What should be done in the future

- Add scheduled/background pruning and retention jobs instead of only operator-triggered APIs.
- Add signed manifests for exported metadata bundles.
- Add object-storage backup targets once production storage is chosen.
- Add a real CI job for the Playwright E2E when devctl services can be provisioned in CI.

### Code review instructions

Start with:

1. `internal/control/maintenance.go` — authorization, path containment, prune and retention behavior.
2. `internal/httpapi/maintenance.go` — API DTOs, downloads, time parsing.
3. `cmd/go-go-host/cmds/maintenance.go` — operator CLI UX.
4. `internal/httpapi/handler.go` — `/readyz` dependency checks and maintenance route wiring.
5. `Dockerfile` and `configs/production.example.yaml` — production packaging/config assumptions.
6. `scripts/final-e2e-playwright.mjs` — live devctl E2E workflow.

Validate with:

```bash
go test ./...
pnpm --dir web/admin build
docmgr doctor --ticket HOST-001-GO-GO-HOST-V1 --stale-after 30
```

Optional live E2E:

```bash
devctl up --force
GO_GO_HOST_E2E=1 GO_GO_HOST_TEST_BUNDLE=/tmp/go-go-host-test-bundle.tar.gz node scripts/final-e2e-playwright.mjs
```

### Technical details

New maintenance APIs:

```http
GET  /api/v1/sites/{site_id}/export/metadata
GET  /api/v1/sites/{site_id}/export/db
GET  /api/v1/deployments/{deployment_id}/bundle
POST /api/v1/sites/{site_id}/deployments/prune
POST /api/v1/admin/audit/retention
```

New maintenance CLI examples:

```bash
go run ./cmd/go-go-host maintenance export metadata --site-id SITE_ID --dev-user OWNER -o site.json
go run ./cmd/go-go-host maintenance export db --site-id SITE_ID --dev-user OWNER -o site.sqlite
go run ./cmd/go-go-host maintenance export bundle --deployment-id DEP_ID --dev-user OWNER -o dep.tar.gz
go run ./cmd/go-go-host maintenance prune-deployments --site-id SITE_ID --dev-user OWNER --older-than 30d --dry-run=false
go run ./cmd/go-go-host maintenance audit-retention --dev-user dev-user --older-than 90d --dry-run=false
```

## Step 29: Added bundled developer, JavaScript API, and agent guides

I added the missing educational layer for app authors and deployment agents. The platform already had command help and agent quick-start pages, but developers needed a more complete path from "what is a bundle?" to "how do I build, deploy, debug, and operate a full app?" The new help entries are written in a textbook style: start with the mental model, ground it in runnable examples, then explain the reference details and failure modes.

The key result is that the docs now travel with the binaries. A developer can run `go-go-host help developer-guide` and `go-go-host help js-api-reference`; an automation engineer can run `go-go-host-agent help agent-guide`. That makes the CLI itself a teaching surface rather than only a command surface.

### Prompt Context

**User prompt (verbatim):** "Create a developer-guide and a js-api-reference and an agent-guide document that will allow a developer/agent to build and deploy full apps against the system, with examples and such. nAnd then a detailed js api reference for the actual in depth JS reference.

Use a textbook writing style for the documents, to make them really nice (see docs). They shoul dbe bundled as glazed help entries into the relevant binaries)"

**Assistant interpretation:** Add rich Glazed help pages to the human and agent CLIs, using textbook-style explanations and concrete examples, covering app development, the hosted JS APIs/DSLs, and machine-agent deployment workflows.

**Inferred user intent:** Make the platform self-teaching for developers and CI agents so they can build and deploy real apps without reverse-engineering source code or relying on chat history.

**Commit (docs):** f526c9f8da97f2734b55840df50898da1a46ee4e — "Add bundled developer and JS API guides"

### What I did

- Read the textbook-authoring and Glazed help page authoring skills before drafting.
- Added `cmd/go-go-host/doc/developer-guide.md` as a top-level tutorial for bundle layout, site creation, deployment, activation, config, capabilities, operations, and debugging.
- Added `cmd/go-go-host/doc/js-api-reference.md` as a detailed reference for manifest fields, `express`, request/response objects, `ui.dsl`, `database`/`db`, `db.guard`, capabilities, forbidden modules, and troubleshooting.
- Added `cmd/go-go-host/doc/agent-guide.md` for human operators who create agents, grants, enrollment tokens, key rotation, and audit review.
- Added `cmd/go-go-host-agent/doc/agent-guide.md` for machine/CI operators who run keygen, enroll, deploy, activate, rotate keys, and troubleshoot signed requests.
- Replaced the older `agent-setup` preview page with a current short bridge to the new agent guides.
- Verified the new pages are bundled and discoverable with:
  - `go run ./cmd/go-go-host help developer-guide`,
  - `go run ./cmd/go-go-host help js-api-reference`,
  - `go run ./cmd/go-go-host help agent-guide`,
  - `go run ./cmd/go-go-host-agent help agent-guide`.

### Why

- The code had a real hosted JS surface, but developers should not have to read Go source files to learn it.
- The agent workflow is security-sensitive; it needs prose that explains why grants, keys, upload tokens, and activation are separated.
- Glazed help entries are the right delivery mechanism because this is CLI-first developer/operator documentation.

### What worked

- The existing CLI roots already load embedded docs via `doc.AddDocToHelpSystem`, so adding Markdown files was enough to bundle the new help entries.
- The Glazed help renderer successfully rendered the long-form pages with tables, code blocks, and section structure.
- Existing `go test ./cmd/go-go-host ./cmd/go-go-host-agent` passed after adding the docs.

### What didn't work

- No code failures occurred in this doc slice. I did notice unrelated dirty files in `cmd/go-go-host*/cmds/support.go` and an untracked `HOST-005-E2E-FIXES` ticket workspace; I left them out of this commit because they were not part of the requested documentation change.

### What I learned

- The old `agent-setup` page still described agent commands as deferred even though Phase 9/10 implemented them. Updating that page was necessary to avoid contradicting the new guides.
- The JS API reference needs to be both tutorial and contract. It should say what is supported, but it should also explicitly say what is not supported so app authors do not infer unsafe capabilities.

### What was tricky to build

- The docs needed to be thorough without inventing APIs. I grounded the JS API reference in the current runtime modules: `express`, `ui.dsl`/`ui`, `database`/`db`, `db.guard`, request/response DTOs, and manifest validation.
- The agent story has two audiences. Human operators create grants and tokens; machines hold keys and sign deploys. I split the docs between the human CLI and agent CLI so each binary teaches the workflow from its own point of view.

### What warrants a second pair of eyes

- Review the `db.guard.configure` wording. The function exists, but the docs intentionally discourage treating app-level guard reconfiguration as the normal app contract.
- Review whether the JS API reference should include every go-go-goja utility function from `time`, `timer`, and `path`, or keep those as high-level utility modules until their upstream docs are linked.
- Review the examples against a live devctl stack to ensure copy/paste commands remain exact as CLI flags evolve.

### What should be done in the future

- Add a `go-go-host site template` or `go-go-host bundle init` command to generate the documented starter bundle.
- Add examples as real files under an examples directory and reference them from the help pages.
- Add a small JS API compatibility test that verifies the documented smoke app routes continue to work.

### Code review instructions

Review:

1. `cmd/go-go-host/doc/developer-guide.md` — complete developer path and operational narrative.
2. `cmd/go-go-host/doc/js-api-reference.md` — exact JS API contract.
3. `cmd/go-go-host/doc/agent-guide.md` — human/operator agent workflow.
4. `cmd/go-go-host-agent/doc/agent-guide.md` — machine/CI workflow.
5. `cmd/go-go-host/doc/agent-setup.md` — now-current bridge page.

Validate:

```bash
go test ./cmd/go-go-host ./cmd/go-go-host-agent
go run ./cmd/go-go-host help developer-guide
go run ./cmd/go-go-host help js-api-reference
go run ./cmd/go-go-host help agent-guide
go run ./cmd/go-go-host-agent help agent-guide
```

### Technical details

New help slugs:

```text
developer-guide
js-api-reference
agent-guide      # available in both go-go-host and go-go-host-agent, with audience-specific content
```
