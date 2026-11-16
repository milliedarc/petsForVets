import {createApp} from 'vue'
import App from './App.vue'
import PrimeVue from 'primevue/config';
import 'primeicons/primeicons.css'
import './assets/main.css'
import ToastService from 'primevue/toastservice';
import ConfirmationService from 'primevue/confirmationservice';

import {loginRoute, router} from './router/router'
import {PetsForVetsTheme} from './theme'
import {setUnauthorizedHandler} from './api/client'
import {useSession} from './stores/session'

// If the session ends while the app is open (it expired, or the password was
// changed in another tab), go to the sign-in page and come back here afterwards.
setUnauthorizedHandler(() => {
    const session = useSession();
    if (!session.isLoggedIn.value) {
        return;
    }
    session.clear();
    router.push(loginRoute(router.currentRoute.value));
});

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
