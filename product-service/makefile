# Membaca variabel dari file .env
include .env

migrate-up:
	migrate -path db/migrations -database ${DB_URL} up

migrate-down:
	migrate -path db/migrations -database ${DB_URL} down

# Menjalankan aplikasi
run:
	go run cmd/api/main.go