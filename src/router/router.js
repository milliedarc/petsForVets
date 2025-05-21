import {createRouter, createWebHistory} from 'vue-router'

import UserLogin from '@/views/UserLogin.vue'
import Dashboard from '@/views/Dashboard.vue'
import PetDetailsView from "@/views/PetDetailsView.vue";
import PrescriptionsView from "@/views/PrescriptionsView.vue";
import UserSettings from "@/views/UserSettings/UserSettings.vue";
import UserDetails from "@/views/UserSettings/UserDetails.vue";
import PasswordChange from "@/views/UserSettings/PasswordChange.vue";

const routes = [
    {path: '/', component: Dashboard, name: 'Dashboard'},
    {path: '/login', component: UserLogin, name: 'UserLogin'},
    {path: '/pets/:id', component: PetDetailsView, name: 'PetDetails'},
    {path: '/pets/:id/prescriptions', component: PrescriptionsView, name: 'Prescriptions'},
    {path: '/user/settings', component: UserSettings, name: 'Settings'},
    {path: '/user/personal-details', component: UserDetails, name: 'UserDetails'},
    {path: '/user/change-password', component: PasswordChange, name: 'PasswordChange'}
    //{ path: '/about', component: AboutView },
]

export const router = createRouter({
    history: createWebHistory(),
    routes,
})