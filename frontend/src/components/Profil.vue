<script setup>
import { ref, onMounted, computed } from 'vue'
import EditUser from './EditUser.vue'
import DeleteUser from './DeleteUser.vue'

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
const UserEdit = ref(false)
const UserDelete = ref(false)

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

        const data = await response.json()

        if (response.ok) {
            userInfo.value = data || []
        } else {
            message.value = data.error || "Error when loading profile"
            if (data.details) {
                console.error("Tech details:", data.details)
            }
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
        <EditUser 
            v-if="UserEdit === true" 
            :username="userInfo.username"
            :id="userInfo.id"
            :picture_url="userInfo.picture_url"
            :email="userInfo.email"
            @user-edited="() => { handleProfilInfo(); UserEdit = false }"
            @cancel="UserEdit = false" 
        />
        <DeleteUser
            v-if="UserDelete === true"
            :id="userInfo.id"
            :username="userInfo.username"
            @user-deleted="emit('logout')"
            @cancel="UserDelete = false"
        />
        <div v-if="userInfo && UserEdit == false && UserDelete == false">
            <h2 class="text-xl font-bold mb-4 text-gray-800">User Profile</h2>
            
            <div class="bg-white rounded-lg flex items-center justify-center">
                <img class="mb-4 w-24 h-24 rounded-full flex items-center justify-center text-white font-bold text-xl bg-indigo-500 uppercase shadow-sm hover:scale-125 transition-all" :src="userInfo.picture_url" alt="U">
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
        </div>

        <div v-else-if="!message && UserEdit == false && UserDelete == false" class="text-center py-4">
            <p class="text-gray-400 animate-pulse">Loading profile...</p>
        </div>

        <div class="flex items-center justify-between mb-2">
            <div class="flex w-full mr-1">
                <button @click="UserEdit = !UserEdit, UserDelete = false" class="w-full py-2 bg-orange-100 text-orange-700 rounded-lg hover:bg-orange-200 transition font-medium">
                    Edit
                </button>
            </div>
            <div class="flex w-full ml-1">
                <button @click="UserDelete = !UserDelete, UserEdit = false" class="w-full py-2 bg-red-100 text-red-700 rounded-lg hover:bg-red-200 transition font-medium">
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