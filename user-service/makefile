# Database Connection URL for User Service (Port 5434)
include .env

.PHONY: migrate-up migrate-down

# Run all 'up' migrations
migrate-up:
	migrate -path db/migrations -database "$(DB_URL)" up

# Rollback all 'down' migrations
migrate-down:
	migrate -path db/migrations -database "$(DB_URL)" down -all

run:
	go run cmd/api/main.go