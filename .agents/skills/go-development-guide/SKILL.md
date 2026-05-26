---
name: go-development-guide
description: Step-by-step development workflow and CRUD implementation guide for Go project — from planning to push. Use when creating new features, implementing CRUD operations, adding permission-based access control, or when user mentions models, DTOs, repositories, services, handlers, or routes.
---

## When to use me

Use this skill when:
- Implementing a new CRUD feature from scratch (Model → DTO → Repository → Service → Handler → Routes)
- Following the mandatory development workflow (plan → issue → branch → PR)
- Adding scoped permission patterns (satker/unit restrictions) to existing features
- You need a pre-push checklist before committing code

> 💡 **If `go-coding-rules` skill is available**, use it for coding conventions, helpers, and full code examples. This guide focuses on workflow and implementation order.

---

# Go Development Guide

## 🔄 Development Workflow (MANDATORY)

Every new feature or bug fix **MUST** follow this flow:

### Phase 1: Planning
1. **Create `plan.md`** in root directory:
   - Feature/bug description
   - Technical analysis (files to create/modify)
   - Implementation steps
   - Impact on existing code
2. **User Approval**: Wait for user confirmation.
3. **Create GitHub Issue** with format: `[BE] [Type] Title`
   - `[Feat]`, `[Bug]`, `[Fix]`, `[Refactor]`, `[Chore]`

### Phase 2: Implementation
1. **Create Branch**: `git checkout -b feat/feature-name`
2. **Implement**: Based on approved plan and issue.
3. **Create Pull Request**: Reference issue in description (`Closes #<id>`).
4. **User Approval**: Summary of changes → approval → merge → close issue.

---

## 🛠️ CRUD Implementation Flow

### 📋 Step-by-Step Order
Every new CRUD feature MUST follow this order:

| Step | Layer | File Path | Purpose |
|------|-------|-----------|---------|
| 1 | Model | `internal/models/{feature}.go` | GORM model + TableName() |
| 2 | DTO | `internal/dtos/{feature}_dto.go` | Request/Response structs + converter |
| 3 | Repository | `internal/repositories/{feature}_repository.go` | Extend GenericRepository |
| 4 | Register Repo | `internal/repositories/00_repository.go` | Add to Repositories struct |
| 5 | Service | `internal/services/{feature}_service.go` | Business logic + logging + notification |
| 6 | Handler | `internal/handlers/{feature}_handler.go` | HTTP handler + validation |
| 7 | Routes | `internal/routes/{feature}_route.go` | Route registration + middleware |
| 8 | Register Routes | `cmd/api/main.go` | Wire routes into app |

> For coding rules (naming, transactions, logging, notification, helpers), see `go-coding-rules` skill.

---

### Step 1: Model

**File:** `internal/models/{feature}.go`

Define GORM model with struct tags. See `go-coding-rules` for conventions.

---

### Step 2: DTOs

**File:** `internal/dtos/{feature}_dto.go`

Create Request struct, Response DTO, and converter functions (`To{Feature}DTO`, `To{Feature}DTOList`). See `go-coding-rules` for conventions.

---

### Step 3: Repository

**File:** `internal/repositories/{feature}_repository.go`

Extend `GenericRepository[Model]`. Only add custom methods for complex queries. See `go-coding-rules` for conventions.

---

### Step 4: Register Repository

**File:** `internal/repositories/00_repository.go`

Add repository field to `Repositories` struct and initialize in `NewRepositories()`.

---

### Step 5: Service

**File:** `internal/services/{feature}_service.go`

Implement CRUD methods on `(s *Services)`. Write ops use transaction + notification. Read ops return DTOs. See `go-coding-rules` for full patterns.

---

### Step 6: Handler

**File:** `internal/handlers/{feature}_handler.go`

Implement HTTP handlers on `(h *Handlers)`. Bind JSON, validate, call service, return response. See `go-coding-rules` for full patterns.

---

### Step 7: Routes

**File:** `internal/routes/{feature}_route.go`

Register routes with permission middleware via `middleware.RequirePermission(acc, "permission.name")`.

---

### Step 8: Register Routes

**File:** `cmd/api/main.go`

Add route registration function call to the `protected` group (after JWT middleware).

---

## ✅ Pre-Push Checklist
- [ ] Model has `TableName()` method
- [ ] DTO variables use feature prefix
- [ ] Request DTO: merged for simple features OR separate (`{Feature}CreateRequest` + `{Feature}UpdateRequest`) for complex features
- [ ] Service functions use feature prefix (`{Feature}{Action}`)
- [ ] Handler functions use feature prefix (`{Feature}{Action}`)
- [ ] Write operations use `TxManager.WithinTransaction()` or `WithinTransactionWithResult()`
- [ ] Read operations use `nil` (NOT `.DB`)
- [ ] Logging on Create/Update/Delete
- [ ] Notification created inside transaction for Create/Update/Delete
- [ ] Repository registered in `00_repository.go`
- [ ] Routes registered in `cmd/api/main.go`
- [ ] Build success (`go build ./...`)
- [ ] Vet clean (`go vet ./...`)
