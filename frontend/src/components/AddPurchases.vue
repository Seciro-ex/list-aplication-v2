<script setup>
import { ref} from 'vue'

const emit = defineEmits(['added-item'])

const cost = ref(0)
const image_url = ref('')
const item_name = ref('')
const item_url = ref('')
const message = ref('')

async function validateBeforeAdd() {
    if (item_name.value.length <= 1) {
        message.value = "Name must have at least 2 symbols!"
        return
    }
    if (cost.value <= 0) {
        message.value = "Cost must be at least 1!"
        return
    }
    addItem()
}

async function addItem() {
    try {
        const token = localStorage.getItem('token')
            
        if (!token) {
            emit('logout')
            return
        }

        const response = await fetch("http://localhost:8081/api/add-purchase", {
            method: "POST",
            headers: { 
                "Content-Type": "application/json", 
                "Authorization": `Bearer ${token}`
            },
            body: JSON.stringify({
                cost: cost.value,
                image_url: image_url.value,
                item_name: item_name.value,
                item_url: item_url.value
            })
        })

        if (response.ok) {
            item_name.value = ''
            cost.value = 0
            item_url.value = ''
            image_url.value = ''
            message.value = "Item added succsesful!"
            emit('added-item')
        } else {
            const data = await response.json()
            message.value = "Ошибка: " + (data.error || "Не удалось добавить")
        }
    } catch(error) {
        message.value = "Error: bad data " + error
    }
}
</script>

<template>
    <div class="bg-green-50 p-5 rounded-lg mb-6 border-2 border-dashed border-green-300 shadow-inner">
        <h3 class="font-bold text-green-800 mb-4">Add new item</h3>
        
        <div class="space-y-3">
            <input v-model="item_name" type="text" placeholder="Item name (e.g. RTX 4090)" 
                   class="w-full p-2 border border-green-200 rounded outline-none focus:ring-2 focus:ring-green-400" />
            
            <input v-model.number="cost" type="number" placeholder="Cost (₽)" 
                   class="w-full p-2 border border-green-200 rounded outline-none focus:ring-2 focus:ring-green-400" />
            
            <input v-model="item_url" type="text" placeholder="Shop link (optional)" 
                   class="w-full p-2 border border-green-200 rounded outline-none focus:ring-2 focus:ring-green-400" />
            
            <input v-model="image_url" type="text" placeholder="Image link (optional)" 
                   class="w-full p-2 border border-green-200 rounded outline-none focus:ring-2 focus:ring-green-400" />
            
            <button @click="validateBeforeAdd" 
                    class="w-full bg-green-500 text-white font-bold py-2 rounded shadow hover:bg-green-600 transition">
                Save to List
            </button>
        </div>

        <p v-if="message" class="text-red-500 text-sm mt-3 text-center">{{ message }}</p>
    </div>
</template>