<script setup lang="ts">
import PocketBase from 'pocketbase'
import {onMounted, ref} from "vue";
import AddPet from "@/components/AddPet.vue";

const pb = new PocketBase('http://127.0.0.1:8090');

const pets = ref(); // creates a global vue reactive variable

onMounted(async () => {
  const result = await pb.collection('pets').getList(1, 20, {  // calls pets collection with its relations
    expand: 'species,breed,breed_secondary'
  })
  pets.value = result.items; // assigns content to 'pets'
})

function toUpperCase(gender) {
  return gender.charAt(0).toUpperCase() + gender.slice(1).toLowerCase();
}

function calculateAge(dob) {
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
  console.log('Pet details:', pet);
  if (pet.expand.breed_secondary) {
    return pet.expand.breed.name + ' x ' + pet.expand.breed_secondary.name;
  }
  return pet.expand.breed.name;
}

</script>

<template>
  <h1>My pets</h1>
  <AddPet/>
  <button type="button" class="btn btn-primary" data-bs-toggle="modal" data-bs-target="#addPet">
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
    </tr>
    </thead>
    <tbody>
    <tr
        v-for="pet in pets"
        :key="pet.id">
      <td></td>
      <td>{{ pet.name }}</td>
      <td>{{ pet.expand.species.name }}</td>
      <td>{{ displayBreed(pet) }}</td>
      <td>{{ calculateAge(pet.date_of_birth) }}</td>
      <td>{{ toUpperCase(pet.gender) }}</td>
      <td>
        <button class="btn btn-primary">Edit details</button>
      </td>
    </tr>
    </tbody>
  </table>
</template>

<style scoped>

</style>