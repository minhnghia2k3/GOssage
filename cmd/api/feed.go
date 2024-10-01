package main

import (
	"context"
	"github.com/minhnghia2k3/GOssage/internal/store"
	"net/http"
)

// @Summary		Fetches the user feed
// @Description	fetched the user feed
// @Tags			feed
// @Accept			json
// @Produce		json
// @Param			limit	query		int		false	"limit"
// @Param			since	query		string	false	"since"
// @Param			until	query		string	false	"until"
// @Param			offset	query		int		false	"offset"
// @Param			sort	query		string	false	"sort"
// @Param			search	query		string	false	"search"
// @Success		200		{object}	store.Post
// @Failure		400		{object}	error
// @Failure		409		{object}	error
// @Failure		404		{object}	error
// @Failure		500		{object}	error
// @Router			/users/feed [get]
func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Get current userID from auth
	userID := int64(1)

	p := store.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
	}

	err := p.Parse(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = Validate.Struct(p)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	feed, err := app.storage.Posts.GetUserFeed(context.Background(), userID, p)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err = app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
	}
}
