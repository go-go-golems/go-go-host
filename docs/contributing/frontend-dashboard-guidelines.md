# Frontend dashboard guidelines

This document defines how to work on the React dashboard in `web/admin`.

The dashboard has two surfaces:

- `/app/*` for organization users and developers.
- `/admin/*` for platform administrators.

The dashboard is a React/Vite/RTK Query/Storybook application using `@go-go-golems/os-core` and the repository's OS1 visual conventions.

## Required references

Read the OS1 playbook before visual or component work:

- [`playbooks/os1-admin-dashboard-ui-work-guidelines.md`](playbooks/os1-admin-dashboard-ui-work-guidelines.md)

Use this file for workflow and architecture rules; use the OS1 playbook for visual rules.

## Key files

| Area | Files |
|---|---|
| App routes | `web/admin/src/app/routes.tsx` |
| Store/provider setup | `web/admin/src/app/store.ts`, `web/admin/src/app/providers` |
| Auth helpers | `web/admin/src/auth/oidc.ts` |
| API client | `web/admin/src/services/goGoHostApi.ts` |
| API types | `web/admin/src/services/types.ts` |
| MSW fixtures/handlers | `web/admin/src/services/msw` |
| Pages | `web/admin/src/pages/*` |
| Components | `web/admin/src/components/*` |
| OS1 bridge CSS | `web/admin/src/app/macos1-bridge.css` |

## API state rules

Use RTK Query for backend API state.

When adding a backend-backed UI feature:

1. Add or update TypeScript types in `web/admin/src/services/types.ts`.
2. Add endpoint definitions in `web/admin/src/services/goGoHostApi.ts`.
3. Choose `providesTags` and `invalidatesTags` deliberately.
4. Add or update MSW fixtures and handlers.
5. Use generated hooks from `goGoHostApi` in pages/components.
6. Add stories for relevant states.

Do not add ad-hoc `fetch` calls in page components unless the API shape requires special handling, such as multipart upload. If special handling is needed, keep it centralized in the service layer.

## Page implementation rules

A page should handle these states explicitly:

- Loading.
- Error.
- Empty data.
- Populated data.
- Permission denied or disabled action when applicable.
- Mutation in progress.
- Mutation failure.

Use existing shared components for loading blocks, error callouts, tables, badges, panels, and buttons where possible. Do not create a new local pattern if a shared component already exists.

## Route rules

Routes live in `web/admin/src/app/routes.tsx`.

When adding a route:

- Put user/org/site pages under `/app`.
- Put platform-admin pages under `/admin`.
- Use existing route guards rather than duplicating auth checks.
- Add navigation only when the page is ready enough to be discoverable.
- Ensure direct navigation and refresh work through the embedded SPA fallback.

## Storybook and MSW

Storybook should exercise deterministic states. Add stories when creating a page or reusable component.

Recommended stories:

- Default/populated.
- Loading.
- Empty.
- Error.
- Permission denied or disabled controls when relevant.

Use MSW fixtures to keep stories close to real API shapes. If a story requires unrealistic mock data, revisit the API type or component design.

## Visual rules

Do not restyle the dashboard independently from the OS1 system. Follow the playbook for:

- Window and panel structure.
- Typography scale.
- Button, checkbox, input, select, and table styling.
- White dashboard workspace background.
- Compact spacing.
- Screenshot review.

Avoid:

- Generic SaaS cards.
- Large border radii.
- Unstructured whitespace.
- Browser-default checkboxes or inputs.
- One-off CSS scales.
- Duplicated copies of OS-core CSS.

## Validation

Run:

```bash
cd web/admin
pnpm build
pnpm storybook:build
```

For embedded dashboard changes, also run from the repository root:

```bash
go run ./cmd/build-web
go test ./internal/webadmin
go build ./...
```

## Review checklist

Before merging dashboard work, verify:

- Backend API state goes through RTK Query.
- Types, endpoints, MSW handlers, and stories are consistent.
- Loading, error, empty, and populated states are represented.
- Mutations invalidate the correct RTK Query tags.
- The page follows the OS1 playbook.
- Direct route refresh works when embedded in the daemon.
- The UI does not enforce a permission that the server fails to enforce.
