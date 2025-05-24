import {createRouter, createWebHistory} from 'vue-router'

import UserLogin from '@/views/UserViews/UserLogin.vue'
import Dashboard from '@/views/Dashboard.vue'
import PetDetailsView from "@/views/PetDetailsView/PetDetailsView.vue";
import PrescriptionsView from "@/views/PrescriptionsView/PrescriptionsView.vue";
import UserSettings from "@/views/UserViews/UserSettings.vue";
import UserDetails from "@/views/UserViews/UserDetails.vue";
import PasswordChange from "@/views/UserViews/PasswordChange.vue";

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