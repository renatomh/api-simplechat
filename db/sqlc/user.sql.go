// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: user.sql

package db

import (
	"context"
	"database/sql"
)

const changeUserPassword = `-- name: ChangeUserPassword :one
UPDATE users
SET
  hash_pass = $2,
  password_changed_at = now()
WHERE id = $1
RETURNING id, created_at, full_name, username, email, avatar_url, last_login_at, hash_pass, password_changed_at
`

type ChangeUserPasswordParams struct {
	ID       int64  `json:"id"`
	HashPass string `json:"hash_pass"`
}

func (q *Queries) ChangeUserPassword(ctx context.Context, arg ChangeUserPasswordParams) (User, error) {
	row := q.db.QueryRowContext(ctx, changeUserPassword, arg.ID, arg.HashPass)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.FullName,
		&i.Username,
		&i.Email,
		&i.AvatarUrl,
		&i.LastLoginAt,
		&i.HashPass,
		&i.PasswordChangedAt,
	)
	return i, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (
  full_name,
  username,
  email,
  hash_pass
) VALUES (
  $1, $2, $3, $4
) RETURNING id, created_at, full_name, username, email, avatar_url, last_login_at, hash_pass, password_changed_at
`

type CreateUserParams struct {
	FullName string         `json:"full_name"`
	Username string         `json:"username"`
	Email    sql.NullString `json:"email"`
	HashPass string         `json:"hash_pass"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.FullName,
		arg.Username,
		arg.Email,
		arg.HashPass,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.FullName,
		&i.Username,
		&i.Email,
		&i.AvatarUrl,
		&i.LastLoginAt,
		&i.HashPass,
		&i.PasswordChangedAt,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const getUser = `-- name: GetUser :one
SELECT id, created_at, full_name, username, email, avatar_url, last_login_at, hash_pass, password_changed_at FROM users
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.FullName,
		&i.Username,
		&i.Email,
		&i.AvatarUrl,
		&i.LastLoginAt,
		&i.HashPass,
		&i.PasswordChangedAt,
	)
	return i, err
}

const getUserByUsername = `-- name: GetUserByUsername :one
SELECT id, created_at, full_name, username, email, avatar_url, last_login_at, hash_pass, password_changed_at FROM users
WHERE username = $1 LIMIT 1
`

func (q *Queries) GetUserByUsername(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByUsername, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.FullName,
		&i.Username,
		&i.Email,
		&i.AvatarUrl,
		&i.LastLoginAt,
		&i.HashPass,
		&i.PasswordChangedAt,
	)
	return i, err
}

const listUsers = `-- name: ListUsers :many
SELECT
  id,
  full_name,
  username,
  email,
  avatar_url,
  last_login_at
FROM users
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListUsersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type ListUsersRow struct {
	ID          int64          `json:"id"`
	FullName    string         `json:"full_name"`
	Username    string         `json:"username"`
	Email       sql.NullString `json:"email"`
	AvatarUrl   sql.NullString `json:"avatar_url"`
	LastLoginAt sql.NullTime   `json:"last_login_at"`
}

func (q *Queries) ListUsers(ctx context.Context, arg ListUsersParams) ([]ListUsersRow, error) {
	rows, err := q.db.QueryContext(ctx, listUsers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListUsersRow
	for rows.Next() {
		var i ListUsersRow
		if err := rows.Scan(
			&i.ID,
			&i.FullName,
			&i.Username,
			&i.Email,
			&i.AvatarUrl,
			&i.LastLoginAt,
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

const updateUser = `-- name: UpdateUser :one
UPDATE users
SET
  full_name = $2,
  username = $3,
  email = $4,
  avatar_url = $5
WHERE id = $1
RETURNING id, created_at, full_name, username, email, avatar_url, last_login_at, hash_pass, password_changed_at
`

type UpdateUserParams struct {
	ID        int64          `json:"id"`
	FullName  string         `json:"full_name"`
	Username  string         `json:"username"`
	Email     sql.NullString `json:"email"`
	AvatarUrl sql.NullString `json:"avatar_url"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUser,
		arg.ID,
		arg.FullName,
		arg.Username,
		arg.Email,
		arg.AvatarUrl,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.FullName,
		&i.Username,
		&i.Email,
		&i.AvatarUrl,
		&i.LastLoginAt,
		&i.HashPass,
		&i.PasswordChangedAt,
	)
	return i, err
}
