-- name: CreateUser :one
INSERT INTO users (
  full_name,
  username,
  email,
  hash_pass
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateUser :one
UPDATE users
SET
  full_name = $2,
  username = $3,
  email = $4,
  avatar_url = $5
WHERE id = $1
RETURNING *;

-- name: ChangeUserPassword :one
UPDATE users
SET
  hash_pass = $2,
  password_changed_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;
