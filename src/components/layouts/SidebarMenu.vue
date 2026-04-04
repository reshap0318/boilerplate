<script setup lang="ts">
export interface IMenuItem {
  icon: unknown
  label: string
  to: string
}

defineProps<{
  appName: string
  menuItems: IMenuItem[]
  isOpen: boolean
}>()

const emit = defineEmits<{
  close: []
}>()
</script>

<template>
  <!-- Overlay for mobile -->
  <Teleport to="body">
    <div
      v-if="isOpen"
      class="fixed inset-0 bg-black/50 z-40 lg:hidden"
      @click="emit('close')"
    />
  </Teleport>

  <!-- Sidebar -->
  <aside
    :class="[
      'fixed top-0 left-0 z-50 h-full bg-white shadow-lg transition-transform duration-300 ease-in-out',
      'w-64',
      isOpen ? 'translate-x-0' : '-translate-x-full',
    ]"
  >
    <!-- Logo / Header -->
    <div class="flex items-center h-12 px-4 border-b border-gray-200">
      <router-link to="/" class="text-xl font-bold text-gray-800">
        {{ appName }}
      </router-link>
    </div>

    <!-- Menu -->
    <nav class="p-4 space-y-1">
      <router-link
        v-for="item in menuItems"
        :key="item.to"
        :to="item.to"
        class="flex items-center gap-3 px-4 py-3 rounded-lg text-gray-700 hover:bg-gray-100 hover:text-gray-900 transition-colors"
        active-class="bg-blue-50 text-blue-600"
        @click="emit('close')"
      >
        <component :is="item.icon" class="w-5 h-5" />
        <span class="font-medium">{{ item.label }}</span>
      </router-link>
    </nav>
  </aside>
</template>
