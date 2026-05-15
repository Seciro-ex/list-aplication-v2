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

    const data = await response.json()

    if (response.ok) {
        emit('item-edited')
    } else {
        message.value = data.error
        if (data.details) {
            console.error("Tech details:", data.details)
        }
    }
    } catch(error) {
        message.value = "Error: " + error
    }
}
</script>

<template>
    <div class="p-5 rounded-xl mb-6 border-2 border-dashed transition-all duration-300
                bg-amber-50 border-amber-200 
                dark:bg-amber-900/10 dark:border-amber-500/30">
        <h3 class="font-bold mb-4 text-amber-800 dark:text-amber-400">Edit item</h3>
        
        <div class="space-y-3">
            <input v-model="item_name" type="text" placeholder="Item name (e.g. RTX 4090)"
                   class="w-full p-2 border rounded-lg outline-none transition-all
                          bg-white text-slate-900 border-amber-200 focus:ring-2 focus:ring-amber-500
                          dark:bg-slate-800 dark:text-white dark:border-slate-700 dark:focus:ring-amber-500 dark:placeholder-amber-700/50" />
            
            <input v-model.number="cost" type="number" placeholder="Cost (€)"
                   class="w-full p-2 border rounded-lg outline-none transition-all
                          bg-white text-slate-900 border-amber-200 focus:ring-2 focus:ring-amber-500
                          dark:bg-slate-800 dark:text-white dark:border-slate-700 dark:focus:ring-amber-500 dark:placeholder-amber-700/50" />
            
            <input v-model="item_url" type="text" placeholder="Shop link (optional)"
                   class="w-full p-2 border rounded-lg outline-none transition-all
                          bg-white text-slate-900 border-amber-200 focus:ring-2 focus:ring-amber-500
                          dark:bg-slate-800 dark:text-white dark:border-slate-700 dark:focus:ring-amber-500 dark:placeholder-amber-700/50" />
            
            <input v-model="image_url" type="text" placeholder="Image link (optional)"
                   class="w-full p-2 border rounded-lg outline-none transition-all
                          bg-white text-slate-900 border-amber-200 focus:ring-2 focus:ring-amber-500
                          dark:bg-slate-800 dark:text-white dark:border-slate-700 dark:focus:ring-amber-500 dark:placeholder-amber-700/50" />
            
            <input v-model="item_description" type="text" placeholder="Item description(optional)"
                   class="w-full p-2 border rounded-lg outline-none transition-all
                          bg-white text-slate-900 border-amber-200 focus:ring-2 focus:ring-amber-500
                          dark:bg-slate-800 dark:text-white dark:border-slate-700 dark:focus:ring-amber-500 dark:placeholder-amber-700/50" />
            
            <div class="flex flex-col gap-2 pt-2">
                <button @click="validateBeforeEdit" 
                        class="w-full font-bold py-2 rounded-lg shadow-md transition-all active:scale-95
                               bg-amber-600 text-white hover:bg-amber-700
                               dark:bg-amber-500/20 dark:text-amber-400 dark:hover:bg-amber-500/30">
                    Save Changes
                </button>
                
                <button @click="emit('cancel')" 
                        class="w-full py-2 rounded-lg font-medium transition-all
                               bg-slate-200 text-slate-700 hover:bg-slate-300
                               dark:bg-slate-700 dark:text-slate-300 dark:hover:bg-slate-600">
                    Cancel
                </button>
            </div>
        </div>

        <p v-if="message" class="text-red-500 dark:text-rose-400 text-sm mt-3 text-center">{{ message }}</p>
    </div>
</template>