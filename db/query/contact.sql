-- name: CreateContact :one
INSERT INTO contacts (
  from_user_id,
  to_user_id,
  status
) VALUES (
  $1, $2, 'Pending'
) RETURNING *;

-- name: GetContact :one
SELECT * FROM contacts
WHERE id = $1 LIMIT 1;

-- name: ListContacts :many
SELECT * FROM contacts
WHERE 
  from_user_id = $1 OR
  to_user_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: ListPendingContacts :many
SELECT * FROM contacts
WHERE 
  (from_user_id = $1 OR
  to_user_id = $1) AND
  status = 'Pending'
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: ListAcceptedContacts :many
SELECT * FROM contacts
WHERE 
  (from_user_id = $1 OR
  to_user_id = $1) AND
  status = 'Accepted'
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: ListRejectedContacts :many
SELECT * FROM contacts
WHERE 
  (from_user_id = $1 OR
  to_user_id = $1) AND
  status = 'Rejected'
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: AcceptContact :one
UPDATE contacts
SET
  status = 'Accepted',
  accepted_at = now()
WHERE id = $1
RETURNING *;

-- name: RejectContact :one
UPDATE contacts
SET
  status = 'Rejected'
WHERE id = $1
RETURNING *;

-- name: DeleteContact :exec
DELETE FROM contacts WHERE id = $1;
