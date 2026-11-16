<script setup lang="ts">
import type {ClinicPet} from "@/types/Clinic";
import {onMounted, ref, watch} from "vue";
import {useRoute, useRouter} from "vue-router";
import {listClinicPets} from "@/api/clinic";
import {errorMessage} from "@/api/client";
import {formatAge} from "@/utils/dates";
import {rememberPetName} from "@/stores/petNames";
import PetIcon from "@/components/PetIcon.vue";

// ****************** CONSTS *********************

const route = useRoute();
const router = useRouter();

// ****************** REFS *********************

/** Starts from ?search= so a search survives going into a patient and back. */
const search = ref<string>((route.query.search as string) ?? '');
const pets = ref<ClinicPet[]>([]);
const isLoading = ref<boolean>(true);
const loadError = ref<string>('');

/** Only the latest search's results are shown, even if an older one answers last. */
let latestSearch = 0;
let debounceTimer: ReturnType<typeof setTimeout> | undefined;

// ****************** FUNCTIONS *********************

async function fetchPets(): Promise<void> {
  const thisSearch = ++latestSearch;
  isLoading.value = true;
  loadError.value = '';
  try {
    const results = await listClinicPets(search.value);
    if (thisSearch === latestSearch) {
      pets.value = results;
      results.forEach(rememberPetName);
    }
  } catch (error) {
    console.error(error);
    if (thisSearch === latestSearch) {
      loadError.value = errorMessage(error, 'Patients could not be loaded.');
    }
  } finally {
    if (thisSearch === latestSearch) {
      isLoading.value = false;
    }
  }
}

function ownerName(pet: ClinicPet): string {
  return [pet.owner.name, pet.owner.surname].filter(Boolean).join(' ') || pet.owner.email;
}

function speciesAndBreed(pet: ClinicPet): string {
  return [pet.species.name, pet.breed?.name].filter(Boolean).join(' · ');
}

// ****************** WATCHERS / HOOKS *********************

// Search as you type, but wait for a pause so we don't send a request per keystroke.
watch(search, (term) => {
  clearTimeout(debounceTimer);
  debounceTimer = setTimeout(() => {
    router.replace({query: term.trim() ? {search: term.trim()} : {}});
    fetchPets();
  }, 300);
});

onMounted(fetchPets);

</script>

<template>
  <main class="page">
    <div class="page-header">
      <h1>Patients</h1>
    </div>

    <IconField class="search">
      <InputIcon class="pi pi-search"/>
      <InputText v-model="search"
                 type="search"
                 placeholder="Search by pet, owner, email or microchip"
                 aria-label="Search patients"
                 autofocus
                 fluid/>
    </IconField>

    <div v-if="loadError" class="state-panel" role="alert">
      <p>{{ loadError }}</p>
      <Button label="Try again" icon="pi pi-refresh" severity="secondary" @click="fetchPets"/>
    </div>

    <div v-else-if="isLoading && pets.length === 0" class="state-panel" aria-live="polite">
      <ProgressSpinner style="width: 40px; height: 40px" strokeWidth="5"/>
      <span>Loading patients…</span>
    </div>

    <div v-else-if="pets.length === 0" class="state-panel">
      <i class="pi pi-search empty-icon" aria-hidden="true"></i>
      <p v-if="search.trim()">No patients match “{{ search.trim() }}”.</p>
      <p v-else>No patients yet. Pets appear here once their owners add them.</p>
    </div>

    <template v-else>
      <p class="result-count text-muted text-sm" aria-live="polite">
        {{ pets.length }} {{ pets.length === 1 ? 'patient' : 'patients' }}<template v-if="pets.length === 100"> (showing the first 100, narrow your search)</template>
      </p>
      <ul class="patient-list" :class="{loading: isLoading}">
        <li v-for="pet in pets" :key="pet.id">
          <RouterLink :to="{name: 'ClinicPatient', params: {id: pet.id}}" class="patient">
            <PetIcon :pet="pet"/>
            <div class="patient-text">
              <div class="patient-name">{{ pet.name }}</div>
              <div class="patient-species">{{ speciesAndBreed(pet) }}<template v-if="pet.date_of_birth"> · {{ formatAge(pet.date_of_birth) }}</template></div>
              <div class="patient-owner">
                <i class="pi pi-user" aria-hidden="true"></i>
                {{ ownerName(pet) }} <span class="text-muted">· {{ pet.owner.email }}</span>
              </div>
            </div>
            <Tag v-if="pet.open_requests > 0"
                 :value="`${pet.open_requests} open ${pet.open_requests === 1 ? 'request' : 'requests'}`"
                 severity="warn"
                 class="patient-tag"/>
            <i class="pi pi-angle-right chevron" aria-hidden="true"></i>
          </RouterLink>
        </li>
      </ul>
    </template>
  </main>
</template>

<style scoped>
.search {
  max-width: 560px;
  margin-bottom: 20px;
}

.result-count {
  margin-bottom: 8px;
}

.patient-list {
  list-style: none;
  margin: 0;
  padding: 0;
  background: var(--p-content-background);
  border: 1px solid var(--p-content-border-color);
  border-radius: var(--pfv-radius);
  overflow: hidden;
  transition: opacity 0.15s;
}

.patient-list.loading {
  opacity: 0.6;
}

.patient-list li + li {
  border-top: 1px solid var(--p-content-border-color);
}

.patient {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 14px 16px;
  color: var(--p-text-color);
}

.patient:hover {
  background: var(--p-content-hover-background);
}

.patient:focus-visible {
  outline: 2px solid var(--p-primary-color);
  outline-offset: -2px;
}

.patient-text {
  flex: 1;
  min-width: 0;
}

.patient-name {
  font-weight: 700;
  font-size: 1.05rem;
}

.patient-species {
  font-size: 0.9rem;
  font-weight: 600;
  color: var(--p-primary-color);
}

.patient-owner {
  font-size: 0.875rem;
  overflow-wrap: anywhere;
}

.patient-owner i {
  font-size: 0.75rem;
  margin-right: 4px;
  color: var(--p-text-muted-color);
}

.chevron {
  color: var(--p-text-muted-color);
}

.empty-icon {
  font-size: 1.75rem;
  color: var(--p-primary-color);
}

@media (max-width: 575px) {
  .patient-tag {
    display: none;
  }
}
</style>
