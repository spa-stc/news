_default:
	@just --list

lint:
	golangci-lint run

create_migration NAME: 
	migrate create -ext sql -dir db/migrations -seq {{NAME}}

run_migrations:
	#!/usr/bin/env bash
	[[ -f .env ]] && export $(cat .env | xargs)
	migrate -path db/migrations -database $NEWSLETTER_DATABASE_URL up

generate:
	sqlc generate

build: 
	go build ./...
