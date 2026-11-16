import {fileURLToPath, URL} from 'node:url'

import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'
import Components from 'unplugin-vue-components/vite';
import {PrimeVueResolver} from '@primevue/auto-import-resolver';

// https://vite.dev/config/
export default defineConfig({
    plugins: [
        vue(),
        vueDevTools(),
        Components({
            resolvers: [
                PrimeVueResolver()
            ]
        })
    ],
    resolve: {
        alias: {
            '@': fileURLToPath(new URL('./src', import.meta.url))
        },
    },
    // Pre-bundle the theme package (imported by src/theme.ts) at startup rather
    // than discovering it mid-session, which forces a reload of the dev server.
    optimizeDeps: {
        include: ['@primeuix/themes'],
    },
    server: {
        // Forward /api/* to the Go server, so the browser sees one origin
        // and the session cookie is sent automatically.
        proxy: {
            '/api': 'http://localhost:8080',
        },
    },
})
