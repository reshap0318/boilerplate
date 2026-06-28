---
name: vue-development-guide
description: Step-by-step development workflow and CRUD implementation guide for Vue project — from planning to push. Use when creating new features, implementing CRUD operations, or when user mentions stores, components, composables, or routes.
---

## When to use me

Use this skill when:
- Implementing a new CRUD feature from scratch (Store → View → Form → Route → Sidebar)
- Following the mandatory development workflow
- You need a pre-push checklist before committing code

> 💡 **If `vue-coding-rule` skill is available**, use it for coding conventions, component patterns, and full code examples. This guide focuses on workflow and implementation order.

---

# Vue Development Guide

## 🛠️ CRUD Implementation Flow

### 📋 Step-by-Step Order
Every new CRUD feature MUST follow this order:

| Step | Layer | File Path | Purpose |
|------|-------|-----------|---------|
| 1 | Store | `src/stores/{feature}.ts` | Pinia store with `useCrud` + `withFile` |
| 2 | View | `src/pages/{feature}/IndexView.vue` | List page with table/grid + pagination + permission-guarded action buttons |
| 3 | Form/Modal | `src/pages/{feature}/FormModal.vue` | Create/Edit form with validation |
| 4 | Route | `src/router/index.ts` | Register route with auth + permission guard |
| 5 | Sidebar | `src/layouts/DefaultLayout.vue` | Add menu item to sidebar |
| 6 | Barrel Exports | `src/stores/index.ts` | Export store + types |

> For coding rules (naming, component structure, validation, styling), see `vue-coding-rule` skill.

---

### Step 1: Store

**File:** `src/stores/{feature}.ts`

Create Pinia store with `useCrud` composable. Add `withFile` HOF if file upload needed. See `vue-coding-rule` for full patterns.

---

### Step 2: View Component

**File:** `src/pages/{feature}/IndexView.vue`

Create list page with data table/grid, pagination, search, and action buttons (Create, Edit, Delete). Call store's `fetchAll()` on mount. Protect action buttons with `usePermission` if required — see `vue-coding-rule` for the pattern.

---

### Step 3: Form/Modal Component

**File:** `src/pages/{feature}/FormModal.vue`

Create modal form for Create/Edit. Use Vuelidate for validation. Two patterns are mandatory — see `vue-coding-rule` for the full code:
- Call `formError.clear()` inside `show()` every time the modal opens
- Call `close()` only on success inside `handleSubmit`, never in `finally`

---

### Step 4: Route

**File:** `src/router/index.ts`

Register route under the authenticated layout with lazy loading. Use plural kebab-case URL (e.g., `/users`, `/uam/roles`). Add `requiresAuth: true` and `permissions: ['feature.index']` meta — see `vue-coding-rule` for the full pattern.

---

### Step 5: Sidebar

**File:** `src/layouts/DefaultLayout.vue`

Add a new item to the `menuItems` computed array. Set `permission` to match the route's `permissions` meta value — items are filtered automatically. See `vue-coding-rule` for flat vs grouped item patterns.

---

### Step 6: Barrel Exports

**File:** `src/stores/index.ts`

Export the new store and its types from the barrel file for clean imports.

---

## ✅ Pre-Push Checklist
- [ ] Store uses `useCrud` composable (or `withFile` for uploads)
- [ ] Store naming: singular file, `useXxxStore` function, `'xxx'` ID
- [ ] View uses `<script setup lang="ts">`
- [ ] Props/Emits use TypeScript syntax
- [ ] Form uses Vuelidate with `await v$.value.$validate()` before submit
- [ ] `formError.clear()` called in modal `show()` method
- [ ] `handleSubmit` calls `close()` only on success, not in `finally`
- [ ] Route registered with `requiresAuth: true` and `permissions: [...]`
- [ ] Action buttons (Create/Edit/Delete) protected with `usePermission` if required — see `vue-coding-rule`
- [ ] Sidebar menu item added to `DefaultLayout.vue` with matching `permission`
- [ ] Barrel exports updated for new store
- [ ] Responsive: `sm` and `md` breakpoints implemented
- [ ] Lint clean: `yarn lint`
- [ ] Format clean: `yarn format`
- [ ] Build success: `yarn build`
