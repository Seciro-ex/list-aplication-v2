<script setup>
import { ref, onMounted, computed } from 'vue'
import AddPurchases from './AddPurchases.vue'
import DeletePurchases from './DeletePurchases.vue'
import EditPurchases from './EditPurchases.vue'
import Profil from './Profil.vue'

const emit = defineEmits(['logout', 'profile'])

const message = ref('')
const purchases = ref([])
const profile_url = ref('')

const isFormVisible = ref(false)
const itemToDeleteId = ref(null)
const itemToEditId = ref(null)
const ProfilOpen = ref(false)

const totalPrice = computed(() => {
    return purchases.value.reduce((sum, item) => {
        return  sum + Number(item.cost || 0)
    }, 0)
})

const bgColors = ['bg-red-400', 'bg-blue-400', 'bg-green-400', 'bg-yellow-400', 'bg-purple-400', 'bg-pink-400']
const getRandomColor = (id) => {
    return bgColors[id % bgColors.length]
}
async function GetProfilePick() {
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
            profile_url.value = data.picture_url || ""
        } else {
            message.value = "Error when loading a list"
        }
    } catch (error) {
        message.value = error
    }
}

async function loadPurchases() {
    try {
        const token = localStorage.getItem('token')
        
        if (!token) {
            emit('logout')
            return
        }

        const response = await fetch("http://localhost:8081/api/user/purchases", {
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
            handleLogout()
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

function handlePurchasesEdit() {
    itemToEditId.value = null
    loadPurchases()
}

onMounted(() => {
    loadPurchases()
    GetProfilePick()
})
</script>

<template>
<div class="max-w-2xl mx-auto mt-10 p-6 bg-white rounded-xl shadow-lg border border-gray-100">
        <div class="flex justify-between items-center mb-4 border-b pb-4">
            <h1 class="text-2xl font-bold text-gray-800">My Shoping List</h1>
            <div @click="emit('profile')" class="cursor-pointer">
                <img v-if="profile_url" 
                    :src="profile_url" 
                    class="w-16 h-16 rounded-full object-cover shadow-sm border border-gray-200 hover:scale-105 transition" 
                    alt="Profile">
                
                <div v-else 
                    class="w-16 h-16 rounded-full bg-indigo-500 flex items-center justify-center text-white font-bold text-xl shadow-sm hover:bg-indigo-600 transition">
                    U
                </div>
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

        <div class="flex justify-end p-4 bg-gray-50 rounded-lg shadow-sm mb-4">
            <h1>Total price: {{ totalPrice }} €</h1>
        </div>

        <div class="space-y-4">
            <div v-for="item in purchases" :key="item.id">
                
                <div v-if="itemToEditId !== item.id" 
                    class="flex items-center justify-between p-4 bg-gray-50 rounded-lg shadow-sm border border-gray-100 hover:shadow-md transition">
                    
                    <div class="flex items-center gap-4">
                        <div class="flex items-center gap-6">
                            <div class="w-32 h-20 shrink-0 bg-gray-50 rounded-xl border border-gray-100 flex items-center justify-center p-2 overflow-hidden shadow-sm">
                                
                                <img v-if="item.image_url" 
                                    :src="item.image_url" 
                                    @error="item.image_url = ''"
                                    alt="item" 
                                    class="w-full h-full object-contain" 
                                >
                                <div v-else 
                                    :class="['w-full h-full rounded-lg flex items-center justify-center text-white font-bold text-xl uppercase', getRandomColor(item.id)]">
                                    {{ item.item_name[0] }}
                                </div>
                            </div>

                            <div class="grow min-w-0">
                                <p class="font-bold text-gray-900 text-lg truncate">{{ item.item_name }}</p>
                                
                                <a v-if="item.item_url" :href="item.item_url" target="_blank" class="text-xs text-blue-500 hover:underline block mb-2">
                                    Go to shop
                                </a>

                                <div class="flex items-center gap-2">
                                    <template v-if="itemToDeleteId !== item.id">
                                        <button @click="itemToEditId = item.id" class="text-[10px] bg-orange-100 text-orange-700 px-2 py-1 rounded uppercase font-bold hover:bg-orange-200">Edit</button>
                                        <button @click="itemToDeleteId = item.id" class="text-[10px] bg-red-100 text-red-700 px-2 py-1 rounded uppercase font-bold hover:bg-red-200">Delete</button>
                                    </template>
                                    <div>
                                        <div class=" flex items-center gap-2">
                                            <DeletePurchases 
                                                v-if="itemToDeleteId === item.id" 
                                                :itemId="item.id" 
                                                @item-deleted="handlePurchaseDeleted"
                                                @cancel="itemToDeleteId = null" 
                                            />
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                    
                    <div class="text-xl font-bold text-green-600 whitespace-nowrap ml-4">
                        {{ item.cost }} €
                    </div>
                </div>

                
                <EditPurchases 
                    v-if="itemToEditId === item.id"
                    :itemId="item.id"
                    :cost="item.cost"
                    :image_url="item.image_url"
                    :item_name="item.item_name"
                    :item_url="item.item_url"
                    @item-edited="handlePurchasesEdit"
                    @cancel="itemToEditId = null"
                />

            </div>
            <div v-if="purchases.length === 0 && !message" class="text-center text-gray-500 py-8">
                Your list is empty. Time to add something!
            </div>
        </div>
    </div>
    <Profil v-if="ProfilOpen" />
</template>