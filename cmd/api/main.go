package main

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"github.com/minhnghia2k3/GOssage/internal/database"
	"github.com/minhnghia2k3/GOssage/internal/env"
	"github.com/minhnghia2k3/GOssage/internal/store"
	"go.uber.org/zap"
	"log"
	"time"
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
		mail: mailConfig{
			exp: 24 * time.Hour, // 1 day
		},
	}

	// Initialize structured logger
	logger := initLogger()
	defer logger.Sync()

	// Initialize connection pool
	db, err := database.New(
		cfg.dbConfig.dsn,
		cfg.dbConfig.maxOpenConns,
		cfg.dbConfig.maxIdleConns,
		cfg.dbConfig.maxIdleTime,
	)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	logger.Info("Database connection pool established\n")

	// Initialize storage layer
	s := store.NewStorage(db)

	app := &application{
		config:  cfg,
		storage: s,
		logger:  logger,
	}

	h := app.mount()

	logger.Fatal(app.serve(h))
}

func initLogger() *zap.SugaredLogger {
	rawJSON := []byte(`{
	  "level": "info",
	  "encoding": "json",
	  "outputPaths": ["stdout", "/tmp/logs"],
	  "errorOutputPaths": ["stderr"],
	  "encoderConfig": {
	    "messageKey": "message",
	    "levelKey": "level",
	    "levelEncoder": "lowercase"
	  }
	}`)

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}
	logger := zap.Must(cfg.Build()).Sugar()

	return logger
}
