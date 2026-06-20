.PHONY: help build test test-integration vet run-api run-worker run-scheduler up down logs tidy

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-18s\033[0m %s\n", $$1, $$2}'

build: ## Build all binaries into ./bin
	go build -o bin/api ./cmd/api
	go build -o bin/worker ./cmd/worker
	go build -o bin/scheduler ./cmd/scheduler

test: ## Run unit tests
	go test ./...

test-integration: ## Run integration tests (needs TEST_DATABASE_URL)
	go test -tags=integration ./tests/...

vet: ## Run go vet
	go vet ./...

load: ## Run a load test (vars: SCENARIO, C, D). Needs the stack running.
	go run ./loadtest -scenario $(or $(SCENARIO),mixed) -c $(or $(C),100) -d $(or $(D),20s)

tidy: ## Tidy modules
	go mod tidy

run-api: ## Run the API locally
	go run ./cmd/api

run-worker: ## Run a worker locally
	go run ./cmd/worker

run-scheduler: ## Run the scheduler locally
	go run ./cmd/scheduler

up: ## Start the whole stack with Docker
	docker compose up --build

down: ## Stop the stack and remove volumes
	docker compose down -v

logs: ## Tail service logs
	docker compose logs -f api worker scheduler
