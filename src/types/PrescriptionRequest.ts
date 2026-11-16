import type {LegacyPet as Pet} from "@/types/Pet";
import type {Prescription} from "@/types/Prescription";
import type {User} from "@/types/User";

export interface PrescriptionRequest {
    id: string;
    status: string;
    rejection_reason: string;
    prescription: string;
    comments: string;
    created: string;
    expand: {
        pet: Pet;
        prescription: Prescription;
        approvedBy: User;
        rejectedBy: User;
    };
}
