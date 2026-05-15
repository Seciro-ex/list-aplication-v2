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
const search = ref('')

const isFormVisible = ref(false)
const itemToDeleteId = ref(null)
const itemToEditId = ref(null)
const ProfilOpen = ref(false)

const filteredPurchases = computed(() => {
    let result = purchases.value

    if (search.value.trim()) {
        const query = search.value.toLowerCase()
        result = result.filter(item => {
            return item.item_name?.toLowerCase().includes(query) ||
                item.item_description?.toLowerCase().includes(query)
        })
    }

    return [...result].sort((a, b) => a.item_name.localeCompare(b.item_name))
})

const totalPrice = computed(() => {
    return filteredPurchases.value.reduce((sum, item) => {
        return sum + Number(item.cost || 0)
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

        const response = await fetch("/api/me", {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${token}`
            },
        })

        const data = await response.json()

        if (response.ok) {
            profile_url.value = data.picture_url || ""
        } else {
            message.value = data.error || "Error when loading a picture"
            if (data.details) {
                console.error("Tech details:", data.details)
            }
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

        const response = await fetch("/api/user/purchases", {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${token}`
            },
        })

        const data = await response.json()

        if (response.ok) {
            purchases.value = data || []
        } else {
            message.value = data.error || "Error when loading a list"
            if (response.status === 401) handleLogout()
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
    <div
        class="max-w-2xl mx-auto mt-10 p-6 bg-white dark:bg-slate-800 rounded-xl shadow-lg border border-slate-100 dark:border-slate-700">
        <div class="flex justify-between items-center mb-4 border-b border-slate-100 dark:border-slate-700 pb-4">
            <h1 class="text-2xl font-bold text-slate-800 dark:text-white">My Shoping List</h1>
            <div @click="emit('profile')" class="cursor-pointer">
                <img v-if="profile_url" :src="profile_url"
                    class="w-16 h-16 rounded-full object-cover shadow-sm border border-slate-200 dark:border-slate-600 hover:scale-105 transition"
                    alt="Profile">

                <div v-else
                    class="w-16 h-16 rounded-full bg-indigo-500 flex items-center justify-center text-white font-bold text-xl shadow-sm hover:bg-indigo-600 transition">
                    U
                </div>
            </div>
        </div>

        <div class="flex gap-2 pb-4 mb-6 border-b border-slate-100 dark:border-slate-700 justify-between">
            <button @click="isFormVisible = !isFormVisible" class="text-sm font-medium px-4 py-2 rounded-lg
                           bg-emerald-50 text-emerald-600 hover:bg-emerald-100 
                           dark:bg-emerald-500/20 dark:text-emerald-400 dark:hover:bg-emerald-500/30">
                {{ isFormVisible ? 'Close' : 'Add New' }}
            </button>
            <button @click="handleLogout" class="text-sm font-medium px-4 py-2 rounded-lg
                           bg-rose-50 text-rose-600 hover:bg-rose-100 
                           dark:bg-rose-500/20 dark:text-rose-400 dark:hover:bg-rose-500/30">
                Sign-out
            </button>
        </div>
        <Transition name="slide-fade">
            <div v-if="isFormVisible">
                <AddPurchases @added-item="handlePurchaseAdded" />
            </div>
        </Transition>

        <div
            class="flex justify-between p-4 bg-slate-50 dark:bg-slate-900/50 rounded-lg shadow-sm mb-4 border border-transparent">
            <input v-model="search" type="text" placeholder="Search..."
                class="p-2 border rounded-lg outline-none
                          bg-white text-slate-900 border-slate-200 focus:ring-2 focus:ring-blue-500 
                          dark:bg-slate-800 dark:text-white dark:border-slate-600 dark:focus:ring-blue-500 dark:placeholder-slate-400">
            <h1 class="flex items-center font-semibold text-slate-800 dark:text-slate-200">
                Total price: <span class="ml-2 text-emerald-600 dark:text-emerald-400">{{ totalPrice }} €</span>
            </h1>
        </div>

        <div class="space-y-4">
            <div v-for="item in filteredPurchases" :key="item.id">
                <Transition name="swap" mode="out-in">
                    <div v-if="itemToEditId !== item.id" key="view"
                        class="group relative flex items-center justify-between p-4 rounded-lg shadow-sm border transition-all duration-300
                        bg-white dark:bg-slate-800 border-slate-100 dark:border-slate-700 hover:shadow-md dark:hover:border-slate-600 gap-4">

                        <div class="flex items-center gap-4 min-w-0 flex-1">
                            <div class="w-24 h-24 sm:w-32 sm:h-32 shrink-0 rounded-xl border flex items-center justify-center p-2 overflow-hidden shadow-sm
                                    bg-slate-50 dark:bg-slate-900 border-slate-100 dark:border-slate-700">
                                <img v-if="item.image_url" :src="item.image_url" alt="item"
                                    class="w-full h-full object-contain">
                                <div v-else
                                    :class="['w-full h-full rounded-lg flex items-center justify-center text-white font-bold text-xl uppercase', getRandomColor(item.id)]">
                                    {{ item.item_name[0] }}
                                </div>
                            </div>

                            <div class="flex-1 min-w-0 flex flex-col justify-center">
                                <p class="font-bold text-lg text-slate-900 dark:text-white truncate"
                                    :title="item.item_name">
                                    {{ item.item_name }}
                                </p>

                                <a v-if="item.item_url" :href="item.item_url" target="_blank"
                                    class="text-xs text-blue-500 dark:text-blue-400 hover:underline block mb-2">
                                    Go to shop
                                </a>

                                <div class="flex flex-wrap items-center gap-2 mt-1">
                                    <template v-if="itemToDeleteId !== item.id">
                                        <button @click="itemToEditId = item.id"
                                            class="text-[10px] px-2 py-1 rounded uppercase font-bold transition-colors
                                                    bg-amber-100 text-amber-700 hover:bg-amber-200
                                                    dark:bg-amber-500/20 dark:text-amber-400 dark:hover:bg-amber-500/30">
                                            Edit
                                        </button>
                                        <button @click="itemToDeleteId = item.id" class="text-[10px] px-2 py-1 rounded uppercase font-bold transition-colors
                                                    bg-rose-100 text-rose-700 hover:bg-rose-200
                                                    dark:bg-rose-500/20 dark:text-rose-400 dark:hover:bg-rose-500/30">
                                            Delete
                                        </button>
                                    </template>
                                    <DeletePurchases v-if="itemToDeleteId === item.id" :itemId="item.id"
                                        @item-deleted="handlePurchaseDeleted" @cancel="itemToDeleteId = null" />
                                </div>

                                <p v-if="item.item_description"
                                    class="text-sm mt-1 max-h-0 overflow-hidden group-hover:max-h-20 transition-all duration-500 text-slate-500 dark:text-slate-400 truncate">
                                    {{ item.item_description }}
                                </p>
                            </div>
                        </div>

                        <div
                            class="text-xl font-bold whitespace-nowrap shrink-0 text-emerald-600 dark:text-emerald-400">
                            {{ item.cost }} €
                        </div>
                    </div>

                    <EditPurchases v-else key="edit" :itemId="item.id" :cost="item.cost" :image_url="item.image_url"
                        :item_name="item.item_name" :item_url="item.item_url" :item_description="item.item_description"
                        @item-edited="handlePurchasesEdit" @cancel="itemToEditId = null" />
                </Transition>

            </div>

            <div v-if="purchases.length === 0 && !message" class="text-center py-8 text-slate-500 dark:text-slate-400">
                Your list is empty. Time to add something!
            </div>
            <div v-if="filteredPurchases.length === 0 && search"
                class="text-center py-8 text-slate-500 dark:text-slate-400">
                Nothing found for "{{ search }}"
            </div>
        </div>
    </div>
    <Profil v-if="ProfilOpen" />
</template>

<style>
.slide-fade-enter-from,
.slide-fade-leave-to {
    transform: translateY(-20px);
    opacity: 0;
    max-height: 0;
    margin-bottom: 0;
}

.slide-fade-enter-active,
.slide-fade-leave-active {
    transition: all 0.5s ease-in-out;
    overflow: hidden;
}

.slide-fade-enter-to,
.slide-fade-leave-from {
    transform: translateY(0);
    opacity: 1;
    max-height: 600px;
}

.swap-enter-active,
.swap-leave-active {
    transition: all 0.3s ease;
}

.swap-enter-from {
    opacity: 0;
    transform: scale(0.98);
}

.swap-leave-to {
    opacity: 0;
    transform: scale(1.02);
}
</style>