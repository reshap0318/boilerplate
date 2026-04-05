# Go Boilerplate Project

## Project Overview

Ini adalah **Go boilerplate** untuk REST API dengan arsitektur berlapis yang modern dan terstruktur. Proyek ini dirancang untuk mempercepat pengembangan API dengan fitur-fitur standar yang sudah terimplementasi.

### Tech Stack

- **Framework**: Gin (HTTP server)
- **Database**: MySQL (default) / PostgreSQL (opsional)
- **ORM**: GORM
- **Cache**: Redis (opsional)
- **Authentication**: JWT (golang-jwt/jwt/v5)
- **Password Hashing**: bcrypt (golang.org/x/crypto)
- **Environment**: godotenv

### Arsitektur

Proyek ini menggunakan **Dependency Injection** pattern dengan struktur berlapis:

```
Routes → Handlers → Services → Repositories → Database/Redis
              ↑
         Middleware (JWT Auth, CORS)
```

#### Layer Components:

| Layer | Lokasi | Deskripsi |
|-------|--------|-----------|
| **Routes** | `internal/routes/` | Mendefinisikan endpoint dan mapping ke handlers |
| **Handlers** | `internal/handlers/` | HTTP handler, request/response handling |
| **Services** | `internal/services/` | Business logic, JWT operations |
| **Repositories** | `internal/repositories/` | Data access layer dengan generic repository pattern |
| **Models** | `internal/models/` | GORM models untuk database entities |
| **DTOs** | `internal/dtos/` | Data Transfer Objects untuk request/response |
| **Middleware** | `internal/middleware/` | JWT auth, CORS, dan middleware lainnya |
| **DI Container** | `internal/di/` | Dependency Injection container |

## Building and Running

### Prerequisites

- Go 1.25.0+
- MySQL 8+ / PostgreSQL 14+ (salah satu)
- Redis (opsional)

### Setup Environment

1. Copy file environment:
```bash
cp .env.example .env
```

2. Edit `.env` sesuai kebutuhan:
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

# Email/SMTP (opsional - untuk password reset)
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

# Run migration (folder tersedia di cmd/migration/)
go run cmd/migration/*.go
```

### Default Server

Server akan berjalan di `http://localhost:8080`

## API Endpoints

### Health Check

| Method | Endpoint | Auth Required | Deskripsi |
|--------|----------|---------------|-----------|
| GET | `/health` | ❌ | Health check endpoint |
| GET | `/api/health` | ❌ | API health check |

### Authentication (Public)

| Method | Endpoint | Auth Required | Deskripsi |
|--------|----------|---------------|-----------|
| POST | `/api/auth/register` | ❌ | Register user baru |
| POST | `/api/auth/login` | ❌ | Login dan dapatkan JWT tokens |
| POST | `/api/auth/refresh` | ❌ | Refresh access token |

### Authentication (Protected)

| Method | Endpoint | Auth Required | Deskripsi |
|--------|----------|---------------|-----------|
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

### ⚠️ Critical Rules (WAJIB DIPATUHI)

#### 1. JANGAN MODIFIKASI File Core Repository

File berikut **TIDAK BOLEH DIUBAH** dalam bentuk apapun:

- `internal/repositories/00_generic.go` - Generic repository pattern (core)
- `internal/repositories/00_transaction.go` - Transaction manager (core)

Kedua file ini adalah **foundation** dari arsitektur repository. Jika butuh custom logic, buat file repository baru (contoh: `user_repository.go`, `product_repository.go`).

#### 2. Model WAJIB Define TableName

Setiap model **WAJIB** memiliki method `TableName()` untuk mendefinisikan nama tabel secara eksplisit:

```go
type User struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    Email     string         `gorm:"uniqueIndex;size:255;not null" json:"email"`
    // ... fields
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// WAJIB ADA
func (User) TableName() string {
    return "users"
}
```

#### 3. Services & Handlers WAJIB Menggunakan Single Struct Pattern

**Semua service method HARUS merupakan method dari struct `Services`**, dan **semua handler method HARUS merupakan method dari struct `Handlers`**.

**✅ BENAR - Service:**
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

**✅ BENAR - Handler:**
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

**❌ SALAH - Jangan buat struct terpisah:**
```go
// JANGAN LAKUKAN INI
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
│   ├── models/           # GORM models (wajib TableName)
│   ├── repositories/     # Data access layer
│   │   ├── 00_generic.go      # ⚠️ JANGAN DIUBAH
│   │   ├── 00_transaction.go  # ⚠️ JANGAN DIUBAH
│   │   └── user_repository.go # Custom repository (aman diubah)
│   ├── routes/           # Route definitions
│   └── services/         # Business logic (single struct: Services)
├── logs/                 # Application logs
├── .env                  # Environment variables (gitignore)
├── .env.example          # Environment template
├── go.mod                # Go module definition
└── main                  # Compiled binary
```

### Coding Patterns

1. **Generic Repository Pattern**: File `internal/repositories/00_generic.go` menyediakan CRUD operations generik yang bisa di-reuse.

2. **Dependency Injection**: Semua dependencies di-inject melalui `internal/di/container.go`.

3. **JWT Middleware**: Protected routes menggunakan middleware JWT di `internal/middleware/jwt.go`.

4. **Error Handling**: Menggunakan custom errors di `internal/helpers/`.

5. **Helpers & Utilities**: Selalu gunakan helper yang ada di `internal/helpers/` untuk operasi umum.

6. **Email Client**: Gunakan `internal/clients/email/` untuk pengiriman email.

7. **Naming Conventions**:
   - Routes: `*_route.go` (contoh: `auth_route.go`, `health_route.go`)
   - DTOs: `*_dto.go` (contoh: `auth_dto.go`, `health_dto.go`)
   - Handlers: `*_handler.go` (contoh: `auth_handler.go`)
   - Services: `*_service.go` (contoh: `auth_service.go`)
   - Repositories: `*_repository.go` (contoh: `user_repository.go`)
   - Models: `snake_case.go` dengan `TableName()` method
   - Functions: `PascalCase` untuk exported, `camelCase` untuk private

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

> ⚠️ **PENTING**: Selalu gunakan helper yang ada sebelum membuat function baru. Jika ada operasi yang bisa digeneralisasi, buat helper baru!

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
  "data": {} // opsional
}
```

---

### **Clients** (`internal/clients/`)

#### **Email Client** (`internal/clients/email/`)

> 📧 Email client di-inject ke Services via DI Container. Akses melalui `s.EmailClient`.

**Struct:** `EmailClient`

**Methods:**

| Method | Description | Example |
|--------|-------------|---------|
| `IsConfigured()` | Check if SMTP is configured | `s.EmailClient.IsConfigured()` |
| `SendEmail(req)` | Send email with custom request | `s.EmailClient.SendEmail(req)` |
| `SendResetPasswordEmail(to, token, resetURL)` | Send reset password email | `s.EmailClient.SendResetPasswordEmail(email, token, url)` |

**DTOs:**

```go
// EmailRequest - untuk custom email
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
| `ResetPasswordEmail(token, resetURL, appName)` | HTML template untuk reset password email |

---

### **Creating New Helpers**

> 💡 **Rule of Thumb**: Jika function bisa digunakan di lebih dari 1 tempat, jadikan helper!

**Kriteria membuat helper baru:**
1. ✅ Operasi general (bukan business logic spesifik)
2. ✅ Bisa di-reuse di multiple services
3. ✅ Tidak bergantung pada dependencies (repo, db, etc)
4. ✅ Pure function (input → output, no side effects)

**Contoh membuat helper baru:**

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

**❌ JANGAN buat helper jika:**
- ❌ Hanya dipakai di 1 service
- ❌ Ada business logic spesifik
- ❌ Butuh dependency injection (db, repo, etc)
- ❌ Ada side effects (write to db, send email, etc)

---

### Testing

Untuk menambahkan testing, buat file `_test.go` di folder yang sama dengan file yang ditest:
```bash
# Contoh
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

### 1. Tambah Model Baru
```go
// internal/models/your_model.go
type YourModel struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    // ... fields
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// WAJIB: Define nama tabel
func (YourModel) TableName() string {
    return "your_table_name"
}
```

### 2. Tambah Repository

**PENTING**: Gunakan **GenericRepository** untuk semua operasi CRUD standar (Create, Read, Update, Delete, Find, Count, dll). Tambah method custom **hanya** untuk operasi yang tidak bisa di-handle generic, seperti:
- Query dengan JOIN kompleks (bukan preload biasa)
- Query dengan subquery khusus
- Operasi batch khusus
- Query dengan aggregasi khusus

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

// ✅ CONTOH CUSTOM METHOD YANG BENAR (JOIN kompleks)
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

// ❌ JANGAN LAKUKAN INI (gunakan generic repo langsung)
// func (r *YourRepository) FindByID(id uint) (*YourModel, error) {
//     return r.FindByID(r.DB, id) // Ini sudah ada di GenericRepository!
// }
```

### 3. Register Repository ke Container

Setelah membuat repository, **WAJIB** daftarkan ke `internal/repositories/00_repository.go` agar bisa digunakan di service:

```go
// internal/repositories/00_repository.go
type Repositories struct {
    TxManager *TransactionManager
    User      *UserRepository
    YourModel *YourRepository  // ← Tambah ini
}

func NewRepositories(db *gorm.DB) (*Repositories, error) {
    txManager := NewTransactionManager(db)
    userRepo := NewUserRepository(db)
    yourRepo := NewYourRepository(db)  // ← Inisialisasi

    return &Repositories{
        TxManager: txManager,
        User:      userRepo,
        YourModel: yourRepo,  // ← Register
    }, nil
}
```

### 4. Tambah Service Method (WAJIB ke struct Services)
```go
// internal/services/your_service.go
func (s *Services) YourMethod(ctx context.Context, param string) error {
    // business logic
    return nil
}
```

### 5. Tambah Handler Method (WAJIB ke struct Handlers)
```go
// internal/handlers/your_handler.go
func (h *Handlers) YourHandler(c *gin.Context) {
    // handle request
}
```

### 6. Tambah Routes
```go
// internal/routes/your_routes.go
func RegisterYourRoutes(r *gin.RouterGroup, handlers *handlers.Handlers) {
    your := r.Group("/your")
    {
        your.GET("/", handlers.YourHandler)
    }
}
```

### 7. Register di DI Container

Update `internal/di/container.go` untuk inject repository baru ke container:

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

Repository sudah otomatis tersedia di service melalui `s.repo.YourModel`.

## Security Notes

- ⚠️ **JWT_SECRET**: Ganti dengan secret key yang kuat di production
- ⚠️ **Database Password**: Jangan commit credentials ke git
- ⚠️ **GIN_MODE**: Set ke `release` di production
- ⚠️ **CORS**: Konfigurasi `ALLOWED_ORIGINS` sesuai domain yang diizinkan
- ⚠️ **TRUSTED_PROXIES**: Set jika menggunakan reverse proxy (nginx, load balancer)

## Documentation

- API documentation tersedia di `docs/` folder
- Environment configuration template di `.env.example`
