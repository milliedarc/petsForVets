<script setup lang="ts">
import type {Pet} from "@/types/Pet";
import type {Prescription} from "@/types/Prescription";
import {openRequestStatuses, type PrescriptionRequest} from "@/types/PrescriptionRequest";
import {onMounted, ref} from "vue";
import {useRoute} from "vue-router";
import {useToast} from "primevue/usetoast";
import {useConfirm} from "primevue/useconfirm";
import {ApiError, errorMessage} from "@/api/client";
import {getPet} from "@/api/pets";
import {cancelRequest, listPrescriptionRequests, listPrescriptions, requestRepeat} from "@/api/prescriptions";
import PrescriptionCard from "@/views/PrescriptionsView/components/PrescriptionCard.vue";
import PrescriptionRequestCard from "@/views/PrescriptionsView/components/PrescriptionRequestCard.vue";

// *********************** CONSTS ***********************

const route = useRoute()
const toast = useToast()
const confirm = useConfirm()
const petId = route.params.id as string

// *********************** REFS ***********************

const pet = ref<Pet>();
const isLoading = ref<boolean>(true)
const petNotFound = ref<boolean>(false)
const loadError = ref<string>('')
const prescriptions = ref<Prescription[]>([])
const prescriptionRequests = ref<PrescriptionRequest[]>([])
/** Id of the prescription whose request button is busy, so only that button spins. */
const requestingId = ref<string>('')

// *********************** FUNCTIONS ***********************

function showError(error: unknown, fallback: string): void {
  console.error(error);
  toast.add({severity: 'error', summary: 'Error', detail: errorMessage(error, fallback), life: 5000});
}

async function fetchPrescriptionRequests(): Promise<void> {
  prescriptionRequests.value = await listPrescriptionRequests(petId);
}

async function loadPage(): Promise<void> {
  isLoading.value = true;
  loadError.value = '';
  try {
    const [loadedPet, loadedPrescriptions] = await Promise.all([
      getPet(petId),
      listPrescriptions(petId),
      fetchPrescriptionRequests(),
    ]);
    pet.value = loadedPet;
    prescriptions.value = loadedPrescriptions;
  } catch (error) {
    if (error instanceof ApiError && error.status === 404) {
      petNotFound.value = true;
    } else {
      console.error(error);
      loadError.value = errorMessage(error, 'Prescriptions could not be loaded.');
    }
  } finally {
    isLoading.value = false;
  }
}

async function makePrescriptionRequest(prescription: Prescription): Promise<void> {
  requestingId.value = prescription.id;
  try {
    await requestRepeat(prescription.id);
    await fetchPrescriptionRequests();
    toast.add({
      severity: 'success',
      summary: 'Request sent',
      detail: `We've asked the clinic for a repeat of ${prescription.medicine.name}.`,
      life: 4000,
    });
  } catch (error) {
    showError(error, 'Could not request a repeat prescription');
  } finally {
    requestingId.value = '';
  }
}

function confirmCancelRequest(prescriptionRequest: PrescriptionRequest): void {
  confirm.require({
    header: 'Cancel request',
    message: `Cancel your request for ${prescriptionRequest.medicine_name}?`,
    icon: 'pi pi-question-circle',
    rejectProps: {label: 'Keep request', severity: 'secondary', outlined: true},
    acceptProps: {label: 'Cancel request', severity: 'danger'},
    defaultFocus: 'reject',
    accept: async () => {
      try {
        await cancelRequest(prescriptionRequest.id);
        await fetchPrescriptionRequests();
        toast.add({severity: 'info', summary: 'Request cancelled', life: 3000});
      } catch (error) {
        showError(error, 'Could not cancel the request');
      }
    },
  });
}

/** A prescription can't be requested again while a request for it is still open. */
function hasOpenRequest(prescription: Prescription): boolean {
  return prescriptionRequests.value.some((prescriptionRequest) => {
    return prescriptionRequest.prescription_id === prescription.id
        && openRequestStatuses.includes(prescriptionRequest.status);
  });
}

// *********************** LIFECYCLE HOOKS ***********************

onMounted(loadPage)
</script>

<template>
  <main class="page">
    <RouterLink to="/" class="back-link">
      <i class="pi pi-arrow-left" aria-hidden="true"></i>My pets
    </RouterLink>

    <div v-if="isLoading" class="state-panel" aria-live="polite">
      <ProgressSpinner style="width: 40px; height: 40px" strokeWidth="5"/>
      <span>Loading prescriptions…</span>
    </div>

    <div v-else-if="petNotFound" class="state-panel">
      <h2>Pet not found</h2>
      <p>This pet doesn't exist, or it isn't on your account.</p>
      <RouterLink to="/">Back to my pets</RouterLink>
    </div>

    <div v-else-if="loadError || !pet" class="state-panel" role="alert">
      <h2>Something went wrong</h2>
      <p>{{ loadError }}</p>
      <Button label="Try again" icon="pi pi-refresh" severity="secondary" @click="loadPage"/>
    </div>

    <template v-else>
      <div class="page-header">
        <h1>{{ pet.name }}'s repeat prescriptions</h1>
      </div>

      <div class="columns">
        <section aria-labelledby="prescriptionsHeading">
          <h2 id="prescriptionsHeading">Prescriptions</h2>
          <p class="text-muted text-sm">Medicines your vet has prescribed. Request a repeat when you're running low.</p>
          <div v-if="prescriptions.length === 0" class="state-panel">
            <p>{{ pet.name }} has no repeat prescriptions yet. Your vet adds them after an appointment.</p>
          </div>
          <div class="card-list">
            <PrescriptionCard
                v-for="prescription in prescriptions"
                :key="prescription.id"
                :prescription="prescription"
                :pet
                :hasOpenRequest="hasOpenRequest(prescription)"
                :loading="requestingId === prescription.id"
                @requestButtonClick="makePrescriptionRequest(prescription)"
            />
          </div>
        </section>

        <section aria-labelledby="requestsHeading">
          <h2 id="requestsHeading">Your requests</h2>
          <p class="text-muted text-sm">Track each request from <em>Pending</em> to <em>Ready</em> for collection.</p>
          <div v-if="prescriptionRequests.length === 0" class="state-panel">
            <p>No requests yet.</p>
          </div>
          <div class="card-list" aria-live="polite">
            <PrescriptionRequestCard
                v-for="prescriptionRequest in prescriptionRequests"
                :key="prescriptionRequest.id"
                :prescriptionRequest="prescriptionRequest"
                @cancelClicked="confirmCancelRequest(prescriptionRequest)"/>
          </div>
        </section>
      </div>
    </template>
  </main>
</template>

<style scoped>
.columns {
  display: grid;
  grid-template-columns: 1fr;
  gap: 32px;
}

@media (min-width: 768px) {
  .columns {
    grid-template-columns: 1fr 1fr;
  }
}

.card-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
</style>
