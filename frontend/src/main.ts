import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import './index.css'
import { useThemeStore } from './stores/theme'

const app = createApp(App)

app.use(createPinia())
app.use(router)

// Initialize the theme store (defaults to dark mode)
const themeStore = useThemeStore()
themeStore.initTheme()

app.mount('#app')
