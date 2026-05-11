-- name: UpsertRuntimeStatus :exec
INSERT INTO runtime_status (site_id, org_id, deployment_id, hosts, status, started_at, last_error, requests_total, errors_total, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
ON CONFLICT (site_id)
DO UPDATE SET
  org_id = EXCLUDED.org_id,
  deployment_id = EXCLUDED.deployment_id,
  hosts = EXCLUDED.hosts,
  status = EXCLUDED.status,
  started_at = EXCLUDED.started_at,
  last_error = EXCLUDED.last_error,
  requests_total = EXCLUDED.requests_total,
  errors_total = EXCLUDED.errors_total,
  updated_at = EXCLUDED.updated_at;

-- name: GetRuntimeStatus :one
SELECT site_id, org_id, deployment_id, hosts, status, started_at, last_error, requests_total, errors_total, updated_at
FROM runtime_status
WHERE site_id = $1;

-- name: ListRuntimeStatuses :many
SELECT site_id, org_id, deployment_id, hosts, status, started_at, last_error, requests_total, errors_total, updated_at
FROM runtime_status
ORDER BY site_id;

-- name: ReconcileStaleRuntimeStatuses :exec
UPDATE runtime_status
SET status = 'stopped', last_error = $1, updated_at = $2
WHERE status IN ('starting', 'ready', 'draining');
