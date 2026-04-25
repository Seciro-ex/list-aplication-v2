<script setup>
import { ref} from 'vue'

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
        message.value = "Server error: " + error.message
    }
}
</script>

<template>
    <div class="mt-2 p-3 bg-red-50 border border-red-200 rounded-lg shadow-inner">
        <p class="text-sm font-bold text-red-800 mb-2">Are you sure?</p>
        <div class="flex gap-2">
            <button @click="DeleteItem" 
                    class="bg-red-500 text-white px-3 py-1 rounded text-xs hover:bg-red-600">
                YES
            </button>
            <button @click="emit('cancel')" 
                    class="bg-gray-300 text-gray-700 px-3 py-1 rounded text-xs hover:bg-gray-400">
                NO
            </button>
        </div>
        <p v-if="message" class="text-xs text-red-600 mt-1">{{ message }}</p>
    </div>
</template>