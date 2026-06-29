<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { type RunOverview, type GameAsset, parseGameAsset } from '@/types'
import { formatNumber } from '@/utils/utils.ts'
import ClassTagComponent from './ClassTagComponent.vue'
import UpgradeListComponent from './UpgradeListComponent.vue'
import WeaponDetailsList from './WeaponDetailsList.vue'

const props = defineProps<{
  rdata: RunOverview
}>()

const rdamage = computed(() => {
  return formatNumber(props.rdata.PlayerDamage)
})

const roverkill = computed(() => {
  return formatNumber(props.rdata.OverkillDamage)
})

const weaponData = ref<GameAsset[]>([])
const upgradeData = ref<GameAsset[]>([])

const upgrades = computed<GameAsset[]>(() => {
  console.log(`Parsing player upgrades:`)
  console.log(upgradeData.value)
  return upgradeData.value.filter((upg) => {
    return upg.Reference.length == 0
  })
})

const fetchWeaponData = () => {
  fetch(`/api/getItemOverview?PlayerId=${props.rdata.PlayerId}&RunId=${props.rdata.RunId}`)
    .then((resp) => {
      if (resp.ok) {
        return resp.json()
      }
    })
    .then((json) => {
      if (Array.isArray(json)) {
        weaponData.value = json.map((item) => parseGameAsset(item))
      } else {
        weaponData.value = [parseGameAsset(json)]
      }
    })
}

const fetchUpgradeData = () => {
  fetch(`/api/getUpgrades?PlayerId=${props.rdata.PlayerId}&RunId=${props.rdata.RunId}`)
    .then((resp) => {
      if (resp.ok) {
        return resp.json()
      }
    })
    .then((json) => {
      if (Array.isArray(json)) {
        upgradeData.value = json.map((upg) => parseGameAsset(upg))
      } else {
        upgradeData.value = [parseGameAsset(json)]
      }
    })
}

onMounted(() => {
  fetchWeaponData()
  fetchUpgradeData()
})

const runTimeStr = computed(() => {
  const minutes = Math.floor(props.rdata.Runtime / 60)
  const seconds = props.rdata.Runtime % 60

  // Pad seconds so 9 seconds becomes "09" instead of "9"
  const paddedSeconds = String(seconds).padStart(2, '0')

  return `${minutes}:${paddedSeconds}`
})
</script>

<template>
  <button @click="$emit('CloseRunDetail')">Back</button>
  <div class="main-comp">
    <div class="twenty">
      <ClassTagComponent />
    </div>
    <div class="data-panel-section">
      <h2 v-if="rdata.Status" class="victory">{{ `Victory -- ${runTimeStr}` }}</h2>
      <h2 v-else class="defeat">Defeat -- {{ runTimeStr }}</h2>
      <div>
        <span>Damage: {{ rdamage }}</span>
        <span class="overkill"> ({{ roverkill }})</span>
      </div>
      <p>Killed: {{ rdata.PlayerKills }} Downed: {{ rdata.PlayerDeaths }}</p>
    </div>
    <div class="forty">
      <UpgradeListComponent :upgrades="upgrades" />
    </div>
    <div class="forty">
      <WeaponDetailsList :items="weaponData" :upgrades="upgradeData" />
    </div>
  </div>
</template>

<style>
.main-comp {
  display: flex;
  flex-direction: row;
}
.twenty {
  width: 10%;
}
.data-panel-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: top;
}
.victory {
  color: green;
}
.defeat {
  color: red;
}
.forty {
  width: 40%;
}
</style>
