---
Title: Dashboard OIDC Refresh Issue and Implementation Guide
Ticket: HOST-012-WEB-OIDC-REFRESH
Status: active
Topics:
  - go-go-host
  - frontend
  - oauth
  - oidc
  - keycloak
  - security
DocType: design-doc
Intent: long-term
RelatedFiles:
  - Path: web/admin/src/auth/oidc.ts
    Note: Browser OIDC login, token storage, logout, and target refresh helpers.
  - Path: web/admin/src/services/goGoHostApi.ts
    Note: RTK Query API transport that currently attaches the stored bearer token synchronously.
  - Path: web/admin/src/services/types.ts
    Note: Public config and OIDC config TypeScript contracts.
  - Path: cmd/go-go-host/cmds/oidc_device.go
    Note: CLI refresh implementation used as the equivalent pattern for browser refresh.
  - Path: internal/httpapi/oidc.go
    Note: Backend OIDC verifier that rejects expired bearer tokens.
ExternalSources:
  - https://datatracker.ietf.org/doc/html/rfc6749
  - https://www.keycloak.org/securing-apps/oidc-layers
Summary: Design and implementation guide for fixing dashboard OIDC expiry by adding browser refresh-token use and 401 recovery.
LastUpdated: 2026-05-13T00:00:00-04:00
WhatFor: Explain why the dashboard currently expires and guide implementation of frontend token refresh.
WhenToUse: Use before changing web dashboard OIDC token storage, RTK Query auth headers, upload auth, or logout behavior.
---

# Dashboard OIDC Refresh Issue and Implementation Guide

## Executive summary

The dashboard currently stores OIDC token expiry and may store a refresh token, but it does not use that refresh token before API calls. The API transport reads the stored access token or ID token synchronously and sends it as `Authorization: Bearer ...`. Once that token expires, the backend correctly rejects it. The user experiences this as an OIDC expiry or a broken dashboard session even though the browser may still hold a valid refresh token.

HOST-012 fixes that gap in the web UI. The target design adds an async refresh path in `web/admin/src/auth/oidc.ts`, wraps RTK Query's base query so requests use a valid access token, retries once after `401` by forcing refresh, and updates the custom upload `fetch()` path to share the same token selection logic.

This is the browser equivalent of the refresh-aware CLI work completed in HOST-011. The CLI now refreshes access tokens in shared config resolution. The web UI needs the same responsibility in its shared API transport layer.

## Problem statement

The current browser OIDC implementation is split across two files.

`web/admin/src/auth/oidc.ts` stores token data after the authorization-code-with-PKCE callback:

```ts
localStorage.setItem(tokenStorageKey, JSON.stringify({
  idToken: tokens.id_token,
  accessToken: tokens.access_token,
  refreshToken: tokens.refresh_token,
  expiresAt: tokens.expires_in ? Date.now() + tokens.expires_in * 1000 : undefined,
}));
```

`web/admin/src/services/goGoHostApi.ts` attaches a bearer token to API requests:

```ts
prepareHeaders: (headers) => {
  const token = bearerToken();
  if (token) headers.set('Authorization', `Bearer ${token}`);
  return headers;
}
```

The problem is that `bearerToken()` only returns the stored token:

```ts
export function bearerToken(): string | undefined {
  const tokens = getStoredTokens();
  return tokens?.accessToken || tokens?.idToken;
}
```

It does not inspect `expiresAt`. It does not call Keycloak's token endpoint with `grant_type=refresh_token`. It does not clear the session when refresh fails. It does not retry a request that receives `401` because the token expired between selection and backend verification.

The backend behavior is correct. An expired OIDC token should be rejected. The frontend behavior is incomplete because it keeps sending a known-expired bearer token.

## Desired behavior

The dashboard should behave like a normal OIDC browser client:

1. Store the refresh token returned by the authorization-code exchange when Keycloak provides one.
2. Before protected API requests, check whether the stored access token expires soon.
3. If the token is still valid, send it without network refresh.
4. If the token expires soon and a refresh token exists, use Keycloak's token endpoint to refresh.
5. Store the refreshed token response, preserving the previous refresh token if Keycloak does not rotate it.
6. If an API request still returns `401`, force one refresh and retry once.
7. If refresh fails, clear local tokens and let the app return to the login path.
8. Use the same logic for the custom deployment-upload `fetch()` path, not only RTK Query's normal endpoints.

The refresh threshold should be short and conservative. A 60-second early refresh window is sufficient for beta and matches the CLI implementation.

## Proposed design

### Token shape

Extend `StoredOIDCTokens` with issuer/client metadata:

```ts
export interface StoredOIDCTokens {
  issuer?: string;
  clientId?: string;
  scopes?: string[];
  idToken: string;
  accessToken?: string;
  refreshToken?: string;
  expiresAt?: number;
}
```

New logins should store `issuer`, `clientId`, and `scopes` from the active config. Existing sessions may not have those fields, so refresh helpers should accept the current `OIDCConfig` as input and use it as the source of truth when available.

### Refresh helper

Add a public helper:

```ts
export async function getValidBearerToken(config?: ConfigResponse, options?: { forceRefresh?: boolean }): Promise<string | undefined>
```

The helper should:

- return `undefined` when no tokens exist,
- return the access token or ID token when it is not close to expiry and refresh is not forced,
- refresh when `forceRefresh` is true,
- refresh when `expiresAt` is within 60 seconds,
- clear tokens and return `undefined` when refresh is impossible or rejected.

Use OIDC discovery to find `token_endpoint`, then post:

```text
grant_type=refresh_token
client_id=<dashboard client id>
refresh_token=<refresh token>
```

The dashboard client is public, so no client secret is sent.

### RTK Query transport

Replace the synchronous `prepareHeaders` approach with a custom async base query:

1. If the request is `/config`, send it without auth.
2. Otherwise obtain public config from a small cached `fetch('/api/v1/config')` helper when refresh might be needed.
3. Call `getValidBearerToken(config)` and attach the token.
4. Run the actual request with `fetchBaseQuery`.
5. If the result is `401`, call `getValidBearerToken(config, { forceRefresh: true })` and retry once.
6. If no token is available after refresh failure, return the original `401` and leave token state cleared.

The config fetch must not recurse through RTK Query's own base query. Use a direct `fetch()` call for this internal bootstrap.

### Upload transport

`uploadDeployment` uses a custom `fetch()` because it posts `FormData`. It must also call `getValidBearerToken()` before sending the upload, and it should retry once after `401` by forcing refresh. Otherwise uploads remain vulnerable to the same expiry issue after normal API endpoints are fixed.

## Alternatives considered

### Always redirect to Keycloak on expiry

This would be simpler, but it would discard valid refresh tokens and interrupt the user unnecessarily. Refresh tokens exist specifically to avoid repeating interactive login during a valid session.

### Refresh on a fixed timer

A timer-based approach can work, but it introduces lifecycle concerns around hidden tabs, multiple dashboard tabs, sleep/wake behavior, and stale timers. Request-time refresh is simpler: the browser refreshes only when it needs a token for an API call.

### Let the backend refresh tokens

The backend does not hold the browser's refresh token, and adding a backend token session layer would change the architecture substantially. The current app is a public browser OIDC client; the browser owns its token storage and sends bearer tokens to the API.

### Use only ID tokens

The current `bearerToken()` falls back to ID token when no access token exists. For API authorization, access tokens are preferred. The fallback should remain for compatibility with existing behavior, but refresh should prefer access tokens and should store new access tokens from Keycloak.

## Implementation plan

1. Create ticket documentation and task list.
2. Add token metadata, refresh-token exchange, refresh locking, token storage helpers, and `getValidBearerToken()` in `web/admin/src/auth/oidc.ts`.
3. Replace RTK Query's synchronous `prepareHeaders` transport with an async base query wrapper in `web/admin/src/services/goGoHostApi.ts`.
4. Update the custom upload `fetch()` path to use the same valid-token helper and one-time 401 retry behavior.
5. Add Vitest coverage for token refresh and base query retry semantics.
6. Run `pnpm exec vitest`, `pnpm build`, and relevant Go tests if embedded assets or Go contracts are touched.
7. Update diary/changelog and commit.

## Review checklist

- Expired access tokens are not sent when a refresh token is available.
- Refresh failure clears local OIDC tokens.
- `401` responses trigger at most one forced refresh and retry.
- `/api/v1/config` is not blocked by auth refresh recursion.
- Uploads use the same auth behavior as normal JSON endpoints.
- Existing dev-auth and logged-out flows continue to work.
- New tests do not depend on real Keycloak.
