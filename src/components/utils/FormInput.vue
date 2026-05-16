<script setup lang="ts">
import { computed } from 'vue'
import { getErrorMessage } from '@/helpers/vuelidate'
import type { FormInputProps } from './types'

interface ValidationLike {
  $error: boolean
  $errors: Array<{ $message: string }>
}

const resolveMessage = (msg: string | { value: string }): string => {
  return typeof msg === 'string' ? msg : msg.value
}

const props = withDefaults(defineProps<FormInputProps>(), {
  type: 'text',
  label: '',
  placeholder: '',
  validation: undefined,
  leadingIcon: undefined,
  trailingIcon: undefined,
  iconSize: 20,
  classes: () => ({}),
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const hasError = computed(() => (props.validation as ValidationLike)?.$error ?? false)
const errorMessage = computed(() => {
  if (!props.validation) return ''
  const raw = getErrorMessage(props.validation as ValidationLike)
  return resolveMessage(raw as any)
})

const hasLeading = computed(() => !!props.leadingIcon)
const hasTrailing = computed(() => !!props.trailingIcon)

function onInput(event: Event) {
  const target = event.target as HTMLInputElement
  emit('update:modelValue', target.value)
}
</script>

<template>
  <div :class="['w-full', props.classes.wrapper]">
    <label
      v-if="props.label"
      :for="props.label"
      :class="['mb-1 block text-sm font-medium text-gray-700', props.classes.label]"
    >
      {{ props.label }}
    </label>

    <div class="relative">
      <component
        :is="props.leadingIcon"
        v-if="hasLeading"
        :size="props.iconSize"
        weight="regular"
        class="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-gray-400"
      />

      <input
        :id="props.label"
        :type="props.type"
        :value="props.modelValue"
        :placeholder="props.placeholder"
        :class="[
          'w-full rounded-md border px-3 py-2 outline-none transition',
          'border-gray-300 focus:border-blue-500 focus:ring-1 focus:ring-blue-500',
          hasError && 'border-red-500 focus:border-red-500 focus:ring-red-500',
          hasLeading && 'pl-10',
          hasTrailing && 'pr-10',
          props.classes.input,
        ]"
        @input="onInput"
      />

      <component
        :is="props.trailingIcon"
        v-if="hasTrailing"
        :size="props.iconSize"
        weight="regular"
        class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400"
      />
    </div>

    <p v-if="hasError && errorMessage" :class="['mt-1 text-sm text-red-500', props.classes.error]">
      {{ errorMessage }}
    </p>
  </div>
</template>
