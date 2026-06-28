// ============================================================
// Store with File Upload (using withFile HOF)
// ============================================================
import { defineStore } from 'pinia'
import { required, email, minLength, helpers, requiredIf } from '@vuelidate/validators'
import { useCrud, withFile } from '@/composables'

export interface IUser {
  id: number
  name: string
  email: string
  avatar: string | null
  created_at: string
}

export interface IUserPayload {
  id?: number
  name: string
  email: string
  password: string
  password_confirmation: string
  avatar: File | null
}

export const useUserStore = defineStore('user', () => {
  const crud = useCrud<IUser, IUserPayload>({
    endpoint: '/users',
    entityName: 'user',
    initialForm: {
      name: '',
      email: '',
      password: '',
      password_confirmation: '',
      avatar: null,
    },
    formRules: {
      name: { required, minLength: minLength(2) },
      email: { required, email },
      password: { required, minLength: minLength(6) },
      password_confirmation: {},
    },
    pageSize: 12,
  })

  const userCrud = withFile<IUser, IUserPayload>(crud, ['avatar'])

  // Extended formRules — defined after useCrud so crud.form is accessible
  const formRules = {
    ...crud.formRules,
    password_confirmation: {
      requiredIf: requiredIf(() => !!crud.form.password),
      sameAsPassword: helpers.withMessage(
        'Password tidak cocok',
        (value: string) => value === crud.form.password,
      ),
    },
  }

  async function create() {
    try {
      await userCrud.createForm(['id'])
    } catch (error: any) {
      console.error('Failed to create user', error)
      throw error
    }
  }

  async function update(id: number) {
    try {
      const excludeFields: (keyof IUserPayload)[] = ['id']
      if (!crud.form.password) {
        excludeFields.push('password', 'password_confirmation')
      }
      await userCrud.updateForm(id, excludeFields)
    } catch (error: any) {
      console.error('Failed to update user', error)
      throw error
    }
  }

  return {
    ...userCrud,
    formRules,
    create,
    update,
  }
})
