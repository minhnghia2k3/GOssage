package main

import (
	"github.com/joho/godotenv"
	"github.com/minhnghia2k3/GOssage/internal/env"
	"log"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v\n", err)
	}
}

func main() {
	cfg := config{
		addr: env.GetString("SERVER_ADDRESS", ":8080"),
	}

	app := &application{
		config: cfg,
	}

	h := app.mount()

	log.Fatal(app.serve(h))
}
