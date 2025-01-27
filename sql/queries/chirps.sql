-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid(),
    NOW() ::TIMESTAMP,
    NOW() ::TIMESTAMP,
    $1,
    $2
)
RETURNING *;

-- name: GetAllChirps :many
SELECT * 
FROM chirps
WHERE
(user_id = sqlc.narg('user_id') OR sqlc.narg('user_id') IS NULL)
ORDER BY created_at ASC;

-- name: GetSingleChirp :one
SELECT *
FROM chirps
WHERE id = $1;

-- name: DeleteChirpById :exec
DELETE 
FROM chirps
WHERE user_id = $1 AND id = $2;

-- name: GetChirpsByAuthorId :many
SELECT * 
FROM chirps
WHERe user_id = $1
ORDER BY created_at ASC;