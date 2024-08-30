set dotenv-load := true

_default:
	@just --list

check: lint test

lint:
	golangci-lint run

create_migration NAME: 
	migrate create -ext sql -dir migrations -seq {{NAME}}

seed_database:
	psql $NEWSLETTER_DATABASE_URL --echo-errors --quiet -c '\timing off' -f db/seeds/reset.sql
	psql $NEWSLETTER_DATABASE_URL --echo-errors --quiet -c '\timing off' -f db/seeds/main.sql

run_migrations:
	migrate -path migrations -database $NEWSLETTER_DATABASE_URL up

down_migrations:
	migrate -path migrations -database $NEWSLETTER_DATABASE_URL down

generate:
	sqlc generate

generate-clean:
	rm -rf ./db/dbsqlc/*.go

build: 
	go build ./...

test: run_migrations seed_database
	go test ./...

run: run_migrations seed_database tailwind
	go run ./main.go

psql_dev: 
	psql $NEWSLETTER_DATABASE_URL

tailwind:
	tailwindcss -i public/tailwind_in.css -o public/assets/main.css
