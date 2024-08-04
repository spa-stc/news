_default:
	@just --list

lint:
	golangci-lint run

create_migration NAME: 
	migrate create -ext sql -dir db/migrations -seq {{NAME}}

generate:
	sqlc generate
