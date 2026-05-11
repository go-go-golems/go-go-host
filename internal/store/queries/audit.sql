-- name: InsertAuditEvent :one
INSERT INTO audit_log (id, org_id, actor_type, actor_id, action, resource_type, resource_id, ip_address, user_agent, metadata_json, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING id, org_id, actor_type, actor_id, action, resource_type, resource_id, ip_address, user_agent, metadata_json, created_at;

-- name: ListAuditEventsForOrg :many
SELECT id, org_id, actor_type, actor_id, action, resource_type, resource_id, ip_address, user_agent, metadata_json, created_at
FROM audit_log
WHERE org_id = $1
ORDER BY created_at DESC
LIMIT $2;

-- name: ListAuditEventsFiltered :many
SELECT id, org_id, actor_type, actor_id, action, resource_type, resource_id, ip_address, user_agent, metadata_json, created_at
FROM audit_log
WHERE org_id = $1
  AND ($2 = '' OR resource_id = $2)
  AND ($3 = '' OR actor_type = $3)
  AND ($4 = '' OR actor_id = $4)
  AND ($5 = '' OR action = $5)
ORDER BY created_at DESC
LIMIT $6;
