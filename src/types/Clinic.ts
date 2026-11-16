// Shapes returned by the clinic endpoints (api/internal/httpapi/clinic.go).

import type {Pet} from "@/types/Pet";
import type {Medicine} from "@/types/Prescription";
import type {PrescriptionRequestStatus} from "@/types/PrescriptionRequest";

export interface Owner {
    id: string;
    name: string;
    surname: string;
    email: string;
}

export interface ClinicSummary {
    pending_requests: number;
    approved_requests: number;
    ready_requests: number;
    patients: number;
    clients: number;
}

export interface ClinicRequest {
    id: string;
    status: PrescriptionRequestStatus;
    comments: string;
    rejection_reason: string;
    created_at: string;   // ISO 8601 timestamp
    prescription: {
        id: string;
        dosage: string;
        quantity: number;
        medicine_name: string;
    };
    pet: {
        id: string;
        name: string;
        species_name: string;
    };
    owner: Owner;
}

/** A pet as staff see it: the owner's view of the pet, plus who owns it. */
export interface ClinicPet extends Pet {
    owner: Owner;
    /** Requests waiting on the clinic (Pending or Approved). */
    open_requests: number;
}

export interface ClinicPrescription {
    id: string;
    dosage: string;
    quantity: number;
    created_at: string;
    medicine: Medicine;
    /** Name of the staff member who wrote it, if known. */
    prescribed_by: string | null;
}

/** Statuses staff can move a request to. */
export type RequestAction = 'Approved' | 'Ready' | 'Rejected';
