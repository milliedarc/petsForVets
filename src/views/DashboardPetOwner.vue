<script setup lang="ts">
import type {Pet} from "@/types/Pet";
// imports, const, models, props, emits, refs, computed, watchers, functions, hooks
import {listPets} from "@/api/pets";
import {onMounted, ref} from "vue";
import {useRouter} from "vue-router";
import DashboardPetCard from "@/components/DashboardPetCard.vue";

const router = useRouter()

const pets = ref<Pet[]>();
const isLoading = ref<boolean>(true)

async function fetchPets(): Promise<void> {
  pets.value = await listPets();
}

function goToPetViewAdd() {
  router.push('/pets/add');
}

onMounted(async () => {
  isLoading.value = true
  await fetchPets()
  isLoading.value = false
})

</script>

<template>
  <main class="container">
    <div class="bg-body-tertiary p-5 rounded">
      <h1 class="mb-4">My pets</h1>
      <div v-if="isLoading" class="d-flex justify-content-center">
        <ProgressSpinner/>
      </div>
      <div v-else>
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
    </div>
  </main>

</template>

<style scoped>

</style>