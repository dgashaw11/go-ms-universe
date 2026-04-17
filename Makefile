help: ## Display available commands
	@grep -E '^[\.a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

lint: ## Run linter
	go mod tidy
	golangci-lint run

lint.fix: ## Auto-fix lint issues
	go mod tidy
	golangci-lint run --fix ./...

test.unit: ## Run unit tests (no Docker needed)
	go test -timeout 1m -count=1 -race -skip Integration ./...

test.integration: ## Run integration tests (requires Docker)
	go test -timeout 2m -count=1 -race -run Integration ./...

test.all: ## Run all tests
	go test -timeout 2m -count=1 -race ./...

test.cover: ## Run all tests with coverage
	go test -timeout 2m -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

generate: ## Generate mocks
	go generate ./...

build: ## Build both services
	go build -o bin/products ./cmd/products
	go build -o bin/notifications ./cmd/notifications

run.products: ## Run the products service locally
	set -o allexport; source .env; set +o allexport && go run ./cmd/products

run.notifications: ## Run the notifications service locally
	set -o allexport; source .env; set +o allexport && go run ./cmd/notifications

docker.up: ## Start local dependencies (Postgres, Kafka, Prometheus)
	docker compose up -d

docker.down: ## Stop local dependencies
	docker compose down -v

modernise: ## Apply modern Go idioms
	go fix ./...
