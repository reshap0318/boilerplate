import axios, { type AxiosInstance, type AxiosRequestConfig, type AxiosResponse } from 'axios'
import storage from '@/helpers/storage'
import { useAuthStore } from '@/stores/auth'

export interface IApiResponse<TData> {
  code: number
  message: string
  data: TData
}

const api: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor
api.interceptors.request.use(
  (config) => {
    const authStore = useAuthStore()
    const token = authStore.token
    if (token && !config.headers.Authorization) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error),
)

// Response interceptor
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      storage.clearAll()
      window.location.href = '/login'
    }
    return Promise.reject(error)
  },
)

// Helper methods
const get = <T = unknown>(url: string, config?: AxiosRequestConfig): Promise<AxiosResponse<T>> => {
  return api.get<T>(url, config)
}

const post = <T = unknown>(
  url: string,
  data?: unknown,
  config?: AxiosRequestConfig,
): Promise<AxiosResponse<T>> => {
  return api.post<T>(url, data, config)
}

const put = <T = unknown>(
  url: string,
  data?: unknown,
  config?: AxiosRequestConfig,
): Promise<AxiosResponse<T>> => {
  return api.put<T>(url, data, config)
}

const patch = <T = unknown>(
  url: string,
  data?: unknown,
  config?: AxiosRequestConfig,
): Promise<AxiosResponse<T>> => {
  return api.patch<T>(url, data, config)
}

const del = <T = unknown>(url: string, config?: AxiosRequestConfig): Promise<AxiosResponse<T>> => {
  return api.delete<T>(url, config)
}

export { get, post, put, patch, del }
export default api
