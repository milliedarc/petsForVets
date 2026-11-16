<script setup lang="ts">
import type {Pet} from "@/types/Pet";
import type {Prescription} from "@/types/Prescription";

const props = defineProps<{
  prescription: Prescription,
  pet: Pet,
  /** A request for this prescription is already pending or approved. */
  hasOpenRequest: boolean,
  loading: boolean,
}>();

const emit = defineEmits<{
  (e: 'requestButtonClick'): void
}>()

</script>

<template>
  <Card>
    <template #title>{{ props.prescription.medicine.name }}</template>
    <template #subtitle>
      {{ props.prescription.dosage }}
      <span v-if="props.prescription.quantity"> · Quantity {{ props.prescription.quantity }}</span>
    </template>
    <template v-if="props.prescription.medicine.description" #content>
      <p class="m-0 text-sm text-muted">{{ props.prescription.medicine.description }}</p>
    </template>
    <template #footer>
      <Button v-if="!props.hasOpenRequest"
              @click="emit('requestButtonClick')"
              label="Request a repeat"
              icon="pi pi-refresh"
              :loading="props.loading"
              :aria-label="`Request repeat prescription of ${props.prescription.medicine.name} for ${pet.name}`"
              class="w-full"/>
      <div v-else class="already-requested">
        <i class="pi pi-clock" aria-hidden="true"></i>
        Repeat already requested
      </div>
    </template>
  </Card>
</template>

<style scoped>
.already-requested {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 8px;
  border-radius: var(--p-border-radius-md);
  background: var(--p-content-hover-background);
  color: var(--p-text-muted-color);
  font-weight: 600;
  font-size: 0.9rem;
}
</style>
