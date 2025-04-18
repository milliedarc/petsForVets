<script setup>
import {computed} from "vue";

const props = defineProps(['name'])

const myPetName = computed(() => {
  if (props.name === '') {
    return 'My Pet'
  }
  return props.name
})

</script>

<template>
  <span class="fw-bold">{{ myPetName }}</span>
</template>

<style scoped>

</style>