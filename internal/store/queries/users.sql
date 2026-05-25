-- name: UpsertUserFromOIDC :one
INSERT INTO users (id, issuer, subject, email, display_name, created_at, last_login_at)
VALUES ($1, $2, $3, $4, $5, $6, $6)
ON CONFLICT (issuer, subject)
DO UPDATE SET
  email = EXCLUDED.email,
  display_name = EXCLUDED.display_name,
  last_login_at = EXCLUDED.last_login_at
RETURNING id, issuer, subject, email, display_name, created_at, last_login_at;

-- name: GetUser :one
SELECT id, issuer, subject, email, display_name, created_at, last_login_at
FROM users
WHERE id = $1;

-- name: AddPlatformAdmin :exec
INSERT INTO platform_admins (user_id, created_at)
VALUES ($1, $2)
ON CONFLICT (user_id) DO NOTHING;

-- name: IsPlatformAdmin :one
SELECT EXISTS (SELECT 1 FROM platform_admins WHERE user_id = $1);
