<template>
  <div>
    <DashboardPage>
      <DashboardPanel grow>
        <PagePortalGeneralSetting/>
      </DashboardPanel>
    </DashboardPage>
  </div>
</template>


<script setup>

definePageMeta({
  name: 'GeneralSettings',
  layout: 'portal',
  middleware: 'auth'
});

</script>

<style lang="scss" scoped>

</style>
