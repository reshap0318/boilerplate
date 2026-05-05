<script setup lang="ts">
import { computed, ref } from 'vue'
import useVuelidate from '@vuelidate/core'
import { useRoleStore } from '@/stores/role'
import UiModal from '@/components/utils/UiModal.vue'
import FormInput from '@/components/utils/FormInput.vue'
import UiButton from '@/components/utils/UiButton.vue'

const roleStore = useRoleStore()
const v$ = useVuelidate(roleStore.formRules, roleStore.form)

const isVisible = ref(false)
const isEdit = computed(() => !!roleStore.form.id)

function show(data?: { id?: number; name: string; description: string }) {
  if (data) {
    roleStore.form.id = data.id
    roleStore.form.name = data.name
    roleStore.form.description = data.description
  } else {
    roleStore.form.id = undefined
    roleStore.form.name = ''
    roleStore.form.description = ''
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
    if (isEdit.value && roleStore.form.id) {
      await roleStore.updateRole(roleStore.form.id)
    } else {
      await roleStore.createRole()
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
    :title="isEdit ? 'Edit Role' : 'Tambah Role'"
    size="md"
    @close="close"
  >
    <form @submit.prevent="handleSubmit">
      <div class="space-y-4">
        <FormInput
          v-model="roleStore.form.name"
          label="Nama Role"
          placeholder="e.g. Admin"
          :validation="v$.name"
        />

        <FormInput
          v-model="roleStore.form.description"
          label="Deskripsi"
          placeholder="Administrator role"
          :validation="v$.description"
        />
      </div>

      <!-- Actions -->
      <div class="mt-6 flex justify-end gap-2">
        <UiButton
          type="button"
          variant="secondary"
          :disabled="roleStore.loading.Form"
          outline
          @click="close"
        >
          Batal
        </UiButton>
        <UiButton type="submit" :loading="roleStore.loading.Form">
          {{ isEdit ? 'Perbarui' : 'Simpan' }}
        </UiButton>
      </div>
    </form>
  </UiModal>
</template>
