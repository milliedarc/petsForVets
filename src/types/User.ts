interface User {
    id: number;
    name: string;
    email: string;
    app_mode: 'petOwner' | 'clinicTeam' | undefined;
    veterinary: boolean;
}