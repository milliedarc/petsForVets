<script setup lang="ts">
// imports, const, models, props, emits, refs, computed, watchers, functions, hooks
import type {Breed, Pet, Species} from "@/types/Pet";
import {useRoute, useRouter} from "vue-router";
import {onMounted, ref, computed, watch} from "vue";
import {ApiError, errorMessage} from "@/api/client";
import {createPet, deletePet, getPet, updatePet, type PetInput} from "@/api/pets";
import {listBreeds, listSpecies} from "@/api/reference";
import {rememberPetName} from "@/stores/petNames";
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
  return route.name === 'PetAdd'
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
  rememberPetName(pet);

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
  const petName = petToEdit.value?.name;
  confirm.require({
    message: `${petName}'s prescriptions and requests will be removed along with their profile. Once removed, the profile can't be restored.`,
    header: `Remove ${petName}'s profile?`,
    icon: 'pi pi-info-circle',
    rejectProps: {
      label: 'Keep profile',
      severity: 'secondary',
      outlined: true
    },
    acceptProps: {
      label: 'Remove profile',
      severity: 'danger'
    },
    defaultFocus: 'reject',
    accept: async () => {
      if (!petToEdit.value) {
        return;
      }
      try {
        await deletePet(petToEdit.value.id);
        toast.add({severity: 'info', summary: 'Profile removed', detail: `${petName}'s profile has been removed.`, life: 4000});
        await router.push('/');
      } catch (error) {
        console.error("Error deleting pet:", error);
        toast.add({severity: 'error', summary: 'Error', detail: errorMessage(error, 'The profile could not be removed. Please try again.'), life: 5000});
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

      <Card v-if="!isNewPet" class="remove-profile">
        <template #content>
          <div class="remove-profile-body">
            <div>
              <h2 class="remove-profile-title">Remove {{ petToEdit?.name }}'s profile</h2>
              <p class="m-0 text-muted text-sm">
                If {{ petToEdit?.name }} is no longer with you, or was added by mistake, you can remove their profile.
              </p>
            </div>
            <Button @click="confirmDelete"
                    label="Remove profile"
                    severity="secondary"
                    variant="outlined"
                    :aria-label="`Remove ${petToEdit?.name}'s profile`"/>
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

.remove-profile {
  margin-top: 40px;
  border: 1px solid var(--p-content-border-color);
  box-shadow: none;
}

.remove-profile-body {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.remove-profile-title {
  font-size: 1rem;
  margin-bottom: 2px;
}
</style>
