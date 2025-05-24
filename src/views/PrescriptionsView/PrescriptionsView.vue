<script setup lang="ts">
import {computed, onMounted, ref} from "vue";
import PocketBase from 'pocketbase'
import {useRoute} from "vue-router";
import PrescriptionCard from "@/views/PrescriptionsView/components/PrescriptionCard.vue";
import PrescriptionRequestCard from "@/views/PrescriptionsView/components/PrescriptionRequestCard.vue";

// *********************** CONSTS ***********************
const pb = new PocketBase('http://127.0.0.1:8090');
const route = useRoute()

// *********************** REFS ***********************

const pet = ref<Pet>();
const prescriptions = ref<Prescription[]>([])
const prescriptionRequests = ref<PrescriptionRequest[]>([])


// *********************** COMPUTED ***********************


// *********************** FUNCTIONS ***********************

async function fetchPetId() {
  pet.value = await pb.collection('pets').getOne(route.params.id as string, {
    expand: 'species,breed'
  })
}

async function fetchPrescriptionRequests() {
  const result = await pb.collection('prescription_requests').getList(1, 100, {
    expand: 'pet,prescription,prescription.medicine',
    filter: pb.filter("pet = {:petId}", {petId: route.params.id})
  })
  prescriptionRequests.value = result.items as any;
  console.log('prescriptionRequests: ', result)
}

async function makePrescriptionRequest(prescription: Prescription) {
  await pb.collection('prescription_requests').create({
    pet: pet.value.id,
    prescription: prescription.id,
    status: 'Pending',
    // comments: ''
  });
  await fetchPrescriptionRequests()
}

async function deletePrescriptionRequest(prescriptionRequest: PrescriptionRequest) {
  await pb.collection('prescription_requests').delete(prescriptionRequest.id);
  await fetchPrescriptionRequests()
}

function isDisabledRequest(prescription: Prescription): boolean {
  const foundPrescriptionRequest = prescriptionRequests.value.find((prescriptionRequest) => {
    return prescriptionRequest.prescription === prescription.id;
  })
  if (foundPrescriptionRequest) {
    return true;
  }
  return false;
}

// *********************** LIFECYCLE HOOKS ***********************

onMounted(async () => {
  await fetchPetId()

  const result = await pb.collection('prescriptions').getList(1, 100, {
    expand: 'pet,medicine',
    filter: pb.filter("pet = {:petId}", {petId: route.params.id})
  })
  prescriptions.value = result.items as any;
  console.log('prescriptions: ', result)

  await fetchPrescriptionRequests()
})
</script>

<template>
  <main class="container" v-if="pet">
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