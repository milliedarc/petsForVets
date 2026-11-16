<script setup lang="ts">
import type {Pet} from "@/types/Pet";
import PetTypeIds from "@/types/PetTypeIds";
import {Cat, Dog, Rabbit, Turtle} from "@lucide/vue";

const props = defineProps<{
  pet: Pet
}>();

</script>

<template>
  <Avatar v-if="props.pet.avatar_url"
          :image="props.pet.avatar_url"
          :alt="`Photo of ${props.pet.name}`"
          size="xlarge" shape="circle"
          class="pet-icon"/>
  <Avatar v-else shape="circle" size="xlarge" class="pet-icon pet-icon--placeholder" aria-hidden="true">
    <Cat v-if="props.pet.species.pet_type_id === PetTypeIds.catId"/>
    <Dog v-else-if="props.pet.species.pet_type_id === PetTypeIds.dogId"/>
    <Rabbit v-else-if="props.pet.species.pet_type_id === PetTypeIds.smallAnimalId"/>
    <Turtle v-else/>
  </Avatar>
</template>

<style scoped>
.pet-icon {
  flex-shrink: 0;
}

.pet-icon--placeholder {
  background: var(--pfv-accent-soft);
  color: var(--p-primary-700);
}
</style>
