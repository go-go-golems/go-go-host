-- name: InsertRuntimeEvent :one
INSERT INTO runtime_events (id, site_id, org_id, deployment_id, event_type, status, message, metadata_json, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id, site_id, org_id, deployment_id, event_type, status, message, metadata_json, created_at;

-- name: ListRuntimeEventsBySite :many
SELECT id, site_id, org_id, deployment_id, event_type, status, message, metadata_json, created_at
FROM runtime_events
WHERE site_id = $1
ORDER BY created_at DESC
LIMIT $2;
