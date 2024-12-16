package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/eliasyoung/go-backend-server-practice/internal/db"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type postKey string

const postCtx postKey = "post"

type CreatePostPayload struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := validateCreatePostPayload(&payload); err != nil {
		app.badRequestResponse(w, r, errors.New("请求参数错误!"))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), db.QueryTimeoutDuration)
	defer cancel()

	row, err := app.service.Post.CreatePost(ctx, payload.Title, payload.Content, payload.Tags, 1)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, row)

}

type GetPostResultWithComments struct {
	ID        int64                       `json:"id"`
	UserID    int64                       `json:"user_id"`
	Title     string                      `json:"title"`
	Content   string                      `json:"content"`
	Version   pgtype.Int4                 `json:"version"`
	CreatedAt pgtype.Timestamptz          `json:"created_at"`
	UpdatedAt pgtype.Timestamptz          `json:"updated_at"`
	Tags      []string                    `json:"tags"`
	Comments  []db.GetCommentsByPostIdRow `json:"comments"`
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	ctx, cancel := context.WithTimeout(r.Context(), db.QueryTimeoutDuration)
	defer cancel()

	comments, err := app.store.Queries.GetCommentsByPostId(ctx, post.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	postRes := GetPostResultWithComments{
		ID:        post.ID,
		Title:     post.Title,
		UserID:    post.UserID,
		Content:   post.Content,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
		Version:   post.Version,
		Tags:      post.Tags,
		Comments:  comments,
	}

	postRes.Tags = responseSliceFormater(postRes.Tags)
	postRes.Comments = responseSliceFormater(postRes.Comments)

	if err := app.jsonResponse(w, http.StatusOK, postRes); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {

	pid := chi.URLParam(r, "postID")
	id, err := strconv.ParseInt(pid, 10, 64)
	if err != nil {
		if errors.Is(err, strconv.ErrSyntax) {
			app.badRequestResponse(w, r, err)
			return
		}
		app.internalServerError(w, r, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), db.QueryTimeoutDuration)
	defer cancel()

	rowsAffected, err := app.service.Post.DeletePostById(ctx, id)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if rowsAffected == 0 {
		errStr := fmt.Sprintf("Delete failed, not found post with id %s", pid)
		app.notFoundResponse(w, r, errors.New(errStr))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type UpdatePostPayload struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	var payload UpdatePostPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := validateUpdatePostPayload(&payload); err != nil {
		app.badRequestResponse(w, r, errors.New("请求参数错误!"))
		return
	}

	if payload.Title != nil {
		post.Title = *payload.Title
	}

	if payload.Content != nil {
		post.Content = *payload.Content
	}

	updatePost := db.UpdatePostByIdParams{
		ID:      post.ID,
		Title:   post.Title,
		Content: post.Content,
		Version: post.Version,
	}

	ctx, cancel := context.WithTimeout(r.Context(), db.QueryTimeoutDuration)
	defer cancel()

	if err := app.service.Post.UpdatePostById(ctx, updatePost); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.notFoundResponse(w, r, db.ErrNotFound)
			return
		}
		app.internalServerError(w, r, err)
		return
	}

	post.Tags = responseSliceFormater(post.Tags)

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) getAllPostsHandler(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), db.QueryTimeoutDuration)
	defer cancel()

	posts, err := app.service.Post.GetAllPosts(ctx)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if len(posts) != 0 {
		for i, post := range posts {
			posts[i].Tags = responseSliceFormater(post.Tags)
		}
	} else {
		posts = responseSliceFormater(posts)
	}

	if err := app.jsonResponse(w, http.StatusOK, posts); err != nil {
		app.internalServerError(w, r, err)
	}

}

func (app *application) postsContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pid := chi.URLParam(r, "postID")
		id, err := strconv.ParseInt(pid, 10, 64)
		if err != nil {
			if errors.Is(err, strconv.ErrSyntax) {
				app.badRequestResponse(w, r, err)
				return
			}
			app.internalServerError(w, r, err)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), db.QueryTimeoutDuration)
		defer cancel()

		post, err := app.service.Post.GetPostById(ctx, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				app.notFoundResponse(w, r, db.ErrNotFound)
				return
			}

			app.badRequestResponse(w, r, err)
			return
		}

		ctx = context.WithValue(r.Context(), postCtx, &post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromCtx(r *http.Request) *db.GetPostByIdRow {
	post, _ := r.Context().Value(postCtx).(*db.GetPostByIdRow)

	return post
}
