import type { Directive, DirectiveBinding } from 'vue'
import { useAuthStore } from '@/stores/auth'

const vPermission: Directive = {
  mounted(el: HTMLElement, binding: DirectiveBinding) {
    checkPermission(el, binding)
  },
  updated(el: HTMLElement, binding: DirectiveBinding) {
    checkPermission(el, binding)
  },
}

function checkPermission(el: HTMLElement, binding: DirectiveBinding) {
  const authStore = useAuthStore()
  const userPermissions = authStore.user?.permissions?.map((p) => p.name) || []

  const requiredPermissions = Array.isArray(binding.value)
    ? binding.value
    : [binding.value]

  const hasAccess = requiredPermissions.every((perm: string) =>
    userPermissions.includes(perm)
  )

  if (!hasAccess) {
    el.parentNode?.removeChild(el)
  }
}

export default vPermission
