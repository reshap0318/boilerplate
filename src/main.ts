import { createApp } from 'vue'
import './assets/styles/main.css'
import App from './App.vue'
import router from './router'
import pinia from './stores'
import vClickOutside from './components/directives/v-click-outside'

const app = createApp(App)

app.use(pinia)
app.use(router)
app.directive('click-outside', vClickOutside)

app.mount('#app')
