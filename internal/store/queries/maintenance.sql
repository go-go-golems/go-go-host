-- name: ListPrunableDeployments :many
SELECT id, site_id, version, status, bundle_ref, unpacked_path, manifest_json, validation_json, created_by_type, created_by_id, created_at, activated_at, bundle_sha256
FROM deployments
WHERE site_id = $1
  AND id <> $2
  AND status = ANY($3::text[])
  AND created_at < $4
ORDER BY created_at ASC;

-- name: DeleteDeployment :exec
DELETE FROM deployments
WHERE id = $1;

-- name: DeleteAuditEventsBefore :one
WITH deleted AS (
  DELETE FROM audit_log
  WHERE created_at < $1
  RETURNING id
)
SELECT COUNT(*)::bigint AS deleted_count FROM deleted;
