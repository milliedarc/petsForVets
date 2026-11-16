import {createApp} from 'vue'
import App from './App.vue'
import PrimeVue from 'primevue/config';
import 'primeicons/primeicons.css'
import './assets/main.css'
import ToastService from 'primevue/toastservice';
import ConfirmationService from 'primevue/confirmationservice';

import {router} from './router/router'
import {PetsForVetsTheme} from './theme'

const app = createApp(App);

app.use(router);
app.use(PrimeVue, {
    theme: {
        preset: PetsForVetsTheme,
        options: {
            darkModeSelector: 'none',
        }
    }
});
app.use(ToastService);
app.use(ConfirmationService);

// Mount the app
app.mount('#app');
