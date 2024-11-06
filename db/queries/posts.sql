-- name: CreatePost :one
INSERT INTO posts (content, title, user_id, tags)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at, updated_at;

-- name: GetPostById :one
SELECT id, user_id, title, content, created_at, updated_at, tags
FROM posts
WHERE id = $1;

-- name: DeletePostById :execrows
DELETE FROM posts WHERE id = $1;
