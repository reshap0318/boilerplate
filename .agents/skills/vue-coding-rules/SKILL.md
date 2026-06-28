---
name: vue-coding-rules
description: Vue Project coding rules and conventions — see references/ for code examples
---

## When to use me

Use this skill when:

- Creating or modifying Vue components, Pinia stores, or route configurations
- Writing TypeScript interfaces, types, or validation rules
- Implementing API calls, error handling, or form validation
- Styling components with Tailwind CSS
- Setting up barrel exports for stores or components
- Reviewing code for adherence to project conventions

Do NOT use this skill for:

- Backend code or API server implementation
- Non-Vue/TypeScript files

---

# Vue Coding Rules

## Project Structure

```
src/
├── assets/           # Static assets
├── components/       # Reusable components
│   ├── directives/   # Custom directives
│   ├── layouts/      # Layout components
│   └── utils/        # UI components (UiButton, UiCard, etc)
├── composables/      # Recomposition functions (useCrud, withFile)
├── helpers/          # Utility functions
├── layouts/          # Page layout wrappers
├── pages/            # Page/view components (route-level)
├── plugins/          # Vue plugins (axios, swal)
├── router/           # Vue Router config
├── stores/           # Pinia state management
└── main.ts           # Entry point
```

**Folder naming:** kebab-case, plural for features (`users/`, `roles/`)
**View files:** `IndexView.vue` as main view per folder

## Code Examples

All examples are in `references/` folder:

- `references/store-crud.ts` — Standard CRUD store with useCrud
- `references/store-with-file.ts` — Store with file upload (withFile HOF)
- `references/component.vue` — Component structure with script setup
- `references/types.ts` — Types file pattern (grouped by component)
- `references/barrel.ts` — Barrel export patterns
- `references/permission.ts` — Permission composable usage

---

## File Naming

| Type            | Pattern             | Example                                     |
| --------------- | ------------------- | ------------------------------------------- |
| Vue Components  | PascalCase          | `HelloWorld.vue`, `SidebarMenu.vue`         |
| UI Utilities    | `Ui` prefix         | `UiButton.vue`, `UiCard.vue`, `UiModal.vue` |
| Form Components | `Form` prefix       | `FormInput.vue`, `FormSelect.vue`           |
| Views           | `View` suffix       | `HomeView.vue`, `IndexView.vue`             |
| Modal/Form      | `FormModal.vue`     | `FormModal.vue`                             |
| TypeScript      | camelCase           | `storage.ts`, `axios.ts`                    |
| Stores          | singular, lowercase | `auth.ts`, `user.ts`, `role.ts`             |
| Interfaces      | `I` prefix          | `IUser`, `IApiResponse`                     |
| Type aliases    | `T` prefix          | `TLoadingKey`                               |

---

## Component Rules

### 1. Always use `<script setup lang="ts">`

See `references/component.vue` for full structure order.

### 2. Props & Emits — TypeScript syntax

```typescript
const props = defineProps<{ userId: number }>()
const emit = defineEmits<{ update: [id: number]; delete: [id: number] }>()
```

### 3. Component Refs

```typescript
const formModalRef = ref<InstanceType<typeof FormModal> | null>(null)
defineExpose({ openModal })
```

### 4. Classes Props Pattern

Use `classes` object prop, NOT multiple individual class props:

```typescript
// GOOD
interface UiCardClasses {
  wrapper?: string
  card?: string
  header?: string
}
const props = defineProps<{ classes?: UiCardClasses }>()

// BAD - Don't do this
const props = defineProps<{ wrapperClass?: string; cardClass?: string }>()
```

See `references/types.ts` for full pattern.

### 5. Types Organization

- Group related types in single `types.ts` per component folder
- If component only uses store interfaces, import from store — don't duplicate in `types.ts`

---

## Pinia Store Rules

### 1. Composable-First Pattern

All CRUD stores MUST use `useCrud` composable. See `references/store-crud.ts`.

### 2. File Upload — use `withFile` HOF

For stores with file uploads, wrap `useCrud` with `withFile`. See `references/store-with-file.ts`.

```typescript
const userCrud = withFile<IUser, IUserPayload>(crud, ['avatar'])
```

### 3. Store Naming

| Item     | Pattern                      | Example        |
| -------- | ---------------------------- | -------------- |
| File     | singular, lowercase          | `user.ts`      |
| Function | `use` + PascalCase + `Store` | `useUserStore` |
| Store ID | singular, lowercase string   | `'user'`       |

### 4. Loading State Pattern

```typescript
loading: Record<TLoadingKey, boolean> // keys: 'Index', 'Form', 'Delete'
// Access: store.loading.Index, store.loading.Form
```

### 5. No Function Redefinition

DO NOT redefine same function across stores. Import and reuse from original store.

### 6. Non-CRUD Stores

For stores that don't fit CRUD pattern (e.g., Profile), write manually but keep same structure.

---

## API & Axios Rules

### 1. API Response Interface

```typescript
interface IApiMetadata {
  total: number
  page: number
  page_size: number
  total_pages: number
}
interface IApiResponse<TData> {
  code: number
  message: string
  data: TData
  metadata?: IApiMetadata
}
```

### 2. API Calls

- All API calls from stores, NOT components
- Use `useCrud` for standard CRUD — handles loading, errors, swal automatically
- Use `withFile` for CRUD with file uploads
- Use try-catch-finally for custom API calls

### 3. Error Handling

The axios interceptor handles all standard HTTP errors (400, 401, 403, 404, 409, 422, 500, etc.) automatically via `swal`. Do **not** duplicate error handling for these cases.

**Default — let interceptor handle it (component `handleSubmit`):**
```typescript
try {
  await store.create()
  close()
} catch {} // interceptor already showed the error
```

**Custom — only when you need a specific message or action beyond what the interceptor does:**
```typescript
catch (error: any) {
  const message = error?.response?.data?.message || 'Default error message.'
  swal.error('Gagal', message)
  throw error
}
```
Use custom catch in store methods for non-CRUD API calls, or when the default interceptor message is not descriptive enough for the context.

---

## Routing Rules

- Lazy load all routes: `component: () => import('@/pages/...')`
- Group by layout using `children`
- Route meta: `requiresAuth: true` or `guest: true`
- URLs: kebab-case, plural resources (`/users`, `/uam/roles`)
- Add `permissions` meta for restricted routes — string format: `{resource}.{action}`

```typescript
{
  path: 'feature',
  name: 'Feature',
  component: () => import('@/pages/feature/IndexView.vue'),
  meta: { requiresAuth: true, permissions: ['feature.index'] },
},
```

---

## Sidebar Menu

Menu items are defined in `src/layouts/DefaultLayout.vue` inside the `menuItems` computed block. Items are filtered automatically by permission.

**Flat item:**
```typescript
{ icon: PhIconName, label: 'Feature', to: '/feature', permission: ['feature.index'] },
```

**Grouped item (parent with children):**
```typescript
{
  icon: PhIconName,
  label: 'Group',
  children: [
    { label: 'Feature', to: '/feature', permission: ['feature.index'] },
  ],
},
```

- Import icon from `@phosphor-icons/vue`
- `permission` must match the route's `permissions` meta value

---

## Permission & Access Control

### 1. Permission Composable

Use `usePermission()` composable for permission checks in components:

```typescript
import { usePermission } from '@/composables'
const permission = usePermission()
```

**Read `references/permission.ts` for the complete API and usage examples.**

### 2. Conditional Render Pattern

Use `usePermission` composable to conditionally render action buttons in IndexView:

```vue
<UiButton v-if="..." @click="openCreate">Tambah</UiButton>
<button v-if="..." @click="openEdit(item)">Edit</button>
<button v-if="..." @click="handleDelete(item.id)">Hapus</button>
```

**Always read `references/permission.ts` first to understand the full API before implementing permission logic.**

### 3. Frontend Permission is UI Only

**IMPORTANT:** Permission checks on frontend are for **UI/UX purposes only** (show/hide buttons, menu filtering).

**DO NOT** add permission validation before sending data to backend. The backend already handles permission validation and will return appropriate errors.

```typescript
// ❌ WRONG - Don't check permission before API call
async function handleSubmit() {
  if (!hasPermission('user.create')) {
    swal.error('Access Denied')
    return
  }
  await store.create() // Backend already validates permission
}

// ✅ CORRECT - Just call the API, let backend handle validation
async function handleSubmit() {
  await store.create() // Backend will reject if no permission
}
```

See `references/permission.ts` for full examples.

---

## Validation (Vuelidate)

```typescript
const v$ = useVuelidate(rules, store.form)
await v$.value.$validate() // before submission
v$.value.$reset() // on form open
```

**Create vs edit rules:** use `computed()` to override rules based on `isEdit`. Common case — password not required on edit:

```typescript
const dynamicRules = computed(() => {
  if (isEdit.value) {
    return { ...store.formRules, password: { minLength: minLength(6) } }
  }
  return store.formRules
})
```

**Custom Validation Rules**

Use `helpers.withMessage` for any validation logic beyond built-in validators. Custom rules belong in the **store**, defined after `useCrud` (or after `form`) so they can reference form state. The validator function is called lazily during validation — it always reads the current reactive value.

**Pattern 1 — Store with `useCrud`:** define extended `formRules` after `useCrud`, referencing `crud.form`:

```typescript
import { helpers, required, requiredIf } from '@vuelidate/validators'

const crud = useCrud<IUser, IUserPayload>({
  ...,
  formRules: {
    password: { required, minLength: minLength(6) },
    password_confirmation: {},
  },
})

// Extended formRules — defined after useCrud so crud.form is accessible
const formRules = {
  ...crud.formRules,

  // Field equality + conditional required (required only when password is filled)
  password_confirmation: {
    requiredIf: requiredIf(() => !!crud.form.password),
    sameAsPassword: helpers.withMessage(
      'Password tidak cocok',
      (value: string) => !crud.form.password || value === crud.form.password,
    ),
  },

  // Array minimum items
  roles: {
    minItems: helpers.withMessage(
      'Minimal pilih 1 role',
      (value: number[]) => value.length > 0,
    ),
  },
}

// Expose formRules (overrides the one from crud spread)
return { ...crud, formRules }
```

**Pattern 2 — Store without `useCrud`:** define `form` first, then `formRules` referencing it:

```typescript
import { helpers, required, requiredIf } from '@vuelidate/validators'

const form = reactive({ password: '', password_confirmation: '', contact_method: '', phone: '' })

const formRules = {
  // Field equality + conditional required
  password_confirmation: {
    requiredIf: requiredIf(() => !!form.password),
    sameAsPassword: helpers.withMessage(
      'Password tidak cocok',
      (value: string) => !form.password || value === form.password,
    ),
  },

  // Conditional required
  phone: {
    conditionalRequired: helpers.withMessage(
      'Nomor telepon wajib diisi jika metode kontak adalah telepon',
      (value: string) => form.contact_method !== 'phone' || !!value,
    ),
  },

  // Custom format (e.g. permission name: resource.action)
  name: {
    format: helpers.withMessage(
      'Format harus resource.action (contoh: user.index)',
      (value: string) => /^[a-z]+\.[a-z]+$/.test(value),
    ),
  },
}
```

---

**Modal open — always clear server errors before showing:**
```typescript
import { useFormError } from '@/composables/useFormError'
const formError = useFormError()

function show(data?: ...) {
  // populate form...
  v$.value.$reset()
  formError.clear()
  isVisible.value = true
}
```

**`handleSubmit` — call `close()` only on success, never in `finally`:**
```typescript
async function handleSubmit() {
  const isValid = await v$.value.$validate()
  if (!isValid) return

  try {
    if (isEdit.value) await store.update(store.form.id)
    else await store.create()
    close() // only reached on success
  } catch {} // error already handled by axios interceptor
}
```

---

## TypeScript Rules

- Always `lang="ts"` in Vue components
- Avoid `any` — use proper interfaces or `unknown`
- Define interfaces in stores when shared, locally when component-specific
- Use generics for reusable functions

---

## Styling (Tailwind CSS)

- Utility-first — avoid custom CSS
- Mobile-first responsive design
- **Always implement at minimum `sm` and `md` breakpoints**
- Always use `scoped` in `<style>` tags

---

## State Management

| Use Stores               | Use Local State                    |
| ------------------------ | ---------------------------------- |
| Shared across components | Local to one component             |
| Server state (API data)  | UI state specific to one component |
| Authentication state     | Temporary form state               |
| Complex form state       |                                    |

---

## Error Handling & Notifications

```typescript
import swal from '@/plugins/swal'
swal.success('Berhasil', 'Message')
swal.error('Gagal', 'Message')
await swal.warning('Title', 'Message') // confirmation
```

- User-facing messages in the app's target language (e.g., Indonesian)
- Use `console.error()` for developer debugging

---

## Barrel Exports

Use `index.ts` files for cleaner imports. See `references/barrel.ts`.

```typescript
// stores/index.ts — export stores + types
// composables/index.ts — export composables
// components/utils/index.ts — export UI components + types
```

**Rules:**

- Export both module and types in barrel
- Use explicit named exports, avoid `export * from`
- Keep alphabetically sorted when possible

---

## Environment Variables

- Prefix: `VITE_` (e.g., `VITE_API_BASE_URL`)
- Access: `import.meta.env.VITE_*`
- Provide fallbacks: `import.meta.env.VITE_APP_NAME || 'Default'`

---

## Anti-Patterns

| ❌ WRONG                             | ✅ CORRECT                                     |
| ------------------------------------ | ---------------------------------------------- |
| API calls in components              | API calls in stores                            |
| `any` type without justification     | Proper interfaces or `unknown`                 |
| Mutate props directly                | Use emits or local state                       |
| Index as key in `v-for`              | Unique identifiers                             |
| Business logic in templates          | Computed properties or methods                 |
| `Page` suffix for views              | `View` suffix                                  |
| Redefine functions across stores     | Import and reuse from original store           |
| Multiple individual class props      | `classes` object prop                          |
| Skip `sm`/`md` breakpoints           | Always implement responsive                    |
| Manual CRUD logic in stores          | `useCrud` composable                           |
| Duplicate file upload logic          | `withFile` HOF                                 |
| `createEntity`/`updateEntity` naming | `create`/`update` or `createForm`/`updateForm` |

---

## Pre-Implementation Checklist

- [ ] Component uses `<script setup lang="ts">`
- [ ] Props/Emits use TypeScript syntax
- [ ] Store uses `useCrud` composable (or `withFile` for uploads)
- [ ] Store naming: singular file, `useXxxStore` function, `'xxx'` ID
- [ ] Interfaces prefixed with `I`, types with `T`
- [ ] Types grouped in `types.ts` per component folder
- [ ] Barrel exports updated for new stores/components
- [ ] API calls in stores, NOT components
- [ ] Error messages in the app's target language (e.g., Indonesian)
- [ ] Responsive: `sm` and `md` breakpoints implemented
- [ ] CRUD buttons protected with `usePermission` composable
- [ ] Lint clean: `yarn lint`
- [ ] Format clean: `yarn format`
