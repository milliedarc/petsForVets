import {api} from "@/api/client";
import type {Pet} from "@/types/Pet";

/** Body for creating or updating a pet (see petRequest in api/internal/httpapi/pets.go). */
export interface PetInput {
    name: string;
    species_id: string;
    breed_id: string | null;
    date_of_birth: string | null;   // "YYYY-MM-DD"
    gender: string;
    neutered: boolean;
    colour: string;
    imported: boolean;
    import_country_code: string;
    microchip_number: string;
}

/** The logged-in user's pets, with species and breed included. */
export function listPets(): Promise<Pet[]> {
    return api<Pet[]>('GET', '/pets');
}

/** One pet. Throws ApiError with status 404 if it doesn't exist or isn't yours. */
export function getPet(id: string): Promise<Pet> {
    return api<Pet>('GET', `/pets/${encodeURIComponent(id)}`);
}

export function createPet(input: PetInput): Promise<Pet> {
    return api<Pet>('POST', '/pets', input);
}

export function updatePet(id: string, input: PetInput): Promise<Pet> {
    return api<Pet>('PUT', `/pets/${encodeURIComponent(id)}`, input);
}

export function deletePet(id: string): Promise<void> {
    return api<void>('DELETE', `/pets/${encodeURIComponent(id)}`);
}
