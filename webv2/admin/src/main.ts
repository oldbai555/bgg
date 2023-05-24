import {createApp} from 'vue'
import App from './App.vue'
import router from './router'
import './assets/css/style.css'
import './assets/css/overflow_y.css'

import Antd from 'ant-design-vue';
import 'ant-design-vue/dist/antd.css';

const app = createApp(App)

app.use(Antd)
app.use(router)

app.mount('#app')
