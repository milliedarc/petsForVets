<script setup lang="ts">
import type {Breed} from "@/types/Pet";
import PetNameFormatted from "@/views/PetDetailsView/components/PetNameFormatted.vue";

const props = defineProps<{
  name: string,
  breedsList: Breed[]   // already filtered to the pet's species
}>()

const breedId = defineModel('breedId', {
  type: String,
  required: true
})

</script>

<template>
  <section class="question">
    <label for="breed" class="question-label">
      What breed is
      <PetNameFormatted :name="props.name"/>?
    </label>
    <Select v-model="breedId" inputId="breed"
            :options="props.breedsList"
            optionLabel="name"
            optionValue="id"
            placeholder="Choose a breed"
            filter
            filterPlaceholder="Search breeds"
            resetFilterOnHide
            showClear
            class="field-input"
            aria-describedby="breedHint"/>
    <small id="breedHint" class="field-hint">
      Not sure? Choose <span class="font-bold">Cross / Mixed breed</span>.
    </small>
  </section>
</template>
