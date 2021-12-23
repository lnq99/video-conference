// register vue composition api globally
import {createApp} from 'vue'
import {createRouter, createWebHistory} from 'vue-router'
import {createPinia} from 'pinia'
import routes from 'virtual:generated-pages'
import App from './App.vue'

import Chat from 'vue3-beautiful-chat'

import '@unocss/reset/tailwind.css'
import './styles/main.css'
import 'uno.css'

const app = createApp(App)
const router = createRouter({
    history: createWebHistory(),
    routes,
})

app.use(Chat)
app.use(router)
app.use(createPinia())
app.mount('#app')
