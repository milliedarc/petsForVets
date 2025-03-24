import {createApp} from 'vue'
import App from './App.vue'
import 'bootstrap/dist/css/bootstrap.min.css'
import * as bootstrap from 'bootstrap/dist/js/bootstrap.bundle'

const app = createApp(App);

// Make bootstrap library available inside the Vue app
app.provide('bootstrap', bootstrap);

// Mount the app
app.mount('#app');