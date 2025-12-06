.PHONY: build run migrate migrate-down

build:
	go build -o bin/server ./cmd/server

run: build
	./bin/server

migrate:
	go run ./cmd/migrate

migrate-down:
	go run ./cmd/migrate down
