<script setup lang="ts">
import PocketBase from 'pocketbase'
import {inject, nextTick, onMounted, ref} from "vue";
import PetModal from "@/components/PetModal.vue";

const bootstrap = inject('bootstrap');

const pb = new PocketBase('http://127.0.0.1:8090');

const pets = ref<Pet[]>(); // creates a global vue reactive variable
const petToEdit = ref<Pet>()

function openModal(pet: Pet | undefined) {
  petToEdit.value = pet;

  nextTick(() => {
    // @ts-ignore
    const addPetModal = bootstrap.Modal.getOrCreateInstance(document.getElementById('petModal'));
    addPetModal.show()
  })
}

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

  return `${years}y ${months}m`;
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
    expand: 'species,breed,breed_secondary'
  })
  // @ts-ignore
  pets.value = result.items; // assigns content to 'pets'
}

onMounted(async () => {
  await fetchPets()
})

function logOut() {
  pb.authStore.clear()
  location.reload()
}

</script>

<template>

  <nav class="navbar navbar-expand-md navbar-dark bg-dark mb-4">
    <div class="container-fluid">
      <a class="navbar-brand" href="#">Pets for Vets</a>
      <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarCollapse"
              aria-controls="navbarCollapse" aria-expanded="false" aria-label="Toggle navigation">
        <span class="navbar-toggler-icon"></span>
      </button>
      <div class="collapse navbar-collapse" id="navbarCollapse">
        <ul class="navbar-nav me-auto mb-2 mb-md-0">
          <li class="nav-item">
            <a class="nav-link active" aria-current="page" href="#">Home</a>
          </li>
          <li class="nav-item">
            <a class="nav-link" href="#">Link</a>
          </li>
          <li class="nav-item">
            <a class="nav-link disabled" aria-disabled="true">Disabled</a>
          </li>
        </ul>
        <div class="me-4 text-white">
          {{ pb.authStore.record.name }}
        </div>
        <!--        <form class="d-flex" role="search">-->
        <!--          <input class="form-control me-2" type="search" placeholder="Search" aria-label="Search">-->
        <!--          <button class="btn btn-outline-success" type="submit">Search</button>-->
        <!--        </form>-->
        <button class="btn btn-outline-light" @click="logOut">Log out</button>
      </div>
    </div>
  </nav>

  <main class="container">
    <div class="bg-body-tertiary p-5 rounded">
      <h1>My pets</h1>
      <PetModal v-model="petToEdit" @savePet="fetchPets"/>
      <button @click="openModal(undefined)" type="button" class="btn btn-primary ">
        Add new pet
      </button>
      <table class="table table-hover">
        <thead>
        <tr>
          <th scope="col">#</th>
          <th scope="col">Name</th>
          <th scope="col">Species</th>
          <th scope="col">Breed</th>
          <th scope="col">Age</th>
          <th scope="col">Gender</th>
          <th scope="col"></th>
          <th scope="col"></th>
        </tr>
        </thead>
        <tbody>
        <tr
            v-for="pet in pets"
            :key="pet.id">
          <td></td>
          <td>{{ pet.name }}</td>
          <td>{{ pet.expand.species.name }}</td>
          <td>
        <span v-if="displayBreed(pet) === ''"
              class="text-muted">
          <small>Unknown</small>
        </span>
            <span v-else>{{ displayBreed(pet) }}</span>
          </td>
          <td>{{ calculateAge(pet.date_of_birth) }}</td>
          <td>{{ toUpperCase(pet.gender) }}</td>
          <td>
            <button @click="openModal(pet)"
                    class="btn btn-primary">
              Edit details
            </button>
          </td>
          <td>
            <button class="btn btn-primary">View</button>
          </td>
        </tr>
        </tbody>
      </table>
    </div>
  </main>

</template>

<style scoped>

</style>