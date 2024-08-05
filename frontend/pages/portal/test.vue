<template>
  <div class="">
<!--    <button @click="generatePdf">Generate PDF</button>-->
<!--    pdfUrl {{pdfUrl}}-->
    <div ref="htmlContent">
      <ResumePreview/>
    </div>
  </div>
</template>
<script setup>
definePageMeta({
  name: 'Test',
  layout: 'blank',
})
</script>

<script>

export default {
  data() {
    return {
      pdfUrl: null,
    };
  },
};
</script>
