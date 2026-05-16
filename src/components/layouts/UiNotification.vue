<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { PhBell, PhCheck, PhTrash, PhCircle } from '@phosphor-icons/vue'

export interface INotificationItem {
  id: number
  title: string
  message: string
  time: string
  isRead: boolean
  type?: 'info' | 'success' | 'warning' | 'error'
}

const isOpen = ref(false)

const notifications = ref<INotificationItem[]>([
  {
    id: 1,
    title: 'New User Registered',
    message: 'A new user has registered to the system.',
    time: '2 minutes ago',
    isRead: false,
    type: 'info',
  },
  {
    id: 2,
    title: 'Role Updated',
    message: 'Admin role has been updated successfully.',
    time: '1 hour ago',
    isRead: false,
    type: 'success',
  },
  {
    id: 3,
    title: 'Permission Warning',
    message: 'Some permissions may conflict with existing rules.',
    time: '3 hours ago',
    isRead: true,
    type: 'warning',
  },
])

const unreadCount = computed(() => notifications.value.filter((n) => !n.isRead).length)

function closeOthers() {
  window.dispatchEvent(new CustomEvent('close-dropdown', { detail: 'notification' }))
}

function close() {
  isOpen.value = false
}

function toggle() {
  if (isOpen.value) {
    close()
  } else {
    closeOthers()
    isOpen.value = true
  }
}

function handleMarkRead(id: number) {
  const notification = notifications.value.find((n) => n.id === id)
  if (notification) {
    notification.isRead = true
  }
}

function handleMarkAllRead() {
  notifications.value.forEach((n) => {
    n.isRead = true
  })
}

function handleDelete(id: number) {
  notifications.value = notifications.value.filter((n) => n.id !== id)
}

function handleDeleteAll() {
  notifications.value = []
}

function handleClick(notification: INotificationItem) {
  if (!notification.isRead) {
    handleMarkRead(notification.id)
  }
}

function typeIcon(type?: string) {
  switch (type) {
    case 'success':
      return PhCheck
    case 'warning':
      return PhCircle
    case 'error':
      return PhCircle
    default:
      return PhCircle
  }
}

function typeColor(type?: string) {
  switch (type) {
    case 'success':
      return 'bg-green-100 text-green-600'
    case 'warning':
      return 'bg-yellow-100 text-yellow-600'
    case 'error':
      return 'bg-red-100 text-red-600'
    default:
      return 'bg-blue-100 text-blue-600'
  }
}

function handleOutsideClose(event: Event) {
  const detail = (event as CustomEvent).detail
  if (detail !== 'notification') {
    close()
  }
}

onMounted(() => {
  window.addEventListener('close-dropdown', handleOutsideClose)
})

onUnmounted(() => {
  window.removeEventListener('close-dropdown', handleOutsideClose)
})
</script>

<template>
  <div class="relative">
    <button
      type="button"
      class="relative p-2 rounded-full hover:bg-gray-100 transition-colors"
      @click.stop="toggle"
    >
      <PhBell class="h-5 w-5 text-gray-600" />
      <span v-if="unreadCount > 0" class="absolute top-1 right-1 w-2 h-2 bg-red-500 rounded-full" />
    </button>

    <Transition
      enter-active-class="transition ease-out duration-200"
      enter-from-class="opacity-0 scale-95 -translate-y-1"
      enter-to-class="opacity-100 scale-100 translate-y-0"
      leave-active-class="transition ease-in duration-150"
      leave-from-class="opacity-100 scale-100 translate-y-0"
      leave-to-class="opacity-0 scale-95 -translate-y-1"
    >
      <div
        v-if="isOpen"
        v-click-outside="close"
        class="absolute right-0 mt-2 w-80 bg-white rounded-2xl shadow-xl border border-gray-100 overflow-hidden"
        @click.stop
      >
        <!-- Header -->
        <div class="px-4 py-3 bg-gradient-to-br from-gray-50 to-white border-b border-gray-100">
          <div class="flex items-center justify-between">
            <h3 class="text-sm font-semibold text-gray-900">
              Notifications
              <span
                v-if="unreadCount > 0"
                class="ml-1.5 px-2 py-0.5 text-xs font-medium bg-red-100 text-red-600 rounded-full"
              >
                {{ unreadCount }}
              </span>
            </h3>
            <div class="flex items-center gap-1">
              <button
                v-if="notifications.length > 0"
                class="p-1.5 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100 transition-colors"
                title="Mark all as read"
                @click="handleMarkAllRead"
              >
                <PhCheck class="h-4 w-4" />
              </button>
              <button
                v-if="notifications.length > 0"
                class="p-1.5 rounded-lg text-gray-400 hover:text-red-600 hover:bg-red-50 transition-colors"
                title="Clear all"
                @click="handleDeleteAll"
              >
                <PhTrash class="h-4 w-4" />
              </button>
            </div>
          </div>
        </div>

        <!-- Notification List -->
        <div class="max-h-96 overflow-y-auto">
          <div v-if="notifications.length === 0" class="px-4 py-8 text-center">
            <PhBell class="h-8 w-8 text-gray-300 mx-auto mb-2" />
            <p class="text-sm text-gray-500">No notifications</p>
          </div>

          <template v-else>
            <button
              v-for="notification in notifications"
              :key="notification.id"
              :class="[
                'w-full px-4 py-3 text-left border-b border-gray-50 hover:bg-gray-50 transition-colors',
                !notification.isRead && 'bg-blue-50/50',
              ]"
              @click="handleClick(notification)"
            >
              <div class="flex items-start gap-3">
                <!-- Icon -->
                <div
                  :class="[
                    'w-8 h-8 rounded-lg flex items-center justify-center shrink-0',
                    typeColor(notification.type),
                  ]"
                >
                  <component :is="typeIcon(notification.type)" class="h-4 w-4" />
                </div>

                <!-- Content -->
                <div class="flex-1 min-w-0">
                  <div class="flex items-start justify-between gap-2">
                    <p
                      :class="[
                        'text-sm truncate',
                        !notification.isRead
                          ? 'font-semibold text-gray-900'
                          : 'font-medium text-gray-700',
                      ]"
                    >
                      {{ notification.title }}
                    </p>
                    <button
                      class="p-1 rounded text-gray-300 hover:text-red-500 hover:bg-red-50 transition-colors shrink-0"
                      @click.stop="handleDelete(notification.id)"
                    >
                      <PhTrash class="h-3.5 w-3.5" />
                    </button>
                  </div>
                  <p class="mt-0.5 text-xs text-gray-500 line-clamp-2">
                    {{ notification.message }}
                  </p>
                  <p class="mt-1 text-xs text-gray-400">
                    {{ notification.time }}
                  </p>
                </div>

                <!-- Unread indicator -->
                <span
                  v-if="!notification.isRead"
                  class="w-2 h-2 bg-blue-500 rounded-full mt-2 shrink-0"
                />
              </div>
            </button>
          </template>
        </div>
      </div>
    </Transition>
  </div>
</template>
