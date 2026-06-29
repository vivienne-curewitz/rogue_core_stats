<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { parseRunOverview, type RunOverview } from '../types'
import RunOverviewComponent from '@/components/RunOverviewComponent.vue'
import RunDetailComponent from '@/components/RunDetailComponent.vue'

const loading = ref(false)
const error = ref<string | null>(null)
const stats = ref<RunOverview[]>([])
const selectedPlayer = ref('Danger')

const selectedRdata = ref<RunOverview | null>(null)

async function fetchStats() {
  loading.value = true
  error.value = null
  try {
    const response = await fetch(`/api/overview?player_id=${selectedPlayer.value}`)
    if (response.ok) {
      const data = await response.json()
      if (Array.isArray(data) && data.length > 0) {
        stats.value = data.map(parseRunOverview)
      }
    }
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchStats()
})

const onRunClicked = (rd: RunOverview) => {
  selectedRdata.value = rd
}

const onCloseRunDetail = () => {
  selectedRdata.value = null
}

const totalRuns = computed(() => stats.value.length)
const winRate = computed(() => {
  if (stats.value.length === 0) return 0
  const wins = stats.value.filter((r) => r.Status).length
  return Math.round((wins / stats.value.length) * 100)
})

const totalKills = computed(() => stats.value.reduce((acc, r) => acc + Number(r.PlayerKills), 0))
const totalMinerals = computed(() =>
  stats.value.reduce((acc, r) => acc + Number(r.MineralsMined), 0).toFixed(1),
)
</script>

<template>
  <main class="animate-fade-in">
    <!-- Hero Header -->
    <div class="hero-section">
      <div class="hero-content">
        <h1>Rogue Core Mission Control</h1>
        <p>Real-time analytics and statistics dashboard for Deep Rock Galactic: Rogue Core.</p>
      </div>
      <div class="player-selector">
        <label for="player-input">Miner ID</label>
        <div class="input-group">
          <input
            id="player-input"
            v-model="selectedPlayer"
            type="text"
            placeholder="Miner Name"
            @keyup.enter="fetchStats"
          />
          <button @click="fetchStats" :disabled="loading">
            <span v-if="loading">Syncing...</span>
            <span v-else>Sync</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Stats Summary Grid -->
    <div class="stats-grid">
      <div class="stats-card">
        <div class="stats-label">Total Runs</div>
        <div class="stats-value">{{ totalRuns }}</div>
      </div>
      <div class="stats-card">
        <div class="stats-label">Win Rate</div>
        <div class="stats-value success">{{ winRate }}%</div>
      </div>
      <div class="stats-card">
        <div class="stats-label">Void Spawn Killed</div>
        <div class="stats-value danger">{{ totalKills }}</div>
      </div>
      <div class="stats-card">
        <div class="stats-label">Expenite Mined</div>
        <div class="stats-value warning">{{ totalMinerals }}</div>
      </div>
    </div>

    <!-- Main Content Panel -->
    <div v-if="!selectedRdata">
      <RunOverviewComponent
        v-for="rd in stats"
        :rdata="rd"
        :key="rd.RunId"
        @click="onRunClicked(rd)"
      />
    </div>
    <div v-else>
      <RunDetailComponent :rdata="selectedRdata" @CloseRunDetail="onCloseRunDetail" />
    </div>
  </main>
</template>

<style scoped>
.hero-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2.5rem;
  gap: 2rem;
  flex-wrap: wrap;
}

.hero-content h1 {
  font-size: 2.5rem;
  font-weight: 800;
  letter-spacing: -0.025em;
  margin-bottom: 0.5rem;
  background: linear-gradient(135deg, var(--text-primary) 30%, var(--text-secondary) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.hero-content p {
  color: var(--text-secondary);
  font-size: 1.1rem;
}

.player-selector {
  background: var(--panel-bg);
  border: 1px solid var(--panel-border);
  padding: 1rem 1.25rem;
  border-radius: 1rem;
  box-shadow: var(--shadow-sm);
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  min-width: 300px;
}

.player-selector label {
  font-size: 0.8rem;
  text-transform: uppercase;
  font-weight: 600;
  letter-spacing: 0.05em;
  color: var(--text-muted);
}

.input-group {
  display: flex;
  gap: 0.5rem;
}

.input-group input {
  flex: 1;
  background: rgba(0, 0, 0, 0.2);
  border: 1px solid var(--panel-border);
  color: var(--text-primary);
  padding: 0.6rem 1rem;
  border-radius: 0.5rem;
  font-family: var(--font-sans);
  font-size: 0.95rem;
  outline: none;
  transition: border-color var(--transition-fast);
}

.input-group input:focus {
  border-color: var(--primary);
}

.input-group button {
  background: var(--primary);
  color: #fff;
  border: none;
  padding: 0.6rem 1.25rem;
  border-radius: 0.5rem;
  font-weight: 600;
  cursor: pointer;
  transition:
    background-color var(--transition-fast),
    transform var(--transition-fast);
}

.input-group button:hover:not(:disabled) {
  background: var(--primary-hover);
  transform: translateY(-1px);
}

.input-group button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* Stats Cards Grid */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2.5rem;
}

.stats-card {
  background: var(--panel-bg);
  border: 1px solid var(--panel-border);
  border-radius: 1rem;
  padding: 1.5rem;
  box-shadow: var(--shadow-sm);
  transition:
    transform var(--transition-fast),
    border-color var(--transition-fast),
    box-shadow var(--transition-fast);
}

.stats-card:hover {
  transform: translateY(-2px);
  border-color: rgba(99, 102, 241, 0.3);
  box-shadow: var(--shadow-md);
}

.stats-label {
  font-size: 0.85rem;
  font-weight: 600;
  text-transform: uppercase;
  color: var(--text-secondary);
  letter-spacing: 0.05em;
  margin-bottom: 0.5rem;
}

.stats-value {
  font-size: 2.25rem;
  font-weight: 800;
  color: var(--text-primary);
  margin-bottom: 0.25rem;
}

.stats-value.success {
  color: var(--secondary);
}

.stats-value.danger {
  color: var(--danger);
}

.stats-value.warning {
  color: var(--warning);
}

.stats-footer {
  font-size: 0.8rem;
  color: var(--text-muted);
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.panel-header h2 {
  font-size: 1.5rem;
  font-weight: 700;
}

.refresh-btn {
  background: transparent;
  border: 1px solid var(--panel-border);
  color: var(--text-secondary);
  width: 2.25rem;
  height: 2.25rem;
  border-radius: 0.5rem;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all var(--transition-fast);
}

.refresh-btn:hover:not(:disabled) {
  color: var(--text-primary);
  background: var(--panel-border);
  transform: rotate(30deg);
}

/* Loading/Empty States */
.loading-state,
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 4rem 2rem;
  color: var(--text-secondary);
  gap: 1rem;
}

.spinner {
  width: 2.5rem;
  height: 2.5rem;
  border: 3px solid rgba(99, 102, 241, 0.1);
  border-top-color: var(--primary);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

/* Table styling */
.table-container {
  overflow-x: auto;
}

table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
}

th {
  color: var(--text-muted);
  font-weight: 600;
  text-transform: uppercase;
  font-size: 0.75rem;
  letter-spacing: 0.05em;
  padding: 1rem;
  border-bottom: 1px solid var(--panel-border);
}

td {
  padding: 1.1rem 1rem;
  border-bottom: 1px solid var(--panel-border);
  font-size: 0.95rem;
}

tr:hover td {
  background: rgba(255, 255, 255, 0.02);
}

@media (max-width: 768px) {
  .hero-section {
    flex-direction: column;
    align-items: stretch;
  }
}
</style>
