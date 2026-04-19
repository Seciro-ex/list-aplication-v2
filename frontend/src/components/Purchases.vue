<script setup>
import { ref } from 'vue'

const emit = defineEmits(['login-success', 'go-to-main', 'go-to-registration'])

const username = ref('')
const password = ref('')
const message = ref('')

async function handleLogin() {
    try {
        const response = await fetch("http://localhost:8081/login", {
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
            message.value = "Success login!"
            setTimeout(() => emit('login-success'), 500)
        } else {
            message.value = "Error: " + (data.error || "wrong credentials")
        }
    } catch (error) {
        message.value = "Server doesn't respond"
    }
}
</script>

<template>
    <div class="max-w-md mx-auto mt-10 p-6 bg-white rounded-xl shadow-lg">
        <div class="mt-6 pt-6 border-t border-gray-100">
          <button @click="emit('go-to-login')" class="w-full mt-4 text-sm text-blue-600 hover:underline">
          Button
          </button>
        </div>
    </div>
</template>