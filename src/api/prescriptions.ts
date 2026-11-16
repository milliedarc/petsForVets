import {api} from "@/api/client";
import type {Prescription} from "@/types/Prescription";
import type {PrescriptionRequest} from "@/types/PrescriptionRequest";

export function listPrescriptions(petId: string): Promise<Prescription[]> {
    return api<Prescription[]>('GET', `/pets/${encodeURIComponent(petId)}/prescriptions`);
}

export function listPrescriptionRequests(petId: string): Promise<PrescriptionRequest[]> {
    return api<PrescriptionRequest[]>('GET', `/pets/${encodeURIComponent(petId)}/prescription-requests`);
}

/** Ask for a repeat. Throws ApiError 409 if there's already an open request. */
export function requestRepeat(prescriptionId: string): Promise<PrescriptionRequest> {
    return api<PrescriptionRequest>('POST', `/prescriptions/${encodeURIComponent(prescriptionId)}/requests`);
}

/** Cancel a request. Only pending requests can be cancelled (ApiError 409 otherwise). */
export function cancelRequest(requestId: string): Promise<void> {
    return api<void>('DELETE', `/prescription-requests/${encodeURIComponent(requestId)}`);
}
