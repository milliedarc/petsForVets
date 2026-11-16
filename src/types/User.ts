export interface User {
    id: string;
    email: string;
    name: string;
    surname: string;
    app_mode: 'petOwner' | 'clinicTeam' | null;
    veterinary: boolean;
}