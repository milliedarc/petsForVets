<script setup lang="ts">
import type {PrescriptionRequest, PrescriptionRequestStatus} from "@/types/PrescriptionRequest";
import {computed} from "vue";

const props = defineProps<{
  prescriptionRequest: PrescriptionRequest
}>();

const emit = defineEmits<{
  (e: 'cancelClicked'): void
}>()

/** How each status looks, and what it means for the owner. */
const statusInfo: Record<PrescriptionRequestStatus, { severity: string, icon: string, help: string }> = {
  Pending: {severity: 'warn', icon: 'pi pi-clock', help: 'Waiting for the clinic to review.'},
  Approved: {severity: 'info', icon: 'pi pi-check', help: 'Approved. The clinic is preparing it.'},
  Ready: {severity: 'success', icon: 'pi pi-check-circle', help: 'Ready to collect from the clinic.'},
  Rejected: {severity: 'danger', icon: 'pi pi-times-circle', help: 'The clinic could not approve this request.'},
};

const status = computed(() => statusInfo[props.prescriptionRequest.status]);

const requestedOn = computed<string>(() => {
  // created_at is a full timestamp with timezone, so new Date() reads it correctly.
  const date = new Date(props.prescriptionRequest.created_at);
  if (date.toDateString() === new Date().toDateString()) {
    return 'Requested today';
  }
  return 'Requested ' + date.toLocaleDateString(undefined, {day: 'numeric', month: 'short', year: 'numeric'});
})

</script>

<template>
  <Card>
    <template #title>
      <div class="title-row">
        <span>{{ props.prescriptionRequest.medicine_name }}</span>
        <Tag :value="props.prescriptionRequest.status" :severity="status.severity" :icon="status.icon"/>
      </div>
    </template>
    <template #subtitle>
      {{ requestedOn }}<span v-if="props.prescriptionRequest.dosage"> · {{ props.prescriptionRequest.dosage }}</span>
    </template>
    <template #content>
      <p class="m-0 text-sm">{{ status.help }}</p>
      <Message v-if="props.prescriptionRequest.rejection_reason"
               severity="error" size="small" class="mt-2">
        <span class="font-bold">Reason:</span> {{ props.prescriptionRequest.rejection_reason }}
      </Message>
      <p v-if="props.prescriptionRequest.comments" class="mt-2 mb-0 text-sm text-muted">
        <span class="font-bold">Note:</span> {{ props.prescriptionRequest.comments }}
      </p>
    </template>
    <template v-if="props.prescriptionRequest.status === 'Pending'" #footer>
      <Button @click="emit('cancelClicked')"
              label="Cancel request"
              :aria-label="`Cancel prescription request of ${props.prescriptionRequest.medicine_name}`"
              icon="pi pi-times"
              size="small"
              severity="secondary"
              variant="text"/>
    </template>
  </Card>
</template>

<style scoped>
.title-row {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
}
</style>
