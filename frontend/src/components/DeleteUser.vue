<script setup>
import { ref } from 'vue'

const props = defineProps(['id', 'username'])
const emit = defineEmits(['user-deleted', 'cancel'])
const password = ref('')

const message = ref('')

async function VerifyUser() {
    try {
        const response = await fetch("http://localhost:8081/api/login", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
                username: props.username, 
                password: password.value
            })
        })

        if (!response.ok) {
            const data = await response.json()
            message.value = "Error: " + (data.error || "Wrong password")
            return
        }
        await DeleteUser()

    } catch (error) {
        message.value = "Wrong password"
    }
}

async function DeleteUser() {
    const url = `http://localhost:8081/api/delete-user?id=${props.id}`
    try {
        const token = localStorage.getItem('token')
        if (!token) return

        const response = await fetch(url, {
            method: "DELETE",
            headers: { 
                "Content-Type": "application/json", 
                "Authorization": `Bearer ${token}`
            },
        }) 
        if (response.ok) {
            message.value = "User deleted. You will be redirecting to login page!"
            setTimeout(() => emit('user-deleted'), 2000)
        } else {
            message.value = "Error to delete user"
        }
    } catch (error) {
        message.value = "Error: " + error
    }
}


</script>

<template>
    <div class="bg-red-50 p-5 rounded-lg mb-6 border-2 border-dashed border-red-300 shadow-inner">
        <h3 class="font-bold text-orange-800 mb-4 flex justify-center">DELETE USER</h3>
        <p class="mb-4">WARNING! This action is irreversible and it's automatically deletes all your purchases and account</p>
        
        <div class="space-y-3">
        <input v-model="password" type="password" placeholder="Enter your password"
            class="w-full p-2 border border-red-200 rounded outline-none focus:ring-2 focus:ring-red-400" />
        
        <button @click="VerifyUser" 
                class="w-full bg-red-500 text-white font-bold py-2 rounded shadow hover:bg-red-600 transition" :disabled="!password">
            Delete User
        </button>
        <button @click="emit('cancel')" 
                class="w-full bg-orange-500 text-white font-bold py-2 rounded shadow hover:bg-orange-600 transition">
            Cancel
        </button>
        </div>
        <p v-if="message" class="text-red-500 text-sm mt-3 text-center">{{ message }}</p>
    </div>
</template>