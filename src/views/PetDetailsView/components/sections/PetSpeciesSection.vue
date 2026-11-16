<script setup lang="ts">
import type {Species} from "@/types/Pet";
import PetTypeIds from "@/types/PetTypeIds";
import PetNameFormatted from "@/views/PetDetailsView/components/PetNameFormatted.vue";
import {computed} from "vue";

const props = defineProps<{
  petType: string,
  name: string,
  speciesList: Species[]
}>()

const species = defineModel('species', {
  type: String,
  required: true
})

/** Cats and dogs have a single species, so the question is only asked for other pet types. */
const isAsked = computed<boolean>(() => {
  return props.petType === PetTypeIds.smallAnimalId || props.petType === PetTypeIds.reptileId;
})

const filteredSpeciesList = computed<Species[]>(() => {
  return props.speciesList.filter((oneSpecies) => props.petType === oneSpecies.pet_type_id);
})

</script>

<template>
  <section v-if="isAsked" class="question">
    <label for="species" class="question-label">
      What species is
      <PetNameFormatted :name="name"/>?
    </label>
    <Select v-model="species" inputId="species"
            :options="filteredSpeciesList"
            optionLabel="name"
            optionValue="id"
            placeholder="Choose a species"
            class="field-input"/>
  </section>
</template>
