// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: chirps.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const createChirp = `-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid(),
    NOW() ::TIMESTAMP,
    NOW() ::TIMESTAMP,
    $1,
    $2
)
RETURNING id, created_at, updated_at, body, user_id
`

type CreateChirpParams struct {
	Body   string
	UserID uuid.UUID
}

func (q *Queries) CreateChirp(ctx context.Context, arg CreateChirpParams) (Chirp, error) {
	row := q.db.QueryRowContext(ctx, createChirp, arg.Body, arg.UserID)
	var i Chirp
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Body,
		&i.UserID,
	)
	return i, err
}

const deleteChirpById = `-- name: DeleteChirpById :exec
DELETE 
FROM chirps
WHERE user_id = $1 AND id = $2
`

type DeleteChirpByIdParams struct {
	UserID uuid.UUID
	ID     uuid.UUID
}

func (q *Queries) DeleteChirpById(ctx context.Context, arg DeleteChirpByIdParams) error {
	_, err := q.db.ExecContext(ctx, deleteChirpById, arg.UserID, arg.ID)
	return err
}

const getAllChirps = `-- name: GetAllChirps :many
SELECT id, created_at, updated_at, body, user_id 
FROM chirps
WHERE
(user_id = $1 OR $1 IS NULL)
ORDER BY created_at ASC
`

func (q *Queries) GetAllChirps(ctx context.Context, userID uuid.NullUUID) ([]Chirp, error) {
	rows, err := q.db.QueryContext(ctx, getAllChirps, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Chirp
	for rows.Next() {
		var i Chirp
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Body,
			&i.UserID,
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

const getChirpsByAuthorId = `-- name: GetChirpsByAuthorId :many
SELECT id, created_at, updated_at, body, user_id 
FROM chirps
WHERe user_id = $1
ORDER BY created_at ASC
`

func (q *Queries) GetChirpsByAuthorId(ctx context.Context, userID uuid.UUID) ([]Chirp, error) {
	rows, err := q.db.QueryContext(ctx, getChirpsByAuthorId, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Chirp
	for rows.Next() {
		var i Chirp
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Body,
			&i.UserID,
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

const getSingleChirp = `-- name: GetSingleChirp :one
SELECT id, created_at, updated_at, body, user_id
FROM chirps
WHERE id = $1
`

func (q *Queries) GetSingleChirp(ctx context.Context, id uuid.UUID) (Chirp, error) {
	row := q.db.QueryRowContext(ctx, getSingleChirp, id)
	var i Chirp
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Body,
		&i.UserID,
	)
	return i, err
}
