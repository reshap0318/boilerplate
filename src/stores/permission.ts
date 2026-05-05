import { defineStore } from 'pinia'
import { reactive, ref } from 'vue'
import { get, post, put, del } from '@/plugins/axios'
import type { IApiResponse } from '@/plugins/axios'
import { required } from '@vuelidate/validators'
import swal from '@/plugins/swal'

export interface IPermission {
  id: number
  name: string
  description: string
}

export interface IPermissionPayload {
  id?: number
  name: string
  description: string
}

export type TLoadingKey = 'Index' | 'Form' | 'Delete'

export const usePermissionStore = defineStore('permission', () => {
  const permissions = ref<IPermission[]>([])
  const loading = ref<Record<TLoadingKey, boolean>>({
    Index: false,
    Form: false,
    Delete: false,
  })

  const form = reactive<IPermissionPayload>({
    name: '',
    description: '',
  })

  const formRules = {
    name: { required },
    description: { required },
  }

  async function fetchPermissions() {
    loading.value.Index = true
    try {
      const response = await get<IApiResponse<IPermission[]>>('/permissions')
      permissions.value = response.data.data || []
    } catch (error: any) {
      console.error('Failed to fetch permissions', error)
      swal.error('Gagal', 'Gagal memuat daftar permission.')
    } finally {
      loading.value.Index = false
    }
  }

  async function createPermission() {
    loading.value.Form = true
    try {
      await post('/permissions', {
        name: form.name,
        description: form.description,
      })
      swal.success('Berhasil', 'Permission berhasil dibuat.')
      await fetchPermissions()
    } catch (error: any) {
      const message = error?.response?.data?.message || 'Gagal membuat permission.'
      swal.error('Gagal', message)
      throw error // Re-throw to let component know if needed
    } finally {
      loading.value.Form = false
    }
  }

  async function updatePermission(id: number) {
    loading.value.Form = true
    try {
      await put(`/permissions/${id}`, {
        name: form.name,
        description: form.description,
      })
      swal.success('Berhasil', 'Permission berhasil diperbarui.')
      await fetchPermissions()
    } catch (error: any) {
      const message = error?.response?.data?.message || 'Gagal memperbarui permission.'
      swal.error('Gagal', message)
      throw error
    } finally {
      loading.value.Form = false
    }
  }

  async function deletePermission(id: number) {
    loading.value.Delete = true
    try {
      await del(`/permissions/${id}`)
      swal.success('Berhasil', 'Permission berhasil dihapus.')
      await fetchPermissions()
    } catch (error: any) {
      const message = error?.response?.data?.message || 'Gagal menghapus permission.'
      swal.error('Gagal', message)
      throw error
    } finally {
      loading.value.Delete = false
    }
  }

  return {
    permissions,
    loading,
    form,
    formRules,
    fetchPermissions,
    createPermission,
    updatePermission,
    deletePermission,
  }
})
