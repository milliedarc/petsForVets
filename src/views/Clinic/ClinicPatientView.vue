<script setup lang="ts">
import type {ClinicPet, ClinicPrescription, ClinicRequest} from "@/types/Clinic";
import {computed, onMounted, ref} from "vue";
import {useRoute} from "vue-router";
import {getClinicPet, listClinicPetPrescriptions, listClinicRequests} from "@/api/clinic";
import {ApiError, errorMessage} from "@/api/client";
import {formatAge, parseIsoDate} from "@/utils/dates";
import {rememberPetName} from "@/stores/petNames";
import countries from "@/types/countries.json";
import PetIcon from "@/components/PetIcon.vue";
import ClinicRequestList from "@/views/Clinic/components/ClinicRequestList.vue";
import NewPrescriptionDialog from "@/views/Clinic/components/NewPrescriptionDialog.vue";

// ****************** CONSTS *********************

const route = useRoute();
const petId = route.params.id as string;

// ****************** REFS *********************

const pet = ref<ClinicPet>();
const prescriptions = ref<ClinicPrescription[]>([]);
const requests = ref<ClinicRequest[]>([]);
const isLoading = ref<boolean>(true);
const petNotFound = ref<boolean>(false);
const loadError = ref<string>('');
const isPrescribing = ref<boolean>(false);

// ****************** COMPUTED *********************

const ownerName = computed<string>(() => {
  if (!pet.value) {
    return '';
  }
  return [pet.value.owner.name, pet.value.owner.surname].filter(Boolean).join(' ') || pet.value.owner.email;
});

/** The pet's details as label/value pairs, skipping anything the owner left blank. */
const details = computed<{label: string, value: string}[]>(() => {
  const p = pet.value;
  if (!p) {
    return [];
  }
  const rows = [
    {label: 'Species', value: p.species.name},
    {label: 'Breed', value: p.breed?.name ?? ''},
    {
      label: 'Age',
      value: p.date_of_birth
          ? `${formatAge(p.date_of_birth)} (born ${parseIsoDate(p.date_of_birth).toLocaleDateString()})`
          : '',
    },
    {label: 'Sex', value: p.gender ? `${p.gender === 'male' ? 'Male' : 'Female'}${p.neutered ? ', neutered' : ''}` : ''},
    {label: 'Colour', value: p.colour},
    {label: 'Microchip', value: p.microchip_number},
    {label: 'Imported from', value: p.imported ? countryName(p.import_country_code) || 'Yes (country not given)' : ''},
  ];
  return rows.filter((row) => row.value !== '');
});

/** Newest first: the history of everything requested for this pet. */
const requestHistory = computed<ClinicRequest[]>(() => requests.value.toReversed());

// ****************** FUNCTIONS *********************

function countryName(code: string): string {
  return countries.countries.find((country) => country.country_code === code)?.name ?? code;
}

async function loadPage(): Promise<void> {
  isLoading.value = true;
  loadError.value = '';
  try {
    const [loadedPet, loadedPrescriptions, loadedRequests] = await Promise.all([
      getClinicPet(petId),
      listClinicPetPrescriptions(petId),
      listClinicRequests({petId}),
    ]);
    pet.value = loadedPet;
    rememberPetName(loadedPet);
    prescriptions.value = loadedPrescriptions;
    requests.value = loadedRequests;
  } catch (error) {
    if (error instanceof ApiError && error.status === 404) {
      petNotFound.value = true;
    } else {
      console.error(error);
      loadError.value = errorMessage(error, 'This patient could not be loaded.');
    }
  } finally {
    isLoading.value = false;
  }
}

function onPrescriptionCreated(prescription: ClinicPrescription): void {
  prescriptions.value = [prescription, ...prescriptions.value];
}

/** Replace the changed request in place: on this page every status stays listed. */
function onRequestUpdated(updated: ClinicRequest): void {
  requests.value = requests.value.map((request) => request.id === updated.id ? updated : request);
}

// ****************** HOOKS *********************

onMounted(loadPage);

</script>

<template>
  <main class="page">
    <div v-if="isLoading" class="state-panel" aria-live="polite">
      <ProgressSpinner style="width: 40px; height: 40px" strokeWidth="5"/>
      <span>Loading patient…</span>
    </div>

    <div v-else-if="petNotFound" class="state-panel">
      <h2>Patient not found</h2>
      <p>This pet may have been removed by its owner.</p>
      <RouterLink :to="{name: 'ClinicPatients'}">Back to patients</RouterLink>
    </div>

    <div v-else-if="loadError || !pet" class="state-panel" role="alert">
      <p>{{ loadError }}</p>
      <Button label="Try again" icon="pi pi-refresh" severity="secondary" @click="loadPage"/>
    </div>

    <template v-else>
      <div class="page-header">
        <div class="patient-header">
          <PetIcon :pet="pet"/>
          <div>
            <h1>{{ pet.name }}</h1>
            <p class="m-0 text-muted">
              {{ [pet.species.name, pet.breed?.name].filter(Boolean).join(' · ') }}
              <template v-if="pet.date_of_birth"> · {{ formatAge(pet.date_of_birth) }}</template>
            </p>
          </div>
        </div>
        <Button label="New prescription" icon="pi pi-plus" @click="isPrescribing = true"/>
      </div>

      <div class="layout">
        <!-- Left: the pet and owner. -->
        <div class="side">
          <Card>
            <template #title>Owner</template>
            <template #content>
              <div class="owner-name">{{ ownerName }}</div>
              <a :href="`mailto:${pet.owner.email}`" class="owner-email">
                <i class="pi pi-envelope" aria-hidden="true"></i>{{ pet.owner.email }}
              </a>
            </template>
          </Card>

          <Card>
            <template #title>Details</template>
            <template #content>
              <dl class="details">
                <template v-for="row in details" :key="row.label">
                  <dt>{{ row.label }}</dt>
                  <dd>{{ row.value }}</dd>
                </template>
              </dl>
            </template>
          </Card>
        </div>

        <!-- Right: prescriptions and requests. -->
        <div class="main">
          <section aria-labelledby="prescriptionsHeading">
            <h2 id="prescriptionsHeading">Repeat prescriptions</h2>
            <div v-if="prescriptions.length === 0" class="state-panel">
              <p>{{ pet.name }} has no repeat prescriptions yet.</p>
              <Button label="Write the first one" icon="pi pi-plus" severity="secondary" @click="isPrescribing = true"/>
            </div>
            <ul v-else class="prescription-list">
              <li v-for="prescription in prescriptions" :key="prescription.id" class="prescription">
                <div class="prescription-medicine">
                  {{ prescription.medicine.name }}
                  <span class="text-muted text-sm">{{ prescription.medicine.description }}</span>
                </div>
                <div>{{ prescription.dosage }} · Quantity {{ prescription.quantity }}</div>
                <div class="text-muted text-sm">
                  Prescribed {{ new Date(prescription.created_at).toLocaleDateString() }}
                  <template v-if="prescription.prescribed_by"> by {{ prescription.prescribed_by }}</template>
                </div>
              </li>
            </ul>
          </section>

          <section aria-labelledby="requestsHeading">
            <h2 id="requestsHeading">Requests</h2>
            <div v-if="requests.length === 0" class="state-panel">
              <p>{{ ownerName }} hasn't requested any repeats for {{ pet.name }} yet.</p>
            </div>
            <ClinicRequestList v-else
                               :requests="requestHistory"
                               show-status
                               :show-patient="false"
                               @updated="onRequestUpdated"/>
          </section>
        </div>
      </div>

      <NewPrescriptionDialog v-model:visible="isPrescribing"
                             :pet="pet"
                             @created="onPrescriptionCreated"/>
    </template>
  </main>
</template>

<style scoped>
.patient-header {
  display: flex;
  align-items: center;
  gap: 16px;
}

.layout {
  display: grid;
  grid-template-columns: 1fr;
  gap: 24px;
}

@media (min-width: 900px) {
  .layout {
    grid-template-columns: 300px 1fr;
    align-items: start;
  }
}

.side,
.main {
  display: flex;
  flex-direction: column;
  gap: 16px;
  min-width: 0;
}

.main {
  gap: 32px;
}

.owner-name {
  font-weight: 700;
}

.owner-email {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  overflow-wrap: anywhere;
}

.details {
  display: grid;
  grid-template-columns: max-content 1fr;
  gap: 8px 16px;
  margin: 0;
  font-size: 0.925rem;
}

.details dt {
  font-weight: 600;
  color: var(--p-text-muted-color);
}

.details dd {
  margin: 0;
  overflow-wrap: anywhere;
}

.prescription-list {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.prescription {
  padding: 16px 20px;
  background: var(--p-content-background);
  border: 1px solid var(--p-content-border-color);
  border-radius: var(--pfv-radius);
}

.prescription-medicine {
  display: flex;
  flex-wrap: wrap;
  align-items: baseline;
  gap: 4px 10px;
  font-weight: 700;
  font-size: 1.05rem;
}
</style>
