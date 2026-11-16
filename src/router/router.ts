import {createRouter, createWebHistory, type RouteLocationNormalized, type RouteRecordRaw} from 'vue-router'
import {useSession} from "@/stores/session";

// Route meta flags, read by the guard below:
//   public:     reachable without signing in (and pointless once signed in)
//   clinicOnly: clinic staff only
declare module 'vue-router' {
    interface RouteMeta {
        public?: boolean
        clinicOnly?: boolean
    }
}

// Pages are loaded on demand (`() => import(...)`), so the first visit only
// downloads the code for the page being opened.
const routes: RouteRecordRaw[] = [
    {path: '/', name: 'DashboardPetOwner', component: () => import('@/views/DashboardPetOwner.vue')},
    {
        path: '/clinic',
        name: 'DashboardClinicTeam',
        component: () => import('@/views/DashboardClinicTeam.vue'),
        meta: {clinicOnly: true},
    },
    {path: '/login', name: 'UserLogin', component: () => import('@/views/UserViews/UserLogin.vue'), meta: {public: true}},
    {
        path: '/register',
        name: 'UserRegistration',
        component: () => import('@/views/UserViews/UserRegistration.vue'),
        meta: {public: true},
    },
    {
        path: '/pets',
        children: [
            {path: 'add', name: 'PetAdd', component: () => import('@/views/PetDetailsView/PetDetailsView.vue')},
            {path: ':id', name: 'PetDetails', component: () => import('@/views/PetDetailsView/PetDetailsView.vue')},
            {
                path: ':id/prescriptions',
                name: 'Prescriptions',
                component: () => import('@/views/PrescriptionsView/PrescriptionsView.vue'),
            },
        ],
    },
    {
        path: '/user',
        children: [
            {path: 'settings', name: 'Settings', component: () => import('@/views/UserViews/UserSettings.vue')},
            {path: 'personal-details', name: 'UserDetails', component: () => import('@/views/UserViews/UserDetails.vue')},
            {path: 'change-password', name: 'EditPassword', component: () => import('@/views/UserViews/EditPassword.vue')},
        ],
    },
    // Must stay last: anything that didn't match above.
    {path: '/:pathMatch(.*)*', name: 'NotFound', component: () => import('@/views/NotFound.vue')},
]

export const router = createRouter({
    history: createWebHistory(),
    routes,
    // Start each new page at the top, but restore the position on back/forward.
    scrollBehavior(_to, _from, savedPosition) {
        return savedPosition ?? {top: 0}
    },
})

/** Only same-site paths, so ?redirect= can't send anyone to another website. */
export function safeRedirect(value: unknown): string {
    return typeof value === 'string' && value.startsWith('/') && !value.startsWith('//') ? value : '/';
}

/** The sign-in page, remembering where to come back to afterwards. */
export function loginRoute(returnTo: RouteLocationNormalized | string) {
    const path = typeof returnTo === 'string' ? returnTo : returnTo.fullPath;
    return {name: 'UserLogin', query: path === '/' ? {} : {redirect: path}};
}

// Runs before every navigation, like an auth middleware: works out who is
// signed in (once), then lets the navigation through or redirects it.
router.beforeEach(async (to) => {
    const session = useSession();
    try {
        await session.ensureLoaded();
    } catch {
        // The API is unreachable. Let the page load; its own requests will show the error.
        return true;
    }
    const user = session.user.value;

    if (!user) {
        return to.meta.public ? true : loginRoute(to);
    }
    if (to.meta.public) {
        return '/';
    }
    if (to.meta.clinicOnly && !user.veterinary) {
        return '/';
    }
    // Staff working in the clinic view start on the clinic dashboard.
    if (to.name === 'DashboardPetOwner' && user.app_mode === 'clinicTeam') {
        return '/clinic';
    }
    return true;
})
