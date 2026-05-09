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
import { required } from '@vuelidate/validators'
import { IPermission } from './permission'
import swal from '@/plugins/swal'

export interface IRole {
  id: number
  name: string
  description: string
  permissions: IPermission[]
}

export interface IRolePayload {
  id?: number
  name: string
  description: string
  permissions: number[]
}

export type TLoadingKey = 'Index' | 'Form' | 'Delete'

export const useRoleStore = defineStore('role', () => {
  const indexData = ref<{ roles: IRole[]; pagination: IApiMetadata }>({
    roles: [],
    pagination: { ...ApiMetadataDefaults },
  })

  const loading = ref<Record<TLoadingKey, boolean>>({
    Index: false,
    Form: false,
    Delete: false,
  })

  const form = reactive<IRolePayload>({
    name: '',
    description: '',
    permissions: [],
  })

  const formRules = {
    name: { required },
    description: { required },
  }

  async function fetchRoles(page?: number) {
    loading.value.Index = true
    const currentPage = page ?? indexData.value.pagination.page
    try {
      const { data } = await get<IApiResponse<IRole[]>>('/roles', {
        params: {
          page: currentPage,
          page_size: indexData.value.pagination.page_size,
        },
      })
      indexData.value.roles = data.data || []
      indexData.value.pagination = data.metadata || ApiMetadataDefaults
    } catch (error: any) {
      console.error('Failed to fetch roles', error)
      swal.error('Gagal', 'Gagal memuat daftar role.')
    } finally {
      loading.value.Index = false
    }
  }

  async function fetchAllRoles(): Promise<IRole[]> {
    try {
      const { data } = await get<IApiResponse<IRole[]>>('/roles')
      return data.data || []
    } catch (error: any) {
      console.error('Failed to fetch all roles', error)
      return []
    }
  }

  async function createRole() {
    loading.value.Form = true
    try {
      await post('/roles', {
        name: form.name,
        description: form.description,
        permissions: form.permissions,
      })
      swal.success('Berhasil', 'Role berhasil dibuat.')
      await fetchRoles()
    } catch (error: any) {
      const message = error?.response?.data?.message || 'Gagal membuat role.'
      swal.error('Gagal', message)
      throw error
    } finally {
      loading.value.Form = false
    }
  }

  async function updateRole(id: number) {
    loading.value.Form = true
    try {
      await put(`/roles/${id}`, {
        name: form.name,
        description: form.description,
        permissions: form.permissions,
      })
      swal.success('Berhasil', 'Role berhasil diperbarui.')
      await fetchRoles()
    } catch (error: any) {
      const message = error?.response?.data?.message || 'Gagal memperbarui role.'
      swal.error('Gagal', message)
      throw error
    } finally {
      loading.value.Form = false
    }
  }

  async function deleteRole(id: number) {
    loading.value.Delete = true
    try {
      await del(`/roles/${id}`)
      swal.success('Berhasil', 'Role berhasil dihapus.')
      await fetchRoles()
    } catch (error: any) {
      const message = error?.response?.data?.message || 'Gagal menghapus role.'
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
    fetchRoles,
    fetchAllRoles,
    createRole,
    updateRole,
    deleteRole,
  }
})
