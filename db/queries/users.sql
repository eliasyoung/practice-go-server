-- name: CreateUser :one
INSERT INTO users (username, password, email)
VALUES ($1, $2, $3)
RETURNING id, created_at;

-- name: GetUserById :one
SELECT id, username, password, email, created_at
FROM users
WHERE id = $1;