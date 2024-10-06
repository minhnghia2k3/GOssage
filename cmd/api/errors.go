package main

import (
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Infow(
		"internal server error",
		"method", r.Method,
		"path", r.URL.Path,
		"error", err.Error(),
	)

	writeJSONError(w, http.StatusInternalServerError, "the server encountered an error")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Infow(
		"bad request",
		"method", r.Method,
		"path", r.URL.Path,
		"error", err.Error(),
	)

	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Infow(
		"conflict error",
		"method", r.Method,
		"path", r.URL.Path,
		"error", err.Error(),
	)
	writeJSONError(w, http.StatusConflict, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Infow(
		"not found",
		"method", r.Method,
		"path", r.URL.Path,
		"error", err.Error(),
	)

	writeJSONError(w, http.StatusNotFound, "not found")
}

func (app *application) unauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Infow(
		"unauthorized error",
		"method", r.Method,
		"path", r.URL.Path,
		"error", err.Error(),
	)

	writeJSONError(w, http.StatusUnauthorized, err.Error())
}

func (app *application) forbiddenResponse(w http.ResponseWriter, r *http.Request) {
	app.logger.Infow(
		"forbidden",
		"method", r.Method,
		"path", r.URL.Path,
		"error",
	)

	writeJSONError(w, http.StatusForbidden, "forbidden")
}

func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	app.logger.Infow(
		"rate limit exceeded",
		"method", r.Method,
		"path", r.URL.Path,
	)

	writeJSONError(w, http.StatusTooManyRequests, "too many requests")
}
