<script setup lang="ts">
// imports, const, models, props, emits, refs, computed, watchers, functions, hooks
import type {Breed, Pet, Species} from "@/types/Pet";
import {useRoute, useRouter} from "vue-router";
import {onMounted, ref, computed, watch} from "vue";
import {ApiError} from "@/api/client";
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

/** The API's error message if there is one, so server-side validation reaches the user. */
function errorMessage(error: unknown, fallback: string): string {
  return error instanceof ApiError ? error.message : fallback;
}

function confirmDelete(): void {
  confirm.require({
    message: "Do you want to delete this pet's profile?",
    header: 'Warning!',
    icon: 'pi pi-info-circle',
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
  <main class="container-fluid my-5">
    <div class="row justify-content-center">
      <div class="col-8">

        <div v-if="!isLoading">
          <div v-if="petNotFound">
            <h4>Pet not found</h4>
          </div>
          <div v-else>
            <div class="my-3">
              <h4 v-if="petToEdit === undefined">Your new pet details</h4>
              <h4 v-else>Edit pet details</h4>
            </div>

            <PetNameSection v-model:name="name"/>

            <section>
              <p>
                <PetNameFormatted :name="name"/>
                is a:
              </p>
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

            <section>
              <Button @click="savePet"
                      :disabled="!canSave"
                      class="myInput"
                      :aria-label="`Save ${name}'s profile`">
                <span v-if="isNewPet">Save pet</span>
                <span v-else>Update pet details</span>
              </Button>
            </section>
            <section>
              <Button v-if="!isNewPet"
                      @click="confirmDelete"
                      label="Delete pet profile"
                      variant="text" severity="danger"
                      class="myInput"
                      :aria-label="`Delete ${name}'s profile`"/>
            </section>
          </div>
        </div>
      </div>
    </div>
  </main>
</template>

<style scoped>

.myInput {
  width: 300px;
}

section {
  margin-top: 30px;
}

</style>