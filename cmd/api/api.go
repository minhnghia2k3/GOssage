package main

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/minhnghia2k3/GOssage/docs"
	"github.com/minhnghia2k3/GOssage/internal"
	"github.com/minhnghia2k3/GOssage/internal/auth"
	"github.com/minhnghia2k3/GOssage/internal/env"
	"github.com/minhnghia2k3/GOssage/internal/store"
	"github.com/minhnghia2k3/GOssage/internal/store/cache"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	limiter     limiterConfig
}

type limiterConfig struct {
	rps     float64
	burst   int64
	enabled bool
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
		AllowedOrigins:   []string{env.GetString("CORS_ALLOW_ORIGIN", "http://localhost:3000")},
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
	r.Use(app.rateLimiter)

	// Define routes
	r.Route("/v1", func(r chi.Router) {

		r.Get("/healthcheck", app.healthCheckHandler)
		r.Get("/debug/vars", expvar.Handler().ServeHTTP)

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

	// GRACEFUL SHUTDOWN
	shutdown := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		app.logger.Infow("signal caught", "signal", s.String())

		shutdown <- srv.Shutdown(ctx)
	}()

	app.logger.Infow("Server is running...", "addr", app.config.addr, "env", app.config.env)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	// Block until receive shutdown channel
	err = <-shutdown
	if err != nil {
		return err
	}

	app.logger.Infow("Server shutdown successfully", "addr", app.config.addr, "env", app.config.env)

	return nil
}
