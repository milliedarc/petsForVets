<script setup lang="ts">
import {computed, onMounted, ref} from "vue";
import PocketBase from 'pocketbase'
import {useRoute} from "vue-router";
import PrescriptionCard from "@/components/Prescriptions/PrescriptionCard.vue";
import PrescriptionRequestCard from "@/components/Prescriptions/PrescriptionRequestCard.vue";

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

// *********************** LIFECYCLE HOOKS ***********************

onMounted(async () => {
  await fetchPetId()
  const result = await pb.collection('prescriptions').getList(1, 100, {
    expand: 'pet,medicine',
    filter: pb.filter("pet = {:petId}", {petId: route.params.id})
  })
  prescriptions.value = result.items as any;
  console.log('prescriptions: ', result)

  const result2 = await pb.collection('prescription_requests').getList(1, 100, {
    expand: 'pet,prescription,prescription.medicine',
    filter: pb.filter("pet = {:petId}", {petId: route.params.id})
  })
  prescriptionRequests.value = result2.items as any;
  console.log('prescriptionRequests: ', result2)
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
            <PrescriptionCard
                class="mb-3"
                v-for="prescription in prescriptions"
                :key="prescription.id"
                :prescription="prescription"
                :pet>
            </PrescriptionCard>
          </div>
          <div class="col">
            <h4 class="mb-4">Requests</h4>
            <PrescriptionRequestCard
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