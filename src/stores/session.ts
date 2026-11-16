// The signed-in user, shared by every component.
//
// This is a plain module with reactive state rather than a library like Pinia:
// a module is only evaluated once, so every `import {useSession}` gets the same
// `user` ref. (Roughly what $_SESSION['user'] is in PHP, but held in the
// browser's memory and refreshed from GET /api/auth/me.)

import {computed, readonly, ref} from "vue";
import type {User} from "@/types/User";
import {ApiError} from "@/api/client";
import * as authApi from "@/api/auth";

const user = ref<User | null>(null);

/** Resolves once we've asked the API who is signed in, so it's only asked once per page load. */
let loaded: Promise<void> | null = null;

/** Asks the API who is signed in (once). A 401 just means nobody is. */
function ensureLoaded(): Promise<void> {
    loaded ??= authApi.getCurrentUser()
        .then((currentUser) => {
            user.value = currentUser;
        })
        .catch((error) => {
            user.value = null;
            if (!(error instanceof ApiError && error.status === 401)) {
                // Server down or similar: allow a retry on the next navigation.
                loaded = null;
                throw error;
            }
        });
    return loaded;
}

async function login(email: string, password: string): Promise<User> {
    user.value = await authApi.login(email, password);
    return user.value;
}

async function register(email: string, password: string, name: string, surname: string): Promise<User> {
    user.value = await authApi.register(email, password, name, surname);
    return user.value;
}

async function logout(): Promise<void> {
    try {
        await authApi.logout();
    } finally {
        user.value = null;
    }
}

async function switchAppMode(appMode: 'petOwner' | 'clinicTeam'): Promise<User> {
    user.value = await authApi.updateAppMode(appMode);
    return user.value;
}

/** For when the server has ended the session (expired, or password changed). */
function clear(): void {
    user.value = null;
}

export function useSession() {
    return {
        user: readonly(user),
        isLoggedIn: computed(() => user.value !== null),
        ensureLoaded,
        login,
        register,
        logout,
        switchAppMode,
        clear,
    };
}
