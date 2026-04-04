<script setup lang="ts">
import { computed, type Component } from 'vue'
import { getErrorMessage } from '@/helpers/vuelidate'

interface ValidationLike {
  $error: boolean
  $errors: Array<{ $message: string | { value: string } }>
}

const resolveMessage = (msg: string | { value: string }): string => {
  return typeof msg === 'string' ? msg : msg.value
}

interface Props {
  modelValue: string
  label?: string
  type?: string
  placeholder?: string
  validation?: ValidationLike
  leadingIcon?: Component
  trailingIcon?: Component
  iconSize?: number
  inputClass?: string
  labelClass?: string
  errorClass?: string
  wrapperClass?: string
}

const props = withDefaults(defineProps<Props>(), {
  type: 'text',
  label: '',
  placeholder: '',
  validation: undefined,
  leadingIcon: undefined,
  trailingIcon: undefined,
  iconSize: 20,
  inputClass: '',
  labelClass: '',
  errorClass: '',
  wrapperClass: '',
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const hasError = computed(() => props.validation?.$error ?? false)
const errorMessage = computed(() => {
  if (!props.validation) return ''
  const raw = getErrorMessage(props.validation as any)
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
  <div :class="['w-full', wrapperClass]">
    <label
      v-if="label"
      :for="label"
      :class="['mb-1 block text-sm font-medium text-gray-700', labelClass]"
    >
      {{ label }}
    </label>

    <div class="relative">
      <component
        :is="leadingIcon"
        v-if="hasLeading"
        :size="iconSize"
        weight="regular"
        class="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-gray-400"
      />

      <input
        :id="label"
        :type="type"
        :value="modelValue"
        :placeholder="placeholder"
        :class="[
          'w-full rounded-md border px-3 py-2 outline-none transition',
          'border-gray-300 focus:border-blue-500 focus:ring-1 focus:ring-blue-500',
          hasError && 'border-red-500 focus:border-red-500 focus:ring-red-500',
          hasLeading && 'pl-10',
          hasTrailing && 'pr-10',
          inputClass,
        ]"
        @input="onInput"
      />

      <component
        :is="trailingIcon"
        v-if="hasTrailing"
        :size="iconSize"
        weight="regular"
        class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400"
      />
    </div>

    <p v-if="hasError && errorMessage" :class="['mt-1 text-sm text-red-500', errorClass]">
      {{ errorMessage }}
    </p>
  </div>
</template>
