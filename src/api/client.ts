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

export async function api<T>(method: Method, path: string, body?: unknown): Promise<T> {
    const response = await fetch(`/api${path}`, {
        method,
        headers: body === undefined ? {} : {'Content-Type': 'application/json'},
        body: body === undefined ? undefined : JSON.stringify(body),
    });

    if (!response.ok) {
        // The API sends {"error": "..."}; fall back to the status text otherwise.
        const data = await response.json().catch(() => null);
        throw new ApiError(response.status, data?.error ?? response.statusText);
    }

    if (response.status === 204) {
        return undefined as T;
    }
    return response.json() as Promise<T>;
}
