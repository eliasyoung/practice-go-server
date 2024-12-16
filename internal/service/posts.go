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

func (s *PostService) GetAllPosts(ctx context.Context) ([]db.GetAllPostsRow, error) {
	posts, err := s.store.Queries.GetAllPosts(ctx)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *PostService) GetPostById(ctx context.Context, pid int64) (db.GetPostByIdRow, error) {
	post, err := s.store.Queries.GetPostById(ctx, pid)
	if err != nil {
		return db.GetPostByIdRow{}, err
	}

	return post, nil
}

func (s *PostService) UpdatePostById(ctx context.Context, updatedPost db.UpdatePostByIdParams) error {
	if _, err := s.store.Queries.UpdatePostById(ctx, updatedPost); err != nil {
		return err
	}

	return nil
}

func (s *PostService) DeletePostById(ctx context.Context, pid int64) (int64, error) {
	rowsAffected, err := s.store.Queries.DeletePostById(ctx, pid)
	if err != nil {
		// maybe logger needed in future.
		return rowsAffected, err
	}
	return rowsAffected, nil
}
