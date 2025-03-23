<script setup lang="ts">
import PocketBase from "pocketbase";
import {computed, onMounted, ref, watch} from "vue";

const pb = new PocketBase('http://127.0.0.1:8090');
const name = ref('')
const species = ref('')
const speciesList = ref()
const breedsList = ref()
const breed = ref('')
const breedSecondary = ref('')
const dob = ref()
const gender = ref('')
const isNeutered = ref(false)
const colour = ref('')
const isImported = ref(false)
const importLocation = ref('')
const importLocationsList = ref()
const microchipNumber = ref('')

const filteredBreedsList = computed(() => {
      if (breedsList.value === undefined) {
        return []
      }
      return breedsList.value.filter((breed) => {
        return breed.species === species.value;
      })
    }
)

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

// function capitaliseInput(input: string) {
//   return input.charAt(0).toUpperCase() + input.slice(1).toLowerCase();
// }

// watch([name, colour], ([newName, newColour], [oldName, oldColour]) => {
//   if (newName !== oldName) name.value = capitaliseInput(newName);
//   if (newColour !== oldColour) colour.value = capitaliseInput(newColour);
// });

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

  const result3 = await pb.collection('importLocations').getList(1, 20, {
    sort: 'name'
  })
  // console.log('Locations:', result3);
  importLocationsList.value = result3.items;
})


</script>

<template>
  <div class="modal fade" id="addPet" tabindex="-1" aria-labelledby="exampleModalLabel" aria-hidden="true">
    <div class="modal-dialog">
      <div class="modal-content">
        <div class="modal-header">
          <h1 class="modal-title fs-5" id="exampleModalLabel">Your new pet details</h1>
          <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
        </div>
        <div class="modal-body">

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

            <!--            breed-->

            <div class="mt-2">
              <label for="petBreed">Breed <span class="text-muted">(or first breed if mixed)</span></label>
              <select v-model="breed" id="petBreed" class="form-select"
                      :class="'form-select ' + (breed === '' ? 'is-invalid' : '')">
                <option disabled selected value="">Select a breed</option>
                <option
                    v-for="breed in filteredBreedsList"
                    :value="breed"
                    :key="breed.id">
                  {{ breed.name }}

                </option>
              </select>
              <div id="validationServerUsernameFeedback" class="invalid-feedback">
                Please choose a breed.
              </div>
              <div
                  v-if="breed && !breed.generic"
                  class="mt-1">
                <label for="petSecondaryBreed">Second breed <span class="text-muted">(if mixed)</span></label>
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
              <label for="petGender">Gender</label>
              <div class="form-check" id="petGender">
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
                <option v-for="importLocation in importLocationsList"
                        :value="importLocation.id"
                        :key="importLocation.id">
                  {{ importLocation.name }}
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

          </form>

        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
          <button type="button" class="btn btn-primary">Save changes</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>

</style>