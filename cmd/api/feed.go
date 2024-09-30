package main

import (
	"context"
	"github.com/minhnghia2k3/GOssage/internal/store"
	"net/http"
)

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
