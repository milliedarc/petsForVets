<script setup lang="ts">
import {computed} from "vue";

const props = defineProps<{
  name: string
}>()

/** The pet's name, or "your pet" while the name field is still empty. */
const myPetName = computed<string>(() => {
  return props.name.trim() === '' ? 'your pet' : props.name.trim()
})

</script>

<template>
  <span class="pet-name">{{ myPetName }}</span>
</template>

<style scoped>
.pet-name {
  color: var(--p-primary-color);
  font-weight: 700;
}
</style>
