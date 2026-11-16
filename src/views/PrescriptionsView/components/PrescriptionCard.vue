<script setup lang="ts">
import type {Pet} from "@/types/Pet";
import type {Prescription} from "@/types/Prescription";

const props = defineProps<{
  prescription: Prescription,
  pet: Pet,
  disabled: boolean,
}>();

const emit = defineEmits([
  'requestButtonClick'
])

</script>

<template>
  <Card>
    <template #title> {{ props.prescription.expand.medicine.name }}</template>
    <template #subtitle> {{ props.prescription.dosage }}</template>
    <template #footer>
      <div class="d-flex gap-4 mt-1">
        <Button @click="emit('requestButtonClick')"
                label="Request"
                :aria-label="`Request repeat prescription of ${props.prescription.expand.medicine.name} for ${pet.name}`"
                class="w-100"
                :disabled="props.disabled"/>
      </div>
    </template>
  </Card>
</template>

<style scoped>

</style>