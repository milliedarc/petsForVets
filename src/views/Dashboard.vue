<script setup lang="ts">
// imports, const, models, props, emits, refs, computed, watchers, functions, hooks

import PocketBase from 'pocketbase'
import {onMounted, ref} from "vue";
import {useRouter} from "vue-router";
import DashboardPetCard from "@/components/DashboardPetCard.vue";
import PrescriptionsView from "@/views/PrescriptionsView/PrescriptionsView.vue";

const pb = new PocketBase('http://127.0.0.1:8090');

const router = useRouter()

const pets = ref<Pet[]>(); // creates a global vue reactive variable

async function fetchPets() {
  const result = await pb.collection('pets').getList(1, 20, {  // calls pets collection with its relations
    expand: 'species,breed'
  })
  // @ts-ignore
  pets.value = result.items; // assigns content to 'pets'
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
      <DashboardPetCard
          v-for="pet in pets"
          :key="pet.id"
          :pet="pet"/>
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