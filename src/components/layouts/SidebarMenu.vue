<script setup lang="ts">
import { ref } from 'vue'
import SidebarMenuItem from './SidebarMenuItem.vue'

export interface IMenuItem {
  icon?: unknown
  label: string
  to?: string
  children?: IMenuItem[]
}

defineProps<{
  appName: string
  menuItems: IMenuItem[]
  isOpen: boolean
}>()

const emit = defineEmits<{
  close: []
}>()

const expandedGroups = ref<Set<string>>(new Set())

function toggleGroup(key: string) {
  if (expandedGroups.value.has(key)) {
    expandedGroups.value.delete(key)
  } else {
    expandedGroups.value.add(key)
  }
}

const handleItemClick = () => {
  if (window.innerWidth < 1024) {
    emit('close')
  }
}
</script>

<template>
  <!-- Overlay for mobile -->
  <Teleport to="body">
    <div
      v-if="isOpen"
      class="fixed inset-0 bg-black/50 backdrop-blur-sm z-40 lg:hidden"
      @click="emit('close')"
    />
  </Teleport>

  <!-- Sidebar -->
  <aside
    :class="[
      'fixed top-0 left-0 z-50 h-full w-64 transition-transform duration-300 ease-in-out',
      'bg-linear-to-b from-slate-900 via-slate-800 to-slate-900',
      'shadow-2xl border-r border-white/5',
      isOpen ? 'translate-x-0' : '-translate-x-full',
    ]"
  >
    <!-- Logo / Header -->
    <div class="flex items-center h-12 px-5 border-b border-white/10">
      <router-link to="/" class="flex items-center group">
        <span
          class="text-lg font-bold bg-linear-to-r from-white to-slate-300 bg-clip-text text-transparent"
        >
          {{ appName }}
        </span>
      </router-link>
    </div>

    <!-- Menu -->
    <nav class="p-3 space-y-1 mt-2 overflow-y-auto" style="max-height: calc(100vh - 3.5rem)">
      <SidebarMenuItem
        v-for="(item, index) in menuItems"
        :key="index"
        :item="item"
        :depth="0"
        :expanded-groups="expandedGroups"
        @toggle-group="toggleGroup"
        @item-click="handleItemClick"
      />
    </nav>
  </aside>
</template>

<style scoped>
.sidebar-item-active {
  box-shadow: inset 3px 0 0 #3b82f6;
}
</style>
