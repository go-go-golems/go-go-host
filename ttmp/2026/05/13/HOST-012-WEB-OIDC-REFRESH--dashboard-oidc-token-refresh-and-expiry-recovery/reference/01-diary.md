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

---

## Step 2: Implement dashboard token refresh and API retry behavior

This step added the missing browser refresh mechanism. The dashboard now treats stored OIDC tokens as expiring credentials rather than permanent bearer strings. When a token is close to expiry, the frontend uses the stored refresh token and Keycloak discovery metadata to request a fresh access token before sending protected API calls.

The API transport now has the same responsibility for all dashboard pages: attach a valid bearer token, retry once after `401` with a forced refresh, and clear tokens if refresh cannot be performed. The custom deployment upload path was updated separately because it uses `fetch()` directly instead of the normal RTK Query JSON transport.

### Prompt Context

**User prompt (verbatim):** (same as Step 1)

**Assistant interpretation:** Implement the planned HOST-012 frontend refresh behavior, add tests, validate, and record the work.

**Inferred user intent:** Stop recurring web dashboard OIDC expiry during normal use without forcing full re-login whenever a refresh token is still valid.

### What I did
- Extended `StoredOIDCTokens` with `issuer`, `clientId`, and `scopes` metadata for new browser sessions.
- Added refresh helpers in `web/admin/src/auth/oidc.ts`:
  - `getValidBearerToken(config, options)`
  - `refreshStoredTokens(config, previous)`
  - refresh-token grant POST to Keycloak token endpoint
  - 60-second early refresh window
  - refresh-token preservation when Keycloak does not rotate it
  - token cleanup on refresh failure
- Replaced synchronous RTK Query `prepareHeaders` with an async `baseQuery` wrapper in `web/admin/src/services/goGoHostApi.ts`.
- Added cached direct fetch of `/api/v1/config` for refresh metadata without recursing through RTK Query's own `/config` endpoint.
- Added one-time forced refresh and retry after `401`.
- Updated the deployment upload `fetch()` path to use the same valid-token logic and one retry after `401`.
- Added Vitest tests:
  - `web/admin/src/auth/oidc.test.ts`
  - `web/admin/src/services/goGoHostApi.test.ts`

### Why
- The backend should reject expired tokens. The bug was that the dashboard kept sending them even when a refresh token was available.
- Centralizing refresh in the API transport prevents every page and mutation from having to understand OIDC expiry.
- Uploads needed explicit handling because they bypass `fetchBaseQuery`.

### What worked
- `pnpm exec vitest run src/auth/oidc.test.ts src/services/goGoHostApi.test.ts` passed.
- `pnpm build` passed.
- `go test ./internal/httpapi ./internal/config -count=1` passed.

### What didn't work
- The first RTK Query test failed in Node because relative base URLs such as `/api/v1/me` are invalid for Node's `Request` constructor. I fixed this by making the base URL absolute through `apiOrigin()`, which uses `window.location.origin` in the browser and `http://localhost` in tests.
- The test mock initially treated `fetch` input as a string, but `fetchBaseQuery` passes a `Request` object. I updated the test to read `input.url` and `input.headers` when the input is a `Request`.

### What I learned
- The production browser behavior is unchanged by the absolute base URL because `window.location.origin` resolves to the same host that served the dashboard.
- Testing RTK Query transport in Node requires handling `Request` objects rather than only URL strings.

### What was tricky to build
- The refresh path needs public config metadata, but `/api/v1/config` is itself an RTK Query endpoint. Using the RTK base query for this internal fetch could recurse into the refresh mechanism. The implementation uses a small direct `fetch()` with a cached promise to avoid recursion.
- `401` retry must be bounded. The implementation retries once after forced refresh and then returns the result.
- Refresh failure should not leave a known-bad token in local storage. The refresh helper clears tokens when refresh is impossible or rejected.

### What warrants a second pair of eyes
- Review whether clearing tokens on missing refresh token is the preferred UX, or whether the app should preserve expired tokens until explicit login redirect.
- Review whether future multi-tab behavior needs a cross-tab lock using `BroadcastChannel` or the Web Locks API.
- Review whether upload retry with the same `FormData` is acceptable across all target browsers. Current browser FormData objects are reusable, but file streams in other environments can have stricter semantics.

### What should be done in the future
- Add an end-to-end browser smoke that forces a short token lifetime and verifies the dashboard refreshes without visible interruption.
- Consider surfacing an unobtrusive "session expired, sign in again" UI when refresh fails.

### Code review instructions
- Start with `web/admin/src/auth/oidc.ts` for refresh semantics.
- Then review `web/admin/src/services/goGoHostApi.ts` for request-time refresh and retry behavior.
- Review tests in `web/admin/src/auth/oidc.test.ts` and `web/admin/src/services/goGoHostApi.test.ts`.
- Validate with:
  - `cd web/admin && pnpm exec vitest run src/auth/oidc.test.ts src/services/goGoHostApi.test.ts`
  - `cd web/admin && pnpm build`

### Technical details
- Refresh threshold: 60 seconds before `expiresAt`.
- Refresh grant body: `grant_type=refresh_token`, `client_id=<dashboard client>`, `refresh_token=<stored refresh token>`.
- Retry policy: at most one forced refresh after a `401`.
