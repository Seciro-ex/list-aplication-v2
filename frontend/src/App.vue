<script setup>
import { ref } from 'vue'
import Registration from './components/Registration.vue' 
import Login from './components/Login.vue'
import Purchases from './components/Purchases.vue'
import Profil from './components/Profil.vue'

const currentScreen = ref('login')

if (localStorage.getItem('token')) {
  currentScreen.value = 'main'
}
</script>

<template>
  <div class="min-h-screen bg-gray-100 py-10">
    
    <Login 
      v-if="currentScreen === 'login'" 
      @login-success="currentScreen = 'main'" 
      @go-to-registration="currentScreen = 'register'" 
    />

    <Registration 
      v-if="currentScreen === 'register'" 
      @go-to-login="currentScreen = 'login'" 
    />

    <Purchases 
      v-if="currentScreen === 'main'" 
      @logout="currentScreen = 'login'" 
      @profile="currentScreen = 'profile'"
    />

    <Profil 
      v-if="currentScreen === 'profile'"
      @go-to-main="currentScreen = 'main'"
      @logout="currentScreen = 'login'" 
    />

  </div>
</template>