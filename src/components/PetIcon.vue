<script setup lang="ts">
import PetTypeIds from "@/types/PetTypeIds";
import {computed} from "vue";
import {pb} from "@/components/Pocketbase"

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