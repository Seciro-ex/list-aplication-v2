<script setup>
import { ref, onMounted, computed } from 'vue'

const emit = defineEmits(['go-to-main', 'logout'])
const message = ref('')
const userInfo = ref()
const formattedDate = computed(() => {
    if (!userInfo.value || !userInfo.value.created_at) return ''
    
    const date = new Date(userInfo.value.created_at)
    return date.toLocaleString('en-EN', {
        day: 'numeric',
        month: 'long',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    })
})

async function handleProfilInfo() {
    try {
        const token = localStorage.getItem('token')
        if (!token) {
            emit('logout')
            return
        }

        const response = await fetch("http://localhost:8081/api/me", {
            method: "GET",
            headers: { 
                "Content-Type": "application/json", 
                "Authorization": `Bearer ${token}`
            },
        })

        if (response.ok) {
            const data = await response.json()
            userInfo.value = data || []
        } else {
            message.value = "Error when loading a list"
        }
    } catch (error) {
        message.value = error
    }
    
}

onMounted(() => {
    handleProfilInfo()
})
</script>

<template>
    <div class="max-w-md mx-auto mt-10 p-6 bg-white rounded-xl shadow-lg border border-gray-100">
        <!-- 2. Используем v-if, чтобы не ловить ошибки, пока данные грузятся -->
        <div v-if="userInfo">
            <h2 class="text-xl font-bold mb-4 text-gray-800">User Profile</h2>
            <div class="bg-gray-50 p-4 rounded-lg mb-6">
                <p class="text-sm text-gray-500">User ID</p>
                <p class="text-lg font-mono text-gray-800">#{{ userInfo.id }}</p>
            </div>

            <div class="bg-gray-50 p-4 rounded-lg mb-4">
                <p class="text-sm text-gray-500">Username</p>
                <p class="text-lg break-all font-semibold text-gray-800">{{ userInfo.username }}</p>
            </div>

            <div class="bg-gray-50 p-4 rounded-lg mb-6">
                <p class="text-sm text-gray-500">Email</p>
                <p class="text-lg font-mono text-gray-800">{{ userInfo.email }}</p>
            </div>

            <div class="bg-gray-50 p-4 rounded-lg mb-6">
                <p class="text-sm text-gray-500">Created at</p>
                <p class="text-lg font-mono text-gray-800">{{ formattedDate }}</p>
            </div>
            <div class="bg-gray-50 p-4 rounded-lg mb-6">
                <img :src="userInfo.picture_url" alt="img">
            </div>
        </div>

        <div v-else-if="!message" class="text-center py-4">
            <p class="text-gray-400 animate-pulse">Loading profile...</p>
        </div>

        <div class="flex flex-rev">
            <div class="flex gap-3">
                <button class="w-full py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition font-medium">
                    Edit
                </button>
            </div>
            <div class="flex gap-3">
                <button class="w-full py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition font-medium">
                    Delete User
                </button> 
            </div>
        </div>

        <div class="flex flex-col gap-3">
            <button @click="emit('go-to-main')" 
                    class="w-full py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition font-medium">
                Back to List
            </button>
            
            <p v-if="message" class="text-red-500 text-center text-sm mt-2">{{ message }}</p>
        </div>
    </div>
</template>