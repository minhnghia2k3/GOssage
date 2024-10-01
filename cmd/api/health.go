package main

import (
	"net/http"
)

// Check system health
//
//	@Summary		Healthcheck
//	@Description	check system health return {status, environment, version}
//	@Tags			Ops
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	string	"ok"
//	@Router			/healthcheck [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"status":      "ok",
		"environment": app.config.env,
		"version":     version,
	}

	if err := app.jsonResponse(w, http.StatusOK, data); err != nil {
		app.internalServerError(w, r, err)
	}
}
