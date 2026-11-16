<script setup lang="ts">
import PetTypeIds from "@/types/PetTypeIds";
import {computed} from "vue";
import {pb} from "@/components/Pocketbase"
import {Cat, Dog, Rabbit, Turtle} from "@lucide/vue";

const props = defineProps<{
  pet: Pet
}>();

const avatarUrl = computed<string>(() => {
  return pb.files.getURL(props.pet, props.pet.avatar);
})

</script>

<template>
  <div>
    <Avatar v-if="props.pet.avatar" :image="avatarUrl" size="xlarge" shape="circle"/>
    <Avatar v-else style="background-color: #f0f0f0; color: #1a2551" shape="circle" size="xlarge">
      <Cat v-if="props.pet.expand.species.pet_type === PetTypeIds.catId"/>
      <Dog v-else-if="props.pet.expand.species.pet_type === PetTypeIds.dogId"/>
      <Rabbit v-else-if="props.pet.expand.species.pet_type === PetTypeIds.smallAnimalId"/>
      <Turtle v-else/>
    </Avatar>
  </div>
</template>

<style scoped>

</style>