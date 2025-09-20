.PHONY: help run build test clean docker-up docker-down docker-build migrate-up migrate-down migrate-reset migrate-create deps lint fmt air frontend frontend-build frontend-preview

APP_NAME=make-it-rain
DOCKER_COMPOSE=docker-compose
GO=go
GOTEST=$(GO) test
GOVET=$(GO) vet
GOFMT=gofmt
BINARY_NAME=main

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run: ## Run the application locally
	$(GO) run main.go

build: ## Build the application binary
	CGO_ENABLED=0 $(GO) build -o $(BINARY_NAME) -v

test: ## Run tests
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html

clean: ## Clean build files
	$(GO) clean
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html

docker-up: ## Start docker containers
	$(DOCKER_COMPOSE) up -d

docker-down: ## Stop docker containers
	$(DOCKER_COMPOSE) down

docker-build: ## Build docker image
	$(DOCKER_COMPOSE) build

docker-dev: ## Start only development database
	$(DOCKER_COMPOSE) -f docker-compose.dev.yml up -d

docker-dev-down: ## Stop development database
	$(DOCKER_COMPOSE) -f docker-compose.dev.yml down

migrate-up: ## Run database migrations up
	@echo "Running migrations..."
	@$(GO) run scripts/migrate.go up

migrate-down: ## Rollback database migrations
	@echo "Rolling back migrations..."
	@$(GO) run scripts/migrate.go down 1

migrate-reset: ## Reset database completely (drops all tables and migrations)
	@echo "⚠️  WARNING: This will DROP ALL TABLES and reset the database completely!"
	@echo "Press Ctrl+C to cancel, or Enter to continue..."
	@read confirm
	@source .env && PGPASSWORD=$$DATABASE_PASSWORD psql -h $$DATABASE_HOST -p $$DATABASE_PORT -U $$DATABASE_USER -d $$DATABASE_NAME -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public; GRANT ALL ON SCHEMA public TO $$DATABASE_USER; GRANT ALL ON SCHEMA public TO public;"
	@echo "✅ Database reset complete. All tables and migrations have been removed."

migrate-create: ## Create a new migration file (usage: make migrate-create NAME=migration_name)
	@if [ -z "$(NAME)" ]; then echo "Please provide a migration name: make migrate-create NAME=your_migration_name"; exit 1; fi
	@touch db/migrations/$$(date +%s)_$(NAME).up.sql
	@touch db/migrations/$$(date +%s)_$(NAME).down.sql
	@echo "Created migration files for $(NAME)"

deps: ## Install dependencies
	$(GO) mod download
	$(GO) mod tidy

lint: ## Run linter
	@if ! which golangci-lint > /dev/null; then \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	golangci-lint run

fmt: ## Format code
	$(GOFMT) -s -w .
	$(GO) fmt ./...

air: ## Run with hot reload (requires air)
	@if ! which air > /dev/null; then \
		echo "Installing air..."; \
		go install github.com/air-verse/air@latest; \
	fi
	air

seed: ## Seed the database with sample data
	$(GO) run scripts/seed.go

swagger: ## Generate swagger documentation
	@if ! which swag > /dev/null; then \
		echo "Installing swag..."; \
		go install github.com/swaggo/swag/cmd/swag@latest; \
	fi
	swag init

env: ## Copy .env.example to .env
	cp .env.example .env

psql: ## Connect to PostgreSQL database
	docker exec -it make-it-rain-db-dev psql -U postgres -d make_it_rain_dev

logs: ## Show application logs
	$(DOCKER_COMPOSE) logs -f app

logs-db: ## Show database logs
	$(DOCKER_COMPOSE) logs -f postgres

frontend: ## Run frontend development server
	@echo "Starting frontend development server..."
	cd frontend && npm install && npm run dev

frontend-build: ## Build frontend for production
	@echo "Building frontend for production..."
	cd frontend && npm install && npm run build

frontend-preview: ## Preview frontend production build
	@echo "Starting frontend preview server..."
	cd frontend && npm run preview