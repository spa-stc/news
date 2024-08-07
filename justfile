set dotenv-load := true

_default:
	@just --list

lint:
	golangci-lint run

create_migration NAME: 
	migrate create -ext sql -dir db/migrations -seq {{NAME}}

run_migrations:
	migrate -path migrations -database $NEWSLETTER_DATABASE_URL up

generate:
	sqlc generate

build: 
	go build ./...
