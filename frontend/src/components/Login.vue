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
            message.value = data.message
            setTimeout(() => emit('login-success'), 500)
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
    <div class="max-w-md mx-auto mt-10 p-6 bg-white rounded-xl shadow-lg">
        <h1 class="text-2xl font-bold mb-6 text-gray-800">Sing-in</h1>
        
        <form @submit.prevent="handleLogin" class="space-y-4">
        <div>
            <label class="block text-sm font-medium text-gray-700">Username:</label>
            <input v-model="username" type="text" required
                class="w-full p-2 border rounded-lg focus:ring-2 focus:ring-blue-500 outline-none">
        </div>

        <div>
            <label class="block text-sm font-medium text-gray-700">Password:</label>
            <input v-model="password" type="password" required
                class="w-full p-2 border rounded-lg focus:ring-2 focus:ring-blue-500 outline-none">
        </div>

        <button type="submit" 
                class="w-full bg-green-600 text-white py-2 rounded-lg font-semibold hover:bg-green-700 transition">
            Jump in
        </button>
        </form>

        <p v-if="message" :class="{'text-green-600': !message.includes('Login'), 'text-red-600': message.includes('Server') || message.includes('Error') || message.includes('Invalid')}" 
        class="mt-4 text-center text-sm font-medium">
        {{ message }}
        </p>

        <div class="mt-6 pt-6 border-t border-gray-100">
          <button @click="emit('go-to-registration')" class="w-full mt-4 text-sm text-blue-600 hover:underline">
          Doesn't have a account? Register
          </button>
        </div>
    </div>
</template>