---
Title: Investigation Diary
Ticket: HOST-009-DOCS-ONBOARDING-PLAYGROUND
DocType: reference
Status: active
Topics: [frontend, documentation, design-system, developer-experience]
---

# Investigation Diary

## 2026-05-12 — Ticket setup and initial brainstorm

User request: create a new ticket for dashboard documentation pages. The dashboard should become welcoming and intuitive for users, not only a management console. Users need to manage sites and agents, but more importantly they need to learn the hosted JavaScript APIs and try them out. User specifically suggested a live playground and a full runthrough walkthrough.

Created ticket:

```text
HOST-009-DOCS-ONBOARDING-PLAYGROUND
```

Created docs:

```text
design-doc/01-dashboard-docs-and-playground-brainstorm.md
reference/01-investigation-diary.md
```

Initial repo evidence:

- Existing CLI/developer docs live under `cmd/go-go-host/doc/`.
- The most relevant source documents are `developer-guide.md`, `js-api-reference.md`, `deploy-workflow.md`, `agent-guide.md`, and `agent-setup.md`.
- A durable runnable example exists under `examples/hello-beta/`.
- The hosted JS runtime modules live under `internal/sitejs/`, especially `web`, `uidsl`, and `dbguard`.

Initial design direction: the dashboard needs a first-class Docs/Learn area that combines prose, examples, copyable bundles, and a safe live playground. The playground should start as a client-side bundle authoring and preview environment, then graduate into controlled server-side dry-run/deploy flows.

## 2026-05-12 — reMarkable upload

Uploaded the brainstorm/design document to reMarkable:

```text
/ai/2026/05/12/HOST-009-DOCS-ONBOARDING-PLAYGROUND/HOST-009_Dashboard_Docs_Playground_Brainstorm.pdf
```

Command output:

```text
OK: uploaded HOST-009_Dashboard_Docs_Playground_Brainstorm.pdf -> /ai/2026/05/12/HOST-009-DOCS-ONBOARDING-PLAYGROUND
```
