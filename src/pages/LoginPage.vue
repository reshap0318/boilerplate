<script setup lang="ts">
import useVuelidate from '@vuelidate/core'
import swal from '@/plugins/swal'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import { PhEnvelope, PhLock } from '@phosphor-icons/vue'
import FormInput from '@/components/utils/FormInput.vue'
import FormPassword from '@/components/utils/FormPassword.vue'
import UiButton from '@/components/utils/UiButton.vue'
import UiCard from '@/components/utils/UiCard.vue'

const router = useRouter()
const authStore = useAuthStore()

const v$ = useVuelidate(authStore.formRules, authStore.form)

async function handleLogin() {
  const isValid = await v$.value.$validate()
  if (!isValid) return

  try {
    await authStore.login()
    swal.success('Login Berhasil')
    router.push('/')
  } catch (error: any) {
    const message = error?.response?.data?.message || 'Login gagal, periksa kembali email dan password'
    swal.error('Login Gagal', message)
  }
}
</script>

<template>
  <div class="flex min-h-screen items-center justify-center bg-gray-100">
    <UiCard :padded="false" card-class="p-8" class="max-w-md">
      <template #header>
        <h2 class="mb-6 text-center text-2xl font-bold text-gray-800">Login</h2>
      </template>
      <form @submit.prevent="handleLogin">
        <FormInput
          v-model="authStore.form.email"
          label="Email"
          type="email"
          class="mb-4"
          placeholder="admin@example.com"
          :validation="v$.email"
          :leading-icon="PhEnvelope"
        />

        <FormPassword
          v-model="authStore.form.password"
          label="Password"
          placeholder="••••••••"
          class="mb-6"
          :validation="v$.password"
          :leading-icon="PhLock"
        />

        <UiButton
          type="submit"
          full-width
          :loading="authStore.isLoading"
        >
          Login
        </UiButton>
      </form>
    </UiCard>
  </div>
</template>
