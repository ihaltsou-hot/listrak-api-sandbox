# default
default:
  just --list

install:
	@brew install golang-migrate
	@go install github.com/a-h/templ/cmd/templ@latest
	@go install github.com/air-verse/air@latest
	@go get ./...
	@go mod vendor
	@go mod tidy
	@go mod download
	@npm install -D tailwindcss
	@npm install -D daisyui@latest

# Database migration up
db-up:
	@go run cmd/migrate/main.go up

# Database reset
db-reset:
	@go run cmd/reset/main.go

# Database migration down
db-down:
	@go run cmd/migrate/main.go down

# Migrations against the database
db-migration name:
	@migrate create -ext sql -dir cmd/migrate/migrations {{name}}
