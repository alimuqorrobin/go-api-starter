.PHONY: tidy build run migrate

tidy:
	go mod tidy

build:
	go build -o bin/server ./cmd/server

run: build
	./bin/server

migrate:
	go run ./cmd/migrate
