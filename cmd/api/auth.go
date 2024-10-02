package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/google/uuid"
	"github.com/minhnghia2k3/GOssage/internal/store"
	"net/http"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,gte=2,lte=255"`
	Email    string `json:"email" validate:"required,email,lte=255"`
	Password string `json:"password" validate:"required,gte=8,lte=255"` // TODO: TEST 72 BYTES
}

type UserWithToken struct {
	*store.User
	Token string `json:"token"`
}

// @Summary		Register user
// @Description	register a user and send activation email to them
// @Tags			authentications
// @Accept			json
// @Produce		json
// @Param			register	body		RegisterUserPayload	true	"Register payload"
// @Success		200			{object}	UserWithToken
// @Failure		400			{object}	error
// @Failure		409			{object}	error
// @Failure		500			{object}	error
// @Router			/authentication/users [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
	}

	// Hash and set the password
	if err := user.Password.Set(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// Generate plain token
	plainToken := uuid.New().String()

	// Hash plain token for storing
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	if err := app.storage.Users.CreateAndInvite(r.Context(), user, hashToken, app.config.mail.exp); err != nil {
		switch {
		case errors.Is(err, store.ErrConflict):
			app.conflictResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	userWthToken := UserWithToken{
		User:  user,
		Token: plainToken,
	}

	// TODO: goroutines to send email

	if err := app.jsonResponse(w, http.StatusOK, userWthToken); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
