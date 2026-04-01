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
GIN_MODE=release

# Database
DB_CONNECTION=mysql
DB_HOST=127.0.0.1
DB_PORT=3306
DB_DATABASE=laravel
DB_USERNAME=root
DB_PASSWORD=

# Redis (optional)
REDIS_ENABLED=true
REDIS_HOST=localhost
REDIS_PORT=6379

# JWT
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRATION=24
JWT_REFRESH_EXPIRATION=168
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
| GET | `/api/health` | ❌ | API health check |

### Authentication (Public)

| Method | Endpoint | Auth Required | Description |
|--------|----------|---------------|-------------|
| POST | `/api/auth/register` | ❌ | Register new user |
| POST | `/api/auth/login` | ❌ | Login and get JWT tokens |
| POST | `/api/auth/refresh` | ❌ | Refresh access token |

### Authentication (Protected)

| Method | Endpoint | Auth Required | Description |
|--------|----------|---------------|-------------|
| GET | `/api/auth/me` | ✅ | Get current user info |
| POST | `/api/auth/logout` | ✅ | Logout user |

### Request/Response Examples

**Register:**
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123","name":"John Doe"}'
```

**Login:**
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

**Protected Route:**
```bash
curl -X GET http://localhost:8080/api/auth/me \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

## Development Conventions

### ⚠️ Critical Rules (MANDATORY)

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

**✅ CORRECT - Service:**
```go
// internal/services/00_services.go
type Services struct {
    repo        *repositories.Repositories
    RedisClient *database.RedisCache
    cfg         *JWTConfig
}

// internal/services/auth_service.go
func (s *Services) AuthLogin(ctx context.Context, email, password string) (*dtos.LoginResponse, error) {
    // business logic
}

func (s *Services) AuthRefreshToken(ctx context.Context, refreshToken string) (*dtos.LoginResponse, error) {
    // business logic
}
```

**✅ CORRECT - Handler:**
```go
// internal/handlers/00_handlers.go
type Handlers struct {
    svcs *services.Services
}

// internal/handlers/auth_handler.go
func (h *Handlers) AuthLogin(c *gin.Context) {
    // handle request
}

func (h *Handlers) AuthLogout(c *gin.Context) {
    // handle request
}
```

**❌ WRONG - Do not create separate structs:**
```go
// DO NOT DO THIS
type AuthService struct {
    repo *repositories.Repositories
}

func (a *AuthService) Login(...) { }
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
├── logs/                 # Application logs
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
   - Models: `snake_case.go` dengan `TableName()` method
   - Functions: `PascalCase` for exported, `camelCase` for private

---

## 📝 Logging Rules (MANDATORY)

> ⚠️ **Logger tersedia di Services**: `s.Logger`

### **Kapan Wajib Log?**

| Operasi | Wajib? | Operasi | Wajib? |
|---------|--------|---------|--------|
| **CREATE** | ✅ YA | **GET Simple** | ❌ TIDAK |
| **UPDATE** | ✅ YA | **GET Complex** | ✅ YA |
| **DELETE** | ✅ YA | **Auth** | ✅ YA |

### **Cara Menggunakan**

```go
s.Logger.LogStart("FuncName", "Message: %s", value)      // Mulai
s.Logger.LogStep("FuncName", "Step: %s", value)          // Step
s.Logger.LogStepWithPrefix("Func", "[OK]", "Done")       // Step dengan prefix
s.Logger.LogEnd("FuncName", "Success: %s", value)        // Selesai
s.Logger.LogEndWithError("Func", "Error: %v", err)       // Selesai + error
s.Logger.LogError("FuncName", "Error: %v", err)          // Error
s.Logger.LogWarn("FuncName", "Warning: %s", value)       // Warning
s.Logger.LogInfo("FuncName", "Info: %s", value)          // Info
```

### **Contoh - CREATE (WAJIB LOG)**

```go
func (s *Services) CreateUser(ctx context.Context, email string) error {
    s.Logger.LogStart("CreateUser", "Creating user: %s", email)
    
    user := &models.User{Email: email}
    if err := s.repo.User.Create(s.repo.User.DB, user); err != nil {
        s.Logger.LogEndWithError("CreateUser", "Failed: %v", err)
        return err
    }
    
    s.Logger.LogEnd("CreateUser", "User created: %s (ID: %d)", email, user.ID)
    return nil
}
```

### **Contoh - GET Simple (TANPA LOG)**

```go
func (s *Services) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
    return s.repo.User.FindByID(s.repo.User.DB, id)
}
```

**Output**: `logs/YYYY-MM-DD.log` | Auto rotation & cleanup (30 hari)

---

## 📚 Available Helpers & Clients

### **Helpers** (`internal/helpers/`)

> ⚠️ **IMPORTANT**: Always use existing helpers before creating new ones. If an operation can be generalized, create a new helper!

#### **Environment Helpers** (`env.go`)
| Function | Description | Example |
|----------|-------------|---------|
| `GetEnv(key, default)` | Get environment variable with default | `helpers.GetEnv("APP_PORT", "8080")` |
| `GetEnvInt(key, default)` | Get environment variable as int | `helpers.GetEnvInt("JWT_EXPIRATION", 24)` |

#### **Error Helpers** (`error.go`)
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

#### **Crypto Helpers** (`crypto.go`)
| Function | Description | Example |
|----------|-------------|---------|
| `GenerateRandomString(length)` | Generate cryptographically secure random string | `helpers.GenerateRandomString(32)` |
| `HashString(str)` | Hash string using bcrypt | `helpers.HashString("mytoken")` |
| `VerifyString(str, hash)` | Verify string against hash | `helpers.VerifyString("token", hash)` |

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
```go
// internal/services/your_service.go
func (s *Services) YourMethod(ctx context.Context, param string) error {
    // business logic
    return nil
}
```

### 5. Add Handler Method (MANDATORY: on Handlers struct)
```go
// internal/handlers/your_handler.go
func (h *Handlers) YourHandler(c *gin.Context) {
    // handle request
}
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
