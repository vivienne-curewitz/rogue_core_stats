<script setup lang="ts">
  defineProps(['className'])
  const isrc  = "/assets/exampleDwarf.png"
</script>


<template>
  <div>
    <img :src="isrc" width="100" height="100" />
    <p>{{ className }}</p>
  </div>
</template>
