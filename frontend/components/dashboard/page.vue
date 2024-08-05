<template>
  <div :class="ui.wrapper">
    <slot />
  </div>
</template>

<script lang="ts" setup>

const ui = {
  wrapper: 'flex flex-1 w-full min-w-0'
}


</script>
