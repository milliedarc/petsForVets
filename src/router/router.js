import {createRouter, createWebHistory} from 'vue-router'

import UserLogin from '@/views/UserLogin.vue'
import Dashboard from '@/views/Dashboard.vue'
import ViewPet from '@/components/ViewPet.vue'
import PetDetailsView from "@/views/PetDetailsView.vue";

const routes = [
    {path: '/', component: Dashboard, name: 'Dashboard'},
    {path: '/login', component: UserLogin, name: 'UserLogin'},
    // {path: '/pets/:id', component: ViewPet, name: 'ViewPet'},
    {path: '/pets/:id', component: PetDetailsView, name: 'PetDetails'},
    //{ path: '/about', component: AboutView },
]

export const router = createRouter({
    history: createWebHistory(),
    routes,
})