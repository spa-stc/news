_default:
	@just --list

lint:
	golangci-lint run
