-- name: ListMembershipsForUser :many
SELECT o.id AS org_id, o.slug AS org_slug, o.name AS org_name, m.role, m.created_at
FROM memberships m
JOIN orgs o ON o.id = m.org_id
WHERE m.user_id = $1
ORDER BY o.slug;
