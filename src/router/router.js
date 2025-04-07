import {createRouter, createWebHistory} from 'vue-router'

import UserLogin from '@/components/UserLogin.vue'
import Dashboard from '@/components/Dashboard.vue'
import ViewPet from '@/components/ViewPet.vue'

const routes = [
    {path: '/', component: Dashboard, name: 'Dashboard'},
    {path: '/login', component: UserLogin, name: 'UserLogin'},
    {path: '/pets/:id', component: ViewPet, name: 'ViewPet'},
    //{ path: '/about', component: AboutView },
]

export const router = createRouter({
    history: createWebHistory(),
    routes,
})