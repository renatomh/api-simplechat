-- name: CreateChat :one
INSERT INTO chats (
  from_user_id,
  to_user_id
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetChat :one
SELECT * FROM chats
WHERE id = $1 LIMIT 1;

-- name: GetChatByUserIDs :one
SELECT * FROM chats
WHERE
  (from_user_id = $1 AND to_user_id = $2) OR 
  (from_user_id = $2 AND to_user_id = $1)
LIMIT 1;

-- name: ListChats :many
SELECT * FROM chats
WHERE 
  from_user_id = $1 OR
  to_user_id = $1
ORDER BY last_message_received_at
LIMIT $2
OFFSET $3;

-- name: UpdateChat :one
UPDATE chats
SET last_message_received_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteChat :exec
DELETE FROM chats WHERE id = $1;
