<script setup lang="ts">
import PetNameFormatted from "@/views/PetDetailsView/components/PetNameFormatted.vue";
import {computed} from "vue";

const props = defineProps<{
  name: string,
  petType: string,
  breedsList: Breed[]
}>()

const breedId = defineModel('breedId', {
  type: String,
  required: true
})

const isValidSelectedBreed = computed<boolean>(() => {
  for (const oneBreed of filteredBreedsList.value) {
    if (oneBreed.id === breedId.value) {
      return true
    }
  }
  return false
})

const filteredBreedsList = computed<[]>(() => {
      if (props.breedsList === undefined) {
        return []
      }
      return props.breedsList.filter((breed) => {
        if (breed.pet_type === props.petType) {
          return true;
        }
        return false;
      })
    }
)

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
                  :options="filteredBreedsList"
                  optionLabel="name"
                  optionValue="id"
                  editable
                  class="myInput"
                  :invalid="!isValidSelectedBreed"/>
          <label for="breed">Breed</label>
        </FloatLabel>
        <Message
            v-if="!isValidSelectedBreed" severity="error" size="small" variant="simple"
            aria-describedby="breedError">
          Please choose a valid breed or select <span class="fw-bold">Cross / Mixed Breed</span> if unknown
        </Message>
      </div>
    </div>
  </section>
</template>

<style scoped>

</style>