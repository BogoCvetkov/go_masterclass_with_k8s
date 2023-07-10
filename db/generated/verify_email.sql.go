// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: verify_email.sql

package db

import (
	"context"
)

const createVerifyEmail = `-- name: CreateVerifyEmail :one
INSERT INTO verify_emails (
    user_id,
    email,
    secret_code
) VALUES (
    $1, $2, $3
) RETURNING id, user_id, email, secret_code, is_used, created_at, expired_at
`

type CreateVerifyEmailParams struct {
	UserID     int64  `json:"user_id"`
	Email      string `json:"email"`
	SecretCode string `json:"secret_code"`
}

func (q *Queries) CreateVerifyEmail(ctx context.Context, arg CreateVerifyEmailParams) (VerifyEmail, error) {
	row := q.db.QueryRowContext(ctx, createVerifyEmail, arg.UserID, arg.Email, arg.SecretCode)
	var i VerifyEmail
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Email,
		&i.SecretCode,
		&i.IsUsed,
		&i.CreatedAt,
		&i.ExpiredAt,
	)
	return i, err
}

const updateVerifyEmail = `-- name: UpdateVerifyEmail :one
UPDATE verify_emails
SET
    is_used = TRUE
WHERE
    email = $1
    AND secret_code = $2
    AND is_used = FALSE
    AND expired_at > now()
RETURNING id, user_id, email, secret_code, is_used, created_at, expired_at
`

type UpdateVerifyEmailParams struct {
	Email      string `json:"email"`
	SecretCode string `json:"secret_code"`
}

func (q *Queries) UpdateVerifyEmail(ctx context.Context, arg UpdateVerifyEmailParams) (VerifyEmail, error) {
	row := q.db.QueryRowContext(ctx, updateVerifyEmail, arg.Email, arg.SecretCode)
	var i VerifyEmail
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Email,
		&i.SecretCode,
		&i.IsUsed,
		&i.CreatedAt,
		&i.ExpiredAt,
	)
	return i, err
}
