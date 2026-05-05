<script setup>
import { ref } from 'vue'

const props = defineProps(['itemId', 'cost', 'image_url', 'item_name', 'item_url', 'item_description'])
const emit = defineEmits(['item-edited', 'cancel'])
const item_name = ref(props.item_name || '')
const cost = ref(props.cost || 0)
const item_url = ref(props.item_url || '')
const image_url = ref(props.image_url || '')
const item_description  = ref(props.item_description || '')

const message = ref('')

async function validateBeforeEdit() {
    if (item_name.value.length <= 1) {
        message.value = "Name must have at least 2 symbols!"
        return
    }
    if (cost.value <= 0) {
        message.value = "Cost must be at least 1!"
        return
    }
    EditItem()
}


async function EditItem() {
    try {
    const token = localStorage.getItem('token')
    if (!token) return

    const url = `http://localhost:8081/api/edit-purchase?id=${props.itemId}`

    const response = await fetch(url, {
        method: "POST",
        headers: { 
            "Content-Type": "application/json", 
            "Authorization": `Bearer ${token}`
        },
        body: JSON.stringify({
            cost: cost.value,
            image_url: image_url.value,
            item_name: item_name.value,
            item_url: item_url.value,
            item_description: item_description.value
        })
    })
    if (response.ok) {
        emit('item-edited')
    } else {
        const errorText = await response.text()
        message.value = "Error: " + errorText
    }
    } catch(error) {
        message.value = "Error: " + error
    }
}
</script>

<template>
    <div class="bg-green-50 p-5 rounded-lg mb-6 border-2 border-dashed border-orange-300 shadow-inner">
        <h3 class="font-bold text-orange-800 mb-4">Edit item</h3>
        
        <div class="space-y-3">
            <input v-model="item_name" type="text" placeholder="Item name (e.g. RTX 4090)"
                   class="w-full p-2 border border-orange-200 rounded outline-none focus:ring-2 focus:ring-orange-400" />
            
            <input v-model.number="cost" type="number" placeholder="Cost (₽)"
                   class="w-full p-2 border border-orange-200 rounded outline-none focus:ring-2 focus:ring-orange-400" />
            
            <input v-model="item_url" type="text" placeholder="Shop link (optional)"
                   class="w-full p-2 border border-orange-200 rounded outline-none focus:ring-2 focus:ring-orange-400" />
            
            <input v-model="image_url" type="text" placeholder="Image link (optional)"
                   class="w-full p-2 border border-orange-200 rounded outline-none focus:ring-2 focus:ring-orange-400" />
            
            <input v-model="item_description" type="text" placeholder="Item description(optional)"
                   class="w-full p-2 border border-orange-200 rounded outline-none focus:ring-2 focus:ring-orange-400" />
            
            <button @click="validateBeforeEdit" 
                    class="w-full bg-orange-500 text-white font-bold py-2 rounded shadow hover:bg-orange-600 transition">
                Save to List
            </button>
            <button @click="emit('cancel')" 
                    class="w-full bg-red-500 text-white font-bold py-2 rounded shadow hover:bg-orange-600 transition">
                Cancel
            </button>
        </div>

        <p v-if="message" class="text-red-500 text-sm mt-3 text-center">{{ message }}</p>
    </div>
</template>