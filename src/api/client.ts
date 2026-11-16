// Thin wrapper around fetch() for talking to the Go API.
// The session cookie is sent automatically because /api is same-origin
// (Vite proxies it to the Go server in development).

export class ApiError extends Error {
    constructor(public status: number, message: string) {
        super(message);
        this.name = 'ApiError';
    }
}

type Method = 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE';

/** Paths whose 401 means "wrong password" or "not signed in yet", not "your session ended". */
const authPaths = ['/auth/login', '/auth/register', '/auth/me'];

let onUnauthorized: (() => void) | null = null;

/**
 * Called when any other request gets a 401, i.e. the session expired while the
 * app was open. Set once in main.ts (to send the user back to sign in).
 */
export function setUnauthorizedHandler(handler: () => void): void {
    onUnauthorized = handler;
}

export async function api<T>(method: Method, path: string, body?: unknown): Promise<T> {
    const response = await fetch(`/api${path}`, {
        method,
        headers: body === undefined ? {} : {'Content-Type': 'application/json'},
        body: body === undefined ? undefined : JSON.stringify(body),
    });

    if (!response.ok) {
        // 502/503/504 come from the Vite proxy, not our API: the Go server
        // isn't running (`npm run dev` starts it with Vite).
        if (response.status >= 502 && response.status <= 504) {
            throw new ApiError(response.status, "Can't reach the server. Is the API running? (npm run dev starts it)");
        }
        if (response.status === 401 && !authPaths.includes(path)) {
            onUnauthorized?.();
        }
        // The API sends {"error": "..."}; fall back to the status text otherwise.
        const data = await response.json().catch(() => null);
        throw new ApiError(response.status, data?.error ?? response.statusText);
    }

    if (response.status === 204) {
        return undefined as T;
    }
    return response.json() as Promise<T>;
}

/** The API's error message if there is one (so server-side validation reaches the user), else the fallback. */
export function errorMessage(error: unknown, fallback: string): string {
    return error instanceof ApiError ? error.message : fallback;
}
