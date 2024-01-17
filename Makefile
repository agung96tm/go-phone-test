## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'


# === Application ===

## run/api: run the api application
.PHONY: run/api
run/api:
	go run ./cmd/api -addr :8000

## run/web: run the web application
.PHONY: run/web
run/web:
	go run ./cmd/web -addr :3000


# === Migrations ===

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migration/up: confirm
	@echo "Running up migrations..."
	migrate -path ./migrations -database ${DB_DSN} up

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migration/new:
	@echo "Create migration files for ${name}"
	migrate create -seq -ext=.sql -dir=./migrations ${name}
