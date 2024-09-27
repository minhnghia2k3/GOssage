package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/minhnghia2k3/GOssage/internal/storage"
	"log"
	"net/http"
	"time"
)

type application struct {
	config  config
	storage storage.Storage
}

type config struct {
	addr     string
	dbConfig dbConfig
}

type dbConfig struct {
	dsn          string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// Middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Define routes
	r.Route("/v1", func(r chi.Router) {
		r.Get("/healthcheck", app.healthCheckHandler)
	})

	return r
}

func (app *application) serve(h http.Handler) error {
	srv := http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Server is running on port %s\n", app.config.addr)
	return srv.ListenAndServe()
}
