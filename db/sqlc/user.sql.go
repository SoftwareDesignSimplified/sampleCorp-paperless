// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: user.sql

package db

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
  username,
  hashed_password,
  full_name,
  email,
  department
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING id, username, hashed_password, full_name, email, department, password_changed_at, created_at
`

type CreateUserParams struct {
	Username       string `json:"username"`
	HashedPassword string `json:"hashed_password"`
	FullName       string `json:"full_name"`
	Email          string `json:"email"`
	Department     string `json:"department"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.Username,
		arg.HashedPassword,
		arg.FullName,
		arg.Email,
		arg.Department,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.HashedPassword,
		&i.FullName,
		&i.Email,
		&i.Department,
		&i.PasswordChangedAt,
		&i.CreatedAt,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM "users" WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteUser, id)
	return err
}

const getUserById = `-- name: GetUserById :one
SELECT id, username, hashed_password, full_name, email, department, password_changed_at, created_at FROM "users" WHERE id = $1
`

func (q *Queries) GetUserById(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRow(ctx, getUserById, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.HashedPassword,
		&i.FullName,
		&i.Email,
		&i.Department,
		&i.PasswordChangedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getUserByUserNameOrEmail = `-- name: GetUserByUserNameOrEmail :one
SELECT id, username, hashed_password, full_name, email, department, password_changed_at, created_at FROM "users" WHERE username = $1 OR email = $2
`

type GetUserByUserNameOrEmailParams struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (q *Queries) GetUserByUserNameOrEmail(ctx context.Context, arg GetUserByUserNameOrEmailParams) (User, error) {
	row := q.db.QueryRow(ctx, getUserByUserNameOrEmail, arg.Username, arg.Email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.HashedPassword,
		&i.FullName,
		&i.Email,
		&i.Department,
		&i.PasswordChangedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getUsers = `-- name: GetUsers :many
SELECT id, username, hashed_password, full_name, email, department, password_changed_at, created_at FROM "users" ORDER BY id DESC LIMIT $1 OFFSET $2
`

type GetUsersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) GetUsers(ctx context.Context, arg GetUsersParams) ([]User, error) {
	rows, err := q.db.Query(ctx, getUsers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.HashedPassword,
			&i.FullName,
			&i.Email,
			&i.Department,
			&i.PasswordChangedAt,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUser = `-- name: UpdateUser :one
UPDATE "users" SET username = $1, email = $2 WHERE id = $3 RETURNING id, username, hashed_password, full_name, email, department, password_changed_at, created_at
`

type UpdateUserParams struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	ID       int64  `json:"id"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, updateUser, arg.Username, arg.Email, arg.ID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.HashedPassword,
		&i.FullName,
		&i.Email,
		&i.Department,
		&i.PasswordChangedAt,
		&i.CreatedAt,
	)
	return i, err
}