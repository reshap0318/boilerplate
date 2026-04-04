<script setup lang="ts">
import { ref, computed, type Component } from 'vue'
import { PhEye, PhEyeSlash } from '@phosphor-icons/vue'
import { getErrorMessage } from '@/helpers/vuelidate'

interface ValidationLike {
  $error: boolean
  $errors: Array<{ $message: string | { value: string } }>
}

interface Props {
  modelValue: string
  label?: string
  placeholder?: string
  validation?: ValidationLike
  leadingIcon?: Component
  iconSize?: number
  inputClass?: string
  labelClass?: string
  errorClass?: string
  wrapperClass?: string
}

const props = withDefaults(defineProps<Props>(), {
  label: '',
  placeholder: '',
  validation: undefined,
  leadingIcon: undefined,
  iconSize: 20,
  inputClass: '',
  labelClass: '',
  errorClass: '',
  wrapperClass: '',
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const showPassword = ref(false)

const inputType = computed(() => (showPassword.value ? 'text' : 'password'))

const hasError = computed(() => props.validation?.$error ?? false)
const errorMessage = computed(() => {
  if (!props.validation) return ''
  const raw = getErrorMessage(props.validation as any)
  return typeof raw === 'string' ? raw : raw.value
})

const hasLeading = computed(() => !!props.leadingIcon)

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
        :type="inputType"
        :value="modelValue"
        :placeholder="placeholder"
        :class="[
          'w-full rounded-md border px-3 py-2 outline-none transition',
          'border-gray-300 focus:border-blue-500 focus:ring-1 focus:ring-blue-500',
          hasError && 'border-red-500 focus:border-red-500 focus:ring-red-500',
          hasLeading && 'pl-10',
          'pr-10',
          inputClass,
        ]"
        @input="onInput"
      />

      <button
        type="button"
        class="absolute inset-y-0 right-0 flex items-center pr-3 text-gray-500 hover:text-gray-700"
        @click="showPassword = !showPassword"
      >
        <PhEye v-if="!showPassword" :size="20" weight="regular" />
        <PhEyeSlash v-else :size="20" weight="regular" />
      </button>
    </div>

    <p v-if="hasError && errorMessage" :class="['mt-1 text-sm text-red-500', errorClass]">
      {{ errorMessage }}
    </p>
  </div>
</template>
