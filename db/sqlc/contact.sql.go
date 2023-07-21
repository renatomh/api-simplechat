// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: contact.sql

package db

import (
	"context"
)

const acceptContact = `-- name: AcceptContact :one
UPDATE contacts
SET
  status = 'Accepted',
  accepted_at = now()
WHERE id = $1
RETURNING id, from_user_id, to_user_id, status, requested_at, accepted_at
`

func (q *Queries) AcceptContact(ctx context.Context, id int64) (Contact, error) {
	row := q.db.QueryRowContext(ctx, acceptContact, id)
	var i Contact
	err := row.Scan(
		&i.ID,
		&i.FromUserID,
		&i.ToUserID,
		&i.Status,
		&i.RequestedAt,
		&i.AcceptedAt,
	)
	return i, err
}

const createContact = `-- name: CreateContact :one
INSERT INTO contacts (
  from_user_id,
  to_user_id,
  status
) VALUES (
  $1, $2, 'Pending'
) RETURNING id, from_user_id, to_user_id, status, requested_at, accepted_at
`

type CreateContactParams struct {
	FromUserID int64 `json:"from_user_id"`
	ToUserID   int64 `json:"to_user_id"`
}

func (q *Queries) CreateContact(ctx context.Context, arg CreateContactParams) (Contact, error) {
	row := q.db.QueryRowContext(ctx, createContact, arg.FromUserID, arg.ToUserID)
	var i Contact
	err := row.Scan(
		&i.ID,
		&i.FromUserID,
		&i.ToUserID,
		&i.Status,
		&i.RequestedAt,
		&i.AcceptedAt,
	)
	return i, err
}

const deleteContact = `-- name: DeleteContact :exec
DELETE FROM contacts WHERE id = $1
`

func (q *Queries) DeleteContact(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteContact, id)
	return err
}

const getContact = `-- name: GetContact :one
SELECT id, from_user_id, to_user_id, status, requested_at, accepted_at FROM contacts
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetContact(ctx context.Context, id int64) (Contact, error) {
	row := q.db.QueryRowContext(ctx, getContact, id)
	var i Contact
	err := row.Scan(
		&i.ID,
		&i.FromUserID,
		&i.ToUserID,
		&i.Status,
		&i.RequestedAt,
		&i.AcceptedAt,
	)
	return i, err
}

const listAcceptedContacts = `-- name: ListAcceptedContacts :many
SELECT id, from_user_id, to_user_id, status, requested_at, accepted_at FROM contacts
WHERE 
  (from_user_id = $1 OR
  to_user_id = $1) AND
  status = 'Accepted'
ORDER BY id
LIMIT $2
OFFSET $3
`

type ListAcceptedContactsParams struct {
	FromUserID int64 `json:"from_user_id"`
	Limit      int32 `json:"limit"`
	Offset     int32 `json:"offset"`
}

func (q *Queries) ListAcceptedContacts(ctx context.Context, arg ListAcceptedContactsParams) ([]Contact, error) {
	rows, err := q.db.QueryContext(ctx, listAcceptedContacts, arg.FromUserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Contact
	for rows.Next() {
		var i Contact
		if err := rows.Scan(
			&i.ID,
			&i.FromUserID,
			&i.ToUserID,
			&i.Status,
			&i.RequestedAt,
			&i.AcceptedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listContacts = `-- name: ListContacts :many
SELECT id, from_user_id, to_user_id, status, requested_at, accepted_at FROM contacts
WHERE 
  from_user_id = $1 OR
  to_user_id = $1
ORDER BY id
LIMIT $2
OFFSET $3
`

type ListContactsParams struct {
	FromUserID int64 `json:"from_user_id"`
	Limit      int32 `json:"limit"`
	Offset     int32 `json:"offset"`
}

func (q *Queries) ListContacts(ctx context.Context, arg ListContactsParams) ([]Contact, error) {
	rows, err := q.db.QueryContext(ctx, listContacts, arg.FromUserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Contact
	for rows.Next() {
		var i Contact
		if err := rows.Scan(
			&i.ID,
			&i.FromUserID,
			&i.ToUserID,
			&i.Status,
			&i.RequestedAt,
			&i.AcceptedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listPendingContacts = `-- name: ListPendingContacts :many
SELECT id, from_user_id, to_user_id, status, requested_at, accepted_at FROM contacts
WHERE 
  (from_user_id = $1 OR
  to_user_id = $1) AND
  status = 'Pending'
ORDER BY id
LIMIT $2
OFFSET $3
`

type ListPendingContactsParams struct {
	FromUserID int64 `json:"from_user_id"`
	Limit      int32 `json:"limit"`
	Offset     int32 `json:"offset"`
}

func (q *Queries) ListPendingContacts(ctx context.Context, arg ListPendingContactsParams) ([]Contact, error) {
	rows, err := q.db.QueryContext(ctx, listPendingContacts, arg.FromUserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Contact
	for rows.Next() {
		var i Contact
		if err := rows.Scan(
			&i.ID,
			&i.FromUserID,
			&i.ToUserID,
			&i.Status,
			&i.RequestedAt,
			&i.AcceptedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listRejectedContacts = `-- name: ListRejectedContacts :many
SELECT id, from_user_id, to_user_id, status, requested_at, accepted_at FROM contacts
WHERE 
  (from_user_id = $1 OR
  to_user_id = $1) AND
  status = 'Rejected'
ORDER BY id
LIMIT $2
OFFSET $3
`

type ListRejectedContactsParams struct {
	FromUserID int64 `json:"from_user_id"`
	Limit      int32 `json:"limit"`
	Offset     int32 `json:"offset"`
}

func (q *Queries) ListRejectedContacts(ctx context.Context, arg ListRejectedContactsParams) ([]Contact, error) {
	rows, err := q.db.QueryContext(ctx, listRejectedContacts, arg.FromUserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Contact
	for rows.Next() {
		var i Contact
		if err := rows.Scan(
			&i.ID,
			&i.FromUserID,
			&i.ToUserID,
			&i.Status,
			&i.RequestedAt,
			&i.AcceptedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const rejectContact = `-- name: RejectContact :one
UPDATE contacts
SET
  status = 'Rejected'
WHERE id = $1
RETURNING id, from_user_id, to_user_id, status, requested_at, accepted_at
`

func (q *Queries) RejectContact(ctx context.Context, id int64) (Contact, error) {
	row := q.db.QueryRowContext(ctx, rejectContact, id)
	var i Contact
	err := row.Scan(
		&i.ID,
		&i.FromUserID,
		&i.ToUserID,
		&i.Status,
		&i.RequestedAt,
		&i.AcceptedAt,
	)
	return i, err
}
