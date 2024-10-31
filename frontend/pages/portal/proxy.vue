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
  layout: 'portal'
});

</script>

<style lang="scss" scoped>

</style>
