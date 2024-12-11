<template>
  <div class="flex min-h-screen items-center justify-center px-4 py-8">
    <UCard class="w-full max-w-md">
      <template #header>
        <div class="flex justify-center">
          <NuxtImg
              src="https://custom-images.strikinglycdn.com/res/hrscywv4p/image/upload/c_limit,fl_lossy,h_1440,w_720,f_auto,q_auto/778352/q7cuvzx29gdhipnlp8xt.jpg"
              alt="Venturas Logo"
              class="h-10 w-auto"
          />
        </div>
        <h2 class="text-center text-2xl font-bold mt-4 text-gray-900 dark:text-white">
          Sign in to your account
        </h2>
      </template>

      <!-- Show login error if exists -->
      <UAlert
          v-if="authStore.loginError"
          color="red"
          :title="authStore.loginError"
          class="mb-4"
      />

      <form @submit.prevent="handleLogin" class="space-y-6">
        <UFormGroup label="Username" name="username">
          <UInput
              v-model="email"
              type="text"
              placeholder="Enter your username"
              icon="i-heroicons-envelope"
              required
          />
        </UFormGroup>

        <UFormGroup label="Password" name="password">
<!--          <div class="flex justify-between items-center mb-1">-->
<!--            <label class="block text-sm font-medium">Password</label>-->
<!--            <ULink-->
<!--                to="#"-->
<!--                disabled-->
<!--                class="text-sm text-primary hover:text-primary-400"-->
<!--            >-->
<!--              Forgot password?-->
<!--            </ULink>-->
<!--          </div>-->
          <UInput
              v-model="password"
              type="password"
              placeholder="Enter your password"
              icon="i-heroicons-lock-closed"
              required
          />
        </UFormGroup>

        <UButton
            type="submit"
            color="primary"
            block
            :loading="isLoading"
        >
          Sign in
        </UButton>
      </form>

      <template #footer>
        <div class="text-center text-sm">
          Not a member?
          <ULink
              to="/register"
              class="text-primary hover:text-primary-400 font-semibold"
          >
            Create an account
          </ULink>
        </div>
      </template>
    </UCard>
  </div>
</template>

<script setup>
const authStore = useAuthStore()
const email = ref('')
const password = ref('')
const isLoading = ref(false)

const handleLogin = async () => {
  isLoading.value = true
  try {
    const success = await authStore.login({
      username: email.value,
      password: password.value
    })

    if (!success) {
      // useToast().add({
      //   title: 'Login Failed',
      //   description: authStore.loginError || 'Unable to login',
      //   color: 'red'
      // })
    }
  } catch (error) {
    useToast().add({
      title: 'Login Error',
      description: 'An unexpected error occurred',
      color: 'red'
    })
  } finally {
    isLoading.value = false
  }
}


</script>