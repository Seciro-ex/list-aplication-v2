<template>
    <div class="flex items-center gap-2">
        <span class="text-xs font-bold text-red-600">Delete this item?</span>
        <button @click="DeleteItem" 
                class="bg-red-500 text-white px-3 py-1 rounded-md text-xs hover:bg-red-600 transition shadow-sm">
            Yes
        </button>
        <button @click="emit('cancel')" 
                class="bg-gray-200 text-gray-700 px-3 py-1 rounded-md text-xs hover:bg-gray-300 transition">
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

        const url = `http://localhost:8081/delete-purchase?id=${props.itemId}`

        const response = await fetch(url, {
            method: "DELETE",
            headers: { 
                "Authorization": `Bearer ${token}`
            }
        })

        if (response.ok) {
            emit('item-deleted')
        } else {
            const errorText = await response.text()
            message.value = "Error: " + errorText
        }
    } catch (error) {
        message.value = "Server error"
    }
}
</script>