<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { PhList, PhHouse, PhGear } from '@phosphor-icons/vue'
import type { IMenuItem } from '@/components/layouts/SidebarMenu.vue'
import SidebarMenu from '@/components/layouts/SidebarMenu.vue'
import TopBar from '@/components/layouts/TopBar.vue'

const sidebarOpen = ref(false)
const appName = import.meta.env.VITE_APP_NAME || 'Admin'

const menuItems: IMenuItem[] = [
  { icon: PhHouse, label: 'Dashboard', to: '/' },
  { icon: PhList, label: 'Data Master', to: '/data' },
  { icon: PhGear, label: 'Settings', to: '/settings' },
]

onMounted(() => {
  sidebarOpen.value = window.innerWidth >= 1024
})
</script>

<template>
  <div class="min-h-screen bg-gray-100">
    <!-- Sidebar -->
    <SidebarMenu
      :app-name="appName"
      :menu-items="menuItems"
      :is-open="sidebarOpen"
      @close="sidebarOpen = false"
    />

    <!-- Main Content -->
    <div
      class="transition-all duration-300 ease-in-out"
      :class="sidebarOpen ? 'lg:ml-64' : ''"
    >
      <!-- Top Bar -->
      <TopBar
        title="Admin Panel"
        :show-hamburger="true"
        @toggle-sidebar="sidebarOpen = !sidebarOpen"
      >
        <template #menu-icon>
          <PhList class="w-6 h-6" />
        </template>
      </TopBar>

      <!-- Page Content -->
      <main class="p-4">
        <router-view />
      </main>
    </div>
  </div>
</template>
