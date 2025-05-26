<script setup lang="ts">
// imports, const, models, props, emits, refs, computed, watchers, functions, hooks

import {useRoute, useRouter} from "vue-router";
import PocketBase from "pocketbase";
import {onMounted, ref, computed, watch} from "vue";
import PetNameFormatted from "@/views/PetDetailsView/components/PetNameFormatted.vue";
import PetTypes from "@/views/PetDetailsView/components/PetTypes.vue";
import PetTypeIds from "@/types/PetTypeIds";
import {useToast} from "primevue/usetoast";
import {useConfirm} from "primevue/useconfirm";
import PetNameSection from "@/views/PetDetailsView/components/PetNameSection.vue";
import PetSpeciesSection from "@/views/PetDetailsView/components/PetSpeciesSection.vue";
import PetBreedSection from "@/views/PetDetailsView/components/PetBreedSection.vue";
import PetAgeSection from "@/views/PetDetailsView/components/PetAgeSection.vue";
import PetGenderSection from "@/views/PetDetailsView/components/PetGenderSection.vue";
import PetIsNeuteredSection from "@/views/PetDetailsView/components/PetIsNeuteredSection.vue";
import PetIsImportedSection from "@/views/PetDetailsView/components/PetIsImportedSection.vue";
import PetColourSection from "@/views/PetDetailsView/components/PetColourSection.vue";
import PetMicrochipNumberSection from "@/views/PetDetailsView/components/PetMicrochipNumberSection.vue";

// ********************** CONST **************************

const pb = new PocketBase('http://127.0.0.1:8090');

const route = useRoute()
const router = useRouter()
const toast = useToast();
const confirm = useConfirm();

// ********************** REF **************************

const petToEdit = ref<Pet | undefined>(undefined);
const petNotFound = ref(false)
const isLoading = ref(true)

const name = ref('')
const petType = ref('')
const petTypes = ref([])
const species = ref('')
const speciesList = ref()
const breedsList = ref()
const breedId = ref('')
const dob = ref<Date | null>(null)
const gender = ref('')
const isNeutered = ref(false)
const colour = ref('')
const isImported = ref(false)
const importCountryCode = ref('')
const microchipNumber = ref('')
const avatar = ref(null)
const avatarFile = ref(null)
const avatarUrl = ref(null)
const deleteAvatarFlag = ref(false)

// ********************** COMPUTED **************************

const isValidName = computed(() => {
  return name.value.trim().length > 0;
})

const canSave = computed(() => {
  if (isValidName.value && petType.value !== '') {
    return true;
  }
  return false;
})

const isNewPet = computed(() => {
  return route.params.id === 'add'
})

const hasAvatarUrl = computed(() => {
  return avatarUrl.value !== null && avatarUrl.value !== '' && !avatar.value
})

// ********************** FUNCTIONS **************************

async function fetchPet() {
  const result = await pb.collection('pets').getOne(route.params.id as string, {
    expand: 'species,breed'
  })
  petToEdit.value = result as any; // assigns content to 'pets'
  console.log(result)

  name.value = petToEdit.value.name
  petType.value = petToEdit.value.expand.species.pet_type
  species.value = petToEdit.value.species
  breedId.value = petToEdit.value.breed
  dob.value = getDateOnly(petToEdit.value.date_of_birth)
  gender.value = petToEdit.value.gender
  isNeutered.value = petToEdit.value.neutered
  colour.value = petToEdit.value.colour
  isImported.value = petToEdit.value.imported
  importCountryCode.value = petToEdit.value.import_country_code
  microchipNumber.value = petToEdit.value.microchip_number
  avatarUrl.value = pb.files.getURL(petToEdit.value, petToEdit.value.avatar, {'thumb': '100x250'});
}

function getDateOnly(dateString: string) {
  if (dateString == '') {
    return null;
  }
  return new Date(dateString.split(' ')[0]);
}

function confirmDelete() {
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
      await pb.collection('pets').delete(petToEdit.value.id);
      toast.add({severity: 'info', summary: 'Confirmed', detail: 'Pet profile deleted', life: 3000});
      await router.push('/');
    }
  });
}

function findSpeciesByName(speciesName: string): string {
  const foundSpecies = speciesList.value.find((species) => {
    return species.name === speciesName;
  })
  return foundSpecies.id;
}

function resetOnClickPetType() {
  breedId.value = '';
  species.value = '';
}

async function savePet() {
  try {
    let mySpecies = species.value;
    if (petType.value === PetTypeIds.catId) {
      mySpecies = findSpeciesByName('Cat')
    } else if (petType.value === PetTypeIds.dogId) {
      mySpecies = findSpeciesByName('Dog')
    }

    const updatedPet = {
      user: pb.authStore.record.id,
      name: name.value,
      species: mySpecies,
      breed: breedId.value,
      date_of_birth: dob.value,
      gender: gender.value,
      neutered: isNeutered.value,
      colour: colour.value,
      imported: isImported.value,
      import_country_code: importCountryCode.value,
      microchip_number: microchipNumber.value,
      avatar: undefined
    };

    if (avatar.value || deleteAvatarFlag.value === true) {
      updatedPet.avatar = avatarFile.value;
    }

    if (petToEdit.value?.id !== undefined) {
      await pb.collection('pets').update(petToEdit.value.id, updatedPet);
    } else {
      await pb.collection('pets').create(updatedPet);
    }
    await router.push('/');
    toast.add({severity: 'success', summary: 'Success!', detail: 'Pet details saved', life: 3000});

  } catch (error) {
    console.error("Error saving pet:", error);
    toast.add({severity: 'error', summary: 'Error', detail: 'Pet details could not be saved', life: 3000});
  }
}

function onFileSelect(event) {
  const file = event.files[0];
  const reader = new FileReader();

  reader.onload = async (e) => {
    avatar.value = e.target.result;
  };
  avatarFile.value = file;
  reader.readAsDataURL(file);
  deleteAvatarFlag.value = false;
}

function deletePetAvatar() {
  avatarUrl.value = ''
  avatarFile.value = null;
  deleteAvatarFlag.value = true;
  avatar.value = null;
}

// ********************** LIFECYCLE HOOKS **************************

onMounted(async () => {
  // console.log('onMounted')
  const result = await pb.collection('species').getList(1, 20, {
    sort: 'name'
  });
  // console.log('species:', result)
  speciesList.value = result.items;

  const result2 = await pb.collection('breeds').getList(1, 100, {
    sort: '-generic,name'
  })
  // console.log('Breeds:', result2)
  breedsList.value = result2.items;

  const result3 = await pb.collection('pet_types').getList(1, 100, {
    sort: 'name'
  })
  petTypes.value = result3.items;

  try {
    if (!isNewPet.value) {
      await fetchPet()
    }
  } catch (e) {
    petNotFound.value = true;
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

            <PetBreedSection v-if="petType === PetTypeIds.catId || petType === PetTypeIds.dogId"
                             :name="name"
                             :pet-type="petType"
                             :breeds-list="breedsList"
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

            <section>
              <div>
                <p class="mb-0">Upload a picture of
                  <PetNameFormatted :name="name"/>
                  <span>:</span>
                </p>
                <p class="fw-light" style="font-size: 0.8rem">
                  You can skip this step and add it later.
                </p>
              </div>

              <div class="card">
                <div class="card-body">
                  <div class="d-flex flex-column align-items-center">
                    <div style="position: relative">
                      <img v-if="hasAvatarUrl" :src="avatarUrl" :alt="`${name}'s uploaded picture`"
                           class="shadow-lg rounded mb-3"
                           style="max-width: 240px"/>
                      <img v-else-if="avatar" :src="avatar" alt="Avatar" class="shadow-lg rounded mb-3"
                           style="max-width: 240px"/>
                      <Button v-if="hasAvatarUrl || avatar" icon="pi pi-times"
                              @click="deletePetAvatar"
                              severity="danger" size="small" rounded
                              :aria-label="`Remove uploaded photo of ${name}`"
                              style="position: absolute; top: -15px; right: -15px;"/>
                    </div>
                    <FileUpload mode="basic" @select="onFileSelect" customUpload severity="secondary"
                                class="p-button-outlined"/>
                  </div>
                </div>
              </div>
            </section>

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