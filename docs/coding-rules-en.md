# Go Boilerplate Project

## Project Overview

This is a **Go boilerplate** for REST API with a modern, layered architecture. The project is designed to accelerate API development with standard features already implemented.

### Tech Stack

- **Framework**: Gin (HTTP server)
- **Database**: MySQL (default) / PostgreSQL (optional)
- **ORM**: GORM
- **Cache**: Redis (optional)
- **Authentication**: JWT (golang-jwt/jwt/v5)
- **Password Hashing**: bcrypt (golang.org/x/crypto)
- **Environment**: godotenv

### Architecture

This project uses **Dependency Injection** pattern with a layered structure:

```
Routes → Handlers → Services → Repositories → Database/Redis
              ↑
         Middleware (JWT Auth, CORS)
```

#### Layer Components:

| Layer | Location | Description |
|-------|----------|-------------|
| **Routes** | `internal/routes/` | Defines endpoints and maps to handlers |
| **Handlers** | `internal/handlers/` | HTTP handlers, request/response handling |
| **Services** | `internal/services/` | Business logic, JWT operations |
| **Repositories** | `internal/repositories/` | Data access layer with generic repository pattern |
| **Models** | `internal/models/` | GORM models for database entities |
| **DTOs** | `internal/dtos/` | Data Transfer Objects for request/response |
| **Middleware** | `internal/middleware/` | JWT auth, CORS, and other middleware |
| **DI Container** | `internal/di/` | Dependency Injection container |

## Building and Running

### Prerequisites

- Go 1.25.0+
- MySQL 8+ / PostgreSQL 14+ (one of them)
- Redis (optional)

### Setup Environment

1. Copy environment file:
```bash
cp .env.example .env
```

2. Edit `.env` as needed:
```env
# Server
APP_HOST=0.0.0.0
APP_PORT=8080
APP_FE_URL=http://localhost:3000

# Gin
GIN_MODE=release

# Trusted Proxies (comma-separated IPs, leave empty to disable)
TRUSTED_PROXIES=

# CORS Allowed Origins (comma-separated or * for all)
ALLOWED_ORIGINS=*

# Database
DB_CONNECTION=mysql
DB_HOST=127.0.0.1
DB_PORT=3306
DB_DATABASE=laravel
DB_USERNAME=root
DB_PASSWORD=

# Redis
REDIS_ENABLED=true
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRATION=24
JWT_REFRESH_EXPIRATION=168

# Rate Limiting
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=60

# Email/SMTP (optional - for password reset)
SMTP_HOST=
SMTP_PORT=587
SMTP_USER=
SMTP_PASSWORD=
SMTP_FROM=
```

### Commands

```bash
# Install dependencies
go mod download

# Run development
go run cmd/api/main.go

# Build binary
go build -o main cmd/api/main.go

# Run compiled binary
./main

# Run migration (folder available at cmd/migration/)
go run cmd/migration/*.go
```

### Default Server

Server will run on `http://localhost:8080`

## API Endpoints

### Health Check

| Method | Endpoint | Auth Required | Description |
|--------|----------|---------------|-------------|
| GET | `/health` | ❌ | Health check endpoint |

### Authentication (Public)

| Method | Endpoint | Auth Required | Description |
|--------|----------|---------------|-------------|
| POST | `/api/auth/login` | ❌ | Login and get JWT tokens |
| POST | `/api/auth/refresh` | ❌ | Refresh access token |
| POST | `/api/auth/forgot-password` | ❌ | Request password reset |
| POST | `/api/auth/reset-password` | ❌ | Reset password with token |

### Authentication (Protected)

| Method | Endpoint | Auth Required | Description |
|--------|----------|---------------|-------------|
| POST | `/api/auth/logout` | ✅ | Logout user |

### Request/Response Examples

**Login:**
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

**Protected Route:**
```bash
curl -X GET http://localhost:8080/api/permissions \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

## Development Conventions

### 🔄 Development Workflow (MANDATORY)

Every new feature development or bug fix **MUST** follow this flow:

#### Phase 1: Planning

1. **Create `plan.md`** in root directory containing:
   - Feature/bug description to be implemented
   - Technical analysis (files to be created/modified)
   - Implementation steps
   - Impact on existing code (if any)

2. **Confirm to User**: Display plan and wait for user approval.
   - ✅ If user approves → proceed to step 3
   - ❌ If user requests revision → go back to step 1, update `plan.md`

3. **Create GitHub Issue** with naming format:
    ```
    [BE] [Type] Title
    ```
    Types to use:
    - `[Feat]` - New feature
    - `[Bug]` - Bug fix
    - `[Fix]` - Small improvement
    - `[Refactor]` - Refactoring code
    - `[Chore]` - Maintenance task

    Example: `[BE] [Feat] Add Product CRUD API`

#### Phase 2: Implementation

1. **Create New Branch** from main/master:
    ```bash
    git checkout -b feat/product-crud
    ```

2. **Implement according to Git Issue**: Work based on the approved plan.

3. **Create Pull Request** to GitHub:
   - Reference git issue in PR description
   - Format: `Closes #<issue-number>` or `Fixes #<issue-number>`

4. **Confirm to User**: Display summary of changes and request approval.
   - ✅ If user approves → merge PR and close git issue
   - ❌ If user requests changes → update according to feedback

> ⚠️ **IMPORTANT**: Git issue **MUST** exist before implementation. No issue = no coding.

### ⚠️ Critical Rules (MANDATORY)

> 📖 **COMPLETE**: For CRUD implementation flow, naming convention, and anti-patterns details, open **[docs/development-guide-en.md](docs/development-guide-en.md)**.

#### 1. DO NOT Modify Core Repository Files

The following files **MUST NOT be modified** in any way:

- `internal/repositories/00_generic.go` - Generic repository pattern (core)
- `internal/repositories/00_transaction.go` - Transaction manager (core)

These files are the **foundation** of the repository architecture. If you need custom logic, create a new repository file (e.g., `user_repository.go`, `product_repository.go`).

#### 2. Models MUST Define TableName

Every model **MUST** have a `TableName()` method to explicitly define the table name:

```go
type User struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    Email     string         `gorm:"uniqueIndex;size:255;not null" json:"email"`
    // ... fields
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// MANDATORY
func (User) TableName() string {
    return "users"
}
```

#### 3. Services & Handlers MUST Use Single Struct Pattern

**All service methods MUST be methods of the `Services` struct**, and **all handler methods MUST be methods of the `Handlers` struct**.

**All functions MUST use feature name as prefix.**

> 📖 **COMPLETE**: See complete example in **[docs/development-guide-en.md](docs/development-guide-en.md)** — Service & Handler section.

```go
// ✅ CORRECT — Service with feature prefix (FeatureName + Action)
func (s *Services) PermissionCreate(ctx context.Context, req dtos.PermissionRequest) (*dtos.PermissionDTO, error) { ... }
func (s *Services) PermissionGetAll(ctx context.Context) ([]dtos.PermissionDTO, error) { ... }

// ✅ CORRECT — Handler with feature prefix
func (h *Handlers) PermissionCreate(c *gin.Context) { ... }
func (h *Handlers) PermissionGetAll(c *gin.Context) { ... }

// ❌ WRONG — Action first instead of feature first
func (s *Services) CreatePermission(...) { }
func (s *Services) GetAllPermissions(...) { }
```

### Project Structure

```
go-boilerplate/
├── cmd/
│   ├── api/              # Main application entry point
│   └── migration/        # Database migration scripts
├── docs/                 # API documentation
├── internal/             # Private application code
│   ├── client/           # External API clients
│   ├── database/         # Database connection & Redis
│   ├── di/               # Dependency Injection
│   ├── dtos/             # Data Transfer Objects
│   ├── handlers/         # HTTP handlers (single struct: Handlers)
│   ├── helpers/          # Helper functions
│   ├── middleware/       # Gin middleware
│   ├── models/           # GORM models (mandatory TableName)
│   ├── repositories/     # Data access layer
│   │   ├── 00_generic.go      # ⚠️ DO NOT MODIFY
│   │   ├── 00_transaction.go  # ⚠️ DO NOT MODIFY
│   │   └── user_repository.go # Custom repository (safe to modify)
│   ├── routes/           # Route definitions
│   └── services/         # Business logic (single struct: Services)
├── storage/              # Application storage (keys, logs)
├── .env                  # Environment variables (gitignore)
├── .env.example          # Environment template
├── go.mod                # Go module definition
└── main                  # Compiled binary
```

### Coding Patterns

1. **Generic Repository Pattern**: File `internal/repositories/00_generic.go` provides generic CRUD operations that can be reused.

2. **Dependency Injection**: All dependencies are injected via `internal/di/container.go`.

3. **JWT Middleware**: Protected routes use JWT middleware at `internal/middleware/jwt.go`.

4. **Error Handling**: Uses custom errors in `internal/helpers/`.

5. **Naming Conventions**:
   - Routes: `*_route.go` (e.g., `auth_route.go`, `health_route.go`)
   - DTOs: `*_dto.go` (e.g., `auth_dto.go`, `health_dto.go`)
   - Handlers: `*_handler.go` (e.g., `auth_handler.go`)
   - Services: `*_service.go` (e.g., `auth_service.go`)
   - Repositories: `*_repository.go` (e.g., `user_repository.go`)
   - Models: `snake_case.go` with `TableName()` method
   - Functions: `PascalCase` for exported, `camelCase` for private

---

## 📝 Logging Rules (MANDATORY)

> ⚠️ **Logger available in Services**: `s.Logger`

### **When to Log?**

| Operation | Required? | Operation | Required? |
|-----------|-----------|-----------|-----------|
| **CREATE** | ✅ YES | **GET Simple** | ❌ NO |
| **UPDATE** | ✅ YES | **GET Complex** | ✅ YES |
| **DELETE** | ✅ YES | **Auth** | ✅ YES |

### **How to Use**

```go
s.Logger.LogStart("FuncName", "Message: %s", value)      // Start
s.Logger.LogStep("FuncName", "Step: %s", value)          // Step
s.Logger.LogStepWithPrefix("Func", "[OK]", "Done")       // Step with prefix
s.Logger.LogEnd("FuncName", "Success: %s", value)        // End
s.Logger.LogEndWithError("Func", "Error: %v", err)       // End + error
s.Logger.LogError("FuncName", "Error: %v", err)          // Error
s.Logger.LogWarn("FuncName", "Warning: %s", value)       // Warning
s.Logger.LogInfo("FuncName", "Info: %s", value)          // Info
```

### **Example - CREATE (MUST LOG)**

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

### **Example - GET Simple (NO LOG)**

```go
func (s *Services) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
    return s.repo.User.FindByID(s.repo.User.DB, id)
}
```

**Output**: `storage/logs/YYYY-MM-DD.log` | Auto rotation & cleanup (30 days)

---

## 📚 Available Helpers & Clients

### **Helpers** (`internal/helpers/`)

> ⚠️ **IMPORTANT**: Always use existing helpers before creating new ones. If an operation can be generalized, create a new helper!

#### **Environment Helpers** (`env_helper.go`)
| Function | Description | Example |
|----------|-------------|---------|
| `GetEnv(key, default)` | Get environment variable with default | `helpers.GetEnv("APP_PORT", "8080")` |
| `GetEnvInt(key, default)` | Get environment variable as int | `helpers.GetEnvInt("JWT_EXPIRATION", 24)` |

#### **Error Helpers** (`error_helper.go`)
| Error | Description |
|-------|-------------|
| `ErrNotFound` | Record not found |
| `ErrInvalidToken` | Invalid JWT token |
| `ErrExpiredToken` | Token has expired |
| `ErrInvalidCredential` | Invalid email or password |
| `ErrUserExists` | User already exists |
| `ErrInvalidEmail` | Invalid email address |
| `ErrTokenExpired` | Reset token has expired |
| `ErrTokenUsed` | Reset token already used |
| `ErrTokenInvalid` | Invalid reset token |

#### **Crypto Helpers** (`crypto_helper.go`)
| Function | Description | Example |
|----------|-------------|---------|
| `GenerateRandomString(length)` | Generate cryptographically secure random string | `helpers.GenerateRandomString(32)` |
| `HashString(str)` | Hash string using bcrypt | `helpers.HashString("mytoken")` |
| `VerifyString(str, hash)` | Verify string against hash | `helpers.VerifyString("token", hash)` |

#### **Response Helpers** (`response_helper.go`)
| Function | Status | Description | Example |
|----------|--------|-------------|---------|
| `OK(c, msg, data)` | 200 | Success response | `helpers.OK(c, "Success", user)` |
| `Created(c, msg, data)` | 201 | Created response | `helpers.Created(c, "Created", user)` |
| `BadRequest(c, msg)` | 400 | Bad request | `helpers.BadRequest(c, "Invalid input")` |
| `Unauthorized(c, msg)` | 401 | Unauthorized | `helpers.Unauthorized(c, "Invalid token")` |
| `Forbidden(c, msg)` | 403 | Forbidden | `helpers.Forbidden(c, "Access denied")` |
| `NotFound(c, msg)` | 404 | Not found | `helpers.NotFound(c, "User not found")` |
| `InternalServerError(c, msg)` | 500 | Server error | `helpers.InternalServerError(c, "Error")` |

**Response format:**
```json
{
  "code": 200,
  "message": "Success message",
  "data": {} // optional
}
```

---

### **Clients** (`internal/clients/`)

#### **Email Client** (`internal/clients/email/`)

> 📧 Email client is injected into Services via DI Container. Access via `s.EmailClient`.

**Struct:** `EmailClient`

**Methods:**

| Method | Description | Example |
|--------|-------------|---------|
| `IsConfigured()` | Check if SMTP is configured | `s.EmailClient.IsConfigured()` |
| `SendEmail(req)` | Send email with custom request | `s.EmailClient.SendEmail(req)` |
| `SendResetPasswordEmail(to, token, resetURL)` | Send reset password email | `s.EmailClient.SendResetPasswordEmail(email, token, url)` |

**DTOs:**

```go
// EmailRequest - for custom email
type EmailRequest struct {
    To      []string  // Recipients
    Subject string    // Email subject
    Body    string    // HTML body
    CC      []string  // CC recipients
    BCC     []string  // BCC recipients
}
```

**Example Usage:**

```go
// Send reset password email (recommended)
err := s.EmailClient.SendResetPasswordEmail(
    user.Email, 
    token, 
    "https://myapp.com/reset-password?token="+token,
)

// Send custom email
req := email.EmailRequest{
    To:      []string{"user@example.com"},
    Subject: "Welcome!",
    Body:    "<h1>Welcome to our app!</h1>",
}
err := s.EmailClient.SendEmail(req)
```

**Templates** (`internal/clients/email/templates/`):

| Template | Description |
|----------|-------------|
| `ResetPasswordEmail(token, resetURL, appName)` | HTML template for reset password email |

---

### **Redis Cache** (`internal/database/redis_cache.go`)

> 🔴 Redis client is injected into Services via DI Container. Access via `s.RedisClient`.

**IMPORTANT**: Always check `s.RedisClient.IsCacheAvailable()` before accessing Redis. Redis errors **MUST NOT** cause operations to fail.

#### Usage

```go
// Check first if Redis is active
if s.RedisClient.IsCacheAvailable() {
    // SET - store data
    err := s.RedisClient.SetJSON("session:1", userDTO, time.Hour*24)
    if err != nil {
        s.Logger.LogWarn("FuncName", "Redis SET failed: %v", err) // fallback, don't return error
    }

    // GET - retrieve data
    var cached dtos.UserDTO
    if err := s.RedisClient.GetJSON("session:1", &cached); err == nil {
        // Cache hit - use cached data
    } else {
        // Cache miss/error - fallback to DB
    }

    // DELETE
    s.RedisClient.Delete("session:1")
}
```

#### Recommended Key Pattern
| Pattern | Example | Description |
|---------|--------|-----------|
| `session:{userID}` | `session:1` | User session cache |

#### Main Methods
| Method | Description |
|--------|-----------|
| `IsCacheAvailable()` | Check if Redis is active |
| `SetJSON(key, value, ttl)` | Store JSON with TTL |
| `GetJSON(key, &dest)` | Get & unmarshal JSON |
| `Delete(keys...)` | Delete key |

---

### **Access - Permission & Role Checking** (`internal/helpers/access.go`)

> 🔐 Access helper is injected into Services via DI Container. Access via `s.Access`.
> Also available in routes as `container.Access` for middleware.

**Struct:** `Access`

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

**Behavior:**
- Redis unavailable or cache miss → fallback to DB → cache result
- If user not found (all tiers) → returns `false`
- `Invalidate()` clears both local cache and Redis session

#### Usage in Services

```go
// Single permission check
func (s *Services) UserDelete(ctx context.Context, id uint) error {
    if !s.Access.HasPermission(ctx, "user.delete") {
        return helpers.ErrForbidden
    }
    // ... delete logic
}

// Multiple permissions - returns true if user has ANY of them
func (s *Services) UserAdminAction(ctx context.Context, req dtos.AdminRequest) error {
    if !s.Access.HasPermission(ctx, "user.delete", "user.update", "user.admin") {
        return helpers.ErrForbidden
    }
    // ...
}

// Role check
func (s *Services) SuperAdminOnly(ctx context.Context) error {
    if !s.Access.HasRole(ctx, "superadmin") {
        return helpers.ErrForbidden
    }
    // ...
}
```

#### Usage in Routes (Middleware)

```go
// internal/routes/user_route.go
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

#### Cache Invalidation

Call `s.Access.Invalidate(userID)` when user's roles/permissions change:

```go
func (s *Services) UserUpdate(ctx context.Context, id uint, req dtos.UserUpdateRequest) (*dtos.UserDTO, string, error) {
    // ... update logic with role changes ...
    
    // Invalidate cached session so next request gets updated permissions
    s.Access.Invalidate(id)
    
    return &dto, oldAvatar, nil
}
```

#### Permission Naming Convention

Use dot separator: `{resource}.{action}`

| Permission | Description |
|------------|-------------|
| `user.index` | View users list |
| `user.create` | Create new user |
| `user.edit` | Update user |
| `user.delete` | Delete user |
| `role.index` | View roles |
| `role.create` | Create role |
| `role.edit` | Update role |
| `role.delete` | Delete role |
| `permission.index` | View permissions |
| `permission.create` | Create permission |
| `permission.edit` | Update permission |
| `permission.delete` | Delete permission |

---

### **Creating New Helpers**

> 💡 **Rule of Thumb**: If a function can be used in more than 1 place, make it a helper!

**Criteria for creating new helpers:**
1. ✅ General operation (not specific business logic)
2. ✅ Reusable across multiple services
3. ✅ No dependencies on repo, db, etc
4. ✅ Pure function (input → output, no side effects)

**Example of creating a new helper:**

```go
// internal/helpers/your_helper.go
package helpers

// FormatCurrency formats number as Indonesian Rupiah
func FormatCurrency(amount int64) string {
    return fmt.Sprintf("Rp %s", strings.ReplaceAll(
        strconv.FormatInt(amount, 10), 
        "000", 
        ".000",
    ))
}
```

**❌ DO NOT create helpers if:**
- ❌ Only used in 1 service
- ❌ Contains specific business logic
- ❌ Requires dependency injection (db, repo, etc)
- ❌ Has side effects (write to db, send email, etc)

---

### Testing

To add testing, create `_test.go` files in the same folder as the file being tested:
```bash
# Examples
internal/handlers/auth_handler_test.go
internal/services/auth_service_test.go
```

## Database Schema

### Users Table

| Column | Type | Constraints |
|--------|------|-------------|
| id | uint | PRIMARY KEY, AUTO_INCREMENT |
| email | string(255) | UNIQUE, NOT NULL |
| password | string(255) | NOT NULL (hashed) |
| name | string(255) | - |
| created_at | timestamp | - |
| updated_at | timestamp | - |
| deleted_at | timestamp | SOFT DELETE index |

## Adding New Features

### 1. Add New Model
```go
// internal/models/your_model.go
type YourModel struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    // ... fields
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// MANDATORY: Define table name
func (YourModel) TableName() string {
    return "your_table_name"
}
```

### 2. Add New Repository

**IMPORTANT**: Use **GenericRepository** for all standard CRUD operations (Create, Read, Update, Delete, Find, Count, etc.). Add custom methods **only** for operations that cannot be handled by the generic repository, such as:
- Queries with complex JOINs (not regular preloads)
- Queries with specific subqueries
- Special batch operations
- Queries with specific aggregations

```go
// internal/repositories/your_repository.go
type YourRepository struct {
    *GenericRepository[YourModel]
}

func NewYourRepository(db *gorm.DB) *YourRepository {
    return &YourRepository{
        GenericRepository: NewGenericRepository(db, &YourModel{}),
    }
}

// ✅ EXAMPLE OF CORRECT CUSTOM METHOD (complex JOIN)
func (r *YourRepository) FindWithComplexJoin(id uint) (*YourModel, error) {
    var result YourModel
    query := r.DB.Joins("LEFT JOIN other_table ON other_table.your_id = your_models.id").
        Where("your_models.id = ?", id).
        First(&result)
    
    if query.Error != nil {
        return nil, query.Error
    }
    return &result, nil
}

// ❌ DO NOT DO THIS (use generic repo directly)
// func (r *YourRepository) FindByID(id uint) (*YourModel, error) {
//     return r.FindByID(r.DB, id) // This already exists in GenericRepository!
// }
```

### 3. Register Repository to Container

After creating a repository, you **MUST** register it in `internal/repositories/00_repository.go` so it can be used in services:

```go
// internal/repositories/00_repository.go
type Repositories struct {
    TxManager *TransactionManager
    User      *UserRepository
    YourModel *YourRepository  // ← Add this
}

func NewRepositories(db *gorm.DB) (*Repositories, error) {
    txManager := NewTransactionManager(db)
    userRepo := NewUserRepository(db)
    yourRepo := NewYourRepository(db)  // ← Initialize

    return &Repositories{
        TxManager: txManager,
        User:      userRepo,
        YourModel: yourRepo,  // ← Register
    }, nil
}
```

### 4. Add Service Method (MANDATORY: on Services struct)

> 📖 **COMPLETE**: See complete write/read operations example in **[docs/development-guide-en.md](docs/development-guide-en.md)** — Service section.

```go
// internal/services/your_feature_service.go
// MANDATORY: Feature name prefix (FeatureName + Action), TxManager for write, nil for read
func (s *Services) YourFeatureCreate(ctx context.Context, req dtos.YourFeatureRequest) (*dtos.YourFeatureDTO, error) { ... }
func (s *Services) YourFeatureGetByID(ctx context.Context, id uint) (*dtos.YourFeatureDTO, error) { ... }
```

**Write Operations** — Use `s.repo.TxManager.WithinTransaction()`:
```go
func (s *Services) YourFeatureCreate(ctx context.Context, req dtos.YourFeatureRequest) (*dtos.YourFeatureDTO, error) {
    s.Logger.LogStart("YourFeatureCreate", "Creating: %s", req.Name)

    entity := &models.YourModel{Name: req.Name}

    var result *models.YourModel
    if err := s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
        var err error
        result, err = s.repo.YourModel.Create(tx, entity)
        return err
    }); err != nil {
        s.Logger.LogEndWithError("YourFeatureCreate", "Failed: %v", err)
        return nil, err
    }

    dto := dtos.ToYourFeatureDTO(result)
    s.Logger.LogEnd("YourFeatureCreate", "Created: %s (ID: %d)", dto.Name, dto.ID)
    return &dto, nil
}
```

**Write with Result Reload** — Use `s.repo.TxManager.WithinTransactionWithResult()` when you need to return the entity after additional operations inside the transaction (e.g., assigning roles, reloading associations):
```go
func (s *Services) UserCreate(ctx context.Context, req dtos.UserCreateRequest) (*dtos.UserDTO, error) {
    s.Logger.LogStart("UserCreate", "Creating user: %s", req.Email)

    user := &models.User{Email: req.Email, Name: req.Name, Password: hashedPassword}

    res, err := s.repo.TxManager.WithinTransactionWithResult(func(tx *gorm.DB) (interface{}, error) {
        result, err := s.repo.User.Create(tx, user)
        if err != nil {
            return nil, err
        }

        // Assign roles inside transaction
        var roles []models.Role
        for _, roleID := range req.Roles {
            roles = append(roles, models.Role{ID: roleID})
        }
        if err := tx.Model(&result).Association("Roles").Append(roles); err != nil {
            return nil, err
        }

        // Reload with associations
        return s.repo.User.FindByID(tx, result.ID, "Roles")
    })
    if err != nil {
        s.Logger.LogEndWithError("UserCreate", "Failed: %v", err)
        return nil, err
    }

    result := res.(*models.User)
    dto := dtos.ToUserDTO(result)
    s.Logger.LogEnd("UserCreate", "User created: %s (ID: %d)", dto.Email, dto.ID)
    return &dto, nil
}
```

**Read Operations** — Use `nil` parameter (NOT `s.repo.Feature.DB`):
```go
func (s *Services) YourFeatureGetAll(ctx context.Context) ([]dtos.YourFeatureDTO, error) {
    entities, err := s.repo.YourModel.FindAll(nil)
    if err != nil {
        return nil, err
    }
    return dtos.ToYourFeatureDTOList(entities), nil
}
```

**Pagination** — Use `FindAllWithOpts` for paginated, sorted, and searchable queries:
```go
func (s *Services) YourFeatureGetAllPaginated(ctx context.Context, opts *repositories.QueryOptions) (*repositories.PagedResult[models.YourModel], error) {
    if opts == nil {
        opts = &repositories.QueryOptions{}
    }
    if opts.SortBy == "" {
        opts.SortBy = "id"
    }
    if opts.Order == "" {
        opts.Order = "ASC"
    }

    return s.repo.YourModel.FindAllWithOpts(nil, opts)
}
```

`QueryOptions` fields:
- `Page`, `PageSize` — pagination
- `SortBy`, `Order` — sorting (`"ASC"` or `"DESC"`)
- `Search`, `SearchFields` — LIKE search across multiple fields
- `Preloads` — relations to preload (e.g., `[]string{"Roles", "Roles.Permissions"}`)

### 5. Add Handler Method (MANDATORY: on Handlers struct)

> 📖 **COMPLETE**: See complete example in **[docs/development-guide-en.md](docs/development-guide-en.md)** — Handler section.

```go
// internal/handlers/your_feature_handler.go
// MANDATORY: Feature name prefix — {Feature}{Action}
func (h *Handlers) YourFeatureCreate(c *gin.Context) { ... }
func (h *Handlers) YourFeatureGetAll(c *gin.Context) { ... }
```

### 6. Add Routes
```go
// internal/routes/your_routes.go
func RegisterYourRoutes(r *gin.RouterGroup, handlers *handlers.Handlers) {
    your := r.Group("/your")
    {
        your.GET("/", handlers.YourHandler)
    }
}
```

### 7. Register in DI Container

Update `internal/di/container.go` to inject the new repository into the container:

```go
// internal/di/container.go
func NewContainer() (*Container, error) {
    // ... existing code ...

    // Initialize repositories
    repos, err := repositories.NewRepositories(container.DB)
    if err != nil {
        return nil, fmt.Errorf("failed to initialize repositories: %w", err)
    }
    container.Repositories = repos

    // ... existing code ...
}
```

The repository is automatically available in services via `s.repo.YourModel`.

## Security Notes

- ⚠️ **JWT_SECRET**: Change to a strong secret key in production
- ⚠️ **Database Password**: Do not commit credentials to git
- ⚠️ **GIN_MODE**: Set to `release` in production
- ⚠️ **CORS**: Configure `ALLOWED_ORIGINS` according to allowed domains
- ⚠️ **TRUSTED_PROXIES**: Set if using reverse proxy (nginx, load balancer)

## Documentation

- API documentation available in `docs/` folder
- Environment configuration template at `.env.example`
