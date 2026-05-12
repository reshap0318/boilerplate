---
name: go-coding-rules
description: Go boilerplate coding rules and conventions including architecture, naming conventions, helpers, clients, and anti-patterns
---

## When to use me

Use this skill when:
- Writing new Go code — to follow naming conventions, struct patterns, and architecture rules
- Reviewing code — to check for anti-patterns and rule violations
- You need a reference for available helpers, clients, or GenericRepository methods
- You need to understand how logging, permissions, or response helpers work

> ⚠️ **For CRUD implementation flow** (step-by-step: Model → DTO → Repository → Service → Handler → Routes with full examples), use the **go-development-guide** skill instead. This skill provides the rules and references that the CRUD flow must follow.

---

# Go Boilerplate — Coding Rules & Conventions

> **READ THIS ENTIRELY before writing any code.** These rules are MANDATORY.

## Project Identity

- **Module**: `github.com/reshap0318/go-boilerplate`
- **Go Version**: 1.25.0+
- **Framework**: Gin (HTTP), GORM (ORM), JWT auth, Redis (optional), bcrypt
- **Database**: MySQL (default) / PostgreSQL
- **Entry Point**: `cmd/api/main.go`

## Architecture

```
Routes → Handlers → Services → Repositories → Database/Redis
              ↑
         Middleware (JWT Auth, CORS, Rate Limit)
```

All layers use **Dependency Injection** via `internal/di/container.go`.

## Project Structure

```
go-boilerplate/
├── cmd/
│   ├── api/main.go            # App entry point
│   └── migration/             # DB migration scripts
├── internal/
│   ├── clients/email/         # Email client (injected as s.EmailClient)
│   ├── database/              # DB connections & Redis cache
│   ├── di/container.go        # DI container — wires everything
│   ├── dtos/                  # Request/Response DTOs
│   ├── handlers/              # HTTP handlers (single struct: Handlers)
│   │   └── 00_handlers.go     # Handlers struct definition
│   ├── helpers/               # Pure utility functions
│   ├── middleware/             # JWT, CORS, rate limit middleware
│   ├── models/                # GORM models (must have TableName)
│   ├── repositories/          # Data access layer
│   │   ├── 00_generic.go      # ⛔ DO NOT MODIFY — Generic CRUD
│   │   ├── 00_transaction.go  # ⛔ DO NOT MODIFY — Transaction manager
│   │   └── 00_repository.go   # Registry — add new repos here
│   ├── routes/                # Route registration
│   └── services/              # Business logic (single struct: Services)
│       └── 00_services.go     # Services struct definition
├── storage/                   # Application storage
│   ├── keys/                  # JWT keys (private.pem, public.pem)
│   ├── logs/                  # Application logs (YYYY-MM-DD.log)
│   ├── tmp/                   # Temporary file uploads
│   └── avatars/               # User avatar files
└── docs/                      # Documentation
```

---

## ⛔ CRITICAL RULES

### 1. NEVER modify these files
- `internal/repositories/00_generic.go`
- `internal/repositories/00_transaction.go`

### 2. Models MUST define `TableName()`
```go
func (Permission) TableName() string {
    return "permissions"
}
```

### 3. Single Struct Pattern (NO separate structs)
- ALL service methods → `func (s *Services) ...`
- ALL handler methods → `func (h *Handlers) ...`
- ❌ NEVER create `type PermissionService struct` or `type PermissionHandler struct`

### 4. Function Naming: `{Feature}{Action}` — Feature FIRST
```go
// ✅ CORRECT
func (s *Services) PermissionCreate(...)
func (s *Services) PermissionGetAll(...)
func (h *Handlers) PermissionCreate(...)
func (h *Handlers) PermissionGetByID(...)

// ❌ WRONG — action first
func (s *Services) CreatePermission(...)
func (s *Services) GetAllPermissions(...)
```

### 5. Write Operations MUST use Transaction
```go
// ✅ Create/Update/Delete — always wrap in transaction
s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
    result, err = s.repo.Permission.Create(tx, permission)
    return err
})
```

### 6. Read Operations MUST use `nil` (NOT `.DB`)
```go
// ✅ CORRECT
s.repo.Permission.FindByID(nil, id)
s.repo.Permission.FindAll(nil)

// ❌ WRONG
s.repo.Permission.FindByID(s.repo.Permission.DB, id)
```

---

## 📖 Reference: CRUD Quick Order

> ⚠️ For full CRUD implementation with detailed examples (Model → DTO → Repository → Service → Handler → Routes), use the **go-development-guide** skill.

### Quick Implementation Order
1. Model          → `internal/models/{feature}.go`
2. DTOs           → `internal/dtos/{feature}_dto.go`
3. Repository     → `internal/repositories/{feature}_repository.go`
4. Register Repo  → `internal/repositories/00_repository.go`
5. Service        → `internal/services/{feature}_service.go`
6. Handler        → `internal/handlers/{feature}_handler.go`
7. Routes         → `internal/routes/{feature}_route.go`
8. Register Routes → `cmd/api/main.go`

---

## File Naming Conventions

| Layer | Pattern | Example |
|-------|---------|---------|
| Model | `{feature}.go` | `permission.go` |
| DTO | `{feature}_dto.go` | `permission_dto.go` |
| Repository | `{feature}_repository.go` | `permission_repository.go` |
| Service | `{feature}_service.go` | `permission_service.go` |
| Handler | `{feature}_handler.go` | `permission_handler.go` |
| Route | `{feature}_route.go` | `permission_route.go` |

---

## Logging Rules

Logger is available in Services via `s.Logger`.

| Operation | Must Log? | Operation | Must Log? |
|-----------|-----------|-----------|-----------|
| CREATE | ✅ YES | GET (simple) | ❌ NO |
| UPDATE | ✅ YES | GET (complex) | ✅ YES |
| DELETE | ✅ YES | Auth ops | ✅ YES |

### Logger Methods
```go
s.Logger.LogStart("FuncName", "Message: %s", value)        // Start
s.Logger.LogStep("FuncName", "Step: %s", value)            // Step
s.Logger.LogStepWithPrefix("Func", "[OK]", "Done")         // Step + prefix
s.Logger.LogEnd("FuncName", "Success: %s", value)          // End success
s.Logger.LogEndWithError("Func", "Error: %v", err)         // End error
s.Logger.LogError("FuncName", "Error: %v", err)            // Error
s.Logger.LogWarn("FuncName", "Warning: %s", value)         // Warning
s.Logger.LogInfo("FuncName", "Info: %s", value)            // Info
```

### Example — CREATE (MUST LOG)
```go
func (s *Services) UserCreate(ctx context.Context, email string) error {
    s.Logger.LogStart("UserCreate", "Creating user: %s", email)
    
    user := &models.User{Email: email}
    if err := s.repo.User.Create(s.repo.User.DB, user); err != nil {
        s.Logger.LogEndWithError("UserCreate", "Failed: %v", err)
        return err
    }
    
    s.Logger.LogEnd("UserCreate", "User created: %s (ID: %d)", email, user.ID)
    return nil
}
```

### Example — GET Simple (NO LOG)
```go
func (s *Services) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
    return s.repo.User.FindByID(s.repo.User.DB, id)
}
```

**Log Output**: `storage/logs/YYYY-MM-DD.log` | Auto rotation & cleanup (30 days)

---

## GenericRepository Available Methods

All repositories extend `GenericRepository[T]`. These methods are ALREADY AVAILABLE — do NOT re-implement:

| Method | Signature | Use Case |
|--------|-----------|----------|
| `FindByID` | `(tx *gorm.DB, id uint, preloads ...string)` | Get by ID |
| `FindByIDWithOpts` | `(tx *gorm.DB, id uint, opts *QueryOptions)` | Get by ID + options |
| `Create` | `(tx *gorm.DB, request *T) (*T, error)` | Create record |
| `CreateMany` | `(tx *gorm.DB, request []T) error` | Batch create |
| `Update` | `(tx *gorm.DB, filter *T, update *T) (*T, error)` | Update by filter |
| `UpdateMap` | `(tx *gorm.DB, filter *T, update map[string]interface{})` | Partial update (supports zero values) |
| `Delete` | `(tx *gorm.DB, id uint) (*T, error)` | Soft delete |
| `FindAll` | `(tx *gorm.DB, preloads ...string)` | Get all |
| `FindAllWithOpts` | `(tx *gorm.DB, opts *QueryOptions) (*PagedResult[T], error)` | Paginated list |
| `FindByField` | `(tx *gorm.DB, filter *T, preloads ...string)` | Filter by struct |
| `FindByFieldWithOpts` | `(tx *gorm.DB, filter *T, opts *QueryOptions)` | Filter + pagination |
| `FindByFieldMap` | `(tx *gorm.DB, filter map[string]interface{}, preloads ...string)` | Filter by map |
| `FindByFieldMapWithOpts` | `(tx *gorm.DB, filter map[string]interface{}, opts *QueryOptions)` | Map filter + pagination |
| `Count` | `(tx *gorm.DB) (int64, error)` | Count all |
| `Exists` | `(tx *gorm.DB, filter map[string]interface{}) (bool, error)` | Check exists (map) |
| `ExistsByField` | `(tx *gorm.DB, filter *T) (bool, error)` | Check exists (struct) |

### QueryOptions
```go
type QueryOptions struct {
    Page           int      // Page number (default: 1)
    PageSize       int      // Items per page (default: 10, 0 = no pagination)
    SortBy         string   // Field to sort by
    Order          string   // "ASC" or "DESC"
    Search         string   // Search keyword
    SearchFields   []string // Fields to search
    Preloads       []string // Relations to preload
}
```

### Transaction Manager
```go
// Simple transaction
s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
    // all operations inside use tx
    return nil
})

// Transaction with result
result, err := s.repo.TxManager.WithinTransactionWithResult(func(tx *gorm.DB) (interface{}, error) {
    return someResult, nil
})
```

---

## Available Helpers

### Environment Helpers (`env_helper.go`)
| Function | Description | Example |
|----------|-------------|---------|
| `GetEnv(key, default)` | Get environment variable with default | `helpers.GetEnv("APP_PORT", "8080")` |
| `GetEnvInt(key, default)` | Get environment variable as int | `helpers.GetEnvInt("JWT_EXPIRATION", 24)` |

### Error Sentinels (`error_helper.go`)
```go
helpers.ErrNotFound          // "record not found"
helpers.ErrInvalidToken      // "invalid token"
helpers.ErrExpiredToken      // "token expired"
helpers.ErrInvalidCredential // "invalid email or password"
helpers.ErrUserExists        // "user already exists"
helpers.ErrInvalidEmail      // "invalid email address"
helpers.ErrTokenExpired      // "reset token has expired"
helpers.ErrTokenUsed         // "reset token has already been used"
helpers.ErrTokenInvalid      // "invalid reset token"
```

### Crypto Helpers (`crypto_helper.go`)
```go
helpers.GenerateRandomString(32)     // secure random string
helpers.HashString("value")          // bcrypt hash
helpers.VerifyString("value", hash)  // verify bcrypt
```

### Response Helpers (`response_helper.go`)
```go
helpers.OK(c, "message", data)                    // 200
helpers.Created(c, "message", data)               // 201
helpers.BadRequest(c, "message")                  // 400 (JSON syntax, invalid params)
helpers.Unauthorized(c, "message")                // 401
helpers.Forbidden(c, "message")                   // 403
helpers.NotFound(c, "message")                    // 404
helpers.InternalServerError(c, "message")         // 500
helpers.ValidationErrorWithMap(c, errorsMap)      // 422 (validation errors map)
helpers.ValidationErrorWithField(c, field, msg)   // 422 (single field error)
```

**Response format:**
```json
{"code": 200, "message": "...", "data": {}}
```

**Validation error format:**
```json
{"code": 422, "message": "The given data was invalid.", "errors": {"field": ["message"]}}
```

**Validation pattern in handlers:**
```go
if err := c.BindJSON(&req); err != nil {
    helpers.BadRequest(c, "Invalid JSON payload")
    return
}
if err := h.Validate.Struct(req); err != nil {
    helpers.ValidationErrorWithMap(c, h.getErrorsMap(err))
    return
}
```

**Field-related business errors (422):**
```go
// Examples: email already exists, password mismatch, etc.
if err == helpers.ErrUserExists {
    helpers.ValidationErrorWithField(c, "email", "Email already exists")
    return
}
```

---

## Available Clients (injected into Services)

### Redis — `s.RedisClient`

> 🔴 Always check `s.RedisClient.IsCacheAvailable()` first. Redis errors **MUST NOT** cause operations to fail.

```go
if s.RedisClient.IsCacheAvailable() {
    s.RedisClient.SetJSON("key", value, time.Hour*24)
    s.RedisClient.GetJSON("key", &dest)
    s.RedisClient.Delete("key")
}
```

**Recommended Key Pattern:**
| Pattern | Example | Description |
|---------|--------|-----------|
| `session:{userID}` | `session:1` | User session cache |

### Email — `s.EmailClient`

```go
s.EmailClient.IsConfigured()
s.EmailClient.SendResetPasswordEmail(email, token, resetURL)
s.EmailClient.SendEmail(email.EmailRequest{To: []string{...}, Subject: "...", Body: "..."})
```

**EmailRequest DTO:**
```go
type EmailRequest struct {
    To      []string  // Recipients
    Subject string    // Email subject
    Body    string    // HTML body
    CC      []string  // CC recipients
    BCC     []string  // BCC recipients
}
```

---

## Access — Permission & Role Checking (`s.Access`)

> 🔐 Access helper is injected into Services via DI Container. Also available in routes as `container.Access` for middleware.

**3-Tier Caching:**
| Layer | Source | Behavior |
|-------|--------|----------|
| **L1** | Local in-memory (`sync.RWMutex` map) | Fastest, checked first |
| **L2** | Redis (`session:{userID}`) | Checked if L1 miss |
| **L3** | Database (user.Roles.Permissions) | Fallback if L1+L2 miss, then caches result to L1+L2 |

**Methods:**
| Method | Description | Example |
|--------|-------------|---------|
| `HasPermission(ctx, permissions...) bool` | Check if user has **ANY** of the specified permissions | `s.Access.HasPermission(ctx, "user.delete", "user.admin")` |
| `HasRole(ctx, role string) bool` | Check if user has the specified role | `s.Access.HasRole(ctx, "admin")` |
| `Invalidate(userID uint)` | Clear cached access data for a user | `s.Access.Invalidate(userID)` |

**Usage in Services:**
```go
func (s *Services) UserDelete(ctx context.Context, id uint) error {
    if !s.Access.HasPermission(ctx, "user.delete") {
        return helpers.ErrForbidden
    }
    // ... delete logic
}
```

**Usage in Routes (Middleware):**
```go
func RegisterUserRoutes(r *gin.RouterGroup, handlers *handlers.Handlers, acc *helpers.Access) {
    users := r.Group("/users")
    {
        users.POST("", middleware.RequirePermission(acc, "user.create"), handlers.UserCreate)
        users.GET("", middleware.RequirePermission(acc, "user.index"), handlers.UserGetAll)
        users.GET("/:id", middleware.RequirePermission(acc, "user.index"), handlers.UserGetByID)
        users.PUT("/:id", middleware.RequirePermission(acc, "user.edit"), handlers.UserUpdate)
        users.DELETE("/:id", middleware.RequirePermission(acc, "user.delete"), handlers.UserDelete)
    }
}
```

**Cache Invalidation:**
Call `s.Access.Invalidate(userID)` when user's roles/permissions change.

**Permission Naming Convention:** `{resource}.{action}`
| Permission | Description |
|------------|-------------|
| `user.index` | View users list |
| `user.create` | Create new user |
| `user.edit` | Update user |
| `user.delete` | Delete user |

---

## Creating New Helpers

> 💡 **Rule of Thumb**: If a function can be used in more than 1 place, make it a helper!

**Criteria for creating new helpers:**
1. ✅ General operation (not specific business logic)
2. ✅ Reusable across multiple services
3. ✅ No dependencies on repo, db, etc
4. ✅ Pure function (input → output, no side effects)

**❌ DO NOT create helpers if:**
- ❌ Only used in 1 service
- ❌ Contains specific business logic
- ❌ Requires dependency injection (db, repo, etc)
- ❌ Has side effects (write to db, send email, etc)

---

## Services Struct (Dependencies)

```go
type Services struct {
    repo        *repositories.Repositories  // Access repos: s.repo.User, s.repo.Permission, etc.
    RedisClient *database.RedisCache        // Redis cache client
    EmailClient *email.EmailClient          // Email client
    Access      *helpers.Access             // Permission/role checker
    Logger      *helpers.Logger             // Logger
    cfg         *JWTConfig                  // JWT config
}
```

## Handlers Struct

```go
type Handlers struct {
    svcs     *services.Services      // Access services: h.svcs.PermissionCreate(...)
    Validate *validator.Validate     // Validator instance with translator
    trans    ut.Translator           // Translator for error messages
}
```

---

## Anti-Patterns (NEVER DO)

| ❌ WRONG | ✅ CORRECT |
|----------|-----------|
| `type PermissionService struct` | Method on `(s *Services)` |
| `func (s *Services) CreatePermission(...)` | `func (s *Services) PermissionCreate(...)` |
| `s.repo.Permission.FindByID(s.repo.Permission.DB, id)` | `s.repo.Permission.FindByID(nil, id)` |
| Direct DB write without transaction | `s.repo.TxManager.WithinTransaction(...)` |
| Separate `CreateRequest` & `UpdateRequest` | Single `{Feature}Request` with `omitempty` |
| `FindByID` check before Update/Delete | Generic repo handles not found automatically |
| `c.ShouldBindJSON(&req)` + `helpers.ValidationError(c, err)` | `c.BindJSON(&req)` + `h.Validate.Struct(req)` + `helpers.ValidationErrorWithMap(c, h.getErrorsMap(err))` |
| `helpers.BadRequest(c, "Email already exists")` | `helpers.ValidationErrorWithField(c, "email", "Email already exists")` |
| File operations in handler layer | File move/delete in service layer |
| Modifying `00_generic.go` or `00_transaction.go` | Create custom repository file |

---

## Pre-Implementation Checklist

- [ ] Model has `TableName()` method
- [ ] DTO structs use feature prefix (`{Feature}DTO`, `{Feature}Request`)
- [ ] Request DTO merged (Create & Update in 1 struct)
- [ ] Service functions: `{Feature}{Action}` on `(s *Services)`
- [ ] Handler functions: `{Feature}{Action}` on `(h *Handlers)`
- [ ] Write operations use `TxManager.WithinTransaction()`
- [ ] Read operations use `nil` parameter
- [ ] CREATE/UPDATE/DELETE have logging
- [ ] Repository registered in `00_repository.go`
- [ ] Routes registered in `cmd/api/main.go`
- [ ] Handler uses `c.BindJSON()` + `h.Validate.Struct()` (NOT `ShouldBindJSON`)
- [ ] Field-related business errors use `ValidationErrorWithField()` (422)
- [ ] Build success: `go build ./...`
- [ ] Vet clean: `go vet ./...`
