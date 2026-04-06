# CRUD Implementation Flow

> ⚠️ **READ THIS** before creating any new CRUD feature. All rules here are **MANDATORY**.

---

## 📋 Step-by-Step CRUD Implementation

Every new CRUD feature MUST follow this order:

```
1. Model          → models/feature.go
2. DTOs           → dtos/feature_dto.go
3. Repository     → repositories/feature_repository.go
4. Register Repo  → repositories/00_repository.go
5. Service        → services/feature_service.go
6. Handler        → handlers/feature_handler.go
7. Routes         → routes/feature_route.go
8. Register Routes→ cmd/api/main.go
```

---

## ⚠️ Naming Convention (MANDATORY)

### 1. DTO Variables — Prefix Feature Name

All structs in DTOs MUST use the feature name as prefix.

```go
// ✅ CORRECT
type PermissionDTO struct { ... }
type PermissionRequest struct { ... }
type RoleDTO struct { ... }
type UserRoleRequest struct { ... }

// ❌ WRONG — No prefix
type DTO struct { ... }
type Request struct { ... }
```

**Rules:**
- Response DTO → `{Feature}DTO`
- Request DTO → `{Feature}Request` (merged Create & Update)
- List converter → `To{Feature}DTO()` and `To{Feature}DTOList()`

### 2. Service Functions — Prefix Feature Name

All functions in the Services struct MUST be prefixed with the feature name.

**Pattern:** `{Feature}{Action}` — Feature FIRST, then Action.

```go
// ✅ CORRECT — Feature name FIRST
func (s *Services) PermissionCreate(ctx context.Context, req dtos.PermissionRequest) (*dtos.PermissionDTO, error)
func (s *Services) PermissionGetAll(ctx context.Context) ([]dtos.PermissionDTO, error)
func (s *Services) PermissionGetByID(ctx context.Context, id uint) (*dtos.PermissionDTO, error)
func (s *Services) PermissionUpdate(ctx context.Context, id uint, req dtos.PermissionRequest) (*dtos.PermissionDTO, error)
func (s *Services) PermissionDelete(ctx context.Context, id uint) error

// ❌ WRONG — Action first (CreatePermission, GetAllPermissions, etc.)
func (s *Services) CreatePermission(...) { }
func (s *Services) GetAllPermissions(...) { }

// ❌ WRONG — No prefix
func (s *Services) Create(...) { }
func (s *Services) GetByID(...) { }
```

### 3. Handler Functions — Prefix Feature Name

All functions in the Handlers struct MUST be prefixed with the feature name.

**Pattern:** `{Feature}{Action}` — Feature FIRST, then Action.

```go
// ✅ CORRECT — Feature name FIRST
func (h *Handlers) PermissionCreate(c *gin.Context)
func (h *Handlers) PermissionGetAll(c *gin.Context)
func (h *Handlers) PermissionGetByID(c *gin.Context)
func (h *Handlers) PermissionUpdate(c *gin.Context)
func (h *Handlers) PermissionDelete(c *gin.Context)

// ❌ WRONG — No prefix
func (h *Handlers) Create(c *gin.Context)
func (h *Handlers) GetByID(c *gin.Context)
```

**Naming pattern for both Services & Handlers:** `{Feature}{Action}` where Action = Create, GetAll, GetByID, Update, Delete

---

## 📁 Step 1: Model

**File:** `internal/models/feature.go`

```go
package models

import (
    "time"
    "gorm.io/gorm"
)

type Permission struct {
    ID          uint           `gorm:"primaryKey" json:"id"`
    Name        string         `gorm:"uniqueIndex;size:100;not null" json:"name"`
    Description *string        `gorm:"size:255" json:"description"`    // nullable
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`                // soft delete
}

// ⚠️ MANDATORY
func (Permission) TableName() string {
    return "permissions"
}
```

**Key points:**
- ✅ `TableName()` is **MANDATORY**
- ✅ `DeletedAt gorm.DeletedAt` for soft delete
- ✅ Use pointer (`*string`) for nullable fields
- ✅ Consistent JSON tags on all fields

---

## 📁 Step 2: DTOs

**File:** `internal/dtos/feature_dto.go`

```go
package dtos

import "github.com/reshap0318/go-boilerplate/internal/models"

// ⚠️ MANDATORY: Prefix with feature name
// ⚠️ MANDATORY: Merge Create & Update into 1 struct
type PermissionRequest struct {
    Name        string  `json:"name" binding:"required,min=3,max=100"`
    Description *string `json:"description" binding:"omitempty,max=255"`
}

// ⚠️ MANDATORY: Prefix with feature name
type PermissionDTO struct {
    ID          uint    `json:"id"`
    Name        string  `json:"name"`
    Description *string `json:"description"`
}

// ⚠️ MANDATORY: Converter with feature prefix
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

**Key points:**
- ✅ **ONE** request struct for both Create & Update (`PermissionRequest`)
- ✅ DTO struct for response (`PermissionDTO`)
- ✅ Converter functions `To{Feature}DTO()` and `To{Feature}DTOList()`
- ✅ Request fields: `required` for Create, `omitempty` for Update
- ✅ Import `models` package

---

## 📁 Step 3: Repository

**File:** `internal/repositories/feature_repository.go`

```go
package repositories

import (
    "github.com/reshap0318/go-boilerplate/internal/models"
    "gorm.io/gorm"
)

// ⚠️ MANDATORY: Extend GenericRepository
type PermissionRepository struct {
    *GenericRepository[models.Permission]
}

func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
    return &PermissionRepository{
        GenericRepository: NewGenericRepository(db, &models.Permission{}),
    }
}
```

**Key points:**
- ✅ **ONLY** extend `GenericRepository[Model]`
- ✅ Custom methods **ONLY** for complex queries (JOINs, subqueries, aggregations)
- ❌ DO NOT duplicate methods already in GenericRepository (Create, FindByID, FindAll, Update, Delete, etc.)
- ❌ DO NOT modify `00_generic.go` or `00_transaction.go`

---

## 📁 Step 4: Register Repository

**File:** `internal/repositories/00_repository.go`

```go
type Repositories struct {
    TxManager     *TransactionManager
    User          *UserRepository
    PasswordReset *PasswordResetRepository
    Permission    *PermissionRepository  // ← Add this
}

func NewRepositories(db *gorm.DB) (*Repositories, error) {
    txManager := NewTransactionManager(db)
    userRepo := NewUserRepository(db)
    passwordResetRepo := NewPasswordResetRepository(db)
    permissionRepo := NewPermissionRepository(db)  // ← Initialize

    return &Repositories{
        TxManager:     txManager,
        User:          userRepo,
        PasswordReset: passwordResetRepo,
        Permission:    permissionRepo,  // ← Register
    }, nil
}
```

---

## 📁 Step 5: Service

**File:** `internal/services/feature_service.go`

### Write Operations (Create/Update/Delete) — MANDATORY Transaction

```go
// ⚠️ MANDATORY: Feature name FIRST ({Feature}{Action})
// ⚠️ MANDATORY: Use TxManager for write operations
func (s *Services) PermissionCreate(ctx context.Context, req dtos.PermissionRequest) (*dtos.PermissionDTO, error) {
    s.Logger.LogStart("PermissionCreate", "Creating permission: %s", req.Name)

    permission := &models.Permission{
        Name:        req.Name,
        Description: req.Description,
    }

    var result *models.Permission
    if err := s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
        var err error
        result, err = s.repo.Permission.Create(tx, permission)
        return err
    }); err != nil {
        s.Logger.LogEndWithError("PermissionCreate", "Failed to create permission: %v", err)
        return nil, err
    }

    dto := dtos.ToPermissionDTO(result)
    s.Logger.LogEnd("PermissionCreate", "Permission created: %s (ID: %d)", dto.Name, dto.ID)
    return &dto, nil
}
```

### Read Operations (Get/Find) — MANDATORY use `nil`

```go
// ⚠️ MANDATORY: Feature name FIRST ({Feature}{Action})
// ⚠️ MANDATORY: Use nil for read operations (NOT s.repo.Permission.DB)
func (s *Services) PermissionGetAll(ctx context.Context) ([]dtos.PermissionDTO, error) {
    permissions, err := s.repo.Permission.FindAll(nil)  // ← nil, NOT .DB
    if err != nil {
        return nil, err
    }
    return dtos.ToPermissionDTOList(permissions), nil
}

// ⚠️ Simple GET does NOT need logging
func (s *Services) PermissionGetByID(ctx context.Context, id uint) (*dtos.PermissionDTO, error) {
    permission, err := s.repo.Permission.FindByID(nil, id)  // ← nil, NOT .DB
    if err != nil {
        return nil, helpers.ErrNotFound
    }
    return dtos.ToPermissionDTO(permission), nil
}
```

### Update with Transaction

```go
func (s *Services) PermissionUpdate(ctx context.Context, id uint, req dtos.PermissionRequest) (*dtos.PermissionDTO, error) {
    s.Logger.LogStart("PermissionUpdate", "Updating permission ID: %d", id)

    permission := &models.Permission{
        ID: id,
    }
    if req.Name != "" {
        permission.Name = req.Name
    }
    if req.Description != nil {
        permission.Description = req.Description
    }

    var result *models.Permission
    if err := s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
        var err error
        result, err = s.repo.Permission.Update(tx, &models.Permission{ID: id}, permission)
        return err
    }); err != nil {
        s.Logger.LogEndWithError("PermissionUpdate", "Failed to update permission: %v", err)
        return nil, err
    }

    dto := dtos.ToPermissionDTO(result)
    s.Logger.LogEnd("PermissionUpdate", "Permission updated: %s (ID: %d)", dto.Name, dto.ID)
    return &dto, nil
}
```

### Delete with Transaction

```go
func (s *Services) PermissionDelete(ctx context.Context, id uint) error {
    s.Logger.LogStart("PermissionDelete", "Deleting permission ID: %d", id)

    if err := s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
        _, err := s.repo.Permission.Delete(tx, id)
        return err
    }); err != nil {
        s.Logger.LogEndWithError("PermissionDelete", "Failed to delete permission: %v", err)
        return err
    }

    s.Logger.LogEnd("PermissionDelete", "Permission deleted: ID: %d", id)
    return nil
}
```

**Service key points:**
- ✅ **Write** (Create/Update/Delete) → `s.repo.TxManager.WithinTransaction()`
- ✅ **Read** (Get/Find) → `nil` parameter (NOT `s.repo.Feature.DB`)
- ✅ **Function name** → Feature FIRST (`PermissionCreate`, `PermissionGetAll`, etc.)
- ✅ **Logging** → Create/Update/Delete MUST log, Simple GET does NOT need logging
- ✅ **Error handling** → use `helpers.ErrNotFound` for record not found
- ✅ **Single struct** → all methods on `(s *Services)`, do NOT create separate structs
- ✅ **Update/Delete** → NO FindByID check before operation (generic repo handles not found)

---

## 📁 Step 6: Handler

**File:** `internal/handlers/feature_handler.go`

```go
package handlers

import (
    "strconv"

    "github.com/gin-gonic/gin"

    "github.com/reshap0318/go-boilerplate/internal/dtos"
    "github.com/reshap0318/go-boilerplate/internal/helpers"
)

// ⚠️ MANDATORY: Feature name FIRST — {Feature}{Action}
func (h *Handlers) PermissionCreate(c *gin.Context) {
    var req dtos.PermissionRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        helpers.BadRequest(c, err.Error())
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
        helpers.BadRequest(c, err.Error())
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

**Handler key points:**
- ✅ **Function name** → `{Feature}{Action}` (`PermissionCreate`, `PermissionGetAll`, etc.)
- ✅ **Parse ID** → `strconv.ParseUint` with error handling
- ✅ **Bind JSON** → `c.ShouldBindJSON()` with validation
- ✅ **Response helper** → `helpers.Created`, `helpers.OK`, `helpers.BadRequest`, etc.
- ✅ **Single struct** → all methods on `(h *Handlers)`, do NOT create separate structs
- ✅ **Service calls** → MUST match service function name (`h.svcs.PermissionCreate`, NOT `h.svcs.CreatePermission`)

---

## 📁 Step 7: Routes

**File:** `internal/routes/feature_route.go`

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

---

## 📁 Step 8: Register Routes

**File:** `cmd/api/main.go`

Add to the `protected` group (after JWT middleware):

```go
protected := apiGroup.Group("")
protected.Use(middleware.JWTAuth(container.Services))
{
    routes.RegisterAuthProtectedRoutes(protected, container.Handlers)
    routes.RegisterPermissionRoutes(protected, container.Handlers)  // ← Add this
}
```

---

## 🚫 Anti-Patterns (DO NOT DO)

| ❌ WRONG | ✅ CORRECT |
|----------|------------|
| `type DTO struct` | `type PermissionDTO struct` |
| `type Request struct` | `type PermissionRequest struct` |
| `func (s *Services) CreatePermission(...)` | `func (s *Services) PermissionCreate(...)` |
| `func (s *Services) GetAllPermissions(...)` | `func (s *Services) PermissionGetAll(...)` |
| `func (h *Handlers) GetByID(...)` | `func (h *Handlers) PermissionGetByID(...)` |
| `h.svcs.CreatePermission(...)` | `h.svcs.PermissionCreate(...)` |
| `s.repo.Permission.FindByID(s.repo.Permission.DB, id)` | `s.repo.Permission.FindByID(nil, id)` |
| Direct DB write without transaction | `s.repo.TxManager.WithinTransaction(...)` |
| Separate struct (`type PermissionService struct`) | Method on `(s *Services)` |
| Separate CreateRequest & UpdateRequest | Single `PermissionRequest` with `omitempty` |
| FindByID before Update/Delete | Generic repo handles not found automatically |

---

## ✅ Pre-Push Checklist

- [ ] Model has `TableName()` method
- [ ] DTO variables use feature prefix
- [ ] Request DTO merged (Create & Update into 1)
- [ ] Service functions use feature prefix ({Feature}{Action})
- [ ] Handler functions use feature prefix (`{Feature}{Action}`)
- [ ] Write operations use `TxManager.WithinTransaction()`
- [ ] Read operations use `nil` (NOT `.DB`)
- [ ] Logging on Create/Update/Delete
- [ ] Repository registered in `00_repository.go`
- [ ] Routes registered in `cmd/api/main.go`
- [ ] Build success (`go build ./...`)
- [ ] Vet clean (`go vet ./...`)
