package main

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/eliasyoung/go-backend-server-practice/internal/db"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type userCtx struct {
}

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload db.CreateUserParams

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := validateCreateUserPayload(&payload); err != nil {
		app.badRequestResponse(w, r, db.ErrNotFound)
		return
	}

	user := &db.CreateUserParams{
		Username: payload.Username,
		Password: payload.Password,
		Email:    payload.Email,
	}

	ctx := r.Context()

	row, err := app.store.Queries.CreateUser(ctx, *user)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, row)

}

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// temp struct for non auth user from ctx
type FollowUser struct {
	UserID int64 `json:"user_id"`
}

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	userToFollow := getUserFromCtx(r)

	// TODO: revert back to auth userID from ctx

	var payload FollowUser

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := validateFollowUserPayload(&payload); err != nil {
		app.badRequestResponse(w, r, errors.New("请求参数错误!"))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), db.QueryTimeoutDuration)
	defer cancel()

	followParams := db.FollowParams{
		UserID:     userToFollow.ID,
		FollowerID: payload.UserID,
	}

	if err := app.store.Queries.Follow(ctx, followParams); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.notFoundResponse(w, r, db.ErrNotFound)
			return
		}
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			app.conflictResponse(w, r, db.ErrConflict)
			return
		}
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23503" {
			app.badRequestResponse(w, r, db.ErrConflict)
			return
		}
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	userToUnfollow := getUserFromCtx(r)

	// TODO: revert back to auth userID from ctx

	var payload FollowUser

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := validateFollowUserPayload(&payload); err != nil {
		app.badRequestResponse(w, r, errors.New("请求参数错误!"))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), db.QueryTimeoutDuration)
	defer cancel()

	unfollowParams := db.UnfollowParams{
		UserID:     userToUnfollow.ID,
		FollowerID: payload.UserID,
	}

	if err := app.store.Queries.Unfollow(ctx, unfollowParams); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.notFoundResponse(w, r, db.ErrNotFound)
			return
		}

		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) usersContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uid, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)

		if err != nil {
			if errors.Is(err, strconv.ErrSyntax) {
				app.badRequestResponse(w, r, err)
				return
			}
			app.internalServerError(w, r, err)
			return
		}

		ctx := r.Context()

		user, err := app.store.Queries.GetUserById(ctx, uid)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				app.notFoundResponse(w, r, db.ErrNotFound)
				return
			}

			app.badRequestResponse(w, r, err)
			return
		}

		ctx = context.WithValue(r.Context(), userCtx{}, &user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromCtx(r *http.Request) *db.GetUserByIdRow {
	user, _ := r.Context().Value(userCtx{}).(*db.GetUserByIdRow)

	return user
}
