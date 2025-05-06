<script setup lang="ts">
import PetTypeIds from "@/components/PetDetails/PetTypeIds";

const props = defineProps<{
  pet: Pet
}>();

</script>

<template>
  <div>
    <Avatar style="background-color: #f0f0f0; color: #1a2551" shape="circle" size="xlarge">
      <font-awesome-icon v-if="props.pet.expand.species.pet_type === PetTypeIds.catId" :icon="['fal', 'cat']"/>
      <font-awesome-icon v-else-if="props.pet.expand.species.pet_type === PetTypeIds.dogId" :icon="['fal', 'dog']"/>
      <font-awesome-icon v-else-if="props.pet.expand.species.pet_type === PetTypeIds.smallAnimalId"
                         :icon="['fal', 'rabbit']"/>
      <font-awesome-icon v-else :icon="['fal', 'turtle']"/>
    </Avatar>
  </div>
</template>

<style scoped>

</style>