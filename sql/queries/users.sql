-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW() ::TIMESTAMP,
    NOW() ::TIMESTAMP,
    $1,
    $2
)
RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users
WHERE 1=1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;