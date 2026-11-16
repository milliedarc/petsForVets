// How each request status is labelled for clinic staff. (Owners see their
// own wording in PrescriptionRequestCard.vue.)

import type {PrescriptionRequestStatus} from "@/types/PrescriptionRequest";

export const requestStatusInfo: Record<PrescriptionRequestStatus, {
    label: string
    severity: 'warn' | 'info' | 'success' | 'danger'
    icon: string
}> = {
    Pending: {label: 'To review', severity: 'warn', icon: 'pi pi-inbox'},
    Approved: {label: 'Preparing', severity: 'info', icon: 'pi pi-check'},
    Ready: {label: 'Ready to collect', severity: 'success', icon: 'pi pi-box'},
    Rejected: {label: 'Rejected', severity: 'danger', icon: 'pi pi-times-circle'},
};
