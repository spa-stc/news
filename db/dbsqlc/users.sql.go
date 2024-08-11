// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: users.sql

package dbsqlc

import (
	"context"

	"github.com/google/uuid"
)

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, name, email, password_hash, is_admin, status, created_ts, updated_ts FROM users WHERE email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, db DBTX, email string) (User, error) {
	row := db.QueryRow(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.PasswordHash,
		&i.IsAdmin,
		&i.Status,
		&i.CreatedTs,
		&i.UpdatedTs,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, name, email, password_hash, is_admin, status, created_ts, updated_ts FROM users WHERE id = $1
`

func (q *Queries) GetUserByID(ctx context.Context, db DBTX, id uuid.UUID) (User, error) {
	row := db.QueryRow(ctx, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.PasswordHash,
		&i.IsAdmin,
		&i.Status,
		&i.CreatedTs,
		&i.UpdatedTs,
	)
	return i, err
}

const insertUser = `-- name: InsertUser :one
INSERT INTO users (
	id,
	name,
	email,
	password_hash
) VALUES (
	$1,
	$2,
	$3,
	$4
) RETURNING id, name, email, password_hash, is_admin, status, created_ts, updated_ts
`

type InsertUserParams struct {
	ID           uuid.UUID
	Name         string
	Email        string
	PasswordHash string
}

func (q *Queries) InsertUser(ctx context.Context, db DBTX, arg InsertUserParams) (User, error) {
	row := db.QueryRow(ctx, insertUser,
		arg.ID,
		arg.Name,
		arg.Email,
		arg.PasswordHash,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.PasswordHash,
		&i.IsAdmin,
		&i.Status,
		&i.CreatedTs,
		&i.UpdatedTs,
	)
	return i, err
}
