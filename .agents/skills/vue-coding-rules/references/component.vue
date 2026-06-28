// ============================================================
// Component Structure — <script setup lang="ts">
// ============================================================
<script setup lang="ts">
// 1. Component imports (FIRST)
import UiCard from '@/components/utils/UiCard.vue'
import UiButton from '@/components/utils/UiButton.vue'
import FormModal from './FormModal.vue'

// 2. Vue core imports
import { ref, computed, onMounted } from 'vue'

// 3. Other imports (stores, composables, helpers, plugins, etc)
import { useUserStore } from '@/stores/user'
import swal from '@/plugins/swal'

// 4. Props & Emits
const props = defineProps<{
  userId: number
}>()

const emit = defineEmits<{
  update: [id: number]
  delete: [id: number]
}>()

// 5. Store/Composable usage
const userStore = useUserStore()

// 6. Reactive state
const isLoading = ref(false)
const localData = ref<string>('')

// 7. Computed properties
const formattedName = computed(() => userStore.user?.name ?? '')

// 8. Methods/functions
function handleAction() {
  // implementation
}

async function handleAsyncAction() {
  isLoading.value = true
  try {
    // async operation
  } finally {
    isLoading.value = false
  }
}

// 9. Lifecycle hooks
onMounted(() => {
  // initialization
})

// 10. Expose (if needed)
defineExpose({ handleAction })
</script>

<template>
  <!-- Template content -->
</template>

<style scoped>
/* Styles - always use scoped unless global styling is intentional */
</style>
