set dotenv-load := true

_default:
	@just --list

lint:
	golangci-lint run

create_migration NAME: 
	migrate create -ext sql -dir db/migrations -seq {{NAME}}

seed_database:
	psql $NEWSLETTER_DATABASE_URL --echo-errors --quiet -c '\timing off' -f seeds/reset.sql
	psql $NEWSLETTER_DATABASE_URL --echo-errors --quiet -c '\timing off' -f seeds/main.sql

run_migrations:
	migrate -path migrations -database $NEWSLETTER_DATABASE_URL up

generate:
	sqlc generate

build: 
	go build ./...

test: run_migrations seed_database
	go test ./...
