<script setup lang="ts">
import swal from '@/plugins/swal'
import UiCard from '@/components/utils/UiCard.vue'
import UiButton from '@/components/utils/UiButton.vue'
import UiPagination from '@/components/utils/UiPagination.vue'
import FormModal from './FormModal.vue'

import { ref, onMounted } from 'vue'
import { usePermissionStore, type IPermission } from '@/stores/permission'
import { PhPlus, PhPencil, PhTrash } from '@phosphor-icons/vue'

const permissionStore = usePermissionStore()
const formModalRef = ref<InstanceType<typeof FormModal> | null>(null)

function openCreate() {
  formModalRef.value?.show()
}

function openEdit(permission: IPermission) {
  formModalRef.value?.show(permission)
}

async function handleDelete(id: number) {
  const result = await swal.warning(
    'Hapus Permission',
    'Apakah Anda yakin ingin menghapus permission ini? Tindakan ini tidak dapat dibatalkan.',
  )

  if (result.isConfirmed) {
    await permissionStore.deletePermission(id)
  }
}

function handlePageChange(page: number) {
  permissionStore.fetchPermissions(page)
}

onMounted(() => {
  permissionStore.fetchPermissions()
})
</script>

<template>
  <div class="mx-auto px-4">
    <!-- Header Section -->
    <div class="mb-6 flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-900">Permissions</h1>
        <p class="hidden sm:block text-sm text-gray-600 mt-1">Kelola daftar akses permission dalam sistem.</p>
      </div>
      <UiButton size="sm" @click="openCreate">
        <template #icon>
          <PhPlus class="w-4 h-4" />
        </template>
        Tambah Permission
      </UiButton>
    </div>

    <!-- Loading State -->
    <div
      v-if="permissionStore.loading.Index"
      class="grid gap-6 md:grid-cols-2 lg:grid-cols-3 animate-pulse"
    >
      <div
        v-for="i in 6"
        :key="i"
        class="h-40 bg-linear-to-br from-gray-200 to-gray-300 rounded-xl"
      ></div>
    </div>

    <!-- Empty State -->
    <div
      v-else-if="permissionStore.indexData.permissions.length === 0"
      class="text-center py-24 bg-white rounded-2xl border-2 border-dashed border-gray-300 shadow-sm"
    >
      <div
        class="inline-flex items-center justify-center w-20 h-20 bg-indigo-100 text-indigo-600 rounded-full mb-4"
      >
        <PhPlus class="w-10 h-10" />
      </div>
      <h3 class="text-2xl font-semibold text-gray-900 mb-2">Belum ada Permission</h3>
      <p class="text-gray-600 mb-6">
        Silakan buat permission baru untuk mulai mengatur hak akses sistem.
      </p>
      <UiButton size="lg" @click="openCreate">
        <template #icon>
          <PhPlus class="w-5 h-5" />
        </template>
        Buat Permission Pertama
      </UiButton>
    </div>

    <!-- No Search Results -->
    <!-- Data List -->
    <div v-else class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
      <UiCard
        v-for="permission in permissionStore.indexData.permissions"
        :key="permission.id"
        class="group hover:shadow-md transition-shadow"
      >
        <div class="flex items-start justify-between">
          <div>
            <h3 class="text-lg font-semibold text-gray-900 mb-1">
              {{ permission.name }}
            </h3>
            <p class="text-sm text-gray-600 line-clamp-2">
              {{ permission.description }}
            </p>
          </div>
          <div
            class="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity duration-300"
          >
            <button
              class="p-1.5 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-md transition-colors"
              title="Edit"
              @click="openEdit(permission)"
            >
              <PhPencil class="w-5 h-5" />
            </button>
            <button
              class="p-1.5 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-md transition-colors"
              title="Hapus"
              :disabled="permissionStore.loading.Delete"
              @click="handleDelete(permission.id)"
            >
              <PhTrash class="w-5 h-5" />
            </button>
          </div>
        </div>
      </UiCard>
    </div>

    <!-- Pagination -->
    <div class="mt-8 flex justify-center">
      <UiPagination
        :page="permissionStore.indexData.pagination.page"
        :total-pages="permissionStore.indexData.pagination.total_pages"
        @update:page="handlePageChange"
      />
    </div>
  </div>

  <FormModal ref="formModalRef" />
</template>
