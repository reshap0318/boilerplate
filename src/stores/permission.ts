import { defineStore } from 'pinia'
import { reactive, ref } from 'vue'
import { 
  get,
  post,
  put,
  del,
  ApiMetadataDefaults,
  type IApiResponse,
  type IApiMetadata 
} from '@/plugins/axios'
import { required, minLength } from '@vuelidate/validators'
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
  const indexData = ref<{ permissions: IPermission[]; pagination: IApiMetadata }>({
    permissions: [],
    pagination: { ...ApiMetadataDefaults },
  })
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
    name: { required, minLength: minLength(3) },
    description: {},
  }

  async function fetchPermissions(page?: number) {
    loading.value.Index = true
    const currentPage = page ?? indexData.value.pagination.page
    try {
      const { data } = await get<IApiResponse<IPermission[]>>('/permissions', {
        params: {
          page: currentPage,
          page_size: indexData.value.pagination.page_size,
        },
      })
      indexData.value.permissions = data.data || []
      indexData.value.pagination = data.metadata || ApiMetadataDefaults
    } catch (error: any) {
      console.error('Failed to fetch permissions', error)
      swal.error('Gagal', 'Gagal memuat daftar permission.')
    } finally {
      loading.value.Index = false
    }
  }

  async function fetchPermissionById(id: number): Promise<IPermission | null> {
    try {
      const { data } = await get<IApiResponse<IPermission>>(`/permissions/${id}`)
      return data.data || null
    } catch (error: any) {
      console.error('Failed to fetch permission', error)
      return null
    }
  }

  async function fetchAllPermissions(): Promise<IPermission[]> {
    try {
      const { data } = await get<IApiResponse<IPermission[]>>('/permissions')
      return data.data || []
    } catch (error: any) {
      console.error('Failed to fetch all permissions', error)
      return []
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
      throw error
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
    indexData,
    loading,
    form,
    formRules,
    fetchPermissions,
    fetchPermissionById,
    fetchAllPermissions,
    createPermission,
    updatePermission,
    deletePermission,
  }
})
