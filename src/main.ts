import {createApp} from 'vue'
import App from './App.vue'
import 'bootstrap/dist/css/bootstrap.min.css'
import * as bootstrap from 'bootstrap/dist/js/bootstrap.bundle'
import PrimeVue from 'primevue/config';
import Aura from '@primeuix/themes/aura';
import 'primeicons/primeicons.css'
import ToastService from 'primevue/toastservice';
import ConfirmationService from 'primevue/confirmationservice';

import {router} from './router/router'

const app = createApp(App);

app.provide('bootstrap', bootstrap);

app.use(router);
app.use(PrimeVue, {
    theme: {
        preset: Aura,
        options: {
            darkModeSelector: 'none',
        }
    }
});
app.use(ToastService);
app.use(ConfirmationService);

// Mount the app
app.mount('#app');