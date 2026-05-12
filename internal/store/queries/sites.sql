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

-- name: UpsertSiteCapability :exec
INSERT INTO site_capabilities (site_id, capability, enabled, config_json, updated_at)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (site_id, capability)
DO UPDATE SET enabled = EXCLUDED.enabled, config_json = EXCLUDED.config_json, updated_at = EXCLUDED.updated_at;

-- name: ListSiteCapabilities :many
SELECT site_id, capability, enabled, config_json, updated_at
FROM site_capabilities
WHERE site_id = $1
ORDER BY capability;

-- name: ListSiteConfig :many
SELECT site_id, key, value_json, updated_at
FROM site_config
WHERE site_id = $1
ORDER BY key;

-- name: UpsertSiteConfig :exec
INSERT INTO site_config (site_id, key, value_json, updated_at)
VALUES ($1, $2, $3, $4)
ON CONFLICT (site_id, key)
DO UPDATE SET value_json = EXCLUDED.value_json, updated_at = EXCLUDED.updated_at;

-- name: DeleteSiteConfig :exec
DELETE FROM site_config
WHERE site_id = $1 AND key = $2;

-- name: CreateSiteDomain :one
INSERT INTO site_domains (id, site_id, hostname, status, verification_token, verified_at, created_at)
VALUES ($1, $2, $3, $4, $5, NULL, $6)
RETURNING id, site_id, hostname, status, verification_token, verified_at, created_at;

-- name: ListSiteDomains :many
SELECT id, site_id, hostname, status, verification_token, verified_at, created_at
FROM site_domains
WHERE site_id = $1
ORDER BY hostname;

-- name: ListVerifiedSiteDomains :many
SELECT id, site_id, hostname, status, verification_token, verified_at, created_at
FROM site_domains
WHERE site_id = $1 AND status = 'verified'
ORDER BY hostname;

-- name: GetSiteDomain :one
SELECT id, site_id, hostname, status, verification_token, verified_at, created_at
FROM site_domains
WHERE id = $1;

-- name: VerifySiteDomain :one
UPDATE site_domains
SET status = 'verified', verified_at = $2
WHERE id = $1
RETURNING id, site_id, hostname, status, verification_token, verified_at, created_at;

-- name: DeleteSiteDomain :exec
DELETE FROM site_domains
WHERE id = $1;
