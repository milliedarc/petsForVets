import {api} from "@/api/client";
import type {User} from "@/types/User";

export function login(email: string, password: string): Promise<User> {
    return api<User>('POST', '/auth/login', {email, password});
}

export function register(email: string, password: string, name: string, surname: string): Promise<User> {
    return api<User>('POST', '/auth/register', {email, password, name, surname});
}

export function logout(): Promise<void> {
    return api<void>('POST', '/auth/logout');
}

/** The logged-in user. Throws ApiError with status 401 if not logged in. */
export function getCurrentUser(): Promise<User> {
    return api<User>('GET', '/auth/me');
}
