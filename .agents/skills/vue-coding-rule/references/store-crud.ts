// ============================================================
// Standard CRUD Store (no file upload)
// ============================================================
import { defineStore } from 'pinia'
import { required, minLength } from '@vuelidate/validators'
import { get, type IApiResponse } from '@/plugins/axios'
import { useCrud } from '@/composables'

export interface IEntity {
  id: number
  name: string
  description: string
}

export interface IEntityPayload {
  id?: number
  name: string
  description: string
}

export const useEntityStore = defineStore('entity', () => {
  const crud = useCrud<IEntity, IEntityPayload>({
    endpoint: '/entities',
    entityName: 'entity',
    initialForm: { name: '', description: '' },
    formRules: {
      name: { required, minLength: minLength(3) },
      description: {},
    },
  })

  // Add custom methods if needed
  async function fetchAllEntities(): Promise<IEntity[]> {
    try {
      const { data } = await get<IApiResponse<IEntity[]>>('/entities')
      return data.data || []
    } catch (error: any) {
      console.error('Failed to fetch all entities', error)
      return []
    }
  }

  return {
    ...crud,
    fetchAllEntities,
  }
})
