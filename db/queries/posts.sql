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

-- name: GetPostsWithMetaData :many
SELECT p.id, p.user_id, p.title, p.content, p.created_at, p.version, p.tags, u.username,
COUNT (c.id) AS comments_count
FROM posts p
LEFT JOIN comments c ON c.post_id = p.id
LEFT JOIN users u ON p.user_id = u.id
JOIN followers f ON f.user_id = p.user_id OR p.user_id = $1
WHERE 
    (p.user_id = $1 OR f.follower_id = $1) AND
    (p.title ILIKE '%' || $5 || '%' OR p.content ILIKE '%' || $5 || '%') AND
    (p.tags @> $6 OR $6 = '{}') AND
    (p.created_at >= $7 OR $7 IS NULL) AND
    (p.created_at <= $8 OR $8 IS NULL)
GROUP BY p.id, u.username
ORDER BY 
    CASE WHEN $2 = 'asc' THEN p.created_at END ASC,
    CASE WHEN $2 = 'desc' THEN p.created_at END DESC
LIMIT $3 OFFSET $4;