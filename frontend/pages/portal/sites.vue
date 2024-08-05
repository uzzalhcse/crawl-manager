<template>
  <div>
    <DashboardPage>
      <DashboardPanel grow>
        <PagePortalSite/>
      </DashboardPanel>
    </DashboardPage>
  </div>
</template>


<script setup>

definePageMeta({
  name: 'Sites',
  layout: 'portal'
});

</script>

<style lang="scss" scoped>

</style>
