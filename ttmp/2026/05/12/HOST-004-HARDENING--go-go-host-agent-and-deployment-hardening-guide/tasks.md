---
Title: Tasks
Ticket: HOST-004-HARDENING
Status: active
Topics:
  - go-go-host
  - security
  - agents
  - deployments
  - hardening
DocType: tasks
Intent: execution
LastUpdated: 2026-05-12
---

# Tasks

## Phase 1: Visibility and key revocation

- [ ] Add `fingerprint` and `last_used_at` support for agent keys.
- [ ] Add `GET /api/v1/orgs/{org_id}/agents/{agent_id}/keys`.
- [ ] Add `POST /api/v1/orgs/{org_id}/agents/{agent_id}/keys/{key_id}/revoke`.
- [ ] Add audit event `agent.key.revoke` with reason and fingerprint metadata.
- [ ] Add dashboard key inventory UI.
- [ ] Add tests for revoked-key signed request denial.

## Phase 2: Safer activation grants

- [ ] Add dashboard grant editor for `canActivate`.
- [ ] Add danger confirmation when enabling `canActivate`.
- [ ] Enforce org-owner-only updates for `canActivate` in backend.
- [ ] Add grant before/after audit metadata.
- [ ] Add Storybook stories for safe deploy-only and dangerous auto-activate grants.

## Phase 3: Security observability

- [ ] Add stable error codes to agent/signature/deploy-run failures.
- [ ] Audit signature failure classes: bad signature, timestamp skew, replay, revoked key, grant denied.
- [ ] Avoid logging raw secrets, signatures, tokens, or private key material.
- [ ] Add dashboard security event filters or page.
- [ ] Add CLI troubleshooting output keyed by error code.

## Phase 4: Upload and deploy-run robustness

- [ ] Add deploy-run statuses `uploading`, `expired`, and `cancelled`.
- [ ] Add atomic `pending -> uploading` transition before multipart parsing/validation.
- [ ] Reject second upload attempt for the same deploy run.
- [ ] Add deploy-run cancel endpoint.
- [ ] Add deploy-run expiry cleanup job or maintenance command.

## Phase 5: Artifact integrity and provenance

- [ ] Store bundle SHA256 on deployments.
- [ ] Include bundle hash in deployment DTOs, CLI output, dashboard detail pages, and audit metadata.
- [ ] Add optional provenance metadata: git SHA, git ref, CI run URL, builder.
- [ ] Add duplicate bundle hash detection or warning.

## Phase 6: Runtime safety after auto-activation

- [ ] Add post-activation health window for auto-activated deployments.
- [ ] Record runtime event or audit event for post-activation health failure.
- [ ] Add dashboard warning for auto-activated deployment health failures.
- [ ] Design rollback policy, but do not auto-rollback until semantics are reviewed.
