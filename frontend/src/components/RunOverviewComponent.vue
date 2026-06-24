<script setup lang="ts">
// import run overview types and store stuff here
import { type RunOverview } from "../types/index.ts"
import ClassTagComponent from "@/components/ClassTagComponent.vue"

// input expects data from the run
const props = defineProps<{
  rdata: RunOverview
}>()
</script>

<template>
  <ClassTagComponent className="DwarfGuy"/>
  <p>{{ rdata.PlayerDamage }} // {{ rdata.OverkillDamage }}</p>
  <p>{{ rdata.PlayerKills }} // {{ rdata.PlayerDeaths }} </p>
</template>
