package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/minhnghia2k3/GOssage/internal"
	"github.com/minhnghia2k3/GOssage/internal/store"
	"net/http"
	"strconv"
	"time"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,gte=2,lte=255"`
	Email    string `json:"email" validate:"required,email,lte=255"`
	Password string `json:"password" validate:"required,gte=8,lte=72"`
}

type UserWithToken struct {
	*store.User
	Token string `json:"token"`
}

// @Summary		Register user
// @Description	register a user and send activation email to them
// @Tags			authentication
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

	// goroutines to send mail
	go app.SendEmail(r.Context(), userWthToken)

	if err := app.jsonResponse(w, http.StatusOK, userWthToken); err != nil {
		app.internalServerError(w, r, err)
	}
}

type CreateUserTokenPayload struct {
	Email    string `json:"email" validate:"required,email,lte=255"`
	Password string `json:"password" validate:"required,gte=8,lte=255"`
}

// @Summary		Creates a token
// @Description	Creates a token for a user
// @Tags			authentication
// @Accept			json
// @Produce		json
// @Param			payload	body		CreateUserTokenPayload	true	"User credentials"
// @Success		200		{string}	string					"Token"
// @Failure		400		{object}	error
// @Failure		401		{object}	error
// @Failure		500		{object}	error
// @Router			/authentication/token [post]
func (app *application) createTokenHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateUserTokenPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user, err := app.storage.Users.GetByEmail(r.Context(), payload.Email)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	claims := jwt.RegisteredClaims{
		Issuer:    "GOssage",
		Subject:   strconv.FormatInt(user.ID, 10),
		Audience:  []string{"GOssage"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(app.config.auth.token.exp)),
		NotBefore: jwt.NewNumericDate(time.Now()),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	jwtToken, err := app.authenticator.GenerateToken(claims)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	response := struct {
		AccessToken string `json:"access_token"`
	}{
		AccessToken: jwtToken,
	}

	if err = app.jsonResponse(w, http.StatusCreated, response); err != nil {
		app.internalServerError(w, r, err)
	}
}

// SendEmail sends an invitation to user Email, with 3 time retries,
// If error occurred, deletes the current user and its invitation.
func (app *application) SendEmail(ctx context.Context, user UserWithToken) {
	activationUrl := fmt.Sprintf("%s/confirm/%s", app.config.frontendURL, user.Token)

	data := struct {
		Username      string
		ActivationURL string
	}{
		Username:      user.Username,
		ActivationURL: activationUrl,
	}

	err := app.mailer.Send(internal.TemplatePath, []string{user.Email}, data)
	if err != nil {
		app.logger.Infow("error sending email", "error", err)

		if err = app.storage.Users.Delete(ctx, user.ID); err != nil {
			app.logger.Infow("error deleting user", "error", err)
		}

	}

	app.logger.Infow("email sent successfully!", "email", user.Email)
}
