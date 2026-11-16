import {api} from "@/api/client";
import type {Breed, Species} from "@/types/Pet";

/** All species, alphabetically. */
export function listSpecies(): Promise<Species[]> {
    return api<Species[]>('GET', '/species');
}

/** All breeds, generic ones ("Mixed breed") first, then alphabetically. */
export function listBreeds(): Promise<Breed[]> {
    return api<Breed[]>('GET', '/breeds');
}
