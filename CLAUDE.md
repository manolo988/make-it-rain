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

## Frontend Quick Reference

### Tech Stack
- **Build**: Vite
- **Framework**: React/Vue/Svelte (choose based on interview)
- **API**: Fetch/Axios to `VITE_API_URL` (localhost:8080/api/v1)

### Directory Structure
```
frontend/
├── src/
│   ├── components/      # UI components
│   ├── hooks/          # Custom React hooks
│   ├── services/       # API calls
│   ├── utils/          # Helpers
│   └── App.tsx         # Main component
```

### Implementation Pattern
**IMPORTANT**: Always reference BE models/endpoints when implementing FE. Never modify BE to conform to FE.

1. **Check BE first**: Review models/ and routes/routes.go for available endpoints
2. **Create API service**: Match exact BE endpoints and payloads
3. **Build component**: Use BE model structure for state/props

### API Integration Example
```typescript
// services/userService.ts - Match BE endpoints exactly
const API_URL = import.meta.env.VITE_API_URL;

export const userService = {
  // GET /api/v1/users
  getUsers: (page = 1, pageSize = 10) =>
    fetch(`${API_URL}/users?page=${page}&page_size=${pageSize}`),

  // POST /api/v1/users - Match models/user.go fields
  createUser: (data: {email: string, name: string, password: string}) =>
    fetch(`${API_URL}/users`, {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify(data)
    })
};
```

### Component Example
```tsx
// components/UserList.tsx - Use BE model structure
interface User {
  id: number;          // Match models/user.go
  email: string;
  name: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}
```

### Frontend Commands
```bash
cd frontend
npm install          # Install dependencies
npm run dev          # Start dev server (port 3000)
npm run build        # Production build
```

### Interview Tips for FE
1. **Always check BE first** - Look at models/ and controllers/ before implementing
2. **Use exact field names** - `is_active` not `isActive`, `created_at` not `createdAt`
3. **Reference BE validation** - Check controller bindings for required fields
4. **Test with BE running** - Ensure `make run` is active on port 8080
5. **CORS is handled** - Middleware already configured in BE

### Quick Component Creation
1. Check BE endpoint: `grep -r "router\." routes/`
2. Check model fields: `cat models/[resource].go`
3. Create service matching BE exactly
4. Build component using BE field names
5. Test integration with running BE