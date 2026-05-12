---
Title: Investigation diary
Ticket: HOST-006-PRODUCTION-READINESS
Status: active
Topics:
    - go-go-host
    - hosting
    - security
    - deployments
    - agents
    - platform-admin
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: "Chronological diary for the HOST-006 production readiness investigation."
LastUpdated: 2026-05-12T13:14:11.750159577-04:00
WhatFor: "Use this to resume or audit the production readiness planning work."
WhenToUse: "When implementing beta launch tasks or reviewing evidence behind the readiness guide."
---

# Diary

## Goal

Capture the investigation and documentation work for the go-go-host production readiness and beta launch plan.

## Step 1: Ticket setup, evidence gathering, and production-readiness guide

I created a new production-readiness ticket and wrote the first full design/implementation guide for moving go-go-host from a local MVP to a beta service. The guide is intentionally ordered by necessity: real browser auth and admin bootstrap first, then release/deploy mechanics, then domains/TLS, runtime isolation, secrets, backups, observability, and beta onboarding.

The main finding is that the core product loop exists, but the service is not production-ready until local/staging exercise the same OIDC path beta users will use. The current dev-user flow is useful for fast backend iteration, but it should not remain the normal browser testing path for beta readiness.

### Prompt Context

**User prompt (verbatim):** "what are all th e missing pieces needed for proper deploy (even just for beta users)?"

**User prompt (verbatim):** "Create a new ticket with a full detail of all the things needed to push into production, in order of necessity.

Create  a detailed analysis / design / implementation guide that is for a new intern, explaining all the parts of the system needed to understand what it is, with prose paragraphs and bullet point sand pseudocode and diagrams and api references and file references. It should be very clear and technical. Store in the ticket and the nupload to remarkable."

**Assistant interpretation:** Create a new docmgr ticket and write a clear intern-facing production-readiness implementation guide, grounded in current code evidence and prioritized by launch necessity, then upload the result to reMarkable.

**Inferred user intent:** Turn the informal question “what is missing for beta/production?” into a durable engineering plan that can drive future implementation work.

**Commit (code):** N/A — documentation/ticket work only.

### What I did

- Loaded the ticket-research, docmgr, reMarkable upload, diary, and textbook-writing instructions.
- Created ticket `HOST-006-PRODUCTION-READINESS` with topics `go-go-host,hosting,security,deployments,agents,platform-admin`.
- Created the primary design doc and investigation diary.
- Gathered evidence from:
  - auth middleware and OIDC verifier,
  - frontend RTK Query API setup,
  - AppShell dev-auth display,
  - dev Docker Compose and devctl plugin,
  - runtime and supervisor,
  - site settings/domain/secrets placeholder APIs,
  - maintenance/export/prune APIs,
  - Dockerfile and production config example,
  - bundled developer/JS/agent docs,
  - gated Playwright E2E script.
- Wrote the production-readiness guide with priority table, architecture diagrams, current-state analysis, gap analysis, implementation phases, API sketches, pseudocode, tests, risks, and file references.
- Rewrote `tasks.md` into a production-readiness backlog.

### Why

- The user asked for a new ticket and a detailed intern-facing design/implementation guide.
- Production readiness spans many subsystems, so the output needed to explain the architecture before prescribing work.
- The order of work matters: browser OIDC and admin bootstrap must come before domains, secrets, or beta onboarding.

### What worked

- The current codebase has enough structure to make the plan concrete: OIDC bearer validation exists, Keycloak Compose services exist, dashboard API access is centralized, and maintenance/runtime/domain APIs are already factored.
- The new developer and agent help docs provide a useful evidence point: user-facing docs are now comparatively strong, while auth/domain/secrets/ops remain the main production gaps.

### What didn't work

- `docmgr vocab list` showed no dedicated `operations` or `production` topic, so I stayed within existing vocabulary topics for the ticket.
- There are unrelated dirty files in `cmd/go-go-host*/cmds/support.go` and an untracked `HOST-005-E2E-FIXES` workspace. I did not touch them because this ticket is documentation-only.

### What I learned

- Keycloak is defined in Compose, but devctl does not start it and the Keycloak config file still has `devAuth: true`; this makes real local browser auth the top production-readiness gap.
- The backend OIDC verifier is already enough to anchor the browser-login implementation, but the frontend needs PKCE, token storage, logout, and Authorization header injection.
- Domain and secrets work are explicit placeholders rather than half-hidden missing features, which is good for planning because the gaps are visible and documentable.

### What was tricky to build

- The guide needed to distinguish “implemented but not production-grade” from “not implemented.” For example, maintenance export exists, but scheduled backup/restore drills do not; OIDC bearer verification exists, but browser login does not; domain rows exist, but DNS verification and TLS do not.
- The priority order needed to be strict. It would be tempting to implement custom domains or secrets first, but beta users cannot safely test anything until real auth and admin bootstrap are in place.

### What warrants a second pair of eyes

- Review the runtime-isolation recommendations, especially whether in-process Goja is acceptable for closed beta after adding interrupts/concurrency limits or whether subprocess isolation should be P0/P1.
- Review the platform-admin bootstrap recommendation: config-based OIDC subject/email/role mapping may be enough for beta, but production may need a stronger bootstrap ceremony.
- Review the TLS recommendation: edge-managed TLS is simpler for beta, but product goals may require daemon-managed ACME later.

### What should be done in the future

- Implement Phase 1 from the guide: Keycloak realm import, devctl OIDC profile, frontend PKCE, Authorization headers, OIDC admin bootstrap, and login E2E.
- Turn the guide's phases into separate implementation tickets if work will be distributed across multiple sessions or engineers.

### Code review instructions

- Start with `design-doc/01-go-go-host-production-readiness-and-beta-launch-implementation-guide.md`.
- Confirm each major claim maps to the referenced files.
- Review `tasks.md` for whether the priority order matches the desired beta launch sequence.
- Validate with `docmgr doctor --ticket HOST-006-PRODUCTION-READINESS --stale-after 30`.

### Technical details

Representative evidence commands:

```bash
nl -ba internal/httpapi/auth.go | sed -n '1,120p'
nl -ba internal/httpapi/oidc.go | sed -n '1,130p'
nl -ba web/admin/src/services/goGoHostApi.ts | sed -n '1,35p'
nl -ba deployments/dev/docker-compose.yaml | sed -n '1,90p'
nl -ba plugins/go-go-host-devctl.py | sed -n '1,150p'
nl -ba internal/runtime/runtime.go | sed -n '1,130p'
nl -ba internal/control/maintenance.go | sed -n '1,240p'
```
