-- name: CreateMessage :one
INSERT INTO messages (
  chat_id,
  from_user_id,
  to_user_id,
  body
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetMessage :one
SELECT * FROM messages
WHERE id = $1 LIMIT 1;

-- name: ListMessages :many
SELECT * FROM messages
WHERE chat_id = $1
ORDER BY sent_at DESC
LIMIT $2
OFFSET $3;

-- name: DeleteMessage :exec
DELETE FROM messages WHERE id = $1;
