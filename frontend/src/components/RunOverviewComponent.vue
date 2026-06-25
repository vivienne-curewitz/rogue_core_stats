<script setup lang="ts">
// import run overview types and store stuff here
import { computed } from 'vue';
import { type RunOverview } from "../types/index.ts"
import ClassTagComponent from "@/components/ClassTagComponent.vue"
import WeaponOverviewComponent from "./WeaponOverviewComponent.vue";

// input expects data from the run
const props = defineProps<{
  rdata: RunOverview
}>();

const runTimeStr = computed(() => {
  const minutes = Math.floor(props.rdata.Runtime / 60);
  const seconds = props.rdata.Runtime % 60;
  
  // Pad seconds so 9 seconds becomes "09" instead of "9"
  const paddedSeconds = String(seconds).padStart(2, '0'); 
  
  return `${minutes}:${paddedSeconds}`;
});
</script>

<template>
  <div class="data-panel">
    <div class="data-panel-internal">
      <ClassTagComponent className="DwarfGuy"/>
      <div class="data-panel-section">
        <h2 v-if="rdata.Status" class="victory">{{ `Victory -- ${runTimeStr}` }}</h2>
        <h2 v-else class="defeat">Defeat -- {{ runTimeStr }}</h2>
        <p>Damage: {{ rdata.PlayerDamage }} Overkill: {{ rdata.OverkillDamage }}</p>
        <p>Killed: {{ rdata.PlayerKills }} Downed: {{ rdata.PlayerDeaths }} </p>
      </div>
      <WeaponOverviewComponent :-run-id="rdata.RunId" :-player-id="rdata.PlayerId" />
    </div>
  </div>
</template>

<style>
  /* Data Panel */
  .data-panel {
    background: var(--panel-bg);
    border: 1px solid var(--panel-border);
    border-radius: 1.25rem;
    padding: 1.75rem;
    box-shadow: var(--shadow-md);
    width: 100%;
    height: 100%;
    margin-bottom: 2rem;
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 5%;
  }

  .data-panel-internal {
    display: flex;
    flex-direction: row;
    justify-content: center;
    align-items: center;
    width: 100%;
    height: 100%;
    gap: 5%;
  }

  .data-panel-section {
    display: flex;
    flex-direction: column;
    width: 20%;
    justify-content: left;
    align-items: left;
  }

  .victory {
    color: green;
  }
  
  .defeat {
    color: red;
  }
</style>