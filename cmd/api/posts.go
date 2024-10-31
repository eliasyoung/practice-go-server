package main

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/eliasyoung/go-backend-server-practice/internal/db"
	"github.com/go-chi/chi/v5"
)

type CreatePostPayload struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload

	if err := readJSON(w, r, &payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	post := &db.CreatePostParams{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		UserID:  1,
	}

	ctx := r.Context()

	row, err := app.store.Queries.CreatePost(ctx, *post)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, row)

}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	pid := chi.URLParam(r, "postID")
	id, err := strconv.ParseInt(pid, 10, 64)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()

	post, err := app.store.Queries.GetPostById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			writeJSONError(w, http.StatusNotFound, err.Error())
			return
		}

		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, post)
}
