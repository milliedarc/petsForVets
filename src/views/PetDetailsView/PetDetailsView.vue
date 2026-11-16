<script setup lang="ts">
// imports, const, models, props, emits, refs, computed, watchers, functions, hooks
import type {Breed, Pet, Species} from "@/types/Pet";
import {useRoute, useRouter} from "vue-router";
import {onMounted, ref, computed, watch} from "vue";
import {ApiError, errorMessage} from "@/api/client";
import {createPet, deletePet, getPet, updatePet, type PetInput} from "@/api/pets";
import {listBreeds, listSpecies} from "@/api/reference";
import {parseIsoDate, toIsoDate} from "@/utils/dates";
import PetNameFormatted from "@/views/PetDetailsView/components/PetNameFormatted.vue";
import PetTypes from "@/views/PetDetailsView/components/PetTypes.vue";
import PetTypeIds from "@/types/PetTypeIds";
import {useToast} from "primevue/usetoast";
import {useConfirm} from "primevue/useconfirm";
import PetNameSection from "@/views/PetDetailsView/components/sections/PetNameSection.vue";
import PetSpeciesSection from "@/views/PetDetailsView/components/sections/PetSpeciesSection.vue";
import PetBreedSection from "@/views/PetDetailsView/components/sections/PetBreedSection.vue";
import PetAgeSection from "@/views/PetDetailsView/components/sections/PetAgeSection.vue";
import PetGenderSection from "@/views/PetDetailsView/components/sections/PetGenderSection.vue";
import PetIsNeuteredSection from "@/views/PetDetailsView/components/sections/PetIsNeuteredSection.vue";
import PetIsImportedSection from "@/views/PetDetailsView/components/sections/PetIsImportedSection.vue";
import PetColourSection from "@/views/PetDetailsView/components/sections/PetColourSection.vue";
import PetMicrochipNumberSection from "@/views/PetDetailsView/components/sections/PetMicrochipNumberSection.vue";

// ********************** CONST **************************

const route = useRoute()
const router = useRouter()
const toast = useToast();
const confirm = useConfirm();

// ********************** REF **************************

const petToEdit = ref<Pet | undefined>(undefined);
const petNotFound = ref<boolean>(false)
const isLoading = ref<boolean>(true)
const isSaving = ref<boolean>(false)

const name = ref<string>('')
const petType = ref<string>('')
const species = ref<string>('')
const speciesList = ref<Species[]>([])
const breedsList = ref<Breed[]>([])
const breedId = ref<string>('')
const dob = ref<Date | null>(null)
const gender = ref<string>('')
const isNeutered = ref<boolean>(false)
const colour = ref<string>('')
const isImported = ref<boolean>(false)
const importCountryCode = ref<string>('')
const microchipNumber = ref<string>('')

// ********************** COMPUTED **************************

const isValidName = computed<boolean>(() => {
  return name.value.trim().length > 0;
})

/** Cats and dogs have exactly one species; other pet types pick one from a list. */
const hasSingleSpecies = computed<boolean>(() => {
  return petType.value === PetTypeIds.catId || petType.value === PetTypeIds.dogId;
})

const canSave = computed<boolean>(() => {
  return isValidName.value
      && petType.value !== ''
      && (hasSingleSpecies.value || species.value !== '')
      && !isSaving.value;
})

/** The pet's species id: implied for cats and dogs, chosen from a list otherwise. */
const selectedSpeciesId = computed<string | undefined>(() => {
  if (petType.value === PetTypeIds.catId) {
    return findSpeciesByName('Cat');
  }
  if (petType.value === PetTypeIds.dogId) {
    return findSpeciesByName('Dog');
  }
  return species.value || undefined;
})

/** Breeds of the selected species. Empty means the species has no breeds, so the question is hidden. */
const breedsForSpecies = computed<Breed[]>(() => {
  return breedsList.value.filter((breed) => breed.species_id === selectedSpeciesId.value);
})

const isNewPet = computed<boolean>(() => {
  return route.params.id === 'add'
})

// ********************** WATCHERS **************************

// If the species changes (e.g. Rabbit -> Guinea pig), drop a breed that no
// longer fits. The database would reject the mismatch anyway.
watch(breedsForSpecies, (breeds) => {
  if (breedId.value !== '' && !breeds.some((breed) => breed.id === breedId.value)) {
    breedId.value = '';
  }
})

// ********************** FUNCTIONS **************************

async function fetchPet(id: string): Promise<void> {
  const pet = await getPet(id)
  petToEdit.value = pet;

  name.value = pet.name
  petType.value = pet.species.pet_type_id
  species.value = pet.species.id
  breedId.value = pet.breed?.id ?? ''
  dob.value = pet.date_of_birth ? parseIsoDate(pet.date_of_birth) : null
  gender.value = pet.gender
  isNeutered.value = pet.neutered
  colour.value = pet.colour
  isImported.value = pet.imported
  importCountryCode.value = pet.import_country_code
  microchipNumber.value = pet.microchip_number
}

function confirmDelete(): void {
  confirm.require({
    message: `Delete ${petToEdit.value?.name}'s profile? Their prescriptions and requests will be deleted too. This can't be undone.`,
    header: 'Delete pet profile',
    icon: 'pi pi-exclamation-triangle',
    rejectLabel: 'Cancel',
    rejectProps: {
      label: 'Cancel',
      severity: 'secondary',
      outlined: true
    },
    acceptProps: {
      label: 'Delete',
      severity: 'danger'
    },
    defaultFocus: 'reject',
    accept: async () => {
      if (!petToEdit.value) {
        return;
      }
      try {
        await deletePet(petToEdit.value.id);
        toast.add({severity: 'info', summary: 'Confirmed', detail: 'Pet profile deleted', life: 3000});
        await router.push('/');
      } catch (error) {
        console.error("Error deleting pet:", error);
        toast.add({severity: 'error', summary: 'Error', detail: errorMessage(error, 'Pet could not be deleted'), life: 5000});
      }
    }
  });
}

function findSpeciesByName(speciesName: string): string | undefined {
  return speciesList.value.find((oneSpecies) => oneSpecies.name === speciesName)?.id;
}

function resetOnClickPetType(): void {
  breedId.value = '';
  species.value = '';
}

async function savePet(): Promise<void> {
  const speciesId = selectedSpeciesId.value;
  if (!speciesId) {
    toast.add({severity: 'error', summary: 'Error', detail: 'Please choose a species', life: 3000});
    return;
  }

  const input: PetInput = {
    name: name.value,
    species_id: speciesId,
    breed_id: breedId.value !== '' ? breedId.value : null,
    date_of_birth: dob.value ? toIsoDate(dob.value) : null,
    gender: gender.value,
    neutered: isNeutered.value,
    colour: colour.value,
    imported: isImported.value,
    import_country_code: importCountryCode.value ?? '',
    microchip_number: microchipNumber.value,
  };

  isSaving.value = true;
  try {
    if (petToEdit.value) {
      await updatePet(petToEdit.value.id, input);
    } else {
      await createPet(input);
    }
    await router.push('/');
    toast.add({severity: 'success', summary: 'Success!', detail: 'Pet details saved', life: 3000});
  } catch (error) {
    console.error("Error saving pet:", error);
    toast.add({severity: 'error', summary: 'Error', detail: errorMessage(error, 'Pet details could not be saved'), life: 5000});
  } finally {
    isSaving.value = false;
  }
}

// ********************** LIFECYCLE HOOKS **************************

onMounted(async () => {
  try {
    const [allSpecies, allBreeds] = await Promise.all([listSpecies(), listBreeds()]);
    speciesList.value = allSpecies;
    breedsList.value = allBreeds;

    if (!isNewPet.value) {
      await fetchPet(route.params.id as string)
    }
  } catch (error) {
    if (error instanceof ApiError && error.status === 404) {
      petNotFound.value = true;
    } else {
      console.error("Error loading pet form:", error);
      toast.add({severity: 'error', summary: 'Error', detail: errorMessage(error, 'Could not load the pet form'), life: 5000});
    }
  } finally {
    isLoading.value = false;
  }
})

</script>

<template>
  <main class="page page--narrow">
    <RouterLink to="/" class="back-link">
      <i class="pi pi-arrow-left" aria-hidden="true"></i>My pets
    </RouterLink>

    <div v-if="isLoading" class="state-panel" aria-live="polite">
      <ProgressSpinner style="width: 40px; height: 40px" strokeWidth="5"/>
      <span>Loading…</span>
    </div>

    <div v-else-if="petNotFound" class="state-panel">
      <h2>Pet not found</h2>
      <p>This pet doesn't exist, or it isn't on your account.</p>
      <RouterLink to="/">Back to my pets</RouterLink>
    </div>

    <template v-else>
      <div class="page-header">
        <h1 v-if="isNewPet">Add a pet</h1>
        <h1 v-else>Edit {{ petToEdit?.name }}'s details</h1>
      </div>

      <form @submit.prevent="savePet">
        <Card>
          <template #content>
            <PetNameSection v-model:name="name"/>

            <section class="question">
              <span class="question-label">
                <PetNameFormatted :name="name"/>
                is a:
              </span>
              <PetTypes v-model="petType" @petTypeClicked="resetOnClickPetType"/>
            </section>

            <PetSpeciesSection :name="name"
                               :pet-type="petType"
                               :species-list="speciesList"
                               v-model:species="species"/>

            <PetBreedSection v-if="breedsForSpecies.length > 0"
                             :name="name"
                             :breeds-list="breedsForSpecies"
                             v-model:breed-id="breedId"/>

            <PetAgeSection :name="name"
                           v-model:date-of-birth="dob"/>

            <PetGenderSection :name="name"
                              v-model:gender="gender"/>

            <PetIsNeuteredSection :name="name"
                                  v-model:is-neutered="isNeutered"/>

            <PetColourSection :name="name"
                              v-model:colour="colour"/>

            <PetIsImportedSection :name="name"
                                  v-model:is-imported="isImported"
                                  v-model:import-country-code="importCountryCode"/>

            <PetMicrochipNumberSection :name="name"
                                       v-model:microchip-number="microchipNumber"/>

            <!-- TODO: PetAvatarSection returns when avatar uploads are ported to the Go API -->
          </template>
        </Card>

        <div class="form-actions">
          <Button type="submit"
                  :disabled="!canSave"
                  :loading="isSaving"
                  :label="isNewPet ? 'Save pet' : 'Save changes'"
                  icon="pi pi-check"/>
          <RouterLink to="/" custom v-slot="{navigate}">
            <Button label="Cancel" severity="secondary" variant="outlined" @click="navigate"/>
          </RouterLink>
          <small v-if="!canSave && !isSaving" class="field-hint form-actions-hint">
            Add a name and choose what type of pet this is to save.
          </small>
        </div>
      </form>

      <Card v-if="!isNewPet" class="danger-zone">
        <template #content>
          <div class="danger-zone-body">
            <div>
              <h2 class="danger-zone-title">Delete this pet</h2>
              <p class="m-0 text-muted text-sm">Removes the profile, prescriptions and requests for good.</p>
            </div>
            <Button @click="confirmDelete"
                    label="Delete"
                    icon="pi pi-trash"
                    severity="danger"
                    variant="outlined"
                    :aria-label="`Delete ${petToEdit?.name}'s profile`"/>
          </div>
        </template>
      </Card>
    </template>
  </main>
</template>

<style scoped>
.form-actions {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 12px;
  margin-top: 20px;
}

.form-actions-hint {
  flex-basis: 100%;
  margin-top: 0;
}

.danger-zone {
  margin-top: 40px;
  border: 1px solid var(--p-red-200);
  box-shadow: none;
}

.danger-zone-body {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.danger-zone-title {
  font-size: 1rem;
  margin-bottom: 2px;
}
</style>
