import type { RouteRecordRaw } from 'vue-router'
import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import DefaultLayout from '@/layouts/DefaultLayout.vue'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/pages/auth/LoginPage.vue'),
    meta: { guest: true },
  },
  {
    path: '/forgot-password',
    name: 'ForgotPassword',
    component: () => import('@/pages/auth/ForgotPasswordPage.vue'),
    meta: { guest: true },
  },
  {
    path: '/reset-password',
    name: 'ResetPassword',
    component: () => import('@/pages/auth/ResetPasswordPage.vue'),
    meta: { guest: true },
  },
  {
    path: '/',
    component: DefaultLayout,
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'Home',
        component: () => import('@/pages/HomePage.vue'),
      },
      {
        path: 'uam/permissions',
        name: 'Permissions',
        component: () => import('@/pages/uam/permissions/IndexPage.vue'),
      },
      {
        path: 'uam/roles',
        name: 'Roles',
        component: () => import('@/pages/uam/roles/IndexPage.vue'),
      },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

// Navigation guard
router.beforeEach((to) => {
  const authStore = useAuthStore()
  const token = authStore.token

  if (to.meta.requiresAuth && !token) {
    return { name: 'Login' }
  }

  if (to.meta.guest && token) {
    return { name: 'Home' }
  }
  return true
})

export default router
