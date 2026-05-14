<script setup>
import { ref } from 'vue'

const emit = defineEmits(['login-success', 'go-to-main', 'go-to-registration'])

const username = ref('')
const password = ref('')
const message = ref('')

async function handleLogin() {
    try {
        const response = await fetch("http://localhost:8081/api/login", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
                username: username.value,
                password: password.value
            })
        })

        const data = await response.json()

        if (response.ok) {
            localStorage.setItem('token', data.token)
            message.value = "Login success!"
            setTimeout(() => emit('login-success'), 1000)
        } else {
            message.value = "Error: " + (data.error || "wrong credentials")
            if (data.details) {
                console.error("Tech details:", data.details)
            }
        }
    } catch (error) {
        message.value = "Invalid username or password!"
    }
}
</script>

<template>
    <div class="max-w-md mx-auto mt-10 p-6 bg-white dark:bg-slate-800 rounded-xl shadow-lg">
        <h1 class="text-2xl font-bold mb-6 text-gray-800 dark:text-orange-100">Sing-in</h1>
        
        <form @submit.prevent="handleLogin" class="space-y-4">
        <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-orange-50">Username:</label>
            <input v-model="username" type="text" required
                class="w-full p-2 border rounded-lg outline-none 
                bg-white text-gray-950
                dark:bg-slate-900  dark:border-gray-600 dark:text-yellow-50 
                focus:ring-2 focus:ring-blue-500 dark:focus:ring-yellow-50">
        </div>

        <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-orange-50">Password:</label>
            <input v-model="password" type="password" required
                class="w-full p-2 border rounded-lg outline-none 
                bg-white text-gray-950
                dark:bg-slate-900  dark:border-gray-600 dark:text-yellow-50 
                focus:ring-2 focus:ring-blue-500 dark:focus:ring-yellow-50">
        </div>

        <button type="submit" 
                class="w-full py-2 rounded-lg font-semibold
                bg-green-600 text-white
                dark:bg-lime-600 dark:text-yellow-50 dark:hover:bg-lime-700 hover:bg-green-700 transition">
            Jump in
        </button>
        </form>

        <p v-if="message" :class="{
            'text-green-600 dark:text-green-400': message.includes('success'), 
            'text-red-600 dark:text-red-400': message.includes('Error') || message.includes('Invalid') || message.includes('Server')
        }" 
        class="mt-4 text-center text-sm font-medium">
        {{ message }}
        </p>

        <div class="mt-6 pt-6 border-t border-gray-100 dark:border-slate-600">
          <button @click="emit('go-to-registration')" class="w-full text-sm text-blue-600 hover:underline">
          Doesn't have a account? Register
          </button>
        </div>
    </div>
</template>