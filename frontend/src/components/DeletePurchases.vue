<template>
    <div class="flex items-center gap-2">
        <span class="text-xs font-bold text-red-600 dark:text-red-400">Delete this item?</span>
        <button @click="DeleteItem"
            class="bg-red-500 text-white px-3 py-1 rounded-md text-xs hover:bg-red-600 transition shadow-sm dark:bg-rose-500/20 dark:text-rose-400 dark:hover:bg-rose-500/30">
            Yes
        </button>
        <button @click="emit('cancel')"
            class="bg-gray-200 text-gray-700 dark:bg-slate-700 dark:text-slate-300 dark:hover:bg-slate-600 px-3 py-1 rounded-md text-xs hover:bg-gray-300 transition">
            No
        </button>
        <span v-if="message" class="text-xs text-red-500 ml-2">{{ message }}</span>
    </div>
</template>

<script setup>
import { ref } from 'vue'

const props = defineProps(['itemId'])
const emit = defineEmits(['item-deleted', 'cancel'])

const message = ref('')

async function DeleteItem() {
    try {
        const token = localStorage.getItem('token')
        if (!token) return

        const url = `/api/delete-purchase?id=${props.itemId}`

        const response = await fetch(url, {
            method: "DELETE",
            headers: {
                "Authorization": `Bearer ${token}`
            }
        })

        const data = await response.json()

        if (response.ok) {
            emit('item-deleted')
        } else {
            message.value = "Error: " + data.error
            if (data.details) {
                console.error("Tech details:", data.details)
            }
        }
    } catch (error) {
        message.value = "Server error"
    }
}
</script>