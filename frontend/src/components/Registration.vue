<script setup>
import { ref } from 'vue'

const emit = defineEmits(['register-success', 'go-to-login'])

const username = ref('')
const password = ref('')
const email = ref('')
const message = ref('')

const emailRegex = /^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$/

async function handleValidation() {
    if (username.value.length < 2) {
        message.value = "Error! Username must be at least 2 symbols long!"
        return
    }
    if (email.value.length < 2 || !emailRegex.test(email.value)) {
        message.value = "Error! Please enter a valid email address."
        return
    }
    if (password.value.length < 6) {
        message.value = "Error! Password must be at least 6 symbols long!"
        return
    }
    handleRegistration()
}

async function handleRegistration() {
    try {
        const response = await fetch("http://localhost:8081/api/register", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
                username: username.value,
                email: email.value,
                password: password.value
            })
        })

        const data = await response.json()

        if (response.ok) {
            message.value = data.message
            setTimeout(() => emit('go-to-login'), 2000)
        } else {
            message.value = data.error || "Unknow error"
            if (data.details) {
                console.error("Tech details:", data.details)
            }
        }
    } catch (error) {
        message.value = "Server Error"
    }
}
</script>

<template>
    <div class="max-w-md mx-auto mt-10 p-6 bg-white dark:bg-slate-800 rounded-xl shadow-lg">
        <h1 class="text-2xl font-bold mb-6 text-gray-800 dark:text-orange-100">Sign-up</h1>
        
        <form @submit.prevent="handleValidation" class="space-y-4">
        <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-orange-50">Username:</label>
            <input v-model="username" type="text" required
                class="w-full p-2 border rounded-lg outline-none 
                bg-white text-gray-950
                dark:bg-slate-900  dark:border-gray-600 dark:text-yellow-50 
                focus:ring-2 focus:ring-blue-500 dark:focus:ring-yellow-50">
        </div>

        <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-orange-50">E-mail:</label>
            <input v-model="email" type="email" required
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
                class="w-full py-2 rounded-lg font-semibold transition
                bg-green-600 text-white hover:bg-green-700 
                dark:bg-lime-600 dark:text-yellow-50 dark:hover:bg-lime-700">
            Create account
        </button>
        </form>

        <p v-if="message" :class="{'text-green-600 dark:text-green-400': !message.includes('Error'), 'text-red-600 dark:text-red-400': message.includes('Error')}" 
        class="mt-4 text-center text-sm font-medium">
        {{ message }}
        </p>

        <button @click="emit('go-to-login')" class="w-full mt-4 text-sm text-blue-600 hover:underline">
        Already have account? Login
        </button>
    </div>
</template>