-- name: NextDeploymentVersion :one
SELECT COALESCE(MAX(version), 0) + 1 AS version
FROM deployments
WHERE site_id = $1;

-- name: CreateDeployment :one
INSERT INTO deployments (id, site_id, version, status, bundle_ref, unpacked_path, manifest_json, validation_json, created_by_type, created_by_id, created_at, activated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NULL)
RETURNING id, site_id, version, status, bundle_ref, unpacked_path, manifest_json, validation_json, created_by_type, created_by_id, created_at, activated_at;

-- name: GetDeployment :one
SELECT id, site_id, version, status, bundle_ref, unpacked_path, manifest_json, validation_json, created_by_type, created_by_id, created_at, activated_at
FROM deployments
WHERE id = $1;

-- name: ListDeploymentsBySite :many
SELECT id, site_id, version, status, bundle_ref, unpacked_path, manifest_json, validation_json, created_by_type, created_by_id, created_at, activated_at
FROM deployments
WHERE site_id = $1
ORDER BY version DESC;

-- name: UpdateDeploymentArtifacts :exec
UPDATE deployments
SET status = $2, bundle_ref = $3, unpacked_path = $4, manifest_json = $5, validation_json = $6
WHERE id = $1;

-- name: UpdateDeploymentStatus :exec
UPDATE deployments
SET status = $2, validation_json = $3
WHERE id = $1;

-- name: ActivateDeployment :exec
UPDATE deployments
SET status = 'active', activated_at = $2
WHERE id = $1;

-- name: SupersedeActiveDeployments :exec
UPDATE deployments
SET status = 'superseded'
WHERE site_id = $1 AND status = 'active' AND id <> $2;

-- name: PreviousValidatedDeployment :one
SELECT id, site_id, version, status, bundle_ref, unpacked_path, manifest_json, validation_json, created_by_type, created_by_id, created_at, activated_at
FROM deployments
WHERE site_id = $1 AND id <> $2 AND status IN ('validated', 'superseded', 'active')
ORDER BY version DESC
LIMIT 1;
