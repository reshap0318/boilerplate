<script setup lang="ts">
export interface IMenuItem {
  icon: unknown
  label: string
  to?: string
  children?: { label: string; to: string }[]
}

defineProps<{
  appName: string
  menuItems: IMenuItem[]
  isOpen: boolean
}>()

const emit = defineEmits<{
  close: []
}>()

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
      <template v-for="(item, index) in menuItems" :key="index">
        <router-link
          v-if="!item.children && item.to"
          :to="item.to"
          class="flex items-center gap-3 px-4 py-3 rounded-lg text-gray-700 hover:bg-gray-100 hover:text-gray-900 transition-colors"
          active-class="bg-blue-50 text-blue-600"
          @click="handleItemClick"
        >
          <component :is="item.icon" class="w-5 h-5" />
          <span class="font-medium">{{ item.label }}</span>
        </router-link>

        <div v-else class="space-y-1">
          <div class="flex items-center gap-3 px-4 py-3 rounded-lg text-gray-700 font-medium">
            <component :is="item.icon" class="w-5 h-5" />
            <span>{{ item.label }}</span>
          </div>
          <div class="pl-11 space-y-1">
            <router-link
              v-for="child in item.children"
              :key="child.to"
              :to="child.to"
              class="block px-4 py-2 rounded-lg text-sm text-gray-600 hover:bg-gray-100 hover:text-gray-900 transition-colors"
              active-class="bg-blue-50 text-blue-600 font-medium"
              @click="handleItemClick"
            >
              {{ child.label }}
            </router-link>
          </div>
        </div>
      </template>
    </nav>
  </aside>
</template>
