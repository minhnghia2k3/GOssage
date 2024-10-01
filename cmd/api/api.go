package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/minhnghia2k3/GOssage/docs"
	_ "github.com/minhnghia2k3/GOssage/docs"
	"github.com/minhnghia2k3/GOssage/internal/store"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type application struct {
	config  config
	storage store.Storage
	logger  *zap.SugaredLogger
}

type config struct {
	addr     string
	dbConfig dbConfig
	env      string
	apiURL   string
}

type dbConfig struct {
	dsn          string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

const version = "0.0.1"

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

		docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))

		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.createPostHandler)

			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.getPostMiddleware)
				r.Get("/", app.getPostHandler)
				r.Patch("/", app.updatePostHandler)
				r.Delete("/", app.deletePostHandler)
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.userContextMiddleware)

				r.Get("/", app.getUserHandler)
				r.Put("/follows", app.followUserHandler)
				r.Put("/unfollows", app.unfollowUserHandler)
			})

			r.Get("/feed", app.getUserFeedHandler)
		})
	})

	return r
}

func (app *application) serve(h http.Handler) error {
	// Swagger config
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.BasePath = "/v1"
	docs.SwaggerInfo.Host = app.config.apiURL

	srv := http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	app.logger.Infow("Server is running...", "addr", app.config.addr, "env", app.config.env)
	return srv.ListenAndServe()
}
