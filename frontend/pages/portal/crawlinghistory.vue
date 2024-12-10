<template>
  <div>
    <DashboardPage>
      <DashboardPanel grow>
        <PagePortalCrawlingHistory/>
      </DashboardPanel>
    </DashboardPage>
  </div>
</template>


<script setup>

definePageMeta({
  name: 'CrawlingHistory',
  layout: 'portal',
  middleware: 'auth'
});

</script>

<style lang="scss" scoped>

</style>
