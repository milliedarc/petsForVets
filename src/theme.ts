// The app's PrimeVue theme: Aura with the Pets for Vets plum as the primary
// colour. Every PrimeVue component (buttons, focus rings, selected states, ...)
// picks it up, and it is also available in CSS as var(--p-primary-color).
import {definePreset} from '@primeuix/themes';
import Aura from '@primeuix/themes/aura';

export const PetsForVetsTheme = definePreset(Aura, {
    semantic: {
        primary: {
            50: '#fbf4f7',
            100: '#f5e4ec',
            200: '#ebc9d8',
            300: '#d9a1bb',
            400: '#c07596',
            500: '#9c4f72',
            600: '#843f61',
            700: '#6d3350',
            800: '#572940',
            900: '#421f31',
            950: '#2b131f',
        },
    },
});
