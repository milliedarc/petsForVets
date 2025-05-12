<script setup lang="ts">
import {computed} from "vue";

const props = defineProps<{
  prescriptionRequest: PrescriptionRequest
}>();

const emit = defineEmits([
  'cancelClicked'
])

const createdAt = computed(() => {
  return props.prescriptionRequest.created.split(' ')[0];
})

const headerBackgroundStyle = computed(() => {
  let colour = ''
  if (props.prescriptionRequest.status === 'Pending') {
    colour = '#f6c7a3'
  } else if (props.prescriptionRequest.status === 'Approved') {
    colour = '#fff3c3'
  } else if (props.prescriptionRequest.status === 'Ready') {
    colour = '#d5eccb'
  } else if (props.prescriptionRequest.status === 'Rejected') {
    colour = '#f3acac'
  }
  return `background-color: ${colour}`;
})
</script>

<template>
  <Card>
    <template #header>
      <div class="w-100 p-3 d-flex justify-content-between align-items-center"
           :style="headerBackgroundStyle">
        <div>
          <h6 class="mb-0 fw-bold">{{ props.prescriptionRequest.status }}</h6>
        </div>
        <Button @click="emit('cancelClicked')" label="Cancel" icon="pi pi-times" size="small" severity="contrast"/>
      </div>
    </template>
    <template #title>{{ props.prescriptionRequest.expand.prescription.expand.medicine.name }}</template>
    <template #subtitle>{{ createdAt }}</template>
    <template #footer>
    </template>
  </Card>
</template>

<style scoped>

</style>