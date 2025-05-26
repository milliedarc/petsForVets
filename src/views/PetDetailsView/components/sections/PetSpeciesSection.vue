<script setup lang="ts">
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

const filteredSpeciesList = computed(() => {
  if (props.speciesList === undefined) {
    return []
  }
  return props.speciesList.filter((species) => {
    if (props.petType === species.pet_type) {
      return true;
    }
    return false;
  })
})

const isValidSelectedSpecies = computed(() => {
  for (const oneSpecies of filteredSpeciesList.value) {
    if (oneSpecies.id === species.value) {
      return true
    }
  }
  return false
})

</script>

<template>
  <section>
    <div v-if="props.petType === PetTypeIds.smallAnimalId || petType === PetTypeIds.reptileId">
      <p>What species is
        <PetNameFormatted :name="name"/>
        <span>?</span>
      </p>
      <div class="mt-4">

        <FloatLabel variant="on">
          <Select v-model="species" inputId="species"
                  :options="filteredSpeciesList"
                  optionLabel="name"
                  optionValue="id"
                  editable
                  class="myInput"
                  :invalid="!isValidSelectedSpecies"
          />
          <label for="species">Species</label>
        </FloatLabel>
        <Message
            v-if="!isValidSelectedSpecies" severity="error" size="small" variant="simple"
            aria-describedby="speciesError">
          Please choose a valid species from the list
        </Message>
      </div>
    </div>

  </section>
</template>

<style scoped>

</style>