<script setup lang="ts">

import PetNameFormatted from "@/views/PetDetailsView/components/PetNameFormatted.vue";

const props = defineProps<{
  name: string
}>()

const colour = defineModel('colour', {
  type: String,
  default: '',
  required: true
})

</script>

<template>
  <section class="question">
    <label for="colour" class="question-label">
      What colour is
      <PetNameFormatted :name="props.name"/>?
    </label>
    <InputText id="colour"
               v-model="colour"
               placeholder="e.g. Black and white"
               class="field-input"/>
  </section>
</template>
