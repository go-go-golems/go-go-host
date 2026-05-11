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
