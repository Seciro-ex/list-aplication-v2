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
    if (response.ok) {
        emit('user-edited')
    } else {
        const errorText = await response.text()
        message.value = "Error: " + errorText
    }

    } catch (error) {
        message.value = error
    }
}

</script>

<template>
<div class="bg-orange-50 p-5 rounded-lg mb-6 border-2 border-dashed border-orange-300 shadow-inner">
        <h3 class="font-bold text-orange-800 mb-4">Edit Profile</h3>
        
        <div class="space-y-3">
            <input v-model="username" type="text" placeholder="Username"
                   class="w-full p-2 border border-orange-200 rounded outline-none focus:ring-2 focus:ring-orange-400" />
            
            <input v-model="picture_url" type="text" placeholder="Profile picture URL"
                   class="w-full p-2 border border-orange-200 rounded outline-none focus:ring-2 focus:ring-orange-400" />
            
            <input v-model="email" type="email" placeholder="Email"
                   class="w-full p-2 border border-orange-200 rounded outline-none focus:ring-2 focus:ring-orange-400" />
            
            <input v-model="old_password" type="password" placeholder="Old password (required to save)"
                    class="w-full p-2 border border-orange-200 rounded outline-none focus:ring-2 focus:ring-orange-400" />

            <input v-model="password" type="password" placeholder="New password (leave blank to keep current)"
                    class="w-full p-2 border border-orange-200 rounded outline-none focus:ring-2 focus:ring-orange-400" />

            <button @click="ValidateInfo" 
                    class="w-full bg-orange-500 text-white font-bold py-2 rounded shadow hover:bg-orange-600 transition">
                Save Profile
            </button>
            <button @click="emit('cancel')" 
                    class="w-full bg-red-500 text-white font-bold py-2 rounded shadow hover:bg-red-600 transition">
                Cancel
            </button>
        </div>

        <p v-if="message" class="text-red-500 text-sm mt-3 text-center">{{ message }}</p>
    </div>
</template>