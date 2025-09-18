# Interview Boilerplate - Quick Reference

## Purpose
Production-ready Go REST API boilerplate for onsite interviews. Optimized for rapid feature development.

## Tech Stack
- **Framework**: Gin (REST API)
- **Database**: PostgreSQL with pgx/v5
- **Config**: Viper (env-based)
- **Logging**: zerolog
- **Migrations**: golang-migrate

## Project Structure
```
├── main.go              # Entry point, graceful shutdown
├── config/              # Configuration (Viper-based)
├── controllers/         # HTTP handlers (request/response)
├── db/                  # Database layer (interface pattern)
│   └── migrations/      # SQL migration files
├── middleware/          # Logger, CORS, Recovery
├── models/              # Data structures
├── routes/              # Route definitions
├── services/            # Business logic
└── utils/              # Helpers
```

## Key Patterns

### Adding New Feature (e.g., Product)
1. **Model**: `models/product.go` - Define struct
2. **DB Interface**: Add methods to `db/db.go` interface
3. **DB Implementation**: `db/product.go` - CRUD operations
4. **Service**: `services/product.go` - Business logic
5. **Controller**: `controllers/product.go` - HTTP handlers
6. **Routes**: Register in `routes/routes.go`
7. **Migration**: `make migrate-create NAME=create_products`

### Database Pattern
```go
// Always use interface
dbService := db.NewDBService()

// Transactions
tx, _ := dbService.BeginTx(ctx)
defer tx.Rollback(ctx)
// ... operations
tx.Commit(ctx)

// Pagination
page, pageSize := utils.ValidatePagination(page, pageSize)
offset := (page - 1) * pageSize
```

### Request Flow
```
HTTP Request → Router → Middleware → Controller → Service → DB → Response
```

## Common Commands
```bash
make run          # Run application
make air          # Hot reload
make migrate-up   # Run migrations
make docker-dev   # Start PostgreSQL
make build        # Build binary
```

## Environment Variables
Key configs in `.env`:
- DATABASE_* (host, port, user, password, name)
- SERVER_PORT (default: 8080)
- JWT_SECRET_KEY
- APP_LOG_LEVEL

## Quick Tips for Interview

### Fast CRUD Implementation
1. Copy `user.go` files from each layer
2. Find/replace "User" with resource name
3. Adjust fields as needed
4. Create migration
5. Register routes

### Common Imports
```go
// Controller
"github.com/gin-gonic/gin"
"github.com/manuel/make-it-rain/db"
"github.com/manuel/make-it-rain/services"
"github.com/rs/zerolog/log"

// DB layer
"github.com/jackc/pgx/v5"
"github.com/manuel/make-it-rain/models"

// Service
"github.com/manuel/make-it-rain/db"
"github.com/manuel/make-it-rain/models"
```

### Error Handling
```go
if err != nil {
    if errors.Is(err, pgx.ErrNoRows) {
        c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
        return
    }
    log.Error().Err(err).Msg("Failed operation")
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
}
```

### Validation
```go
type CreateRequest struct {
    Name  string `json:"name" binding:"required"`
    Email string `json:"email" binding:"required,email"`
    Age   int    `json:"age" binding:"min=0,max=150"`
}
```

## Testing Connection
```bash
# Health check
curl localhost:8080/health

# Create user
curl -X POST localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","name":"Test","password":"password123"}'
```

## Interview Strategy
1. **Start with model** - Define data structure
2. **Implement DB layer** - Show interface pattern knowledge
3. **Add business logic** - Demonstrate separation of concerns
4. **Wire up controller** - Handle HTTP properly
5. **Test with curl** - Verify implementation

## Gotchas
- Viper needs explicit `BindEnv()` for env vars
- Use context for all DB operations
- Remember to handle pagination limits
- Always validate/sanitize user input
- Use transactions for multi-step operations