package main

import (
	"context"
	"errors"
	"github.com/minhnghia2k3/GOssage/internal/store"
	"net/http"
	"strconv"
	"strings"
)

const userCtx = "user"

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := int64(1)
		//if err != nil {
		//	app.badRequestResponse(w, r, err)
		//	return
		//}

		user, err := app.storage.Users.GetByID(r.Context(), userID)

		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx := context.WithValue(r.Context(), userCtx, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			app.unauthorizedErrorResponse(w, r, errors.New("missing Authorization header"))
			return
		}

		parts := strings.Split(header, " ")
		if len(parts) != 2 && parts[0] != "Bearer" {
			app.unauthorizedErrorResponse(w, r, errors.New("invalid Authorization header"))
			return
		}

		token, err := app.authenticator.ValidateToken(parts[1])
		if err != nil {
			app.unauthorizedErrorResponse(w, r, err)
			return
		}

		subject, err := token.Claims.GetSubject()
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		userID, err := strconv.ParseInt(subject, 10, 64)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		user, err := app.storage.Users.GetByID(r.Context(), userID)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx := context.WithValue(r.Context(), userCtx, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
