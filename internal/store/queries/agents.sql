-- name: CreateAgent :one
INSERT INTO agents (id, org_id, name, status, created_by_user_id, created_at, last_seen_at)
VALUES ($1, $2, $3, $4, $5, $6, NULL)
RETURNING id, org_id, name, status, created_by_user_id, created_at, last_seen_at;

-- name: GetAgent :one
SELECT id, org_id, name, status, created_by_user_id, created_at, last_seen_at
FROM agents
WHERE id = $1;

-- name: ListAgentsByOrg :many
SELECT id, org_id, name, status, created_by_user_id, created_at, last_seen_at
FROM agents
WHERE org_id = $1
ORDER BY created_at DESC;

-- name: UpdateAgentStatus :exec
UPDATE agents
SET status = $2
WHERE id = $1;

-- name: UpsertAgentSiteGrant :one
INSERT INTO agent_site_grants (agent_id, site_id, can_deploy, can_rollback, allowed_channels, allowed_paths, expires_at, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $8)
ON CONFLICT (agent_id, site_id)
DO UPDATE SET can_deploy = EXCLUDED.can_deploy, can_rollback = EXCLUDED.can_rollback, allowed_channels = EXCLUDED.allowed_channels, allowed_paths = EXCLUDED.allowed_paths, expires_at = EXCLUDED.expires_at, updated_at = EXCLUDED.updated_at
RETURNING agent_id, site_id, can_deploy, can_rollback, allowed_channels, allowed_paths, expires_at, created_at, updated_at;

-- name: ListAgentSiteGrants :many
SELECT agent_id, site_id, can_deploy, can_rollback, allowed_channels, allowed_paths, expires_at, created_at, updated_at
FROM agent_site_grants
WHERE agent_id = $1
ORDER BY site_id;

-- name: CreateAgentEnrollmentToken :exec
INSERT INTO agent_enrollment_tokens (token_hash, agent_id, org_id, status, expires_at, created_at)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetAgentEnrollmentToken :one
SELECT token_hash, agent_id, org_id, status, expires_at, created_at, used_at
FROM agent_enrollment_tokens
WHERE token_hash = $1;

-- name: MarkAgentEnrollmentTokenUsed :exec
UPDATE agent_enrollment_tokens
SET status = 'used', used_at = $2
WHERE token_hash = $1;

-- name: CreateAgentKey :one
INSERT INTO agent_keys (id, agent_id, public_key, status, created_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, agent_id, public_key, status, created_at, revoked_at;

-- name: GetAgentKey :one
SELECT id, agent_id, public_key, status, created_at, revoked_at
FROM agent_keys
WHERE id = $1;

-- name: ListAgentKeys :many
SELECT id, agent_id, public_key, status, created_at, revoked_at
FROM agent_keys
WHERE agent_id = $1
ORDER BY created_at DESC;

-- name: TouchAgentLastSeen :exec
UPDATE agents
SET last_seen_at = $2
WHERE id = $1;

-- name: InsertAgentNonce :exec
INSERT INTO agent_nonces (agent_id, nonce, seen_at)
VALUES ($1, $2, $3);

-- name: CreateDeployRun :one
INSERT INTO deploy_runs (id, site_id, actor_type, actor_id, agent_id, status, allowed_actions, allowed_channels, allowed_paths, upload_token_hash, expires_at, created_at)
VALUES ($1, $2, 'agent', $3, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING id, site_id, actor_type, actor_id, agent_id, requested_by_user_id, status, allowed_actions, allowed_channels, allowed_paths, upload_token_hash, expires_at, created_at, finished_at;

-- name: GetDeployRun :one
SELECT id, site_id, actor_type, actor_id, agent_id, requested_by_user_id, status, allowed_actions, allowed_channels, allowed_paths, upload_token_hash, expires_at, created_at, finished_at
FROM deploy_runs
WHERE id = $1;

-- name: FinishDeployRun :exec
UPDATE deploy_runs
SET status = $2, finished_at = $3
WHERE id = $1;
