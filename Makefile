.PHONY: run build clean db-up db-down sqlc migrate-up migrate-down migrate-create

include .env
export

run:
	go run ./cmd/server

build:
	go build -o bin/server ./cmd/server

clean:
	rm -rf bin/

db-up:
	docker compose up -d postgres

db-down:
	docker compose down

sqlc:
	sqlc generate

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down 1

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name
