_default:
  just --list

# Run development server. 
devel: 
  go run cmd/newsletter/main.go serve

# Setup default development env.
setup-devel: migrate add-default-admin

# Apply migrations.
migrate:
  go run cmd/newsletter/main.go migrate up

# Add a default admin for development.
add-default-admin:
  go run cmd/newsletter/main.go admin create test@gmail.com 1234567890

build-docker:
  nix build .#bin-docker
