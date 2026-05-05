-- name: CreateUser :one
INSERT INTO users (email, password_hash)
VALUES ($1, $2)
RETURNING id, email, is_active, created_at, updated_at;

-- name: GetUserByEmail :one
SELECT id, email, password_hash, is_active, created_at, updated_at
FROM users
WHERE email = $1;