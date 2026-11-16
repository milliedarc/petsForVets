// ---------------------------------------------------------------------------
// Shapes returned by the Go API (api/internal/httpapi/pets.go).
// ---------------------------------------------------------------------------

export interface SpeciesRef {
    id: string;
    name: string;
    pet_type_id: string;
}

export interface BreedRef {
    id: string;
    name: string;
}

export interface Pet {
    id: string;
    name: string;
    date_of_birth: string | null;   // "YYYY-MM-DD"
    gender: string;
    neutered: boolean;
    colour: string;
    imported: boolean;
    import_country_code: string;
    microchip_number: string;
    avatar_url: string | null;
    species: SpeciesRef;
    breed: BreedRef | null;
}

/** A row of GET /api/species. */
export interface Species {
    id: string;
    name: string;
    pet_type_id: string;
}

/** A row of GET /api/breeds. */
export interface Breed {
    id: string;
    name: string;
    species_id: string;
    pet_type_id: string;
    generic: boolean;
}
