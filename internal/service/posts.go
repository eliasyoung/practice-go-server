package service

import (
	"context"

	"github.com/eliasyoung/go-backend-server-practice/internal/db"
)

type PostService struct {
	store *db.Store
}

func NewPostService(store *db.Store) *PostService {
	return &PostService{
		store: store,
	}
}

func (s *PostService) CreatePost(ctx context.Context, title string, content string, tags []string, userId int64) (db.CreatePostRow, error) {
	post := db.CreatePostParams{
		Title:   title,
		Content: content,
		Tags:    tags,
		UserID:  userId,
	}

	row, err := s.store.Queries.CreatePost(ctx, post)
	if err != nil {
		return db.CreatePostRow{}, err
	}

	return row, nil
}
