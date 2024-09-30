package main

import (
	"context"
	"errors"
	"github.com/minhnghia2k3/GOssage/internal/store"
	"net/http"
)

const userCtx = "users"

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := getUserFromContext(r)

	followedID, err := parseID(r, "userID")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.storage.Followers.Follow(r.Context(), followerUser.ID, followedID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrConflict):
			app.conflictResponse(w, r, err)
		case errors.Is(err, store.ErrFollowSelf):
			app.badRequestResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := getUserFromContext(r)

	unfollowedID, err := parseID(r, "userID")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.storage.Followers.Unfollow(r.Context(), followerUser.ID, unfollowedID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

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

func getUserFromContext(r *http.Request) *store.User {
	return r.Context().Value(userCtx).(*store.User)
}
