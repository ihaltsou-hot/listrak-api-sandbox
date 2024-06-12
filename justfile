# default
default:
  just --list

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
