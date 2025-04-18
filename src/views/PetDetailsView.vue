<script setup lang="ts">
import {useRoute} from "vue-router";
import PocketBase from "pocketbase";
import {onMounted, ref, computed} from "vue";
import countries from '../components/countries.json'
import PetNameFormatted from "@/components/PetDetails/PetNameFormatted.vue";
import PetTypes from "@/components/PetDetails/PetTypes.vue";

const route = useRoute()

const petToEdit = ref<Pet | undefined>(undefined);
const petNotFound = ref(false)
const isLoading = ref(true)

const name = ref('')
const petType = ref('')
const petTypes = ref([])
const species = ref('')
const speciesList = ref()
const breedsList = ref()
// const breedGroup = ref('')
const breed = ref<Breed>()
// const breedSecondary = ref('')
const dob = ref()
const gender = ref('')
const isNeutered = ref(false)
const colour = ref('')
const isImported = ref(false)
const importLocation = ref('')
const microchipNumber = ref('')

const pb = new PocketBase('http://127.0.0.1:8090');

async function fetchPet() {
  const result = await pb.collection('pets').getOne(route.params.id as string, {
    expand: 'species,breed,breed_secondary'
  })
  petToEdit.value = result as any; // assigns content to 'pets'
  console.log(result)

  name.value = petToEdit.value.name
  petType.value = petToEdit.value.expand.species.pet_type
  species.value = petToEdit.value.species
  // breedGroup.value = petToEdit.value.breed_group
  breed.value = petToEdit.value.expand.breed
  // breedSecondary.value = petToEdit.value.breed_secondary
  dob.value = getDateOnly()
  gender.value = petToEdit.value.gender
  isNeutered.value = petToEdit.value.neutered
  colour.value = petToEdit.value.colour
  isImported.value = petToEdit.value.imported
  importLocation.value = petToEdit.value.import_country_code
  microchipNumber.value = petToEdit.value.microchip_number
}

function getDateOnly() {
  if (petToEdit.value.date_of_birth == '') {
    return undefined;
  }
  return petToEdit.value.date_of_birth.split(' ')[0];
}

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
    await fetchPet()
  } catch (e) {
    petNotFound.value = true;
  } finally {
    isLoading.value = false;
  }
})

</script>

<template>
  <main class="container-fluid mt-5">
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

            <section class="mt-4">
              <p>My pet's name is:</p>
              <FloatLabel variant="on">
                <InputText v-model="name" id="petName"/>
                <label for="petName">Pet name *</label>
              </FloatLabel>
              <Message v-if="name === ''" severity="error" size="small" variant="simple">
                Please enter a pet's name to continue
              </Message>
            </section>

            <section class="mt-4">
              <p>
                <PetNameFormatted :name="name"/>
                is a:
              </p>
              <PetTypes v-model="petType"/>


            </section>

            <section class="mt-4">
              <div>
                <p>What breed is
                  <PetNameFormatted :name="name"/>
                  <span>?</span>
                </p>
              </div>
            </section>
          </div>
        </div>
      </div>
    </div>
  </main>
</template>

<style scoped>

</style>