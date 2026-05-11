-- name: CreateSite :one
INSERT INTO sites (id, org_id, slug, name, primary_host, status, active_deployment_id, created_at)
VALUES ($1, $2, $3, $4, $5, $6, '', $7)
RETURNING id, org_id, slug, name, primary_host, status, active_deployment_id, created_at;

-- name: GetSite :one
SELECT id, org_id, slug, name, primary_host, status, active_deployment_id, created_at
FROM sites
WHERE id = $1;

-- name: ListSitesByOrg :many
SELECT id, org_id, slug, name, primary_host, status, active_deployment_id, created_at
FROM sites
WHERE org_id = $1
ORDER BY slug;

-- name: UpdateSiteStatus :exec
UPDATE sites
SET status = $2
WHERE id = $1;

-- name: UpdateSiteActiveDeployment :exec
UPDATE sites
SET active_deployment_id = $2, status = 'active'
WHERE id = $1;

-- name: CreateDefaultSiteQuota :exec
INSERT INTO site_quotas (site_id, bundle_max_bytes, db_soft_max_bytes, db_hard_max_bytes, request_timeout_ms, updated_at)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetSiteQuota :one
SELECT site_id, bundle_max_bytes, db_soft_max_bytes, db_hard_max_bytes, request_timeout_ms, updated_at
FROM site_quotas
WHERE site_id = $1;
