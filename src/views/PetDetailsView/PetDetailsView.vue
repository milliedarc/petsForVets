<script setup lang="ts">
// imports, const, models, props, emits, refs, computed, watchers, functions, hooks
import {useRoute, useRouter} from "vue-router";
import {pb} from "@/components/Pocketbase"
import {onMounted, ref, computed} from "vue";
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
import PetAvatarSection from "@/views/PetDetailsView/components/sections/PetAvatarSection.vue";

// ********************** CONST **************************

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
const speciesList = ref([])
const breedsList = ref([])
const breedId = ref('')
const dob = ref<Date | null>(null)
const gender = ref('')
const isNeutered = ref(false)
const colour = ref('')
const isImported = ref(false)
const importCountryCode = ref('')
const microchipNumber = ref('')
const avatarFile = ref(null)
const avatarUrl = ref('')
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

    if (deleteAvatarFlag.value === true) {
      updatedPet.avatar = null;
    }
    if (avatarFile.value) {
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

// ********************** LIFECYCLE HOOKS **************************

onMounted(async () => {
  // console.log('onMounted')
  const promises = [
    pb.collection('species').getList(1, 20, {
      sort: 'name'
    }).then(result => speciesList.value = result.items),
    pb.collection('breeds').getList(1, 100, {
      sort: '-generic,name'
    }).then(result => breedsList.value = result.items),
    pb.collection('pet_types').getList(1, 100, {
      sort: 'name'
    }).then(result => petTypes.value = result.items)
  ]

  await Promise.all(promises)

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

            <PetAvatarSection :name="name"
                              :avatar-url="avatarUrl"
                              @file-selected="deleteAvatarFlag = false; avatarFile = $event"
                              @delete-clicked="deleteAvatarFlag = true; avatarFile = null"/>

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