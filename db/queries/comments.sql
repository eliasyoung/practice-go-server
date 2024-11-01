-- name: GetCommentsByPostId :many
SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, users.username, users.id
FROM comments c
JOIN users on users.id = c.user_id
WHERE c.post_id = 2
ORDER BY c.created_at DESC;
