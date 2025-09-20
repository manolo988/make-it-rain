.PHONY: help run build test clean docker-up docker-down docker-build migrate-up migrate-down migrate-reset migrate-force migrate-fix-dirty migrate-create deps lint fmt air frontend-install frontend-dev frontend-build dev

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

migrate-force: ## Force migration version and clear dirty flag (usage: make migrate-force VERSION=123456)
	@if [ -z "$(VERSION)" ]; then echo "Please provide a version: make migrate-force VERSION=<version_number>"; exit 1; fi
	@echo "Forcing migration version to $(VERSION) and clearing dirty flag..."
	@$(GO) run scripts/migrate.go force $(VERSION)

migrate-fix-dirty: ## Manually fix dirty migration flag in database
	@echo "⚠️  This will manually clear the dirty flag in the schema_migrations table"
	@echo "Current migration status:"
	@source .env && PGPASSWORD=$$DATABASE_PASSWORD psql -h $$DATABASE_HOST -p $$DATABASE_PORT -U $$DATABASE_USER -d $$DATABASE_NAME -t -c "SELECT version, dirty FROM schema_migrations;" || echo "No migrations table found"
	@echo ""
	@echo "To clear the dirty flag, press Enter. To cancel, press Ctrl+C"
	@read confirm
	@source .env && PGPASSWORD=$$DATABASE_PASSWORD psql -h $$DATABASE_HOST -p $$DATABASE_PORT -U $$DATABASE_USER -d $$DATABASE_NAME -c "UPDATE schema_migrations SET dirty = false;" && echo "✅ Dirty flag cleared successfully"
	@echo "New migration status:"
	@source .env && PGPASSWORD=$$DATABASE_PASSWORD psql -h $$DATABASE_HOST -p $$DATABASE_PORT -U $$DATABASE_USER -d $$DATABASE_NAME -t -c "SELECT version, dirty FROM schema_migrations;"

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

frontend-install: ## Install frontend dependencies
	cd frontend && npm install

frontend-dev: ## Run frontend development server
	cd frontend && npm run dev

frontend-build: ## Build frontend for production
	cd frontend && npm run build

frontend-clean: ## Clean frontend build files
	rm -rf frontend/dist frontend/node_modules

dev: ## Run both backend and frontend in development mode (requires two terminals)
	@echo "Starting backend and frontend servers..."
	@echo "Backend will run on http://localhost:8080"
	@echo "Frontend will run on http://localhost:3000"
	@echo ""
	@echo "Run these commands in separate terminals:"
	@echo "  make air        # Backend with hot reload"
	@echo "  make frontend-dev  # Frontend with hot reload"

full-build: build frontend-build ## Build both backend and frontend for production

full-clean: clean frontend-clean ## Clean both backend and frontend build files