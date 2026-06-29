<script setup lang="ts">
import { type GameAsset } from '@/types'

defineProps<{
  upgrades: GameAsset[]
}>()
</script>

<template>
  <h2>Upgrades:</h2>
  <div class="listVert">
    <div v-for="upg in upgrades" :key="upg.UUID" class="rowItem">
      <img :src="`/assets/${upg.Asset}`" class="square-icon" />
      <span>{{ upg.Name }} : {{ upg.Description }}</span>
    </div>
  </div>
</template>

<style>
.listVert {
  display: flex;
  flex-direction: column;
  justify-content: left;
  align-items: center;
}

.rowItem {
  display: flex;
  flex-direction: row;
  justify-content: center;
  align-items: center;
  gap: 2em;
}

.square-icon {
  width: 2em;
  height: 2em;
  /* 'cover' crops the image slightly so it fills the box without distorting */
  object-fit: cover;
}
</style>
