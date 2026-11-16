// Pet names by id, for labels outside the page that loaded the pet (the
// breadcrumbs show "My pets › Binki › Prescriptions"). Pages that load pets
// remember their names here; reading it from a computed or template is
// reactive, so the breadcrumb updates as soon as the name is known.

import {reactive} from "vue";

const names = reactive<Record<string, string>>({});

export function rememberPetName(pet: { id: string, name: string }): void {
    names[pet.id] = pet.name;
}

export function petName(id: string): string | undefined {
    return names[id];
}
