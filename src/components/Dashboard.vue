<script setup>
import PocketBase from 'pocketbase'
import {onMounted, ref} from "vue";

const pb = new PocketBase('http://127.0.0.1:8090');

const pets = ref(); // creates a global vue reactive variable

onMounted(async () => { // global vue variable that executes before rendering component
  const result = await pb.collection('pets').getList(1, 20, {  // calls pets collection with its relations
    expand: 'species,breed'

  });
  console.log(result);
  pets.value = result.items; // assigns content to 'pets'
})

function toUpperCase(gender) {
  return gender.charAt(0).toUpperCase() + gender.slice(1).toLowerCase();
}

function calculateAge(dob) {
  const birthDate = new Date(dob);
  const today = new Date();

  let age = today.getFullYear() - birthDate.getFullYear();
  const monthDiff = today.getMonth() - birthDate.getMonth();
  const dayDiff = today.getDate() - birthDate.getDate();

  // Adjust age if birthday hasn't occurred yet this year
  if (monthDiff < 0 || (monthDiff === 0 && dayDiff < 0)) {
    age--;
  }
  return age;
}
</script>

<template>
  <table class="table table-hover">
    <thead>
    <tr>
      <th scope="col">#</th>
      <th scope="col">Name</th>
      <th scope="col">Species</th>
      <th scope="col">Breed</th>
      <th scope="col">Age</th>
      <th scope="col">Gender</th>
    </tr>
    </thead>
    <tbody>
    <tr
        v-for="pet in pets"
        :key="pet.id">
      <th scope="row">1</th>
      <td>{{ pet.name }}</td>
      <td>{{ pet.expand.species.name }}</td>
      <td>{{ pet.expand.breed.name }}</td>
      <td>{{ calculateAge(pet.date_of_birth) }}</td>
      <td>{{ toUpperCase(pet.gender) }}</td>
    </tr>
    </tbody>
  </table>
</template>

<style scoped>

</style>