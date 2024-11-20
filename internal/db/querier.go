// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type Querier interface {
	CreatePost(ctx context.Context, arg CreatePostParams) (CreatePostRow, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (CreateUserRow, error)
	DeletePostById(ctx context.Context, id int64) (int64, error)
	Follow(ctx context.Context, arg FollowParams) error
	GetCommentsByPostId(ctx context.Context, postID int64) ([]GetCommentsByPostIdRow, error)
	GetPostById(ctx context.Context, id int64) (GetPostByIdRow, error)
	GetPostsWithMetaData(ctx context.Context, arg GetPostsWithMetaDataParams) ([]GetPostsWithMetaDataRow, error)
	GetUserById(ctx context.Context, id int64) (GetUserByIdRow, error)
	Unfollow(ctx context.Context, arg UnfollowParams) error
	UpdatePostById(ctx context.Context, arg UpdatePostByIdParams) (pgtype.Int4, error)
}

var _ Querier = (*Queries)(nil)
