<script setup>
import { ref, onMounted, computed } from 'vue'
import AddPurchases from './AddPurchases.vue'
import DeletePurchases from './DeletePurchases.vue'

const emit = defineEmits(['logout'])

const message = ref('')
const purchases = ref([])
const isFormVisible = ref(false)
const itemToDeleteId = ref(null)

const totalPrice = computed(() => {
    return purchases.value.reduce((sum, item) => {
        return  sum + Number(item.cost || 0)
    }, 0)
})

const bgColors = ['bg-red-400', 'bg-blue-400', 'bg-green-400', 'bg-yellow-400', 'bg-purple-400', 'bg-pink-400']
const getRandomColor = (id) => {
    return bgColors[id % bgColors.length]
}

async function loadPurchases() {
    try {
        const token = localStorage.getItem('token')
        
        if (!token) {
            emit('logout')
            return
        }

        const response = await fetch("http://localhost:8081/user/purchases", {
            method: "GET",
            headers: { 
                "Content-Type": "application/json", 
                "Authorization": `Bearer ${token}`
            },
        })

        if (response.ok) {
            const data = await response.json()
            purchases.value = data || []
        } else if (response.status === 401) {
            message.value = "Session ended. Please login again."
            localStorage.removeItem('token')
            emit('logout')
        } else {
            message.value = "Error when loading a list"
        }
    } catch (error) {
        message.value = "Server doesn't respond"
    }
}

function handleLogout() {
    localStorage.removeItem('token')
    emit('logout')
}

function handlePurchaseAdded() {
    isFormVisible.value = false
    loadPurchases()
}

function handlePurchaseDeleted() {
    itemToDeleteId.value = null 
    loadPurchases()
}

onMounted(() => {
    loadPurchases()
})
</script>

<template>
<div class="max-w-2xl mx-auto mt-10 p-6 bg-white rounded-xl shadow-lg border border-gray-100">
        <div class="flex justify-between items-center mb-4 border-b pb-4">
            <h1 class="text-2xl font-bold text-gray-800">My Shoping List</h1>
            <div @click="" class="w-16 h-16 rounded-full flex items-center justify-center text-white font-bold text-xl bg-pink-100 uppercase shadow-sm hover:bg-red-100">
            </div>
        </div>

        <div class="flex gap-2 pb-4 mb-6 border-b justify-between">
            <button @click="isFormVisible = !isFormVisible" 
                    class="text-sm bg-green-50 text-green-600 px-4 py-2 rounded-lg hover:bg-green-100 transition">
                {{ isFormVisible ? 'Close' : 'Add New' }}
            </button>
            <button @click="handleLogout" 
                    class="text-sm bg-red-50 text-red-600 px-4 py-2 rounded-lg hover:bg-red-100 transition">
                Sign-out
            </button>
        </div>
        <AddPurchases v-if="isFormVisible" @added-item="handlePurchaseAdded" />

        <p v-if="message" class="text-red-500 text-center mb-4">{{ message }}</p>

        <div class="flex justify-end p-4 bg-gray-50 rounded-lg shadow-sm mb-4">
            <h1>Total price: {{ totalPrice }} €</h1>
        </div>

        <div class="space-y-4">
            <div v-for="item in purchases" :key="item.id" 
                class="flex items-center justify-between p-4 bg-gray-100 rounded-lg shadow-sm">
                <div>
                <div class="flex items-center gap-4">
                    <img v-if="item.image_url" :src="item.image_url" alt="item" class="w-16 h-16 object-cover rounded-md border">
                    <div v-else 
                        :class="['w-16 h-16 rounded-md flex items-center justify-center text-white font-bold text-xl uppercase shadow-sm', getRandomColor(item.id)]">
                        {{ item.item_name[0] }}
                    </div>

                    <div>
                        <p class="font-semibold text-gray-800 text-lg">{{ item.item_name }}</p>
                        <a v-if="item.item_url" :href="item.item_url" target="_blank" class="text-sm text-blue-500 hover:underline">
                            Go to shop
                        </a>
                        <button @click="" 
                            class="text-sm  bg-orange-200 rounded-lg p-2 hover:bg-orange-500 transition">
                            Edit
                        </button>
                        <button @click="itemToDeleteId = item.id" 
                                class="text-sm ml-2 bg-red-200 rounded-lg p-2 hover:bg-red-500 transition">
                            Delete
                        </button>
                    </div>
                </div>
                <DeletePurchases 
                    v-if="itemToDeleteId === item.id" 
                    :itemId="item.id" 
                    @item-deleted="handlePurchaseDeleted"
                    @cancel="itemToDeleteId = null" 
                />
                </div>
                
                <div class="text-xl font-bold text-green-600">
                    {{ item.cost }} €
                </div>
            </div>
            <div v-if="purchases.length === 0 && !message" class="text-center text-gray-500 py-8">
                Your list is empty. Time to add something!
            </div>
        </div>
    </div>
</template>