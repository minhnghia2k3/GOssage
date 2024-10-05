package main

import (
	"github.com/minhnghia2k3/GOssage/internal/database"
	"github.com/minhnghia2k3/GOssage/internal/env"
	"github.com/minhnghia2k3/GOssage/internal/store"
	"log"
)

func main() {
	addr := env.GetString("DATABASE_ADDR", "postgres://root:secret@localhost:5432/gossage?sslmode=disable")
	conn, err := database.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	storage := store.NewStorage(conn)

	database.Seed(storage, conn)
}
