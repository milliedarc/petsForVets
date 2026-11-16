// Shapes returned by the Go API (api/internal/httpapi/prescriptions.go).

export interface Medicine {
    id: string;
    name: string;
    description: string;
}

export interface Prescription {
    id: string;
    dosage: string;
    quantity: number;
    medicine: Medicine;
}
