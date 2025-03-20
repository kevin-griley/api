run: build
	@./bin/api

build: swag
	@go build -o bin/api cmd/api/main.go

test:
	@go test -v ./...

swag:
	@swag init -g ./api/main.go -d cmd,internal

up: 
	@godotenv -f .env goose up

down:
	@godotenv -f .env goose down

create: 
	@godotenv -f .env goose create $(name) sql

.PHONY: native

native:
	@cd ./native && npm install && npm run ios
