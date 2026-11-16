// Shape returned by the Go API (api/internal/httpapi/prescriptions.go).

export type PrescriptionRequestStatus = 'Pending' | 'Approved' | 'Ready' | 'Rejected';

/** Statuses that block a new request for the same prescription (enforced by the database). */
export const openRequestStatuses: PrescriptionRequestStatus[] = ['Pending', 'Approved'];

export interface PrescriptionRequest {
    id: string;
    status: PrescriptionRequestStatus;
    comments: string;
    rejection_reason: string;
    created_at: string;   // ISO 8601 timestamp, e.g. "2026-11-16T10:50:47Z"
    prescription_id: string;
    medicine_name: string;
    dosage: string;
}
