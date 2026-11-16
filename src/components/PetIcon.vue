<script setup lang="ts">
import type {Pet} from "@/types/Pet";
import PetTypeIds from "@/types/PetTypeIds";
import {Cat, Dog, Rabbit, Turtle} from "@lucide/vue";

const props = defineProps<{
  pet: Pet
}>();


</script>

<template>
  <div>
    <Avatar v-if="props.pet.avatar_url" :image="props.pet.avatar_url" size="xlarge" shape="circle"/>
    <Avatar v-else style="background-color: #f0f0f0; color: #1a2551" shape="circle" size="xlarge">
      <Cat v-if="props.pet.species.pet_type_id === PetTypeIds.catId"/>
      <Dog v-else-if="props.pet.species.pet_type_id === PetTypeIds.dogId"/>
      <Rabbit v-else-if="props.pet.species.pet_type_id === PetTypeIds.smallAnimalId"/>
      <Turtle v-else/>
    </Avatar>
  </div>
</template>

<style scoped>

</style>