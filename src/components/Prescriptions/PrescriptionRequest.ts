interface PrescriptionRequest {
    id: string;
    status: string;
    rejection_reason: string;
    prescription: string;
    comments: string;
    created: string;
    expand: Expand;
}

interface Expand {
    pet: Pet;
    prescription: Prescription;
    approvedBy: User;
    rejectedBy: User;
}