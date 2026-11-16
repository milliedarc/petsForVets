<script setup lang="ts">

import PetNameFormatted from "@/views/PetDetailsView/components/PetNameFormatted.vue";

const props = defineProps<{
  name: string
}>()

const isNeutered = defineModel('isNeutered', {
  type: Boolean,
  required: true
})

</script>

<template>
  <section class="question">
    <span class="question-label">
      Is
      <PetNameFormatted :name="props.name"/>
      neutered?
    </span>
    <div class="flex items-center gap-3">
      <ToggleSwitch v-model="isNeutered" inputId="neutered"/>
      <label for="neutered">{{ isNeutered ? 'Yes, neutered' : 'No, not neutered' }}</label>
    </div>
  </section>
</template>
