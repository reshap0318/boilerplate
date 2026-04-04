import { defineStore } from 'pinia'
import { reactive, ref } from 'vue'
import { post } from '@/plugins/axios'
import storage from '@/helpers/storage'
import type { IApiResponse } from '@/plugins/axios'
import { required, email } from '@vuelidate/validators'

export interface ILoginPayload {
  email: string
  password: string
}

export interface ILoginResponse {
  token: string
  user: {
    id: number
    name: string
    email: string
  }
}

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(storage.getItem<string>('token'))
  const user = ref<ILoginResponse['user'] | null>(storage.getItem<ILoginResponse['user']>('user'))
  const isLoading = ref(false)

  const form = reactive<ILoginPayload>({
    email: '',
    password: '',
  })

  const formRules = {
    email: {
      required,
      email,
    },
    password: {
      required,
    },
  }

  async function login(): Promise<void> {
    isLoading.value = true
    try {
      const response = await post<IApiResponse<ILoginResponse>>('/auth/login', {
        email: form.email,
        password: form.password,
      })
      const { token: newToken, user: userData } = response.data.data

      token.value = newToken
      user.value = userData

      storage.setItem('token', newToken)
      storage.setItem('user', userData)
    } finally {
      isLoading.value = false
    }
  }

  function logout() {
    token.value = null
    user.value = null
    storage.removeItem('token')
    storage.removeItem('user')
  }

  const isAuthenticated = (): boolean => !!token.value

  return { 
    token, 
    user, 
    isLoading, 
    form, 
    formRules, 
    login, 
    logout, 
    isAuthenticated
  }
})
