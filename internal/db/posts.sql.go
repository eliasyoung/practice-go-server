// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: posts.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createPost = `-- name: CreatePost :one
INSERT INTO posts (content, title, user_id, tags)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at, updated_at
`

type CreatePostParams struct {
	Content string   `json:"content"`
	Title   string   `json:"title"`
	UserID  int64    `json:"user_id"`
	Tags    []string `json:"tags"`
}

type CreatePostRow struct {
	ID        int64              `json:"id"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) (CreatePostRow, error) {
	row := q.db.QueryRow(ctx, createPost,
		arg.Content,
		arg.Title,
		arg.UserID,
		arg.Tags,
	)
	var i CreatePostRow
	err := row.Scan(&i.ID, &i.CreatedAt, &i.UpdatedAt)
	return i, err
}

const deletePostById = `-- name: DeletePostById :execrows
DELETE FROM posts WHERE id = $1
`

func (q *Queries) DeletePostById(ctx context.Context, id int64) (int64, error) {
	result, err := q.db.Exec(ctx, deletePostById, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const getAllPosts = `-- name: GetAllPosts :many
SELECT p.id, p.user_id, p.title, p.content, p.created_at, p.tags, u.username,
COUNT (c.id) AS comments_count
FROM posts p
LEFT JOIN comments c ON c.post_id = p.id
LEFT JOIN users u ON p.user_id = u.id
GROUP BY p.id, u.username
`

type GetAllPostsRow struct {
	ID            int64              `json:"id"`
	UserID        int64              `json:"user_id"`
	Title         string             `json:"title"`
	Content       string             `json:"content"`
	CreatedAt     pgtype.Timestamptz `json:"created_at"`
	Tags          []string           `json:"tags"`
	Username      pgtype.Text        `json:"username"`
	CommentsCount int64              `json:"comments_count"`
}

func (q *Queries) GetAllPosts(ctx context.Context) ([]GetAllPostsRow, error) {
	rows, err := q.db.Query(ctx, getAllPosts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllPostsRow
	for rows.Next() {
		var i GetAllPostsRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Title,
			&i.Content,
			&i.CreatedAt,
			&i.Tags,
			&i.Username,
			&i.CommentsCount,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPostById = `-- name: GetPostById :one
SELECT id, user_id, title, content, created_at, updated_at, tags, version
FROM posts
WHERE id = $1
`

type GetPostByIdRow struct {
	ID        int64              `json:"id"`
	UserID    int64              `json:"user_id"`
	Title     string             `json:"title"`
	Content   string             `json:"content"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
	Tags      []string           `json:"tags"`
	Version   pgtype.Int4        `json:"version"`
}

func (q *Queries) GetPostById(ctx context.Context, id int64) (GetPostByIdRow, error) {
	row := q.db.QueryRow(ctx, getPostById, id)
	var i GetPostByIdRow
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Title,
		&i.Content,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Tags,
		&i.Version,
	)
	return i, err
}

const getPostsWithMetaData = `-- name: GetPostsWithMetaData :many
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
LIMIT $3 OFFSET $4
`

type GetPostsWithMetaDataParams struct {
	UserID      int64              `json:"user_id"`
	Column2     interface{}        `json:"column_2"`
	Limit       int32              `json:"limit"`
	Offset      int32              `json:"offset"`
	Column5     pgtype.Text        `json:"column_5"`
	Tags        []string           `json:"tags"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	CreatedAt_2 pgtype.Timestamptz `json:"created_at_2"`
}

type GetPostsWithMetaDataRow struct {
	ID            int64              `json:"id"`
	UserID        int64              `json:"user_id"`
	Title         string             `json:"title"`
	Content       string             `json:"content"`
	CreatedAt     pgtype.Timestamptz `json:"created_at"`
	Version       pgtype.Int4        `json:"version"`
	Tags          []string           `json:"tags"`
	Username      pgtype.Text        `json:"username"`
	CommentsCount int64              `json:"comments_count"`
}

func (q *Queries) GetPostsWithMetaData(ctx context.Context, arg GetPostsWithMetaDataParams) ([]GetPostsWithMetaDataRow, error) {
	rows, err := q.db.Query(ctx, getPostsWithMetaData,
		arg.UserID,
		arg.Column2,
		arg.Limit,
		arg.Offset,
		arg.Column5,
		arg.Tags,
		arg.CreatedAt,
		arg.CreatedAt_2,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPostsWithMetaDataRow
	for rows.Next() {
		var i GetPostsWithMetaDataRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Title,
			&i.Content,
			&i.CreatedAt,
			&i.Version,
			&i.Tags,
			&i.Username,
			&i.CommentsCount,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updatePostById = `-- name: UpdatePostById :one
UPDATE posts
SET title = $1, content = $2, version = version + 1
WHERE id = $3 AND version = $4
RETURNING version
`

type UpdatePostByIdParams struct {
	Title   string      `json:"title"`
	Content string      `json:"content"`
	ID      int64       `json:"id"`
	Version pgtype.Int4 `json:"version"`
}

func (q *Queries) UpdatePostById(ctx context.Context, arg UpdatePostByIdParams) (pgtype.Int4, error) {
	row := q.db.QueryRow(ctx, updatePostById,
		arg.Title,
		arg.Content,
		arg.ID,
		arg.Version,
	)
	var version pgtype.Int4
	err := row.Scan(&version)
	return version, err
}
