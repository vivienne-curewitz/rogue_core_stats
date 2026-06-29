<script setup lang="ts">
import { type GameAsset } from '@/types'
import { computed } from 'vue'
const props = defineProps<{
  item: GameAsset
  upgrades: GameAsset[]
}>()

const refUpgrades = computed<GameAsset[]>(() => {
  return props.upgrades.filter((upg) => {
    return upg.Reference == props.item.Reference
  })
})
</script>

<template>
  <div class="main-comp">
    <img :src="`/assets/${item.Asset}`" class="big-icon" />
    <div class="upgrade-list">
      <p v-for="upg in refUpgrades" :key="upg.UUID">{{ upg.Name }} :: {{ upg.Description }}</p>
    </div>
  </div>
</template>

<style>
.big-icon {
  height: 10em;
  width: auto;
}

.upgrade-list {
  display: flex;
  flex-direction: column;
}

.main-comp {
  display: flex;
  flex-direction: row;
}
</style>
