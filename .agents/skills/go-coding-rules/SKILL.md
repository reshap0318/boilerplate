---
name: go-coding-rules
description: Go project coding rules and conventions — see references/ for code examples
---

## When to use me

Use this skill when:

- Writing or modifying Go code — follow rules below, see `references/` for examples
- Adding features to existing codebase — ensure consistency with project conventions
- Reviewing code — check for anti-patterns and rule violations
- Need reference for helpers, packages, or GenericRepository methods

---

# Go Coding Rules

> **READ THIS ENTIRELY before writing any code.** These rules are MANDATORY.

## Project Identity

| Item        | Value                              |
| ----------- | ---------------------------------- |
| Module      | `github.com/reshap0318/go-project` |
| Go Version  | 1.25.0+                            |
| Framework   | Gin, GORM, JWT, Redis, bcrypt      |
| Database    | MySQL (default) / PostgreSQL       |
| Entry Point | `cmd/api/main.go`                  |

## Architecture

```
Routes → Handlers → Services → Repositories → Database/Redis
              ↑
         Middleware (JWT, CORS, Rate Limit)
```

All layers use **Dependency Injection** via `internal/di/container.go`.

## Code Examples

All examples are in `references/` folder:

- `references/model.go` — Model with TableName()
- `references/dto.go` — Request/Response DTOs
- `references/repository.go` — Repository registration + custom repo
- `references/service.go` — CRUD methods with logging + notification
- `references/handler.go` — HTTP handlers with validation
- `references/route.go` — Routes with permission middleware

---

## CRITICAL RULES

### 1. NEVER modify

- `internal/repositories/00_generic.go`
- `internal/repositories/00_transaction.go`

### 2. Models MUST have `TableName()`

```go
func (Permission) TableName() string { return "permissions" }
```

### 3. Single Struct Pattern

- ALL service methods → `func (s *Services) ...`
- ALL handler methods → `func (h *Handlers) ...`
- NEVER create `type PermissionService struct`

### 4. Naming: `{Feature}{Action}` — Feature FIRST

```go
// CORRECT
func (s *Services) PermissionCreate(...)
func (h *Handlers) PermissionGetByID(...)

// WRONG
func (s *Services) CreatePermission(...)
```

### 5. Write Operations → MUST use Transaction

```go
s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
    result, err = s.repo.Permission.Create(tx, permission)
    return err
})
```

### 6. Read Operations → MUST use `nil`

```go
s.repo.Permission.FindByID(nil, id)  // CORRECT
s.repo.Permission.FindByID(s.repo.Permission.DB, id)  // WRONG
```

### 7. Read Operations → Return DTOs, NOT Models

```go
// Convert models to DTOs using To{Feature}DTO()
func (s *Services) PermissionGetByID(ctx context.Context, id uint) (*dtos.PermissionDTO, error) {
    permission, err := s.repo.Permission.FindByID(nil, id, "Roles")
    if err != nil {
        return nil, err
    }
    return dtos.ToPermissionDTO(permission), nil
}
```

### 8. Paginated Read → Handle nil opts

```go
func (s *Services) PermissionGetAllPaginated(ctx context.Context, opts *repositories.QueryOptions) (*repositories.PagedResult[dtos.PermissionDTO], error) {
    if opts == nil {
        opts = &repositories.QueryOptions{}
    }
    // ... convert to DTOs, return PagedResult[DTO]
}
```

### 9. File Naming

| Layer      | Pattern                   | Example                    |
| ---------- | ------------------------- | -------------------------- |
| Model      | `{feature}.go`            | `permission.go`            |
| DTO        | `{feature}_dto.go`        | `permission_dto.go`        |
| Repository | `{feature}_repository.go` | `permission_repository.go` |
| Service    | `{feature}_service.go`    | `permission_service.go`    |
| Handler    | `{feature}_handler.go`    | `permission_handler.go`    |
| Route      | `{feature}_route.go`      | `permission_route.go`      |

### 10. DTO Field Naming — Request vs Response

**Request DTO** — use entity name without `_id` suffix:

```go
type RoleRequest struct {
    Permissions []uint `json:"permissions"` // NOT "permission_ids"
}

type UserRequest struct {
    Roles []uint `json:"roles"` // NOT "role_ids"
}
```

**Response DTO** — use `_id` suffix for foreign keys:

```go
type NotificationDTO struct {
    UserID uint `json:"user_id"` // keep _id in response
}
```

---

## Logging

Logger available via `s.Logger`. CREATE/UPDATE/DELETE MUST log. GET simple: NO log.

```go
s.Logger.LogStart("FuncName", "Message: %s", value)
s.Logger.LogStep("FuncName", "Step: %s", value)
s.Logger.LogEnd("FuncName", "Success: %s", value)
s.Logger.LogEndWithError("Func", "Error: %v", err)
s.Logger.LogError("FuncName", "Error: %v", err)
s.Logger.LogWarn("FuncName", "Warning: %s", value)
s.Logger.LogInfo("FuncName", "Info: %s", value)
```

See `references/service.go` for full examples.

---

## Notification

Every CREATE/UPDATE/DELETE MUST create notification via `s.NotificationCreate()`:

```go
_ = s.NotificationCreate(ctx, &services.NotificationCreateParams{
    Type:    "success",  // info, warning, success, error
    Title:   "Title",
    Message: "Message",
    Data:    map[string]interface{}{"id": result.ID, "name": result.Name}, // identifier data only
})
```

**Rules:**

- Errors MUST NOT fail the transaction → use `_ = s.NotificationCreate(...)`
- Created **AFTER** transaction completes — NOT inside transaction
- Data = identifier data only (id, name, status, etc) — NOT full object

See `references/service.go` for full examples.

---

## GenericRepository Methods

All repos extend `GenericRepository[T]`. Do NOT re-implement:

| Method                   | Signature                                                                        |
| ------------------------ | -------------------------------------------------------------------------------- |
| `FindByID`               | `(tx *gorm.DB, id uint, preloads ...string)`                                     |
| `FindByIDWithOpts`       | `(tx *gorm.DB, id uint, opts *QueryOptions) (*T, error)`                         |
| `Create`                 | `(tx *gorm.DB, request *T) (*T, error)`                                          |
| `CreateMany`             | `(tx *gorm.DB, request []T) error`                                               |
| `Update`                 | `(tx *gorm.DB, filter *T, update *T) (*T, error)`                                |
| `UpdateMap`              | `(tx *gorm.DB, filter *T, update map[string]interface{})`                        |
| `Delete`                 | `(tx *gorm.DB, id uint) (*T, error)`                                             |
| `FindAll`                | `(tx *gorm.DB, preloads ...string)`                                               |
| `FindAllWithOpts`        | `(tx *gorm.DB, opts *QueryOptions) (*PagedResult[T], error)`                     |
| `FindByField`            | `(tx *gorm.DB, filter *T, preloads ...string)`                                   |
| `FindByFieldWithOpts`    | `(tx *gorm.DB, filter *T, opts *QueryOptions) (*PagedResult[T], error)`          |
| `FindByFieldMap`         | `(tx *gorm.DB, filter map[string]interface{}, preloads ...string)`               |
| `FindByFieldMapWithOpts` | `(tx *gorm.DB, filter map[string]interface{}, opts *QueryOptions) (*PagedResult[T], error)` |
| `Count`                  | `(tx *gorm.DB) (int64, error)`                                                   |
| `Exists`                 | `(tx *gorm.DB, filter map[string]interface{}) (bool, error)`                     |
| `ExistsByField`          | `(tx *gorm.DB, filter *T) (bool, error)`                                         |

### QueryOptions Fields

```go
type QueryOptions struct {
    Page            int              // Page number (default: 1)
    PageSize        int              // Items per page (default: 10)
    SortBy          string           // Field to sort by
    Order           string           // "ASC" | "DESC" (default: "ASC")
    Preloads        []string         // Relations to preload
    Omits           []string         // Columns to omit from SELECT
    ConditionGroups []ConditionGroup // WHERE groups, AND-ed together
}
```

`ConditionGroups` — groups joined by AND; conditions within each group use group `Logic`:

```go
// Example: (name LIKE '%john%' OR email LIKE '%john%') AND (status = 1)
opts := &repositories.QueryOptions{
    ConditionGroups: []repositories.ConditionGroup{
        {
            Logic: "OR",
            Conditions: []repositories.QueryCondition{
                {Column: "name", Operator: "LIKE", Value: "%john%"},
                {Column: "email", Operator: "LIKE", Value: "%john%"},
            },
        },
        {
            Logic: "AND",
            Conditions: []repositories.QueryCondition{
                {Column: "status", Operator: "=", Value: 1},
            },
        },
    },
}
```

### Priority: Use Generic Methods First

**Prefer** using existing `GenericRepository` methods over creating custom repository methods.

**Guidelines:**

1. Check if generic method can solve the problem first (see table above)
2. Create custom method when generic methods cannot handle the requirement
3. Common scenarios that may need custom methods:
   - Complex JOIN queries
   - Aggregate functions (SUM, COUNT with GROUP BY)
   - Subqueries or nested conditions
   - Database-specific features not covered by generic methods

See `references/repository.go` for examples.

### Transaction Manager

```go
s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error { ... })
result, err := s.repo.TxManager.WithinTransactionWithResult(func(tx *gorm.DB) (interface{}, error) { ... })
```

---

## Helpers

### Error Sentinels

```go
helpers.ErrNotFound, helpers.ErrInvalidToken, helpers.ErrExpiredToken
helpers.ErrInvalidCredential, helpers.ErrUserExists, helpers.ErrInvalidEmail
helpers.ErrForbidden
```

### Custom Error — `CustomError`

Use `helpers.CustomError` to return custom HTTP status and message from service layer:

```go
// Default status 400 (Bad Request)
return &helpers.CustomError{
    Message: "ID tidak valid",
}

// With custom status
return &helpers.CustomError{
    Status:  http.StatusConflict, // 409
    Message: "Data sudah ada",
}
```

Handler handles it automatically via `HandleError()`:

```go
err := s.Services.SomeAction(id)
if helpers.HandleError(c, err, "Fallback message") {
    return
}
```

| Scenario | Status | Message |
|----------|--------|---------|
| `Status` not set | 400 | Custom message |
| `Status` set | Custom | Custom message |

### Response Helpers

```go
helpers.OK(c, "msg", data)               // 200
helpers.OKWithMetadata(c, "msg", result) // 200 — for paginated PagedResult
helpers.Created(c, "msg", data)          // 201
helpers.BadRequest(c, "msg")             // 400
helpers.Unauthorized(c, "msg")           // 401
helpers.Forbidden(c, "msg")              // 403
helpers.NotFound(c, "msg")               // 404
helpers.InternalServerError(c, "msg")    // 500
helpers.ValidationResponse(c, errs)      // 422
helpers.ValidationError(c, field, msg)   // 422
```

### Validation Pattern (Handler)

```go
if err := c.BindJSON(&req); err != nil {
    helpers.BadRequest(c, "Invalid JSON payload")
    return
}
if err := h.Validate.Struct(&req); err != nil {
    helpers.ValidationResponse(c, h.getErrorsMap(err))
    return
}

// ALL service errors MUST go through HandleError:
if helpers.HandleError(c, err, "Failed to process request") {
    return
}
```

See `references/handler.go` for full examples.

---

## Packages (`internal/pkg/`)

### Redis — `s.RedisClient`

Always check `IsCacheAvailable()` first. Redis errors MUST NOT fail operations.

```go
if s.RedisClient.IsCacheAvailable() {
    s.RedisClient.SetJSON("key", value, time.Hour*24)
    s.RedisClient.GetJSON("key", &dest)
    s.RedisClient.Delete("key")
}
```

### Email — `s.EmailClient`

```go
s.EmailClient.IsConfigured()
s.EmailClient.SendResetPasswordEmail(email, token, resetURL)
```

### Access — `s.Access` (Permission/Role)

3-tier cache: L1 (in-memory) → L2 (Redis) → L3 (DB)

```go
s.Access.HasPermission(ctx, "user.delete")  // Check ANY permission
s.Access.HasRole(ctx, "admin")              // Check role
s.Access.Invalidate(userID)                 // Clear cache
```

Permission naming: `{resource}.{action}` (e.g., `user.create`, `user.delete`)

### Scoped Permissions (Satker/Unit Restriction)

When a role must only access data within their own unit/satker, use scoped permission naming:

| Pattern                      | Example             | Description                                  |
| ---------------------------- | ------------------- | -------------------------------------------- |
| `{resource}.{action}`        | `user.index`        | Full access — view all data                  |
| `{resource}.{action}-satker` | `user.index-satker` | Scoped access — view data in own satker only |

**Route Usage (Access Control):**

```go
// Allow BOTH permissions — middleware checks if user has ANY of them
users.GET("", middleware.RequirePermission(acc, "user.index", "user.index-satker"), handlers.UserGetAll)
```

**Service Usage (Data Filtering):**

```go
func (s *Services) UserGetAll(ctx context.Context) ([]dtos.UserDTO, error) {
    // Scoped access — filter by caller's satker
    if s.Access.HasPermission(ctx, "user.index-satker") {
        callerID := helpers.GetCallerID(ctx)
        caller, _ := s.repo.User.FindByID(nil, callerID, "Satker")
        return s.repo.User.FindByFieldMap(nil, map[string]interface{}{
            "satker_id": caller.SatkerID,
        }, "Roles")
    }
    // Default: return all (full access or no restriction)
    return s.repo.User.FindAll(nil, "Roles")
}
```

**Rules:**

- Routes allow BOTH permissions (full OR scoped)
- Service decides data scope based on which permission user actually has
- Filtering logic belongs in Service layer, NOT handler or route

See `references/route.go` for full examples.

---

## Services & Handlers Struct

```go
type Services struct {
    repo         *repositories.Repositories
    RedisClient  *database.RedisCache
    EmailClient  *pkgEmail.EmailClient    // internal/pkg/email
    Access       *helpers.Access
    Logger       *helpers.Logger
    // ...
}

type Handlers struct {
    svcs     *services.Services
    Validate *validator.Validate
    trans    ut.Translator
}
```

---

## Anti-Patterns

| WRONG                                           | CORRECT                                               |
| ----------------------------------------------- | ----------------------------------------------------- |
| `type PermissionService struct`                 | Method on `(s *Services)`                             |
| `func (s *Services) CreatePermission(...)`      | `func (s *Services) PermissionCreate(...)`            |
| `FindByID(s.repo.Permission.DB, id)`            | `FindByID(nil, id)`                                   |
| Direct DB write without transaction             | `TxManager.WithinTransaction(...)`                    |
| `c.ShouldBindJSON(&req)`                        | `c.BindJSON(&req)` + `h.Validate.Struct(&req)`        |
| `helpers.BadRequest(c, "msg")` for field errors | `helpers.ValidationError(c, "field", "msg")` |
| Manual error check `if err == helpers.ErrX`     | `helpers.HandleError(c, err, "fallback")`            |

---

## Pre-Implementation Checklist

- [ ] Model has `TableName()`
- [ ] DTO variables use feature prefix
- [ ] Request DTO: no `_id` suffix (e.g. `roles`), Response DTO: keep `_id` (e.g. `user_id`)
- [ ] Service/Handler: `{Feature}{Action}` on single struct
- [ ] Write ops use `TxManager.WithinTransaction()`
- [ ] Read ops use `nil` parameter
- [ ] Read ops return DTOs, NOT models
- [ ] Paginated reads handle nil opts
- [ ] CREATE/UPDATE/DELETE have logging + notification
- [ ] Notification Data = identifier only (id, name, status)
- [ ] Repository registered in `00_repository.go`
- [ ] Custom repo method: prefer generic methods if they can handle the requirement
- [ ] Handler uses `c.BindJSON()` + `h.Validate.Struct()`
- [ ] ALL errors handled via `HandleError()`, NO manual error checks
- [ ] Build: `go build ./...` | Vet: `go vet ./...`
