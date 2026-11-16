// Clinic staff endpoints. Non-staff get ApiError 403 from all of these.

import {api} from "@/api/client";
import type {Medicine} from "@/types/Prescription";
import type {PrescriptionRequestStatus} from "@/types/PrescriptionRequest";
import type {
    ClinicPet,
    ClinicPrescription,
    ClinicRequest,
    ClinicSummary,
    RequestAction,
} from "@/types/Clinic";

export function getClinicSummary(): Promise<ClinicSummary> {
    return api<ClinicSummary>('GET', '/clinic/summary');
}

/** Requests, oldest first. Both filters are optional. */
export function listClinicRequests(filters: {status?: PrescriptionRequestStatus, petId?: string} = {}): Promise<ClinicRequest[]> {
    const query = new URLSearchParams();
    if (filters.status) {
        query.set('status', filters.status);
    }
    if (filters.petId) {
        query.set('pet_id', filters.petId);
    }
    const suffix = query.size > 0 ? `?${query}` : '';
    return api<ClinicRequest[]>('GET', `/clinic/prescription-requests${suffix}`);
}

/**
 * Approve, mark ready or reject a request. Throws ApiError 409 if it has
 * already moved on (e.g. another staff member approved it first).
 */
export function updateClinicRequest(id: string, status: RequestAction, rejectionReason = ''): Promise<ClinicRequest> {
    return api<ClinicRequest>('PATCH', `/clinic/prescription-requests/${encodeURIComponent(id)}`, {
        status,
        rejection_reason: rejectionReason,
    });
}

/** Patients matching pet name, owner name or email, or microchip number. */
export function listClinicPets(search = ''): Promise<ClinicPet[]> {
    const suffix = search.trim() ? `?${new URLSearchParams({search: search.trim()})}` : '';
    return api<ClinicPet[]>('GET', `/clinic/pets${suffix}`);
}

export function getClinicPet(id: string): Promise<ClinicPet> {
    return api<ClinicPet>('GET', `/clinic/pets/${encodeURIComponent(id)}`);
}

export function listClinicPetPrescriptions(petId: string): Promise<ClinicPrescription[]> {
    return api<ClinicPrescription[]>('GET', `/clinic/pets/${encodeURIComponent(petId)}/prescriptions`);
}

export interface PrescriptionInput {
    medicine_id: string;
    dosage: string;
    quantity: number;
}

export function createPrescription(petId: string, input: PrescriptionInput): Promise<ClinicPrescription> {
    return api<ClinicPrescription>('POST', `/clinic/pets/${encodeURIComponent(petId)}/prescriptions`, input);
}

export function listMedicines(): Promise<Medicine[]> {
    return api<Medicine[]>('GET', '/clinic/medicines');
}
