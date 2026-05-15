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

        const data = await response.json()

        if (!response.ok) {
            message.value = "Error: " + (data.error || "Wrong password")
            if (data.details) {
                console.error("Tech details:", data.details)
            }
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

        const data = await response.json()

        if (response.ok) {
            message.value = data.message
            setTimeout(() => emit('user-deleted'), 2000)
        } else {
            message.value = data.error
        }
    } catch (error) {
        message.value = "Error: " + error
    }
}


</script>

<template>
    <div class="p-5 rounded-xl border-2 border-dashed transition-all duration-300 shadow-inner
                bg-rose-50 border-rose-200 
                dark:bg-rose-950/20 dark:border-rose-500/30">
        <h3 class="font-bold mb-4 flex justify-center items-center gap-2 text-rose-800 dark:text-rose-400">DELETE USER</h3>
        <p class="mb-6 text-sm text-center leading-relaxed text-slate-600 dark:text-slate-300">
            <span class="font-bold text-rose-600 dark:text-rose-400">WARNING!</span> 
            This action is irreversible. It will permanently delete your profile and all associated purchases.
        </p>
        
        <div class="space-y-3">
        <input v-model="password" type="password" placeholder="Enter your password"
            class="w-full p-2 border rounded-lg outline-none transition-all
                       bg-white text-slate-900 border-rose-200 focus:ring-2 focus:ring-rose-500
                       dark:bg-slate-800 dark:text-white dark:border-slate-700 dark:focus:ring-rose-500 dark:placeholder-rose-900/30" />
        
        <button @click="VerifyUser" 
                class="w-full font-bold py-2 rounded-lg shadow-md transition-all active:scale-95
                               bg-rose-600 text-white hover:bg-rose-700 disabled:opacity-50 disabled:cursor-not-allowed
                               dark:bg-rose-500/20 dark:text-rose-400 dark:hover:bg-rose-500/30 " :disabled="!password">
            Confirm Destruction
        </button>
        <button @click="emit('cancel')" 
                class="w-full py-2 rounded-lg font-medium transition-all
                               bg-slate-200 text-slate-700 hover:bg-slate-300
                               dark:bg-slate-700 dark:text-slate-300 dark:hover:bg-slate-600">
            I changed my mind
        </button>
        </div>
        <p v-if="message" class="text-rose-500 dark:text-rose-400 text-sm mt-3 text-center font-medium">{{ message }}</p>
    </div>
</template>