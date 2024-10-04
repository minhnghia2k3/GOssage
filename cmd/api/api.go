package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/minhnghia2k3/GOssage/docs"
	_ "github.com/minhnghia2k3/GOssage/docs"
	"github.com/minhnghia2k3/GOssage/internal"
	"github.com/minhnghia2k3/GOssage/internal/auth"
	"github.com/minhnghia2k3/GOssage/internal/store"
	"github.com/minhnghia2k3/GOssage/internal/store/cache"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type application struct {
	config        config
	storage       store.Storage
	cacheStorage  *cache.Storage
	logger        *zap.SugaredLogger
	mailer        internal.Client
	authenticator auth.Authenticator
}

type config struct {
	addr        string
	dbConfig    dbConfig
	env         string
	apiURL      string
	mail        mailConfig
	frontendURL string
	auth        authConfig
	redisConfig redisConfig
}

type redisConfig struct {
	addr    string
	pw      string
	db      int
	enabled bool
}

type authConfig struct {
	token tokenConfig
}

type tokenConfig struct {
	secret string
	exp    time.Duration
	iss    string
}

type mailConfig struct {
	exp       time.Duration
	fromEmail string
	dialer    dialer
}

type dialer struct {
	host     string
	port     int
	username string
	password string
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

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

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
			r.Use(app.AuthMiddleware)
			r.Post("/", app.createPostHandler)

			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.getPostMiddleware)
				r.Get("/", app.getPostHandler)
				r.Patch("/", app.checkPostOwnerShip("moderator", app.updatePostHandler))
				r.Delete("/", app.checkPostOwnerShip("admin", app.deletePostHandler))
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Put("/activate/{token}", app.activeUserHandler)

			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.AuthMiddleware)
				r.Get("/", app.getUserHandler)
				r.Put("/follows", app.followUserHandler)
				r.Put("/unfollows", app.unfollowUserHandler)
			})

			r.Group(func(r chi.Router) {
				r.Use(app.AuthMiddleware)
				r.Get("/feed", app.getUserFeedHandler)
			})
		})

		r.Route("/authentication", func(r chi.Router) {
			r.Post("/users", app.registerUserHandler)
			r.Post("/token", app.createTokenHandler)
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
