<script setup lang="ts">
import { ref, computed } from 'vue'
import { PhUploadSimple, PhX } from '@phosphor-icons/vue'

interface Props {
  modelValue: File | null
  currentAvatar?: string | null
  label?: string
  accept?: string
  maxSize?: number
  disabled?: boolean
  labelClass?: string
  wrapperClass?: string
}

const props = withDefaults(defineProps<Props>(), {
  currentAvatar: null,
  label: 'Avatar',
  accept: 'image/jpeg,image/png,image/gif,image/webp',
  maxSize: 5,
  disabled: false,
  labelClass: '',
  wrapperClass: '',
})

const emit = defineEmits<{
  'update:modelValue': [value: File | null]
  'remove': []
}>()

const fileInput = ref<HTMLInputElement | null>(null)
const previewUrl = ref<string | null>(null)
const errorMessage = ref('')

const hasNewFile = computed(() => props.modelValue !== null)
const previewImage = computed(() => {
  if (previewUrl.value) return previewUrl.value
  if (props.currentAvatar) return props.currentAvatar
  return null
})

function handleClick() {
  if (props.disabled) return
  fileInput.value?.click()
}

function handleFileSelect(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return

  if (props.maxSize > 0 && file.size > props.maxSize * 1024 * 1024) {
    errorMessage.value = `File terlalu besar (max ${props.maxSize} MB)`
    return
  }

  errorMessage.value = ''
  if (previewUrl.value) {
    URL.revokeObjectURL(previewUrl.value)
  }
  previewUrl.value = URL.createObjectURL(file)
  emit('update:modelValue', file)
  input.value = ''
}

function handleRemove() {
  if (previewUrl.value) {
    URL.revokeObjectURL(previewUrl.value)
    previewUrl.value = null
  }
  emit('update:modelValue', null)
  emit('remove')
}

function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}
</script>

<template>
  <div :class="['w-full', wrapperClass]">
    <label v-if="label" :class="['mb-2 block text-sm font-medium text-gray-700', labelClass]">
      {{ label }}
    </label>

    <div class="flex items-start gap-4">
      <!-- Avatar Preview -->
      <div class="relative shrink-0">
        <div
          v-if="previewImage"
          class="w-20 h-20 rounded-full overflow-hidden border-2 border-gray-200"
        >
          <img :src="previewImage" class="w-full h-full object-cover" />
        </div>
        <div
          v-else
          class="w-20 h-20 rounded-full bg-gray-100 border-2 border-dashed border-gray-300 flex items-center justify-center"
        >
          <PhUploadSimple :size="24" class="text-gray-400" />
        </div>

        <!-- Remove button -->
        <button
          v-if="previewImage"
          type="button"
          class="absolute -top-1 -right-1 p-1 bg-red-500 text-white rounded-full hover:bg-red-600 transition-colors shadow-sm"
          @click="handleRemove"
        >
          <PhX :size="12" />
        </button>
      </div>

      <!-- Upload Area -->
      <div class="flex-1">
        <input
          ref="fileInput"
          type="file"
          :accept="accept"
          :disabled="disabled"
          class="hidden"
          @change="handleFileSelect"
        />

        <button
          type="button"
          :disabled="disabled"
          :class="[
            'inline-flex items-center gap-2 px-4 py-2 text-sm font-medium rounded-md border transition-colors',
            disabled
              ? 'border-gray-200 bg-gray-50 text-gray-400 cursor-not-allowed'
              : 'border-gray-300 bg-white text-gray-700 hover:bg-gray-50 hover:border-gray-400',
          ]"
          @click="handleClick"
        >
          <PhUploadSimple :size="16" />
          {{ hasNewFile ? 'Ganti Avatar' : 'Upload Avatar' }}
        </button>

        <p class="mt-1.5 text-xs text-gray-500">
          JPG, PNG, GIF, WEBP. Max {{ maxSize }}MB
        </p>

        <!-- File info -->
        <div v-if="hasNewFile && modelValue" class="mt-2 flex items-center gap-2 text-xs text-gray-600">
          <span class="truncate max-w-[200px]">{{ modelValue.name }}</span>
          <span class="text-gray-400">•</span>
          <span>{{ formatFileSize(modelValue.size) }}</span>
        </div>

        <!-- Error message -->
        <p v-if="errorMessage" class="mt-1.5 text-xs text-red-500">
          {{ errorMessage }}
        </p>
      </div>
    </div>
  </div>
</template>
