interface Breed {
    id: string;
    name: string;
    species: string;
    generic: boolean;
    pet_type: string;
}

interface Species {
    id: string;
    name: string;
    pet_type: string;
}

interface Expand {
    breed: Breed;
    species: Species;
}

interface Pet {
    id: string | undefined;
    name: string;
    species: string;
    breed_group: string;
    breed: string;
    date_of_birth: string;
    gender: string;
    neutered: boolean;
    colour: string;
    imported: boolean;
    import_country_code: string;
    microchip_number: string;
    avatar: string;
    expand: Expand;
}
