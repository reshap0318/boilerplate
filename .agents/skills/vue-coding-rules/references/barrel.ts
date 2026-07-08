// ============================================================
// Barrel Exports — PATTERN EXAMPLE ONLY, not an exhaustive list.
// The real barrels (src/stores/index.ts, src/composables/index.ts,
// src/components/utils/index.ts) will have more entries than shown
// here — always read the real file before adding a new export.
// ============================================================

// ============================================================
// Barrel Exports — stores/index.ts
// ============================================================
export { useAuthStore } from './auth'
export type {
  ILoginPayload,
  ILoginResponse,
} from './auth'

export { useUserStore } from './user'
export type { IUser, IUserPayload } from './user'

export { useRoleStore } from './role'
export type { IRole, IRolePayload } from './role'

export { usePermissionStore } from './permission'
export type { IPermission, IPermissionPayload } from './permission'

// Usage:
// import { useUserStore, useRoleStore } from '@/stores'
// import type { IUser } from '@/stores'

// ============================================================
// Barrel Exports — composables/index.ts
// ============================================================
export { useCrud } from './useCrud'
export type { UseCrudOptions } from './useCrud'
export { withFile } from './withFile'

// Usage:
// import { useCrud, withFile } from '@/composables'

// ============================================================
// Barrel Exports — components/utils/index.ts
// ============================================================
export { default as UiButton } from './UiButton.vue'
export { default as UiCard } from './UiCard.vue'
export { default as UiModal } from './UiModal.vue'
export { default as UiPagination } from './UiPagination.vue'
export { default as UiEmptyState } from './UiEmptyState.vue'
export { default as FormInput } from './FormInput.vue'
export { default as FormSelect } from './FormSelect.vue'
export { default as FormPassword } from './FormPassword.vue'
export { default as FormAvatar } from './FormAvatar.vue'
export { default as FormFile } from './FormFile.vue'

// Export types
export type { UiCardClasses, UiCardProps } from './types'
export type { UiModalProps } from './types'

// Usage:
// import { UiButton, UiCard, FormInput } from '@/components/utils'
