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

## CRUD Implementation Flow (Step-by-Step)

When creating a new feature, follow this exact order:

### Step 1: Model — `internal/models/{feature}.go`
```go
package models

import (
    "time"
    "gorm.io/gorm"
)

type Permission struct {
    ID          uint           `gorm:"primaryKey" json:"id"`
    Name        string         `gorm:"uniqueIndex;size:100;not null" json:"name"`
    Description *string        `gorm:"size:255" json:"description"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Permission) TableName() string {
    return "permissions"
}
```
- Use `*string` for nullable fields
- `DeletedAt gorm.DeletedAt` for soft delete
- `TableName()` is MANDATORY

### Step 2: DTOs — `internal/dtos/{feature}_dto.go`
```go
package dtos

import "github.com/reshap0318/go-boilerplate/internal/models"

// ONE request struct for both Create & Update
type PermissionRequest struct {
    Name        string  `json:"name" binding:"required,min=3,max=100"`
    Description *string `json:"description" binding:"omitempty,max=255"`
}

type PermissionDTO struct {
    ID          uint    `json:"id"`
    Name        string  `json:"name"`
    Description *string `json:"description"`
}

func ToPermissionDTO(p *models.Permission) PermissionDTO {
    return PermissionDTO{
        ID:          p.ID,
        Name:        p.Name,
        Description: p.Description,
    }
}

func ToPermissionDTOList(permissions []models.Permission) []PermissionDTO {
    result := make([]PermissionDTO, len(permissions))
    for i, p := range permissions {
        result[i] = ToPermissionDTO(&p)
    }
    return result
}
```
- Prefix ALL structs with feature name: `PermissionDTO`, `PermissionRequest`
- Merge Create & Update into ONE `{Feature}Request` struct
- Converter functions: `To{Feature}DTO()` and `To{Feature}DTOList()`

### Step 3: Repository — `internal/repositories/{feature}_repository.go`
```go
package repositories

import (
    "github.com/reshap0318/go-boilerplate/internal/models"
    "gorm.io/gorm"
)

type PermissionRepository struct {
    *GenericRepository[models.Permission]
}

func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
    return &PermissionRepository{
        GenericRepository: NewGenericRepository(db, &models.Permission{}),
    }
}
```
- ONLY extend `GenericRepository[Model]`
- Add custom methods ONLY for complex JOINs, subqueries, aggregations
- ❌ DO NOT re-implement FindByID, Create, Update, Delete — they exist in GenericRepository

### Step 4: Register Repository — `internal/repositories/00_repository.go`

Add field to `Repositories` struct + initialize in `NewRepositories()`:
```go
type Repositories struct {
    TxManager  *TransactionManager
    User       *UserRepository
    Permission *PermissionRepository  // ← ADD
}

func NewRepositories(db *gorm.DB) (*Repositories, error) {
    permissionRepo := NewPermissionRepository(db)  // ← INIT
    return &Repositories{
        Permission: permissionRepo,  // ← REGISTER
    }, nil
}
```

### Step 5: Service — `internal/services/{feature}_service.go`
```go
package services

import (
    "context"
    "gorm.io/gorm"
    "github.com/reshap0318/go-boilerplate/internal/dtos"
    "github.com/reshap0318/go-boilerplate/internal/helpers"
    "github.com/reshap0318/go-boilerplate/internal/models"
)

// CREATE — transaction + logging
func (s *Services) PermissionCreate(ctx context.Context, req dtos.PermissionRequest) (*dtos.PermissionDTO, error) {
    s.Logger.LogStart("PermissionCreate", "Creating permission: %s", req.Name)

    permission := &models.Permission{Name: req.Name, Description: req.Description}

    var result *models.Permission
    if err := s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
        var err error
        result, err = s.repo.Permission.Create(tx, permission)
        return err
    }); err != nil {
        s.Logger.LogEndWithError("PermissionCreate", "Failed: %v", err)
        return nil, err
    }

    dto := dtos.ToPermissionDTO(result)
    s.Logger.LogEnd("PermissionCreate", "Created: %s (ID: %d)", dto.Name, dto.ID)
    return &dto, nil
}

// GET ALL — nil, no logging
func (s *Services) PermissionGetAll(ctx context.Context) ([]dtos.PermissionDTO, error) {
    permissions, err := s.repo.Permission.FindAll(nil)
    if err != nil {
        return nil, err
    }
    return dtos.ToPermissionDTOList(permissions), nil
}

// GET BY ID — nil, no logging, return ErrNotFound
func (s *Services) PermissionGetByID(ctx context.Context, id uint) (*dtos.PermissionDTO, error) {
    permission, err := s.repo.Permission.FindByID(nil, id)
    if err != nil {
        return nil, helpers.ErrNotFound
    }
    dto := dtos.ToPermissionDTO(permission)
    return &dto, nil
}

// UPDATE — transaction + logging
func (s *Services) PermissionUpdate(ctx context.Context, id uint, req dtos.PermissionRequest) (*dtos.PermissionDTO, error) {
    s.Logger.LogStart("PermissionUpdate", "Updating permission ID: %d", id)

    permission := &models.Permission{ID: id}
    if req.Name != "" { permission.Name = req.Name }
    if req.Description != nil { permission.Description = req.Description }

    var result *models.Permission
    if err := s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
        var err error
        result, err = s.repo.Permission.Update(tx, &models.Permission{ID: id}, permission)
        return err
    }); err != nil {
        s.Logger.LogEndWithError("PermissionUpdate", "Failed: %v", err)
        return nil, err
    }

    dto := dtos.ToPermissionDTO(result)
    s.Logger.LogEnd("PermissionUpdate", "Updated: %s (ID: %d)", dto.Name, dto.ID)
    return &dto, nil
}

// DELETE — transaction + logging
func (s *Services) PermissionDelete(ctx context.Context, id uint) error {
    s.Logger.LogStart("PermissionDelete", "Deleting permission ID: %d", id)

    if err := s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
        _, err := s.repo.Permission.Delete(tx, id)
        return err
    }); err != nil {
        s.Logger.LogEndWithError("PermissionDelete", "Failed: %v", err)
        return err
    }

    s.Logger.LogEnd("PermissionDelete", "Deleted: ID: %d", id)
    return nil
}
```

### Step 6: Handler — `internal/handlers/{feature}_handler.go`
```go
package handlers

import (
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/reshap0318/go-boilerplate/internal/dtos"
    "github.com/reshap0318/go-boilerplate/internal/helpers"
)

func (h *Handlers) PermissionCreate(c *gin.Context) {
    var req dtos.PermissionRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        helpers.ValidationError(c, err)
        return
    }
    dto, err := h.svcs.PermissionCreate(c.Request.Context(), req)
    if err != nil {
        helpers.InternalServerError(c, "Failed to create permission")
        return
    }
    helpers.Created(c, "Permission created successfully", dto)
}

func (h *Handlers) PermissionGetAll(c *gin.Context) {
    dtos, err := h.svcs.PermissionGetAll(c.Request.Context())
    if err != nil {
        helpers.InternalServerError(c, "Failed to fetch permissions")
        return
    }
    helpers.OK(c, "Permissions fetched successfully", dtos)
}

func (h *Handlers) PermissionGetByID(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        helpers.BadRequest(c, "Invalid permission ID")
        return
    }
    dto, err := h.svcs.PermissionGetByID(c.Request.Context(), uint(id))
    if err != nil {
        helpers.NotFound(c, "Permission not found")
        return
    }
    helpers.OK(c, "Permission fetched successfully", dto)
}

func (h *Handlers) PermissionUpdate(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        helpers.BadRequest(c, "Invalid permission ID")
        return
    }
    var req dtos.PermissionRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        helpers.ValidationError(c, err)
        return
    }
    dto, err := h.svcs.PermissionUpdate(c.Request.Context(), uint(id), req)
    if err != nil {
        helpers.NotFound(c, "Permission not found")
        return
    }
    helpers.OK(c, "Permission updated successfully", dto)
}

func (h *Handlers) PermissionDelete(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        helpers.BadRequest(c, "Invalid permission ID")
        return
    }
    err = h.svcs.PermissionDelete(c.Request.Context(), uint(id))
    if err != nil {
        helpers.NotFound(c, "Permission not found")
        return
    }
    helpers.OK(c, "Permission deleted successfully", nil)
}
```
- Use `helpers.ValidationError(c, err)` for binding errors (handles both validation & JSON errors)
- Use `helpers.Created()` for POST, `helpers.OK()` for GET/PUT/DELETE
- Parse ID with `strconv.ParseUint(c.Param("id"), 10, 64)`

### Step 7: Routes — `internal/routes/{feature}_route.go`
```go
package routes

import (
    "github.com/gin-gonic/gin"
    "github.com/reshap0318/go-boilerplate/internal/handlers"
)

func RegisterPermissionRoutes(r *gin.RouterGroup, handlers *handlers.Handlers) {
    permissions := r.Group("/permissions")
    {
        permissions.POST("", handlers.PermissionCreate)
        permissions.GET("", handlers.PermissionGetAll)
        permissions.GET("/:id", handlers.PermissionGetByID)
        permissions.PUT("/:id", handlers.PermissionUpdate)
        permissions.DELETE("/:id", handlers.PermissionDelete)
    }
}
```

### Step 8: Register Routes — `cmd/api/main.go`

Add to the `protected` group (after JWT middleware):
```go
protected := apiGroup.Group("")
protected.Use(middleware.JWTAuth(container.Services))
{
    routes.RegisterPermissionRoutes(protected, container.Handlers)  // ← ADD
}
```

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
    IncludeDeleted bool     // Include soft-deleted records
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

### Response Helpers (`helpers`)
```go
helpers.OK(c, "message", data)            // 200
helpers.Created(c, "message", data)       // 201
helpers.BadRequest(c, "message")          // 400
helpers.Unauthorized(c, "message")        // 401
helpers.Forbidden(c, "message")           // 403
helpers.NotFound(c, "message")            // 404
helpers.InternalServerError(c, "message") // 500
helpers.ValidationError(c, err)           // 422 (validation) or 400 (JSON syntax)
```

Response format:
```json
{"code": 200, "message": "...", "data": {}}
```

Validation error format:
```json
{"code": 422, "message": "The given data was invalid.", "errors": {"field": ["message"]}}
```

### Error Sentinels (`helpers`)
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

### Environment Helpers
```go
helpers.GetEnv("KEY", "default")     // string
helpers.GetEnvInt("KEY", 42)         // int
```

### Crypto Helpers
```go
helpers.GenerateRandomString(32)     // secure random string
helpers.HashString("value")          // bcrypt hash
helpers.VerifyString("value", hash)  // verify bcrypt
```

---

## Available Clients (injected into Services)

### Redis — `s.RedisClient`
```go
if s.RedisClient.IsCacheAvailable() {
    s.RedisClient.SetJSON("key", value, time.Hour*24)
    s.RedisClient.GetJSON("key", &dest)
    s.RedisClient.Delete("key")
}
```
⚠️ Always check `IsCacheAvailable()` first. Redis errors must NOT cause operations to fail.

### Email — `s.EmailClient`
```go
s.EmailClient.IsConfigured()
s.EmailClient.SendResetPasswordEmail(email, token, resetURL)
s.EmailClient.SendEmail(email.EmailRequest{To: []string{...}, Subject: "...", Body: "..."})
```

---

## Services Struct (Dependencies)

```go
type Services struct {
    repo        *repositories.Repositories  // Access repos: s.repo.User, s.repo.Permission, etc.
    RedisClient *database.RedisCache        // Redis cache client
    EmailClient *email.EmailClient          // Email client
    Logger      *helpers.Logger             // Logger
    cfg         *JWTConfig                  // JWT config
}
```

## Handlers Struct

```go
type Handlers struct {
    svcs *services.Services  // Access services: h.svcs.PermissionCreate(...)
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
| Separate `CreateRequest` & `UpdateRequest` | Single `PermissionRequest` with `omitempty` |
| `FindByID` check before Update/Delete | Generic repo handles not found automatically |
| `helpers.BadRequest(c, err.Error())` for binding | `helpers.ValidationError(c, err)` |
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
- [ ] Build success: `go build ./...`
- [ ] Vet clean: `go vet ./...`
