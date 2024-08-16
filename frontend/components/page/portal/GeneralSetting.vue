<template>
  <div class="">

    <DashboardToolbar>
      <template #left>
        <h1 class="flex items-center gap-1.5 font-semibold text-gray-900 dark:text-white min-w-0">
          <span class="truncate">General Settings</span>
        </h1>
      </template>

    </DashboardToolbar>



    <UTabs :items="items" class="w-full">
      <template #default_env="{ item }">
        <UCard @submit.prevent="onSubmitEnv">
          <template #header>
            <p class="text-base font-semibold leading-6 text-gray-900 dark:text-white">
              {{ item.label }}
            </p>
            <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
              All Default Env will be set here.
            </p>
          </template>

          <UFormGroup label="ENV" name="env" class="mb-3">
            <UTextarea
                :rows="10"
                color="primary"
                variant="outline"
                placeholder="Dedault Env in JSON {}"
                model-value="envForm.env"
            />
          </UFormGroup>

          <template #footer>
            <UButton type="submit" color="black">
              Save
            </UButton>
          </template>
        </UCard>
      </template>

      <template #others="{ item }">
        <UCard @submit.prevent="onSubmitPassword">
          <template #header>
            <h3 class="text-base font-semibold leading-6 text-gray-900 dark:text-white">
              {{ item.label }}
            </h3>
            <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
              Upcoming feature
            </p>
          </template>
        </UCard>
      </template>
    </UTabs>


  </div>
</template>

<script setup lang="ts">
const items = [{
  slot: 'default_env',
  label: 'Default ENV'
}, {
  slot: 'others',
  label: 'Others'
}]

const envForm = reactive({ env: ''})

function onSubmitEnv () {
  console.log('Submitted form:', envForm)
}

function onSubmitPassword () {
}
</script>
