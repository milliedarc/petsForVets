interface Prescription {
    id: string | undefined;
    dosage: string;
    quantity: number;
    expand: Expand;
}

interface Medicine {
    id: string;
    name: string;
    description: string;
}

interface Expand {
    pet: Pet;
    medicine: Medicine;
}