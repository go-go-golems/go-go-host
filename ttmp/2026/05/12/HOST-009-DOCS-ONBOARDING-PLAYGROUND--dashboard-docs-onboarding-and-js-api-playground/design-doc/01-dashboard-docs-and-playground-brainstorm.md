---
Title: Dashboard Docs and Playground Brainstorm
Ticket: HOST-009-DOCS-ONBOARDING-PLAYGROUND
Status: active
Topics:
    - frontend
    - documentation
    - design-system
    - developer-experience
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: cmd/go-go-host/doc/agent-guide.md
      Note: Existing agent guide to adapt into dashboard agent onboarding
    - Path: cmd/go-go-host/doc/developer-guide.md
      Note: Existing full developer guide to adapt into dashboard Learn and Quickstart pages
    - Path: cmd/go-go-host/doc/js-api-reference.md
      Note: Existing JS API reference to break into interactive module docs
    - Path: examples/hello-beta
      Note: Durable demo bundle source recommended as the first playground template
    - Path: internal/sitejs
      Note: Runtime JS API implementation source for docs/playground truth
    - Path: web/admin/src/app/routes.tsx
      Note: Frontend route table that will host Learn/Docs/Playground pages
ExternalSources: []
Summary: ""
LastUpdated: 0001-01-01T00:00:00Z
WhatFor: ""
WhenToUse: ""
---


# Dashboard Docs and Playground Brainstorm

## Executive summary

The dashboard should not only be a control panel for sites, deployments, agents, domains, and quotas. It should also be the place where a new developer learns the platform, writes their first hosted JavaScript app, tests the API surface, packages a bundle, deploys it, and understands what happened. Today the product has useful CLI docs and a working demo bundle, but the dashboard experience is still mostly operational. HOST-009 should add a **Docs / Learn / Playground** area that makes the platform feel welcoming and self-explanatory.

My recommendation is to build this in three layers:

1. **Guided onboarding path** — a checklist-style runthrough from first org/site to first live app.
2. **JS API docs pages** — structured, searchable, example-heavy docs for `express`, `ui.dsl`, `database`, `db.guard`, assets, config, and platform context.
3. **Live playground** — an in-dashboard code editor with runnable examples, manifest preview, bundle preview, validation feedback, and eventually a safe dry-run/deploy path.

The playground is probably the most important differentiator. It turns documentation from passive reference into a learning loop: edit code, see what it renders, understand required capabilities, package it, deploy it, verify it.

## Problem statement

A hosting platform has two very different audiences:

- **Operators**, who need to manage sites, agents, deployments, domains, capabilities, and audit trails.
- **Developers**, who need to understand what code they can write, what APIs exist, how bundles work, how capabilities are requested, and how to debug failures.

The current dashboard is becoming better as an operator console, but a first-time developer still has to infer too much:

- What does a hosted app look like?
- Which files go in a bundle?
- What does `go-go-host.json` mean?
- What JavaScript modules are available?
- What is the difference between `express`, `ui.dsl`, `database`, and `assets`?
- How do I test locally or in the dashboard?
- Why was my bundle rejected?
- How do agents differ from human uploads?
- What is safe/unsafe in hosted v1?

If we want users to succeed quickly, the dashboard should teach these answers directly and interactively.

## Proposed top-level navigation

Add a first-class navigation item to the app shell:

```text
Learn
```

or:

```text
Docs
```

I prefer **Learn** for the user-facing nav label because it feels more welcoming. The route can still be `/app/docs` or `/app/learn`.

Suggested routes:

```text
/app/learn
/app/learn/quickstart
/app/learn/playground
/app/learn/js-api
/app/learn/js-api/express
/app/learn/js-api/ui-dsl
/app/learn/js-api/database
/app/learn/js-api/assets
/app/learn/js-api/platform-context
/app/learn/deployments
/app/learn/agents
/app/learn/troubleshooting
```

For site-specific docs, we can include contextual links:

```text
/app/orgs/:orgId/sites/:siteId/learn
/app/orgs/:orgId/sites/:siteId/playground
```

But I would start with a global Learn area that can optionally preselect an org/site.

## Recommended page set

### 1. Welcome / learning dashboard

Goal: orient the user in one screen.

Sections:

- **Build your first hosted app** — start quickstart.
- **Try the JS API playground** — open playground with hello app loaded.
- **Deploy with a human upload** — upload bundle flow.
- **Deploy with an agent** — CI/deploy-agent path.
- **Read the JS API reference** — module cards.
- **Troubleshoot a rejected deployment** — validation and runtime failure guide.

This should be friendly and task-based, not a wall of docs.

Good card titles:

```text
Start with Hello
Render HTML with ui.dsl
Store visits with database
Serve CSS from /assets
Inspect req.platform
Deploy from the dashboard
Deploy from CI with an agent
Fix a rejected bundle
```

### 2. Full runthrough walkthrough

This should be the canonical tutorial. It should be runnable end-to-end.

Proposed title:

```text
From zero to live site in 10 minutes
```

Steps:

1. Create or choose an organization.
2. Create a site.
3. Open the playground with the `hello-beta` template.
4. Read the manifest.
5. Run/preview the app locally in the playground.
6. Inspect routes: `/`, `/platform`, `/db`, `/assets/style.css`.
7. Package the bundle.
8. Upload and validate deployment.
9. Activate deployment.
10. Visit generated public host.
11. Create an agent.
12. Enroll agent key.
13. Deploy same bundle via agent with `bundlePath`.
14. Revoke temporary test agent.

This is more than a tutorial; it becomes a smoke-test teaching flow.

### 3. JS API overview

A map of available hosted modules:

| Module | Teaches | User question |
|---|---|---|
| `express` | routes, handlers, JSON/html responses | How do I handle URLs? |
| `ui.dsl` / `ui` | HTML nodes and safe rendering | How do I render a page? |
| `database` / `db` | per-site SQLite | How do I store data? |
| `db.guard` | quota stats | How do I understand DB limits? |
| `assets` | static files under `/assets` | How do I include CSS/images? |
| `req.platform` | host/site/deployment metadata | How do I know where I'm running? |
| config | non-secret site config | How do operators tune my app? |

Each module page should include:

- mental model,
- minimal example,
- full example,
- common mistakes,
- capability requirements,
- playground button: "Open this example".

### 4. Live playground

This should be the centerpiece.

#### MVP playground

The first playground can be safe and purely dashboard-side:

- CodeMirror editor for `scripts/app.js`.
- CodeMirror editor for `go-go-host.json`.
- optional `assets/style.css` editor.
- template selector.
- manifest/capabilities preview.
- generated bundle file tree preview.
- static lint/checks:
  - manifest parses,
  - entrypoint exists,
  - `assetsDir` exists if assets used,
  - capabilities include modules used by code,
  - dangerous modules (`fs`, `exec`, `process.env`) highlighted as unsupported.
- copy/download bundle button.
- upload to selected site button.

This MVP does **not** need to execute untrusted JS in the browser. It can teach the bundle model and reduce friction.

#### Better playground

Next phase: server-side dry run API.

Add an endpoint like:

```text
POST /api/v1/playground/validate
POST /api/v1/sites/:siteId/playground/dry-run
```

Inputs:

```json
{
  "manifest": {...},
  "files": {
    "scripts/app.js": "...",
    "assets/style.css": "..."
  },
  "smokePath": "/"
}
```

Outputs:

```json
{
  "valid": true,
  "manifest": {...},
  "capabilities": ["express", "ui.dsl", "database", "assets"],
  "routes": ["/", "/platform", "/db"],
  "smoke": {
    "status": 200,
    "contentType": "text/html",
    "bodyPreview": "<html>..."
  },
  "warnings": [],
  "errors": []
}
```

This is the ideal learning loop: write code, run smoke, see rendered preview and validation messages before creating a deployment.

#### Best playground

Eventually:

- split-pane editor + rendered preview,
- selectable fake request path,
- fake platform context editor,
- in-memory or disposable SQLite DB for preview,
- route explorer,
- request/response inspector,
- one-click "Deploy to this site".

This would make go-go-host feel like a mini hosted-app IDE.

### 5. Agent onboarding docs

Agents are powerful but conceptually tricky. They deserve a friendly path.

Pages:

```text
What is an agent?
Create a CI deploy agent
Enrollment token vs signing key
Bundle path policy
Rotate and revoke keys
Audit agent deploys
```

The docs must make the identity split clear:

```text
Humans use OIDC.
Agents use enrollment tokens + Ed25519 signed requests.
```

Add diagrams showing:

```text
Create agent -> issue one-time enrollment token -> agent enrolls public key -> signed deploy request -> upload token -> bundle upload -> validation -> activation
```

### 6. Troubleshooting docs

These should be highly practical.

Topics:

- Bundle rejected: missing manifest.
- Bundle rejected: path traversal / unsafe archive.
- Capability denied.
- Route returns 500.
- Asset 404.
- Database quota warning.
- Runtime missing/stopped/failed.
- Agent signature rejected.
- Bundle path not allowed.
- Access token expired.

Every troubleshooting page should answer:

```text
What you see
What it means
How to fix it
Where to look in the dashboard
CLI command equivalent
```

## Playground templates

The first template set should be small but high quality.

### Template 1: Hello UI

Capabilities:

```json
["express", "ui.dsl"]
```

Teaches:

- route registration,
- `ui.page`,
- `ui.h1`, `ui.p`, links.

### Template 2: Hello with assets

Capabilities:

```json
["express", "ui.dsl", "assets"]
```

Teaches:

- `assetsDir`,
- `/assets/style.css`,
- `ui.link`.

### Template 3: Visits counter

Capabilities:

```json
["express", "ui.dsl", "database"]
```

Teaches:

- DB initialization,
- insert/query,
- state across requests.

### Template 4: Platform inspector

Capabilities:

```json
["express"]
```

Teaches:

- `req.platform`,
- JSON responses,
- runtime metadata.

### Template 5: Full hello-beta

Capabilities:

```json
["express", "ui.dsl", "database", "assets"]
```

Teaches the full happy path and should match `examples/hello-beta`.

## Implementation plan

### Phase 1 — Static docs shell

- Add Learn/Docs nav item.
- Add docs routes and page shells.
- Convert existing markdown docs into dashboard-readable content.
- Add module cards and quickstart cards.
- Add Storybook/MSW states.

Files likely touched:

```text
web/admin/src/app/routes.tsx
web/admin/src/components/organisms/OrgSidebar/OrgSidebar.tsx
web/admin/src/pages/LearnPage/
web/admin/src/pages/DocsPage/
web/admin/src/pages/PlaygroundPage/
cmd/go-go-host/doc/*.md
```

### Phase 2 — Playground MVP

- CodeMirror editors for manifest, app JS, CSS.
- Template selector.
- Bundle tree preview.
- Static validation and capability inference.
- Download/copy bundle.
- Optional upload to selected site using existing deployment API.

### Phase 3 — Server-side dry run

- Add API endpoint for validation/dry-run.
- Reuse existing bundle validator and runtime dry-run behavior.
- Return structured diagnostics and response preview.

### Phase 4 — Guided full runthrough

- Checklist that integrates with live org/site state.
- "Mark done" / auto-detect completed steps.
- Links from each step to the relevant dashboard page or playground template.

## Design recommendations

1. **Start with Learn, not Reference.** New users need a path before they need an encyclopedia.
2. **Make every docs page actionable.** Each section should include "Open in playground", "Copy", or "Go to page".
3. **Use templates as docs.** The best docs are runnable examples.
4. **Teach capabilities early.** Capabilities explain why code works or fails.
5. **Keep agents separate but connected.** Agent onboarding should come after first human deploy, not before.
6. **Surface validation as teaching.** Bundle errors should link to docs pages explaining the failure.
7. **Keep playground safe by default.** Start with static/client-side checks, then add controlled server dry-run.

## Open questions

1. Should the playground execute JavaScript in-browser, or only send to server dry-run?
2. Should docs content be stored as Markdown files and rendered, or authored as React content/components?
3. Should playground templates live under `examples/` and be imported, or duplicated as frontend fixtures?
4. Should docs be global (`/app/learn`) or contextual under an org/site?
5. How much of agent setup should be runnable in-dashboard versus copy/paste CLI guidance?
6. Should generated playground bundles be deployable immediately to beta sites?

## My recommendation

Build the first slice as:

```text
/app/learn
/app/learn/quickstart
/app/learn/js-api
/app/learn/playground
```

with three polished experiences:

1. **Welcome Learn page** with task cards.
2. **Quickstart runthrough** based on `examples/hello-beta`.
3. **Playground MVP** with CodeMirror editors, template selector, manifest/capability checks, file tree preview, and upload/download actions.

Then add server-side dry-run once the UI has proven useful.

This sequence gives users value quickly without creating an unsafe remote-code execution surface before we have designed the sandbox/dry-run API properly.
