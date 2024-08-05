<template>
  <span>
    <template v-if="item.is_active">
      {{item.text}}
    </template>
  </span>
</template>

<script setup>
const props = defineProps({
  item:{
    type:Object,
    required:true,
    default:()=>{
      return {
        text:'',
        is_active:false,
      }

    }
  },
})
</script>

<style lang="scss" scoped>

</style>
