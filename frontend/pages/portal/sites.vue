<template>
  <div>
    <DashboardPage>
      <DashboardPanel grow>
        <LazyPagePortalSite/>
      </DashboardPanel>
    </DashboardPage>
  </div>
</template>


<script setup>

definePageMeta({
  name: 'Sites',
  layout: 'portal',
  middleware: 'auth'
});

</script>

<style lang="scss" scoped>

</style>
