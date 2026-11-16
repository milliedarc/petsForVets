import type {Pet} from "@/types/Pet";

export interface Medicine {
    id: string;
    name: string;
    description: string;
}

export interface Prescription {
    id: string | undefined;
    dosage: string;
    quantity: number;
    expand: {
        pet: Pet;
        medicine: Medicine;
    };
}
