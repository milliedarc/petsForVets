<script setup lang="ts">

import PetNameFormatted from "@/views/PetDetailsView/components/PetNameFormatted.vue";

const props = defineProps<{
  name: string
}>()

const gender = defineModel('gender', {
  type: String,
  default: '',
  required: true
})

const options = [
  {label: 'Male', value: 'male'},
  {label: 'Female', value: 'female'},
]

</script>

<template>
  <section class="question">
    <span id="genderLabel" class="question-label">
      What's
      <PetNameFormatted :name="props.name"/>'s gender?
    </span>
    <SelectButton v-model="gender"
                  :options="options"
                  optionLabel="label"
                  optionValue="value"
                  aria-labelledby="genderLabel"/>
  </section>
</template>
