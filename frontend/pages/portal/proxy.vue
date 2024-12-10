<template>
  <div>
    <DashboardPage>
      <DashboardPanel grow>
        <LazyPagePortalProxy/>
      </DashboardPanel>
    </DashboardPage>
  </div>
</template>


<script setup>

definePageMeta({
  name: 'Proxy',
  layout: 'portal',
  middleware: 'auth'
});

</script>

<style lang="scss" scoped>

</style>
