---
Title: Tasks
Ticket: HOST-012-WEB-OIDC-REFRESH
Status: active
Topics:
  - go-go-host
  - frontend
  - oauth
  - oidc
  - keycloak
  - security
DocType: tasks
---

# Tasks

- [x] 1. Document the dashboard OIDC expiry issue, current token flow, and desired refresh/retry behavior.
- [x] 2. Implement browser-side token metadata and refresh helpers in `web/admin/src/auth/oidc.ts`.
- [x] 3. Wrap RTK Query baseQuery so normal API requests refresh before expiry and retry once after 401.
- [x] 4. Update custom upload fetch auth to use the same async valid-token path.
- [x] 5. Add focused frontend tests for refresh behavior, 401 retry behavior, and refresh failure cleanup.
- [x] 6. Run frontend build/tests and relevant Go tests.
- [ ] 7. Update diary/changelog, relate files, and commit implementation at appropriate milestones.
