<script setup>
import { ref, onMounted } from 'vue'
import Registration from './components/Registration.vue'
import Login from './components/Login.vue'
import Purchases from './components/Purchases.vue'
import Profil from './components/Profil.vue'

const currentScreen = ref('login')
const isDark = ref(false)

function toggleTheme() {
  isDark.value = !isDark.value
  if (isDark.value) {
    document.documentElement.classList.add('dark')
    localStorage.setItem('theme', 'dark')
  } else {
    document.documentElement.classList.remove('dark')
    localStorage.setItem('theme', 'light')
  }
}

if (localStorage.getItem('token')) {
  currentScreen.value = 'main'
}

onMounted(() => {
  isDark.value = document.documentElement.classList.contains('dark')
})
</script>

<template>
  <div class="min-h-screen bg-gray-100 dark:bg-slate-900 py-10 transition-colors duration-300">
    <div class="max-w-screen mx-auto flex justify-end mb-4 px-12">
      <button @click="toggleTheme"
        class="relative inline-flex items-center p-2 rounded-full bg-white dark:bg-gray-800 shadow-md border border-gray-200 dark:border-gray-700 hover:scale-110 transition-all duration-300 group"
        title="Toggle Theme">
        <span v-if="isDark" class="text-yellow-400">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364-6.364l-.707.707M6.343 17.657l-.707.707m12.728 0l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707-.707M6.343 6.343l-.707-.707M12 8a4 4 0 100 8 4 4 0 000-8z" />
          </svg>
        </span>

        <span v-else class="text-indigo-600 group-hover:text-indigo-400">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
          </svg>
        </span>
      </button>
    </div>

    <Login v-if="currentScreen === 'login'" @login-success="currentScreen = 'main'"
      @go-to-registration="currentScreen = 'register'" />

    <Registration v-if="currentScreen === 'register'" @go-to-login="currentScreen = 'login'" />

    <Purchases v-if="currentScreen === 'main'" @logout="currentScreen = 'login'" @profile="currentScreen = 'profile'" />

    <Profil v-if="currentScreen === 'profile'" @go-to-main="currentScreen = 'main'" @logout="currentScreen = 'login'" />

  </div>
</template>