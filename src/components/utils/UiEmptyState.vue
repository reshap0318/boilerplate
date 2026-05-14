<script setup lang="ts">
import type { Component } from 'vue'

interface Props {
  icon?: Component
  title: string
  description?: string
  variant?: 'default' | 'compact'
}

withDefaults(defineProps<Props>(), {
  icon: undefined,
  description: '',
  variant: 'default',
})
</script>

<template>
  <div
    :class="[
      'text-center bg-white rounded-2xl border-2 border-dashed border-gray-300 shadow-sm',
      variant === 'compact' ? 'py-12' : 'py-24',
    ]"
  >
    <div
      v-if="icon"
      :class="[
        'inline-flex items-center justify-center rounded-full mb-4',
        variant === 'compact' ? 'w-16 h-16' : 'w-20 h-20',
        'bg-indigo-100 text-indigo-600',
      ]"
    >
      <component
        :is="icon"
        :class="variant === 'compact' ? 'w-8 h-8' : 'w-10 h-10'"
      />
    </div>
    <h3
      :class="[
        'font-semibold text-gray-900 mb-2',
        variant === 'compact' ? 'text-xl' : 'text-2xl',
      ]"
    >
      {{ title }}
    </h3>
    <p
      v-if="description"
      class="text-gray-600 mb-6"
    >
      {{ description }}
    </p>
    <slot />
  </div>
</template>
