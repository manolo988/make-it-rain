# Make It Rain - Go Interview Boilerplate

A production-ready Go REST API boilerplate optimized for onsite interviews. Built with best practices and a familiar structure for quick feature implementation.

## Features

- **Gin Framework** - Fast HTTP router and middleware
- **PostgreSQL** with pgx/v5 - High-performance database driver
- **Database Migrations** - Built-in migration system
- **Configuration Management** - Environment-based config with Viper
- **Structured Logging** - JSON logging with zerolog
- **Graceful Shutdown** - Proper cleanup on termination
- **Docker Support** - Development and production containers
- **Hot Reload** - Development with Air
- **CRUD Example** - Complete User management implementation

## Project Structure

```
make-it-rain/
├── main.go                 # Application entry point
├── config/                 # Configuration management
├── controllers/            # HTTP request handlers
├── db/                     # Database layer & migrations
│   ├── migrations/         # SQL migration files
│   └── *.go               # Database operations
├── middleware/            # HTTP middleware (logging, CORS, recovery)
├── models/                # Data structures
├── routes/                # API route definitions
├── services/              # Business logic layer
├── utils/                 # Helper functions
└── scripts/               # Utility scripts
```

## Quick Start

### Prerequisites

- Go 1.23+
- PostgreSQL 15+
- Docker & Docker Compose (optional)

### Setup

1. **Clone and setup environment:**
```bash
cp .env.example .env
# Edit .env with your database credentials
```

2. **Start PostgreSQL (using Docker):**
```bash
make docker-dev
```

3. **Install dependencies:**
```bash
make deps
```

4. **Run migrations:**
```bash
make migrate-up
```

5. **Run the application:**
```bash
make run
# OR with hot reload
make air
```

The API will be available at `http://localhost:8080`

## API Endpoints

### Health Checks
- `GET /health` - Application health status
- `GET /ready` - Readiness check

### User Management (Example CRUD)
- `POST /api/v1/users` - Create user
- `GET /api/v1/users/:id` - Get user by ID
- `GET /api/v1/users` - List users (paginated)
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user

### Example Requests

```bash
# Create user
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","name":"John Doe","password":"password123"}'

# Get users (with pagination)
curl "http://localhost:8080/api/v1/users?page=1&page_size=10&sort_by=created_at&sort_order=desc"
```

## Development

### Available Make Commands

```bash
make help          # Show all available commands
make run           # Run application
make build         # Build binary
make test          # Run tests
make docker-up     # Start all containers
make docker-down   # Stop containers
make migrate-up    # Run migrations
make migrate-down  # Rollback migrations
make lint          # Run linter
make fmt           # Format code
```

### Adding New Features

1. **Create Model** (`models/resource.go`)
2. **Add DB Interface** (`db/db.go`)
3. **Implement DB Methods** (`db/resource.go`)
4. **Create Service** (`services/resource.go`)
5. **Create Controller** (`controllers/resource.go`)
6. **Register Routes** (`routes/routes.go`)
7. **Create Migration** (`make migrate-create NAME=create_resources_table`)

## Database Migrations

```bash
# Create new migration
make migrate-create NAME=add_user_status

# Run migrations
make migrate-up

# Rollback last migration
make migrate-down
```

## Configuration

Configuration is managed via environment variables:

- `SERVER_PORT` - API server port (default: 8080)
- `DATABASE_HOST` - PostgreSQL host
- `DATABASE_PORT` - PostgreSQL port
- `DATABASE_NAME` - Database name
- `DATABASE_USER` - Database user
- `DATABASE_PASSWORD` - Database password
- `JWT_SECRET_KEY` - JWT signing key
- `APP_LOG_LEVEL` - Log level (debug/info/warn/error)

## Best Practices Implemented

- ✅ Interface-based database layer
- ✅ Context propagation for request cancellation
- ✅ Structured error handling
- ✅ Request validation with binding tags
- ✅ Connection pooling with pgx
- ✅ Transaction support
- ✅ Pagination helpers
- ✅ Graceful shutdown
- ✅ Middleware chain (logging, recovery, CORS)
- ✅ Environment-based configuration

## Interview Tips

This boilerplate is designed for rapid feature development during interviews:

1. **Familiar Structure** - Follows common Go project patterns
2. **Quick CRUD** - Copy the User example for new resources
3. **Database Ready** - Migrations and connection pooling configured
4. **No Auth Complexity** - Add JWT/OAuth only if required
5. **Production Patterns** - Shows knowledge of best practices

## License

MIT