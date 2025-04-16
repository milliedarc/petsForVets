<script setup lang="ts">
// imports, consts, models, props, emits, refs, computed, watchers, functions, hooks
import {useRoute} from "vue-router";
import PocketBase from "pocketbase";
import {onMounted, ref, computed, inject} from "vue";
import countries from './countries.json'


const pb = new PocketBase('http://127.0.0.1:8090');

const petToEdit = defineModel<Pet>()

const emit = defineEmits(['savePet']);

const pet = ref<Pet | undefined>(undefined);

const petNotFound = ref(false)
const isLoading = ref(true)
const route = useRoute()

const name = ref('')
const species = ref('')
const speciesList = ref()
const breedsList = ref()
const breedGroup = ref('')
const breed = ref<Breed>()
const breedSecondary = ref('')
const dob = ref()
const gender = ref('')
const isNeutered = ref(false)
const colour = ref('')
const isImported = ref(false)
const importLocation = ref('')
const microchipNumber = ref('')

const filteredBreedsList = computed(() => {
      if (breedsList.value === undefined) {
        return []
      }
      return breedsList.value.filter((breed) => {
        if (breed.species === species.value) {
          if (breedGroup.value === 'generic' && breed.generic) {
            return true;
          } else if (breedGroup.value !== 'generic' && !breed.generic) {
            return true;
          }
          return false;
        }
      })
    }
)

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

const calculatedAge = computed(() => {
  if (dob.value === undefined) {
    return ''
  }
  const birthDate = new Date(dob.value);
  const today = new Date();

  let years = today.getFullYear() - birthDate.getFullYear();
  let months = today.getMonth() - birthDate.getMonth();
  const dayDiff = today.getDate() - birthDate.getDate();

  // Adjust if birthday hasn't occurred yet this month
  if (dayDiff < 0) {
    months--;
  }

  // Adjust if month difference is negative
  if (months < 0) {
    years--;
    months += 12;
  }

  return `${years}y ${months}m`;
})

const isValidName = computed(() => {
  return name.value.trim().length > 0;
})

const canSave = computed(() => {
  if (isValidName.value && species.value !== '') {
    return true;
  }
  return false;
})

async function fetchPet() {
  const result = await pb.collection('pets').getOne(route.params.id, {
    expand: 'species,breed,breed_secondary'
  })
  // @ts-ignore
  pet.value = result; // assigns content to 'pets'
  console.log(result)
}

onMounted(async () => {
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
})

function resetForm() {
  name.value = ''
  species.value = ''
  breedGroup.value = ''
  breed.value = undefined
  breedSecondary.value = ''
  dob.value = undefined
  gender.value = ''
  isNeutered.value = false
  colour.value = ''
  isImported.value = false
  importLocation.value = ''
  microchipNumber.value = ''
}

async function savePet() {
  try {
    const updatedPet = {
      user: pb.authStore.record.id,
      name: name.value,
      species: species.value,
      breed_group: breedGroup.value,
      breed: breed.value?.id,
      breed_secondary: breedSecondary.value,
      date_of_birth: dob.value,
      gender: gender.value,
      neutered: isNeutered.value,
      colour: colour.value,
      imported: isImported.value,
      import_country_code: importLocation.value,
      microchip_number: microchipNumber.value
    };
    if (petToEdit.value?.id !== undefined) {
      await pb.collection('pets').update(petToEdit.value.id, updatedPet);
    } else {
      await pb.collection('pets').create(updatedPet);
    }
    // Get the modal instance using Bootstrap's library and hide it
    // @ts-ignore
  } catch (error) {
    console.error("Error saving pet:", error);
  }
}

function getDateOnly() {
  if (petToEdit.value.date_of_birth == '') {
    return undefined;
  }
  return petToEdit.value.date_of_birth.split(' ')[0];
}

function fetchPetToEdit() {
  name.value = petToEdit.value.name
  species.value = petToEdit.value.species
  breedGroup.value = petToEdit.value.breed_group
  breed.value = petToEdit.value.expand.breed
  breedSecondary.value = petToEdit.value.breed_secondary
  dob.value = getDateOnly()
  gender.value = petToEdit.value.gender
  isNeutered.value = petToEdit.value.neutered
  colour.value = petToEdit.value.colour
  isImported.value = petToEdit.value.imported
  importLocation.value = petToEdit.value.import_country_code
  microchipNumber.value = petToEdit.value.microchip_number
}

</script>

<template>

  <div v-if="!isLoading">
    <div v-if="petNotFound">
      <h1>Pet not found</h1>
    </div>
    <div v-else>
      <h1>View & edit your pet</h1>
      <h2> {{ pet?.name }}</h2>
    </div>
  </div>

  <h1 v-if="petToEdit === undefined">Your new pet details</h1>
  <h1 v-else>Edit pet details</h1>
  <button type="button" class="btn-close"
          aria-label="Close"></button>

  <form>
    <!--             name-->
    <div class="mt-1">
      <label for="petName">Your pet's name</label>
      <input v-model.trim="name"
             autofocus
             :class="'form-control ' + (!isValidName ? 'is-invalid' : '')"
             type="text" id="petName" placeholder="e.g. Toby"

             required/>
      <div v-if="!isValidName"
           id="validationServerUsernameFeedback" class="invalid-feedback">
        Please choose a name for your pet.
      </div>
    </div>

    <!--            species-->

    <div class="mt-2">
      <label for="petSpecies" class="mt-2">Species</label>
      <select v-model="species" id="petSpecies"
              :class="'form-select ' + (species === '' ? 'is-invalid' : '')">
        <option disabled selected value="">Select a species</option>
        <option
            v-for="specie in speciesList"
            :value="specie.id"
            :key="specie.id">
          {{ specie.name }}
        </option>
      </select>
      <div id="validationServerUsernameFeedback" class="invalid-feedback">
        Please choose a species.
      </div>
    </div>


    <section v-if="species !== ''">

      <!--            breed group-->

      <div class="mt-2">
        <label for="petBreedGroup">Breed Group</label>
        <select v-model="breedGroup" id="petBreedGroup" class="form-select">
          <option disabled selected value="">Select a breed group</option>
          <option value="generic">Generic</option>
          <option value="pedigree">Pedigree</option>
          <option value="crossbreed">Crossbreed</option>
        </select>
      </div>


      <section v-if="breedGroup !== ''" class="mt-2">

        <!--            breed -->

        <div class="mt-2">
          <label for="petBreed">Breed</label>
          <select v-model="breed" id="petBreed" class="form-select">
            <option disabled selected value="">Select a breed</option>
            <option
                v-for="breed in filteredBreedsList"
                :value="breed"
                :key="breed.id">
              {{ breed.name }}

            </option>
          </select>

          <div
              v-if="breedGroup === 'crossbreed'"
              class="mt-1">
            <label for="petSecondaryBreed">Second breed</label>
            <select v-model="breedSecondary" id="petSecondaryBreed" class="form-select">
              <option
                  v-for="breed in filteredBreedsList"
                  :value="breed.id"
                  :key="breed.id"
                  :disabled="breed.generic">
                {{ breed.name }}

              </option>
            </select>
          </div>
        </div>


        <!--            date of birth-->

        <div class="mt-2">
          <label for="dob">Date of Birth</label>
          <input v-model="dob" id="dob" class="d-inline form-control" type="date"/>
          <div v-if="calculatedAge !== ''">
            Age: {{ calculatedAge }}
          </div>
        </div>

        <!--            gender-->

        <div class="mt-2">
          <label>Gender</label>
          <div class="form-check">
            <input v-model="gender" class="form-check-input" type="radio" name="flexRadioDefault"
                   id="flexRadioDefault1" value="male">
            <label class="form-check-label" for="flexRadioDefault1">
              Male
            </label>
          </div>
          <div class="form-check">
            <input v-model="gender" class="form-check-input" type="radio" name="flexRadioDefault"
                   id="flexRadioDefault2" value="female">
            <label class="form-check-label" for="flexRadioDefault2">
              Female
            </label>
          </div>
        </div>

        <!--            neutered-->

        <div class="form-check mt-4">
          <input v-model="isNeutered" class="form-check-input" type="checkbox" id="flexCheckDefault3">
          <label class="form-check-label" for="flexCheckDefault3">
            Neutered
          </label>
        </div>

        <!--            colour-->

        <div class="mt-4">
          <label for="petColour">Your pet's colour</label>
          <input
              v-model="colour"
              class="form-control"
              type="text" id="petColour" placeholder="e.g. Black & White"
          />
        </div>

        <!--            imported-->

        <div class="form-check mt-4">
          <input v-model="isImported"
                 class="form-check-input" type="checkbox" value="" id="flexCheckDefault4"
          >
          <label class="form-check-label" for="flexCheckDefault4">
            Imported
          </label>
        </div>

        <!--            import location-->

        <div class="mt-2" v-if="isImported">
          <label for="importLocation">From where?</label>
          <select v-model="importLocation" id="importLocation" class="form-select">
            <option v-for="country in sortedCountries"
                    :value="country.country_code"
                    :key="country.country_code">
              {{ country.name }}
            </option>
          </select>
        </div>

        <!--            microchip number-->

        <div class="mt-3">
          <label for="microchipNumber">Microchip number</label>
          <input
              v-model="microchipNumber"
              class="form-control"
              type="text" id="microchipNumber"/>
        </div>
      </section>
    </section>
  </form>

  <div class="modal-footer">
    <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
    <button @click="savePet" type="button" class="btn btn-primary"
            :disabled="!canSave">Save changes
    </button>
  </div>

</template>

<style scoped>

</style>