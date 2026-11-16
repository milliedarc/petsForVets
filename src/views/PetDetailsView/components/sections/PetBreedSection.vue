<script setup lang="ts">
import type {Breed} from "@/types/Pet";
import PetNameFormatted from "@/views/PetDetailsView/components/PetNameFormatted.vue";
import {computed} from "vue";

const props = defineProps<{
  name: string,
  breedsList: Breed[]   // already filtered to the pet's species
}>()

const breedId = defineModel('breedId', {
  type: String,
  required: true
})

const isValidSelectedBreed = computed<boolean>(() => {
  return props.breedsList.some((breed) => breed.id === breedId.value);
})

</script>

<template>
  <section>
    <div>
      <p>What breed is
        <PetNameFormatted :name="props.name"/>
        <span>?</span>
      </p>

      <div class="mt-4">

        <FloatLabel variant="on">
          <Select v-model="breedId" inputId="breed"
                  :options="props.breedsList"
                  optionLabel="name"
                  optionValue="id"
                  filter
                  filterPlaceholder="Search breeds"
                  resetFilterOnHide
                  class="myInput"
                  :invalid="!isValidSelectedBreed"/>
          <label for="breed">Breed</label>
        </FloatLabel>
        <Message
            v-if="!isValidSelectedBreed" severity="error" size="small" variant="simple"
            aria-describedby="breedError">
          Please choose a valid breed or select <span class="fw-bold">Cross / Mixed breed</span> if unknown
        </Message>
      </div>
    </div>
  </section>
</template>

<style scoped>

/* A non-editable Select has no input inside to give it a width. */
.myInput {
  width: 300px;
}

</style>