<template>
  <div>
    <AuthLogin />
  </div>
</template>

<script setup>

// Optional: Add page metadata
definePageMeta({
  name: 'Login',
  layout: 'auth',
  middleware: ['auth']
})

</script>

<style lang="scss" scoped>

</style>
