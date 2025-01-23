-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email)
VALUES (
    gen_random_uuid(),
    NOW() ::TIMESTAMP,
    NOW() ::TIMESTAMP,
    $1
)
RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users
WHERE 1=1;