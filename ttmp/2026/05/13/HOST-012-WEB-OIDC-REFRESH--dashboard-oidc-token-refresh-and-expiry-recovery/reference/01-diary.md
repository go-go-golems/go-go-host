---
Title: Diary
Ticket: HOST-012-WEB-OIDC-REFRESH
Status: active
Topics:
  - go-go-host
  - frontend
  - oauth
  - oidc
  - keycloak
  - security
DocType: reference
Intent: long-term
Summary: Chronological implementation diary for dashboard OIDC refresh and expiry recovery.
---

# Diary

## Goal

Capture the investigation, implementation, validation, and deployment notes for fixing dashboard OIDC expiry by adding refresh-token use and request retry behavior to the web UI.

## Step 1: Create ticket, explain the issue, and define implementation tasks

The issue is that the dashboard stores OIDC expiry and refresh-token data, but the API transport keeps sending the originally stored bearer token until the backend rejects it. This step created a dedicated ticket and design guide so the fix is explicit: token refresh belongs in the shared frontend auth/API layer, not in individual pages.

The design mirrors the HOST-011 CLI refresh work. The browser dashboard should refresh access tokens shortly before expiry, retry once after `401`, clear local tokens when refresh fails, and apply the same behavior to both RTK Query endpoints and the custom deployment upload path.

### Prompt Context

**User prompt (verbatim):** "ok, create a ticket + implementation guide (that also explains the issue), and then add tasks and implement it, committing at appropriate intervals, keeping a diary."

**Assistant interpretation:** Create a docmgr ticket for the dashboard OIDC expiry problem, write a clear implementation guide, break the work into tasks, then implement with diary/changelog updates and commits.

**Inferred user intent:** Turn a recurring dashboard logout/expiry annoyance into a tracked, documented, implemented feature rather than an ad hoc patch.

### What I did
- Created ticket `HOST-012-WEB-OIDC-REFRESH`.
- Created design doc `design-doc/01-dashboard-oidc-refresh-issue-and-implementation-guide.md`.
- Created diary `reference/01-diary.md`.
- Replaced the generated task list with detailed implementation tasks.
- Inspected current files:
  - `web/admin/src/auth/oidc.ts`
  - `web/admin/src/services/goGoHostApi.ts`
  - `web/admin/src/services/types.ts`

### Why
- The frontend had all the evidence needed to explain the bug: it stored `expiresAt` and `refreshToken`, but `bearerToken()` ignored expiry and returned the old token.
- A ticket-level guide makes the intended refresh/retry semantics reviewable before code changes.

### What worked
- `docmgr ticket create-ticket` and `docmgr doc add` created the workspace and documents.
- The design doc now records the current issue, desired behavior, alternatives, and implementation plan.

### What didn't work
- N/A.

### What I learned
- The dashboard already has enough token metadata to implement refresh without backend changes.
- The custom upload `fetch()` path is a separate auth path and must be handled explicitly.

### What was tricky to build
- The key design constraint is avoiding recursion: RTK Query's `/config` request cannot depend on the same refresh flow that may need `/config` to know issuer/client metadata.

### What warrants a second pair of eyes
- Confirm the retry-on-401 behavior should clear tokens immediately if refresh fails, rather than redirecting directly from the API layer.

### What should be done in the future
- After implementation, consider a cross-tab token refresh lock if multiple dashboard tabs become common.

### Code review instructions
- Start with the design doc for intended behavior.
- Then review `web/admin/src/auth/oidc.ts` and `web/admin/src/services/goGoHostApi.ts` together.

### Technical details
- Ticket path: `ttmp/2026/05/13/HOST-012-WEB-OIDC-REFRESH--dashboard-oidc-token-refresh-and-expiry-recovery/`
