<template>
  <div
    ref="el"
    :class="[ui.wrapper, grow ? ui.grow : ui.border, collapsible ? 'hidden lg:flex' : 'flex']"
  >
    <slot />

  </div>

</template>

<script lang="ts" setup>

import type { PropType } from 'vue';

const ui = {
  wrapper: 'flex-col items-stretch relative w-full',
  border: 'border-b lg:border-b-0 lg:border-r border-gray-200 dark:border-gray-800 lg:w-[--width] flex-shrink-0',
  grow: 'flex-1'
}
const props = defineProps({
  collapsible: {
    type: Boolean,
    default: false
  },
  grow: {
    type: Boolean,
    default: false
  },
  class: {
    type: [String, Object, Array] as PropType<any>,
    default: undefined
  }
})


</script>
