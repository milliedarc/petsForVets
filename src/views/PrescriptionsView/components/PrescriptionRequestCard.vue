<script setup lang="ts">
import {computed} from "vue";

const props = defineProps<{
  prescriptionRequest: PrescriptionRequest
}>();

const emit = defineEmits([
  'cancelClicked'
])

const prescriptionMedicine = props.prescriptionRequest.expand.prescription.expand.medicine.name;

const createdAt = computed<string>(() => {
  const date = new Date(props.prescriptionRequest.created.split(' ')[0]);
  const today = new Date();
  if (date.toLocaleDateString() === today.toLocaleDateString()) {
    return 'Requested today'
  }
  return 'Requested on ' + date.toLocaleDateString();
})

const headerBackgroundStyle = computed<string>(() => {
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
        <Button v-if="props.prescriptionRequest.status === 'Pending'"
                @click="emit('cancelClicked')"
                label="Cancel"
                :aria-label="`Cancel prescription request of ${prescriptionMedicine}`"
                icon="pi pi-times"
                size="small"
                severity="contrast"/>
      </div>
      <div v-if="props.prescriptionRequest.rejection_reason != ''" class="px-2 pt-2">
        <Message severity="error" size="small" aria-label="Prescription rejection reason">
          {{ props.prescriptionRequest.rejection_reason }}
        </Message>
      </div>
    </template>
    <template #title>{{ prescriptionMedicine }}</template>
    <template #subtitle> {{ createdAt }}</template>
    <template #footer>

    </template>
  </Card>
</template>

<style scoped>

</style>