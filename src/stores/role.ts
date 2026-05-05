import { defineStore } from 'pinia'
import { reactive, ref } from 'vue'
import { get, post, put, del, IApiResponse } from '@/plugins/axios'
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
}

export interface IPagination {
  page: number
  pageSize: number
  total: number
  totalPages: number
}

interface IRolesApiResponse {
  data: IRole[]
  total: number
  page: number
  page_size: number
  total_pages: number
}

export type TLoadingKey = 'Index' | 'Form' | 'Delete'

export const useRoleStore = defineStore('role', () => {
  const indexData = ref<{ roles: IRole[]; pagination: IPagination }>({
    roles: [],
    pagination: { page: 1, pageSize: 10, total: 0, totalPages: 1 },
  })

  const loading = ref<Record<TLoadingKey, boolean>>({
    Index: false,
    Form: false,
    Delete: false,
  })

  const form = reactive<IRolePayload>({
    name: '',
    description: '',
  })

  const formRules = {
    name: { required },
    description: { required },
  }

  async function fetchRoles(page?: number) {
    loading.value.Index = true
    const currentPage = page ?? indexData.value.pagination.page
    try {
      const { data } = await get<IApiResponse<IRolesApiResponse>>('/roles', {
        params: {
          page: currentPage,
          page_size: indexData.value.pagination.pageSize,
        },
      })
      const body = data.data
      indexData.value.roles = body.data || []
      indexData.value.pagination = {
        page: body.page,
        pageSize: body.page_size,
        total: body.total,
        totalPages: body.total_pages,
      }
    } catch (error: any) {
      console.error('Failed to fetch roles', error)
      swal.error('Gagal', 'Gagal memuat daftar role.')
    } finally {
      loading.value.Index = false
    }
  }

  async function createRole() {
    loading.value.Form = true
    try {
      await post('/roles', {
        name: form.name,
        description: form.description,
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
    createRole,
    updateRole,
    deleteRole,
  }
})
