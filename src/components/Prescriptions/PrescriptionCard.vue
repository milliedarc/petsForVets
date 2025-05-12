<script setup lang="ts">

const props = defineProps<{
  prescription: Prescription,
  pet: Pet
}>();

</script>

<template>
  <Card>
    <template #title> {{ props.prescription.expand.medicine.name }}</template>
    <template #subtitle> {{ props.prescription.dosage }}</template>
    <template #footer>
      <div class="d-flex gap-4 mt-1">
        <Button label="Request" class="w-full"/>
      </div>
    </template>
  </Card>
</template>

<style scoped>

</style>