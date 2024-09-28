include .env
MIGRATION_PATH= ./cmd/migrate/migrations

PHONY: compose.up
compose.up:
	docker compose up --build -d

# WARNING
PHONY: compose.down
compose.down:
	docker compose down -v

PHONY: migrate
migrate:
	migrate create -ext sql -dir $(MIGRATION_PATH) -seq $(NAME)

PHONY: migrate.force
migrate.force:
	migrate -path $(MIGRATION_PATH) -database $(DATABASE_ADDR) force $(VERSION)

PHONY: migrate.up
migrate.up:
	migrate -database $(DATABASE_ADDR) -path $(MIGRATION_PATH) up

PHONY: migrate.down
migrate.down:
	migrate -database $(DATABASE_ADDR) -path $(MIGRATION_PATH) down

PHONY: seed
seed:
	go run cmd/migrate/seed/main.go