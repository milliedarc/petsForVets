<script setup lang="ts">
import type {Pet} from "@/types/Pet";
import type {Prescription} from "@/types/Prescription";
import {openRequestStatuses, type PrescriptionRequest} from "@/types/PrescriptionRequest";
import {onMounted, ref} from "vue";
import {useRoute} from "vue-router";
import {useToast} from "primevue/usetoast";
import {ApiError} from "@/api/client";
import {getPet} from "@/api/pets";
import {cancelRequest, listPrescriptionRequests, listPrescriptions, requestRepeat} from "@/api/prescriptions";
import PrescriptionCard from "@/views/PrescriptionsView/components/PrescriptionCard.vue";
import PrescriptionRequestCard from "@/views/PrescriptionsView/components/PrescriptionRequestCard.vue";

// *********************** CONSTS ***********************

const route = useRoute()
const toast = useToast()
const petId = route.params.id as string

// *********************** REFS ***********************

const pet = ref<Pet>();
const petNotFound = ref<boolean>(false)
const prescriptions = ref<Prescription[]>([])
const prescriptionRequests = ref<PrescriptionRequest[]>([])

// *********************** FUNCTIONS ***********************

function showError(error: unknown, fallback: string): void {
  console.error(error);
  const detail = error instanceof ApiError ? error.message : fallback;
  toast.add({severity: 'error', summary: 'Error', detail, life: 5000});
}

async function fetchPrescriptionRequests(): Promise<void> {
  prescriptionRequests.value = await listPrescriptionRequests(petId);
}

async function makePrescriptionRequest(prescription: Prescription): Promise<void> {
  try {
    await requestRepeat(prescription.id);
    await fetchPrescriptionRequests();
  } catch (error) {
    showError(error, 'Could not request a repeat prescription');
  }
}

async function deletePrescriptionRequest(prescriptionRequest: PrescriptionRequest): Promise<void> {
  try {
    await cancelRequest(prescriptionRequest.id);
    await fetchPrescriptionRequests();
  } catch (error) {
    showError(error, 'Could not cancel the request');
  }
}

/** A prescription can't be requested again while a request for it is still open. */
function isDisabledRequest(prescription: Prescription): boolean {
  return prescriptionRequests.value.some((prescriptionRequest) => {
    return prescriptionRequest.prescription_id === prescription.id
        && openRequestStatuses.includes(prescriptionRequest.status);
  });
}

// *********************** LIFECYCLE HOOKS ***********************

onMounted(async () => {
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
      showError(error, 'Could not load prescriptions');
    }
  }
})
</script>

<template>
  <main class="container" v-if="petNotFound">
    <section class="bg-body-tertiary p-5 rounded">
      <h4>Pet not found</h4>
    </section>
  </main>
  <main class="container" v-else-if="pet">
    <section class="bg-body-tertiary p-5 rounded">
      <div class="d-flex">
        <h1>{{ pet.name }}'s repeat prescriptions</h1>
      </div>
      <section class="container">
        <div class="row mt-4">
          <div class="col">
            <h4 class="mb-4">Prescriptions</h4>
            <div v-if="prescriptions.length === 0">
              <Card class="p-2 text-center">
                <template #subtitle>{{ pet.name }} has no repeat prescriptions.</template>
              </Card>
            </div>
            <PrescriptionCard
                class="mb-3"
                v-for="prescription in prescriptions"
                :key="prescription.id"
                :prescription="prescription"
                :pet
                :disabled="isDisabledRequest(prescription)"
                @requestButtonClick="makePrescriptionRequest(prescription)"
            />
          </div>
          <div class="col">
            <h4 class="mb-4">Requests</h4>
            <div v-if="prescriptionRequests.length === 0">
              <Card class="p-2 text-center">
                <template #subtitle>No requests yet</template>
              </Card>
            </div>
            <PrescriptionRequestCard
                v-else
                @cancelClicked="deletePrescriptionRequest(prescriptionRequest)"
                class="mb-3"
                v-for="prescriptionRequest in prescriptionRequests"
                :key="prescriptionRequest.id"
                :prescriptionRequest="prescriptionRequest"/>
          </div>
        </div>
      </section>

    </section>
  </main>
</template>

<style scoped>

</style>