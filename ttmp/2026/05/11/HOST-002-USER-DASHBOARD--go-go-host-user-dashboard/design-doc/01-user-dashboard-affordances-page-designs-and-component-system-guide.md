---
Title: User Dashboard Affordances, Page Designs, and Component System Guide
Ticket: HOST-002-USER-DASHBOARD
Status: active
Topics:
    - goja
    - hosting
    - go-go-host
    - rtk-query
    - storybook
    - vm-runtime
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: internal/control/deployments.go
      Note: Deployment validation and activation service behavior
    - Path: internal/httpapi/agents_audit.go
      Note: Agent and audit API contracts
    - Path: internal/httpapi/api.go
      Note: Session
    - Path: internal/httpapi/deployments.go
      Note: Deployment API contracts for upload
    - Path: internal/httpapi/handler.go
      Note: Dashboard mount points and API route registration
    - Path: internal/httpapi/runtime.go
      Note: Runtime status API contract
    - Path: internal/store/migrations/001_initial_schema.sql
      Note: Control-plane schema backing dashboard entities
    - Path: internal/store/migrations/003_runtime_status.sql
      Note: Runtime status schema for dashboard counters
    - Path: internal/webadmin/handler.go
      Note: Current placeholder dashboard handler to replace with embedded SPA
    - Path: plugins/go-go-host-devctl.py
      Note: devctl launch plan for Storybook and Vite
    - Path: web/admin/.storybook/preview.tsx
      Note: Storybook MSW and mock provider setup
    - Path: web/admin/package.json
      Note: Dashboard package scripts and dependencies
    - Path: web/admin/src/services/goGoHostApi.ts
      Note: Initial RTK Query API slice
    - Path: web/admin/src/services/msw/handlers.ts
      Note: MSW handlers for page stories
ExternalSources: []
Summary: 'Design for the /app user dashboard: affordances, pages, ASCII wireframes, component taxonomy, Storybook/MSW requirements, and intern-oriented implementation plan.'
LastUpdated: 2026-05-11T18:51:11.491639125-04:00
WhatFor: Use this as the implementation guide for the Phase 7 React dashboard and Storybook component system.
WhenToUse: Read before creating web/admin or implementing any /app dashboard page, widget, RTK Query endpoint, MSW handler, or Storybook story.
---



# User Dashboard Affordances, Page Designs, and Component System Guide

## Executive Summary

Phase 7 builds the first real **user/org developer dashboard** for `go-go-host`. The dashboard lives at `/app/*` and replaces the current placeholder HTML with an embedded React/Vite application. Its job is to make the control-plane features that now exist in Go usable without the CLI: organization selection, site management, deployment upload/activation/rollback, runtime inspection, agent visibility, audit exploration, and quota/status awareness.

This document is intentionally written for a new intern. It explains:

1. What the product is.
2. Which backend APIs already exist.
3. Which user affordances the dashboard must provide.
4. Which pages to build, including ASCII screenshots.
5. Which reusable widgets belong in the component system.
6. How to organize each widget directory.
7. How to write Storybook stories for every widget and every page.
8. How to use MSW/fake RTK Query state so Storybook works without a live daemon.
9. How to phase implementation so the dashboard stays testable and reviewable.

The design assumes the dashboard is built on top of:

- React + TypeScript + Vite.
- Redux Toolkit + RTK Query.
- Storybook.
- MSW for API mocks.
- `@go-go-golems/os-core` for base theme/layout primitives.
- A local component taxonomy: **atoms**, **molecules**, **organisms**, and **pages**.

The dashboard should not be a single large app file. Every reusable UI affordance should live in its own directory, with tests where appropriate and at least one Storybook story. Pages also need Storybook stories, because page-level composition and loading/error/empty states are where most dashboard regressions happen.

## Problem Statement

The backend and CLI now support the core v1 workflow, but the web dashboard is still a placeholder. Current server routing mounts `/app/` and `/admin/` to `webadmin.NewPlaceholderHandler()` in `internal/httpapi/handler.go:19-20`, and that placeholder explicitly says the dashboard will be embedded later in `internal/webadmin/handler.go:8-20`.

The system needs a dashboard that supports these jobs:

- A developer can see who they are logged in as and which organizations they can access.
- A developer can create and inspect sites.
- A developer can upload a deployment bundle and read validation reports.
- A developer can activate a deployment and roll back to a previous one.
- A developer can see runtime state, hosts, counters, and errors.
- A developer can list/create deployment agents.
- A developer can inspect audit events for organization activity.
- An operator or power user can understand when a failure is auth-related, validation-related, runtime-related, or infrastructure-related.

The dashboard must be designed before implementation because it touches many system boundaries: auth state, org membership, site state, runtime supervisor state, deployments, file upload, agents, audit, and future admin surfaces.

## Current-State Analysis with Evidence

### Server mounts and routing

Observed current state:

- `/app/` and `/admin/` are mounted in `internal/httpapi/handler.go:19-20`.
- The mounted handler is a placeholder in `internal/webadmin/handler.go:8-20`.
- Authenticated API routes are registered in `internal/httpapi/handler.go:29-45`.
- API routes are wrapped with auth middleware in `internal/httpapi/handler.go:46-53`.
- The final handler adds request IDs and falls back to the runtime supervisor in `internal/httpapi/handler.go:55`.

Implication for Phase 7:

- The React app should be served from `/app/*`.
- The platform admin console remains `/admin/*` and should not be mixed into the user dashboard.
- The SPA must avoid stealing `/api/*` routes.
- The app can call relative `/api/v1/...` URLs and does not need a separate API origin in production.

### User/session/org APIs

Observed current state:

- `GET /api/v1/me` returns user, memberships, and platform admin boolean in `internal/httpapi/api.go:40-61`.
- `POST /api/v1/orgs` creates an org in `internal/httpapi/api.go:65-83`.
- `GET /api/v1/orgs/{org_id}/sites` lists sites in `internal/httpapi/api.go:86-108`.
- `POST /api/v1/orgs/{org_id}/sites` creates a site in `internal/httpapi/api.go:111-133`.
- Membership roles are stored as `org_owner`, `org_developer`, and `org_viewer` in `internal/store/migrations/001_initial_schema.sql:25-31`.

Implication for Phase 7:

- The app should bootstrap session state from `/api/v1/me`.
- Organization selection should be based on `memberships` from `/me`.
- Route guards must distinguish:
  - unauthenticated user,
  - user with no orgs,
  - user with org access,
  - platform admin user.

### Site/deployment/runtime APIs

Observed current state:

- Deployment upload endpoint exists at `POST /api/v1/sites/{site_id}/deployments` in `internal/httpapi/handler.go:40` and `internal/httpapi/deployments.go:33-72`.
- Deployment list endpoint exists at `GET /api/v1/sites/{site_id}/deployments` in `internal/httpapi/handler.go:41` and `internal/httpapi/deployments.go:75-92`.
- Deployment detail endpoint exists at `GET /api/v1/deployments/{deployment_id}` in `internal/httpapi/handler.go:43` and `internal/httpapi/deployments.go:95-108`.
- Activation endpoint exists at `POST /api/v1/deployments/{deployment_id}/activate` in `internal/httpapi/handler.go:44` and `internal/httpapi/deployments.go:111-124`.
- Rollback endpoint exists at `POST /api/v1/sites/{site_id}/rollback` in `internal/httpapi/handler.go:42` and `internal/httpapi/deployments.go:127-140`.
- Runtime endpoint exists at `GET /api/v1/sites/{site_id}/runtime` in `internal/httpapi/runtime.go:9-32`.

The deployment service performs validation and dry-run runtime loading:

- Placeholder deployment row is created before validation in `internal/control/deployments.go:56-64`.
- Bundle validation/storage runs in `internal/control/deployments.go:65-68`.
- Dry-run runtime load and smoke check run in `internal/control/deployments.go:76-89`.
- Rejected deployments get status `rejected` in `internal/control/deployments.go:91-94`.
- Activation creates a runtime spec and calls the supervisor in `internal/control/deployments.go:135-166`.

Implication for Phase 7:

- Deployment upload UI must show nested validation reports.
- Deployment detail UI must distinguish `uploaded`, `validated`, `rejected`, `active`, and `superseded`.
- Activation and rollback are mutating actions and must be confirmation-gated.
- Runtime status is not just a badge; it includes hosts, deployment ID, request counters, and last error.

### Agents and audit APIs

Observed current state:

- Agents can be listed, created, and revoked under org-scoped routes in `internal/httpapi/handler.go:35-38` and `internal/httpapi/agents_audit.go:42-96`.
- Audit can be listed with filters in `internal/httpapi/agents_audit.go:98-117`.
- Agent schema exists in `internal/store/migrations/001_initial_schema.sql:111-147`.
- Audit schema exists in `internal/store/migrations/001_initial_schema.sql:150-162`.
- Agent service currently supports create/list/revoke only in `internal/control/agents.go:14-51`.
- Audit service lists org-scoped events in `internal/control/agents.go:54-59`.

Implication for Phase 7:

- Agent UI should be honest about current capabilities: list/create/revoke is supported; key enrollment/grants/deploy-run flows are future work.
- Audit UI should support filters that map directly to current query params: `resource_id`, `actor_type`, `actor_id`, `action`, `limit`.

### Data model surface relevant to UI

The dashboard needs to understand these backend entities:

```text
User
  id, email, displayName

Membership
  orgId, orgSlug, orgName, role

Org
  id, slug, name

Site
  id, orgId, slug, name, primaryHost, status, activeDeploymentId

Deployment
  id, siteId, version, status, bundleRef, unpackedPath,
  manifestJson, validationJson, createdByType, createdById,
  createdAt, activatedAt

RuntimeStatus
  siteId, orgId, deploymentId, hosts, status, startedAt,
  lastError, requestsTotal, errorsTotal

Agent
  id, orgId, name, status, createdByUserId, createdAt, lastSeenAt

AuditEvent
  id, orgId, actorType, actorId, action, resourceType,
  resourceId, ipAddress, userAgent, metadataJson, createdAt
```

The Postgres schema anchors this model in `internal/store/migrations/001_initial_schema.sql:38-162` and runtime counters in `internal/store/migrations/003_runtime_status.sql:1-14`.

## Product Affordances

An **affordance** is an action or perception the UI makes possible. Phase 7 should be evaluated by whether these affordances are obvious and reliable.

### Global affordances

1. **Know who I am**
   - Show current display name/email.
   - Show current org and role.
   - Show whether dev auth is active if `/api/v1/config` says `devAuth: true`.

2. **Switch organizations**
   - Use memberships from `/api/v1/me`.
   - Preserve selected org in URL when possible: `/app/orgs/:orgId/...`.
   - If the user has one org, auto-select it.
   - If the user has zero orgs, show onboarding.

3. **Understand system status**
   - Health/config status in the shell.
   - Runtime badges on site cards and detail pages.
   - Clear empty/error/loading states.

4. **Navigate by task, not by database table**
   - Sites, Deployments, Agents, Audit, Usage, Members.
   - Site detail should act as the hub for deployment/runtime/host workflows.

5. **Use safe confirmations for destructive or traffic-changing actions**
   - Activate deployment changes traffic.
   - Rollback changes traffic.
   - Revoke agent removes a future automation identity.

### Organization affordances

1. List accessible orgs.
2. Create an org if allowed by current v1 policy.
3. See org role: owner/developer/viewer.
4. Restrict owner/developer actions when viewer.
5. Link to org-scoped audit.

### Site affordances

1. List sites in selected org.
2. Create a site.
3. See site status and primary host.
4. Copy primary host.
5. Open public site URL or show curl command with Host header for localhost.
6. See active deployment version and ID.
7. See runtime status badge.
8. See runtime counters and last error.
9. Navigate to deployments, agents, audit filtered for site activity.

### Deployment affordances

1. Upload `.tar.gz` or `.zip` bundle.
2. Set deployment message/channel when backend supports/records those fields.
3. Read validation report.
4. Understand why validation failed.
5. Compare deployment versions.
6. Activate a validated deployment.
7. Roll back to previous validated/superseded deployment.
8. See manifest and effective validation output.
9. Know which deployment is active.
10. Know when activation failed because runtime load failed.

### Runtime affordances

1. See status: stopped, starting, ready, failed, draining.
2. See deployment ID attached to runtime.
3. See hosts served by runtime.
4. See started time.
5. See total requests and errors.
6. See last error.
7. Provide a refresh button.
8. Link runtime deployment ID to deployment detail.

### Agent affordances

1. List agents in org.
2. Create an agent record.
3. Revoke an agent.
4. Show future-work state for keys/grants/enrollment if not implemented yet.
5. Explain that full deploy-run token flow is not yet available.

### Audit affordances

1. List org audit events.
2. Filter by action.
3. Filter by actor type/id.
4. Filter by resource ID.
5. Show timestamps and metadata JSON.
6. Link resource IDs to site/deployment/agent pages when possible.
7. Preserve filter state in query params.

### Usage/quota affordances

Some quota fields exist in backend schema, but no quota API exists yet. The UI should still reserve a page/section with explicit placeholder behavior.

1. Display bundle max bytes when API exists.
2. Display DB soft/hard limits when API exists.
3. Display runtime request/error counters now.
4. Show a clear "quota API pending" panel until endpoints exist.

## Information Architecture

The `/app` user dashboard should use this route structure:

```text
/app
  /login-or-dev-session
  /orgs/new
  /orgs/:orgId
    /sites
    /sites/new
    /sites/:siteId
      /overview
      /deployments
      /deployments/:deploymentId
      /runtime
      /usage
    /agents
    /audit
    /members
```

Recommended React Router object sketch:

```ts
const routes = [
  { path: "/app", element: <AppRoot />,
    children: [
      { index: true, element: <OrgRedirectOrOnboarding /> },
      { path: "orgs/new", element: <CreateOrgPage /> },
      { path: "orgs/:orgId", element: <OrgLayout />,
        children: [
          { index: true, element: <Navigate to="sites" replace /> },
          { path: "sites", element: <SitesPage /> },
          { path: "sites/new", element: <CreateSitePage /> },
          { path: "sites/:siteId", element: <SiteLayout />,
            children: [
              { index: true, element: <SiteOverviewPage /> },
              { path: "deployments", element: <DeploymentsPage /> },
              { path: "deployments/:deploymentId", element: <DeploymentDetailPage /> },
              { path: "runtime", element: <RuntimePage /> },
              { path: "usage", element: <UsagePage /> },
            ]
          },
          { path: "agents", element: <AgentsPage /> },
          { path: "audit", element: <AuditPage /> },
          { path: "members", element: <MembersPage /> },
        ]
      },
    ]
  }
];
```

## Page Descriptions and ASCII Screenshots

ASCII screenshots define layout intent. They are not pixel-perfect. They are contracts for information hierarchy.

### 1. App Bootstrap / Loading State

Purpose: show the shell while `/api/v1/me` and `/api/v1/config` are loading.

```text
┌──────────────────────────────────────────────────────────────────────────────┐
│ go-go-host                                                        Loading... │
├──────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│                         ┌────────────────────────┐                           │
│                         │  Loading session...    │                           │
│                         │  /api/v1/me            │                           │
│                         └────────────────────────┘                           │
│                                                                              │
└──────────────────────────────────────────────────────────────────────────────┘
```

Storybook states:

- `LoadingSession`
- `ConfigLoadedSessionLoading`
- `SessionError401`
- `StoreHydratedWithOneOrg`

### 2. No Organizations / Onboarding Page

Purpose: guide a user with no memberships.

```text
┌──────────────────────────────────────────────────────────────────────────────┐
│ go-go-host                                      dev-user@dev.local  [Config] │
├──────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Welcome to go-go-host                                                       │
│  You do not belong to any organizations yet.                                 │
│                                                                              │
│  ┌────────────────────────────────────────────────────────────────────────┐  │
│  │ Create your first organization                                         │  │
│  │ Slug  [ my-org________________ ]                                      │  │
│  │ Name  [ My Organization________ ]                                     │  │
│  │                                                        [Create org]    │  │
│  └────────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└──────────────────────────────────────────────────────────────────────────────┘
```

Primary APIs:

- `GET /api/v1/me`
- `POST /api/v1/orgs`

Storybook states:

- empty memberships,
- create success,
- slug validation error,
- backend duplicate slug error.

### 3. Org Sites Page

Purpose: list sites and make the active deployment/runtime visible.

```text
┌──────────────────────────────────────────────────────────────────────────────┐
│ go-go-host   Org: Demo Org ▾                       Alice      Dev auth ON    │
├───────────────┬──────────────────────────────────────────────────────────────┤
│ Sites         │ Sites                                             [New site] │
│ Deployments   │ ┌────────────┬───────────────┬─────────┬─────────┬────────┐ │
│ Agents        │ │ Name       │ Host          │ Runtime │ Active  │ Status │ │
│ Audit         │ ├────────────┼───────────────┼─────────┼─────────┼────────┤ │
│ Members       │ │ hello      │ hello.local…  │ Ready   │ v4      │ active │ │
│ Usage         │ │ docs       │ docs.local…   │ Stopped │ —       │ prov…  │ │
│               │ └────────────┴───────────────┴─────────┴─────────┴────────┘ │
│               │                                                              │
│               │ Empty state: "No sites yet. Create one or deploy from CLI." │
└───────────────┴──────────────────────────────────────────────────────────────┘
```

Primary APIs:

- `GET /api/v1/orgs/{org_id}/sites`
- optionally fan out `GET /api/v1/sites/{site_id}/runtime` for visible cards.

Storybook states:

- empty list,
- two sites with mixed runtime status,
- loading runtime badges,
- site list error,
- viewer role with disabled create button.

### 4. Create Site Page

Purpose: create a site and preview its generated hostname.

```text
┌──────────────────────────────────────────────────────────────────────────────┐
│ Demo Org / Sites / New                                                       │
├──────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Create site                                                                 │
│                                                                              │
│  Slug        [ hello____________________ ]                                    │
│  Name        [ Hello Site_______________ ]                                    │
│                                                                              │
│  Preview host: hello.localhost                                               │
│                                                                              │
│  [Cancel]                                                  [Create site]     │
│                                                                              │
└──────────────────────────────────────────────────────────────────────────────┘
```

Primary API:

- `POST /api/v1/orgs/{org_id}/sites`

Storybook states:

- valid form,
- invalid slug,
- submit loading,
- duplicate host/backend error,
- viewer role disabled.

### 5. Site Overview Page

Purpose: central hub for a single site.

```text
┌──────────────────────────────────────────────────────────────────────────────┐
│ Demo Org / Sites / hello                                      [Open public]  │
├──────────────────────────────────────────────────────────────────────────────┤
│ [Overview] [Deployments] [Runtime] [Usage]                                  │
├──────────────────────────────────────────────────────────────────────────────┤
│ ┌────────────────────────────┐ ┌───────────────────────────────────────────┐ │
│ │ Site                       │ │ Runtime                                   │ │
│ │ Name: Hello Site           │ │ Status: Ready ●                           │ │
│ │ Host: hello.localhost [⧉]  │ │ Deployment: dep_abc / v4                 │ │
│ │ Status: active             │ │ Requests: 1,234    Errors: 2             │ │
│ └────────────────────────────┘ └───────────────────────────────────────────┘ │
│                                                                              │
│ ┌──────────────────────────────────────────────────────────────────────────┐ │
│ │ Active deployment                                                       │ │
│ │ v4  dep_abc  activated 2026-05-11 18:20                                │ │
│ │ [View deployment] [Upload new bundle] [Rollback]                        │ │
│ └──────────────────────────────────────────────────────────────────────────┘ │
└──────────────────────────────────────────────────────────────────────────────┘
```

Primary APIs:

- `GET /api/v1/sites/{site_id}/runtime`
- `GET /api/v1/sites/{site_id}/deployments`

Storybook states:

- active runtime,
- stopped runtime,
- failed runtime with last error,
- no deployment yet,
- deployment list fetch failed.

### 6. Deployments Page

Purpose: upload bundles and list deployment history.

```text
┌──────────────────────────────────────────────────────────────────────────────┐
│ hello / Deployments                                             [Upload]     │
├──────────────────────────────────────────────────────────────────────────────┤
│ ┌──────────────────────────────────────────────────────────────────────────┐ │
│ │ Drop bundle here or choose file                                         │ │
│ │ Supported: .tar.gz, .zip     Manifest: go-go-host.json                  │ │
│ │ Channel [ default____ ] Message [ ________________________________ ]    │ │
│ └──────────────────────────────────────────────────────────────────────────┘ │
│                                                                              │
│ ┌────┬────────────┬────────────┬──────────────┬───────────────────────────┐ │
│ │ v  │ Status     │ Created    │ Activated    │ Actions                   │ │
│ ├────┼────────────┼────────────┼──────────────┼───────────────────────────┤ │
│ │ 4  │ active     │ 18:10      │ 18:20        │ [View]                    │ │
│ │ 3  │ superseded │ 17:55      │ 18:00        │ [View] [Activate]         │ │
│ │ 2  │ rejected   │ 17:40      │ —            │ [View report]             │ │
│ └────┴────────────┴────────────┴──────────────┴───────────────────────────┘ │
└──────────────────────────────────────────────────────────────────────────────┘
```

Primary APIs:

- `POST /api/v1/sites/{site_id}/deployments`
- `GET /api/v1/sites/{site_id}/deployments`
- `POST /api/v1/deployments/{deployment_id}/activate`

Storybook states:

- empty deployment history,
- upload drag active,
- upload validating,
- validation success,
- validation rejected with nested errors,
- activation confirmation dialog,
- activation error.

### 7. Deployment Detail Page

Purpose: show manifest, validation report, artifact paths, and actions.

```text
┌──────────────────────────────────────────────────────────────────────────────┐
│ hello / Deployments / v4                                      [Activate]    │
├──────────────────────────────────────────────────────────────────────────────┤
│ Status: validated   Deployment: dep_abc   Created by: user usr_123          │
│                                                                              │
│ ┌──────────────────────────────┐ ┌─────────────────────────────────────────┐ │
│ │ Manifest                     │ │ Validation report                       │ │
│ │ scriptsDir: scripts          │ │ Valid: true                             │ │
│ │ assetsDir: assets            │ │ Files: 12                               │ │
│ │ smokePath: /                 │ │ Bytes: 44,203                           │ │
│ │ capabilities: time,timer     │ │ Effective caps: time,timer              │ │
│ └──────────────────────────────┘ └─────────────────────────────────────────┘ │
│                                                                              │
│ Validation errors                                                            │
│ ┌──────────────────────────────────────────────────────────────────────────┐ │
│ │ none                                                                     │ │
│ └──────────────────────────────────────────────────────────────────────────┘ │
└──────────────────────────────────────────────────────────────────────────────┘
```

Primary APIs:

- `GET /api/v1/deployments/{deployment_id}`
- `POST /api/v1/deployments/{deployment_id}/activate`

Storybook states:

- active deployment,
- validated but inactive deployment,
- rejected deployment with errors,
- malformed legacy validation JSON fallback.

### 8. Runtime Page

Purpose: inspect runtime status and counters.

```text
┌──────────────────────────────────────────────────────────────────────────────┐
│ hello / Runtime                                                  [Refresh]   │
├──────────────────────────────────────────────────────────────────────────────┤
│ ┌───────────────┬──────────────────────────────────────────────────────────┐ │
│ │ Status        │ Ready ●                                                  │ │
│ │ Deployment    │ dep_abc [View]                                           │ │
│ │ Hosts         │ hello.localhost                                          │ │
│ │ Started       │ 2026-05-11 18:20:00                                      │ │
│ │ Requests      │ 1,234                                                    │ │
│ │ Errors        │ 2                                                        │ │
│ │ Last error    │ —                                                        │ │
│ └───────────────┴──────────────────────────────────────────────────────────┘ │
└──────────────────────────────────────────────────────────────────────────────┘
```

Primary API:

- `GET /api/v1/sites/{site_id}/runtime`

Storybook states:

- ready,
- stopped,
- failed with last error,
- loading,
- unauthorized/forbidden.

### 9. Agents Page

Purpose: list and create automation identities while making future gaps clear.

```text
┌──────────────────────────────────────────────────────────────────────────────┐
│ Demo Org / Agents                                           [Create agent]  │
├──────────────────────────────────────────────────────────────────────────────┤
│ ┌──────────────────────────────────────────────────────────────────────────┐ │
│ │ Agent deployment is in preview. Keys/grants/enrollment are coming later. │ │
│ └──────────────────────────────────────────────────────────────────────────┘ │
│                                                                              │
│ ┌────────────┬──────────┬──────────────┬────────────┬────────────────────┐ │
│ │ Name       │ Status   │ Created      │ Last seen  │ Actions            │ │
│ ├────────────┼──────────┼──────────────┼────────────┼────────────────────┤ │
│ │ ci-bot     │ active   │ 18:35        │ —          │ [Revoke]           │ │
│ └────────────┴──────────┴──────────────┴────────────┴────────────────────┘ │
└──────────────────────────────────────────────────────────────────────────────┘
```

Primary APIs:

- `GET /api/v1/orgs/{org_id}/agents`
- `POST /api/v1/orgs/{org_id}/agents`
- `POST /api/v1/orgs/{org_id}/agents/{agent_id}/revoke`

Storybook states:

- no agents,
- active/revoked agents,
- create dialog,
- revoke confirmation,
- viewer role disabled create/revoke.

### 10. Audit Page

Purpose: org-scoped activity timeline with filters.

```text
┌──────────────────────────────────────────────────────────────────────────────┐
│ Demo Org / Audit                                                             │
├──────────────────────────────────────────────────────────────────────────────┤
│ Action [ deployment.activate ▾ ] Actor [ user usr_123____ ] Resource [____] │
│ [Apply filters] [Reset]                                                      │
│                                                                              │
│ ┌────────────┬─────────────────────┬──────────────┬──────────────────────┐ │
│ │ Time       │ Action              │ Actor        │ Resource             │ │
│ ├────────────┼─────────────────────┼──────────────┼──────────────────────┤ │
│ │ 18:35      │ agent.create        │ user usr_123 │ agent agt_123        │ │
│ │ 18:20      │ deployment.activate │ user usr_123 │ deployment dep_abc   │ │
│ └────────────┴─────────────────────┴──────────────┴──────────────────────┘ │
│                                                                              │
│ Selected event metadata                                                       │
│ { "ip": "", "userAgent": "" }                                             │
└──────────────────────────────────────────────────────────────────────────────┘
```

Primary API:

- `GET /api/v1/orgs/{org_id}/audit?resource_id=&actor_type=&actor_id=&action=&limit=`

Storybook states:

- populated timeline,
- no matching events,
- filter error,
- selected metadata expanded.

### 11. Usage Page

Purpose: show runtime counters now and reserve space for quota APIs.

```text
┌──────────────────────────────────────────────────────────────────────────────┐
│ hello / Usage                                                                │
├──────────────────────────────────────────────────────────────────────────────┤
│ ┌──────────────────────┐ ┌──────────────────────┐ ┌──────────────────────┐ │
│ │ Requests             │ │ Errors               │ │ Error rate           │ │
│ │ 1,234                │ │ 2                    │ │ 0.16%                │ │
│ └──────────────────────┘ └──────────────────────┘ └──────────────────────┘ │
│                                                                              │
│ ┌──────────────────────────────────────────────────────────────────────────┐ │
│ │ Storage quotas                                                          │ │
│ │ Quota API pending. Backend has site_quotas table but no endpoint yet.   │ │
│ └──────────────────────────────────────────────────────────────────────────┘ │
└──────────────────────────────────────────────────────────────────────────────┘
```

Storybook states:

- runtime counters loaded,
- stopped runtime,
- quota unavailable.

### 12. Members Page

Purpose: reserve org membership UI and show current user roles. Membership mutation endpoints do not exist yet.

```text
┌──────────────────────────────────────────────────────────────────────────────┐
│ Demo Org / Members                                                           │
├──────────────────────────────────────────────────────────────────────────────┤
│ ┌──────────────────────────────────────────────────────────────────────────┐ │
│ │ Membership management API pending.                                      │ │
│ │ Current membership data comes from /api/v1/me.                          │ │
│ └──────────────────────────────────────────────────────────────────────────┘ │
│                                                                              │
│ ┌──────────────┬────────────────────────────┬──────────────┐               │
│ │ User         │ Email                      │ Role         │               │
│ ├──────────────┼────────────────────────────┼──────────────┤               │
│ │ Alice        │ alice@dev.local            │ org_owner    │               │
│ └──────────────┴────────────────────────────┴──────────────┘               │
└──────────────────────────────────────────────────────────────────────────────┘
```

Storybook states:

- owner,
- developer,
- viewer,
- API pending message.

## Component System Design

Use this directory layout:

```text
web/admin/src/
  app/
    store.ts
    App.tsx
    routes.tsx
    providers/
      AppProviders.tsx
      MockAppProviders.tsx
  services/
    goGoHostApi.ts
    msw/
      handlers.ts
      fixtures.ts
      browser.ts
      server.ts
  components/
    atoms/
      RuntimeStatusDot/
        RuntimeStatusDot.tsx
        RuntimeStatusDot.stories.tsx
        index.ts
      CopyButton/
      EmptyState/
      ErrorCallout/
      LoadingBlock/
      Timestamp/
      RoleBadge/
      StatusPill/
      CodeBlock/
      JsonTree/
    molecules/
      RuntimeBadge/
      OrgSwitcher/
      SiteHostCopy/
      DeploymentStatusPill/
      ValidationSummary/
      FileDropZone/
      ConfirmActionDialog/
      FilterToolbar/
      MetricCard/
      ManifestSummary/
      AuditEventRow/
      AgentStatusBadge/
    organisms/
      AppShell/
      OrgSidebar/
      SiteHeader/
      SiteTabs/
      SitesTable/
      DeploymentTimeline/
      DeploymentUploadPanel/
      ValidationReportPanel/
      RuntimeStatusPanel/
      AgentsTable/
      AuditTimeline/
      QuotaPanel/
      MembersTable/
    pages/
      AppBootstrapPage/
      NoOrgsPage/
      SitesPage/
      CreateSitePage/
      SiteOverviewPage/
      DeploymentsPage/
      DeploymentDetailPage/
      RuntimePage/
      AgentsPage/
      AuditPage/
      UsagePage/
      MembersPage/
```

Every directory must contain:

```text
WidgetName/
  WidgetName.tsx
  WidgetName.stories.tsx
  index.ts
```

For pages, use the same rule:

```text
pages/SitesPage/
  SitesPage.tsx
  SitesPage.stories.tsx
  index.ts
```

If a component has nontrivial behavior, add tests later:

```text
WidgetName.test.tsx
```

### Atoms

Atoms are tiny primitives. They should not call RTK Query directly.

| Atom | Purpose | Props sketch | Story states |
|---|---|---|---|
| `RuntimeStatusDot` | colored dot for runtime status | `{ status: RuntimeStatus }` | ready, failed, stopped, starting |
| `StatusPill` | labeled status chip | `{ status, tone? }` | deployment statuses, site statuses |
| `RoleBadge` | displays org role | `{ role }` | owner, developer, viewer |
| `CopyButton` | copy text to clipboard | `{ value, label? }` | idle, copied, error |
| `EmptyState` | empty list placeholder | `{ title, body, action? }` | with/without action |
| `ErrorCallout` | visible error block | `{ title, error, retry? }` | auth, validation, network |
| `LoadingBlock` | skeleton/loading placeholder | `{ lines? }` | small/large |
| `Timestamp` | consistent time rendering | `{ value, mode }` | absolute, relative, empty |
| `CodeBlock` | monospace block | `{ code, language? }` | JSON, shell |
| `JsonTree` | readable JSON object | `{ value }` | manifest, validation report |

Example atom pseudocode:

```tsx
export function RuntimeStatusDot({ status }: { status: RuntimeStatus }) {
  const tone = statusTone[status] ?? "neutral";
  return <span className={styles.dot} data-tone={tone} aria-label={`runtime ${status}`} />;
}
```

### Molecules

Molecules combine atoms into self-contained controls. They should still avoid page-level data fetching unless the data is local to the molecule.

| Molecule | Purpose | Depends on | Story states |
|---|---|---|---|
| `RuntimeBadge` | status dot + text + counters | `RuntimeStatusDot`, `StatusPill` | ready, failed, stopped |
| `OrgSwitcher` | select org membership | `RoleBadge` | one org, many orgs, no orgs |
| `SiteHostCopy` | host text + copy/open actions | `CopyButton` | localhost, public domain |
| `DeploymentStatusPill` | deployment status chip | `StatusPill` | uploaded/validated/rejected/active/superseded |
| `ValidationSummary` | valid/errors/files/bytes | `JsonTree`, `ErrorCallout` | valid, invalid, warnings |
| `FileDropZone` | drag/drop bundle input | atoms | idle, drag, selected, invalid extension |
| `ConfirmActionDialog` | confirmation modal | os-core dialog primitives | activate, rollback, revoke |
| `FilterToolbar` | audit filters | form atoms | no filters, active filters |
| `MetricCard` | numeric dashboard card | atom typography | request/error/error-rate |
| `ManifestSummary` | manifest fields | `JsonTree` | minimal, full |
| `AuditEventRow` | one audit event | `Timestamp`, `StatusPill` | selected/unselected |
| `AgentStatusBadge` | active/revoked status | `StatusPill` | active/revoked |

### Organisms

Organisms compose molecules and usually receive data from pages. They should be easy to story with fixtures.

| Organism | Purpose | Story states |
|---|---|---|
| `AppShell` | top bar, side nav, content region | logged in, dev auth, error banner |
| `OrgSidebar` | org-scoped navigation | sites active, agents active, audit active |
| `SiteHeader` | site title, host, actions | active, provisioning, suspended |
| `SiteTabs` | overview/deployments/runtime/usage nav | each selected tab |
| `SitesTable` | site list | empty, many, mixed runtime statuses |
| `DeploymentTimeline` | deployment history | active/superseded/rejected mix |
| `DeploymentUploadPanel` | upload form and result | idle, uploading, success, rejected |
| `ValidationReportPanel` | nested validation report | valid, invalid, malformed JSON |
| `RuntimeStatusPanel` | runtime details | ready, stopped, failed |
| `AgentsTable` | agent list/actions | empty, active/revoked, viewer readonly |
| `AuditTimeline` | filterable audit list | empty, many, selected metadata |
| `QuotaPanel` | quotas and usage | API pending, counters only |
| `MembersTable` | membership display | owner/developer/viewer |

### Pages

Pages are route-level components. Pages may call RTK Query hooks. Every page needs Storybook stories with MSW/fake store.

| Page | API dependencies | Required stories |
|---|---|---|
| `AppBootstrapPage` | `/me`, `/config` | loading, unauthorized, one org, no orgs |
| `NoOrgsPage` | `/me`, `POST /orgs` | empty, create success, create error |
| `SitesPage` | `/orgs/:orgId/sites`, runtime per site | empty, populated, loading, error |
| `CreateSitePage` | `POST /orgs/:orgId/sites` | valid, invalid, forbidden |
| `SiteOverviewPage` | runtime, deployments | active, no deployment, failed runtime |
| `DeploymentsPage` | upload/list/activate | upload success, rejected, activation confirm |
| `DeploymentDetailPage` | get/activate | active, validated, rejected |
| `RuntimePage` | runtime | ready, stopped, failed |
| `AgentsPage` | agents list/create/revoke | empty, populated, create, revoke |
| `AuditPage` | audit list | populated, filtered empty, metadata selected |
| `UsagePage` | runtime now, quota later | counters, quota pending |
| `MembersPage` | `/me` | owner/developer/viewer |

## Storybook and MSW Requirements

### Rule: every widget and page gets a story

For every component directory:

- `ComponentName.stories.tsx` must exist.
- Stories must include at least:
  - default state,
  - loading state if async data is relevant,
  - error state if data can fail,
  - empty state if the component renders collections,
  - permission-restricted state if actions depend on role.

### MSW structure

Use MSW to mock backend APIs in page stories:

```text
web/admin/src/services/msw/
  fixtures.ts
  handlers.ts
  browser.ts
  server.ts
```

Fixtures should be typed and reusable:

```ts
export const fixtures = {
  user: {
    id: "usr_123",
    email: "alice@dev.local",
    displayName: "Alice",
  },
  memberships: [
    { orgId: "org_123", orgSlug: "demo", orgName: "Demo Org", role: "org_owner" },
  ],
  sites: [
    { id: "site_123", orgId: "org_123", slug: "hello", name: "Hello", primaryHost: "hello.localhost", status: "active", activeDeploymentId: "dep_4" },
  ],
};
```

Handlers should mirror current Go routes:

```ts
export const handlers = [
  http.get("/api/v1/me", () => HttpResponse.json({
    user: fixtures.user,
    memberships: fixtures.memberships,
    platformAdmin: false,
  })),
  http.get("/api/v1/orgs/:orgId/sites", ({ params }) => HttpResponse.json(fixtures.sites)),
  http.get("/api/v1/sites/:siteId/runtime", () => HttpResponse.json(fixtures.runtime.ready)),
  http.get("/api/v1/sites/:siteId/deployments", () => HttpResponse.json(fixtures.deployments)),
  http.get("/api/v1/orgs/:orgId/agents", () => HttpResponse.json(fixtures.agents)),
  http.get("/api/v1/orgs/:orgId/audit", () => HttpResponse.json(fixtures.auditEvents)),
];
```

Per-story overrides should model edge states:

```ts
export const RuntimeFailed: Story = {
  parameters: {
    msw: {
      handlers: [
        http.get("/api/v1/sites/:siteId/runtime", () => HttpResponse.json(fixtures.runtime.failed)),
      ],
    },
  },
};
```

### Fake store / RTK Query store

Use a Storybook decorator that creates a fresh Redux store per story:

```tsx
export function withMockStore(story: StoryFn) {
  const store = makeStore({ preloadedState: {} });
  return <Provider store={store}>{story()}</Provider>;
}
```

Do not share a singleton store between stories. Shared stores cause cache bleed and make stories flaky.

### Storybook preview

`.storybook/preview.tsx` should:

1. Import `@go-go-golems/os-core` theme CSS.
2. Import the selected desktop theme.
3. Start MSW Storybook addon.
4. Wrap all stories in app providers.

Pseudocode:

```tsx
import "@go-go-golems/os-core/theme.css";
import "@go-go-golems/os-core/themes/desktop.css";
import { initialize, mswLoader } from "msw-storybook-addon";
import { withAppProviders } from "../src/app/providers/MockAppProviders";

initialize();

export const decorators = [withAppProviders];
export const loaders = [mswLoader];
```

## RTK Query API Design

Create `web/admin/src/services/goGoHostApi.ts` with one API slice:

```ts
export const goGoHostApi = createApi({
  reducerPath: "goGoHostApi",
  baseQuery: fetchBaseQuery({ baseUrl: "/api/v1" }),
  tagTypes: [
    "Me", "Org", "Site", "Deployment", "Runtime", "Agent", "Audit", "Config",
  ],
  endpoints: (build) => ({
    getConfig: build.query<ConfigResponse, void>({ query: () => "/config", providesTags: ["Config"] }),
    getMe: build.query<MeResponse, void>({ query: () => "/me", providesTags: ["Me", "Org"] }),
    createOrg: build.mutation<Org, CreateOrgRequest>({ query: (body) => ({ url: "/orgs", method: "POST", body }), invalidatesTags: ["Me", "Org"] }),
    listSites: build.query<Site[], string>({ query: (orgId) => `/orgs/${orgId}/sites`, providesTags: (result) => siteTags(result) }),
    createSite: build.mutation<Site, { orgId: string; body: CreateSiteRequest }>({ query: ({ orgId, body }) => ({ url: `/orgs/${orgId}/sites`, method: "POST", body }), invalidatesTags: (_r, _e, arg) => [{ type: "Site", id: `ORG:${arg.orgId}` }] }),
    getRuntime: build.query<RuntimeStatus, string>({ query: (siteId) => `/sites/${siteId}/runtime`, providesTags: (_r, _e, siteId) => [{ type: "Runtime", id: siteId }] }),
    listDeployments: build.query<Deployment[], string>({ query: (siteId) => `/sites/${siteId}/deployments`, providesTags: (result, _e, siteId) => deploymentTags(result, siteId) }),
    uploadDeployment: build.mutation<UploadDeploymentResponse, UploadDeploymentArgs>({ query: ({ siteId, formData }) => ({ url: `/sites/${siteId}/deployments`, method: "POST", body: formData }), invalidatesTags: (_r, _e, arg) => [{ type: "Deployment", id: `SITE:${arg.siteId}` }] }),
    getDeployment: build.query<Deployment, string>({ query: (deploymentId) => `/deployments/${deploymentId}`, providesTags: (_r, _e, id) => [{ type: "Deployment", id }] }),
    activateDeployment: build.mutation<Deployment, string>({ query: (id) => ({ url: `/deployments/${id}/activate`, method: "POST" }), invalidatesTags: ["Deployment", "Runtime", "Site"] }),
    rollbackSite: build.mutation<Deployment, string>({ query: (siteId) => ({ url: `/sites/${siteId}/rollback`, method: "POST" }), invalidatesTags: ["Deployment", "Runtime", "Site"] }),
    listAgents: build.query<Agent[], string>({ query: (orgId) => `/orgs/${orgId}/agents`, providesTags: ["Agent"] }),
    createAgent: build.mutation<Agent, { orgId: string; name: string }>({ query: ({ orgId, name }) => ({ url: `/orgs/${orgId}/agents`, method: "POST", body: { name } }), invalidatesTags: ["Agent", "Audit"] }),
    revokeAgent: build.mutation<{ status: string; agentId: string }, { orgId: string; agentId: string }>({ query: ({ orgId, agentId }) => ({ url: `/orgs/${orgId}/agents/${agentId}/revoke`, method: "POST" }), invalidatesTags: ["Agent", "Audit"] }),
    listAudit: build.query<AuditEvent[], AuditFilter>({ query: ({ orgId, ...query }) => ({ url: `/orgs/${orgId}/audit`, params: query }), providesTags: ["Audit"] }),
  }),
});
```

Important implementation notes:

- Keep generated hook names stable.
- Parse `manifestJson` and `validationJson` in UI selectors, not inside components repeatedly.
- Treat malformed JSON as displayable error state, not as a page crash.
- Use RTK Query tags to invalidate deployments/runtime after activation and rollback.

## TypeScript Types

Define API types in `web/admin/src/services/types.ts`:

```ts
export type Role = "org_owner" | "org_developer" | "org_viewer";
export type DeploymentStatus = "uploaded" | "validated" | "rejected" | "active" | "superseded";
export type RuntimeState = "starting" | "ready" | "failed" | "stopped" | "draining";

export interface MeResponse {
  user: { id: string; email: string; displayName: string };
  memberships: Membership[];
  platformAdmin: boolean;
}

export interface Site {
  id: string;
  orgId: string;
  slug: string;
  name: string;
  primaryHost: string;
  status: string;
  activeDeploymentId: string;
}

export interface Deployment {
  id: string;
  siteId: string;
  version: number;
  status: DeploymentStatus;
  bundleRef: string;
  unpackedPath: string;
  manifestJson: string;
  validationJson: string;
  createdByType: string;
  createdById: string;
  createdAt: string;
  activatedAt?: string;
}

export interface ValidationReport {
  valid: boolean;
  errors?: string[];
  warnings?: string[];
  files: number;
  bytes: number;
  requestedCapabilities?: string[];
  effectiveCapabilities?: string[];
}
```

## Theming and go-go-os-core Integration

The exact `@go-go-golems/os-core` API should be verified during implementation, but the intended contract is:

1. Use os-core theme tokens as the base design language.
2. Do not hardcode one-off colors in widgets.
3. Component CSS should use semantic variables:

```css
.runtimeStatusDot[data-tone="success"] {
  background: var(--os-color-success-500);
}

.card {
  background: var(--os-surface-panel);
  border: 1px solid var(--os-border-subtle);
  border-radius: var(--os-radius-md);
}
```

4. Page layout should use os-core shell/panel primitives if available.
5. If os-core lacks a primitive, create local components that still use os-core tokens.

## Implementation Plan

### Phase 7.0: Project scaffold

Create:

```text
web/admin/package.json
web/admin/vite.config.ts
web/admin/tsconfig.json
web/admin/index.html
web/admin/src/main.tsx
web/admin/src/app/store.ts
web/admin/src/app/App.tsx
web/admin/src/app/routes.tsx
web/admin/.storybook/main.ts
web/admin/.storybook/preview.tsx
```

Backend serving options:

- Development: Vite dev server proxies `/api` to `go-go-hostd`.
- Production: `go generate` or build target copies `web/admin/dist` into an embedded Go filesystem served by `internal/webadmin`.

### Phase 7.1: API and mock foundation

Implement:

- `services/types.ts`
- `services/goGoHostApi.ts`
- `services/msw/fixtures.ts`
- `services/msw/handlers.ts`
- Storybook MSW setup.

Acceptance criteria:

- `getMe`, `listSites`, `listDeployments`, `getRuntime`, `listAgents`, `listAudit` work in Storybook without a daemon.
- One sample page story renders entirely from MSW.

### Phase 7.2: Component atoms and molecules

Implement atoms first, then molecules:

1. `StatusPill`
2. `RuntimeStatusDot`
3. `RoleBadge`
4. `CopyButton`
5. `ErrorCallout`
6. `LoadingBlock`
7. `RuntimeBadge`
8. `ValidationSummary`
9. `FileDropZone`
10. `ConfirmActionDialog`

Acceptance criteria:

- Every component has stories.
- Stories include non-happy states.
- Components do not call RTK Query directly.

### Phase 7.3: Shell and routing

Implement:

- `AppShell`
- `OrgSwitcher`
- `OrgSidebar`
- `RequireSession`
- `RequireOrgAccess`
- route layout components.

Acceptance criteria:

- `/app` loads session.
- no-org user sees onboarding.
- one-org user reaches sites page.
- org selector changes route.

### Phase 7.4: Sites and site overview

Implement:

- `SitesPage`
- `CreateSitePage`
- `SiteOverviewPage`
- `SiteHeader`
- `SitesTable`

Acceptance criteria:

- user can create a site from UI.
- site list shows primary host and runtime badge.
- all states are storybooked.

### Phase 7.5: Deployments and runtime

Implement:

- `DeploymentsPage`
- `DeploymentDetailPage`
- `RuntimePage`
- `DeploymentUploadPanel`
- `DeploymentTimeline`
- `ValidationReportPanel`
- `RuntimeStatusPanel`

Acceptance criteria:

- upload success shows validation report.
- rejected upload shows validation errors.
- activation confirmation works.
- rollback confirmation works.
- runtime page refreshes status.

### Phase 7.6: Agents, audit, usage, members

Implement:

- `AgentsPage`
- `AuditPage`
- `UsagePage`
- `MembersPage`
- corresponding organisms.

Acceptance criteria:

- agents can be listed and created.
- agent revoke is confirmation-gated.
- audit filters map to API query params.
- usage/members pages clearly mark API-pending features.

### Phase 7.7: Backend embed and CI

Implement:

- Makefile targets:
  - `make web-install`
  - `make web-dev`
  - `make web-build`
  - `make storybook`
  - `make storybook-build`
- Go embed path in `internal/webadmin`.
- GitHub Actions or existing CI integration for `pnpm build` and `pnpm storybook:build`.

Acceptance criteria:

- `go-go-hostd` serves built app at `/app/`.
- Storybook builds in CI.
- SPA fallback works for nested `/app/orgs/:orgId/...` routes.

## Testing Strategy

### Unit/component tests

Use Vitest + React Testing Library for components with behavior:

- `FileDropZone`
- `ConfirmActionDialog`
- `FilterToolbar`
- JSON parsing helpers
- route guard helpers

### Storybook tests

Storybook is mandatory for:

- all atoms,
- all molecules,
- all organisms,
- all pages.

Each page story should run with MSW and fake store. Do not require a live Go daemon.

### Playwright smoke tests

Add at least one dashboard smoke test:

1. Start `go-go-hostd` in dev mode.
2. Visit `/app`.
3. Use dev auth mechanism or test fixture headers if supported.
4. Assert site list renders.
5. Create a site or use seeded data.
6. Navigate to deployments page.

If browser auth setup is not ready, Storybook interaction tests may cover more UI states than Playwright initially.

### Backend compatibility tests

As the dashboard is implemented, add Go tests only when server serving changes:

- `/app/` serves `index.html`.
- `/app/orgs/org_123/sites` also serves `index.html`.
- `/api/v1/...` is not swallowed by SPA fallback.

## API Reference for Dashboard Implementation

### Session and config

```http
GET /api/v1/config
GET /api/v1/me
```

### Organizations and sites

```http
GET  /api/v1/orgs
POST /api/v1/orgs
GET  /api/v1/orgs/{org_id}/sites
POST /api/v1/orgs/{org_id}/sites
```

### Runtime

```http
GET /api/v1/sites/{site_id}/runtime
```

### Deployments

```http
POST /api/v1/sites/{site_id}/deployments   multipart form: bundle, message, channel
GET  /api/v1/sites/{site_id}/deployments
GET  /api/v1/deployments/{deployment_id}
POST /api/v1/deployments/{deployment_id}/activate
POST /api/v1/sites/{site_id}/rollback
```

### Agents

```http
GET  /api/v1/orgs/{org_id}/agents
POST /api/v1/orgs/{org_id}/agents          JSON: { "name": "ci-bot" }
POST /api/v1/orgs/{org_id}/agents/{agent_id}/revoke
```

### Audit

```http
GET /api/v1/orgs/{org_id}/audit?resource_id=&actor_type=&actor_id=&action=&limit=
```

## Design Decisions

### Decision 1: Pages also get Storybook stories

Page stories are mandatory because most bugs happen in composition states: loading plus partial data, empty lists, permission restrictions, nested validation errors, and malformed JSON. Component stories alone are not enough.

### Decision 2: Data fetching lives in pages and hooks, not atoms/molecules

Atoms and molecules should be reusable and easy to story with plain props. RTK Query belongs in page components, layout guards, or dedicated route-level hooks.

### Decision 3: Use MSW as the default Storybook data source

MSW keeps stories close to real API contracts. Fake props are fine for atoms/molecules, but page stories should exercise request/response shape.

### Decision 4: Keep `/app` and `/admin` separate

The existing server already reserves both surfaces. `/app` is the org/user developer dashboard. `/admin` is the future platform operator console. Mixing them will create authorization and navigation confusion.

### Decision 5: Be honest about incomplete backend features

Agent keys, grants, deploy-run token issuance, membership mutation, and quota APIs are not complete yet. The UI should show preview/API-pending states rather than inventing fake working controls.

## Alternatives Considered

### Alternative: Build only pages first, component system later

Rejected. The user explicitly requested a component system with atoms/molecules/organisms and Storybook stories. Building pages first would produce duplicated status badges, tables, dialogs, and validation panels.

### Alternative: Use static fixtures only, no MSW

Rejected. Static props are useful for atoms, but page stories should validate backend contracts. MSW catches route and JSON-shape drift earlier.

### Alternative: Put all widgets in one `components/` folder

Rejected. One flat folder becomes hard to navigate. One directory per widget gives clean ownership and enforces stories.

### Alternative: Implement `/admin` together with `/app`

Rejected for Phase 7. The requested scope is the user dashboard. Admin pages have different authorization and operational workflows.

## Risks and Open Questions

1. **os-core exact APIs need verification.** This document assumes theme/layout primitives exist, but implementation must inspect installed package exports.
2. **Auth in browser is still dev/OIDC foundation, not full OAuth UI.** The dashboard should start with current cookie/header/bearer development constraints and evolve with Phase 2 auth completion.
3. **Deployment DTO encodes `manifestJson` and `validationJson` as strings.** UI must parse safely and handle malformed JSON.
4. **Agent grants and keys are schema-present but not API-complete.** UI should avoid promising full automation setup.
5. **Quota APIs are missing.** Usage page should start with runtime counters and a clear pending panel.
6. **Large validation reports may be unwieldy.** `ValidationReportPanel` should support collapsible sections.
7. **Runtime status is partly in-memory.** The backend reconciles stale statuses on startup; UI should explain stopped-after-restart behavior when appropriate.

## Intern Implementation Checklist

Before opening a PR:

- [ ] Every new component/page has a Storybook story.
- [ ] Page stories use MSW or a fake store where relevant.
- [ ] Components are split into atoms/molecules/organisms/pages.
- [ ] No component uses hardcoded colors when an os-core token exists.
- [ ] RTK Query endpoints match the Go API routes listed above.
- [ ] Upload flow handles rejected validation reports as a normal result, not a crash.
- [ ] Activation, rollback, and revoke require confirmation.
- [ ] Empty/loading/error/permission states are visible in Storybook.
- [ ] `pnpm build` passes.
- [ ] `pnpm storybook:build` passes.
- [ ] Go server still serves `/api/*` as API and `/app/*` as SPA.

## File References

Existing backend files that define dashboard contracts:

- `internal/httpapi/handler.go` — route registration and `/app` mount.
- `internal/webadmin/handler.go` — current placeholder handler to replace with embedded SPA serving.
- `internal/httpapi/api.go` — session/org/site HTTP handlers.
- `internal/httpapi/deployments.go` — deployment upload/list/detail/activate/rollback handlers.
- `internal/httpapi/runtime.go` — runtime status handlers.
- `internal/httpapi/agents_audit.go` — agent and audit handlers.
- `internal/control/deployments.go` — deployment validation/activation service behavior.
- `internal/control/agents.go` — agent/audit service behavior.
- `internal/store/migrations/001_initial_schema.sql` — org/site/deployment/agent/audit schema.
- `internal/store/migrations/003_runtime_status.sql` — runtime status/counter schema.
- `cmd/go-go-host/cmds/deployments.go` — CLI workflow names that should align with dashboard terminology.
- `cmd/go-go-host/doc/deploy-workflow.md` — user-facing deployment workflow help.
- `cmd/go-go-host/doc/rollback-workflow.md` — user-facing rollback workflow help.

## Handoff Summary

Phase 7 should produce a Storybook-first dashboard implementation. The app should be useful even before every backend feature is complete, because Storybook/MSW lets us model future and edge states safely. The first production-ready slice should be the shell, session bootstrap, org selection, site list, site detail, deployment upload/list/detail, and runtime status. Agents, audit, usage, and members can follow immediately after because the design and component slots are already defined.
