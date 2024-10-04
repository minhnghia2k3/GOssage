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
			app.unauthorizedErrorResponse(w, r, err)
			return
		}

		userID, err := strconv.ParseInt(subject, 10, 64)
		if err != nil {
			app.unauthorizedErrorResponse(w, r, err)
			return
		}

		user, err := app.getUser(r.Context(), userID)
		if err != nil {
			app.unauthorizedErrorResponse(w, r, err)
			return
		}

		ctx := context.WithValue(r.Context(), userCtx, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) getUser(ctx context.Context, userID int64) (*store.User, error) {
	if !app.config.redisConfig.enabled {
		return app.storage.Users.GetByID(ctx, userID)
	}

	// Get user from cache
	user, err := app.cacheStorage.Users.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Set user to cache server if not existed
	if user == nil {
		user, err = app.storage.Users.GetByID(ctx, userID)
		if err != nil {
			return nil, err
		}

		err = app.cacheStorage.Users.Set(ctx, user)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (app *application) checkPostOwnerShip(role string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(userCtx).(*store.User)
		post := r.Context().Value(postCtx).(*store.Post)

		if post.UserID == user.ID {
			next.ServeHTTP(w, r)
			return
		}

		allow, err := app.checkPriority(r.Context(), user, role)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		if !allow {
			app.forbiddenResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func (app *application) checkPriority(ctx context.Context, user *store.User, roleName string) (bool, error) {
	role, err := app.storage.Roles.GetByName(ctx, roleName)
	if err != nil {
		return false, err
	}

	// Check user's role >= base level => allow update
	// moderator  			>=  moderator
	return user.Role.Level >= role.Level, nil
}
