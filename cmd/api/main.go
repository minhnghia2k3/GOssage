package main

import (
	"github.com/joho/godotenv"
	"github.com/minhnghia2k3/GOssage/internal/database"
	"github.com/minhnghia2k3/GOssage/internal/env"
	"github.com/minhnghia2k3/GOssage/internal/store"
	"log"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v\n", err)
	}
}

// @title						GopherSocial API
// @description				API for GopherSocial, a social network for gophers
// @termsOfService				http://swagger.io/terms/
// @contact.name				API Support
// @contact.url				http://www.swagger.io/support
// @contact.email				support@swagger.io
// @license.name				Apache 2.0
// @license.url				http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath					/v1
//
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
func main() {
	cfg := config{
		addr:   env.GetString("SERVER_ADDR", ":8080"),
		apiURL: env.GetString("EXTERNAL_URL", "localhost:8080"),
		dbConfig: dbConfig{
			dsn:          env.GetString("DATABASE_ADDR", "postgres://root:secret@localhost:5432/gossage?sslmode=disable"),
			maxOpenConns: env.GetInt("MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
	}

	// Initialize connection pool
	db, err := database.New(
		cfg.dbConfig.dsn,
		cfg.dbConfig.maxOpenConns,
		cfg.dbConfig.maxIdleConns,
		cfg.dbConfig.maxIdleTime,
	)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	log.Printf("Database connection pool established\n")

	// Initialize storage layer
	s := store.NewStorage(db)

	app := &application{
		config:  cfg,
		storage: s,
	}

	h := app.mount()

	log.Fatal(app.serve(h))
}
