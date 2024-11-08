# make variables from .env available.
-include .env
NEWSLETTER_PUBLIC_DIR ?= ./public
NEWSLETTER_DEVELOPMENT=1

export


.PHONY: default
default:
	@echo "Nothing Here..."

.PHONY: lint
lint:
	@golangci-lint run 

.PHONY: seed-db
seed-db:
	@psql ${NEWSLETTER_DATABASE_URL} --echo-errors --quiet -c '\timing off' -f db/seeds/reset.sql
	@psql ${NEWSLETTER_DATABASE_URL} --echo-errors --quiet -c '\timing off' -f db/seeds/main.sql

.PHONY: migrate
migrate:
	@migrate -path migrations -database ${NEWSLETTER_DATABASE_URL} up

.PHONY: migrate-down
migrate-down:
	@migrate -path migrations -database ${NEWSLETTER_DATABASE_URL} down

.PHONY: generate
generate:
	sqlc generate

.PHONY: build
build:
	go build ./...

.PHONY: test
test: migrate seed-db
	go test ./...

.PHONY: dev 
dev: migrate 
	go run ./main.go

.PHONY: tailwind
tailwind:
	cd $(NEWSLETTER_PUBLIC_DIR) && tailwindcss -c tailwind.config.js -i tailwind_in.css -o assets/main.css

.PHONY: tailwind-watch
tailwind-watch:
	cd $(NEWSLETTER_PUBLIC_DIR) && tailwindcss -w -c tailwind.config.js -i tailwind_in.css -o assets/main.css




