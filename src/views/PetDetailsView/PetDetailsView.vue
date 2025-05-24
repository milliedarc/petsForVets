<script setup lang="ts">
// imports, const, models, props, emits, refs, computed, watchers, functions, hooks
import {useRoute, useRouter} from "vue-router";
import PocketBase from "pocketbase";
import {onMounted, ref, computed, watch} from "vue";
import countries from '@/types/countries.json'
import PetNameFormatted from "@/views/PetDetailsView/components/PetNameFormatted.vue";
import PetTypes from "@/views/PetDetailsView/components/PetTypes.vue";
import PetTypeIds from "@/types/PetTypeIds";
import {useToast} from "primevue/usetoast";
import {useConfirm} from "primevue/useconfirm";
import PetNameSection from "@/views/PetDetailsView/components/PetNameSection.vue";
import PetSpeciesSection from "@/views/PetDetailsView/components/PetSpeciesSection.vue";

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
const dob = ref()
const years = ref()
const months = ref()
const gender = ref()
const isNeutered = ref(false)
const colour = ref('')
const isImported = ref(false)
const importCountryCode = ref('')
const microchipNumber = ref('')
const avatar = ref(null)
const avatarFile = ref(null)
const avatarUrl = ref(null)
const deleteAvatarFlag = ref(false)

const dobTab = ref<'dobDate' | 'dobAge'>('dobDate')

// ********************** COMPUTED **************************

const sortedCountries = computed(() => {
  return countries.countries.toSorted(function (countryA, countryB) {
    if (countryA.name < countryB.name) {
      return -1;
    }
    if (countryA.name > countryB.name) {
      return 1;
    }
    return 0;
  });
})

const isValidSelectedBreed = computed(() => {
  for (const oneBreed of filteredBreedsList.value) {
    if (oneBreed.id === breedId.value) {
      return true
    }
  }
  return false
})

const filteredBreedsList = computed(() => {
      if (breedsList.value === undefined) {
        return []
      }
      return breedsList.value.filter((breed) => {
        if (breed.pet_type === petType.value) {
          return true;
        }
        return false;
      })
    }
)

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
  dob.value = getDateOnly()
  gender.value = petToEdit.value.gender
  isNeutered.value = petToEdit.value.neutered
  colour.value = petToEdit.value.colour
  isImported.value = petToEdit.value.imported
  importCountryCode.value = petToEdit.value.import_country_code
  microchipNumber.value = petToEdit.value.microchip_number
  avatarUrl.value = pb.files.getURL(petToEdit.value, petToEdit.value.avatar, {'thumb': '100x250'});
  console.log('avatar URL:', avatarUrl.value)
}

function getDateOnly() {
  if (petToEdit.value.date_of_birth == '') {
    return undefined;
  }
  return petToEdit.value.date_of_birth.split(' ')[0];
}

function calculateBirthDate(years: number, months: number): Date {
  const today = new Date();

  // Clone today's date
  const birthDate = new Date(today);

  // Subtract years and months
  birthDate.setFullYear(birthDate.getFullYear() - years);
  birthDate.setMonth(birthDate.getMonth() - months);

  // Set to the first day of the month
  birthDate.setDate(1);

  // Set time to 00:00:00
  birthDate.setHours(0, 0, 0, 0);

  return birthDate;
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
    if (dobTab.value === 'dobAge') {
      dob.value = calculateBirthDate(years.value, months.value)
    }

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

            <PetSpeciesSection v-model:species="species" :name="name" :petType="petType" :speciesList="speciesList"/>

            <section v-if="petType === PetTypeIds.catId || petType === PetTypeIds.dogId">
              <div>
                <p>What breed is
                  <PetNameFormatted :name="name"/>
                  <span>?</span>
                </p>

                <div class="mt-4">

                  <FloatLabel variant="on">
                    <Select v-model="breedId" inputId="breed"
                            :options="filteredBreedsList"
                            optionLabel="name"
                            optionValue="id"
                            editable
                            class="myInput"
                            :invalid="!isValidSelectedBreed"/>
                    <label for="breed">Breed</label>
                  </FloatLabel>
                  <Message
                      v-if="!isValidSelectedBreed" severity="error" size="small" variant="simple"
                      aria-describedby="breedError">
                    Please choose a valid breed or select <span class="fw-bold">Cross / Mixed Breed</span> if unknown
                  </Message>
                </div>
              </div>
            </section>

            <section>
              <div>
                <p class="mb-1">How old is
                  <PetNameFormatted :name="name"/>
                  <span>?</span>
                </p>
              </div>
              <div>

                <Tabs v-model:value="dobTab" style="width: 350px">
                  <TabList>
                    <Tab value="dobDate">Date of birth</Tab>
                    <Tab value="dobAge">Age in years</Tab>
                  </TabList>

                  <TabPanels>
                    <TabPanel value="dobDate">
                      <div class="d-flex">
                        <label for="dob" class="font-bold block mb-2"></label>
                        <DatePicker v-model="dob" showIcon fluid iconDisplay="input" inputId="dob" class="myInput"/>
                      </div>
                    </TabPanel>

                    <TabPanel value="dobAge">
                      <div class="d-flex gap-3">
                        <FloatLabel variant="on">
                          <InputText type="number" min="0" max="30" id="years" style="width: 150px"
                                     v-model="years"/>
                          <label for="years">Years *</label>
                        </FloatLabel>
                        <FloatLabel variant="on">
                          <InputText type="number" min="0" max="11" id="months" style="width: 150px"
                                     v-model="months"/>
                          <label for="months">Months *</label>
                        </FloatLabel>
                      </div>
                    </TabPanel>
                  </TabPanels>
                </Tabs>
                <p class="fw-light" style="font-size: 0.8rem">
                  Don't worry if you don't know the exact date. Choose
                  <span class="fw-bold">Age in years</span>
                  instead and give us an estimate!
                </p>
              </div>
            </section>

            <section>
              <div>
                <p>What's
                  <PetNameFormatted :name="name"/>
                  <span>'s gender?</span>
                </p>
              </div>

              <div class="d-flex gap-4">
                <div class="d-flex align-items-center">
                  <RadioButton v-model="gender" inputId="male" name="gender" value="male"/>
                  <label class="ms-3" for="male">Male</label>
                </div>
                <div class="d-flex align-items-center">
                  <RadioButton v-model="gender" inputId="female" name="gender" value="female"/>
                  <label class="ms-3" for="female">Female</label>
                </div>
              </div>
            </section>

            <section>
              <div>
                <p>Is
                  <PetNameFormatted :name="name"/>
                  <span> neutered?</span>
                </p>
              </div>

              <div class="flex items-center gap-2">
                <Checkbox v-model="isNeutered"
                          binary
                          inputId="neutered" name="neutered" value="true"/>
                <label class="ms-3" for="neutered">
                  Please tick this box if your pet is neutered
                </label>
              </div>
            </section>

            <section>
              <div>
                <p>What colour is
                  <PetNameFormatted :name="name"/>
                  <span>?</span>
                </p>
              </div>

              <div>
                <FloatLabel class="myInput" variant="on">
                  <InputText id="colour" v-model="colour" class="myInput"/>
                  <label for="colour">Colour</label>
                </FloatLabel>
              </div>
            </section>

            <section>
              <div>
                <p>Is
                  <PetNameFormatted :name="name"/>
                  <span> imported?</span>
                </p>
              </div>

              <div class="flex items-center gap-2">
                <Checkbox v-model="isImported"
                          binary
                          inputId="imported" name="imported" value="true"/>
                <label class="ms-3" for="imported">
                  Please tick this box if your pet is imported
                </label>
              </div>

              <div class="mt-3" v-if="isImported">

                <FloatLabel variant="on">
                  <Select v-model="importCountryCode" editable inputId="country" :options="sortedCountries"
                          optionLabel="name"
                          optionValue="country_code"
                          class="myInput"/>
                  <label for="country">Select import country</label>
                </FloatLabel>
              </div>
            </section>

            <section>

              <div>
                <p>What's
                  <PetNameFormatted :name="name"/>
                  <span>'s microchip number?</span>
                </p>
              </div>

              <div>
                <FloatLabel variant="on">
                  <InputText id="microchip" v-model="microchipNumber" class="myInput"/>
                  <label for="microchip">Microchip Number</label>
                </FloatLabel>
              </div>
            </section>

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