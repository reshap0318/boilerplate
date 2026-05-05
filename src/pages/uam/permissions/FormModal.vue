<script setup lang="ts">
import { computed, ref } from 'vue'
import useVuelidate from '@vuelidate/core'
import { usePermissionStore } from '@/stores/permission'
import UiModal from '@/components/utils/UiModal.vue'
import FormInput from '@/components/utils/FormInput.vue'
import UiButton from '@/components/utils/UiButton.vue'

const permissionStore = usePermissionStore()
const v$ = useVuelidate(permissionStore.formRules, permissionStore.form)

const isVisible = ref(false)
const isEdit = computed(() => !!permissionStore.form.id)

function show(data?: { id?: number; name: string; description: string }) {
  if (data) {
    permissionStore.form.id = data.id
    permissionStore.form.name = data.name
    permissionStore.form.description = data.description
  } else {
    permissionStore.form.id = undefined
    permissionStore.form.name = ''
    permissionStore.form.description = ''
  }
  v$.value.$reset()
  isVisible.value = true
}

function close() {
  isVisible.value = false
}

async function handleSubmit() {
  const isValid = await v$.value.$validate()
  if (!isValid) return

  try {
    if (isEdit.value && permissionStore.form.id) {
      await permissionStore.updatePermission(permissionStore.form.id)
    } else {
      await permissionStore.createPermission()
    }
  } finally {
    close()
  }
}

defineExpose({ show, close })
</script>

<template>
  <UiModal
    v-model="isVisible"
    :title="isEdit ? 'Edit Permission' : 'Tambah Permission'"
    size="md"
    @close="close"
  >
    <form @submit.prevent="handleSubmit">
      <div class="space-y-4">
        <FormInput
          v-model="permissionStore.form.name"
          label="Nama Permission"
          placeholder="e.g. users.index"
          :validation="v$.name"
        />

        <FormInput
          v-model="permissionStore.form.description"
          label="Deskripsi"
          placeholder="Melihat daftar pengguna"
          :validation="v$.description"
        />
      </div>

      <!-- Actions -->
      <div class="mt-6 flex justify-end gap-2">
        <UiButton
          type="button"
          variant="secondary"
          :disabled="permissionStore.loading.Form"
          outline
          @click="close"
        >
          Batal
        </UiButton>
        <UiButton type="submit" :loading="permissionStore.loading.Form">
          {{ isEdit ? 'Perbarui' : 'Simpan' }}
        </UiButton>
      </div>
    </form>
  </UiModal>
</template>
