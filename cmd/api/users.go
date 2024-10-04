package main

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/minhnghia2k3/GOssage/internal/store"
	"net/http"
	"strconv"
)

// getUserHandler gets user by id
//
//	@Summary		Get user
//	@Description	get user by given id
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path	int	true	"User ID"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	store.User
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Router			/users/{userID} [get]
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "userID")

	userID, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user, err := app.getUser(r.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err = app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
	}
}

// @Summary		Follow user
// @Description	authenticated user follow provided user
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			userID	path	int	true	"User ID"
// @Security		ApiKeyAuth
// @Success		204
// @Failure		400	{object}	error
// @Failure		409	{object}	error
// @Failure		404	{object}	error
// @Failure		500	{object}	error
// @Router			/users/{userID}/follows [put]
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

// @Summary		Unfollow user
// @Description	authenticated user unfollow provided user
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			userID	path	int	true	"User ID"
// @Security		ApiKeyAuth
// @Success		204
// @Failure		400	{object}	error
// @Failure		404	{object}	error
// @Failure		500	{object}	error
// @Router			/users/{userID}/unfollows [put]
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

// @Summary		Active user
// @Description	active user by using given token
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			token	path	string	true	"Activation token"
// @Success		204
// @Failure		400	{object}	error
// @Failure		404	{object}	error
// @Failure		500	{object}	error
// @Router			/users/activate/{token} [put]
func (app *application) activeUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	if err := app.storage.Users.Activate(r.Context(), token); err != nil {
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

func getUserFromContext(r *http.Request) *store.User {
	return r.Context().Value(userCtx).(*store.User)
}
