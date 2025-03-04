import { createApp } from 'vue'
import App from './App.vue'
import './index.css';
import './assets/tailwind.css'
import router from './routes/index.js'

console.log("✅ Main.js is loading...");

const app = createApp(App);
app.use(router)
app.mount('#app');
