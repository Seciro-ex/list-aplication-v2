import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import { useTheme } from './composables/useTheme'

const { initializeTheme } = useTheme()
initializeTheme()

createApp(App).mount('#app')
