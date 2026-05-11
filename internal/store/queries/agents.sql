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
