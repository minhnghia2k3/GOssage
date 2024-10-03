package main

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/minhnghia2k3/GOssage/internal/store"
	"net/http"
	"strconv"
	"time"
)

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,lte=255"`
	Content string   `json:"content" validate:"required,lte=255"`
	Tags    []string `json:"tags" validate:"omitempty,dive,lte=100"`
}
type UpdatePostPayload struct {
	Title   *string  `json:"title" validate:"omitempty,lte=255"`
	Content *string  `json:"content" validate:"omitempty,lte=255"`
	Tags    []string `json:"tags" validate:"omitempty,dive,lte=100"`
}

const postCtx = "post"

// getPostHandler gets post by provided post ID,
// then get all its comments,
// response result with http.StatusOK
//
//	@Summary		Get post
//	@Description	get post by id
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			postID	path	int	true	"Post ID"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	store.Post
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Router			/posts/{postID} [get]
func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	post := r.Context().Value(postCtx).(*store.Post)

	// fetch comments of the post
	comments, err := app.storage.Comments.GetByPostID(context.Background(), post.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	post.Comments = comments

	if err = app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
	}
}

// createPostHandler creates new post using request body data.
//
//	@Summary		Create a post
//	@Description	create a new post
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			post	body	CreatePostPayload	true	"Create post payload"
//	@Security		ApiKeyAuth
//	@Success		201	{object}	store.Post
//	@Failure		400	{object}	error
//	@Failure		500	{object}	error
//	@Router			/posts [post]
func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload

	user := r.Context().Value(userCtx).(*store.User)

	err := readJSON(w, r, &payload)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err = Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	post := store.Post{
		UserID:  user.ID,
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
	}

	err = app.storage.Posts.Create(context.Background(), &post)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err = app.jsonResponse(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
	}
}

// updatePostHandler updates an existing post with provided data
//
//	@Summary		Update post
//	@Description	update post by id
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			postID	path	int					true	"Post ID"
//	@Param			post	body	UpdatePostPayload	true	"Update post"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	store.Post
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Router			/posts/{postID} [patch]
func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	var payload UpdatePostPayload
	post := r.Context().Value(postCtx).(*store.Post)

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if payload.Title != nil {
		post.Title = *payload.Title
	}
	if payload.Content != nil {
		post.Content = *payload.Content
	}
	if payload.Tags != nil {
		post.Tags = payload.Tags
	}
	post.UpdatedAt = time.Now()

	if err := app.storage.Posts.Update(context.Background(), post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
	}
}

// deletePostHandler delete a post by given post ID.
//
//	@Summary		Delete post
//	@Description	delete post by id
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			postID	path	int	true	"Post ID"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	store.Post
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Router			/posts/{postID} [delete]
func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	postID, err := parseID(r, "postID")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.storage.Posts.Delete(context.Background(), postID)
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

// parseID parse given key into base64 integer.
func parseID(r *http.Request, key string) (int64, error) {
	id, err := strconv.ParseInt(chi.URLParam(r, key), 10, 64)
	if err != nil {
		return -1, err
	}
	return id, nil
}

// getPostMiddleware gets post by given postID and used as a middleware.
func (app *application) getPostMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postID, err := parseID(r, "postID")
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		post, err := app.storage.Posts.GetByID(context.Background(), postID)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx := context.WithValue(r.Context(), postCtx, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
