import {createRouter, createWebHistory} from 'vue-router'

import UserLogin from '@/views/UserViews/UserLogin.vue'
import DashboardPetOwner from '@/views/DashboardPetOwner.vue'
import PetDetailsView from "@/views/PetDetailsView/PetDetailsView.vue";
import PrescriptionsView from "@/views/PrescriptionsView/PrescriptionsView.vue";
import UserSettings from "@/views/UserViews/UserSettings.vue";
import UserDetails from "@/views/UserViews/UserDetails.vue";
import EditPassword from "@/views/UserViews/EditPassword.vue";
import UserRegistration from "@/views/UserViews/UserRegistration.vue";
import DashboardClinicTeam from "@/views/DashboardClinicTeam.vue";
import NotFound from "@/views/NotFound.vue";

// TODO: optimise routes by nesting ?

const routes = [
    {path: '/', component: DashboardPetOwner, name: 'DashboardPetOwner'},
    {path: '/clinic', component: DashboardClinicTeam, name: 'DashboardClinicTeam'},
    {path: '/login', component: UserLogin, name: 'UserLogin'},
    {path: '/register', component: UserRegistration, name: 'UserRegistration'},
    {path: '/pets/:id', component: PetDetailsView, name: 'PetDetails'},
    {path: '/pets/:id/prescriptions', component: PrescriptionsView, name: 'Prescriptions'},
    {path: '/user/settings', component: UserSettings, name: 'Settings'},
    {path: '/user/personal-details', component: UserDetails, name: 'UserDetails'},
    {path: '/user/change-password', component: EditPassword, name: 'EditPassword'},
    // Must stay last: anything that didn't match above.
    {path: '/:pathMatch(.*)*', component: NotFound, name: 'NotFound'},
]

export const router = createRouter({
    history: createWebHistory(),
    routes,
    // Start each new page at the top, but restore the position on back/forward.
    scrollBehavior(_to, _from, savedPosition) {
        return savedPosition ?? {top: 0}
    },
})