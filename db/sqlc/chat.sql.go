// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: chat.sql

package db

import (
	"context"
)

const createChat = `-- name: CreateChat :one
INSERT INTO chats (
  from_user_id,
  to_user_id
) VALUES (
  $1, $2
) RETURNING id, from_user_id, to_user_id, last_message_received_at
`

type CreateChatParams struct {
	FromUserID int64 `json:"from_user_id"`
	ToUserID   int64 `json:"to_user_id"`
}

func (q *Queries) CreateChat(ctx context.Context, arg CreateChatParams) (Chat, error) {
	row := q.db.QueryRowContext(ctx, createChat, arg.FromUserID, arg.ToUserID)
	var i Chat
	err := row.Scan(
		&i.ID,
		&i.FromUserID,
		&i.ToUserID,
		&i.LastMessageReceivedAt,
	)
	return i, err
}

const deleteChat = `-- name: DeleteChat :exec
DELETE FROM chats WHERE id = $1
`

func (q *Queries) DeleteChat(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteChat, id)
	return err
}

const getChat = `-- name: GetChat :one
SELECT id, from_user_id, to_user_id, last_message_received_at FROM chats
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetChat(ctx context.Context, id int64) (Chat, error) {
	row := q.db.QueryRowContext(ctx, getChat, id)
	var i Chat
	err := row.Scan(
		&i.ID,
		&i.FromUserID,
		&i.ToUserID,
		&i.LastMessageReceivedAt,
	)
	return i, err
}

const getChatByUserIDs = `-- name: GetChatByUserIDs :one
SELECT id, from_user_id, to_user_id, last_message_received_at FROM chats
WHERE
  (from_user_id = $1 AND to_user_id = $2) OR 
  (from_user_id = $2 AND to_user_id = $1)
LIMIT 1
`

type GetChatByUserIDsParams struct {
	FromUserID int64 `json:"from_user_id"`
	ToUserID   int64 `json:"to_user_id"`
}

func (q *Queries) GetChatByUserIDs(ctx context.Context, arg GetChatByUserIDsParams) (Chat, error) {
	row := q.db.QueryRowContext(ctx, getChatByUserIDs, arg.FromUserID, arg.ToUserID)
	var i Chat
	err := row.Scan(
		&i.ID,
		&i.FromUserID,
		&i.ToUserID,
		&i.LastMessageReceivedAt,
	)
	return i, err
}

const listChats = `-- name: ListChats :many
SELECT id, from_user_id, to_user_id, last_message_received_at FROM chats
WHERE 
  from_user_id = $1 OR
  to_user_id = $1
ORDER BY last_message_received_at
LIMIT $2
OFFSET $3
`

type ListChatsParams struct {
	FromUserID int64 `json:"from_user_id"`
	Limit      int32 `json:"limit"`
	Offset     int32 `json:"offset"`
}

func (q *Queries) ListChats(ctx context.Context, arg ListChatsParams) ([]Chat, error) {
	rows, err := q.db.QueryContext(ctx, listChats, arg.FromUserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Chat{}
	for rows.Next() {
		var i Chat
		if err := rows.Scan(
			&i.ID,
			&i.FromUserID,
			&i.ToUserID,
			&i.LastMessageReceivedAt,
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

const updateChat = `-- name: UpdateChat :one
UPDATE chats
SET last_message_received_at = now()
WHERE id = $1
RETURNING id, from_user_id, to_user_id, last_message_received_at
`

func (q *Queries) UpdateChat(ctx context.Context, id int64) (Chat, error) {
	row := q.db.QueryRowContext(ctx, updateChat, id)
	var i Chat
	err := row.Scan(
		&i.ID,
		&i.FromUserID,
		&i.ToUserID,
		&i.LastMessageReceivedAt,
	)
	return i, err
}
