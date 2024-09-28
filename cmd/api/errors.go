package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Internal server error: [%s] path: %s, err: %v", r.Method, r.URL.Path, err)

	_ = writeJSONError(w, http.StatusInternalServerError, "the server encountered an error")
}

func (app *application) badRequest(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Bad request: [%s] path: %s, err: %v", r.Method, r.URL.Path, err)

	_ = writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFound(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Not found: [%s] path: %s, err: %v", r.Method, r.URL.Path, err)

	_ = writeJSONError(w, http.StatusNotFound, "not found")
}
