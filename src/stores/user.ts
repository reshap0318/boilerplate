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
import { required, email, minLength, sameAs } from '@vuelidate/validators'
import { IRole } from './role'
import { IPermission } from './permission'
import { uploadFile } from '@/helpers/upload'
import swal from '@/plugins/swal'

export interface IUser {
  id: number
  email: string
  name: string
  avatar: string | null
  created_at: string
  roles: IRole[]
  permissions: IPermission[]
}

export interface IUserPayload {
  id?: number
  name: string
  email: string
  password: string
  password_confirmation: string
  roles: number[]
  avatar: File | null
}

type TLoadingKey = 'Index' | 'Form' | 'Delete'

export const useUserStore = defineStore('user', () => {
  const indexData = ref<{ users: IUser[]; pagination: IApiMetadata }>({
    users: [],
    pagination: { ...ApiMetadataDefaults },
  })

  const loading = ref<Record<TLoadingKey, boolean>>({
    Index: false,
    Form: false,
    Delete: false,
  })

  const form = reactive<IUserPayload>({
    name: '',
    email: '',
    password: '',
    password_confirmation: '',
    roles: [],
    avatar: null,
  })

  const formRules = {
    name: { required, minLength: minLength(2) },
    email: { required, email },
    password: { required, minLength: minLength(6) },
    password_confirmation: { required, sameAsPassword: sameAs('password') },
  }

  async function fetchUsers(page?: number) {
    loading.value.Index = true
    const currentPage = page ?? indexData.value.pagination.page
    try {
      const { data } = await get<IApiResponse<IUser[]>>('/users', {
        params: {
          page: currentPage,
          page_size: 12,
        },
      })
      indexData.value.users = data.data || []
      indexData.value.pagination = data.metadata || ApiMetadataDefaults
    } catch (error: any) {
      console.error('Failed to fetch users', error)
      swal.error('Gagal', 'Gagal memuat daftar user.')
    } finally {
      loading.value.Index = false
    }
  }

  async function fetchUserById(id: number): Promise<IUser | null> {
    try {
      const { data } = await get<IApiResponse<IUser>>(`/users/${id}`)
      return data.data || null
    } catch (error: any) {
      console.error('Failed to fetch user', error)
      return null
    }
  }

  async function createUser() {
    loading.value.Form = true
    try {
      let avatarUuid: string | null = null
      if (form.avatar) {
        const uploaded = await uploadFile(form.avatar)
        avatarUuid = uploaded.uuid
      }

      const payload: any = {
        name: form.name,
        email: form.email,
        password: form.password,
        password_confirmation: form.password_confirmation,
        roles: form.roles,
      }
      if (avatarUuid) {
        payload.avatar = avatarUuid
      }

      await post('/users', payload)
      swal.success('Berhasil', 'User berhasil dibuat.')
      await fetchUsers()
    } catch (error: any) {
      const message = error?.response?.data?.message || 'Gagal membuat user.'
      swal.error('Gagal', message)
      throw error
    } finally {
      loading.value.Form = false
    }
  }

  async function updateUser(id: number) {
    loading.value.Form = true
    try {
      let avatarUuid: string | null = null
      if (form.avatar) {
        const uploaded = await uploadFile(form.avatar)
        avatarUuid = uploaded.uuid
      }

      const payload: any = {
        name: form.name,
        email: form.email,
        roles: form.roles,
      }
      if (form.password) {
        payload.password = form.password
        payload.password_confirmation = form.password_confirmation
      }
      if (avatarUuid) {
        payload.avatar = avatarUuid
      }

      await put(`/users/${id}`, payload)
      swal.success('Berhasil', 'User berhasil diperbarui.')
      await fetchUsers()
    } catch (error: any) {
      const message = error?.response?.data?.message || 'Gagal memperbarui user.'
      swal.error('Gagal', message)
      throw error
    } finally {
      loading.value.Form = false
    }
  }

  async function deleteUser(id: number) {
    loading.value.Delete = true
    try {
      await del(`/users/${id}`)
      swal.success('Berhasil', 'User berhasil dihapus.')
      await fetchUsers()
    } catch (error: any) {
      const message = error?.response?.data?.message || 'Gagal menghapus user.'
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
    fetchUsers,
    fetchUserById,
    createUser,
    updateUser,
    deleteUser,
  }
})
