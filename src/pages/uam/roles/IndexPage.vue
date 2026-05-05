<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoleStore } from '@/stores/role'
import type { IRole } from '@/stores/role'
import { PhPlus, PhPencil, PhTrash } from '@phosphor-icons/vue'
import swal from '@/plugins/swal'
import UiCard from '@/components/utils/UiCard.vue'
import UiButton from '@/components/utils/UiButton.vue'
import UiPagination from '@/components/utils/UiPagination.vue'
import FormModal from './FormModal.vue'

const roleStore = useRoleStore()
const formModalRef = ref<InstanceType<typeof FormModal> | null>(null)

function openCreate() {
  formModalRef.value?.show()
}

function openEdit(role: IRole) {
  formModalRef.value?.show(role)
}

async function handleDelete(id: number) {
  const result = await swal.warning(
    'Hapus Role',
    'Apakah Anda yakin ingin menghapus role ini? Tindakan ini tidak dapat dibatalkan.',
  )

  if (result.isConfirmed) {
    await roleStore.deleteRole(id)
  }
}

function handlePageChange(page: number) {
  roleStore.fetchRoles(page)
}

onMounted(() => {
  roleStore.fetchRoles()
})
</script>

<template>
  <div class="max-w-7xl mx-auto px-4">
    <!-- Header Section -->
    <div class="mb-6 flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-900">Roles</h1>
        <p class="text-sm text-gray-600 mt-1">Kelola daftar role dalam sistem.</p>
      </div>
      <UiButton size="sm" @click="openCreate">
        <template #icon>
          <PhPlus class="w-4 h-4" />
        </template>
        Tambah Role
      </UiButton>
    </div>

    <!-- Loading State -->
    <div
      v-if="roleStore.loading.Index"
      class="grid gap-6 md:grid-cols-2 lg:grid-cols-3 animate-pulse"
    >
      <div
        v-for="i in 6"
        :key="i"
        class="h-48 bg-linear-to-br from-gray-200 to-gray-300 rounded-xl"
      ></div>
    </div>

    <!-- Empty State -->
    <div
      v-else-if="roleStore.indexData.roles.length === 0"
      class="text-center py-24 bg-white rounded-2xl border-2 border-dashed border-gray-300 shadow-sm"
    >
      <div
        class="inline-flex items-center justify-center w-20 h-20 bg-indigo-100 text-indigo-600 rounded-full mb-4"
      >
        <PhPlus class="w-10 h-10" />
      </div>
      <h3 class="text-2xl font-semibold text-gray-900 mb-2">Belum ada Role</h3>
      <p class="text-gray-600 mb-6">
        Silakan buat role baru untuk mulai mengatur hak akses sistem.
      </p>
      <UiButton size="lg" @click="openCreate">
        <template #icon>
          <PhPlus class="w-5 h-5" />
        </template>
        Buat Role Pertama
      </UiButton>
    </div>

    <!-- Data List -->
    <template v-else>
      <div class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        <UiCard
          v-for="role in roleStore.indexData.roles"
          :key="role.id"
          class="group hover:shadow-md transition-shadow"
        >
          <div class="flex items-start justify-between">
            <div class="min-w-0 flex-1">
              <h3 class="text-lg font-semibold text-gray-900 mb-1">
                {{ role.name }}
              </h3>
              <p class="text-sm text-gray-600 line-clamp-2">
                {{ role.description }}
              </p>
            </div>
            <div
              class="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity duration-300 shrink-0 ml-2"
            >
              <button
                class="p-1.5 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-md transition-colors"
                title="Edit"
                @click="openEdit(role)"
              >
                <PhPencil class="w-5 h-5" />
              </button>
              <button
                class="p-1.5 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-md transition-colors"
                title="Hapus"
                :disabled="roleStore.loading.Delete"
                @click="handleDelete(role.id)"
              >
                <PhTrash class="w-5 h-5" />
              </button>
            </div>
          </div>

          <!-- Permissions badges -->
          <div v-if="role.permissions?.length" class="mt-3 flex flex-wrap gap-1.5">
            <span
              v-for="perm in role.permissions.slice(0, 5)"
              :key="perm.id"
              class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-blue-50 text-blue-700"
            >
              {{ perm.name }}
            </span>
            <span
              v-if="role.permissions.length > 5"
              class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-500"
            >
              +{{ role.permissions.length - 5 }} lainnya
            </span>
          </div>
        </UiCard>
      </div>

      <!-- Pagination -->
      <div class="mt-8 flex justify-center">
        <UiPagination
          :page="roleStore.indexData.pagination.page"
          :total-pages="roleStore.indexData.pagination.totalPages"
          @update:page="handlePageChange"
        />
      </div>
    </template>
  </div>

  <FormModal ref="formModalRef" />
</template>
