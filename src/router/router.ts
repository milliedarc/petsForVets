import {
    createRouter,
    createWebHistory,
    type RouteLocationNormalized,
    type RouteLocationNormalizedLoaded,
    type RouteLocationRaw,
    type RouteRecordRaw,
} from 'vue-router'
import {useSession} from "@/stores/session";
import {petName} from "@/stores/petNames";

/** One step of the breadcrumb trail; the last one (the current page) has no link. */
export interface Crumb {
    label: string
    to?: RouteLocationRaw
}

// Route meta, read by the guard below and by AppBreadcrumbs:
//   public:      reachable without signing in (and pointless once signed in)
//   clinicOnly:  clinic staff only
//   breadcrumbs: the trail shown in the navbar for this page
declare module 'vue-router' {
    interface RouteMeta {
        public?: boolean
        clinicOnly?: boolean
        breadcrumbs?: (route: RouteLocationNormalizedLoaded) => Crumb[]
    }
}

const myPets: Crumb = {label: 'My pets', to: '/'};
const account: Crumb = {label: 'Account'};
const clinic: Crumb = {label: 'Clinic', to: {name: 'DashboardClinicTeam'}};
const patients: Crumb = {label: 'Patients', to: {name: 'ClinicPatients'}};

/** The pet's name once a page has loaded it, else a neutral placeholder. */
function petCrumb(route: RouteLocationNormalizedLoaded, linked: boolean): Crumb {
    const id = route.params.id as string;
    return {label: petName(id) ?? 'Pet', to: linked ? {name: 'PetDetails', params: {id}} : undefined};
}

// Pages are loaded on demand (`() => import(...)`), so the first visit only
// downloads the code for the page being opened.
const routes: RouteRecordRaw[] = [
    {
        path: '/',
        name: 'DashboardPetOwner',
        component: () => import('@/views/DashboardPetOwner.vue'),
        meta: {breadcrumbs: () => [{label: 'My pets'}]},
    },
    {
        // Child routes inherit meta, so every clinic page is staff-only.
        path: '/clinic',
        meta: {clinicOnly: true},
        children: [
            {
                path: '',
                name: 'DashboardClinicTeam',
                component: () => import('@/views/Clinic/ClinicDashboard.vue'),
                meta: {breadcrumbs: () => [{label: 'Clinic dashboard'}]},
            },
            {
                path: 'patients',
                name: 'ClinicPatients',
                component: () => import('@/views/Clinic/ClinicPatients.vue'),
                meta: {breadcrumbs: () => [clinic, {label: 'Patients'}]},
            },
            {
                path: 'patients/:id',
                name: 'ClinicPatient',
                component: () => import('@/views/Clinic/ClinicPatientView.vue'),
                meta: {
                    breadcrumbs: (route) => [clinic, patients, {label: petName(route.params.id as string) ?? 'Patient'}],
                },
            },
        ],
    },
    {
        path: '/login',
        name: 'UserLogin',
        component: () => import('@/views/UserViews/UserLogin.vue'),
        meta: {public: true, breadcrumbs: () => [{label: 'Sign in'}]},
    },
    {
        path: '/register',
        name: 'UserRegistration',
        component: () => import('@/views/UserViews/UserRegistration.vue'),
        meta: {public: true, breadcrumbs: () => [{label: 'Create an account'}]},
    },
    {
        path: '/pets',
        children: [
            {
                path: 'add',
                name: 'PetAdd',
                component: () => import('@/views/PetDetailsView/PetDetailsView.vue'),
                meta: {breadcrumbs: () => [myPets, {label: 'Add a pet'}]},
            },
            {
                path: ':id',
                name: 'PetDetails',
                component: () => import('@/views/PetDetailsView/PetDetailsView.vue'),
                meta: {breadcrumbs: (route) => [myPets, petCrumb(route, false)]},
            },
            {
                path: ':id/prescriptions',
                name: 'Prescriptions',
                component: () => import('@/views/PrescriptionsView/PrescriptionsView.vue'),
                meta: {breadcrumbs: (route) => [myPets, petCrumb(route, true), {label: 'Prescriptions'}]},
            },
        ],
    },
    {
        path: '/user',
        children: [
            {
                path: 'settings',
                name: 'Settings',
                component: () => import('@/views/UserViews/UserSettings.vue'),
                meta: {breadcrumbs: () => [account, {label: 'Settings'}]},
            },
            {
                path: 'personal-details',
                name: 'UserDetails',
                component: () => import('@/views/UserViews/UserDetails.vue'),
                meta: {breadcrumbs: () => [account, {label: 'Personal details'}]},
            },
            {
                path: 'change-password',
                name: 'EditPassword',
                component: () => import('@/views/UserViews/EditPassword.vue'),
                meta: {breadcrumbs: () => [account, {label: 'Settings', to: {name: 'Settings'}}, {label: 'Change password'}]},
            },
        ],
    },
    // Must stay last: anything that didn't match above.
    {
        path: '/:pathMatch(.*)*',
        name: 'NotFound',
        component: () => import('@/views/NotFound.vue'),
        meta: {breadcrumbs: () => [{label: 'Page not found'}]},
    },
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
