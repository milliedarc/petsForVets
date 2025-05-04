<script setup lang="ts">
import PocketBase from 'pocketbase'
import {inject, nextTick, onMounted, ref} from "vue";
import PetModal from "@/components/PetModal.vue";
import DeletePetConfirmModal from "@/components/DeletePetConfirmModal.vue";
import {useRouter} from "vue-router";

const bootstrap = inject('bootstrap');

const pb = new PocketBase('http://127.0.0.1:8090');

const router = useRouter()

const pets = ref<Pet[]>(); // creates a global vue reactive variable
const petToEdit = ref<Pet>()

function toUpperCase(gender: string): string {
  return gender.charAt(0).toUpperCase() + gender.slice(1).toLowerCase();
}

function calculateAge(dob: string): string {
  if (dob === '') {
    return ''
  }

  const birthDate = new Date(dob);
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

  return `${years} years ${months} months`;
}

function displayBreed(pet) {
  // console.log('Pet details:', pet);
  if (pet.expand.breed === undefined) {
    return '';
  }
  if (pet.expand.breed_secondary) {
    return pet.expand.breed.name + ' x ' + pet.expand.breed_secondary.name;
  }
  return pet.expand.breed.name;
}

async function fetchPets() {
  const result = await pb.collection('pets').getList(1, 20, {  // calls pets collection with its relations
    expand: 'species,breed'
  })
  // @ts-ignore
  pets.value = result.items; // assigns content to 'pets'
}

function goToPetView(pet: Pet) {
  router.push(`/pets/${pet.id}`);
}

function goToPetViewAdd() {
  router.push('/pets/add');
}

onMounted(async () => {
  await fetchPets()
})


</script>

<template>
  <main class="container">
    <div class="bg-body-tertiary p-5 rounded">
      <h1 class="mb-4">My pets</h1>
      <div v-for="pet in pets"
           :key="pet.id"
           class="card mb-3">
        <div class="card-body">
          <div class="d-flex justify-content-between align-items-center">
            <div class="d-flex">
              <div>
                <PetIcon :pet="pet"/>
              </div>
              <div class="ms-4">
                <h5>{{ pet.name }}</h5>
                <div>
                  <div class="mb-2">{{ calculateAge(pet.date_of_birth) }}</div>
                </div>
              </div>
            </div>
            <div>
              <button @click="goToPetView(pet)" class="btn btn-link" style="color: black"><i
                  class="pi pi-pencil me-2"></i><span class="fw-bold">Edit</span>
              </button>
              <!--              <button @click="goToPetView(pet)" class="btn btn-secondary ms-2">-->
              <!--                View-->
              <!--              </button>-->
              <!--              <button @click="petToEdit=pet" class="btn btn-danger ms-2"-->
              <!--                      data-bs-toggle="modal" data-bs-target="#deletePetModal">-->
              <!--                Delete-->
              <!--              </button>-->
            </div>
          </div>
        </div>
      </div>
      <div>
        <Button
            class="w-100"
            @click="goToPetViewAdd">
          Add new pet
          <i class="pi pi-plus"></i>
        </Button>
      </div>
    </div>
  </main>

</template>

<style scoped>

</style>