package service

import (
	"context"

	"github.com/eliasyoung/go-backend-server-practice/internal/db"
)

type CommentService struct {
	store *db.Store
}

func NewCommentServce(store *db.Store) *CommentService {
	return &CommentService{
		store: store,
	}
}

func (s *CommentService) GetCommentsByPostId(ctx context.Context, pid int64) ([]db.GetCommentsByPostIdRow, error) {
	ctx, cancel := context.WithTimeout(ctx, db.QueryTimeoutDuration)
	defer cancel()

	comments, err := s.store.Queries.GetCommentsByPostId(ctx, pid)
	if err != nil {
		return make([]db.GetCommentsByPostIdRow, 0), err
	}

	return comments, nil
}
