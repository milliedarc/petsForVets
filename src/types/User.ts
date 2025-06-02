interface User {
    id: number;
    email: string;
    name: string;
    surname: string;
    app_mode: 'petOwner' | 'clinicTeam' | null;
    veterinary: boolean;
}