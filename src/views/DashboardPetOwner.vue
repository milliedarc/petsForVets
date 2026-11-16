<script setup lang="ts">
import type {Pet} from "@/types/Pet";
// imports, const, models, props, emits, refs, computed, watchers, functions, hooks
import {listPets} from "@/api/pets";
import {rememberPetName} from "@/stores/petNames";
import {errorMessage} from "@/api/client";
import {onMounted, ref} from "vue";
import DashboardPetCard from "@/components/DashboardPetCard.vue";
import {PawPrint} from "@lucide/vue";

// ****************** REFS *********************

const pets = ref<Pet[]>([]);
const isLoading = ref<boolean>(true)
const loadError = ref<string>('')

// ****************** FUNCTIONS *********************

async function fetchPets(): Promise<void> {
  isLoading.value = true
  loadError.value = ''
  try {
    pets.value = await listPets();
    pets.value.forEach(rememberPetName);
  } catch (error) {
    console.error("Error loading pets:", error);
    loadError.value = errorMessage(error, 'Your pets could not be loaded.');
  } finally {
    isLoading.value = false
  }
}

// ****************** LIFECYCLE HOOKS *********************

onMounted(fetchPets)

</script>

<template>
  <main class="page">
    <div class="page-header">
      <h1>My pets</h1>
      <RouterLink v-if="pets.length > 0" to="/pets/add" custom v-slot="{navigate}">
        <Button label="Add a pet" icon="pi pi-plus" @click="navigate"/>
      </RouterLink>
    </div>

    <div v-if="isLoading" class="state-panel" aria-live="polite">
      <ProgressSpinner style="width: 40px; height: 40px" strokeWidth="5"/>
      <span>Loading your pets…</span>
    </div>

    <div v-else-if="loadError" class="state-panel" role="alert">
      <i class="pi pi-exclamation-triangle state-icon" aria-hidden="true"></i>
      <h2>Something went wrong</h2>
      <p>{{ loadError }}</p>
      <Button label="Try again" icon="pi pi-refresh" severity="secondary" @click="fetchPets"/>
    </div>

    <div v-else-if="pets.length === 0" class="state-panel">
      <PawPrint :size="48" class="state-icon" aria-hidden="true"/>
      <h2>No pets yet</h2>
      <p>Add your first pet to keep their details in one place and request repeat prescriptions.</p>
      <RouterLink to="/pets/add" custom v-slot="{navigate}">
        <Button label="Add a pet" icon="pi pi-plus" @click="navigate"/>
      </RouterLink>
    </div>

    <div v-else class="pet-grid">
      <DashboardPetCard
          v-for="pet in pets"
          :key="pet.id"
          :pet="pet"/>
    </div>
  </main>
</template>

<style scoped>
.pet-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(min(100%, 320px), 1fr));
  gap: 16px;
}

.state-icon {
  font-size: 2rem;
  color: var(--p-primary-color);
}
</style>
