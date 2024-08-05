<template>
  <div :class="ui.wrapper">
    <div :class="ui.container">
      <slot>
        <div :class="ui.left">
          <slot name="left" />
        </div>

        <div :class="ui.right">
          <slot name="right" />
        </div>
      </slot>
    </div>
  </div>
</template>

<script lang="ts" setup>

const ui = {
  wrapper: 'min-h-[49px] flex-shrink-0 flex items-center border-b border-gray-200 dark:border-gray-800 px-4 py-2 gap-x-4',
  container: 'flex items-center justify-between flex-1 gap-x-1.5',
  left: 'flex items-stretch gap-1.5',
  right: 'flex items-stretch gap-1.5'
}
</script>
