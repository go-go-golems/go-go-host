-- name: ListAdminOrgs :many
SELECT
  o.id,
  o.slug,
  o.name,
  o.created_at,
  COUNT(DISTINCT m.user_id)::bigint AS member_count,
  COUNT(DISTINCT s.id)::bigint AS site_count,
  COUNT(DISTINCT d.id)::bigint AS deployment_count
FROM orgs o
LEFT JOIN memberships m ON m.org_id = o.id
LEFT JOIN sites s ON s.org_id = o.id
LEFT JOIN deployments d ON d.site_id = s.id
GROUP BY o.id, o.slug, o.name, o.created_at
ORDER BY o.slug;

-- name: ListAdminUsers :many
SELECT
  u.id,
  u.email,
  u.display_name,
  u.created_at,
  u.last_login_at,
  EXISTS (SELECT 1 FROM platform_admins pa WHERE pa.user_id = u.id) AS platform_admin,
  COUNT(DISTINCT m.org_id)::bigint AS org_count
FROM users u
LEFT JOIN memberships m ON m.user_id = u.id
GROUP BY u.id, u.email, u.display_name, u.created_at, u.last_login_at
ORDER BY u.email, u.id;

-- name: ListAdminSites :many
SELECT
  s.id,
  s.org_id,
  o.slug AS org_slug,
  o.name AS org_name,
  s.slug,
  s.name,
  s.primary_host,
  s.status,
  s.active_deployment_id,
  s.created_at,
  COALESCE(rs.status, 'stopped') AS runtime_status,
  COALESCE(rs.requests_total, 0)::bigint AS requests_total,
  COALESCE(rs.errors_total, 0)::bigint AS errors_total,
  COALESCE(rs.last_error, '') AS last_error
FROM sites s
JOIN orgs o ON o.id = s.org_id
LEFT JOIN runtime_status rs ON rs.site_id = s.id
ORDER BY o.slug, s.slug;

-- name: ListAdminDeployments :many
SELECT
  d.id,
  d.site_id,
  s.slug AS site_slug,
  s.primary_host,
  s.org_id,
  o.slug AS org_slug,
  o.name AS org_name,
  d.version,
  d.status,
  d.bundle_ref,
  d.unpacked_path,
  d.manifest_json,
  d.validation_json,
  d.created_by_type,
  d.created_by_id,
  d.created_at,
  d.activated_at
FROM deployments d
JOIN sites s ON s.id = d.site_id
JOIN orgs o ON o.id = s.org_id
WHERE (sqlc.narg('org_id')::text IS NULL OR s.org_id = sqlc.narg('org_id')::text)
  AND (sqlc.narg('site_id')::text IS NULL OR d.site_id = sqlc.narg('site_id')::text)
  AND (sqlc.narg('status')::text IS NULL OR d.status = sqlc.narg('status')::text)
ORDER BY d.created_at DESC
LIMIT sqlc.arg('limit')::int;

-- name: ListAdminAgents :many
SELECT
  a.id,
  a.org_id,
  o.slug AS org_slug,
  o.name AS org_name,
  a.name,
  a.status,
  a.created_by_user_id,
  a.created_at,
  a.last_seen_at,
  COUNT(DISTINCT g.site_id)::bigint AS grant_count
FROM agents a
JOIN orgs o ON o.id = a.org_id
LEFT JOIN agent_site_grants g ON g.agent_id = a.id
WHERE (sqlc.narg('org_id')::text IS NULL OR a.org_id = sqlc.narg('org_id')::text)
  AND (sqlc.narg('status')::text IS NULL OR a.status = sqlc.narg('status')::text)
GROUP BY a.id, a.org_id, o.slug, o.name, a.name, a.status, a.created_by_user_id, a.created_at, a.last_seen_at
ORDER BY o.slug, a.name;

-- name: ListAdminAuditEvents :many
SELECT id, org_id, actor_type, actor_id, action, resource_type, resource_id, ip_address, user_agent, metadata_json, created_at
FROM audit_log
WHERE (sqlc.narg('org_id')::text IS NULL OR org_id = sqlc.narg('org_id')::text)
  AND (sqlc.narg('resource_id')::text IS NULL OR resource_id = sqlc.narg('resource_id')::text)
  AND (sqlc.narg('actor_type')::text IS NULL OR actor_type = sqlc.narg('actor_type')::text)
  AND (sqlc.narg('actor_id')::text IS NULL OR actor_id = sqlc.narg('actor_id')::text)
  AND (sqlc.narg('action')::text IS NULL OR action = sqlc.narg('action')::text)
ORDER BY created_at DESC
LIMIT sqlc.arg('limit')::int;
