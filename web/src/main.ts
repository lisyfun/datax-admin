import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ArcoVue from '@arco-design/web-vue'
import ArcoVueIcon from '@arco-design/web-vue/es/icon'
import router from './router'
import App from './App.vue'
import permissionDirective from './directives/permission'

import '@arco-design/web-vue/dist/arco.css'
import './styles/index.css'
import 'default-passive-events'

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(ArcoVue)
app.use(ArcoVueIcon)
app.use(permissionDirective)

app.mount('#app')
