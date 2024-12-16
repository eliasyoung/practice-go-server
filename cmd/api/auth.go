package main

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/eliasyoung/go-backend-server-practice/internal/db"
	"github.com/eliasyoung/go-backend-server-practice/internal/utils"
	"github.com/google/uuid"
)

type RegisterUserPayload struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterUserRes struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}

// registerUserHandler godoc
//
//	@Summary		Registers a user
//	@Description	Registers a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterUserPayload	true	"User credentials"
//	@Success		201		{object}	RegisterUserRes		"User registered"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/user [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := validateRegisterUserPayload(&payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	tempUser := db.CreateUserParams{
		Username: payload.Username,
		Email:    payload.Email,
	}

	// hash password
	hashPw, err := utils.PasswordHasher(payload.Password)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	tempUser.Password = hashPw

	// start the transaction
	ctx := r.Context()

	plainToken := uuid.New().String()

	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	err = app.service.User.CreateAndInviteUser(ctx, tempUser, hashToken, app.config.mail.exp)
	if err != nil {
		switch err {
		case db.ErrDuplicateEmail:
			app.badRequestResponse(w, r, err)
		case db.ErrDuplicateUsername:
			app.badRequestResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	res := RegisterUserRes{
		Username: tempUser.Username,
		Email:    tempUser.Email,
		Token:    plainToken,
	}

	if err := app.jsonResponse(w, http.StatusCreated, res); err != nil {
		app.internalServerError(w, r, err)
	}

}
