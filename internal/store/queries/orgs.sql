-- name: CreateOrg :one
INSERT INTO orgs (id, slug, name, created_at)
VALUES ($1, $2, $3, $4)
RETURNING id, slug, name, created_at;

-- name: GetOrg :one
SELECT id, slug, name, created_at
FROM orgs
WHERE id = $1;

-- name: AddMembership :exec
INSERT INTO memberships (org_id, user_id, role, created_at)
VALUES ($1, $2, $3, $4)
ON CONFLICT (org_id, user_id)
DO UPDATE SET role = EXCLUDED.role;

-- name: GetMembershipRole :one
SELECT role
FROM memberships
WHERE org_id = $1 AND user_id = $2;

-- name: ListOrgsForUser :many
SELECT o.id, o.slug, o.name, o.created_at
FROM orgs o
JOIN memberships m ON m.org_id = o.id
WHERE m.user_id = $1
ORDER BY o.slug;
