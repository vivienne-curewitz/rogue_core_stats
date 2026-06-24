<script setup lang="ts">
</script>

<template>
  <main class="animate-fade-in">
    <div class="about-container">
      <div class="about-header">
        <h1>About Rogue Stats</h1>
        <p class="subtitle">Deep Rock Galactic: Rogue Core Parser & Stats Tracker</p>
      </div>

      <div class="about-content">
        <div class="info-card">
          <h2>Project Overview</h2>
          <p>
            Rogue Stats is a stats tracker specifically built to read local save files from 
            <strong>Deep Rock Galactic: Rogue Core</strong>, parse key run data, and present structured metrics, 
            run history, and unlock progress in an interactive interface.
          </p>
        </div>

        <div class="info-card">
          <h2>How it Works</h2>
          <div class="steps">
            <div class="step">
              <div class="step-num">1</div>
              <div>
                <h3>Upload Saves</h3>
                <p>Submit your local <code>.sav</code> data to the stats parser.</p>
              </div>
            </div>
            <div class="step">
              <div class="step-num">2</div>
              <div>
                <h3>Automatic Processing</h3>
                <p>The Go backend parses the binary structures, extracting item links, ranks, and logs.</p>
              </div>
            </div>
            <div class="step">
              <div class="step-num">3</div>
              <div>
                <h3>View Analytics</h3>
                <p>Inspect minerals mined, kills, survival rate, and overall mission performance.</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </main>
</template>

<style scoped>
.about-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 1rem 0;
}

.about-header {
  text-align: center;
  margin-bottom: 3rem;
}

.about-header h1 {
  font-size: 2.5rem;
  font-weight: 800;
  margin-bottom: 0.5rem;
  background: linear-gradient(135deg, var(--text-primary) 30%, var(--text-secondary) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.subtitle {
  color: var(--text-secondary);
  font-size: 1.15rem;
}

.about-content {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.info-card {
  background: var(--panel-bg);
  border: 1px solid var(--panel-border);
  border-radius: 1.25rem;
  padding: 2rem;
  box-shadow: var(--shadow-md);
}

.info-card h2 {
  font-size: 1.4rem;
  font-weight: 700;
  margin-bottom: 1rem;
  color: var(--primary);
}

.info-card p {
  color: var(--text-secondary);
  line-height: 1.7;
  font-size: 1.05rem;
}

.steps {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  margin-top: 1rem;
}

.step {
  display: flex;
  gap: 1.25rem;
  align-items: flex-start;
}

.step-num {
  background: var(--primary-glow);
  color: var(--primary);
  border: 1.5px solid var(--panel-border);
  font-size: 1.1rem;
  font-weight: 700;
  width: 2.5rem;
  height: 2.5rem;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.step h3 {
  font-size: 1.1rem;
  font-weight: 600;
  margin-bottom: 0.25rem;
  color: var(--text-primary);
}

.step p {
  font-size: 0.95rem;
  color: var(--text-secondary);
  line-height: 1.5;
}

code {
  font-family: var(--font-mono);
  background: rgba(0, 0, 0, 0.2);
  padding: 0.2rem 0.4rem;
  border-radius: 0.25rem;
  font-size: 0.9rem;
  color: var(--accent);
}
</style>
