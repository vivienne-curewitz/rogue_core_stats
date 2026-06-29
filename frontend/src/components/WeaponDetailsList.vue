<script setup lang="ts">
import { type GameAsset } from '@/types/'
import WeaponDetail from './WeaponDetail.vue'

defineProps<{
  items: GameAsset[]
  upgrades: GameAsset[]
}>()
</script>

<template>
  <div class="weapListMain">
    <WeaponDetail v-for="item in items" :key="item.UUID" :item="item" :upgrades="upgrades" />
  </div>
</template>

<style>
.weapListMain {
  display: flex;
  flex-direction: column;
  gap: 2em;
}
</style>
