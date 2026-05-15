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

        const response = await fetch("/api/me", {
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
    <div class="max-w-md mx-auto mt-10 p-6 bg-white rounded-xl shadow-lg border border-slate-100 transition-all duration-300
                dark:bg-slate-800 dark:border-slate-700">

        <Transition name="swap" mode="out-in">

            <div v-if="UserEdit" key="edit">
                <EditUser :username="userInfo.username" :id="userInfo.id" :picture_url="userInfo.picture_url"
                    :email="userInfo.email" @user-edited="() => { handleProfilInfo(); UserEdit = false }"
                    @cancel="UserEdit = false" />
            </div>

            <div v-else-if="UserDelete" key="delete">
                <DeleteUser :id="userInfo.id" :username="userInfo.username" @user-deleted="emit('logout')"
                    @cancel="UserDelete = false" />
            </div>

            <div v-else-if="userInfo" key="info">
                <h2 class="text-2xl font-bold mb-6 text-slate-800 dark:text-white">User Profile</h2>

                <div class="flex flex-col items-center mb-6">
                    <div class="p-1 rounded-full border-2 border-dashed border-slate-200 dark:border-slate-600">
                        <img v-if="userInfo.picture_url" :src="userInfo.picture_url"
                            class="w-24 h-24 rounded-full object-cover shadow-md hover:scale-110 transition-transform duration-300"
                            alt="Avatar">
                        <div v-else
                            class="w-24 h-24 rounded-full bg-indigo-500 flex items-center justify-center text-white font-bold text-2xl uppercase">
                            {{ userInfo.username[0] }}
                        </div>
                    </div>
                </div>

                <div class="space-y-4 mb-6">
                    <div
                        class="bg-slate-50 dark:bg-slate-900/50 p-4 rounded-xl border border-transparent dark:border-slate-700/50">
                        <p
                            class="text-xs font-semibold uppercase tracking-wider text-slate-400 dark:text-slate-500 mb-1">
                            Username</p>
                        <p class="text-lg font-medium text-slate-800 dark:text-white break-all">{{ userInfo.username }}
                        </p>
                    </div>

                    <div
                        class="bg-slate-50 dark:bg-slate-900/50 p-4 rounded-xl border border-transparent dark:border-slate-700/50">
                        <p
                            class="text-xs font-semibold uppercase tracking-wider text-slate-400 dark:text-slate-500 mb-1">
                            Email</p>
                        <p class="text-lg font-medium text-slate-800 dark:text-white">{{ userInfo.email }}</p>
                    </div>

                    <div
                        class="bg-slate-50 dark:bg-slate-900/50 p-4 rounded-xl border border-transparent dark:border-slate-700/50">
                        <p
                            class="text-xs font-semibold uppercase tracking-wider text-slate-400 dark:text-slate-500 mb-1">
                            Created at</p>
                        <p class="text-sm font-mono text-slate-600 dark:text-slate-300">{{ formattedDate }}</p>
                    </div>
                </div>

                <div class="flex gap-3 mb-3">
                    <button @click="UserEdit = true" class="flex-1 py-2 font-bold rounded-lg transition-all active:scale-95
                                   bg-amber-100 text-amber-700 hover:bg-amber-200 
                                   dark:bg-amber-500/20 dark:text-amber-400 dark:hover:bg-amber-500/30">
                        Edit
                    </button>
                    <button @click="UserDelete = true" class="flex-1 py-2 font-bold rounded-lg transition-all active:scale-95
                                   bg-rose-100 text-rose-700 hover:bg-rose-200 
                                   dark:bg-rose-500/20 dark:text-rose-400 dark:hover:bg-rose-500/30">
                        Delete
                    </button>
                </div>

                <button @click="emit('go-to-main')" class="w-full py-2 font-medium rounded-lg transition-all
                               bg-slate-100 text-slate-600 hover:bg-slate-200
                               dark:bg-slate-700 dark:text-slate-300 dark:hover:bg-slate-600">
                    Back to List
                </button>
            </div>

            <div v-else class="text-center py-12" key="loading">
                <p class="text-slate-400 animate-pulse">Loading profile...</p>
            </div>
        </Transition>

        <p v-if="message" class="text-rose-500 dark:text-rose-400 text-center text-sm mt-4 font-medium">{{ message }}
        </p>
    </div>
</template>