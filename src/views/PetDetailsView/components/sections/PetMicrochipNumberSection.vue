<script setup lang="ts">

import PetNameFormatted from "@/views/PetDetailsView/components/PetNameFormatted.vue";

const props = defineProps<{
  name: string
}>()

const microchipNumber = defineModel('microchipNumber', {
  type: String,
  default: '',
  required: true
})

</script>

<template>
  <section class="question">
    <label for="microchip" class="question-label">
      What's
      <PetNameFormatted :name="props.name"/>'s microchip number?
    </label>
    <InputText id="microchip"
               v-model="microchipNumber"
               inputmode="numeric"
               autocomplete="off"
               class="field-input"
               aria-describedby="microchipHint"/>
    <small id="microchipHint" class="field-hint">
      Usually 15 digits. You'll find it on your pet's vaccination card or passport.
    </small>
  </section>
</template>
