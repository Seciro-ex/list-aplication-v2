<script setup>
import { ref} from 'vue'

const emit = defineEmits(['added-item'])

const cost = ref(0)
const image_url = ref('')
const item_name = ref('')
const item_url = ref('')
const item_description = ref('')
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
                item_url: item_url.value,
                item_description: item_description.value
            })
        })

        const data = await response.json()

        if (response.ok) {
            item_name.value = ''
            cost.value = 0
            item_url.value = ''
            image_url.value = ''
            item_description.value = ''
            message.value = data.message
            emit('added-item')
        } else {
            message.value = data.error
        }
    } catch(error) {
        message.value = "Error: bad data " + error
    }
}
</script>

<template>
    <div class="p-5 rounded-xl mb-6 border-2 border-dashed transition-all duration-300
                bg-green-50 border-green-300 
                dark:bg-emerald-950/20 dark:border-emerald-500/40">
        <h3 class="font-bold mb-4 text-green-800 dark:text-emerald-400">Add new item</h3>
        
        <div class="space-y-3">
            <input v-model="item_name" type="text" placeholder="Item name (e.g. RTX 4090)" 
                   class="w-full p-2 border rounded-lg outline-none transition-all
                          bg-white text-slate-900 border-green-200 focus:ring-2 focus:ring-green-500
                          dark:bg-slate-800 dark:text-white dark:border-emerald-900/50 dark:focus:ring-emerald-500" />
            
            <input v-model.number="cost" type="number" placeholder="Cost (₽)" 
                   class="w-full p-2 border rounded-lg outline-none transition-all
                          bg-white text-slate-900 border-green-200 focus:ring-2 focus:ring-green-500
                          dark:bg-slate-800 dark:text-white dark:border-emerald-900/50 dark:focus:ring-emerald-500" />
            
            <input v-model="item_url" type="text" placeholder="Shop link (optional)" 
                   class="w-full p-2 border rounded-lg outline-none transition-all
                          bg-white text-slate-900 border-green-200 focus:ring-2 focus:ring-green-500
                          dark:bg-slate-800 dark:text-white dark:border-emerald-900/50 dark:focus:ring-emerald-500" />
            
            <input v-model="image_url" type="text" placeholder="Image link (optional)" 
                   class="w-full p-2 border rounded-lg outline-none transition-all
                          bg-white text-slate-900 border-green-200 focus:ring-2 focus:ring-green-500
                          dark:bg-slate-800 dark:text-white dark:border-emerald-900/50 dark:focus:ring-emerald-500" />
            
            <input v-model="item_description" type="textbox" placeholder="Purchases description(optional)" 
                   class="w-full p-2 border rounded-lg outline-none transition-all
                          bg-white text-slate-900 border-green-200 focus:ring-2 focus:ring-green-500
                          dark:bg-slate-800 dark:text-white dark:border-emerald-900/50 dark:focus:ring-emerald-500" />
            
            <button @click="validateBeforeAdd" 
                    class="w-full font-bold py-2 rounded-lg shadow-md transition-all active:scale-95
                           bg-green-600 text-white hover:bg-green-700
                           dark:bg-emerald-600 dark:hover:bg-emerald-500 dark:text-white">
                Save to List
            </button>
        </div>

        <p v-if="message" class="text-red-500 text-sm mt-3 text-center">
           {{ message }}
        </p>
    </div>
</template>