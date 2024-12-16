-- name: CreateUser :one
INSERT INTO users (username, password, email)
VALUES ($1, $2, $3)
RETURNING id, created_at;

-- name: GetUserById :one
SELECT id, username, password, email, created_at
FROM users
WHERE id = $1;


-- name: GetUsers :many
SELECT id, username, email, created_at
FROM users;

-- name: GetUserFromInvitation :one
SELECT u.id, u.username, u.email, u.created_at, u.is_active
FROM users u
JOIN user_invitations ui
ON u.id = ui.user_id
WHERE ui.token = $1 AND ui.expiry > $2;

-- name: UpdateUserInfoById :exec
UPDATE users
SET username = $1, email = $2, is_active = $3
WHERE id = $4;
