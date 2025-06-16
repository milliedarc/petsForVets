<script setup lang="ts">
import {computed} from "vue";

const props = defineProps(['name'])

const myPetName = computed<string>(() => {
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