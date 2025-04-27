import {createApp} from 'vue'
import App from './App.vue'
import 'bootstrap/dist/css/bootstrap.min.css'
import * as bootstrap from 'bootstrap/dist/js/bootstrap.bundle'
import PrimeVue from 'primevue/config';
import Aura from '@primeuix/themes/aura';
import 'primeicons/primeicons.css'
import ToastService from 'primevue/toastservice';

/* import the fontawesome core */
import {library} from '@fortawesome/fontawesome-svg-core'

/* import font awesome icon component */
import {FontAwesomeIcon, FontAwesomeLayers, FontAwesomeLayersText} from '@fortawesome/vue-fontawesome'

/* import specific icons */
import {fadr} from '@fortawesome/duotone-regular-svg-icons'
import {fadl} from '@fortawesome/duotone-light-svg-icons'
import {fas} from '@fortawesome/pro-solid-svg-icons'

/* add icons to the library */
library.add(fadr, fadl, fas)


import {router} from './router/router.js'

const app = createApp(App);

app.provide('bootstrap', bootstrap);

app.component('font-awesome-icon', FontAwesomeIcon)
app.component('font-awesome-layers', FontAwesomeLayers)
app.component('font-awesome-layer-text', FontAwesomeLayersText)

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
// Mount the app
app.mount('#app');