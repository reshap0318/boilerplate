<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { PhUser, PhSignOut, PhCaretDown, PhGear } from '@phosphor-icons/vue'
import { useAuthStore } from '@/stores/auth'

defineProps<{
  title: string
  showHamburger: boolean
}>()

const emit = defineEmits<{
  toggleSidebar: []
}>()

const router = useRouter()
const authStore = useAuthStore()
const profileOpen = ref(false)

const toggleProfile = () => {
  profileOpen.value = !profileOpen.value
}

const closeProfile = () => {
  profileOpen.value = false
}

const handleLogout = () => {
  authStore.logout()
  router.push({ name: 'Login' })
}
</script>

<template>
  <header class="sticky top-0 z-30 bg-white shadow-sm">
    <div class="flex items-center justify-between px-4 h-12">
      <!-- Left: Hamburger + Title -->
      <div class="flex items-center gap-4">
        <button
          v-if="showHamburger"
          class="p-2 rounded-md hover:bg-gray-100"
          @click="emit('toggleSidebar')"
        >
          <slot name="menu-icon" />
        </button>
        <h2 class="text-lg font-semibold text-gray-800">{{ title }}</h2>
      </div>

      <!-- Right: Profile Dropdown -->
      <div class="relative">
        <button
          type="button"
          class="group flex items-center gap-2 p-2 rounded-lg hover:bg-gray-100 transition-all duration-200 hover:shadow-sm"
          @click.stop="toggleProfile"
        >
          <div class="w-8 h-8 rounded-full bg-blue-500 flex items-center justify-center transition-transform duration-200 group-hover:scale-110">
            <PhUser class="w-5 h-5 text-white" />
          </div>
          <span class="hidden sm:inline-block text-sm font-medium text-gray-700 group-hover:text-gray-900">
            {{ authStore.user?.name || 'User' }}
          </span>
          <PhCaretDown class="w-4 h-4 text-gray-500" />
        </button>

        <!-- Dropdown Menu -->
        <div
          v-if="profileOpen"
          v-click-outside="closeProfile"
          class="absolute right-0 mt-2 w-56 bg-white rounded-lg shadow-lg border border-gray-200 py-1"
          @click.stop
        >
          <!-- User Info -->
          <div class="px-4 py-3 border-b border-gray-100">
            <p class="text-sm font-semibold text-gray-800">{{ authStore.user?.name }}</p>
            <p class="text-xs text-gray-500 truncate">{{ authStore.user?.email }}</p>
          </div>

          <!-- Menu Items -->
          <div class="py-1">
            <router-link
              to="/profile"
              class="flex items-center gap-3 px-4 py-2.5 text-sm text-gray-700 hover:bg-blue-50 hover:text-blue-600 transition-colors"
              @click="closeProfile"
            >
              <PhUser class="w-4 h-4" />
              <span>Profile</span>
            </router-link>
            <router-link
              to="/settings"
              class="flex items-center gap-3 px-4 py-2.5 text-sm text-gray-700 hover:bg-blue-50 hover:text-blue-600 transition-colors"
              @click="closeProfile"
            >
              <PhGear class="w-4 h-4" />
              <span>Settings</span>
            </router-link>
            <hr class="my-1 border-gray-100" />
            <button
              class="flex items-center gap-3 w-full px-4 py-2.5 text-sm text-red-600 hover:bg-red-50 hover:text-red-700 transition-colors"
              @click="handleLogout"
            >
              <PhSignOut class="w-4 h-4" />
              <span>Logout</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </header>
</template>
