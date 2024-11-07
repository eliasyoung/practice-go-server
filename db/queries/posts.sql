-- name: CreatePost :one
INSERT INTO posts (content, title, user_id, tags)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at, updated_at;

-- name: GetPostById :one
SELECT id, user_id, title, content, created_at, updated_at, tags, version
FROM posts
WHERE id = $1;

-- name: DeletePostById :execrows
DELETE FROM posts WHERE id = $1;

-- name: UpdatePostById :one
UPDATE posts
SET title = $1, content = $2, version = version + 1
WHERE id = $3 AND version = $4
RETURNING version;
