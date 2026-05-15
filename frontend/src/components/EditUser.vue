<script setup>
import { ref } from 'vue';

const props = defineProps(['id', 'username', 'picture_url', 'email', 'password'])
const emit = defineEmits(['user-edited', 'cancel'])

const username = ref(props.username || '')
const picture_url = ref(props.picture_url || '')
const email = ref(props.email || '')

const password = ref(props.password || '')
const old_password = ref('')

const emailRegex = /^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$/
const message = ref('')

async function ValidateInfo() {
    message.value = ""

    if (username.value.length < 2) {
        message.value = "Error! Username must be at least 2 symbols long!"
        return
    }
    if (email.value.length < 2 || !emailRegex.test(email.value)) {
        message.value = "Error! Please enter a valid email address."
        return
    }
    if (password.value !== "" && password.value.length < 6) {
        message.value = "Error! New password must be at least 6 symbols long!"
        return
    }
    try {
        const response = await fetch("http://localhost:8081/api/login", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
                username: props.username, 
                password: old_password.value
            })
        })

        if (!response.ok) {
            const data = await response.json()
            message.value = "Error: " + (data.error || "Wrong old password")
            return
        }
        await EditUser()

    } catch (error) {
        message.value = "Wrong old password"
    }
}

async function EditUser() {
    try {
    const token = localStorage.getItem('token')
    if (!token) return

    const url = `http://localhost:8081/api/edit-user?id=${props.id}`

    const response = await fetch(url, {
        method: "POST",
        headers: { 
            "Content-Type": "application/json", 
            "Authorization": `Bearer ${token}`
        },
        body: JSON.stringify({
            username: username.value,
            picture_url: picture_url.value,
            email: email.value,
            password: password.value
        })
    })

    const data = await response.json()

    if (response.ok) {
        emit('user-edited')
    } else {
        message.value = data.error
        if (data.details) {
            console.error("Tech details:", data.details)
        }
    }

    } catch (error) {
        message.value = error
    }
}

</script>

<template>
<div class="bg-amber-50 border-amber-200 p-5 rounded-lg border-2 border-dashed dark:bg-amber-900/10 dark:border-amber-500/30 shadow-inner">
        <h3 class="font-bold text-orange-800 dark:text-amber-400 mb-4">Edit Profile</h3>
        
        <div class="space-y-3">
            <input v-model="username" type="text" placeholder="Username"
                   class="w-full p-2 border rounded-lg outline-none transition-all
                          bg-white text-slate-900 border-amber-200 focus:ring-2 focus:ring-amber-500
                          dark:bg-slate-800 dark:text-white dark:border-slate-700 dark:focus:ring-amber-500 dark:placeholder-amber-700/50" />
            
            <input v-model="picture_url" type="text" placeholder="Profile picture URL"
                   class="w-full p-2 border rounded-lg outline-none transition-all
                          bg-white text-slate-900 border-amber-200 focus:ring-2 focus:ring-amber-500
                          dark:bg-slate-800 dark:text-white dark:border-slate-700 dark:focus:ring-amber-500 dark:placeholder-amber-700/50" />
            
            <input v-model="email" type="email" placeholder="Email"
                   class="w-full p-2 border rounded-lg outline-none transition-all
                          bg-white text-slate-900 border-amber-200 focus:ring-2 focus:ring-amber-500
                          dark:bg-slate-800 dark:text-white dark:border-slate-700 dark:focus:ring-amber-500 dark:placeholder-amber-700/50" />
            
            <input v-model="old_password" type="password" placeholder="Old password (required to save)"
                    class="w-full p-2 border rounded-lg outline-none transition-all
                          bg-white text-slate-900 border-amber-200 focus:ring-2 focus:ring-amber-500
                          dark:bg-slate-800 dark:text-white dark:border-slate-700 dark:focus:ring-amber-500 dark:placeholder-amber-700/50" />

            <input v-model="password" type="password" placeholder="New password (leave blank to keep current)"
                    class="w-full p-2 border rounded-lg outline-none transition-all
                          bg-white text-slate-900 border-amber-200 focus:ring-2 focus:ring-amber-500
                          dark:bg-slate-800 dark:text-white dark:border-slate-700 dark:focus:ring-amber-500 dark:placeholder-amber-700/50" />

            <button @click="ValidateInfo" 
                    class="w-full font-bold py-2 rounded-lg shadow-md transition-all active:scale-95
                               bg-amber-600 text-white hover:bg-amber-700
                               dark:bg-amber-500/20 dark:text-amber-400 dark:hover:bg-amber-500/30">
                Save Profile
            </button>
            <button @click="emit('cancel')" 
                    class="w-full py-2 rounded-lg font-medium transition-all
                               bg-slate-200 text-slate-700 hover:bg-slate-300
                               dark:bg-slate-700 dark:text-slate-300 dark:hover:bg-slate-600">
                Cancel
            </button>
        </div>

        <p v-if="message" class="text-red-500 dark:text-rose-400 text-sm mt-3 text-center">{{ message }}</p>
    </div>
</template>