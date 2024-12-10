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
          Create an account
        </h2>
      </template>

      <!-- Show registration error if exists -->
      <UAlert
          v-if="authStore.loginError"
          color="red"
          :title="authStore.loginError"
          class="mb-4"
      />

      <form @submit.prevent="handleRegister" class="space-y-6">
        <UFormGroup label="Full Name" name="name">
          <UInput
              v-model="name"
              type="text"
              placeholder="Enter your full name"
              icon="i-heroicons-user"
              required
          />
        </UFormGroup>
        <UFormGroup label="Username" name="username">
          <UInput
              v-model="username"
              type="text"
              placeholder="Enter your username"
              icon="i-heroicons-user"
              required
          />
        </UFormGroup>

        <UFormGroup label="Email address" name="email">
          <UInput
              v-model="email"
              type="email"
              placeholder="Enter your email"
              icon="i-heroicons-envelope"
              required
          />
        </UFormGroup>

        <UFormGroup label="Password" name="password">
          <UInput
              v-model="password"
              type="password"
              placeholder="Create a password"
              icon="i-heroicons-lock-closed"
              required
          />
        </UFormGroup>

        <UFormGroup label="Confirm Password" name="confirmPassword">
          <UInput
              v-model="confirmPassword"
              type="password"
              placeholder="Confirm your password"
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
          Create Account
        </UButton>
      </form>

      <template #footer>
        <div class="text-center text-sm">
          Already have an account?
          <ULink
              to="/login"
              class="text-primary hover:text-primary-400 font-semibold"
          >
            Sign in
          </ULink>
        </div>
      </template>
    </UCard>
  </div>
</template>

<script setup>
const authStore = useAuthStore()
const name = ref('')
const email = ref('')
const username = ref('')
const password = ref('')
const confirmPassword = ref('')
const isLoading = ref(false)

const handleRegister = async () => {
  // Basic password confirmation
  if (password.value !== confirmPassword.value) {
    useToast().add({
      title: 'Registration Failed',
      description: 'Passwords do not match',
      color: 'red'
    })
    return
  }

  isLoading.value = true
  try {
    const success = await authStore.register({
      name: name.value,
      email: email.value,
      username:username.value,
      password: password.value
    })

    if (!success) {
      useToast().add({
        title: 'Registration Failed',
        description: authStore.loginError || 'Unable to create account',
        color: 'red'
      })
    }
  } catch (error) {
    useToast().add({
      title: 'Registration Error',
      description: 'An unexpected error occurred',
      color: 'red'
    })
  } finally {
    isLoading.value = false
  }
}
</script>