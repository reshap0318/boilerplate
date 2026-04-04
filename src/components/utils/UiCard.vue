<script setup lang="ts">
interface Props {
  headerClass?: string
  bodyClass?: string
  footerClass?: string
  cardClass?: string
  wrapperClass?: string
  padded?: boolean
}

withDefaults(defineProps<Props>(), {
  headerClass: '',
  bodyClass: '',
  footerClass: '',
  cardClass: '',
  wrapperClass: '',
  padded: true,
})

defineSlots<{
  header?: () => any
  default?: () => any
  footer?: () => any
}>()
</script>

<template>
  <div :class="['w-full', wrapperClass]">
    <div :class="['overflow-hidden rounded-lg bg-white shadow-md', cardClass]">
      <!-- Header -->
      <div v-if="$slots.header" :class="[headerClass]">
        <slot name="header" />
      </div>

      <!-- Body -->
      <div :class="[padded && 'p-6', bodyClass]">
        <slot />
      </div>

      <!-- Footer -->
      <div v-if="$slots.footer" :class="['bg-gray-50 px-6 py-4', footerClass]">
        <slot name="footer" />
      </div>
    </div>
  </div>
</template>
