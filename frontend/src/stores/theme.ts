import { ref } from 'vue'
import { defineStore } from 'pinia'

export const useThemeStore = defineStore('theme', () => {
  // Default theme is 'dark'
  const currentTheme = ref<'light' | 'dark'>(
    (localStorage.getItem('theme') as 'light' | 'dark') || 'dark'
  )

  function initTheme() {
    applyTheme(currentTheme.value)
  }

  function toggleTheme() {
    currentTheme.value = currentTheme.value === 'light' ? 'dark' : 'light'
    localStorage.setItem('theme', currentTheme.value)
    applyTheme(currentTheme.value)
  }

  function applyTheme(theme: 'light' | 'dark') {
    const root = document.documentElement
    if (theme === 'light') {
      root.classList.add('light-mode')
    } else {
      root.classList.remove('light-mode')
    }
  }

  return {
    currentTheme,
    initTheme,
    toggleTheme,
  }
})
