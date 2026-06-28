// Permission Composable Usage
import { usePermission } from '@/composables'

const { hasPermission, hasAnyPermission, hasAllPermissions } = usePermission()

// Single permission
if (hasPermission('user.create')) {
  // proceed
}

// OR logic (any of)
if (hasAnyPermission(['role.index', 'permission.index'])) {
  // proceed
}

// AND logic (all of)
if (hasAllPermissions(['user.edit', 'user.delete'])) {
  // proceed
}

// In template (conditional render with v-if)
// <button v-if="hasPermission('user.create')">Create</button>
// <button v-if="hasAnyPermission(['user.edit', 'user.delete'])">Actions</button>

// IndexView pattern (protect CRUD buttons)
// <UiButton v-if="hasPermission('user.create')" @click="openCreate">Tambah</UiButton>
// <button v-if="hasPermission('user.edit')" @click="openEdit(item)">Edit</button>
// <button v-if="hasPermission('user.delete')" @click="handleDelete(item.id)">Hapus</button>

// Router pattern (route-level protection)
// { path: 'users', meta: { requiresAuth: true, permissions: ['user.index'] } }

// Menu pattern (auto-filtered by permission)
// { icon: PhUsers, label: 'Users', to: '/users', permission: ['user.index'] }
