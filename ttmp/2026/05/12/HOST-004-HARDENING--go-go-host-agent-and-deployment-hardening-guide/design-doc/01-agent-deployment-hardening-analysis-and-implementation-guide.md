---
Title: Agent Deployment Hardening Analysis and Implementation Guide
Ticket: HOST-004-HARDENING
Status: active
Topics:
    - go-go-host
    - security
    - agents
    - deployments
    - hardening
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: cmd/go-go-host-agent/cmds/deploy.go
      Note: agent CLI signed deploy and --activate workflow
    - Path: internal/control/agent_runs.go
      Note: agent enrollment
    - Path: internal/control/deployments.go
      Note: deployment upload/activation service logic and ActivateAsAgent path
    - Path: internal/deploy/bundle.go
      Note: bundle validation
    - Path: internal/httpapi/agent_signed_integration_test.go
      Note: signed agent security and auto-activation regression coverage
    - Path: internal/httpapi/agents_audit.go
      Note: agent creation
    - Path: internal/httpapi/deployments.go
      Note: deploy-run upload endpoint and optional auto-activation behavior
ExternalSources: []
Summary: Intern-ready guide to go-go-host agent deployment hardening, signed deploy runs, auto-activation risks, key rotation, upload safety, audit, and operational controls.
LastUpdated: 2026-05-12T10:45:12.574701914-04:00
WhatFor: Use this guide to understand the current go-go-host agent/deployment security model and implement the next hardening phase safely.
WhenToUse: Read before changing agent enrollment, signatures, deploy runs, upload tokens, activation, grants, audit, or deployment validation.
---


# Agent Deployment Hardening Analysis and Implementation Guide

## Executive Summary

`go-go-host` now has the core pieces of a hosting platform: users create organizations and sites, upload deployment bundles, validate those bundles, activate them into a Goja runtime, and inspect audit/runtime state through CLI and dashboards. The newest part of the system is the headless deployment-agent flow. An agent can be created by a human, enrolled with an Ed25519 key, and then use signed requests to create short-lived deploy runs. The deploy run gives the agent an upload token for exactly one deployment attempt. If the agent has an explicit `can_activate` grant and asks for activation, the upload endpoint can also promote the validated deployment to live traffic.

That is enough to make CI deployments possible, but it is not yet enough to call the system production-hardened. The hardening work is about reducing blast radius, increasing observability, and making dangerous actions understandable to operators. The important idea is not “add more security features everywhere.” The important idea is to preserve the control boundaries already present in the design:

- A human creates the agent and its grants.
- The agent owns a private key; the server stores only the public key.
- A signed deploy-run request authorizes intent.
- A short-lived upload token authorizes the large multipart upload.
- Bundle validation decides whether the uploaded code is acceptable.
- Activation swaps live traffic and must be audited distinctly.

The next hardening phase should make these boundaries visible and enforceable. The most useful immediate work is: a grant UI that makes `canActivate` hard to enable accidentally, key inventory and revoke/rotation, security-failure audit events, upload-token one-time semantics, bundle SHA256 storage, and cleanup jobs for nonces/deploy runs. These changes make the existing system easier to operate and safer without redesigning the platform.

## 1. The System in One Page

The easiest way to understand this area is to separate three concepts that are easy to confuse: **identity**, **authorization**, and **deployment state**.

An **agent identity** is a machine identity. It is represented by an `agents` row and one or more `agent_keys` rows. The agent proves identity by signing a canonical string with an Ed25519 private key. The server verifies the signature against the stored public key.

An **agent authorization** is a grant. It is represented by `agent_site_grants`. A grant says what the agent may do for a particular site: deploy, rollback, activate, use specific channels, and use specific paths. The grant is the blast-radius control.

A **deployment state transition** is what happens to hosted code. A bundle is uploaded, validated, recorded as a deployment, and optionally activated. Activation swaps runtime traffic. That action is more dangerous than upload, so it has a separate grant bit: `can_activate`.

```text
+-------------------+       creates        +-------------------+
| Human user/owner  | -------------------> | Agent record      |
| org_owner/dev     |                      | agents            |
+-------------------+                      +-------------------+
          |                                           |
          | creates grant                             | enrolls key
          v                                           v
+-------------------+                      +-------------------+
| Site grant        |                      | Agent key         |
| agent_site_grants |                      | agent_keys        |
+-------------------+                      +-------------------+
          |                                           |
          | authorizes                                | signs
          v                                           v
+-------------------+       upload token    +-------------------+
| Deploy run        | -------------------> | Bundle upload     |
| deploy_runs       |                      | deployment row    |
+-------------------+                      +-------------------+
                                                    |
                                                    | optional can_activate
                                                    v
                                         +----------------------+
                                         | Runtime activation   |
                                         | deployment.activate  |
                                         +----------------------+
```

The key rule: **the agent never gets human credentials**. It gets its own cryptographic identity and scoped grants.

## 2. Current Implementation Map

This section tells you where to look in the code. If you are new to this repository, do not start with the CLI. Start with the service layer, because it defines the invariants. Then read HTTP handlers and CLI commands as transport adapters.

| Area | Files | What to learn there |
|---|---|---|
| Agent service logic | `internal/control/agent_runs.go` | Enrollment token hashing, key registration, signature verification, nonce checks, grant checks, deploy-run creation. |
| Deployment service logic | `internal/control/deployments.go` | Upload, validation, activation, agent activation, runtime swap, audit. |
| Agent HTTP endpoints | `internal/httpapi/agents_audit.go` | Human agent APIs, enrollment endpoint, signed deploy-run endpoint. |
| Upload HTTP endpoint | `internal/httpapi/deployments.go` | Multipart upload, deploy-run upload token validation, optional auto-activation. |
| Routes | `internal/httpapi/handler.go` | Which endpoints are user-authenticated vs agent-signed. |
| Store wrappers | `internal/store/agents.go` | Store-level methods for agents, keys, grants, nonces, deploy runs. |
| SQL queries | `internal/store/queries/agents.sql` | sqlc source for agent/deploy-run operations. |
| Schema | `internal/store/migrations/*.sql` | Tables and migration history. |
| Bundle validation | `internal/deploy/bundle.go` | Archive path checks, manifest parsing, quota/capability validation. |
| Agent CLI | `cmd/go-go-host-agent/cmds/*.go` | Local key config, request signing, enroll/deploy commands. |
| Human CLI | `cmd/go-go-host/cmds/agents.go` | Human creation of agents and immediate grants. |
| Security tests | `internal/httpapi/agent_signed_integration_test.go` | Executable security contract for signed deploy runs. |

When implementing hardening, keep the layering intact:

```text
CLI / HTTP handler
    -> control service
        -> store/sqlc
        -> runtime/deploy service
```

Avoid putting security invariants only in CLI code or only in React UI. UI can warn, but the server must enforce.

## 3. Data Model Reference

### 3.1 `agents`

An agent is the machine identity container. It has a name, org, status, creator, created time, and last-seen time.

```sql
CREATE TABLE agents (
  id TEXT PRIMARY KEY,
  org_id TEXT NOT NULL REFERENCES orgs(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  status TEXT NOT NULL,
  created_by_user_id TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL,
  last_seen_at TIMESTAMPTZ
);
```

Important states:

- `active`: signed requests may be accepted if the key and grant also allow them.
- `revoked`: signed requests must be denied.

Hardening gap: `last_seen_at` is currently agent-level. We should also track last use per key.

### 3.2 `agent_keys`

An agent key stores the public key. The server never stores the private key.

```sql
CREATE TABLE agent_keys (
  id TEXT PRIMARY KEY,
  agent_id TEXT NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
  public_key TEXT NOT NULL,
  status TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  revoked_at TIMESTAMPTZ
);
```

Hardening gaps:

- Add `fingerprint` or compute it consistently in API/UI.
- Add `last_used_at`.
- Add explicit key revoke endpoint.
- Add key rotation flow.

### 3.3 `agent_site_grants`

A grant is where authorization lives. It limits an agent to a site, channel set, path set, expiry, and specific actions.

```sql
CREATE TABLE agent_site_grants (
  agent_id TEXT NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
  site_id TEXT NOT NULL REFERENCES sites(id) ON DELETE CASCADE,
  can_deploy BOOLEAN NOT NULL DEFAULT false,
  can_rollback BOOLEAN NOT NULL DEFAULT false,
  can_activate BOOLEAN NOT NULL DEFAULT false,
  allowed_channels TEXT[] NOT NULL DEFAULT '{}',
  allowed_paths TEXT[] NOT NULL DEFAULT '{}',
  expires_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  PRIMARY KEY (agent_id, site_id)
);
```

The most dangerous field is `can_activate`. It means the agent can promote code to live traffic. It should be visible in the dashboard with warning language, role restrictions, and audit metadata.

### 3.4 `agent_nonces`

A nonce prevents replay. If an attacker captures a signed request, they cannot submit it again because the `(agent_id, nonce)` pair already exists.

```sql
CREATE TABLE agent_nonces (
  agent_id TEXT NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
  nonce TEXT NOT NULL,
  seen_at TIMESTAMPTZ NOT NULL,
  PRIMARY KEY (agent_id, nonce)
);
```

Hardening gap: add cleanup. Nonces do not need to live forever. They need to outlive the accepted timestamp skew window plus a safety buffer.

### 3.5 `deploy_runs`

A deploy run is an authorization envelope for an upload attempt.

```sql
CREATE TABLE deploy_runs (
  id TEXT PRIMARY KEY,
  site_id TEXT NOT NULL REFERENCES sites(id) ON DELETE CASCADE,
  actor_type TEXT NOT NULL,
  actor_id TEXT NOT NULL,
  agent_id TEXT NOT NULL DEFAULT '',
  requested_by_user_id TEXT NOT NULL DEFAULT '',
  status TEXT NOT NULL,
  allowed_actions TEXT[] NOT NULL DEFAULT '{}',
  allowed_channels TEXT[] NOT NULL DEFAULT '{}',
  allowed_paths TEXT[] NOT NULL DEFAULT '{}',
  upload_token_hash TEXT NOT NULL DEFAULT '',
  expires_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  finished_at TIMESTAMPTZ
);
```

The deploy run intentionally stores a hash of the upload token. The raw upload token is returned once to the agent. This mirrors the enrollment-token pattern.

Hardening gaps:

- Add `uploading` status or a compare-and-swap transition to prevent double upload races.
- Add `expired` or `cancelled` status.
- Add cleanup jobs.
- Store request metadata: key ID, bundle hash, CI URL, git SHA.

## 4. API Reference

### 4.1 Create an agent and optional immediate grant

Human-authenticated endpoint:

```http
POST /api/v1/orgs/{org_id}/agents
Content-Type: application/json
X-Go-Go-Host-User: alice
```

Request:

```json
{
  "name": "ci-main",
  "siteId": "site_123",
  "allowedChannels": ["default"],
  "allowedPaths": ["**"],
  "canActivate": true
}
```

Response:

```json
{
  "agent": {
    "id": "agt_...",
    "orgId": "org_...",
    "name": "ci-main",
    "status": "active"
  },
  "enrollmentToken": "enroll_...",
  "grant": {
    "agentId": "agt_...",
    "siteId": "site_123",
    "canDeploy": true,
    "canActivate": true,
    "allowedChannels": ["default"],
    "allowedPaths": ["**"]
  }
}
```

Security note: the enrollment token is a one-time secret. UI should show it once, mask it by default, and never log it.

### 4.2 Enroll an agent key

Unauthenticated endpoint. The token is the credential.

```http
POST /api/v1/agent/enroll
Content-Type: application/json
```

Request:

```json
{
  "token": "enroll_...",
  "publicKey": "base64-ed25519-public-key"
}
```

Response:

```json
{
  "agent": { "id": "agt_...", "orgId": "org_...", "status": "active" },
  "keyId": "ak_...",
  "grant": { "siteId": "site_123", "canDeploy": true }
}
```

### 4.3 Create a signed deploy run

Agent-signed endpoint:

```http
POST /api/v1/agent/deploy-runs
Content-Type: application/json
X-Go-Go-Agent-ID: agt_...
X-Go-Go-Agent-Key-ID: ak_...
X-Go-Go-Agent-Timestamp: 2026-05-12T14:00:00Z
X-Go-Go-Agent-Nonce: random-nonce
X-Go-Go-Agent-Signature: base64-signature
```

Request:

```json
{
  "siteId": "site_123",
  "channel": "default",
  "path": "bundles/app.tar.gz",
  "action": "deploy",
  "activate": true
}
```

Response:

```json
{
  "id": "dr_...",
  "siteId": "site_123",
  "status": "pending",
  "uploadToken": "upload_...",
  "expiresAt": "2026-05-12T14:30:00Z",
  "allowedPaths": ["**"]
}
```

If `activate` is true but the grant does not have `can_activate`, this endpoint must fail before any upload token exists.

### 4.4 Upload a bundle for a deploy run

Upload-token endpoint:

```http
POST /api/v1/agent/deploy-runs/{run_id}/upload
Content-Type: multipart/form-data
X-Go-Go-Upload-Token: upload_...
```

Response without auto-activation:

```json
{
  "deployRunId": "dr_...",
  "activated": false,
  "deployment": { "id": "dep_...", "status": "validated", "createdByType": "agent" },
  "report": { "valid": true }
}
```

Response with auto-activation:

```json
{
  "deployRunId": "dr_...",
  "activated": true,
  "deployment": { "id": "dep_...", "status": "active", "createdByType": "agent" },
  "report": { "valid": true }
}
```

## 5. Signed Request Canonicalization

The signature contract is deliberately small. The CLI and server both build this string:

```text
METHOD\nPATH?QUERY\nSHA256_HEX_BODY\nRFC3339_TIMESTAMP\nNONCE
```

For example:

```text
POST
/api/v1/agent/deploy-runs
9b4f2f...e8
2026-05-12T14:00:00Z
9YfrpJr8OFkO5DcY
```

Pseudocode:

```text
function sign_request(method, url, body, private_key):
    timestamp = now_utc_rfc3339()
    nonce = random_128_bits_base64url()
    body_hash = sha256_hex(body)
    canonical = upper(method) + "\n" + url.request_uri + "\n" + body_hash + "\n" + timestamp + "\n" + nonce
    signature = ed25519_sign(private_key, canonical)
    set_headers(agent_id, key_id, timestamp, nonce, base64(signature))
```

Server verification:

```text
function verify_signed_request(request):
    parse timestamp
    reject if timestamp older/newer than allowed skew

    load agent by X-Go-Go-Agent-ID
    reject if agent is not active

    load key by X-Go-Go-Agent-Key-ID
    reject if key does not belong to agent or is not active

    body_hash = sha256_hex(exact_request_body_bytes)
    canonical = canonical_string(method, request_uri, body_hash, timestamp, nonce)
    reject if ed25519_verify(public_key, canonical, signature) fails

    insert (agent_id, nonce)
    reject if insert violates primary key

    mark agent last seen
    return agent
```

The body hash signs the exact bytes. Equivalent JSON with different whitespace signs differently. That is acceptable and simpler than canonical JSON, but it must be documented.

## 6. Deployment and Activation Flow

The current flow separates upload and activation. Auto-activation is an optional extension of the upload path.

```text
Agent CLI
  |
  | signed JSON: create deploy run { activate: true }
  v
HTTP /api/v1/agent/deploy-runs
  |
  | verifies signature, nonce, timestamp
  | checks grant: can_deploy, can_activate, channel, path, expiry
  | stores deploy_run allowed_actions=[deploy, activate]
  | returns upload token
  v
Agent CLI
  |
  | multipart upload + X-Go-Go-Upload-Token
  v
HTTP /api/v1/agent/deploy-runs/{id}/upload
  |
  | validates upload token and run status
  | stores and validates bundle
  | creates deployment created_by_type=agent
  | if report.valid and allowed_actions contains activate:
  |     ActivateAsAgent(agent_id, deployment_id)
  v
Runtime supervisor swaps traffic
```

The important invariant is that the upload endpoint does not accept `activate=true`. The decision was made earlier by the signed deploy-run endpoint and stored in `allowed_actions`.

## 7. Threat Model

A useful hardening guide needs to say what it is defending against. The table below is not exhaustive, but it covers the highest-value cases.

| Threat | Example | Current defense | Hardening still needed |
|---|---|---|---|
| Stolen enrollment token | Someone copies `enroll_...` from logs. | Token is hashed at rest and one-time use. | Token expiry UI, token revoke, one-active-token limit, redaction guidance. |
| Stolen agent private key | Attacker signs deploy runs as CI. | Agent/key status checks and site grants. | Per-key revoke, rotation, last-used, IP anomaly audit. |
| Replay attack | Captured signed request is submitted twice. | Nonce primary key and timestamp skew. | Nonce retention cleanup and security-failure audit. |
| Overbroad grant | Agent has `**` and `can_activate` for production. | Human must create grant. | Danger UX, owner-only changes, channel-specific activation. |
| Upload token leak | Raw `upload_...` is copied. | Token is short-lived and hashed at rest. | One-time/in-progress upload state and failure audit. |
| Archive bomb | Tiny compressed file expands massively. | File count/size quota validation. | Compression-ratio limits, per-file limit, HTTP-layer limit. |
| Runtime bad deploy | Code validates but fails under traffic. | Health check during activation. | Post-activation health window and auto-rollback policy. |
| Audit blind spot | Denied signatures leave no record. | Successful actions audited. | Security-failure audit events with stable error codes. |

## 8. Hardening Work Packages

The rest of this document turns the hardening ideas into implementation-sized packages. Each package has a goal, rationale, implementation sketch, API shape, files to touch, tests, and operator impact.

### Package A: Dashboard `canActivate` Grant UI

#### Goal

Make auto-activation permission visible and hard to enable by accident. `canActivate` grants should feel dangerous in the UI because they are dangerous: they allow an agent to promote code to live traffic.

#### Why this matters

The backend now enforces `can_activate`, but backend enforcement alone does not make an operator understand the consequence of the flag. A user can still set the flag through CLI/API without seeing the risk. The dashboard should explain what the permission means before it is granted.

#### Proposed UX

On an agent detail page or grant editor:

```text
[ ] Allow this agent to activate deployments automatically

Warning: this lets the agent move validated code to live traffic for this site.
Only enable this for trusted CI pipelines. The action will be audited as the agent.

Allowed channels: [ default v ]
Allowed paths:    [ **        ]
Expires at:       [ 2026-06-12 ]
```

When enabling `canActivate`, show a confirmation dialog:

```text
Title: Grant auto-activation permission?

This agent will be able to promote deployments to live traffic for site <site>.
If the agent key is stolen, an attacker may deploy and activate code within this grant.

Type ACTIVATE to confirm.
```

#### Server rules

UI warnings are not enough. The server should enforce:

- Only `org_owner` can set `canActivate`.
- `org_developer` may create deploy-only agents, but not activation-capable agents.
- Every grant update audit event should include before/after metadata.

#### Pseudocode

```text
function upsert_agent_grant(actor, org_id, grant_request):
    role = membership_role(actor.user_id, org_id)

    if grant_request.can_activate and role != org_owner:
        reject permission_denied

    old_grant = get_grant(agent_id, site_id)
    new_grant = upsert_grant(...)

    audit("agent.grant.upsert", metadata={
        before: old_grant,
        after: new_grant,
        changed_fields: diff(old_grant, new_grant)
    })
```

#### Files to touch

- `internal/control/agent_runs.go`
- `internal/httpapi/agents_audit.go`
- `internal/store/agents.go`
- `web/admin/src/pages/AgentsPage/` or a new `AgentDetailPage`
- `web/admin/src/components/.../AgentGrantEditor/`
- `web/admin/src/services/goGoHostApi.ts`
- `web/admin/src/services/types.ts`
- `web/admin/src/services/msw/fixtures.ts`
- `web/admin/src/services/msw/handlers.ts`

#### Tests

- Org developer cannot set `canActivate`.
- Org owner can set `canActivate`.
- Audit metadata records before/after.
- Storybook interaction story shows danger confirmation.

### Package B: Agent Key Inventory, Fingerprints, and Last-Used Tracking

#### Goal

Operators should know which keys exist, which are active, when they were created, when they were last used, and how to identify them without copying raw public keys.

#### Data model changes

Add columns:

```sql
ALTER TABLE agent_keys
  ADD COLUMN IF NOT EXISTS fingerprint TEXT NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS last_used_at TIMESTAMPTZ;
```

The fingerprint can be SHA256 of the decoded public key, encoded as a short display string:

```text
SHA256:3b:91:af:...:10
```

#### API shape

```http
GET /api/v1/orgs/{org_id}/agents/{agent_id}/keys
```

Response:

```json
[
  {
    "id": "ak_...",
    "agentId": "agt_...",
    "fingerprint": "SHA256:3b91af...",
    "status": "active",
    "createdAt": "2026-05-12T14:00:00Z",
    "lastUsedAt": "2026-05-12T14:30:00Z",
    "revokedAt": null
  }
]
```

#### Pseudocode

```text
function verify_signed_request(request):
    ... existing checks ...
    if signature_valid:
        touch_agent_last_seen(agent_id)
        touch_agent_key_last_used(key_id)
```

#### Files to touch

- `internal/store/migrations/006_agent_key_inventory.sql`
- `internal/store/queries/agents.sql`
- `internal/store/agents.go`
- `internal/control/agent_runs.go`
- `internal/httpapi/agents_audit.go`
- dashboard agent pages and stories

#### Tests

- Key list returns active key after enrollment.
- Signed request updates `last_used_at`.
- Fingerprint is stable and does not expose the raw private key.

### Package C: Key Rotation and Key Revoke

#### Goal

Allow an operator to replace a key without revoking the whole agent. This is necessary when a CI worker is replaced, a private key might have leaked, or a key is simply old.

#### Why this matters

Right now revocation is coarse. Revoking the whole agent works, but it disrupts every pipeline using that agent. Key-level revoke lets operators retire one credential while preserving the agent identity, grants, and history.

#### API shape

```http
POST /api/v1/orgs/{org_id}/agents/{agent_id}/keys/{key_id}/revoke
```

Request:

```json
{ "reason": "rotated after CI runner rebuild" }
```

Response:

```json
{
  "keyId": "ak_...",
  "status": "revoked",
  "revokedAt": "2026-05-12T15:00:00Z"
}
```

Key rotation can reuse enrollment tokens or introduce a narrower key-add token:

```http
POST /api/v1/orgs/{org_id}/agents/{agent_id}/key-rotation-tokens
```

A conservative v1 can use enrollment token semantics but bind the token to an existing agent.

#### Pseudocode

```text
function revoke_agent_key(actor, org_id, agent_id, key_id, reason):
    require org_owner or org_developer
    agent = get_agent(agent_id)
    require agent.org_id == org_id
    key = get_key(key_id)
    require key.agent_id == agent_id

    update key set status=revoked, revoked_at=now()
    audit("agent.key.revoke", metadata={ key_id, fingerprint, reason })
```

Verification already checks key status. The missing pieces are endpoint, UI, and tests.

#### Tests

- Active key can sign.
- Revoked key cannot sign.
- Replacement key can sign after old key is revoked.
- Revoke event appears in audit.

### Package D: Security Failure Audit Events

#### Goal

Successful security-sensitive actions are audited today. Failed attempts should also be visible, especially for signatures and grants.

#### Events to add

```text
agent.signature.invalid
agent.signature.timestamp_skew
agent.signature.nonce_replay
agent.signature.revoked_agent
agent.signature.revoked_key
agent.grant.denied
agent.upload_token.invalid
agent.deploy_run.expired
agent.deploy_run.wrong_site
agent.deploy_run.wrong_channel
agent.deploy_run.wrong_path
```

#### Design choice

Do not log secrets. Audit metadata may include:

- agent ID if known,
- key ID if known,
- site ID if requested,
- channel/path requested,
- reason code,
- request ID,
- remote IP,
- user agent.

Do not include:

- raw signature,
- raw upload token,
- raw enrollment token,
- private key material.

#### Pseudocode

```text
function deny_agent_request(reason, context):
    audit_security_event(reason, scrub(context))
    return error(code=reason, status=403 or 400)
```

#### Operator impact

A security-events page can show these separately from normal audit. Operators can distinguish “CI misconfigured” from “someone is replaying requests.”

### Package E: Stable Error Codes

#### Goal

Return machine-readable error codes so CLI, UI, and docs can explain failures consistently.

Current errors are mostly free-form strings. Add stable codes:

```json
{
  "error": "permission denied",
  "code": "agent_grant_denied",
  "requestId": "..."
}
```

Suggested codes:

| Code | Meaning |
|---|---|
| `agent_signature_invalid` | Signature verification failed. |
| `agent_timestamp_skew` | Timestamp too old or too far in future. |
| `agent_nonce_replay` | Nonce already used. |
| `agent_revoked` | Agent status is revoked. |
| `agent_key_revoked` | Key status is revoked. |
| `agent_grant_denied` | Site/channel/path/action not allowed. |
| `deploy_run_expired` | Deploy run expired before upload. |
| `upload_token_invalid` | Upload token missing or invalid. |
| `bundle_validation_failed` | Bundle failed validation. |

#### Files to touch

- `internal/httpapi/*.go`
- control error types
- CLI error formatting
- dashboard error callouts

### Package F: Upload Token One-Time and In-Progress Semantics

#### Goal

Prevent concurrent or repeated uploads to the same deploy run.

Current status values include pending/completed/rejected. A robust flow should transition atomically:

```text
pending -> uploading -> completed
pending -> uploading -> rejected
pending -> expired
pending -> cancelled
```

The important operation is compare-and-swap:

```sql
UPDATE deploy_runs
SET status = 'uploading'
WHERE id = $1
  AND status = 'pending'
  AND expires_at > now()
RETURNING *;
```

If no row is returned, the upload is denied. This prevents two simultaneous requests from both using the same upload token.

#### Pseudocode

```text
function begin_upload(run_id, token):
    run = get_run(run_id)
    reject if token hash mismatch
    updated = transition_pending_to_uploading(run_id)
    reject if updated row count == 0
    return run
```

### Package G: Bundle Hash and Provenance

#### Goal

Store what was uploaded in a stable, inspectable way. Operators should be able to answer: “Which exact archive is running?”

Add to deployments:

```sql
ALTER TABLE deployments
  ADD COLUMN IF NOT EXISTS bundle_sha256 TEXT NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS provenance_json JSONB NOT NULL DEFAULT '{}';
```

The agent CLI can send optional metadata:

```json
{
  "gitSha": "abc123",
  "gitRef": "main",
  "ciRunUrl": "https://github.com/org/repo/actions/runs/123",
  "builder": "github-actions"
}
```

Audit metadata should include `bundleSha256` and `deployRunId`.

### Package H: HTTP-Layer Upload Size Limits

#### Goal

Reject oversized uploads before writing the full body to disk.

Bundle validation already checks quota after the temp file exists. That is not enough to protect disk under malicious upload load. Wrap request body:

```go
r.Body = http.MaxBytesReader(w, r.Body, maxUploadBytes)
```

The limit should come from site quota when possible. For agent deploy-run upload, fetch the run, site, and quota before parsing multipart.

### Package I: Nonce and Deploy-Run Cleanup

#### Goal

Keep security tables from growing forever.

Suggested cleanup policy:

```text
agent_nonces: delete rows older than 24 hours
expired pending deploy_runs: mark expired after expires_at + 1 hour
deploy_runs: retain 90 days, or longer if linked to deployments/audit
```

Implementation options:

- daemon startup cleanup,
- periodic goroutine in `go-go-hostd`,
- admin CLI maintenance command,
- future background jobs table.

For v1, a startup cleanup and CLI maintenance command is simplest.

### Package J: Auto-Activation Health Window and Rollback Policy

#### Goal

A deployment can pass dry-run validation and still fail after live activation. Auto-activation raises the cost of this failure because no human is watching.

Add a post-activation health window:

```text
activate deployment
for 30 seconds:
    call runtime health endpoint every 5 seconds
    if failures exceed threshold:
        mark runtime/deployment failed
        optionally rollback if policy allows
```

Start with observe-only:

- Mark runtime failed.
- Emit audit/runtime event.
- Show dashboard alarm.

Do not auto-rollback until rollback semantics are carefully designed.

## 9. Recommended Immediate Additions to the Current Phase

The user asked which hardening items should be added to the current phase because they make the system easier and more robust immediately. These are the highest-value items to pull forward now.

### 9.1 Add `canActivate` grant UI with danger UX

This should be first because scoped auto-activation already exists. Without UI, operators can use CLI/API, but the permission is not visible where people inspect agents. It improves usability and safety at the same time.

Immediate scope:

- Show `canActivate` on agent tables/details.
- Add grant editor checkbox.
- Require confirmation when enabling it.
- Show warning text.
- If possible, enforce org-owner-only for `canActivate` in backend now.

### 9.2 Add key inventory and key revoke

Auto-activation makes key compromise more important. Operators need to see keys and revoke a key without deleting the whole agent.

Immediate scope:

- `GET /agents/{agent_id}/keys`.
- `POST /agents/{agent_id}/keys/{key_id}/revoke`.
- Dashboard key list.
- Test revoked key denial.

Rotation can come next, but revoke should come first.

### 9.3 Add security failure audit events for signed requests

This helps debugging immediately. It also turns security failures into visible signals.

Immediate scope:

- Audit nonce replay.
- Audit timestamp skew.
- Audit bad signature when agent/key is known.
- Audit grant denied with site/channel/path/action metadata.

### 9.4 Add upload token one-time transition

This closes a real race class. It is backend-only and does not need much UI.

Immediate scope:

- Add `uploading` deploy-run status.
- Atomic `pending -> uploading` query.
- Reject second upload.
- Tests for two upload attempts.

### 9.5 Add bundle SHA256 storage

This improves audit, debugging, and future provenance. It is low-risk and useful in dashboards/CLI.

Immediate scope:

- Compute SHA256 while copying upload to temp file or during validation.
- Store on deployment.
- Include in deployment DTO and audit metadata.

## 10. Suggested Task Breakdown

### Milestone 1: Visibility and revocation

- Add key list API.
- Add key fingerprint and last-used tracking.
- Add key revoke endpoint.
- Add dashboard key inventory.
- Add revoked-key signed request test.

### Milestone 2: Safer activation grants

- Add dashboard grant editor for `canActivate`.
- Add owner-only backend rule for `canActivate`.
- Add audit before/after metadata for grant changes.
- Add Storybook interactions for enabling/disabling `canActivate`.

### Milestone 3: Security observability

- Add structured error codes.
- Add security failure audit events.
- Add dashboard security-events view or filters.
- Add CLI troubleshooting output using error codes.

### Milestone 4: Upload/deploy-run robustness

- Add `uploading`, `expired`, and `cancelled` statuses.
- Add atomic begin-upload transition.
- Add upload token one-time-use tests.
- Add deploy-run expiry cleanup.

### Milestone 5: Artifact integrity

- Store bundle SHA256.
- Add provenance metadata fields.
- Include deploy-run ID and bundle hash in activation audit.
- Show bundle hash in deployment detail pages.

## 11. Intern Implementation Checklist

Before coding:

- Read `internal/control/agent_runs.go` from top to bottom.
- Read `internal/httpapi/agent_signed_integration_test.go` and explain each denial test in your own words.
- Run `go test ./...`.
- Run the live devctl smoke once if your environment supports it.

During coding:

- Put authorization checks in `internal/control`, not only in HTTP or UI.
- Add sqlc queries before store wrappers.
- Regenerate sqlc after query changes.
- Add audit events for both success and security-sensitive denial paths.
- Add at least one integration test for each new security invariant.

Before review:

- Run `sqlc generate`.
- Run `go test ./...`.
- Run `make web-build` if dashboard code changed.
- Run `make storybook-build` if stories changed.
- Update ticket diary/changelog/tasks.
- If changing agent deploy behavior, run live devctl smoke.

## 12. Open Questions

- Should `canActivate` be globally owner-only, or configurable per org?
- Should auto-activation be channel-specific with a dedicated `can_activate_channels` field?
- Should deploy-run request `path` represent a source artifact path, a deploy target path, or just an arbitrary CI label?
- Should upload tokens be accepted after validation failure for retry, or should retry always require a new deploy run?
- How long should nonces be retained?
- Should failed signature attempts be written to the main audit log or a separate security-events table?
- Should auto-activation ever auto-rollback on post-activation health failure?

## 13. Glossary

- **Agent**: A machine identity used by CI or automation to deploy without human credentials.
- **Agent key**: An Ed25519 public key registered for an agent. The private key stays on the agent machine.
- **Enrollment token**: A one-time secret created by a human to let an agent register a public key.
- **Grant**: A scoped authorization row connecting an agent to a site and allowed actions/channels/paths.
- **Deploy run**: A short-lived authorization envelope for one upload attempt.
- **Upload token**: A one-time-ish bearer secret returned after signed deploy-run creation and used for multipart upload.
- **Validation**: Bundle checks before a deployment can be activated.
- **Activation**: Runtime traffic swap to a deployment.
- **Auto-activation**: Agent-requested activation after upload/validation, allowed only by `can_activate`.
- **Nonce**: A unique value in a signed request that prevents replay.

## 14. References

Primary implementation files:

- `internal/control/agent_runs.go`
- `internal/control/deployments.go`
- `internal/httpapi/agents_audit.go`
- `internal/httpapi/deployments.go`
- `internal/httpapi/agent_signed_integration_test.go`
- `internal/store/queries/agents.sql`
- `internal/store/migrations/004_agent_enrollment_runs.sql`
- `internal/store/migrations/005_agent_auto_activate.sql`
- `internal/deploy/bundle.go`
- `cmd/go-go-host-agent/cmds/support.go`
- `cmd/go-go-host-agent/cmds/deploy.go`
- `cmd/go-go-host/cmds/agents.go`

Related ticket docs:

- `HOST-001-GO-GO-HOST-V1` implementation diary and tasks.
- `HOST-003-ADMIN-DASHBOARD` admin dashboard docs for operator UI patterns.
